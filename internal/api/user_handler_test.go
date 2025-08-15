package api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/espennoreng/go-http-rental-server/internal/api"
	"github.com/espennoreng/go-http-rental-server/internal/auth"
	"github.com/espennoreng/go-http-rental-server/internal/middleware"
	"github.com/espennoreng/go-http-rental-server/internal/models"
	"github.com/espennoreng/go-http-rental-server/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

type mockUserService struct {
	createUserFunc  func(ctx context.Context, params services.CreateUserParams) (*models.User, error)
	getUserByIDFunc func(ctx context.Context, params services.GetUserByIDParams) (*models.User, error)
}

func (m *mockUserService) CreateUser(ctx context.Context, params services.CreateUserParams) (*models.User, error) {
	return m.createUserFunc(ctx, params)
}

func (m *mockUserService) GetUserByID(ctx context.Context, params services.GetUserByIDParams) (*models.User, error) {
	return m.getUserByIDFunc(ctx, params)
}

func TestUserHandler_CreateUser(t *testing.T) {
	t.Run("successful user creation", func(t *testing.T) {
		mockService := &mockUserService{
			createUserFunc: func(ctx context.Context, params services.CreateUserParams) (*models.User, error) {
				return &models.User{
					ID:       "user-001",
					Username: params.Username,
					Email:    params.Email,
				}, nil
			},
		}

		r := chi.NewRouter()
		handler := api.NewUserHandler(mockService)
		r.Post("/users", handler.CreateUser)

		reqBody := `{"username": "John Doe", "email": "john.doe@example.com"}`
		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(reqBody))
		res := httptest.NewRecorder()

		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusCreated, res.Code)
		api.AssertJSONContentType(t, res)

		var user models.User
		err := json.NewDecoder(res.Body).Decode(&user)
		assert.NoError(t, err)
		assert.NotEmpty(t, user.ID)
		assert.Equal(t, "user-001", user.ID)
		assert.Equal(t, "John Doe", user.Username)
		assert.Equal(t, "john.doe@example.com", user.Email)

	})

	t.Run("invalid request body", func(t *testing.T) {
		mockService := &mockUserService{}

		r := chi.NewRouter()
		handler := api.NewUserHandler(mockService)
		r.Post("/users", handler.CreateUser)

		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString("invalid json"))
		res := httptest.NewRecorder()

		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)
		api.AssertJSONContentType(t, res)
	})

	t.Run("service returns validation error", func(t *testing.T) {
		mockService := &mockUserService{
			createUserFunc: func(ctx context.Context, params services.CreateUserParams) (*models.User, error) {
				return nil, services.ErrInvalidInput
			},
		}

		r := chi.NewRouter()
		handler := api.NewUserHandler(mockService)
		r.Post("/users", handler.CreateUser)

		reqBody := `{"email": "john.doe@example.com"}`
		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(reqBody))
		res := httptest.NewRecorder()

		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)
		api.AssertJSONContentType(t, res)
	})

	t.Run("service returns internal server error", func(t *testing.T) {
		mockService := &mockUserService{
			createUserFunc: func(ctx context.Context, params services.CreateUserParams) (*models.User, error) {
				return nil, services.ErrInternalServer
			},
		}

		r := chi.NewRouter()
		handler := api.NewUserHandler(mockService)
		r.Post("/users", handler.CreateUser)

		reqBody := `{"username": "John Doe", "email": "john.doe@example.com"}`
		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(reqBody))
		res := httptest.NewRecorder()

		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)
		api.AssertJSONContentType(t, res)
	})

	t.Run("service returns user with duplicate details error", func(t *testing.T) {
		mockService := &mockUserService{
			createUserFunc: func(ctx context.Context, params services.CreateUserParams) (*models.User, error) {
				return nil, services.ErrUserWithDuplicateDetailsExists
			},
		}

		r := chi.NewRouter()
		handler := api.NewUserHandler(mockService)
		r.Post("/users", handler.CreateUser)

		reqBody := `{"username": "John Doe", "email": "john.doe@example.com"}`
		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(reqBody))
		res := httptest.NewRecorder()

		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusConflict, res.Code)
		api.AssertJSONContentType(t, res)
	})
}

func TestUserHandler_GetUserByID(t *testing.T) {
	actingUser := auth.Identity{UserID: "user-001"}
	t.Run("successful user retrieval", func(t *testing.T) {

		mockService := &mockUserService{
			getUserByIDFunc: func(ctx context.Context, params services.GetUserByIDParams) (*models.User, error) {
				if params.ID == "user-001" {
					return &models.User{
						ID:       "user-001",
						Username: "John Doe",
						Email:    "john.doe@example.com",
					}, nil
				}
				return nil, nil
			},
		}

		r := chi.NewRouter()
		handler := api.NewUserHandler(mockService)
		authedHandler := middleware.TestAuthMiddleware(http.HandlerFunc(handler.GetUserByID), actingUser)

		r.Method(http.MethodGet, "/users/{id}", authedHandler)

		req := httptest.NewRequest(http.MethodGet, "/users/user-001", nil)
		res := httptest.NewRecorder()

		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
		api.AssertJSONContentType(t, res)
	})

	t.Run("service returns not found error", func(t *testing.T) {
		mockService := &mockUserService{
			getUserByIDFunc: func(ctx context.Context, params services.GetUserByIDParams) (*models.User, error) {
				return nil, services.ErrUserNotFound
			},
		}

		r := chi.NewRouter()
		handler := api.NewUserHandler(mockService)

		authedHandler := middleware.TestAuthMiddleware(http.HandlerFunc(handler.GetUserByID), actingUser)

		r.Method(http.MethodGet, "/users/{id}", authedHandler)

		req := httptest.NewRequest(http.MethodGet, "/users/non-existent-id", nil)
		res := httptest.NewRecorder()

		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusNotFound, res.Code)
		api.AssertJSONContentType(t, res)
	})

	t.Run("service returns internal server error", func(t *testing.T) {
		mockService := &mockUserService{
			getUserByIDFunc: func(ctx context.Context, params services.GetUserByIDParams) (*models.User, error) {
				return nil, services.ErrInternalServer
			},
		}

		r := chi.NewRouter()
		handler := api.NewUserHandler(mockService)
		authedHandler := middleware.TestAuthMiddleware(http.HandlerFunc(handler.GetUserByID), actingUser)
		r.Method(http.MethodGet, "/users/{id}", authedHandler)

		req := httptest.NewRequest(http.MethodGet, "/users/user-001", nil)
		res := httptest.NewRecorder()

		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)
		api.AssertJSONContentType(t, res)
	})

	t.Run("unauthorized access", func(t *testing.T) {
		mockService := &mockUserService{}

		r := chi.NewRouter()
		handler := api.NewUserHandler(mockService)

		r.Get("/users/{id}", handler.GetUserByID)

		req := httptest.NewRequest(http.MethodGet, "/users/user-001", nil)
		res := httptest.NewRecorder()

		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusUnauthorized, res.Code)
	})
}
