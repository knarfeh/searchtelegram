package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/knarfeh/searchtelegram/server/domain"
	"github.com/olebedev/config"

	"github.com/RedisLabs/redisearch-go/redisearch"
	tb "github.com/tucnak/telebot"
	elastic "gopkg.in/olivere/elastic.v5"
)

// TeleBot encapsulation telebot, redis(for search)
type TeleBot struct {
	tb               *tb.Bot
	redisClient      *RedisClient
	redisearchClient *RedisearchClient
	esClient         *ESClient
}

// CreateTeleBot create Telebot
func CreateTeleBot(conf *config.Config) (*TeleBot, error) {
	ESHOSTPORT, _ := conf.String("ESHOSTPORT")
	REDISHOST, _ := conf.String("REDISHOST")
	REDISPORT, _ := conf.String("REDISPORT")
	TGBOTTOKEN, _ := conf.String("TGBOTTOKEN")
	es, _ := NewESClient(
		ElasticConfig{
			Endpoint:           ESHOSTPORT,
			Username:           "",
			Password:           "",
			Retries:            3,
			HealthCheckTimeout: 3 * time.Second,
		},
	)
	b, err := tb.NewBot(tb.Settings{
		Token:  TGBOTTOKEN,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		fmt.Println(err)
	}
	redisClient := NewRedisClient(REDISHOST, REDISPORT)
	redisearchClient := NewRedisearchClient(REDISHOST, REDISPORT)

	telebot := &TeleBot{
		esClient:         es,
		tb:               b,
		redisClient:      redisClient,
		redisearchClient: redisearchClient,
	}

	b.Handle("/start", telebot.start)
	b.Handle("/get", telebot.get)
	b.Handle("/submit", telebot.submit)
	b.Handle("/search", telebot.search)
	b.Handle("/s", telebot.search)
	b.Handle("/help", telebot.start)
	// b.Handle("/tips", telebot.tips)

	b.Handle("/search_group", telebot.searchGroup)
	b.Handle("/s_group", telebot.searchGroup)

	b.Handle("/search_bot", telebot.searchBot)
	b.Handle("/s_bot", telebot.searchBot)

	b.Handle("/search_channel", telebot.searchChannel)
	b.Handle("/s_channel", telebot.searchChannel)

	b.Handle("/search_people", telebot.searchChannel)
	b.Handle("/s_people", telebot.searchChannel)

	// b.Handle("/top", telebot.top)
	b.Handle("/delete", telebot.delete)
	b.Handle("/ping", telebot.pong)
	b.Handle("/stats", telebot.stats)

	return telebot, nil
}

// get detail of an tg_ID
func (telebot *TeleBot) start(m *tb.Message) {
	fmt.Printf("[start]sender: %s, user id: %d, payload: %s\n", m.Sender.Username, m.Sender.ID, m.Payload)
	telebot.tb.Send(m.Sender, StartInfo())
	telebot.redisClient.Client.SAdd("status:unique-user", m.Sender.Username)
}

// get detail of an tg_ID
func (telebot *TeleBot) get(m *tb.Message) {
	fmt.Printf("[detail]sender: %s, user id: %d, payload: %s\n", m.Sender.Username, m.Sender.ID, m.Payload)

	tgIDExist, _ := telebot.redisClient.Client.Get("tgid:" + m.Payload).Result()
	if tgIDExist == "" {
		telebot.tb.Send(m.Sender, "Ops, this id does not exist, perhaps you could submit with /submit "+m.Payload)
		return
	}

	resourceResult, _ := telebot.esClient.Client.Get().Index("telegram").Type("resource").Id(m.Payload).Do(context.TODO())
	tgResource := domain.NewTgResource()
	json.Unmarshal(*resourceResult.Source, tgResource)

	message, emoji := TgResource2Str(*tgResource)
	channelMessage := "\n" + emoji + "\n \n @" + message
	telebot.tb.Send(m.Sender, channelMessage)
	telebot.redisClient.Client.SAdd("status:get-unique-user", m.Sender.Username)
}

// submit new group, channel, bot
func (telebot *TeleBot) submit(m *tb.Message) {
	fmt.Printf("[submit]sender: %s, user id: %d, payload: %s\n", m.Sender.Username, m.Sender.ID, m.Payload)
	if m.Payload == "" {
		telebot.tb.Send(m.Sender, "Please input telegram ID, like: /submit telegram")
		return
	}
	tgIDExist, _ := telebot.redisClient.Client.Get("tgid:" + m.Payload).Result()
	if tgIDExist != "" {
		telebot.tb.Send(m.Sender, "Ha, this id already exist, you could get detailed information with /get "+m.Payload)
		return
	}
	tgResource := domain.NewTgResource()
	tgResource.TgID = m.Payload
	tgResouceString, _ := json.Marshal(tgResource)
	fmt.Printf("Telegram, %s submit resource: %s\n", m.Sender.Username, tgResouceString)
	err := telebot.redisClient.Client.Publish("st_submit", string(1))
	telebot.redisClient.Client.LPush("st_submit_list", string(tgResouceString))
	if err != nil {
		panic(err)
	}

	telebot.tb.Send(m.Sender, "👏Successfully submitted. If everything goes well, you will be able to search for it after a while.")
	telebot.redisClient.Client.SAdd("status:submit-unique-user", m.Sender.Username)
}

func (telebot *TeleBot) search(m *tb.Message) {
	fmt.Printf("[search]username: %s, payload: %s\n", m.Sender.Username, m.Payload)
	if m.Payload == "" {
		telebot.tb.Send(m.Sender, "Please input search string, like: /search telegram")
		return
	}

	simpleQuery, boolQuery, queryString, tagsSlice := BuildESQuery(m.Payload)
	val := telebot.redisClient.Client.SIsMember("redisearch:cached-search-string", queryString).Val()
	pipe := telebot.redisClient.Client.Pipeline()
	pipe.SAdd("status:search-unique-user", m.Sender.Username)
	pipe.Publish("st_search", string(queryString))
	if val {
		rediQueryStr := queryString
		if len(tagsSlice) != 0 {
			rediQueryStr = queryString + fmt.Sprintf(" @tags:{%s}", strings.Join(tagsSlice, "|"))
		}
		fmt.Printf("Go queryString %s in redis, rediQueryStr: %s\n", queryString, rediQueryStr)
		q := redisearch.NewQuery(rediQueryStr)
		docs, total, _ := telebot.redisearchClient.Client.Search(q)
		result := Redisearch2Str(docs, total)
		telebot.tb.Send(m.Sender, result)

		if _, err := pipe.Exec(); err != nil {
			fmt.Println(err)
		}
		return
	}

	fmt.Printf("No cache in redis, search in elasticsearch, search str: %s\n", m.Payload)
	search := telebot.esClient.Client.Search().Index("telegram").Type("resource").Size(20).PostFilter(boolQuery)
	searchResult, err := search.Query(simpleQuery).Do(context.TODO())
	if err != nil {
		panic(err)
	}

	result := Hits2Str(*searchResult.Hits)
	telebot.tb.Send(m.Sender, result)

	if _, err := pipe.Exec(); err != nil {
		fmt.Println(err)
	}
}

func (telebot *TeleBot) tags(m *tb.Message) {
	// search group balabala
	fmt.Println(m.Sender)
	telebot.tb.Send(m.Sender, "tags, TODO\n")
}

func (telebot *TeleBot) help(m *tb.Message) {
	// search group balabala
	fmt.Println(m.Sender)
	telebot.tb.Send(m.Sender, "help, TODO\n")
}

func (telebot *TeleBot) searchGroup(m *tb.Message) {
	// search balabala
	fmt.Println(m.Sender)
	telebot.tb.Send(m.Sender, "search group, TODO\n")
}

func (telebot *TeleBot) searchBot(m *tb.Message) {
	fmt.Println(m.Sender)
	telebot.tb.Send(m.Sender, "search bot, TODO\n")
}

func (telebot *TeleBot) searchChannel(m *tb.Message) {
	fmt.Println(m.Sender)
	telebot.tb.Send(m.Sender, "search channel, TODO\n")
}

// Private. For test purpose
func (telebot *TeleBot) pong(m *tb.Message) {
	if m.Sender.Username != "knarfeh" {
		return
	}

	fmt.Println(m.Sender)
	telebot.tb.Send(m.Sender, "pong "+m.Payload)
	telebot.redisClient.Client.SAdd("status:ping-unique-user", m.Sender.Username)
}

// Private. Gathering server info
func (telebot *TeleBot) stats(m *tb.Message) {
	if m.Sender.Username != "knarfeh" {
		return
	}
	result := telebot.serverStats()
	telebot.tb.Send(m.Sender, result)
	telebot.redisClient.Client.SAdd("status:status-unique-user", m.Sender.Username)
}

// serverStatus ...
func (telebot *TeleBot) serverStats() string {
	pipe := telebot.redisClient.Client.Pipeline()
	uniqueUserPipe := pipe.SCard("status:unique-user")
	searchUniqueUserPipe := pipe.SCard("status:search-unique-user")
	getUniqueUserPipe := pipe.SCard("status:get-unique-user")
	submitUniqueUserPipe := pipe.SCard("status:submit-unique-user")
	pingUniqueUserPipe := pipe.SCard("status:ping-unique-user")
	statusUniqueUserPipe := pipe.SCard("status:status-unique-user")
	cachedStringsPipe := pipe.SMembers("redisearch:cached-search-string")
	if _, err := pipe.Exec(); err != nil {
		fmt.Println(err)
	}

	uniqueUserStr := fmt.Sprintf("Unique user: %d\n", uniqueUserPipe.Val())
	searchUniqueUserStr := fmt.Sprintf("Unique user who input /search: %d\n", searchUniqueUserPipe.Val())
	getUniqueUserStr := fmt.Sprintf("Unique user who input /get: %d\n", getUniqueUserPipe.Val())
	submitUniqueUserStr := fmt.Sprintf("Unique user who input /submit: %d\n", submitUniqueUserPipe.Val())
	pingUniqueUserStr := fmt.Sprintf("Unique user who input /ping: %d\n", pingUniqueUserPipe.Val())
	statusUniqueUserStr := fmt.Sprintf("Unique user who input /status: %d\n", statusUniqueUserPipe.Val())
	cachedSearchStr := fmt.Sprintf("Cached search string:\n %s\n", strings.Join(cachedStringsPipe.Val(), ", "))

	// TODO: leaderboard

	docCount, _ := telebot.esClient.Client.Count("telegram").Do(context.TODO())
	esDocCountStr := fmt.Sprintf("ES Document count: %d\n", docCount)

	tagCountAgg := elastic.NewCardinalityAggregation().Field("tags.name.keyword")
	aggBuilder := telebot.esClient.Client.Search().Index("telegram").Type("resource").Query(elastic.NewMatchAllQuery()).Size(0)
	aggBuilder = aggBuilder.Aggregation("tagsCardinality", tagCountAgg).Size(0)
	searchResult, _ := aggBuilder.Do(context.TODO())
	tagCountResult, _ := searchResult.Aggregations.Cardinality("tagsCardinality")
	tagsCountStr := fmt.Sprintf("Tags count: %v\n", *tagCountResult.Value)

	return uniqueUserStr + searchUniqueUserStr + getUniqueUserStr + submitUniqueUserStr + pingUniqueUserStr + statusUniqueUserStr + cachedSearchStr + "\n" + esDocCountStr + tagsCountStr
}

// Private. Delete an item
func (telebot *TeleBot) delete(m *tb.Message) {
	fmt.Println(m.Sender)
	telebot.esClient.Client.Delete().Index("telegram").Type("resource").Id(m.Payload).Do(context.TODO())
	telebot.redisClient.Client.Expire("tgid:"+m.Payload, 1)
	telebot.redisearchClient.Client.Drop()
	telebot.redisClient.Client.Expire("redisearch:cached-search-string", 1)
	result := fmt.Sprintf("%s had been deleted", m.Payload)
	telebot.tb.Send(m.Sender, result)
}
