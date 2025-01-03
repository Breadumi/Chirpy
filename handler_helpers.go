package main

import (
	"encoding/json"
	"fmt"
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

	w.Header().Set("Content-Type", "application/json")

	type response struct {
		Error string `json:"error"`
	}

	w.WriteHeader(code)

	dat, err := json.Marshal(response{
		Error: msg,
	})

	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
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
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(dat)

}

//lint:ignore U1000 debugging function
func prettyprint(payload interface{}) {
	p, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		log.Fatal("prettyprint failed")
	}
	fmt.Printf("%+v\n", string(p))
}
