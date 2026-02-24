package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joe-maitan/chirpy/internal/api"
	"github.com/joe-maitan/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
} // End handlerReadiness() func

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

	apiCfg := api.Config{
		FileServerHits: atomic.Int32{},
		DB:             dbQueries,
		Platform:       os.Getenv("PLATFORM"),
		JWTSecret:      os.Getenv("SECRET"),
	}

	mux := http.NewServeMux()

	mux.Handle("/app/", apiCfg.MiddlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))

	// Method specific routing. [METHOD ][HOST]/[PATH]
	mux.HandleFunc("GET  /api/healthz", handlerReadiness)

	// Creates a new user
	mux.HandleFunc("POST /api/users", apiCfg.HandleCreateUser)

	// Handles a user login
	mux.HandleFunc("POST /api/login", apiCfg.HandleUserLogin)

	// Handles chirp creation, retrieval, and listing
	mux.HandleFunc("POST /api/chirps", apiCfg.HandleCreateChirp)
	mux.HandleFunc("GET  /api/chirps", apiCfg.HandleGetAllChirps)
	mux.HandleFunc("GET  /api/chirps/{chirpID}", apiCfg.HandleGetChirp)

	// Handles refresh token rotation and revocation
	mux.HandleFunc("POST /api/refresh", apiCfg.CheckRefreshToken)
	mux.HandleFunc("POST /api/revoke", apiCfg.RevokeRefreshToken)

	// Admin routes
	mux.HandleFunc("GET  /admin/metrics", apiCfg.HandlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.HandlerReset)

	/* A http.Server is a struct that describes a server configuration */
	server := http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	fmt.Printf("Server started on: %v...\n", port)

	/* ListenAndServe() blocks the main function until the server shuts down or an
	unexpected error crashes it. */
	log.Fatal(server.ListenAndServe())
} // End main() func
