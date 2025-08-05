package api_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/espennoreng/go-http-rental-server/internal/api"
	"github.com/espennoreng/go-http-rental-server/internal/models"
	"github.com/espennoreng/go-http-rental-server/internal/services"
	"github.com/go-chi/chi/v5"
)

type mockUserService struct {
	createUserFunc func(ctx context.Context, input models.CreateUserInput) (*models.User, error)
	getUserByIDFunc func(ctx context.Context, id string) (*models.User, error)
}

func (m *mockUserService) CreateUser(ctx context.Context, input models.CreateUserInput) (*models.User, error) {
	return m.createUserFunc(ctx, input)
}

func (m *mockUserService) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	return m.getUserByIDFunc(ctx, id)
}

func TestUserHandler_CreateUser(t *testing.T) {
	t.Run("successful user creation", func(t *testing.T) {
		mockService := &mockUserService{
			createUserFunc: func(ctx context.Context, input models.CreateUserInput) (*models.User, error) {
				return &models.User{
					ID:       "user-001",
					Username: input.Username,
					Email:    input.Email,
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

		if res.Code != http.StatusCreated {
			t.Errorf("expected status %d, got %d", http.StatusCreated, res.Code)
		}
	})

	t.Run("invalid request body", func(t *testing.T) {
		mockService := &mockUserService{}

		r := chi.NewRouter()
		handler := api.NewUserHandler(mockService)
		r.Post("/users", handler.CreateUser)

		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString("invalid json"))
		res := httptest.NewRecorder()

		r.ServeHTTP(res, req)

		if res.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, res.Code)
		}
	})

	t.Run("service returns validation error", func(t *testing.T) {
		mockService := &mockUserService{
			createUserFunc: func(ctx context.Context, input models.CreateUserInput) (*models.User, error) {
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

		if res.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, res.Code)
		}
	})

	t.Run("service returns internal server error", func(t *testing.T) {
		mockService := &mockUserService{
			createUserFunc: func(ctx context.Context, input models.CreateUserInput) (*models.User, error) {
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

		if res.Code != http.StatusInternalServerError {
			t.Errorf("expected status %d, got %d", http.StatusInternalServerError, res.Code)
		}
	})

	t.Run("service returns user already exists error", func(t *testing.T) {
		mockService := &mockUserService{
			createUserFunc: func(ctx context.Context, input models.CreateUserInput) (*models.User, error) {
				return nil, services.ErrUserAlreadyExists
			},
		}

		r := chi.NewRouter()
		handler := api.NewUserHandler(mockService)
		r.Post("/users", handler.CreateUser)

		reqBody := `{"username": "John Doe", "email": "john.doe@example.com"}`
		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(reqBody))
		res := httptest.NewRecorder()

		r.ServeHTTP(res, req)

		if res.Code != http.StatusConflict {
			t.Errorf("expected status %d, got %d", http.StatusConflict, res.Code)
		}
	})
}

func TestUserHandler_GetUserByID(t *testing.T) {
	t.Run("successful user retrieval", func(t *testing.T) {
		mockService := &mockUserService{
			getUserByIDFunc: func(ctx context.Context, id string) (*models.User, error) {
				if id == "user-001" {
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
		r.Get("/users/{id}", handler.GetUserByID)


		req := httptest.NewRequest(http.MethodGet, "/users/user-001", nil)
		res := httptest.NewRecorder()

		r.ServeHTTP(res, req)

		if res.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, res.Code)
		}
	})

	t.Run("service returns not found error", func(t *testing.T) {
		mockService := &mockUserService{
			getUserByIDFunc: func(ctx context.Context, id string) (*models.User, error) {
				return nil, services.ErrUserNotFound
			},
		}

		r := chi.NewRouter()
		handler := api.NewUserHandler(mockService)
		r.Get("/users/{id}", handler.GetUserByID)

		req := httptest.NewRequest(http.MethodGet, "/users/non-existent-id", nil)
		res := httptest.NewRecorder()

		r.ServeHTTP(res, req)

		if res.Code != http.StatusNotFound {
			t.Errorf("expected status %d, got %d", http.StatusNotFound, res.Code)
		}
	})

	t.Run("service returns internal server error", func(t *testing.T) {
		mockService := &mockUserService{
			getUserByIDFunc: func(ctx context.Context, id string) (*models.User, error) {
				return nil, services.ErrInternalServer
			},
		}

		r := chi.NewRouter()
		handler := api.NewUserHandler(mockService)
		r.Get("/users/{id}", handler.GetUserByID)

		req := httptest.NewRequest(http.MethodGet, "/users/user-001", nil)
		res := httptest.NewRecorder()

		r.ServeHTTP(res, req)

		if res.Code != http.StatusInternalServerError {
			t.Errorf("expected status %d, got %d", http.StatusInternalServerError, res.Code)
		}
	})

}
