package mysql

import (
	"context"

	"github.com/diegoclair/go_utils/mysqlutils"
	"github.com/diegoclair/leaderpro/internal/domain/contract"
	"github.com/diegoclair/leaderpro/internal/domain/entity"
)

type feedbackRepo struct {
	db dbConn
}

func newFeedbackRepo(db dbConn) contract.FeedbackRepo {
	return &feedbackRepo{
		db: db,
	}
}

const feedbackSelectBase string = `
	SELECT 
		f.feedback_id,
		f.feedback_uuid,
		f.company_id,
		f.person_id,
		f.given_by,
		f.one_on_one_id,
		f.type,
		f.category,
		f.content,
		f.mentioned_from,
		f.mentioned_date,
		f.created_at,
		f.updated_at,
		f.is_private
	
	FROM tab_feedback f
`

func (r *feedbackRepo) parseFeedback(row scanner) (feedback entity.Feedback, err error) {
	err = row.Scan(
		&feedback.ID,
		&feedback.UUID,
		&feedback.CompanyID,
		&feedback.PersonID,
		&feedback.GivenBy,
		&feedback.OneOnOneID,
		&feedback.Type,
		&feedback.Category,
		&feedback.Content,
		&feedback.MentionedFrom,
		&feedback.MentionedDate,
		&feedback.CreatedAt,
		&feedback.UpdatedAt,
		&feedback.IsPrivate,
	)

	if err != nil {
		return feedback, err
	}

	return feedback, nil
}

func (r *feedbackRepo) CreateFeedback(ctx context.Context, feedback entity.Feedback) (createdID int64, err error) {
	query := `
		INSERT INTO tab_feedback (
			feedback_uuid,
			company_id,
			person_id,
			given_by,
			one_on_one_id,
			type,
			category,
			content,
			mentioned_from,
			mentioned_date,
			is_private
		) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
	`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return createdID, mysqlutils.HandleMySQLError(err)
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx,
		feedback.UUID,
		feedback.CompanyID,
		feedback.PersonID,
		feedback.GivenBy,
		feedback.OneOnOneID,
		feedback.Type,
		feedback.Category,
		feedback.Content,
		feedback.MentionedFrom,
		feedback.MentionedDate,
		feedback.IsPrivate,
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

func (r *feedbackRepo) GetFeedbackByUUID(ctx context.Context, feedbackUUID string) (feedback entity.Feedback, err error) {
	query := feedbackSelectBase + `
		WHERE f.feedback_uuid = ?
	`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return feedback, mysqlutils.HandleMySQLError(err)
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, feedbackUUID)
	feedback, err = r.parseFeedback(row)
	if err != nil {
		return feedback, mysqlutils.HandleMySQLError(err)
	}

	return feedback, nil
}

func (r *feedbackRepo) GetFeedbackByPerson(ctx context.Context, personID int64, take, skip int64) (feedback []entity.Feedback, totalRecords int64, err error) {
	var params = []any{personID}

	query := feedbackSelectBase + `
		WHERE f.person_id = ?
		ORDER BY f.created_at DESC
	`

	totalRecords, err = getTotalRecordsFromQuery(ctx, r.db, query, params...)
	if err != nil {
		return feedback, totalRecords, mysqlutils.HandleMySQLError(err)
	}

	if totalRecords < 1 {
		return feedback, totalRecords, nil
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
		return feedback, totalRecords, mysqlutils.HandleMySQLError(err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, params...)
	if err != nil {
		return feedback, totalRecords, mysqlutils.HandleMySQLError(err)
	}
	defer rows.Close()

	for rows.Next() {
		fb, err := r.parseFeedback(rows)
		if err != nil {
			return feedback, totalRecords, mysqlutils.HandleMySQLError(err)
		}
		feedback = append(feedback, fb)
	}

	return feedback, totalRecords, nil
}

func (r *feedbackRepo) GetFeedbackByGiver(ctx context.Context, giverID int64, take, skip int64) (feedback []entity.Feedback, totalRecords int64, err error) {
	var params = []any{giverID}

	query := feedbackSelectBase + `
		WHERE f.given_by = ?
		ORDER BY f.created_at DESC
	`

	totalRecords, err = getTotalRecordsFromQuery(ctx, r.db, query, params...)
	if err != nil {
		return feedback, totalRecords, mysqlutils.HandleMySQLError(err)
	}

	if totalRecords < 1 {
		return feedback, totalRecords, nil
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
		return feedback, totalRecords, mysqlutils.HandleMySQLError(err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, params...)
	if err != nil {
		return feedback, totalRecords, mysqlutils.HandleMySQLError(err)
	}
	defer rows.Close()

	for rows.Next() {
		fb, err := r.parseFeedback(rows)
		if err != nil {
			return feedback, totalRecords, mysqlutils.HandleMySQLError(err)
		}
		feedback = append(feedback, fb)
	}

	return feedback, totalRecords, nil
}

func (r *feedbackRepo) UpdateFeedback(ctx context.Context, feedbackID int64, feedback entity.Feedback) (err error) {
	query := `
		UPDATE tab_feedback
		SET 
			type = ?,
			category = ?,
			content = ?,
			mentioned_from = ?,
			mentioned_date = ?,
			is_private = ?,
			updated_at = NOW()
		WHERE feedback_id = ?
	`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return mysqlutils.HandleMySQLError(err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx,
		feedback.Type,
		feedback.Category,
		feedback.Content,
		feedback.MentionedFrom,
		feedback.MentionedDate,
		feedback.IsPrivate,
		feedbackID,
	)
	if err != nil {
		return mysqlutils.HandleMySQLError(err)
	}

	return nil
}

func (r *feedbackRepo) DeleteFeedback(ctx context.Context, feedbackID int64) (err error) {
	query := `
		DELETE FROM tab_feedback
		WHERE feedback_id = ?
	`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return mysqlutils.HandleMySQLError(err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, feedbackID)
	if err != nil {
		return mysqlutils.HandleMySQLError(err)
	}

	return nil
}

func (r *feedbackRepo) GetFeedbackSummary(ctx context.Context, personID int64, period string) (summary entity.FeedbackSummary, err error) {
	// Get basic counts
	query := `
		SELECT 
			COUNT(*) as total_count,
			SUM(CASE WHEN type = 'positive' THEN 1 ELSE 0 END) as positive_count,
			SUM(CASE WHEN type = 'constructive' THEN 1 ELSE 0 END) as constructive_count
		FROM tab_feedback
		WHERE person_id = ?
		AND created_at >= ?
		AND created_at <= ?
	`

	// Note: This is a simplified implementation. In a real application, you'd need to
	// parse the period parameter and convert it to proper date ranges.
	// For now, we'll use a basic implementation.
	
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return summary, mysqlutils.HandleMySQLError(err)
	}
	defer stmt.Close()

	// Using NOW() as placeholder - in real implementation, parse period parameter
	row := stmt.QueryRowContext(ctx, personID, "2024-01-01", "2024-12-31")
	
	err = row.Scan(&summary.TotalCount, &summary.PositiveCount, &summary.Constructive)
	if err != nil {
		return summary, mysqlutils.HandleMySQLError(err)
	}

	summary.PersonID = personID
	summary.Period = period

	// Get top categories
	categoryQuery := `
		SELECT category, COUNT(*) as count
		FROM tab_feedback
		WHERE person_id = ?
		AND created_at >= ?
		AND created_at <= ?
		GROUP BY category
		ORDER BY count DESC
		LIMIT 5
	`

	stmt2, err := r.db.PrepareContext(ctx, categoryQuery)
	if err != nil {
		return summary, mysqlutils.HandleMySQLError(err)
	}
	defer stmt2.Close()

	rows, err := stmt2.QueryContext(ctx, personID, "2024-01-01", "2024-12-31")
	if err != nil {
		return summary, mysqlutils.HandleMySQLError(err)
	}
	defer rows.Close()

	for rows.Next() {
		var category string
		var count int
		err = rows.Scan(&category, &count)
		if err != nil {
			return summary, mysqlutils.HandleMySQLError(err)
		}
		summary.TopCategories = append(summary.TopCategories, category)
	}

	// Get highlights (sample feedback content)
	highlightQuery := `
		SELECT content
		FROM tab_feedback
		WHERE person_id = ?
		AND created_at >= ?
		AND created_at <= ?
		AND type = 'positive'
		ORDER BY created_at DESC
		LIMIT 3
	`

	stmt3, err := r.db.PrepareContext(ctx, highlightQuery)
	if err != nil {
		return summary, mysqlutils.HandleMySQLError(err)
	}
	defer stmt3.Close()

	rows2, err := stmt3.QueryContext(ctx, personID, "2024-01-01", "2024-12-31")
	if err != nil {
		return summary, mysqlutils.HandleMySQLError(err)
	}
	defer rows2.Close()

	for rows2.Next() {
		var content string
		err = rows2.Scan(&content)
		if err != nil {
			return summary, mysqlutils.HandleMySQLError(err)
		}
		summary.Highlights = append(summary.Highlights, content)
	}

	return summary, nil
}