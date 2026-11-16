package repository

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/your-org/used-tools-marketplace/internal/models"
)

type SMSRepo struct {
	db *pgxpool.Pool
}

func NewSMSRepo(db *pgxpool.Pool) *SMSRepo {
	return &SMSRepo{db: db}
}

// GetActiveProviderConfig returns the default active SMS provider configuration.
func (r *SMSRepo) GetActiveProviderConfig(ctx context.Context) (*models.SMSProviderConfig, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, provider_code, provider_name, api_endpoint, sender_id, credentials, is_active, is_default, daily_limit, daily_sent, created_at, updated_at
		FROM sms_provider_configs
		WHERE is_active = true
		ORDER BY is_default DESC, id ASC
		LIMIT 1
	`)

	var cfg models.SMSProviderConfig
	var credBytes []byte

	err := row.Scan(
		&cfg.ID,
		&cfg.ProviderCode,
		&cfg.ProviderName,
		&cfg.APIEndpoint,
		&cfg.SenderID,
		&credBytes,
		&cfg.IsActive,
		&cfg.IsDefault,
		&cfg.DailyLimit,
		&cfg.DailySent,
		&cfg.CreatedAt,
		&cfg.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // No active provider configured
		}
		return nil, err
	}

	if len(credBytes) > 0 {
		_ = json.Unmarshal(credBytes, &cfg.Credentials)
	}
	if cfg.Credentials == nil {
		cfg.Credentials = make(map[string]interface{})
	}

	return &cfg, nil
}

// SaveDeliveryLog records an SMS delivery attempt.
func (r *SMSRepo) SaveDeliveryLog(ctx context.Context, entry *models.SMSDeliveryLog) (int64, error) {
	var id int64
	err := r.db.QueryRow(ctx, `
		INSERT INTO sms_delivery_logs (phone_masked, purpose, provider, provider_ref, status, error_category, retry_count, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`, entry.PhoneMasked, entry.Purpose, entry.Provider, entry.ProviderRef, entry.Status, entry.ErrorCategory, entry.RetryCount, time.Now()).Scan(&id)
	if err != nil {
		return 0, err
	}
	entry.ID = id
	return id, nil
}

// UpdateDeliveryLogStatus updates an existing SMS delivery log with provider response.
func (r *SMSRepo) UpdateDeliveryLogStatus(ctx context.Context, id int64, status string, providerRef *string, errCategory *string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE sms_delivery_logs
		SET status = $1, provider_ref = COALESCE($2, provider_ref), error_category = $3
		WHERE id = $4
	`, status, providerRef, errCategory, id)
	return err
}

// GetDeliveryLogs returns paginated SMS audit logs.
func (r *SMSRepo) GetDeliveryLogs(ctx context.Context, limit, offset int) ([]models.SMSDeliveryLog, int, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	var total int
	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM sms_delivery_logs`).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.db.Query(ctx, `
		SELECT id, phone_masked, purpose, provider, provider_ref, status, error_category, retry_count, created_at
		FROM sms_delivery_logs
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []models.SMSDeliveryLog
	for rows.Next() {
		var l models.SMSDeliveryLog
		if err := rows.Scan(
			&l.ID,
			&l.PhoneMasked,
			&l.Purpose,
			&l.Provider,
			&l.ProviderRef,
			&l.Status,
			&l.ErrorCategory,
			&l.RetryCount,
			&l.CreatedAt,
		); err != nil {
			return nil, 0, err
		}
		list = append(list, l)
	}

	return list, total, nil
}

// IncrementDailySent increments the sent counter for an SMS provider config.
func (r *SMSRepo) IncrementDailySent(ctx context.Context, providerID int) error {
	_, err := r.db.Exec(ctx, `
		UPDATE sms_provider_configs
		SET daily_sent = daily_sent + 1, updated_at = NOW()
		WHERE id = $1
	`, providerID)
	return err
}

// GetProviderConfigs returns all SMS provider configs.
func (r *SMSRepo) GetProviderConfigs(ctx context.Context) ([]models.SMSProviderConfig, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, provider_code, provider_name, api_endpoint, sender_id, credentials, is_active, is_default, daily_limit, daily_sent, created_at, updated_at
		FROM sms_provider_configs
		ORDER BY id ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.SMSProviderConfig
	for rows.Next() {
		var cfg models.SMSProviderConfig
		var credBytes []byte
		if err := rows.Scan(
			&cfg.ID,
			&cfg.ProviderCode,
			&cfg.ProviderName,
			&cfg.APIEndpoint,
			&cfg.SenderID,
			&credBytes,
			&cfg.IsActive,
			&cfg.IsDefault,
			&cfg.DailyLimit,
			&cfg.DailySent,
			&cfg.CreatedAt,
			&cfg.UpdatedAt,
		); err != nil {
			return nil, err
		}
		if len(credBytes) > 0 {
			_ = json.Unmarshal(credBytes, &cfg.Credentials)
		}
		list = append(list, cfg)
	}
	return list, nil
}
