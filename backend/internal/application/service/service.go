package service

import (
	"errors"
	"time"

	"github.com/diegoclair/leaderpro/internal/domain"
	"github.com/diegoclair/leaderpro/internal/domain/contract"
)

type Apps struct {
	User    contract.UserApp
	Auth    contract.AuthApp
	Company contract.CompanyApp
	Person  contract.PersonApp
}

// New to get instance of all services
func New(infra domain.Infrastructure, accessTokenDuration time.Duration) (*Apps, error) {
	if err := validateInfrastructure(infra); err != nil {
		return nil, err
	}

	userApp := newUserSvc(infra)
	authApp := newAuthApp(infra, userApp, accessTokenDuration)

	return &Apps{
		User:    userApp,
		Auth:    authApp,
		Company: newCompanyService(infra, authApp),
		Person:  newPersonService(infra, authApp),
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
