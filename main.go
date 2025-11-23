package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
  fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    cfg.fileserverHits.Add(1)
    next.ServeHTTP(w, r)
  })
}

func (cfg *apiConfig) handleMetrics(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "text/plain; charset=utf-8")
  w.WriteHeader(http.StatusOK)
  hits := fmt.Sprintf("Hits: %d", cfg.fileserverHits.Load())
  w.Write([]byte(hits))
}

func (cfg *apiConfig) handleReset(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "text/plain; charset=utf-8")
  w.WriteHeader(http.StatusOK)
  cfg.fileserverHits.Store(0)
  w.Write([]byte("Hits reset to 0"))
}

func main() {
  mux := http.NewServeMux()

  apiConfig := &apiConfig{}

  handler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))

  mux.Handle("/app/", apiConfig.middlewareMetricsInc(handler))

  mux.HandleFunc("/healthz", handleRediness)
  mux.HandleFunc("/metrics", apiConfig.handleMetrics)
  mux.HandleFunc("/reset", apiConfig.handleReset)

  server := http.Server {
    Addr: ":8080",
    Handler: mux,
  }

  server.ListenAndServe()
}

func handleRediness(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(http.StatusText(http.StatusOK)))
}