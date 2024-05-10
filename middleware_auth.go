package main

import (
	"net/http"

	"github.com/vicentefiorito/rss_feed_aggregator/internal/auth"
	"github.com/vicentefiorito/rss_feed_aggregator/internal/database"
)

// type for handlers that require authentication
type authedHandler func(http.ResponseWriter, *http.Request, database.User)

// this middleware authenticates a request
// gets the user
// and calls the next authed handler
func (cfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get the api key from the header
		apiKey, err := auth.GetApiKey(r.Header)
		if err != nil {
			jsonError(w, http.StatusUnauthorized, "Couldn't find Api Key")
			return
		}

		user, err := cfg.DB.GetUserByApiKey(r.Context(), apiKey)
		if err != nil {
			jsonError(w, http.StatusInternalServerError, "Couldn't get user")
			return
		}

		handler(w, r, user)
	}
}
