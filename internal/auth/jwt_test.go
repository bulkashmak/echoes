package auth_test

import (
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/bulkashmak/echoes/internal/auth"
)

func TestJWTLifecycle(t *testing.T) {
	secret := "test-secret"
	userID := uuid.New()

	token, err := auth.MakeJWT(userID, secret, time.Minute)
	if err != nil {
		t.Fatalf("expected no error creating JWT, got: %v", err)
	}
	if token == "" {
		t.Fatal("expected token to be non-empty")
	}

	parsedID, err := auth.ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("expected no error validating JWT, got: %v", err)
	}
	if parsedID != userID {
		t.Errorf("expected userID %s, got %s", userID, parsedID)
	}
}

func TestJWTExpired(t *testing.T) {
	secret := "test-secret"
	userID := uuid.New()

	token, err := auth.MakeJWT(userID, secret, -1*time.Minute)
	if err != nil {
		t.Fatalf("expected no error creating expired JWT, got: %v", err)
	}

	_, err = auth.ValidateJWT(token, secret)
	if err == nil {
		t.Fatal("expected error for expired token, got none")
	}
}

func TestJWTWrongSecret(t *testing.T) {
	secret := "test-secret"
	wrongSecret := "wrong-secret"
	userID := uuid.New()

	token, err := auth.MakeJWT(userID, secret, time.Minute)
	if err != nil {
		t.Fatalf("expected no error creating JWT, got: %v", err)
	}

	_, err = auth.ValidateJWT(token, wrongSecret)
	if err == nil {
		t.Fatal("expected error for wrong secret, got none")
	}
}
