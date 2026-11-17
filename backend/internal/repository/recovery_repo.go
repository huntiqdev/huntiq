package repository

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/your-org/used-tools-marketplace/internal/models"
)

type RecoveryRepo struct {
	db *pgxpool.Pool
}

func NewRecoveryRepo(db *pgxpool.Pool) *RecoveryRepo {
	return &RecoveryRepo{db: db}
}

// CreateRecoveryRequest checks the mandatory 3-hour rolling cooldown and creates an OTP recovery session.
func (r *RecoveryRepo) CreateRecoveryRequest(ctx context.Context, userID int, phone, otpHash string) (*models.PasswordRecoveryRequest, error) {
	// 1. Enforce rolling 3-hour cooldown
	var cooldownUntil time.Time
	err := r.db.QueryRow(ctx, `
		SELECT cooldown_until
		FROM password_recovery_requests
		WHERE phone = $1 AND cooldown_until > NOW()
		ORDER BY cooldown_until DESC
		LIMIT 1
	`, phone).Scan(&cooldownUntil)
	if err == nil {
		remaining := time.Until(cooldownUntil).Round(time.Minute)
		return nil, fmt.Errorf("recovery_cooldown_active: password recovery OTP may only be requested once per 3 hours (cooldown remaining: %v)", remaining)
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	// 2. Insert new recovery record with 3-hour cooldown
	req := &models.PasswordRecoveryRequest{
		UserID:        userID,
		Phone:         phone,
		OTPHash:       otpHash,
		Attempts:      0,
		MaxAttempts:   5,
		Status:        models.RecoveryStatusPending,
		CooldownUntil: time.Now().Add(3 * time.Hour),
	}

	query := `
		INSERT INTO password_recovery_requests (
			user_id, phone, otp_hash, attempts, max_attempts, status, cooldown_until, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, NOW(), NOW()
		) RETURNING id, created_at, updated_at
	`
	err = r.db.QueryRow(ctx, query,
		req.UserID, req.Phone, req.OTPHash, req.Attempts, req.MaxAttempts,
		req.Status, req.CooldownUntil,
	).Scan(&req.ID, &req.CreatedAt, &req.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// ValidateRecoveryOTP validates OTP attempts, enforces 5-minute validity, and issues a single-use reset authorization token.
func (r *RecoveryRepo) ValidateRecoveryOTP(ctx context.Context, phone, otpHash string) (string, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return "", err
	}
	defer tx.Rollback(ctx)

	var req models.PasswordRecoveryRequest
	query := `
		SELECT id, user_id, phone, otp_hash, attempts, max_attempts, status, created_at
		FROM password_recovery_requests
		WHERE phone = $1 AND status = 'pending' AND created_at > NOW() - INTERVAL '5 minutes'
		ORDER BY id DESC
		LIMIT 1
		FOR UPDATE
	`
	err = tx.QueryRow(ctx, query, phone).Scan(
		&req.ID, &req.UserID, &req.Phone, &req.OTPHash, &req.Attempts, &req.MaxAttempts,
		&req.Status, &req.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", errors.New("no active password recovery request found or OTP has expired (5-minute window)")
		}
		return "", err
	}

	if req.Attempts >= req.MaxAttempts {
		_, _ = tx.Exec(ctx, `UPDATE password_recovery_requests SET status = 'expired', updated_at = NOW() WHERE id = $1`, req.ID)
		_ = tx.Commit(ctx)
		return "", errors.New("maximum verification attempts exceeded — please wait for the 3-hour cooldown window")
	}

	if req.OTPHash != otpHash {
		_, _ = tx.Exec(ctx, `UPDATE password_recovery_requests SET attempts = attempts + 1, updated_at = NOW() WHERE id = $1`, req.ID)
		_ = tx.Commit(ctx)
		return "", errors.New("invalid verification code")
	}

	// Generate single-use reset token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", err
	}
	rawToken := hex.EncodeToString(tokenBytes)
	tokenHash := sha256.Sum256([]byte(rawToken))
	tokenHashHex := hex.EncodeToString(tokenHash[:])
	tokenExpiresAt := time.Now().Add(15 * time.Minute)

	_, err = tx.Exec(ctx, `
		UPDATE password_recovery_requests
		SET status = 'verified', reset_token_hash = $1, reset_token_expires_at = $2, updated_at = NOW()
		WHERE id = $3
	`, tokenHashHex, tokenExpiresAt, req.ID)
	if err != nil {
		return "", err
	}

	if err := tx.Commit(ctx); err != nil {
		return "", err
	}

	return rawToken, nil
}

// ConsumeResetToken checks authorization validity, atomically burns the token, and returns the target user ID.
func (r *RecoveryRepo) ConsumeResetToken(ctx context.Context, rawToken string) (int, string, error) {
	tokenHash := sha256.Sum256([]byte(rawToken))
	tokenHashHex := hex.EncodeToString(tokenHash[:])

	var userID int
	var phone string

	query := `
		UPDATE password_recovery_requests
		SET status = 'consumed', updated_at = NOW()
		WHERE reset_token_hash = $1 AND status = 'verified' AND reset_token_expires_at > NOW()
		RETURNING user_id, phone
	`
	err := r.db.QueryRow(ctx, query, tokenHashHex).Scan(&userID, &phone)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, "", errors.New("invalid or expired password-reset authorization token")
		}
		return 0, "", err
	}

	return userID, phone, nil
}

