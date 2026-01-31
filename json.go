package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type validResponse struct {
	Body         string `json:"body"`
	Cleaned_body string `json:"cleaned_body"`
	Valid        bool   `json:"valid"`
}
type errorResponse struct {
	Err string `json:"error"`
}

func postJSON(resp interface{}, statusCode int, w http.ResponseWriter) {
	dat, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(dat)
}
