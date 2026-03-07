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
		respondWithError(w, http.StatusUnauthorized, "Missing API key", err)
		return
	}

	if fetchedKey != cfg.polkaAPIKey {
		respondWithError(w, http.StatusUnauthorized, "Invalid API key", nil)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Error decoding parameters", err)
		return
	}

	log.Printf("Received Polka webhook - Event: %s, Data: %s", params.Event, params.Data)

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		w.Write([]byte("Webhook received but no action taken"))
	}

	userID, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID format", err)
		return
	}

	fetchedUser, err := cfg.db.GetUserByID(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Error fetching user", err)
		return
	}

	if fetchedUser.ID != userID {
		respondWithError(w, http.StatusInternalServerError, "User ID mismatch", nil)
		return
	}

	_, err = cfg.db.UpgradeUserToChirpyRed(r.Context(), userID)
	if err != nil {
		log.Printf("api.go - HandlePolkaWebhook() - Error upgrading user to Chirpy Red: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Error upgrading user to Chirpy Red", err)
		return
	}

	w.WriteHeader(204)
	w.Write([]byte("User upgraded to Chirpy Red successfully"))
} // End HandlePolkaWebhook() func