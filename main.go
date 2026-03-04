package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joe-maitan/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileServerHits atomic.Int32
	db             *database.Queries // DB driver
	platform       string            // platform the api is running on (e.g. "DEV", "PROD")
	jwtSecret      string            // secret key for signing JWT tokens
	polkaAPIKey	   string            // secret key for validating Polka webhooks
} // End api config struct

func main() {
	godotenv.Load()

	const filepathRoot = "."
	const port = "8080"

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}

	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database: %s", err)
	}

	dbQueries := database.New(dbConn)

	cfg := apiConfig{
		fileServerHits: atomic.Int32{},
		db:             dbQueries,
		platform:       os.Getenv("PLATFORM"),
		jwtSecret:      os.Getenv("SECRET"),
		polkaAPIKey: 	os.Getenv("POLKA_API_KEY"),
	}

	mux := http.NewServeMux()

	// TODO: Add middleware for logging and rate limiting as needed.
	mux.Handle("/app/", cfg.MiddlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))

	// Method specific routing. [METHOD ][HOST]/[PATH]
	mux.HandleFunc("GET  /api/healthz", cfg.HandlerReadiness) // Basic health check

	mux.HandleFunc("POST /api/users", cfg.HandleCreateUser) // Creates a new user.
	mux.HandleFunc("PUT /api/users", cfg.HandleUpdateUser) // Update a user with new email, password, or both.
	mux.HandleFunc("POST /api/login", cfg.HandleUserLogin) // Handles a user login

	mux.HandleFunc("POST /api/chirps", cfg.HandleCreateChirp) // Create a chirp.
	mux.HandleFunc("GET  /api/chirps", cfg.HandleGetAllChirps) // Get all chirps.
	mux.HandleFunc("GET  /api/chirps/{chirpID}", cfg.HandleGetChirp) // Get a single chirp.
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", cfg.HandleDeleteChirp) // Delete a chirp.

	mux.HandleFunc("POST /api/refresh", cfg.CheckRefreshToken) // Check refresh token
	mux.HandleFunc("POST /api/revoke", cfg.RevokeRefreshToken) // Revoke refresh token

	// Admin routes
	mux.HandleFunc("GET  /admin/metrics", cfg.HandlerMetrics)
	mux.HandleFunc("POST /admin/reset", cfg.HandlerReset)

	mux.HandleFunc("POST /api/polka/webhooks", cfg.HandlePolkaWebhook)

	/* A http.Server is a struct that describes a server configuration */
	server := http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	fmt.Printf("Server started on: localhost:%v/app/\n", port)

	/* ListenAndServe() blocks the main function until the server shuts down or an
	unexpected error crashes it. */
	log.Fatal(server.ListenAndServe())
} // End main() func
