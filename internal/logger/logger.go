package logger

import (
	"log/slog"
	"os"

	"github.com/espennoreng/go-http-rental-server/internal/config"
)

func New(env config.Env) *slog.Logger {
	var handler slog.Handler

	switch env {
	case config.Development:
		// Use a more readable handler for development.
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug, // Log all levels in development.
		})
	default:
		// Use JSON for production, which is better for log collectors.
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	}

	logger := slog.New(handler)

	slog.SetDefault(logger)

	return logger
}
