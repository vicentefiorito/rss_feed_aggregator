package main

import (
	"net/http"

	"github.com/vicentefiorito/rss_feed_aggregator/internal/database"
)

// this function gets all the posts by an user
func (cfg *apiConfig) handleGetPostByUser(w http.ResponseWriter, r *http.Request, user database.User) {

}
