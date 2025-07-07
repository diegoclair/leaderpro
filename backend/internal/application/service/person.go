package service

import (
	"context"

	"github.com/diegoclair/leaderpro/internal/domain"
	"github.com/diegoclair/leaderpro/internal/domain/contract"
	"github.com/diegoclair/leaderpro/internal/domain/entity"
)

type personService struct {
	infra domain.Infrastructure
}

func newPersonService(infra domain.Infrastructure) contract.PersonApp {
	return &personService{
		infra: infra,
	}
}

func (s *personService) CreatePerson(ctx context.Context, person entity.Person) error {
	// TODO: Implement
	return nil
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