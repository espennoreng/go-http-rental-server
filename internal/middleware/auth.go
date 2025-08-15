package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/espennoreng/go-http-rental-server/internal/auth"
	"github.com/espennoreng/go-http-rental-server/internal/services"
)

// AuthMiddleware authenticates a request and injects the user's Identity into the context.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// In a real app, you would parse a JWT or session token here.
		identity := auth.Identity{
			UserID: "a53e4b0c-9d6c-4f7f-8c3b-5a1e2f3g4h5i",
		}

		// Use the canonical function from the auth package to add to the context.
		ctx := auth.ToContext(r.Context(), identity)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// NewAuthMiddleware creates a new authentication middleware.
func NewAuthMiddleware(verifier auth.TokenVerifier, userService services.UserService, audience string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 1. Get token from header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header required", http.StatusUnauthorized)
				return
			}
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				http.Error(w, "Authorization header must be 'Bearer {token}'", http.StatusUnauthorized)
				return
			}
			tokenString := parts[1]

			// 2. Verify the token with Google
			payload, err := verifier.Verify(r.Context(), tokenString, audience)
			if err != nil {
				log.Printf("Token verification failed: %v", err)
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			// 3. Find or create a user in your database
			googleID := payload.Subject
			email := payload.Claims["email"].(string)

			user, err := userService.FindOrCreateByGoogleID(r.Context(), googleID, email)
			if err != nil {
				log.Printf("Failed to find or create user: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			// 4. Create identity and inject it into the context
			identity := auth.Identity{UserID: user.ID} // Use your internal user ID
			ctx := auth.ToContext(r.Context(), identity)

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
