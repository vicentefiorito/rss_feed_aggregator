package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/vicentefiorito/rss_feed_aggregator/internal/auth"
	"github.com/vicentefiorito/rss_feed_aggregator/internal/database"
)

// function that handles the creation of a user
func (cfg *apiConfig) handleUserCreate(w http.ResponseWriter, r *http.Request) {
	// request body
	type parameters struct {
		Name string `json:"name"`
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

	// creates a new user in the database
	user, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      params.Name,
	})

	if err != nil {
		jsonError(w, http.StatusInternalServerError, "Couldn't create user")
		fmt.Println(err)
		return
	}

	// Valid response an valid user created
	jsonResponse(w, http.StatusCreated, user)

}

// gets the user by apikey
func (cfg *apiConfig) handleGetUser(w http.ResponseWriter, r *http.Request) {
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

	// respond with a valid created user
	jsonResponse(w, http.StatusOK, databaseUserToUser(user))
}
