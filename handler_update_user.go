package main

import (
	"time"
	"net/http"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/joe-maitan/chirpy/internal/auth"
	"github.com/joe-maitan/chirpy/internal/database"
)

func (cfg *apiConfig) HandleUpdateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type Response struct {
		ID	uuid.UUID `json:"id"`
		Email string `json:"email"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		IsChirpyRed bool `json:"is_chirpy_red"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Missing authentication token", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Could not validate JWT", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Error decoding parameters", err)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error hashing password", err)
		return
	}

	user, err := cfg.db.UpdateUserPassword(r.Context(), database.UpdateUserPasswordParams{
		ID:             userID,
		Email:          params.Email,
		HashedPassword: hashedPassword,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error updating user", err)
		return
	}

	respondWithJSON(w, http.StatusOK, Response{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
		IsChirpyRed: user.IsChirpyRed,
	})
} // End HandleUpdateUser() func