package api

import (
	"github.com/bulkashmak/echoes/internal/auth"
	"net/http"
)

type RefreshResponse struct {
	Token string `json:"token"`
}

func (cfg *APIConfig) HandleRefresh(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	refreshTokenStr, err := auth.GetBearerToken(r.Header)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	refreshToken, err := cfg.DB.GetRefreshToken(r.Context(), refreshTokenStr)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	if refreshToken.RevokedAt.Valid {
		RespondWithError(w, http.StatusUnauthorized, "refresh token is revoked")
		return
	}

	user, err := cfg.DB.GetUserFromRefreshToken(r.Context(), refreshTokenStr)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	accessToken, err := auth.MakeJWT(user.ID, cfg.AuthSecret)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := RefreshResponse{
		Token: accessToken,
	}
	RespondWithJSON(w, http.StatusOK, resp)
}
