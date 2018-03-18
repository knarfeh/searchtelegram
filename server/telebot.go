package main

import (
	"log"
	"time"

	tb "github.com/tucnak/telebot"
)

func testss() {
	b, err := tb.NewBot(tb.Settings{
		Token:  "TODO",
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	b.Handle("/hello", func(m *tb.Message) {
		b.Send(m.Sender, "hello world")
	})

	b.Start()
}
