package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/espennoreng/go-http-rental-server/internal/api"
	"github.com/espennoreng/go-http-rental-server/internal/auth"
	"github.com/espennoreng/go-http-rental-server/internal/config"
	"github.com/espennoreng/go-http-rental-server/internal/logger"
	"github.com/espennoreng/go-http-rental-server/internal/repositories/postgres"
	"github.com/espennoreng/go-http-rental-server/internal/services"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {

	log := logger.New(config.Env(os.Getenv("APP_ENV")))
	// 1. Load application configuration
	cfg, err := config.Load()
	if err != nil {
		log.Error("Could not load config: %v", slog.Any("error", err))
	}

	log.Info("Configuration loaded successfully",
		slog.String("app_env", os.Getenv("APP_ENV")),
		slog.String("port", cfg.Port),
		slog.Bool("database_url_set", cfg.DatabaseURL != ""),
		slog.Bool("google_client_id_set", cfg.GoogleOAuthClientID != ""),
	)

	// 2. Establish database connection and run migration
	dbpool := connectToDB(cfg.DatabaseURL)
	defer dbpool.Close()

	// 3. Set up dependencies (repositories, services)
	userRepo := postgres.NewUserRepository(dbpool, log)
	organizationRepo := postgres.NewOrganizationRepository(dbpool, log)
	organizationUserRepo := postgres.NewOrganizationUserRepository(dbpool, log)

	accessService := services.NewAccessService(organizationUserRepo, log)
	organizationUserService := services.NewOrganizationUserService(organizationUserRepo, accessService)
	userService := services.NewUserService(userRepo, organizationUserRepo, log)
	organizationService := services.NewOrganizationService(organizationRepo, log)

	tokenVerifier := &auth.GoogleTokenVerifier{}

	// 4. Set up the HTTP server
	server := api.NewServer(cfg, tokenVerifier, log, userService, organizationService, organizationUserService, accessService)

	// 5. Start the server using the port from the config
	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Info("Starting server", slog.String("address", addr))
	if err := http.ListenAndServe(addr, server); err != nil {
		log.Error("Could not start server", slog.Any("error", err))
		os.Exit(1)
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
		slog.Error("Unable to create connection pool", slog.Any("error", err))
		os.Exit(1)
	}

	// Ping the database to verify that a connection has been established.
	if err := dbpool.Ping(context.Background()); err != nil {
		slog.Error("Unable to connect to the database", slog.Any("error", err))
		os.Exit(1)
	}

	slog.Info("Connected to the database successfully")
	return dbpool
}

// runMigrations applies all pending database migrations.
func runMigrations(databaseURL string) {
	log.Println("Running database migrations...")
	m, err := migrate.New("file://migrations", databaseURL)
	if err != nil {
		slog.Error("Failed to create migration instance", slog.Any("error", err))
		os.Exit(1)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		slog.Error("Failed to apply migrations", slog.Any("error", err))
		os.Exit(1)
	}

	slog.Info("Migrations applied successfully")
}
