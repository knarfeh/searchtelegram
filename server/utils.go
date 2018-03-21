package main

import (
	"encoding/json"
	"fmt"
	"github.com/knarfeh/searchtelegram/server/domain"
	"strings"

	elastic "gopkg.in/olivere/elastic.v5"
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
	emoji := emojiWithType(tgResource.Type)
	rawMessage := tgResource.TgID + "\nDescription: " + description + "\nTags: " + tagsString + "\n\n"
	return rawMessage + sigStr(), emoji
}

func sigStr() string {
	return "\n\nBy searchtelegram \n@searchtelegramdotcombot   Robot, index of telegram \n@searchtelegramchannel         searchtelegram updates \n@searchtelegrampublic            Public group of searchtelegram"
}

// emojiWithType ...
func emojiWithType(typeEmoji string) string {
	result := ""
	switch typeEmoji {
	case "bot":
		result = "🤖"
	case "channel":
		result = "🔊"
	case "group":
		result = "👥"
	case "people":
		result = "👤"
	}
	return result
}

// Tags2String ...
func Tags2String(tags []domain.Tag) string {
	result := ""
	for _, entry := range tags {
		result = result + "#" + entry.Name + " "
	}
	return result
}

// String2TagSlice ...
func String2TagSlice(tagstring string) []string {
	if !strings.Contains(tagstring, "#") {
		tagstring = ""
	} else {
		i := strings.Index(tagstring, "#")
		tagstring = tagstring[i:]
	}
	noSpaceString := strings.Replace(tagstring, " ", "", -1)
	notags := strings.Replace(noSpaceString, "#", " ", -1)
	justSpaceString := strings.TrimSpace(notags)
	result := strings.Split(justSpaceString, " ")
	return result
}

// TgResources2Str ...
func Hits2Str(hits elastic.SearchHits) string {
	result := "🎉🎉🎉 " + fmt.Sprintf("%d", hits.TotalHits) + " results\n\n"
	if hits.TotalHits == 1 {
		result = "🎉🎉🎉 " + fmt.Sprintf("%d", hits.TotalHits) + " result\n\n"
	}
	hitStr := ""
	if hits.TotalHits == 0 {
		return "😱Sorry, but we don't find any result"
	}
	for _, hit := range hits.Hits {
		hitStr = ""
		instance := domain.NewTgResource()
		json.Unmarshal(*hit.Source, instance)

		description := instance.Desc
		if description == "" {
			description = "None"
		}
		hitStr = emojiWithType(instance.Type) + "  @" + instance.TgID + "\nDescription: " + strings.TrimSpace(description) + "\nTags: " + Tags2String(instance.Tags) + "\n\n"
		result = result + hitStr
	}
	return result + sigStr()
}

// StartInfo ...
func StartInfo() string {
	result := `
/start Get help information.

/search [telegramID] [tags] Search group, channel, bot, people. i.e /search telegram #group#people#tag3

/get [telegramID] Get details by telegram ID. i.e /get searchtelegramdotcombot

/submit [telegramID] Submit an item, i.e /submit searchtelegramchannel
`
	return result
}
