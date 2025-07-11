package mysql

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/diegoclair/leaderpro/internal/domain/entity"
	"github.com/stretchr/testify/require"
	"github.com/twinj/uuid"
)

func validateTwoUsers(t *testing.T, userExpected entity.User, userToCompare entity.User) {
	require.NotZero(t, userToCompare.ID)
	require.Equal(t, userExpected.UUID, userToCompare.UUID)
	require.Equal(t, userExpected.Email, userToCompare.Email)
	require.Equal(t, userExpected.Name, userToCompare.Name)
	require.Equal(t, userExpected.Phone, userToCompare.Phone)
	require.Equal(t, userExpected.ProfilePhoto, userToCompare.ProfilePhoto)
	require.Equal(t, userExpected.Plan, userToCompare.Plan)
	require.Equal(t, userExpected.Active, userToCompare.Active)
	require.Equal(t, userExpected.EmailVerified, userToCompare.EmailVerified)
}

func createRandomUser(t *testing.T) entity.User {
	user := entity.User{
		UUID:          uuid.NewV4().String(),
		Email:         "test" + uuid.NewV4().String()[:8] + "@example.com", // unique email
		Name:          "Test User",
		Password:      "hashedpassword",
		Phone:         "+1234567890",
		ProfilePhoto:  "https://example.com/photo.jpg",
		Plan:          "basic",
		Active:        true,
		EmailVerified: false,
	}

	userID, err := testMysql.User().CreateUser(context.Background(), user)
	require.NoError(t, err)
	require.NotZero(t, userID)

	user.ID = userID
	return user
}

// createRandomUserForTests creates a simple user for use in other tests
func createRandomUserForTests(t *testing.T) entity.User {
	user := entity.User{
		UUID:     uuid.NewV4().String(),
		Email:    "testuser" + uuid.NewV4().String()[:8] + "@example.com",
		Name:     "Test User",
		Password: "hashedpassword",
		Phone:    "+1234567890",
		Active:   true,
	}

	userID, err := testMysql.User().CreateUser(context.Background(), user)
	require.NoError(t, err)
	require.NotZero(t, userID)

	user.ID = userID
	return user
}

func createRandomCompanyForTests(t *testing.T) entity.Company {
	user := createRandomUserForTests(t)

	company := entity.Company{
		UUID:        uuid.NewV4().String(),
		Name:        "Test Company",
		Description: "A test company for tests",
		Industry:    "Technology",
		Size:        "small",
		CreatedBy:   user.ID,
		Active:      true,
	}

	companyID, err := testMysql.Company().CreateCompany(context.Background(), company)
	require.NoError(t, err)
	require.NotZero(t, companyID)

	company.ID = companyID
	return company
}

func createRandomPersonForTests(t *testing.T) entity.Person {
	user := createRandomUserForTests(t)
	company := createRandomCompanyForTests(t)

	birthday := time.Now().AddDate(-30, 0, 0) // 30 years ago
	startDate := time.Now().AddDate(-2, 0, 0) // 2 years ago

	person := entity.Person{
		UUID:       uuid.NewV4().String(),
		CompanyID:  company.ID,
		Name:       "Test Person",
		Email:      "testperson" + uuid.NewV4().String()[:8] + "@example.com",
		Phone:      "+1987654321",
		Position:   "Software Engineer",
		Department: "Engineering",
		Birthday:   &birthday,
		StartDate:  &startDate,
		ManagerID:  nil, // No manager for test person
		Notes:      "Test person",
		CreatedBy:  user.ID,
		Active:     true,
	}

	personID, err := testMysql.Person().CreatePerson(context.Background(), person)
	require.NoError(t, err)
	require.NotZero(t, personID)

	person.ID = personID
	return person
}

func TestCreateUser(t *testing.T) {
	ctx := context.Background()

	user := entity.User{
		UUID:          uuid.NewV4().String(),
		Email:         "test@example.com",
		Name:          "Test User",
		Password:      "hashedpassword",
		Phone:         "+1234567890",
		ProfilePhoto:  "https://example.com/photo.jpg",
		Plan:          "basic",
		TrialEndsAt:   nil,
		SubscribedAt:  nil,
		Active:        true,
		EmailVerified: false,
	}

	userID, err := testMysql.User().CreateUser(ctx, user)
	require.NoError(t, err)
	require.NotZero(t, userID)

	user.ID = userID
}

func TestGetUserByEmail(t *testing.T) {
	ctx := context.Background()
	user := createRandomUser(t)

	user2, err := testMysql.User().GetUserByEmail(ctx, user.Email)
	require.NoError(t, err)
	require.NotEmpty(t, user2)
	validateTwoUsers(t, user, user2)
}

func TestGetUserByUUID(t *testing.T) {
	ctx := context.Background()
	user := createRandomUser(t)

	user2, err := testMysql.User().GetUserByUUID(ctx, user.UUID)
	require.NoError(t, err)
	require.NotEmpty(t, user2)
	validateTwoUsers(t, user, user2)
}

func TestGetUserIDByUUID(t *testing.T) {
	ctx := context.Background()
	user := createRandomUser(t)

	userID, err := testMysql.User().GetUserIDByUUID(ctx, user.UUID)
	require.NoError(t, err)
	require.Equal(t, user.ID, userID)
}

func TestUpdateUser(t *testing.T) {
	ctx := context.Background()
	user := createRandomUser(t)

	// Update user fields
	user.Name = "Updated Name"
	user.Phone = "+9876543210"
	user.ProfilePhoto = "https://example.com/updated-photo.jpg"
	user.Plan = "premium"
	user.EmailVerified = true

	err := testMysql.User().UpdateUser(ctx, user.ID, user)
	require.NoError(t, err)

	// Get updated user and validate
	updatedUser, err := testMysql.User().GetUserByUUID(ctx, user.UUID)
	require.NoError(t, err)
	require.Equal(t, user.Name, updatedUser.Name)
	require.Equal(t, user.Phone, updatedUser.Phone)
	require.Equal(t, user.ProfilePhoto, updatedUser.ProfilePhoto)
	require.Equal(t, user.Plan, updatedUser.Plan)
	require.Equal(t, user.EmailVerified, updatedUser.EmailVerified)
}

func TestUpdateLastLogin(t *testing.T) {
	ctx := context.Background()
	user := createRandomUser(t)

	err := testMysql.User().UpdateLastLogin(ctx, user.ID)
	require.NoError(t, err)

	// Get updated user and validate last login was set
	updatedUser, err := testMysql.User().GetUserByUUID(ctx, user.UUID)
	require.NoError(t, err)
	require.NotNil(t, updatedUser.LastLoginAt)
	require.WithinDuration(t, time.Now(), *updatedUser.LastLoginAt, 5*time.Second)
}

// Error tests with mocks
func TestCreateUserErrorsWithMock(t *testing.T) {
	testForInsertErrorsWithMock(t, func(db *sql.DB) error {
		_, err := newUserRepo(db).CreateUser(context.Background(), entity.User{})
		return err
	})
}

func TestGetUserByEmailErrorsWithMock(t *testing.T) {
	testForSelectErrorsWithMock(t, "user_id", func(db *sql.DB) error {
		_, err := newUserRepo(db).GetUserByEmail(context.Background(), "test@example.com")
		return err
	})
}

func TestGetUserByUUIDErrorsWithMock(t *testing.T) {
	testForSelectErrorsWithMock(t, "user_id", func(db *sql.DB) error {
		_, err := newUserRepo(db).GetUserByUUID(context.Background(), "user-uuid")
		return err
	})
}

func TestGetUserIDByUUIDErrorsWithMock(t *testing.T) {
	testForSelectErrorsWithMock(t, "user_id", func(db *sql.DB) error {
		_, err := newUserRepo(db).GetUserIDByUUID(context.Background(), "user-uuid")
		return err
	})
}

func TestUpdateUserErrorsWithMock(t *testing.T) {
	testForUpdateDeleteErrorsWithMock(t, func(db *sql.DB) error {
		return newUserRepo(db).UpdateUser(context.Background(), 1, entity.User{})
	})
}

func TestUpdateLastLoginErrorsWithMock(t *testing.T) {
	testForUpdateDeleteErrorsWithMock(t, func(db *sql.DB) error {
		return newUserRepo(db).UpdateLastLogin(context.Background(), 1)
	})
}
