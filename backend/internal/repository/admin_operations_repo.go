package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/your-org/used-tools-marketplace/internal/models"
	"strconv"
	"strings"
)

func (r *AdminRepo) GetOperationsHealth(ctx context.Context) (*models.OperationsHealth, error) {
	var h models.OperationsHealth
	err := r.db.QueryRow(ctx, `SELECT * FROM operations_health LIMIT 1`).Scan( //nolint:gosec
		&h.TotalJobs24h, &h.FailedJobs24h, &h.StuckJobs,
		&h.InactiveUsers, &h.StalePendingProducts, &h.ExpiredOTPs, &h.StaleUpgradeRequests,
		&h.DatabaseSize, &h.ProductsTableSize, &h.UsersTableSize)
	if err != nil {
		return nil, err
	}
	return &h, nil
}

// ListJobStatusSummary returns job execution statistics
func (r *AdminRepo) ListJobStatusSummary(ctx context.Context) ([]models.JobStatusSummary, error) {
	rows, err := r.db.Query(ctx, `SELECT * FROM job_status_summary`) //nolint:gosec
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []models.JobStatusSummary
	for rows.Next() {
		var j models.JobStatusSummary
		if err := rows.Scan(&j.JobName, &j.TotalRuns, &j.SuccessfulRuns, &j.FailedRuns,
			&j.CurrentlyRunning, &j.LastRunAt, &j.LastSuccessAt, &j.LastFailureAt,
			&j.AvgDurationMS, &j.TotalRecordsProcessed); err != nil {
			return nil, err
		}
		jobs = append(jobs, j)
	}
	return jobs, rows.Err()
}

// ListRecentFailedJobs returns recent failed job executions
func (r *AdminRepo) ListRecentFailedJobs(ctx context.Context, limit int) ([]models.JobExecutionLog, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	rows, err := r.db.Query(ctx, `SELECT * FROM recent_failed_jobs LIMIT $1`, limit) //nolint:gosec
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []models.JobExecutionLog
	for rows.Next() {
		var j models.JobExecutionLog
		if err := rows.Scan(&j.ID, &j.JobName, &j.StartedAt, &j.CompletedAt,
			&j.DurationMS, &j.RecordsProcessed, &j.ErrorMessage, &j.CreatedAt); err != nil {
			return nil, err
		}
		jobs = append(jobs, j)
	}
	return jobs, rows.Err()
}

// ── Seller Verification (Role 08) ─────────────────────────────────────────────

// ListVerificationDocuments returns seller verification documents with pagination
func (r *AdminRepo) ListVerificationDocuments(ctx context.Context, status string, page, limit int) ([]models.SellerVerificationDocument, int, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	var wb strings.Builder
	wb.WriteString("1=1")
	args := []interface{}{}

	if status != "" && status != "all" {
		args = appendCond(&wb, args, " AND d.status=", status)
	}
	whereStr := wb.String()

	var total int
	_ = r.db.QueryRow(ctx, `SELECT COUNT(*) FROM seller_verification_documents d WHERE `+whereStr, args...).Scan(&total) //nolint:gosec

	i := len(args) + 1
	args = append(args, limit, offset)
	rows, err := r.db.Query(ctx, //nolint:gosec
		`SELECT d.id, d.seller_id, d.document_type, d.file_url, COALESCE(d.file_name,''),
		        d.file_size_bytes, d.status, d.reviewed_by, d.reviewed_at,
		        COALESCE(d.reject_reason,''), COALESCE(d.notes,''), d.created_at, d.updated_at,
		        u.name AS seller_name, COALESCE(r.name,'') AS reviewer_name
		 FROM seller_verification_documents d
		 JOIN users u ON u.id = d.seller_id
		 LEFT JOIN users r ON r.id = d.reviewed_by
		 WHERE `+whereStr+`
		 ORDER BY d.created_at DESC
		 LIMIT $`+strconv.Itoa(i)+` OFFSET $`+strconv.Itoa(i+1),
		args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var docs []models.SellerVerificationDocument
	for rows.Next() {
		var d models.SellerVerificationDocument
		if err := rows.Scan(&d.ID, &d.SellerID, &d.DocumentType, &d.FileURL, &d.FileName,
			&d.FileSizeBytes, &d.Status, &d.ReviewedBy, &d.ReviewedAt,
			&d.RejectReason, &d.Notes, &d.CreatedAt, &d.UpdatedAt,
			&d.SellerName, &d.ReviewerName); err != nil {
			return nil, 0, err
		}
		docs = append(docs, d)
	}
	return docs, total, rows.Err()
}

// ReviewVerificationDocument approves or rejects a verification document
func (r *AdminRepo) ReviewVerificationDocument(ctx context.Context, docID, adminID int, status, reason, notes string) error {
	if status != "approved" && status != "rejected" {
		return fmt.Errorf("status must be approved or rejected")
	}

	tag, err := r.db.Exec(ctx, //nolint:gosec
		`UPDATE seller_verification_documents
		 SET status=$1, reviewed_by=$2, reviewed_at=NOW(), reject_reason=$3, notes=$4, updated_at=NOW()
		 WHERE id=$5`,
		status, adminID, reason, notes, docID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("document not found")
	}
	return nil
}

// ── Trust & Safety (Role 21) ──────────────────────────────────────────────────

// ListAbuseReports returns abuse reports with pagination and filtering
func (r *AdminRepo) ListAbuseReports(ctx context.Context, f models.AbuseReportFilter) ([]models.AbuseReport, int, error) {
	if f.Limit <= 0 || f.Limit > 100 {
		f.Limit = 50
	}
	if f.Page < 1 {
		f.Page = 1
	}
	offset := (f.Page - 1) * f.Limit

	var wb strings.Builder
	wb.WriteString("1=1")
	args := []interface{}{}

	if f.Status != "" && f.Status != "all" {
		args = appendCond(&wb, args, " AND a.status=", f.Status)
	}
	if f.Priority != "" && f.Priority != "all" {
		args = appendCond(&wb, args, " AND a.priority=", f.Priority)
	}
	if f.Reason != "" && f.Reason != "all" {
		args = appendCond(&wb, args, " AND a.reason=", f.Reason)
	}
	if f.ReportedUserID > 0 {
		args = appendCond(&wb, args, " AND a.reported_user_id=", f.ReportedUserID)
	}
	whereStr := wb.String()

	var total int
	_ = r.db.QueryRow(ctx, `SELECT COUNT(*) FROM abuse_reports a WHERE `+whereStr, args...).Scan(&total) //nolint:gosec

	i := len(args) + 1
	args = append(args, f.Limit, offset)
	rows, err := r.db.Query(ctx, //nolint:gosec
		`SELECT a.id, a.reporter_id, a.reporter_role, a.reported_user_id, a.reported_product_id,
		        a.reason, a.description, a.evidence_urls, a.status, a.priority,
		        a.reviewed_by, a.reviewed_at, COALESCE(a.action_taken,''), COALESCE(a.resolution_notes,''),
		        a.created_at, a.updated_at,
		        COALESCE(rep.name,'Anonymous') AS reporter_name,
		        COALESCE(ru.name,'') AS reported_user_name,
		        COALESCE(p.title,'') AS reported_product_title,
		        COALESCE(rev.name,'') AS reviewer_name
		 FROM abuse_reports a
		 LEFT JOIN users rep ON rep.id = a.reporter_id
		 LEFT JOIN users ru ON ru.id = a.reported_user_id
		 LEFT JOIN products p ON p.id = a.reported_product_id
		 LEFT JOIN users rev ON rev.id = a.reviewed_by
		 WHERE `+whereStr+`
		 ORDER BY a.priority DESC, a.created_at DESC
		 LIMIT $`+strconv.Itoa(i)+` OFFSET $`+strconv.Itoa(i+1),
		args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var reports []models.AbuseReport
	for rows.Next() {
		var a models.AbuseReport
		if err := rows.Scan(&a.ID, &a.ReporterID, &a.ReporterRole, &a.ReportedUserID, &a.ReportedProductID,
			&a.Reason, &a.Description, &a.EvidenceURLs, &a.Status, &a.Priority,
			&a.ReviewedBy, &a.ReviewedAt, &a.ActionTaken, &a.ResolutionNotes,
			&a.CreatedAt, &a.UpdatedAt,
			&a.ReporterName, &a.ReportedUserName, &a.ReportedProductTitle, &a.ReviewerName); err != nil {
			return nil, 0, err
		}
		reports = append(reports, a)
	}
	return reports, total, rows.Err()
}

// ReviewAbuseReport updates the status and resolution of an abuse report and returns reported_user_id if any
func (r *AdminRepo) ReviewAbuseReport(ctx context.Context, reportID, adminID int, status, actionTaken, notes string) (*int, error) {
	validStatus := map[string]bool{"investigating": true, "resolved": true, "dismissed": true}
	if !validStatus[status] {
		return nil, fmt.Errorf("invalid status: must be investigating, resolved, or dismissed")
	}

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var reportedUserID *int
	var reportedProductID *int
	err = tx.QueryRow(ctx, //nolint:gosec
		`UPDATE abuse_reports
		 SET status=$1, reviewed_by=$2, reviewed_at=NOW(), action_taken=$3, resolution_notes=$4, updated_at=NOW()
		 WHERE id=$5
		 RETURNING reported_user_id, reported_product_id`,
		status, adminID, actionTaken, notes, reportID).Scan(&reportedUserID, &reportedProductID)
	if err != nil {
		return nil, err
	}

	// If action is removed or content_removed, reject the reported product to hide from marketplace
	if (actionTaken == "removed" || actionTaken == "content_removed") && reportedProductID != nil && *reportedProductID > 0 {
		_, _ = tx.Exec(ctx, `UPDATE products SET status='rejected', rejection_reason=$1, updated_at=NOW() WHERE id=$2`, notes, *reportedProductID)
	}

	// If action is suspended, atomically deactivate user and suspend capabilities
	if actionTaken == "suspended" && reportedUserID != nil && *reportedUserID > 0 {
		reason := "Account suspended via abuse report #" + strconv.Itoa(reportID)
		if notes != "" {
			reason += ": " + notes
		}
		_, err = tx.Exec(ctx,
			`UPDATE users 
			 SET is_active=FALSE, is_banned=FALSE, referral_suspended=TRUE, deactivation_reason=$1, updated_at=NOW() 
			 WHERE id=$2`,
			reason, *reportedUserID)
		if err != nil {
			return nil, fmt.Errorf("failed to suspend user: %w", err)
		}
		_, err = tx.Exec(ctx,
			`UPDATE user_capabilities 
			 SET status='suspended', 
			     metadata=jsonb_set(COALESCE(metadata, '{}'::jsonb), '{suspend_reason}', to_jsonb($2::text)) 
			 WHERE user_id=$1`,
			*reportedUserID, reason)
		if err != nil {
			return nil, fmt.Errorf("failed to suspend user capabilities: %w", err)
		}
	}

	// If action is banned, atomically ban user and revoke capabilities
	if actionTaken == "banned" && reportedUserID != nil && *reportedUserID > 0 {
		reason := "Account banned via abuse report #" + strconv.Itoa(reportID)
		if notes != "" {
			reason += ": " + notes
		}
		_, err = tx.Exec(ctx,
			`UPDATE users 
			 SET is_banned=TRUE, is_active=FALSE, referral_suspended=TRUE, referral_permanently_blocked=TRUE, referral_block_reason=$1, deactivation_reason=$1, updated_at=NOW() 
			 WHERE id=$2`,
			reason, *reportedUserID)
		if err != nil {
			return nil, fmt.Errorf("failed to ban user: %w", err)
		}
		_, err = tx.Exec(ctx,
			`UPDATE user_capabilities 
			 SET status='revoked', 
			     metadata=jsonb_set(COALESCE(metadata, '{}'::jsonb), '{ban_reason}', to_jsonb($2::text)) 
			 WHERE user_id=$1`,
			*reportedUserID, reason)
		if err != nil {
			return nil, fmt.Errorf("failed to revoke user capabilities: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return reportedUserID, nil
}

// UpdateAbuseReportPriority changes the priority of an abuse report
func (r *AdminRepo) UpdateAbuseReportPriority(ctx context.Context, reportID int, priority string) error {
	validPriority := map[string]bool{"low": true, "normal": true, "high": true, "critical": true}
	if !validPriority[priority] {
		return fmt.Errorf("invalid priority")
	}

	_, err := r.db.Exec(ctx, //nolint:gosec
		`UPDATE abuse_reports SET priority=$1, updated_at=NOW() WHERE id=$2`,
		priority, reportID)
	return err
}

// ── Refund Management (Role 20) ───────────────────────────────────────────────

// ListRefundRequests returns refund requests with pagination and filtering
func (r *AdminRepo) ListRefundRequests(ctx context.Context, f models.RefundRequestFilter) ([]models.RefundRequest, int, error) {
	if f.Limit <= 0 || f.Limit > 100 {
		f.Limit = 50
	}
	if f.Page < 1 {
		f.Page = 1
	}
	offset := (f.Page - 1) * f.Limit

	var wb strings.Builder
	wb.WriteString("1=1")
	args := []interface{}{}

	if f.Status != "" && f.Status != "all" {
		args = appendCond(&wb, args, " AND r.status=", f.Status)
	}
	if f.RequesterRole != "" && f.RequesterRole != "all" {
		args = appendCond(&wb, args, " AND r.requester_role=", f.RequesterRole)
	}
	whereStr := wb.String()

	var total int
	_ = r.db.QueryRow(ctx, `SELECT COUNT(*) FROM refund_requests r WHERE `+whereStr, args...).Scan(&total) //nolint:gosec

	i := len(args) + 1
	args = append(args, f.Limit, offset)
	rows, err := r.db.Query(ctx, //nolint:gosec
		`SELECT r.id, r.transaction_id, r.product_id, r.requester_id, r.requester_role,
		        r.amount, r.reason, r.description, r.evidence_urls, r.status,
		        r.reviewed_by, r.reviewed_at, COALESCE(r.reject_reason,''),
		        COALESCE(r.refund_reference,''), r.processed_at, r.created_at, r.updated_at,
		        u.name AS requester_name,
		        COALESCE(p.title,'') AS product_title,
		        COALESCE(rev.name,'') AS reviewer_name
		 FROM refund_requests r
		 JOIN users u ON u.id = r.requester_id
		 LEFT JOIN products p ON p.id = r.product_id
		 LEFT JOIN users rev ON rev.id = r.reviewed_by
		 WHERE `+whereStr+`
		 ORDER BY r.created_at DESC
		 LIMIT $`+strconv.Itoa(i)+` OFFSET $`+strconv.Itoa(i+1),
		args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var refunds []models.RefundRequest
	for rows.Next() {
		var rf models.RefundRequest
		if err := rows.Scan(&rf.ID, &rf.TransactionID, &rf.ProductID, &rf.RequesterID, &rf.RequesterRole,
			&rf.Amount, &rf.Reason, &rf.Description, &rf.EvidenceURLs, &rf.Status,
			&rf.ReviewedBy, &rf.ReviewedAt, &rf.RejectReason,
			&rf.RefundReference, &rf.ProcessedAt, &rf.CreatedAt, &rf.UpdatedAt,
			&rf.RequesterName, &rf.ProductTitle, &rf.ReviewerName); err != nil {
			return nil, 0, err
		}
		refunds = append(refunds, rf)
	}
	return refunds, total, rows.Err()
}

// ReviewRefundRequest approves or rejects a refund request
func (r *AdminRepo) ReviewRefundRequest(ctx context.Context, refundID, adminID int, approved bool, reason string) error {
	status := "approved"
	if !approved {
		status = "rejected"
	}

	var requesterID int
	err := r.db.QueryRow(ctx, `SELECT requester_id FROM refund_requests WHERE id=$1`, refundID).Scan(&requesterID)
	if err != nil {
		return fmt.Errorf("refund request not found")
	}
	if requesterID == adminID {
		return errors.New("maker-checker separation: administrator cannot review their own refund request")
	}

	tag, err := r.db.Exec(ctx, //nolint:gosec
		`UPDATE refund_requests
		 SET status=$1, reviewed_by=$2, reviewed_at=NOW(), reject_reason=$3, updated_at=NOW()
		 WHERE id=$4 AND status='pending'`,
		status, adminID, reason, refundID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("refund request not found or already processed")
	}
	return nil
}

// MarkRefundProcessed marks a refund as processed with reference
func (r *AdminRepo) MarkRefundProcessed(ctx context.Context, refundID int, reference string) error {
	tag, err := r.db.Exec(ctx, //nolint:gosec
		`UPDATE refund_requests
		 SET status='processed', refund_reference=$1, processed_at=NOW(), updated_at=NOW()
		 WHERE id=$2 AND status='approved'`,
		reference, refundID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("refund request not found or not approved")
	}
	return nil
}

// ── Account Appeals (Role 30) ─────────────────────────────────────────────────

// ListAccountAppeals returns account appeals with pagination and filtering
func (r *AdminRepo) ListAccountAppeals(ctx context.Context, f models.AccountAppealFilter) ([]models.AccountAppeal, int, error) {
	if f.Limit <= 0 || f.Limit > 100 {
		f.Limit = 50
	}
	if f.Page < 1 {
		f.Page = 1
	}
	offset := (f.Page - 1) * f.Limit

	var wb strings.Builder
	wb.WriteString("1=1")
	args := []interface{}{}

	if f.Status != "" && f.Status != "all" {
		args = appendCond(&wb, args, " AND a.status=", f.Status)
	}
	if f.AppealType != "" && f.AppealType != "all" {
		args = appendCond(&wb, args, " AND a.appeal_type=", f.AppealType)
	}
	if f.UserID > 0 {
		args = appendCond(&wb, args, " AND a.user_id=", f.UserID)
	}
	whereStr := wb.String()

	var total int
	_ = r.db.QueryRow(ctx, `SELECT COUNT(*) FROM account_appeals a WHERE `+whereStr, args...).Scan(&total) //nolint:gosec

	i := len(args) + 1
	args = append(args, f.Limit, offset)
	rows, err := r.db.Query(ctx, //nolint:gosec
		`SELECT a.id, a.user_id, a.appeal_type, a.related_product_id, a.related_document_id,
		        a.reason, a.evidence_urls, a.status, a.reviewed_by, a.reviewed_at,
		        COALESCE(a.decision,''), COALESCE(a.action_taken,''), a.created_at, a.updated_at,
		        u.name AS user_name,
		        COALESCE(rev.name,'') AS reviewer_name,
		        COALESCE(p.title,'') AS product_title
		 FROM account_appeals a
		 JOIN users u ON u.id = a.user_id
		 LEFT JOIN users rev ON rev.id = a.reviewed_by
		 LEFT JOIN products p ON p.id = a.related_product_id
		 WHERE `+whereStr+`
		 ORDER BY a.created_at DESC
		 LIMIT $`+strconv.Itoa(i)+` OFFSET $`+strconv.Itoa(i+1),
		args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var appeals []models.AccountAppeal
	for rows.Next() {
		var ap models.AccountAppeal
		if err := rows.Scan(&ap.ID, &ap.UserID, &ap.AppealType, &ap.RelatedProductID, &ap.RelatedDocumentID,
			&ap.Reason, &ap.EvidenceURLs, &ap.Status, &ap.ReviewedBy, &ap.ReviewedAt,
			&ap.Decision, &ap.ActionTaken, &ap.CreatedAt, &ap.UpdatedAt,
			&ap.UserName, &ap.ReviewerName, &ap.ProductTitle); err != nil {
			return nil, 0, err
		}
		appeals = append(appeals, ap)
	}
	return appeals, total, rows.Err()
}

// ReviewAccountAppeal approves or rejects an account appeal
func (r *AdminRepo) ReviewAccountAppeal(ctx context.Context, appealID, adminID int, approved bool, decision, actionTaken string) error {
	status := "approved"
	if !approved {
		status = "rejected"
	}

	tag, err := r.db.Exec(ctx, //nolint:gosec
		`UPDATE account_appeals
		 SET status=$1, reviewed_by=$2, reviewed_at=NOW(), decision=$3, action_taken=$4, updated_at=NOW()
		 WHERE id=$5 AND status IN ('pending', 'under_review')`,
		status, adminID, decision, actionTaken, appealID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("appeal not found or already reviewed")
	}
	return nil
}

// UpdateAppealStatus changes the status of an appeal (e.g. to under_review)
func (r *AdminRepo) UpdateAppealStatus(ctx context.Context, appealID int, status string) error {
	validStatus := map[string]bool{"pending": true, "under_review": true, "approved": true, "rejected": true}
	if !validStatus[status] {
		return fmt.Errorf("invalid status")
	}

	_, err := r.db.Exec(ctx, //nolint:gosec
		`UPDATE account_appeals SET status=$1, updated_at=NOW() WHERE id=$2`,
		status, appealID)
	return err
}

// CreateVerificationDocument allows a seller to submit an identity/business document
func (r *AdminRepo) CreateVerificationDocument(ctx context.Context, doc *models.SellerVerificationDocument) (int, error) {
	var id int
	err := r.db.QueryRow(ctx,
		`INSERT INTO seller_verification_documents (seller_id, document_type, file_url, file_name, file_size_bytes, status)
		 VALUES ($1, $2, $3, $4, $5, 'pending')
		 RETURNING id`,
		doc.SellerID, doc.DocumentType, doc.FileURL, doc.FileName, doc.FileSizeBytes,
	).Scan(&id)
	return id, err
}

// ListSellerVerificationDocuments returns all documents uploaded by a specific seller
func (r *AdminRepo) ListSellerVerificationDocuments(ctx context.Context, sellerID int) ([]models.SellerVerificationDocument, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, seller_id, document_type, file_url, COALESCE(file_name,''),
		        file_size_bytes, status, reviewed_by, reviewed_at,
		        COALESCE(reject_reason,''), COALESCE(notes,''), created_at, updated_at
		 FROM seller_verification_documents
		 WHERE seller_id = $1
		 ORDER BY created_at DESC`,
		sellerID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var docs []models.SellerVerificationDocument
	for rows.Next() {
		var d models.SellerVerificationDocument
		if err := rows.Scan(
			&d.ID, &d.SellerID, &d.DocumentType, &d.FileURL, &d.FileName,
			&d.FileSizeBytes, &d.Status, &d.ReviewedBy, &d.ReviewedAt,
			&d.RejectReason, &d.Notes, &d.CreatedAt, &d.UpdatedAt,
		); err != nil {
			return nil, err
		}
		docs = append(docs, d)
	}
	return docs, rows.Err()
}

// CreateAbuseReport allows buyers or sellers to report abuse/fraud
func (r *AdminRepo) CreateAbuseReport(ctx context.Context, report *models.AbuseReport) (int, error) {
	priority := report.Priority
	if priority == "" {
		priority = "normal"
	}
	var id int
	err := r.db.QueryRow(ctx,
		`INSERT INTO abuse_reports (reporter_id, reporter_role, reported_user_id, reported_product_id, reason, description, evidence_urls, priority, status)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, 'pending')
		 RETURNING id`,
		report.ReporterID, report.ReporterRole, report.ReportedUserID, report.ReportedProductID,
		report.Reason, report.Description, report.EvidenceURLs, priority,
	).Scan(&id)
	return id, err
}

// CreateRefundRequest allows buyers/sellers to submit a refund/dispute request
func (r *AdminRepo) CreateRefundRequest(ctx context.Context, req *models.RefundRequest) (int, error) {
	if req.Amount <= 0 {
		return 0, fmt.Errorf("refund amount must be greater than zero")
	}
	if req.Reason == "" {
		return 0, fmt.Errorf("reason is required")
	}

	// If transaction_id is provided, validate ownership, status, and maximum amount
	if req.TransactionID != nil && *req.TransactionID > 0 {
		var txSellerID int
		var txBuyerPhone string
		var txAmount float64
		var txProductID int
		err := r.db.QueryRow(ctx, `
			SELECT seller_id, buyer_phone, amount, product_id 
			FROM transactions 
			WHERE id = $1`, *req.TransactionID).Scan(&txSellerID, &txBuyerPhone, &txAmount, &txProductID)
		if err != nil {
			return 0, fmt.Errorf("transaction not found")
		}

		// Verify authorization against transaction parties
		switch req.RequesterRole {
		case "seller":
			if txSellerID != req.RequesterID {
				return 0, fmt.Errorf("you are not authorized to request a refund for this transaction")
			}
		case "buyer":
			var requesterPhone string
			_ = r.db.QueryRow(ctx, `SELECT phone FROM users WHERE id = $1`, req.RequesterID).Scan(&requesterPhone)
			if requesterPhone == "" || requesterPhone != txBuyerPhone {
				return 0, fmt.Errorf("you are not authorized to request a refund for this transaction")
			}
		}

		if req.Amount > txAmount {
			return 0, fmt.Errorf("refund amount (ETB %.2f) cannot exceed transaction amount (ETB %.2f)", req.Amount, txAmount)
		}

		// Prevent duplicate pending or processed refund requests for the same transaction
		var hasExisting bool
		_ = r.db.QueryRow(ctx, `
			SELECT EXISTS (
				SELECT 1 FROM refund_requests 
				WHERE transaction_id = $1 AND status IN ('pending', 'approved', 'processed')
			)`, *req.TransactionID).Scan(&hasExisting)
		if hasExisting {
			return 0, fmt.Errorf("a refund request has already been submitted for this transaction")
		}

		if req.ProductID == nil && txProductID > 0 {
			req.ProductID = &txProductID
		}
	}

	var id int
	err := r.db.QueryRow(ctx,
		`INSERT INTO refund_requests (transaction_id, product_id, requester_id, requester_role, amount, reason, description, evidence_urls, status)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, 'pending')
		 RETURNING id`,
		req.TransactionID, req.ProductID, req.RequesterID, req.RequesterRole,
		req.Amount, req.Reason, req.Description, req.EvidenceURLs,
	).Scan(&id)
	return id, err
}

// CreateAccountAppeal allows a user to appeal an account ban or action
func (r *AdminRepo) CreateAccountAppeal(ctx context.Context, appeal *models.AccountAppeal) (int, error) {
	// Verify user exists in database
	var userExists bool
	_ = r.db.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM users WHERE id=$1)`, appeal.UserID).Scan(&userExists)
	if !userExists {
		return 0, fmt.Errorf("user not found")
	}

	// Prevent duplicate pending appeal spam
	var duplicatePending bool
	_ = r.db.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM account_appeals WHERE user_id=$1 AND appeal_type=$2 AND status='pending')`,
		appeal.UserID, appeal.AppealType).Scan(&duplicatePending)
	if duplicatePending {
		return 0, fmt.Errorf("an appeal of this type is already pending review")
	}

	var id int
	err := r.db.QueryRow(ctx,
		`INSERT INTO account_appeals (user_id, appeal_type, related_product_id, related_document_id, reason, evidence_urls, status)
		 VALUES ($1, $2, $3, $4, $5, $6, 'pending')
		 RETURNING id`,
		appeal.UserID, appeal.AppealType, appeal.RelatedProductID, appeal.RelatedDocumentID,
		appeal.Reason, appeal.EvidenceURLs,
	).Scan(&id)
	return id, err
}
