package entity

import (
	"time"
)

// PersonAttribute represents a flexible person attribute
type PersonAttribute struct {
	ID                  int64     `db:"id" json:"id"`
	PersonID            int64     `db:"person_id" json:"person_id"`
	AttributeKey        string    `db:"attribute_key" json:"attribute_key"`
	AttributeValue      string    `db:"attribute_value" json:"attribute_value"`
	Source              string    `db:"source" json:"source"` // 'manual', 'ai_extracted', 'imported'
	ExtractedFromNoteID *int64    `db:"extracted_from_note_id" json:"extracted_from_note_id,omitempty"`
	CreatedAt           time.Time `db:"created_at" json:"created_at"`
	UpdatedAt           time.Time `db:"updated_at" json:"updated_at"`
}

// AIPrompt represents a versioned AI prompt
type AIPrompt struct {
	ID          int64   `db:"id" json:"id"`
	Type        string  `db:"type" json:"type"` // 'leadership_coach', 'attribute_extraction'
	Version     int     `db:"version" json:"version"`
	Prompt      string  `db:"prompt" json:"prompt"`
	Model       string  `db:"model" json:"model"`
	Temperature float64 `db:"temperature" json:"temperature"`
	MaxTokens   int     `db:"max_tokens" json:"max_tokens"`
	IsActive    bool    `db:"is_active" json:"is_active"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	CreatedBy   int64   `db:"created_by" json:"created_by"`
}

// AIUsageTracker represents AI usage tracking
type AIUsageTracker struct {
	ID             int64     `db:"id" json:"id"`
	UserID         int64     `db:"user_id" json:"user_id"`
	CompanyID      int64     `db:"company_id" json:"company_id"`
	PromptID       int64     `db:"prompt_id" json:"prompt_id"`
	PersonID       *int64    `db:"person_id" json:"person_id,omitempty"`
	RequestType    string    `db:"request_type" json:"request_type"` // 'chat', 'extraction', 'suggestion'
	TokensUsed     int       `db:"tokens_used" json:"tokens_used"`
	CostUSD        float64   `db:"cost_usd" json:"cost_usd"`
	ResponseTimeMs int       `db:"response_time_ms" json:"response_time_ms"`
	Feedback       *string   `db:"feedback" json:"feedback,omitempty"` // 'helpful', 'not_helpful', 'neutral'
	FeedbackComment *string  `db:"feedback_comment" json:"feedback_comment,omitempty"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
}

// AIConversation represents the content of an AI conversation
type AIConversation struct {
	ID          int64     `db:"id" json:"id"`
	UsageID     int64     `db:"usage_id" json:"usage_id"`
	UserMessage string    `db:"user_message" json:"user_message"`
	AIResponse  string    `db:"ai_response" json:"ai_response"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	ExpiresAt   time.Time `db:"expires_at" json:"expires_at"`
}

// ChatRequest represents a chat request to AI
type ChatRequest struct {
	Message    string  `json:"message" validate:"required,min=1,max=1000"`
	PersonUUID *string `json:"-"` // Set from URL parameter, not from request body
}

// ChatResponse represents the AI chat response
type ChatResponse struct {
	Response string    `json:"response"`
	UsageID  int64     `json:"usage_id"`
	Usage    AIUsage   `json:"-"` // Not exposed to client, used internally
}

// AIUsage represents usage information from AI provider
type AIUsage struct {
	InputTokens  int     `json:"input_tokens"`
	OutputTokens int     `json:"output_tokens"`
	TotalTokens  int     `json:"total_tokens"`
	CostUSD      float64 `json:"cost_usd"`
}

// ExtractionRequest represents an attribute extraction request
type ExtractionRequest struct {
	PersonID int64  `json:"person_id" validate:"required"`
	NoteID   int64  `json:"note_id" validate:"required"`
	Content  string `json:"content" validate:"required,min=1"`
}

// AttributesResponse represents an attribute extraction response
type AttributesResponse struct {
	Attributes []PersonAttribute `json:"attributes"`
	UsageID    int64             `json:"usage_id"`
	Usage      AIUsage           `json:"-"` // Not exposed to client, used internally
}

// PersonAIContext represents a person's context for AI
type PersonAIContext struct {
	Person      Person                    `json:"person"`
	Attributes  map[string]string         `json:"attributes"`
	RecentNotes []Note                    `json:"recent_notes"`
	LastMeeting *PersonLastMeeting        `json:"last_meeting,omitempty"`
}

// PersonLastMeeting represents information from the last meeting
type PersonLastMeeting struct {
	Date  time.Time `json:"date"`
	Notes string    `json:"notes"`
}

// AIUsageReport represents an AI usage report
type AIUsageReport struct {
	TotalRequests     int     `json:"total_requests"`
	TotalTokens       int     `json:"total_tokens"`
	TotalCostUSD      float64 `json:"total_cost_usd"`
	AverageResponseMs int     `json:"average_response_ms"`
	HelpfulFeedback   int     `json:"helpful_feedback"`
	UnhelpfulFeedback int     `json:"unhelpful_feedback"`
	Period            string  `json:"period"`
}