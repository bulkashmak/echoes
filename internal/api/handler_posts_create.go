package api

import (
	"encoding/json"
	"github.com/bulkashmak/echoes/internal/auth"
	"github.com/bulkashmak/echoes/internal/database"
	"github.com/google/uuid"
	"net/http"
	"strings"
	"time"
)

type CreatePostRequest struct {
	Body string `json:"body"`
}

type Post struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

var profaneWords = map[string]bool{
	"kerfuffle": true,
	"sharbert":  true,
	"fornax":    true,
}

func (cfg *APIConfig) HandleCreatePost(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	userID := authenticate(w, r, cfg)

	var req CreatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if len(req.Body) > 140 {
		RespondWithError(w, http.StatusBadRequest, "Post is too long")
		return
	}

	user, err := cfg.DB.GetUserByID(r.Context(), userID)
	if err != nil {
		RespondWithError(w, http.StatusNotFound, "User not found")
		return
	}

	cleanedBody := cleanPost(req.Body)
	post, err := cfg.DB.CreatePost(r.Context(), database.CreatePostParams{
		Body:   cleanedBody,
		UserID: user.ID,
	})
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondWithJSON(w, http.StatusCreated, Post{
		ID:        post.ID,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
		Body:      post.Body,
		UserID:    post.UserID,
	})
}

func cleanPost(body string) string {
	words := strings.Split(body, " ")
	for i, word := range words {
		lower := strings.ToLower(word)
		if profaneWords[lower] {
			words[i] = "****"
		}
	}
	return strings.Join(words, " ")
}

func authenticate(w http.ResponseWriter, r *http.Request, cfg *APIConfig) uuid.UUID {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, err.Error())
	}

	userID, err := auth.ValidateJWT(token, cfg.AuthSecret)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, err.Error())
	}

	return userID
}
