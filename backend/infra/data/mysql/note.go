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

func (r *noteRepo) GetPersonTimeline(ctx context.Context, personID int64, filters entity.TimelineFilters, take, skip int64) (timeline []entity.UnifiedTimelineEntry, totalRecords int64, err error) {
	// Build base queries for timeline (direct notes) and mentions
	timelineBaseQuery := `
		SELECT 
			n.note_uuid as uuid,
			n.type,
			n.content,
			u.name as author_name,
			n.created_at,
			n.feedback_type,
			n.feedback_category,
			NULL as person_name,
			NULL as source_person_name,
			'direct' as entry_source
		FROM tab_note n
		INNER JOIN tab_user u ON n.user_id = u.user_id
		WHERE n.person_id = ? AND n.deleted_at IS NULL`
	
	mentionsBaseQuery := `
		SELECT 
			n.note_uuid as uuid,
			'mention' as type,
			n.content,
			u.name as author_name,
			n.created_at,
			n.feedback_type,
			n.feedback_category,
			p.name as person_name,
			sp.name as source_person_name,
			'mention' as entry_source
		FROM tab_note_mention nm
		INNER JOIN tab_note n ON nm.note_id = n.note_id
		INNER JOIN tab_user u ON n.user_id = u.user_id
		INNER JOIN tab_person p ON n.person_id = p.person_id
		LEFT JOIN tab_person sp ON nm.source_person_id = sp.person_id
		WHERE nm.mentioned_person_id = ? AND n.deleted_at IS NULL`
	
	timelineArgs := []interface{}{personID}
	mentionsArgs := []interface{}{personID}
	
	includeTimeline := true
	includeMentions := true
	
	// Apply filters
	if filters.SearchQuery != "" {
		searchCondition := ` AND (n.content LIKE ? OR u.name LIKE ? OR n.feedback_category LIKE ? OR n.feedback_type LIKE ?)`
		searchValue := "%" + filters.SearchQuery + "%"
		
		timelineBaseQuery += searchCondition
		mentionsBaseQuery += searchCondition
		
		timelineArgs = append(timelineArgs, searchValue, searchValue, searchValue, searchValue)
		mentionsArgs = append(mentionsArgs, searchValue, searchValue, searchValue, searchValue)
	}
	
	// Apply type filters
	if len(filters.Types) > 0 {
		var timelineTypes []string
		hasMentionType := false
		
		for _, t := range filters.Types {
			if t == "mention" {
				hasMentionType = true
			} else {
				timelineTypes = append(timelineTypes, t)
				timelineArgs = append(timelineArgs, t)
			}
		}
		
		// If no direct types, exclude timeline
		if len(timelineTypes) == 0 {
			includeTimeline = false
		} else {
			placeholders := make([]string, len(timelineTypes))
			for i := range placeholders {
				placeholders[i] = "?"
			}
			timelineBaseQuery += ` AND n.type IN (` + joinStringSlice(placeholders, ",") + `)`
		}
		
		// If no mention type, exclude mentions
		if !hasMentionType {
			includeMentions = false
		}
	}
	
	// Apply feedback type filters
	if len(filters.FeedbackTypes) > 0 {
		placeholders := make([]string, len(filters.FeedbackTypes))
		for i := range placeholders {
			placeholders[i] = "?"
		}
		feedbackFilter := ` AND n.feedback_type IN (` + joinStringSlice(placeholders, ",") + `)`
		
		if includeTimeline {
			timelineBaseQuery += feedbackFilter
			for _, s := range filters.FeedbackTypes {
				timelineArgs = append(timelineArgs, s)
			}
		}
		
		if includeMentions {
			mentionsBaseQuery += feedbackFilter
			for _, s := range filters.FeedbackTypes {
				mentionsArgs = append(mentionsArgs, s)
			}
		}
	}
	
	// Apply period filter
	if filters.Period != "" && filters.Period != "all" {
		var dateCondition string
		switch filters.Period {
		case "7d":
			dateCondition = ` AND n.created_at >= DATE_SUB(NOW(), INTERVAL 7 DAY)`
		case "30d":
			dateCondition = ` AND n.created_at >= DATE_SUB(NOW(), INTERVAL 30 DAY)`
		case "3m":
			dateCondition = ` AND n.created_at >= DATE_SUB(NOW(), INTERVAL 3 MONTH)`
		case "6m":
			dateCondition = ` AND n.created_at >= DATE_SUB(NOW(), INTERVAL 6 MONTH)`
		case "1y":
			dateCondition = ` AND n.created_at >= DATE_SUB(NOW(), INTERVAL 1 YEAR)`
		}
		
		if dateCondition != "" {
			if includeTimeline {
				timelineBaseQuery += dateCondition
			}
			if includeMentions {
				mentionsBaseQuery += dateCondition
			}
		}
	}
	
	// Build combined query
	var queriesParts []string
	var allArgs []interface{}
	
	if includeTimeline {
		queriesParts = append(queriesParts, "("+timelineBaseQuery+")")
		allArgs = append(allArgs, timelineArgs...)
	}
	
	if includeMentions {
		queriesParts = append(queriesParts, "("+mentionsBaseQuery+")")
		allArgs = append(allArgs, mentionsArgs...)
	}
	
	if len(queriesParts) == 0 {
		// No results due to filtering
		return timeline, 0, nil
	}
	
	combinedQuery := joinStringSlice(queriesParts, " UNION ALL ") + " ORDER BY created_at DESC"
	
	// Count query (simplified - just count both sources without complex filtering for now)
	countQuery := `
		SELECT (
			SELECT COUNT(*) FROM tab_note n 
			WHERE n.person_id = ? AND n.deleted_at IS NULL
		) + (
			SELECT COUNT(*) FROM tab_note_mention nm 
			INNER JOIN tab_note n ON nm.note_id = n.note_id
			WHERE nm.mentioned_person_id = ? AND n.deleted_at IS NULL
		) as total
	`
	
	// Execute count query
	stmt, err := r.db.PrepareContext(ctx, countQuery)
	if err != nil {
		return timeline, totalRecords, mysqlutils.HandleMySQLError(err)
	}
	defer stmt.Close()
	
	row := stmt.QueryRowContext(ctx, personID, personID)
	err = row.Scan(&totalRecords)
	if err != nil {
		return timeline, totalRecords, mysqlutils.HandleMySQLError(err)
	}
	
	// Add pagination
	if take > 0 {
		combinedQuery += ` LIMIT ?`
		allArgs = append(allArgs, take)
		if skip > 0 {
			combinedQuery += ` OFFSET ?`
			allArgs = append(allArgs, skip)
		}
	}
	
	// Execute main query
	stmt2, err := r.db.PrepareContext(ctx, combinedQuery)
	if err != nil {
		return timeline, totalRecords, mysqlutils.HandleMySQLError(err)
	}
	defer stmt2.Close()
	
	rows, err := stmt2.QueryContext(ctx, allArgs...)
	if err != nil {
		return timeline, totalRecords, mysqlutils.HandleMySQLError(err)
	}
	defer rows.Close()
	
	for rows.Next() {
		var entry entity.UnifiedTimelineEntry
		var feedbackType, feedbackCategory sql.NullString
		var personName, sourcePersonName sql.NullString
		
		err = rows.Scan(
			&entry.UUID, &entry.Type, &entry.Content,
			&entry.AuthorName, &entry.CreatedAt, &feedbackType,
			&feedbackCategory, &personName, &sourcePersonName,
			&entry.EntrySource,
		)
		if err != nil {
			return timeline, totalRecords, mysqlutils.HandleMySQLError(err)
		}
		
		// Handle nullable fields
		if feedbackType.Valid {
			entry.FeedbackType = &feedbackType.String
		}
		if feedbackCategory.Valid {
			entry.FeedbackCategory = &feedbackCategory.String
		}
		if personName.Valid {
			entry.PersonName = &personName.String
		}
		if sourcePersonName.Valid {
			entry.SourcePersonName = &sourcePersonName.String
		}
		
		timeline = append(timeline, entry)
	}
	
	if err = rows.Err(); err != nil {
		return timeline, totalRecords, mysqlutils.HandleMySQLError(err)
	}
	
	return timeline, totalRecords, nil
}

// Helper function to join string slice
func joinStringSlice(slice []string, separator string) string {
	if len(slice) == 0 {
		return ""
	}
	if len(slice) == 1 {
		return slice[0]
	}
	
	result := slice[0]
	for i := 1; i < len(slice); i++ {
		result += separator + slice[i]
	}
	return result
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