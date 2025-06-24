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
	db, err := sql.Open("postgres", dbURL)
	dbQueries := database.New(db)

	apiCfg := api.APIConfig{
		DB: dbQueries,
	}

	mux := http.NewServeMux()

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	mux.Handle("/app/", apiCfg.MiddlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))))
	mux.HandleFunc("POST /api/users", apiCfg.HandleCreateUser)
	mux.HandleFunc("POST /api/chirps", apiCfg.HandleCreatePost)
	mux.HandleFunc("GET  /api/healthz", api.HandleReadiness)
	mux.HandleFunc("GET  /admin/metrics", apiCfg.HandleMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.HandleReset)

	log.Printf("Listening on port %s\n", port)
	log.Fatal(server.ListenAndServe())
}
