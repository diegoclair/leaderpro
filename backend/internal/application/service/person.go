package service

import (
	"context"
	"time"

	"github.com/diegoclair/go_utils/logger"
	"github.com/diegoclair/go_utils/mysqlutils"
	"github.com/diegoclair/go_utils/resterrors"
	"github.com/diegoclair/go_utils/validator"
	"github.com/diegoclair/leaderpro/internal/domain"
	"github.com/diegoclair/leaderpro/internal/domain/contract"
	"github.com/diegoclair/leaderpro/internal/domain/entity"
	"github.com/twinj/uuid"
)

type personService struct {
	dm        contract.DataManager
	log       logger.Logger
	validator validator.Validator
	authApp   contract.AuthApp
}

func newPersonService(infra domain.Infrastructure, authApp contract.AuthApp) contract.PersonApp {
	return &personService{
		dm:        infra.DataManager(),
		log:       infra.Logger(),
		validator: infra.Validator(),
		authApp:   authApp,
	}
}

// validateUserCompanyAccess checks if the logged user has access to a specific company
// Returns the company entity if access is granted, or error if not
func (s *personService) validateUserCompanyAccess(ctx context.Context, userID, companyID int64) (entity.Company, error) {
	company, err := s.dm.Company().GetCompanyByID(ctx, companyID)
	if err != nil {
		if mysqlutils.SQLNotFound(err.Error()) {
			return company, resterrors.NewNotFoundError("company not found")
		}
		s.log.Errorw(ctx, "error getting company by ID", logger.Err(err))
		return company, err
	}

	// Check if the company belongs to the logged user
	if company.UserOwnerID != userID {
		s.log.Errorw(ctx, "user trying to access company they don't own",
			logger.Int64("company_id", companyID),
			logger.Int64("company_owner_id", company.UserOwnerID),
			logger.Int64("logged_user_id", userID),
		)
		return company, resterrors.NewUnauthorizedError("you don't have permission to access this company")
	}

	return company, nil
}

func (s *personService) CreatePerson(ctx context.Context, person entity.Person, companyUUID string) (entity.Person, error) {
	s.log.Info(ctx, "Process Started")
	defer s.log.Info(ctx, "Process Finished")

	// Generate UUID for the person
	person.UUID = uuid.NewV4().String()

	// Get company by UUID and validate it belongs to the logged user
	company, err := s.dm.Company().GetCompanyByUUID(ctx, companyUUID)
	if err != nil {
		if mysqlutils.SQLNotFound(err.Error()) {
			return person, resterrors.NewNotFoundError("company not found")
		}
		s.log.Errorw(ctx, "error getting company by UUID", logger.Err(err))
		return person, err
	}

	// Validate that the company belongs to the logged user
	userID, err := s.authApp.GetLoggedUserID(ctx)
	if err != nil {
		return person, err
	}

	if company.UserOwnerID != userID {
		s.log.Errorw(ctx, "user trying to create person in company they don't own",
			logger.Int64("company_owner_id", company.UserOwnerID),
			logger.Int64("logged_user_id", userID),
		)
		return person, resterrors.NewUnauthorizedError("you don't have permission to add people to this company")
	}

	// Set the company ID and creator in the person entity
	person.CompanyID = company.ID
	person.CreatedBy = userID
	person.Active = true

	// Create the person in database
	personID, err := s.dm.Person().CreatePerson(ctx, person)
	if err != nil {
		s.log.Errorw(ctx, "error creating person", logger.Err(err))
		return person, err
	}

	// Set the ID and timestamps for the response
	person.ID = personID
	person.CreatedAt = time.Now()
	person.UpdatedAt = time.Now()

	s.log.Infow(ctx, "person created successfully",
		logger.Int64("person_id", personID),
		logger.String("person_name", person.Name),
		logger.Int64("company_id", company.ID),
		logger.String("company_name", company.Name),
	)

	return person, nil
}

func (s *personService) GetPersonByUUID(ctx context.Context, personUUID string) (entity.Person, error) {
	s.log.Info(ctx, "Process Started")
	defer s.log.Info(ctx, "Process Finished")

	person, err := s.dm.Person().GetPersonByUUID(ctx, personUUID)
	if err != nil {
		if mysqlutils.SQLNotFound(err.Error()) {
			return person, resterrors.NewNotFoundError("person not found")
		}
		s.log.Errorw(ctx, "error getting person by UUID", logger.Err(err))
		return person, err
	}

	// Validate that the person belongs to a company owned by the logged user
	userID, err := s.authApp.GetLoggedUserID(ctx)
	if err != nil {
		return person, err
	}

	// Validate user has access to the person's company
	_, err = s.validateUserCompanyAccess(ctx, userID, person.CompanyID)
	if err != nil {
		return person, err
	}

	return person, nil
}

func (s *personService) GetCompanyPeople(ctx context.Context, companyUUID string) ([]entity.Person, error) {
	s.log.Info(ctx, "Process Started")
	defer s.log.Info(ctx, "Process Finished")

	// Get company by UUID and validate it belongs to the logged user
	company, err := s.dm.Company().GetCompanyByUUID(ctx, companyUUID)
	if err != nil {
		if mysqlutils.SQLNotFound(err.Error()) {
			return nil, resterrors.NewNotFoundError("company not found")
		}
		s.log.Errorw(ctx, "error getting company by UUID", logger.Err(err))
		return nil, err
	}

	// Validate that the company belongs to the logged user
	userID, err := s.authApp.GetLoggedUserID(ctx)
	if err != nil {
		return nil, err
	}

	if company.UserOwnerID != userID {
		s.log.Errorw(ctx, "user trying to get people from company they don't own",
			logger.Int64("company_owner_id", company.UserOwnerID),
			logger.Int64("logged_user_id", userID),
		)
		return nil, resterrors.NewUnauthorizedError("you don't have permission to access people from this company")
	}

	// Get people by company ID
	people, err := s.dm.Person().GetPersonsByCompany(ctx, company.ID)
	if err != nil {
		s.log.Errorw(ctx, "error getting people by company", logger.Err(err))
		return nil, err
	}

	s.log.Infow(ctx, "people retrieved successfully",
		logger.Int("people_count", len(people)),
		logger.Int64("company_id", company.ID),
		logger.String("company_name", company.Name),
	)

	return people, nil
}

func (s *personService) UpdatePerson(ctx context.Context, personUUID string, person entity.Person) error {
	s.log.Info(ctx, "Process Started")
	defer s.log.Info(ctx, "Process Finished")

	// Get existing person to validate ownership
	existingPerson, err := s.dm.Person().GetPersonByUUID(ctx, personUUID)
	if err != nil {
		if mysqlutils.SQLNotFound(err.Error()) {
			return resterrors.NewNotFoundError("person not found")
		}
		s.log.Errorw(ctx, "error getting person by UUID", logger.Err(err))
		return err
	}

	// Validate that the person belongs to a company owned by the logged user
	userID, err := s.authApp.GetLoggedUserID(ctx)
	if err != nil {
		return err
	}

	// Validate user has access to the person's company
	_, err = s.validateUserCompanyAccess(ctx, userID, existingPerson.CompanyID)
	if err != nil {
		return err
	}

	// Update the person
	err = s.dm.Person().UpdatePerson(ctx, existingPerson.ID, person)
	if err != nil {
		s.log.Errorw(ctx, "error updating person", logger.Err(err))
		return err
	}

	s.log.Infow(ctx, "person updated successfully",
		logger.String("person_uuid", personUUID),
		logger.String("person_name", person.Name),
		logger.Int64("company_id", existingPerson.CompanyID),
	)

	return nil
}

func (s *personService) DeletePerson(ctx context.Context, personUUID string) error {
	s.log.Info(ctx, "Process Started")
	defer s.log.Info(ctx, "Process Finished")

	// Get existing person to validate ownership
	existingPerson, err := s.dm.Person().GetPersonByUUID(ctx, personUUID)
	if err != nil {
		if mysqlutils.SQLNotFound(err.Error()) {
			return resterrors.NewNotFoundError("person not found")
		}
		s.log.Errorw(ctx, "error getting person by UUID", logger.Err(err))
		return err
	}

	// Validate that the person belongs to a company owned by the logged user
	userID, err := s.authApp.GetLoggedUserID(ctx)
	if err != nil {
		return err
	}

	// Validate user has access to the person's company and get company for logging
	company, err := s.validateUserCompanyAccess(ctx, userID, existingPerson.CompanyID)
	if err != nil {
		return err
	}

	// Delete the person (soft delete)
	err = s.dm.Person().DeletePerson(ctx, existingPerson.ID)
	if err != nil {
		s.log.Errorw(ctx, "error deleting person", logger.Err(err))
		return err
	}

	s.log.Infow(ctx, "person deleted successfully",
		logger.String("person_uuid", personUUID),
		logger.String("person_name", existingPerson.Name),
		logger.Int64("company_id", existingPerson.CompanyID),
		logger.String("company_name", company.Name),
	)

	return nil
}

func (s *personService) SearchPeople(ctx context.Context, companyUUID string, search string) ([]entity.Person, error) {
	s.log.Info(ctx, "Process Started")
	defer s.log.Info(ctx, "Process Finished")

	// Get company by UUID and validate it belongs to the logged user
	company, err := s.dm.Company().GetCompanyByUUID(ctx, companyUUID)
	if err != nil {
		if mysqlutils.SQLNotFound(err.Error()) {
			return nil, resterrors.NewNotFoundError("company not found")
		}
		s.log.Errorw(ctx, "error getting company by UUID", logger.Err(err))
		return nil, err
	}

	// Validate that the company belongs to the logged user
	userID, err := s.authApp.GetLoggedUserID(ctx)
	if err != nil {
		return nil, err
	}

	if company.UserOwnerID != userID {
		s.log.Errorw(ctx, "user trying to search people from company they don't own",
			logger.Int64("company_owner_id", company.UserOwnerID),
			logger.Int64("logged_user_id", userID),
		)
		return nil, resterrors.NewUnauthorizedError("you don't have permission to search people from this company")
	}

	// Search people by company ID and search term
	people, err := s.dm.Person().SearchPeople(ctx, company.ID, search)
	if err != nil {
		s.log.Errorw(ctx, "error searching people", logger.Err(err))
		return nil, err
	}

	s.log.Infow(ctx, "people search completed successfully",
		logger.Int("people_count", len(people)),
		logger.String("search_term", search),
		logger.Int64("company_id", company.ID),
	)

	return people, nil
}