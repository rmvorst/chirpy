package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/rmvorst/chirpy/internal/auth"
	"github.com/rmvorst/chirpy/internal/database"
)

type loginRequest struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}
type loginValidUser struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
}

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, req *http.Request) {

	decoder := json.NewDecoder(req.Body)
	params := loginRequest{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Println("Error in handlerLogin: Could not decode JSON")
		postErr := errorResponse{Err: fmt.Sprintf("%s", err)}
		postJSON(postErr, http.StatusInternalServerError, w)
		return
	}

	user, err := cfg.db.EmailLookup(req.Context(), params.Email)
	if err != nil {
		log.Println("Error in handlerLogin: Could not lookup user email:", params.Email)
		postErr := errorResponse{Err: "Incorrect email or password"}
		postJSON(postErr, http.StatusUnauthorized, w)
		return
	}

	match, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil || !match {
		log.Println("Error in handlerLogin: Failed to check password hash for user:", user.ID)
		postErr := errorResponse{Err: "Incorrect email or password"}
		postJSON(postErr, http.StatusUnauthorized, w)
		return
	}

	authTokenString, err := auth.MakeJWT(user.ID, cfg.secret, time.Duration(authExpiryTime)*time.Second)
	if err != nil {
		log.Println("Error in handlerLogin: Failed to make JWT for user:", user.ID)
		postErr := errorResponse{Err: fmt.Sprintf("%s", err)}
		postJSON(postErr, http.StatusInternalServerError, w)
		return
	}

	refreshTokenString, err := auth.MakeRefreshToken()
	if err != nil {
		log.Println("Error in handlerLogin: Failed to make refresh token for user:", user.ID)
		postErr := errorResponse{Err: fmt.Sprintf("%s", err)}
		postJSON(postErr, http.StatusInternalServerError, w)
		return
	}

	refreshTokenParams := database.CreateTokenParams{
		Token:  refreshTokenString,
		UserID: user.ID,
	}
	_, err = cfg.db.CreateToken(req.Context(), refreshTokenParams)
	if err != nil {
		log.Println("Error in handlerLogin: Failed to add refresh token to server for user:", user.ID)
		postErr := errorResponse{Err: fmt.Sprintf("%s", err)}
		postJSON(postErr, http.StatusInternalServerError, w)
		return
	}

	log.Println("Success: logged in user:", user.ID)
	postUser := loginValidUser{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		Token:        authTokenString,
		RefreshToken: refreshTokenString,
	}
	postJSON(postUser, http.StatusOK, w)
}
