package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/ewpt3ch/chirpy/internal/auth"
	"github.com/ewpt3ch/chirpy/internal/database"
	"github.com/google/uuid"
)

func (c *apiConfig) handlerChirpyRedUpgrade(w http.ResponseWriter, req *http.Request) {

	token, err := auth.GetApiKey(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid auth header")
		return
	}

	if token != c.polka_key {
		respondWithError(w, http.StatusUnauthorized, "invalid token")
		return
	}

	type reqParameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(req.Body)
	reqParams := reqParameters{}
	err = decoder.Decode(&reqParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "json decode failed")
		return
	}

	if reqParams.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	userID, err := uuid.Parse(reqParams.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	dbParams := database.UpgradeUserIdChirpyRedParams{
		ID:        userID,
		UpdatedAt: time.Now(),
	}

	err = c.db.UpgradeUserIdChirpyRed(req.Context(), dbParams)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, "user not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "database error")
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
