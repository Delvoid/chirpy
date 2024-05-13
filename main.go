package main

import (
	"log"
	"net/http"
)

func healthzHandlert(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func main() {
	const port = "8080"

	mux := http.NewServeMux()

	mux.Handle("/app/*", http.StripPrefix("/app", http.FileServer(http.Dir("."))))

	assetsDir := http.Dir("./assets")
	assetsHandler := http.StripPrefix("/app/assets/", http.FileServer(assetsDir))
	mux.Handle("/app/assets/", assetsHandler)

	mux.HandleFunc("/healthz", healthzHandlert)

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
