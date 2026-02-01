package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rmvorst/chirpy/internal/auth"
	"github.com/rmvorst/chirpy/internal/database"
)

type createChirpRequest struct {
	Body   string    `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handleCreateChirp(w http.ResponseWriter, req *http.Request) {
	type validResponse struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}

	const maxChirpLength = 140

	userID, err := validateUser(cfg, req)
	if err != nil {
		postErr := errorResponse{Err: fmt.Sprintf("%s", err)}
		postJSON(postErr, http.StatusUnauthorized, w)
		return
	}

	decoder := json.NewDecoder(req.Body)
	params := createChirpRequest{}
	err = decoder.Decode(&params)
	if err != nil {
		postErr := errorResponse{Err: fmt.Sprintf("%s", err)}
		postJSON(postErr, http.StatusInternalServerError, w)
		return
	}

	err = checkChirpLength(maxChirpLength, params)
	if err != nil {
		postErr := errorResponse{Err: fmt.Sprintf("%s", err)}
		postJSON(postErr, http.StatusBadRequest, w)
		return
	}

	chirp, err := generateChirp(params, userID, cfg, req)
	if err != nil {
		postErr := errorResponse{Err: fmt.Sprintf("%s", err)}
		postJSON(postErr, http.StatusInternalServerError, w)
		return
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

func validateUser(cfg *apiConfig, req *http.Request) (uuid.UUID, error) {
	authToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		log.Println("Error in validateUser: Failed to validate user")
		return uuid.Nil, fmt.Errorf("Authentication error")
	}

	validID, err := auth.ValidateJWT(authToken, cfg.secret)
	if err != nil {
		log.Println("Error in validateUser: Failed to validate user")
		return uuid.Nil, fmt.Errorf("Authentication error")
	}
	return validID, nil
}

func checkChirpLength(maxChirpLength int, params createChirpRequest) error {
	if len(params.Body) > maxChirpLength {
		return fmt.Errorf("Chirp is too long")
	}
	return nil
}

func generateChirp(params createChirpRequest, userID uuid.UUID, cfg *apiConfig, req *http.Request) (database.Chirp, error) {
	cleanedBody := cleanBody(params.Body)

	chirpParams := database.CreateChirpParams{
		Body:   cleanedBody,
		UserID: userID,
	}
	chirp, err := cfg.db.CreateChirp(req.Context(), chirpParams)
	if err != nil {
		return database.Chirp{}, fmt.Errorf("Error creating chirp")
	}
	return chirp, nil
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
