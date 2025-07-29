package service

import (
	"errors"
	"time"

	"github.com/diegoclair/leaderpro/internal/domain"
	"github.com/diegoclair/leaderpro/internal/domain/contract"
)

type Apps struct {
	User      contract.UserApp
	Auth      contract.AuthApp
	Company   contract.CompanyApp
	Person    contract.PersonApp
	Dashboard contract.DashboardApp
	AI        contract.AIApp
}

// New to get instance of all services
func New(infra domain.Infrastructure, aiProvider contract.AIProvider, accessTokenDuration time.Duration) (*Apps, error) {
	if err := validateInfrastructure(infra); err != nil {
		return nil, err
	}

	userApp := newUserApp(infra)
	authApp := newAuthApp(infra, userApp, accessTokenDuration)
	personApp := newPersonApp(infra, authApp)

	// Initialize AI service if AI Provider is provided
	var aiApp contract.AIApp
	if aiProvider != nil {
		aiApp = newAIApp(infra, aiProvider, authApp)
		// Inject AI service into person service for automatic extraction
		personApp.SetAIApp(aiApp)
	}

	return &Apps{
		User:      userApp,
		Auth:      authApp,
		Company:   newCompanyApp(infra, authApp),
		Person:    personApp,
		Dashboard: newDashboardService(infra, authApp, personApp),
		AI:        aiApp,
	}, nil
}

// validateInfrastructure validate the dependencies needed to initialize the services
func validateInfrastructure(infra domain.Infrastructure) error {
	if infra.Logger() == nil {
		return errors.New("logger is required")
	}

	if infra.DataManager() == nil {
		return errors.New("data manager is required")
	}

	if infra.CacheManager() == nil {
		return errors.New("cache manager is required")
	}

	if infra.Crypto() == nil {
		return errors.New("crypto is required")
	}

	if infra.Validator() == nil {
		return errors.New("validator is required")
	}

	return nil
}
