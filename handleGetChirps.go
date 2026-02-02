package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type validResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) getChirps(w http.ResponseWriter, req *http.Request) {
	chirps, err := cfg.db.ReturnChirps(req.Context())
	if err != nil {
		log.Println("Error in getChirps: Cannot get chirps")
		postErr := errorResponse{Err: fmt.Sprintf("%s\n", err)}
		postJSON(postErr, http.StatusInternalServerError, w)
		return
	}
	postChirps := []validResponse{}
	for _, chirp := range chirps {
		postChirps = append(postChirps, validResponse{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}
	log.Printf("Success: Retrieved chirps")
	postJSON(postChirps, http.StatusOK, w)
}

func (cfg *apiConfig) getChirp(w http.ResponseWriter, req *http.Request) {
	rawID := req.PathValue("chirp_id")
	chirpID, err := uuid.Parse(rawID)
	if err != err {
		log.Println("Error parsing chirp id:", rawID)
		postErr := errorResponse{Err: fmt.Sprintf("%s\n", err)}
		postJSON(postErr, http.StatusBadRequest, w)
		return
	}

	chirp, err := cfg.db.ReturnChirp(req.Context(), chirpID)
	if err != nil {
		log.Println("Error returning the chirp:", rawID)
		postErr := errorResponse{Err: fmt.Sprintf("%s\n", err)}
		postJSON(postErr, http.StatusNotFound, w)
		return
	}

	log.Println("Success: Returned chirp:", rawID)
	postChirp := validResponse{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}
	postJSON(postChirp, http.StatusOK, w)
}
