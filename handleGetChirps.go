package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) getChirps(w http.ResponseWriter, req *http.Request) {
	chirps, err := cfg.db.ReturnChirps(req.Context())
	if err != nil {
		postErr := errorResponse{
			Err: fmt.Sprintf("Error getting chirps: %s\n", err),
		}
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
	postJSON(postChirps, http.StatusOK, w)
}
