package dashboardroute

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
	dashboardService contract.DashboardApp
}

func NewHandler(dashboardService contract.DashboardApp) *Handler {
	Once.Do(func() {
		instance = &Handler{
			dashboardService: dashboardService,
		}
	})

	return instance
}

func (s *Handler) handleGetDashboard(c echo.Context) error {
	ctx := routeutils.GetContext(c)

	companyUUID, err := routeutils.GetRequiredStringPathParam(c, "company_uuid", "company_uuid path parameter is required")
	if err != nil {
		return routeutils.HandleError(c, err)
	}

	dashboard, err := s.dashboardService.GetDashboardData(ctx, companyUUID)
	if err != nil {
		return routeutils.HandleError(c, err)
	}

	response := viewmodel.DashboardResponse{}
	response.FillFromEntity(dashboard)

	return routeutils.ResponseAPIOk(c, response)
}