package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gocolly/colly"
)

// PictureInformation contains the information
// about the location.
type PictureInformation struct {
	Date  string `json:"date"`
	Link  string `json:"link"`
	Title string `json:"title"`
}

func writeToFile(fileName string, content string) {
	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := f.Write([]byte(content)); err != nil {
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	c := colly.NewCollector()

	var pictures []PictureInformation

	// Find and visit all links
	c.OnHTML("article", func(e *colly.HTMLElement) {
		var pictureURL string

		linkSrc, _ := e.DOM.Find("div").Find("a").Find("img").Attr("srcset")
		links := strings.Split(linkSrc, " ")

		// Refers to the sizes 300, 1024, 1920. We only care about 1920.
		if len(links) == 6 {
			pictureURL = links[len(links)-2]
		}

		// Add the data to a struct.
		newPicture := PictureInformation{
			Date:  e.DOM.Find(".date").Text(),
			Link:  pictureURL,
			Title: e.DOM.Find(".entry-title").Text(),
		}

		// Append the struct to the array.
		pictures = append(pictures, newPicture)

	})

	for i := 1; i < 5; i++ {
		if i == 1 {
			c.Visit("https://windows10spotlight.com/")
		} else {
			c.Visit("https://windows10spotlight.com/page/" + fmt.Sprintf("%d", i))
		}
		jsonMarshsal, err := json.Marshal(pictures)
		if err != nil {
			fmt.Println(err)
		}
		writeToFile("results.json", string(jsonMarshsal))
		pictures = nil
	}
}
