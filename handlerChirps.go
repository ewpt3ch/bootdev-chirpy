package main

import (
	"encoding/json"
	"net/http"
)

func handlerValidateChirp(w http.ResponseWriter, req *http.Request) {

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

	type respParams struct {
		Valid bool `json:"valid"`
	}

	respondWithJSON(w, 200, respParams{Valid: true})
}
