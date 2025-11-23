package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/schwarzschild/relay/internal/config"
	"github.com/schwarzschild/relay/internal/handlers"
)

// corsMiddleware adds CORS headers to the response
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Max-Age", "3600")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	mux := http.NewServeMux()

	// Initialize handlers with config
	relayHandler := handlers.NewRelayHandler(cfg)

	// Register endpoints
	mux.HandleFunc("POST /relay", relayHandler.HandleRelay)

	// Health check endpoint
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status":"ok"}`)
	})

	// Wrap mux with CORS middleware
	handler := corsMiddleware(mux)

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Server starting on %s", addr)

	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
