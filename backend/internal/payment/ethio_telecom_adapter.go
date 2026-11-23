package payment

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/your-org/used-tools-marketplace/internal/models"
)

type EthioTelecomAdapterConfig struct {
	MerchantID  string
	APIKey      string
	APIEndpoint string
	NotifyURL   string
	ReturnURL   string
	IsProd      bool
}

type EthioTelecomAdapter struct {
	config EthioTelecomAdapterConfig
}

func NewEthioTelecomAdapter(cfg EthioTelecomAdapterConfig) *EthioTelecomAdapter {
	return &EthioTelecomAdapter{config: cfg}
}

func (a *EthioTelecomAdapter) ProviderCode() string {
	return "ethio_telecom"
}

func (a *EthioTelecomAdapter) ProviderName() string {
	return "Ethio Telecom TelePay / Direct Carrier Billing"
}

func (a *EthioTelecomAdapter) HealthCheck(ctx context.Context) error {
	if a.config.IsProd && (a.config.MerchantID == "" || a.config.APIKey == "") {
		return errors.New("ethio telecom production merchant credentials incomplete (FAIL-CLOSED)")
	}
	return nil
}

func (a *EthioTelecomAdapter) InitiatePayment(ctx context.Context, order *models.PaymentOrder) (*InitiateResult, error) {
	if a.config.IsProd {
		if a.config.MerchantID == "" || a.config.APIKey == "" {
			return nil, errors.New("ethio telecom payment provider is not configured with production credentials")
		}
	}

	checkoutURL := fmt.Sprintf("%s/checkout?merchant_id=%s&reference=%s&amount=%.2f",
		a.config.APIEndpoint, a.config.MerchantID, order.OrderID, order.Amount)

	return &InitiateResult{
		ProviderTxID: fmt.Sprintf("ET-PRE-%s", order.OrderID),
		CheckoutURL:  checkoutURL,
		Instructions: "Approve transaction on your Ethio Telecom mobile SIM prompt or telepay checkout.",
		ExtraData: map[string]interface{}{
			"merchant_id": a.config.MerchantID,
		},
	}, nil
}

func (a *EthioTelecomAdapter) VerifyPayment(ctx context.Context, order *models.PaymentOrder) (*VerificationResult, error) {
	if a.config.IsProd && a.config.APIKey == "" {
		return nil, errors.New("cannot verify transaction: ethio telecom production credentials missing")
	}

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

func (a *EthioTelecomAdapter) ProcessWebhook(ctx context.Context, rawBody []byte, signature string) (*WebhookResult, error) {
	var payload struct {
		TransactionID string  `json:"transaction_id"`
		Reference     string  `json:"reference"`
		Status        string  `json:"status"` // SUCCESS, FAILED
		Amount        float64 `json:"amount"`
		MSISDN        string  `json:"msisdn"`
		AuthSignature string  `json:"signature"`
	}

	if err := json.Unmarshal(rawBody, &payload); err != nil {
		return nil, fmt.Errorf("malformed ethio telecom webhook payload: %w", err)
	}

	if payload.Reference == "" {
		return nil, errors.New("missing reference in ethio telecom webhook")
	}

	status := models.PaymentStatusFailed
	if strings.EqualFold(payload.Status, "SUCCESS") || strings.EqualFold(payload.Status, "COMPLETED") {
		status = models.PaymentStatusSuccessful
	}

	var rawMap map[string]interface{}
	_ = json.Unmarshal(rawBody, &rawMap)

	return &WebhookResult{
		OrderID:      payload.Reference,
		ProviderTxID: payload.TransactionID,
		Status:       status,
		PayerPhone:   payload.MSISDN,
		AmountPaid:   payload.Amount,
		Currency:     "ETB",
		RawPayload:   rawMap,
	}, nil
}
