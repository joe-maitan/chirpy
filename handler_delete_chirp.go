package main

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/joe-maitan/chirpy/internal/auth"
	"github.com/joe-maitan/chirpy/internal/database"
)

func (cfg *apiConfig) HandleDeleteChirp(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Missing authentication token", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		log.Printf("api.go - HandleUpdateUser() - Error trying to call auth.ValidateJWT: %v", err)
		respondWithError(w, http.StatusUnauthorized, "Could not validate JWT", err)
		return
	}

	chirpID := r.PathValue("chirpID")
	if chirpID == "" {
		respondWithError(w, http.StatusNotFound, "Invalid chirp ID", nil)
		return
	}

	data, err := uuid.Parse(chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Error parsing id as a uuid", err)
		return
	}

	chirp, err := cfg.db.GetChirp(r.Context(), data)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Error finding chirp", err)
		return
	}

	if chirp.UserID != userID {
		respondWithError(w, http.StatusForbidden, "You don't have permission to delete this chirp", nil)
		return
	}

	err = cfg.db.DeleteChirp(r.Context(), database.DeleteChirpParams{
		ID: chirp.ID,
		UserID: userID,
	})
	if err != nil {
		log.Printf("api.go - HandleDeleteChirp() - Error deleting chirp: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Error deleting chirp", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
} // End HandleDeleteChirp() func