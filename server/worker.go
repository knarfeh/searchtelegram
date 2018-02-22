package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"

	"github.com/knarfeh/searchtelegram/server/domain"

	elastic "gopkg.in/olivere/elastic.v5"
)

// Hauler ...
type Hauler struct {
	esClient    *ESClient
	redisClient *RedisClient
}

type tgMeInfo struct {
	Title       string
	Description string
	ImgSrc      string
	Type        string
}

// CreateConsumer create consumer ...
func CreateConsumer(esURL, redisHost, redisPort string) (*Hauler, error) {
	es, _ := NewESClient(esURL, "", "", 3)
	redisClient := NewRedisClient(redisHost, redisPort)
	return &Hauler{
		esClient:    es,
		redisClient: redisClient,
	}, nil
}

// Query2ES subscribe redis channel, get data from t.me, save it to es
// func
func (hauler *Hauler) Query2ES() {
	pubsub := hauler.redisClient.Client.Subscribe("searchtelegram")
	defer pubsub.Close()

	substr, err := pubsub.ReceiveTimeout(time.Second)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Subscribe string: %s", substr)

	ch := pubsub.Channel()
	for {
		select {
		case msg := <-ch:
			hauler.handleQuery(msg.Payload)
		}
	}
}

// handleQuery ...
func (hauler *Hauler) handleQuery(query string) {
	tgResource := domain.NewTgResource()
	if err := json.Unmarshal([]byte(query), &tgResource); err != nil {
		panic(err)
	}

	fmt.Println(tgResource)
	tgInfo, err := hauler.getData(tgResource.TgID)
	if err != nil {
		fmt.Println("Got error!!!")
		fmt.Println(err.Error())
		return
	}
	fmt.Println("TODO, handle 404")
	fmt.Println("tgInfo????")
	fmt.Println(tgInfo)

	tgResource.Title = tgInfo.Title
	tgResource.Type = tgInfo.Type // TODO: add tags
	tagItem := &domain.Tag{
		Count: 1,
		Name:  tgInfo.Type,
	}
	tgResource.Tags = append(tgResource.Tags, *tagItem)
	tgResource.Imgsrc = tgInfo.ImgSrc
	if tgResource.Desc == "" {
		fmt.Println("tgResource.Desc is nil, got from tdotme")
		tgResource.Desc = tgInfo.Description
	}
	_, err = hauler.esClient.Client.Index().OpType("create").Index("telegram").Type("resource").Id(tgResource.TgID).BodyJson(tgResource).Do(context.TODO())

	if err != nil {
		// Please make sure domain not exist
		e, _ := err.(*elastic.Error)
		if e.Status == 409 {
			errorItem := make(map[string]string)
			errorItem["code"] = "resource_already_exist"
			errorItem["message"] = e.Details.Reason
			errorItem["source"] = "10001"
			fmt.Println("Conflict, dont panic, error: ", errorItem)
		}
		// Should not happen...
		// panic(err)
	}
}

// getTgType ...
func (hauler *Hauler) getTgType(content string) string {
	if strings.Contains(content, "GROUP") {
		return "group"
	}
	if strings.Contains(content, "CHANNEL") {
		return "channel"
	}
	return "unknown"
}

// downloadPic ...
func (hauler *Hauler) downloadPic(imgSrc, tgID string) bool {
	response, e := http.Get(imgSrc)
	if e != nil {
		fmt.Print(e)
	}

	defer response.Body.Close()

	// open a file for writing
	file, err := os.Create("/tmp/images/" + tgID + ".jpg")
	if err != nil {
		fmt.Print(err)
		return false
	}
	// Use io.Copy to just dump the response body to the file. This supports huge files
	_, err = io.Copy(file, response.Body)
	if err != nil {
		fmt.Print(err)
		return false
	}
	file.Close()
	return true
}

// getData ...
func (hauler *Hauler) getData(tgID string) (*tgMeInfo, error) {
	url := "https://t.me/" + tgID
	fmt.Printf("Getting data, url: %s", url)
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return nil, errors.New("Got error from tdotme...")
	}

	title := strings.TrimSpace(doc.Find(".tgme_page_title").Text())
	description := strings.TrimSpace(doc.Find(".tgme_page_description").Text())
	imgSrc, _ := doc.Find(".tgme_page_photo_image").Attr("src")
	buttonContent := doc.Find(".tgme_action_button_new").Text()
	tgType := hauler.getTgType(strings.ToUpper(buttonContent))

	if tgType == "unknown" {
		if strings.HasSuffix(strings.ToUpper(tgID), "BOT") {
			tgType = "bot"
		} else {
			tgType = "people"
		}
	}

	imgPath := ""
	if imgSrc != "" {
		hauler.downloadPic(imgSrc, tgID)
		imgPath = "/images/" + tgID + ".jpg"
	} else {
		imgPath = "/images/telegram.jpg"
	}
	if title == "" && strings.HasPrefix(description, "If you have Telegram, you can contact") {
		return nil, errors.New("404")
	}

	return &tgMeInfo{
		Title:       title,
		Description: description,
		ImgSrc:      imgPath,
		Type:        tgType,
	}, nil
}
