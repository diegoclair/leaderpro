package viewmodel

import (
	"time"

	"github.com/diegoclair/leaderpro/internal/domain/entity"
)

type CreateUser struct {
	Email    string `json:"email" validate:"required,email"`
	Name     string `json:"name" validate:"required,min=2"`
	Password string `json:"password" validate:"required,min=8"`
	Phone    string `json:"phone"`
}

func (c *CreateUser) ToEntity() entity.User {
	return entity.User{
		Email:    c.Email,
		Name:     c.Name,
		Password: c.Password,
		Phone:    c.Phone,
		Active:   true,
	}
}

type UpdateUser struct {
	Name         string `json:"name" validate:"required,min=2"`
	Phone        string `json:"phone"`
	ProfilePhoto string `json:"profile_photo"`
}

func (u *UpdateUser) ToEntity() entity.User {
	return entity.User{
		Name:         u.Name,
		Phone:        u.Phone,
		ProfilePhoto: u.ProfilePhoto,
	}
}

type User struct {
	UUID          string     `json:"uuid"`
	Email         string     `json:"email"`
	Name          string     `json:"name"`
	Phone         string     `json:"phone"`
	ProfilePhoto  string     `json:"profile_photo"`
	Plan          string     `json:"plan"`
	TrialEndsAt   *time.Time `json:"trial_ends_at"`
	SubscribedAt  *time.Time `json:"subscribed_at"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	LastLoginAt   *time.Time `json:"last_login_at"`
	EmailVerified bool       `json:"email_verified"`
}

func FromEntityUser(user entity.User) User {
	return User{
		UUID:          user.UUID,
		Email:         user.Email,
		Name:          user.Name,
		Phone:         user.Phone,
		ProfilePhoto:  user.ProfilePhoto,
		Plan:          user.Plan,
		TrialEndsAt:   user.TrialEndsAt,
		SubscribedAt:  user.SubscribedAt,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
		LastLoginAt:   user.LastLoginAt,
		EmailVerified: user.EmailVerified,
	}
}

type UpdateUserPreferences struct {
	Theme string `json:"theme" validate:"required,oneof=light dark"`
}

func (u *UpdateUserPreferences) ToEntity() entity.UserPreferences {
	return entity.UserPreferences{
		Theme: u.Theme,
	}
}

type UserPreferences struct {
	Theme string `json:"theme"`
}

func FromEntityUserPreferences(preferences entity.UserPreferences) UserPreferences {
	return UserPreferences{
		Theme: preferences.Theme,
	}
}
