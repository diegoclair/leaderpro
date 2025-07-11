package userroute_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/diegoclair/leaderpro/infra/contract"
	"github.com/diegoclair/leaderpro/internal/application/dto"
	"github.com/diegoclair/leaderpro/internal/domain/entity"
	"github.com/diegoclair/leaderpro/internal/transport/rest/routes/test"
	"github.com/diegoclair/leaderpro/internal/transport/rest/routes/userroute"
	"github.com/diegoclair/leaderpro/internal/transport/rest/viewmodel"
	echo "github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestHandler_handleCreateUser(t *testing.T) {
	type args struct {
		body any
	}

	tests := []struct {
		name          string
		args          args
		buildMocks    func(ctx context.Context, m test.SvcMocks, args args)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Should complete request with no error and auto-login",
			args: args{
				body: viewmodel.CreateUser{
					Name:     "John Doe",
					Email:    "john.doe@example.com",
					Password: "password123",
					Phone:    "+1234567890",
				},
			},
			buildMocks: func(ctx context.Context, m test.SvcMocks, args args) {
				body := args.body.(viewmodel.CreateUser)
				
				// Mock CreateUser
				mockUser := entity.User{
					ID:    1,
					UUID:  "user-uuid-123",
					Name:  body.Name,
					Email: body.Email,
					Phone: body.Phone,
				}
				m.UserAppMock.EXPECT().CreateUser(ctx, body.ToEntity()).Return(mockUser, nil).Times(1)

				// Mock auto-login after user creation
				loginInput := dto.LoginInput{
					Email:    body.Email,
					Password: body.Password,
				}
				m.AuthAppMock.EXPECT().Login(ctx, loginInput).Return(mockUser, nil).Times(1)

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
						require.NotEmpty(t, req.RefreshTokenExpiredAt)
						return nil
					}).Times(1)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				
				// Verify AuthResponse structure
				var response viewmodel.AuthResponse
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				require.NoError(t, err)
				
				// Check user data
				require.Equal(t, "user-uuid-123", response.User.UUID)
				require.Equal(t, "John Doe", response.User.Name)
				require.Equal(t, "john.doe@example.com", response.User.Email)
				
				// Check auth tokens
				require.Equal(t, "access-token-123", response.Auth.AccessToken)
				require.Equal(t, "refresh-token-123", response.Auth.RefreshToken)
				require.NotEmpty(t, response.Auth.AccessTokenExpiresAt)
				require.NotEmpty(t, response.Auth.RefreshTokenExpiresAt)
			},
		},
		{
			name: "Should return error when body is invalid",
			args: args{
				body: "invalid body",
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, resp.Code)
				require.Contains(t, resp.Body.String(), "Unmarshal type error")
			},
		},
		{
			name: "Should return error when create user fails",
			args: args{
				body: viewmodel.CreateUser{
					Name:     "John Doe",
					Email:    "john.doe@example.com",
					Password: "password123",
				},
			},
			buildMocks: func(ctx context.Context, m test.SvcMocks, args args) {
				body := args.body.(viewmodel.CreateUser)
				m.UserAppMock.EXPECT().CreateUser(ctx, body.ToEntity()).Return(entity.User{}, fmt.Errorf("error to create user")).Times(1)
				// No login mocks needed since CreateUser fails before reaching login
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusServiceUnavailable, resp.Code)
				require.Contains(t, resp.Body.String(), "error to create user")
			},
		},
		{
			name: "Should return error when auto-login fails after user creation",
			args: args{
				body: viewmodel.CreateUser{
					Name:     "John Doe",
					Email:    "john.doe@example.com",
					Password: "password123",
					Phone:    "+1234567890",
				},
			},
			buildMocks: func(ctx context.Context, m test.SvcMocks, args args) {
				body := args.body.(viewmodel.CreateUser)
				
				// Mock CreateUser - succeeds
				mockUser := entity.User{
					ID:    1,
					UUID:  "user-uuid-123",
					Name:  body.Name,
					Email: body.Email,
					Phone: body.Phone,
				}
				m.UserAppMock.EXPECT().CreateUser(ctx, body.ToEntity()).Return(mockUser, nil).Times(1)

				// Mock auto-login - fails
				loginInput := dto.LoginInput{
					Email:    body.Email,
					Password: body.Password,
				}
				m.AuthAppMock.EXPECT().Login(ctx, loginInput).Return(entity.User{}, fmt.Errorf("login failed")).Times(1)
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusServiceUnavailable, resp.Code)
				require.Contains(t, resp.Body.String(), "login failed")
			},
		},
		{
			name: "Should return error when token creation fails in auto-login",
			args: args{
				body: viewmodel.CreateUser{
					Name:     "John Doe",
					Email:    "john.doe@example.com",
					Password: "password123",
					Phone:    "+1234567890",
				},
			},
			buildMocks: func(ctx context.Context, m test.SvcMocks, args args) {
				body := args.body.(viewmodel.CreateUser)
				
				// Mock CreateUser - succeeds
				mockUser := entity.User{
					ID:    1,
					UUID:  "user-uuid-123",
					Name:  body.Name,
					Email: body.Email,
					Phone: body.Phone,
				}
				m.UserAppMock.EXPECT().CreateUser(ctx, body.ToEntity()).Return(mockUser, nil).Times(1)

				// Mock auto-login steps
				loginInput := dto.LoginInput{
					Email:    body.Email,
					Password: body.Password,
				}
				m.AuthAppMock.EXPECT().Login(ctx, loginInput).Return(mockUser, nil).Times(1)

				// Mock token creation - fails
				m.AuthTokenMock.EXPECT().CreateAccessToken(ctx, gomock.Any()).
					Return("", contract.TokenPayload{}, fmt.Errorf("token creation failed")).Times(1)
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusServiceUnavailable, resp.Code)
				require.Contains(t, resp.Body.String(), "token creation failed")
			},
		},
		{
			name: "Should return error when session creation fails in auto-login",
			args: args{
				body: viewmodel.CreateUser{
					Name:     "John Doe",
					Email:    "john.doe@example.com",
					Password: "password123",
					Phone:    "+1234567890",
				},
			},
			buildMocks: func(ctx context.Context, m test.SvcMocks, args args) {
				body := args.body.(viewmodel.CreateUser)
				
				// Mock CreateUser - succeeds
				mockUser := entity.User{
					ID:    1,
					UUID:  "user-uuid-123",
					Name:  body.Name,
					Email: body.Email,
					Phone: body.Phone,
				}
				m.UserAppMock.EXPECT().CreateUser(ctx, body.ToEntity()).Return(mockUser, nil).Times(1)

				// Mock auto-login steps - all succeed until session
				loginInput := dto.LoginInput{
					Email:    body.Email,
					Password: body.Password,
				}
				m.AuthAppMock.EXPECT().Login(ctx, loginInput).Return(mockUser, nil).Times(1)

				// Mock token creation - succeeds
				tokenPayload := contract.TokenPayload{
					ExpiredAt: time.Now().Add(15 * time.Minute),
				}
				refreshTokenPayload := contract.TokenPayload{
					ExpiredAt: time.Now().Add(24 * time.Hour),
				}

				m.AuthTokenMock.EXPECT().CreateAccessToken(ctx, gomock.Any()).
					Return("access-token-123", tokenPayload, nil).Times(1)

				m.AuthTokenMock.EXPECT().CreateRefreshToken(ctx, gomock.Any()).
					Return("refresh-token-123", refreshTokenPayload, nil).Times(1)

				// Mock session creation - fails
				m.AuthAppMock.EXPECT().CreateSession(ctx, gomock.Any()).
					Return(fmt.Errorf("session creation failed")).Times(1)
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusServiceUnavailable, resp.Code)
				require.Contains(t, resp.Body.String(), "session creation failed")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userroute.Once = sync.Once{}
			m, server, ctrl := test.GetServerTest(t)
			defer ctrl.Finish()

			recorder := httptest.NewRecorder()
			url := "/users"

			body, err := json.Marshal(tt.args.body)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
			require.NoError(t, err)

			ctx := test.GetTestContext(t, req, recorder, false)

			if tt.buildMocks != nil {
				tt.buildMocks(ctx, m, tt.args)
			}

			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			server.Echo().ServeHTTP(recorder, req)
			if tt.checkResponse != nil {
				tt.checkResponse(t, recorder)
			}
		})
	}
}

func TestHandler_handleGetProfile(t *testing.T) {
	tests := append(test.PrivateEndpointValidations,
		test.PrivateEndpointTest{
			Name: "Should complete request with no error",
			SetupAuth: func(ctx context.Context, t *testing.T, req *http.Request, m test.SvcMocks) {
				test.AddAuthorization(ctx, t, req, m)
			},
			BuildMocks: func(ctx context.Context, m test.SvcMocks, body any) {
				mockUser := entity.User{
					UUID:  "user-uuid-123",
					Name:  "John Doe",
					Email: "john.doe@example.com",
				}
				m.UserAppMock.EXPECT().GetProfile(ctx).Return(mockUser, nil).Times(1)
			},
			CheckResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				require.Contains(t, recorder.Body.String(), "John Doe")
				require.Contains(t, recorder.Body.String(), "john.doe@example.com")
			},
		},
		test.PrivateEndpointTest{
			Name: "Should return error when get profile fails",
			SetupAuth: func(ctx context.Context, t *testing.T, req *http.Request, m test.SvcMocks) {
				test.AddAuthorization(ctx, t, req, m)
			},
			BuildMocks: func(ctx context.Context, m test.SvcMocks, body any) {
				m.UserAppMock.EXPECT().GetProfile(ctx).Return(entity.User{}, fmt.Errorf("error to get profile")).Times(1)
			},
			CheckResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusServiceUnavailable, recorder.Code)
				require.Contains(t, recorder.Body.String(), "error to get profile")
			},
		},
	)

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			userroute.Once = sync.Once{}
			m, server, ctrl := test.GetServerTest(t)
			defer ctrl.Finish()

			recorder := httptest.NewRecorder()
			url := "/users/profile"

			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			ctx := test.GetTestContext(t, req, recorder, true)

			if tt.SetupAuth != nil {
				tt.SetupAuth(ctx, t, req, m)
			}

			if tt.BuildMocks != nil {
				tt.BuildMocks(ctx, m, tt.Body)
			}

			server.Echo().ServeHTTP(recorder, req)
			if tt.CheckResponse != nil {
				tt.CheckResponse(t, recorder)
			}
		})
	}
}

func TestHandler_handleUpdateProfile(t *testing.T) {
	tests := append(test.PrivateEndpointValidations,
		test.PrivateEndpointTest{
			Name: "Should complete request with no error",
			Body: viewmodel.UpdateUser{
				Name:  "John Updated",
				Phone: "+9876543210",
			},
			SetupAuth: func(ctx context.Context, t *testing.T, req *http.Request, m test.SvcMocks) {
				test.AddAuthorization(ctx, t, req, m)
			},
			BuildMocks: func(ctx context.Context, m test.SvcMocks, body any) {
				if body != nil {
					updateUser := body.(viewmodel.UpdateUser)
					mockUser := entity.User{
						UUID:  "user-uuid-123",
						Name:  updateUser.Name,
						Phone: updateUser.Phone,
					}
					m.UserAppMock.EXPECT().UpdateProfile(ctx, updateUser.ToEntity()).Return(mockUser, nil).Times(1)
				}
			},
			CheckResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				require.Contains(t, recorder.Body.String(), "John Updated")
			},
		},
		test.PrivateEndpointTest{
			Name: "Should return error when body is invalid",
			Body: "invalid body", // This will cause Bind() to fail
			SetupAuth: func(ctx context.Context, t *testing.T, req *http.Request, m test.SvcMocks) {
				test.AddAuthorization(ctx, t, req, m)
			},
			BuildMocks: func(ctx context.Context, m test.SvcMocks, body any) {
				// No mocks needed since Bind() fails before reaching service
			},
			CheckResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				require.Contains(t, recorder.Body.String(), "Unmarshal type error")
			},
		},
		test.PrivateEndpointTest{
			Name: "Should return error when update profile fails",
			Body: viewmodel.UpdateUser{
				Name: "John Updated",
			},
			SetupAuth: func(ctx context.Context, t *testing.T, req *http.Request, m test.SvcMocks) {
				test.AddAuthorization(ctx, t, req, m)
			},
			BuildMocks: func(ctx context.Context, m test.SvcMocks, body any) {
				if body != nil {
					updateUser := body.(viewmodel.UpdateUser)
					m.UserAppMock.EXPECT().UpdateProfile(ctx, updateUser.ToEntity()).Return(entity.User{}, fmt.Errorf("error to update profile")).Times(1)
				}
			},
			CheckResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusServiceUnavailable, recorder.Code)
				require.Contains(t, recorder.Body.String(), "error to update profile")
			},
		},
	)

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			userroute.Once = sync.Once{}
			m, server, ctrl := test.GetServerTest(t)
			defer ctrl.Finish()

			recorder := httptest.NewRecorder()
			url := "/users/profile"

			var body []byte
			var err error
			if tt.Body != nil {
				body, err = json.Marshal(tt.Body)
				require.NoError(t, err)
			}

			req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(body))
			require.NoError(t, err)

			ctx := test.GetTestContext(t, req, recorder, true)

			if tt.SetupAuth != nil {
				tt.SetupAuth(ctx, t, req, m)
			}

			if tt.BuildMocks != nil {
				tt.BuildMocks(ctx, m, tt.Body)
			}

			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			server.Echo().ServeHTTP(recorder, req)
			if tt.CheckResponse != nil {
				tt.CheckResponse(t, recorder)
			}
		})
	}
}