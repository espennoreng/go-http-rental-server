package api

import (
	"errors"

	"github.com/espennoreng/go-http-rental-server/internal/models"
)

type AddUserToOrganizationRequest struct {
	UserID string      `json:"user_id"`
	Role   models.Role `json:"role"`
}

func (r *AddUserToOrganizationRequest) Validate() error {
	if r.UserID == "" {
		return errors.New("user_id is required")
	}
	if r.Role == "" {
		return errors.New("role is required")
	}
	return nil
}

type UpdateUserRoleRequest struct {
	UserID string      `json:"user_id"`
	Role   models.Role `json:"role"`
}

func (r *UpdateUserRoleRequest) Validate() error {
	if r.UserID == "" {
		return errors.New("user_id is required")
	}
	if r.Role == "" {
		return errors.New("role is required")
	}
	return nil
}

type CreateOrganizationRequest struct {
	Name string `json:"name"`
}

func (r *CreateOrganizationRequest) Validate() error {
	if r.Name == "" {
		return errors.New("name is required")
	}
	return nil
}

