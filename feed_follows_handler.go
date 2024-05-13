package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/vicentefiorito/rss_feed_aggregator/internal/database"
)

// handles the creation of a feed follow object
func (cfg *apiConfig) handleCreateFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		FeedID uuid.UUID `json:"feed_id"`
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

	//  create the feed follow
	feedFollow, err := cfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		FeedID:    params.FeedID,
		UserID:    user.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		jsonError(w, http.StatusInternalServerError, "Couldn't create feed follow")
		fmt.Println(err)
		return
	}

	// return a valid json
	jsonResponse(w, http.StatusCreated, databaseFeedFollowToFeedFollow(feedFollow))
}
