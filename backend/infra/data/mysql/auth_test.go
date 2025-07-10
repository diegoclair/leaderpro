package mysql

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/diegoclair/leaderpro/internal/application/dto"
	"github.com/diegoclair/leaderpro/internal/domain/entity"
	"github.com/stretchr/testify/require"
	"github.com/twinj/uuid"
)

func validateTwoSessions(t *testing.T, sessionExpected dto.Session, sessionToCompare dto.Session) {
	require.NotZero(t, sessionToCompare.UserID)
	require.Equal(t, sessionExpected.SessionUUID, sessionToCompare.SessionUUID)
	require.Equal(t, sessionExpected.RefreshToken, sessionToCompare.RefreshToken)
	require.Equal(t, sessionExpected.UserAgent, sessionToCompare.UserAgent)
	require.Equal(t, sessionExpected.ClientIP, sessionToCompare.ClientIP)
	require.Equal(t, sessionExpected.IsBlocked, sessionToCompare.IsBlocked)
	require.WithinDuration(t, sessionExpected.RefreshTokenExpiredAt, sessionToCompare.RefreshTokenExpiredAt, 2*time.Second)
}

func createRandomUserForAuth(t *testing.T) int64 {
	// Create a user using the user repository for auth tests
	userRepo := newUserRepo(testMysql.(*MysqlConn).db)
	
	user := entity.User{
		UUID:     uuid.NewV4().String(),
		Email:    "authtest" + uuid.NewV4().String()[:8] + "@example.com",
		Name:     "Auth Test User",
		Password: "hashedpassword",
		Phone:    "+1234567890",
		Active:   true,
	}

	userID, err := userRepo.CreateUser(context.Background(), user)
	require.NoError(t, err)
	require.NotZero(t, userID)

	return userID
}

func TestCreateAndGetSession(t *testing.T) {
	ctx := context.Background()
	userID := createRandomUserForAuth(t)

	session := dto.Session{
		SessionUUID:           uuid.NewV4().String(),
		UserID:               userID,
		RefreshToken:          uuid.NewV4().String(),
		UserAgent:             "user-agent",
		ClientIP:              "client-ip",
		IsBlocked:             false,
		RefreshTokenExpiredAt: time.Now().Add(24 * time.Hour),
	}

	sessionID, err := testMysql.Auth().CreateSession(ctx, session)
	require.NoError(t, err)
	require.NotZero(t, sessionID)

	session2, err := testMysql.Auth().GetSessionByUUID(ctx, session.SessionUUID)
	require.NoError(t, err)
	require.NotEmpty(t, session2)
	validateTwoSessions(t, session, session2)
}

func TestCreateSessionErrorsWithMock(t *testing.T) {
	testForInsertErrorsWithMock(t, func(db *sql.DB) error {
		_, err := newAuthRepo(db).CreateSession(context.Background(), dto.Session{})
		return err
	})
}

func TestGetSessionErrorsWithMock(t *testing.T) {
	testForSelectErrorsWithMock(t, "session_id", func(db *sql.DB) error {
		_, err := newAuthRepo(db).GetSessionByUUID(context.Background(), "session-uuid")
		return err
	})
}

func TestSetSessionAsBlocked(t *testing.T) {
	ctx := context.Background()
	userID := createRandomUserForAuth(t)

	session := dto.Session{
		SessionUUID:           uuid.NewV4().String(),
		UserID:               userID,
		RefreshToken:          uuid.NewV4().String(),
		UserAgent:             "user-agent",
		ClientIP:              "client-ip",
		IsBlocked:             false,
		RefreshTokenExpiredAt: time.Now().Add(24 * time.Hour),
	}

	sessionID, err := testMysql.Auth().CreateSession(ctx, session)
	require.NoError(t, err)
	require.NotZero(t, sessionID)

	err = testMysql.Auth().SetSessionAsBlocked(ctx, session.UserID)
	require.NoError(t, err)

	session2, err := testMysql.Auth().GetSessionByUUID(ctx, session.SessionUUID)
	require.NoError(t, err)
	require.NotEmpty(t, session2)
	require.True(t, session2.IsBlocked)
}

func TestSetSessionAsBlockedErrorsWithMock(t *testing.T) {
	testForUpdateDeleteErrorsWithMock(t, func(db *sql.DB) error {
		return newAuthRepo(db).SetSessionAsBlocked(context.Background(), 1)
	})
}
