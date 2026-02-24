package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/joe-maitan/chirpy/internal/auth"
	"github.com/joe-maitan/chirpy/internal/database"
)

type Config struct {
	FileServerHits atomic.Int32
	DB             *database.Queries // DB driver
	Platform       string            // platform the api is running on (e.g. "DEV", "PROD")
	JWTSecret      string            // secret key for signing JWT tokens
} // End api config struct

type User struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
} // End User struct

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
} // End Chirp struct

func RespondWithError(w http.ResponseWriter, code int, msg string, err error) {
	if err != nil {
		log.Println(err)
	}

	if code > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}

	type errorResponse struct {
		Error string `json:"error"`
	}

	RespondWithJSON(w, code, errorResponse{
		Error: msg,
	})
} // End RespondWithError() func

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
	}

	w.WriteHeader(code)
	w.Write(data)
} // End RespondWithJSON() func

func (cfg *Config) HandlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte(fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", cfg.FileServerHits.Load())))
} // End handlerMetrics() func

func (cfg *Config) HandlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.Platform != "dev" {
		RespondWithError(w, 403, "Forbidden", nil)
		return
	}

	cfg.DB.DeleteUsers(r.Context())
	cfg.FileServerHits.Store(0)
	w.WriteHeader(200)
	w.Write([]byte("Hits reset to 0"))
} // End HandlerReset() func

func (cfg *Config) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.FileServerHits.Add(1)
		next.ServeHTTP(w, r)
	})
} // End MiddlewareMetricsInc() func

func (cfg *Config) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Error decoding parameters", err)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Error hashing password", err)
		return
	}

	user, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPassword,
	})

	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Error creating user", err)
		return
	}

	RespondWithJSON(w, http.StatusCreated, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	})
} // End HandleCreateUser() func

func (cfg *Config) HandleUserLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email     string         `json:"email"`
		Password  string         `json:"password"`
		ExpiresIn *time.Duration `json:"expires_in_seconds"` // Optional parameter
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Error decoding parameters", err)
		return
	}

	user, err := cfg.DB.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	passwordMatch, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	if passwordMatch == false {
		RespondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	expiresIn := 1 * time.Hour
	if params.ExpiresIn != nil && *params.ExpiresIn < time.Hour {
		expiresIn = *params.ExpiresIn
	}

	token, err := auth.MakeJWT(user.ID, cfg.JWTSecret, expiresIn)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Error creating token", err)
		return
	}

	refresh_token, err := auth.MakeRefreshToken()
	if err != nil {
		RespondWithError(w, 500, "Error making refresh token", err)
		return
	}

	cfg.DB.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refresh_token,
		UserID:    user.ID,
		ExpiresAt: time.Now().AddDate(0, 0, 60),
	})

	RespondWithJSON(w, http.StatusOK, User{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		Token:        token,
		RefreshToken: refresh_token,
	})
} // End HandleUserLogin() func

func (cfg *Config) HandleCreateChirp(w http.ResponseWriter, r *http.Request) {
	badWords := []string{
		"kerfuffle",
		"sharbert",
		"fornax",
	}

	type parameters struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Error decoding parameters", err)
		return
	}

	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Missing or invalid token", err)
		return
	}

	userID, err := auth.ValidateJWT(tokenString, cfg.JWTSecret)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "JWT Token is not valid", err)
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
		RespondWithError(w, http.StatusNotFound, "Invalid request", nil)
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

func (cfg *Config) CheckRefreshToken(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		RespondWithError(w, 401, "authorization header missing", nil)
		return
	}

	const prefix = "Bearer "
	if !strings.HasPrefix(authHeader, prefix) {
		RespondWithError(w, 401, "authorization header is not a bearer token", nil)
	}

	token := strings.TrimPrefix(authHeader, prefix)
	if token == "" {
		RespondWithError(w, 401, "bearer token is empty", nil)
		return
	}

	fetchedRefreshToken, err := cfg.DB.GetRefreshToken(r.Context(), token)
	if err != nil {
		RespondWithError(w, 401, "token did not exist", nil)
	}

	// if the token expires now revoke it.
	if fetchedRefreshToken.ExpiresAt.Equal(time.Now()) {
		cfg.RevokeRefreshToken(w, r)
	}

	if fetchedRefreshToken.RevokedAt.Valid {
		RespondWithError(w, 401, "refresh token has been revoked.", nil)
	}

} // End CheckRefreshToken() func

func (cfg *Config) RevokeRefreshToken(w http.ResponseWriter, r *http.Request) {

} // End RevokeRefreshToken() func
