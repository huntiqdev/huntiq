package repository

import (
	"context"
	"encoding/base64"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/your-org/used-tools-marketplace/internal/models"
)

// AdminPasskeyRepo provides database persistence for passkeys, sessions, challenges, and abuse monitoring
type AdminPasskeyRepo struct {
	db *pgxpool.Pool
}

// NewAdminPasskeyRepo creates an instance of AdminPasskeyRepo
func NewAdminPasskeyRepo(db *pgxpool.Pool) *AdminPasskeyRepo {
	return &AdminPasskeyRepo{db: db}
}

// SavePasskey persists a new WebAuthn passkey credential
func (r *AdminPasskeyRepo) SavePasskey(ctx context.Context, cred *models.PasskeyCredential) error {
	query := `
		INSERT INTO admin_passkeys (
			user_id, credential_id, public_key, attestation_type, aaguid,
			sign_count, attachment, name, transports, backup_eligible, backup_state, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, NOW())
		RETURNING id, created_at
	`
	return r.db.QueryRow(ctx, query,
		cred.UserID,
		cred.CredentialID,
		cred.PublicKey,
		cred.AttestationType,
		cred.AAGUID,
		cred.SignCount,
		cred.Attachment,
		cred.Name,
		cred.Transports,
		cred.BackupEligible,
		cred.BackupState,
	).Scan(&cred.ID, &cred.CreatedAt)
}

// ListPasskeysByUserID retrieves all active passkeys for an administrator
func (r *AdminPasskeyRepo) ListPasskeysByUserID(ctx context.Context, userID int) ([]models.PasskeyCredential, error) {
	query := `
		SELECT id, user_id, credential_id, public_key, attestation_type, aaguid,
		       sign_count, clone_warning, attachment, name, transports, backup_eligible,
		       backup_state, created_at, last_used_at, revoked_at
		FROM admin_passkeys
		WHERE user_id = $1 AND revoked_at IS NULL
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.PasskeyCredential
	for rows.Next() {
		var c models.PasskeyCredential
		var transports []string
		if err := rows.Scan(
			&c.ID, &c.UserID, &c.CredentialID, &c.PublicKey, &c.AttestationType, &c.AAGUID,
			&c.SignCount, &c.CloneWarning, &c.Attachment, &c.Name, &transports,
			&c.BackupEligible, &c.BackupState, &c.CreatedAt, &c.LastUsedAt, &c.RevokedAt,
		); err != nil {
			return nil, err
		}
		c.Transports = transports
		c.CredentialIDB64 = base64.RawURLEncoding.EncodeToString(c.CredentialID)
		list = append(list, c)
	}
	return list, rows.Err()
}

// GetPasskeyByCredentialID looks up a credential by its raw binary ID
func (r *AdminPasskeyRepo) GetPasskeyByCredentialID(ctx context.Context, credentialID []byte) (*models.PasskeyCredential, error) {
	query := `
		SELECT id, user_id, credential_id, public_key, attestation_type, aaguid,
		       sign_count, clone_warning, attachment, name, transports, backup_eligible,
		       backup_state, created_at, last_used_at, revoked_at
		FROM admin_passkeys
		WHERE credential_id = $1 AND revoked_at IS NULL
	`
	var c models.PasskeyCredential
	var transports []string
	err := r.db.QueryRow(ctx, query, credentialID).Scan(
		&c.ID, &c.UserID, &c.CredentialID, &c.PublicKey, &c.AttestationType, &c.AAGUID,
		&c.SignCount, &c.CloneWarning, &c.Attachment, &c.Name, &transports,
		&c.BackupEligible, &c.BackupState, &c.CreatedAt, &c.LastUsedAt, &c.RevokedAt,
	)
	if err != nil {
		return nil, err
	}
	c.Transports = transports
	c.CredentialIDB64 = base64.RawURLEncoding.EncodeToString(c.CredentialID)
	return &c, nil
}

// UpdatePasskeySignCount updates sign count and marks last_used_at
func (r *AdminPasskeyRepo) UpdatePasskeySignCount(ctx context.Context, credentialID []byte, signCount uint32) error {
	query := `
		UPDATE admin_passkeys
		SET sign_count = $1, last_used_at = NOW()
		WHERE credential_id = $2
	`
	_, err := r.db.Exec(ctx, query, signCount, credentialID)
	return err
}

// FlagPasskeyCloneWarning sets the clone warning flag when counter rollback is detected
func (r *AdminPasskeyRepo) FlagPasskeyCloneWarning(ctx context.Context, credentialID []byte) error {
	query := `
		UPDATE admin_passkeys
		SET clone_warning = TRUE
		WHERE credential_id = $1
	`
	_, err := r.db.Exec(ctx, query, credentialID)
	return err
}

// RevokePasskey marks a passkey credential as revoked
func (r *AdminPasskeyRepo) RevokePasskey(ctx context.Context, passkeyID, userID int, reason string) error {
	query := `
		UPDATE admin_passkeys
		SET revoked_at = NOW(), revocation_reason = $1
		WHERE id = $2 AND user_id = $3 AND revoked_at IS NULL
	`
	res, err := r.db.Exec(ctx, query, reason, passkeyID, userID)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return errors.New("passkey not found or already revoked")
	}
	return nil
}

// SaveChallenge stores a single-use authentication/registration challenge
func (r *AdminPasskeyRepo) SaveChallenge(ctx context.Context, c *models.AuthChallenge) error {
	query := `
		INSERT INTO admin_auth_challenges (
			challenge, user_id, portal, ceremony_type, session_data, expires_at, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, NOW())
		RETURNING id, created_at
	`
	return r.db.QueryRow(ctx, query,
		c.Challenge,
		c.UserID,
		c.Portal,
		c.CeremonyType,
		c.SessionData,
		c.ExpiresAt,
	).Scan(&c.ID, &c.CreatedAt)
}

// GetAndConsumeChallenge retrieves and marks a challenge as consumed atomically
func (r *AdminPasskeyRepo) GetAndConsumeChallenge(ctx context.Context, challengeStr string) (*models.AuthChallenge, error) {
	query := `
		UPDATE admin_auth_challenges
		SET used_at = NOW()
		WHERE challenge = $1 AND used_at IS NULL AND expires_at > NOW()
		RETURNING id, challenge, user_id, portal, ceremony_type, session_data, expires_at, used_at, created_at
	`
	var c models.AuthChallenge
	err := r.db.QueryRow(ctx, query, challengeStr).Scan(
		&c.ID,
		&c.Challenge,
		&c.UserID,
		&c.Portal,
		&c.CeremonyType,
		&c.SessionData,
		&c.ExpiresAt,
		&c.UsedAt,
		&c.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("challenge expired, already used, or invalid")
		}
		return nil, err
	}
	return &c, nil
}

// CreateSession records a new active administrative session
func (r *AdminPasskeyRepo) CreateSession(ctx context.Context, s *models.AdminSession) error {
	query := `
		INSERT INTO admin_sessions (
			session_token_hash, user_id, portal, role, ip_address, user_agent,
			last_active_at, idle_timeout_seconds, absolute_expires_at, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, NOW(), $7, $8, NOW())
		RETURNING id, created_at, last_active_at
	`
	return r.db.QueryRow(ctx, query,
		s.SessionTokenHash,
		s.UserID,
		s.Portal,
		s.Role,
		s.IPAddress,
		s.UserAgent,
		s.IdleTimeoutSeconds,
		s.AbsoluteExpiresAt,
	).Scan(&s.ID, &s.CreatedAt, &s.LastActiveAt)
}

// GetSessionByTokenHash retrieves an unrevoked, unexpired session
func (r *AdminPasskeyRepo) GetSessionByTokenHash(ctx context.Context, tokenHash string) (*models.AdminSession, error) {
	query := `
		SELECT id, session_token_hash, user_id, portal, role, ip_address, user_agent,
		       last_active_at, idle_timeout_seconds, absolute_expires_at, revoked_at,
		       COALESCE(revocation_reason, ''), created_at
		FROM admin_sessions
		WHERE session_token_hash = $1 AND revoked_at IS NULL
	`
	var s models.AdminSession
	err := r.db.QueryRow(ctx, query, tokenHash).Scan(
		&s.ID,
		&s.SessionTokenHash,
		&s.UserID,
		&s.Portal,
		&s.Role,
		&s.IPAddress,
		&s.UserAgent,
		&s.LastActiveAt,
		&s.IdleTimeoutSeconds,
		&s.AbsoluteExpiresAt,
		&s.RevokedAt,
		&s.RevocationReason,
		&s.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	// Check absolute expiry
	if now.After(s.AbsoluteExpiresAt) {
		_ = r.RevokeSession(ctx, s.ID, "absolute_expiry_exceeded")
		return nil, errors.New("session expired")
	}

	// Check idle timeout
	if now.Sub(s.LastActiveAt) > time.Duration(s.IdleTimeoutSeconds)*time.Second {
		_ = r.RevokeSession(ctx, s.ID, "idle_timeout_exceeded")
		return nil, errors.New("session timed out due to inactivity")
	}

	return &s, nil
}

// GetSessionByID retrieves an unrevoked, unexpired session by primary key ID
func (r *AdminPasskeyRepo) GetSessionByID(ctx context.Context, sessionID int) (*models.AdminSession, error) {
	query := `
		SELECT id, session_token_hash, user_id, portal, role, ip_address, user_agent,
		       last_active_at, idle_timeout_seconds, absolute_expires_at, revoked_at,
		       COALESCE(revocation_reason, ''), created_at
		FROM admin_sessions
		WHERE id = $1 AND revoked_at IS NULL
	`
	var s models.AdminSession
	err := r.db.QueryRow(ctx, query, sessionID).Scan(
		&s.ID,
		&s.SessionTokenHash,
		&s.UserID,
		&s.Portal,
		&s.Role,
		&s.IPAddress,
		&s.UserAgent,
		&s.LastActiveAt,
		&s.IdleTimeoutSeconds,
		&s.AbsoluteExpiresAt,
		&s.RevokedAt,
		&s.RevocationReason,
		&s.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	// Check absolute expiry
	if now.After(s.AbsoluteExpiresAt) {
		_ = r.RevokeSession(ctx, s.ID, "absolute_expiry_exceeded")
		return nil, errors.New("session expired")
	}

	// Check idle timeout
	if now.Sub(s.LastActiveAt) > time.Duration(s.IdleTimeoutSeconds)*time.Second {
		_ = r.RevokeSession(ctx, s.ID, "idle_timeout_exceeded")
		return nil, errors.New("session timed out due to inactivity")
	}

	return &s, nil
}

// UpdateSessionActivity refreshes the last_active_at timestamp of a session
func (r *AdminPasskeyRepo) UpdateSessionActivity(ctx context.Context, sessionID int) error {
	query := `UPDATE admin_sessions SET last_active_at = NOW() WHERE id = $1`
	_, err := r.db.Exec(ctx, query, sessionID)
	return err
}

// RevokeSession marks a session as revoked
func (r *AdminPasskeyRepo) RevokeSession(ctx context.Context, sessionID int, reason string) error {
	query := `UPDATE admin_sessions SET revoked_at = NOW(), revocation_reason = $1 WHERE id = $2 AND revoked_at IS NULL`
	_, err := r.db.Exec(ctx, query, reason, sessionID)
	return err
}

// RevokeAllUserSessions invalidates all active sessions for a user
func (r *AdminPasskeyRepo) RevokeAllUserSessions(ctx context.Context, userID int, reason string) error {
	query := `UPDATE admin_sessions SET revoked_at = NOW(), revocation_reason = $1 WHERE user_id = $2 AND revoked_at IS NULL`
	_, err := r.db.Exec(ctx, query, reason, userID)
	return err
}

// ListActiveUserSessions lists active sessions for an administrator
func (r *AdminPasskeyRepo) ListActiveUserSessions(ctx context.Context, userID int) ([]models.AdminSession, error) {
	query := `
		SELECT id, session_token_hash, user_id, portal, role, ip_address, user_agent,
		       last_active_at, idle_timeout_seconds, absolute_expires_at, created_at
		FROM admin_sessions
		WHERE user_id = $1 AND revoked_at IS NULL AND absolute_expires_at > NOW()
		ORDER BY last_active_at DESC
	`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.AdminSession
	for rows.Next() {
		var s models.AdminSession
		if err := rows.Scan(
			&s.ID, &s.SessionTokenHash, &s.UserID, &s.Portal, &s.Role, &s.IPAddress, &s.UserAgent,
			&s.LastActiveAt, &s.IdleTimeoutSeconds, &s.AbsoluteExpiresAt, &s.CreatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, s)
	}
	return list, rows.Err()
}

// RecordAuthAttempt logs an authentication attempt for progressive throttling
func (r *AdminPasskeyRepo) RecordAuthAttempt(ctx context.Context, attempt *models.AuthAttempt) error {
	query := `
		INSERT INTO admin_auth_attempts (identifier, ip_address, portal, success, failure_reason, attempted_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
	`
	_, err := r.db.Exec(ctx, query,
		attempt.Identifier,
		attempt.IPAddress,
		attempt.Portal,
		attempt.Success,
		attempt.FailureReason,
	)
	return err
}

// CountRecentFailedAttempts counts failed attempts within the given lookback window
func (r *AdminPasskeyRepo) CountRecentFailedAttempts(ctx context.Context, identifier, ip string, window time.Duration) (int, error) {
	since := time.Now().Add(-window)
	query := `
		SELECT COUNT(*)
		FROM admin_auth_attempts
		WHERE (identifier = $1 OR ip_address = $2)
		  AND success = FALSE
		  AND attempted_at >= $3
	`
	var count int
	err := r.db.QueryRow(ctx, query, identifier, ip, since).Scan(&count)
	return count, err
}

// SaveReauthToken stores a sensitive action reauth token
func (r *AdminPasskeyRepo) SaveReauthToken(ctx context.Context, token *models.ReauthToken) error {
	query := `
		INSERT INTO admin_reauth_tokens (user_id, token_hash, action_scope, expires_at, created_at)
		VALUES ($1, $2, $3, $4, NOW())
		RETURNING id, created_at
	`
	return r.db.QueryRow(ctx, query,
		token.UserID,
		token.TokenHash,
		token.ActionScope,
		token.ExpiresAt,
	).Scan(&token.ID, &token.CreatedAt)
}

// ValidateAndConsumeReauthToken checks that a reauth token is valid and marks it used
func (r *AdminPasskeyRepo) ValidateAndConsumeReauthToken(ctx context.Context, tokenHash string, userID int, actionScope string) (bool, error) {
	query := `
		UPDATE admin_reauth_tokens
		SET used_at = NOW()
		WHERE token_hash = $1 AND user_id = $2 AND action_scope = $3 AND used_at IS NULL AND expires_at > NOW()
		RETURNING id
	`
	var id int
	err := r.db.QueryRow(ctx, query, tokenHash, userID, actionScope).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
