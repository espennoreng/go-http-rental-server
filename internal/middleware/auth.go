package middleware

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/espennoreng/go-http-rental-server/internal/auth"
	"github.com/espennoreng/go-http-rental-server/internal/services"
)

// NewAuthMiddleware creates a new authentication middleware.
func NewAuthMiddleware(log *slog.Logger, verifier auth.TokenVerifier, userService services.UserService, audience string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 1. Get token from header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				log.Error("Authorization header is missing")
				http.Error(w, "Authorization header required", http.StatusUnauthorized)
				return
			}
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				log.Error("Invalid Authorization header format")
				http.Error(w, "Authorization header must be 'Bearer {token}'", http.StatusUnauthorized)
				return
			}
			tokenString := parts[1]

			// 2. Verify the token with Google
			payload, err := verifier.Verify(r.Context(), tokenString, audience)
			if err != nil {
				log.Error("Token verification failed", slog.Any("error", err))
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			// 3. Find or create a user in your database
			googleID := payload.Subject
			email := payload.Claims["email"].(string)

			user, err := userService.FindOrCreateByGoogleID(r.Context(), googleID, email)
			if err != nil {
				log.Error("Failed to find or create user", slog.Any("error", err))
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			// 4. Create identity and inject it into the context
			identity := auth.Identity{UserID: user.ID} // Use your internal user ID
			ctx := auth.ToContext(r.Context(), identity)

			log.Info("Authenticated user", slog.String("user_id", user.ID), slog.String("email", email))

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// TestAuthMiddleware is a helper for setting up authenticated tests.
func NewTestAuthMiddleware(next http.Handler, identity auth.Identity) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := auth.ToContext(r.Context(), identity)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
