package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func handlerChirpsValidation(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithErr(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithErr(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	cleaned := applyProfaneFilter(params.Body)
	respondWithJSON(w, http.StatusOK, returnVals{CleanedBody: cleaned})
}

func applyProfaneFilter(body string) string {
	profaneWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	words := strings.Split(body, " ")
	for i, word := range words {
		if _, ok := profaneWords[strings.ToLower(word)]; ok {
			words[i] = "****"
		}
	}
	return strings.Join(words, " ")
}
