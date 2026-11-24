package payment

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/your-org/used-tools-marketplace/internal/models"
)

type TelebirrAdapterConfig struct {
	AppID       string
	AppKey      string
	ShortCode   string
	APIEndpoint string
	NotifyURL   string
	ReturnURL   string
	IsProd      bool
}

type TelebirrAdapter struct {
	config TelebirrAdapterConfig
}

func NewTelebirrAdapter(cfg TelebirrAdapterConfig) *TelebirrAdapter {
	return &TelebirrAdapter{config: cfg}
}

func (a *TelebirrAdapter) ProviderCode() string {
	return "telebirr"
}

func (a *TelebirrAdapter) ProviderName() string {
	return "Ethio Telecom Telebirr SuperApp"
}

func (a *TelebirrAdapter) HealthCheck(ctx context.Context) error {
	if a.config.IsProd {
		if a.config.AppID == "" || a.config.AppKey == "" || a.config.ShortCode == "" {
			return errors.New("telebirr production credentials incomplete (FAIL-CLOSED)")
		}
	}
	return nil
}

func (a *TelebirrAdapter) ComputeHMACSignature(payload []byte) string {
	h := hmac.New(sha256.New, []byte(a.config.AppKey))
	h.Write(payload)
	return hex.EncodeToString(h.Sum(nil))
}

func (a *TelebirrAdapter) InitiatePayment(ctx context.Context, order *models.PaymentOrder) (*InitiateResult, error) {
	if a.config.IsProd {
		if a.config.AppID == "" || a.config.AppKey == "" || a.config.ShortCode == "" {
			return nil, errors.New("telebirr provider is not configured with production merchant credentials")
		}
	}

	checkoutURL := fmt.Sprintf("%s/checkout?out_trade_no=%s&amount=%.2f&short_code=%s",
		a.config.APIEndpoint, order.OrderID, order.Amount, a.config.ShortCode)

	return &InitiateResult{
		ProviderTxID: fmt.Sprintf("TB-PRE-%s", order.OrderID),
		CheckoutURL:  checkoutURL,
		Instructions: "Pay securely via your Telebirr mobile wallet app or web checkout.",
		ExtraData: map[string]interface{}{
			"short_code": a.config.ShortCode,
			"app_id":     a.config.AppID,
		},
	}, nil
}

func (a *TelebirrAdapter) VerifyPayment(ctx context.Context, order *models.PaymentOrder) (*VerificationResult, error) {
	if a.config.IsProd && (a.config.AppID == "" || a.config.AppKey == "") {
		return nil, errors.New("cannot verify transaction: telebirr production credentials missing")
	}

	// In real production, calls Telebirr Query API
	// If order already marked successful, confirm status
	if order.Status == models.PaymentStatusSuccessful {
		return &VerificationResult{
			Status:       models.PaymentStatusSuccessful,
			ProviderTxID: order.OrderID,
			AmountPaid:   order.Amount,
			Currency:     order.Currency,
		}, nil
	}

	return &VerificationResult{
		Status:       models.PaymentStatusPending,
		ProviderTxID: order.OrderID,
		AmountPaid:   order.Amount,
		Currency:     order.Currency,
	}, nil
}

func (a *TelebirrAdapter) ProcessWebhook(ctx context.Context, rawBody []byte, signature string) (*WebhookResult, error) {
	var payload struct {
		OutTradeNo  string  `json:"out_trade_no"`
		TradeNo     string  `json:"trade_no"`
		TradeStatus string  `json:"trade_status"`
		TotalAmount float64 `json:"total_amount"`
		PayerPhone  string  `json:"payer_phone"`
		Signature   string  `json:"signature"`
	}

	if err := json.Unmarshal(rawBody, &payload); err != nil {
		return nil, fmt.Errorf("malformed telebirr webhook payload: %w", err)
	}

	if payload.OutTradeNo == "" {
		return nil, errors.New("missing out_trade_no in webhook payload")
	}

	sigToVerify := signature
	if sigToVerify == "" {
		sigToVerify = payload.Signature
	}

	// Verify cryptographic signature
	if a.config.IsProd && a.config.AppKey == "" {
		return nil, errors.New("telebirr app key unconfigured in production: refusing unverified webhook")
	}
	if a.config.AppKey != "" {
		expectedSig := a.ComputeHMACSignature(rawBody)
		if !hmac.Equal([]byte(sigToVerify), []byte(expectedSig)) {
			return nil, errors.New("invalid telebirr cryptographic webhook signature")
		}
	}

	status := models.PaymentStatusFailed
	isSuccess := strings.EqualFold(payload.TradeStatus, "COMPLETED") ||
		strings.EqualFold(payload.TradeStatus, "SUCCESS") ||
		strings.EqualFold(payload.TradeStatus, "PAID")
	if isSuccess {
		status = models.PaymentStatusSuccessful
	}

	var rawMap map[string]interface{}
	_ = json.Unmarshal(rawBody, &rawMap)

	return &WebhookResult{
		OrderID:      payload.OutTradeNo,
		ProviderTxID: payload.TradeNo,
		Status:       status,
		PayerPhone:   payload.PayerPhone,
		AmountPaid:   payload.TotalAmount,
		Currency:     "ETB",
		RawPayload:   rawMap,
	}, nil
}
