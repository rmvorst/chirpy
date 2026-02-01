package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/rmvorst/chirpy/internal/database"
	"github.com/rmvorst/chirpy/internal/database/auth"
)

type parameters struct {
	Body   string    `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handleCreateChirp(w http.ResponseWriter, req *http.Request) {
	const maxChirpLength = 140

	params, err := decodeJson(req)
	if err != nil {
		postErr := writeErrorPost(err)
		postJSON(postErr, http.StatusInternalServerError, w)
		return
	}

	userID, err := validateUser(cfg, req)
	if err != nil {
		postErr := writeErrorPost(err)
		postJSON(postErr, http.StatusUnauthorized, w)
		return
	}

	err = checkChirpLength(maxChirpLength, params)
	if err != nil {
		postErr := writeErrorPost(err)
		postJSON(postErr, http.StatusBadRequest, w)
		return
	}

	chirp, err := generateChirp(params, userID, cfg, req)
	if err != nil {
		postErr := writeErrorPost(err)
		postJSON(postErr, http.StatusInternalServerError, w)
		return
	}

	postValid := writeValidPost(chirp)
	postJSON(postValid, http.StatusCreated, w)
}

func writeErrorPost(err error) errorResponse {
	postErr := errorResponse{
		Err: fmt.Sprintf("%s", err),
	}
	return postErr
}

func writeValidPost(chirp database.Chirp) validResponse {
	postValid := validResponse{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}
	return postValid
}

func decodeJson(req *http.Request) (parameters, error) {
	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	return params, err
}

func validateUser(cfg *apiConfig, req *http.Request) (uuid.UUID, error) {
	authToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		return uuid.Nil, fmt.Errorf("Authentication error\n")
	}

	validID, err := auth.ValidateJWT(authToken, cfg.secret)
	if err != nil {
		return uuid.Nil, fmt.Errorf("Authentication error\n")
	}
	return validID, nil
}

func checkChirpLength(maxChirpLength int, params parameters) error {
	if len(params.Body) > maxChirpLength {
		return fmt.Errorf("Chirp is too long\n")
	}
	return nil
}

func generateChirp(params parameters, userID uuid.UUID, cfg *apiConfig, req *http.Request) (database.Chirp, error) {
	cleanedBody := cleanBody(params.Body)

	chirpParams := database.CreateChirpParams{
		Body:   cleanedBody,
		UserID: userID,
	}
	chirp, err := cfg.db.CreateChirp(req.Context(), chirpParams)
	if err != nil {
		return database.Chirp{}, fmt.Errorf("Error creating chirp\n")
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
