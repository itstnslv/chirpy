package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/itstnslv/chirpy/internal/auth"
	"net/http"
)

func (cfg *apiConfig) polkaWebhookHandler(w http.ResponseWriter, r *http.Request) {
	type event struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		}
	}

	key, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithErr(w, http.StatusUnauthorized, "API KEY missing", err)
		return
	}
	if key != cfg.polkaKey {
		respondWithErr(w, http.StatusUnauthorized, "invalid API KEY", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	incEvent := event{}
	if err := decoder.Decode(&incEvent); err != nil {
		respondWithErr(w, http.StatusInternalServerError, "couldn't parse event data", err)
		return
	}

	if incEvent.Event != "user.upgraded" {
		respondWithErr(w, http.StatusNoContent, "this type of event is not handled", nil)
		return
	}

	if err := cfg.db.UpgradeUserToChirpyRed(r.Context(), incEvent.Data.UserID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithErr(w, http.StatusNotFound, "user not found", err)
			return
		}
		respondWithErr(w, http.StatusInternalServerError, "couldn't upgrade user", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}
