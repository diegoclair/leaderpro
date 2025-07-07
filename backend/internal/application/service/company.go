package service

import (
	"context"

	"github.com/diegoclair/leaderpro/internal/domain"
	"github.com/diegoclair/leaderpro/internal/domain/contract"
	"github.com/diegoclair/leaderpro/internal/domain/entity"
)

type companyService struct {
	infra domain.Infrastructure
}

func newCompanyService(infra domain.Infrastructure) contract.CompanyApp {
	return &companyService{
		infra: infra,
	}
}

func (s *companyService) CreateCompany(ctx context.Context, company entity.Company) error {
	// TODO: Implement
	return nil
}

func (s *companyService) GetCompanyByUUID(ctx context.Context, companyUUID string) (entity.Company, error) {
	// TODO: Implement
	return entity.Company{}, nil
}

func (s *companyService) GetUserCompanies(ctx context.Context) ([]entity.Company, error) {
	// TODO: Implement
	return []entity.Company{}, nil
}

func (s *companyService) UpdateCompany(ctx context.Context, companyUUID string, company entity.Company) error {
	// TODO: Implement
	return nil
}

func (s *companyService) DeleteCompany(ctx context.Context, companyUUID string) error {
	// TODO: Implement
	return nil
}

func (s *companyService) AddUserToCompany(ctx context.Context, companyUUID string, userEmail string, role string) error {
	// TODO: Implement
	return nil
}