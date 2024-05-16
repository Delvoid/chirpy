package main

import (
	"encoding/json"
	"net/http"
)

type chirpRequest struct {
	Body string `json:"body"`
}

type chirpResponse struct {
	Error string `json:"error,omitempty"`
	Valid bool   `json:"valid,omitempty"`
}

func validateChirpHandler(w http.ResponseWriter, r *http.Request) {
	resBody := chirpRequest{}

	if err := json.NewDecoder(r.Body).Decode(&resBody); err != nil {
		respondWithError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// check length of body
	if len(resBody.Body) > 140 {
		respondWithError(w, "Chirp is too long", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(chirpResponse{Valid: true})

}

func respondWithError(w http.ResponseWriter, errorMessage string, statusCode int) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(chirpResponse{Error: errorMessage})
}
