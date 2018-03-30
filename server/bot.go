package main

import (
	"encoding/json"
	"fmt"

	"gopkg.in/telegram-bot-api.v4"
)

// IncommingUpdate ...
func IncommingUpdate(upd *tgbotapi.Update, app *App) string {
	messageString, _ := json.Marshal(upd)
	fmt.Printf("messageString: %s", messageString)
	err := app.RedisClient.Client.Publish("searchtelegram", "knarfeh").Err()
	if err != nil {
		panic(err)
	}
	return "lalal"
}
