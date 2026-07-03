package main

import (
	"database/sql"
	"errors"
	"net/http"
	"sort"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerChirpsGet(w http.ResponseWriter, r *http.Request) {
	chirpIDString := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}

	// Fetch chirp by ID, return 404 if not found, 500 for database issues
	dbChirp, err := cfg.db.GetChirp(r.Context(), chirpID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "Couldn't get chirp", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve chirp", err)
		return
	}

	respondWithJSON(w, http.StatusOK, Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		UserID:    dbChirp.UserID,
		Body:      dbChirp.Body,
	})
}

func (cfg *apiConfig) handlerChirpsRetrieve(w http.ResponseWriter, r *http.Request) {
	authorIDString := r.URL.Query().Get("author_id")
	sortingOrder := r.URL.Query().Get("sort")

	if sortingOrder == "" {
		sortingOrder = "asc"
	}

	if authorIDString != "" {
		authorID, err := uuid.Parse(authorIDString)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "failed to parse", err)
			return
		}
		dbChirps, err := cfg.db.GetAuthorChirps(r.Context(), authorID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps", err)
			return
		}

		chirps := []Chirp{}
		for _, dbChirp := range dbChirps {
			chirps = append(chirps, Chirp{
				ID:        dbChirp.ID,
				CreatedAt: dbChirp.CreatedAt,
				UpdatedAt: dbChirp.UpdatedAt,
				UserID:    dbChirp.UserID,
				Body:      dbChirp.Body,
			})
		}

		sortedChirps(chirps, sortingOrder)

		respondWithJSON(w, http.StatusOK, chirps)
		return

	}

	dbChirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps", err)
		return
	}

	chirps := []Chirp{}
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, Chirp{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			UserID:    dbChirp.UserID,
			Body:      dbChirp.Body,
		})
	}
	sortedChirps(chirps, sortingOrder)

	respondWithJSON(w, http.StatusOK, chirps)
}

func sortedChirps(chirps []Chirp, order string) {
	if order == "asc" {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].CreatedAt.Before(chirps[j].CreatedAt)
		})
	} else if order == "desc" {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].CreatedAt.After(chirps[j].CreatedAt)
		})
	}
}
