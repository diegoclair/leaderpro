package mysql

import (
	"context"
	"fmt"

	"github.com/diegoclair/go_utils/mysqlutils"
	"github.com/diegoclair/leaderpro/internal/domain/contract"
	"github.com/diegoclair/leaderpro/internal/domain/entity"
)

type aiRepo struct {
	db dbConn
}

func newAIRepo(db dbConn) contract.AIRepo {
	return &aiRepo{
		db: db,
	}
}

// ========== AI Prompts ==========

func (r *aiRepo) GetActivePromptByType(ctx context.Context, promptType string) (entity.AIPrompt, error) {
	query := `
		SELECT 	id, 
				type, 
				version,
				prompt,
				model,
				temperature,
				max_tokens,
				is_active,
				created_at,
				created_by

		FROM  ai_prompts
		WHERE type 		= ? 
		  AND is_active = TRUE
		ORDER BY version DESC
		LIMIT 1
	`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return entity.AIPrompt{}, mysqlutils.HandleMySQLError(err)
	}
	defer stmt.Close()

	var prompt entity.AIPrompt
	err = stmt.QueryRowContext(ctx, promptType).Scan(
		&prompt.ID,
		&prompt.Type,
		&prompt.Version,
		&prompt.Prompt,
		&prompt.Model,
		&prompt.Temperature,
		&prompt.MaxTokens,
		&prompt.IsActive,
		&prompt.CreatedAt,
		&prompt.CreatedBy,
	)
	if err != nil {
		return entity.AIPrompt{}, mysqlutils.HandleMySQLError(err)
	}

	return prompt, nil
}

// ========== AI Usage Tracker ==========

func (r *aiRepo) CreateUsage(ctx context.Context, usage entity.AIUsageTracker) (entity.AIUsageTracker, error) {
	query := `
		INSERT INTO ai_usage_tracker (
			user_id,
			company_id,
			prompt_id,
			person_id,
			request_type,
			tokens_used,
			cost_usd,
			response_time_ms
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return entity.AIUsageTracker{}, mysqlutils.HandleMySQLError(err)
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx,
		usage.UserID,
		usage.CompanyID,
		usage.PromptID,
		usage.PersonID,
		usage.RequestType,
		usage.TokensUsed,
		usage.CostUSD,
		usage.ResponseTimeMs,
	)
	if err != nil {
		return entity.AIUsageTracker{}, mysqlutils.HandleMySQLError(err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return entity.AIUsageTracker{}, err
	}

	usage.ID = id
	return usage, nil
}

func (r *aiRepo) UpdateUsageFeedback(ctx context.Context, usageID int64, feedback string, comment string) error {
	query := `
		UPDATE ai_usage_tracker
		  SET feedback 			= ?, 
		      feedback_comment 	= ?

		WHERE id = ?
	`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return mysqlutils.HandleMySQLError(err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, feedback, comment, usageID)
	if err != nil {
		return mysqlutils.HandleMySQLError(err)
	}

	return nil
}

func (r *aiRepo) GetUsageReport(ctx context.Context, userID int64, period string) (entity.AIUsageReport, error) {
	var whereClause string
	var args []interface{}

	// build WHERE clause based on time.
	switch period {
	case "today":
		whereClause = "WHERE user_id = ? AND DATE(created_at) = CURDATE()"
		args = []interface{}{userID}
	case "week":
		whereClause = "WHERE user_id = ? AND created_at >= DATE_SUB(NOW(), INTERVAL 7 DAY)"
		args = []interface{}{userID}
	case "month":
		whereClause = "WHERE user_id = ? AND created_at >= DATE_SUB(NOW(), INTERVAL 30 DAY)"
		args = []interface{}{userID}
	case "year":
		whereClause = "WHERE user_id = ? AND created_at >= DATE_SUB(NOW(), INTERVAL 365 DAY)"
		args = []interface{}{userID}
	default:
		whereClause = "WHERE user_id = ?"
		args = []interface{}{userID}
	}

	query := fmt.Sprintf(`
		SELECT 
			COUNT(*) as total_requests,
			COALESCE(SUM(tokens_used), 0) as total_tokens,
			COALESCE(SUM(cost_usd), 0) as total_cost_usd,
			COALESCE(AVG(response_time_ms), 0) as average_response_ms,
			COALESCE(SUM(CASE WHEN feedback = 'helpful' THEN 1 ELSE 0 END), 0) as helpful_feedback,
			COALESCE(SUM(CASE WHEN feedback = 'not_helpful' THEN 1 ELSE 0 END), 0) as unhelpful_feedback
		FROM ai_usage_tracker
		%s
	`, whereClause)

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return entity.AIUsageReport{}, mysqlutils.HandleMySQLError(err)
	}
	defer stmt.Close()

	var report entity.AIUsageReport
	err = stmt.QueryRowContext(ctx, args...).Scan(
		&report.TotalRequests,
		&report.TotalTokens,
		&report.TotalCostUSD,
		&report.AverageResponseMs,
		&report.HelpfulFeedback,
		&report.UnhelpfulFeedback,
	)
	if err != nil {
		return entity.AIUsageReport{}, mysqlutils.HandleMySQLError(err)
	}

	report.Period = period
	return report, nil
}

// ========== AI Conversations ==========

func (r *aiRepo) CreateConversation(ctx context.Context, conversation entity.AIConversation) (entity.AIConversation, error) {
	query := `
		INSERT INTO ai_conversations (
			usage_id,
			user_message,
			ai_response
		)
		VALUES (?, ?, ?)
	`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return entity.AIConversation{}, mysqlutils.HandleMySQLError(err)
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx,
		conversation.UsageID,
		conversation.UserMessage,
		conversation.AIResponse,
	)
	if err != nil {
		return entity.AIConversation{}, mysqlutils.HandleMySQLError(err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return entity.AIConversation{}, err
	}

	conversation.ID = id
	return conversation, nil
}
