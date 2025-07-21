package entity

import (
	"time"
)

// User represents the leader/manager using the platform
type User struct {
	ID           int64
	UUID         string
	Email        string
	Name         string
	Password     string
	Phone        string
	ProfilePhoto string

	// Subscription info
	Plan         string // basic, standard, unlimited
	TrialEndsAt  *time.Time
	SubscribedAt *time.Time

	// Metadata
	CreatedAt     time.Time
	UpdatedAt     time.Time
	LastLoginAt   *time.Time
	Active        bool
	EmailVerified bool
}

// IsTrialActive returns true if the user is in trial period
func (u *User) IsTrialActive() bool {
	if u.TrialEndsAt == nil {
		return false
	}
	return u.TrialEndsAt.After(time.Now())
}

// HasActiveSubscription returns true if user has an active paid subscription
func (u *User) HasActiveSubscription() bool {
	return u.SubscribedAt != nil && u.Plan != ""
}

// UserPreferences represents user preferences and settings
type UserPreferences struct {
	ID     int64
	UserID int64
	
	// Appearance
	Theme string // system, light, dark
	
	// Metadata
	CreatedAt time.Time
	UpdatedAt time.Time
}

// SetDefaults sets default values for user preferences
func (p *UserPreferences) SetDefaults() {
	if p.Theme == "" {
		p.Theme = "light"
	}
}
