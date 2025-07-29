package mysql

import (
	"context"

	"github.com/diegoclair/go_utils/mysqlutils"
	"github.com/diegoclair/leaderpro/internal/domain/contract"
	"github.com/diegoclair/leaderpro/internal/domain/entity"
)

type companyRepo struct {
	db dbConn
}

func newCompanyRepo(db dbConn) contract.CompanyRepo {
	return &companyRepo{
		db: db,
	}
}

const companySelectBase string = `
	SELECT 
		c.company_id,
		c.company_uuid,
		c.name,
		c.industry,
		c.size,
		c.role,
		c.is_default,
		c.created_at,
		c.updated_at,
		c.user_owner_id,
		c.active
	
	FROM tab_company c
`

func (r *companyRepo) parseCompany(row scanner) (company entity.Company, err error) {
	err = row.Scan(
		&company.ID,
		&company.UUID,
		&company.Name,
		&company.Industry,
		&company.Size,
		&company.Role,
		&company.IsDefault,
		&company.CreatedAt,
		&company.UpdatedAt,
		&company.UserOwnerID,
		&company.Active,
	)

	if err != nil {
		return company, err
	}

	return company, nil
}


func (r *companyRepo) CreateCompany(ctx context.Context, company entity.Company) (createdID int64, err error) {
	query := `
		INSERT INTO tab_company (
			company_uuid,
			name,
			industry,
			size,
			role,
			is_default,
			user_owner_id,
			active
		) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?);
	`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return createdID, mysqlutils.HandleMySQLError(err)
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx,
		company.UUID,
		company.Name,
		company.Industry,
		company.Size,
		company.Role,
		company.IsDefault,
		company.UserOwnerID,
		company.Active,
	)
	if err != nil {
		return createdID, mysqlutils.HandleMySQLError(err)
	}

	createdID, err = result.LastInsertId()
	if err != nil {
		return createdID, mysqlutils.HandleMySQLError(err)
	}

	return createdID, nil
}

func (r *companyRepo) GetCompanyByID(ctx context.Context, companyID int64) (company entity.Company, err error) {
	query := companySelectBase + `
		WHERE c.company_id = ?
		  AND c.active     = 1
	`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return company, mysqlutils.HandleMySQLError(err)
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, companyID)
	company, err = r.parseCompany(row)
	if err != nil {
		return company, mysqlutils.HandleMySQLError(err)
	}

	return company, nil
}

func (r *companyRepo) GetCompanyByUUID(ctx context.Context, companyUUID string) (company entity.Company, err error) {
	query := companySelectBase + `
		WHERE c.company_uuid = ?
		  AND c.active       = 1
	`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return company, mysqlutils.HandleMySQLError(err)
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, companyUUID)
	company, err = r.parseCompany(row)
	if err != nil {
		return company, mysqlutils.HandleMySQLError(err)
	}

	return company, nil
}

func (r *companyRepo) GetCompaniesByUser(ctx context.Context, userID int64) (companies []entity.Company, err error) {
	query := companySelectBase + `
		WHERE c.user_owner_id = ?
		  AND c.active        = 1
		ORDER BY c.created_at DESC
	`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return companies, mysqlutils.HandleMySQLError(err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, userID)
	if err != nil {
		return companies, mysqlutils.HandleMySQLError(err)
	}
	defer rows.Close()

	for rows.Next() {
		company, err := r.parseCompany(rows)
		if err != nil {
			return companies, mysqlutils.HandleMySQLError(err)
		}
		companies = append(companies, company)
	}

	if err = rows.Err(); err != nil {
		return companies, mysqlutils.HandleMySQLError(err)
	}

	return companies, nil
}


func (r *companyRepo) UpdateCompany(ctx context.Context, companyID int64, company entity.Company) (err error) {
	query := `
		UPDATE tab_company
		  SET  name       = ?,
		       industry   = ?,
		       size       = ?,
		       role       = ?,
		       is_default = ?,
		       updated_at = NOW()

		WHERE company_id = ?
		  AND active     = 1
	`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return mysqlutils.HandleMySQLError(err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx,
		company.Name,
		company.Industry,
		company.Size,
		company.Role,
		company.IsDefault,
		companyID,
	)
	if err != nil {
		return mysqlutils.HandleMySQLError(err)
	}

	return nil
}

func (r *companyRepo) DeleteCompany(ctx context.Context, companyID int64) (err error) {
	query := `
		UPDATE tab_company
		  SET  active     = 0,
		       updated_at = NOW()

		WHERE company_id = ?
	`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return mysqlutils.HandleMySQLError(err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, companyID)
	if err != nil {
		return mysqlutils.HandleMySQLError(err)
	}

	return nil
}

