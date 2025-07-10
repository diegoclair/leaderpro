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

	"github.com/diegoclair/leaderpro/internal/domain/entity"
	"github.com/diegoclair/leaderpro/internal/transport/rest/routes/test"
	"github.com/diegoclair/leaderpro/internal/transport/rest/routes/userroute"
	"github.com/diegoclair/leaderpro/internal/transport/rest/viewmodel"
	echo "github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
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
			name: "Should complete request with no error",
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
				mockUser := entity.User{
					UUID:  "user-uuid-123",
					Name:  body.Name,
					Email: body.Email,
					Phone: body.Phone,
				}
				m.UserAppMock.EXPECT().CreateUser(ctx, body.ToEntity()).Return(mockUser, nil).Times(1)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				require.Contains(t, recorder.Body.String(), "John Doe")
				require.Contains(t, recorder.Body.String(), "john.doe@example.com")
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
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusServiceUnavailable, resp.Code)
				require.Contains(t, resp.Body.String(), "error to create user")
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