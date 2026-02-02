package main

import (
	"log"
	"net/http"
	"time"

	"github.com/rmvorst/chirpy/internal/auth"
)

func (cfg *apiConfig) handleRefresh(w http.ResponseWriter, req *http.Request) {
	type validTokenResponse struct {
		Token string `json:"token"`
	}

	refreshTokenString, err := auth.GetBearerToken(req.Header)
	if err != nil {
		log.Printf("Error in handleRefresh: Could not get bearer token")
		postErr := errorResponse{Err: "Authorization error"}
		postJSON(postErr, http.StatusUnauthorized, w)
		return
	}

	refreshToken, err := cfg.db.GetToken(req.Context(), refreshTokenString)
	if err != nil {
		log.Printf("Error in handleRefresh: Could not get user refresh token")
		postErr := errorResponse{Err: "Authorization error"}
		postJSON(postErr, http.StatusUnauthorized, w)
		return
	}

	authTokenString, err := auth.MakeJWT(refreshToken.UserID, cfg.secret, time.Duration(authExpiryTime)*time.Second)
	if err != nil {
		log.Printf("Error in handleRefresh: Could not make new user JWT")
		postErr := errorResponse{Err: "Authorization error"}
		postJSON(postErr, http.StatusUnauthorized, w)
		return
	}

	log.Printf("Success: User tokens refreshed")
	postValid := validTokenResponse{Token: authTokenString}
	postJSON(postValid, http.StatusOK, w)
}
