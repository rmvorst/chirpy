package main

import (
	"encoding/json"
	"fmt"
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

	params, err := loginDecodeJSON(req)
	if err != nil {
		postErr := errorResponse{Err: fmt.Sprintf("%s", err)}
		postJSON(postErr, http.StatusInternalServerError, w)
		return
	}

	user, err := cfg.db.EmailLookup(req.Context(), params.Email)
	if err != nil {
		postErr := errorResponse{Err: "Incorrect email or password"}
		postJSON(postErr, http.StatusUnauthorized, w)
		return
	}

	match, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil || !match {
		postErr := errorResponse{Err: "Incorrect email or password"}
		postJSON(postErr, http.StatusUnauthorized, w)
		return
	}

	authTokenString, err := auth.MakeJWT(user.ID, cfg.secret, time.Duration(authExpiryTime)*time.Second)
	if err != nil {
		postErr := errorResponse{Err: fmt.Sprintf("%s", err)}
		postJSON(postErr, http.StatusInternalServerError, w)
		return
	}

	refreshTokenString, err := auth.MakeRefreshToken()
	if err != nil {
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
		postErr := errorResponse{Err: fmt.Sprintf("%s", err)}
		postJSON(postErr, http.StatusInternalServerError, w)
		return
	}

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

func loginDecodeJSON(req *http.Request) (loginRequest, error) {
	decoder := json.NewDecoder(req.Body)
	params := loginRequest{}
	err := decoder.Decode(&params)
	return params, err
}
