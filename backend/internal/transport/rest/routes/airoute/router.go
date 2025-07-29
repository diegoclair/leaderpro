package airoute

import (
	"net/http"

	"github.com/diegoclair/goswag"
	"github.com/diegoclair/goswag/models"
	"github.com/diegoclair/leaderpro/infra"
	"github.com/diegoclair/leaderpro/internal/transport/rest/routeutils"
	"github.com/diegoclair/leaderpro/internal/transport/rest/viewmodel"
)

const GroupRouteName = "companies/:company_uuid/ai"

const (
	ChatRoute        = "/chat"
	FeedbackRoute    = "/usage/:usage_id/feedback"
	UsageReportRoute = "/usage"
)

type AIRouter struct {
	ctrl *Handler
}

func NewRouter(ctrl *Handler) *AIRouter {
	return &AIRouter{
		ctrl: ctrl,
	}
}

func (r *AIRouter) RegisterRoutes(g *routeutils.EchoGroups) {
	router := g.CompanyGroup.Group(GroupRouteName)

	router.POST(ChatRoute, r.ctrl.handleChatWithAI).
		Summary("Chat with AI Leadership Coach").
		Description("Send a message to the AI leadership coach for advice and guidance").
		Read(viewmodel.AIChatRequest{}).
		Returns([]models.ReturnType{{StatusCode: http.StatusOK, Body: viewmodel.AIChatResponse{}}}).
		PathParam("company_uuid", "Company UUID", goswag.StringType, true).
		HeaderParam(infra.TokenKey.String(), infra.TokenKeyDescription, goswag.StringType, true)

	router.GET(UsageReportRoute, r.ctrl.handleGetUsageReport).
		Summary("Get AI usage report").
		Description("Get usage statistics and costs for AI features").
		Returns([]models.ReturnType{{StatusCode: http.StatusOK, Body: viewmodel.AIUsageReportResponse{}}}).
		PathParam("company_uuid", "Company UUID", goswag.StringType, true).
		QueryParam("period", "Period (today, week, month, year, all)", goswag.StringType, false).
		HeaderParam(infra.TokenKey.String(), infra.TokenKeyDescription, goswag.StringType, true)

	router.POST(FeedbackRoute, r.ctrl.handleSendFeedback).
		Summary("Send feedback about AI response").
		Description("Provide feedback (helpful/not helpful) about an AI response").
		Read(viewmodel.AIFeedbackRequest{}).
		Returns([]models.ReturnType{{StatusCode: http.StatusOK, Body: viewmodel.MessageResponse{}}}).
		PathParam("company_uuid", "Company UUID", goswag.StringType, true).
		PathParam("usage_id", "Usage ID", goswag.StringType, true).
		HeaderParam(infra.TokenKey.String(), infra.TokenKeyDescription, goswag.StringType, true)

	r.registerPersonAIRoutes(g)
}

func (r *AIRouter) registerPersonAIRoutes(g *routeutils.EchoGroups) {
	personRouter := g.CompanyGroup.Group("companies/:company_uuid/people")

	personRouter.POST("/:person_uuid/ai/chat", r.ctrl.handleChatAboutPerson).
		Summary("Chat with AI about specific person").
		Description("Send a message to the AI with context about a specific person").
		Read(viewmodel.AIPersonChatRequest{}).
		Returns([]models.ReturnType{{StatusCode: http.StatusOK, Body: viewmodel.AIChatResponse{}}}).
		PathParam("company_uuid", "Company UUID", goswag.StringType, true).
		PathParam("person_uuid", "Person UUID", goswag.StringType, true).
		HeaderParam(infra.TokenKey.String(), infra.TokenKeyDescription, goswag.StringType, true)

}
