package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ewpt3ch/chirpy/internal/auth"
	"github.com/ewpt3ch/chirpy/internal/database"
)

func (c *apiConfig) handlerUpdateUser(w http.ResponseWriter, req *http.Request) {
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
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(req.Body)
	reqParams := reqParameters{}
	err = decoder.Decode(&reqParams)
	if err != nil {
		respondWithError(w, 500, "Somethings went wrong json")
		return
	}

	if len(reqParams.Email) == 0 {
		respondWithError(w, 400, "No email in request")
		return
	}

	if len(reqParams.Password) == 0 {
		respondWithError(w, 400, "No password in request")
		return
	}

	password, err := auth.HashPassword(reqParams.Password)
	if err != nil {
		respondWithError(w, 500, "password issue")
		return
	}

	updateParams := database.UpdateUserParams{
		ID:             userID,
		Email:          reqParams.Email,
		HashedPassword: password,
		UpdatedAt:      time.Now(),
	}

	dbuser, err := c.db.UpdateUser(req.Context(), updateParams)

	user := User{
		ID:        dbuser.ID,
		CreatedAt: dbuser.CreatedAt,
		UpdatedAt: dbuser.UpdatedAt,
		Email:     dbuser.Email,
	}

	respondWithJSON(w, 200, user)

}
