package dto

import (
	"context"
	"time"

	"github.com/diegoclair/go_utils/validator"
)

type LoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (l *LoginInput) Validate(ctx context.Context, v validator.Validator) error {
	return v.ValidateStruct(ctx, l)
}

type Session struct {
	SessionID             int64
	SessionUUID           string `validate:"required,uuid"`
	UserID                int64  `validate:"required"`
	RefreshToken          string `validate:"required"`
	UserAgent             string
	ClientIP              string
	IsBlocked             bool
	RefreshTokenExpiredAt time.Time
}

func (s *Session) Validate(ctx context.Context, v validator.Validator) error {
	return v.ValidateStruct(ctx, s)
}
