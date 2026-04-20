package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/ewpt3ch/chirpy/internal/auth"
	"github.com/ewpt3ch/chirpy/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Email       string    `json:"email"`
	IsChirpyRed bool      `json:"is_chirpy_red"`
}

func (c *apiConfig) handlerCreateUser(w http.ResponseWriter, req *http.Request) {
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

	id := uuid.New()
	timeNow := time.Now()
	email := reqParams.Email
	password, err := auth.HashPassword(reqParams.Password)
	if err != nil {
		respondWithError(w, 500, "password issue")
		return
	}

	newUser := database.CreateUserParams{
		ID:             id,
		CreatedAt:      timeNow,
		UpdatedAt:      timeNow,
		Email:          email,
		HashedPassword: password,
	}

	dbuser, err := c.db.CreateUser(req.Context(), newUser)
	if err != nil {
		log.Printf("Database error: %v", err)
		respondWithError(w, 500, "Failed to create user in database")
		return
	}

	user := User{
		ID:          dbuser.ID,
		CreatedAt:   dbuser.CreatedAt,
		UpdatedAt:   dbuser.UpdatedAt,
		Email:       dbuser.Email,
		IsChirpyRed: dbuser.IsChirpyRed,
	}

	respondWithJSON(w, 201, user)

}
