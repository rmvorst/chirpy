package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/rmvorst/chirpy/internal/database"
)

func (cfg *apiConfig) handleCreateChirp(w http.ResponseWriter, req *http.Request) {
	const maxChirpLength = 140

	type parameters struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		postErr := errorResponse{
			Err: fmt.Sprintf("Error decoding chirp body.: %s\n", err),
		}
		postJSON(postErr, http.StatusInternalServerError, w)
		return
	}

	if len(params.Body) > maxChirpLength {
		postErr := errorResponse{
			Err: "Chirp is too long\n",
		}
		postJSON(postErr, http.StatusBadRequest, w)
		return
	}
	cleanedBody := cleanBody(params.Body)
	createChirpParameters := database.CreateChirpParams{
		Body:   cleanedBody,
		UserID: params.UserID,
	}
	chirp, err := cfg.db.CreateChirp(req.Context(), createChirpParameters)
	if err != nil {
		postErr := errorResponse{
			Err: "Error creating chirp\n",
		}
		postJSON(postErr, http.StatusInternalServerError, w)
	}
	postValid := validResponse{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}
	postJSON(postValid, http.StatusCreated, w)
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
