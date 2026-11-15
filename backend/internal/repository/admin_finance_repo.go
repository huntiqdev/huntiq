package repository

import (
	"context"
	"fmt"
)

func (r *AdminRepo) AnalyticsOverview(ctx context.Context) (map[string]interface{}, error) {
	out := map[string]interface{}{}

	// Total revenue (broker fees from completed transactions)
	var totalRevenue, totalVolume float64
	_ = r.db.QueryRow(ctx, `SELECT COALESCE(SUM(broker_fee),0), COALESCE(SUM(amount),0) FROM transactions WHERE status='completed'`).Scan(&totalRevenue, &totalVolume) //nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
	out["total_revenue"] = totalRevenue
	out["total_volume"] = totalVolume

	// Listing fee revenue (verified payments on approved/expired/sold products only)
	var listingFeeRevenue, totalPlatformProfit float64
	_ = r.db.QueryRow(ctx, `SELECT COALESCE(SUM(total_fee_amount),0) FROM products WHERE fee_paid=TRUE AND status IN ('approved','expired','sold')`).Scan(&listingFeeRevenue) //nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
	_ = r.db.QueryRow(ctx, `SELECT COALESCE(total_platform_profit,0) FROM admin_profit_summary LIMIT 1`).Scan(&totalPlatformProfit)
	out["listing_fee_revenue"] = listingFeeRevenue
	out["net_profit"] = totalPlatformProfit

	// Counts
	var totalProducts, pendingProducts, approvedProducts, expiredProducts int
	_ = r.db.QueryRow(ctx, `SELECT COUNT(*), COUNT(*) FILTER (WHERE status='pending'), COUNT(*) FILTER (WHERE status='approved'), COUNT(*) FILTER (WHERE status='expired') FROM products`).Scan(&totalProducts, &pendingProducts, &approvedProducts, &expiredProducts) //nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
	out["total_products"] = totalProducts
	out["pending_products"] = pendingProducts
	out["approved_products"] = approvedProducts
	out["expired_products"] = expiredProducts

	// Seller counts by tier
	if rows, err := r.db.Query(ctx, `SELECT COALESCE(tier::text,'normal'), COUNT(*) FROM users WHERE role='seller' GROUP BY tier`); err == nil { //nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
		defer rows.Close()
		tierCounts := map[string]int{}
		for rows.Next() {
			var t string
			var c int
			if rows.Scan(&t, &c) == nil {
				tierCounts[t] = c
			}
		}
		out["seller_tiers"] = tierCounts
	}

	// Buyer counts by tier
	if brows, err := r.db.Query(ctx, `SELECT COALESCE(buyer_tier::text,'normal'), COUNT(*) FROM users WHERE role='buyer' GROUP BY buyer_tier`); err == nil { //nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
		defer brows.Close()
		btierCounts := map[string]int{}
		for brows.Next() {
			var t string
			var c int
			if brows.Scan(&t, &c) == nil {
				btierCounts[t] = c
			}
		}
		out["buyer_tiers"] = btierCounts
	}

	// Orders by status
	if orows, err := r.db.Query(ctx, `SELECT status, COUNT(*) FROM order_requests GROUP BY status`); err == nil { //nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
		defer orows.Close()
		orderCounts := map[string]int{}
		for orows.Next() {
			var s string
			var c int
			if orows.Scan(&s, &c) == nil {
				orderCounts[s] = c
			}
		}
		out["order_counts"] = orderCounts
	}

	// Seller totals
	var totalSellers, activeSellers int
	_ = r.db.QueryRow(ctx, `SELECT COUNT(*), COUNT(*) FILTER (WHERE is_active=TRUE) FROM users WHERE role='seller'`).Scan(&totalSellers, &activeSellers) //nolint:gosec
	out["total_sellers"] = totalSellers
	out["active_sellers"] = activeSellers

	// Buyer totals
	var totalBuyers, premiumBuyers int
	_ = r.db.QueryRow(ctx, `SELECT COUNT(*), COUNT(*) FILTER (WHERE buyer_tier IN ('gold','premium')) FROM users WHERE role='buyer'`).Scan(&totalBuyers, &premiumBuyers) //nolint:gosec
	out["total_buyers"] = totalBuyers
	out["premium_buyers"] = premiumBuyers

	return out, nil
}

// AnalyticsWoredaHeatmap returns product counts grouped by woreda.
func (r *AdminRepo) AnalyticsWoredaHeatmap(ctx context.Context) ([]map[string]interface{}, error) {
	rows, err := //nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
		r.db.Query(ctx, `
		SELECT l.region, COALESCE(l.zone,'') AS zone, COALESCE(l.woreda,'') AS woreda,
		       COUNT(p.id) AS product_count,
		       COUNT(p.id) FILTER (WHERE p.status='approved') AS active_count
		FROM locations l
		LEFT JOIN products p ON p.location_id = l.id
		GROUP BY l.region, l.zone, l.woreda
		HAVING COUNT(p.id) > 0
		ORDER BY product_count DESC
		LIMIT 50`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []map[string]interface{}
	for rows.Next() {
		var region, zone, woreda string
		var total, active int
		if err := rows.Scan(&region, &zone, &woreda, &total, &active); err != nil {
			return nil, err
		}
		out = append(out, map[string]interface{}{
			"region": region, "zone": zone, "woreda": woreda,
			"total": total, "active": active,
		})
	}
	return out, rows.Err()
}

// AnalyticsMonthly returns monthly revenue and listing counts for the last N months.
func (r *AdminRepo) AnalyticsMonthly(ctx context.Context, months int) ([]map[string]interface{}, error) {
	rows, err := //nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
		r.db.Query(ctx, `
		SELECT
		  TO_CHAR(DATE_TRUNC('month', gs), 'Mon YYYY') AS label,
		  DATE_TRUNC('month', gs) AS month_start,
		  COALESCE(SUM(t.broker_fee), 0) AS revenue,
		  COALESCE(SUM(t.amount), 0) AS volume,
		  COUNT(DISTINCT t.id) AS transactions,
		  COUNT(DISTINCT p.id) AS new_listings
		FROM generate_series(
		  DATE_TRUNC('month', NOW()) - (($1::int - 1) * INTERVAL '1 month'),
		  DATE_TRUNC('month', NOW()),
		  INTERVAL '1 month'
		) AS gs
		LEFT JOIN transactions t ON DATE_TRUNC('month', t.created_at) = gs AND t.status='completed'
		LEFT JOIN products p ON DATE_TRUNC('month', p.created_at) = gs
		GROUP BY gs
		ORDER BY gs`, months)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []map[string]interface{}
	for rows.Next() {
		var label string
		var monthStart interface{}
		var revenue, volume float64
		var txCount, listings int
		if err := rows.Scan(&label, &monthStart, &revenue, &volume, &txCount, &listings); err != nil {
			return nil, err
		}
		out = append(out, map[string]interface{}{
			"label": label, "revenue": revenue, "volume": volume,
			"transactions": txCount, "new_listings": listings,
		})
	}
	return out, rows.Err()
}

// AnalyticsGantt returns product listing durations for Gantt chart.
func (r *AdminRepo) AnalyticsGantt(ctx context.Context, limit int) ([]map[string]interface{}, error) {
	rows, err := //nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
		r.db.Query(ctx, `
		SELECT p.id, p.title, p.status,
		       COALESCE(tier::text,'normal') AS seller_tier,
		       p.created_at, p.expires_at,
		       p.listing_duration_days,
		       COALESCE(l.region,'') AS region
		FROM products p
		JOIN users u ON u.id = p.seller_id
		LEFT JOIN locations l ON l.id = p.location_id
		WHERE p.listing_duration_days IS NOT NULL
		ORDER BY p.created_at DESC
		LIMIT $1`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []map[string]interface{}
	for rows.Next() {
		var id, durationDays int
		var title, status, sellerTier, region string
		var createdAt interface{}
		var expiresAt interface{}
		if err := rows.Scan(&id, &title, &status, &sellerTier, &createdAt, &expiresAt, &durationDays, &region); err != nil {
			return nil, err
		}
		out = append(out, map[string]interface{}{
			"id": id, "title": title, "status": status, "seller_tier": sellerTier,
			"created_at": createdAt, "expires_at": expiresAt,
			"duration_days": durationDays, "region": region,
		})
	}
	return out, rows.Err()
}

// AnalyticsFeeRevenue returns fee revenue grouped by region.
func (r *AdminRepo) AnalyticsFeeRevenue(ctx context.Context) ([]map[string]interface{}, error) {
	rows, err := //nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
		r.db.Query(ctx, `
		SELECT COALESCE(l.region,'Unknown') AS region,
		       COALESCE(SUM(p.total_fee_amount),0) AS total_fees,
		       COALESCE(SUM(p.total_fee_amount) FILTER (WHERE p.fee_paid=TRUE),0) AS paid_fees,
		       COUNT(p.id) AS listings
		FROM products p
		LEFT JOIN locations l ON l.id = p.location_id
		WHERE p.total_fee_amount IS NOT NULL
		GROUP BY l.region
		ORDER BY total_fees DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []map[string]interface{}
	for rows.Next() {
		var region string
		var totalFees, paidFees float64
		var listings int
		if err := rows.Scan(&region, &totalFees, &paidFees, &listings); err != nil {
			return nil, err
		}
		out = append(out, map[string]interface{}{
			"region": region, "total_fees": totalFees, "paid_fees": paidFees, "listings": listings,
		})
	}
	return out, rows.Err()
}

// ── Platform Settings ─────────────────────────────────────────────────────────

func (r *AdminRepo) GetAllSettings(ctx context.Context) (map[string]string, error) {
	rows, err := r.db.Query(ctx, `SELECT key, value FROM platform_settings ORDER BY key`) //nolint:gosec // G201: fully static query, no user input
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := map[string]string{}
	for rows.Next() {
		var k, v string
		if err := rows.Scan(&k, &v); err != nil {
			return nil, err
		}
		out[k] = v
	}
	return out, rows.Err()
}

// GetSetting returns a single platform setting value, or "" if not found.
func (r *AdminRepo) GetSetting(ctx context.Context, key string) (string, error) {
	var value string
	err := r.db.QueryRow(ctx, `SELECT value FROM platform_settings WHERE key=$1`, key).Scan(&value) //nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
	if err != nil {
		return "", nil // key not found — return empty, not an error
	}
	return value, nil
}

func (r *AdminRepo) UpdateSetting(ctx context.Context, key, value string, adminID int) error {
	// Read old value for audit log
	var oldValue string
	_ = r.db.QueryRow(ctx, `SELECT value FROM platform_settings WHERE key=$1`, key).Scan(&oldValue) //nolint:gosec // G201/G202: parameterized query, no user input in SQL structure
	_, err := r.db.Exec(ctx, `
		INSERT INTO platform_settings (key, value, updated_by, updated_at)
		VALUES ($1, $2, $3, NOW())
		ON CONFLICT (key) DO UPDATE SET value=$2, updated_by=$3, updated_at=NOW()`,
		key, value, adminID)
	if err != nil {
		return err
	}
	// Write audit log (non-fatal — don't fail the update if log insert fails)
	_, _ = r.db.Exec(ctx, `
		INSERT INTO settings_change_log (admin_id, setting_key, old_value, new_value, changed_at)
		VALUES ($1, $2, $3, $4, NOW())`,
		adminID, key, oldValue, value)
	return nil
}

// UpdateSettingWithReason updates a setting and records the reason in the audit log.
func (r *AdminRepo) UpdateSettingWithReason(ctx context.Context, key, value, reason string, adminID int) error {
	var oldValue string
	_ = r.db.QueryRow(ctx, `SELECT value FROM platform_settings WHERE key=$1`, key).Scan(&oldValue) //nolint:gosec // G201/G202: parameterized query, no user input in SQL structure
	_, err := r.db.Exec(ctx, `
		INSERT INTO platform_settings (key, value, updated_by, updated_at)
		VALUES ($1, $2, $3, NOW())
		ON CONFLICT (key) DO UPDATE SET value=$2, updated_by=$3, updated_at=NOW()`,
		key, value, adminID)
	if err != nil {
		return err
	}
	var reasonVal interface{}
	if reason != "" {
		reasonVal = reason
	}
	_, _ = r.db.Exec(ctx, `
		INSERT INTO settings_change_log (admin_id, setting_key, old_value, new_value, reason, changed_at)
		VALUES ($1, $2, $3, $4, $5, NOW())`,
		adminID, key, oldValue, value, reasonVal)
	return nil
}

// GetSettingsChangeLog returns the audit log for a specific setting key (or all keys if empty).
func (r *AdminRepo) GetSettingsChangeLog(ctx context.Context, key string, limit int) ([]map[string]interface{}, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	var rows interface {
		Next() bool
		Scan(...interface{}) error
		Err() error
		Close()
	}
	var err error
	if key != "" {
		rows, err = //nolint:gosec // G201/G202: parameterized query, no user input in SQL structure
			r.db.Query(ctx, `
			SELECT l.id, l.setting_key, l.old_value, l.new_value, COALESCE(l.reason,''),
			       COALESCE(u.name,''), l.admin_id, l.changed_at
			FROM settings_change_log l
			LEFT JOIN users u ON u.id = l.admin_id
			WHERE l.setting_key = $1
			ORDER BY l.changed_at DESC LIMIT $2`, key, limit)
	} else {
		rows, err = //nolint:gosec // G201/G202: parameterized query, no user input in SQL structure
			r.db.Query(ctx, `
			SELECT l.id, l.setting_key, l.old_value, l.new_value, COALESCE(l.reason,''),
			       COALESCE(u.name,''), l.admin_id, l.changed_at
			FROM settings_change_log l
			LEFT JOIN users u ON u.id = l.admin_id
			ORDER BY l.changed_at DESC LIMIT $1`, limit)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []map[string]interface{}
	for rows.Next() {
		var id, adminID int
		var settingKey, oldVal, newVal, reason, adminName string
		var changedAt interface{}
		if err := rows.Scan(&id, &settingKey, &oldVal, &newVal, &reason, &adminName, &adminID, &changedAt); err != nil {
			return nil, err
		}
		out = append(out, map[string]interface{}{
			"id": id, "setting_key": settingKey,
			"old_value": oldVal, "new_value": newVal, "reason": reason,
			"admin_name": adminName, "admin_id": adminID, "changed_at": changedAt,
		})
	}
	if out == nil {
		out = []map[string]interface{}{}
	}
	return out, rows.Err()
}

// ── Admin Notifications ───────────────────────────────────────────────────────

func (r *AdminRepo) ListAdminNotifications(ctx context.Context, adminID int, adminRole string) ([]map[string]interface{}, error) {
	rows, err := //nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
		r.db.Query(ctx, `
		SELECT id, title, message, type,
		       COALESCE(target_role::text,'') AS target_role,
		       ($1 = ANY(read_by)) AS is_read_by_me,
		       created_at
		FROM admin_notifications
		WHERE target_role IS NULL OR target_role::text = $2
		ORDER BY created_at DESC
		LIMIT 50`, adminID, adminRole)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []map[string]interface{}
	for rows.Next() {
		var id int
		var title, message, ntype, targetRole string
		var isRead bool
		var createdAt interface{}
		if err := rows.Scan(&id, &title, &message, &ntype, &targetRole, &isRead, &createdAt); err != nil {
			return nil, err
		}
		out = append(out, map[string]interface{}{
			"id": id, "title": title, "message": message, "type": ntype,
			"target_role": targetRole, "is_read": isRead, "created_at": createdAt,
		})
	}
	return out, rows.Err()
}

func (r *AdminRepo) MarkNotificationRead(ctx context.Context, notifID, adminID int) error {
	_, err := r.db.Exec(ctx, `
		UPDATE admin_notifications
		SET read_by = array_append(read_by, $1),
		    is_read  = TRUE
		WHERE id = $2 AND NOT ($1 = ANY(read_by))`,
		adminID, notifID)
	return err
}

func (r *AdminRepo) CreateAdminNotification(ctx context.Context, title, message, ntype, targetRole string) error {
	var tr interface{}
	if targetRole != "" {
		tr = targetRole
	}
	priority := "normal"
	switch ntype {
	case "error", "critical":
		priority = "high"
	case "warning":
		priority = "medium"
	case "low":
		priority = "low"
	}
	_, err := r.db.Exec(ctx, `
		INSERT INTO admin_notifications (title, message, type, priority, target_role)
		VALUES ($1, $2, $3, $4, $5)`,
		title, message, ntype, priority, tr)
	return err
}

// ── Profit Analytics (Super Admin Only) ──────────────────────────────────────

// ProfitSummary returns the full platform profit breakdown from the view.
func (r *AdminRepo) ProfitSummary(ctx context.Context) (map[string]interface{}, error) {
	row := r.db.QueryRow(ctx, `SELECT * FROM admin_profit_summary LIMIT 1`) //nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
	var totalTx int
	var grossVolume, brokerFeeRevenue, brokerFeePct float64
	var listingFeeCollected, listingFeeOutstanding, listingFeeTotal float64
	var sellerUpgradeRevenue, buyerUpgradeRevenue, totalUpgradeRevenue, totalPlatformProfit float64
	var referralCommissionRate string
	if err := row.Scan(
		&totalTx, &grossVolume, &brokerFeeRevenue, &brokerFeePct,
		&listingFeeCollected, &listingFeeOutstanding, &listingFeeTotal,
		&sellerUpgradeRevenue, &buyerUpgradeRevenue, &totalUpgradeRevenue,
		&totalPlatformProfit, &referralCommissionRate,
	); err != nil {
		// View may return no rows if no data yet
		return map[string]interface{}{
			"total_transactions": 0, "gross_volume": 0,
			"broker_fee_revenue": 0, "broker_fee_pct": 2.0,
			"listing_fee_collected": 0, "listing_fee_outstanding": 0, "listing_fee_total": 0,
			"seller_upgrade_revenue": 0, "buyer_upgrade_revenue": 0,
			"total_upgrade_revenue": 0, "total_platform_profit": 0,
			"referral_commission_rate": "1.2",
		}, nil
	}
	return map[string]interface{}{
		"total_transactions":       totalTx,
		"gross_volume":             grossVolume,
		"broker_fee_revenue":       brokerFeeRevenue,
		"broker_fee_pct":           brokerFeePct,
		"listing_fee_collected":    listingFeeCollected,
		"listing_fee_outstanding":  listingFeeOutstanding,
		"listing_fee_total":        listingFeeTotal,
		"seller_upgrade_revenue":   sellerUpgradeRevenue,
		"buyer_upgrade_revenue":    buyerUpgradeRevenue,
		"total_upgrade_revenue":    totalUpgradeRevenue,
		"total_platform_profit":    totalPlatformProfit,
		"referral_commission_rate": referralCommissionRate,
	}, nil
}

// ProfitByTier returns profit breakdown per seller tier.
func (r *AdminRepo) ProfitByTier(ctx context.Context) ([]map[string]interface{}, error) {
	rows, err := r.db.Query(ctx, `SELECT * FROM profit_by_seller_tier ORDER BY total_profit DESC`) //nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []map[string]interface{}
	for rows.Next() {
		var tier string
		var products int
		var volume, brokerFees, listingFees, upgradeFees, totalProfit float64
		if err := rows.Scan(&tier, &products, &volume, &brokerFees, &listingFees, &upgradeFees, &totalProfit); err != nil {
			return nil, err
		}
		out = append(out, map[string]interface{}{
			"tier": tier, "products": products, "volume": volume,
			"broker_fees": brokerFees, "listing_fees": listingFees,
			"upgrade_fees": upgradeFees, "total_profit": totalProfit,
		})
	}
	return out, rows.Err()
}

// MonthlyProfitTrend returns monthly profit trend.
func (r *AdminRepo) MonthlyProfitTrend(ctx context.Context) ([]map[string]interface{}, error) {
	rows, err := r.db.Query(ctx, `SELECT month, transactions, volume, broker_fees, listing_fees, total_profit, avg_fee_pct FROM monthly_profit_trend LIMIT 12`) //nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []map[string]interface{}
	for rows.Next() {
		var month string
		var txCount int
		var volume, brokerFees, listingFees, totalProfit, avgFeePct float64
		if err := rows.Scan(&month, &txCount, &volume, &brokerFees, &listingFees, &totalProfit, &avgFeePct); err != nil {
			return nil, err
		}
		out = append(out, map[string]interface{}{
			"month": month, "transactions": txCount,
			"volume": volume, "broker_fees": brokerFees,
			"listing_fees": listingFees, "total_profit": totalProfit,
			"avg_fee_pct": avgFeePct,
		})
	}
	return out, rows.Err()
}

// WeeklyProfitTrend returns the last 13 weeks of profit data.
func (r *AdminRepo) WeeklyProfitTrend(ctx context.Context) ([]map[string]interface{}, error) {
	rows, err := //nolint:gosec // G201/G202: parameterized query, no user input in SQL structure
		r.db.Query(ctx, `
		SELECT week, volume, broker_fees, listing_fees, total_profit, transactions, avg_fee_pct
		FROM weekly_profit_trend ORDER BY week_start ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []map[string]interface{}
	for rows.Next() {
		var week string
		var txCount int
		var volume, brokerFees, listingFees, totalProfit, avgFeePct float64
		if err := rows.Scan(&week, &volume, &brokerFees, &listingFees, &totalProfit, &txCount, &avgFeePct); err != nil {
			return nil, err
		}
		out = append(out, map[string]interface{}{
			"week": week, "transactions": txCount,
			"volume": volume, "broker_fees": brokerFees,
			"listing_fees": listingFees, "total_profit": totalProfit,
			"avg_fee_pct": avgFeePct,
		})
	}
	if out == nil {
		out = []map[string]interface{}{}
	}
	return out, rows.Err()
}

// AnnualProfitTrend returns year-over-year profit data.
func (r *AdminRepo) AnnualProfitTrend(ctx context.Context) ([]map[string]interface{}, error) {
	rows, err := //nolint:gosec // G201/G202: parameterized query, no user input in SQL structure
		r.db.Query(ctx, `
		SELECT year, transactions, volume, broker_fees, listing_fees, total_profit, avg_fee_pct
		FROM annual_profit_trend ORDER BY year ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []map[string]interface{}
	for rows.Next() {
		var year, txCount int
		var volume, brokerFees, listingFees, totalProfit, avgFeePct float64
		if err := rows.Scan(&year, &txCount, &volume, &brokerFees, &listingFees, &totalProfit, &avgFeePct); err != nil {
			return nil, err
		}
		out = append(out, map[string]interface{}{
			"year": year, "transactions": txCount, "volume": volume,
			"broker_fees": brokerFees, "listing_fees": listingFees,
			"total_profit": totalProfit, "avg_fee_pct": avgFeePct,
		})
	}
	if out == nil {
		out = []map[string]interface{}{}
	}
	return out, rows.Err()
}

// ── Sub-Admin Management (Super Admin Only) ───────────────────────────────────

func (r *AdminRepo) ListPendingFeePayments(ctx context.Context) ([]map[string]interface{}, error) {
	rows, err := //nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
		r.db.Query(ctx, `
		SELECT p.id, p.title, p.price, p.total_fee_amount, p.status,
		       COALESCE(p.admin_notes,'') AS payment_claim,
		       u.id AS seller_id, u.name AS seller_name, u.phone AS seller_phone,
		       p.created_at, p.updated_at
		FROM products p
		JOIN users u ON u.id = p.seller_id
		WHERE p.fee_payment_pending = TRUE AND p.fee_paid = FALSE
		ORDER BY p.updated_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []map[string]interface{}
	for rows.Next() {
		var productID, sellerID int
		var title, paymentClaim, status, sellerName, sellerPhone string
		var price, feeAmount float64
		var createdAt, updatedAt interface{}
		if err := rows.Scan(&productID, &title, &price, &feeAmount, &status,
			&paymentClaim, &sellerID, &sellerName, &sellerPhone,
			&createdAt, &updatedAt); err != nil {
			return nil, err
		}
		out = append(out, map[string]interface{}{
			"product_id": productID, "title": title, "price": price,
			"fee_amount": feeAmount, "status": status,
			"payment_claim": paymentClaim,
			"seller_id":     sellerID, "seller_name": sellerName, "seller_phone": sellerPhone,
			"created_at": createdAt, "updated_at": updatedAt,
		})
	}
	if out == nil {
		out = []map[string]interface{}{}
	}
	return out, rows.Err()
}

// VerifyListingFeePayment marks a seller's payment claim as verified.
// Sets fee_paid=TRUE and fee_payment_pending=FALSE.
// Only then does the payment count as collected revenue in profit reports.
func (r *AdminRepo) VerifyListingFeePayment(ctx context.Context, productID, adminID int) error {
	tag, err := r.db.Exec(ctx, `
		UPDATE products
		SET fee_paid = TRUE,
		    fee_paid_at = NOW(),
		    fee_payment_pending = FALSE,
		    reviewed_by = $2,
		    updated_at = NOW()
		WHERE id = $1 AND fee_payment_pending = TRUE AND fee_paid = FALSE`,
		productID, adminID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("product not found or payment already verified/rejected")
	}
	return nil
}

// RejectListingFeePayment rejects a seller's payment claim.
// Clears fee_payment_pending so the seller can resubmit with correct details.
func (r *AdminRepo) RejectListingFeePayment(ctx context.Context, productID, adminID int, reason string) error {
	note := "[FEE PAYMENT REJECTED]"
	if reason != "" {
		note += " " + reason
	}
	tag, err := r.db.Exec(ctx, `
		UPDATE products
		SET fee_payment_pending = FALSE,
		    admin_notes = $3,
		    reviewed_by = $2,
		    updated_at = NOW()
		WHERE id = $1 AND fee_payment_pending = TRUE AND fee_paid = FALSE`,
		productID, adminID, note)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("product not found or no pending payment claim")
	}
	return nil
}

