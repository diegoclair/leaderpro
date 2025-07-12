package contract

import (
	"context"

	"github.com/diegoclair/leaderpro/internal/application/dto"
	"github.com/diegoclair/leaderpro/internal/domain/entity"
)

type UserApp interface {
	CreateUser(ctx context.Context, user entity.User) (createdUser entity.User, err error)
	GetUserByEmail(ctx context.Context, email string) (user entity.User, err error)
	GetUserByUUID(ctx context.Context, userUUID string) (user entity.User, err error)
	GetLoggedUser(ctx context.Context) (user entity.User, err error)
	GetLoggedUserID(ctx context.Context) (userID int64, err error)
	GetProfile(ctx context.Context) (user entity.User, err error)
	UpdateProfile(ctx context.Context, user entity.User) (updatedUser entity.User, err error)
	UpdateUser(ctx context.Context, userUUID string, user entity.User) (err error)
}

type AuthApp interface {
	Login(ctx context.Context, input dto.LoginInput) (user entity.User, err error)
	CreateSession(ctx context.Context, session dto.Session) (err error)
	GetSessionByUUID(ctx context.Context, sessionUUID string) (session dto.Session, err error)
	Logout(ctx context.Context, accessToken string) (err error)
	GetLoggedUserID(ctx context.Context) (userID int64, err error)
}

type CompanyApp interface {
	CreateCompany(ctx context.Context, company entity.Company) (createdCompany entity.Company, err error)
	GetCompanyByUUID(ctx context.Context, companyUUID string) (company entity.Company, err error)
	GetUserCompanies(ctx context.Context) (companies []entity.Company, err error)
	UpdateCompany(ctx context.Context, companyUUID string, company entity.Company) (err error)
	DeleteCompany(ctx context.Context, companyUUID string) (err error)
}

type PersonApp interface {
	CreatePerson(ctx context.Context, person entity.Person, companyUUID string) (createdPerson entity.Person, err error)
	GetPersonByUUID(ctx context.Context, personUUID string) (person entity.Person, err error)
	GetCompanyPeople(ctx context.Context, companyUUID string) (people []entity.Person, err error)
	UpdatePerson(ctx context.Context, personUUID string, person entity.Person) (err error)
	DeletePerson(ctx context.Context, personUUID string) (err error)
	SearchPeople(ctx context.Context, companyUUID string, search string) (people []entity.Person, err error)
}

type OneOnOneApp interface {
	CreateOneOnOne(ctx context.Context, oneOnOne entity.OneOnOne) (err error)
	GetOneOnOneByUUID(ctx context.Context, oneOnOneUUID string) (oneOnOne entity.OneOnOne, err error)
	GetPersonOneOnOnes(ctx context.Context, personUUID string, take, skip int64) (oneOnOnes []entity.OneOnOne, totalRecords int64, err error)
	GetManagerOneOnOnes(ctx context.Context, take, skip int64) (oneOnOnes []entity.OneOnOne, totalRecords int64, err error)
	UpdateOneOnOne(ctx context.Context, oneOnOneUUID string, oneOnOne entity.OneOnOne) (err error)
	DeleteOneOnOne(ctx context.Context, oneOnOneUUID string) (err error)
	GetUpcomingOneOnOnes(ctx context.Context) (oneOnOnes []entity.OneOnOne, err error)
	GetOverdueOneOnOnes(ctx context.Context) (oneOnOnes []entity.OneOnOne, err error)
}

type FeedbackApp interface {
	CreateFeedback(ctx context.Context, feedback entity.Feedback) (err error)
	GetFeedbackByUUID(ctx context.Context, feedbackUUID string) (feedback entity.Feedback, err error)
	GetPersonFeedback(ctx context.Context, personUUID string, take, skip int64) (feedback []entity.Feedback, totalRecords int64, err error)
	UpdateFeedback(ctx context.Context, feedbackUUID string, feedback entity.Feedback) (err error)
	DeleteFeedback(ctx context.Context, feedbackUUID string) (err error)
	GetFeedbackSummary(ctx context.Context, personUUID string, period string) (summary entity.FeedbackSummary, err error)
}
