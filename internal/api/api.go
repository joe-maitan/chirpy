package api

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/joe-maitan/chirpy/internal/auth"
	"github.com/joe-maitan/chirpy/internal/database"
)



type User struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	// Token        string    `json:"token"`
	// RefreshToken string    `json:"refresh_token"`
	IsChirpyRed    bool		`json:"is_chirpy_red"`
} // End User struct

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
} // End Chirp struct













func (cfg *Config) HandleUpdateUser(w http.ResponseWriter, r *http.Request) {
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
		RespondWithError(w, http.StatusUnauthorized, "Missing authentication token", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.JWTSecret)
	if err != nil {
		log.Printf("api.go - HandleUpdateUser() - Error trying to call auth.ValidateJWT: %v", err)
		RespondWithError(w, http.StatusUnauthorized, "Could not validate JWT", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		log.Printf("api.go - HandleUpdateUser() - Error decoding parameters: %v", err)
		RespondWithError(w, http.StatusUnauthorized, "Error decoding parameters", err)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		log.Printf("api.go - HandleUpdateUser() - Error hashing password: %v", err)
		RespondWithError(w, http.StatusInternalServerError, "Error hashing password", err)
		return
	}

	user, err := cfg.DB.UpdateUserPassword(r.Context(), database.UpdateUserPasswordParams{
		ID:             userID,
		Email:          params.Email,
		HashedPassword: hashedPassword,
	})

	if err != nil {
		log.Printf("api.go - HandleUpdateUser() - Error updating user: %v", err)
		RespondWithError(w, http.StatusInternalServerError, "Error updating user", err)
		return
	}

	RespondWithJSON(w, http.StatusOK, Response{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
		IsChirpyRed: user.IsChirpyRed,
	})
} // End HandleUpdateUser() func

func (cfg *Config) HandleCreateChirp(w http.ResponseWriter, r *http.Request) {
	badWords := []string{
		"kerfuffle",
		"sharbert",
		"fornax",
	}

	type parameters struct {
		Body string `json:"body"`
		// UserID uuid.UUID `json:"user_id"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Missing JWT", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.JWTSecret)
	if err != nil {
		log.Printf("api.go - HandleCreateChirp() - Error trying to call auth.ValidateJWT: %v", err)
		RespondWithError(w, http.StatusUnauthorized, "Could not validate JWT", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		log.Printf("api.go - HandleCreateChirp() - Error decoding parameters: %v", err)
		RespondWithError(w, http.StatusInternalServerError, "Error decoding parameters", err)
		return
	}

	// 1. Validate the length of the chirp.
	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		RespondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	// 2. Clean the chirp for any bad words.
	words := strings.Split(params.Body, " ")
	for _, word := range words {
		for _, pottyWord := range badWords {
			if strings.Contains(strings.ToLower(word), strings.ToLower(pottyWord)) {
				params.Body = strings.ReplaceAll(params.Body, word, "****")
			}
		}
	}

	chirp, err := cfg.DB.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   params.Body,
		UserID: userID,
	})

	if err != nil {
		log.Printf("api.go - HandleCreateChirp() - Error creating chirp: %v", err)
		RespondWithError(w, http.StatusInternalServerError, "Error creating chirp", err)
		return
	}

	RespondWithJSON(w, http.StatusCreated, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
} // End HandleCreateChirp() func

func (cfg *Config) HandleGetAllChirps(w http.ResponseWriter, r *http.Request) {
	DBChirps, err := cfg.DB.GetAllChirps(r.Context())
	if err != nil {
		log.Printf("api.go - HandleGetAllChirps() - Error getting all chirps: %v", err)
		RespondWithError(w, http.StatusInternalServerError, "Error getting all chirps", err)
		return
	}

	chirps := []Chirp{}
	for _, c := range DBChirps {
		chirps = append(chirps, Chirp{
			ID:        c.ID,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
			Body:      c.Body,
			UserID:    c.UserID,
		})
	}

	RespondWithJSON(w, http.StatusOK, chirps)
} // End HandleGetAllChirps() func

func (cfg *Config) HandleGetChirp(w http.ResponseWriter, r *http.Request) {
	log.Println(r.PathValue("chirpID"))
	chirpID := r.PathValue("chirpID")
	if chirpID == "" {
		RespondWithError(w, http.StatusNotFound, "Invalid chirp ID", nil)
		return
	}

	data, err := uuid.Parse(chirpID)
	if err != nil {
		RespondWithError(w, http.StatusNotFound, "Error parsing id as a uuid", err)
		return
	}

	chirp, err := cfg.DB.GetChirp(r.Context(), data)
	if err != nil {
		RespondWithError(w, http.StatusNotFound, "Error finding chirp", err)
		return
	}

	RespondWithJSON(w, http.StatusOK, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
} // End HandleGetChirp() func

func (cfg *Config) HandleDeleteChirp(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Missing authentication token", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.JWTSecret)
	if err != nil {
		log.Printf("api.go - HandleUpdateUser() - Error trying to call auth.ValidateJWT: %v", err)
		RespondWithError(w, http.StatusUnauthorized, "Could not validate JWT", err)
		return
	}

	chirpID := r.PathValue("chirpID")
	if chirpID == "" {
		RespondWithError(w, http.StatusNotFound, "Invalid chirp ID", nil)
		return
	}

	data, err := uuid.Parse(chirpID)
	if err != nil {
		RespondWithError(w, http.StatusNotFound, "Error parsing id as a uuid", err)
		return
	}

	chirp, err := cfg.DB.GetChirp(r.Context(), data)
	if err != nil {
		RespondWithError(w, http.StatusNotFound, "Error finding chirp", err)
		return
	}

	if chirp.UserID != userID {
		RespondWithError(w, http.StatusForbidden, "You don't have permission to delete this chirp", nil)
		return
	}

	err = cfg.DB.DeleteChirp(r.Context(), database.DeleteChirpParams{
		ID: chirp.ID,
		UserID: userID,
	})
	if err != nil {
		log.Printf("api.go - HandleDeleteChirp() - Error deleting chirp: %v", err)
		RespondWithError(w, http.StatusInternalServerError, "Error deleting chirp", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
} // End HandleDeleteChirp() func

func (cfg *Config) CheckRefreshToken(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Could not find token", err)
		return
	}

	user, err := cfg.DB.GetUserFromRefreshToken(r.Context(), refreshToken)
	if err != nil {
		log.Printf("api.go - CheckRefreshToken() - Error trying to get user from refresh token: %v", err)
		RespondWithError(w, http.StatusUnauthorized, "Could not get user from refresh token", err)
		return
	}

	accessToken, err := auth.MakeJWT(
		user.ID,
		cfg.JWTSecret,
		time.Hour,
	)

	if err != nil {
		log.Printf("api.go - CheckRefreshToken() - Error trying to create access JWT: %v", err)
		RespondWithError(w, http.StatusInternalServerError, "Couldn't validate token", err)
		return
	}

	RespondWithJSON(w, http.StatusOK, response{
		Token: accessToken,
	})
} // End CheckRefreshToken() func

func (cfg *Config) RevokeRefreshToken(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("api.go - RevokeRefreshToken() - Error trying to get refresh token from header: %v", err)
		RespondWithError(w, http.StatusBadRequest, "Couldn't find token", err)
		return
	}

	_, err = cfg.DB.RevokeRefreshToken(r.Context(), refreshToken)
	if err != nil {
		log.Printf("api.go - RevokeRefreshToken() - Error trying to revoke refresh token: %v", err)
		RespondWithError(w, http.StatusInternalServerError, "Couldn't revoke session", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
} // End RevokeRefreshToken() func

func (cfg *Config) HandlePolkaWebhook(w http.ResponseWriter, r *http.Request) {
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
		RespondWithError(w, 401, "Missing API key", err)
		return
	}

	if fetchedKey != cfg.PolkaAPIKey {
		log.Printf("api.go - HandlePolkaWebhook() - Invalid API key: %s", fetchedKey)
		RespondWithError(w, 401, "Invalid API key", nil)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		log.Printf("api.go - HandleUpdateUser() - Error decoding parameters: %v", err)
		RespondWithError(w, http.StatusUnauthorized, "Error decoding parameters", err)
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
		RespondWithError(w, 400, "Invalid user ID format", err)
		return
	}

	fetchedUser, err := cfg.DB.GetUserByID(r.Context(), userID)
	if err != nil {
		log.Printf("api.go - HandlePolkaWebhook() - Error fetching user by ID: %v", err)
		RespondWithError(w, 404, "Error fetching user", err)
		return
	}

	if fetchedUser.ID != userID {
		log.Printf("api.go - HandlePolkaWebhook() - User ID mismatch: expected %s, got %s", userID, fetchedUser.ID)
		RespondWithError(w, 500, "User ID mismatch", nil)
		return
	}

	_, err = cfg.DB.UpgradeUserToChirpyRed(r.Context(), userID)
	if err != nil {
		log.Printf("api.go - HandlePolkaWebhook() - Error upgrading user to Chirpy Red: %v", err)
		RespondWithError(w, 500, "Error upgrading user to Chirpy Red", err)
		return
	}

	w.WriteHeader(204)
	w.Write([]byte("User upgraded to Chirpy Red successfully"))
} // End HandlePolkaWebhook() func
