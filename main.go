package main

import (
	"database/sql"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/itstnslv/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"
)

type apiConfig struct {
	fileServerHits atomic.Int32
	db             *database.Queries
}

type User struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
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
	dbQueries := database.New(dbConn)

	apiCfg := apiConfig{
		fileServerHits: atomic.Int32{},
		db:             dbQueries,
	}
	mux := http.NewServeMux()
	fsHandler := http.StripPrefix("/app", apiCfg.middlewareMetricInc(http.FileServer(http.Dir("."))))
	mux.Handle("/app/", fsHandler)
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("POST /api/validate_chirp", handlerChirpsValidation)

	mux.HandleFunc("POST /api/users", func(w http.ResponseWriter, r *http.Request) {
		type parameters struct {
			Email string `json:"email"`
		}

		decoder := json.NewDecoder(r.Body)
		params := parameters{}
		err := decoder.Decode(&params)
		if err != nil {
			respondWithErr(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
			return
		}

		user, err := apiCfg.db.CreateUser(r.Context(), params.Email)
		if err != nil {
			respondWithErr(w, http.StatusInternalServerError, "Couldn't create user", err)
			return
		}
		respondWithJSON(w, http.StatusCreated, User{
			user.ID,
			user.CreatedAt,
			user.UpdatedAt,
			user.Email,
		})
	})

	mux.HandleFunc("GET /admin/metrics", apiCfg.handleMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handleReset)

	server := &http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	log.Printf("Listening on port %s", server.Addr)
	log.Fatal(server.ListenAndServe())
}
