package api

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/bulkashmak/echoes/internal/database"
	"github.com/bulkashmak/echoes/internal/auth"
	"log"
	"net/http"
	"os"
	"sync/atomic"
)

type APIConfig struct {
	DB             *database.Queries
	AuthSecret     string
	FileServerHits atomic.Int32
	PolkaKey       string
}

func (cfg *APIConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.FileServerHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *APIConfig) HandleMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	html := fmt.Sprintf(`
		<html>
		  <body>
		    <h1>Welcome, Chirpy Admin</h1>
		    <p>Chirpy has been visited %d times!</p>
		  </body>
		</html>
	`, cfg.FileServerHits.Load())

	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(html))
	if err != nil {
		log.Fatal(err)
	}
}

func (cfg *APIConfig) HandleReset(w http.ResponseWriter, r *http.Request) {
	cfg.FileServerHits.Store(0)
	log.Println("file server hits reset")

	env := os.Getenv("ENV")
	if env == "" {
		RespondWithError(w, http.StatusInternalServerError, "Environment variable not set")
		return
	} else if env != "dev" {
		w.WriteHeader(http.StatusForbidden)
	}

	err := cfg.DB.DeleteAllUsers(r.Context())
	if err != nil {
		log.Fatal(err)
	}
	log.Println("all users deleted")
}

func RespondWithError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(map[string]string{"error": msg})
	if err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

func RespondWithJSON(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(payload)
	if err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

func authenticate(w http.ResponseWriter, r *http.Request, cfg *APIConfig) uuid.UUID {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, err.Error())
		return uuid.Nil
	}

	userID, err := auth.ValidateJWT(token, cfg.AuthSecret)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, err.Error())
		return uuid.Nil
	}

	log.Printf("user with id '%s' authenticated successfuly", userID)

	return userID
}

