package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/espennoreng/go-http-rental-server/internal/auth"
	"github.com/espennoreng/go-http-rental-server/internal/logger"
	"github.com/espennoreng/go-http-rental-server/internal/middleware"
	"github.com/espennoreng/go-http-rental-server/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

type mockAccessService struct {
	isAdminFunc  func(ctx context.Context, params services.OrgAccessParams) error
	isMemberFunc func(ctx context.Context, params services.OrgAccessParams) error
}

func (m *mockAccessService) IsAdmin(ctx context.Context, params services.OrgAccessParams) error {
	return m.isAdminFunc(ctx, params)
}

func (m *mockAccessService) IsMember(ctx context.Context, params services.OrgAccessParams) error {
	return m.isMemberFunc(ctx, params)
}

func TestAccessMiddleware_RequireAdmin_AllowsAdminUser(t *testing.T) {

	finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	mockSvc := &mockAccessService{
		isAdminFunc: func(ctx context.Context, params services.OrgAccessParams) error {
			return nil
		},
	}

	accessMiddleware := middleware.NewAccessMiddleware(mockSvc, logger.NewTestLogger(t))
	handlerChain := middleware.NewTestAuthMiddleware(accessMiddleware.RequireAdmin(finalHandler), auth.Identity{UserID: "admin-user"})

	router := chi.NewRouter()
	router.Method(http.MethodGet, "/orgs/{orgID}", handlerChain)

	req := httptest.NewRequest("GET", "/orgs/org-123", nil)
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		assert.Fail(t, "expected status OK (200); got %d", res.Code)
	}
}

// in middleware/access_middleware_test.go

func TestAccessMiddleware_RequireAdmin_BlocksNonAdminUser(t *testing.T) {
	// --- Arrange ---
	// A dummy handler that should NOT be called
	finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("final handler should not be called for a non-admin user")
	})

	// Configure the mock to return "false" for IsAdmin
	mockSvc := &mockAccessService{
		isAdminFunc: func(ctx context.Context, params services.OrgAccessParams) error {
			return services.ErrUnauthorized
		},
	}

	accessMiddleware := middleware.NewAccessMiddleware(mockSvc, logger.NewTestLogger(t))
	handlerChain := middleware.NewTestAuthMiddleware(accessMiddleware.RequireAdmin(finalHandler), auth.Identity{UserID: "non-admin-user"})

	router := chi.NewRouter()
	router.Method(http.MethodGet, "/orgs/{orgID}", handlerChain)
	req := httptest.NewRequest("GET", "/orgs/org-123", nil)
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Code != http.StatusForbidden {
		assert.Fail(t, "expected status Forbidden (403); got %d", res.Code)
	}

	assert.Contains(t, res.Body.String(), "Forbidden")
}
