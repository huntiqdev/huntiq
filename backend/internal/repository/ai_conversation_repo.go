package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/your-org/used-tools-marketplace/internal/models"
)

// AIConversationRepo manages secure multi-turn conversation persistence in PostgreSQL.
type AIConversationRepo struct {
	db *pgxpool.Pool
}

// NewAIConversationRepo creates a new repository instance.
func NewAIConversationRepo(db *pgxpool.Pool) *AIConversationRepo {
	return &AIConversationRepo{db: db}
}

// CreateConversation inserts a new conversation owned by the authenticated user.
func (r *AIConversationRepo) CreateConversation(ctx context.Context, userID int, title string) (*models.AIConversation, error) {
	if title == "" {
		title = "New Chat"
	}
	query := `
		INSERT INTO ai_conversations (user_id, title, created_at, updated_at)
		VALUES ($1, $2, NOW(), NOW())
		RETURNING id, user_id, title, created_at, updated_at`

	var conv models.AIConversation
	err := r.db.QueryRow(ctx, query, userID, title).Scan(
		&conv.ID,
		&conv.UserID,
		&conv.Title,
		&conv.CreatedAt,
		&conv.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create ai conversation: %w", err)
	}
	return &conv, nil
}

// GetConversationByID retrieves a conversation verifying authenticated user ownership.
func (r *AIConversationRepo) GetConversationByID(ctx context.Context, conversationID string, userID int) (*models.AIConversation, error) {
	query := `
		SELECT id, user_id, title, created_at, updated_at
		FROM ai_conversations
		WHERE id = $1 AND user_id = $2`

	var conv models.AIConversation
	err := r.db.QueryRow(ctx, query, conversationID, userID).Scan(
		&conv.ID,
		&conv.UserID,
		&conv.Title,
		&conv.CreatedAt,
		&conv.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("conversation not found or unauthorized: %w", err)
	}

	messages, err := r.GetMessagesByConversationID(ctx, conversationID)
	if err != nil {
		return nil, err
	}
	conv.Messages = messages

	return &conv, nil
}

// ListConversationsByUserID retrieves conversations for the user ordered by recent activity.
func (r *AIConversationRepo) ListConversationsByUserID(ctx context.Context, userID int, limit, offset int) ([]models.AIConversation, error) {
	if limit <= 0 || limit > 50 {
		limit = 30
	}
	if offset < 0 {
		offset = 0
	}

	query := `
		SELECT id, user_id, title, created_at, updated_at
		FROM ai_conversations
		WHERE user_id = $1
		ORDER BY updated_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list conversations: %w", err)
	}
	defer rows.Close()

	var conversations []models.AIConversation
	for rows.Next() {
		var conv models.AIConversation
		if err := rows.Scan(&conv.ID, &conv.UserID, &conv.Title, &conv.CreatedAt, &conv.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan conversation row: %w", err)
		}
		conversations = append(conversations, conv)
	}
	return conversations, nil
}

// UpdateConversationTitle updates a conversation's title with user ownership enforcement.
func (r *AIConversationRepo) UpdateConversationTitle(ctx context.Context, conversationID string, userID int, title string) error {
	query := `
		UPDATE ai_conversations
		SET title = $1, updated_at = NOW()
		WHERE id = $2 AND user_id = $3`

	cmd, err := r.db.Exec(ctx, query, title, conversationID, userID)
	if err != nil {
		return fmt.Errorf("failed to update conversation title: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("conversation not found or unauthorized")
	}
	return nil
}

// DeleteConversation deletes a conversation with cascade deleting messages.
func (r *AIConversationRepo) DeleteConversation(ctx context.Context, conversationID string, userID int) error {
	query := `
		DELETE FROM ai_conversations
		WHERE id = $1 AND user_id = $2`

	cmd, err := r.db.Exec(ctx, query, conversationID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete conversation: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("conversation not found or unauthorized")
	}
	return nil
}

// CreateMessage inserts a new dialogue turn and touches the parent conversation's updated_at.
func (r *AIConversationRepo) CreateMessage(ctx context.Context, msg *models.AIMessage) error {
	var toolCallsJSON []byte
	if len(msg.ToolCalls) > 0 {
		var err error
		toolCallsJSON, err = json.Marshal(msg.ToolCalls)
		if err != nil {
			return fmt.Errorf("failed to marshal tool calls: %w", err)
		}
	}

	var metadataJSON []byte
	if len(msg.Metadata) > 0 {
		var err error
		metadataJSON, err = json.Marshal(msg.Metadata)
		if err != nil {
			return fmt.Errorf("failed to marshal metadata: %w", err)
		}
	}

	query := `
		INSERT INTO ai_messages (conversation_id, role, content, tool_calls, tool_call_id, metadata, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW())
		RETURNING id, created_at`

	err := r.db.QueryRow(
		ctx,
		query,
		msg.ConversationID,
		msg.Role,
		msg.Content,
		toolCallsJSON,
		msg.ToolCallID,
		metadataJSON,
	).Scan(&msg.ID, &msg.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to insert ai message: %w", err)
	}

	// Update conversation's updated_at timestamp
	_, _ = r.db.Exec(ctx, `UPDATE ai_conversations SET updated_at = NOW() WHERE id = $1`, msg.ConversationID)
	return nil
}

// GetMessagesByConversationID retrieves all messages in chronological sequence.
func (r *AIConversationRepo) GetMessagesByConversationID(ctx context.Context, conversationID string) ([]models.AIMessage, error) {
	query := `
		SELECT id, conversation_id, role, content, tool_calls, tool_call_id, metadata, created_at
		FROM ai_messages
		WHERE conversation_id = $1
		ORDER BY created_at ASC`

	rows, err := r.db.Query(ctx, query, conversationID)
	if err != nil {
		return nil, fmt.Errorf("failed to query messages: %w", err)
	}
	defer rows.Close()

	var messages []models.AIMessage
	for rows.Next() {
		var msg models.AIMessage
		var toolCallsRaw []byte
		var metadataRaw []byte

		err := rows.Scan(
			&msg.ID,
			&msg.ConversationID,
			&msg.Role,
			&msg.Content,
			&toolCallsRaw,
			&msg.ToolCallID,
			&metadataRaw,
			&msg.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan message row: %w", err)
		}

		if len(toolCallsRaw) > 0 {
			_ = json.Unmarshal(toolCallsRaw, &msg.ToolCalls)
		}
		if len(metadataRaw) > 0 {
			_ = json.Unmarshal(metadataRaw, &msg.Metadata)
		}
		messages = append(messages, msg)
	}
	return messages, nil
}
