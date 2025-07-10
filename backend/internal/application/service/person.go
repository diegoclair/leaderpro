package service

import (
	"context"

	"github.com/diegoclair/go_utils/logger"
	"github.com/diegoclair/go_utils/validator"
	"github.com/diegoclair/leaderpro/internal/domain"
	"github.com/diegoclair/leaderpro/internal/domain/contract"
	"github.com/diegoclair/leaderpro/internal/domain/entity"
)

type personService struct {
	dm        contract.DataManager
	log       logger.Logger
	validator validator.Validator
}

func newPersonService(infra domain.Infrastructure) contract.PersonApp {
	return &personService{
		dm:        infra.DataManager(),
		log:       infra.Logger(),
		validator: infra.Validator(),
	}
}

func (s *personService) CreatePerson(ctx context.Context, companyUUID string, person entity.Person) error {
	// Get company ID from UUID
	company, err := s.dm.Company().GetCompanyByUUID(ctx, companyUUID)
	if err != nil {
		return err
	}
	
	// Set the company ID in the person entity
	person.CompanyID = company.ID
	
	// TODO: Get the logged user ID from context - for now using 0
	// userID := context.GetUserIDFromContext(ctx)
	person.CreatedBy = 0 // Will be updated when authentication is implemented
	
	// Create the person
	_, err = s.dm.Person().CreatePerson(ctx, person)
	return err
}

func (s *personService) GetPersonByUUID(ctx context.Context, personUUID string) (entity.Person, error) {
	// TODO: Implement
	return entity.Person{}, nil
}

func (s *personService) GetCompanyPeople(ctx context.Context, companyUUID string) ([]entity.Person, error) {
	// TODO: Implement
	return []entity.Person{}, nil
}

func (s *personService) UpdatePerson(ctx context.Context, personUUID string, person entity.Person) error {
	// TODO: Implement
	return nil
}

func (s *personService) DeletePerson(ctx context.Context, personUUID string) error {
	// TODO: Implement
	return nil
}

func (s *personService) SearchPeople(ctx context.Context, companyUUID string, search string) ([]entity.Person, error) {
	// TODO: Implement
	return []entity.Person{}, nil
}