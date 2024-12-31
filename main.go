package main

import (
	"log"
	"net/http"
)

const (
	filepathRoot = "/"
)

func main() {

	port := "8080"

	mux := http.NewServeMux()
	server := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	mux.Handle("/", http.FileServer(http.Dir(".")))
	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())

}
