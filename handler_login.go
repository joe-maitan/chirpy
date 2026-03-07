package main

import (
	"time"
	"net/http"
	"encoding/json"

	"github.com/joe-maitan/chirpy/internal/auth"
	"github.com/joe-maitan/chirpy/internal/database"
)

func (cfg *apiConfig) HandleUserLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	type response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		logStatement("handler_login.go", "HandleUserLogin()", "Error decoding parameters", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		logStatement("handler_login.go", "HandleUserLogin()", "Error getting user by email", err)
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	match, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil || !match {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	accessToken, err := auth.MakeJWT(
		user.ID,
		cfg.jwtSecret,
		time.Hour,
	)
	if err != nil {
		logStatement("handler_login.go", "HandleUserLogin()", "Error creating access JWT", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't create access JWT", err)
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		logStatement("handler_login.go", "HandleUserLogin()", "Error creating refresh token", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't create refresh token", err)
		return
	}

	_, err = cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().UTC().Add(time.Hour * 24 * 60),
	})

	if err != nil {
		logStatement("handler_login.go", "HandleUserLogin()", "Error saving refresh token in db", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't save refresh token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:        		user.ID,
			CreatedAt: 		user.CreatedAt,
			UpdatedAt: 		user.UpdatedAt,
			Email:     		user.Email,
			IsChirpyRed: 	user.IsChirpyRed,
		},
		Token:        accessToken,
		RefreshToken: refreshToken,
	})
} // End HandleUserLogin() func