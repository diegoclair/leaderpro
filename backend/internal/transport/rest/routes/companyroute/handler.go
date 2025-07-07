package companyroute

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
	companyService contract.CompanyApp
}

func NewHandler(companyService contract.CompanyApp) *Handler {
	Once.Do(func() {
		instance = &Handler{
			companyService: companyService,
		}
	})

	return instance
}

func (s *Handler) handleCreateCompany(c echo.Context) error {
	ctx := routeutils.GetContext(c)

	input := viewmodel.CompanyRequest{}
	err := c.Bind(&input)
	if err != nil {
		return routeutils.ResponseInvalidRequestBody(c, err)
	}

	err = s.companyService.CreateCompany(ctx, input.ToEntity())
	if err != nil {
		return routeutils.HandleError(c, err)
	}

	return routeutils.ResponseCreated(c)
}

func (s *Handler) handleGetCompanies(c echo.Context) error {
	ctx := routeutils.GetContext(c)

	companies, err := s.companyService.GetUserCompanies(ctx)
	if err != nil {
		return routeutils.HandleError(c, err)
	}

	response := []viewmodel.CompanyResponse{}
	for _, company := range companies {
		item := viewmodel.CompanyResponse{}
		item.FillFromEntity(company)
		response = append(response, item)
	}

	return routeutils.ResponseAPIOk(c, response)
}

func (s *Handler) handleGetCompanyByUUID(c echo.Context) error {
	ctx := routeutils.GetContext(c)

	companyUUID, err := routeutils.GetRequiredStringPathParam(c, "company_uuid", "Invalid company_uuid")
	if err != nil {
		return routeutils.HandleError(c, err)
	}

	company, err := s.companyService.GetCompanyByUUID(ctx, companyUUID)
	if err != nil {
		return routeutils.HandleError(c, err)
	}

	response := viewmodel.CompanyResponse{}
	response.FillFromEntity(company)

	return routeutils.ResponseAPIOk(c, response)
}

func (s *Handler) handleUpdateCompany(c echo.Context) error {
	ctx := routeutils.GetContext(c)

	companyUUID, err := routeutils.GetRequiredStringPathParam(c, "company_uuid", "Invalid company_uuid")
	if err != nil {
		return routeutils.HandleError(c, err)
	}

	input := viewmodel.CompanyRequest{}
	err = c.Bind(&input)
	if err != nil {
		return routeutils.ResponseInvalidRequestBody(c, err)
	}

	err = s.companyService.UpdateCompany(ctx, companyUUID, input.ToEntity())
	if err != nil {
		return routeutils.HandleError(c, err)
	}

	return routeutils.ResponseNoContent(c)
}

func (s *Handler) handleDeleteCompany(c echo.Context) error {
	ctx := routeutils.GetContext(c)

	companyUUID, err := routeutils.GetRequiredStringPathParam(c, "company_uuid", "Invalid company_uuid")
	if err != nil {
		return routeutils.HandleError(c, err)
	}

	err = s.companyService.DeleteCompany(ctx, companyUUID)
	if err != nil {
		return routeutils.HandleError(c, err)
	}

	return routeutils.ResponseNoContent(c)
}