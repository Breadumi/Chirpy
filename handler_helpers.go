package main

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strings"
)

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

	type response struct {
		Error string `json:"error"`
	}

	w.WriteHeader(code)

	dat, err := json.Marshal(response{
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	dat, err := json.Marshal(payload)

	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}

	w.Write(dat)

}
