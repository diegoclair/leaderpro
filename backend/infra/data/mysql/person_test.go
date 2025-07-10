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

func validateTwoPeople(t *testing.T, personExpected entity.Person, personToCompare entity.Person) {
	require.NotZero(t, personToCompare.ID)
	require.Equal(t, personExpected.UUID, personToCompare.UUID)
	require.Equal(t, personExpected.CompanyID, personToCompare.CompanyID)
	require.Equal(t, personExpected.Name, personToCompare.Name)
	require.Equal(t, personExpected.Email, personToCompare.Email)
	require.Equal(t, personExpected.Position, personToCompare.Position)
	require.Equal(t, personExpected.Department, personToCompare.Department)
	require.Equal(t, personExpected.Phone, personToCompare.Phone)
	require.Equal(t, personExpected.IsManager, personToCompare.IsManager)
	require.Equal(t, personExpected.HasKids, personToCompare.HasKids)
	require.Equal(t, personExpected.Interests, personToCompare.Interests)
	require.Equal(t, personExpected.Personality, personToCompare.Personality)
	require.Equal(t, personExpected.CreatedBy, personToCompare.CreatedBy)
	require.Equal(t, personExpected.Active, personToCompare.Active)
}



func createRandomPerson(t *testing.T) entity.Person {
	// Create a user and company first
	user := createRandomUserForTests(t)
	company := createRandomCompanyForTests(t)

	birthday := time.Date(1990, 5, 15, 0, 0, 0, 0, time.UTC)
	startDate := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

	person := entity.Person{
		UUID:        uuid.NewV4().String(),
		CompanyID:   company.ID,
		Name:        "John Doe",
		Email:       "john.doe" + uuid.NewV4().String()[:8] + "@example.com", // unique email
		Position:    "Software Engineer",
		Department:  "Engineering",
		Phone:       "+1234567890",
		Birthday:    &birthday,
		StartDate:   &startDate,
		IsManager:   false,
		ManagerID:   nil,
		Notes:       "Great team player",
		HasKids:     true,
		Interests:   "Programming, Reading, Gaming",
		Personality: "Analytical, Creative",
		CreatedBy:   user.ID,
		Active:      true,
	}

	personID, err := testMysql.Person().CreatePerson(context.Background(), person)
	require.NoError(t, err)
	require.NotZero(t, personID)

	person.ID = personID
	return person
}

func TestCreatePerson(t *testing.T) {
	ctx := context.Background()
	user := createRandomUserForTests(t)
	company := createRandomCompanyForTests(t)

	birthday := time.Date(1985, 12, 25, 0, 0, 0, 0, time.UTC)
	startDate := time.Date(2021, 6, 1, 0, 0, 0, 0, time.UTC)

	person := entity.Person{
		UUID:        uuid.NewV4().String(),
		CompanyID:   company.ID,
		Name:        "Jane Smith",
		Email:       "jane.smith@example.com",
		Position:    "Product Manager",
		Department:  "Product",
		Phone:       "+9876543210",
		Birthday:    &birthday,
		StartDate:   &startDate,
		IsManager:   true,
		ManagerID:   nil,
		Notes:       "Strong leadership skills",
		HasKids:     false,
		Interests:   "UX Design, Strategy",
		Personality: "Strategic, Empathetic",
		CreatedBy:   user.ID,
		Active:      true,
	}

	personID, err := testMysql.Person().CreatePerson(ctx, person)
	require.NoError(t, err)
	require.NotZero(t, personID)
}

func TestGetPersonByUUID(t *testing.T) {
	ctx := context.Background()
	person := createRandomPerson(t)

	person2, err := testMysql.Person().GetPersonByUUID(ctx, person.UUID)
	require.NoError(t, err)
	require.NotEmpty(t, person2)
	validateTwoPeople(t, person, person2)
}

func TestGetPersonsByCompany(t *testing.T) {
	ctx := context.Background()
	user := createRandomUserForTests(t)
	company := createRandomCompanyForTests(t)

	// Create multiple people for the company
	person1 := entity.Person{
		UUID:       uuid.NewV4().String(),
		CompanyID:  company.ID,
		Name:       "Alice Johnson",
		Email:      "alice" + uuid.NewV4().String()[:8] + "@example.com",
		Position:   "Frontend Developer",
		Department: "Engineering",
		CreatedBy:  user.ID,
		Active:     true,
	}

	person2 := entity.Person{
		UUID:       uuid.NewV4().String(),
		CompanyID:  company.ID,
		Name:       "Bob Wilson",
		Email:      "bob" + uuid.NewV4().String()[:8] + "@example.com",
		Position:   "Backend Developer",
		Department: "Engineering",
		CreatedBy:  user.ID,
		Active:     true,
	}

	// Create people
	_, err := testMysql.Person().CreatePerson(ctx, person1)
	require.NoError(t, err)

	_, err = testMysql.Person().CreatePerson(ctx, person2)
	require.NoError(t, err)

	// Get people by company
	people, err := testMysql.Person().GetPersonsByCompany(ctx, company.ID)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(people), 2) // At least the 2 we created

	// Find our created people
	foundPerson1 := false
	foundPerson2 := false
	for _, person := range people {
		if person.UUID == person1.UUID {
			foundPerson1 = true
		}
		if person.UUID == person2.UUID {
			foundPerson2 = true
		}
	}
	require.True(t, foundPerson1)
	require.True(t, foundPerson2)
}

func TestUpdatePerson(t *testing.T) {
	ctx := context.Background()
	person := createRandomPerson(t)

	// Update person fields
	updatedPerson := entity.Person{
		Name:        "Updated Name",
		Email:       "updated@example.com",
		Position:    "Senior Software Engineer",
		Department:  "Platform",
		Phone:       "+9999999999",
		IsManager:   true,
		Notes:       "Updated notes",
		HasKids:     false,
		Interests:   "Updated interests",
		Personality: "Updated personality",
	}

	err := testMysql.Person().UpdatePerson(ctx, person.ID, updatedPerson)
	require.NoError(t, err)

	// Get updated person and validate
	retrievedPerson, err := testMysql.Person().GetPersonByUUID(ctx, person.UUID)
	require.NoError(t, err)
	require.Equal(t, updatedPerson.Name, retrievedPerson.Name)
	require.Equal(t, updatedPerson.Email, retrievedPerson.Email)
	require.Equal(t, updatedPerson.Position, retrievedPerson.Position)
	require.Equal(t, updatedPerson.Department, retrievedPerson.Department)
	require.Equal(t, updatedPerson.Phone, retrievedPerson.Phone)
	require.Equal(t, updatedPerson.IsManager, retrievedPerson.IsManager)
	require.Equal(t, updatedPerson.Notes, retrievedPerson.Notes)
	require.Equal(t, updatedPerson.HasKids, retrievedPerson.HasKids)
	require.Equal(t, updatedPerson.Interests, retrievedPerson.Interests)
	require.Equal(t, updatedPerson.Personality, retrievedPerson.Personality)
}

func TestDeletePerson(t *testing.T) {
	ctx := context.Background()
	person := createRandomPerson(t)

	err := testMysql.Person().DeletePerson(ctx, person.ID)
	require.NoError(t, err)

	// Try to get the deleted person - should return error (not found)
	_, err = testMysql.Person().GetPersonByUUID(ctx, person.UUID)
	require.Error(t, err)
}

func TestSearchPeople(t *testing.T) {
	ctx := context.Background()
	user := createRandomUserForTests(t)
	company := createRandomCompanyForTests(t)

	// Create a person with specific searchable data
	uniqueName := "SearchableTestPerson" + uuid.NewV4().String()[:8]
	person := entity.Person{
		UUID:       uuid.NewV4().String(),
		CompanyID:  company.ID,
		Name:       uniqueName,
		Email:      "searchable@example.com",
		Position:   "Search Engineer",
		Department: "Search Team",
		CreatedBy:  user.ID,
		Active:     true,
	}

	_, err := testMysql.Person().CreatePerson(ctx, person)
	require.NoError(t, err)

	// Test search by name
	people, err := testMysql.Person().SearchPeople(ctx, company.ID, "SearchableTest")
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(people), 1)

	found := false
	for _, p := range people {
		if p.UUID == person.UUID {
			found = true
			break
		}
	}
	require.True(t, found)

	// Test search by position
	people, err = testMysql.Person().SearchPeople(ctx, company.ID, "Search Engineer")
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(people), 1)

	// Test search by department
	people, err = testMysql.Person().SearchPeople(ctx, company.ID, "Search Team")
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(people), 1)
}

// Error tests with mocks
func TestCreatePersonErrorsWithMock(t *testing.T) {
	testForInsertErrorsWithMock(t, func(db *sql.DB) error {
		_, err := newPersonRepo(db).CreatePerson(context.Background(), entity.Person{})
		return err
	})
}

func TestGetPersonByUUIDErrorsWithMock(t *testing.T) {
	testForSelectErrorsWithMock(t, "person_id", func(db *sql.DB) error {
		_, err := newPersonRepo(db).GetPersonByUUID(context.Background(), "person-uuid")
		return err
	})
}

func TestGetPersonsByCompanyErrorsWithMock(t *testing.T) {
	testForSelectErrorsWithMock(t, "person_id", func(db *sql.DB) error {
		_, err := newPersonRepo(db).GetPersonsByCompany(context.Background(), 1)
		return err
	})
}

func TestUpdatePersonErrorsWithMock(t *testing.T) {
	testForUpdateDeleteErrorsWithMock(t, func(db *sql.DB) error {
		return newPersonRepo(db).UpdatePerson(context.Background(), 1, entity.Person{})
	})
}

func TestDeletePersonErrorsWithMock(t *testing.T) {
	testForUpdateDeleteErrorsWithMock(t, func(db *sql.DB) error {
		return newPersonRepo(db).DeletePerson(context.Background(), 1)
	})
}

func TestSearchPeopleErrorsWithMock(t *testing.T) {
	testForSelectErrorsWithMock(t, "person_id", func(db *sql.DB) error {
		_, err := newPersonRepo(db).SearchPeople(context.Background(), 1, "search")
		return err
	})
}