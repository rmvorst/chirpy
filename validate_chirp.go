package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type parameters struct {
	Body string `json:"body"`
}
type postVals struct {
	Body         string `json:"body"`
	Cleaned_body string `json:"cleaned_body"`
	Valid        bool   `json:"valid"`
}
type errVals struct {
	Err string `json:"error"`
}

var profaneWords = []string{"kerfuffle", "sharbert", "fornax"}

func postError(resp errVals, w http.ResponseWriter) {
	dat, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(400)
	w.Write(dat)
}

func postValid(resp postVals, w http.ResponseWriter) {
	dat, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(dat)
}

func cleanBody(body string) string {
	bodyWords := strings.Fields(body)
	for idx, substring := range bodyWords {
		for _, profaneWord := range profaneWords {
			if strings.ToLower(substring) == profaneWord {
				bodyWords[idx] = "****"
			}
		}
	}
	return strings.Join(bodyWords, " ")
}

func validateChirp(w http.ResponseWriter, req *http.Request) {

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}

	if len(params.Body) > 140 {
		postErr := errVals{
			Err: "Chirp is too long",
		}
		postError(postErr, w)
		return
	}

	postVal := postVals{
		Body:         params.Body,
		Cleaned_body: cleanBody(params.Body),
		Valid:        true,
	}
	postValid(postVal, w)
}
