package companyroute_test

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
	"github.com/diegoclair/leaderpro/internal/transport/rest/routes/companyroute"
	"github.com/diegoclair/leaderpro/internal/transport/rest/routes/test"
	"github.com/diegoclair/leaderpro/internal/transport/rest/viewmodel"
	echo "github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func TestHandler_handleCreateCompany(t *testing.T) {
	tests := append(test.PrivateEndpointValidations,
		test.PrivateEndpointTest{
			Name: "Should complete request with no error",
			Body: viewmodel.CompanyRequest{
				Name: "Test Company",
				Role: "Tech Lead",
			},
			SetupAuth: func(ctx context.Context, t *testing.T, req *http.Request, m test.SvcMocks) {
				test.AddAuthorization(ctx, t, req, m)
			},
			BuildMocks: func(ctx context.Context, m test.SvcMocks, body any) {
				if body != nil {
					companyReq := body.(viewmodel.CompanyRequest)
					m.CompanyAppMock.EXPECT().CreateCompany(ctx, companyReq.ToEntity()).Return(entity.Company{}, nil).Times(1)
				}
			},
			CheckResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
			},
		},
		test.PrivateEndpointTest{
			Name: "Should return error when create company fails",
			Body: viewmodel.CompanyRequest{
				Name: "Test Company",
				Role: "Manager",
			},
			SetupAuth: func(ctx context.Context, t *testing.T, req *http.Request, m test.SvcMocks) {
				test.AddAuthorization(ctx, t, req, m)
			},
			BuildMocks: func(ctx context.Context, m test.SvcMocks, body any) {
				if body != nil {
					companyReq := body.(viewmodel.CompanyRequest)
					m.CompanyAppMock.EXPECT().CreateCompany(ctx, companyReq.ToEntity()).Return(entity.Company{}, fmt.Errorf("error to create company")).Times(1)
				}
			},
			CheckResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusServiceUnavailable, recorder.Code)
				require.Contains(t, recorder.Body.String(), "error to create company")
			},
		},
	)

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			companyroute.Once = sync.Once{}
			m, server, ctrl := test.GetServerTest(t)
			defer ctrl.Finish()

			recorder := httptest.NewRecorder()
			url := "/companies"

			var body []byte
			var err error
			if tt.Body != nil {
				body, err = json.Marshal(tt.Body)
				require.NoError(t, err)
			}

			req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
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

func TestHandler_handleGetCompanies(t *testing.T) {
	tests := append(test.PrivateEndpointValidations,
		test.PrivateEndpointTest{
			Name: "Should complete request with no error",
			SetupAuth: func(ctx context.Context, t *testing.T, req *http.Request, m test.SvcMocks) {
				test.AddAuthorization(ctx, t, req, m)
			},
			BuildMocks: func(ctx context.Context, m test.SvcMocks, body any) {
				mockCompanies := []entity.Company{
					{UUID: "company-1", Name: "Company 1"},
					{UUID: "company-2", Name: "Company 2"},
				}
				m.CompanyAppMock.EXPECT().GetUserCompanies(ctx).Return(mockCompanies, nil).Times(1)
			},
			CheckResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				require.Contains(t, recorder.Body.String(), "Company 1")
				require.Contains(t, recorder.Body.String(), "Company 2")
			},
		},
		test.PrivateEndpointTest{
			Name: "Should return error when get companies fails",
			SetupAuth: func(ctx context.Context, t *testing.T, req *http.Request, m test.SvcMocks) {
				test.AddAuthorization(ctx, t, req, m)
			},
			BuildMocks: func(ctx context.Context, m test.SvcMocks, body any) {
				m.CompanyAppMock.EXPECT().GetUserCompanies(ctx).Return(nil, fmt.Errorf("error to get companies")).Times(1)
			},
			CheckResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusServiceUnavailable, recorder.Code)
				require.Contains(t, recorder.Body.String(), "error to get companies")
			},
		},
	)

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			companyroute.Once = sync.Once{}
			m, server, ctrl := test.GetServerTest(t)
			defer ctrl.Finish()

			recorder := httptest.NewRecorder()
			url := "/companies"

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

func TestHandler_handleGetCompanyByUUID(t *testing.T) {
	type args struct {
		companyUUID string
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
				companyUUID: "company-uuid-123",
			},
			buildMocks: func(ctx context.Context, m test.SvcMocks, args args) {
				mockCompany := entity.Company{
					UUID: args.companyUUID,
					Name: "Test Company",
				}
				m.CompanyAppMock.EXPECT().GetCompanyByUUID(ctx, args.companyUUID).Return(mockCompany, nil).Times(1)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				require.Contains(t, recorder.Body.String(), "Test Company")
			},
		},
		{
			name: "Should return error when get company fails",
			args: args{
				companyUUID: "company-uuid-123",
			},
			buildMocks: func(ctx context.Context, m test.SvcMocks, args args) {
				m.CompanyAppMock.EXPECT().GetCompanyByUUID(ctx, args.companyUUID).Return(entity.Company{}, fmt.Errorf("company not found")).Times(1)
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusServiceUnavailable, resp.Code)
				require.Contains(t, resp.Body.String(), "company not found")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			companyroute.Once = sync.Once{}
			m, server, ctrl := test.GetServerTest(t)
			defer ctrl.Finish()

			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/companies/%s", tt.args.companyUUID)

			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			ctx := test.GetTestContext(t, req, recorder, true)

			if tt.buildMocks != nil {
				tt.buildMocks(ctx, m, tt.args)
			}

			server.Echo().ServeHTTP(recorder, req)
			if tt.checkResponse != nil {
				tt.checkResponse(t, recorder)
			}
		})
	}
}

func TestHandler_handleUpdateCompany(t *testing.T) {
	type args struct {
		companyUUID string
		body        any
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
				companyUUID: "company-uuid-123",
				body: viewmodel.CompanyRequest{
					Name:        "Updated Company",
				},
			},
			buildMocks: func(ctx context.Context, m test.SvcMocks, args args) {
				body := args.body.(viewmodel.CompanyRequest)
				m.CompanyAppMock.EXPECT().UpdateCompany(ctx, args.companyUUID, body.ToEntity()).Return(nil).Times(1)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNoContent, recorder.Code)
			},
		},
		{
			name: "Should return error when update fails",
			args: args{
				companyUUID: "company-uuid-123",
				body: viewmodel.CompanyRequest{
					Name: "Updated Company",
				},
			},
			buildMocks: func(ctx context.Context, m test.SvcMocks, args args) {
				body := args.body.(viewmodel.CompanyRequest)
				m.CompanyAppMock.EXPECT().UpdateCompany(ctx, args.companyUUID, body.ToEntity()).Return(fmt.Errorf("error to update company")).Times(1)
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusServiceUnavailable, resp.Code)
				require.Contains(t, resp.Body.String(), "error to update company")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			companyroute.Once = sync.Once{}
			m, server, ctrl := test.GetServerTest(t)
			defer ctrl.Finish()

			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/companies/%s", tt.args.companyUUID)

			body, err := json.Marshal(tt.args.body)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(body))
			require.NoError(t, err)

			ctx := test.GetTestContext(t, req, recorder, true)

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

func TestHandler_handleDeleteCompany(t *testing.T) {
	type args struct {
		companyUUID string
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
				companyUUID: "company-uuid-123",
			},
			buildMocks: func(ctx context.Context, m test.SvcMocks, args args) {
				m.CompanyAppMock.EXPECT().DeleteCompany(ctx, args.companyUUID).Return(nil).Times(1)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNoContent, recorder.Code)
			},
		},
		{
			name: "Should return error when delete fails",
			args: args{
				companyUUID: "company-uuid-123",
			},
			buildMocks: func(ctx context.Context, m test.SvcMocks, args args) {
				m.CompanyAppMock.EXPECT().DeleteCompany(ctx, args.companyUUID).Return(fmt.Errorf("error to delete company")).Times(1)
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusServiceUnavailable, resp.Code)
				require.Contains(t, resp.Body.String(), "error to delete company")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			companyroute.Once = sync.Once{}
			m, server, ctrl := test.GetServerTest(t)
			defer ctrl.Finish()

			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/companies/%s", tt.args.companyUUID)

			req, err := http.NewRequest(http.MethodDelete, url, nil)
			require.NoError(t, err)

			ctx := test.GetTestContext(t, req, recorder, true)

			if tt.buildMocks != nil {
				tt.buildMocks(ctx, m, tt.args)
			}

			server.Echo().ServeHTTP(recorder, req)
			if tt.checkResponse != nil {
				tt.checkResponse(t, recorder)
			}
		})
	}
}
