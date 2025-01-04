package main

import (
	"bytes"
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/Breadumi/Chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

const (
	filepathRoot = "/"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
	secret         string
	polkaKey       string
}

func main() {

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	if err = upMigrateDB(db); err != nil {
		log.Fatal(err)
	}

	dbQueries := database.New(db)

	cfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
		platform:       os.Getenv("PLATFORM"),
		secret:         os.Getenv("SECRET"),
		polkaKey:       os.Getenv("POLKA_KEY"),
	}

	port := "8080"

	mux := http.NewServeMux()
	server := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	mux.HandleFunc("GET /api/healthz", readinessEndpoint)             // check readiness
	mux.HandleFunc("POST /api/users", cfg.createUser)                 // create one user
	mux.HandleFunc("POST /api/chirps", cfg.createChirp)               // create one chirp
	mux.HandleFunc("GET /api/chirps", cfg.getChirps)                  // return all chirps
	mux.HandleFunc("GET /api/chirps/{chirpID}", cfg.getChirp)         // return one chirp
	mux.HandleFunc("POST /api/login", cfg.login)                      // log in a user
	mux.HandleFunc("POST /api/refresh", cfg.refresh)                  // create new JWT if refresh valid
	mux.HandleFunc("POST /api/revoke", cfg.revoke)                    // revoke expired refresh tokens
	mux.HandleFunc("PUT /api/users", cfg.updateLogin)                 // update email and password
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", cfg.deleteChirp)   // delete chirp with chirpID
	mux.HandleFunc("POST /api/polka/webhooks", cfg.registerChirpyRed) // subscribe user to Chirpy Red

	fileServer := http.FileServer(http.Dir("."))
	mux.Handle("/app/", http.StripPrefix("/app", cfg.mwMetricInc(fileServer))) // serve files from root directory

	mux.HandleFunc("GET /admin/metrics", cfg.reqCount)  // return hit counter
	mux.HandleFunc("POST /admin/reset", cfg.resetCount) // reset hit counter and re-init database

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())

}

func upMigrateDB(db *sql.DB) error {
	old := os.Stdout
	r, w, _ := os.Pipe()
	var buf bytes.Buffer
	log.SetOutput(&buf)

	err := goose.Up(db, "sql/schema") // successful output is not printed to console

	// close new output containing unwanted logs
	w.Close()
	r.Close()
	os.Stdout = old
	log.SetOutput(os.Stderr)

	// log migration errors in console if necessary
	if err != nil {
		return err
	}

	return nil
}
