package middleware

import (
	"context"
	"errors"
	"net/http"
)

type UserCtxKeyType string

const UserCtxKey = UserCtxKeyType("userID")

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: Use a real user ID from the request context or a token
		userID := "a53e4b0c-9d6c-4f7f-8c3b-5a1e2f3g4h5i"

		ctx := context.WithValue(r.Context(), UserCtxKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserIDFromContext is a helper function to safely retrieve the userID from the context.
func GetUserIDFromContext(ctx context.Context) (string, error) {
	userID, ok := ctx.Value(UserCtxKey).(string)
	if !ok || userID == "" {
		return "", errors.New("UNAUTHORIZED: user ID not found in context")
	}
	return userID, nil
}
