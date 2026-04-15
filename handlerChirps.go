package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/ewpt3ch/chirpy/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (c *apiConfig) handlerChirps(w http.ResponseWriter, req *http.Request) {
	type reqParameters struct {
		Body    string    `json:"body"`
		User_ID uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(req.Body)
	reqParams := reqParameters{}
	err := decoder.Decode(&reqParams)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}

	chirpBody, err := validateChirp(reqParams.Body)
	if err != nil {
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	id := uuid.New()
	timeNow := time.Now()
	body := chirpBody
	user_id := reqParams.User_ID

	newChirp := database.CreateChirpParams{
		ID:        id,
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
		Body:      body,
		UserID:    user_id,
	}

	dbChirp, err := c.db.CreateChirp(req.Context(), newChirp)
	if err != nil {
		log.Printf("Database error: %v", err)
		respondWithError(w, 500, "Failed to create chirp")
		return
	}

	chirp := Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	}

	respondWithJSON(w, 201, chirp)

}

func validateChirp(rawChirp string) (string, error) {
	// placeholder for profane words
	theProfane := []string{"kerfuffle", "sharbert", "fornax"}

	if len(rawChirp) > 140 {
		return "", errors.New("Chirp is too long")
	}

	type cleanedParams struct {
		Cleaned_Body string `json:"cleaned_body"`
	}

	censoredChirp := replaceProfane(rawChirp, theProfane)
	return censoredChirp, nil
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
