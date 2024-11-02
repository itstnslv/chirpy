package main

import (
	"errors"
	"github.com/itstnslv/chirpy/internal/auth"
	"net/http"
)

func (cfg *apiConfig) revokeTokenHandler(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if errors.Is(err, auth.ErrNoAuthHeaderIncluded) {
		respondWithErr(w, http.StatusBadRequest, "Couldn't find token", err)
		return
	}
	err = cfg.db.RevokeToken(r.Context(), refreshToken)
	if err != nil {
		respondWithErr(w, http.StatusInternalServerError, "couldn't revoke refreshToken in db", err)
		return
	}
	respondWithJSON(w, http.StatusNoContent, nil)
}
