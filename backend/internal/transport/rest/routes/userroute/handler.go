package userroute

import (
	"sync"

	"github.com/diegoclair/leaderpro/internal/domain/contract"
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
}

func NewHandler(userService contract.UserApp) *Handler {
	Once.Do(func() {
		instance = &Handler{
			userService: userService,
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

	user, err := s.userService.CreateUser(ctx, input.ToEntity())
	if err != nil {
		return routeutils.HandleError(c, err)
	}

	response := viewmodel.FromEntityUser(user)
	return routeutils.ResponseAPIOk(c, response)
}

func (s *Handler) handleGetProfile(c echo.Context) error {
	ctx := routeutils.GetContext(c)

	user, err := s.userService.GetProfile(ctx)
	if err != nil {
		return routeutils.HandleError(c, err)
	}

	response := viewmodel.FromEntityUser(user)
	return routeutils.ResponseAPIOk(c, response)
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