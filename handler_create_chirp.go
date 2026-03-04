package main

import (
	"strings"
	"net/http"
	"encoding/json"

	"github.com/joe-maitan/chirpy/internal/auth"
	"github.com/joe-maitan/chirpy/internal/database"
)

func validateChirpLength(chirpBody string) (bool) {
	const maxChirpLength = 140
	if len(chirpBody) > maxChirpLength {
		return false
	}

	return true
} // End validateChirpLength() func

func cleanChirp(chirpBody string) (string) {
	badWords := []string{
		"kerfuffle",
		"sharbert",
		"fornax",
	}

	words := strings.Split(chirpBody, " ")
	for _, word := range words {
		for _, pottyWord := range badWords {
			if strings.Contains(strings.ToLower(word), strings.ToLower(pottyWord)) {
				chirpBody = strings.ReplaceAll(chirpBody, word, "****")
			}
		}
	}

	return chirpBody
} // End cleanChirp() func

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

	if !validateChirpLength(params.Body) {
		logStatement("handler_create_chirp.go","HandleCreateChirp()", "Chirp is too long", err)
		respondWithError(w, http.StatusBadRequest, "", err)
	}

	chirpBody := cleanChirp(params.Body)
	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   chirpBody,
		UserID: userID,
	})

	if err != nil {
		logStatement("handler_create_chirp.go","HandleCreateChirp()", "Error creating chirp", err)
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