package personroute

import (
	"sync"

	"github.com/diegoclair/leaderpro/internal/domain/contract"
	"github.com/diegoclair/leaderpro/internal/domain/entity"
	"github.com/diegoclair/leaderpro/internal/transport/rest/routeutils"
	"github.com/diegoclair/leaderpro/internal/transport/rest/viewmodel"

	echo "github.com/labstack/echo/v4"
)

var (
	instance *Handler
	Once     sync.Once
)

type Handler struct {
	personService contract.PersonApp
}

func NewHandler(personService contract.PersonApp) *Handler {
	Once.Do(func() {
		instance = &Handler{
			personService: personService,
		}
	})

	return instance
}

func (s *Handler) handleCreatePerson(c echo.Context) error {
	ctx := routeutils.GetContext(c)

	companyUUID, err := routeutils.GetRequiredStringPathParam(c, "company_uuid", "Invalid company_uuid")
	if err != nil {
		return routeutils.HandleError(c, err)
	}

	input := viewmodel.PersonRequest{}
	err = c.Bind(&input)
	if err != nil {
		return routeutils.ResponseInvalidRequestBody(c, err)
	}

	person := input.ToEntity()

	createdPerson, err := s.personService.CreatePerson(ctx, person, companyUUID)
	if err != nil {
		return routeutils.HandleError(c, err)
	}

	response := viewmodel.PersonResponse{}
	response.FillFromEntity(createdPerson)

	return routeutils.ResponseCreated(c, response)
}

func (s *Handler) handleGetCompanyPeople(c echo.Context) error {
	ctx := routeutils.GetContext(c)

	companyUUID, err := routeutils.GetRequiredStringPathParam(c, "company_uuid", "Invalid company_uuid")
	if err != nil {
		return routeutils.HandleError(c, err)
	}

	search := c.QueryParam("search")
	
	var people []entity.Person
	if search != "" {
		people, err = s.personService.SearchPeople(ctx, companyUUID, search)
	} else {
		people, err = s.personService.GetCompanyPeople(ctx, companyUUID)
	}
	
	if err != nil {
		return routeutils.HandleError(c, err)
	}

	response := []viewmodel.PersonResponse{}
	for _, person := range people {
		item := viewmodel.PersonResponse{}
		item.FillFromEntity(person)
		response = append(response, item)
	}

	return routeutils.ResponseAPIOk(c, response)
}

func (s *Handler) handleGetPersonByUUID(c echo.Context) error {
	ctx := routeutils.GetContext(c)

	personUUID, err := routeutils.GetRequiredStringPathParam(c, "person_uuid", "Invalid person_uuid")
	if err != nil {
		return routeutils.HandleError(c, err)
	}

	person, err := s.personService.GetPersonByUUID(ctx, personUUID)
	if err != nil {
		return routeutils.HandleError(c, err)
	}

	response := viewmodel.PersonResponse{}
	response.FillFromEntity(person)

	return routeutils.ResponseAPIOk(c, response)
}

func (s *Handler) handleUpdatePerson(c echo.Context) error {
	ctx := routeutils.GetContext(c)

	personUUID, err := routeutils.GetRequiredStringPathParam(c, "person_uuid", "Invalid person_uuid")
	if err != nil {
		return routeutils.HandleError(c, err)
	}

	input := viewmodel.PersonRequest{}
	err = c.Bind(&input)
	if err != nil {
		return routeutils.ResponseInvalidRequestBody(c, err)
	}

	err = s.personService.UpdatePerson(ctx, personUUID, input.ToEntity())
	if err != nil {
		return routeutils.HandleError(c, err)
	}

	return routeutils.ResponseNoContent(c)
}

func (s *Handler) handleDeletePerson(c echo.Context) error {
	ctx := routeutils.GetContext(c)

	personUUID, err := routeutils.GetRequiredStringPathParam(c, "person_uuid", "Invalid person_uuid")
	if err != nil {
		return routeutils.HandleError(c, err)
	}

	err = s.personService.DeletePerson(ctx, personUUID)
	if err != nil {
		return routeutils.HandleError(c, err)
	}

	return routeutils.ResponseNoContent(c)
}