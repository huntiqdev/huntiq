package repository

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/your-org/used-tools-marketplace/internal/models"
)

type ReferralTokenRepo struct {
	db *pgxpool.Pool
}

func NewReferralTokenRepo(db *pgxpool.Pool) *ReferralTokenRepo {
	return &ReferralTokenRepo{db: db}
}

// generateSecureToken creates a 64-character hexadecimal random token.
func generateSecureToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// GetOrCreateReferralToken ensures a high-entropy referral token exists for the given user.
func (r *ReferralTokenRepo) GetOrCreateReferralToken(ctx context.Context, userID int, referralCode string) (*models.ReferralToken, error) {
	var t models.ReferralToken
	query := `
		SELECT id, user_id, token, referral_code, click_count, conversion_count, is_active, created_at, updated_at
		FROM referral_tokens
		WHERE user_id = $1
		LIMIT 1
	`
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&t.ID,
		&t.UserID,
		&t.Token,
		&t.ReferralCode,
		&t.ClickCount,
		&t.ConversionCount,
		&t.IsActive,
		&t.CreatedAt,
		&t.UpdatedAt,
	)
	if err == nil {
		return &t, nil
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	// Generate new token
	tokenStr, err := generateSecureToken()
	if err != nil {
		return nil, err
	}

	insertQuery := `
		INSERT INTO referral_tokens (user_id, token, referral_code, click_count, conversion_count, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, 0, 0, true, NOW(), NOW())
		RETURNING id, user_id, token, referral_code, click_count, conversion_count, is_active, created_at, updated_at
	`
	err = r.db.QueryRow(ctx, insertQuery, userID, tokenStr, referralCode).Scan(
		&t.ID,
		&t.UserID,
		&t.Token,
		&t.ReferralCode,
		&t.ClickCount,
		&t.ConversionCount,
		&t.IsActive,
		&t.CreatedAt,
		&t.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// GetByToken fetches a referral token record by its 64-character token string.
func (r *ReferralTokenRepo) GetByToken(ctx context.Context, token string) (*models.ReferralToken, error) {
	var t models.ReferralToken
	query := `
		SELECT id, user_id, token, referral_code, click_count, conversion_count, is_active, created_at, updated_at
		FROM referral_tokens
		WHERE token = $1
	`
	err := r.db.QueryRow(ctx, query, token).Scan(
		&t.ID,
		&t.UserID,
		&t.Token,
		&t.ReferralCode,
		&t.ClickCount,
		&t.ConversionCount,
		&t.IsActive,
		&t.CreatedAt,
		&t.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// GetByReferralCode fetches a token record by user referral code.
func (r *ReferralTokenRepo) GetByReferralCode(ctx context.Context, code string) (*models.ReferralToken, error) {
	var t models.ReferralToken
	query := `
		SELECT id, user_id, token, referral_code, click_count, conversion_count, is_active, created_at, updated_at
		FROM referral_tokens
		WHERE referral_code = $1
		ORDER BY id DESC
		LIMIT 1
	`
	err := r.db.QueryRow(ctx, query, code).Scan(
		&t.ID,
		&t.UserID,
		&t.Token,
		&t.ReferralCode,
		&t.ClickCount,
		&t.ConversionCount,
		&t.IsActive,
		&t.CreatedAt,
		&t.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// RecordAttribution inserts an attribution event.
func (r *ReferralTokenRepo) RecordAttribution(ctx context.Context, attr *models.ReferralAttribution) (int64, error) {
	var id int64
	query := `
		INSERT INTO referral_attributions (
			token_id, inviter_id, ip_hash, user_agent_hash, landing_url, status, attributed_at, expires_at
		) VALUES ($1, $2, $3, $4, $5, 'clicked', NOW(), NOW() + INTERVAL '30 days')
		RETURNING id
	`
	err := r.db.QueryRow(ctx, query, attr.TokenID, attr.InviterID, attr.IPHash, attr.UserAgentHash, attr.LandingURL).Scan(&id)
	if err != nil {
		return 0, err
	}
	attr.ID = id
	return id, nil
}

// IncrementClickCount increments the click count of a token.
func (r *ReferralTokenRepo) IncrementClickCount(ctx context.Context, tokenID int) error {
	_, err := r.db.Exec(ctx, `
		UPDATE referral_tokens
		SET click_count = click_count + 1, updated_at = NOW()
		WHERE id = $1
	`, tokenID)
	return err
}

// IncrementConversionCount increments the conversion counter.
func (r *ReferralTokenRepo) IncrementConversionCount(ctx context.Context, tokenID int) error {
	_, err := r.db.Exec(ctx, `
		UPDATE referral_tokens
		SET conversion_count = conversion_count + 1, updated_at = NOW()
		WHERE id = $1
	`, tokenID)
	return err
}

// FindLatestActiveAttribution retrieves the most recent active attribution for an IP/token within the 30-day window.
func (r *ReferralTokenRepo) FindLatestActiveAttribution(ctx context.Context, tokenID int, ipHash string) (*models.ReferralAttribution, error) {
	var a models.ReferralAttribution
	query := `
		SELECT id, token_id, inviter_id, ip_hash, user_agent_hash, landing_url, status, registered_user_id, attributed_at, expires_at
		FROM referral_attributions
		WHERE token_id = $1 AND ip_hash = $2 AND status = 'clicked' AND expires_at > NOW()
		ORDER BY attributed_at DESC
		LIMIT 1
	`
	err := r.db.QueryRow(ctx, query, tokenID, ipHash).Scan(
		&a.ID,
		&a.TokenID,
		&a.InviterID,
		&a.IPHash,
		&a.UserAgentHash,
		&a.LandingURL,
		&a.Status,
		&a.RegisteredUserID,
		&a.AttributedAt,
		&a.ExpiresAt,
	)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

// LockAttributionOnRegistration finalizes attribution upon phone verification.
func (r *ReferralTokenRepo) LockAttributionOnRegistration(ctx context.Context, attributionID int64, newUserID int) error {
	_, err := r.db.Exec(ctx, `
		UPDATE referral_attributions
		SET status = 'registered', registered_user_id = $1
		WHERE id = $2
	`, newUserID, attributionID)
	return err
}

// SetUserInvitedBy permanently associates the invited user with the inviter.
func (r *ReferralTokenRepo) SetUserInvitedBy(ctx context.Context, userID int, inviterID int) error {
	_, err := r.db.Exec(ctx, `
		UPDATE users
		SET invited_by = $1, invitation_accepted_at = NOW(), updated_at = NOW()
		WHERE id = $2 AND invited_by IS NULL
	`, inviterID, userID)
	return err
}

// GetUserInviter returns the inviter ID for a user if present.
func (r *ReferralTokenRepo) GetUserInviter(ctx context.Context, userID int) (int, error) {
	var inviterID *int
	err := r.db.QueryRow(ctx, `SELECT invited_by FROM users WHERE id = $1`, userID).Scan(&inviterID)
	if err != nil {
		return 0, err
	}
	if inviterID == nil {
		return 0, nil
	}
	return *inviterID, nil
}
