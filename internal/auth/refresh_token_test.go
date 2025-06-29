package auth_test

import (
	"github.com/bulkashmak/echoes/internal/auth"
	"testing"
)

func TestMakeRefreshToken(t *testing.T) {
	_, err := auth.MakeRefreshToken()
	if err != nil {
		t.Fatalf("expected err to be nil, got %v", err)
	}
}
