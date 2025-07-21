package userroute

import (
	"sync"

	"github.com/diegoclair/leaderpro/internal/application/dto"
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
	userService contract.UserApp
	authHelper  *shared.AuthHelper
}

func NewHandler(userService contract.UserApp, authHelper *shared.AuthHelper) *Handler {
	Once.Do(func() {
		instance = &Handler{
			userService: userService,
			authHelper:  authHelper,
		}
	})

	return instance
}

func (s *Handler) handleCreateUser(c echo.Context) error {
	ctx := routeutils.GetContext(c)

	input := viewmodel.CreateUser{}
	err := c.Bind(&input)
	if err != nil {
		return routeutils.ResponseInvalidRequestBody(c, err)
	}

	_, err = s.userService.CreateUser(ctx, input.ToEntity())
	if err != nil {
		return routeutils.HandleError(c, err)
	}

	// auto login after user creation to get access token and refresh token
	loginInput := dto.LoginInput{
		Email:    input.Email,
		Password: input.Password,
	}

	authResponse, err := s.authHelper.DoLogin(ctx, c, loginInput)
	if err != nil {
		return routeutils.HandleError(c, err)
	}

	return routeutils.ResponseAPIOk(c, authResponse)
}

func (s *Handler) handleGetProfile(c echo.Context) error {
	ctx := routeutils.GetContext(c)

	user, err := s.userService.GetProfile(ctx)
	if err != nil {
		return routeutils.HandleError(c, err)
	}

	return routeutils.ResponseAPIOk(c, viewmodel.FromEntityUser(user))
}

func (s *Handler) handleUpdateProfile(c echo.Context) error {
	ctx := routeutils.GetContext(c)

	input := viewmodel.UpdateUser{}
	err := c.Bind(&input)
	if err != nil {
		return routeutils.ResponseInvalidRequestBody(c, err)
	}

	user, err := s.userService.UpdateProfile(ctx, input.ToEntity())
	if err != nil {
		return routeutils.HandleError(c, err)
	}

	response := viewmodel.FromEntityUser(user)
	return routeutils.ResponseAPIOk(c, response)
}

func (s *Handler) handleGetUserPreferences(c echo.Context) error {
	ctx := routeutils.GetContext(c)

	preferences, err := s.userService.GetUserPreferences(ctx)
	if err != nil {
		return routeutils.HandleError(c, err)
	}

	return routeutils.ResponseAPIOk(c, viewmodel.FromEntityUserPreferences(preferences))
}

func (s *Handler) handleUpdateUserPreferences(c echo.Context) error {
	ctx := routeutils.GetContext(c)

	input := viewmodel.UpdateUserPreferences{}
	err := c.Bind(&input)
	if err != nil {
		return routeutils.ResponseInvalidRequestBody(c, err)
	}

	preferences, err := s.userService.UpdateUserPreferences(ctx, input.ToEntity())
	if err != nil {
		return routeutils.HandleError(c, err)
	}

	response := viewmodel.FromEntityUserPreferences(preferences)
	return routeutils.ResponseAPIOk(c, response)
}
