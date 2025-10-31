package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/Raj7442/url-shortener/internal/handlers"
	"github.com/Raj7442/url-shortener/internal/storage"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	store := storage.NewInMemoryStore()
	h := handlers.NewHandler(store, "http://localhost:"+port)

	mux := http.NewServeMux()
	mux.HandleFunc("/api/shorten", h.ShortenHandler)
	mux.HandleFunc("/api/metrics", h.MetricsHandler)
	mux.HandleFunc("/", h.RedirectHandler)

	log.Printf("Server running on port %s", port)
	if err := http.ListenAndServe(":"+port, loggingMiddleware(mux)); err != nil {
		log.Fatal(err)
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

