package main

import "net/http"

func main() {
  mux := http.NewServeMux()

  mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir("."))))
  mux.HandleFunc("/healthz", handleRediness)

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
