package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/vicentefiorito/rss_feed_aggregator/internal/database"
)

// this struct holds the creation of a feed and a feed follow
type response struct {
	Feed       Feed       `json:"feed"`
	FeedFollow FeedFollow `json:"feed_follow"`
}

// this function creates a feed into the db
func (cfg *apiConfig) handleFeedCreate(w http.ResponseWriter, r *http.Request, user database.User) {
	// request body
	type parameters struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}

	// empty params to generate the response object
	params := parameters{}

	// decodes the json from the request
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, "Couldn't decode json")
		return
	}

	// creating the feed
	feed, err := cfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      params.Name,
		Url:       params.URL,
		UserID:    user.ID,
	})

	if err != nil {
		jsonError(w, http.StatusInternalServerError, "Couldn't create feed")
		fmt.Println(err)
		return
	}

	//  create the feed follow
	feedFollow, err := cfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		FeedID:    feed.ID,
		UserID:    user.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		jsonError(w, http.StatusInternalServerError, "Couldn't create feed follow")
		fmt.Println(err)
		return
	}

	jsonResponse(w, http.StatusCreated, response{
		Feed:       databaseFeedToFeed(feed),
		FeedFollow: databaseFeedFollowToFeedFollow(feedFollow),
	})

}

// handler that returns all the feeds created
func (cfg *apiConfig) handleGetFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := cfg.DB.GetFeeds(r.Context())
	if err != nil {
		jsonError(w, http.StatusInternalServerError, "Couldn't get feeds")
		return
	}
	jsonResponse(w, http.StatusOK, databaseFeedsToFeeds(feeds))
}
