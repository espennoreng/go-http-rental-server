package models

import "time"

type OrganizationUser struct {
	OrgID     string    `json:"org_id"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	Role      Role      `json:"role"`
}

type Role string

const (
	RoleAdmin  Role = "admin"
	RoleMember Role = "member"
)
