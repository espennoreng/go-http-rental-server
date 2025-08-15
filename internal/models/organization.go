package models

import "time"

type Organization struct {
	ID        string  
	Name      string  
	CreatedBy string  
	CreatedAt time.Time
	UpdatedAt time.Time
}


