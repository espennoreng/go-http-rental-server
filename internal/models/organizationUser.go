package models

import "time"

type OrganizationUser struct {
	ID        string
	OrgID     string
	UserID    string
	CreatedAt time.Time
	Role      Role
}
