package models

type Role string

const (
	RoleAdmin  Role = "admin"
	RoleMember Role = "member"
)

var ValidRoles = map[Role]bool{
	RoleMember: true,
	RoleAdmin:  true,
}
