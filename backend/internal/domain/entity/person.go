package entity

import (
	"time"
)

type Person struct {
	ID          int64
	UUID        string
	CompanyID   int64
	Name        string
	Email       string
	Position    string
	Department  string
	Phone       string
	Birthday    *time.Time
	StartDate   *time.Time
	IsManager   bool
	ManagerID   *int64
	Notes       string
	
	// Personal information
	HasKids     bool
	Gender      *string // "male", "female", "other"
	Interests   string
	Personality string
	
	// One-on-One information
	LastOneOnOneDate *time.Time `json:"last_one_on_one_date"`
	
	// Address information (loaded separately)
	PrimaryAddress *Address `json:"primary_address,omitempty"`
	
	// Metadata
	CreatedAt   time.Time
	UpdatedAt   time.Time
	CreatedBy   int64
	Active      bool
}

// GetAge returns the person's age based on their birthday
func (p *Person) GetAge() *int {
	if p.Birthday == nil {
		return nil
	}
	
	now := time.Now()
	age := now.Year() - p.Birthday.Year()
	
	// Adjust if birthday hasn't occurred this year
	if now.YearDay() < p.Birthday.YearDay() {
		age--
	}
	
	return &age
}

// GetTenure returns how long the person has been with the company in months
func (p *Person) GetTenure() *int {
	if p.StartDate == nil {
		return nil
	}
	
	months := int(time.Since(*p.StartDate).Hours() / 24 / 30)
	return &months
}