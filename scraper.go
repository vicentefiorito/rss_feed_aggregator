package main

import (
	"encoding/xml"
	"io"
	"net/http"
	"time"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Language    string    `xml:"language"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	PubDate     string `xml:"pubDate"`
	Description string `xml:"description"`
}

// function that accepts the url of a live rss feed
// and returns the parsed data into a go struct
func fetchRssFeed(url string) (*RSSFeed, error) {
	// create a http client to get the info from the feeds
	httpClient := http.Client{
		Timeout: 10 * time.Second,
	}

	// gets the data from the url
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// reads the data from the url
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// unmarshalling the xml
	rssFeed := RSSFeed{}
	err = xml.Unmarshal(data, &rssFeed)
	if err != nil {
		return nil, err
	}

	return &rssFeed, nil

}
