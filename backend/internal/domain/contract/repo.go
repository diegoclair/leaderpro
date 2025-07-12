package contract

import (
	"context"

	"github.com/diegoclair/leaderpro/internal/application/dto"
	"github.com/diegoclair/leaderpro/internal/domain/entity"
)

// DataManager holds the methods that manipulates the main data.
type DataManager interface {
	WithTransaction(ctx context.Context, fn func(r DataManager) error) error
	User() UserRepo
	Company() CompanyRepo
	Person() PersonRepo
	OneOnOne() OneOnOneRepo
	Feedback() FeedbackRepo
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
	GetCompanyByUUID(ctx context.Context, companyUUID string) (company entity.Company, err error)
	GetCompaniesByUser(ctx context.Context, userID int64) (companies []entity.Company, err error)
	UpdateCompany(ctx context.Context, companyID int64, company entity.Company) (err error)
	DeleteCompany(ctx context.Context, companyID int64) (err error)
}

type PersonRepo interface {
	CreatePerson(ctx context.Context, person entity.Person) (createdID int64, err error)
	GetPersonByUUID(ctx context.Context, personUUID string) (person entity.Person, err error)
	GetPersonsByCompany(ctx context.Context, companyID int64) (people []entity.Person, err error)
	UpdatePerson(ctx context.Context, personID int64, person entity.Person) (err error)
	DeletePerson(ctx context.Context, personID int64) (err error)
	SearchPeople(ctx context.Context, companyID int64, search string) (people []entity.Person, err error)
}

type OneOnOneRepo interface {
	CreateOneOnOne(ctx context.Context, oneOnOne entity.OneOnOne) (createdID int64, err error)
	GetOneOnOneByUUID(ctx context.Context, oneOnOneUUID string) (oneOnOne entity.OneOnOne, err error)
	GetOneOnOnesByPerson(ctx context.Context, personID int64, take, skip int64) (oneOnOnes []entity.OneOnOne, totalRecords int64, err error)
	GetOneOnOnesByManager(ctx context.Context, managerID int64, take, skip int64) (oneOnOnes []entity.OneOnOne, totalRecords int64, err error)
	UpdateOneOnOne(ctx context.Context, oneOnOneID int64, oneOnOne entity.OneOnOne) (err error)
	DeleteOneOnOne(ctx context.Context, oneOnOneID int64) (err error)
	GetUpcomingOneOnOnes(ctx context.Context, managerID int64) (oneOnOnes []entity.OneOnOne, err error)
	GetOverdueOneOnOnes(ctx context.Context, managerID int64) (oneOnOnes []entity.OneOnOne, err error)
}

type FeedbackRepo interface {
	CreateFeedback(ctx context.Context, feedback entity.Feedback) (createdID int64, err error)
	GetFeedbackByUUID(ctx context.Context, feedbackUUID string) (feedback entity.Feedback, err error)
	GetFeedbackByPerson(ctx context.Context, personID int64, take, skip int64) (feedback []entity.Feedback, totalRecords int64, err error)
	GetFeedbackByGiver(ctx context.Context, giverID int64, take, skip int64) (feedback []entity.Feedback, totalRecords int64, err error)
	UpdateFeedback(ctx context.Context, feedbackID int64, feedback entity.Feedback) (err error)
	DeleteFeedback(ctx context.Context, feedbackID int64) (err error)
	GetFeedbackSummary(ctx context.Context, personID int64, period string) (summary entity.FeedbackSummary, err error)
}
