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

const port = "8080"
const filepathRoot = "."

func errorHandle(introString string, err error) {
	if err != nil {
		fmt.Printf("%s: %v\n", introString, err)
	}
}

func main() {
	server := run()
	log.Fatal(server.ListenAndServe())
}

func run() *http.Server {
	godotenv.Load()

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

	mux := apiCfg.routes(filepathRoot)
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	log.Printf("Serving on port %s\n", port)

	return server
}
