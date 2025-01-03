package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync/atomic"
)

func (cfg *apiConfig) reqCount(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	body := fmt.Sprintf(`<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %v times!</p>
  </body>
</html>`, strconv.Itoa(int(cfg.fileserverHits.Load())))

	w.Write([]byte(body))
}

func (cfg *apiConfig) resetCount(w http.ResponseWriter, req *http.Request) {

	type response struct {
		Body string `json:"body"`
	}

	if cfg.platform != "dev" {
		respondWithError(w, http.StatusForbidden, "")
		return
	}

	err := cfg.db.DeleteUsers(req.Context())
	if err != nil {
		log.Printf("error deleting users: %s", err)
	}

	cfg.fileserverHits = atomic.Int32{}
	respondWithJSON(w, http.StatusOK, response{
		Body: "Hit count reset to 0\nDatabase re-initialized\n",
	})
}
