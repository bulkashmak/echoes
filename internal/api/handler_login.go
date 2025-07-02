package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/bulkashmak/echoes/internal/auth"
	"github.com/bulkashmak/echoes/internal/database"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	User
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
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

	token, err := auth.MakeJWT(user.ID, cfg.AuthSecret)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	refreshToken, err := createRefreshToken(r.Context(), cfg, user.ID)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondWithJSON(w, http.StatusOK, LoginResponse{
		User: User{
			ID:          user.ID,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
			Email:       user.Email,
			IsEchoesRed: user.IsEchoesRed,
		},
		Token:        token,
		RefreshToken: refreshToken.Token,
	})
}

func createRefreshToken(ctx context.Context, cfg *APIConfig, userID uuid.UUID) (database.RefreshToken, error) {
	ttl := 60 * 24 * time.Hour
	refreshTokenStr, err := auth.MakeRefreshToken()
	if err != nil {
		return database.RefreshToken{}, err
	}

	refreshToken, err := cfg.DB.CreateRefreshToken(ctx, database.CreateRefreshTokenParams{
		Token:     refreshTokenStr,
		UserID:    userID,
		ExpiresAt: time.Now().Add(ttl),
		RevokedAt: sql.NullTime{},
	})
	if err != nil {
		return database.RefreshToken{}, err
	}

	return refreshToken, nil
}
