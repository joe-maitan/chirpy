package main

import (
	"net/http"
)

func (cfg *apiConfig) HandlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(w, 403, "Forbidden", nil)
		return
	}

	cfg.db.DeleteUsers(r.Context())
	cfg.fileServerHits.Store(0)
	w.WriteHeader(200)
	w.Write([]byte("Hits reset to 0"))
} // End HandlerReset() func