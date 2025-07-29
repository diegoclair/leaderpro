package viewmodel

type AIChatRequest struct {
	Message string `json:"message" validate:"required,min=1,max=1000"`
}

type AIPersonChatRequest struct {
	Message string `json:"message" validate:"required,min=1,max=1000"`
}

type AIChatResponse struct {
	Response string `json:"response"`
	UsageID  int64  `json:"usage_id"`
}

type AIFeedbackRequest struct {
	Feedback string `json:"feedback" validate:"required,oneof=helpful not_helpful neutral"`
	Comment  string `json:"comment,omitempty" validate:"max=500"`
}

type AIUsageReportResponse struct {
	Period            string  `json:"period"`
	TotalRequests     int     `json:"total_requests"`
	TotalTokens       int     `json:"total_tokens"`
	TotalCostUSD      float64 `json:"total_cost_usd"`
	AverageResponseMs int     `json:"average_response_ms"`
	HelpfulFeedback   int     `json:"helpful_feedback"`
	UnhelpfulFeedback int     `json:"unhelpful_feedback"`
}

// ========== Shared Response ==========

type MessageResponse struct {
	Message string `json:"message"`
}
