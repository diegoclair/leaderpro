package dashboardroute

import (
	"net/http"

	"github.com/diegoclair/goswag"
	"github.com/diegoclair/goswag/models"
	"github.com/diegoclair/leaderpro/internal/transport/rest/routeutils"
	"github.com/diegoclair/leaderpro/internal/transport/rest/viewmodel"
	"github.com/diegoclair/leaderpro/infra"
)

const GroupRouteName = "dashboard"

const (
	RootRoute = ""
)

type DashboardRouter struct {
	ctrl *Handler
}

func NewRouter(ctrl *Handler) *DashboardRouter {
	return &DashboardRouter{
		ctrl: ctrl,
	}
}

func (r *DashboardRouter) RegisterRoutes(g *routeutils.EchoGroups) {
	router := g.PrivateGroup.Group(GroupRouteName)

	router.GET(RootRoute, r.ctrl.handleGetDashboard).
		Summary("Get dashboard data").
		Description("Get dashboard data with people and statistics for a specific company").
		Returns([]models.ReturnType{
			{
				StatusCode: http.StatusOK,
				Body:       viewmodel.DashboardResponse{},
			},
		}).
		QueryParam("company_uuid", "Company UUID to get dashboard data", goswag.StringType, true).
		HeaderParam(infra.TokenKey.String(), infra.TokenKeyDescription, goswag.StringType, true)
}