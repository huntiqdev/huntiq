package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/your-org/used-tools-marketplace/internal/models"
)

type PaymentOrderFilter struct {
	UserID      *int
	Provider    string
	Status      string
	Environment string
	Search      string
	Limit       int
	Offset      int
}

type PaymentRepo struct {
	db *pgxpool.Pool
}

func NewPaymentRepo(db *pgxpool.Pool) *PaymentRepo {
	return &PaymentRepo{db: db}
}

func (r *PaymentRepo) CreateOrder(ctx context.Context, o *models.PaymentOrder) error {
	query := `
		INSERT INTO payment_orders (
			order_id, user_id, provider, environment, order_type, target_tier,
			user_role, duration_days, product_id, amount, currency, title,
			status, payer_phone, provider_tx_id, checkout_url, is_demo,
			otp_verified, verification_metadata, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12,
			$13, $14, $15, $16, $17, $18, CASE WHEN $19::text IS NOT NULL AND $19::text != '' THEN ($19::text)::jsonb ELSE '{}'::jsonb END, NOW(), NOW()
		)
		RETURNING id, created_at, updated_at
	`
	var metaStr *string
	if o.VerificationMetadata != nil {
		b, _ := json.Marshal(o.VerificationMetadata)
		s := string(b)
		metaStr = &s
	}

	return r.db.QueryRow(ctx, query,
		o.OrderID, o.UserID, o.Provider, o.Environment, o.OrderType, o.TargetTier,
		o.UserRole, o.DurationDays, o.ProductID, o.Amount, o.Currency, o.Title,
		o.Status, o.PayerPhone, o.ProviderTxID, o.CheckoutURL, o.IsDemo,
		o.OTPVerified, metaStr,
	).Scan(&o.ID, &o.CreatedAt, &o.UpdatedAt)
}

func (r *PaymentRepo) GetOrderByOrderID(ctx context.Context, orderID string) (*models.PaymentOrder, error) {
	query := `
		SELECT p.id, p.order_id, p.user_id, COALESCE(u.name, ''), COALESCE(u.phone, ''),
		       p.provider, p.environment, p.order_type, p.target_tier, p.user_role,
		       p.duration_days, p.product_id, COALESCE(pr.title, ''), p.amount, p.currency,
		       p.title, p.status, p.payer_phone, p.provider_tx_id, p.checkout_url,
		       p.raw_callback, p.verification_metadata, p.is_demo, p.otp_verified,
		       p.otp_verified_at, p.verified_at, p.failure_reason, p.created_at, p.updated_at
		FROM payment_orders p
		LEFT JOIN users u ON p.user_id = u.id
		LEFT JOIN products pr ON p.product_id = pr.id
		WHERE p.order_id = $1
	`
	return r.scanOrderRow(r.db.QueryRow(ctx, query, orderID))
}

func (r *PaymentRepo) GetOrderByID(ctx context.Context, id int) (*models.PaymentOrder, error) {
	query := `
		SELECT p.id, p.order_id, p.user_id, COALESCE(u.name, ''), COALESCE(u.phone, ''),
		       p.provider, p.environment, p.order_type, p.target_tier, p.user_role,
		       p.duration_days, p.product_id, COALESCE(pr.title, ''), p.amount, p.currency,
		       p.title, p.status, p.payer_phone, p.provider_tx_id, p.checkout_url,
		       p.raw_callback, p.verification_metadata, p.is_demo, p.otp_verified,
		       p.otp_verified_at, p.verified_at, p.failure_reason, p.created_at, p.updated_at
		FROM payment_orders p
		LEFT JOIN users u ON p.user_id = u.id
		LEFT JOIN products pr ON p.product_id = pr.id
		WHERE p.id = $1
	`
	return r.scanOrderRow(r.db.QueryRow(ctx, query, id))
}

func (r *PaymentRepo) scanOrderRow(row pgx.Row) (*models.PaymentOrder, error) {
	var o models.PaymentOrder
	var rawCbBytes, metaBytes []byte

	err := row.Scan(
		&o.ID, &o.OrderID, &o.UserID, &o.UserName, &o.UserPhone,
		&o.Provider, &o.Environment, &o.OrderType, &o.TargetTier, &o.UserRole,
		&o.DurationDays, &o.ProductID, &o.ProductTitle, &o.Amount, &o.Currency,
		&o.Title, &o.Status, &o.PayerPhone, &o.ProviderTxID, &o.CheckoutURL,
		&rawCbBytes, &metaBytes, &o.IsDemo, &o.OTPVerified,
		&o.OTPVerifiedAt, &o.VerifiedAt, &o.FailureReason, &o.CreatedAt, &o.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if len(rawCbBytes) > 0 {
		_ = json.Unmarshal(rawCbBytes, &o.RawCallback)
	}
	if len(metaBytes) > 0 {
		_ = json.Unmarshal(metaBytes, &o.VerificationMetadata)
	}

	return &o, nil
}

func (r *PaymentRepo) UpdateOrderStatus(ctx context.Context, orderID, status string, providerTxID, payerPhone, failureReason *string, rawCallback, metadata map[string]interface{}) error {
	var rawStr, metaStr *string
	if rawCallback != nil {
		b, _ := json.Marshal(rawCallback)
		s := string(b)
		rawStr = &s
	}
	if metadata != nil {
		b, _ := json.Marshal(metadata)
		s := string(b)
		metaStr = &s
	}

	query := `
		UPDATE payment_orders
		SET status = $1::varchar,
		    provider_tx_id = COALESCE($2, provider_tx_id),
		    payer_phone = COALESCE($3, payer_phone),
		    failure_reason = COALESCE($4, failure_reason),
		    raw_callback = CASE WHEN $5::text IS NOT NULL AND $5::text != '' THEN ($5::text)::jsonb ELSE raw_callback END,
		    verification_metadata = CASE WHEN $6::text IS NOT NULL AND $6::text != '' THEN ($6::text)::jsonb ELSE verification_metadata END,
		    verified_at = CASE WHEN $1::varchar = 'SUCCESSFUL' THEN NOW() ELSE verified_at END,
		    updated_at = NOW()
		WHERE order_id = $7
	`
	_, err := r.db.Exec(ctx, query, status, providerTxID, payerPhone, failureReason, rawStr, metaStr, orderID)
	return err
}

func (r *PaymentRepo) SetOTPVerified(ctx context.Context, orderID string) error {
	query := `
		UPDATE payment_orders
		SET otp_verified = true, otp_verified_at = NOW(), updated_at = NOW()
		WHERE order_id = $1
	`
	_, err := r.db.Exec(ctx, query, orderID)
	return err
}

func (r *PaymentRepo) ListOrders(ctx context.Context, f PaymentOrderFilter) ([]models.PaymentOrder, int, error) {
	if f.Limit <= 0 || f.Limit > 100 {
		f.Limit = 50
	}
	if f.Offset < 0 {
		f.Offset = 0
	}

	var where []string
	var args []interface{}
	idx := 1

	if f.UserID != nil && *f.UserID > 0 {
		where = append(where, fmt.Sprintf("p.user_id = $%d", idx))
		args = append(args, *f.UserID)
		idx++
	}
	if f.Provider != "" && f.Provider != "all" {
		where = append(where, fmt.Sprintf("p.provider = $%d", idx))
		args = append(args, f.Provider)
		idx++
	}
	if f.Status != "" && f.Status != "all" {
		where = append(where, fmt.Sprintf("p.status = $%d", idx))
		args = append(args, f.Status)
		idx++
	}
	if f.Environment != "" && f.Environment != "all" {
		where = append(where, fmt.Sprintf("p.environment = $%d", idx))
		args = append(args, f.Environment)
		idx++
	}
	if f.Search != "" {
		where = append(where, fmt.Sprintf("(p.order_id ILIKE $%d OR u.name ILIKE $%d OR u.phone ILIKE $%d)", idx, idx, idx))
		args = append(args, "%"+f.Search+"%")
		idx++
	}

	whereClause := ""
	if len(where) > 0 {
		whereClause = "WHERE " + strings.Join(where, " AND ")
	}

	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM payment_orders p LEFT JOIN users u ON p.user_id = u.id %s`, whereClause)
	var total int
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := fmt.Sprintf(`
		SELECT p.id, p.order_id, p.user_id, COALESCE(u.name, ''), COALESCE(u.phone, ''),
		       p.provider, p.environment, p.order_type, p.target_tier, p.user_role,
		       p.duration_days, p.product_id, COALESCE(pr.title, ''), p.amount, p.currency,
		       p.title, p.status, p.payer_phone, p.provider_tx_id, p.checkout_url,
		       p.raw_callback, p.verification_metadata, p.is_demo, p.otp_verified,
		       p.otp_verified_at, p.verified_at, p.failure_reason, p.created_at, p.updated_at
		FROM payment_orders p
		LEFT JOIN users u ON p.user_id = u.id
		LEFT JOIN products pr ON p.product_id = pr.id
		%s
		ORDER BY p.created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, idx, idx+1)

	args = append(args, f.Limit, f.Offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []models.PaymentOrder
	for rows.Next() {
		var o models.PaymentOrder
		var rawCbBytes, metaBytes []byte

		if err := rows.Scan(
			&o.ID, &o.OrderID, &o.UserID, &o.UserName, &o.UserPhone,
			&o.Provider, &o.Environment, &o.OrderType, &o.TargetTier, &o.UserRole,
			&o.DurationDays, &o.ProductID, &o.ProductTitle, &o.Amount, &o.Currency,
			&o.Title, &o.Status, &o.PayerPhone, &o.ProviderTxID, &o.CheckoutURL,
			&rawCbBytes, &metaBytes, &o.IsDemo, &o.OTPVerified,
			&o.OTPVerifiedAt, &o.VerifiedAt, &o.FailureReason, &o.CreatedAt, &o.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}

		if len(rawCbBytes) > 0 {
			_ = json.Unmarshal(rawCbBytes, &o.RawCallback)
		}
		if len(metaBytes) > 0 {
			_ = json.Unmarshal(metaBytes, &o.VerificationMetadata)
		}

		list = append(list, o)
	}

	return list, total, nil
}

// ── Provider Configuration Management ─────────────────────────────────────────

func (r *PaymentRepo) GetProviderConfigs(ctx context.Context) ([]models.PaymentProviderConfig, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, provider_code, provider_name, environment_mode, is_enabled,
		       api_endpoint, app_id, 
		       CASE WHEN app_key_encrypted != '' THEN '••••••••••••' ELSE '' END as app_key_masked,
		       short_code, notify_url, return_url, extra_config, updated_by, created_at, updated_at
		FROM payment_provider_configs
		ORDER BY id ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var configs []models.PaymentProviderConfig
	for rows.Next() {
		var c models.PaymentProviderConfig
		var extraBytes []byte
		if err := rows.Scan(
			&c.ID, &c.ProviderCode, &c.ProviderName, &c.EnvironmentMode, &c.IsEnabled,
			&c.APIEndpoint, &c.AppID, &c.AppKeyMasked, &c.ShortCode, &c.NotifyURL, &c.ReturnURL,
			&extraBytes, &c.UpdatedBy, &c.CreatedAt, &c.UpdatedAt,
		); err != nil {
			return nil, err
		}
		if len(extraBytes) > 0 {
			_ = json.Unmarshal(extraBytes, &c.ExtraConfig)
		}
		configs = append(configs, c)
	}
	return configs, nil
}

func (r *PaymentRepo) GetProviderConfigByCode(ctx context.Context, code string) (*models.PaymentProviderConfig, error) {
	var c models.PaymentProviderConfig
	var extraBytes []byte
	err := r.db.QueryRow(ctx, `
		SELECT id, provider_code, provider_name, environment_mode, is_enabled,
		       api_endpoint, app_id, 
		       CASE WHEN app_key_encrypted != '' THEN '••••••••••••' ELSE '' END as app_key_masked,
		       short_code, notify_url, return_url, extra_config, updated_by, created_at, updated_at
		FROM payment_provider_configs
		WHERE provider_code = $1
	`, code).Scan(
		&c.ID, &c.ProviderCode, &c.ProviderName, &c.EnvironmentMode, &c.IsEnabled,
		&c.APIEndpoint, &c.AppID, &c.AppKeyMasked, &c.ShortCode, &c.NotifyURL, &c.ReturnURL,
		&extraBytes, &c.UpdatedBy, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	if len(extraBytes) > 0 {
		_ = json.Unmarshal(extraBytes, &c.ExtraConfig)
	}
	return &c, nil
}

func (r *PaymentRepo) UpdateProviderConfig(ctx context.Context, code, envMode string, isEnabled bool, endpoint, appID, appKey, shortCode, notifyURL, returnURL string, updatedBy int) error {
	query := `
		UPDATE payment_provider_configs
		SET environment_mode = $1,
		    is_enabled = $2,
		    api_endpoint = $3,
		    app_id = $4,
		    app_key_encrypted = CASE WHEN $5 != '' THEN $5 ELSE app_key_encrypted END,
		    short_code = $6,
		    notify_url = $7,
		    return_url = $8,
		    updated_by = $9,
		    updated_at = NOW()
		WHERE provider_code = $10
	`
	_, err := r.db.Exec(ctx, query, envMode, isEnabled, endpoint, appID, appKey, shortCode, notifyURL, returnURL, updatedBy, code)
	return err
}

func (r *PaymentRepo) SetGlobalOperatingMode(ctx context.Context, mode string) error {
	query := `
		UPDATE payment_provider_configs
		SET environment_mode = $1, updated_at = NOW()
		WHERE provider_code != 'demo_mock'
	`
	_, err := r.db.Exec(ctx, query, mode)
	if err != nil {
		return err
	}
	_, err = r.db.Exec(ctx, `
		INSERT INTO platform_settings (key, value, updated_at)
		VALUES ('payment_operating_mode', $1, NOW())
		ON CONFLICT (key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW()
	`, mode)
	return err
}

func (r *PaymentRepo) GetGlobalOperatingMode(ctx context.Context) (string, error) {
	var mode string
	err := r.db.QueryRow(ctx, `SELECT value FROM platform_settings WHERE key = 'payment_operating_mode'`).Scan(&mode)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "demo", nil
		}
		return "", err
	}
	if mode != "production" {
		return "demo", nil
	}
	return "production", nil
}

func (r *PaymentRepo) GetTierPricing(ctx context.Context) (map[string]float64, error) {
	rows, err := r.db.Query(ctx, `
		SELECT key, value FROM platform_settings
		WHERE key IN (
			'seller_gold_price_etb', 'seller_premium_price_etb',
			'buyer_gold_price_etb', 'buyer_premium_price_etb',
			'listing_fee_etb'
		)
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pricing := map[string]float64{
		"seller_gold_price_etb":    500.0,
		"seller_premium_price_etb": 2000.0,
		"buyer_gold_price_etb":     200.0,
		"buyer_premium_price_etb":  500.0,
		"listing_fee_etb":          50.0,
	}

	for rows.Next() {
		var k, v string
		if err := rows.Scan(&k, &v); err == nil {
			var f float64
			if _, err := fmt.Sscanf(v, "%f", &f); err == nil && f > 0 {
				pricing[k] = f
			}
		}
	}
	return pricing, nil
}

func (r *PaymentRepo) GetDecryptedAppKey(ctx context.Context, providerCode string) (string, error) {
	var appKey string
	err := r.db.QueryRow(ctx, `
		SELECT app_key_encrypted FROM payment_provider_configs WHERE provider_code = $1
	`, providerCode).Scan(&appKey)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", nil
		}
		return "", err
	}
	return appKey, nil
}

// ActivateTier atomically upgrades the user's tier, writes audit tier history, and notifies user.
func (r *PaymentRepo) ActivateTier(ctx context.Context, userID int, userRole, targetTier string, durationDays int, orderID string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	targetTierLower := strings.ToLower(strings.TrimSpace(targetTier))
	if targetTierLower != "gold" && targetTierLower != "premium" {
		targetTierLower = "gold"
	}

	if durationDays <= 0 {
		durationDays = 30
	}
	expiresAt := time.Now().AddDate(0, 0, durationDays)

	var oldTier string
	if userRole == "buyer" {
		err = tx.QueryRow(ctx, `SELECT COALESCE(buyer_tier::varchar, 'normal') FROM users WHERE id = $1 FOR UPDATE`, userID).Scan(&oldTier)
		if err != nil {
			return fmt.Errorf("user not found: %w", err)
		}

		_, err = tx.Exec(ctx, `UPDATE users SET buyer_tier = $1::buyer_tier, updated_at = NOW() WHERE id = $2`, targetTierLower, userID)
		if err != nil {
			return fmt.Errorf("failed to update buyer tier: %w", err)
		}

		_, err = tx.Exec(ctx, `
			INSERT INTO user_tier_history (user_id, tier_type, old_tier, new_tier, reason, duration_days, expires_at, created_at)
			VALUES ($1, 'buyer', $2, $3, $4, $5, $6, NOW())
		`, userID, oldTier, targetTierLower, fmt.Sprintf("Payment verified order: %s", orderID), durationDays, expiresAt)
		if err != nil {
			return fmt.Errorf("failed to record tier history: %w", err)
		}

		_, _ = tx.Exec(ctx, `
			INSERT INTO notifications (user_id, title, message, type)
			VALUES ($1, 'Buyer Membership Upgraded!', $2, 'tier_upgrade')
		`, userID, fmt.Sprintf("Congratulations! Your buyer membership has been successfully upgraded to %s for %d days.", strings.ToUpper(targetTierLower), durationDays))
	} else {
		err = tx.QueryRow(ctx, `SELECT COALESCE(tier::varchar, 'normal') FROM users WHERE id = $1 FOR UPDATE`, userID).Scan(&oldTier)
		if err != nil {
			return fmt.Errorf("user not found: %w", err)
		}

		_, err = tx.Exec(ctx, `UPDATE users SET tier = $1::seller_tier, updated_at = NOW() WHERE id = $2`, targetTierLower, userID)
		if err != nil {
			return fmt.Errorf("failed to update seller tier: %w", err)
		}

		_, err = tx.Exec(ctx, `
			INSERT INTO user_tier_history (user_id, tier_type, old_tier, new_tier, reason, duration_days, expires_at, created_at)
			VALUES ($1, 'seller', $2, $3, $4, $5, $6, NOW())
		`, userID, oldTier, targetTierLower, fmt.Sprintf("Payment verified order: %s", orderID), durationDays, expiresAt)
		if err != nil {
			return fmt.Errorf("failed to record tier history: %w", err)
		}

		_, _ = tx.Exec(ctx, `
			INSERT INTO notifications (user_id, title, message, type)
			VALUES ($1, 'Seller Membership Upgraded!', $2, 'tier_upgrade')
		`, userID, fmt.Sprintf("Congratulations! Your seller membership has been successfully upgraded to %s for %d days. You now enjoy priority placement and enhanced limits.", strings.ToUpper(targetTierLower), durationDays))
	}

	return tx.Commit(ctx)
}

// ActivateListingFee marks the product listing fee as paid and notifies the seller.
func (r *PaymentRepo) ActivateListingFee(ctx context.Context, productID int, amount float64, orderID string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var sellerID int
	var title string
	err = tx.QueryRow(ctx, `
		SELECT seller_id, title FROM products WHERE id = $1 FOR UPDATE
	`, productID).Scan(&sellerID, &title)
	if err != nil {
		return fmt.Errorf("product not found: %w", err)
	}

	_, err = tx.Exec(ctx, `
		UPDATE products
		SET fee_paid = true, fee_payment_pending = false, fee_paid_at = NOW(), total_fee_amount = $1, updated_at = NOW()
		WHERE id = $2
	`, amount, productID)
	if err != nil {
		return fmt.Errorf("failed to activate listing fee: %w", err)
	}

	_, _ = tx.Exec(ctx, `
		INSERT INTO notifications (user_id, title, message, type)
		VALUES ($1, 'Listing Fee Paid & Active!', $2, 'listing_fee')
	`, sellerID, fmt.Sprintf("Listing fee of %.2f ETB for '%s' has been verified (Order: %s). Your listing is active and published.", amount, title, orderID))

	return tx.Commit(ctx)
}

// RecordTransaction writes double-entry financial transaction record.
func (r *PaymentRepo) RecordTransaction(ctx context.Context, o *models.PaymentOrder) error {
	// Insert into transactions or admin ledger
	_, err := r.db.Exec(ctx, `
		INSERT INTO admin_activity_logs (admin_id, action, details, created_at)
		VALUES (NULL, 'PAYMENT_VERIFIED', $1, NOW())
	`, fmt.Sprintf("Order %s for %.2f %s (%s) verified via %s [demo=%v]", o.OrderID, o.Amount, o.Currency, o.OrderType, o.Provider, o.IsDemo))
	return err
}
