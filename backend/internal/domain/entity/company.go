package entity

import (
	"time"
)

type Company struct {
	ID          int64
	UUID        string
	Name        string
	Industry    string
	Size        string // small (10-49), medium (50-99), large (100-499), enterprise (500+)
	Role        string
	IsDefault   bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	UserOwnerID int64
	Active      bool
}

