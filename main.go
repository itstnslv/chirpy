package main

import (
	"database/sql"
	"github.com/itstnslv/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"sync/atomic"
)

type apiConfig struct {
	fileServerHits atomic.Int32
	db             *database.Queries
	platform       string
	secret         string
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL environment variable not set")
	}
	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error connecting to database: %s", err)
	}
	tokenSecret := os.Getenv("secret")
	if tokenSecret == "" {
		log.Fatal("secret environment variable not set")
	}

	apiCfg := apiConfig{
		fileServerHits: atomic.Int32{},
		db:             database.New(dbConn),
		platform:       os.Getenv("PLATFORM"),
		secret:         tokenSecret,
	}

	mux := http.NewServeMux()
	fsHandler := http.StripPrefix("/app", apiCfg.middlewareMetricInc(http.FileServer(http.Dir("."))))
	mux.Handle("/app/", fsHandler)
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("POST /api/chirps", apiCfg.createChirpHandler)
	mux.HandleFunc("GET /api/chirps", apiCfg.getChirpsHandler)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.getChirpByIdHandler)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.deleteChirpHandler)

	mux.HandleFunc("POST /api/users", apiCfg.createUserHandler)
	mux.HandleFunc("PUT /api/users", apiCfg.updateUserHandler)
	mux.HandleFunc("POST /api/login", apiCfg.loginUserHandler)
	mux.HandleFunc("POST /api/refresh", apiCfg.refreshTokenHandler)
	mux.HandleFunc("POST /api/revoke", apiCfg.revokeTokenHandler)

	mux.HandleFunc("GET /admin/metrics", apiCfg.handleMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handleReset)

	mux.HandleFunc("POST /api/polka/webhooks", apiCfg.polkaWebhookHandler)

	server := &http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	log.Printf("Listening on port %s", server.Addr)
	log.Fatal(server.ListenAndServe())
}
