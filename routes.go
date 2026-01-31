package main

import "net/http"

func (cfg *apiConfig) routes(fpathRoot string) *http.ServeMux {
	mux := http.NewServeMux()
	fileServerHandler := http.StripPrefix("/app", http.FileServer(http.Dir(fpathRoot)))
	mux.Handle("/app/", cfg.serverHitInc(fileServerHandler))

	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("POST /api/validate_chirp", validateChirp)

	mux.HandleFunc("GET /admin/metrics", cfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", cfg.handlerResetServerHits)

	return mux
}
