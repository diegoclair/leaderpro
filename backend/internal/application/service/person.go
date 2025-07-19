package service

import (
	"context"
	"fmt"
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

// CreateNote creates a new note for a person
func (s *personService) CreateNote(ctx context.Context, note entity.Note, companyUUID string, personUUID string) (entity.Note, error) {
	s.log.Info(ctx, "Process Started")
	defer s.log.Info(ctx, "Process Finished")

	// Get and validate company ownership
	company, err := s.dm.Company().GetCompanyByUUID(ctx, companyUUID)
	if err != nil {
		if mysqlutils.SQLNotFound(err.Error()) {
			return note, resterrors.NewNotFoundError("company not found")
		}
		s.log.Errorw(ctx, "error getting company by UUID", logger.Err(err))
		return note, err
	}

	// Validate user ownership
	userID, err := s.authApp.GetLoggedUserID(ctx)
	if err != nil {
		return note, err
	}

	if company.UserOwnerID != userID {
		return note, resterrors.NewUnauthorizedError("you don't have permission to access this company")
	}

	// Get and validate person
	person, err := s.dm.Person().GetPersonByUUID(ctx, personUUID)
	if err != nil {
		if mysqlutils.SQLNotFound(err.Error()) {
			return note, resterrors.NewNotFoundError("person not found")
		}
		s.log.Errorw(ctx, "error getting person by UUID", logger.Err(err))
		return note, err
	}

	// Validate person belongs to this company
	if person.CompanyID != company.ID {
		return note, resterrors.NewBadRequestError("person does not belong to this company")
	}

	// Validate note data
	if err := s.validator.ValidateStruct(ctx, note); err != nil {
		return note, err
	}

	// Set note fields
	note.UUID = uuid.NewV4().String()
	note.CompanyID = company.ID
	note.PersonID = person.ID
	note.UserID = userID
	note.CreatedAt = time.Now()
	note.UpdatedAt = time.Now()

	// Create note
	noteID, err := s.dm.Note().CreateNote(ctx, note)
	if err != nil {
		s.log.Errorw(ctx, "error creating note", logger.Err(err))
		return note, err
	}

	note.ID = noteID

	// Process mentions from content
	mentionedUUIDs := note.ExtractMentionUUIDs()
	for _, mentionedUUID := range mentionedUUIDs {
		// Get mentioned person to validate they exist and belong to same company
		mentionedPerson, err := s.dm.Person().GetPersonByUUID(ctx, mentionedUUID)
		if err != nil {
			s.log.Warnw(ctx, "mentioned person not found, skipping mention",
				logger.String("mentioned_uuid", mentionedUUID),
				logger.Err(err),
			)
			continue
		}

		// Validate mentioned person belongs to same company
		if mentionedPerson.CompanyID != company.ID {
			s.log.Warnw(ctx, "mentioned person not in same company, skipping mention",
				logger.String("mentioned_uuid", mentionedUUID),
				logger.Int64("mentioned_person_company", mentionedPerson.CompanyID),
				logger.Int64("note_company", company.ID),
			)
			continue
		}

		// Create mention
		mention := entity.NoteMention{
			UUID:              uuid.NewV4().String(),
			NoteID:            noteID,
			MentionedPersonID: mentionedPerson.ID,
			SourcePersonID:    person.ID,
			FullContent:       note.Content,
			CreatedAt:         time.Now(),
		}

		_, err = s.dm.Note().CreateNoteMention(ctx, mention)
		if err != nil {
			s.log.Errorw(ctx, "error creating note mention", 
				logger.Err(err),
				logger.String("mentioned_person_uuid", mentionedUUID),
			)
			// Continue processing other mentions even if one fails
		}
	}

	s.log.Infow(ctx, "note created successfully",
		logger.String("note_uuid", note.UUID),
		logger.String("note_type", note.Type),
		logger.String("person_uuid", personUUID),
		logger.Int("mentions_count", len(mentionedUUIDs)),
	)

	return note, nil
}

// GetPersonTimeline gets the complete timeline (notes + mentions) for a person
func (s *personService) GetPersonTimeline(ctx context.Context, personUUID string, filters entity.TimelineFilters, take, skip int64) ([]entity.UnifiedTimelineEntry, int64, error) {
	s.log.Info(ctx, "Process Started")
	defer s.log.Info(ctx, "Process Finished")

	// Get person and validate ownership
	person, err := s.dm.Person().GetPersonByUUID(ctx, personUUID)
	if err != nil {
		if mysqlutils.SQLNotFound(err.Error()) {
			return nil, 0, resterrors.NewNotFoundError("person not found")
		}
		s.log.Errorw(ctx, "error getting person by UUID", logger.Err(err))
		return nil, 0, err
	}

	// Validate user has access to this person's company
	userID, err := s.authApp.GetLoggedUserID(ctx)
	if err != nil {
		return nil, 0, err
	}

	_, err = s.validateUserCompanyAccess(ctx, userID, person.CompanyID)
	if err != nil {
		return nil, 0, err
	}

	// Get unified timeline from repository
	timeline, totalRecords, err := s.dm.Note().GetPersonTimeline(ctx, person.ID, filters, take, skip)
	if err != nil {
		s.log.Errorw(ctx, "error getting person timeline", logger.Err(err))
		return nil, 0, err
	}

	s.log.Infow(ctx, "unified timeline retrieved successfully",
		logger.String("person_uuid", personUUID),
		logger.Int64("total_records", totalRecords),
		logger.Int("returned_records", len(timeline)),
		logger.String("search_query", filters.SearchQuery),
		logger.String("types", fmt.Sprintf("%v", filters.Types)),
		logger.String("period", filters.Period),
	)

	return timeline, totalRecords, nil
}

// GetPersonMentions gets notes where this person was mentioned (feedbacks received)
func (s *personService) GetPersonMentions(ctx context.Context, personUUID string, take, skip int64) ([]entity.MentionEntry, int64, error) {
	s.log.Info(ctx, "Process Started")
	defer s.log.Info(ctx, "Process Finished")

	// Get person and validate ownership
	person, err := s.dm.Person().GetPersonByUUID(ctx, personUUID)
	if err != nil {
		if mysqlutils.SQLNotFound(err.Error()) {
			return nil, 0, resterrors.NewNotFoundError("person not found")
		}
		s.log.Errorw(ctx, "error getting person by UUID", logger.Err(err))
		return nil, 0, err
	}

	// Validate user has access to this person's company
	userID, err := s.authApp.GetLoggedUserID(ctx)
	if err != nil {
		return nil, 0, err
	}

	_, err = s.validateUserCompanyAccess(ctx, userID, person.CompanyID)
	if err != nil {
		return nil, 0, err
	}

	// Get mentions from repository
	mentions, totalRecords, err := s.dm.Note().GetPersonMentions(ctx, person.ID, take, skip)
	if err != nil {
		s.log.Errorw(ctx, "error getting person mentions", logger.Err(err))
		return nil, 0, err
	}

	s.log.Infow(ctx, "mentions retrieved successfully",
		logger.String("person_uuid", personUUID),
		logger.Int64("total_records", totalRecords),
		logger.Int("returned_records", len(mentions)),
	)

	return mentions, totalRecords, nil
}

// UpdateNote updates an existing note
func (s *personService) UpdateNote(ctx context.Context, noteUUID string, updatedNote entity.Note) error {
	s.log.Info(ctx, "Process Started")
	defer s.log.Info(ctx, "Process Finished")

	// Get existing note
	existingNote, err := s.dm.Note().GetNoteByUUID(ctx, noteUUID)
	if err != nil {
		if mysqlutils.SQLNotFound(err.Error()) {
			return resterrors.NewNotFoundError("note not found")
		}
		s.log.Errorw(ctx, "error getting note by UUID", logger.Err(err))
		return err
	}

	// Validate user ownership through company
	userID, err := s.authApp.GetLoggedUserID(ctx)
	if err != nil {
		return err
	}

	_, err = s.validateUserCompanyAccess(ctx, userID, existingNote.CompanyID)
	if err != nil {
		return err
	}

	// Validate note data
	if err := s.validator.ValidateStruct(ctx, updatedNote); err != nil {
		return err
	}

	// Preserve original fields and update only allowed fields
	updatedNote.ID = existingNote.ID
	updatedNote.UUID = existingNote.UUID
	updatedNote.CompanyID = existingNote.CompanyID
	updatedNote.PersonID = existingNote.PersonID
	updatedNote.UserID = existingNote.UserID
	updatedNote.CreatedAt = existingNote.CreatedAt
	updatedNote.UpdatedAt = time.Now()

	// Update note
	err = s.dm.Note().UpdateNote(ctx, existingNote.ID, updatedNote)
	if err != nil {
		s.log.Errorw(ctx, "error updating note", logger.Err(err))
		return err
	}

	// Handle mention updates: delete old mentions and create new ones
	err = s.dm.Note().DeleteMentionsByNote(ctx, existingNote.ID)
	if err != nil {
		s.log.Errorw(ctx, "error deleting old mentions", logger.Err(err))
		// Continue with update even if mention deletion fails
	}

	// Process new mentions
	mentionedUUIDs := updatedNote.ExtractMentionUUIDs()
	for _, mentionedUUID := range mentionedUUIDs {
		mentionedPerson, err := s.dm.Person().GetPersonByUUID(ctx, mentionedUUID)
		if err != nil {
			s.log.Warnw(ctx, "mentioned person not found during update, skipping mention",
				logger.String("mentioned_uuid", mentionedUUID),
				logger.Err(err),
			)
			continue
		}

		if mentionedPerson.CompanyID != existingNote.CompanyID {
			s.log.Warnw(ctx, "mentioned person not in same company during update, skipping mention",
				logger.String("mentioned_uuid", mentionedUUID),
			)
			continue
		}

		mention := entity.NoteMention{
			UUID:              uuid.NewV4().String(),
			NoteID:            existingNote.ID,
			MentionedPersonID: mentionedPerson.ID,
			SourcePersonID:    existingNote.PersonID,
			FullContent:       updatedNote.Content,
			CreatedAt:         time.Now(),
		}

		_, err = s.dm.Note().CreateNoteMention(ctx, mention)
		if err != nil {
			s.log.Errorw(ctx, "error creating mention during update", 
				logger.Err(err),
				logger.String("mentioned_person_uuid", mentionedUUID),
			)
		}
	}

	s.log.Infow(ctx, "note updated successfully",
		logger.String("note_uuid", noteUUID),
		logger.String("note_type", updatedNote.Type),
	)

	return nil
}

// DeleteNote deletes a note and its mentions
func (s *personService) DeleteNote(ctx context.Context, noteUUID string) error {
	s.log.Info(ctx, "Process Started")
	defer s.log.Info(ctx, "Process Finished")

	// Get existing note
	existingNote, err := s.dm.Note().GetNoteByUUID(ctx, noteUUID)
	if err != nil {
		if mysqlutils.SQLNotFound(err.Error()) {
			return resterrors.NewNotFoundError("note not found")
		}
		s.log.Errorw(ctx, "error getting note by UUID", logger.Err(err))
		return err
	}

	// Validate user ownership
	userID, err := s.authApp.GetLoggedUserID(ctx)
	if err != nil {
		return err
	}

	_, err = s.validateUserCompanyAccess(ctx, userID, existingNote.CompanyID)
	if err != nil {
		return err
	}

	// Delete note (this should cascade delete mentions via foreign key)
	err = s.dm.Note().DeleteNote(ctx, existingNote.ID)
	if err != nil {
		s.log.Errorw(ctx, "error deleting note", logger.Err(err))
		return err
	}

	s.log.Infow(ctx, "note deleted successfully",
		logger.String("note_uuid", noteUUID),
		logger.String("note_type", existingNote.Type),
	)

	return nil
}