package mysql

import (
	"context"
	"database/sql"
	"testing"

	"github.com/diegoclair/leaderpro/internal/domain/entity"
	"github.com/stretchr/testify/require"
	"github.com/twinj/uuid"
)

func validateTwoCompanies(t *testing.T, companyExpected entity.Company, companyToCompare entity.Company) {
	require.NotZero(t, companyToCompare.ID)
	require.Equal(t, companyExpected.UUID, companyToCompare.UUID)
	require.Equal(t, companyExpected.Name, companyToCompare.Name)
	require.Equal(t, companyExpected.Industry, companyToCompare.Industry)
	require.Equal(t, companyExpected.Size, companyToCompare.Size)
	require.Equal(t, companyExpected.Role, companyToCompare.Role)
	require.Equal(t, companyExpected.IsDefault, companyToCompare.IsDefault)
	require.Equal(t, companyExpected.UserOwnerID, companyToCompare.UserOwnerID)
	require.Equal(t, companyExpected.Active, companyToCompare.Active)
}


func createRandomCompany(t *testing.T) entity.Company {
	// Create a user first to be the creator
	user := createRandomUserForTests(t)

	company := entity.Company{
		UUID:        uuid.NewV4().String(),
		Name:        "Test Company",
		Industry:    "Technology",
		Size:        "small",
		Role:        "Tech Lead",
		IsDefault:   true,
		UserOwnerID: user.ID,
		Active:      true,
	}

	companyID, err := testMysql.Company().CreateCompany(context.Background(), company)
	require.NoError(t, err)
	require.NotZero(t, companyID)

	company.ID = companyID
	return company
}

func TestCreateCompany(t *testing.T) {
	ctx := context.Background()
	user := createRandomUserForTests(t)

	company := entity.Company{
		UUID:        uuid.NewV4().String(),
		Name:        "New Company",
		Industry:    "Finance",
		Size:        "medium",
		Role:        "Manager",
		IsDefault:   false,
		UserOwnerID: user.ID,
		Active:      true,
	}

	companyID, err := testMysql.Company().CreateCompany(ctx, company)
	require.NoError(t, err)
	require.NotZero(t, companyID)
}

func TestGetCompanyByUUID(t *testing.T) {
	ctx := context.Background()
	company := createRandomCompany(t)

	company2, err := testMysql.Company().GetCompanyByUUID(ctx, company.UUID)
	require.NoError(t, err)
	require.NotEmpty(t, company2)
	validateTwoCompanies(t, company, company2)
}

func TestGetCompaniesByUser(t *testing.T) {
	ctx := context.Background()
	user := createRandomUserForTests(t)

	// Create multiple companies for the user
	company1 := entity.Company{
		UUID:        uuid.NewV4().String(),
		Name:        "Company 1",
		Industry:    "Technology",
		Size:        "small",
		Role:        "CTO",
		IsDefault:   true,
		UserOwnerID: user.ID,
		Active:      true,
	}

	company2 := entity.Company{
		UUID:        uuid.NewV4().String(),
		Name:        "Company 2",
		Industry:    "Healthcare",
		Size:        "large",
		Role:        "Lead Developer",
		IsDefault:   false,
		UserOwnerID: user.ID,
		Active:      true,
	}

	// Create companies
	companyID1, err := testMysql.Company().CreateCompany(ctx, company1)
	require.NoError(t, err)
	require.NotZero(t, companyID1)

	companyID2, err := testMysql.Company().CreateCompany(ctx, company2)
	require.NoError(t, err)
	require.NotZero(t, companyID2)

	// Companies are already owned by the user through UserOwnerID

	// Get companies by user
	companies, err := testMysql.Company().GetCompaniesByUser(ctx, user.ID)
	require.NoError(t, err)
	require.Len(t, companies, 2)

	// Verify we got both companies
	companyUUIDs := make(map[string]bool)
	for _, company := range companies {
		companyUUIDs[company.UUID] = true
	}
	require.True(t, companyUUIDs[company1.UUID])
	require.True(t, companyUUIDs[company2.UUID])
}

func TestUpdateCompany(t *testing.T) {
	ctx := context.Background()
	company := createRandomCompany(t)

	// Update company fields
	updatedCompany := entity.Company{
		Name:      "Updated Company Name",
		Industry:  "Updated Industry",
		Size:      "enterprise",
		Role:      "Senior Manager",
		IsDefault: false,
	}

	err := testMysql.Company().UpdateCompany(ctx, company.ID, updatedCompany)
	require.NoError(t, err)

	// Get updated company and validate
	retrievedCompany, err := testMysql.Company().GetCompanyByUUID(ctx, company.UUID)
	require.NoError(t, err)
	require.Equal(t, updatedCompany.Name, retrievedCompany.Name)
	require.Equal(t, updatedCompany.Industry, retrievedCompany.Industry)
	require.Equal(t, updatedCompany.Size, retrievedCompany.Size)
}

func TestDeleteCompany(t *testing.T) {
	ctx := context.Background()
	company := createRandomCompany(t)

	err := testMysql.Company().DeleteCompany(ctx, company.ID)
	require.NoError(t, err)

	// Try to get the deleted company - should return error (not found)
	_, err = testMysql.Company().GetCompanyByUUID(ctx, company.UUID)
	require.Error(t, err)
}


// Error tests with mocks
func TestCreateCompanyErrorsWithMock(t *testing.T) {
	testForInsertErrorsWithMock(t, func(db *sql.DB) error {
		_, err := newCompanyRepo(db).CreateCompany(context.Background(), entity.Company{})
		return err
	})
}

func TestGetCompanyByUUIDErrorsWithMock(t *testing.T) {
	testForSelectErrorsWithMock(t, "company_id", func(db *sql.DB) error {
		_, err := newCompanyRepo(db).GetCompanyByUUID(context.Background(), "company-uuid")
		return err
	})
}

func TestGetCompaniesByUserErrorsWithMock(t *testing.T) {
	testForSelectErrorsWithMock(t, "company_id", func(db *sql.DB) error {
		_, err := newCompanyRepo(db).GetCompaniesByUser(context.Background(), 1)
		return err
	})
}

func TestUpdateCompanyErrorsWithMock(t *testing.T) {
	testForUpdateDeleteErrorsWithMock(t, func(db *sql.DB) error {
		return newCompanyRepo(db).UpdateCompany(context.Background(), 1, entity.Company{})
	})
}

func TestDeleteCompanyErrorsWithMock(t *testing.T) {
	testForUpdateDeleteErrorsWithMock(t, func(db *sql.DB) error {
		return newCompanyRepo(db).DeleteCompany(context.Background(), 1)
	})
}

