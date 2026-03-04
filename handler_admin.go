package main

import (
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