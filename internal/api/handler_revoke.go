package api

import (
	"database/sql"
	"github.com/bulkashmak/echoes/internal/auth"
	"github.com/bulkashmak/echoes/internal/database"
	"net/http"
	"time"
)

func (cfg *APIConfig) HandleRevoke(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	refreshTokenStr, err := auth.GetBearerToken(r.Header)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	refreshToken, err := cfg.DB.GetRefreshToken(r.Context(), refreshTokenStr)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	if refreshToken.RevokedAt.Valid {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	now := sql.NullTime{Time: time.Now(), Valid: true}
	err = cfg.DB.UpdateRefreshTokenRevokedAtByToken(r.Context(), database.UpdateRefreshTokenRevokedAtByTokenParams{
		RevokedAt: now,
		Token:     refreshTokenStr,
	})
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
	}

	w.WriteHeader(http.StatusNoContent)
}
