package entity

import (
	"time"
)

type Address struct {
	ID        int64
	UUID      string
	PersonID  int64
	City      string
	State     string
	Country   string
	IsPrimary bool
	CreatedAt time.Time
	UpdatedAt time.Time
	Active    bool
}

// GetFullLocation returns a formatted string with city, state, country
func (a *Address) GetFullLocation() string {
	parts := []string{}
	
	if a.City != "" {
		parts = append(parts, a.City)
	}
	
	if a.State != "" {
		parts = append(parts, a.State)
	}
	
	if a.Country != "" && a.Country != "Brazil" {
		parts = append(parts, a.Country)
	}
	
	if len(parts) == 0 {
		return ""
	}
	
	result := parts[0]
	for i := 1; i < len(parts); i++ {
		result += ", " + parts[i]
	}
	
	return result
}

// GetCityState returns just city and state
func (a *Address) GetCityState() string {
	if a.City != "" && a.State != "" {
		return a.City + ", " + a.State
	}
	if a.City != "" {
		return a.City
	}
	if a.State != "" {
		return a.State
	}
	return ""
}