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
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/knarfeh/searchtelegram/server/domain"

	elastic "gopkg.in/olivere/elastic.v5"
)

// Hauler ...
type Hauler struct {
	esClient    *ESClient
	redisClient *RedisClient
	s3Client    *S3Client
}

type tgMeInfo struct {
	Title       string
	Description string
	ImgSrc      string
	Type        string
}

// CreateConsumer create consumer ...
func CreateConsumer(esURL, redisHost, redisPort, accessKey, secretKey, region string) (*Hauler, error) {
	es, _ := NewESClient(esURL, "", "", 3)
	redisClient := NewRedisClient(redisHost, redisPort)
	s3Client := NewS3Client(accessKey, secretKey, region)
	return &Hauler{
		esClient:    es,
		redisClient: redisClient,
		s3Client:    s3Client,
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

	if strings.HasPrefix(tgResource.TgID, "https://t.me/") {
		tgResource.TgID = tgResource.TgID[13:]
	} else if strings.HasPrefix(tgResource.TgID, "@") {
		tgResource.TgID = tgResource.TgID[1:]
	}
	tgID := tgResource.TgID
	tgInfo, err := hauler.getData(tgID)
	if err != nil {
		fmt.Printf("Got error when getting data from t.me, error: %s\n", err.Error())
		return
	}
	fmt.Printf("tgInfo: %s\n", tgInfo)

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
