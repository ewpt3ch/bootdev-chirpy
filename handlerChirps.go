package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/ewpt3ch/chirpy/internal/auth"
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

	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, 401, "invalid auth header")
		return
	}

	userID, err := auth.ValidateJWT(token, c.jwt_secret)
	if err != nil {
		respondWithError(w, 401, "invalid jwt")
		return
	}

	type reqParameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(req.Body)
	reqParams := reqParameters{}
	err = decoder.Decode(&reqParams)
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
	user_id := userID

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

	var dbChirps = []database.Chirp{}
	var err error

	if len(req.URL.Query().Get("author_id")) != 0 {
		UserID := req.URL.Query().Get("author_id")
		userID, err := uuid.Parse(UserID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "invalid uuid")
			return
		}
		dbChirps, err = c.db.GetChirpsByUserID(req.Context(), userID)
		if err != nil {
			if err == sql.ErrNoRows {
				respondWithError(w, http.StatusNotFound, "author not found")
				return
			}
			log.Printf("database error: %v", err)
			respondWithError(w, http.StatusInternalServerError, "database error")
			return
		}
	} else {
		dbChirps, err = c.db.GetChirps(req.Context())
		if err != nil {
			log.Printf("Database error %v", err)
			respondWithError(w, 500, "Failed to get chirps")
			return
		}
	}

	if len(req.URL.Query().Get("sort")) != 0 {
		order := req.URL.Query().Get("sort")
		if order == "desc" {
			sort.Slice(dbChirps, func(i, j int) bool { return dbChirps[i].CreatedAt.After(dbChirps[j].CreatedAt) })
		}
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
