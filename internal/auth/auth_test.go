package auth_test

import (
	"github.com/bulkashmak/echoes/internal/auth"
	"net/http"
	"testing"
)

func TestGetBearerToken(t *testing.T) {
	headers := http.Header{}
	headers["Authorization"] = []string{"Bearer test-token"}

	token, err := auth.GetBearerToken(headers)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if token != "test-token" {
		t.Fatalf("expected token 'test-token', got: %s", token)
	}
}

func TestGetBearerTokenEmpty(t *testing.T) {
	headers := http.Header{}
	headers["Authorization"] = []string{""}

	_, err := auth.GetBearerToken(headers)
	if err == nil {
		t.Fatal("expected an error")
	}
}

func TestGetBearerTokenNoHeader(t *testing.T) {
	headers := http.Header{}

	_, err := auth.GetBearerToken(headers)
	if err == nil {
		t.Fatal("expected an error")
	}
}

func TestGetBearerTokenNoPrefix(t *testing.T) {
	headers := http.Header{}
	headers.Set("Authorization", "test-token")

	_, err := auth.GetBearerToken(headers)
	if err == nil {
		t.Fatal("expected an error")
	}
}
