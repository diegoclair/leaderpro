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

	input := viewmodel.PersonRequest{}
	err := c.Bind(&input)
	if err != nil {
		return routeutils.ResponseInvalidRequestBody(c, err)
	}

	person := input.ToEntity()

	createdPerson, err := s.personService.CreatePerson(ctx, person)
	if err != nil {
		return routeutils.HandleError(c, err)
	}

	response := viewmodel.PersonResponse{}
	response.FillFromEntity(createdPerson)

	return routeutils.ResponseCreated(c, response)
}

func (s *Handler) handleGetCompanyPeople(c echo.Context) error {
	ctx := routeutils.GetContext(c)

	search := c.QueryParam("search")

	var people []entity.Person
	var err error
	if search != "" {
		people, err = s.personService.SearchPeople(ctx, search)
	} else {
		people, err = s.personService.GetCompanyPeople(ctx)
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

func (s *Handler) handleCreateNote(c echo.Context) error {
	ctx := routeutils.GetContext(c)

	personUUID, err := routeutils.GetRequiredStringPathParam(c, "person_uuid", "Invalid person_uuid")
	if err != nil {
		return routeutils.HandleError(c, err)
	}

	input := viewmodel.CreateNoteRequest{}
	err = c.Bind(&input)
	if err != nil {
		return routeutils.ResponseInvalidRequestBody(c, err)
	}

	note := input.ToEntity()

	createdNote, err := s.personService.CreateNote(ctx, note, personUUID)
	if err != nil {
		return routeutils.HandleError(c, err)
	}

	response := viewmodel.NoteResponse{}
	response.FillFromEntity(createdNote)

	return routeutils.ResponseCreated(c, response)
}

func (s *Handler) handleGetPersonTimeline(c echo.Context) error {
	ctx := routeutils.GetContext(c)

	personUUID, err := routeutils.GetRequiredStringPathParam(c, "person_uuid", "Invalid person_uuid")
	if err != nil {
		return routeutils.HandleError(c, err)
	}

	// Parse filters from query parameters using routeutils helpers
	types := routeutils.GetStringArrayQueryParam(c, "types", ",")
	feedbackTypes := routeutils.GetStringArrayQueryParam(c, "feedback_types", ",")

	filtersReq := viewmodel.TimelineFiltersRequest{
		SearchQuery:   c.QueryParam("search_query"),
		Types:         types,
		FeedbackTypes: feedbackTypes,
		Direction:     c.QueryParam("direction"),
		Period:        c.QueryParam("period"),
	}

	filters := filtersReq.ToEntity()
	take, skip := routeutils.GetPagingParams(c, "", "")

	timeline, totalRecords, err := s.personService.GetPersonTimeline(ctx, personUUID, filters, take, skip)
	if err != nil {
		return routeutils.HandleError(c, err)
	}

	response := []viewmodel.UnifiedTimelineResponse{}
	for _, entry := range timeline {
		item := viewmodel.UnifiedTimelineResponse{}
		item.FillFromUnifiedTimelineEntry(entry)
		response = append(response, item)
	}

	paginatedResponse := viewmodel.BuildPaginatedResponse(response, skip, take, totalRecords)

	return routeutils.ResponseAPIOk(c, paginatedResponse)
}

func (s *Handler) handleGetPersonMentions(c echo.Context) error {
	ctx := routeutils.GetContext(c)

	personUUID, err := routeutils.GetRequiredStringPathParam(c, "person_uuid", "Invalid person_uuid")
	if err != nil {
		return routeutils.HandleError(c, err)
	}

	take, skip := routeutils.GetPagingParams(c, "", "")

	mentions, totalRecords, err := s.personService.GetPersonMentions(ctx, personUUID, take, skip)
	if err != nil {
		return routeutils.HandleError(c, err)
	}

	response := []viewmodel.MentionResponse{}
	for _, mention := range mentions {
		item := viewmodel.MentionResponse{}
		item.FillFromMentionEntry(mention)
		response = append(response, item)
	}

	paginatedResponse := viewmodel.BuildPaginatedResponse(response, skip, take, totalRecords)

	return routeutils.ResponseAPIOk(c, paginatedResponse)
}

func (s *Handler) handleUpdateNote(c echo.Context) error {
	ctx := routeutils.GetContext(c)

	_, err := routeutils.GetRequiredStringPathParam(c, "person_uuid", "Invalid person_uuid")
	if err != nil {
		return routeutils.HandleError(c, err)
	}

	noteUUID, err := routeutils.GetRequiredStringPathParam(c, "note_uuid", "Invalid note_uuid")
	if err != nil {
		return routeutils.HandleError(c, err)
	}

	input := viewmodel.UpdateNoteRequest{}
	err = c.Bind(&input)
	if err != nil {
		return routeutils.ResponseInvalidRequestBody(c, err)
	}

	note := input.ToEntity()
	note.UUID = noteUUID

	err = s.personService.UpdateNote(ctx, noteUUID, note)
	if err != nil {
		return routeutils.HandleError(c, err)
	}

	return routeutils.ResponseNoContent(c)
}

func (s *Handler) handleDeleteNote(c echo.Context) error {
	ctx := routeutils.GetContext(c)

	_, err := routeutils.GetRequiredStringPathParam(c, "person_uuid", "Invalid person_uuid")
	if err != nil {
		return routeutils.HandleError(c, err)
	}

	noteUUID, err := routeutils.GetRequiredStringPathParam(c, "note_uuid", "Invalid note_uuid")
	if err != nil {
		return routeutils.HandleError(c, err)
	}

	err = s.personService.DeleteNote(ctx, noteUUID)
	if err != nil {
		return routeutils.HandleError(c, err)
	}

	return routeutils.ResponseNoContent(c)
}
