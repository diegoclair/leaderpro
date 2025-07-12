package entity

import (
	"time"
)

type Company struct {
	ID          int64
	UUID        string
	Name        string
	Description string
	Industry    string
	Size        string // small, medium, large, enterprise
	CreatedAt   time.Time
	UpdatedAt   time.Time
	CreatedBy   int64
	Active      bool
}

// CompanyUser represents the relationship between a user and a company
type CompanyUser struct {
	ID        int64
	CompanyID int64
	UserID    int64
	Role      string // owner, admin, member
	IsDefault bool   // indicates if this is the user's default company
	JoinedAt  time.Time
}

// UserCompany combines company data with user-specific information
type UserCompany struct {
	Company
	IsDefault bool
	Role      string
	JoinedAt  time.Time
}