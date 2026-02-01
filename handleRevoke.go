package main

import (
	"net/http"

	"github.com/rmvorst/chirpy/internal/auth"
)

func (cfg *apiConfig) handleRevoke(w http.ResponseWriter, req *http.Request) {
	authTokenString, err := auth.GetBearerToken(req.Header)
	if err != nil {
		postErr := errorResponse{Err: "Authorization error"}
		postJSON(postErr, http.StatusUnauthorized, w)
		return
	}

	_, err = cfg.db.RevokeToken(req.Context(), authTokenString)
	if err != nil {
		postErr := errorResponse{Err: "Authorization error"}
		postJSON(postErr, http.StatusUnauthorized, w)
		return
	}

	postJSON(nil, http.StatusNoContent, w)
}
