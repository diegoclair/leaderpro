package ai

import (
	"fmt"

	"github.com/diegoclair/leaderpro/infra/ai/openai"
	"github.com/diegoclair/leaderpro/internal/domain/contract"
)

// NewManager creates a new instance of the AI manager
func NewManager(cfg Config) (*Manager, error) {
	manager := &Manager{
		providers:       make(map[string]contract.AIProvider),
		defaultProvider: cfg.DefaultProvider,
	}

	// Initialize OpenAI provider if configured
	if cfg.OpenAI.APIKey != "" {
		openaiConfig := openai.Config{
			APIKey:      cfg.OpenAI.APIKey,
			Model:       cfg.OpenAI.Model,
			Temperature: cfg.OpenAI.Temperature,
			MaxTokens:   cfg.OpenAI.MaxTokens,
			BaseURL:     cfg.OpenAI.BaseURL,
		}
		openaiProvider := openai.NewProvider(openaiConfig)
		manager.providers["openai"] = openaiProvider
	}

	// TODO: add other providers if needed
	// if cfg.Anthropic.APIKey != "" {
	//     anthropicProvider := anthropic.NewProvider(cfg.Anthropic)
	//     manager.providers["anthropic"] = anthropicProvider
	// }

	// Check if at least one provider was provided
	if len(manager.providers) == 0 {
		return nil, fmt.Errorf("no AI providers configured")
	}

	// check if default provider exists
	if _, exists := manager.providers[cfg.DefaultProvider]; !exists {
		// If existing provider is not the default, then use the first configured
		for name := range manager.providers {
			manager.defaultProvider = name
			break
		}
	}

	return manager, nil
}

// GetProvider returns a specific provider or default one
func (m *Manager) GetProvider(providerName string) contract.AIProvider {
	if providerName == "" {
		providerName = m.defaultProvider
	}

	if provider, exists := m.providers[providerName]; exists {
		return provider
	}

	// If the requested provider doesn't exist, return the default one
	return m.providers[m.defaultProvider]
}

// GetDefaultProvider returns the default provider
func (m *Manager) GetDefaultProvider() contract.AIProvider {
	return m.providers[m.defaultProvider]
}

// ListProviders returns the list of available providers
func (m *Manager) ListProviders() []string {
	var providers []string
	for name := range m.providers {
		providers = append(providers, name)
	}
	return providers
}

// GetDefaultProviderName returns the name of the default provider
func (m *Manager) GetDefaultProviderName() string {
	return m.defaultProvider
}

// IsProviderAvailable checks if a provider is available
func (m *Manager) IsProviderAvailable(providerName string) bool {
	_, exists := m.providers[providerName]
	return exists
}
