package services_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/espennoreng/go-http-rental-server/internal/models"
	"github.com/espennoreng/go-http-rental-server/internal/services"
)

type mockUserRepository struct {
	createFunc  func(ctx context.Context, user *models.User) error
	getByIDFunc func(ctx context.Context, id string) (*models.User, error)
}

func (m *mockUserRepository) Create(ctx context.Context, user *models.User) error {
	return m.createFunc(ctx, user)
}

func (m *mockUserRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	return m.getByIDFunc(ctx, id)
}

func TestUserService_CreateUser(t *testing.T) {
	t.Run("create user successfully", func(t *testing.T) {
		repo := &mockUserRepository{
			createFunc: func(ctx context.Context, user *models.User) error {
				return nil // Simulate successful creation
			},
			getByIDFunc: func(ctx context.Context, id string) (*models.User, error) {
				return &models.User{ID: id, Username: "John Doe", Email: "john.doe@example.com"}, nil
			},
		}

		service := services.NewUserService(repo)

		user := models.CreateUserInput{Username: "John Doe", Email: "john.doe@example.com"}
		createdUser, err := service.CreateUser(context.Background(), user)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if createdUser.ID == "" {
			t.Fatal("expected user ID to be set")
		}
	})

	t.Run("create user with empty username", func(t *testing.T) {
		repo := &mockUserRepository{}

		service := services.NewUserService(repo)

		_, err := service.CreateUser(context.Background(), models.CreateUserInput{Username: "", Email: "john.doe@example.com"})
		if err == nil {
			t.Fatal("expected error, got none")
		}
	})

	t.Run("create user with empty email", func(t *testing.T) {
		repo := &mockUserRepository{}

		service := services.NewUserService(repo)

		_, err := service.CreateUser(context.Background(), models.CreateUserInput{Username: "John Doe", Email: ""})
		if err == nil {
			t.Fatal("expected error, got none")
		}
	})

	t.Run("create user repository error", func(t *testing.T) {
		repo := &mockUserRepository{
			createFunc: func(ctx context.Context, user *models.User) error {
				return fmt.Errorf("repository error")
			},
		}

		service := services.NewUserService(repo)

		_, err := service.CreateUser(context.Background(), models.CreateUserInput{Username: "John Doe", Email: "john.doe@example.com"})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

}

func TestUserService_GetUserByID(t *testing.T) {
	t.Run("get user by ID successfully", func(t *testing.T) {
		repo := &mockUserRepository{
			getByIDFunc: func(ctx context.Context, id string) (*models.User, error) {
				return &models.User{ID: id, Username: "John Doe", Email: "john.doe@example.com"}, nil
			},
		}

		service := services.NewUserService(repo)

		user, err := service.GetUserByID(context.Background(), "user-001")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if user.ID == "" {
			t.Fatal("expected user ID to be set")
		}
	})
}

func TestUserService_GetUserByID_NotFound(t *testing.T) {
	t.Run("get user by ID not found", func(t *testing.T) {
		repo := &mockUserRepository{
			getByIDFunc: func(ctx context.Context, id string) (*models.User, error) {
				return nil, nil // Simulate user not found
			},
		}

		service := services.NewUserService(repo)

		user, err := service.GetUserByID(context.Background(), "non-existent-id")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if user != nil {
			t.Fatal("expected user to be nil for non-existent ID")
		}
	})
}