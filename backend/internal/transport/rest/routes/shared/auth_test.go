package shared_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/diegoclair/leaderpro/infra/contract"
	infraMocks "github.com/diegoclair/leaderpro/infra/mocks"
	"github.com/diegoclair/leaderpro/internal/application/dto"
	"github.com/diegoclair/leaderpro/internal/domain/entity"
	"github.com/diegoclair/leaderpro/internal/transport/rest/routes/shared"
	"github.com/diegoclair/leaderpro/internal/transport/rest/routes/test"
	"github.com/diegoclair/leaderpro/mocks"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestAuthHelper_DoLogin(t *testing.T) {
	type args struct {
		loginInput dto.LoginInput
		userAgent  string
		clientIP   string
	}

	tests := []struct {
		name       string
		args       args
		buildMocks func(ctx context.Context, m test.SvcMocks, args args, echoContext echo.Context)
		wantErr    bool
		errMsg     string
	}{
		{
			name: "Should complete login with no error",
			args: args{
				loginInput: dto.LoginInput{
					Email:    "test@test.com",
					Password: "12345678",
				},
				userAgent: "TestAgent/1.0",
				clientIP:  "192.168.1.100",
			},
			buildMocks: func(ctx context.Context, m test.SvcMocks, args args, echoContext echo.Context) {
				// Mock Login
				user := entity.User{
					ID:   1,
					UUID: "user-uuid-123",
					Name: "Test User",
				}
				m.AuthAppMock.EXPECT().Login(ctx, args.loginInput).Return(user, nil).Times(1)

				// Mock token creation
				tokenPayload := contract.TokenPayload{
					ExpiredAt: time.Now().Add(15 * time.Minute),
				}
				refreshTokenPayload := contract.TokenPayload{
					ExpiredAt: time.Now().Add(24 * time.Hour),
				}

				m.AuthTokenMock.EXPECT().CreateAccessToken(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, req contract.TokenPayloadInput) (string, contract.TokenPayload, error) {
						require.Equal(t, "user-uuid-123", req.UserUUID)
						require.NotEmpty(t, req.SessionUUID)
						return "access-token-123", tokenPayload, nil
					}).Times(1)

				m.AuthTokenMock.EXPECT().CreateRefreshToken(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, req contract.TokenPayloadInput) (string, contract.TokenPayload, error) {
						require.Equal(t, "user-uuid-123", req.UserUUID)
						require.NotEmpty(t, req.SessionUUID)
						return "refresh-token-123", refreshTokenPayload, nil
					}).Times(1)

				// Mock session creation
				m.AuthAppMock.EXPECT().CreateSession(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, req dto.Session) error {
						require.NotEmpty(t, req.SessionUUID)
						require.Equal(t, int64(1), req.UserID)
						require.Equal(t, "refresh-token-123", req.RefreshToken)
						require.Equal(t, args.userAgent, req.UserAgent)
						require.Equal(t, args.clientIP, req.ClientIP)
						require.NotEmpty(t, req.RefreshTokenExpiredAt)
						return nil
					}).Times(1)
			},
			wantErr: false,
		},
		{
			name: "Should return error when login fails",
			args: args{
				loginInput: dto.LoginInput{
					Email:    "test@test.com",
					Password: "wrongpassword",
				},
			},
			buildMocks: func(ctx context.Context, m test.SvcMocks, args args, echoContext echo.Context) {
				m.AuthAppMock.EXPECT().Login(ctx, args.loginInput).
					Return(entity.User{}, fmt.Errorf("invalid credentials")).Times(1)
			},
			wantErr: true,
			errMsg:  "invalid credentials",
		},
		{
			name: "Should return error when create access token fails",
			args: args{
				loginInput: dto.LoginInput{
					Email:    "test@test.com",
					Password: "12345678",
				},
			},
			buildMocks: func(ctx context.Context, m test.SvcMocks, args args, echoContext echo.Context) {
				user := entity.User{ID: 1, UUID: "user-uuid-123"}
				m.AuthAppMock.EXPECT().Login(ctx, args.loginInput).Return(user, nil).Times(1)
				m.AuthTokenMock.EXPECT().CreateAccessToken(ctx, gomock.Any()).
					Return("", contract.TokenPayload{}, fmt.Errorf("error creating access token")).Times(1)
			},
			wantErr: true,
			errMsg:  "error creating access token",
		},
		{
			name: "Should return error when create refresh token fails",
			args: args{
				loginInput: dto.LoginInput{
					Email:    "test@test.com",
					Password: "12345678",
				},
			},
			buildMocks: func(ctx context.Context, m test.SvcMocks, args args, echoContext echo.Context) {
				user := entity.User{ID: 1, UUID: "user-uuid-123"}
				tokenPayload := contract.TokenPayload{ExpiredAt: time.Now().Add(15 * time.Minute)}

				m.AuthAppMock.EXPECT().Login(ctx, args.loginInput).Return(user, nil).Times(1)
				m.AuthTokenMock.EXPECT().CreateAccessToken(ctx, gomock.Any()).
					Return("access-token-123", tokenPayload, nil).Times(1)
				m.AuthTokenMock.EXPECT().CreateRefreshToken(ctx, gomock.Any()).
					Return("", contract.TokenPayload{}, fmt.Errorf("error creating refresh token")).Times(1)
			},
			wantErr: true,
			errMsg:  "error creating refresh token",
		},
		{
			name: "Should return error when create session fails",
			args: args{
				loginInput: dto.LoginInput{
					Email:    "test@test.com",
					Password: "12345678",
				},
			},
			buildMocks: func(ctx context.Context, m test.SvcMocks, args args, echoContext echo.Context) {
				user := entity.User{ID: 1, UUID: "user-uuid-123"}
				tokenPayload := contract.TokenPayload{ExpiredAt: time.Now().Add(15 * time.Minute)}
				refreshTokenPayload := contract.TokenPayload{ExpiredAt: time.Now().Add(24 * time.Hour)}

				m.AuthAppMock.EXPECT().Login(ctx, args.loginInput).Return(user, nil).Times(1)
				m.AuthTokenMock.EXPECT().CreateAccessToken(ctx, gomock.Any()).
					Return("access-token-123", tokenPayload, nil).Times(1)
				m.AuthTokenMock.EXPECT().CreateRefreshToken(ctx, gomock.Any()).
					Return("refresh-token-123", refreshTokenPayload, nil).Times(1)
				m.AuthAppMock.EXPECT().CreateSession(ctx, gomock.Any()).
					Return(fmt.Errorf("error creating session")).Times(1)
			},
			wantErr: true,
			errMsg:  "error creating session",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := test.SvcMocks{
				AuthAppMock:   mocks.NewMockAuthApp(ctrl),
				UserAppMock:   mocks.NewMockUserApp(ctrl),
				AuthTokenMock: infraMocks.NewMockAuthToken(ctrl),
			}

			// Create echo context with request
			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/", nil)
			req.Header.Set("User-Agent", tt.args.userAgent)
			req.Header.Set("X-Forwarded-For", tt.args.clientIP)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			ctx := context.Background()

			// Setup mocks
			if tt.buildMocks != nil {
				tt.buildMocks(ctx, m, tt.args, c)
			}

			// Create AuthHelper
			authHelper := shared.NewAuthHelper(m.AuthAppMock, m.UserAppMock, m.AuthTokenMock)

			// Execute
			result, err := authHelper.DoLogin(ctx, c, tt.args.loginInput)

			// Assert
			if tt.wantErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.errMsg)
				require.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				require.NotNil(t, result.User)
				require.NotNil(t, result.Auth)
				require.Equal(t, "access-token-123", result.Auth.AccessToken)
				require.Equal(t, "refresh-token-123", result.Auth.RefreshToken)
				require.NotEmpty(t, result.Auth.AccessTokenExpiresAt)
				require.NotEmpty(t, result.Auth.RefreshTokenExpiresAt)
			}
		})
	}
}

func TestNewAuthHelper(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	authApp := mocks.NewMockAuthApp(ctrl)
	userApp := mocks.NewMockUserApp(ctrl)
	authToken := infraMocks.NewMockAuthToken(ctrl)

	helper := shared.NewAuthHelper(authApp, userApp, authToken)

	require.NotNil(t, helper)
}