package main

import (
	"encoding/json"
	"net/http"

	"github.com/Delvoid/chirpy/database"
)

type userRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	resBody := userRequest{}

	if err := json.NewDecoder(r.Body).Decode(&resBody); err != nil {
		respondWithError(w, "Invalid request body", http.StatusBadRequest)
		defer r.Body.Close()
		return
	}
	defer r.Body.Close()

	user, err := database.CreateUser(resBody.Email, resBody.Password)
	if err != nil {
		respondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	respondWithJSON(w, struct {
		ID    int    `json:"id"`
		Email string `json:"email"`
	}{
		ID:    user.ID,
		Email: user.Email,
	}, http.StatusCreated)

}
