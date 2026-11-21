package models

import (
	"time"
)

// Administrative Enterprise Role Constants — 30 Enterprise Roles
const (
	// Department 1: Executive & Strategic Governance
	RolePlatformOwner        = "platform_owner"
	RoleGovComplianceOfficer = "gov_compliance_officer"
	RoleLegalCounsel         = "legal_counsel"
	RoleHRDirector           = "hr_director"

	// Department 2: Cybersecurity, SOC & Trust & Safety
	RoleCISO               = "ciso"
	RoleSecurityAdmin      = "security_admin"
	RoleSOCAnalyst         = "soc_analyst"
	RoleFraudInvestigator  = "fraud_investigator"
	RoleComplianceAuditor  = "compliance_auditor"

	// Department 3: Banking, Treasury, Tax & Financial Operations
	RoleFinanceDirector     = "finance_director"
	RoleTransactionAuditor  = "transaction_auditor"
	RoleTaxController       = "tax_controller"
	RoleEscrowReconciler    = "escrow_reconciler"
	RoleFinanceSupportAgent = "finance_support_agent"

	// Department 4: VIP Concierge, Premium & Growth Services
	RoleVIPAccountDirector = "vip_account_director"
	RoleGoldConcierge      = "gold_concierge"
	RolePremiumConcierge   = "premium_concierge"
	RoleReferralDirector   = "referral_director"
	RoleRetentionSpecialist = "retention_specialist"

	// Department 5: Marketplace Operations & KYC Verification
	RoleOperationsAdmin        = "operations_admin"
	RoleKYCVerificationOfficer = "kyc_verification_officer"
	RoleProductInspectionOfficer = "product_inspection_officer"
	RoleDisputeRefundOfficer   = "dispute_refund_officer"
	RoleAppealsOfficer         = "appeals_officer"

	// Department 6: Catalog, Taxonomy, Logistics & Support
	RolePlatformAdmin        = "platform_admin"
	RoleTaxonomyArchitect    = "taxonomy_architect"
	RoleGISLocationSpecialist = "gis_location_specialist"
	RoleContentModerator     = "content_moderator"
	RoleModerator            = "moderator" // alias for content_moderator
	RoleCustomerSupportLead  = "customer_support_lead"
	RoleSupportAgent         = "support_agent"

	// Legacy backward-compatibility constants
	LegacyRoleSuper  = "super"
	LegacyRoleNormal = "normal"
)

// PasskeyCredential represents a registered WebAuthn / FIDO2 credential
type PasskeyCredential struct {
	ID              int        `json:"id"`
	UserID          int        `json:"user_id"`
	CredentialID    []byte     `json:"-"`
	CredentialIDB64 string     `json:"credential_id"` // Base64URL encoded for JSON
	PublicKey       []byte     `json:"-"`
	AttestationType string     `json:"attestation_type"`
	AAGUID          []byte     `json:"-"`
	SignCount       uint32     `json:"sign_count"`
	CloneWarning    bool       `json:"clone_warning"`
	Attachment      string     `json:"attachment"` // platform or cross-platform
	Name            string     `json:"name"`
	Transports      []string   `json:"transports"`
	BackupEligible  bool       `json:"backup_eligible"`
	BackupState     bool       `json:"backup_state"`
	CreatedAt       time.Time  `json:"created_at"`
	LastUsedAt      *time.Time `json:"last_used_at,omitempty"`
	RevokedAt       *time.Time `json:"revoked_at,omitempty"`
}

// AuthChallenge represents a server-issued single-use WebAuthn or reauth challenge
type AuthChallenge struct {
	ID           int        `json:"id"`
	Challenge    string     `json:"challenge"`
	UserID       *int       `json:"user_id,omitempty"`
	Portal       string     `json:"portal"`
	CeremonyType string     `json:"ceremony_type"` // registration, authentication, reauth
	SessionData  string     `json:"-"`             // JSON serialized WebAuthn session data
	ExpiresAt    time.Time  `json:"expires_at"`
	UsedAt       *time.Time `json:"used_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
}

// AdminSession represents an active server-tracked administrative session
type AdminSession struct {
	ID                 int        `json:"id"`
	SessionTokenHash   string     `json:"-"`
	UserID             int        `json:"user_id"`
	Portal             string     `json:"portal"`
	Role               string     `json:"role"`
	DeviceID           string     `json:"device_id,omitempty"`
	IPAddress          string     `json:"ip_address"`
	UserAgent          string     `json:"user_agent"`
	LastActiveAt       time.Time  `json:"last_active_at"`
	IdleTimeoutSeconds int        `json:"idle_timeout_seconds"`
	AbsoluteExpiresAt  time.Time  `json:"absolute_expires_at"`
	RevokedAt          *time.Time `json:"revoked_at,omitempty"`
	RevocationReason   string     `json:"revocation_reason,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
}

// AdminDeviceRegistration tracks authorized hardware computers/workstations per administrator
type AdminDeviceRegistration struct {
	ID               int        `json:"id"`
	DeviceID         string     `json:"device_id"`
	UserID           int        `json:"user_id"`
	DeviceName       string     `json:"device_name"`
	Browser          string     `json:"browser"`
	OS               string     `json:"os"`
	DeviceType       string     `json:"device_type"`
	Status           string     `json:"status"` // active, pending_verification, revoked, suspended
	IsTrusted        bool       `json:"is_trusted"`
	LastIP           string     `json:"last_ip"`
	LastActiveAt     time.Time  `json:"last_active_at"`
	RegisteredAt     time.Time  `json:"registered_at"`
	RevokedAt        *time.Time `json:"revoked_at,omitempty"`
	RevocationReason string     `json:"revocation_reason,omitempty"`
}

// AuthAttempt tracks failed and successful logins for progressive rate limiting & abuse detection
type AuthAttempt struct {
	ID            int       `json:"id"`
	Identifier    string    `json:"identifier"`
	IPAddress     string    `json:"ip_address"`
	Portal        string    `json:"portal"`
	Success       bool      `json:"success"`
	FailureReason string    `json:"failure_reason,omitempty"`
	AttemptedAt   time.Time `json:"attempted_at"`
}

// ReauthToken tracks high-assurance re-authentication tokens for sensitive actions
type ReauthToken struct {
	ID          int        `json:"id"`
	UserID      int        `json:"user_id"`
	TokenHash   string     `json:"-"`
	ActionScope string     `json:"action_scope"`
	ExpiresAt   time.Time  `json:"expires_at"`
	UsedAt      *time.Time `json:"used_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}
