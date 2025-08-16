package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

// NewSlogMiddleware returns a middleware that logs requests using the provided slog.Logger.
func NewSlogMiddleware(logger *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get the request ID from the context (thanks to middleware.RequestID)
			reqID := middleware.GetReqID(r.Context())

			// Create a logger with request-specific fields
			requestLogger := logger.With(
				slog.String("request_id", reqID),
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("remote_addr", r.RemoteAddr),
				slog.String("user_agent", r.UserAgent()),
			)

			requestLogger.Info("Request started")

			// Use chi's WrapResponseWriter to capture status code and bytes written
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			start := time.Now()

			defer func() {
				requestLogger.Info("Request finished",
					slog.Int("status", ww.Status()),
					slog.Int("bytes_written", ww.BytesWritten()),
					slog.Duration("latency", time.Since(start)),
				)
			}()

			next.ServeHTTP(ww, r)
		})
	}
}
