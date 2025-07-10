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

func validateTwoFeedbacks(t *testing.T, feedbackExpected entity.Feedback, feedbackToCompare entity.Feedback) {
	require.NotZero(t, feedbackToCompare.ID)
	require.Equal(t, feedbackExpected.UUID, feedbackToCompare.UUID)
	require.Equal(t, feedbackExpected.CompanyID, feedbackToCompare.CompanyID)
	require.Equal(t, feedbackExpected.PersonID, feedbackToCompare.PersonID)
	require.Equal(t, feedbackExpected.GivenBy, feedbackToCompare.GivenBy)
	require.Equal(t, feedbackExpected.OneOnOneID, feedbackToCompare.OneOnOneID)
	require.Equal(t, feedbackExpected.Type, feedbackToCompare.Type)
	require.Equal(t, feedbackExpected.Category, feedbackToCompare.Category)
	require.Equal(t, feedbackExpected.Content, feedbackToCompare.Content)
	require.Equal(t, feedbackExpected.MentionedFrom, feedbackToCompare.MentionedFrom)
	require.WithinDuration(t, feedbackExpected.MentionedDate, feedbackToCompare.MentionedDate, time.Second)
	require.Equal(t, feedbackExpected.IsPrivate, feedbackToCompare.IsPrivate)
}




func createRandomFeedback(t *testing.T) entity.Feedback {
	// Create user, company, and person first
	user := createRandomUserForTests(t)
	company := createRandomCompanyForTests(t)
	person := createRandomPersonForTests(t)

	mentionedDate := time.Now().Add(-24 * time.Hour) // Yesterday

	feedback := entity.Feedback{
		UUID:          uuid.NewV4().String(),
		CompanyID:     company.ID,
		PersonID:      person.ID,
		GivenBy:       user.ID,
		OneOnOneID:    nil, // Direct feedback, not from a 1:1
		Type:          "positive",
		Category:      "performance",
		Content:       "Excellent work on the latest project. Shows great attention to detail.",
		MentionedFrom: "Daily standup",
		MentionedDate: mentionedDate,
		IsPrivate:     false,
	}

	feedbackID, err := testMysql.Feedback().CreateFeedback(context.Background(), feedback)
	require.NoError(t, err)
	require.NotZero(t, feedbackID)

	feedback.ID = feedbackID
	return feedback
}

func TestCreateFeedback(t *testing.T) {
	ctx := context.Background()
	user := createRandomUserForTests(t)
	company := createRandomCompanyForTests(t)
	person := createRandomPersonForTests(t)
	// Create one-on-one for feedback test
	oneOnOne := entity.OneOnOne{
		UUID:          uuid.NewV4().String(),
		CompanyID:     company.ID,
		PersonID:      person.ID,
		ManagerID:     user.ID,
		ScheduledDate: time.Now().Add(24 * time.Hour),
		Duration:      30,
		Status:        "scheduled",
		Agenda:        "Test meeting",
	}
	
	oneOnOneID, err := testMysql.OneOnOne().CreateOneOnOne(context.Background(), oneOnOne)
	require.NoError(t, err)
	oneOnOne.ID = oneOnOneID

	feedback := entity.Feedback{
		UUID:          uuid.NewV4().String(),
		CompanyID:     company.ID,
		PersonID:      person.ID,
		GivenBy:       user.ID,
		OneOnOneID:    &oneOnOne.ID, // From a 1:1 meeting
		Type:          "constructive",
		Category:      "communication",
		Content:       "Could improve clarity in written communication.",
		MentionedFrom: "1:1 Meeting",
		MentionedDate: time.Now(),
		IsPrivate:     true,
	}

	feedbackID, err := testMysql.Feedback().CreateFeedback(ctx, feedback)
	require.NoError(t, err)
	require.NotZero(t, feedbackID)
}

func TestGetFeedbackByUUID(t *testing.T) {
	ctx := context.Background()
	feedback := createRandomFeedback(t)

	feedback2, err := testMysql.Feedback().GetFeedbackByUUID(ctx, feedback.UUID)
	require.NoError(t, err)
	require.NotEmpty(t, feedback2)
	validateTwoFeedbacks(t, feedback, feedback2)
}

func TestGetFeedbackByPerson(t *testing.T) {
	ctx := context.Background()
	feedback := createRandomFeedback(t)

	feedbacks, totalRecords, err := testMysql.Feedback().GetFeedbackByPerson(ctx, feedback.PersonID, 10, 0)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(feedbacks), 1)
	require.GreaterOrEqual(t, totalRecords, int64(1))

	// Find our created feedback
	found := false
	for _, f := range feedbacks {
		if f.UUID == feedback.UUID {
			found = true
			validateTwoFeedbacks(t, feedback, f)
			break
		}
	}
	require.True(t, found)
}

func TestGetFeedbackByGiver(t *testing.T) {
	ctx := context.Background()
	feedback := createRandomFeedback(t)

	feedbacks, totalRecords, err := testMysql.Feedback().GetFeedbackByGiver(ctx, feedback.GivenBy, 10, 0)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(feedbacks), 1)
	require.GreaterOrEqual(t, totalRecords, int64(1))

	// Find our created feedback
	found := false
	for _, f := range feedbacks {
		if f.UUID == feedback.UUID {
			found = true
			validateTwoFeedbacks(t, feedback, f)
			break
		}
	}
	require.True(t, found)
}

func TestUpdateFeedback(t *testing.T) {
	ctx := context.Background()
	feedback := createRandomFeedback(t)

	// Update feedback fields
	updatedFeedback := entity.Feedback{
		Type:          "observation",
		Category:      "collaboration",
		Content:       "Updated feedback content with more specific examples.",
		MentionedFrom: "Updated context",
		MentionedDate: time.Now(),
		IsPrivate:     true, // Changed to private
	}

	err := testMysql.Feedback().UpdateFeedback(ctx, feedback.ID, updatedFeedback)
	require.NoError(t, err)

	// Get updated feedback and validate
	retrievedFeedback, err := testMysql.Feedback().GetFeedbackByUUID(ctx, feedback.UUID)
	require.NoError(t, err)
	require.Equal(t, updatedFeedback.Type, retrievedFeedback.Type)
	require.Equal(t, updatedFeedback.Category, retrievedFeedback.Category)
	require.Equal(t, updatedFeedback.Content, retrievedFeedback.Content)
	require.Equal(t, updatedFeedback.MentionedFrom, retrievedFeedback.MentionedFrom)
	require.Equal(t, updatedFeedback.IsPrivate, retrievedFeedback.IsPrivate)
}

func TestDeleteFeedback(t *testing.T) {
	ctx := context.Background()
	feedback := createRandomFeedback(t)

	err := testMysql.Feedback().DeleteFeedback(ctx, feedback.ID)
	require.NoError(t, err)

	// Try to get the deleted feedback - should return error (not found)
	_, err = testMysql.Feedback().GetFeedbackByUUID(ctx, feedback.UUID)
	require.Error(t, err)
}

func TestGetFeedbackSummary(t *testing.T) {
	ctx := context.Background()
	user := createRandomUserForTests(t)
	company := createRandomCompanyForTests(t)
	person := createRandomPersonForTests(t)

	// Create multiple feedbacks for the person
	positiveFeedback := entity.Feedback{
		UUID:          uuid.NewV4().String(),
		CompanyID:     company.ID,
		PersonID:      person.ID,
		GivenBy:       user.ID,
		Type:          "positive",
		Category:      "performance",
		Content:       "Great job on project delivery!",
		MentionedFrom: "Performance review",
		MentionedDate: time.Now(),
		IsPrivate:     false,
	}

	constructiveFeedback := entity.Feedback{
		UUID:          uuid.NewV4().String(),
		CompanyID:     company.ID,
		PersonID:      person.ID,
		GivenBy:       user.ID,
		Type:          "constructive",
		Category:      "communication",
		Content:       "Could improve meeting facilitation skills.",
		MentionedFrom: "1:1 meeting",
		MentionedDate: time.Now(),
		IsPrivate:     true,
	}

	// Create feedbacks
	_, err := testMysql.Feedback().CreateFeedback(ctx, positiveFeedback)
	require.NoError(t, err)

	_, err = testMysql.Feedback().CreateFeedback(ctx, constructiveFeedback)
	require.NoError(t, err)

	// Get feedback summary
	summary, err := testMysql.Feedback().GetFeedbackSummary(ctx, person.ID, "2024")
	require.NoError(t, err)
	
	require.Equal(t, person.ID, summary.PersonID)
	require.Equal(t, "2024", summary.Period)
	require.GreaterOrEqual(t, summary.TotalCount, 2)
	require.GreaterOrEqual(t, summary.PositiveCount, 1)
	require.GreaterOrEqual(t, summary.Constructive, 1)
}

func TestMultipleFeedbackTypes(t *testing.T) {
	ctx := context.Background()
	user := createRandomUserForTests(t)
	company := createRandomCompanyForTests(t)
	person := createRandomPersonForTests(t)

	// Create different types of feedback
	feedbackTypes := []struct {
		feedbackType string
		category     string
		content      string
		isPrivate    bool
	}{
		{"positive", "performance", "Excellent problem-solving skills", false},
		{"constructive", "behavior", "Could be more punctual to meetings", true},
		{"observation", "skill", "Shows good leadership potential", false},
		{"positive", "collaboration", "Great team player", false},
		{"constructive", "performance", "Needs to improve code review quality", true},
	}

	createdFeedbacks := make([]entity.Feedback, 0, len(feedbackTypes))

	for _, ft := range feedbackTypes {
		feedback := entity.Feedback{
			UUID:          uuid.NewV4().String(),
			CompanyID:     company.ID,
			PersonID:      person.ID,
			GivenBy:       user.ID,
			Type:          ft.feedbackType,
			Category:      ft.category,
			Content:       ft.content,
			MentionedFrom: "Test context",
			MentionedDate: time.Now(),
			IsPrivate:     ft.isPrivate,
		}

		feedbackID, err := testMysql.Feedback().CreateFeedback(ctx, feedback)
		require.NoError(t, err)
		feedback.ID = feedbackID
		createdFeedbacks = append(createdFeedbacks, feedback)
	}

	// Get all feedback for the person and verify types
	feedbacks, totalRecords, err := testMysql.Feedback().GetFeedbackByPerson(ctx, person.ID, 100, 0)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(feedbacks), len(feedbackTypes))
	require.GreaterOrEqual(t, totalRecords, int64(len(feedbackTypes)))

	// Count feedback types
	typeCounts := make(map[string]int)
	for _, f := range feedbacks {
		// Only count our test feedbacks
		for _, created := range createdFeedbacks {
			if f.UUID == created.UUID {
				typeCounts[f.Type]++
				break
			}
		}
	}

	require.GreaterOrEqual(t, typeCounts["positive"], 2)
	require.GreaterOrEqual(t, typeCounts["constructive"], 2)
	require.GreaterOrEqual(t, typeCounts["observation"], 1)
}

// Error tests with mocks
func TestCreateFeedbackErrorsWithMock(t *testing.T) {
	testForInsertErrorsWithMock(t, func(db *sql.DB) error {
		_, err := newFeedbackRepo(db).CreateFeedback(context.Background(), entity.Feedback{})
		return err
	})
}

func TestGetFeedbackByUUIDErrorsWithMock(t *testing.T) {
	testForSelectErrorsWithMock(t, "feedback_id", func(db *sql.DB) error {
		_, err := newFeedbackRepo(db).GetFeedbackByUUID(context.Background(), "feedback-uuid")
		return err
	})
}

func TestGetFeedbackByPersonErrorsWithMock(t *testing.T) {
	testForSelectErrorsWithMock(t, "feedback_id", func(db *sql.DB) error {
		_, _, err := newFeedbackRepo(db).GetFeedbackByPerson(context.Background(), 1, 10, 0)
		return err
	})
}

func TestGetFeedbackByGiverErrorsWithMock(t *testing.T) {
	testForSelectErrorsWithMock(t, "feedback_id", func(db *sql.DB) error {
		_, _, err := newFeedbackRepo(db).GetFeedbackByGiver(context.Background(), 1, 10, 0)
		return err
	})
}

func TestUpdateFeedbackErrorsWithMock(t *testing.T) {
	testForUpdateDeleteErrorsWithMock(t, func(db *sql.DB) error {
		return newFeedbackRepo(db).UpdateFeedback(context.Background(), 1, entity.Feedback{})
	})
}

func TestDeleteFeedbackErrorsWithMock(t *testing.T) {
	testForUpdateDeleteErrorsWithMock(t, func(db *sql.DB) error {
		return newFeedbackRepo(db).DeleteFeedback(context.Background(), 1)
	})
}

func TestGetFeedbackSummaryErrorsWithMock(t *testing.T) {
	testForSelectErrorsWithMock(t, "feedback_id", func(db *sql.DB) error {
		_, err := newFeedbackRepo(db).GetFeedbackSummary(context.Background(), 1, "2024")
		return err
	})
}