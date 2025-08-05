package api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/espennoreng/go-http-rental-server/internal/api"
	"github.com/espennoreng/go-http-rental-server/internal/models"
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
	t.Run("create user successfully", func(t *testing.T) {
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

		var createdUser models.User
		if err := json.NewDecoder(res.Body).Decode(&createdUser); err != nil {
			t.Fatalf("failed to decode response body: %v", err)
		}

		if createdUser.Username != "John Doe" {
			t.Errorf("expected username 'John Doe', got '%s'", createdUser.Username)
		}

		if createdUser.Email != "john.doe@example.com" {
			t.Errorf("expected email 'john.doe@example.com', got '%s'", createdUser.Email)
		}
	})
}

func TestUserHandler_GetUserByID(t *testing.T) {
	t.Run("get user by ID successfully", func(t *testing.T) {
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

		var user models.User
		if err := json.NewDecoder(res.Body).Decode(&user); err != nil {
			t.Fatalf("failed to decode response body: %v", err)
		}

		if user.ID != "user-001" {
			t.Errorf("expected user ID 'user-001', got '%s'", user.ID)
		}
	})
}
