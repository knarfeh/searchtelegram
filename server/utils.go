package main

import (
	"github.com/knarfeh/searchtelegram/server/domain"
)

// Must raises an error if it not nil
func Must(e error) {
	if e != nil {
		panic(e)
	}
}

// TgResource2Str ...
func TgResource2Str(tgResource domain.TgResource) (string, string) {
	description := tgResource.Desc
	if description == "" {
		description = "None"
	}
	tagsString := Tags2String(tgResource.Tags)
	// Copy from https://emojifinder.com/
	typeEmoji := ""
	switch tgResource.Type {
	case "bot":
		typeEmoji = "🤖"
	case "channel":
		typeEmoji = "🔊"
	case "group":
		typeEmoji = "👥"
	case "people":
		typeEmoji = "👤"
	}
	rawMessage := tgResource.TgID + "\n\nType: " + tgResource.Type + "\nDescription: " + description + "\nTags: " + tagsString
	sigStr := "\n\nBy searchtelegram \n@searchtelegramdotcombot   Robot, index of telegram \n@searchtelegramchannel         searchtelegram updates \n@searchtelegrampublic            Public group of searchtelegram"
	return rawMessage + sigStr, typeEmoji
}

// Tags2String ...
func Tags2String(tags []domain.Tag) string {
	result := ""
	for _, entry := range tags {
		result = result + "#" + entry.Name + " "
	}
	return result
}
