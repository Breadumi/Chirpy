package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func readinessEndpoint(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func validateChirp(w http.ResponseWriter, req *http.Request) {

	type params struct {
		Body string `json:"body"`
	}

	type response struct {
		CleanedBody string `json:"cleaned_body"`
	}

	w.Header().Set("Content-Type", "application/json") // set content type to JSON

	decoder := json.NewDecoder(req.Body)
	r := params{}
	err := decoder.Decode(&r)

	// set http status codes and body
	if err != nil {
		log.Printf("Error decoding paramters: %s", err)
		w.WriteHeader(500)
		return
	}

	// check if length is acceptable
	if len(r.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	} else {
		respondWithJSON(w, 200, response{
			CleanedBody: cleanText(r.Body),
		})
		return
	}

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
		w.WriteHeader(500)
		return
	}

	user, err := cfg.db.CreateUser(req.Context(), params.Email)
	if err != nil {
		log.Printf("error creating user: %s", err)
		w.WriteHeader(500)
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
