package openai

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/diegoclair/leaderpro/internal/domain/entity"
	"github.com/sashabaranov/go-openai"
)

// Config represents OpenAI configuration
type Config struct {
	APIKey      string
	Model       string
	Temperature float64
	MaxTokens   int
	BaseURL     string
}

// Provider implements the AIProvider interface for OpenAI
type Provider struct {
	client *openai.Client
	config Config
}

// NewProvider creates a new instance of the OpenAI provider
func NewProvider(cfg Config) *Provider {
	clientConfig := openai.DefaultConfig(cfg.APIKey)
	if cfg.BaseURL != "" {
		clientConfig.BaseURL = cfg.BaseURL
	}
	
	return &Provider{
		client: openai.NewClientWithConfig(clientConfig),
		config: cfg,
	}
}

// Chat performs a conversation with the AI
func (p *Provider) Chat(ctx context.Context, req entity.ChatRequest, systemPrompt string, contextPrompt string) (entity.ChatResponse, error) {
	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: systemPrompt,
		},
	}
	
	if contextPrompt != "" {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: contextPrompt,
		})
	}
	
	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: req.Message,
	})
	
	request := openai.ChatCompletionRequest{
		Model:       p.config.Model,
		Messages:    messages,
		Temperature: float32(p.config.Temperature),
		MaxTokens:   p.config.MaxTokens,
	}
	
	resp, err := p.client.CreateChatCompletion(ctx, request)
	if err != nil {
		return entity.ChatResponse{}, fmt.Errorf("openai chat completion error: %w", err)
	}
	
	if len(resp.Choices) == 0 {
		return entity.ChatResponse{}, fmt.Errorf("no response choices from openai")
	}
	
	// Calculate usage information
	usage := entity.AIUsage{
		InputTokens:  resp.Usage.PromptTokens,
		OutputTokens: resp.Usage.CompletionTokens,
		TotalTokens:  resp.Usage.TotalTokens,
		CostUSD:      p.CalculateCost(resp.Usage.PromptTokens, resp.Usage.CompletionTokens),
	}
	
	response := entity.ChatResponse{
		Response: resp.Choices[0].Message.Content,
		Usage:    usage,
	}
	
	// Note: UsageID will be filled by the service layer after saving to database
	
	return response, nil
}

// ExtractAttributes extracts text attributes using AI
func (p *Provider) ExtractAttributes(ctx context.Context, req entity.ExtractionRequest, extractionPrompt string) (map[string]string, entity.AIUsage, error) {
	fullPrompt := fmt.Sprintf("%s\n\nNOTAS:\n%s", extractionPrompt, req.Content)
	
	request := openai.ChatCompletionRequest{
		Model: p.config.Model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: fullPrompt,
			},
		},
		Temperature: 0.3, // Lower for data extraction
		MaxTokens:   1000,
	}
	
	resp, err := p.client.CreateChatCompletion(ctx, request)
	if err != nil {
		return nil, entity.AIUsage{}, fmt.Errorf("openai extraction error: %w", err)
	}
	
	if len(resp.Choices) == 0 {
		return nil, entity.AIUsage{}, fmt.Errorf("no response choices from openai")
	}
	
	// Calculate usage information
	usage := entity.AIUsage{
		InputTokens:  resp.Usage.PromptTokens,
		OutputTokens: resp.Usage.CompletionTokens,
		TotalTokens:  resp.Usage.TotalTokens,
		CostUSD:      p.CalculateCost(resp.Usage.PromptTokens, resp.Usage.CompletionTokens),
	}
	
	var attributes map[string]string
	content := resp.Choices[0].Message.Content
	
	err = json.Unmarshal([]byte(content), &attributes)
	if err != nil {
		// If unable to parse, return empty (AI didn't return valid JSON)
		return map[string]string{}, usage, nil
	}
	
	return attributes, usage, nil
}

// GetProviderName returns the provider name
func (p *Provider) GetProviderName() string {
	return "openai"
}

// CalculateCost calculates cost based on tokens used
// Approximate prices for gpt-4o-mini (January 2025)
func (p *Provider) CalculateCost(inputTokens, outputTokens int) float64 {
	const (
		inputCostPer1KTokens  = 0.00015  // $0.15 per 1M input tokens
		outputCostPer1KTokens = 0.0006   // $0.60 per 1M output tokens
	)
	
	inputCost := float64(inputTokens) * inputCostPer1KTokens / 1000
	outputCost := float64(outputTokens) * outputCostPer1KTokens / 1000
	
	return inputCost + outputCost
}

// GetTokenUsage returns token usage information from the last request
// This is a helper function that can be useful for logging/debugging
func (p *Provider) GetTokenUsage(resp openai.ChatCompletionResponse) (inputTokens, outputTokens int) {
	return resp.Usage.PromptTokens, resp.Usage.CompletionTokens
}