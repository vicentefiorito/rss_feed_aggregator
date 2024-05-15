package main

import (
	"net/http"

	"github.com/vicentefiorito/rss_feed_aggregator/internal/database"
)

// this function gets all the posts by an user
func (cfg *apiConfig) handleGetPostByUser(w http.ResponseWriter, r *http.Request, user database.User) {
	posts, err := cfg.DB.GetPostForUser(r.Context(), database.GetPostForUserParams{
		UserID: user.ID,
		Limit:  10,
	})

	if err != nil {
		jsonError(w, http.StatusInternalServerError, "couldn't get posts")
		return
	}

	// return the valid json
	jsonResponse(w, http.StatusOK, databasePostsToPosts(posts))
}
