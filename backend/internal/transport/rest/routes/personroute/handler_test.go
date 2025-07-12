package personroute_test

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
	"github.com/diegoclair/leaderpro/internal/transport/rest/routes/personroute"
	"github.com/diegoclair/leaderpro/internal/transport/rest/routes/test"
	"github.com/diegoclair/leaderpro/internal/transport/rest/viewmodel"
	echo "github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func TestHandler_handleCreatePerson(t *testing.T) {
	tests := append(test.PrivateEndpointValidations,
		test.PrivateEndpointTest{
			Name: "Should complete request with no error",
			Body: viewmodel.PersonRequest{
				Name:       "John Doe",
				Email:      "john.doe@example.com",
				Position:   "Software Engineer",
				Department: "Engineering",
				Phone:      "+1234567890",
			},
			SetupAuth: func(ctx context.Context, t *testing.T, req *http.Request, m test.AppMocks) {
				test.AddAuthorization(ctx, t, req, m)
			},
			BuildMocks: func(ctx context.Context, m test.AppMocks, body any) {
				if body != nil {
					personReq := body.(viewmodel.PersonRequest)
					m.PersonAppMock.EXPECT().CreatePerson(ctx, "company-uuid-123", personReq.ToEntity()).Return(nil).Times(1)
				}
			},
			CheckResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
			},
		},
		test.PrivateEndpointTest{
			Name: "Should return error when create person fails",
			Body: viewmodel.PersonRequest{
				Name:       "John Doe",
				Email:      "john.doe@example.com",
				Position:   "Software Engineer",
				Department: "Engineering",
			},
			SetupAuth: func(ctx context.Context, t *testing.T, req *http.Request, m test.AppMocks) {
				test.AddAuthorization(ctx, t, req, m)
			},
			BuildMocks: func(ctx context.Context, m test.AppMocks, body any) {
				if body != nil {
					personReq := body.(viewmodel.PersonRequest)
					m.PersonAppMock.EXPECT().CreatePerson(ctx, "company-uuid-123", personReq.ToEntity()).Return(fmt.Errorf("error to create person")).Times(1)
				}
			},
			CheckResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusServiceUnavailable, recorder.Code)
				require.Contains(t, recorder.Body.String(), "error to create person")
			},
		},
	)

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			personroute.Once = sync.Once{}
			m, server, ctrl := test.GetServerTest(t)
			defer ctrl.Finish()

			recorder := httptest.NewRecorder()
			url := "/companies/company-uuid-123/people"

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

func TestHandler_handleGetCompanyPeople(t *testing.T) {
	type args struct {
		companyUUID string
		search      string
	}

	tests := []struct {
		name          string
		args          args
		buildMocks    func(ctx context.Context, m test.AppMocks, args args)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Should complete request with no error - no search",
			args: args{
				companyUUID: "company-uuid-123",
			},
			buildMocks: func(ctx context.Context, m test.AppMocks, args args) {
				mockPeople := []entity.Person{
					{UUID: "person-1", Name: "John Doe"},
					{UUID: "person-2", Name: "Jane Smith"},
				}
				m.PersonAppMock.EXPECT().GetCompanyPeople(ctx, args.companyUUID).Return(mockPeople, nil).Times(1)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				require.Contains(t, recorder.Body.String(), "John Doe")
				require.Contains(t, recorder.Body.String(), "Jane Smith")
			},
		},
		{
			name: "Should complete request with search",
			args: args{
				companyUUID: "company-uuid-123",
				search:      "John",
			},
			buildMocks: func(ctx context.Context, m test.AppMocks, args args) {
				mockPeople := []entity.Person{
					{UUID: "person-1", Name: "John Doe"},
				}
				m.PersonAppMock.EXPECT().SearchPeople(ctx, args.companyUUID, args.search).Return(mockPeople, nil).Times(1)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				require.Contains(t, recorder.Body.String(), "John Doe")
			},
		},
		{
			name: "Should return error when get company people fails",
			args: args{
				companyUUID: "company-uuid-123",
			},
			buildMocks: func(ctx context.Context, m test.AppMocks, args args) {
				m.PersonAppMock.EXPECT().GetCompanyPeople(ctx, args.companyUUID).Return(nil, fmt.Errorf("error to get people")).Times(1)
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusServiceUnavailable, resp.Code)
				require.Contains(t, resp.Body.String(), "error to get people")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			personroute.Once = sync.Once{}
			m, server, ctrl := test.GetServerTest(t)
			defer ctrl.Finish()

			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/companies/%s/people", tt.args.companyUUID)
			if tt.args.search != "" {
				url += fmt.Sprintf("?search=%s", tt.args.search)
			}

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

func TestHandler_handleGetPersonByUUID(t *testing.T) {
	type args struct {
		companyUUID string
		personUUID  string
	}

	tests := []struct {
		name          string
		args          args
		buildMocks    func(ctx context.Context, m test.AppMocks, args args)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Should complete request with no error",
			args: args{
				companyUUID: "company-uuid-123",
				personUUID:  "person-uuid-456",
			},
			buildMocks: func(ctx context.Context, m test.AppMocks, args args) {
				mockPerson := entity.Person{
					UUID: args.personUUID,
					Name: "John Doe",
				}
				m.PersonAppMock.EXPECT().GetPersonByUUID(ctx, args.personUUID).Return(mockPerson, nil).Times(1)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				require.Contains(t, recorder.Body.String(), "John Doe")
			},
		},
		{
			name: "Should return error when get person fails",
			args: args{
				companyUUID: "company-uuid-123",
				personUUID:  "person-uuid-456",
			},
			buildMocks: func(ctx context.Context, m test.AppMocks, args args) {
				m.PersonAppMock.EXPECT().GetPersonByUUID(ctx, args.personUUID).Return(entity.Person{}, fmt.Errorf("person not found")).Times(1)
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusServiceUnavailable, resp.Code)
				require.Contains(t, resp.Body.String(), "person not found")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			personroute.Once = sync.Once{}
			m, server, ctrl := test.GetServerTest(t)
			defer ctrl.Finish()

			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/companies/%s/people/%s", tt.args.companyUUID, tt.args.personUUID)

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

func TestHandler_handleUpdatePerson(t *testing.T) {
	type args struct {
		companyUUID string
		personUUID  string
		body        any
	}

	tests := []struct {
		name          string
		args          args
		buildMocks    func(ctx context.Context, m test.AppMocks, args args)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Should complete request with no error",
			args: args{
				companyUUID: "company-uuid-123",
				personUUID:  "person-uuid-456",
				body: viewmodel.PersonRequest{
					Name:     "John Updated",
					Position: "Senior Engineer",
				},
			},
			buildMocks: func(ctx context.Context, m test.AppMocks, args args) {
				body := args.body.(viewmodel.PersonRequest)
				m.PersonAppMock.EXPECT().UpdatePerson(ctx, args.personUUID, body.ToEntity()).Return(nil).Times(1)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNoContent, recorder.Code)
			},
		},
		{
			name: "Should return error when update fails",
			args: args{
				companyUUID: "company-uuid-123",
				personUUID:  "person-uuid-456",
				body: viewmodel.PersonRequest{
					Name: "John Updated",
				},
			},
			buildMocks: func(ctx context.Context, m test.AppMocks, args args) {
				body := args.body.(viewmodel.PersonRequest)
				m.PersonAppMock.EXPECT().UpdatePerson(ctx, args.personUUID, body.ToEntity()).Return(fmt.Errorf("error to update person")).Times(1)
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusServiceUnavailable, resp.Code)
				require.Contains(t, resp.Body.String(), "error to update person")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			personroute.Once = sync.Once{}
			m, server, ctrl := test.GetServerTest(t)
			defer ctrl.Finish()

			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/companies/%s/people/%s", tt.args.companyUUID, tt.args.personUUID)

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

func TestHandler_handleDeletePerson(t *testing.T) {
	type args struct {
		companyUUID string
		personUUID  string
	}

	tests := []struct {
		name          string
		args          args
		buildMocks    func(ctx context.Context, m test.AppMocks, args args)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Should complete request with no error",
			args: args{
				companyUUID: "company-uuid-123",
				personUUID:  "person-uuid-456",
			},
			buildMocks: func(ctx context.Context, m test.AppMocks, args args) {
				m.PersonAppMock.EXPECT().DeletePerson(ctx, args.personUUID).Return(nil).Times(1)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNoContent, recorder.Code)
			},
		},
		{
			name: "Should return error when delete fails",
			args: args{
				companyUUID: "company-uuid-123",
				personUUID:  "person-uuid-456",
			},
			buildMocks: func(ctx context.Context, m test.AppMocks, args args) {
				m.PersonAppMock.EXPECT().DeletePerson(ctx, args.personUUID).Return(fmt.Errorf("error to delete person")).Times(1)
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusServiceUnavailable, resp.Code)
				require.Contains(t, resp.Body.String(), "error to delete person")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			personroute.Once = sync.Once{}
			m, server, ctrl := test.GetServerTest(t)
			defer ctrl.Finish()

			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/companies/%s/people/%s", tt.args.companyUUID, tt.args.personUUID)

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
