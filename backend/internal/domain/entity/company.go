package entity

import (
	"time"
)

type Company struct {
	ID          int64
	UUID        string
	Name        string
	Industry    string
	Size        string // small, medium, large, enterprise
	Role        string
	IsDefault   bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	UserOwnerID int64
	Active      bool
}

