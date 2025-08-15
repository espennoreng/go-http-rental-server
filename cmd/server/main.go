package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/espennoreng/go-http-rental-server/internal/api"
	"github.com/espennoreng/go-http-rental-server/internal/auth"
	"github.com/espennoreng/go-http-rental-server/internal/config"
	"github.com/espennoreng/go-http-rental-server/internal/repositories/postgres"
	"github.com/espennoreng/go-http-rental-server/internal/services"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// 1. Load application configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}

	log.Printf("Loaded configuration for APP_ENV: %s", os.Getenv("APP_ENV"))
	log.Printf("Database URL: %s", cfg.DatabaseURL)
	log.Printf("Server Port: %s", cfg.Port)
	log.Printf("Google OAuth Client ID: %s", cfg.GoogleOAuthClientID)

	// 2. Establish database connection and run migrations
	dbpool := connectToDB(cfg.DatabaseURL)
	defer dbpool.Close()

	// 3. Set up dependencies (repositories, services)
	userRepo := postgres.NewUserRepository(dbpool)
	organizationRepo := postgres.NewOrganizationRepository(dbpool)
	organizationUserRepo := postgres.NewOrganizationUserRepository(dbpool)

	accessService := services.NewAccessService(organizationUserRepo)
	userService := services.NewUserService(userRepo)
	organizationService := services.NewOrganizationService(organizationRepo)
	organizationUserService := services.NewOrganizationUserService(organizationUserRepo, accessService)

	tokenVerifier := &auth.GoogleTokenVerifier{}

	// 4. Set up the HTTP server
	server := api.NewServer(cfg, tokenVerifier, userService, organizationService, organizationUserService, accessService)

	// 5. Start the server using the port from the config
	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Starting server on %s", addr)
	if err := http.ListenAndServe(addr, server); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}

// connectToDB establishes a connection to the database, runs migrations,
// and returns a connection pool. It will exit the application on any error.
func connectToDB(databaseURL string) *pgxpool.Pool {
	// Run migrations first to ensure the database schema is ready.
	runMigrations(databaseURL)

	// Create a new connection pool from the URL.
	dbpool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v", err)
	}

	// Ping the database to verify that a connection has been established.
	if err := dbpool.Ping(context.Background()); err != nil {
		log.Fatalf("Unable to connect to the database: %v", err)
	}

	log.Println("Connected to the database successfully")
	return dbpool
}

// runMigrations applies all pending database migrations.
func runMigrations(databaseURL string) {
	log.Println("Running database migrations...")
	m, err := migrate.New("file://migrations", databaseURL)
	if err != nil {
		log.Fatalf("Failed to create migrator: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Failed to apply migrations: %v", err)
	}

	log.Println("Migrations applied successfully")
}
