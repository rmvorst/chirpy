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

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, req *http.Request) {
	type userCreateRequest struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	type userCreateResponse struct {
		ID         uuid.UUID `json:"id"`
		Created_at time.Time `json:"created_at"`
		Updated_at time.Time `json:"updated_at"`
		Email      string    `json:"email"`
	}

	decoder := json.NewDecoder(req.Body)
	params := userCreateRequest{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Println("Error decoding incoming JSON")
		postErr := errorResponse{Err: fmt.Sprintf("%s\n", err)}
		postJSON(postErr, http.StatusInternalServerError, w)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		log.Println("Error hashing user password")
		postErr := errorResponse{Err: fmt.Sprintf("%s\n", err)}
		postJSON(postErr, http.StatusInternalServerError, w)
		return
	}

	createUserParameters := database.CreateUserParams{
		HashedPassword: hashedPassword,
		Email:          params.Email,
	}
	user, err := cfg.db.CreateUser(req.Context(), createUserParameters)
	if err != nil {
		log.Println("Error creating user")
		postErr := errorResponse{Err: fmt.Sprintf("%s\n", err)}
		postJSON(postErr, http.StatusInternalServerError, w)
		return
	}

	createdUser := userCreateResponse{
		ID:         user.ID,
		Created_at: user.CreatedAt,
		Updated_at: user.UpdatedAt,
		Email:      user.Email,
	}
	postJSON(createdUser, http.StatusCreated, w)
}

func (cfg *apiConfig) handleUpdateAccount(w http.ResponseWriter, req *http.Request) {
	type userUpdateRequest struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	type userUpdateResponse struct {
		ID         uuid.UUID `json:"id"`
		Created_at time.Time `json:"created_at"`
		Updated_at time.Time `json:"updated_at"`
		Email      string    `json:"email"`
	}

	userID, err := validateUser(cfg, req)
	if err != nil {
		postErr := errorResponse{Err: fmt.Sprintf("%s\n", err)}
		postJSON(postErr, http.StatusUnauthorized, w)
		return
	}

	decoder := json.NewDecoder(req.Body)
	params := userUpdateRequest{}
	err = decoder.Decode(&params)
	if err != nil {
		postErr := errorResponse{Err: fmt.Sprintf("%s\n", err)}
		postJSON(postErr, http.StatusInternalServerError, w)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		log.Println("Error hashing user password")
		postErr := errorResponse{Err: fmt.Sprintf("%s\n", err)}
		postJSON(postErr, http.StatusInternalServerError, w)
		return
	}

	updateUserParameters := database.UpdateUserParams{
		ID:             userID,
		Email:          params.Email,
		HashedPassword: hashedPassword,
	}
	user, err := cfg.db.UpdateUser(req.Context(), updateUserParameters)
	if err != nil {
		log.Println("Error updating user")
		postErr := errorResponse{Err: fmt.Sprintf("%s\n", err)}
		postJSON(postErr, http.StatusInternalServerError, w)
		return
	}

	createdUser := userUpdateResponse{
		ID:         user.ID,
		Created_at: user.CreatedAt,
		Updated_at: user.UpdatedAt,
		Email:      user.Email,
	}
	postJSON(createdUser, http.StatusOK, w)
}
