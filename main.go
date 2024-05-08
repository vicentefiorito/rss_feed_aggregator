package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("Port variable is not set!")
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
