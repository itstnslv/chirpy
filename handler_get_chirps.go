package main

import (
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/itstnslv/chirpy/internal/database"
	"net/http"
)

func (cfg *apiConfig) getChirpsHandler(w http.ResponseWriter, r *http.Request) {
	authorId := r.URL.Query().Get("author_id")
	if authorId != "" {
		userID, _ := uuid.Parse(authorId)
		chirpsByAuthor, err := cfg.db.ListChirpsByAuthor(r.Context(), userID)
		processDbChirps(w, chirpsByAuthor, err)
		return
	}
	dbChirps, err := cfg.db.ListChirps(r.Context())
	processDbChirps(w, dbChirps, err)
}

func processDbChirps(w http.ResponseWriter, dbChirps []database.Chirp, err error) {
	if err != nil {
		respondWithErr(w, http.StatusInternalServerError, "couldn't retrieve chirps", err)
		return
	}
	var chirpsJson []Chirp
	for _, dbChirp := range dbChirps {
		chirpsJson = append(chirpsJson, wrapWithJsonTags(dbChirp))
	}
	respondWithJSON(w, http.StatusOK, chirpsJson)
}

func (cfg *apiConfig) getChirpByIdHandler(w http.ResponseWriter, r *http.Request) {
	pathValue := r.PathValue("chirpID")
	id, err := uuid.Parse(pathValue)
	if err != nil {
		respondWithErr(w, http.StatusBadRequest, "Invalid chirp id", nil)
		return
	}
	chirp, err := cfg.db.GetChirpById(r.Context(), id)
	if errors.Is(err, sql.ErrNoRows) {
		respondWithErr(w, http.StatusNotFound, "There's no such record", err)
		return
	}
	respondWithJSON(w, http.StatusOK, wrapWithJsonTags(chirp))
}

func wrapWithJsonTags(dbChirp database.Chirp) Chirp {
	return Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	}
}
