package main

import (
	"errors"
	"github.com/itstnslv/chirpy/internal/auth"
	"net/http"
	"time"
)

func (cfg *apiConfig) refreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	type newAccessToken struct {
		NewToken string `json:"token"`
	}
	token, err := auth.GetBearerToken(r.Header)
	if errors.Is(err, auth.ErrNoAuthHeaderIncluded) {
		respondWithErr(w, http.StatusBadRequest, "no auth header found", err)
		return
	}
	tokenData, err := cfg.db.FindRefreshToken(r.Context(), token)
	if err != nil {
		respondWithErr(w, http.StatusUnauthorized, "refresh token not found in db", err)
		return
	}
	if time.Now().After(tokenData.ExpiresAt) {
		respondWithErr(w, http.StatusUnauthorized, "refresh token expired", err)
		return
	}
	if tokenData.RevokedAt.Valid {
		respondWithErr(w, http.StatusUnauthorized, "refresh token revoked", err)
		return
	}
	jwt, err := auth.MakeJWT(tokenData.UserID, cfg.secret, time.Hour)
	if err != nil {
		respondWithErr(w, http.StatusInternalServerError, "couldn't generate new access token", err)
		return
	}
	respondWithJSON(w, http.StatusOK, newAccessToken{NewToken: jwt})
}
