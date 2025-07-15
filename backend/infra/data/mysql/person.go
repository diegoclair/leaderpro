package mysql

import (
	"context"

	"github.com/diegoclair/go_utils/mysqlutils"
	"github.com/diegoclair/leaderpro/internal/domain"
	"github.com/diegoclair/leaderpro/internal/domain/contract"
	"github.com/diegoclair/leaderpro/internal/domain/entity"
)

type personRepo struct {
	db dbConn
}

func newPersonRepo(db dbConn) contract.PersonRepo {
	return &personRepo{
		db: db,
	}
}

func getPersonSelectBase() string {
	return `
	SELECT 
		p.person_id,
		p.person_uuid,
		p.company_id,
		p.name,
		p.email,
		p.position,
		p.department,
		p.phone,
		p.birthday,
		p.start_date,
		p.is_manager,
		p.manager_id,
		p.notes,
		p.has_kids,
		p.gender,
		p.interests,
		p.personality,
		(
			SELECT MAX(n.created_at) 
			FROM tab_note n 
			WHERE n.person_id = p.person_id 
			AND n.type = '` + domain.NoteTypeOneOnOne + `' 
			AND n.deleted_at IS NULL
		) as last_one_on_one_date,
		p.created_at,
		p.updated_at,
		p.created_by,
		p.active
	
	FROM tab_person p
`
}

func (r *personRepo) parsePerson(row scanner) (person entity.Person, err error) {
	err = row.Scan(
		&person.ID,
		&person.UUID,
		&person.CompanyID,
		&person.Name,
		&person.Email,
		&person.Position,
		&person.Department,
		&person.Phone,
		&person.Birthday,
		&person.StartDate,
		&person.IsManager,
		&person.ManagerID,
		&person.Notes,
		&person.HasKids,
		&person.Gender,
		&person.Interests,
		&person.Personality,
		&person.LastOneOnOneDate,
		&person.CreatedAt,
		&person.UpdatedAt,
		&person.CreatedBy,
		&person.Active,
	)

	if err != nil {
		return person, err
	}

	return person, nil
}

func (r *personRepo) CreatePerson(ctx context.Context, person entity.Person) (createdID int64, err error) {
	query := `
		INSERT INTO tab_person (
			person_uuid,
			company_id,
			name,
			email,
			position,
			department,
			phone,
			birthday,
			start_date,
			is_manager,
			manager_id,
			notes,
			has_kids,
			gender,
			interests,
			personality,
			created_by,
			active
		) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
	`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return createdID, mysqlutils.HandleMySQLError(err)
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx,
		person.UUID,
		person.CompanyID,
		person.Name,
		person.Email,
		person.Position,
		person.Department,
		person.Phone,
		person.Birthday,
		person.StartDate,
		person.IsManager,
		person.ManagerID,
		person.Notes,
		person.HasKids,
		person.Gender,
		person.Interests,
		person.Personality,
		person.CreatedBy,
		person.Active,
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

func (r *personRepo) GetPersonByUUID(ctx context.Context, personUUID string) (person entity.Person, err error) {
	query := getPersonSelectBase() + `
		WHERE p.person_uuid = ? AND p.active = 1
	`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return person, mysqlutils.HandleMySQLError(err)
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, personUUID)
	person, err = r.parsePerson(row)
	if err != nil {
		return person, mysqlutils.HandleMySQLError(err)
	}

	return person, nil
}

func (r *personRepo) GetPersonsByCompany(ctx context.Context, companyID int64) (people []entity.Person, err error) {
	query := getPersonSelectBase() + `
		WHERE p.company_id = ? AND p.active = 1
		ORDER BY p.name ASC
	`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return people, mysqlutils.HandleMySQLError(err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, companyID)
	if err != nil {
		return people, mysqlutils.HandleMySQLError(err)
	}
	defer rows.Close()

	for rows.Next() {
		person, err := r.parsePerson(rows)
		if err != nil {
			return people, mysqlutils.HandleMySQLError(err)
		}
		people = append(people, person)
	}

	return people, nil
}

func (r *personRepo) UpdatePerson(ctx context.Context, personID int64, person entity.Person) (err error) {
	query := `
		UPDATE tab_person
		SET 
			name = ?,
			email = ?,
			position = ?,
			department = ?,
			phone = ?,
			birthday = ?,
			start_date = ?,
			is_manager = ?,
			manager_id = ?,
			notes = ?,
			has_kids = ?,
			gender = ?,
			interests = ?,
			personality = ?,
			updated_at = NOW()
		WHERE person_id = ? AND active = 1
	`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return mysqlutils.HandleMySQLError(err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx,
		person.Name,
		person.Email,
		person.Position,
		person.Department,
		person.Phone,
		person.Birthday,
		person.StartDate,
		person.IsManager,
		person.ManagerID,
		person.Notes,
		person.HasKids,
		person.Gender,
		person.Interests,
		person.Personality,
		personID,
	)
	if err != nil {
		return mysqlutils.HandleMySQLError(err)
	}

	return nil
}

func (r *personRepo) DeletePerson(ctx context.Context, personID int64) (err error) {
	query := `
		UPDATE tab_person
		SET 
			active = 0,
			updated_at = NOW()
		WHERE person_id = ?
	`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return mysqlutils.HandleMySQLError(err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, personID)
	if err != nil {
		return mysqlutils.HandleMySQLError(err)
	}

	return nil
}

func (r *personRepo) SearchPeople(ctx context.Context, companyID int64, search string) (people []entity.Person, err error) {
	query := getPersonSelectBase() + `
		WHERE p.company_id = ? AND p.active = 1
		AND (
			p.name LIKE ? OR 
			p.email LIKE ? OR 
			p.position LIKE ? OR 
			p.department LIKE ?
		)
		ORDER BY p.name ASC
	`

	searchTerm := "%" + search + "%"

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return people, mysqlutils.HandleMySQLError(err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, companyID, searchTerm, searchTerm, searchTerm, searchTerm)
	if err != nil {
		return people, mysqlutils.HandleMySQLError(err)
	}
	defer rows.Close()

	for rows.Next() {
		person, err := r.parsePerson(rows)
		if err != nil {
			return people, mysqlutils.HandleMySQLError(err)
		}
		people = append(people, person)
	}

	return people, nil
}

func (r *personRepo) GetPeopleCountByCompany(ctx context.Context, companyID int64) (count int64, err error) {
	query := `
		SELECT COUNT(*) 
		FROM tab_person p
		WHERE p.company_id = ? AND p.active = 1
	`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return count, mysqlutils.HandleMySQLError(err)
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, companyID)
	err = row.Scan(&count)
	if err != nil {
		return count, mysqlutils.HandleMySQLError(err)
	}

	return count, nil
}