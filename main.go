package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/vicentefiorito/rss_feed_aggregator/internal/database"
)

// holds the stateful data
type apiConfig struct {
	DB *database.Queries
}

func main() {
	feed, err := fetchRssFeed("https://wagslane.dev/index.xml")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(feed)

	godotenv.Load(".env")

	// load the database
	dbUrl := os.Getenv("DB_CONNECTION")
	if dbUrl == "" {
		log.Fatal("DB_CONNECTION variable is not set!")
	}

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal("Problem connecting to the db!")
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("Port variable is not set!")
	}

	// stores the queries in a database package
	dbQueries := database.New(db)

	// config struct
	apiConfig := apiConfig{
		DB: dbQueries,
	}

	// creating a server
	mux := http.NewServeMux()

	// route handling begins here
	mux.HandleFunc("GET /v1/readiness", handleReadiness)
	mux.HandleFunc("GET /v1/err", handleErr)

	// user endpoints
	mux.HandleFunc("POST /v1/users", apiConfig.handleUserCreate)
	mux.HandleFunc("GET /v1/users", apiConfig.middlewareAuth(apiConfig.handleGetUser))

	// feed endpoints
	mux.HandleFunc("POST /v1/feeds", apiConfig.middlewareAuth(apiConfig.handleFeedCreate))
	mux.HandleFunc("GET /v1/feeds", apiConfig.handleGetFeeds)

	// feedFollow endpoints
	mux.HandleFunc("GET /v1/feed_follows", apiConfig.middlewareAuth(apiConfig.handleGetFeedFollows))
	mux.HandleFunc("POST /v1/feed_follows", apiConfig.middlewareAuth(apiConfig.handleCreateFeedFollow))
	mux.HandleFunc("DELETE /v1/feed_follows/{feedFollowID}", apiConfig.middlewareAuth(apiConfig.handleDeleteFeedFollow))

	s := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	log.Printf("Listening on port: %s\n", port)
	log.Fatal(s.ListenAndServe())
}
