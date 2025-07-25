package userroute

import (
	"net/http"

	"github.com/diegoclair/goswag"
	"github.com/diegoclair/goswag/models"
	"github.com/diegoclair/leaderpro/infra"
	"github.com/diegoclair/leaderpro/internal/transport/rest/routeutils"
	"github.com/diegoclair/leaderpro/internal/transport/rest/viewmodel"
)

const GroupRouteName = "users"

const (
	RootUserRoute           = ""
	GetProfileRoute         = "/profile"
	UpdateProfileRoute      = "/profile"
	GetPreferencesRoute     = "/preferences"
	UpdatePreferencesRoute  = "/preferences"
)

type UserRouter struct {
	ctrl *Handler
}

func NewRouter(ctrl *Handler) *UserRouter {
	return &UserRouter{
		ctrl: ctrl,
	}
}

func (r *UserRouter) RegisterRoutes(g *routeutils.EchoGroups) {
	router := g.AppGroup.Group(GroupRouteName)
	privateRouter := g.PrivateGroup.Group(GroupRouteName)

	router.POST(RootUserRoute, r.ctrl.handleCreateUser).
		Summary("Create User").
		Description("Create a new user account and return authentication tokens").
		Read(viewmodel.CreateUser{}).
		Returns([]models.ReturnType{
			{
				StatusCode: http.StatusOK,
				Body:       viewmodel.AuthResponse{},
			},
		})

	privateRouter.GET(GetProfileRoute, r.ctrl.handleGetProfile).
		Summary("Get User Profile").
		Description("Get the current user's profile").
		Returns([]models.ReturnType{
			{
				StatusCode: http.StatusOK,
				Body:       viewmodel.User{},
			},
		}).
		HeaderParam(infra.TokenKey.String(), infra.TokenKeyDescription, goswag.StringType, true)

	privateRouter.PUT(UpdateProfileRoute, r.ctrl.handleUpdateProfile).
		Summary("Update User Profile").
		Description("Update the current user's profile").
		Read(viewmodel.UpdateUser{}).
		Returns([]models.ReturnType{
			{
				StatusCode: http.StatusOK,
				Body:       viewmodel.User{},
			},
		}).
		HeaderParam(infra.TokenKey.String(), infra.TokenKeyDescription, goswag.StringType, true)

	privateRouter.GET(GetPreferencesRoute, r.ctrl.handleGetUserPreferences).
		Summary("Get User Preferences").
		Description("Get the current user's preferences").
		Returns([]models.ReturnType{
			{
				StatusCode: http.StatusOK,
				Body:       viewmodel.UserPreferences{},
			},
		}).
		HeaderParam(infra.TokenKey.String(), infra.TokenKeyDescription, goswag.StringType, true)

	privateRouter.PUT(UpdatePreferencesRoute, r.ctrl.handleUpdateUserPreferences).
		Summary("Update User Preferences").
		Description("Update the current user's preferences").
		Read(viewmodel.UpdateUserPreferences{}).
		Returns([]models.ReturnType{
			{
				StatusCode: http.StatusOK,
				Body:       viewmodel.UserPreferences{},
			},
		}).
		HeaderParam(infra.TokenKey.String(), infra.TokenKeyDescription, goswag.StringType, true)
}
