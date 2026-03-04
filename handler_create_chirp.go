package main

import (
	"log"
	"errors"
	"strings"
	"net/http"
	"encoding/json"

	// "github.com/google/uuid"
	"github.com/joe-maitan/chirpy/internal/auth"
	"github.com/joe-maitan/chirpy/internal/database"
)

func validateAndCleanChirp(chirpBody string) (string, error) {
	badWords := []string{
		"kerfuffle",
		"sharbert",
		"fornax",
	}

	// 1. Validate the length of the chirp.
	const maxChirpLength = 140
	if len(chirpBody) > maxChirpLength {
		// respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return "", errors.New("Chirp is too long")
	}

	// 2. Clean the chirp for any bad words.
	words := strings.Split(chirpBody, " ")
	for _, word := range words {
		for _, pottyWord := range badWords {
			if strings.Contains(strings.ToLower(word), strings.ToLower(pottyWord)) {
				chirpBody = strings.ReplaceAll(chirpBody, word, "****")
			}
		}
	}

	return chirpBody, nil
} // End validateAndCleanChirp()

func (cfg *apiConfig) HandleCreateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
		// UserID uuid.UUID `json:"user_id"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Missing JWT", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		logStatement("handler_create_chirp.go","HandleCreateChirp()", "Error trying to call auth.ValidateJWT", err)
		respondWithError(w, http.StatusUnauthorized, "Could not validate JWT", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		logStatement("handler_create_chirp.go","HandleCreateChirp()", "Error decoding parameters", err)
		respondWithError(w, http.StatusInternalServerError, "Error decoding parameters", err)
		return
	}

	chirpBody, err := validateAndCleanChirp(params.Body)
	if err != nil {
		log.Printf("handler_create_chirp.go")
		respondWithError(w, http.StatusBadRequest, "", err)
	}

	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   chirpBody,
		UserID: userID,
	})

	if err != nil {
		log.Printf("api.go - HandleCreateChirp() - Error creating chirp: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Error creating chirp", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, database.Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
} // End HandleCreateChirp() func