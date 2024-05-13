package main

import (
	"log"
	"net/http"
)

func main() {
	const port = "8080"

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./"))
	mux.Handle("/", fileServer)

	assetsDir := http.Dir("./assets")
	assetsHandler := http.StripPrefix("/assets/", http.FileServer(assetsDir))
	mux.Handle("/assets/", assetsHandler)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Starting server on port: %s\n", port)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
