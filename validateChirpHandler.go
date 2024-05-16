package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type chirpRequest struct {
	Body string `json:"body"`
}

type chirpResponse struct {
	Error       string `json:"error,omitempty"`
	CleanedBody string `json:"cleaned_body,omitempty"`
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

	cleanedBody := replaceProfaneWords(resBody.Body)
	respondWithJSON(w, chirpResponse{CleanedBody: cleanedBody}, http.StatusOK)

}

func respondWithError(w http.ResponseWriter, errorMessage string, statusCode int) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(chirpResponse{Error: errorMessage})
}

func respondWithJSON(w http.ResponseWriter, payload interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(statusCode)
	w.Write(dat)
}
