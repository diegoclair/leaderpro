package personroute

import (
	"net/http"

	"github.com/diegoclair/goswag"
	"github.com/diegoclair/goswag/models"
	"github.com/diegoclair/leaderpro/internal/transport/rest/routeutils"
	"github.com/diegoclair/leaderpro/internal/transport/rest/viewmodel"
	"github.com/diegoclair/leaderpro/infra"
)

const GroupRouteName = "companies/:company_uuid/people"

const (
	RootRoute            = ""
	PersonByUUIDRoute    = "/:person_uuid"
	PersonNotesRoute     = "/:person_uuid/notes"
	PersonNoteByUUIDRoute = "/:person_uuid/notes/:note_uuid"
	PersonTimelineRoute  = "/:person_uuid/timeline"
	PersonMentionsRoute  = "/:person_uuid/mentions"
)

type PersonRouter struct {
	ctrl *Handler
}

func NewRouter(ctrl *Handler) *PersonRouter {
	return &PersonRouter{
		ctrl: ctrl,
	}
}

func (r *PersonRouter) RegisterRoutes(g *routeutils.EchoGroups) {
	router := g.PrivateGroup.Group(GroupRouteName)

	router.POST(RootRoute, r.ctrl.handleCreatePerson).
		Summary("Create a new person").
		Description("Create a new person in the company").
		Read(viewmodel.PersonRequest{}).
		Returns([]models.ReturnType{{StatusCode: http.StatusCreated}}).
		PathParam("company_uuid", "company uuid", goswag.StringType, true).
		HeaderParam(infra.TokenKey.String(), infra.TokenKeyDescription, goswag.StringType, true)

	router.GET(RootRoute, r.ctrl.handleGetCompanyPeople).
		Summary("Get company people").
		Description("Get all people in the company, optionally filtered by search").
		Returns([]models.ReturnType{
			{
				StatusCode: http.StatusOK,
				Body:       []viewmodel.PersonResponse{},
			},
		}).
		PathParam("company_uuid", "company uuid", goswag.StringType, true).
		QueryParam("search", "search term to filter people", goswag.StringType, false).
		HeaderParam(infra.TokenKey.String(), infra.TokenKeyDescription, goswag.StringType, true)

	router.GET(PersonByUUIDRoute, r.ctrl.handleGetPersonByUUID).
		Summary("Get person by UUID").
		Description("Get person details by UUID").
		Returns([]models.ReturnType{
			{
				StatusCode: http.StatusOK,
				Body:       viewmodel.PersonResponse{},
			},
		}).
		PathParam("company_uuid", "company uuid", goswag.StringType, true).
		PathParam("person_uuid", "person uuid", goswag.StringType, true).
		HeaderParam(infra.TokenKey.String(), infra.TokenKeyDescription, goswag.StringType, true)

	router.PUT(PersonByUUIDRoute, r.ctrl.handleUpdatePerson).
		Summary("Update person").
		Description("Update person by UUID").
		Read(viewmodel.PersonRequest{}).
		Returns([]models.ReturnType{{StatusCode: http.StatusNoContent}}).
		PathParam("company_uuid", "company uuid", goswag.StringType, true).
		PathParam("person_uuid", "person uuid", goswag.StringType, true).
		HeaderParam(infra.TokenKey.String(), infra.TokenKeyDescription, goswag.StringType, true)

	router.DELETE(PersonByUUIDRoute, r.ctrl.handleDeletePerson).
		Summary("Delete person").
		Description("Delete person by UUID").
		Returns([]models.ReturnType{{StatusCode: http.StatusNoContent}}).
		PathParam("company_uuid", "company uuid", goswag.StringType, true).
		PathParam("person_uuid", "person uuid", goswag.StringType, true).
		HeaderParam(infra.TokenKey.String(), infra.TokenKeyDescription, goswag.StringType, true)

	router.POST(PersonNotesRoute, r.ctrl.handleCreateNote).
		Summary("Create a note for a person").
		Description("Create a new note (1:1, feedback, or observation) for a person").
		Read(viewmodel.CreateNoteRequest{}).
		Returns([]models.ReturnType{{StatusCode: http.StatusCreated}}).
		PathParam("company_uuid", "company uuid", goswag.StringType, true).
		PathParam("person_uuid", "person uuid", goswag.StringType, true).
		HeaderParam(infra.TokenKey.String(), infra.TokenKeyDescription, goswag.StringType, true)

	router.GET(PersonTimelineRoute, r.ctrl.handleGetPersonTimeline).
		Summary("Get person timeline").
		Description("Get timeline of direct notes for a person (1:1s and observations, excluding feedbacks/mentions)").
		Returns([]models.ReturnType{
			{
				StatusCode: http.StatusOK,
				Body:       viewmodel.PaginatedResponse[[]viewmodel.TimelineResponse]{},
			},
		}).
		PathParam("company_uuid", "company uuid", goswag.StringType, true).
		PathParam("person_uuid", "person uuid", goswag.StringType, true).
		QueryParam("page", "page number", goswag.NumberType, false).
		QueryParam("quantity", "items per page", goswag.NumberType, false).
		HeaderParam(infra.TokenKey.String(), infra.TokenKeyDescription, goswag.StringType, true)

	router.GET(PersonMentionsRoute, r.ctrl.handleGetPersonMentions).
		Summary("Get person mentions").
		Description("Get notes where this person was mentioned (feedbacks received)").
		Returns([]models.ReturnType{
			{
				StatusCode: http.StatusOK,
				Body:       viewmodel.PaginatedResponse[[]viewmodel.MentionResponse]{},
			},
		}).
		PathParam("company_uuid", "company uuid", goswag.StringType, true).
		PathParam("person_uuid", "person uuid", goswag.StringType, true).
		QueryParam("page", "page number", goswag.NumberType, false).
		QueryParam("quantity", "items per page", goswag.NumberType, false).
		HeaderParam(infra.TokenKey.String(), infra.TokenKeyDescription, goswag.StringType, true)

	router.PUT(PersonNoteByUUIDRoute, r.ctrl.handleUpdateNote).
		Summary("Update a note").
		Description("Update an existing note (1:1, feedback, or observation) by UUID").
		Read(viewmodel.UpdateNoteRequest{}).
		Returns([]models.ReturnType{{StatusCode: http.StatusNoContent}}).
		PathParam("company_uuid", "company uuid", goswag.StringType, true).
		PathParam("person_uuid", "person uuid", goswag.StringType, true).
		PathParam("note_uuid", "note uuid", goswag.StringType, true).
		HeaderParam(infra.TokenKey.String(), infra.TokenKeyDescription, goswag.StringType, true)

	router.DELETE(PersonNoteByUUIDRoute, r.ctrl.handleDeleteNote).
		Summary("Delete a note").
		Description("Delete an existing note (1:1, feedback, or observation) by UUID").
		Returns([]models.ReturnType{{StatusCode: http.StatusNoContent}}).
		PathParam("company_uuid", "company uuid", goswag.StringType, true).
		PathParam("person_uuid", "person uuid", goswag.StringType, true).
		PathParam("note_uuid", "note uuid", goswag.StringType, true).
		HeaderParam(infra.TokenKey.String(), infra.TokenKeyDescription, goswag.StringType, true)
}