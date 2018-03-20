package main

import (
	"fmt"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

// TeleBot encapsulation telebot, redis(for search)
type TeleBot struct {
	tb          *tb.Bot
	redisClient *RedisClient
}

// CreateTeleBot create Telebot
func CreateTeleBot(tgBotToken, redisHost, redisPort string) (*TeleBot, error) {
	b, err := tb.NewBot(tb.Settings{
		Token:  tgBotToken,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		fmt.Println(err)
	}
	redisClient := NewRedisClient(redisHost, redisPort)

	telebot := &TeleBot{
		tb:          b,
		redisClient: redisClient,
	}

	b.Handle("/ping", telebot.pong)
	b.Handle("/start", telebot.start)
	b.Handle("/detail", telebot.detail)
	b.Handle("/submit", telebot.detail)
	b.Handle("/search", telebot.search)
	b.Handle("/search_group", telebot.searchGroup)
	b.Handle("/search_bot", telebot.searchBot)
	b.Handle("/search_channel", telebot.searchChannel)
	b.Handle("/search_people", telebot.searchChannel)
	// b.Handle("/top", telebot.pong)     // TODO

	return telebot, nil
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

func (telebot *TeleBot) start(m *tb.Message) {
	// get detail of an tg_ID
	fmt.Println(m.Sender)
	telebot.tb.Send(m.Sender, "start, TODO")
}

func (telebot *TeleBot) detail(m *tb.Message) {
	// get detail of an tg_ID
	fmt.Println(m.Sender)
	telebot.tb.Send(m.Sender, "detail, TODO")
}

func (telebot *TeleBot) submit(m *tb.Message) {
	// submit new group, channel, bot
	fmt.Println(m.Sender)
	telebot.tb.Send(m.Sender, "submit, TODO")
}

func (telebot *TeleBot) search(m *tb.Message) {
	// search group balabala
	fmt.Println(m.Sender)
	telebot.tb.Send(m.Sender, "search, TODO")
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
