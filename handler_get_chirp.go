package main

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/joe-maitan/chirpy/internal/database"
)

func (cfg *apiConfig) HandleGetAllChirps(w http.ResponseWriter, r *http.Request) {
	DBChirps, err := cfg.db.GetAllChirps(r.Context())
	if err != nil {
		log.Printf("api.go - HandleGetAllChirps() - Error getting all chirps: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Error getting all chirps", err)
		return
	}

	chirps := []database.Chirp{}
	for _, c := range DBChirps {
		chirps = append(chirps, database.Chirp{
			ID:        c.ID,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
			Body:      c.Body,
			UserID:    c.UserID,
		})
	}

	respondWithJSON(w, http.StatusOK, chirps)
} // End HandleGetAllChirps() func

func (cfg *apiConfig) HandleGetChirp(w http.ResponseWriter, r *http.Request) {
	log.Println(r.PathValue("chirpID"))
	chirpID := r.PathValue("chirpID")
	if chirpID == "" {
		respondWithError(w, http.StatusNotFound, "Invalid chirp ID", nil)
		return
	}

	data, err := uuid.Parse(chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Error parsing id as a uuid", err)
		return
	}

	chirp, err := cfg.db.GetChirp(r.Context(), data)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Error finding chirp", err)
		return
	}

	respondWithJSON(w, http.StatusOK, database.Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
} // End HandleGetChirp() func