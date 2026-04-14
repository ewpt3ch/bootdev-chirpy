package main

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"
)

func handlerValidateChirp(w http.ResponseWriter, req *http.Request) {
	// placeholder for profane words
	theProfane := []string{"kerfuffle", "sharbert", "fornax"}

	type reqParameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(req.Body)
	reqParams := reqParameters{}
	err := decoder.Decode(&reqParams)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}

	if len(reqParams.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	type cleanedParams struct {
		Cleaned_Body string `json:"cleaned_body"`
	}

	censoredChirp := replaceProfane(reqParams.Body, theProfane)
	respondWithJSON(w, 200, cleanedParams{Cleaned_Body: censoredChirp})

}

func replaceProfane(chirp string, theProfane []string) string {
	words := strings.Split(chirp, " ")
	for i, word := range words {
		if slices.Contains(theProfane, strings.ToLower(word)) {
			words[i] = "****"
		}
	}
	return strings.Join(words, " ")
}
