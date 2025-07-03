package api

import (
	"sort"
	"net/http"

	"github.com/bulkashmak/echoes/internal/database"
	"github.com/google/uuid"
)

func (cfg *APIConfig) HandleRetrievePosts(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	authorID := r.URL.Query().Get("author_id")
	ssort := r.URL.Query().Get("sort")

  var posts []database.Post
	var err error

	if authorID != "" {
		userID, parseErr := uuid.Parse(authorID)
		if parseErr != nil {
			RespondWithError(w, http.StatusBadRequest, "Invalid author_id")
			return
		}
		posts, err = cfg.DB.ListPostsByAuthor(r.Context(), userID)
	} else {
		posts, err = cfg.DB.ListPosts(r.Context())
	}
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if ssort != "" {
		if ssort == "asc" {
			sort.Slice(posts, func(i, j int) bool {
				return posts[i].CreatedAt.Before(posts[j].CreatedAt)
			})
		} else if ssort == "desc" { 
			sort.Slice(posts, func(i, j int) bool {
				return posts[j].CreatedAt.Before(posts[i].CreatedAt)
			})
		}
	}

	var resp []Post

	for _, post := range posts {
		resp = append(resp, Post{
			ID:        post.ID,
			CreatedAt: post.CreatedAt,
			UpdatedAt: post.UpdatedAt,
			Body:      post.Body,
			UserID:    post.UserID,
		})
	}

	RespondWithJSON(w, http.StatusOK, resp)
}

func (cfg *APIConfig) HandleRetrievePostByID(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	postID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid chirp ID")
		return
	}

	post, err := cfg.DB.RetrievePostByID(r.Context(), postID)
	if err != nil {
		RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	RespondWithJSON(w, http.StatusOK, Post{
		ID:        post.ID,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
		Body:      post.Body,
		UserID:    post.UserID,
	})
}

