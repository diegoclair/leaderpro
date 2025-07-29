package contract

import (
	"context"

	"github.com/diegoclair/leaderpro/internal/domain/entity"
)

// AIProvider defines the interface for AI providers (OpenAI, Anthropic, etc)
type AIProvider interface {
	// Chat performs a conversation with the AI and returns response with usage info
	Chat(ctx context.Context, req entity.ChatRequest, systemPrompt string, contextPrompt string) (entity.ChatResponse, error)
	
	// ExtractAttributes extracts attributes from text using AI and returns usage info
	ExtractAttributes(ctx context.Context, req entity.ExtractionRequest, extractionPrompt string) (map[string]string, entity.AIUsage, error)
	
	// GetProviderName returns the provider name
	GetProviderName() string
}