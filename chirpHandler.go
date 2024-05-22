package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/Delvoid/chirpy/database"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type chirpRequest struct {
	Body string `json:"body"`
}

type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}

func createChirpHandler(w http.ResponseWriter, r *http.Request) {
	resBody := chirpRequest{}

	if err := json.NewDecoder(r.Body).Decode(&resBody); err != nil {
		respondWithError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	chirp, err := database.CreateChirp(resBody.Body)
	if err != nil {
		respondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	respondWithJSON(w, chirp, http.StatusCreated)

}

func getChirpsHandler(w http.ResponseWriter, r *http.Request) {
	chirps, err := database.GetChirps()
	if err != nil {
		respondWithError(w, "Failed to retrieve chirps", http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, chirps, http.StatusOK)
}

func getChirpByIDHandler(w http.ResponseWriter, r *http.Request) {
	chirpID, err := strconv.Atoi(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, "Invalid chirp ID", http.StatusBadRequest)
		return
	}
	chirp, err := database.GetChirpByID(chirpID)
	if err != nil {
		if errors.Is(err, database.ErrChirpNotFound) {
			respondWithError(w, "Chirp not found", http.StatusNotFound)
		} else {
			respondWithError(w, "Failed to retrieve chirp", http.StatusInternalServerError)
		}
		return
	}

	respondWithJSON(w, chirp, http.StatusOK)
}

func respondWithError(w http.ResponseWriter, errorMessage string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{Error: errorMessage})
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
