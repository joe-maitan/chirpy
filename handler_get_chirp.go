package main

import (
	"sort"
	"net/http"

	"github.com/google/uuid"
	"github.com/joe-maitan/chirpy/internal/database"
)

func (cfg *apiConfig) HandleGetAllChirps(w http.ResponseWriter, r *http.Request) {
	var DBChirps []database.Chirp
	var err error

	authorID := r.URL.Query().Get("author_id")
	sortOrder := r.URL.Query().Get("sort")
	if authorID != "" {
		// logStatement("handler_get_chirp.go", "HandleGetAllChirps()", fmt.Sprintf("authorID=%v", authorID), nil)
		data, err := uuid.Parse(authorID)
		if err != nil {
			respondWithError(w, http.StatusNotFound, "Error parsing authorID as a uuid", err)
			return
		}

		DBChirps, err = cfg.db.GetAllChripsByUser(r.Context(), data)
		if err != nil {
			// logStatement("handler_get_chirp", "HandleGetAllChirps()", "Error getting all chirps", err)
			respondWithError(w, http.StatusInternalServerError, "Error getting all chirps", err)
			return
		}
	} else {
		DBChirps, err = cfg.db.GetAllChirps(r.Context())
		if err != nil {
			logStatement("handler_get_chirp", "HandleGetAllChirps()", "Error getting all chirps", err)
			respondWithError(w, http.StatusInternalServerError, "Error getting all chirps", err)
			return
		}
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

	if sortOrder == "asc" {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].CreatedAt.Before(chirps[j].CreatedAt)
		})
	} else if sortOrder == "desc" {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].CreatedAt.After(chirps[j].CreatedAt)
		})
	}

	respondWithJSON(w, http.StatusOK, chirps)
} // End HandleGetAllChirps() func

func (cfg *apiConfig) HandleGetChirp(w http.ResponseWriter, r *http.Request) {
	// log.Println(r.PathValue("chirpID"))
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
		logStatement("handler_get_chirp", "HandleGetChirp()", "Error finding chirp", err)
		respondWithError(w, http.StatusNotFound, "Error finding chirp", err)
		return
	}

	respondWithJSON(w, http.StatusOK, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
} // End HandleGetChirp() func