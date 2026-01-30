package main

import "net/http"

func (cfg *apiConfig) handlerResetServerHits(w http.ResponseWriter, req *http.Request) {
	cfg.fileserverHits.Store(0)
}
