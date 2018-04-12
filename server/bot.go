package main

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/RedisLabs/redisearch-go/redisearch"
	"github.com/knarfeh/searchtelegram/server/domain"
	tb "github.com/tucnak/telebot"
	elastic "gopkg.in/olivere/elastic.v5"
	"gopkg.in/telegram-bot-api.v4"
)

// Bot ...
type Bot struct {
	Token    string
	Tgbot    *tgbotapi.BotAPI
	app      *App
	handlers map[string]interface{}
	result   string
}

// NewBot ...
func NewBot(token string) (*Bot, error) {
	tgbot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		fmt.Println(err)
	}

	bot := &Bot{
		Token:    token,
		handlers: make(map[string]interface{}),
		Tgbot:    tgbot,
	}

	bot.Handle("/start", bot.start)
	bot.Handle("/get", bot.get)
	bot.Handle("/submit", bot.submit)
	bot.Handle("/search", bot.search)
	bot.Handle("/s", bot.search)
	// bot.Handle("/tags", bot.tags)
	bot.Handle("/help", bot.help)
	// b.Handle("/tips", telebot.tips)

	bot.Handle("/search_group", bot.searchGroup)
	bot.Handle("/s_group", bot.searchGroup)

	bot.Handle("/search_bot", bot.searchBot)
	bot.Handle("/s_bot", bot.searchBot)

	bot.Handle("/search_channel", bot.searchChannel)
	bot.Handle("/s_channel", bot.searchChannel)

	bot.Handle("/search_people", bot.searchChannel)
	bot.Handle("/s_people", bot.searchChannel)

	bot.Handle("/delete", bot.delete)
	bot.Handle("/ping", bot.pong)
	bot.Handle("/stats", bot.stats)
	// bot.Handle(("/update", bot.update)
	// bot.Handle("/echo", bot.echo)
	return bot, nil
}

func (b *Bot) handle(end string, m *tb.Message) bool {
	handler, ok := b.handlers[end]
	if !ok {
		return false
	}

	if handler, ok := handler.(func(*tb.Message)); ok {
		go func(b *Bot, handler func(*tb.Message), m *tb.Message) {
			handler(m)
		}(b, handler, m)
		return true
	}
	return false
}

// Handle ...
func (b *Bot) Handle(endpoint interface{}, handler interface{}) {
	switch end := endpoint.(type) {
	case string:
		b.handlers[end] = handler
	default:
		panic("telebot: unsupported endpoint")
	}
}

var (
	cmdRx = regexp.MustCompile(`^(\/\w+)(@(\w+))?(\s|$)(.+)?`)
)

// incommingUpdate ...
func (b *Bot) incommingUpdate(upd *tb.Update, app *App) string {
	messageString, _ := json.Marshal(upd)
	fmt.Printf("messageString: %s", messageString)

	b.app = app
	if upd.Message != nil {
		m := upd.Message

		// Commands
		if m.Text != "" {
			// Filtering malicious messages
			if m.Text[0] == '\a' {
				return ""
			}
			// Command found, handle and return
			match := cmdRx.FindAllStringSubmatch(m.Text, -1)
			if match != nil {
				// Syntax: "</command>@<bot> <payload"
				command, _ := match[0][1], match[0][3]
				m.Payload = match[0][5]

				// if botName != "" && !strings.EqualFold(b.Me, t string)

				if b.handle(command, m) {
					return b.result
				}
			}
		}
	}

	return ""
}

// Get start info
func (b *Bot) start(m *tb.Message) {
	fmt.Printf("[start]sender: %s, user id: %d, payload: %s\n", m.Sender.Username, m.Sender.ID, m.Payload)
	result := StartInfo()
	b.handleResult(m.Chat.ID, result)
	b.app.RedisClient.Client.SAdd("status:unique-user", m.Sender.Username)
}

// get detail of an tg_ID
func (b *Bot) get(m *tb.Message) {
	fmt.Printf("[detail]sender: %s, user id: %d, payload: %s\n", m.Sender.Username, m.Sender.ID, m.Payload)

	tgIDExist, _ := b.app.RedisClient.Client.Get("tgid:" + m.Payload).Result()
	if tgIDExist == "" {
		result := "Ops, this id does not exist, perhaps you could submit with /submit " + m.Payload
		b.handleResult(m.Chat.ID, result)
		return
	}

	resourceResult, _ := b.app.ESClient.Client.Get().Index("telegram").Type("resource").Id(m.Payload).Do(context.TODO())
	tgResource := domain.NewTgResource()
	json.Unmarshal(*resourceResult.Source, tgResource)

	message, emoji := TgResource2Str(*tgResource)
	result := "\n" + emoji + "\n \n @" + message
	b.handleResult(m.Chat.ID, result)
	b.app.RedisClient.Client.SAdd("status:get-unique-user", m.Sender.Username)
}

// submit new group, channel, bot
func (b *Bot) submit(m *tb.Message) {
	fmt.Printf("[submit]sender: %s, user id: %d, payload: %s\n", m.Sender.Username, m.Sender.ID, m.Payload)
	if m.Payload == "" {
		result := "Please input telegram ID, like: /submit telegram"
		b.handleResult(m.Chat.ID, result)
		b.app.RedisClient.Client.SAdd("status:submit-unique-user", m.Sender.Username)
		return
	}
	tgIDExist, _ := b.app.RedisClient.Client.Get("tgid:" + m.Payload).Result()
	if tgIDExist != "" {
		result := "Ha, this id already exist, you could get detailed information with /get " + m.Payload
		b.handleResult(m.Chat.ID, result)
		b.app.RedisClient.Client.SAdd("status:submit-unique-user", m.Sender.Username)
		return
	}
	tgResource := domain.NewTgResource()
	tgResource.TgID = m.Payload
	tgResouceString, _ := json.Marshal(tgResource)
	fmt.Printf("Telegram, %s submit resource: %s\n", m.Sender.Username, tgResouceString)
	err := b.app.RedisClient.Client.Publish("st_submit", string(1)).Err()
	b.app.RedisClient.Client.LPush("st_submit_list", string(tgResouceString))
	if err != nil {
		panic(err)
	}
	result := "👏Successfully submitted. If everything goes well, you will be able to search for it after a while."
	b.handleResult(m.Chat.ID, result)
	b.app.RedisClient.Client.SAdd("status:submit-unique-user", m.Sender.Username)
}

func (b *Bot) search(m *tb.Message) {
	fmt.Printf("[search]username: %s, payload: %s\n", m.Sender.Username, m.Payload)
	if m.Payload == "" {
		result := "Please input search string, like: /search telegram"
		b.handleResult(m.Chat.ID, result)
		return
	}
	simpleQuery, boolQuery, queryString, tagsSlice := BuildESQuery(m.Payload)
	val := b.app.RedisClient.Client.SIsMember("redisearch:cached-search-string", queryString).Val()
	pipe := b.app.RedisClient.Client.Pipeline()
	pipe.SAdd("status:search-unique-user", m.Sender.Username)
	pipe.Publish("st_search", string(queryString))
	if val {
		rediQueryStr := queryString
		tagsStr := strings.Join(tagsSlice, "|")
		if tagsStr != "" {
			rediQueryStr = queryString + fmt.Sprintf(" @tags:{%s}", tagsStr)
		}
		fmt.Printf("Go queryString %s in redis, rediQueryStr: %s\n", queryString, rediQueryStr)
		q := redisearch.NewQuery(rediQueryStr)
		docs, total, _ := b.app.RedisearchClient.Client.Search(q)
		result := Redisearch2Str(docs, total)
		b.handleResult(m.Chat.ID, result)

		if _, err := pipe.Exec(); err != nil {
			fmt.Println(err)
		}
		return
	}

	fmt.Printf("No cache in redis, search in elasticsearch, search str: %s\n", m.Payload)
	search := b.app.ESClient.Client.Search().Index("telegram").Type("resource").Size(20).PostFilter(boolQuery)
	searchResult, err := search.Query(simpleQuery).Do(context.TODO())
	if err != nil {
		panic(err)
	}

	result := Hits2Str(*searchResult.Hits)
	b.handleResult(m.Chat.ID, result)
	if _, err := pipe.Exec(); err != nil {
		fmt.Println(err)
	}
}

func (b *Bot) help(m *tb.Message) {
	fmt.Println(m.Sender)
	result := "help, TODO"
	b.handleResult(m.Chat.ID, result)
}

func (b *Bot) searchGroup(m *tb.Message) {
	fmt.Println(m.Sender)
	result := "search group, TODO"
	b.handleResult(m.Chat.ID, result)
}

func (b *Bot) searchBot(m *tb.Message) {
	fmt.Println(m.Sender)
	result := "search bot, TODO"
	b.handleResult(m.Chat.ID, result)
}

func (b *Bot) searchChannel(m *tb.Message) {
	fmt.Println(m.Sender)
	result := "search channel, TODO"
	b.handleResult(m.Chat.ID, result)
}

// ------------------------------- Private ------------------------------------
func (b *Bot) delete(m *tb.Message) {
	if m.Sender.Username != "knarfeh" {
		return
	}
	fmt.Printf("delete %s by %s", m.Payload, m.Sender.Username)
	b.app.ESClient.Client.Delete().Index("telegram").Type("resource").Id(m.Payload).Do(context.TODO())
	pipe := b.app.RedisClient.Client.Pipeline()
	pipe.Expire("tgid:"+m.Payload, 1*time.Second)
	pipe.Expire("redisearch:cached-search-string", 1*time.Second)
	if _, err := pipe.Exec(); err != nil {
		fmt.Println(err)
	}
	b.app.RedisearchClient.Client.Drop()
	result := fmt.Sprintf("%s had been deleted", m.Payload)
	b.handleResult(m.Chat.ID, result)
}

// For test purpose
func (b *Bot) pong(m *tb.Message) {
	if m.Sender.Username != "knarfeh" {
		return
	}

	fmt.Printf("ping! %s", m.Sender.Username)
	result := "pong " + m.Payload
	b.handleResult(m.Chat.ID, result)
	b.app.RedisClient.Client.SAdd("status:ping-unique-user", m.Sender.Username)
}

// Private. Gathering server info
func (b *Bot) stats(m *tb.Message) {
	if m.Sender.Username != "knarfeh" {
		return
	}
	result := b.serverStats()
	b.handleResult(m.Chat.ID, result)
	b.app.RedisClient.Client.SAdd("status:status-unique-user", m.Sender.Username)
}

// serverStatus ...
func (b *Bot) serverStats() string {
	pipe := b.app.RedisClient.Client.Pipeline()
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

	docCount, _ := b.app.ESClient.Client.Count("telegram").Do(context.TODO())
	esDocCountStr := fmt.Sprintf("ES Document count: %d\n", docCount)

	tagCountAgg := elastic.NewCardinalityAggregation().Field("tags.name.keyword")
	aggBuilder := b.app.ESClient.Client.Search().Index("telegram").Type("resource").Query(elastic.NewMatchAllQuery()).Size(0)
	aggBuilder = aggBuilder.Aggregation("tagsCardinality", tagCountAgg).Size(0)
	searchResult, _ := aggBuilder.Do(context.TODO())
	tagCountResult, _ := searchResult.Aggregations.Cardinality("tagsCardinality")
	tagsCountStr := fmt.Sprintf("Tags count: %v\n", *tagCountResult.Value)

	return uniqueUserStr + searchUniqueUserStr + getUniqueUserStr + submitUniqueUserStr + pingUniqueUserStr + statusUniqueUserStr + cachedSearchStr + "\n" + esDocCountStr + tagsCountStr
}

// handleResult ...
func (b *Bot) handleResult(chatID int64, result string) {
	msg := tgbotapi.NewMessage(chatID, result)
	go b.Tgbot.Send(msg)
	b.result = result
}
