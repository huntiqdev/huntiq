package service

import (
	"context"
	"strings"
	"testing"
)

func TestTelebirrSignatureAndOutTradeNo(t *testing.T) {
	cfg := &TelebirrConfig{
		AppID:       "TEST_APP",
		AppKey:      "super_secret_test_key_123",
		ShortCode:   "100200",
		APIEndpoint: "https://test.telebirr.et",
	}

	svc := NewTelebirrService(nil, nil, nil, cfg, false)

	payload := []byte(`{"out_trade_no":"TB-123456","trade_status":"COMPLETED"}`)
	sig := svc.ComputeHMACSignature(payload)
	if len(sig) != 64 {
		t.Errorf("HMAC-SHA256 signature length = %d, want 64", len(sig))
	}

	// Verify idempotency of signature calculation
	sig2 := svc.ComputeHMACSignature(payload)
	if sig != sig2 {
		t.Errorf("Signature computation should be deterministic")
	}

	// Verify distinct payloads produce distinct signatures
	sigDiff := svc.ComputeHMACSignature([]byte(`{"out_trade_no":"TB-999999"}`))
	if sig == sigDiff {
		t.Errorf("Different payloads should produce different signatures")
	}

	outTradeNo := generateOutTradeNo()
	if !strings.HasPrefix(outTradeNo, "TB-") {
		t.Errorf("outTradeNo prefix mismatch: %s", outTradeNo)
	}
}

func TestTelebirrWebhookPayloadParsing(t *testing.T) {
	cfg := &TelebirrConfig{
		AppID:       "TEST_APP",
		AppKey:      "test_key",
		ShortCode:   "100200",
		APIEndpoint: "https://test.telebirr.et",
	}

	svc := NewTelebirrService(nil, nil, nil, cfg, true)

	// Malformed payload
	_, err := svc.HandleWebhook(context.Background(), []byte(`invalid json`), "")
	if err == nil {
		t.Errorf("Expected error for malformed json")
	}

	// Missing out_trade_no
	_, err = svc.HandleWebhook(context.Background(), []byte(`{"trade_status":"SUCCESS"}`), "")
	if err == nil {
		t.Errorf("Expected error for missing out_trade_no")
	}
}

