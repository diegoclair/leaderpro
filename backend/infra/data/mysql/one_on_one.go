package mysql

import (
	"context"

	"github.com/diegoclair/go_utils/mysqlutils"
	"github.com/diegoclair/leaderpro/internal/domain/contract"
	"github.com/diegoclair/leaderpro/internal/domain/entity"
)

type oneOnOneRepo struct {
	db dbConn
}

func newOneOnOneRepo(db dbConn) contract.OneOnOneRepo {
	return &oneOnOneRepo{
		db: db,
	}
}

const oneOnOneSelectBase string = `
	SELECT 
		o.one_on_one_id,
		o.one_on_one_uuid,
		o.company_id,
		o.person_id,
		o.manager_id,
		o.scheduled_date,
		o.actual_date,
		o.duration,
		o.location,
		o.status,
		o.agenda,
		o.discussion_notes,
		o.action_items,
		o.private_notes,
		o.ai_context,
		o.ai_suggestions,
		o.created_at,
		o.updated_at,
		o.completed_at
	
	FROM tab_one_on_one o
`

func (r *oneOnOneRepo) parseOneOnOne(row scanner) (oneOnOne entity.OneOnOne, err error) {
	err = row.Scan(
		&oneOnOne.ID,
		&oneOnOne.UUID,
		&oneOnOne.CompanyID,
		&oneOnOne.PersonID,
		&oneOnOne.ManagerID,
		&oneOnOne.ScheduledDate,
		&oneOnOne.ActualDate,
		&oneOnOne.Duration,
		&oneOnOne.Location,
		&oneOnOne.Status,
		&oneOnOne.Agenda,
		&oneOnOne.DiscussionNotes,
		&oneOnOne.ActionItems,
		&oneOnOne.PrivateNotes,
		&oneOnOne.AIContext,
		&oneOnOne.AISuggestions,
		&oneOnOne.CreatedAt,
		&oneOnOne.UpdatedAt,
		&oneOnOne.CompletedAt,
	)

	if err != nil {
		return oneOnOne, err
	}

	return oneOnOne, nil
}

func (r *oneOnOneRepo) CreateOneOnOne(ctx context.Context, oneOnOne entity.OneOnOne) (createdID int64, err error) {
	query := `
		INSERT INTO tab_one_on_one (
			one_on_one_uuid,
			company_id,
			person_id,
			manager_id,
			scheduled_date,
			actual_date,
			duration,
			location,
			status,
			agenda,
			discussion_notes,
			action_items,
			private_notes,
			ai_context,
			ai_suggestions
		) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
	`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return createdID, mysqlutils.HandleMySQLError(err)
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx,
		oneOnOne.UUID,
		oneOnOne.CompanyID,
		oneOnOne.PersonID,
		oneOnOne.ManagerID,
		oneOnOne.ScheduledDate,
		oneOnOne.ActualDate,
		oneOnOne.Duration,
		oneOnOne.Location,
		oneOnOne.Status,
		oneOnOne.Agenda,
		oneOnOne.DiscussionNotes,
		oneOnOne.ActionItems,
		oneOnOne.PrivateNotes,
		oneOnOne.AIContext,
		oneOnOne.AISuggestions,
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

func (r *oneOnOneRepo) GetOneOnOneByUUID(ctx context.Context, oneOnOneUUID string) (oneOnOne entity.OneOnOne, err error) {
	query := oneOnOneSelectBase + `
		WHERE o.one_on_one_uuid = ?
	`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return oneOnOne, mysqlutils.HandleMySQLError(err)
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, oneOnOneUUID)
	oneOnOne, err = r.parseOneOnOne(row)
	if err != nil {
		return oneOnOne, mysqlutils.HandleMySQLError(err)
	}

	return oneOnOne, nil
}

func (r *oneOnOneRepo) GetOneOnOnesByPerson(ctx context.Context, personID int64, take, skip int64) (oneOnOnes []entity.OneOnOne, totalRecords int64, err error) {
	var params = []any{personID}

	query := oneOnOneSelectBase + `
		WHERE o.person_id = ?
		ORDER BY o.scheduled_date DESC
	`

	totalRecords, err = getTotalRecordsFromQuery(ctx, r.db, query, params...)
	if err != nil {
		return oneOnOnes, totalRecords, mysqlutils.HandleMySQLError(err)
	}

	if totalRecords < 1 {
		return oneOnOnes, totalRecords, nil
	}

	if take > 0 {
		query += ` LIMIT ?`
		params = append(params, take)
	}

	if skip > 0 {
		query += ` OFFSET ?`
		params = append(params, skip)
	}

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return oneOnOnes, totalRecords, mysqlutils.HandleMySQLError(err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, params...)
	if err != nil {
		return oneOnOnes, totalRecords, mysqlutils.HandleMySQLError(err)
	}
	defer rows.Close()

	for rows.Next() {
		oneOnOne, err := r.parseOneOnOne(rows)
		if err != nil {
			return oneOnOnes, totalRecords, mysqlutils.HandleMySQLError(err)
		}
		oneOnOnes = append(oneOnOnes, oneOnOne)
	}

	return oneOnOnes, totalRecords, nil
}

func (r *oneOnOneRepo) GetOneOnOnesByManager(ctx context.Context, managerID int64, take, skip int64) (oneOnOnes []entity.OneOnOne, totalRecords int64, err error) {
	var params = []any{managerID}

	query := oneOnOneSelectBase + `
		WHERE o.manager_id = ?
		ORDER BY o.scheduled_date DESC
	`

	totalRecords, err = getTotalRecordsFromQuery(ctx, r.db, query, params...)
	if err != nil {
		return oneOnOnes, totalRecords, mysqlutils.HandleMySQLError(err)
	}

	if totalRecords < 1 {
		return oneOnOnes, totalRecords, nil
	}

	if take > 0 {
		query += ` LIMIT ?`
		params = append(params, take)
	}

	if skip > 0 {
		query += ` OFFSET ?`
		params = append(params, skip)
	}

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return oneOnOnes, totalRecords, mysqlutils.HandleMySQLError(err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, params...)
	if err != nil {
		return oneOnOnes, totalRecords, mysqlutils.HandleMySQLError(err)
	}
	defer rows.Close()

	for rows.Next() {
		oneOnOne, err := r.parseOneOnOne(rows)
		if err != nil {
			return oneOnOnes, totalRecords, mysqlutils.HandleMySQLError(err)
		}
		oneOnOnes = append(oneOnOnes, oneOnOne)
	}

	return oneOnOnes, totalRecords, nil
}

func (r *oneOnOneRepo) UpdateOneOnOne(ctx context.Context, oneOnOneID int64, oneOnOne entity.OneOnOne) (err error) {
	query := `
		UPDATE tab_one_on_one
		SET 
			scheduled_date = ?,
			actual_date = ?,
			duration = ?,
			location = ?,
			status = ?,
			agenda = ?,
			discussion_notes = ?,
			action_items = ?,
			private_notes = ?,
			ai_context = ?,
			ai_suggestions = ?,
			updated_at = NOW(),
			completed_at = CASE WHEN ? = 'completed' THEN NOW() ELSE completed_at END
		WHERE one_on_one_id = ?
	`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return mysqlutils.HandleMySQLError(err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx,
		oneOnOne.ScheduledDate,
		oneOnOne.ActualDate,
		oneOnOne.Duration,
		oneOnOne.Location,
		oneOnOne.Status,
		oneOnOne.Agenda,
		oneOnOne.DiscussionNotes,
		oneOnOne.ActionItems,
		oneOnOne.PrivateNotes,
		oneOnOne.AIContext,
		oneOnOne.AISuggestions,
		oneOnOne.Status, // for the CASE statement
		oneOnOneID,
	)
	if err != nil {
		return mysqlutils.HandleMySQLError(err)
	}

	return nil
}

func (r *oneOnOneRepo) DeleteOneOnOne(ctx context.Context, oneOnOneID int64) (err error) {
	query := `
		DELETE FROM tab_one_on_one
		WHERE one_on_one_id = ?
	`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return mysqlutils.HandleMySQLError(err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, oneOnOneID)
	if err != nil {
		return mysqlutils.HandleMySQLError(err)
	}

	return nil
}

func (r *oneOnOneRepo) GetUpcomingOneOnOnes(ctx context.Context, managerID int64) (oneOnOnes []entity.OneOnOne, err error) {
	query := oneOnOneSelectBase + `
		WHERE o.manager_id = ? 
		AND o.status = 'scheduled'
		AND o.scheduled_date >= NOW()
		ORDER BY o.scheduled_date ASC
	`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return oneOnOnes, mysqlutils.HandleMySQLError(err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, managerID)
	if err != nil {
		return oneOnOnes, mysqlutils.HandleMySQLError(err)
	}
	defer rows.Close()

	for rows.Next() {
		oneOnOne, err := r.parseOneOnOne(rows)
		if err != nil {
			return oneOnOnes, mysqlutils.HandleMySQLError(err)
		}
		oneOnOnes = append(oneOnOnes, oneOnOne)
	}

	return oneOnOnes, nil
}

func (r *oneOnOneRepo) GetOverdueOneOnOnes(ctx context.Context, managerID int64) (oneOnOnes []entity.OneOnOne, err error) {
	query := oneOnOneSelectBase + `
		WHERE o.manager_id = ? 
		AND o.status = 'scheduled'
		AND o.scheduled_date < NOW()
		ORDER BY o.scheduled_date ASC
	`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return oneOnOnes, mysqlutils.HandleMySQLError(err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, managerID)
	if err != nil {
		return oneOnOnes, mysqlutils.HandleMySQLError(err)
	}
	defer rows.Close()

	for rows.Next() {
		oneOnOne, err := r.parseOneOnOne(rows)
		if err != nil {
			return oneOnOnes, mysqlutils.HandleMySQLError(err)
		}
		oneOnOnes = append(oneOnOnes, oneOnOne)
	}

	return oneOnOnes, nil
}