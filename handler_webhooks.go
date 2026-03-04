package main

import (
	"log"
	"net/http"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/joe-maitan/chirpy/internal/auth"
)

func (cfg *apiConfig) HandlePolkaWebhook(w http.ResponseWriter, r *http.Request) {
	type EventData struct {
		UserID string `json:"user_id"`
	}

	type parameters struct {
		Event string `json:"event"`
		Data  EventData `json:"data"`
	}

	fetchedKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		log.Printf("api.go - HandlePolkaWebhook() - Error trying to get API key from header: %v", err)
		respondWithError(w, 401, "Missing API key", err)
		return
	}

	if fetchedKey != cfg.polkaAPIKey {
		log.Printf("api.go - HandlePolkaWebhook() - Invalid API key: %s", fetchedKey)
		respondWithError(w, 401, "Invalid API key", nil)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		log.Printf("api.go - HandleUpdateUser() - Error decoding parameters: %v", err)
		respondWithError(w, http.StatusUnauthorized, "Error decoding parameters", err)
		return
	}

	log.Printf("Received Polka webhook - Event: %s, Data: %s", params.Event, params.Data)

	if params.Event != "user.upgraded" {
		w.WriteHeader(204)
		w.Write([]byte("Webhook received but no action taken"))
	}

	userID, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		log.Printf("api.go - HandlePolkaWebhook() - Error parsing user ID: %v", err)
		respondWithError(w, 400, "Invalid user ID format", err)
		return
	}

	fetchedUser, err := cfg.db.GetUserByID(r.Context(), userID)
	if err != nil {
		log.Printf("api.go - HandlePolkaWebhook() - Error fetching user by ID: %v", err)
		respondWithError(w, 404, "Error fetching user", err)
		return
	}

	if fetchedUser.ID != userID {
		log.Printf("api.go - HandlePolkaWebhook() - User ID mismatch: expected %s, got %s", userID, fetchedUser.ID)
		respondWithError(w, 500, "User ID mismatch", nil)
		return
	}

	_, err = cfg.db.UpgradeUserToChirpyRed(r.Context(), userID)
	if err != nil {
		log.Printf("api.go - HandlePolkaWebhook() - Error upgrading user to Chirpy Red: %v", err)
		respondWithError(w, 500, "Error upgrading user to Chirpy Red", err)
		return
	}

	w.WriteHeader(204)
	w.Write([]byte("User upgraded to Chirpy Red successfully"))
} // End HandlePolkaWebhook() func