package models

import "time"

type OrganizationUser struct {
	ID        string    `json:"id"`
	OrgID     string    `json:"org_id"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	Role      Role      `json:"role"`
}
