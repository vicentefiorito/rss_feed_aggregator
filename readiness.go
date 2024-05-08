package main

import "net/http"

// checks that the sever has a status code of 200

func handleReadiness(w http.ResponseWriter, r *http.Request) {
	type resp struct {
		Status string `json:"status"`
	}
	statusOK := "ok"
	jsonResponse(w, http.StatusOK, resp{
		Status: statusOK,
	})
}

// checks that the sever has a status code of 500
func handleErr(w http.ResponseWriter, r *http.Request) {
	jsonError(w, http.StatusInternalServerError, "Internal Server Error")
}
