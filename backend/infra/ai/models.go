package ai

import "github.com/diegoclair/leaderpro/internal/domain/contract"

// OpenAIConfig represents OpenAI configuration
type OpenAIConfig struct {
	APIKey      string
	Model       string
	Temperature float64
	MaxTokens   int
	BaseURL     string
}

// AnthropicConfig represents Anthropic configuration
type AnthropicConfig struct {
	APIKey      string
	Model       string
	Temperature float64
	MaxTokens   int
}

// RateLimitConfig represents rate limiting configuration
type RateLimitConfig struct {
	RequestsPerMinute int
	RequestsPerHour   int
	RequestsPerDay    int
}

// Config represents AI configuration
type Config struct {
	DefaultProvider string
	OpenAI          OpenAIConfig
	Anthropic       AnthropicConfig
	RateLimit       RateLimitConfig
}

// Manager manage multiple ai provider
type Manager struct {
	providers       map[string]contract.AIProvider
	defaultProvider string
}
