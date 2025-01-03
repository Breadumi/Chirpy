package main

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strings"
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

func cleanText(msg string) string {

	wordList := []string{"kerfuffle", "sharbert", "fornax"}
	replWord := "****"
	words := strings.Split(msg, " ")

	for i, word := range words {
		if slices.Contains(wordList, strings.ToLower(word)) {
			words[i] = replWord
		}
	}

	return strings.Join(words, " ")

}

func respondWithError(w http.ResponseWriter, code int, msg string) {

	w.WriteHeader(code)

	dat, err := json.Marshal(eS{
		Error: msg,
	})

	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}

	w.Write(dat)

}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.WriteHeader(code)

	dat, err := json.Marshal(payload)

	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}

	w.Write(dat)

}
