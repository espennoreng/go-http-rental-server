package services

import (
	"context"
	"fmt"

	"github.com/espennoreng/go-http-rental-server/internal/models"
	"github.com/espennoreng/go-http-rental-server/internal/repositories"
	"github.com/google/uuid"
)

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) *userService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) CreateUser(ctx context.Context, input models.CreateUserInput) (*models.User, error) {
	if input.Username == "" {
		return nil, fmt.Errorf("username cannot be empty")
	}
	if input.Email == "" {
		return nil, fmt.Errorf("email cannot be empty")
	}

	newUser := models.User{
		ID:       uuid.New().String(),
		Username: input.Username,
		Email:    input.Email,
	}

	if err := s.userRepo.Create(ctx, &newUser); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &newUser, nil
}

func (s *userService) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	if id == "" {
		return nil, fmt.Errorf("user ID cannot be empty")
	}

	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return user, nil
}