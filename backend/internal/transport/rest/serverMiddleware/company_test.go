package servermiddleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/diegoclair/go_utils/resterrors"
	"github.com/diegoclair/leaderpro/infra"
	"github.com/diegoclair/leaderpro/mocks"
	echo "github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCompanyOwnershipMiddleware(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockCompanyService := mocks.NewMockCompanyApp(ctrl)
	middleware := CompanyOwnershipMiddleware(mockCompanyService)

	t.Run("Should complete the middleware without errors when company ownership is valid", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/companies/company-uuid-123/people", nil)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)
		c.SetPath("/companies/:company_uuid/people")
		c.SetParamNames("company_uuid")
		c.SetParamValues("company-uuid-123")
		c.Set(infra.UserUUIDKey.String(), "user-uuid-456")

		mockCompanyService.EXPECT().ValidateCompanyOwnership(
			gomock.Any(),
			"company-uuid-123",
			"user-uuid-456",
		).Return(nil)

		err := middleware(func(c echo.Context) error {
			return nil
		})(c)

		assert.Nil(t, err)
		assert.Equal(t, "company-uuid-123", c.Get(infra.CompanyUUIDKey.String()))
	})

	t.Run("Should return error when company_uuid is missing", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/companies//people", nil)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)
		c.SetPath("/companies/:company_uuid/people")
		c.SetParamNames("company_uuid")
		c.SetParamValues("")
		c.Set(infra.UserUUIDKey.String(), "user-uuid-456")

		err := middleware(func(c echo.Context) error {
			return nil
		})(c)

		assert.NotNil(t, err)
		assert.Equal(t, http.StatusBadRequest, err.(resterrors.RestErr).StatusCode())
		assert.Equal(t, "company_uuid is required", err.(resterrors.RestErr).Message())
	})

	t.Run("Should return error when user is not authenticated", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/companies/company-uuid-123/people", nil)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)
		c.SetPath("/companies/:company_uuid/people")
		c.SetParamNames("company_uuid")
		c.SetParamValues("company-uuid-123")
		// User UUID not set in context

		err := middleware(func(c echo.Context) error {
			return nil
		})(c)

		assert.NotNil(t, err)
		assert.Equal(t, http.StatusUnauthorized, err.(resterrors.RestErr).StatusCode())
		assert.Equal(t, "user not authenticated", err.(resterrors.RestErr).Message())
	})

	t.Run("Should return error when user UUID is invalid type", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/companies/company-uuid-123/people", nil)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)
		c.SetPath("/companies/:company_uuid/people")
		c.SetParamNames("company_uuid")
		c.SetParamValues("company-uuid-123")
		c.Set(infra.UserUUIDKey.String(), 123) // Invalid type (int instead of string)

		err := middleware(func(c echo.Context) error {
			return nil
		})(c)

		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.(resterrors.RestErr).StatusCode())
		assert.Equal(t, "invalid user context", err.(resterrors.RestErr).Message())
	})

	t.Run("Should return error when company validation fails with not found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/companies/company-uuid-123/people", nil)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)
		c.SetPath("/companies/:company_uuid/people")
		c.SetParamNames("company_uuid")
		c.SetParamValues("company-uuid-123")
		c.Set(infra.UserUUIDKey.String(), "user-uuid-456")

		mockCompanyService.EXPECT().ValidateCompanyOwnership(
			gomock.Any(),
			"company-uuid-123",
			"user-uuid-456",
		).Return(resterrors.NewNotFoundError("company not found"))

		err := middleware(func(c echo.Context) error {
			return nil
		})(c)

		assert.NotNil(t, err)
		assert.Equal(t, http.StatusNotFound, err.(resterrors.RestErr).StatusCode())
		assert.Equal(t, "company not found", err.(resterrors.RestErr).Message())
	})

	t.Run("Should return error when user doesn't own the company", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/companies/company-uuid-123/people", nil)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)
		c.SetPath("/companies/:company_uuid/people")
		c.SetParamNames("company_uuid")
		c.SetParamValues("company-uuid-123")
		c.Set(infra.UserUUIDKey.String(), "user-uuid-456")

		mockCompanyService.EXPECT().ValidateCompanyOwnership(
			gomock.Any(),
			"company-uuid-123",
			"user-uuid-456",
		).Return(resterrors.NewUnauthorizedError("you don't have permission to access this company"))

		err := middleware(func(c echo.Context) error {
			return nil
		})(c)

		assert.NotNil(t, err)
		assert.Equal(t, http.StatusUnauthorized, err.(resterrors.RestErr).StatusCode())
		assert.Equal(t, "you don't have permission to access this company", err.(resterrors.RestErr).Message())
	})

	t.Run("Should return error when company service returns internal error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/companies/company-uuid-123/people", nil)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)
		c.SetPath("/companies/:company_uuid/people")
		c.SetParamNames("company_uuid")
		c.SetParamValues("company-uuid-123")
		c.Set(infra.UserUUIDKey.String(), "user-uuid-456")

		mockCompanyService.EXPECT().ValidateCompanyOwnership(
			gomock.Any(),
			"company-uuid-123",
			"user-uuid-456",
		).Return(resterrors.NewInternalServerError("database connection failed"))

		err := middleware(func(c echo.Context) error {
			return nil
		})(c)

		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.(resterrors.RestErr).StatusCode())
		assert.Equal(t, "database connection failed", err.(resterrors.RestErr).Message())
	})

	t.Run("Should call next handler when middleware passes", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/companies/company-uuid-123/people", nil)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)
		c.SetPath("/companies/:company_uuid/people")
		c.SetParamNames("company_uuid")
		c.SetParamValues("company-uuid-123")
		c.Set(infra.UserUUIDKey.String(), "user-uuid-456")

		mockCompanyService.EXPECT().ValidateCompanyOwnership(
			gomock.Any(),
			"company-uuid-123",
			"user-uuid-456",
		).Return(nil)

		handlerCalled := false
		err := middleware(func(c echo.Context) error {
			handlerCalled = true
			return c.String(http.StatusOK, "success")
		})(c)

		assert.Nil(t, err)
		assert.True(t, handlerCalled)
		assert.Equal(t, "company-uuid-123", c.Get(infra.CompanyUUIDKey.String()))
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "success", rec.Body.String())
	})
}