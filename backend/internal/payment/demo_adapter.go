package payment

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/your-org/used-tools-marketplace/internal/models"
)

// DemoMockAdapter provides isolated, controlled simulation of payments for development & testing.
type DemoMockAdapter struct {
	baseURL string
}

func NewDemoMockAdapter(baseURL string) *DemoMockAdapter {
	if baseURL == "" {
		baseURL = "https://huntiq.et"
	}
	return &DemoMockAdapter{baseURL: baseURL}
}

func (a *DemoMockAdapter) ProviderCode() string {
	return "demo_mock"
}

func (a *DemoMockAdapter) ProviderName() string {
	return "HUNTIQ FinTech Demo Simulator"
}

func (a *DemoMockAdapter) HealthCheck(ctx context.Context) error {
	return nil // Always healthy
}

func (a *DemoMockAdapter) InitiatePayment(ctx context.Context, order *models.PaymentOrder) (*InitiateResult, error) {
	mockTxID := fmt.Sprintf("SIM-TX-%d", time.Now().UnixNano())
	checkoutURL := fmt.Sprintf("%s/payment/demo-checkout?order_id=%s&amount=%.2f&currency=%s",
		a.baseURL, order.OrderID, order.Amount, order.Currency)

	return &InitiateResult{
		ProviderTxID: mockTxID,
		CheckoutURL:  checkoutURL,
		Instructions: "DEMO MODE: Complete this simulated payment in the developer checkout window.",
		ExtraData: map[string]interface{}{
			"is_simulated": true,
			"environment":  "demo",
		},
	}, nil
}

func (a *DemoMockAdapter) VerifyPayment(ctx context.Context, order *models.PaymentOrder) (*VerificationResult, error) {
	// If the order already has metadata indicating simulation outcome
	if order.Status == models.PaymentStatusSuccessful {
		return &VerificationResult{
			Status:       models.PaymentStatusSuccessful,
			ProviderTxID: fmt.Sprintf("SIM-VERIFIED-%s", order.OrderID),
			PayerPhone:   "+251911000000",
			AmountPaid:   order.Amount,
			Currency:     order.Currency,
			RawResponse: map[string]interface{}{
				"simulated": true,
				"verified":  true,
			},
		}, nil
	}

	return &VerificationResult{
		Status:       models.PaymentStatusPending,
		ProviderTxID: fmt.Sprintf("SIM-%s", order.OrderID),
		AmountPaid:   order.Amount,
		Currency:     order.Currency,
		RawResponse: map[string]interface{}{
			"simulated": true,
			"status":    "pending_user_action",
		},
	}, nil
}

func (a *DemoMockAdapter) ProcessWebhook(ctx context.Context, rawBody []byte, signature string) (*WebhookResult, error) {
	var payload struct {
		OrderID       string  `json:"order_id"`
		ProviderTxID  string  `json:"provider_tx_id"`
		Status        string  `json:"status"` // SUCCESSFUL, FAILED, CANCELLED, EXPIRED
		Amount        float64 `json:"amount"`
		Currency      string  `json:"currency"`
		PayerPhone    string  `json:"payer_phone"`
		FailureReason string  `json:"failure_reason"`
	}

	if err := json.Unmarshal(rawBody, &payload); err != nil {
		return nil, fmt.Errorf("malformed demo webhook payload: %w", err)
	}

	if payload.OrderID == "" {
		return nil, errors.New("missing order_id in demo webhook")
	}

	status := strings.ToUpper(payload.Status)
	if status == "" {
		status = models.PaymentStatusSuccessful
	}

	var rawMap map[string]interface{}
	_ = json.Unmarshal(rawBody, &rawMap)

	return &WebhookResult{
		OrderID:       payload.OrderID,
		ProviderTxID:  payload.ProviderTxID,
		Status:        status,
		PayerPhone:    payload.PayerPhone,
		AmountPaid:    payload.Amount,
		Currency:      payload.Currency,
		FailureReason: payload.FailureReason,
		RawPayload:    rawMap,
	}, nil
}

// SimulateOutcome allows tests to drive explicit lifecycle outcomes
func (a *DemoMockAdapter) SimulateOutcome(order *models.PaymentOrder, action string) (*VerificationResult, error) {
	normAction := strings.ToLower(strings.TrimSpace(action))
	mockTx := fmt.Sprintf("SIM-ACT-%d", time.Now().UnixNano())

	switch normAction {
	case "success", "successful":
		return &VerificationResult{
			Status:       models.PaymentStatusSuccessful,
			ProviderTxID: mockTx,
			PayerPhone:   "+251911000000",
			AmountPaid:   order.Amount,
			Currency:     order.Currency,
			RawResponse: map[string]interface{}{
				"simulated_action": "success",
				"timestamp":        time.Now().Format(time.RFC3339),
			},
		}, nil
	case "fail", "failed":
		return &VerificationResult{
			Status:        models.PaymentStatusFailed,
			ProviderTxID:  mockTx,
			AmountPaid:    order.Amount,
			Currency:      order.Currency,
			FailureReason: "Simulated payment failure (insufficient balance or provider decline)",
			RawResponse: map[string]interface{}{
				"simulated_action": "failed",
			},
		}, nil
	case "cancel", "cancelled":
		return &VerificationResult{
			Status:        models.PaymentStatusCancelled,
			ProviderTxID:  mockTx,
			AmountPaid:    order.Amount,
			Currency:      order.Currency,
			FailureReason: "User cancelled simulated payment",
			RawResponse: map[string]interface{}{
				"simulated_action": "cancelled",
			},
		}, nil
	case "expire", "expired":
		return &VerificationResult{
			Status:        models.PaymentStatusExpired,
			ProviderTxID:  mockTx,
			AmountPaid:    order.Amount,
			Currency:      order.Currency,
			FailureReason: "Payment transaction expired (exceeded 30-minute window)",
			RawResponse: map[string]interface{}{
				"simulated_action": "expired",
			},
		}, nil
	case "timeout":
		return nil, errors.New("simulated network connection timeout to payment provider gateway")
	default:
		return nil, fmt.Errorf("unsupported simulation action: %s", action)
	}
}
