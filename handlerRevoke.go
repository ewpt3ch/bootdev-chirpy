package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ewpt3ch/chirpy/internal/auth"
	"github.com/ewpt3ch/chirpy/internal/database"
)

func (c *apiConfig) handlerRevoke(w http.ResponseWriter, req *http.Request) {
	intoken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, 401, "invalid auth header")
		return
	}

	revokeParams := database.RevokeRefreshTokenParams{
		Token:     intoken,
		UpdatedAt: time.Now(),
	}

	revoked, err := c.db.RevokeRefreshToken(req.Context(), revokeParams)
	fmt.Println(revoked)
	respondWithJSON(w, 204, "")

}
