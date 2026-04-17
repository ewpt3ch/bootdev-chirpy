package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestPasswordsMatch(t *testing.T) {
	hashed, err := HashPassword("password")
	if err != nil {
		t.Fatalf("hash library failed: %v", err)
	}

	good, err := CheckPasswordHash("password", hashed)
	if err != nil {
		t.Fatalf("hash library failed: %v", err)
	}
	if !good {
		t.Error("expected passwords to match, they did not")
	}
}

func TestPasswordsNoMatch(t *testing.T) {
	hashed, err := HashPassword("password")
	if err != nil {
		t.Fatalf("hash library failed: %v", err)
	}

	good, err := CheckPasswordHash("BadWord", hashed)
	if err != nil {
		t.Fatalf("hash library failed: %v", err)
	}
	if good {
		t.Error("expected passwords to not match, they did")
	}
}

func TestJwtMatch(t *testing.T) {
	inUserID := uuid.New()
	jwt, err := MakeJWT(inUserID, "fakestring", 5*time.Minute)
	if err != nil {
		t.Fatalf("jwt failed to create: %v", err)
	}

	outUserID, err := ValidateJWT(jwt, "fakestring")
	if err != nil {
		t.Fatalf("jwt failed to validate: %v", err)
	}

	if inUserID != outUserID {
		t.Error("expected uuids to match, they did not")
	}
}
