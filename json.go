package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// this function responds with a valid json response
func jsonResponse(w http.ResponseWriter, statusCode int, payload interface{}) {
	// marshalls the response from the json
	res, err := json.Marshal(payload)

	// problem marshalling the data
	if err != nil {
		log.Printf("Error marshalling the json %s", err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(statusCode)
	w.Write(res)

}

// this function responds with a json error
func jsonError(w http.ResponseWriter, statusCode int, msg string) {
	// problem with the server
	if statusCode > 499 {
		log.Printf("Responding with 5xx status error: %s", msg)
	}

	// error response type
	type ErrorResp struct {
		Error string `json:"error"`
	}

	// calls the json response with a new error
	jsonResponse(w, statusCode, ErrorResp{
		Error: msg,
	})
}
