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
	
	// User Preferences
	GetUserPreferences(ctx context.Context) (preferences entity.UserPreferences, err error)
	UpdateUserPreferences(ctx context.Context, preferences entity.UserPreferences) (updatedPreferences entity.UserPreferences, err error)
}

type AuthApp interface {
	Login(ctx context.Context, input dto.LoginInput) (user entity.User, err error)
	CreateSession(ctx context.Context, session dto.Session) (err error)
	GetSessionByUUID(ctx context.Context, sessionUUID string) (session dto.Session, err error)
	Logout(ctx context.Context, accessToken string) (err error)
	GetLoggedUserID(ctx context.Context) (userID int64, err error)
	GetCompanyFromContext(ctx context.Context) (companyUUID string, err error)
}

type CompanyApp interface {
	CreateCompany(ctx context.Context, company entity.Company) (createdCompany entity.Company, err error)
	GetCompanyByUUID(ctx context.Context, companyUUID string) (company entity.Company, err error)
	GetLoggedUserCompany(ctx context.Context) (company entity.Company, err error)
	GetUserCompanies(ctx context.Context) (companies []entity.Company, err error)
	UpdateCompany(ctx context.Context, companyUUID string, company entity.Company) (err error)
	UpdateLoggedUserCompany(ctx context.Context, company entity.Company) (err error)
	DeleteCompany(ctx context.Context, companyUUID string) (err error)
	DeleteLoggedUserCompany(ctx context.Context) (err error)
	ValidateCompanyOwnership(ctx context.Context, companyUUID string, userUUID string) (err error)
}

type PersonApp interface {
	CreatePerson(ctx context.Context, person entity.Person) (createdPerson entity.Person, err error)
	GetPersonByUUID(ctx context.Context, personUUID string) (person entity.Person, err error)
	GetCompanyPeople(ctx context.Context) (people []entity.Person, err error)
	UpdatePerson(ctx context.Context, personUUID string, person entity.Person) (err error)
	DeletePerson(ctx context.Context, personUUID string) (err error)
	SearchPeople(ctx context.Context, search string) (people []entity.Person, err error)
	
	// Note management methods
	CreateNote(ctx context.Context, note entity.Note, personUUID string) (createdNote entity.Note, err error)
	GetPersonTimeline(ctx context.Context, personUUID string, filters entity.TimelineFilters, take, skip int64) (timeline []entity.UnifiedTimelineEntry, totalRecords int64, err error)
	GetPersonMentions(ctx context.Context, personUUID string, take, skip int64) (mentions []entity.MentionEntry, totalRecords int64, err error)
	UpdateNote(ctx context.Context, noteUUID string, note entity.Note) (err error)
	DeleteNote(ctx context.Context, noteUUID string) (err error)
}

type DashboardApp interface {
	GetDashboardData(ctx context.Context, companyUUID string) (dashboard entity.Dashboard, err error)
}

type AIApp interface {
	// ChatWithLeadershipCoach performs chat with leadership context
	ChatWithLeadershipCoach(ctx context.Context, req entity.ChatRequest) (entity.ChatResponse, error)
	
	// ExtractAttributesFromNote extracts attributes from a note
	ExtractAttributesFromNote(ctx context.Context, noteID int64) (entity.AttributesResponse, error)
	
	// GetPersonContext retrieves complete person context for AI
	GetPersonContext(ctx context.Context, personID int64) (entity.PersonAIContext, error)
	
	// SendFeedback records feedback about an AI response
	SendFeedback(ctx context.Context, usageID int64, feedback string, comment string) error
	
	// GetUsageReport returns AI usage report
	GetUsageReport(ctx context.Context, period string) (entity.AIUsageReport, error)
}
