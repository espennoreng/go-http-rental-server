package services

import (
	"context"

	"github.com/espennoreng/go-http-rental-server/internal/models"
	"github.com/espennoreng/go-http-rental-server/internal/repositories"
)

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) *userService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) CreateUser(ctx context.Context, input repositories.CreateUserParams) (*models.User, error) {
	if input.Username == "" {
		return nil, ErrInvalidInput
	}
	if input.Email == "" {
		return nil, ErrInvalidInput
	}

	newUser, err := s.userRepo.Create(ctx, &repositories.CreateUserParams{
		Username: input.Username,
		Email:    input.Email,
	})

	if err != nil {
		return nil, ErrInternalServer
	}

	return newUser, nil
}

func (s *userService) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	if id == "" {
		return nil, ErrInvalidInput
	}

	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrInternalServer
	}

	return user, nil
}
