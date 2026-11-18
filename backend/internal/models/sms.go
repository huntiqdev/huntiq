package models

import "time"

// SMSProviderConfig represents a configured SMS gateway provider.
type SMSProviderConfig struct {
	ID           int                    `json:"id"`
	ProviderCode string                 `json:"provider_code"`
	ProviderName string                 `json:"provider_name"`
	APIEndpoint  string                 `json:"api_endpoint"`
	SenderID     string                 `json:"sender_id"`
	Credentials  map[string]interface{} `json:"credentials"`
	IsActive     bool                   `json:"is_active"`
	IsDefault    bool                   `json:"is_default"`
	DailyLimit   int                    `json:"daily_limit"`
	DailySent    int                    `json:"daily_sent"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
}

// SMSDeliveryLog captures audit records of SMS attempts with masked identifiers.
type SMSDeliveryLog struct {
	ID            int64     `json:"id"`
	PhoneMasked   string    `json:"phone_masked"`
	Purpose       string    `json:"purpose"`
	Provider      string    `json:"provider"`
	ProviderRef   *string   `json:"provider_ref,omitempty"`
	Status        string    `json:"status"` // queued, delivered, failed, rejected
	ErrorCategory *string   `json:"error_category,omitempty"`
	RetryCount    int       `json:"retry_count"`
	CreatedAt     time.Time `json:"created_at"`
}

