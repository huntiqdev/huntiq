package repository

import "time"

// ReferralStats contains referral dashboard data
type ReferralStats struct {
	ReferralCode          string
	ReferralLink          string
	ReferralSuspended     bool
	InvitedCount          int
	TotalCommission       float64
	InvitationAcceptedAt  *time.Time
	InviterCode           string
	InviterName           string
}

// CommissionEntry represents a single commission record
type CommissionEntry struct {
	ID                int
	UserID            int
	SourceUserID      int
	SourceUserName    string
	InvitedUserID     int
	InvitedUserName   string
	TransactionType   string
	TransactionAmount float64
	Amount            float64
	CommissionRate    float64
	CommissionAmount  float64
	CreatedAt         time.Time
}

// WithdrawalRequest represents a referral withdrawal request
type WithdrawalRequest struct {
	ID            int
	UserID        int
	UserName      string
	UserPhone     string
	Amount        float64
	Status        string
	PaymentMethod string
	PaymentPhone  string
	PayerPhone    string
	Reference     string
	RejectionNote string
	AdminNote     string
	RequestedAt   time.Time
	ProcessedAt   *time.Time
	ReviewedAt    *time.Time
	ReviewedBy    *int
	PaidAt        *time.Time
}

// OTPRecord represents an OTP log entry
type OTPRecord struct {
	ID        int
	Phone     string
	Purpose   string
	Role      string
	SentAt    time.Time
	ExpiresAt time.Time
	Verified  bool
}

// ActivityLog represents a user activity log entry
type ActivityLog struct {
	ID         int
	Event      string
	Phone      string
	Email      string
	UserType   string
	IPAddress  string
	Meta       string
	CreatedAt  time.Time
}

// ActivityLogEntry is an alias for ActivityLog (used in admin panel)
type ActivityLogEntry = ActivityLog

// ActivityLogFilter for querying activity logs
type ActivityLogFilter struct {
	UserID    *int
	Action    string
	Event     string
	Phone     string
	UserType  string
	DateFrom  *time.Time
	DateTo    *time.Time
	StartDate *time.Time
	EndDate   *time.Time
	Page      int
	Limit     int
	Offset    int
}

// ReferralBalance represents commission balance info
type ReferralBalance struct {
	AvailableBalance float64
	PendingAmount    float64
	TotalEarned      float64
}

// InvitationRequest represents a pending invitation
type InvitationRequest struct {
	ID                  int
	UserID              int
	InviterCode         string
	InviterCodeSubmitted string
	Status              string
	RejectReason        string
	RequestedAt         time.Time
	CreatedAt           time.Time
	ProcessedAt         *time.Time
}

// UserBalance for referral balance
type UserBalance struct {
	AvailableBalance float64
	PendingAmount    float64
	TotalEarned      float64
	TotalWithdrawn   float64
}

// OTPFilter for filtering OTP logs
type OTPFilter struct {
	Phone     string
	Purpose   string
	Role      string
	StartDate *time.Time
	EndDate   *time.Time
	Limit     int
	Offset    int
}

// OTPLogFilter for admin OTP log queries
type OTPLogFilter struct {
	Phone        string
	Email        string
	UserType     string
	Status       string
	ShowPlainOTP bool
	Page         int
	Limit        int
}

// OTPLogEntry represents an OTP log entry for admin panel
type OTPLogEntry struct {
	ID            int
	Phone         string
	Email         string
	Purpose       string
	UserType      string
	Used          bool
	Attempts      int
	ExpiresAt     time.Time
	CreatedAt     time.Time
	Status        string
	EncryptedCode string
	PlainOTP      string
}

// AdminReferralRow represents referral data in admin list
type AdminReferralRow struct {
	InvitedUserID    int
	InvitedUserName  string
	InvitedUserPhone string
	InvitedUserRole  string
	InviterID        int
	InviterName      string
	InviterCode      string
	InviterSuspended bool
	AcceptedAt       *time.Time
	TotalCommission  float64
}

// Repository type aliases for backwards compatibility
type UserRepository = *UserRepo
type AdminRepository = *AdminRepo
type ReferralRepository = *ReferralRepo
type ProductRepository = *ProductRepo
type CategoryRepository = *CategoryRepo
type LocationRepository = *LocationRepo
type TransactionRepository = *TransactionRepo
type BannerRepository = *BannerRepo
