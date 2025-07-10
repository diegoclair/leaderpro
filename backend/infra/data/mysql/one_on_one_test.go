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

func validateTwoOneOnOnes(t *testing.T, oneOnOneExpected entity.OneOnOne, oneOnOneToCompare entity.OneOnOne) {
	require.NotZero(t, oneOnOneToCompare.ID)
	require.Equal(t, oneOnOneExpected.UUID, oneOnOneToCompare.UUID)
	require.Equal(t, oneOnOneExpected.CompanyID, oneOnOneToCompare.CompanyID)
	require.Equal(t, oneOnOneExpected.PersonID, oneOnOneToCompare.PersonID)
	require.Equal(t, oneOnOneExpected.ManagerID, oneOnOneToCompare.ManagerID)
	require.WithinDuration(t, oneOnOneExpected.ScheduledDate, oneOnOneToCompare.ScheduledDate, time.Second)
	require.Equal(t, oneOnOneExpected.Duration, oneOnOneToCompare.Duration)
	require.Equal(t, oneOnOneExpected.Location, oneOnOneToCompare.Location)
	require.Equal(t, oneOnOneExpected.Status, oneOnOneToCompare.Status)
	require.Equal(t, oneOnOneExpected.Agenda, oneOnOneToCompare.Agenda)
	require.Equal(t, oneOnOneExpected.DiscussionNotes, oneOnOneToCompare.DiscussionNotes)
	require.Equal(t, oneOnOneExpected.ActionItems, oneOnOneToCompare.ActionItems)
	require.Equal(t, oneOnOneExpected.PrivateNotes, oneOnOneToCompare.PrivateNotes)
	require.Equal(t, oneOnOneExpected.AIContext, oneOnOneToCompare.AIContext)
	require.Equal(t, oneOnOneExpected.AISuggestions, oneOnOneToCompare.AISuggestions)
}




func createRandomOneOnOne(t *testing.T) entity.OneOnOne {
	// Create user, company, and person first
	user := createRandomUserForTests(t)
	company := createRandomCompanyForTests(t)
	person := createRandomPersonForTests(t)

	scheduledDate := time.Now().Add(24 * time.Hour) // Tomorrow

	oneOnOne := entity.OneOnOne{
		UUID:            uuid.NewV4().String(),
		CompanyID:       company.ID,
		PersonID:        person.ID,
		ManagerID:       user.ID,
		ScheduledDate:   scheduledDate,
		ActualDate:      nil,
		Duration:        30, // 30 minutes
		Location:        "Conference Room A",
		Status:          "scheduled",
		Agenda:          "Weekly check-in",
		DiscussionNotes: "",
		ActionItems:     "",
		PrivateNotes:    "Prepare questions about career goals",
		AIContext:       `{"recent_projects": ["Feature X"], "mood": "positive"}`,
		AISuggestions:   `["Ask about Feature X progress", "Discuss career development"]`,
	}

	oneOnOneID, err := testMysql.OneOnOne().CreateOneOnOne(context.Background(), oneOnOne)
	require.NoError(t, err)
	require.NotZero(t, oneOnOneID)

	oneOnOne.ID = oneOnOneID
	return oneOnOne
}

func TestCreateOneOnOne(t *testing.T) {
	ctx := context.Background()
	user := createRandomUserForTests(t)
	company := createRandomCompanyForTests(t)
	person := createRandomPersonForTests(t)

	scheduledDate := time.Now().Add(48 * time.Hour) // Day after tomorrow

	oneOnOne := entity.OneOnOne{
		UUID:            uuid.NewV4().String(),
		CompanyID:       company.ID,
		PersonID:        person.ID,
		ManagerID:       user.ID,
		ScheduledDate:   scheduledDate,
		Duration:        45,
		Location:        "Office",
		Status:          "scheduled",
		Agenda:          "Monthly review",
		DiscussionNotes: "",
		ActionItems:     "",
		PrivateNotes:    "Focus on performance feedback",
		AIContext:       `{"context": "monthly_review"}`,
		AISuggestions:   `["Review last month's goals"]`,
	}

	oneOnOneID, err := testMysql.OneOnOne().CreateOneOnOne(ctx, oneOnOne)
	require.NoError(t, err)
	require.NotZero(t, oneOnOneID)
}

func TestGetOneOnOneByUUID(t *testing.T) {
	ctx := context.Background()
	oneOnOne := createRandomOneOnOne(t)

	oneOnOne2, err := testMysql.OneOnOne().GetOneOnOneByUUID(ctx, oneOnOne.UUID)
	require.NoError(t, err)
	require.NotEmpty(t, oneOnOne2)
	validateTwoOneOnOnes(t, oneOnOne, oneOnOne2)
}

func TestGetOneOnOnesByPerson(t *testing.T) {
	ctx := context.Background()
	oneOnOne := createRandomOneOnOne(t)

	oneOnOnes, totalRecords, err := testMysql.OneOnOne().GetOneOnOnesByPerson(ctx, oneOnOne.PersonID, 10, 0)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(oneOnOnes), 1)
	require.GreaterOrEqual(t, totalRecords, int64(1))

	// Find our created one-on-one
	found := false
	for _, o := range oneOnOnes {
		if o.UUID == oneOnOne.UUID {
			found = true
			validateTwoOneOnOnes(t, oneOnOne, o)
			break
		}
	}
	require.True(t, found)
}

func TestGetOneOnOnesByManager(t *testing.T) {
	ctx := context.Background()
	oneOnOne := createRandomOneOnOne(t)

	oneOnOnes, totalRecords, err := testMysql.OneOnOne().GetOneOnOnesByManager(ctx, oneOnOne.ManagerID, 10, 0)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(oneOnOnes), 1)
	require.GreaterOrEqual(t, totalRecords, int64(1))

	// Find our created one-on-one
	found := false
	for _, o := range oneOnOnes {
		if o.UUID == oneOnOne.UUID {
			found = true
			validateTwoOneOnOnes(t, oneOnOne, o)
			break
		}
	}
	require.True(t, found)
}

func TestUpdateOneOnOne(t *testing.T) {
	ctx := context.Background()
	oneOnOne := createRandomOneOnOne(t)

	// Update one-on-one fields
	actualDate := time.Now()
	updatedOneOnOne := entity.OneOnOne{
		ScheduledDate:   oneOnOne.ScheduledDate,
		ActualDate:      &actualDate,
		Duration:        60, // Updated to 60 minutes
		Location:        "Virtual Meeting",
		Status:          "completed",
		Agenda:          "Updated agenda",
		DiscussionNotes: "Great discussion about project progress",
		ActionItems:     "1. Complete feature Y by next week\n2. Review documentation",
		PrivateNotes:    "Employee seems motivated and engaged",
		AIContext:       `{"updated_context": "completed_meeting"}`,
		AISuggestions:   `["Follow up on action items"]`,
	}

	err := testMysql.OneOnOne().UpdateOneOnOne(ctx, oneOnOne.ID, updatedOneOnOne)
	require.NoError(t, err)

	// Get updated one-on-one and validate
	retrievedOneOnOne, err := testMysql.OneOnOne().GetOneOnOneByUUID(ctx, oneOnOne.UUID)
	require.NoError(t, err)
	require.Equal(t, updatedOneOnOne.Duration, retrievedOneOnOne.Duration)
	require.Equal(t, updatedOneOnOne.Location, retrievedOneOnOne.Location)
	require.Equal(t, updatedOneOnOne.Status, retrievedOneOnOne.Status)
	require.Equal(t, updatedOneOnOne.Agenda, retrievedOneOnOne.Agenda)
	require.Equal(t, updatedOneOnOne.DiscussionNotes, retrievedOneOnOne.DiscussionNotes)
	require.Equal(t, updatedOneOnOne.ActionItems, retrievedOneOnOne.ActionItems)
	require.Equal(t, updatedOneOnOne.PrivateNotes, retrievedOneOnOne.PrivateNotes)
	require.Equal(t, updatedOneOnOne.AIContext, retrievedOneOnOne.AIContext)
	require.Equal(t, updatedOneOnOne.AISuggestions, retrievedOneOnOne.AISuggestions)
	require.NotNil(t, retrievedOneOnOne.ActualDate)
	require.NotNil(t, retrievedOneOnOne.CompletedAt) // Should be set when status = completed
}

func TestDeleteOneOnOne(t *testing.T) {
	ctx := context.Background()
	oneOnOne := createRandomOneOnOne(t)

	err := testMysql.OneOnOne().DeleteOneOnOne(ctx, oneOnOne.ID)
	require.NoError(t, err)

	// Try to get the deleted one-on-one - should return error (not found)
	_, err = testMysql.OneOnOne().GetOneOnOneByUUID(ctx, oneOnOne.UUID)
	require.Error(t, err)
}

func TestGetUpcomingOneOnOnes(t *testing.T) {
	ctx := context.Background()
	user := createRandomUserForTests(t)
	company := createRandomCompanyForTests(t)
	person := createRandomPersonForTests(t)

	// Create an upcoming one-on-one
	futureDate := time.Now().Add(72 * time.Hour) // 3 days from now
	upcomingOneOnOne := entity.OneOnOne{
		UUID:          uuid.NewV4().String(),
		CompanyID:     company.ID,
		PersonID:      person.ID,
		ManagerID:     user.ID,
		ScheduledDate: futureDate,
		Duration:      30,
		Status:        "scheduled",
		Agenda:        "Upcoming meeting",
	}

	_, err := testMysql.OneOnOne().CreateOneOnOne(ctx, upcomingOneOnOne)
	require.NoError(t, err)

	// Get upcoming one-on-ones
	oneOnOnes, err := testMysql.OneOnOne().GetUpcomingOneOnOnes(ctx, user.ID)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(oneOnOnes), 1)

	// Find our created upcoming one-on-one
	found := false
	for _, o := range oneOnOnes {
		if o.UUID == upcomingOneOnOne.UUID {
			found = true
			require.Equal(t, "scheduled", o.Status)
			require.True(t, o.ScheduledDate.After(time.Now()))
			break
		}
	}
	require.True(t, found)
}

func TestGetOverdueOneOnOnes(t *testing.T) {
	ctx := context.Background()
	user := createRandomUserForTests(t)
	company := createRandomCompanyForTests(t)
	person := createRandomPersonForTests(t)

	// Create an overdue one-on-one
	pastDate := time.Now().Add(-24 * time.Hour) // 1 day ago
	overdueOneOnOne := entity.OneOnOne{
		UUID:          uuid.NewV4().String(),
		CompanyID:     company.ID,
		PersonID:      person.ID,
		ManagerID:     user.ID,
		ScheduledDate: pastDate,
		Duration:      30,
		Status:        "scheduled", // Still scheduled but past date = overdue
		Agenda:        "Overdue meeting",
	}

	_, err := testMysql.OneOnOne().CreateOneOnOne(ctx, overdueOneOnOne)
	require.NoError(t, err)

	// Get overdue one-on-ones
	oneOnOnes, err := testMysql.OneOnOne().GetOverdueOneOnOnes(ctx, user.ID)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(oneOnOnes), 1)

	// Find our created overdue one-on-one
	found := false
	for _, o := range oneOnOnes {
		if o.UUID == overdueOneOnOne.UUID {
			found = true
			require.Equal(t, "scheduled", o.Status)
			require.True(t, o.ScheduledDate.Before(time.Now()))
			break
		}
	}
	require.True(t, found)
}

// Error tests with mocks
func TestCreateOneOnOneErrorsWithMock(t *testing.T) {
	testForInsertErrorsWithMock(t, func(db *sql.DB) error {
		_, err := newOneOnOneRepo(db).CreateOneOnOne(context.Background(), entity.OneOnOne{})
		return err
	})
}

func TestGetOneOnOneByUUIDErrorsWithMock(t *testing.T) {
	testForSelectErrorsWithMock(t, "one_on_one_id", func(db *sql.DB) error {
		_, err := newOneOnOneRepo(db).GetOneOnOneByUUID(context.Background(), "one-on-one-uuid")
		return err
	})
}

func TestGetOneOnOnesByPersonErrorsWithMock(t *testing.T) {
	testForSelectErrorsWithMock(t, "one_on_one_id", func(db *sql.DB) error {
		_, _, err := newOneOnOneRepo(db).GetOneOnOnesByPerson(context.Background(), 1, 10, 0)
		return err
	})
}

func TestGetOneOnOnesByManagerErrorsWithMock(t *testing.T) {
	testForSelectErrorsWithMock(t, "one_on_one_id", func(db *sql.DB) error {
		_, _, err := newOneOnOneRepo(db).GetOneOnOnesByManager(context.Background(), 1, 10, 0)
		return err
	})
}

func TestUpdateOneOnOneErrorsWithMock(t *testing.T) {
	testForUpdateDeleteErrorsWithMock(t, func(db *sql.DB) error {
		return newOneOnOneRepo(db).UpdateOneOnOne(context.Background(), 1, entity.OneOnOne{})
	})
}

func TestDeleteOneOnOneErrorsWithMock(t *testing.T) {
	testForUpdateDeleteErrorsWithMock(t, func(db *sql.DB) error {
		return newOneOnOneRepo(db).DeleteOneOnOne(context.Background(), 1)
	})
}

func TestGetUpcomingOneOnOnesErrorsWithMock(t *testing.T) {
	testForSelectErrorsWithMock(t, "one_on_one_id", func(db *sql.DB) error {
		_, err := newOneOnOneRepo(db).GetUpcomingOneOnOnes(context.Background(), 1)
		return err
	})
}

func TestGetOverdueOneOnOnesErrorsWithMock(t *testing.T) {
	testForSelectErrorsWithMock(t, "one_on_one_id", func(db *sql.DB) error {
		_, err := newOneOnOneRepo(db).GetOverdueOneOnOnes(context.Background(), 1)
		return err
	})
}