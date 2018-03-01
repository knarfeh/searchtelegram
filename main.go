package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func getTgType(content string) string {
	if strings.Contains(content, "GROUP") {
		return "group"
	}
	if strings.Contains(content, "CHANNEL") {
		return "channel"
	}
	return "unknown"
}

func downloadPic(imgSrc, tgID string) bool {
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
func uploadPic2S3(tgID string) {
	awsAccessKeyID := ""
	awsSecretAccessKey := ""
	token := ""
	creds := credentials.NewStaticCredentials(awsAccessKeyID, awsSecretAccessKey, token)
	_, err := creds.Get()
	if err != nil {
		fmt.Printf("bad credentials: %s", err)
	}
	cfg := aws.NewConfig().WithRegion("us-east-1").WithCredentials(creds)
	svc := s3.New(session.New(), cfg)

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
	resp, err := svc.PutObject(params)
	if err != nil {
		fmt.Printf("bad response: %s", err)
	}
	fmt.Printf("response %s", awsutil.StringValue(resp))
}

func ExampleScrape() {
	tgID := "knarfeh"
	doc, err := goquery.NewDocument("http://t.me/" + tgID)
	if err != nil {
		log.Fatal(err)
	}

	// Find the review items
	doc.Find(".sidebar-reviews article .content-block").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		band := s.Find("a").Text()
		title := s.Find("i").Text()
		fmt.Printf("Review %d: %s - %s\n", i, band, title)
	})

	title := strings.TrimSpace(doc.Find(".tgme_page_title").Text())
	description := strings.TrimSpace(doc.Find(".tgme_page_description").Text())
	imgSrc, _ := doc.Find(".tgme_page_photo_image").Attr("src")
	buttonContent := doc.Find(".tgme_action_button_new").Text()
	tgType := getTgType(strings.ToUpper(buttonContent))
	if tgType == "unknown" {
		if strings.HasSuffix(strings.ToUpper("knarfeh"), "BOT") {
			tgType = "bot"
		} else {
			tgType = "people"
		}
	}
	fmt.Printf("\n\n\n title: %s, description: %s, src: %s, tgType: %s", title, description, imgSrc, tgType)

	imgPath := ""
	if imgSrc != "" {
		downloadPic(imgSrc, tgID)
		uploadPic2S3(tgID)
		imgPath = "https://s3.amazonaws.com/searchtelegram/media/images/" + tgID + ".jpg"
		// imgPath = "/media/images/" + "telegram" + ".jpg"
	} else {
		imgPath = "https://s3.amazonaws.com/searchtelegram/media/images/" + tgID + ".jpg"
	}
	fmt.Printf("imgPath???%s", imgPath)
	if title == "" && strings.HasPrefix(description, "If you have Telegram, you can contact") {
		fmt.Printf("Oh, tragedy")
	}

	fmt.Println("Success!")
}

func main() {
	ExampleScrape()
}
