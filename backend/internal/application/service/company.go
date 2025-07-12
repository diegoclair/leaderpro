package service

import (
	"context"

	"github.com/diegoclair/go_utils/logger"
	"github.com/diegoclair/go_utils/mysqlutils"
	"github.com/diegoclair/go_utils/resterrors"
	"github.com/diegoclair/go_utils/validator"
	"github.com/diegoclair/leaderpro/internal/domain"
	"github.com/diegoclair/leaderpro/internal/domain/contract"
	"github.com/diegoclair/leaderpro/internal/domain/entity"
	"github.com/twinj/uuid"
)

type companyService struct {
	cache     contract.CacheManager
	dm        contract.DataManager
	log       logger.Logger
	validator validator.Validator
	authApp   contract.AuthApp
}

func newCompanyService(infra domain.Infrastructure, authApp contract.AuthApp) contract.CompanyApp {
	return &companyService{
		cache:     infra.CacheManager(),
		dm:        infra.DataManager(),
		log:       infra.Logger(),
		validator: infra.Validator(),
		authApp:   authApp,
	}
}

func (s *companyService) CreateCompany(ctx context.Context, company entity.Company, isDefault bool) error {
	s.log.Info(ctx, "Process Started")
	defer s.log.Info(ctx, "Process Finished")

	// Generate UUID for the company
	company.UUID = uuid.NewV4().String()

	// Set default values
	if company.Size == "" {
		company.Size = "medium"
	}
	company.Active = true

	// Get logged user ID to set as creator
	userID, err := s.authApp.GetLoggedUserID(ctx)
	if err != nil {
		return err
	}
	company.CreatedBy = userID

	// Create company in database
	companyID, err := s.dm.Company().CreateCompany(ctx, company)
	if err != nil {
		s.log.Errorw(ctx, "error creating company", logger.Err(err))
		return err
	}

	// Add user to company as owner
	err = s.dm.Company().AddUserToCompany(ctx, companyID, userID, "owner")
	if err != nil {
		s.log.Errorw(ctx, "error adding user to company", logger.Err(err))
		return err
	}

	// Set as default company if requested
	if isDefault {
		err = s.dm.Company().SetCompanyAsDefault(ctx, companyID, userID)
		if err != nil {
			s.log.Errorw(ctx, "error setting company as default", logger.Err(err))
			return err
		}
	}

	s.log.Infow(ctx, "company created successfully",
		logger.Int64("company_id", companyID),
		logger.String("company_name", company.Name),
		logger.Bool("is_default", isDefault),
	)

	return nil
}

func (s *companyService) GetCompanyByUUID(ctx context.Context, companyUUID string) (entity.Company, error) {
	s.log.Info(ctx, "Process Started")
	defer s.log.Info(ctx, "Process Finished")

	company, err := s.dm.Company().GetCompanyByUUID(ctx, companyUUID)
	if err != nil {
		if mysqlutils.SQLNotFound(err.Error()) {
			return company, resterrors.NewNotFoundError("company not found")
		}
		s.log.Errorw(ctx, "error getting company by UUID", logger.Err(err))
		return company, err
	}

	return company, nil
}

func (s *companyService) GetUserCompanies(ctx context.Context) ([]entity.Company, error) {
	s.log.Info(ctx, "Process Started")
	defer s.log.Info(ctx, "Process Finished")

	userID, err := s.authApp.GetLoggedUserID(ctx)
	if err != nil {
		return nil, err
	}

	companies, err := s.dm.Company().GetCompaniesByUser(ctx, userID)
	if err != nil {
		s.log.Errorw(ctx, "error getting user companies", logger.Err(err))
		return nil, err
	}

	return companies, nil
}

func (s *companyService) UpdateCompany(ctx context.Context, companyUUID string, company entity.Company) error {
	s.log.Info(ctx, "Process Started")
	defer s.log.Info(ctx, "Process Finished")

	// Get company by UUID to get the ID
	existingCompany, err := s.dm.Company().GetCompanyByUUID(ctx, companyUUID)
	if err != nil {
		if mysqlutils.SQLNotFound(err.Error()) {
			return resterrors.NewNotFoundError("company not found")
		}
		s.log.Errorw(ctx, "error getting company by UUID", logger.Err(err))
		return err
	}

	// Update company
	err = s.dm.Company().UpdateCompany(ctx, existingCompany.ID, company)
	if err != nil {
		s.log.Errorw(ctx, "error updating company", logger.Err(err))
		return err
	}

	s.log.Infow(ctx, "company updated successfully",
		logger.Int64("company_id", existingCompany.ID),
		logger.String("company_uuid", companyUUID),
	)

	return nil
}

func (s *companyService) DeleteCompany(ctx context.Context, companyUUID string) error {
	s.log.Info(ctx, "Process Started")
	defer s.log.Info(ctx, "Process Finished")

	// Get company by UUID to get the ID
	company, err := s.dm.Company().GetCompanyByUUID(ctx, companyUUID)
	if err != nil {
		if mysqlutils.SQLNotFound(err.Error()) {
			return resterrors.NewNotFoundError("company not found")
		}
		s.log.Errorw(ctx, "error getting company by UUID", logger.Err(err))
		return err
	}

	// Delete company (soft delete)
	err = s.dm.Company().DeleteCompany(ctx, company.ID)
	if err != nil {
		s.log.Errorw(ctx, "error deleting company", logger.Err(err))
		return err
	}

	s.log.Infow(ctx, "company deleted successfully",
		logger.Int64("company_id", company.ID),
		logger.String("company_uuid", companyUUID),
	)

	return nil
}

func (s *companyService) GetUserCompaniesWithDefault(ctx context.Context) ([]entity.UserCompany, error) {
	s.log.Info(ctx, "Process Started")
	defer s.log.Info(ctx, "Process Finished")

	userID, err := s.authApp.GetLoggedUserID(ctx)
	if err != nil {
		return nil, err
	}

	companies, err := s.dm.Company().GetUserCompaniesWithDefault(ctx, userID)
	if err != nil {
		s.log.Errorw(ctx, "error getting user companies with default", logger.Err(err))
		return nil, err
	}

	return companies, nil
}

func (s *companyService) AddUserToCompany(ctx context.Context, companyUUID string, userEmail string, role string) error {
	s.log.Info(ctx, "Process Started")
	defer s.log.Info(ctx, "Process Finished")

	// Get company by UUID
	company, err := s.dm.Company().GetCompanyByUUID(ctx, companyUUID)
	if err != nil {
		if mysqlutils.SQLNotFound(err.Error()) {
			return resterrors.NewNotFoundError("company not found")
		}
		s.log.Errorw(ctx, "error getting company by UUID", logger.Err(err))
		return err
	}

	// Get user by email
	user, err := s.dm.User().GetUserByEmail(ctx, userEmail)
	if err != nil {
		if mysqlutils.SQLNotFound(err.Error()) {
			return resterrors.NewNotFoundError("user not found")
		}
		s.log.Errorw(ctx, "error getting user by email", logger.Err(err))
		return err
	}

	// Add user to company
	err = s.dm.Company().AddUserToCompany(ctx, company.ID, user.ID, role)
	if err != nil {
		s.log.Errorw(ctx, "error adding user to company", logger.Err(err))
		return err
	}

	s.log.Infow(ctx, "user added to company successfully",
		logger.Int64("company_id", company.ID),
		logger.Int64("user_id", user.ID),
		logger.String("role", role),
	)

	return nil
}