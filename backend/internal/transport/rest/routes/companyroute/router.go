package companyroute

import (
	"net/http"

	"github.com/diegoclair/goswag"
	"github.com/diegoclair/goswag/models"
	"github.com/diegoclair/leaderpro/internal/transport/rest/routeutils"
	"github.com/diegoclair/leaderpro/internal/transport/rest/viewmodel"
	"github.com/diegoclair/leaderpro/infra"
)

const GroupRouteName = "companies"

const (
	RootRoute           = ""
	CompanyByUUIDRoute  = "/:company_uuid"
)

type CompanyRouter struct {
	ctrl *Handler
}

func NewRouter(ctrl *Handler) *CompanyRouter {
	return &CompanyRouter{
		ctrl: ctrl,
	}
}

func (r *CompanyRouter) RegisterRoutes(g *routeutils.EchoGroups) {
	// Routes that don't need company validation (create and list all user companies)
	router := g.PrivateGroup.Group(GroupRouteName)

	router.POST(RootRoute, r.ctrl.handleCreateCompany).
		Summary("Create a new company").
		Read(viewmodel.CompanyRequest{}).
		Returns([]models.ReturnType{{StatusCode: http.StatusCreated}}).
		HeaderParam(infra.TokenKey.String(), infra.TokenKeyDescription, goswag.StringType, true)

	router.GET(RootRoute, r.ctrl.handleGetCompanies).
		Summary("Get user companies").
		Description("Get all companies for the authenticated user").
		Returns([]models.ReturnType{
			{
				StatusCode: http.StatusOK,
				Body:       []viewmodel.CompanyResponse{},
			},
		}).
		HeaderParam(infra.TokenKey.String(), infra.TokenKeyDescription, goswag.StringType, true)

	// Routes that need company ownership validation (specific company operations)
	companyRouter := g.CompanyGroup.Group(GroupRouteName)

	companyRouter.GET(CompanyByUUIDRoute, r.ctrl.handleGetCompanyByUUID).
		Summary("Get company by UUID").
		Description("Get company details by UUID").
		Returns([]models.ReturnType{
			{
				StatusCode: http.StatusOK,
				Body:       viewmodel.CompanyResponse{},
			},
		}).
		PathParam("company_uuid", "company uuid", goswag.StringType, true).
		HeaderParam(infra.TokenKey.String(), infra.TokenKeyDescription, goswag.StringType, true)

	companyRouter.PUT(CompanyByUUIDRoute, r.ctrl.handleUpdateCompany).
		Summary("Update company").
		Description("Update company by UUID").
		Read(viewmodel.CompanyRequest{}).
		Returns([]models.ReturnType{{StatusCode: http.StatusNoContent}}).
		PathParam("company_uuid", "company uuid", goswag.StringType, true).
		HeaderParam(infra.TokenKey.String(), infra.TokenKeyDescription, goswag.StringType, true)

	companyRouter.DELETE(CompanyByUUIDRoute, r.ctrl.handleDeleteCompany).
		Summary("Delete company").
		Description("Delete company by UUID").
		Returns([]models.ReturnType{{StatusCode: http.StatusNoContent}}).
		PathParam("company_uuid", "company uuid", goswag.StringType, true).
		HeaderParam(infra.TokenKey.String(), infra.TokenKeyDescription, goswag.StringType, true)
}