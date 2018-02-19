package main

import (
	"context"
	"encoding/json"
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
	tgInfo := hauler.getData(tgResource.TgID)
	fmt.Println("TODO, handle 404")
	fmt.Println("tgInfo????")
	fmt.Println(tgInfo)
	// TODO write to es
	tgResource.Title = tgInfo.Title
	tgResource.Type = "TODO, add tags!!!"
	tgResource.Imgsrc = tgInfo.ImgSrc
	if tgResource.Desc == "" {
		fmt.Println("tgResource.Desc is nil")
		tgResource.Desc = tgInfo.Description
	}
	_, err := hauler.esClient.Client.Index().OpType("create").Index("telegram").Type("resource").Id(tgResource.TgID).BodyJson(tgResource).Do(context.TODO())

	if err != nil {
		// Please make sure domain not exist
		e, _ := err.(*elastic.Error)
		if e.Status == 409 {
			errorItem := make(map[string]string)
			errorItem["code"] = "resource_already_exist"
			errorItem["message"] = e.Details.Reason
			errorItem["source"] = "10001"
			panic(e)
		}
		// Should not happen...
		panic(err)
	}
}

// getData
func (hauler *Hauler) getData(tgID string) *tgMeInfo {
	url := "https://t.me/" + tgID
	fmt.Printf("Getting data, url: %s", url)
	doc, err := goquery.NewDocument(url)

	title := strings.TrimSpace(doc.Find(".tgme_page_title").Text())
	description := strings.TrimSpace(doc.Find(".tgme_page_description").Text())
	imgSrc, _ := doc.Find(".tgme_page_photo_image").Attr("src")
	// TODO get type from button

	// don't worry about errors
	response, e := http.Get(imgSrc)
	if e != nil {
		fmt.Print(e)
	}

	defer response.Body.Close()

	// open a file for writing
	file, err := os.Create("/media/images/" + tgID + ".jpg")
	if err != nil {
		fmt.Print(err)
	}
	// Use io.Copy to just dump the response body to the file. This supports huge files
	_, err = io.Copy(file, response.Body)
	if err != nil {
		fmt.Print(err)
	}
	file.Close()
	fmt.Println("Success!")

	return &tgMeInfo{
		Title:       title,
		Description: description,
		ImgSrc:      "/media/images/" + tgID + ".jpg",
		// Type
	}
}
