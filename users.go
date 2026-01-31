package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, req *http.Request) {
	type userCreateRequest struct {
		Email string `json:"email"`
	}
	type userCreateSuccess struct {
		ID         uuid.UUID `json:"id"`
		Created_at time.Time `json:"created_at"`
		Updated_at time.Time `json:"updated_at"`
		Email      string    `json:"email"`
	}

	decoder := json.NewDecoder(req.Body)
	params := userCreateRequest{}
	err := decoder.Decode(&params)
	if err != nil {
		postErr := errorResponse{
			Err: fmt.Sprintf("Error decoding user creation request: %s\n", err),
		}
		postJSON(postErr, http.StatusInternalServerError, w)
		return
	}

	user, err := cfg.db.CreateUser(req.Context(), params.Email)
	if err != nil {
		postErr := errorResponse{
			Err: fmt.Sprintf("Error creating new user: %s\n", err),
		}
		postJSON(postErr, http.StatusInternalServerError, w)
		return
	}
	createdUser := userCreateSuccess{
		ID:         user.ID,
		Created_at: user.CreatedAt,
		Updated_at: user.UpdatedAt,
		Email:      user.Email,
	}
	postJSON(createdUser, http.StatusCreated, w)

}
