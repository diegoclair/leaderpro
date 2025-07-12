package authroute

import (
	"net/http"

	"github.com/diegoclair/goswag"
	"github.com/diegoclair/goswag/models"
	"github.com/diegoclair/leaderpro/internal/transport/rest/routeutils"
	"github.com/diegoclair/leaderpro/internal/transport/rest/viewmodel"

	"github.com/diegoclair/leaderpro/infra"
)

const GroupRouteName = "auth"

const (
	LoginRoute        = "/login"
	LogoutRoute       = "/logout"
	RefreshTokenRoute = "/refresh-token"
)

type AuthRouter struct {
	ctrl *Handler
}

func NewRouter(ctrl *Handler) *AuthRouter {
	return &AuthRouter{
		ctrl: ctrl,
	}
}

func (r *AuthRouter) RegisterRoutes(g *routeutils.EchoGroups) {
	router := g.AppGroup.Group(GroupRouteName)
	privateRouter := g.PrivateGroup.Group(GroupRouteName)

	router.POST(LoginRoute, r.ctrl.handleLogin).
		Summary("Login").
		Description("Login user and return user data with authentication tokens").
		Read(viewmodel.Login{}).
		Returns([]models.ReturnType{
			{
				StatusCode: http.StatusOK,
				Body:       viewmodel.AuthResponse{},
			},
		})

	router.POST(RefreshTokenRoute, r.ctrl.handleRefreshToken).
		Summary("Refresh Token").
		Description("Generate a new access token using the refresh token").
		Read(viewmodel.RefreshTokenRequest{}).
		Returns([]models.ReturnType{
			{
				StatusCode: http.StatusOK,
				Body:       viewmodel.RefreshTokenResponse{},
			},
		})

	privateRouter.POST(LogoutRoute, r.ctrl.handleLogout).
		Summary("Logout").
		Description("Logout the user").
		Returns([]models.ReturnType{
			{
				StatusCode: http.StatusOK,
			},
		}).
		HeaderParam(infra.TokenKey.String(), infra.TokenKeyDescription, goswag.StringType, true)
}
