package models

import "time"

// JobExecutionLog tracks background job runs for monitoring
type JobExecutionLog struct {
	ID                int        `json:"id"`
	JobName           string     `json:"job_name"`
	Status            string     `json:"status"` // 'success', 'failed', 'running'
	StartedAt         time.Time  `json:"started_at"`
	CompletedAt       *time.Time `json:"completed_at,omitempty"`
	DurationMS        *int       `json:"duration_ms,omitempty"`
	RecordsProcessed  int        `json:"records_processed"`
	ErrorMessage      string     `json:"error_message,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
}

// JobStatusSummary aggregates job execution stats
type JobStatusSummary struct {
	JobName                string     `json:"job_name"`
	TotalRuns              int        `json:"total_runs"`
	SuccessfulRuns         int        `json:"successful_runs"`
	FailedRuns             int        `json:"failed_runs"`
	CurrentlyRunning       int        `json:"currently_running"`
	LastRunAt              *time.Time `json:"last_run_at,omitempty"`
	LastSuccessAt          *time.Time `json:"last_success_at,omitempty"`
	LastFailureAt          *time.Time `json:"last_failure_at,omitempty"`
	AvgDurationMS          *float64   `json:"avg_duration_ms,omitempty"`
	TotalRecordsProcessed  int64      `json:"total_records_processed"`
}

// OperationsHealth contains system health metrics
type OperationsHealth struct {
	TotalJobs24h             int    `json:"total_jobs_24h"`
	FailedJobs24h            int    `json:"failed_jobs_24h"`
	StuckJobs                int    `json:"stuck_jobs"`
	InactiveUsers            int    `json:"inactive_users"`
	StalePendingProducts     int    `json:"stale_pending_products"`
	ExpiredOTPs              int    `json:"expired_otps"`
	StaleUpgradeRequests     int    `json:"stale_upgrade_requests"`
	DatabaseSize             string `json:"database_size"`
	ProductsTableSize        string `json:"products_table_size"`
	UsersTableSize           string `json:"users_table_size"`
}

// SellerVerificationDocument tracks seller identity verification documents
type SellerVerificationDocument struct {
	ID             int        `json:"id"`
	SellerID       int        `json:"seller_id"`
	DocumentType   string     `json:"document_type"` // 'national_id', 'business_license', etc.
	FileURL        string     `json:"file_url"`
	FileName       string     `json:"file_name,omitempty"`
	FileSizeBytes  int        `json:"file_size_bytes,omitempty"`
	Status         string     `json:"status"` // 'pending', 'approved', 'rejected'
	ReviewedBy     *int       `json:"reviewed_by,omitempty"`
	ReviewedAt     *time.Time `json:"reviewed_at,omitempty"`
	RejectReason   string     `json:"reject_reason,omitempty"`
	Notes          string     `json:"notes,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	// Joined fields
	SellerName     string     `json:"seller_name,omitempty"`
	ReviewerName   string     `json:"reviewer_name,omitempty"`
}

// AbuseReport tracks user-reported incidents
type AbuseReport struct {
	ID                  int        `json:"id"`
	ReporterID          *int       `json:"reporter_id,omitempty"`
	ReporterRole        string     `json:"reporter_role"` // 'seller', 'buyer', 'admin', 'anonymous'
	ReportedUserID      *int       `json:"reported_user_id,omitempty"`
	ReportedProductID   *int       `json:"reported_product_id,omitempty"`
	Reason              string     `json:"reason"` // 'fraud', 'spam', 'inappropriate', etc.
	Description         string     `json:"description"`
	EvidenceURLs        []string   `json:"evidence_urls,omitempty"`
	Status              string     `json:"status"` // 'pending', 'investigating', 'resolved', 'dismissed'
	Priority            string     `json:"priority"` // 'low', 'normal', 'high', 'critical'
	ReviewedBy          *int       `json:"reviewed_by,omitempty"`
	ReviewedAt          *time.Time `json:"reviewed_at,omitempty"`
	ActionTaken         string     `json:"action_taken,omitempty"`
	ResolutionNotes     string     `json:"resolution_notes,omitempty"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
	// Joined fields
	ReporterName        string     `json:"reporter_name,omitempty"`
	ReportedUserName    string     `json:"reported_user_name,omitempty"`
	ReportedProductTitle string    `json:"reported_product_title,omitempty"`
	ReviewerName        string     `json:"reviewer_name,omitempty"`
}

// RefundRequest tracks refund and dispute requests
type RefundRequest struct {
	ID                int        `json:"id"`
	TransactionID     *int       `json:"transaction_id,omitempty"`
	ProductID         *int       `json:"product_id,omitempty"`
	RequesterID       int        `json:"requester_id"`
	RequesterRole     string     `json:"requester_role"` // 'seller', 'buyer'
	Amount            float64    `json:"amount"`
	Reason            string     `json:"reason"`
	Description       string     `json:"description"`
	EvidenceURLs      []string   `json:"evidence_urls,omitempty"`
	Status            string     `json:"status"` // 'pending', 'approved', 'rejected', 'processed'
	ReviewedBy        *int       `json:"reviewed_by,omitempty"`
	ReviewedAt        *time.Time `json:"reviewed_at,omitempty"`
	RejectReason      string     `json:"reject_reason,omitempty"`
	RefundReference   string     `json:"refund_reference,omitempty"`
	ProcessedAt       *time.Time `json:"processed_at,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	// Joined fields
	RequesterName     string     `json:"requester_name,omitempty"`
	ProductTitle      string     `json:"product_title,omitempty"`
	ReviewerName      string     `json:"reviewer_name,omitempty"`
}

// AccountAppeal tracks user appeals for account actions
type AccountAppeal struct {
	ID                  int        `json:"id"`
	UserID              int        `json:"user_id"`
	AppealType          string     `json:"appeal_type"` // 'ban', 'deactivation', 'listing_rejection', etc.
	RelatedProductID    *int       `json:"related_product_id,omitempty"`
	RelatedDocumentID   *int       `json:"related_document_id,omitempty"`
	Reason              string     `json:"reason"`
	EvidenceURLs        []string   `json:"evidence_urls,omitempty"`
	Status              string     `json:"status"` // 'pending', 'under_review', 'approved', 'rejected'
	ReviewedBy          *int       `json:"reviewed_by,omitempty"`
	ReviewedAt          *time.Time `json:"reviewed_at,omitempty"`
	Decision            string     `json:"decision,omitempty"`
	ActionTaken         string     `json:"action_taken,omitempty"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
	// Joined fields
	UserName            string     `json:"user_name,omitempty"`
	ReviewerName        string     `json:"reviewer_name,omitempty"`
	ProductTitle        string     `json:"product_title,omitempty"`
}

// AbuseReportFilter for filtering abuse reports
type AbuseReportFilter struct {
	Status         string `json:"status"`
	Priority       string `json:"priority"`
	Reason         string `json:"reason"`
	ReportedUserID int    `json:"reported_user_id"`
	Page           int    `json:"page"`
	Limit          int    `json:"limit"`
}

// RefundRequestFilter for filtering refund requests
type RefundRequestFilter struct {
	Status        string `json:"status"`
	RequesterRole string `json:"requester_role"`
	Page          int    `json:"page"`
	Limit         int    `json:"limit"`
}

// AccountAppealFilter for filtering appeals
type AccountAppealFilter struct {
	Status     string `json:"status"`
	AppealType string `json:"appeal_type"`
	UserID     int    `json:"user_id"`
	Page       int    `json:"page"`
	Limit      int    `json:"limit"`
}
