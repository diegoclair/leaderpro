package auth

import (
	"time"

	"github.com/diegoclair/leaderpro/infra/contract"
)

type tokenPayloadInput struct {
	UserUUID    string
	SessionUUID string
}

func fromContractTokenPayloadInput(input contract.TokenPayloadInput) tokenPayloadInput {
	return tokenPayloadInput{
		UserUUID:    input.UserUUID,
		SessionUUID: input.SessionUUID,
	}
}

// tokenPayload represents the payload of a JWT token
type tokenPayload struct {
	UserUUID     string
	SessionUUID  string
	RefreshToken string
	IssuedAt     time.Time
	ExpiredAt    time.Time
}

func (t *tokenPayload) toContract() contract.TokenPayload {
	return contract.TokenPayload{
		UserUUID:     t.UserUUID,
		SessionUUID:  t.SessionUUID,
		RefreshToken: t.RefreshToken,
		IssuedAt:     t.IssuedAt,
		ExpiredAt:    t.ExpiredAt,
	}
}

func newPayload(input tokenPayloadInput, duration time.Duration) *tokenPayload {
	return &tokenPayload{
		SessionUUID: input.SessionUUID,
		UserUUID:    input.UserUUID,
		IssuedAt:    time.Now(),
		ExpiredAt:   time.Now().Add(duration),
	}
}

// Valid checks if the token payload is valid or not
func (p *tokenPayload) Valid() error {
	if time.Now().After(p.ExpiredAt) {
		return errExpiredToken
	}
	return nil
}
