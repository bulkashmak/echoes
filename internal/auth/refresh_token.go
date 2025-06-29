package auth

import (
	"crypto/rand"
	"encoding/hex"
)

func MakeRefreshToken() (string, error) {
	random := make([]byte, 32)
	_, err := rand.Read(random)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(random), nil
}
