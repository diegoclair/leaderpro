package servermiddleware

import (
	"github.com/diegoclair/go_utils/resterrors"
	"github.com/diegoclair/leaderpro/infra"
	"github.com/diegoclair/leaderpro/internal/domain/contract"
	echo "github.com/labstack/echo/v4"
)

// CompanyOwnershipMiddleware validates that the company_uuid in the route belongs to the authenticated user
func CompanyOwnershipMiddleware(companyService contract.CompanyApp) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			// Get company_uuid from path parameters
			companyUUID := ctx.Param("company_uuid")
			if companyUUID == "" {
				return resterrors.NewBadRequestError("company_uuid is required")
			}

			// Get authenticated user UUID from context (set by auth middleware)
			userUUID := ctx.Get(infra.UserUUIDKey.String())
			if userUUID == nil {
				return resterrors.NewUnauthorizedError("user not authenticated")
			}

			userUUIDStr, ok := userUUID.(string)
			if !ok {
				return resterrors.NewInternalServerError("invalid user context")
			}

			// Validate company ownership
			err := companyService.ValidateCompanyOwnership(ctx.Request().Context(), companyUUID, userUUIDStr)
			if err != nil {
				return err
			}

			// Add company UUID to context for handlers to use if needed
			ctx.Set(infra.CompanyUUIDKey.String(), companyUUID)

			return next(ctx)
		}
	}
}