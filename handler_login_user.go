package main

import (
	"encoding/json"
	"github.com/itstnslv/chirpy/internal/auth"
	"net/http"
)

func (cfg *apiConfig) loginUserHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithErr(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}
	user, err := cfg.db.FindUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithErr(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}
	err = auth.CheckPasswordHash(user.HashedPassword, params.Password)
	if err != nil {
		respondWithErr(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}
	respondWithJSON(w, http.StatusOK, User{
		Id:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})
}
