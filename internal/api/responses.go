package api

import "github.com/espennoreng/go-http-rental-server/internal/models"

type OrganizationUserResponse struct {
	UserID string      `json:"user_id"`
	OrgID  string      `json:"org_id"`
	Role   models.Role `json:"role"`
}

func toUserResponse(user *models.OrganizationUser) *OrganizationUserResponse {
	return &OrganizationUserResponse{
		UserID: user.UserID,
		OrgID:  user.OrgID,
		Role:   user.Role,
	}
}

type UserResponse struct {
	ID       string      `json:"id"`
	Username string      `json:"username"`
	Email    string      `json:"email"`
	Role     models.Role `json:"role"`
}

type UsersResponse struct {
	Users []UserResponse `json:"users"`
}

func toUsersResponse(users []*models.UserWithRole) *UsersResponse {
	userResponses := make([]UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = UserResponse{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Role:     user.Role,
		}
	}
	return &UsersResponse{Users: userResponses}
}

type OrganizationResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func toOrganizationResponse(org *models.Organization) *OrganizationResponse {
	return &OrganizationResponse{
		ID:   org.ID,
		Name: org.Name,
	}
}
