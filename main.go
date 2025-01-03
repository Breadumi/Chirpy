package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/Breadumi/Chirpy/internal/database"
	"github.com/joho/godotenv"
)

const (
	filepathRoot = "/"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
}

func main() {

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	dbQueries := database.New(db)

	cfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
	}

	port := "8080"

	mux := http.NewServeMux()
	server := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	mux.HandleFunc("GET /api/healthz", readinessEndpoint)     // register readiness endpoint
	mux.HandleFunc("POST /api/validate_chirp", validateChirp) // register chirp validation endpoint
	mux.HandleFunc("POST /api/users", createUser)             // register user creation endpoint

	fileServer := http.FileServer(http.Dir("."))
	mux.Handle("/app/", http.StripPrefix("/app", cfg.mwMetricInc(fileServer))) // serve files from root directory

	mux.HandleFunc("GET /admin/metrics", cfg.reqCount)  // register hit counter endpoint
	mux.HandleFunc("POST /admin/reset", cfg.resetCount) // register hit counter reset

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())

}
