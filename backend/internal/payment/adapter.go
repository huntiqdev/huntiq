package payment

import (
	"context"

	"github.com/your-org/used-tools-marketplace/internal/models"
)

// InitiateResult contains output of payment creation at the provider
type InitiateResult struct {
	ProviderTxID string                 `json:"provider_tx_id"`
	CheckoutURL  string                 `json:"checkout_url"`
	Instructions string                 `json:"instructions,omitempty"`
	ExtraData    map[string]interface{} `json:"extra_data,omitempty"`
}

// VerificationResult contains authoritative verification status from provider
type VerificationResult struct {
	Status        string                 `json:"status"` // SUCCESSFUL, FAILED, CANCELLED, EXPIRED, PENDING
	ProviderTxID  string                 `json:"provider_tx_id"`
	PayerPhone    string                 `json:"payer_phone"`
	AmountPaid    float64                `json:"amount_paid"`
	Currency      string                 `json:"currency"`
	FailureReason string                 `json:"failure_reason,omitempty"`
	RawResponse   map[string]interface{} `json:"raw_response,omitempty"`
}

// WebhookResult contains parsed and verified inbound webhook callback
type WebhookResult struct {
	OrderID       string                 `json:"order_id"`
	ProviderTxID  string                 `json:"provider_tx_id"`
	Status        string                 `json:"status"`
	PayerPhone    string                 `json:"payer_phone"`
	AmountPaid    float64                `json:"amount_paid"`
	Currency      string                 `json:"currency"`
	FailureReason string                 `json:"failure_reason,omitempty"`
	RawPayload    map[string]interface{} `json:"raw_payload"`
}

// Provider defines the standard contract for any payment provider adapter
type Provider interface {
	ProviderCode() string
	ProviderName() string
	InitiatePayment(ctx context.Context, order *models.PaymentOrder) (*InitiateResult, error)
	VerifyPayment(ctx context.Context, order *models.PaymentOrder) (*VerificationResult, error)
	ProcessWebhook(ctx context.Context, rawBody []byte, signature string) (*WebhookResult, error)
	HealthCheck(ctx context.Context) error
}
