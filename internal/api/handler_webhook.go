package api

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

const (
	eventUserUpgraded = "user.upgraded"
)

func (cfg *APIConfig) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	type request struct {
    Event string `json:"event"`
		Data  struct{
      UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	defer r.Body.Close()

	req := request{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if req.Event != eventUserUpgraded {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	_, err := cfg.DB.UpdateEchoesRed(r.Context(), req.Data.UserID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

