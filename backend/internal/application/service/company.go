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

type companyApp struct {
	cache     contract.CacheManager
	dm        contract.DataManager
	log       logger.Logger
	validator validator.Validator
	authApp   contract.AuthApp
}

func newCompanyApp(infra domain.Infrastructure, authApp contract.AuthApp) contract.CompanyApp {
	return &companyApp{
		cache:     infra.CacheManager(),
		dm:        infra.DataManager(),
		log:       infra.Logger(),
		validator: infra.Validator(),
		authApp:   authApp,
	}
}

func (s *companyApp) CreateCompany(ctx context.Context, company entity.Company) (entity.Company, error) {
	s.log.Info(ctx, "Process Started")
	defer s.log.Info(ctx, "Process Finished")

	// Generate UUID for the company
	company.UUID = uuid.NewV4().String()

	// Set default values
	if company.Size == "" {
		company.Size = "medium"
	}
	company.Active = true

	// Get logged user ID to set as owner
	userID, err := s.authApp.GetLoggedUserID(ctx)
	if err != nil {
		return company, err
	}
	company.UserOwnerID = userID

	// If this company is being set as default, unset other defaults for this user
	if company.IsDefault {
		err = s.unsetUserDefaultCompanies(ctx, userID)
		if err != nil {
			s.log.Errorw(ctx, "error unsetting default companies", logger.Err(err))
			return company, err
		}
	}

	// Create company in database
	companyID, err := s.dm.Company().CreateCompany(ctx, company)
	if err != nil {
		s.log.Errorw(ctx, "error creating company", logger.Err(err))
		return company, err
	}

	// Set the created ID
	company.ID = companyID

	s.log.Infow(ctx, "company created successfully",
		logger.Int64("company_id", companyID),
		logger.String("company_name", company.Name),
		logger.Int64("user_owner_id", userID),
		logger.String("role", company.Role),
		logger.Bool("is_default", company.IsDefault),
	)

	return company, nil
}

func (s *companyApp) GetCompanyByUUID(ctx context.Context, companyUUID string) (entity.Company, error) {
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

func (s *companyApp) GetLoggedUserCompany(ctx context.Context) (entity.Company, error) {
	s.log.Info(ctx, "Process Started")
	defer s.log.Info(ctx, "Process Finished")

	// Get company UUID from context
	companyUUID, err := s.authApp.GetCompanyFromContext(ctx)
	if err != nil {
		return entity.Company{}, err
	}

	return s.GetCompanyByUUID(ctx, companyUUID)
}

func (s *companyApp) GetUserCompanies(ctx context.Context) ([]entity.Company, error) {
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

func (s *companyApp) UpdateCompany(ctx context.Context, companyUUID string, company entity.Company) error {
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

func (s *companyApp) UpdateLoggedUserCompany(ctx context.Context, company entity.Company) error {
	s.log.Info(ctx, "Process Started")
	defer s.log.Info(ctx, "Process Finished")

	// Get company UUID from context
	companyUUID, err := s.authApp.GetCompanyFromContext(ctx)
	if err != nil {
		return err
	}

	return s.UpdateCompany(ctx, companyUUID, company)
}

func (s *companyApp) DeleteCompany(ctx context.Context, companyUUID string) error {
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

// unsetUserDefaultCompanies removes the default flag from all companies of a user
func (s *companyApp) unsetUserDefaultCompanies(ctx context.Context, userID int64) error {
	companies, err := s.dm.Company().GetCompaniesByUser(ctx, userID)
	if err != nil {
		return err
	}

	for _, company := range companies {
		if company.IsDefault {
			company.IsDefault = false
			err = s.dm.Company().UpdateCompany(ctx, company.ID, company)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *companyApp) DeleteLoggedUserCompany(ctx context.Context) error {
	s.log.Info(ctx, "Process Started")
	defer s.log.Info(ctx, "Process Finished")

	// Get company UUID from context
	companyUUID, err := s.authApp.GetCompanyFromContext(ctx)
	if err != nil {
		return err
	}

	return s.DeleteCompany(ctx, companyUUID)
}

// ValidateCompanyOwnership validates that the user owns the specified company
func (s *companyApp) ValidateCompanyOwnership(ctx context.Context, companyUUID string, userUUID string) error {
	s.log.Infow(ctx, "Process Started: validating company ownership",
		logger.String("company_uuid", companyUUID),
		logger.String("user_uuid", userUUID),
	)
	defer s.log.Infow(ctx, "Process Finished")

	// Get company by UUID
	company, err := s.dm.Company().GetCompanyByUUID(ctx, companyUUID)
	if err != nil {
		if mysqlutils.SQLNotFound(err.Error()) {
			return resterrors.NewNotFoundError("company not found")
		}
		s.log.Errorw(ctx, "error getting company by UUID", logger.Err(err))
		return err
	}

	// Get user by UUID to get user ID
	user, err := s.dm.User().GetUserByUUID(ctx, userUUID)
	if err != nil {
		if mysqlutils.SQLNotFound(err.Error()) {
			return resterrors.NewUnauthorizedError("user not found")
		}
		s.log.Errorw(ctx, "error getting user by UUID", logger.Err(err))
		return err
	}

	// Check if the company belongs to the user
	if company.UserOwnerID != user.ID {
		s.log.Warnw(ctx, "user trying to access company they don't own",
			logger.Int64("company_owner_id", company.UserOwnerID),
			logger.Int64("requesting_user_id", user.ID),
			logger.String("company_uuid", companyUUID),
		)
		return resterrors.NewUnauthorizedError("you don't have permission to access this company")
	}

	return nil
}
