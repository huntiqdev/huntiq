package payment_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/your-org/used-tools-marketplace/internal/models"
	"github.com/your-org/used-tools-marketplace/internal/payment"
)

func TestDemoMockAdapter(t *testing.T) {
	ctx := context.Background()
	adapter := payment.NewDemoMockAdapter("https://huntiq.et")

	order := &models.PaymentOrder{
		OrderID:     "ORD-TEST-12345",
		UserID:      1,
		Amount:      500.0,
		Currency:    "ETB",
		Status:      models.PaymentStatusCreated,
		Environment: models.PaymentEnvDemo,
		IsDemo:      true,
	}

	// 1. Test Initiation
	initRes, err := adapter.InitiatePayment(ctx, order)
	if err != nil {
		t.Fatalf("InitiatePayment failed: %v", err)
	}
	if initRes.CheckoutURL == "" || initRes.ProviderTxID == "" {
		t.Errorf("Expected checkout URL and provider TX ID, got empty")
	}

	// 2. Test Verification Pending
	verRes, err := adapter.VerifyPayment(ctx, order)
	if err != nil {
		t.Fatalf("VerifyPayment failed: %v", err)
	}
	if verRes.Status != models.PaymentStatusPending {
		t.Errorf("Expected pending status, got %s", verRes.Status)
	}

	// 3. Test Verification of Successful Order
	order.Status = models.PaymentStatusSuccessful
	verResSuccess, err := adapter.VerifyPayment(ctx, order)
	if err != nil {
		t.Fatalf("VerifyPayment on successful order failed: %v", err)
	}
	if verResSuccess.Status != models.PaymentStatusSuccessful {
		t.Errorf("Expected successful status, got %s", verResSuccess.Status)
	}

	// 4. Test Webhook
	payload := map[string]interface{}{
		"order_id":       order.OrderID,
		"provider_tx_id": "SIM-WEBHOOK-999",
		"status":         "SUCCESSFUL",
		"amount":         500.0,
		"currency":       "ETB",
	}
	bytes, _ := json.Marshal(payload)
	webhookRes, err := adapter.ProcessWebhook(ctx, bytes, "")
	if err != nil {
		t.Fatalf("ProcessWebhook failed: %v", err)
	}
	if webhookRes.Status != models.PaymentStatusSuccessful || webhookRes.OrderID != order.OrderID {
		t.Errorf("Webhook result mismatch: %+v", webhookRes)
	}
}

func TestTelebirrAdapterFailClosed(t *testing.T) {
	ctx := context.Background()

	// Incomplete production configuration MUST fail closed
	prodAdapter := payment.NewTelebirrAdapter(payment.TelebirrAdapterConfig{
		IsProd: true,
		AppID:  "", // missing credentials
	})

	if err := prodAdapter.HealthCheck(ctx); err == nil {
		t.Errorf("Expected Telebirr health check to fail closed in production with missing credentials")
	}

	order := &models.PaymentOrder{
		OrderID: "ORD-PROD-999",
		Amount:  500.0,
	}
	_, err := prodAdapter.InitiatePayment(ctx, order)
	if err == nil {
		t.Errorf("Expected InitiatePayment to fail closed when credentials missing in production")
	}

	// Valid Signature computation
	configuredAdapter := payment.NewTelebirrAdapter(payment.TelebirrAdapterConfig{
		AppID:     "TEST_APP",
		AppKey:    "SUPER_SECRET_KEY_12345",
		ShortCode: "123456",
		IsProd:    false,
	})

	sig := configuredAdapter.ComputeHMACSignature([]byte(`{"out_trade_no":"ORD-1"}`))
	if len(sig) != 64 { // SHA256 hex length
		t.Errorf("Expected 64-character hex signature, got %s", sig)
	}
}

func TestEthioTelecomAdapterFailClosed(t *testing.T) {
	ctx := context.Background()

	prodAdapter := payment.NewEthioTelecomAdapter(payment.EthioTelecomAdapterConfig{
		IsProd: true,
	})

	if err := prodAdapter.HealthCheck(ctx); err == nil {
		t.Errorf("Expected Ethio Telecom health check to fail closed with missing production credentials")
	}

	order := &models.PaymentOrder{
		OrderID: "ORD-PROD-ET-1",
		Amount:  200.0,
	}
	_, err := prodAdapter.InitiatePayment(ctx, order)
	if err == nil {
		t.Errorf("Expected InitiatePayment to fail closed in production without credentials")
	}
}
