package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/schwarzschild/relay/internal/config"
	"github.com/schwarzschild/relay/internal/handlers"
)

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

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Server starting on %s", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
