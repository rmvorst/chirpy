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
	}
	type validUser struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
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

	postUser := validUser{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}
	postJSON(postUser, http.StatusOK, w)

}
