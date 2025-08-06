package services

import (
	"context"
	"errors"

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
		return nil, ErrInvalidInput
	}
	if input.Email == "" {
		return nil, ErrInvalidInput
	}

	newUser := models.User{
		ID:       uuid.New().String(),
		Username: input.Username,
		Email:    input.Email,
	}

	if err := s.userRepo.Create(ctx, &newUser); err != nil {
		if errors.Is(err, repositories.ErrUniqueConstraint) {
			return nil, ErrUserWithDuplicateDetailsExists
		}
		return nil, ErrInternalServer
	}

	return &newUser, nil
}

func (s *userService) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	if id == "" {
		return nil, ErrInvalidInput
	}

	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repositories.ErrNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, ErrInternalServer
	}

	return user, nil
}