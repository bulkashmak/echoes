package api

import (
	"encoding/json"
	"github.com/bulkashmak/echoes/internal/auth"
	"net/http"
	"time"
)

type LoginRequest struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	TokenTTLSeconds int    `json:"expires_in_seconds"`
}

type LoginResponse struct {
	User
	Token string `json:"token"`
}

func (cfg *APIConfig) HandleLogin(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	user, err := cfg.DB.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	err = auth.CheckPasswordHash(req.Password, user.PasswordHash)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	ttl := parseTTL(req.TokenTTLSeconds)
	token, err := auth.MakeJWT(user.ID, cfg.AuthSecret, ttl)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	RespondWithJSON(w, http.StatusOK, LoginResponse{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
		Token: token,
	})
}

func parseTTL(ttlSeconds int) time.Duration {
	if ttlSeconds <= 0 || ttlSeconds > 60 {
		return time.Hour
	}

	return time.Duration(ttlSeconds)
}
