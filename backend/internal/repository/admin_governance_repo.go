package repository

import (
	"context"
	"fmt"
	"github.com/your-org/used-tools-marketplace/internal/models"
	"strconv"
	"strings"
	"time"
)

func (r *AdminRepo) ListAdmins(ctx context.Context) ([]models.User, error) {
	rows, err := //nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
		r.db.Query(ctx, `
		SELECT id, name, phone, COALESCE(email,''), COALESCE(admin_role::text,'normal') AS admin_role,
		       is_active, created_at, updated_at
		FROM users WHERE role='admin'
		ORDER BY created_at ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Phone, &u.Email, &u.AdminRole,
			&u.IsActive, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		u.Role = "admin"
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	// Load permissions for each normal admin (filtered to valid keys only)
	validPerms := map[string]bool{}
	for _, p := range models.AllAdminPermissions {
		validPerms[p] = true
	}
	for i, u := range users {
		if u.AdminRole != "platform_owner" && u.AdminRole != "super" {
			raw, _ := r.GetAdminPermissions(ctx, u.ID)
			filtered := make([]string, 0, len(raw))
			for _, p := range raw {
				if validPerms[p] {
					filtered = append(filtered, p)
				}
			}
			users[i].Permissions = filtered
		}
	}
	return users, nil
}

func (r *AdminRepo) CreateAdmin(ctx context.Context, name, phone, email, passwordHash, adminRole string) (int, error) {
	var id int
	var emailVal interface{}
	if email != "" {
		emailVal = email
	}
	err := //nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
		r.db.QueryRow(ctx, `
		INSERT INTO users (name, phone, email, password_hash, role, admin_role, is_active, referral_code)
		VALUES ($1, $2, $3, $4, 'admin', $5::admin_role, true,
		  'REF' || translate(upper(substring(md5(random()::text || $6::text), 1, 8)), 'IO01', 'ABCD')
		)
		RETURNING id`,
			name, phone, emailVal, passwordHash, adminRole, phone).Scan(&id)
	return id, err
}

func (r *AdminRepo) UpdateAdminRole(ctx context.Context, adminID int, adminRole string) error {
	tag, err := r.db.Exec(ctx, `
		UPDATE users SET admin_role=$1::admin_role, updated_at=NOW()
		WHERE id=$2 AND role='admin'`,
		adminRole, adminID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("admin with id %d not found", adminID)
	}
	return nil
}

func (r *AdminRepo) DeleteAdmin(ctx context.Context, adminID int) error {
	_, err := r.db.Exec(ctx, `DELETE FROM users WHERE id=$1 AND role='admin'`, adminID) //nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
	return err
}

func (r *AdminRepo) CountSuperAdmins(ctx context.Context) (int, error) {
	var count int
	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM users WHERE role='admin' AND (admin_role='super' OR admin_role='platform_owner')`).Scan(&count) //nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
	return count, err
}

func (r *AdminRepo) IsAdminSuper(ctx context.Context, adminID int) (bool, error) {
	var adminRole string
	err := r.db.QueryRow(ctx, `SELECT COALESCE(admin_role::text,'') FROM users WHERE id=$1 AND role='admin'`, adminID).Scan(&adminRole) //nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
	if err != nil {
		// No row = not an admin at all
		return false, nil
	}
	return adminRole == "super" || adminRole == "platform_owner", nil
}

// ── Granular Permissions ──────────────────────────────────────────────────────

func (r *AdminRepo) GetAdminPermissions(ctx context.Context, adminID int) ([]string, error) {
	rows, err := //nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
		r.db.Query(ctx, `
		SELECT permission FROM admin_permissions
		WHERE admin_id=$1 AND enabled=TRUE
		ORDER BY permission`, adminID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var perms []string
	for rows.Next() {
		var p string
		if err := rows.Scan(&p); err != nil {
			return nil, err
		}
		perms = append(perms, p)
	}
	return perms, rows.Err()
}

// SetAdminPermissions replaces all permissions for an admin with the given enabled list.
func (r *AdminRepo) SetAdminPermissions(ctx context.Context, adminID int, permissions []string) error {
	// Build a set for quick lookup
	enabled := map[string]bool{}
	for _, p := range permissions {
		enabled[p] = true
	}

	// Upsert all known permissions
	for _, perm := range models.AllAdminPermissions {
		_, err := r.db.Exec(ctx, `
			INSERT INTO admin_permissions (admin_id, permission, enabled, updated_at)
			VALUES ($1, $2, $3, NOW())
			ON CONFLICT (admin_id, permission) DO UPDATE SET enabled=$3, updated_at=NOW()`,
			adminID, perm, enabled[perm])
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *AdminRepo) HasPermission(ctx context.Context, adminID int, permission string) (bool, error) {
	var enabled bool
	err := //nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
		r.db.QueryRow(ctx, `
		SELECT COALESCE(enabled, FALSE) FROM admin_permissions
		WHERE admin_id=$1 AND permission=$2`, adminID, permission).Scan(&enabled)
	if err != nil {
		// No row = no permission
		return false, nil
	}
	return enabled, nil
}

// ── Admin Activity Logs ───────────────────────────────────────────────────────

func (r *AdminRepo) InsertAdminLog(ctx context.Context, l models.AdminLog) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO admin_logs (admin_id, product_id, target_user_id, action, notes, details, ip_address, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())`,
		l.AdminID, l.ProductID, l.TargetUserID, l.Action, l.Notes, l.Details, l.IPAddress)
	return err
}

// validAdminLogActions is the allowlist of action values accepted in log filters.
var validAdminLogActions = map[string]bool{
	"approve": true, "reject": true, "delete": true, "activate": true, "deactivate": true,
	"tier_change": true, "warn": true, "ban": true, "verify_fee": true, "reject_fee": true,
	"repair_missing_email": true, "create": true, "update": true,
}

func (r *AdminRepo) ListAdminLogs(ctx context.Context, f models.AdminLogFilter) ([]models.AdminLog, int, error) {
	// Allowlist action to prevent arbitrary values reaching the SQL builder
	if f.Action != "" && !validAdminLogActions[f.Action] {
		f.Action = ""
	}
	// Validate date format to prevent malformed strings reaching the DB
	const dateLayout = "2006-01-02"
	if f.DateFrom != "" {
		if _, err := time.Parse(dateLayout, f.DateFrom); err != nil {
			return nil, 0, fmt.Errorf("invalid date_from format — use YYYY-MM-DD")
		}
	}
	if f.DateTo != "" {
		if _, err := time.Parse(dateLayout, f.DateTo); err != nil {
			return nil, 0, fmt.Errorf("invalid date_to format — use YYYY-MM-DD")
		}
	}

	var wb strings.Builder
	wb.WriteString("1=1")
	args := []interface{}{}

	if f.AdminID > 0 {
		args = appendCond(&wb, args, " AND al.admin_id=", f.AdminID)
	}
	if f.Action != "" {
		args = appendCond(&wb, args, " AND al.action=", f.Action)
	}
	if f.IPAddress != "" {
		args = appendCond(&wb, args, " AND al.ip_address ILIKE ", likeWrap(f.IPAddress))
	}
	if f.DateFrom != "" {
		args = appendCond(&wb, args, " AND al.created_at>=", f.DateFrom)
	}
	if f.DateTo != "" {
		args = appendCond(&wb, args, " AND al.created_at<=", f.DateTo+" 23:59:59")
	}

	whereStr := wb.String()

	// Count
	var total int
	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM admin_logs al WHERE `+whereStr, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	limit := f.Limit
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	page := f.Page
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	i := len(args) + 1
	args = append(args, limit, offset)
	rows, err := r.db.Query(ctx,
		`SELECT al.id, al.admin_id, COALESCE(u.name,'') AS admin_name,
		       al.product_id, al.target_user_id, COALESCE(tu.name,'') AS target_name,
		       al.action, COALESCE(al.notes,''), COALESCE(al.details,''), COALESCE(al.ip_address,''),
		       al.created_at
		FROM admin_logs al
		LEFT JOIN users u  ON u.id = al.admin_id
		LEFT JOIN users tu ON tu.id = al.target_user_id
		WHERE `+whereStr+`
		ORDER BY al.created_at DESC
		LIMIT $`+strconv.Itoa(i)+` OFFSET $`+strconv.Itoa(i+1),
		args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var out []models.AdminLog
	for rows.Next() {
		var l models.AdminLog
		if err := rows.Scan(&l.ID, &l.AdminID, &l.AdminName, &l.ProductID, &l.TargetUserID,
			&l.TargetName, &l.Action, &l.Notes, &l.Details, &l.IPAddress, &l.CreatedAt); err != nil {
			return nil, 0, err
		}
		out = append(out, l)
	}
	return out, total, rows.Err()
}

// ── Tier Management ───────────────────────────────────────────────────────────

func (r *AdminRepo) SetSellerTier(ctx context.Context, userID int, tier string, changedBy int, reason string) error {
	var oldTier string
	err := r.db.QueryRow(ctx, `SELECT tier FROM users WHERE id=$1`, userID).Scan(&oldTier) //nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}
	_, err = r.db.Exec(ctx, `UPDATE users SET tier=$1::seller_tier, updated_at=NOW() WHERE id=$2`, tier, userID) //nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
	if err != nil {
		return err
	}
	// The DB trigger trg_update_products_on_seller_tier handles cascading to products.
	_, _ = r.db.Exec(ctx, `
		INSERT INTO user_tier_history (user_id, tier_type, old_tier, new_tier, changed_by, reason)
		VALUES ($1, 'seller', $2, $3, $4, $5)`,
		userID, oldTier, tier, changedBy, reason)
	return nil
}

func (r *AdminRepo) SetBuyerTier(ctx context.Context, userID int, tier string, changedBy int, reason string) error {
	var oldTier string
	err := r.db.QueryRow(ctx, `SELECT buyer_tier FROM users WHERE id=$1`, userID).Scan(&oldTier) //nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}
	_, err = r.db.Exec(ctx, `UPDATE users SET buyer_tier=$1::buyer_tier, updated_at=NOW() WHERE id=$2`, tier, userID) //nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
	if err != nil {
		return err
	}
	_, _ = r.db.Exec(ctx, `
		INSERT INTO user_tier_history (user_id, tier_type, old_tier, new_tier, changed_by, reason)
		VALUES ($1, 'buyer', $2, $3, $4, $5)`,
		userID, oldTier, tier, changedBy, reason)
	return nil
}

func (r *AdminRepo) ListTierHistory(ctx context.Context, userID int) ([]models.TierHistory, error) {
	rows, err := //nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
		r.db.Query(ctx, `
		SELECT h.id, h.user_id, COALESCE(u.name,'') AS user_name,
		       h.tier_type, h.old_tier, h.new_tier,
		       h.changed_by, COALESCE(a.name,'') AS changed_by_name,
		       COALESCE(h.reason,''), h.created_at
		FROM user_tier_history h
		LEFT JOIN users u ON u.id = h.user_id
		LEFT JOIN users a ON a.id = h.changed_by
		WHERE h.user_id=$1
		ORDER BY h.created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.TierHistory
	for rows.Next() {
		var h models.TierHistory
		if err := rows.Scan(&h.ID, &h.UserID, &h.UserName, &h.TierType, &h.OldTier, &h.NewTier,
			&h.ChangedBy, &h.ChangedByName, &h.Reason, &h.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, h)
	}
	return out, rows.Err()
}

func (r *AdminRepo) ListPremiumUsers(ctx context.Context, role string, tierType string, tier string) ([]models.User, error) {
	var rows interface {
		Next() bool
		Scan(...interface{}) error
		Close()
		Err() error
	}
	var err error

	if tierType == "buyer" {
		q := `
			SELECT id, name, phone, COALESCE(business_name,''), COALESCE(email,''), rating, role,
			       COALESCE(admin_role::text,''), is_active, tier, total_value_sold, buyer_tier, otp_verified, created_at, updated_at
			FROM users
			WHERE role='buyer'`
		args := []interface{}{}
		if tier != "" && tier != "all" {
			q += ` AND buyer_tier=$1`
			args = append(args, tier)
		}
		q += ` ORDER BY created_at DESC`
		rows, err = r.db.Query(ctx, q, args...) //nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
	} else {
		q := `
			SELECT id, name, phone, COALESCE(business_name,''), COALESCE(email,''), rating, role,
			       COALESCE(admin_role::text,''), is_active, tier, total_value_sold, buyer_tier, otp_verified, created_at, updated_at
			FROM users
			WHERE role='seller'`
		args := []interface{}{}
		if tier != "" && tier != "all" {
			q += ` AND tier=$1`
			args = append(args, tier)
		}
		q += ` ORDER BY total_value_sold DESC`
		rows, err = r.db.Query(ctx, q, args...) //nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Phone, &u.BusinessName, &u.Email, &u.Rating,
			&u.Role, &u.AdminRole, &u.IsActive, &u.Tier, &u.TotalValueSold, &u.BuyerTier,
			&u.OTPVerified, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, u)
	}
	return out, rows.Err()
}

// DeleteAdminLogs removes admin activity log records to free storage.
// before="" → deletes all logs; before="YYYY-MM-DD" → deletes logs older than that date.
func (r *AdminRepo) DeleteAdminLogs(ctx context.Context, before string) (int64, error) {
	var tag interface{ RowsAffected() int64 }
	var err error
	if before == "" {
		tag, err = r.db.Exec(ctx, `DELETE FROM admin_logs`) //nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
	} else {
		tag, err = r.db.Exec(ctx, `DELETE FROM admin_logs WHERE created_at < $1::date`, before) //nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
	}
	if err != nil {
		return 0, err
	}
	return tag.RowsAffected(), nil
}

// ── Seller Upgrade Requests (Admin) ──────────────────────────────────────────

// ListUpgradeRequests returns all seller upgrade requests, optionally filtered by status.
// ── Upgrade Requests (unified seller + buyer) ─────────────────────────────────

// ListUpgradeRequests returns upgrade requests filtered by role and/or status.
// role: "seller" | "buyer" | "" (all)
// status: "pending" | "approved" | "rejected" | "" (all)
func (r *AdminRepo) ListUpgradeRequests(ctx context.Context, status string) ([]map[string]interface{}, error) {
	return r.listUpgradeRequests(ctx, "", status, 1, 100)
}

func (r *AdminRepo) ListBuyerUpgradeRequests(ctx context.Context, status string, page, limit int) ([]map[string]interface{}, int, error) {
	return r.listUpgradeRequestsPaged(ctx, "buyer", status, page, limit)
}

// listUpgradeRequests is the internal helper used by both seller and buyer list methods.
func (r *AdminRepo) listUpgradeRequests(ctx context.Context, role, status string, page, limit int) ([]map[string]interface{}, error) {
	results, _, err := r.listUpgradeRequestsPaged(ctx, role, status, page, limit)
	return results, err
}

// validUpgradeRole and validUpgradeStatus are allowlists that prevent arbitrary
// values from reaching the SQL query builder.
var validUpgradeRole = map[string]bool{"seller": true, "buyer": true}
var validUpgradeStatus = map[string]bool{"pending": true, "approved": true, "rejected": true}

func (r *AdminRepo) listUpgradeRequestsPaged(ctx context.Context, role, status string, page, limit int) ([]map[string]interface{}, int, error) {
	if limit <= 0 || limit > 200 {
		limit = 100
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	// Allowlist role and status — reject unknown values before they reach SQL
	if role != "" && !validUpgradeRole[role] {
		return nil, 0, fmt.Errorf("invalid role filter: must be seller or buyer")
	}
	if status != "" && status != "all" && !validUpgradeStatus[status] {
		return nil, 0, fmt.Errorf("invalid status filter: must be pending, approved, or rejected")
	}

	var wb strings.Builder
	wb.WriteString("1=1")
	args := []interface{}{}
	if role != "" {
		args = appendCond(&wb, args, " AND r.role=", role)
	}
	if status != "" && status != "all" {
		args = appendCond(&wb, args, " AND r.status=", status)
	}
	whereStr := wb.String()

	var total int
	_ = r.db.QueryRow(ctx, `SELECT COUNT(*) FROM upgrade_requests r WHERE `+whereStr, args...).Scan(&total)

	i := len(args) + 1
	args = append(args, limit, offset)
	rows, err := r.db.Query(ctx,
		`SELECT r.id, r.user_id, r.role, u.name, u.phone,
		       CASE r.role
		         WHEN 'seller' THEN COALESCE(u.tier::text,'normal')
		         ELSE COALESCE(u.buyer_tier::text,'normal')
		       END AS current_tier,
		       r.requested_tier, r.payment_method, r.payer_phone,
		       r.reference, r.amount, COALESCE(r.notes,''),
		       r.status, COALESCE(r.reject_reason,''), r.reviewed_by, r.reviewed_at, r.created_at
		FROM upgrade_requests r
		JOIN users u ON u.id = r.user_id
		WHERE `+whereStr+` ORDER BY r.created_at DESC LIMIT $`+strconv.Itoa(i)+` OFFSET $`+strconv.Itoa(i+1),
		args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var out []map[string]interface{}
	for rows.Next() {
		var id, userID int
		var role, name, phone, currentTier, requestedTier, method, payerPhone, ref, notes, st, rejectReason string
		var amount float64
		var reviewedBy *int
		var reviewedAt, createdAt *time.Time
		if err := rows.Scan(&id, &userID, &role, &name, &phone, &currentTier,
			&requestedTier, &method, &payerPhone, &ref, &amount, &notes,
			&st, &rejectReason, &reviewedBy, &reviewedAt, &createdAt); err != nil {
			continue
		}
		out = append(out, map[string]interface{}{
			"id": id, "user_id": userID, "role": role,
			"user_name": name, "user_phone": phone, "current_tier": currentTier,
			"requested_tier": requestedTier, "payment_method": method,
			"payer_phone": payerPhone, "reference": ref, "amount": amount,
			"notes": notes, "status": st, "reject_reason": rejectReason,
			"reviewed_by": reviewedBy, "reviewed_at": reviewedAt, "created_at": createdAt,
		})
	}
	if out == nil {
		out = []map[string]interface{}{}
	}
	return out, total, rows.Err()
}

// ApproveUpgradeRequest approves a seller or buyer upgrade request and applies the tier.
func (r *AdminRepo) ApproveUpgradeRequest(ctx context.Context, requestID, adminID int) error {
	return r.approveUpgradeRequest(ctx, requestID, adminID)
}

func (r *AdminRepo) ApproveBuyerUpgradeRequest(ctx context.Context, requestID, adminID int) error {
	return r.approveUpgradeRequest(ctx, requestID, adminID)
}

func (r *AdminRepo) approveUpgradeRequest(ctx context.Context, requestID, adminID int) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	var userID int
	var role, requestedTier, status string
	var amount float64
	var reference, payerPhone, paymentMethod string
	err = tx.QueryRow(ctx, `
		SELECT user_id, role, requested_tier, status, amount, reference, payer_phone, payment_method
		FROM upgrade_requests WHERE id=$1`,
		requestID).Scan(&userID, &role, &requestedTier, &status, &amount, &reference, &payerPhone, &paymentMethod)
	if err != nil {
		return fmt.Errorf("upgrade request not found")
	}
	if status != "pending" {
		return fmt.Errorf("request is already %s", status)
	}

	// Apply tier to the correct column based on role
	if role == "seller" {
		_, err = tx.Exec(ctx, `UPDATE users SET tier=$1::seller_tier, updated_at=NOW() WHERE id=$2`, requestedTier, userID) //nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
	} else {
		_, err = tx.Exec(ctx, `UPDATE users SET buyer_tier=$1::buyer_tier, buyer_tier_paid_at=NOW(), updated_at=NOW() WHERE id=$2`, requestedTier, userID) //nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
		if err == nil {
			_, err = tx.Exec(ctx, `
				INSERT INTO buyer_tier_payments (buyer_id, tier, amount, reference, payer_phone, payment_method)
				VALUES ($1, $2::buyer_tier, $3, $4, $5, $6)`,
				userID, requestedTier, amount, reference, payerPhone, paymentMethod)
		}
	}
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `
		UPDATE upgrade_requests SET status='approved', reviewed_by=$1, reviewed_at=NOW(), updated_at=NOW()
		WHERE id=$2`, adminID, requestID)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// RejectUpgradeRequest rejects a seller or buyer upgrade request.
func (r *AdminRepo) RejectUpgradeRequest(ctx context.Context, requestID, adminID int) error {
	return r.rejectUpgradeRequest(ctx, requestID, adminID, "")
}

func (r *AdminRepo) RejectBuyerUpgradeRequest(ctx context.Context, requestID, adminID int, reason string) error {
	return r.rejectUpgradeRequest(ctx, requestID, adminID, reason)
}

func (r *AdminRepo) rejectUpgradeRequest(ctx context.Context, requestID, adminID int, reason string) error {
	tag, err := r.db.Exec(ctx, `
		UPDATE upgrade_requests
		SET status='rejected', reviewed_by=$1, reviewed_at=NOW(),
		    reject_reason=NULLIF($2,''), updated_at=NOW()
		WHERE id=$3 AND status='pending'`,
		adminID, reason, requestID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("request not found or already reviewed")
	}
	return nil
}

// (unified upgrade request methods are above)

// GetUpgradeRequestDetails returns (userID, role, requestedTier, amount, error) for commission recording.
func (r *AdminRepo) GetUpgradeRequestDetails(ctx context.Context, requestID int) (int, string, string, float64, error) {
	var userID int
	var role, requestedTier string
	var amount float64
	err := r.db.QueryRow(ctx, //nolint:gosec // G201/G202: parameterized query, no user input in SQL structure
		`SELECT user_id, role, requested_tier, amount FROM upgrade_requests WHERE id=$1`, requestID,
	).Scan(&userID, &role, &requestedTier, &amount)
	return userID, role, requestedTier, amount, err
}

// ── Listing Fee Payment Verification ─────────────────────────────────────────

// ListPendingFeePayments returns products where the seller has submitted a payment
// claim (fee_payment_pending=TRUE) but admin has not yet verified it (fee_paid=FALSE).
func (r *AdminRepo) IssueWarning(ctx context.Context, userID, adminID int, reason string) error {
	// 1. Insert warning record
	//nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
	_, err := r.db.Exec(ctx, `INSERT INTO user_warnings (user_id, admin_id, reason) VALUES ($1, $2, $3)`, userID, adminID, reason)
	return err
}

func (r *AdminRepo) SuspendUser(ctx context.Context, userID, adminID int, reason string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// 1. Mark user as inactive, explicitly NOT banned, set referral_suspended and deactivation_reason
	//nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
	tag, err := tx.Exec(ctx,
		`UPDATE users 
		 SET is_active=FALSE, is_banned=FALSE, referral_suspended=TRUE, deactivation_reason=$1, updated_at=NOW() 
		 WHERE id=$2`,
		reason, userID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("user with id %d not found", userID)
	}

	// 2. Suspend capabilities in user_capabilities
	//nolint:gosec
	_, err = tx.Exec(ctx,
		`UPDATE user_capabilities 
		 SET status='suspended', 
		     metadata=jsonb_set(COALESCE(metadata, '{}'::jsonb), '{suspend_reason}', to_jsonb($2::text)) 
		 WHERE user_id=$1`,
		userID, reason)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *AdminRepo) BanUser(ctx context.Context, userID, adminID int, reason string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// 1. Mark as permanently banned, inactive, referral permanently blocked
	//nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
	tag, err := tx.Exec(ctx,
		`UPDATE users 
		 SET is_banned=TRUE, is_active=FALSE, referral_suspended=TRUE, referral_permanently_blocked=TRUE, referral_block_reason=$1, deactivation_reason=$1, updated_at=NOW() 
		 WHERE id=$2`,
		reason, userID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("user with id %d not found", userID)
	}

	// 2. Revoke user capabilities
	//nolint:gosec
	_, err = tx.Exec(ctx,
		`UPDATE user_capabilities 
		 SET status='revoked', 
		     metadata=jsonb_set(COALESCE(metadata, '{}'::jsonb), '{ban_reason}', to_jsonb($2::text)) 
		 WHERE user_id=$1`,
		userID, reason)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *AdminRepo) CheckSystemHealth(ctx context.Context) (map[string]interface{}, error) {
	report := make(map[string]interface{})

	// 1. Sub-admins with no permissions
	var noPermAdmins []map[string]interface{}
	rows, err := r.db.Query(ctx, `
		SELECT id, name, email 
		FROM users 
		WHERE role='admin' AND admin_role NOT IN ('super', 'platform_owner') 
		AND id NOT IN (SELECT admin_id FROM admin_permissions)`)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var id int
			var name, email string
			if err := rows.Scan(&id, &name, &email); err == nil {
				noPermAdmins = append(noPermAdmins, map[string]interface{}{"id": id, "name": name, "email": email})
			}
		}
	}
	report["no_permission_admins"] = noPermAdmins

	// 2. Duplicate phones across same roles
	var dupPhones []map[string]interface{}
	rows, err = r.db.Query(ctx, `
		SELECT phone, role, COUNT(*) 
		FROM users 
		GROUP BY phone, role 
		HAVING COUNT(*) > 1`)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var phone, role string
			var count int
			if err := rows.Scan(&phone, &role, &count); err == nil {
				dupPhones = append(dupPhones, map[string]interface{}{"phone": phone, "role": role, "count": count})
			}
		}
	}
	report["duplicate_phones"] = dupPhones

	// 3. Duplicate emails
	var dupEmails []map[string]interface{}
	rows, err = r.db.Query(ctx, `
		SELECT email, COUNT(*) 
		FROM users 
		WHERE email IS NOT NULL 
		GROUP BY email 
		HAVING COUNT(*) > 1`)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var email string
			var count int
			if err := rows.Scan(&email, &count); err == nil {
				dupEmails = append(dupEmails, map[string]interface{}{"email": email, "count": count})
			}
		}
	}
	report["duplicate_emails"] = dupEmails

	// 4. Admins missing emails (security risk)
	var adminsNoEmail []map[string]interface{}
	rows, err = r.db.Query(ctx, "SELECT id, name, phone FROM users WHERE role='admin' AND (email IS NULL OR email = '')")
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var id int
			var name, phone string
			if err := rows.Scan(&id, &name, &phone); err == nil {
				adminsNoEmail = append(adminsNoEmail, map[string]interface{}{"id": id, "name": name, "phone": phone})
			}
		}
	}
	report["admins_missing_email"] = adminsNoEmail

	// 5. Duplicate referral codes
	var dupRefCodes []map[string]interface{}
	rows, err = r.db.Query(ctx, "SELECT referral_code, COUNT(*) FROM users GROUP BY referral_code HAVING COUNT(*) > 1")
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var code string
			var count int
			if err := rows.Scan(&code, &count); err == nil {
				dupRefCodes = append(dupRefCodes, map[string]interface{}{"referral_code": code, "count": count})
			}
		}
	}
	report["duplicate_referral_codes"] = dupRefCodes

	// 5. Circular Referrals (A -> B -> A)
	var circular []map[string]interface{}
	cRows, err := //nolint:gosec // G201/G202: parameterized query, no user input in SQL structure
		r.db.Query(ctx, `
		WITH RECURSIVE referral_chains AS (
			-- Base case: everyone who has an inviter
			SELECT id, invited_by, ARRAY[id] as path, false as is_cycle
			FROM users
			WHERE invited_by IS NOT NULL
			UNION ALL
			-- Recursive step: walk up the chain
			SELECT u.id, u.invited_by, rc.path || u.id, u.id = ANY(rc.path)
			FROM users u
			INNER JOIN referral_chains rc ON rc.invited_by = u.id
			WHERE u.invited_by IS NOT NULL AND NOT rc.is_cycle
		)
		SELECT id, path FROM referral_chains WHERE is_cycle = true`)
	if err == nil {
		defer cRows.Close()
		for cRows.Next() {
			var id int
			var path []int
			if err := cRows.Scan(&id, &path); err == nil {
				circular = append(circular, map[string]interface{}{
					"user_id": id,
					"path":    path,
				})
			}
		}
	}
	report["circular_referrals"] = circular

	return report, nil
}
func (r *AdminRepo) RepairMissingAdminEmails(ctx context.Context, adminID int, emailMap map[int]string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for id, email := range emailMap {
		_, err := tx.Exec(ctx, `UPDATE users SET email = $1, updated_at = NOW() WHERE id = $2 AND role = 'admin'`, email, id)
		if err != nil {
			return err
		}
		_, _ = tx.Exec(ctx,
			`INSERT INTO admin_logs (admin_id, target_user_id, action, notes, created_at)
			VALUES ($1, $2, 'repair_missing_email', $3, NOW())`,
			adminID, id, "Repaired missing admin email to: "+email)
	}

	return tx.Commit(ctx)
}

// ── Operations Monitoring (Role 02) ───────────────────────────────────────────

// GetOperationsHealth returns system health metrics
