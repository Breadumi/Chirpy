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
