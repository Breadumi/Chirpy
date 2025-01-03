package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Breadumi/Chirpy/internal/database"
	"github.com/google/uuid"
)

func readinessEndpoint(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (cfg *apiConfig) createChirp(w http.ResponseWriter, req *http.Request) {

	type parameters struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	type response struct {
		Chirp
	}

	params := parameters{}

	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&params); err != nil {
		log.Printf("Error decoding parameters: %s", err)
		respondWithError(w, http.StatusInternalServerError, "error decoding parameters")
		return
	}

	//prettyprint(params) // debugging line

	// check if length is acceptable
	if len(params.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	chirp, err := cfg.db.CreateChirp(req.Context(), database.CreateChirpParams{
		Body:   params.Body,
		UserID: params.UserID,
	})
	//prettyprint(chirp) // debugging line
	if err != nil {
		log.Printf("Error querying database: %s", err)
		respondWithError(w, http.StatusInternalServerError, "error querying database")
		return
	}

	respondWithJSON(w, http.StatusCreated, response{
		Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      cleanText(chirp.Body),
			UserID:    chirp.UserID,
		},
	})

}

func (cfg *apiConfig) getChirps(w http.ResponseWriter, req *http.Request) {
	chirpsSQLC, err := cfg.db.GetChirps(req.Context())
	if err != nil {
		log.Printf("error querying database: %s", err)
		respondWithError(w, http.StatusInternalServerError, "error querying database")
		return
	}

	var chirps []Chirp

	for _, c := range chirpsSQLC {
		chirps = append(chirps, Chirp{
			ID:        c.ID,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
			Body:      c.Body,
			UserID:    c.UserID,
		})
	}

	respondWithJSON(w, http.StatusOK, chirps)

}

func (cfg *apiConfig) getChirp(w http.ResponseWriter, req *http.Request) {
	chirpID := req.PathValue("chirpID")

}

func (cfg *apiConfig) createUser(w http.ResponseWriter, req *http.Request) {

	type parameters struct {
		Email string `json:"email"`
	}

	type response struct {
		User
	}

	params := parameters{}
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&params); err != nil {
		log.Printf("error decoding parameters: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user, err := cfg.db.CreateUser(req.Context(), params.Email)
	//prettyprint(user) // debugging line
	if err != nil {
		log.Printf("error creating user: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, 201, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
	})

}
