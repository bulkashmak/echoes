package api

import (
	"net/http"
	"github.com/google/uuid"
)

func (cfg *APIConfig) HandleDeletePost(w http.ResponseWriter, r *http.Request) {
  defer r.Body.Close()

	userID := authenticate(w, r, cfg)
	if userID == uuid.Nil {
		return
	}

	postID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "path var is invalid")
		return
	}

	post, err := cfg.DB.RetrievePostByID(r.Context(), postID)
	if err != nil {
		RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	if post.UserID != userID {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	err = cfg.DB.DeletePostByID(r.Context(), post.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

