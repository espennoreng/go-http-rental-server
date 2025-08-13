package middleware

import (
	"net/http"

	"github.com/espennoreng/go-http-rental-server/internal/auth"
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

// TestAuthMiddleware is a helper for setting up authenticated tests.
func TestAuthMiddleware(next http.Handler, identity auth.Identity) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := auth.ToContext(r.Context(), identity)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}