package main

import (
	"encoding/json"
	"github.com/itstnslv/chirpy/internal/auth"
	"github.com/itstnslv/chirpy/internal/database"
	"net/http"
	"time"
)

func (cfg *apiConfig) loginUserHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	type loginResponse struct {
		User
		Token        string `json:"accessToken"`
		RefreshToken string `json:"refresh_token"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		respondWithErr(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	user, err := cfg.db.FindUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithErr(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}
	if err = auth.CheckPasswordHash(user.HashedPassword, params.Password); err != nil {
		respondWithErr(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	accessToken, err := auth.MakeJWT(user.ID, cfg.secret, time.Hour)
	if err != nil {
		respondWithErr(w, http.StatusInternalServerError, "Couldn't generate JWT", err)
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithErr(w, http.StatusInternalServerError, "Couldn't generate refresh token", err)
	}

	if _, err := cfg.db.CreateRefreshToken(r.Context(),
		database.CreateRefreshTokenParams{
			Token:     refreshToken,
			UserID:    user.ID,
			ExpiresAt: time.Now().Add(60 * 24 * time.Hour),
		}); err != nil {
		respondWithErr(w, http.StatusInternalServerError, "Couldn't create refresh token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, loginResponse{
		User{
			Id:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
		accessToken,
		refreshToken,
	})
}
