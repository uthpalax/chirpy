package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	"github.com/utphalax/chirpy/internal/database"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
	jwtSecret      string
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handleMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	hits := fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", cfg.fileserverHits.Load())
	w.Write([]byte(hits))
}

func main() {
	godotenv.Load()

	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	jwtSecret := os.Getenv("JWT_SECRET")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("cannot connect to database")
	}

	dbQueries := database.New(db)

	mux := http.NewServeMux()

	apiConfig := &apiConfig{}
	apiConfig.db = dbQueries
	apiConfig.platform = platform
	apiConfig.jwtSecret = jwtSecret

	handler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))

	mux.Handle("/app/", apiConfig.middlewareMetricsInc(handler))

	mux.HandleFunc("GET /api/healthz", handleRediness)
	mux.HandleFunc("GET /admin/metrics", apiConfig.handleMetrics)
	mux.HandleFunc("POST /admin/reset", apiConfig.handleReset)
	mux.HandleFunc("POST /api/users", apiConfig.handleCreateUser)
	mux.HandleFunc("POST /api/login", apiConfig.handleLogin)
	mux.HandleFunc("POST /api/refresh", apiConfig.handleRefreshToken)
	mux.HandleFunc("POST /api/revoke", apiConfig.handleRevokeToken)
	mux.HandleFunc("POST /api/chirps", apiConfig.handleCreateChirps)
	mux.HandleFunc("GET /api/chirps", apiConfig.handleGetChirps)
	mux.HandleFunc("GET /api/chirps/{id}", apiConfig.handleGetChirp)

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	server.ListenAndServe()
}
