package models

type Organization struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CreateOrganizationInput struct {
	Name string `json:"name"`
}
