package mysql

import (
	"context"
	"database/sql"
	"time"

	"github.com/diegoclair/go_utils/mysqlutils"
	"github.com/diegoclair/leaderpro/internal/domain"
	"github.com/diegoclair/leaderpro/internal/domain/contract"
	"github.com/diegoclair/leaderpro/internal/domain/entity"
)

type noteRepo struct {
	db dbConn
}

func newNoteRepo(db dbConn) contract.NoteRepo {
	return &noteRepo{
		db: db,
	}
}

func (r *noteRepo) CreateNote(ctx context.Context, note entity.Note) (createdID int64, err error) {
	query := `
		INSERT INTO tab_note (
			note_uuid, company_id, person_id, user_id, type, content, 
			feedback_type, feedback_category, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return createdID, mysqlutils.HandleMySQLError(err)
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx,
		note.UUID, note.CompanyID, note.PersonID, note.UserID, note.Type,
		note.Content, note.FeedbackType, note.FeedbackCategory,
		note.CreatedAt, note.UpdatedAt,
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

func (r *noteRepo) GetNoteByUUID(ctx context.Context, noteUUID string) (note entity.Note, err error) {
	query := `
		SELECT note_id, note_uuid, company_id, person_id, user_id, type, content,
			   feedback_type, feedback_category, created_at, updated_at
		FROM tab_note 
		WHERE note_uuid = ? AND deleted_at IS NULL
	`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return note, mysqlutils.HandleMySQLError(err)
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, noteUUID)
	err = row.Scan(
		&note.ID, &note.UUID, &note.CompanyID, &note.PersonID, &note.UserID,
		&note.Type, &note.Content, &note.FeedbackType, &note.FeedbackCategory,
		&note.CreatedAt, &note.UpdatedAt,
	)
	if err != nil {
		return note, mysqlutils.HandleMySQLError(err)
	}

	return note, nil
}

func (r *noteRepo) GetNotesByPerson(ctx context.Context, personID int64, take, skip int64) (notes []entity.Note, totalRecords int64, err error) {
	// Count query
	countQuery := `
		SELECT COUNT(*) 
		FROM tab_note 
		WHERE person_id = ? AND deleted_at IS NULL
	`

	stmt, err := r.db.PrepareContext(ctx, countQuery)
	if err != nil {
		return notes, 0, mysqlutils.HandleMySQLError(err)
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, personID)
	err = row.Scan(&totalRecords)
	if err != nil {
		return notes, 0, mysqlutils.HandleMySQLError(err)
	}

	// Data query
	query := `
		SELECT note_id, note_uuid, company_id, person_id, user_id, type, content,
			   feedback_type, feedback_category, created_at, updated_at
		FROM tab_note 
		WHERE person_id = ? AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	stmt2, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return notes, totalRecords, mysqlutils.HandleMySQLError(err)
	}
	defer stmt2.Close()

	rows, err := stmt2.QueryContext(ctx, personID, take, skip)
	if err != nil {
		return notes, totalRecords, mysqlutils.HandleMySQLError(err)
	}
	defer rows.Close()

	for rows.Next() {
		var note entity.Note
		err = rows.Scan(
			&note.ID, &note.UUID, &note.CompanyID, &note.PersonID, &note.UserID,
			&note.Type, &note.Content, &note.FeedbackType, &note.FeedbackCategory,
			&note.CreatedAt, &note.UpdatedAt,
		)
		if err != nil {
			return notes, totalRecords, mysqlutils.HandleMySQLError(err)
		}
		notes = append(notes, note)
	}

	if err = rows.Err(); err != nil {
		return notes, totalRecords, mysqlutils.HandleMySQLError(err)
	}

	return notes, totalRecords, nil
}

func (r *noteRepo) UpdateNote(ctx context.Context, noteID int64, note entity.Note) (err error) {
	query := `
		UPDATE tab_note 
		SET type = ?, content = ?, feedback_type = ?, feedback_category = ?, updated_at = ?
		WHERE note_id = ? AND deleted_at IS NULL
	`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return mysqlutils.HandleMySQLError(err)
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx,
		note.Type, note.Content, note.FeedbackType, note.FeedbackCategory,
		note.UpdatedAt, noteID,
	)
	if err != nil {
		return mysqlutils.HandleMySQLError(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *noteRepo) DeleteNote(ctx context.Context, noteID int64) (err error) {
	query := `
		UPDATE tab_note 
		SET deleted_at = NOW()
		WHERE note_id = ? AND deleted_at IS NULL
	`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return mysqlutils.HandleMySQLError(err)
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, noteID)
	if err != nil {
		return mysqlutils.HandleMySQLError(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *noteRepo) CreateNoteMention(ctx context.Context, mention entity.NoteMention) (createdID int64, err error) {
	query := `
		INSERT INTO tab_note_mention (
			mention_uuid, note_id, mentioned_person_id, source_person_id, 
			full_content, created_at
		) VALUES (?, ?, ?, ?, ?, ?)
	`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return createdID, mysqlutils.HandleMySQLError(err)
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx,
		mention.UUID, mention.NoteID, mention.MentionedPersonID,
		mention.SourcePersonID, mention.FullContent, mention.CreatedAt,
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

func (r *noteRepo) GetMentionsByPerson(ctx context.Context, mentionedPersonID int64, take, skip int64) (mentions []entity.NoteMention, totalRecords int64, err error) {
	// Count query
	countQuery := `
		SELECT COUNT(*) 
		FROM tab_note_mention 
		WHERE mentioned_person_id = ?
	`

	stmt, err := r.db.PrepareContext(ctx, countQuery)
	if err != nil {
		return mentions, 0, mysqlutils.HandleMySQLError(err)
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, mentionedPersonID)
	err = row.Scan(&totalRecords)
	if err != nil {
		return mentions, 0, mysqlutils.HandleMySQLError(err)
	}

	// Data query
	query := `
		SELECT mention_id, mention_uuid, note_id, mentioned_person_id, 
			   source_person_id, full_content, created_at
		FROM tab_note_mention 
		WHERE mentioned_person_id = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	stmt2, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return mentions, totalRecords, mysqlutils.HandleMySQLError(err)
	}
	defer stmt2.Close()

	rows, err := stmt2.QueryContext(ctx, mentionedPersonID, take, skip)
	if err != nil {
		return mentions, totalRecords, mysqlutils.HandleMySQLError(err)
	}
	defer rows.Close()

	for rows.Next() {
		var mention entity.NoteMention
		err = rows.Scan(
			&mention.ID, &mention.UUID, &mention.NoteID, &mention.MentionedPersonID,
			&mention.SourcePersonID, &mention.FullContent, &mention.CreatedAt,
		)
		if err != nil {
			return mentions, totalRecords, mysqlutils.HandleMySQLError(err)
		}
		mentions = append(mentions, mention)
	}

	if err = rows.Err(); err != nil {
		return mentions, totalRecords, mysqlutils.HandleMySQLError(err)
	}

	return mentions, totalRecords, nil
}

func (r *noteRepo) GetPersonTimeline(ctx context.Context, personID int64, take, skip int64) (timeline []entity.TimelineEntry, totalRecords int64, err error) {
	// Count query - count only notes directly about this person (excluding mentions/feedbacks)
	countQuery := `
		SELECT COUNT(*) 
		FROM tab_note 
		WHERE person_id = ? AND deleted_at IS NULL
	`

	stmt, err := r.db.PrepareContext(ctx, countQuery)
	if err != nil {
		return timeline, 0, mysqlutils.HandleMySQLError(err)
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, personID)
	err = row.Scan(&totalRecords)
	if err != nil {
		return timeline, 0, mysqlutils.HandleMySQLError(err)
	}

	// Main query - timeline with only direct notes (1:1s and observations, excluding mentions/feedbacks)
	query := `
		SELECT 
			n.note_uuid as uuid,
			n.type,
			n.content,
			u.name as author_name,
			n.created_at,
			n.feedback_type,
			n.feedback_category,
			NULL as source_person_name
		FROM tab_note n
		INNER JOIN tab_user u ON n.user_id = u.user_id
		WHERE n.person_id = ? AND n.deleted_at IS NULL
		ORDER BY n.created_at DESC
		LIMIT ? OFFSET ?
	`

	stmt2, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return timeline, totalRecords, mysqlutils.HandleMySQLError(err)
	}
	defer stmt2.Close()

	rows, err := stmt2.QueryContext(ctx, personID, take, skip)
	if err != nil {
		return timeline, totalRecords, mysqlutils.HandleMySQLError(err)
	}
	defer rows.Close()

	for rows.Next() {
		var entry entity.TimelineEntry
		err = rows.Scan(
			&entry.UUID, &entry.Type, &entry.Content,
			&entry.AuthorName, &entry.CreatedAt, &entry.FeedbackType,
			&entry.FeedbackCategory, &entry.SourcePersonName,
		)
		if err != nil {
			return timeline, totalRecords, mysqlutils.HandleMySQLError(err)
		}
		timeline = append(timeline, entry)
	}

	if err = rows.Err(); err != nil {
		return timeline, totalRecords, mysqlutils.HandleMySQLError(err)
	}

	return timeline, totalRecords, nil
}

func (r *noteRepo) GetPersonMentions(ctx context.Context, mentionedPersonID int64, take, skip int64) (mentions []entity.MentionEntry, totalRecords int64, err error) {
	// Count query - count notes where this person was mentioned
	countQuery := `
		SELECT COUNT(*)
		FROM tab_note_mention nm
		INNER JOIN tab_note n ON nm.note_id = n.note_id
		WHERE nm.mentioned_person_id = ? AND n.deleted_at IS NULL
	`

	stmt, err := r.db.PrepareContext(ctx, countQuery)
	if err != nil {
		return mentions, 0, mysqlutils.HandleMySQLError(err)
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, mentionedPersonID)
	err = row.Scan(&totalRecords)
	if err != nil {
		return mentions, 0, mysqlutils.HandleMySQLError(err)
	}

	// Main query - get notes where this person was mentioned
	query := `
		SELECT 
			n.note_uuid,
			n.type,
			n.content,
			n.feedback_type,
			n.feedback_category,
			n.created_at,
			p.person_uuid as person_id,
			p.name as person_name
		FROM tab_note_mention nm
		INNER JOIN tab_note n ON nm.note_id = n.note_id
		INNER JOIN tab_person p ON n.person_id = p.person_id
		WHERE nm.mentioned_person_id = ? AND n.deleted_at IS NULL
		ORDER BY n.created_at DESC
		LIMIT ? OFFSET ?
	`

	stmt2, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return mentions, totalRecords, mysqlutils.HandleMySQLError(err)
	}
	defer stmt2.Close()

	rows, err := stmt2.QueryContext(ctx, mentionedPersonID, take, skip)
	if err != nil {
		return mentions, totalRecords, mysqlutils.HandleMySQLError(err)
	}
	defer rows.Close()

	for rows.Next() {
		var mention entity.MentionEntry
		var feedbackType, feedbackCategory sql.NullString
		
		err = rows.Scan(
			&mention.UUID, &mention.Type, &mention.Content,
			&feedbackType, &feedbackCategory, &mention.CreatedAt,
			&mention.PersonID, &mention.PersonName,
		)
		if err != nil {
			return mentions, totalRecords, mysqlutils.HandleMySQLError(err)
		}
		
		// Handle nullable fields
		if feedbackType.Valid {
			mention.FeedbackType = &feedbackType.String
		}
		if feedbackCategory.Valid {
			mention.FeedbackCategory = &feedbackCategory.String
		}
		
		// TODO: In future, we can populate mention.Mentions with detailed mention info
		// For now, leave it empty as the frontend renders mentions from content parsing
		
		mentions = append(mentions, mention)
	}

	if err = rows.Err(); err != nil {
		return mentions, totalRecords, mysqlutils.HandleMySQLError(err)
	}

	return mentions, totalRecords, nil
}

func (r *noteRepo) DeleteMentionsByNote(ctx context.Context, noteID int64) (err error) {
	query := `DELETE FROM tab_note_mention WHERE note_id = ?`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return mysqlutils.HandleMySQLError(err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, noteID)
	if err != nil {
		return mysqlutils.HandleMySQLError(err)
	}

	return nil
}

func (r *noteRepo) GetOneOnOnesCountThisMonth(ctx context.Context, companyID int64) (count int64, err error) {
	query := `
		SELECT COUNT(*) 
		FROM tab_note n
		INNER JOIN tab_person p ON n.person_id = p.person_id
		WHERE p.company_id = ? 
		AND p.active = 1
		AND n.type = ?
		AND n.deleted_at IS NULL
		AND YEAR(n.created_at) = YEAR(NOW()) 
		AND MONTH(n.created_at) = MONTH(NOW())
	`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return count, mysqlutils.HandleMySQLError(err)
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, companyID, domain.NoteTypeOneOnOne)
	err = row.Scan(&count)
	if err != nil {
		return count, mysqlutils.HandleMySQLError(err)
	}

	return count, nil
}

func (r *noteRepo) GetAverageFrequencyDays(ctx context.Context, companyID int64) (avgDays float64, err error) {
	query := `
		SELECT COALESCE(AVG(day_diff), 0) as avg_frequency
		FROM (
			SELECT DATEDIFF(
				LEAD(n.created_at) OVER (PARTITION BY n.person_id ORDER BY n.created_at),
				n.created_at
			) as day_diff
			FROM tab_note n
			INNER JOIN tab_person p ON n.person_id = p.person_id
			WHERE p.company_id = ? 
			AND p.active = 1
			AND n.type = ?
			AND n.deleted_at IS NULL
		) as frequency_data
		WHERE day_diff IS NOT NULL
	`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return avgDays, mysqlutils.HandleMySQLError(err)
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, companyID, domain.NoteTypeOneOnOne)
	err = row.Scan(&avgDays)
	if err != nil {
		return avgDays, mysqlutils.HandleMySQLError(err)
	}

	return avgDays, nil
}

func (r *noteRepo) GetLastMeetingDate(ctx context.Context, companyID int64) (lastDate *time.Time, err error) {
	query := `
		SELECT MAX(n.created_at) as last_meeting_date
		FROM tab_note n
		INNER JOIN tab_person p ON n.person_id = p.person_id
		WHERE p.company_id = ? 
		AND p.active = 1
		AND n.type = ?
		AND n.deleted_at IS NULL
	`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return lastDate, mysqlutils.HandleMySQLError(err)
	}
	defer stmt.Close()

	var nullableDate sql.NullTime
	row := stmt.QueryRowContext(ctx, companyID, domain.NoteTypeOneOnOne)
	err = row.Scan(&nullableDate)
	if err != nil {
		return lastDate, mysqlutils.HandleMySQLError(err)
	}

	if nullableDate.Valid {
		lastDate = &nullableDate.Time
	}

	return lastDate, nil
}