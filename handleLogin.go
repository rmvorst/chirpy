package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/rmvorst/chirpy/internal/database/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, req *http.Request) {
	type requestBody struct {
		Password string `json:"password"`
		Email    string `json:"email"`
		Expires  int    `json:"expires_in_seconds"`
	}
	type validUser struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
		Token     string    `json:"token"`
	}

	decoder := json.NewDecoder(req.Body)
	params := requestBody{}
	err := decoder.Decode(&params)
	if err != nil {
		postErr := errorResponse{
			Err: "Error decoding JSON",
		}
		postJSON(postErr, http.StatusInternalServerError, w)
		return
	}
	if params.Expires <= 0 || params.Expires > 3600 {
		params.Expires = 3600
	}

	user, err := cfg.db.EmailLookup(req.Context(), params.Email)
	if err != nil {
		postErr := errorResponse{
			Err: "Incorrect email or password",
		}
		postJSON(postErr, http.StatusUnauthorized, w)
		return
	}

	match, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil || !match {
		postErr := errorResponse{
			Err: "Incorrect email or password",
		}
		postJSON(postErr, http.StatusUnauthorized, w)
		return
	}

	tokenString, err := auth.MakeJWT(user.ID, cfg.secret, time.Duration(params.Expires)*time.Second)

	postUser := validUser{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		Token:     tokenString,
	}
	postJSON(postUser, http.StatusOK, w)

}
