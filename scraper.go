package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/vicentefiorito/rss_feed_aggregator/internal/database"
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

// this function is going to run concurrently
// along with our server
func startScraping(db *database.Queries,
	concurrency int,
	timeBetweenRequest time.Duration,
) {
	log.Printf("Scraping on %v goroutines every %s duration", concurrency, timeBetweenRequest)
	// ticker executes based on timeBetweenRequest
	ticker := time.NewTicker(timeBetweenRequest)
	// for loop made this way so it automatically executes
	// the first time
	for ; ; <-ticker.C {
		// fetch the feeds
		feeds, err := db.GetNextFeedsToFetch(context.Background(), int32(concurrency))
		if err != nil {
			log.Print("Couldn't get next feeds to fetch", err)
			continue
		}

		// fetches each individual feed at the same time
		// uses synchronization
		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1) //adds 1 go routine for each feed

			// this calls the function concurrently
			// for each feed
			go scrapeFeed(db, wg, feed)
		}
		wg.Wait()

	}

}

// this function is going to take information
// from each individual feed
func scrapeFeed(
	db *database.Queries,
	wg *sync.WaitGroup,
	feed database.Feed,
) {
	defer wg.Done()

	// mark the feed as fetched
	_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Println("Eror marking feed as fetched:", err)
		return
	}

	// go fetch the feed
	feedData, err := fetchRssFeed(feed.Url)
	if err != nil {
		log.Println("Error fetching feed: ", err)
		return
	}

	// iterating through all items of the feed
	// and saving them into the table
	for _, item := range feedData.Channel.Item {
		// creates a description to create a post
		description := sql.NullString{}
		if item.Description != "" {
			description.String = item.Description
			description.Valid = true
		}

		// parses the string into a pubdate for the post
		pubDate := sql.NullTime{}
		if t, err := time.Parse(time.RFC1123Z, item.PubDate); err == nil {
			pubDate = sql.NullTime{
				Time:  t,
				Valid: true,
			}
		}

		// creates the posts
		_, err := db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Url:         item.Link,
			Description: description,
			PublishedAt: pubDate,
			FeedID:      feed.ID,
		})

		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				continue
			}
			log.Println("Couldn't create post:", err)
			continue
		}
	}

	log.Printf("Feed %s collected, %v posts found", feed.Name, len(feedData.Channel.Item))

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
