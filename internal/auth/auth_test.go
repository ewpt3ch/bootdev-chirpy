package auth

import (
	"net/http"
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

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name    string
		headers http.Header
		want    string
		wantErr bool
	}{
		{
			name:    "valid bearer token",
			headers: http.Header{"Authorization": []string{"Bearer mytoken123"}},
			want:    "mytoken123",
			wantErr: false,
		},
		{
			name:    "missing header",
			headers: http.Header{},
			want:    "",
			wantErr: true,
		},
		{
			name:    "no bearer prefix",
			headers: http.Header{"Authorization": []string{"mytoken123"}},
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetBearerToken(tt.headers)
			if (err != nil) != tt.wantErr {
				t.Errorf("got err %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
