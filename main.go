package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// Define structs to match the structure of the RSS feed
type Rss struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Items       []Item `xml:"item"`
}

type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	GUID        string `xml:"guid"`
	PubDate     string `xml:"pubDate"`
}

func fetchRssFeed(url string) (*Rss, error) {
	// Make an HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Unmarshal the XML into the Rss struct
	var rss Rss
	if err := xml.Unmarshal(body, &rss); err != nil {
		return nil, err
	}

	return &rss, nil
}

func main() {
	url := "https://fortiguard.fortinet.com/rss/ir.xml" // Your RSS feed URL
	rss, err := fetchRssFeed(url)
	if err != nil {
		fmt.Println("Error fetching RSS feed:", err)
		return
	}

	// Custom layout to match the pubDate format including the numeric timezone
	const layout = "Mon, 02 Jan 2006 15:04:05 -0700"
	currentTime := time.Now()

	for _, item := range rss.Channel.Items {
		pubDate, err := time.Parse(layout, item.PubDate)
		if err != nil {
			fmt.Println("Error parsing date:", err)
			continue
		}

		if currentTime.Sub(pubDate).Hours() <= 168 { // 168 hours in 7 days
			fmt.Printf("Title: %s\nLink: %s\nDescription: %s\nGUID: %s\nPubDate: %s\n\n", item.Title, item.Link, item.Description, item.GUID, item.PubDate)
		}
	}
}
