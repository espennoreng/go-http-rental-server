package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type Env string

const (
	Development Env = "development"
	Staging     Env = "staging"
	Production  Env = "production"
)

// AppConfig holds the configuration for the application.
// We use yaml tags to map the YAML keys to our struct fields.
type AppConfig struct {
	Port                string `yaml:"port"`
	DatabaseURL         string `yaml:"database_url"`
	GoogleOAuthClientID string `yaml:"google_oauth_client_id"`
}

// file holds the structure of the entire YAML file.
type file struct {
	Default AppConfig `yaml:"default"`
	Dev     AppConfig `yaml:"dev"`
	Staging AppConfig `yaml:"staging"`
	Prod    AppConfig `yaml:"prod"`
}

// Load reads the configuration from a YAML file and environment variables.
func Load() (*AppConfig, error) {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev"
	}

	// **LOAD .ENV FILE HERE**
	// Load the general .env file first, then the environment-specific one.
	// This allows for a general .env and specific overrides.
	godotenv.Load() // Loads .env file if it exists
	envFile := fmt.Sprintf(".env.%s", env)
	if err := godotenv.Load(envFile); err != nil {
		log.Printf("Warning: Could not load %s file. Proceeding with existing environment variables.", envFile)
	}

	// Now, the rest of your function can proceed as it did before.
	// os.Getenv("DATABASE_URL") will now be populated if it was in the .env file.

	f, err := os.ReadFile("config.yaml")
	if err != nil {
		return nil, fmt.Errorf("could not read config file 'config.yaml': %w", err)
	}

	var cfgFile file
	if err := yaml.Unmarshal(f, &cfgFile); err != nil {
		return nil, fmt.Errorf("could not parse config file: %w", err)
	}

	appConfig := cfgFile.Default
	switch env {
	case "development":
		merge(&appConfig, cfgFile.Dev)
	case "staging":
		merge(&appConfig, cfgFile.Staging)
	case "production":
		merge(&appConfig, cfgFile.Prod)
	default:
		return nil, fmt.Errorf("invalid APP_ENV specified: '%s'", env)
	}

	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		appConfig.DatabaseURL = dbURL
	}

	if appConfig.DatabaseURL == "" {
		return nil, fmt.Errorf("database_url is a required config field")
	}
	if appConfig.GoogleOAuthClientID == "" {
		return nil, fmt.Errorf("google_oauth_client_id is a required config field")
	}

	return &appConfig, nil
}

// merge overwrites fields in the base config with fields from the override config.
// This allows environments to only specify the settings they want to change.
func merge(base *AppConfig, override AppConfig) {
	if override.Port != "" {
		base.Port = override.Port
	}
	if override.DatabaseURL != "" {
		base.DatabaseURL = override.DatabaseURL
	}
	if override.GoogleOAuthClientID != "" {
		base.GoogleOAuthClientID = override.GoogleOAuthClientID
	}
}