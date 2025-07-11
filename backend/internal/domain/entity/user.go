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
