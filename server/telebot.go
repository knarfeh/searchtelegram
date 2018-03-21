package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/knarfeh/searchtelegram/server/domain"
	"github.com/olebedev/config"

	tb "github.com/tucnak/telebot"
	elastic "gopkg.in/olivere/elastic.v5"
)

// TeleBot encapsulation telebot, redis(for search)
type TeleBot struct {
	tb          *tb.Bot
	redisClient *RedisClient
	esClient    *ESClient
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

	telebot := &TeleBot{
		esClient:    es,
		tb:          b,
		redisClient: redisClient,
	}

	b.Handle("/ping", telebot.pong)
	b.Handle("/start", telebot.start)
	b.Handle("/get", telebot.get)
	b.Handle("/submit", telebot.submit)
	b.Handle("/search", telebot.search)
	b.Handle("/help", telebot.help)
	b.Handle("/search_group", telebot.searchGroup)
	b.Handle("/search_bot", telebot.searchBot)
	b.Handle("/search_channel", telebot.searchChannel)
	b.Handle("/search_people", telebot.searchChannel)
	// b.Handle("/top", telebot.pong)     // TODO

	return telebot, nil
}

func (telebot *TeleBot) start(m *tb.Message) {
	// get detail of an tg_ID
	fmt.Println(m.Sender)
	telebot.tb.Send(m.Sender, "start, TODO")
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
	err := telebot.redisClient.Client.Publish("searchtelegram", string(tgResouceString)).Err()
	if err != nil {
		panic(err)
	}

	telebot.tb.Send(m.Sender, "Successfully submitted. If everything goes well, you will be able to search for it after a while.")
}

func (telebot *TeleBot) search(m *tb.Message) {
	if m.Payload == "" {
		telebot.tb.Send(m.Sender, "Please input search string, like: /search telegram")
		return
	}
	splitPayload := strings.SplitN(m.Payload, " ", 2)
	queryString := splitPayload[0]
	var tagstring string
	if len(splitPayload) == 2 {
		tagstring = splitPayload[1]
	}
	if !strings.Contains(tagstring, "#") {
		fmt.Printf("No # in tagstring, sender: %s", m.Sender.Username)
		tagstring = ""
	}
	boolQuery := elastic.NewBoolQuery()
	for _, item := range String2TagSlice(tagstring) {
		if item == " " || item == "" {
			continue
		}
		boolQuery = boolQuery.Should(elastic.NewTermQuery("tags.name.keyword", item))
	}
	simpleQuery := elastic.NewSimpleQueryStringQuery(queryString)
	// fmt.Println(simpleQuery.Source()) // TODO: add tests
	// fmt.Println(boolQuery.Source())
	search := telebot.esClient.Client.Search().Index("telegram").Type("resource").Size(20).PostFilter(boolQuery)
	searchResult, err := search.Query(simpleQuery).Do(context.TODO())
	if err != nil {
		panic(err)
	}

	result := Hits2Str(*searchResult.Hits)
	telebot.tb.Send(m.Sender, result)
}

func (telebot *TeleBot) help(m *tb.Message) {
	// search group balabala
	fmt.Println(m.Sender)
	telebot.tb.Send(m.Sender, "help, TODO")
}

func (telebot *TeleBot) searchGroup(m *tb.Message) {
	// search balabala
	fmt.Println(m.Sender)
	telebot.tb.Send(m.Sender, "search group, TODO")
}

func (telebot *TeleBot) searchBot(m *tb.Message) {
	fmt.Println(m.Sender)
	telebot.tb.Send(m.Sender, "search bot, TODO")
}

func (telebot *TeleBot) searchChannel(m *tb.Message) {
	fmt.Println(m.Sender)
	telebot.tb.Send(m.Sender, "search channel, TODO")
}

// Private. For test purpose
func (telebot *TeleBot) pong(m *tb.Message) {
	if m.Sender.Username != "knarfeh" {
		return
	}

	stChannel := &tb.Chat{
		Type:     tb.ChatChannel,
		Username: "knarfehDebug",
	}
	telebot.tb.Send(stChannel, "pong "+m.Payload)

	fmt.Println(m.Sender)
	telebot.tb.Send(m.Sender, "pong "+m.Payload)
}

// Utils

// func (telebot *TeleBot) getTgResourceFromString(s string) {
// fmt.Println("Dont need it now")
// }
