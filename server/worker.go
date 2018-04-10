package main

import (
	"bytes"
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
	"github.com/RedisLabs/redisearch-go/redisearch"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/knarfeh/searchtelegram/server/domain"
	"github.com/olebedev/config"

	tb "github.com/tucnak/telebot"
	elastic "gopkg.in/olivere/elastic.v5"
)

// Hauler ...
type Hauler struct {
	esClient         *ESClient
	redisClient      *RedisClient
	redisearchClient *RedisearchClient
	s3Client         *S3Client
	tb               *tb.Bot
	conf             *config.Config
}

type tgMeInfo struct {
	Title       string
	Description string
	ImgSrc      string
	Type        string
}

// CreateConsumer create consumer ...
func CreateConsumer(conf *config.Config) (*Hauler, error) {
	ESHOSTPORT, _ := conf.String("ESHOSTPORT")
	REDISHOST, _ := conf.String("REDISHOST")
	REDISPORT, _ := conf.String("REDISPORT")
	AWSACCESSKEY, _ := conf.String("AWSACCESSKEY")
	AWSSECRETKEY, _ := conf.String("AWSSECRETKEY")
	AWSREGION, _ := conf.String("AWSREGION")
	TGBOTTOKEN, _ := conf.String("TGBOTTOKEN")
	es, _ := NewESClient(
		ElasticConfig{
			Endpoint:           ESHOSTPORT,
			Username:           "",
			Password:           "",
			Retries:            3,
			HealthCheckTimeout: 3 * time.Second,
		},
	)
	redisClient := NewRedisClient(REDISHOST, REDISPORT)
	redisearchClient := NewRedisearchClient(REDISHOST, REDISPORT)
	s3Client := NewS3Client(strings.TrimSpace(AWSACCESSKEY), strings.TrimSpace(AWSSECRETKEY), strings.TrimSpace(AWSREGION))

	b, err := tb.NewBot(tb.Settings{
		Token:  strings.TrimSpace(TGBOTTOKEN),
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		fmt.Println(err)
	}

	return &Hauler{
		esClient:         es,
		redisClient:      redisClient,
		redisearchClient: redisearchClient,
		s3Client:         s3Client,
		tb:               b,
		conf:             conf,
	}, nil
}

// Submit2ES subscribe redis channel st_submit, get data from t.me, save it to es
func (hauler *Hauler) Submit2ES() {
	pubsub := hauler.redisClient.Client.Subscribe("st_submit")
	defer pubsub.Close()

	substr, err := pubsub.ReceiveTimeout(time.Second)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Redis subscribe: %s\n", substr)

	ch := pubsub.Channel()
	for {
		select {
		case <-ch:
			// hauler.handleSubmit(msg.Payload)
			fmt.Println("Got submit")
			submitStr := hauler.redisClient.Client.RPop("st_submit_list").Val()
			if submitStr != "" {
				fmt.Printf("Got submit str from list: %s\n", submitStr)
				hauler.handleSubmit(submitStr)
			}
		}
	}
}

// handleSubmit ...
func (hauler *Hauler) handleSubmit(submitStr string) {
	tgResource := domain.NewTgResource()
	if err := json.Unmarshal([]byte(submitStr), &tgResource); err != nil {
		panic(err)
	}

	fmt.Printf("Got TgID: %s\n", tgResource.TgID)
	tgResource.TgID = hauler.getRealtgID(tgResource.TgID)
	tgID := tgResource.TgID
	tgInfo, err := hauler.getData(tgID)
	if err != nil {
		fmt.Printf("Got error when getting data from t.me, error: %s\n", err.Error())
		return
	}
	fmt.Printf("tgInfo, tgID: %s, tgType: %s\n", tgInfo.Title, tgInfo.Type)

	tgResource.Title = tgInfo.Title
	tgResource.Type = tgInfo.Type
	tagItem := &domain.Tag{
		Count: 1,
		Name:  tgInfo.Type,
	}
	tgResource.Tags = append(tgResource.Tags, *tagItem)
	tgResource.Tags = hauler.rmDuplicateTags(tgResource.Tags)
	tgResource.Imgsrc = tgInfo.ImgSrc

	if tgResource.Desc == "" {
		fmt.Println("No description input, got from tdotme")
		tgResource.Desc = tgInfo.Description
	}
	_, err = hauler.esClient.Client.Index().OpType("create").Index("telegram").Type("resource").Id(tgID).BodyJson(tgResource).Do(context.TODO())
	hauler.redisClient.Client.Set("tgid:"+tgID, "1", 0)
	doc := redisearch.NewDocument(tgID, float32(1)/float32(10)).Set("tgid", tgID).Set("desc", tgResource.Desc).Set("title", tgResource.Title).Set("type", tgResource.Type).Set("tags", Tags2String(tgResource.Tags)).Set("tagsforsearch", Tags2String(tgResource.Tags))
	if err := hauler.redisearchClient.Client.IndexOptions(redisearch.DefaultIndexingOptions, doc); err != nil {
		fmt.Println(err)
	}
	hauler.send2stChannel(*tgResource)

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

// getRealtgID ...
func (hauler *Hauler) getRealtgID(tgID string) string {
	if strings.HasPrefix(tgID, "https://t.me/") {
		return tgID[13:]
	} else if strings.HasPrefix(tgID, "https://telegram.me/") {
		return tgID[20:]
	} else if strings.HasPrefix(tgID, "@") {
		return tgID[1:]
	}
	return tgID
}

// send2stChannel ...
func (hauler *Hauler) send2stChannel(tgResource domain.TgResource) {
	channelName, _ := hauler.conf.String("TGCHANNELNAME")
	stChannel := &tb.Chat{
		Type:     tb.ChatChannel,
		Username: channelName,
	}

	message, emoji := TgResource2Str(tgResource)
	channelMessage := "\n🆕 " + emoji + "\n \n @" + message
	hauler.tb.Send(stChannel, channelMessage)
}

// rmDuplicateTags remove duplicate tags
func (hauler *Hauler) rmDuplicateTags(tags []domain.Tag) []domain.Tag {
	keys := make(map[string]bool)
	list := []domain.Tag{}
	for _, entry := range tags {
		if _, value := keys[entry.Name]; !value {
			keys[entry.Name] = true
			list = append(list, entry)
		}
	}
	return list
}

// getTgType ...
func (hauler *Hauler) getTgType(content, tgID string) string {
	if strings.Contains(content, "GROUP") {
		return "group"
	}
	if strings.Contains(content, "CHANNEL") {
		return "channel"
	}
	if strings.HasSuffix(tgID, "BOT") {
		return "bot"
	}
	return "people"
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

// https://medium.com/@questhenkart/s3-image-uploads-via-aws-sdk-with-golang-63422857c548
// uploadPic2S3 upload picture to s3
func (hauler *Hauler) uploadPic2S3(tgID string) {
	file, err := os.Open("/tmp/images/" + tgID + ".jpg") // DEBUG mode
	if err != nil {
		fmt.Printf("err opening file: %s", err)
	}
	defer file.Close()
	fileInfo, _ := file.Stat()
	size := fileInfo.Size()
	buffer := make([]byte, size) // read file content to buffer

	file.Read(buffer)
	fileBytes := bytes.NewReader(buffer)
	fileType := http.DetectContentType(buffer)
	path := "/images/" + tgID + ".jpg"
	params := &s3.PutObjectInput{
		ACL:           aws.String("public-read"),
		Bucket:        aws.String("searchtelegram"),
		Key:           aws.String(path),
		Body:          fileBytes,
		ContentLength: aws.Int64(size),
		ContentType:   aws.String(fileType),
	}
	resp, err := hauler.s3Client.Client.PutObject(params)
	if err != nil {
		fmt.Printf("bad response: %s", err)
	}
	fmt.Printf("Response from s3 %s\n", awsutil.StringValue(resp))
}

// getData ...
func (hauler *Hauler) getData(tgID string) (*tgMeInfo, error) {
	url := "https://t.me/" + tgID

	tgIDExist, _ := hauler.redisClient.Client.Get("tgid:" + tgID).Result()
	if tgIDExist != "" {
		return nil, errors.New("Already exist in redis")
	}

	fmt.Printf("\nGetting tgID, url: %s \n", url)
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return nil, errors.New("got error from tdotme")
	}

	title := strings.TrimSpace(doc.Find(".tgme_page_title").Text())
	description := strings.TrimSpace(doc.Find(".tgme_page_description").Text())
	imgSrc, _ := doc.Find(".tgme_page_photo_image").Attr("src")
	buttonContent := doc.Find(".tgme_action_button_new").Text()
	tgType := hauler.getTgType(strings.ToUpper(buttonContent), strings.ToUpper(tgID))

	imgPath := ""
	if imgSrc != "" {
		hauler.downloadPic(imgSrc, tgID)
		hauler.uploadPic2S3(tgID)
		os.Remove("/tmp/images/" + tgID + ".jpg")
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

// Search2Redisearch subscribe redis channel st_search
func (hauler *Hauler) Search2Redisearch() {
	pubsub := hauler.redisClient.Client.Subscribe("st_search")
	defer pubsub.Close()

	substr, err := pubsub.ReceiveTimeout(time.Second)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Redis subscribe: %s\n", substr)

	ch := pubsub.Channel()
	for {
		select {
		case msg := <-ch:
			hauler.handleSearch(msg.Payload)
		}
	}
}

// handleSearch ...
func (hauler *Hauler) handleSearch(searchStr string) {
	fmt.Printf("Handle search string: %s\n", searchStr)
	val := hauler.redisClient.Client.SIsMember("redisearch:cached-search-string", searchStr).Val()
	if val {
		fmt.Printf("search string: %s already in cache set\n", searchStr)
		return
	}

	size := 500 // change to 5 when debug
	simpleQuery, _, _, _ := BuildESQuery(searchStr)

	search := hauler.esClient.Client.Search().Index("telegram").Type("resource").Size(size)
	searchResult, _ := search.Query(simpleQuery).Do(context.TODO())
	if searchResult.TotalHits() < int64(size) {
		size = int(searchResult.TotalHits())
	}
	hauler.hits2redisearch(*searchResult.Hits, size)

	for from := size; int64(from) < searchResult.TotalHits(); {
		search = hauler.esClient.Client.Search().Index("telegram").Type("resource").From(from).Size(size)
		searchResult, _ = search.Query(simpleQuery).Do(context.TODO())
		nowSize := size
		if searchResult.TotalHits()-int64(from) < int64(size) {
			nowSize = int(searchResult.TotalHits() - int64(from))
		}
		hauler.hits2redisearch(*searchResult.Hits, nowSize)
		from = from + size
	}
	hauler.redisClient.Client.SAdd("redisearch:cached-search-string", searchStr)
}

// hits2redisearch ...
func (hauler *Hauler) hits2redisearch(hits elastic.SearchHits, size int) {
	docs := make([]redisearch.Document, size)
	for i, hit := range hits.Hits {
		instance := domain.NewTgResource()
		json.Unmarshal(*hit.Source, instance)
		docs[i] = redisearch.NewDocument(instance.TgID, float32(hits.TotalHits-int64(i))/float32(hits.TotalHits)).Set("tgid", instance.TgID).Set("desc", instance.Desc).Set("title", instance.Title).Set("type", instance.Type).Set("tags", Tags2String(instance.Tags)).Set("tagsforsearch", Tags2String(instance.Tags))
		if i >= size {
			break
		}
	}
	if err := hauler.redisearchClient.Client.IndexOptions(redisearch.DefaultIndexingOptions, docs...); err != nil {
		fmt.Println(err)
	}
}

// Scheduler ...
func (hauler *Hauler) Scheduler(refresh time.Duration) {
	for i := 1; ; i++ {
		fmt.Println("wait 30 seconds")

		select {
		case <-time.After(refresh):
		}
	}
}
