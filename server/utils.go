package main

import (
	"encoding/json"
	"fmt"
	"github.com/knarfeh/searchtelegram/server/domain"
	"strings"

	"github.com/RedisLabs/redisearch-go/redisearch"
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
	fmt.Printf("Got tagstring: %s\n", tagstring)
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

// Hits2Str ...
func Hits2Str(hits elastic.SearchHits) string {
	result := "🎉🎉🎉 " + fmt.Sprintf("%d", hits.TotalHits) + " results\n\n"
	if hits.TotalHits == 1 {
		result = "🎉🎉🎉 " + fmt.Sprintf("%d", hits.TotalHits) + " result\n\n"
	}
	if hits.TotalHits == 0 {
		return "😱Sorry, but we don't find any result"
	}
	hitStr := ""
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

// Redisearch2Str ...
func Redisearch2Str(docs []redisearch.Document, total int) string {
	result := "🎉🎉🎉 " + fmt.Sprintf("%d", total) + " results\n\n"
	if total == 1 {
		result = "🎉🎉🎉 " + fmt.Sprintf("%d", total) + " result\n\n"
	}
	if total == 0 {
		return "😱Sorry, but we don't find any result"
	}
	for _, doc := range docs {
		tgType := fmt.Sprintf("%s", doc.Properties["type"])
		tgID := fmt.Sprintf("%s", doc.Properties["tgid"])
		tgDesc := fmt.Sprintf("%s", doc.Properties["desc"])
		tgTags := fmt.Sprintf("%s", doc.Properties["tags"])

		if tgDesc == "" {
			tgDesc = "None"
		}
		hitStr := emojiWithType(tgType) + "  @" + tgID + "\nDescription: " + strings.TrimSpace(tgDesc) + "\nTags: " + tgTags + "\n\n"
		result = result + hitStr
	}
	return result + sigStr()
}

// StartInfo ...
func StartInfo() string {
	result := `
🇬🇧
I will help you search telegram group, channel, bot, people. You can also submit new item, get details with telegram ID

🇨🇳
我可以帮助您搜索电报群组，频道，机器人，用户。您也可以提交新的电报 ID，根据 ID 获取详细信息。

/start Get help information

/search [searchstring] [tagstring] Search group, channel, bot, people
  e.g. /search telegram #group#people#tag3

/get [telegramID] Get details with telegram ID
  e.g. /get searchtelegramdotcombot

/submit [telegramID] Submit new item
  e.g. /submit searchtelegramchannel

/s_channel [channelID] Search channel
  e.g. /s_channel telegram

/s_group [groupID] Search group
  e.g. /s_group python

/s_bot [channelID] Search bot
  e.g. /s_bot picture

Our website: https://searchtelegram.com
`
	return result
}

// BuildESQuery get payload string, return boolquery, simpleQuery, queryString
func BuildESQuery(payload string) (*elastic.SimpleQueryStringQuery, *elastic.BoolQuery, string, []string) {
	splitPayload := strings.SplitN(payload, "#", 2)
	queryString := splitPayload[0]
	tagstring := ""
	if len(splitPayload) == 2 {
		tagstring = "#" + splitPayload[1]
	}
	boolQuery := elastic.NewBoolQuery()
	tagSlice := String2TagSlice(tagstring)
	for _, item := range tagSlice {
		if item == " " || item == "" {
			continue
		}
		boolQuery = boolQuery.Should(elastic.NewTermQuery("tags.name.keyword", item))
	}
	simpleQuery := elastic.NewSimpleQueryStringQuery(queryString)

	return simpleQuery, boolQuery, queryString, tagSlice
}
