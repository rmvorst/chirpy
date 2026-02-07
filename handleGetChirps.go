package main

import (
	"fmt"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/rmvorst/chirpy/internal/database"
)

type validResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) getChirps(w http.ResponseWriter, req *http.Request) {
	var chirps []database.Chirp
	var err error
	authorID := req.URL.Query().Get("author_id")
	if authorID == "" {
		chirps, err = cfg.db.ReturnChirps(req.Context())
		if err != nil {
			log.Println("Error in getChirps: Cannot get chirps")
			postErr := errorResponse{Err: fmt.Sprintf("%s\n", err)}
			postJSON(postErr, http.StatusInternalServerError, w)
			return
		}
	} else {
		authorUUID, err := uuid.Parse(authorID)
		if err != nil {
			log.Println("Error in getChirps: Could not parse author id")
			postErr := errorResponse{Err: fmt.Sprintf("%s\n", err)}
			postJSON(postErr, http.StatusNotFound, w)
			return
		}
		chirps, err = cfg.db.ReturnUserChirps(req.Context(), authorUUID)
		if err != nil {
			log.Println("Error in getChirps: Cannot get chirps for user:", authorUUID)
			postErr := errorResponse{Err: fmt.Sprintf("%s\n", err)}
			postJSON(postErr, http.StatusInternalServerError, w)
			return
		}
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
	sortOrder := req.URL.Query().Get("sort")
	if sortOrder == "desc" {
		sort.Slice(postChirps, func(i, j int) bool { return postChirps[i].CreatedAt.After(postChirps[j].CreatedAt) })
	} else {
		sort.Slice(postChirps, func(i, j int) bool { return postChirps[i].CreatedAt.Before(postChirps[j].CreatedAt) })
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
