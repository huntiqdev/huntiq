package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/your-org/used-tools-marketplace/internal/models"
)

type IdentityRepo struct {
	db *pgxpool.Pool
}

func NewIdentityRepo(db *pgxpool.Pool) *IdentityRepo {
	return &IdentityRepo{db: db}
}

// CreateIdentityVerification registers an official national identity verification record.
func (r *IdentityRepo) CreateIdentityVerification(ctx context.Context, v *models.IdentityVerification) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Check if credential hash is already registered to another user (anti-duplicate account rule)
	var existingUserID int
	err = tx.QueryRow(ctx, `SELECT user_id FROM identity_verifications WHERE credential_id_hash = $1`, v.CredentialIDHash).Scan(&existingUserID)
	if err == nil {
		if existingUserID != v.UserID {
			return fmt.Errorf("credential already registered to another individual account")
		}
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return err
	}

	metaBytes, _ := json.Marshal(v.VerificationMetadata)

	query := `
		INSERT INTO identity_verifications (
			user_id, provider_code, credential_type, credential_id_hash, id_masked,
			full_legal_name, phone_matched, phone_match_status, auth_reference,
			status, consent_recorded, consent_timestamp, verification_metadata,
			verified_at, expires_at, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, NOW(), NOW()
		)
		ON CONFLICT (credential_id_hash) DO UPDATE SET
			status = EXCLUDED.status,
			phone_matched = EXCLUDED.phone_matched,
			phone_match_status = EXCLUDED.phone_match_status,
			verified_at = EXCLUDED.verified_at,
			verification_metadata = EXCLUDED.verification_metadata,
			updated_at = NOW()
		RETURNING id
	`
	err = tx.QueryRow(ctx, query,
		v.UserID, v.ProviderCode, v.CredentialType, v.CredentialIDHash, v.IDMasked,
		v.FullLegalName, v.PhoneMatched, v.PhoneMatchStatus, v.AuthReference,
		v.Status, v.ConsentRecorded, v.ConsentTimestamp, metaBytes,
		v.VerifiedAt, v.ExpiresAt,
	).Scan(&v.ID)
	if err != nil {
		return err
	}

	// Update user's identity status if verified
	if v.Status == models.IdentityStatusVerified {
		_, err = tx.Exec(ctx, `UPDATE users SET identity_verified = true, name = $1, updated_at = NOW() WHERE id = $2`, v.FullLegalName, v.UserID)
		if err != nil {
			return fmt.Errorf("failed to mark user identity verified: %w", err)
		}
	}

	return tx.Commit(ctx)
}

// GetIdentityVerificationByUserID fetches the verification record for a user.
func (r *IdentityRepo) GetIdentityVerificationByUserID(ctx context.Context, userID int) (*models.IdentityVerification, error) {
	var v models.IdentityVerification
	var metaBytes []byte

	query := `
		SELECT id, user_id, provider_code, credential_type, id_masked, full_legal_name,
		       phone_matched, phone_match_status, auth_reference, status, consent_recorded,
		       consent_timestamp, verification_metadata, exception_reason, reviewed_by,
		       verified_at, expires_at, created_at, updated_at
		FROM identity_verifications
		WHERE user_id = $1
		ORDER BY id DESC
		LIMIT 1
	`
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&v.ID, &v.UserID, &v.ProviderCode, &v.CredentialType, &v.IDMasked, &v.FullLegalName,
		&v.PhoneMatched, &v.PhoneMatchStatus, &v.AuthReference, &v.Status, &v.ConsentRecorded,
		&v.ConsentTimestamp, &metaBytes, &v.ExceptionReason, &v.ReviewedBy,
		&v.VerifiedAt, &v.ExpiresAt, &v.CreatedAt, &v.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	if len(metaBytes) > 0 {
		_ = json.Unmarshal(metaBytes, &v.VerificationMetadata)
	}

	return &v, nil
}

// UpdateIdentityStatus allows KYC Officers to resolve identity exceptions with an audit trail.
func (r *IdentityRepo) UpdateIdentityStatus(ctx context.Context, id int64, status string, reviewerID int, reason string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var userID int
	err = tx.QueryRow(ctx, `
		UPDATE identity_verifications
		SET status = $1, reviewed_by = $2, exception_reason = $3,
		    verified_at = CASE WHEN $1 = 'verified' THEN NOW() ELSE verified_at END,
		    updated_at = NOW()
		WHERE id = $4
		RETURNING user_id
	`, status, reviewerID, reason, id).Scan(&userID)
	if err != nil {
		return err
	}
	if userID == reviewerID {
		return errors.New("maker-checker separation: reviewer cannot approve or resolve their own identity verification")
	}

	isVerified := status == models.IdentityStatusVerified
	_, err = tx.Exec(ctx, `UPDATE users SET identity_verified = $1, updated_at = NOW() WHERE id = $2`, isVerified, userID)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// ListPendingIdentityVerifications returns paginated pending identity records for review.
func (r *IdentityRepo) ListPendingIdentityVerifications(ctx context.Context, limit, offset int) ([]models.IdentityVerification, int, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	var total int
	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM identity_verifications WHERE status = 'pending' OR status = 'escalated'`).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.db.Query(ctx, `
		SELECT id, user_id, provider_code, credential_type, id_masked, full_legal_name,
		       phone_matched, phone_match_status, auth_reference, status, consent_recorded,
		       consent_timestamp, exception_reason, reviewed_by, verified_at, expires_at, created_at, updated_at
		FROM identity_verifications
		WHERE status = 'pending' OR status = 'escalated'
		ORDER BY created_at ASC
		LIMIT $1 OFFSET $2
	`, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []models.IdentityVerification
	for rows.Next() {
		var v models.IdentityVerification
		if err := rows.Scan(
			&v.ID, &v.UserID, &v.ProviderCode, &v.CredentialType, &v.IDMasked, &v.FullLegalName,
			&v.PhoneMatched, &v.PhoneMatchStatus, &v.AuthReference, &v.Status, &v.ConsentRecorded,
			&v.ConsentTimestamp, &v.ExceptionReason, &v.ReviewedBy, &v.VerifiedAt, &v.ExpiresAt,
			&v.CreatedAt, &v.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		list = append(list, v)
	}

	return list, total, nil
}
