package main

import (
	"database/sql"
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

func (c *apiConfig) handlerGetChirps(w http.ResponseWriter, req *http.Request) {
	dbChirps, err := c.db.GetChirps(req.Context())
	if err != nil {
		log.Printf("Database error %v", err)
		respondWithError(w, 500, "Failed to get chirps")
		return
	}

	var respChirps = []Chirp{}
	for _, chirp := range dbChirps {
		respChirps = append(respChirps, Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}

	respondWithJSON(w, 200, respChirps)

}

func (c *apiConfig) handlerGetChirpByID(w http.ResponseWriter, req *http.Request) {

	chirpID, err := uuid.Parse(req.PathValue("chirpID"))
	if err != nil {
		log.Printf("uuid error: %v", err)
		respondWithError(w, 400, "invalid chirp id")
		return
	}

	dbChirp, err := c.db.GetChirpByID(req.Context(), chirpID)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, 404, "not found")
			return
		}

		log.Printf("Database error %v", err)
		respondWithError(w, 500, "failed to get chirp")
		return
	}

	chirp := Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	}

	respondWithJSON(w, 200, chirp)

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
