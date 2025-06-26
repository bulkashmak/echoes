package api

import (
	"encoding/json"
	"github.com/bulkashmak/echoes/internal/database"
	"github.com/google/uuid"
	"net/http"
	"strings"
	"time"
)

type CreatePostRequest struct {
	Body   string    `json:"body"`
	UserID uuid.UUID `json:"user_id"`
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

	var req CreatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if len(req.Body) > 140 {
		RespondWithError(w, http.StatusBadRequest, "Post is too long")
		return
	}

	_, err := cfg.DB.GetUserByID(r.Context(), req.UserID)
	if err != nil {
		RespondWithError(w, http.StatusNotFound, "User not found")
	}

	cleanedBody := cleanPost(req.Body)
	post, err := cfg.DB.CreatePost(r.Context(), database.CreatePostParams{
		Body:   cleanedBody,
		UserID: req.UserID,
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
