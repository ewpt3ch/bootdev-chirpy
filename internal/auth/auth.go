package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}

	return hashedPassword, nil
}

func CheckPasswordHash(password, hash string) (bool, error) {
	return argon2id.ComparePasswordAndHash(password, hash)
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {

	timeNow := time.Now()

	t := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.RegisteredClaims{
			Issuer:    "chirpy-access",
			IssuedAt:  jwt.NewNumericDate(timeNow),
			ExpiresAt: jwt.NewNumericDate(timeNow.Add(expiresIn)),
			Subject:   userID.String(),
		})

	s, err := t.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}
	return s, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	id, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}

	userID, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil, err
	}

	return userID, nil

}

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	prefix := "Bearer "
	if !strings.HasPrefix(authHeader, prefix) {
		return "", errors.New("invalid authorization header")
	}
	token := strings.TrimPrefix(authHeader, prefix)
	if token == "" {
		return "", errors.New("no authorization header")
	}

	return token, nil
}

func MakeRefreshToken() string {
	rseed := make([]byte, 32)
	rand.Read(rseed)
	rtoken := hex.EncodeToString(rseed)
	return rtoken
}

func GetApiKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	prefix := "ApiKey "
	if !strings.HasPrefix(authHeader, prefix) {
		return "", errors.New("invalid authorization header")
	}
	token := strings.TrimPrefix(authHeader, prefix)
	if token == "" {
		return "", errors.New("no authorization header")
	}

	return token, nil
}
