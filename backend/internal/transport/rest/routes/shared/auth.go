package shared

import (
	"context"

	infraContract "github.com/diegoclair/leaderpro/infra/contract"
	"github.com/diegoclair/leaderpro/internal/application/dto"
	"github.com/diegoclair/leaderpro/internal/domain/contract"
	"github.com/diegoclair/leaderpro/internal/transport/rest/viewmodel"
	"github.com/twinj/uuid"

	echo "github.com/labstack/echo/v4"
)

type AuthHelper struct {
	authService contract.AuthApp
	userService contract.UserApp
	authToken   infraContract.AuthToken
}

func NewAuthHelper(authService contract.AuthApp, userService contract.UserApp, authToken infraContract.AuthToken) *AuthHelper {
	return &AuthHelper{
		authService: authService,
		userService: userService,
		authToken:   authToken,
	}
}

func (h *AuthHelper) DoLogin(ctx context.Context, c echo.Context, loginInput dto.LoginInput) (*viewmodel.AuthResponse, error) {
	user, err := h.authService.Login(ctx, loginInput)
	if err != nil {
		return nil, err
	}

	sessionUUID := uuid.NewV4().String()
	req := infraContract.TokenPayloadInput{
		UserUUID:    user.UUID,
		SessionUUID: sessionUUID,
	}

	accessToken, tokenPayload, err := h.authToken.CreateAccessToken(ctx, req)
	if err != nil {
		return nil, err
	}

	refreshToken, refreshTokenPayload, err := h.authToken.CreateRefreshToken(ctx, req)
	if err != nil {
		return nil, err
	}

	// create session with user agent and client ip
	sessionReq := dto.Session{
		SessionUUID:           sessionUUID,
		UserID:                user.ID,
		RefreshToken:          refreshToken,
		UserAgent:             c.Request().UserAgent(),
		ClientIP:              c.RealIP(),
		RefreshTokenExpiredAt: refreshTokenPayload.ExpiredAt,
	}

	err = h.authService.CreateSession(ctx, sessionReq)
	if err != nil {
		return nil, err
	}

	// build response
	authResponse := &viewmodel.AuthResponse{
		User: viewmodel.FromEntityUser(user),
		Auth: viewmodel.LoginResponse{
			AccessToken:           accessToken,
			AccessTokenExpiresAt:  tokenPayload.ExpiredAt,
			RefreshToken:          refreshToken,
			RefreshTokenExpiresAt: refreshTokenPayload.ExpiredAt,
		},
	}

	return authResponse, nil
}
