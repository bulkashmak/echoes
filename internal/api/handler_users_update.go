package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/bulkashmak/echoes/internal/auth"
	"github.com/bulkashmak/echoes/internal/database"
	"github.com/google/uuid"
)

func (cfg *APIConfig) HandleUpdateUser(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type response struct {
		User
	}

	log.Println("received update user request")

  defer r.Body.Close()  

	userID := authenticate(w, r, cfg)
	if userID == uuid.Nil {
		w.WriteHeader(http.StatusUnauthorized)
	}

  log.Println("decoding request")

	req := request{}
  if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
    RespondWithError(w, http.StatusBadRequest, "failed to decode request")
		return
	}

  log.Println("hashing password")

	passwordHash, err := auth.HashPassword(req.Password)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "failed to hash password")
		return
	}

	log.Println("updating user")

	user, err := cfg.DB.UpdateUserEmailAndPasswordByID(r.Context(), database.UpdateUserEmailAndPasswordByIDParams{
		ID:           userID,
		Email:        req.Email,
		PasswordHash: passwordHash,
	})
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("failed to update user: %w", err).Error())
		return
	}
	
	log.Printf("user with id '%s' updated successfuly", user.ID)

  RespondWithJSON(w, http.StatusOK, response{
		User: User{
      ID:          user.ID,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
			Email:       user.Email,
			IsEchoesRed: user.IsEchoesRed,
		}, 
	})
}

