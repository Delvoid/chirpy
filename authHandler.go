package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Delvoid/chirpy/database"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type loginRequest struct {
	Email            string `json:"email"`
	Password         string `json:"password"`
	ExpiresInSeconds int    `json:"expires_in_seconds,omitempty"`
}

func loginHandler(jwtSecret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var req loginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondWithError(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		user, err := database.GetUserByEmail(req.Email)
		if err != nil {
			respondWithError(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
		if err != nil {
			respondWithError(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		expiresInSeconds := int64(24 * 60 * 60) // Default to 24 hours
		if req.ExpiresInSeconds > 0 {
			if req.ExpiresInSeconds > 24*60*60 {
				expiresInSeconds = int64(24 * 60 * 60) // Limit to 24 hours
			} else {
				expiresInSeconds = int64(req.ExpiresInSeconds)
			}
		}

		println("req expiresInSeconds: ", req.ExpiresInSeconds)
		println("expiresInSeconds: ", expiresInSeconds)

		claims := &jwt.StandardClaims{
			Issuer:    "chirpy",
			IssuedAt:  jwt.TimeFunc().Unix(),
			ExpiresAt: jwt.TimeFunc().Unix() + expiresInSeconds,
			Subject:   fmt.Sprintf("%d", user.ID),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(jwtSecret))
		if err != nil {
			respondWithError(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		respondWithJSON(w, struct {
			ID    int    `json:"id"`
			Email string `json:"email"`
			Token string `json:"token"`
		}{
			ID:    user.ID,
			Email: user.Email,
			Token: tokenString,
		}, http.StatusOK)
	}

}
