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
		c.description,
		c.industry,
		c.size,
		c.created_at,
		c.updated_at,
		c.created_by,
		c.active
	
	FROM tab_company c
`

func (r *companyRepo) parseCompany(row scanner) (company entity.Company, err error) {
	err = row.Scan(
		&company.ID,
		&company.UUID,
		&company.Name,
		&company.Description,
		&company.Industry,
		&company.Size,
		&company.CreatedAt,
		&company.UpdatedAt,
		&company.CreatedBy,
		&company.Active,
	)

	if err != nil {
		return company, err
	}

	return company, nil
}

func (r *companyRepo) parseUserCompany(row scanner) (userCompany entity.UserCompany, err error) {
	var companyUserID, companyID, userID int64
	
	err = row.Scan(
		// Company fields
		&userCompany.Company.ID,
		&userCompany.Company.UUID,
		&userCompany.Company.Name,
		&userCompany.Company.Description,
		&userCompany.Company.Industry,
		&userCompany.Company.Size,
		&userCompany.Company.CreatedAt,
		&userCompany.Company.UpdatedAt,
		&userCompany.Company.CreatedBy,
		&userCompany.Company.Active,
		// UserCompany fields (we ignore IDs for UserCompany as they're embedded)
		&companyUserID,
		&companyID,
		&userID,
		&userCompany.Role,
		&userCompany.IsDefault,
		&userCompany.JoinedAt,
	)

	if err != nil {
		return userCompany, err
	}

	return userCompany, nil
}

func (r *companyRepo) CreateCompany(ctx context.Context, company entity.Company) (createdID int64, err error) {
	query := `
		INSERT INTO tab_company (
			company_uuid,
			name,
			description,
			industry,
			size,
			created_by,
			active
		) 
		VALUES (?, ?, ?, ?, ?, ?, ?);
	`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return createdID, mysqlutils.HandleMySQLError(err)
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx,
		company.UUID,
		company.Name,
		company.Description,
		company.Industry,
		company.Size,
		company.CreatedBy,
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

func (r *companyRepo) GetCompanyByUUID(ctx context.Context, companyUUID string) (company entity.Company, err error) {
	query := companySelectBase + `
		WHERE c.company_uuid = ? AND c.active = 1
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
	query := `
		SELECT 
			c.company_id,
			c.company_uuid,
			c.name,
			c.description,
			c.industry,
			c.size,
			c.created_at,
			c.updated_at,
			c.created_by,
			c.active
		
		FROM tab_company c
		INNER JOIN tab_company_user cu ON c.company_id = cu.company_id
		WHERE cu.user_id = ? AND c.active = 1
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

func (r *companyRepo) GetUserCompaniesWithDefault(ctx context.Context, userID int64) (companies []entity.UserCompany, err error) {
	query := `
		SELECT 
			c.company_id,
			c.company_uuid,
			c.name,
			c.description,
			c.industry,
			c.size,
			c.created_at,
			c.updated_at,
			c.created_by,
			c.active,
			cu.company_user_id,
			cu.company_id,
			cu.user_id,
			cu.role,
			cu.is_default,
			cu.joined_at
		FROM tab_company c
		INNER JOIN tab_company_user cu ON c.company_id = cu.company_id
		WHERE cu.user_id = ? AND c.active = 1
		ORDER BY cu.is_default DESC, c.created_at DESC
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
		userCompany, err := r.parseUserCompany(rows)
		if err != nil {
			return companies, mysqlutils.HandleMySQLError(err)
		}
		companies = append(companies, userCompany)
	}

	if err = rows.Err(); err != nil {
		return companies, mysqlutils.HandleMySQLError(err)
	}

	return companies, nil
}

func (r *companyRepo) UpdateCompany(ctx context.Context, companyID int64, company entity.Company) (err error) {
	query := `
		UPDATE tab_company
		SET 
			name = ?,
			description = ?,
			industry = ?,
			size = ?,
			updated_at = NOW()
		WHERE company_id = ? AND active = 1
	`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return mysqlutils.HandleMySQLError(err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx,
		company.Name,
		company.Description,
		company.Industry,
		company.Size,
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
		SET 
			active = 0,
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

func (r *companyRepo) AddUserToCompany(ctx context.Context, companyID, userID int64, role string) (err error) {
	query := `
		INSERT INTO tab_company_user (
			company_id,
			user_id,
			role
		) 
		VALUES (?, ?, ?)
		ON DUPLICATE KEY UPDATE
			role = VALUES(role)
	`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return mysqlutils.HandleMySQLError(err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, companyID, userID, role)
	if err != nil {
		return mysqlutils.HandleMySQLError(err)
	}

	return nil
}

// SetCompanyAsDefault sets a company as default for a user
func (r *companyRepo) SetCompanyAsDefault(ctx context.Context, companyID, userID int64) (err error) {
	// First, remove default flag from all user companies
	removeDefaultQuery := `
		UPDATE tab_company_user
		SET is_default = 0
		WHERE user_id = ?
	`

	stmt, err := r.db.PrepareContext(ctx, removeDefaultQuery)
	if err != nil {
		return mysqlutils.HandleMySQLError(err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, userID)
	if err != nil {
		return mysqlutils.HandleMySQLError(err)
	}

	// Then set the specified company as default
	setDefaultQuery := `
		UPDATE tab_company_user
		SET is_default = 1
		WHERE company_id = ? AND user_id = ?
	`

	stmt2, err := r.db.PrepareContext(ctx, setDefaultQuery)
	if err != nil {
		return mysqlutils.HandleMySQLError(err)
	}
	defer stmt2.Close()

	_, err = stmt2.ExecContext(ctx, companyID, userID)
	if err != nil {
		return mysqlutils.HandleMySQLError(err)
	}

	return nil
}

func (r *companyRepo) RemoveUserFromCompany(ctx context.Context, companyID, userID int64) (err error) {
	query := `
		DELETE FROM tab_company_user
		WHERE company_id = ? AND user_id = ?
	`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return mysqlutils.HandleMySQLError(err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, companyID, userID)
	if err != nil {
		return mysqlutils.HandleMySQLError(err)
	}

	return nil
}