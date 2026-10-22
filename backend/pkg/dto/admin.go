package dto

// ReviewProductRequest is the body for PUT /api/admin/products/:id/review
type ReviewProductRequest struct {
	Status string `json:"status"` // approved | rejected
	Notes  string `json:"notes"`
}

// SetSellerActiveRequest is the body for PUT /api/admin/sellers/:id/active
type SetSellerActiveRequest struct {
	Active bool   `json:"active"`
	Reason string `json:"reason"` // required when active=false
}

// CreateNotificationRequest is the body for POST /api/admin/notifications
type CreateNotificationRequest struct {
	Title      string `json:"title"`
	Message    string `json:"message"`
	Type       string `json:"type"`        // info | warning | success | error
	TargetRole string `json:"target_role"` // "" = all | "super" | "normal"
}

// UpdateSettingsWithReasonRequest is the body for PUT /api/admin/settings/with-reason
type UpdateSettingsWithReasonRequest struct {
	Settings map[string]string `json:"settings"`
	Reason   string            `json:"reason"`
}

// CreateAdminRequest is the body for POST /api/admin/admins
type CreateAdminRequest struct {
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UpdateAdminRoleRequest is the body for PUT /api/admin/admins/:id/role
type UpdateAdminRoleRequest struct {
	AdminRole string `json:"admin_role"` // normal | super
}

// SetAdminPermissionsRequest is the body for PUT /api/admin/admins/:id/permissions
type SetAdminPermissionsRequest struct {
	Permissions []string `json:"permissions"`
}

// SetTierRequest is the body for PUT /api/admin/users/:id/seller-tier and buyer-tier
type SetTierRequest struct {
	Tier   string `json:"tier"`   // normal | gold | premium
	Reason string `json:"reason"` // required audit note
}

// ReasonRequest is a shared body for endpoints that only need a reason string.
// Used by: reject upgrade, reject buyer upgrade, warn user, ban user, reject fee payment.
type ReasonRequest struct {
	Reason string `json:"reason"`
}

// SystemHealthRepairRequest is the body for POST /api/admin/system/health/repair
type SystemHealthRepairRequest struct {
	Action             string         `json:"action"`
	MissingAdminEmails map[int]string `json:"missing_admin_emails"`
}
