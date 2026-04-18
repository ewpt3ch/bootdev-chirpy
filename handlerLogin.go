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

func (c *apiConfig) handlerLogin(w http.ResponseWriter, req *http.Request) {
	type reqParameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(req.Body)
	reqParams := reqParameters{}
	err := decoder.Decode(&reqParams)
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

	dbuser, err := c.db.GetUserByEmail(req.Context(), reqParams.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, 401, "incorrect email or password")
			return
		}
		respondWithError(w, 500, "internal server error")
		return
	}

	loginGood, err := auth.CheckPasswordHash(reqParams.Password, dbuser.HashedPassword)
	if err != nil {
		respondWithError(w, 500, "Something went wrong checkhash")
		return
	}

	if !loginGood {
		respondWithError(w, 401, "incorrect email or password")
		return
	}

	expiresIn := time.Hour

	token, err := auth.MakeJWT(dbuser.ID, c.jwt_secret, expiresIn)
	if err != nil {
		respondWithError(w, 500, "failed to create jwt")
		return
	}

	rtoken := auth.MakeRefreshToken()
	timeNow := time.Now()
	expiresAt := time.Now().Add(time.Hour * 60 * 24)

	rtokenParams := database.CreateRefreshTokenParams{
		Token:     rtoken,
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
		UserID:    dbuser.ID,
		ExpiresAt: expiresAt,
	}
	_, err = c.db.CreateRefreshToken(req.Context(), rtokenParams)
	if err != nil {
		respondWithError(w, 500, "failed to insert refreshtoken")
		return
	}

	type loginResponse struct {
		ID           uuid.UUID `json:"id"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
		Email        string    `json:"email"`
		Token        string    `json:"token"`
		RefreshToken string    `json:"refresh_token"`
	}

	user := loginResponse{
		ID:           dbuser.ID,
		CreatedAt:    dbuser.CreatedAt,
		UpdatedAt:    dbuser.UpdatedAt,
		Email:        dbuser.Email,
		Token:        token,
		RefreshToken: rtoken,
	}

	respondWithJSON(w, 200, user)

}
