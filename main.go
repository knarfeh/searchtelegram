package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
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

func ExampleScrape() {
	doc, err := goquery.NewDocument("http://t.me/asdfasdffffffffffffff")
	html, _ := doc.Html()
	fmt.Printf("doc??? %s", html)
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
		downloadPic(imgSrc, "telegram")
		imgPath = "/media/images/" + "telegram" + ".jpg"
	} else {
		imgPath = "/media/images/" + "telegram" + ".jpg"
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
