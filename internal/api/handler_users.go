package api

import (
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type CreateUserRequest struct {
	Email string `json:"email"`
}

type CreateUserResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *APIConfig) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := cfg.DB.CreateUser(r.Context(), req.Email)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondWithJSON(w, http.StatusCreated, CreateUserResponse{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	})
}
