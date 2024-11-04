package main

import (
	"encoding/json"
	"github.com/itstnslv/chirpy/internal/auth"
	"github.com/itstnslv/chirpy/internal/database"
	"net/http"
)

func (cfg *apiConfig) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithErr(w, http.StatusUnauthorized, "auth token not parsed", err)
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithErr(w, http.StatusUnauthorized, "auth token not valid", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		respondWithErr(w, http.StatusBadRequest, "invalid json", err)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithErr(w, http.StatusInternalServerError, "unable to hash password", err)
		return
	}

	user, err := cfg.db.UpdateUser(r.Context(), database.UpdateUserParams{
		ID:             userID,
		Email:          params.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respondWithErr(w, http.StatusInternalServerError, "unable to update user", err)
		return
	}

	respondWithJSON(w, http.StatusOK, User{
		Id:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	})

}
