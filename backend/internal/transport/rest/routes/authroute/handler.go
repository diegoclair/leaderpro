package authroute

import (
	"context"
	"sync"
	"time"

	"github.com/diegoclair/go_utils/logger"
	"github.com/diegoclair/leaderpro/infra"
	infraContract "github.com/diegoclair/leaderpro/infra/contract"
	"github.com/diegoclair/leaderpro/internal/domain/contract"
	"github.com/diegoclair/leaderpro/internal/transport/rest/routes/shared"
	"github.com/diegoclair/leaderpro/internal/transport/rest/routeutils"
	"github.com/diegoclair/leaderpro/internal/transport/rest/viewmodel"

	echo "github.com/labstack/echo/v4"
)

var (
	instance *Handler
	Once     sync.Once
)

type Handler struct {
	authService contract.AuthApp
	authToken   infraContract.AuthToken
	authHelper  *shared.AuthHelper
	log         logger.Logger
}

func NewHandler(authService contract.AuthApp, authToken infraContract.AuthToken, authHelper *shared.AuthHelper, log logger.Logger) *Handler {
	Once.Do(func() {
		instance = &Handler{
			authService: authService,
			authToken:   authToken,
			authHelper:  authHelper,
			log:         log,
		}
	})

	return instance
}

func (s *Handler) handleLogin(c echo.Context) error {
	ctx := routeutils.GetContext(c)

	input := viewmodel.Login{}
	err := c.Bind(&input)
	if err != nil {
		return routeutils.ResponseInvalidRequestBody(c, err)
	}

	authResponse, err := s.authHelper.DoLogin(ctx, c, input.ToDto())
	if err != nil {
		return routeutils.HandleError(c, err)
	}

	return routeutils.ResponseAPIOk(c, authResponse)
}

func (s *Handler) handleRefreshToken(c echo.Context) error {
	ctx := routeutils.GetContext(c)
	s.log.Info(ctx, "Starting refresh token")
	defer s.log.Info(ctx, "Finished refresh token")

	input := viewmodel.RefreshTokenRequest{}

	err := c.Bind(&input)
	if err != nil {
		return routeutils.ResponseInvalidRequestBody(c, err)
	}

	refreshPayload, err := s.authToken.VerifyToken(ctx, input.RefreshToken)
	if err != nil {
		return routeutils.HandleError(c, err)
	}

	ctx = context.WithValue(ctx, infra.UserUUIDKey, refreshPayload.UserUUID)
	ctx = context.WithValue(ctx, infra.SessionKey, refreshPayload.SessionUUID)

	session, err := s.authService.GetSessionByUUID(ctx, refreshPayload.SessionUUID)
	if err != nil {
		return routeutils.HandleError(c, err)
	}

	if session.IsBlocked {
		return routeutils.ResponseUnauthorizedError(c, "session blocked")
	}

	if session.RefreshToken != input.RefreshToken {
		return routeutils.ResponseUnauthorizedError(c, "mismatched session token")
	}

	if time.Now().After(session.RefreshTokenExpiredAt) {
		return routeutils.ResponseUnauthorizedError(c, "expired session")
	}

	req := infraContract.TokenPayloadInput{
		UserUUID:    refreshPayload.UserUUID,
		SessionUUID: refreshPayload.SessionUUID,
	}
	accessToken, accessPayload, err := s.authToken.CreateAccessToken(ctx, req)
	if err != nil {
		return routeutils.HandleError(c, err)
	}

	response := viewmodel.RefreshTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayload.ExpiredAt,
	}

	return routeutils.ResponseAPIOk(c, response)
}

func (s *Handler) handleLogout(c echo.Context) error {
	accessToken := c.Request().Header.Get(infra.TokenKey.String())
	ctx := routeutils.GetContext(c)

	err := s.authService.Logout(ctx, accessToken)
	if err != nil {
		return routeutils.HandleError(c, err)
	}

	return routeutils.ResponseAPIOk(c, struct{}{})
}
