package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/rmvorst/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
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
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		errorHandle("Issue Connecting to Database", err)
	}
	dbQueries := database.New(db)
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
	}

	serverMux := http.NewServeMux()
	fileServerHandler := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))
	serverMux.Handle("/app/", apiCfg.serverHitInc(fileServerHandler))

	serverMux.HandleFunc("GET /api/healthz", handlerReadiness)
	serverMux.HandleFunc("POST /api/validate_chirp", validateChirp)

	serverMux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	serverMux.HandleFunc("POST /admin/reset", apiCfg.handlerResetServerHits)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: serverMux,
	}
	log.Printf("Servign on port %s\n", port)
	log.Fatal(server.ListenAndServe())
}
