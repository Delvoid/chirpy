package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/Delvoid/chirpy/database"
)

type apiConfig struct {
	fileserverHits int
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) metricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	html := `
    <html>
    <body>
        <h1>Welcome, Chirpy Admin</h1>
        <p>Chirpy has been visited %d times!</p>
    </body>
    </html>
`
	fmt.Fprintf(w, html, cfg.fileserverHits)
}

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits = 0
	w.WriteHeader(http.StatusOK)
}

func healthzHandlert(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func main() {
	debug := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()

	const port = "8080"
	cfg := &apiConfig{}

	if *debug {
		log.Println("Debug mode enabled")
		err := database.RemoveDatabase()
		if err != nil {
			log.Fatalf("Failed to remove database: %v", err)
		}
	}

	err := database.Init()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("."))
	appHandler := cfg.middlewareMetricsInc(http.StripPrefix("/app/", fileServer))
	mux.Handle("/app/", appHandler)
	mux.HandleFunc("GET /admin/metrics", cfg.metricsHandler)

	mux.HandleFunc("GET /api/healthz", healthzHandlert)
	mux.HandleFunc("GET /api/reset", cfg.resetHandler)

	mux.HandleFunc("POST /api/chirps", createChirpHandler)
	mux.HandleFunc("GET /api/chirps", getChirpsHandler)
	mux.HandleFunc("GET /api/chirps/{chirpID}", getChirpByIDHandler)

	mux.HandleFunc("POST /api/users", createUserHandler)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Starting server on port: %s\n", port)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
