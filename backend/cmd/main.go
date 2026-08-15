package main

import (
	"log"
	"os"

	"github.com/holyflow/backend/internal/config"
	"github.com/holyflow/backend/internal/server"
)

// @title HolyFlow API
// @version 1.0
// @description HolyFlow is a service for storing Christian song texts, their rhythm and mp3 examples
// @BasePath /api/v1

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize server
	srv := server.NewServer(cfg)

	// Run server
	if err := srv.Run(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
		os.Exit(1)
	}
}
