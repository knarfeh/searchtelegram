package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/knarfeh/searchtelegram/server/domain"

	redis "github.com/go-redis/redis"
	// elastic "gopkg.in/olivere/elastic.v5"
)

type Hauler struct {
	// esClient    *elastic.Client
	redisClient *redis.Client
}

// CreateConsumer create consumer ...
func CreateConsumer(esURL, redisHost string) (*Hauler, error) {
	// esClient, err := elastic.NewClient(elastic.SetURL(esURL))
	// if err != nil {
	// return nil, err
	// }
	fmt.Printf("WTF is host??? %s", redisHost)
	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: "",
		DB:       0,
	})
	return &Hauler{
		// esClient:    esClient,
		redisClient: redisClient,
	}, nil
}

// Query2ES subscribe redis channel, get data from t.me, save it to es
// func
func (hauler *Hauler) Query2ES() {
	pubsub := hauler.redisClient.Subscribe("searchtelegram")
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

}

// getData
func (hauler *Hauler) getData(url string) string {
	fmt.Printf("Getting data, url: %s", url)
	doc, err := goquery.NewDocument(url)
	html, _ := doc.Html()
	title := strings.TrimSpace(doc.Find(".tgme_page_title").Text())
	description := strings.TrimSpace(doc.Find(".tgme_page_description").Text())
	imgSrc, _ := doc.Find(".tgme_page_photo_image").Attr("src")
	fmt.Printf("\n title: %s \n description: %s \n src: %s \n", title, description, imgSrc)
}
