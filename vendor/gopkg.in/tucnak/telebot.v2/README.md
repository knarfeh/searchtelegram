# Telebot
>Telebot is a Telegram bot framework in Go.

[![GoDoc](https://godoc.org/github.com/tucnak/telebot?status.svg)](https://godoc.org/github.com/tucnak/telebot)
[![Travis](https://travis-ci.org/tucnak/telebot.svg?branch=master)](https://travis-ci.org/tucnak/telebot)

**NOTE:** We are currently working on completing both v1 and v2 of our API.
You can track the progress [here](https://github.com/tucnak/telebot/milestone/1)
and [here](https://github.com/tucnak/telebot/issues/93).

Bots are special Telegram accounts designed to handle messages automatically.
Users can interact with bots by sending them command messages in private or
via group chats / channels. These accounts serve as an interface to your code.

Telebot offers a pretty convenient interface to Bots API and uses default HTTP
client. Ideally, you wouldn't need to worry about actual networking at all.

	go get gopkg.in/tucnak/telebot.v1

(after setting up your `GOPATH` properly).

We highly recommend you to keep your bot access token outside the code base,
preferably in an environmental variable:

	export BOT_TOKEN=<your token here>

Take a look at a minimal functional bot setup:
```go
package main

import (
	"log"
	"os"
	"time"

	tb "gopkg.in/tucnak/telebot.v1"
)

func main() {
	bot, err := tb.NewBot(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Fatalln(err)
	}

	messages := make(chan tb.Message, 100)
	bot.Listen(messages, 10 * time.Second)

	for message := range messages {
		if message.Text == "/hi" {
			bot.SendMessage(message.Chat,
				"Hello, "+message.Sender.FirstName+"!", nil)
		}
	}
}
```

## Inline mode
As of January 4, 2016, Telegram added inline mode support for bots. Here's
a nice way to handle both incoming messages and inline queries in the meantime:

```go
package main

import (
	"log"
	"time"
	"os"

	tb "gopkg.in/tucnak/telebot.v1"
)

func main() {
	bot, err := tb.NewBot(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Fatalln(err)
	}

	bot.Messages = make(chan tb.Message, 100)
	bot.Queries = make(chan tb.Query, 1000)

	go messages(bot)
	go queries(bot)

	bot.Start(10 * time.Second)
}

func messages(bot *tb.Bot) {
	for message := range bot.Messages {
		log.Printf("Received a message from %s with the text: %s\n",
			message.Sender.Username, message.Text)
	}
}

func queries(bot *tb.Bot) {
	for query := range bot.Queries {
		log.Println("--- new query ---")
		log.Println("from:", query.From.Username)
		log.Println("text:", query.Text)

		// Create an article (a link) object to show in results.
		article := &tb.InlineQueryResultArticle{
			Title: "Telebot",
			URL:   "https://github.com/tucnak/telebot",
			InputMessageContent: &tb.InputTextMessageContent{
				Text:		   "Telebot is a Telegram bot framework.",
				DisablePreview: false,
			},
		}

		// Build the list of results (make sure to pass pointers!).
		results := []tb.InlineQueryResult{article}

		// Build a response object to answer the query.
		response := tb.QueryResponse{
			Results:	results,
			IsPersonal: true,
		}

		// Send it.
		if err := bot.AnswerInlineQuery(&query, &response); err != nil {
			log.Println("Failed to respond to query:", err)
		}
	}
}
```

## Files
Telebot lets you upload files from the file system:

```go
boom, err := tb.NewFile("boom.ogg")
if err != nil {
	return err
}

audio := tb.Audio{File: boom}

// Next time you send &audio, telebot won't issue
// an upload, but would re-use existing file.
err = bot.SendAudio(recipient, &audio, nil)
```

## Reply markup
Sometimes you wanna send a little complicated messages with some optional parameters. The third argument of all `Send*` methods accepts `telebot.SendOptions`, capable of defining an advanced reply markup:

```go
// Send a selective force reply message.
bot.SendMessage(user, "pong", &tb.SendOptions{
		ReplyMarkup: tb.ReplyMarkup{
			ForceReply: true,
			Selective: true,
			CustomKeyboard: [][]string{
				[]string{"1", "2", "3"},
				[]string{"4", "5", "6"},
				[]string{"7", "8", "9"},
				[]string{"*", "0", "#"},
			},
		},
	},
)
```
