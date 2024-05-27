package main

import (
	"encoding/json"
	"net/http"

	"github.com/Delvoid/chirpy/database"
)

type PolkaWebhookEvent struct {
	Event string `json:"event"`
	Data  struct {
		UserID int `json:"user_id"`
	} `json:"data"`
}

func polkaWebhookHandler(w http.ResponseWriter, r *http.Request) {

	var req PolkaWebhookEvent
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Event != "user.upgraded" {
		respondWithError(w, "Invalid event type", http.StatusNoContent)
		return
	}

	err := database.UpgradeUserToChirpyRed(req.Data.UserID)
	if err != nil {
		if err == database.ErrUserNotFound {
			respondWithError(w, "User not found", http.StatusNotFound)
		} else {
			respondWithError(w, "Failed to upgrade user", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
