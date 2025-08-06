package config

import (
	"log"
	"os"
)

// Config holds all application configuration.
type Config struct {
	DatabaseURL string
	Port        string
}

// New loads configuration from environment variables. It will exit if a required
// variable like DATABASE_URL is not set.
func New() *Config {
	// Get DATABASE_URL from environment; exit if not found.
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("FATAL: DATABASE_URL environment variable is not set")
	}

	// Get PORT from environment, or use a default.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
	}

	return &Config{
		DatabaseURL: dbURL,
		Port:        port,
	}
}
