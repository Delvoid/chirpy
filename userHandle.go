package main

import (
	"encoding/json"
	"net/http"

	"github.com/Delvoid/chirpy/database"
)

type userRequest struct {
	Email string `json:"email"`
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	resBody := userRequest{}

	if err := json.NewDecoder(r.Body).Decode(&resBody); err != nil {
		respondWithError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	chirp, err := database.CreateUser(resBody.Email)
	if err != nil {
		respondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	respondWithJSON(w, chirp, http.StatusCreated)

}
