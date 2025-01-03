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

	w.Header().Set("Content-Type", "application/json") // set content type to JSON

	decoder := json.NewDecoder(req.Body)
	r := reqS{}
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
		respondWithJSON(w, 200, cleanedBody{
			Cleaned_Body: cleanText(r.Body),
		})
		return
	}

}

func (cfg *apiConfig) createUser(w http.ResponseWriter, req *http.Request) {
	e := emailReq{}
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&e); err != nil {
		log.Printf("error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}

	email, err := cfg.db.CreateUser(req.Context(), e.Email)
	if err != nil {
		log.Printf("error creating user: %s", err)
		w.WriteHeader(500)
		return
	}

	customUser := User{
		ID:        email.ID,
		CreatedAt: email.CreatedAt,
		UpdatedAt: email.UpdatedAt,
		Email:     email.Email,
	}

	respondWithJSON(w, 201, customUser)

}
