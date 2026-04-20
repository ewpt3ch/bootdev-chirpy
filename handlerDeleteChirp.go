package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/ewpt3ch/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (c *apiConfig) handlerDeleteChirp(w http.ResponseWriter, req *http.Request) {

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

	if userID != dbChirp.UserID {
		respondWithError(w, 403, "unauthorized")
		return
	}

	err = c.db.DeleteChirp(req.Context(), chirpID)
	if err != nil {
		respondWithError(w, 500, "delete failed")
		return
	}

	respondWithJSON(w, 204, "")

}
