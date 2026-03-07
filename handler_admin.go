package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) HandlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(w, http.StatusForbidden, "Forbidden", nil)
		return
	}

	cfg.db.DeleteUsers(r.Context())
	cfg.fileServerHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
} // End HandlerReset() func

func (cfg *apiConfig) HandlerMetrics(w http.ResponseWriter, r *http.Request) {
	// Display metrics to the admin user
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte(fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", cfg.fileServerHits.Load())))
} // End handlerMetrics() func