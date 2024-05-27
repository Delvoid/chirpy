package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/Delvoid/chirpy/database"
	"github.com/dgrijalva/jwt-go"
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
		ID          int    `json:"id"`
		Email       string `json:"email"`
		IsChirpyRed bool   `json:"is_chirpy_red"`
	}{
		ID:          user.ID,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	}, http.StatusCreated)

}

func updateUserHandler(jwtSecret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			respondWithError(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			respondWithError(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		claims := &jwt.StandardClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})

		if err != nil {
			respondWithError(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			respondWithError(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		userID, err := strconv.Atoi(claims.Subject)
		if err != nil {
			respondWithError(w, "Invalid user ID", http.StatusUnauthorized)
			return
		}

		var req userRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondWithError(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		user, err := database.UpdateUser(userID, req.Email, req.Password)
		if err != nil {
			respondWithError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		respondWithJSON(w, struct {
			ID    int    `json:"id"`
			Email string `json:"email"`
		}{
			ID:    user.ID,
			Email: user.Email,
		}, http.StatusOK)
	}
}
