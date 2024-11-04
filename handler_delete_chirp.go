package main

import (
	"github.com/google/uuid"
	"github.com/itstnslv/chirpy/internal/auth"
	"net/http"
)

func (cfg *apiConfig) deleteChirpHandler(w http.ResponseWriter, r *http.Request) {
	pathValue := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(pathValue)
	if err != nil {
		respondWithErr(w, http.StatusBadRequest, "couldn't parse chirp id", err)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithErr(w, http.StatusUnauthorized, "token missing", err)
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithErr(w, http.StatusUnauthorized, "invalid token", err)
		return
	}

	chirp, err := cfg.db.GetChirpById(r.Context(), chirpID)
	if err != nil {
		respondWithErr(w, http.StatusNotFound, "chirp not found", err)
		return
	}

	if chirp.UserID != userID {
		respondWithErr(w, http.StatusForbidden, "action not allowed", err)
		return
	}

	if err := cfg.db.DeleteChirpById(r.Context(), chirpID); err != nil {
		respondWithErr(w, http.StatusInternalServerError, "failed to delete chirp", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}
