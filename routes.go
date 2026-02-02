package main

import "net/http"

func (cfg *apiConfig) routes(fpathRoot string) *http.ServeMux {
	mux := http.NewServeMux()
	fileServerHandler := http.StripPrefix("/app", http.FileServer(http.Dir(fpathRoot)))
	mux.Handle("/app/", cfg.serverHitInc(fileServerHandler))

	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /api/chirps", cfg.getChirps)
	mux.HandleFunc("GET /api/chirps/{chirp_id}", cfg.getChirp)
	mux.HandleFunc("POST /api/chirps", cfg.handleCreateChirp)
	mux.HandleFunc("POST /api/users", cfg.handlerCreateUser)
	mux.HandleFunc("POST /api/login", cfg.handlerLogin)
	mux.HandleFunc("POST /api/refresh", cfg.handleRefresh)
	mux.HandleFunc("POST /api/revoke", cfg.handleRevoke)
	mux.HandleFunc("PUT /api/users", cfg.handleUpdateAccount)
	mux.HandleFunc("DELETE /api/chirps/{chirp_id}", cfg.handleDeleteChirps)

	mux.HandleFunc("GET /admin/metrics", cfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", cfg.handlerResetServerHits)

	return mux
}
