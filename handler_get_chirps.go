package main

import (
	"github.com/itstnslv/chirpy/internal/database"
	"net/http"
)

func (cfg *apiConfig) getChirpsHandler(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.ListChirps(r.Context())
	if err != nil {
		respondWithErr(w, http.StatusInternalServerError, "Couldn't retrieve chirps", err)
		return
	}
	chirpsJSON := wrapWithJsonTags(chirps)
	respondWithJSON(w, http.StatusOK, chirpsJSON)
}

func wrapWithJsonTags(dbChirps []database.Chirp) []Chirp {
	var chirps []Chirp
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, Chirp{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			Body:      dbChirp.Body,
			UserID:    dbChirp.UserID,
		})
	}
	return chirps
}
