package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/your-org/used-tools-marketplace/internal/models"
)

var (
	ErrSessionNotFound = errors.New("refresh token not found")
	ErrSessionRevoked  = errors.New("refresh token has been revoked")
	ErrSessionExpired  = errors.New("refresh token has expired")
)

type SessionRepo struct {
	db *pgxpool.Pool
}

func NewSessionRepo(db *pgxpool.Pool) *SessionRepo {
	return &SessionRepo{db: db}
}

// CreateRefreshToken stores a newly issued refresh token in the database
func (r *SessionRepo) CreateRefreshToken(ctx context.Context, token *models.UserRefreshToken) error {
	query := `
		INSERT INTO user_refresh_tokens (
			id, user_id, family_id, token_hash, parent_token_hash,
			is_revoked, expires_at, user_agent, ip_address, created_at, last_used_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	if token.ID == uuid.Nil {
		token.ID = uuid.New()
	}
	now := time.Now()
	if token.CreatedAt.IsZero() {
		token.CreatedAt = now
	}
	if token.LastUsedAt.IsZero() {
		token.LastUsedAt = now
	}

	_, err := r.db.Exec(ctx, query,
		token.ID,
		token.UserID,
		token.FamilyID,
		token.TokenHash,
		token.ParentTokenHash,
		token.IsRevoked,
		token.ExpiresAt,
		token.UserAgent,
		token.IPAddress,
		token.CreatedAt,
		token.LastUsedAt,
	)
	return err
}

// GetRefreshTokenByHash retrieves a refresh token record by its SHA-256 hash
func (r *SessionRepo) GetRefreshTokenByHash(ctx context.Context, tokenHash string) (*models.UserRefreshToken, error) {
	query := `
		SELECT id, user_id, family_id, token_hash, COALESCE(parent_token_hash, ''),
		       is_revoked, COALESCE(revocation_reason, ''), revoked_at,
		       expires_at, COALESCE(user_agent, ''), COALESCE(ip_address, ''),
		       created_at, last_used_at
		FROM user_refresh_tokens
		WHERE token_hash = $1
	`
	row := r.db.QueryRow(ctx, query, tokenHash)

	var t models.UserRefreshToken
	err := row.Scan(
		&t.ID,
		&t.UserID,
		&t.FamilyID,
		&t.TokenHash,
		&t.ParentTokenHash,
		&t.IsRevoked,
		&t.RevocationReason,
		&t.RevokedAt,
		&t.ExpiresAt,
		&t.UserAgent,
		&t.IPAddress,
		&t.CreatedAt,
		&t.LastUsedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrSessionNotFound
		}
		return nil, err
	}
	return &t, nil
}

// RotateRefreshToken atomically consumes the old refresh token and inserts the new one
func (r *SessionRepo) RotateRefreshToken(ctx context.Context, oldTokenID uuid.UUID, newToken *models.UserRefreshToken) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// 1. Mark old token as revoked/rotated
	updateQuery := `
		UPDATE user_refresh_tokens
		SET is_revoked = TRUE,
		    revocation_reason = 'rotated',
		    revoked_at = NOW(),
		    last_used_at = NOW()
		WHERE id = $1 AND is_revoked = FALSE
	`
	tag, err := tx.Exec(ctx, updateQuery, oldTokenID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrSessionRevoked
	}

	// 2. Insert new token in the same transaction
	insertQuery := `
		INSERT INTO user_refresh_tokens (
			id, user_id, family_id, token_hash, parent_token_hash,
			is_revoked, expires_at, user_agent, ip_address, created_at, last_used_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	if newToken.ID == uuid.Nil {
		newToken.ID = uuid.New()
	}
	now := time.Now()
	if newToken.CreatedAt.IsZero() {
		newToken.CreatedAt = now
	}
	if newToken.LastUsedAt.IsZero() {
		newToken.LastUsedAt = now
	}

	_, err = tx.Exec(ctx, insertQuery,
		newToken.ID,
		newToken.UserID,
		newToken.FamilyID,
		newToken.TokenHash,
		newToken.ParentTokenHash,
		newToken.IsRevoked,
		newToken.ExpiresAt,
		newToken.UserAgent,
		newToken.IPAddress,
		newToken.CreatedAt,
		newToken.LastUsedAt,
	)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// RevokeFamily invalidates all refresh tokens belonging to a compromised lineage
func (r *SessionRepo) RevokeFamily(ctx context.Context, familyID uuid.UUID, reason string) error {
	query := `
		UPDATE user_refresh_tokens
		SET is_revoked = TRUE,
		    revocation_reason = $1,
		    revoked_at = NOW()
		WHERE family_id = $2 AND is_revoked = FALSE
	`
	_, err := r.db.Exec(ctx, query, reason, familyID)
	return err
}

// RevokeByUserID invalidates all active sessions for a user (e.g. on logout-all, password change, account suspension)
func (r *SessionRepo) RevokeByUserID(ctx context.Context, userID int, reason string) error {
	query := `
		UPDATE user_refresh_tokens
		SET is_revoked = TRUE,
		    revocation_reason = $1,
		    revoked_at = NOW()
		WHERE user_id = $2 AND is_revoked = FALSE
	`
	_, err := r.db.Exec(ctx, query, reason, userID)
	return err
}

// RevokeByTokenHash invalidates a specific refresh token (e.g. single session logout)
func (r *SessionRepo) RevokeByTokenHash(ctx context.Context, tokenHash string, reason string) error {
	query := `
		UPDATE user_refresh_tokens
		SET is_revoked = TRUE,
		    revocation_reason = $1,
		    revoked_at = NOW()
		WHERE token_hash = $2 AND is_revoked = FALSE
	`
	_, err := r.db.Exec(ctx, query, reason, tokenHash)
	return err
}
