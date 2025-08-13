package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/espennoreng/go-http-rental-server/internal/auth"
	"github.com/espennoreng/go-http-rental-server/internal/services"
	"github.com/go-chi/chi/v5"
)

// AccessMiddleware holds the dependencies for our authorization middleware.
type AccessMiddleware struct {
	accessService services.AccessService
}

func NewAccessMiddleware(accessService services.AccessService) *AccessMiddleware {
	return &AccessMiddleware{
		accessService: accessService,
	}
}

// accessCheckFunc defines a function signature for any organization-level permission check.
type accessCheckFunc func(context.Context, services.OrgAccessParams) error

func (am *AccessMiddleware) requireAccess(next http.Handler, check accessCheckFunc, forbiddenMsg string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		identity, err := auth.FromContext(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		orgID := chi.URLParam(r, "orgID")
		if orgID == "" {
			http.Error(w, "Bad Request: Missing organization ID", http.StatusBadRequest)
			return
		}

		err = check(r.Context(), services.OrgAccessParams{
			OrgID:  orgID,
			UserID: identity.UserID,
		})
		if err != nil {
			if errors.Is(err, services.ErrUnauthorized) {
				http.Error(w, forbiddenMsg, http.StatusForbidden)
				return
			}
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (am *AccessMiddleware) RequireAdmin(next http.Handler) http.Handler {
	return am.requireAccess(
		next,
		am.accessService.IsAdmin,
		"Forbidden: You are not an admin of this organization",
	)
}

func (am *AccessMiddleware) RequireMember(next http.Handler) http.Handler {
	return am.requireAccess(
		next,
		am.accessService.IsMember,
		"Forbidden: You are not a member of this organization",
	)
}
