package middleware

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/espennoreng/go-http-rental-server/internal/logger"
	"github.com/espennoreng/go-http-rental-server/internal/models"
	"github.com/espennoreng/go-http-rental-server/internal/services"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"google.golang.org/api/idtoken"
)

// mockTokenVerifier lets us control the outcome of token verification.
type mockTokenVerifier struct {
	payload *idtoken.Payload
	err     error
}
func (m *mockTokenVerifier) Verify(ctx context.Context, _, _ string) (*idtoken.Payload, error) {
	return m.payload, m.err
}

type mockUserService struct {
	findOrCreateByGoogleIDFunc func(ctx context.Context, googleID, email string) (*models.User, error)
}

func (m *mockUserService) FindOrCreateByGoogleID(ctx context.Context, googleID, email string) (*models.User, error) {
	if m.findOrCreateByGoogleIDFunc != nil {
		return m.findOrCreateByGoogleIDFunc(ctx, googleID, email)
	}
	return nil, errors.New("not implemented")
}

func (m *mockUserService) GetUserByID(ctx context.Context, params services.GetUserByIDParams) (*models.User, error) {
	return nil, errors.New("not implemented")
}

func (m *mockUserService) CreateUser(ctx context.Context, params services.CreateUserParams) (*models.User, error) {
	return nil, errors.New("not implemented")
}

// TestAuthMiddleware contains all test cases for the authentication middleware.
func TestAuthMiddleware(t *testing.T) {
	// --- Common Test Data ---
	sampleUserID := uuid.New().String()
	validPayload := &idtoken.Payload{
		Subject: "123456789", // Google User ID
		Claims: map[string]interface{}{
			"email": "test@example.com",
		},
	}
	// A simple handler that will be protected by the middleware.
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// --- Test Cases ---

	t.Run("should succeed with a valid token and existing user", func(t *testing.T) {
		// Arrange
		verifier := &mockTokenVerifier{payload: validPayload}
		userService := &mockUserService{
			findOrCreateByGoogleIDFunc: func(ctx context.Context, googleID, email string) (*models.User, error) {
				return &models.User{ID: sampleUserID}, nil
			},
		}

		req := httptest.NewRequest("GET", "/private", nil)
		req.Header.Set("Authorization", "Bearer valid-token")
		rr := httptest.NewRecorder()

		middleware := NewAuthMiddleware(logger.NewTestLogger(t), verifier, userService, "test-audience")
		handler := middleware(testHandler)

		// Act
		handler.ServeHTTP(rr, req)

		// Assert
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "OK", rr.Body.String())
	})

	t.Run("should fail when token verification fails", func(t *testing.T) {
		// Arrange
		verifier := &mockTokenVerifier{err: errors.New("invalid signature")}
		userService := &mockUserService{}

		req := httptest.NewRequest("GET", "/private", nil)
		req.Header.Set("Authorization", "Bearer invalid-token")
		rr := httptest.NewRecorder()

		middleware := NewAuthMiddleware(logger.NewTestLogger(t), verifier, userService, "test-audience")
		handler := middleware(testHandler)

		// Act
		handler.ServeHTTP(rr, req)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		assert.Contains(t, rr.Body.String(), "Invalid token")
	})

	t.Run("should fail when authorization header is missing", func(t *testing.T) {
		// Arrange
		req := httptest.NewRequest("GET", "/private", nil)
		rr := httptest.NewRecorder()

		middleware := NewAuthMiddleware(logger.NewTestLogger(t), &mockTokenVerifier{}, &mockUserService{}, "test-audience")
		handler := middleware(testHandler)

		// Act
		handler.ServeHTTP(rr, req)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})
	
	t.Run("should fail when user service returns an error", func(t *testing.T) {
		// Arrange
		verifier := &mockTokenVerifier{payload: validPayload}
		userService := &mockUserService{
			findOrCreateByGoogleIDFunc: func(ctx context.Context, googleID, email string) (*models.User, error) {
				return nil, errors.New("user not found")
			},
		}

		req := httptest.NewRequest("GET", "/private", nil)
		req.Header.Set("Authorization", "Bearer valid-token")
		rr := httptest.NewRecorder()

		middleware := NewAuthMiddleware(logger.NewTestLogger(t), verifier, userService, "test-audience")
		handler := middleware(testHandler)

		// Act
		handler.ServeHTTP(rr, req)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}