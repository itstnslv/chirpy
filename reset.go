package main

import "net/http"

func (cfg *apiConfig) handleReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Reset is only allowed in dev environment."))
		return
	}
	cfg.fileServerHits.Store(0)
	err := cfg.db.DeleteAllUsers(r.Context())
	if err != nil {
		respondWithErr(w, http.StatusInternalServerError, "Couldn't delete users from database", err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0. Users deleted."))
}
