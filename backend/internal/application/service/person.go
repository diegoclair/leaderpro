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

func (s *personService) CreatePerson(ctx context.Context, person entity.Person, companyUUID string) error {
	s.log.Info(ctx, "Process Started")
	defer s.log.Info(ctx, "Process Finished")

	// Generate UUID for the person
	person.UUID = uuid.NewV4().String()

	// Get company by UUID and validate it belongs to the logged user
	company, err := s.dm.Company().GetCompanyByUUID(ctx, companyUUID)
	if err != nil {
		if mysqlutils.SQLNotFound(err.Error()) {
			return resterrors.NewNotFoundError("company not found")
		}
		s.log.Errorw(ctx, "error getting company by UUID", logger.Err(err))
		return err
	}

	// Validate that the company belongs to the logged user
	userID, err := s.authApp.GetLoggedUserID(ctx)
	if err != nil {
		return err
	}

	if company.UserOwnerID != userID {
		s.log.Errorw(ctx, "user trying to create person in company they don't own",
			logger.Int64("company_owner_id", company.UserOwnerID),
			logger.Int64("logged_user_id", userID),
		)
		return resterrors.NewUnauthorizedError("you don't have permission to add people to this company")
	}

	// Set the company ID and creator in the person entity
	person.CompanyID = company.ID
	person.CreatedBy = userID
	person.Active = true

	// Create the person in database
	personID, err := s.dm.Person().CreatePerson(ctx, person)
	if err != nil {
		s.log.Errorw(ctx, "error creating person", logger.Err(err))
		return err
	}

	s.log.Infow(ctx, "person created successfully",
		logger.Int64("person_id", personID),
		logger.String("person_name", person.Name),
		logger.Int64("company_id", company.ID),
		logger.String("company_name", company.Name),
	)

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