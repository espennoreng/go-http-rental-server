package logger

import (
	"bytes"
	"log/slog"
	"os"
	"testing"

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

// testWriter adapts *testing.T to the io.Writer interface.
type testWriter struct {
    t *testing.T
}

// Write logs the provided bytes to the test's log.
func (tw *testWriter) Write(p []byte) (n int, err error) {
    // Trim space is useful to remove the trailing newline slog often adds.
    tw.t.Log(string(bytes.TrimSpace(p)))
    return len(p), nil
}

// NewTestLogger creates a slog.Logger that writes to the test's log buffer.
func NewTestLogger(t *testing.T) *slog.Logger {
    return slog.New(slog.NewTextHandler(&testWriter{t: t}, &slog.HandlerOptions{
        // Use a low level to ensure all logs are captured.
        Level: slog.LevelDebug, 
    }))
}