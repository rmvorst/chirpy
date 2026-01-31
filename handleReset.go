package main

import (
	"log"
	"net/http"
)

func (cfg *apiConfig) handlerResetServerHits(w http.ResponseWriter, req *http.Request) {
	err := cfg.db.DeleteUsers(req.Context())
	if err != nil {
		log.Printf("Error deleting users: %s", err)
	}

	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}
