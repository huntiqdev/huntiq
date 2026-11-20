package ai

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AIRepository interface {
	GetOrCreateConversation(ctx context.Context, sessionID string, userID *int, lang string) (*Conversation, error)
	GetConversation(ctx context.Context, sessionID string) (*Conversation, error)
	GetRecentMessages(ctx context.Context, conversationID int64, limit int) ([]Message, error)
	SaveMessage(ctx context.Context, msg *Message) error
	RecordAuditEvent(ctx context.Context, evt *AuditEvent) error
	ArchiveConversation(ctx context.Context, sessionID string, userID *int) error
	GetAdminMetrics(ctx context.Context) (*AIMetricsReport, error)
}

type PostgresAIRepo struct {
	db *pgxpool.Pool
}

func NewAIRepository(db *pgxpool.Pool) *PostgresAIRepo {
	return &PostgresAIRepo{db: db}
}

func (r *PostgresAIRepo) GetOrCreateConversation(ctx context.Context, sessionID string, userID *int, lang string) (*Conversation, error) {
	if sessionID == "" {
		return nil, errors.New("empty session id")
	}
	if !IsSupportedLanguage(lang) {
		lang = LangEN
	}

	conv := &Conversation{}
	err := r.db.QueryRow(ctx, `
		SELECT id, session_id, user_id, title, language_code, created_at, updated_at, expires_at, is_archived
		FROM ai_conversations
		WHERE session_id = $1 AND is_archived = FALSE
	`, sessionID).Scan(
		&conv.ID, &conv.SessionID, &conv.UserID, &conv.Title,
		&conv.LanguageCode, &conv.CreatedAt, &conv.UpdatedAt, &conv.ExpiresAt, &conv.IsArchived,
	)

	if err == nil {
		// Update user_id if authenticated later in session
		if userID != nil && conv.UserID == nil {
			_, _ = r.db.Exec(ctx, `UPDATE ai_conversations SET user_id = $1, updated_at = NOW() WHERE id = $2`, *userID, conv.ID)
			conv.UserID = userID
		}
		return conv, nil
	}

	if !errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("lookup ai conversation: %w", err)
	}

	// Create new conversation
	err = r.db.QueryRow(ctx, `
		INSERT INTO ai_conversations (session_id, user_id, title, language_code, expires_at)
		VALUES ($1, $2, 'HUNTIQ AI Discovery', $3, NOW() + INTERVAL '7 days')
		RETURNING id, session_id, user_id, title, language_code, created_at, updated_at, expires_at, is_archived
	`, sessionID, userID, lang).Scan(
		&conv.ID, &conv.SessionID, &conv.UserID, &conv.Title,
		&conv.LanguageCode, &conv.CreatedAt, &conv.UpdatedAt, &conv.ExpiresAt, &conv.IsArchived,
	)
	if err != nil {
		return nil, fmt.Errorf("create ai conversation: %w", err)
	}

	return conv, nil
}

func (r *PostgresAIRepo) GetConversation(ctx context.Context, sessionID string) (*Conversation, error) {
	conv := &Conversation{}
	err := r.db.QueryRow(ctx, `
		SELECT id, session_id, user_id, title, language_code, created_at, updated_at, expires_at, is_archived
		FROM ai_conversations
		WHERE session_id = $1
	`, sessionID).Scan(
		&conv.ID, &conv.SessionID, &conv.UserID, &conv.Title,
		&conv.LanguageCode, &conv.CreatedAt, &conv.UpdatedAt, &conv.ExpiresAt, &conv.IsArchived,
	)
	if err != nil {
		return nil, err
	}
	return conv, nil
}

func (r *PostgresAIRepo) GetRecentMessages(ctx context.Context, conversationID int64, limit int) ([]Message, error) {
	if limit <= 0 || limit > 50 {
		limit = 10
	}

	rows, err := r.db.Query(ctx, `
		SELECT id, conversation_id, role, content, intent, tool_calls, tool_results, metadata, token_count, latency_ms, created_at
		FROM ai_messages
		WHERE conversation_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`, conversationID, limit)
	if err != nil {
		return nil, fmt.Errorf("query ai messages: %w", err)
	}
	defer rows.Close()

	var msgs []Message
	for rows.Next() {
		var m Message
		var toolCallsJSON, toolResultsJSON, metaJSON []byte
		if err := rows.Scan(
			&m.ID, &m.ConversationID, &m.Role, &m.Content, &m.Intent,
			&toolCallsJSON, &toolResultsJSON, &metaJSON,
			&m.TokenCount, &m.LatencyMs, &m.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan ai message: %w", err)
		}
		if len(toolCallsJSON) > 0 {
			_ = json.Unmarshal(toolCallsJSON, &m.ToolCalls)
		}
		if len(toolResultsJSON) > 0 {
			_ = json.Unmarshal(toolResultsJSON, &m.ToolResults)
		}
		if len(metaJSON) > 0 {
			_ = json.Unmarshal(metaJSON, &m.Metadata)
		}
		msgs = append(msgs, m)
	}

	// Reverse to ascending order (chronological)
	for i, j := 0, len(msgs)-1; i < j; i, j = i+1, j-1 {
		msgs[i], msgs[j] = msgs[j], msgs[i]
	}

	return msgs, nil
}

func (r *PostgresAIRepo) SaveMessage(ctx context.Context, msg *Message) error {
	toolCallsJSON, _ := json.Marshal(msg.ToolCalls)
	toolResultsJSON, _ := json.Marshal(msg.ToolResults)
	metaJSON, _ := json.Marshal(msg.Metadata)

	err := r.db.QueryRow(ctx, `
		INSERT INTO ai_messages (
			conversation_id, role, content, intent, tool_calls, tool_results, metadata, token_count, latency_ms, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW())
		RETURNING id, created_at
	`, msg.ConversationID, msg.Role, msg.Content, msg.Intent, toolCallsJSON, toolResultsJSON, metaJSON, msg.TokenCount, msg.LatencyMs).Scan(&msg.ID, &msg.CreatedAt)

	if err != nil {
		return fmt.Errorf("insert ai message: %w", err)
	}

	// Update conversation updated_at
	_, _ = r.db.Exec(ctx, `UPDATE ai_conversations SET updated_at = NOW() WHERE id = $1`, msg.ConversationID)
	return nil
}

func (r *PostgresAIRepo) RecordAuditEvent(ctx context.Context, evt *AuditEvent) error {
	detailsJSON, _ := json.Marshal(evt.Details)
	_, err := r.db.Exec(ctx, `
		INSERT INTO ai_audit_events (session_id, user_id, event_type, details, ip_address, created_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
	`, evt.SessionID, evt.UserID, evt.EventType, detailsJSON, evt.IPAddress)
	return err
}

func (r *PostgresAIRepo) ArchiveConversation(ctx context.Context, sessionID string, userID *int) error {
	var err error
	if userID != nil {
		_, err = r.db.Exec(ctx, `
			UPDATE ai_conversations
			SET is_archived = TRUE, updated_at = NOW()
			WHERE session_id = $1 AND (user_id = $2 OR user_id IS NULL)
		`, sessionID, *userID)
	} else {
		_, err = r.db.Exec(ctx, `
			UPDATE ai_conversations
			SET is_archived = TRUE, updated_at = NOW()
			WHERE session_id = $1
		`, sessionID)
	}
	return err
}

func (r *PostgresAIRepo) GetAdminMetrics(ctx context.Context) (*AIMetricsReport, error) {
	report := &AIMetricsReport{}

	// Aggregated counts from ai_audit_events and ai_messages
	_ = r.db.QueryRow(ctx, `
		SELECT
			COALESCE(COUNT(*), 0) AS total_queries,
			COALESCE(COUNT(*) FILTER (WHERE event_type = 'OUT_OF_DOMAIN'), 0) AS out_of_domain,
			COALESCE(COUNT(*) FILTER (WHERE event_type = 'PROMPT_INJECTION_BLOCKED'), 0) AS adversarial_blocks,
			COALESCE(COUNT(*) FILTER (WHERE event_type = 'PROVIDER_ERROR'), 0) AS provider_errors
		FROM ai_audit_events
		WHERE created_at >= NOW() - INTERVAL '30 days'
	`).Scan(&report.TotalQueries, &report.OutOfDomainQueries, &report.AdversarialBlocks, &report.ProviderErrors)

	_ = r.db.QueryRow(ctx, `
		SELECT
			COALESCE(SUM(token_count), 0),
			COALESCE(AVG(latency_ms), 0)
		FROM ai_messages
		WHERE created_at >= NOW() - INTERVAL '30 days'
	`).Scan(&report.TotalTokens, &report.AverageLatencyMs)

	_ = r.db.QueryRow(ctx, `
		SELECT COUNT(DISTINCT session_id)
		FROM ai_conversations
		WHERE updated_at >= NOW() - INTERVAL '24 hours' AND is_archived = FALSE
	`).Scan(&report.ActiveSessions24h)

	// Estimated cost for Gemini 1.5 Flash (~$0.075 per 1M input tokens, ~$0.30 per 1M output tokens, blend ~$0.15 / 1M)
	report.EstimatedCostUSD = (float64(report.TotalTokens) / 1_000_000.0) * 0.15

	return report, nil
}
