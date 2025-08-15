package api

import (
	"time"

	"github.com/espennoreng/go-http-rental-server/internal/models"
)

type OrganizationUserResponse struct {
	UserID string      `json:"user_id"`
	OrgID  string      `json:"org_id"`
	Role   models.Role `json:"role"`
}

func NewOrganizationUserResponse(user *models.OrganizationUser) *OrganizationUserResponse {
	return &OrganizationUserResponse{
		UserID: user.UserID,
		OrgID:  user.OrgID,
		Role:   user.Role,
	}
}

type OrganizationMemberResponse struct {
	ID       string      `json:"id"`
	Username string      `json:"username"`
	Email    string      `json:"email"`
	Role     models.Role `json:"role"`
}

type OrganizationMembersResponse struct {
	Users []OrganizationMemberResponse `json:"users"`
}

func NewOrganizationMembersResponse(users []*models.UserWithRole) *OrganizationMembersResponse {
	memberResponses := make([]OrganizationMemberResponse, len(users))
	for i, user := range users {
		memberResponses[i] = OrganizationMemberResponse{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Role:     user.Role,
		}
	}
	return &OrganizationMembersResponse{Users: memberResponses}
}

type OrganizationResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func NewOrganizationResponse(org *models.Organization) *OrganizationResponse {
	return &OrganizationResponse{
		ID:   org.ID,
		Name: org.Name,
	}
}

type UserResponse struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func NewUserResponse(user *models.User) *UserResponse {
	return &UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}
}
