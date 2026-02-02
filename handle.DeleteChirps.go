package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handleDeleteChirps(w http.ResponseWriter, req *http.Request) {
	userID, err := validateUser(cfg, req)
	if err != nil {
		log.Println("Error in handleDeleteChirps: Cannot authorize user")
		postErr := errorResponse{Err: fmt.Sprintf("%s\n", err)}
		postJSON(postErr, http.StatusUnauthorized, w)
		return
	}

	rawID := req.PathValue("chirp_id")
	chirpID, err := uuid.Parse(rawID)
	if err != err {
		log.Println("Error in handleDeleteChirps: Cannot parse chirp id:", rawID)
		postErr := errorResponse{Err: fmt.Sprintf("%s\n", err)}
		postJSON(postErr, http.StatusBadRequest, w)
		return
	}

	chirp, err := cfg.db.ReturnChirp(req.Context(), chirpID)
	if err != err {
		log.Println("Error in handleDeleteChirps: Cannot get chirp id:", rawID)
		postErr := errorResponse{Err: fmt.Sprintf("%s\n", err)}
		postJSON(postErr, http.StatusNotFound, w)
		return
	}
	if chirp.UserID != userID {
		log.Println("Error in handleDeleteChirps: User does not have authorization to delete chirp")
		postErr := errorResponse{Err: "Incorrect user\n"}
		postJSON(postErr, http.StatusForbidden, w)
		return
	}

	err = cfg.db.DeleteChirp(req.Context(), chirpID)
	if err != err {
		log.Println("Error in handleDeleteChirps: Could not delete chirp:", rawID)
		postErr := errorResponse{Err: fmt.Sprintf("%s\n", err)}
		postJSON(postErr, http.StatusNotFound, w)
		return
	}

	postJSON(validResponse{}, http.StatusNoContent, w)
	log.Println("Success: Chirp", rawID, "deleted")
}
