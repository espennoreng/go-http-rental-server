package services_test

import (
	"context"
	"testing"

	"github.com/espennoreng/go-http-rental-server/internal/models"
	"github.com/espennoreng/go-http-rental-server/internal/repositories"
	"github.com/espennoreng/go-http-rental-server/internal/services"
)

type mockUserRepository struct {
	createFunc  func(ctx context.Context, user *repositories.CreateUserParams) (*models.User, error)
	getByIDFunc func(ctx context.Context, id string) (*models.User, error)
}

func (m *mockUserRepository) Create(ctx context.Context, user *repositories.CreateUserParams) (*models.User, error) {
	return m.createFunc(ctx, user)
}

func (m *mockUserRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	return m.getByIDFunc(ctx, id)
}

func TestUserService_CreateUser(t *testing.T) {
	t.Run("create user successfully", func(t *testing.T) {
		repo := &mockUserRepository{
			createFunc: func(ctx context.Context, user *repositories.CreateUserParams) (*models.User, error) {
				return &models.User{ID: "user-001", Username: user.Username, Email: user.Email}, nil
			},
			getByIDFunc: func(ctx context.Context, id string) (*models.User, error) {
				return &models.User{ID: id, Username: "John Doe", Email: "john.doe@example.com"}, nil
			},
		}

		service := services.NewUserService(repo)

		createdUser, err := service.CreateUser(context.Background(), repositories.CreateUserParams{Username: "John Doe", Email: "john.doe@example.com"})
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

		_, err := service.CreateUser(context.Background(), repositories.CreateUserParams{Username: "", Email: "john.doe@example.com"})
		if err == nil {
			t.Fatal("expected error, got none")
		}
	})

	t.Run("create user with empty email", func(t *testing.T) {
		repo := &mockUserRepository{}

		service := services.NewUserService(repo)

		_, err := service.CreateUser(context.Background(), repositories.CreateUserParams{Username: "John Doe", Email: ""})
		if err == nil {
			t.Fatal("expected error, got none")
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

	t.Run("get user by empty ID", func(t *testing.T) {
		repo := &mockUserRepository{}

		service := services.NewUserService(repo)

		_, err := service.GetUserByID(context.Background(), "")
		if err == nil {
			t.Fatal("expected error, got none")
		}
	})
}
