package service

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/your-org/used-tools-marketplace/internal/models"
	"github.com/your-org/used-tools-marketplace/internal/repository"
	"github.com/your-org/used-tools-marketplace/pkg/fees"
)

type TelebirrConfig struct {
	AppID       string
	AppKey      string
	ShortCode   string
	APIEndpoint string
	NotifyURL   string
	ReturnURL   string
}

type TelebirrService struct {
	repo         *repository.TelebirrRepo
	tokenSvc     *ReferralTokenService
	ledgerSvc    *LedgerService
	config       *TelebirrConfig
	devMode      bool
}

func NewTelebirrService(
	repo *repository.TelebirrRepo,
	tokenSvc *ReferralTokenService,
	ledgerSvc *LedgerService,
	cfg *TelebirrConfig,
	devMode bool,
) *TelebirrService {
	if cfg == nil {
		cfg = &TelebirrConfig{
			AppID:       "HUNTIQ_TELEBIRR_APP",
			AppKey:      "dev_telebirr_secret_key_2026",
			ShortCode:   "100200",
			APIEndpoint: "https://app.telebirr.et/payment",
			NotifyURL:   "https://huntiq.et/api/telebirr/webhook",
			ReturnURL:   "https://huntiq.et/payment/callback",
		}
	}
	return &TelebirrService{
		repo:      repo,
		tokenSvc:  tokenSvc,
		ledgerSvc: ledgerSvc,
		config:    cfg,
		devMode:   devMode,
	}
}

// generateOutTradeNo creates an unambiguous merchant reference ID.
func generateOutTradeNo() string {
	b := make([]byte, 4)
	_, _ = rand.Read(b)
	return fmt.Sprintf("TB-%d-%s", time.Now().Unix(), hex.EncodeToString(b))
}

// ComputeHMACSignature computes the HMAC-SHA256 signature for Telebirr payload verification.
func (s *TelebirrService) ComputeHMACSignature(payload []byte) string {
	h := hmac.New(sha256.New, []byte(s.config.AppKey))
	h.Write(payload)
	return hex.EncodeToString(h.Sum(nil))
}

// CreatePaymentOrder initiates a Telebirr order and prepares payment details.
func (s *TelebirrService) CreatePaymentOrder(ctx context.Context, userID int, amount float64, orderType, title string) (*models.TelebirrOrder, string, error) {
	if userID <= 0 || amount <= 0 {
		return nil, "", errors.New("invalid user ID or payment amount")
	}

	cleanAmount := roundCents(amount)
	outTradeNo := generateOutTradeNo()

	order := &models.TelebirrOrder{
		OutTradeNo: outTradeNo,
		UserID:     userID,
		Amount:     cleanAmount,
		Currency:   "ETB",
		Title:      title,
		OrderType:  orderType,
		Status:     "pending",
	}

	if err := s.repo.CreateOrder(ctx, order); err != nil {
		return nil, "", fmt.Errorf("failed to create Telebirr order: %w", err)
	}

	// Prepare checkout payload or redirection URL
	checkoutURL := fmt.Sprintf("%s/checkout?out_trade_no=%s&amount=%.2f&short_code=%s",
		s.config.APIEndpoint, outTradeNo, cleanAmount, s.config.ShortCode)

	return order, checkoutURL, nil
}

// TelebirrWebhookPayload represents inbound payment notifications from Telebirr.
type TelebirrWebhookPayload struct {
	OutTradeNo  string  `json:"out_trade_no"`
	TradeNo     string  `json:"trade_no"`
	TradeStatus string  `json:"trade_status"` // COMPLETED, SUCCESS, FAILED
	TotalAmount float64 `json:"total_amount"`
	PayerPhone  string  `json:"payer_phone"`
	Signature   string  `json:"signature"`
}

// HandleWebhook validates signature, updates order status, and triggers Reward B commission.
func (s *TelebirrService) HandleWebhook(ctx context.Context, rawBody []byte, signatureHeader string) (*models.TelebirrOrder, error) {
	var payload TelebirrWebhookPayload
	if err := json.Unmarshal(rawBody, &payload); err != nil {
		return nil, fmt.Errorf("malformed webhook json: %w", err)
	}

	if payload.OutTradeNo == "" {
		return nil, errors.New("missing out_trade_no in webhook payload")
	}

	// Verify signature
	sigToVerify := signatureHeader
	if sigToVerify == "" {
		sigToVerify = payload.Signature
	}

	if !s.devMode && s.config.AppKey != "" {
		expectedSig := s.ComputeHMACSignature(rawBody)
		if !hmac.Equal([]byte(sigToVerify), []byte(expectedSig)) {
			return nil, errors.New("invalid webhook signature")
		}
	}

	// Retrieve existing order
	order, err := s.repo.GetOrderByOutTradeNo(ctx, payload.OutTradeNo)
	if err != nil {
		return nil, fmt.Errorf("order not found: %s", payload.OutTradeNo)
	}

	// Check if already completed (idempotency)
	if order.Status == "completed" {
		return order, nil
	}

	isSuccess := strings.EqualFold(payload.TradeStatus, "COMPLETED") ||
		strings.EqualFold(payload.TradeStatus, "SUCCESS") ||
		strings.EqualFold(payload.TradeStatus, "PAID")

	// Verify that the amount paid matches the registered order amount
	if isSuccess && !fees.AmountsMatch(payload.TotalAmount, order.Amount) {
		return nil, fmt.Errorf("telebirr amount mismatch: payload %.2f != order %.2f", payload.TotalAmount, order.Amount)
	}

	newStatus := "failed"
	if isSuccess {
		newStatus = "completed"
	}

	var rawMap map[string]interface{}
	_ = json.Unmarshal(rawBody, &rawMap)

	var txNoPtr, phonePtr *string
	if payload.TradeNo != "" {
		txNoPtr = &payload.TradeNo
	}
	if payload.PayerPhone != "" {
		phonePtr = &payload.PayerPhone
	}

	if err := s.repo.UpdateOrderStatus(ctx, order.OutTradeNo, newStatus, txNoPtr, phonePtr, rawMap); err != nil {
		return nil, fmt.Errorf("failed to update order status: %w", err)
	}

	order.Status = newStatus
	order.TransactionNo = txNoPtr
	order.PayerPhone = phonePtr

	// If successfully completed, credit Reward B (Qualifying Transaction Reward) to the inviter
	if isSuccess && s.tokenSvc != nil && s.ledgerSvc != nil {
		inviterID, err := s.tokenSvc.GetUserInviter(ctx, order.UserID)
		if err == nil && inviterID > 0 && inviterID != order.UserID {
			log.Printf("[TelebirrService] Qualifying transaction completed for user %d (Order %s, %.2f ETB). Crediting Reward B to inviter %d",
				order.UserID, order.OutTradeNo, order.Amount, inviterID)
			_ = s.ledgerSvc.CreditTransactionCommission(ctx, inviterID, order.UserID, order.ID, order.Amount)
		}
	}

	return order, nil
}

// GetOrderStatus fetches current state of an order.
func (s *TelebirrService) GetOrderStatus(ctx context.Context, outTradeNo string) (*models.TelebirrOrder, error) {
	return s.repo.GetOrderByOutTradeNo(ctx, outTradeNo)
}

// GetUserOrders returns paginated Telebirr orders for a user.
func (s *TelebirrService) GetUserOrders(ctx context.Context, userID, limit, offset int) ([]models.TelebirrOrder, int, error) {
	return s.repo.GetOrdersByUser(ctx, userID, limit, offset)
}

