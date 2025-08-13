package api

import (
	"context"
	"net/http"

	"github.com/espennoreng/go-http-rental-server/internal/middleware"
)

func TestAuthMiddleware(next http.Handler, userID string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), middleware.UserCtxKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
