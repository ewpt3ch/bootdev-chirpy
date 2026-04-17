package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/ewpt3ch/chirpy/internal/auth"
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

	user := User{
		ID:        dbuser.ID,
		CreatedAt: dbuser.CreatedAt,
		UpdatedAt: dbuser.UpdatedAt,
		Email:     dbuser.Email,
	}

	respondWithJSON(w, 200, user)

}
