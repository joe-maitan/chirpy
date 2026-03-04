package main

import (
	"net/http"
	"encoding/json"

	"github.com/joe-maitan/chirpy/internal/auth"
	"github.com/joe-maitan/chirpy/internal/database"
)

func (cfg *apiConfig) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		logStatement("handler_create_user.go", "HandleCreateUser()", "Error decoding parameters", err)
		respondWithError(w, http.StatusInternalServerError, "Error decoding parameters", err)
		return
	}

	if params.Email == "" || params.Password == "" {
		respondWithError(w, http.StatusBadRequest, "Email and password are required", nil)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		logStatement("handler_create_user.go", "HandleCreateUser()", "Error hashing password", err)
		respondWithError(w, http.StatusInternalServerError, "Error hashing password", err)
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPassword,
	})

	if err != nil {
		logStatement("handler_create_user.go", "HandleCreateUser()", "Error creating user", err)
		respondWithError(w, http.StatusInternalServerError, "Error creating user", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, User{
		ID: user.ID,
		Email: user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		IsChirpyRed: user.IsChirpyRed,
	})
} // End HandleCreateUser() func