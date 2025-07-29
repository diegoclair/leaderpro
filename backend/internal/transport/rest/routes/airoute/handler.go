package airoute

import (
	"errors"
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
	aiApp contract.AIApp
}

func NewHandler(aiApp contract.AIApp) *Handler {
	Once.Do(func() {
		instance = &Handler{
			aiApp: aiApp,
		}
	})

	return instance
}

func (h *Handler) handleChatWithAI(c echo.Context) error {
	ctx := routeutils.GetContext(c)

	var request viewmodel.AIChatRequest
	if err := c.Bind(&request); err != nil {
		return routeutils.ResponseInvalidRequestBody(c, err)
	}

	if request.Message == "" {
		return routeutils.ResponseInvalidRequestBody(c, errors.New("Message is required"))
	}

	chatReq := entity.ChatRequest{
		Message: request.Message,
	}

	response, err := h.aiApp.ChatWithLeadershipCoach(ctx, chatReq)
	if err != nil {
		return routeutils.HandleError(c, err)
	}

	result := viewmodel.AIChatResponse{
		Response: response.Response,
		UsageID:  response.UsageID,
	}

	return routeutils.ResponseAPIOk(c, result)
}

func (h *Handler) handleChatAboutPerson(c echo.Context) error {
	ctx := routeutils.GetContext(c)

	personUUID, err := routeutils.GetRequiredStringPathParam(c, "person_uuid", "Invalid person_uuid")
	if err != nil {
		return routeutils.HandleError(c, err)
	}

	var request viewmodel.AIPersonChatRequest
	if err := c.Bind(&request); err != nil {
		return routeutils.ResponseInvalidRequestBody(c, err)
	}

	if request.Message == "" {
		return routeutils.ResponseInvalidRequestBody(c, errors.New("Message is required"))
	}

	chatReq := entity.ChatRequest{
		Message:    request.Message,
		PersonUUID: &personUUID,
	}

	response, err := h.aiApp.ChatWithLeadershipCoach(ctx, chatReq)
	if err != nil {
		return routeutils.HandleError(c, err)
	}

	result := viewmodel.AIChatResponse{
		Response: response.Response,
		UsageID:  response.UsageID,
	}

	return routeutils.ResponseAPIOk(c, result)
}

func (h *Handler) handleSendFeedback(c echo.Context) error {
	ctx := routeutils.GetContext(c)

	usageIDStr, err := routeutils.GetRequiredStringPathParam(c, "usage_id", "Invalid usage_id")
	if err != nil {
		return routeutils.HandleError(c, err)
	}

	usageID, err := routeutils.GetRequiredParam(usageIDStr, routeutils.Int64Converter, "usage_id must be a valid integer")
	if err != nil {
		return routeutils.HandleError(c, err)
	}

	var request viewmodel.AIFeedbackRequest
	if err := c.Bind(&request); err != nil {
		return routeutils.ResponseInvalidRequestBody(c, err)
	}

	if request.Feedback != "helpful" && request.Feedback != "not_helpful" && request.Feedback != "neutral" {
		return routeutils.ResponseInvalidRequestBody(c, errors.New("Invalid feedback value. Must be: helpful, not_helpful, or neutral"))
	}

	err = h.aiApp.SendFeedback(ctx, usageID, request.Feedback, request.Comment)
	if err != nil {
		return routeutils.HandleError(c, err)
	}

	return routeutils.ResponseAPIOk(c, viewmodel.MessageResponse{
		Message: "Feedback sent successfully",
	})
}

func (h *Handler) handleGetUsageReport(c echo.Context) error {
	ctx := routeutils.GetContext(c)

	period := c.QueryParam("period")
	if period == "" {
		period = "month"
	}

	report, err := h.aiApp.GetUsageReport(ctx, period)
	if err != nil {
		return routeutils.HandleError(c, err)
	}

	result := viewmodel.AIUsageReportResponse{
		Period:            report.Period,
		TotalRequests:     report.TotalRequests,
		TotalTokens:       report.TotalTokens,
		TotalCostUSD:      report.TotalCostUSD,
		AverageResponseMs: report.AverageResponseMs,
		HelpfulFeedback:   report.HelpfulFeedback,
		UnhelpfulFeedback: report.UnhelpfulFeedback,
	}

	return routeutils.ResponseAPIOk(c, result)
}
