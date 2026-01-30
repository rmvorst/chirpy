package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/rmvorst/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries      *database.Queries
}

func errorHandle(introString string, err error) {
	if err != nil {
		fmt.Printf("%s: %v\n", introString, err)
	}
}

func main() {
	godotenv.Load()
	const filepathRoot = "."
	const port = "8080"

	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		errorHandle("Issue Connecting to Database", err)
	}
	dbQueries := database.New(db)
	apiCfg := &apiConfig{
		dbQueries: dbQueries,
	}

	serverMux := http.NewServeMux()
	server := &http.Server{
		Addr:    ":" + port,
		Handler: serverMux,
	}
	fileServerHandler := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))
	serverMux.Handle("/app/", apiCfg.serverHitInc(fileServerHandler))
	serverMux.HandleFunc("GET /api/healthz", handlerReadiness)
	serverMux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	serverMux.HandleFunc("POST /admin/reset", apiCfg.handlerResetServerHits)
	serverMux.HandleFunc("POST /api/validate_chirp", validateChirp)
	err = server.ListenAndServe()
	errorHandle("Server Error", err)
}
