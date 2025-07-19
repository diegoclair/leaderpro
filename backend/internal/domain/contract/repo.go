package contract

import (
	"context"
	"time"

	"github.com/diegoclair/leaderpro/internal/application/dto"
	"github.com/diegoclair/leaderpro/internal/domain/entity"
)

// DataManager holds the methods that manipulates the main data.
type DataManager interface {
	WithTransaction(ctx context.Context, fn func(r DataManager) error) error
	User() UserRepo
	Company() CompanyRepo
	Person() PersonRepo
	Note() NoteRepo
	Auth() AuthRepo
}

type AuthRepo interface {
	CreateSession(ctx context.Context, session dto.Session) (sessionID int64, err error)
	GetSessionByUUID(ctx context.Context, sessionUUID string) (session dto.Session, err error)
	SetSessionAsBlocked(ctx context.Context, userID int64) (err error)
}

type UserRepo interface {
	CreateUser(ctx context.Context, user entity.User) (createdID int64, err error)
	GetUserByEmail(ctx context.Context, email string) (user entity.User, err error)
	GetUserByUUID(ctx context.Context, userUUID string) (user entity.User, err error)
	GetUserIDByUUID(ctx context.Context, userUUID string) (userID int64, err error)
	UpdateUser(ctx context.Context, userID int64, user entity.User) (err error)
	UpdateLastLogin(ctx context.Context, userID int64) (err error)
}

type CompanyRepo interface {
	CreateCompany(ctx context.Context, company entity.Company) (createdID int64, err error)
	GetCompanyByID(ctx context.Context, companyID int64) (company entity.Company, err error)
	GetCompanyByUUID(ctx context.Context, companyUUID string) (company entity.Company, err error)
	GetCompaniesByUser(ctx context.Context, userID int64) (companies []entity.Company, err error)
	UpdateCompany(ctx context.Context, companyID int64, company entity.Company) (err error)
	DeleteCompany(ctx context.Context, companyID int64) (err error)
}

type PersonRepo interface {
	CreatePerson(ctx context.Context, person entity.Person) (createdID int64, err error)
	GetPersonByUUID(ctx context.Context, personUUID string) (person entity.Person, err error)
	GetPersonsByCompany(ctx context.Context, companyID int64) (people []entity.Person, err error)
	GetPeopleCountByCompany(ctx context.Context, companyID int64) (count int64, err error)
	UpdatePerson(ctx context.Context, personID int64, person entity.Person) (err error)
	DeletePerson(ctx context.Context, personID int64) (err error)
	SearchPeople(ctx context.Context, companyID int64, search string) (people []entity.Person, err error)
}

type NoteRepo interface {
	CreateNote(ctx context.Context, note entity.Note) (createdID int64, err error)
	GetNoteByUUID(ctx context.Context, noteUUID string) (note entity.Note, err error)
	GetNotesByPerson(ctx context.Context, personID int64, take, skip int64) (notes []entity.Note, totalRecords int64, err error)
	UpdateNote(ctx context.Context, noteID int64, note entity.Note) (err error)
	DeleteNote(ctx context.Context, noteID int64) (err error)

	// Note mention methods
	CreateNoteMention(ctx context.Context, mention entity.NoteMention) (createdID int64, err error)
	GetMentionsByPerson(ctx context.Context, mentionedPersonID int64, take, skip int64) (mentions []entity.NoteMention, totalRecords int64, err error)
	GetPersonTimeline(ctx context.Context, personID int64, filters entity.TimelineFilters, take, skip int64) (timeline []entity.UnifiedTimelineEntry, totalRecords int64, err error)
	GetPersonMentions(ctx context.Context, mentionedPersonID int64, take, skip int64) (mentions []entity.MentionEntry, totalRecords int64, err error)
	DeleteMentionsByNote(ctx context.Context, noteID int64) (err error)

	// Dashboard stats methods (based on one-on-one notes)
	GetOneOnOnesCountThisMonth(ctx context.Context, companyID int64) (count int64, err error)
	GetAverageFrequencyDays(ctx context.Context, companyID int64) (avgDays float64, err error)
	GetLastMeetingDate(ctx context.Context, companyID int64) (lastDate *time.Time, err error)
}
