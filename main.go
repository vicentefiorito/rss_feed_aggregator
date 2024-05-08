package main

import (
	"database/sql"
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
	// load the database
	dbUrl := os.Getenv("DB_CONNECTION")
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal("Problem connecting to the db!")
	}
	godotenv.Load(".env")
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

	s := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	log.Printf("Listening on port: %s\n", port)
	log.Fatal(s.ListenAndServe())
}
