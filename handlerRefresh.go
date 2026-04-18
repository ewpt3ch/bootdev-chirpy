package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ewpt3ch/chirpy/internal/auth"
)

func (c *apiConfig) handlerRefresh(w http.ResponseWriter, req *http.Request) {
	intoken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, 401, "invalid auth header")
		return
	}

	userID, err := c.db.GetUserFromRefreshToken(req.Context(), intoken)
	if err != nil {
		fmt.Printf("database error: %v", err)
		respondWithError(w, 401, "invalid refresh token")
		return
	}

	token, err := auth.MakeJWT(userID, c.jwt_secret, time.Hour)
	if err != nil {
		respondWithError(w, 500, "failed to create jwt")
		return
	}

	type refreshResponse struct {
		Token string `json:"token"`
	}
	refresh := refreshResponse{
		Token: token,
	}

	respondWithJSON(w, 200, refresh)
}
