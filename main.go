package main

import (
	"database/sql"
	"github.com/bulkashmak/echoes/internal/api"
	"github.com/bulkashmak/echoes/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
)

const (
	port = "9000"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		panic("DB_URL environment variable not set")
	}
	db, err := sql.Open("postgres", dbURL)
	dbQueries := database.New(db)

	authSecret := os.Getenv("AUTH_SECRET")
	if authSecret == "" {
		panic("AUTH_SECRET environment variable not set")
	}

	apiCfg := api.APIConfig{
		DB:         dbQueries,
		AuthSecret: authSecret,
	}

	mux := http.NewServeMux()

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	mux.Handle("/app/", apiCfg.MiddlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))))
	mux.HandleFunc("POST /api/login", apiCfg.HandleLogin)
	mux.HandleFunc("POST /api/refresh", apiCfg.HandleRefresh)
	mux.HandleFunc("POST /api/revoke", apiCfg.HandleRevoke)
	mux.HandleFunc("POST /api/users", apiCfg.HandleCreateUser)
	mux.HandleFunc("PUT  /api/users", apiCfg.HandleUpdateUser)
	mux.HandleFunc("POST /api/chirps", apiCfg.HandleCreatePost)
	mux.HandleFunc("GET  /api/chirps", apiCfg.HandleRetrievePosts)
	mux.HandleFunc("GET  /api/chirps/{chirpID}", apiCfg.HandleRetrievePostByID)
	mux.HandleFunc("GET  /api/healthz", api.HandleReadiness)
	mux.HandleFunc("GET  /admin/metrics", apiCfg.HandleMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.HandleReset)

	log.Printf("listening on port %s\n", port)
	log.Fatal(server.ListenAndServe())
}

