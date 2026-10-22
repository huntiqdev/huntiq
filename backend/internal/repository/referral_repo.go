package repository

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/your-org/used-tools-marketplace/pkg/constants"
)

type ReferralRepo struct{ db *pgxpool.Pool }

func NewReferralRepo(db *pgxpool.Pool) *ReferralRepo { return &ReferralRepo{db} }

// GetStats returns referral dashboard data for a user.
func (r *ReferralRepo) GetStats(ctx context.Context, userID int) (*ReferralStats, error) {
	stats := &ReferralStats{}

	//nolint:gosec // G201/G202: parameterized query, no user input in SQL structure
	err := r.db.QueryRow(ctx, `
		SELECT
			u.referral_code,
			u.referral_suspended,
			u.invitation_accepted_at,
			inv.referral_code,
			inv.name
		FROM users u
		LEFT JOIN users inv ON inv.id = u.invited_by
		WHERE u.id = $1`, userID,
	).Scan(
		&stats.ReferralCode,
		&stats.ReferralSuspended,
		&stats.InvitationAcceptedAt,
		&stats.InviterCode,
		&stats.InviterName,
	)
	if err != nil {
		return nil, err
	}

	_ = r.db.QueryRow(ctx, //nolint:gosec // G201/G202: parameterized query, no user input in SQL structure
		`SELECT COUNT(*) FROM users WHERE invited_by = $1`, userID,
	).Scan(&stats.InvitedCount)

	_ = r.db.QueryRow(ctx, //nolint:gosec // G201/G202: parameterized query, no user input in SQL structure
		`SELECT COALESCE(SUM(commission_amount),0) FROM referral_commissions WHERE referrer_id = $1`, userID,
	).Scan(&stats.TotalCommission)

	return stats, nil
}

// GetCommissions returns paginated commission history for an inviter.
func (r *ReferralRepo) GetCommissions(ctx context.Context, inviterID int, page, limit int) ([]CommissionEntry, int, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	var total int
	_ = r.db.QueryRow(ctx, //nolint:gosec // G201/G202: parameterized query, no user input in SQL structure
		`SELECT COUNT(*) FROM referral_commissions WHERE referrer_id = $1`, inviterID,
	).Scan(&total)

	rows, err := r.db.Query(ctx, `
		SELECT rc.id, rc.referred_id, COALESCE(u.name,''), rc.transaction_amount,
		       rc.commission_amount, rc.transaction_type, rc.created_at
		FROM referral_commissions rc
		LEFT JOIN users u ON u.id = rc.referred_id
		WHERE rc.referrer_id = $1
		ORDER BY rc.created_at DESC
		LIMIT $2 OFFSET $3`, inviterID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var out []CommissionEntry
	for rows.Next() {
		var e CommissionEntry
		if err := rows.Scan(&e.ID, &e.InvitedUserID, &e.InvitedUserName,
			&e.TransactionAmount, &e.CommissionAmount, &e.TransactionType, &e.CreatedAt); err != nil {
			return nil, 0, err
		}
		out = append(out, e)
	}
	if out == nil {
		out = []CommissionEntry{}
	}
	return out, total, rows.Err()
}

// RecordCommission inserts a commission row using the platform's configured rate.
// Rate is loaded from platform_settings.referral_commission_rate (default 1.2%).
// Uses ON CONFLICT DO NOTHING to prevent double-recording the same source event.
func (r *ReferralRepo) RecordCommission(ctx context.Context, inviterID, invitedUserID int, txAmount float64, txType, invitedUserRole string, sourceID *int) error {
	// Load live rate from DB — falls back to 1.2% if not set
	rateStr := ""
	_ = r.db.QueryRow(ctx, `SELECT value FROM platform_settings WHERE key='referral_commission_rate'`).Scan(&rateStr) //nolint:gosec // G201/G202: parameterized query, no user input in SQL structure
	rate := 0.012
	if v, err := strconv.ParseFloat(rateStr, 64); err == nil && v >= 0 {
		rate = v / 100 // stored as percentage (e.g. "1.2" → 0.012)
	}
	commission := math.Round(txAmount*rate*100) / 100
	//nolint:gosec // G201/G202: parameterized query, no user input in SQL structure
	_, err := r.db.Exec(ctx, `
		INSERT INTO referral_commissions
		  (referrer_id, referred_id, transaction_amount, commission_amount, transaction_type, transaction_id)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT ON CONSTRAINT uq_referral_commission_source DO NOTHING`,
		inviterID, invitedUserID, txAmount, commission, txType, sourceID)
	return err
}

// RecordPendingCommission inserts a commission with 'pending' status.
// The commission will be confirmed after pendingUntil timestamp passes.
// This prevents withdrawal before payment finalization (chargeback protection).
func (r *ReferralRepo) RecordPendingCommission(
	ctx context.Context,
	inviterID, invitedUserID int,
	txAmount float64,
	txType, invitedUserRole string,
	sourceID *int,
	pendingUntil time.Time,
) error {
	// Load live rate from DB — falls back to 1.2% if not set
	rateStr := ""
	_ = r.db.QueryRow(ctx, `SELECT value FROM platform_settings WHERE key='referral_commission_rate'`).Scan(&rateStr) //nolint:gosec // G201/G202: parameterized query
	rate := 0.012
	if v, err := strconv.ParseFloat(rateStr, 64); err == nil && v >= 0 {
		rate = v / 100 // stored as percentage (e.g. "1.2" → 0.012)
	}
	commission := math.Round(txAmount*rate*100) / 100

	//nolint:gosec // G201/G202: parameterized query, no user input in SQL structure
	_, err := r.db.Exec(ctx, `
		INSERT INTO referral_commissions
		  (referrer_id, referred_id, transaction_amount, commission_amount, 
		   transaction_type, transaction_id, 
		   status, pending_until, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, 'pending', $7, NOW())
		ON CONFLICT ON CONSTRAINT uq_referral_commission_source DO NOTHING`,
		inviterID, invitedUserID, txAmount, commission, txType, sourceID, pendingUntil)
	return err
}

// GetInviterID returns the inviter's user ID for a given user (0 if none).
// Also returns whether the inviter's referral link is suspended or permanently blocked.
func (r *ReferralRepo) GetInviterID(ctx context.Context, userID int) (inviterID int, suspended bool, err error) {
	var id *int
	var susp, blocked bool
	//nolint:gosec // G201/G202: parameterized query, no user input in SQL structure
	err = r.db.QueryRow(ctx, `
		SELECT u.invited_by,
		       COALESCE(inv.referral_suspended, FALSE),
		       COALESCE(inv.referral_permanently_blocked, FALSE)
		FROM users u
		LEFT JOIN users inv ON inv.id = u.invited_by
		WHERE u.id = $1`, userID,
	).Scan(&id, &susp, &blocked)
	if err != nil {
		return 0, false, err
	}
	if id == nil {
		return 0, false, nil
	}
	return *id, susp || blocked, nil // treat blocked same as suspended for commission purposes
}

// AcceptInvitationByCode links a registered user to an inviter via referral code.
func (r *ReferralRepo) AcceptInvitationByCode(ctx context.Context, userID int, referralCode string) error {
	// Check if the user accepting is permanently blocked
	var userBlocked bool
	_ = r.db.QueryRow(ctx, //nolint:gosec // G201/G202: parameterized query, no user input in SQL structure
		`SELECT referral_permanently_blocked FROM users WHERE id = $1`, userID,
	).Scan(&userBlocked)
	if userBlocked {
		return fmt.Errorf("your referral account has been permanently blocked — contact support")
	}

	// Find the inviter — also check suspension, block, and global account status
	var inviterID int
	var suspended, blocked, isActive bool
	err := r.db.QueryRow(ctx, //nolint:gosec // G201/G202: parameterized query, no user input in SQL structure
		`SELECT id, referral_suspended, referral_permanently_blocked, is_active FROM users WHERE referral_code = $1`, referralCode,
	).Scan(&inviterID, &suspended, &blocked, &isActive)
	if err != nil {
		return fmt.Errorf("referral code not found")
	}
	if suspended || blocked || !isActive {
		return fmt.Errorf("this referral code is no longer active (account inactive or suspended)")
	}
	if inviterID == userID {
		return fmt.Errorf("you cannot use your own referral code")
	}

	// ── Deep Scan: Circular Referral Check ─────────────────────────────────────
	// Prevent A refers B, B refers A (and deeper chains A→B→C→A).
	// We check if the current user (userID) is an ancestor of the intended inviter.
	var isCircular bool
	_ = r.db.QueryRow(ctx, `
		WITH RECURSIVE referral_chain AS (
			-- Start with the intended inviter's own inviter
			SELECT invited_by FROM users WHERE id = $1
			UNION ALL
			-- Walk up the chain of inviters
			SELECT u.invited_by FROM users u
			INNER JOIN referral_chain rc ON rc.invited_by = u.id
			WHERE u.invited_by IS NOT NULL
		)
		SELECT EXISTS (SELECT 1 FROM referral_chain WHERE invited_by = $2 OR id = $2)`,
		inviterID, userID).Scan(&isCircular)

	if isCircular {
		return fmt.Errorf("circular referral detected — you cannot accept an invitation from someone you already refer (directly or indirectly)")
	}
	// ───────────────────────────────────────────────────────────────────────────

	// Check if user already has an inviter
	var existingInviter *int
	_ = r.db.QueryRow(ctx, //nolint:gosec // G201/G202: parameterized query, no user input in SQL structure
		`SELECT invited_by FROM users WHERE id = $1`, userID,
	).Scan(&existingInviter)
	if existingInviter != nil {
		return fmt.Errorf("you have already accepted an invitation — this cannot be changed")
	}

	//nolint:gosec // G201/G202: parameterized query, no user input in SQL structure
	tag, err := r.db.Exec(ctx, `
		UPDATE users
		SET invited_by = $1, invitation_accepted_at = NOW(), updated_at = NOW()
		WHERE id = $2 AND invited_by IS NULL`,
		inviterID, userID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("invitation already accepted — cannot change inviter")
	}

	// Audit trail — insert an accepted record. Use INSERT ... ON CONFLICT DO NOTHING
	// so repeated calls (e.g. idempotent retries) don't create duplicate rows.
	// The partial unique index only covers pending rows, so we guard with a manual check.
	//nolint:gosec // G201/G202: parameterized query, no user input in SQL structure
	_, _ = r.db.Exec(ctx, `
		INSERT INTO invitation_requests (user_id, inviter_code_submitted, status)
		VALUES ($1, $2, 'accepted')
		ON CONFLICT (user_id) WHERE status = 'accepted' DO NOTHING`,
		userID, referralCode)

	return nil
}

// SetInvitedByOnRegister links a new user to an inviter at registration time.
func (r *ReferralRepo) SetInvitedByOnRegister(ctx context.Context, userID int, referralCode string) error {
	if referralCode == "" {
		return nil
	}
	var inviterID int
	var suspended, blocked bool
	err := r.db.QueryRow(ctx, //nolint:gosec // G201/G202: parameterized query, no user input in SQL structure
		`SELECT id, referral_suspended, referral_permanently_blocked FROM users WHERE referral_code = $1`, referralCode,
	).Scan(&inviterID, &suspended, &blocked)
	if err != nil || suspended || blocked {
		return nil // invalid, suspended, or blocked code — silently ignore at registration
	}
	if inviterID == userID {
		return nil
	}
	//nolint:gosec // G201/G202: parameterized query, no user input in SQL structure
	_, err = r.db.Exec(ctx, `
		UPDATE users
		SET invited_by = $1, invitation_accepted_at = NOW(), updated_at = NOW()
		WHERE id = $2 AND invited_by IS NULL`,
		inviterID, userID)
	return err
}

// GetUserIDByReferralCode looks up a user by their referral code.
func (r *ReferralRepo) GetUserIDByReferralCode(ctx context.Context, code string) (int, error) {
	var id int
	err := r.db.QueryRow(ctx, //nolint:gosec // G201/G202: parameterized query, no user input in SQL structure
		`SELECT id FROM users WHERE referral_code = $1`, code,
	).Scan(&id)
	return id, err
}

// GetPendingInvitationRequest returns the pending request for a user (if any).
func (r *ReferralRepo) GetPendingInvitationRequest(ctx context.Context, userID int) (*InvitationRequest, error) {
	req := &InvitationRequest{}
	//nolint:gosec // G201/G202: parameterized query, no user input in SQL structure
	err := r.db.QueryRow(ctx, `
		SELECT id, user_id, inviter_code_submitted, status, reject_reason, created_at
		FROM invitation_requests
		WHERE user_id = $1 AND status = 'pending'
		ORDER BY created_at DESC LIMIT 1`, userID,
	).Scan(&req.ID, &req.UserID, &req.InviterCodeSubmitted, &req.Status, &req.RejectReason, &req.CreatedAt)
	if err != nil {
		return nil, err
	}
	return req, nil
}

// SetReferralSuspended suspends or unsuspends a user's referral code (admin only).
func (r *ReferralRepo) SetReferralSuspended(ctx context.Context, userID int, suspended bool) error {
	tag, err := r.db.Exec(ctx, //nolint:gosec // G201/G202: parameterized query, no user input in SQL structure
		`UPDATE users SET referral_suspended = $1, updated_at = NOW() WHERE id = $2`,
		suspended, userID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}

// AdminListReferrals returns paginated referral relationships for the admin panel.
func (r *ReferralRepo) AdminListReferrals(ctx context.Context, page, limit int, search string) ([]AdminReferralRow, int, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	// Cap search length to prevent oversized ILIKE patterns (DoS via regex backtracking)
	if len(search) > 100 {
		search = search[:100]
	}

	// Build WHERE using independent counter — SQL string built before args populated.
	i := 0
	searchClause := ""
	var args []interface{}
	if search != "" {
		i += 3
		searchClause = fmt.Sprintf(" AND (u.name ILIKE $%d ESCAPE '\\' OR u.phone ILIKE $%d ESCAPE '\\' OR inv.name ILIKE $%d ESCAPE '\\')", i-2, i-1, i) //nolint:gosec // G201: Building static SQL fragment with computed $N parameter placeholders only
		escaped := likeWrap(search)
		args = append(args, escaped, escaped, escaped)
	}

	var total int
	//nolint:gosec // G201: searchClause built from static SQL fragments with $N placeholders only
	_ = r.db.QueryRow(ctx, "SELECT COUNT(*) FROM users u LEFT JOIN users inv ON inv.id = u.invited_by WHERE u.invited_by IS NOT NULL"+searchClause, args...).Scan(&total)

	args = append(args, limit, offset)
	n := len(args)
	selectSQL := "SELECT u.id, u.name, u.phone, u.role," +
		" inv.id, inv.name, inv.referral_code, inv.referral_suspended," +
		" u.invitation_accepted_at," +
		" COALESCE((SELECT SUM(commission_amount) FROM referral_commissions WHERE referrer_id = inv.id AND referred_id = u.id), 0)" +
		" FROM users u JOIN users inv ON inv.id = u.invited_by" +
		" WHERE u.invited_by IS NOT NULL" + searchClause +
		fmt.Sprintf(" ORDER BY u.invitation_accepted_at DESC LIMIT $%d OFFSET $%d", n-1, n)
	//nolint:gosec // G201: SQL built from static strings with computed $N placeholders, no user input in structure
	rows, err := r.db.Query(ctx, selectSQL, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var out []AdminReferralRow
	for rows.Next() {
		var row AdminReferralRow
		if err := rows.Scan(
			&row.InvitedUserID, &row.InvitedUserName, &row.InvitedUserPhone, &row.InvitedUserRole,
			&row.InviterID, &row.InviterName, &row.InviterCode, &row.InviterSuspended,
			&row.AcceptedAt, &row.TotalCommission,
		); err != nil {
			return nil, 0, err
		}
		out = append(out, row)
	}
	if out == nil {
		out = []AdminReferralRow{}
	}
	return out, total, rows.Err()
}

// AdminGetReferralOverview returns aggregate referral stats for the admin profit page.
func (r *ReferralRepo) AdminGetReferralOverview(ctx context.Context) (map[string]interface{}, error) {
	result := map[string]interface{}{}

	var totalCommissions float64
	var totalRows int
	_ = r.db.QueryRow(ctx, //nolint:gosec // G201/G202: parameterized query, no user input in SQL structure
		`SELECT COALESCE(SUM(commission_amount),0), COUNT(*) FROM referral_commissions`,
	).Scan(&totalCommissions, &totalRows)
	result["total_commissions_paid"] = totalCommissions
	result["total_commission_rows"] = totalRows

	var activeReferrals int
	_ = r.db.QueryRow(ctx, //nolint:gosec // G201/G202: parameterized query, no user input in SQL structure
		`SELECT COUNT(*) FROM users WHERE invited_by IS NOT NULL`,
	).Scan(&activeReferrals)
	result["active_referrals"] = activeReferrals

	var suspendedCodes int
	_ = r.db.QueryRow(ctx, //nolint:gosec // G201/G202: parameterized query, no user input in SQL structure
		`SELECT COUNT(*) FROM users WHERE referral_suspended = TRUE`,
	).Scan(&suspendedCodes)
	result["suspended_codes"] = suspendedCodes

	// Top inviters — close rows before opening typeRows
	{
		rows, err := r.db.Query(ctx, `
			SELECT rc.referrer_id, COALESCE(u.name,''), COUNT(*) AS commission_count,
			       SUM(rc.commission_amount) AS total_earned
			FROM referral_commissions rc
			LEFT JOIN users u ON u.id = rc.referrer_id
			GROUP BY rc.referrer_id, u.name
			ORDER BY total_earned DESC
			LIMIT 10`)
		var topInviters []map[string]interface{}
		if err == nil {
			for rows.Next() {
				var inviterID, count int
				var name string
				var earned float64
				if err := rows.Scan(&inviterID, &name, &count, &earned); err == nil {
					topInviters = append(topInviters, map[string]interface{}{
						"inviter_id": inviterID, "name": name,
						"commission_count": count, "total_earned": earned,
					})
				}
			}
			rows.Close()
		}
		if topInviters == nil {
			topInviters = []map[string]interface{}{}
		}
		result["top_inviters"] = topInviters
	}

	// Commissions by transaction type
	{
		typeRows, err := r.db.Query(ctx, `
			SELECT transaction_type, COUNT(*), SUM(commission_amount)
			FROM referral_commissions
			GROUP BY transaction_type
			ORDER BY SUM(commission_amount) DESC`)
		var byType []map[string]interface{}
		if err == nil {
			for typeRows.Next() {
				var txType string
				var count int
				var total float64
				if err := typeRows.Scan(&txType, &count, &total); err == nil {
					byType = append(byType, map[string]interface{}{
						"type": txType, "count": count, "total": total,
					})
				}
			}
			typeRows.Close()
		}
		if byType == nil {
			byType = []map[string]interface{}{}
		}
		result["by_type"] = byType
	}

	return result, nil
}

// SetPermanentlyBlocked permanently blocks or unblocks a user from the referral system.
func (r *ReferralRepo) SetPermanentlyBlocked(ctx context.Context, userID int, blocked bool, reason string) error {
	var reasonVal interface{}
	if reason != "" {
		reasonVal = reason
	}
	//nolint:gosec // G201/G202: parameterized query, no user input in SQL structure
	tag, err := r.db.Exec(ctx, `
		UPDATE users
		SET referral_permanently_blocked = $1,
		    referral_block_reason = $2,
		    updated_at = NOW()
		WHERE id = $3`,
		blocked, reasonVal, userID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}

// GetUserBalance returns a user's referral balance.
// Only counts commissions with status='confirmed' (excludes pending and reversed).
// available_balance subtracts pending withdrawals to prevent double-spending.
func (r *ReferralRepo) GetUserBalance(ctx context.Context, userID int) (*UserBalance, error) {
	bal := &UserBalance{}

	// Calculate balance using only confirmed commissions
	//nolint:gosec // G201/G202: parameterized query, no user input in SQL structure
	err := r.db.QueryRow(ctx, `
		SELECT
			COALESCE(SUM(commission_amount), 0) AS total_earned,
			COALESCE((SELECT SUM(amount) FROM referral_withdrawals 
			          WHERE user_id = $1 AND status IN ('approved', 'paid')), 0) AS total_withdrawn,
			COALESCE(SUM(commission_amount), 0) - 
			COALESCE((SELECT SUM(amount) FROM referral_withdrawals 
			          WHERE user_id = $1 AND status IN ('approved', 'paid', 'pending')), 0) AS available_balance
		FROM referral_commissions
		WHERE referrer_id = $1 AND status = 'confirmed'`, userID,
	).Scan(&bal.TotalEarned, &bal.TotalWithdrawn, &bal.AvailableBalance)
	if err != nil {
		return &UserBalance{}, nil
	}

	// Fetch pending commission amount for display (not withdrawable yet)
	//nolint:gosec // G201/G202: parameterized query, no user input in SQL structure
	_ = r.db.QueryRow(ctx, `
		SELECT COALESCE(SUM(commission_amount), 0)
		FROM referral_commissions
		WHERE referrer_id = $1 AND status = 'pending'`, userID,
	).Scan(&bal.PendingAmount)

	return bal, nil
}

// RequestWithdrawal creates a pending withdrawal request.
// The balance check and INSERT are done atomically using a CTE to prevent
// race conditions where two concurrent requests could both pass the balance check.
func (r *ReferralRepo) RequestWithdrawal(ctx context.Context, userID int, amount float64, paymentMethod, payerPhone string) (int, error) {
	minWithdrawal := constants.DefaultMinReferralWithdrawal
	var minVal string
	_ = r.db.QueryRow(ctx, `SELECT value FROM platform_settings WHERE key=$1`, constants.SettingMinReferralWithdrawal).Scan(&minVal) //nolint:gosec
	if minVal != "" {
		if f, err2 := strconv.ParseFloat(minVal, 64); err2 == nil && f > 0 {
			minWithdrawal = f
		}
	}
	maxWithdrawal := constants.DefaultMaxReferralWithdrawal
	var maxVal string
	_ = r.db.QueryRow(ctx, `SELECT value FROM platform_settings WHERE key=$1`, constants.SettingMaxReferralWithdrawal).Scan(&maxVal) //nolint:gosec
	if maxVal != "" {
		if f, err2 := strconv.ParseFloat(maxVal, 64); err2 == nil && f > 0 {
			maxWithdrawal = f
		}
	}

	if amount < minWithdrawal {
		return 0, fmt.Errorf("minimum withdrawal amount is ETB %.2f", minWithdrawal)
	}
	if amount > maxWithdrawal {
		return 0, fmt.Errorf("maximum withdrawal amount is ETB %.2f per request", maxWithdrawal)
	}

	// Validate payment method and phone before touching the DB
	validMethods := map[string]bool{"telebirr": true, "safaricom": true, "cbe": true, "bank": true}
	if !validMethods[paymentMethod] {
		return 0, fmt.Errorf("payment_method must be telebirr, safaricom, cbe, or bank")
	}
	if len(strings.TrimSpace(payerPhone)) < 7 {
		return 0, fmt.Errorf("payer_phone is required and must be valid")
	}

	// Check not permanently blocked
	var blocked bool
	_ = r.db.QueryRow(ctx, //nolint:gosec // G201/G202: parameterized query, no user input in SQL structure
		`SELECT referral_permanently_blocked FROM users WHERE id = $1`, userID,
	).Scan(&blocked)
	if blocked {
		return 0, fmt.Errorf("your referral account has been permanently blocked — contact support")
	}

	// Atomic balance check + INSERT using a CTE.
	// This prevents the TOCTOU race: the INSERT only proceeds if the subquery
	// confirms sufficient balance at the moment of insertion.
	// The unique partial index uq_referral_withdrawals_one_pending prevents
	// two concurrent pending rows even if both pass the balance check simultaneously.
	var id int
	//nolint:gosec // G201/G202: parameterized query, no user input in SQL structure
	err := r.db.QueryRow(ctx, `
		WITH balance AS (
			SELECT COALESCE(available_balance, 0) AS avail
			FROM referral_balance_view
			WHERE user_id = $1
		)
		INSERT INTO referral_withdrawals (user_id, amount, payment_method, payer_phone)
		SELECT $1, $2, $3, $4
		FROM balance
		WHERE avail >= $2
		RETURNING id`,
		userID, amount, paymentMethod, strings.TrimSpace(payerPhone),
	).Scan(&id)
	if err != nil {
		msg := err.Error()
		if strings.Contains(msg, "uq_referral_withdrawals_one_pending") ||
			strings.Contains(msg, "duplicate key") {
			return 0, fmt.Errorf("you already have a pending withdrawal request — wait for it to be processed")
		}
		if strings.Contains(msg, "no rows") || strings.Contains(msg, "ErrNoRows") {
			return 0, fmt.Errorf("insufficient balance — you need at least ETB %.2f available to withdraw that amount", amount)
		}
		return 0, err
	}
	return id, nil
}

func (r *ReferralRepo) GetMyWithdrawals(ctx context.Context, userID int) ([]WithdrawalRequest, error) {
	//nolint:gosec // G201/G202: parameterized query, no user input in SQL structure
	rows, err := r.db.Query(ctx, `
		SELECT id, user_id, amount, status, payment_method, payer_phone,
		       reference, admin_note, created_at, reviewed_at, reviewed_by, paid_at
		FROM referral_withdrawals
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT 50`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanWithdrawals(rows)
}

// AdminListWithdrawals returns paginated withdrawal requests for admin review.
func (r *ReferralRepo) AdminListWithdrawals(ctx context.Context, status string, page, limit int) ([]WithdrawalRequest, int, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	// Allowlist status to prevent arbitrary values reaching the SQL builder
	validWithdrawalStatus := map[string]bool{"pending": true, "approved": true, "rejected": true, "paid": true}
	if status != "" && status != "all" && !validWithdrawalStatus[status] {
		status = ""
	}

	// Build WHERE using independent counter — SQL string built before args populated.
	i := 0
	statusClause := ""
	if status != "" && status != "all" {
		i++
		statusClause = fmt.Sprintf(" AND w.status = $%d", i) //nolint:gosec // G201: Building static SQL fragment with computed $N parameter placeholder only
	}
	var args []interface{}
	if status != "" && status != "all" {
		args = append(args, status)
	}

	var total int
	//nolint:gosec // G201: statusClause built from static SQL fragments with $N placeholders only
	_ = r.db.QueryRow(ctx, "SELECT COUNT(*) FROM referral_withdrawals w WHERE 1=1"+statusClause, args...).Scan(&total)

	args = append(args, limit, offset)
	n := len(args)
	selectSQL := "SELECT w.id, w.user_id, COALESCE(u.name,''), COALESCE(u.phone,'')," +
		" w.amount, w.status, w.payment_method, w.payer_phone," +
		" w.reference, w.admin_note, w.created_at, w.reviewed_at, w.reviewed_by, w.paid_at" +
		" FROM referral_withdrawals w LEFT JOIN users u ON u.id = w.user_id" +
		" WHERE 1=1" + statusClause +
		fmt.Sprintf(" ORDER BY w.created_at DESC LIMIT $%d OFFSET $%d", n-1, n)
	//nolint:gosec // G201: SQL built from static strings with computed $N placeholders, no user input in structure
	rows, err := r.db.Query(ctx, selectSQL, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var out []WithdrawalRequest
	for rows.Next() {
		var w WithdrawalRequest
		if err := rows.Scan(
			&w.ID, &w.UserID, &w.UserName, &w.UserPhone,
			&w.Amount, &w.Status, &w.PaymentMethod, &w.PayerPhone,
			&w.Reference, &w.AdminNote, &w.RequestedAt, &w.ReviewedAt, &w.ReviewedBy, &w.PaidAt,
		); err != nil {
			return nil, 0, err
		}
		out = append(out, w)
	}
	if out == nil {
		out = []WithdrawalRequest{}
	}
	return out, total, rows.Err()
}

// AdminApproveWithdrawal marks a withdrawal approved.
func (r *ReferralRepo) AdminApproveWithdrawal(ctx context.Context, withdrawalID, adminID int, note string) error {
	var noteVal interface{}
	if note != "" {
		noteVal = note
	}
	//nolint:gosec // G201/G202: parameterized query, no user input in SQL structure
	tag, err := r.db.Exec(ctx, `
		UPDATE referral_withdrawals
		SET status = 'approved', reviewed_by = $1, reviewed_at = NOW(), admin_note = $2
		WHERE id = $3 AND status = 'pending'`,
		adminID, noteVal, withdrawalID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("withdrawal not found or already processed")
	}
	return nil
}

// AdminRejectWithdrawal rejects a withdrawal with a reason.
func (r *ReferralRepo) AdminRejectWithdrawal(ctx context.Context, withdrawalID, adminID int, reason string) error {
	if reason == "" {
		return fmt.Errorf("rejection reason is required")
	}
	//nolint:gosec // G201/G202: parameterized query, no user input in SQL structure
	tag, err := r.db.Exec(ctx, `
		UPDATE referral_withdrawals
		SET status = 'rejected', reviewed_by = $1, reviewed_at = NOW(), admin_note = $2
		WHERE id = $3 AND status = 'pending'`,
		adminID, reason, withdrawalID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("withdrawal not found or already processed")
	}
	return nil
}

// AdminMarkWithdrawalPaid marks a withdrawal as paid and records the reference.
func (r *ReferralRepo) AdminMarkWithdrawalPaid(ctx context.Context, withdrawalID, adminID int, reference string) error {
	if reference == "" {
		return fmt.Errorf("payment reference is required")
	}
	//nolint:gosec // G201/G202: parameterized query, no user input in SQL structure
	tag, err := r.db.Exec(ctx, `
		UPDATE referral_withdrawals
		SET status = 'paid', reviewed_by = $1, paid_at = NOW(), reference = $2
		WHERE id = $3 AND status = 'approved'`,
		adminID, reference, withdrawalID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("withdrawal not found, not yet approved, or already paid/rejected — approve it first")
	}
	return nil
}

// scanWithdrawals scans rows from referral_withdrawals (user-facing, no name/phone columns).
func scanWithdrawals(rows interface {
	Next() bool
	Scan(...interface{}) error
	Err() error
}) ([]WithdrawalRequest, error) {
	var out []WithdrawalRequest
	for rows.Next() {
		var w WithdrawalRequest
		if err := rows.Scan(
			&w.ID, &w.UserID, &w.Amount, &w.Status, &w.PaymentMethod, &w.PayerPhone,
			&w.Reference, &w.AdminNote, &w.RequestedAt, &w.ReviewedAt, &w.ReviewedBy, &w.PaidAt,
		); err != nil {
			return nil, err
		}
		out = append(out, w)
	}
	if out == nil {
		out = []WithdrawalRequest{}
	}
	return out, rows.Err()
}
