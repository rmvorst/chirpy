package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/rmvorst/chirpy/internal/auth"
)

func (cfg *apiConfig) handleUpgradeToRed(w http.ResponseWriter, req *http.Request) {
	type accountUpgradeRequest struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	apiKeyString, err := auth.GetAPIKey(req.Header)
	if err != nil {
		log.Println("Error in handleUpgradeToRed: Could not get API Key")
		postErr := errorResponse{Err: fmt.Sprintf("%s", err)}
		postJSON(postErr, http.StatusUnauthorized, w)
		return
	}
	if apiKeyString != cfg.polka_key {
		log.Println("Error in handleUpgradeToRed: Received API Key does not match expected")
		postErr := errorResponse{Err: fmt.Sprintf("%s", err)}
		postJSON(postErr, http.StatusUnauthorized, w)
		return
	}

	decoder := json.NewDecoder(req.Body)
	params := accountUpgradeRequest{}
	err = decoder.Decode(&params)
	if err != nil {
		log.Println("Error in handleUpgradeToRed: Could not decode JSON")
		postErr := errorResponse{Err: fmt.Sprintf("%s", err)}
		postJSON(postErr, http.StatusBadRequest, w)
		return
	}

	if params.Event != "user.upgraded" {
		log.Println("Warning in handleUpgradeToRed: Handler called, but event is not user.upgraded - no action is being taken")
		postJSON(nil, http.StatusNoContent, w)
		return
	}

	userID, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		log.Println("Error in handleUpgradeToRed: Cannot parse user id:", params.Data.UserID)
		postErr := errorResponse{Err: fmt.Sprintf("%s\n", err)}
		postJSON(postErr, http.StatusBadRequest, w)
		return
	}

	_, err = cfg.db.UpgradeToRed(req.Context(), userID)
	if err != nil {
		log.Println("Error in handleUpgradeToRed: Cannot find user:", params.Data.UserID)
		postErr := errorResponse{Err: fmt.Sprintf("%s\n", err)}
		postJSON(postErr, http.StatusNotFound, w)
		return
	}

	log.Println("Success: User account successfully upgraded to Red:", userID)
	postJSON(nil, http.StatusNoContent, w)
}
