package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type parameters struct {
	Body string `json:"body"`
}
type validResponse struct {
	Body         string `json:"body"`
	Cleaned_body string `json:"cleaned_body"`
	Valid        bool   `json:"valid"`
}
type errorResponse struct {
	Err string `json:"error"`
}

func validateChirp(w http.ResponseWriter, req *http.Request) {
	const maxChirpLength = 140

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		postErr := errorResponse{
			Err: fmt.Sprintf("Error decoding parameters: %s\n", err),
		}
		postJSON(postErr, http.StatusInternalServerError, w)
		return
	}

	if len(params.Body) > maxChirpLength {
		postErr := errorResponse{
			Err: "Chirp is too long",
		}
		postJSON(postErr, http.StatusBadRequest, w)
		return
	}

	postVal := validResponse{
		Body:         params.Body,
		Cleaned_body: cleanBody(params.Body),
		Valid:        true,
	}
	postJSON(postVal, http.StatusOK, w)
}

func cleanBody(body string) string {
	profaneWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	bodyWords := strings.Fields(body)
	for idx, substring := range bodyWords {
		if _, ok := profaneWords[strings.ToLower(substring)]; ok {
			bodyWords[idx] = "****"
		}
	}
	return strings.Join(bodyWords, " ")
}
