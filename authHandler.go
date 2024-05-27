package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Delvoid/chirpy/database"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type loginRequest struct {
	Email            string `json:"email"`
	Password         string `json:"password"`
	ExpiresInSeconds int    `json:"expires_in_seconds,omitempty"`
}

func validateToken(r *http.Request, jwtSecret string) (int, error) {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		return 0, errors.New("Missing Authorization header")
	}

	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	claims := &jwt.StandardClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return 0, errors.New("Invalid token")
	}

	if !token.Valid {
		return 0, errors.New("Invalid token")
	}

	userID, err := strconv.Atoi(claims.Subject)
	if err != nil {
		return 0, errors.New("Invalid user ID")
	}

	return userID, nil
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

		refreshToken, err := database.CreateRefreshToken(user.ID, 60*24*60*60*time.Second) // Expire in 60 days
		if err != nil {
			respondWithError(w, "Failed to generate refresh token", http.StatusInternalServerError)
			return
		}

		respondWithJSON(w, struct {
			ID           int    `json:"id"`
			Email        string `json:"email"`
			Token        string `json:"token"`
			RefreshToken string `json:"refresh_token"`
			IsChirpyRed  bool   `json:"is_chirpy_red"`
		}{
			ID:           user.ID,
			Email:        user.Email,
			Token:        tokenString,
			RefreshToken: refreshToken,
			IsChirpyRed:  user.IsChirpyRed,
		}, http.StatusOK)
	}

}

func refreshHandler(jwtSecret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			respondWithError(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		refreshToken, err := database.GetRefreshToken(tokenString)
		if err != nil {
			respondWithError(w, "Invalid refresh token", http.StatusUnauthorized)
			return
		}

		if refreshToken.ExpiresAt.Before(time.Now()) {
			respondWithError(w, "Refresh token has expired", http.StatusUnauthorized)
			return
		}

		claims := &jwt.StandardClaims{
			Issuer:    "chirpy",
			IssuedAt:  jwt.TimeFunc().Unix(),
			ExpiresAt: jwt.TimeFunc().Unix() + 3600, // Expire in 1 hour
			Subject:   fmt.Sprintf("%d", refreshToken.UserID),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err = token.SignedString([]byte(jwtSecret))
		if err != nil {
			respondWithError(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		respondWithJSON(w, struct {
			Token string `json:"token"`
		}{
			Token: tokenString,
		}, http.StatusOK)
	}
}

func revokeHandler(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		respondWithError(w, "Missing Authorization header", http.StatusUnauthorized)
		return
	}

	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	err := database.DeleteRefreshToken(tokenString)
	if err != nil {
		respondWithError(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
