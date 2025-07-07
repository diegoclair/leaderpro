package dto

import "time"

type LoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type Session struct {
	SessionUUID string    `json:"session_uuid"`
	UserID      int64     `json:"user_id"`
	AccessToken string    `json:"access_token"`
	CreatedAt   time.Time `json:"created_at"`
	ExpiresAt   time.Time `json:"expires_at"`
	Blocked     bool      `json:"blocked"`
}