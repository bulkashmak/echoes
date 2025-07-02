package auth

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	log.Println("password hashed successfuly")

	return string(hash), nil
}

func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func GetBearerToken(headers http.Header) (string, error) {
	token := headers.Get("Authorization")
	if token == "" {
		return "", errors.New("'Authorization' header not found")
	}
	if !strings.HasPrefix(token, "Bearer ") {
		return "", errors.New("invalid token")
	}
	return strings.TrimPrefix(token, "Bearer "), nil
}
