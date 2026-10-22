package models

import "time"

// AdminLog represents a single admin action for audit purposes.
type AdminLog struct {
	ID           int       `json:"id"`
	AdminID      int       `json:"admin_id"`
	AdminName    string    `json:"admin_name,omitempty"`
	ProductID    *int      `json:"product_id,omitempty"`
	TargetUserID *int      `json:"target_user_id,omitempty"`
	TargetName   string    `json:"target_name,omitempty"`
	Action       string    `json:"action"`
	Notes        string    `json:"notes,omitempty"`
	Details      string    `json:"details,omitempty"`
	IPAddress    string    `json:"ip_address,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

// AdminLogFilter holds query parameters for filtering admin logs.
type AdminLogFilter struct {
	AdminID   int
	Action    string
	IPAddress string
	DateFrom  string
	DateTo    string
	Page      int
	Limit     int
}

// TierHistory records a user tier change.
type TierHistory struct {
	ID            int       `json:"id"`
	UserID        int       `json:"user_id"`
	UserName      string    `json:"user_name,omitempty"`
	TierType      string    `json:"tier_type"` // seller | buyer
	OldTier       string    `json:"old_tier"`
	NewTier       string    `json:"new_tier"`
	ChangedBy     *int      `json:"changed_by,omitempty"`
	ChangedByName string    `json:"changed_by_name,omitempty"`
	Reason        string    `json:"reason,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}
