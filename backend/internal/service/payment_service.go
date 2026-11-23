package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/your-org/used-tools-marketplace/internal/models"
	"github.com/your-org/used-tools-marketplace/internal/payment"
	"github.com/your-org/used-tools-marketplace/internal/repository"
)

type PaymentService struct {
	paymentRepo *repository.PaymentRepo
	userRepo    *repository.UserRepo
	smsManager  *SMSManager
	baseURL     string
}

func NewPaymentService(
	paymentRepo *repository.PaymentRepo,
	userRepo *repository.UserRepo,
	smsManager *SMSManager,
	baseURL string,
) *PaymentService {
	if baseURL == "" {
		baseURL = "http://localhost:3000"
	}
	return &PaymentService{
		paymentRepo: paymentRepo,
		userRepo:    userRepo,
		smsManager:  smsManager,
		baseURL:     baseURL,
	}
}

func (s *PaymentService) generateOrderID(prefix string) string {
	bytes := make([]byte, 4)
	_, _ = rand.Read(bytes)
	randSuffix := hex.EncodeToString(bytes)
	return fmt.Sprintf("ORD-%s-%d-%s", strings.ToUpper(prefix), time.Now().Unix(), strings.ToUpper(randSuffix))
}

// resolveAdapter instantiates provider adapter with fresh database credentials
func (s *PaymentService) resolveAdapter(ctx context.Context, providerCode string, isProd bool) (payment.Provider, error) {
	switch providerCode {
	case "demo_mock":
		if isProd {
			return nil, errors.New("demo_mock adapter is strictly forbidden in production mode")
		}
		return payment.NewDemoMockAdapter(s.baseURL), nil

	case "telebirr":
		cfg, err := s.paymentRepo.GetProviderConfigByCode(ctx, "telebirr")
		if err != nil || cfg == nil {
			return nil, errors.New("telebirr provider configuration not found")
		}
		decryptedKey, _ := s.paymentRepo.GetDecryptedAppKey(ctx, "telebirr")
		return payment.NewTelebirrAdapter(payment.TelebirrAdapterConfig{
			AppID:       cfg.AppID,
			AppKey:      decryptedKey,
			ShortCode:   cfg.ShortCode,
			APIEndpoint: cfg.APIEndpoint,
			NotifyURL:   cfg.NotifyURL,
			ReturnURL:   cfg.ReturnURL,
			IsProd:      isProd,
		}), nil

	case "ethio_telecom":
		cfg, err := s.paymentRepo.GetProviderConfigByCode(ctx, "ethio_telecom")
		if err != nil || cfg == nil {
			return nil, errors.New("ethio telecom provider configuration not found")
		}
		decryptedKey, _ := s.paymentRepo.GetDecryptedAppKey(ctx, "ethio_telecom")
		return payment.NewEthioTelecomAdapter(payment.EthioTelecomAdapterConfig{
			MerchantID:  cfg.AppID,
			APIKey:      decryptedKey,
			APIEndpoint: cfg.APIEndpoint,
			NotifyURL:   cfg.NotifyURL,
			ReturnURL:   cfg.ReturnURL,
			IsProd:      isProd,
		}), nil

	default:
		return nil, fmt.Errorf("unsupported payment provider: %s", providerCode)
	}
}

// InitiateUpgradePayment initiates tier upgrade payment with authoritative server pricing
func (s *PaymentService) InitiateUpgradePayment(
	ctx context.Context,
	userID int,
	role string,
	targetTier string,
	durationDays int,
	providerCode string,
	payerPhone string,
	otpVerified bool,
) (*models.PaymentOrder, error) {
	globalMode, err := s.paymentRepo.GetGlobalOperatingMode(ctx)
	if err != nil {
		globalMode = "demo"
	}
	isProd := globalMode == "production"

	// Provider selection validation
	if providerCode == "" {
		if isProd {
			providerCode = "telebirr"
		} else {
			providerCode = "demo_mock"
		}
	}

	if isProd && providerCode == "demo_mock" {
		return nil, errors.New("demo_mock provider is forbidden in production mode")
	}

	prices, err := s.paymentRepo.GetTierPricing(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch authoritative pricing: %w", err)
	}

	targetTierClean := strings.ToLower(strings.TrimSpace(targetTier))
	if targetTierClean != "gold" && targetTierClean != "premium" {
		return nil, errors.New("target tier must be 'gold' or 'premium'")
	}

	if durationDays <= 0 {
		durationDays = 30
	}

	var amount float64
	var orderTitle string
	if role == "buyer" {
		basePrice := prices["buyer_gold_price_etb"]
		if targetTierClean == "premium" {
			basePrice = prices["buyer_premium_price_etb"]
		}
		amount = basePrice * (float64(durationDays) / 30.0)
		orderTitle = fmt.Sprintf("Buyer %s Upgrade (%d Days)", strings.ToUpper(targetTierClean), durationDays)
	} else {
		// Seller upgrade
		if targetTierClean == "premium" {
			amount = prices["seller_premium_price_etb"]
		} else {
			amount = prices["seller_gold_price_etb"]
		}
		orderTitle = fmt.Sprintf("Seller %s Upgrade (%d Days)", strings.ToUpper(targetTierClean), durationDays)
	}

	orderID := s.generateOrderID(role[:1] + "-" + targetTierClean[:1])

	envMode := models.PaymentEnvDemo
	if isProd {
		envMode = models.PaymentEnvProduction
	}

	order := &models.PaymentOrder{
		OrderID:      orderID,
		UserID:       userID,
		Provider:     providerCode,
		Environment:  envMode,
		OrderType:    models.OrderTypeTierUpgrade,
		TargetTier:   &targetTierClean,
		UserRole:     role,
		DurationDays: durationDays,
		Amount:       amount,
		Currency:     "ETB",
		Title:        orderTitle,
		Status:       models.PaymentStatusCreated,
		IsDemo:       !isProd,
		OTPVerified:  otpVerified,
	}

	if payerPhone != "" {
		order.PayerPhone = &payerPhone
	}
	if otpVerified {
		now := time.Now()
		order.OTPVerifiedAt = &now
	}

	// Resolve provider adapter
	adapter, err := s.resolveAdapter(ctx, providerCode, isProd)
	if err != nil {
		return nil, err
	}

	initResult, err := adapter.InitiatePayment(ctx, order)
	if err != nil {
		return nil, fmt.Errorf("provider initiation failed: %w", err)
	}

	order.Status = models.PaymentStatusPending
	order.CheckoutURL = &initResult.CheckoutURL
	order.ProviderTxID = &initResult.ProviderTxID
	order.VerificationMetadata = initResult.ExtraData

	if err := s.paymentRepo.CreateOrder(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to persist payment order: %w", err)
	}

	return order, nil
}

// InitiateListingFeePayment initiates a pay-per-listing fee payment
func (s *PaymentService) InitiateListingFeePayment(
	ctx context.Context,
	userID int,
	productID int,
	providerCode string,
	payerPhone string,
) (*models.PaymentOrder, error) {
	globalMode, err := s.paymentRepo.GetGlobalOperatingMode(ctx)
	if err != nil {
		globalMode = "demo"
	}
	isProd := globalMode == "production"

	if providerCode == "" {
		if isProd {
			providerCode = "telebirr"
		} else {
			providerCode = "demo_mock"
		}
	}

	if isProd && providerCode == "demo_mock" {
		return nil, errors.New("demo_mock provider is forbidden in production mode")
	}

	prices, err := s.paymentRepo.GetTierPricing(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load pricing: %w", err)
	}

	amount := prices["listing_fee_etb"]
	if amount <= 0 {
		amount = 50.0
	}

	orderID := s.generateOrderID(fmt.Sprintf("LST-%d", productID))
	envMode := models.PaymentEnvDemo
	if isProd {
		envMode = models.PaymentEnvProduction
	}

	order := &models.PaymentOrder{
		OrderID:      orderID,
		UserID:       userID,
		Provider:     providerCode,
		Environment:  envMode,
		OrderType:    models.OrderTypeListingFee,
		ProductID:    &productID,
		UserRole:     "seller",
		DurationDays: 30,
		Amount:       amount,
		Currency:     "ETB",
		Title:        fmt.Sprintf("Listing Fee for Product #%d", productID),
		Status:       models.PaymentStatusCreated,
		IsDemo:       !isProd,
		OTPVerified:  true,
	}

	if payerPhone != "" {
		order.PayerPhone = &payerPhone
	}

	adapter, err := s.resolveAdapter(ctx, providerCode, isProd)
	if err != nil {
		return nil, err
	}

	initResult, err := adapter.InitiatePayment(ctx, order)
	if err != nil {
		return nil, fmt.Errorf("provider initiation failed: %w", err)
	}

	order.Status = models.PaymentStatusPending
	order.CheckoutURL = &initResult.CheckoutURL
	order.ProviderTxID = &initResult.ProviderTxID
	order.VerificationMetadata = initResult.ExtraData

	if err := s.paymentRepo.CreateOrder(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to persist listing fee order: %w", err)
	}

	return order, nil
}

// VerifyAndActivatePayment verifies the transaction status and automatically activates membership/fee if successful
func (s *PaymentService) VerifyAndActivatePayment(ctx context.Context, orderID string) (*models.PaymentOrder, error) {
	order, err := s.paymentRepo.GetOrderByOrderID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("order not found: %w", err)
	}

	// Idempotency: If already completed or failed terminal, return immediately
	if order.Status == models.PaymentStatusSuccessful {
		return order, nil
	}

	isProd := order.Environment == models.PaymentEnvProduction
	adapter, err := s.resolveAdapter(ctx, order.Provider, isProd)
	if err != nil {
		return nil, err
	}

	verResult, err := adapter.VerifyPayment(ctx, order)
	if err != nil {
		return nil, fmt.Errorf("verification call failed: %w", err)
	}

	if verResult.Status == models.PaymentStatusSuccessful {
		now := time.Now()
		order.Status = models.PaymentStatusSuccessful
		order.VerifiedAt = &now
		if verResult.ProviderTxID != "" {
			order.ProviderTxID = &verResult.ProviderTxID
		}
		if verResult.PayerPhone != "" {
			order.PayerPhone = &verResult.PayerPhone
		}

		// Update order in database
		if err := s.paymentRepo.UpdateOrderStatus(
			ctx,
			order.OrderID,
			models.PaymentStatusSuccessful,
			order.ProviderTxID,
			order.PayerPhone,
			nil,
			verResult.RawResponse,
			order.VerificationMetadata,
		); err != nil {
			return nil, fmt.Errorf("failed to update order status: %w", err)
		}

		// Execute business fulfillment
		if order.OrderType == models.OrderTypeTierUpgrade && order.TargetTier != nil {
			if err := s.paymentRepo.ActivateTier(
				ctx,
				order.UserID,
				order.UserRole,
				*order.TargetTier,
				order.DurationDays,
				order.OrderID,
			); err != nil {
				return nil, fmt.Errorf("tier activation failed: %w", err)
			}
		} else if order.OrderType == models.OrderTypeListingFee && order.ProductID != nil {
			if err := s.paymentRepo.ActivateListingFee(
				ctx,
				*order.ProductID,
				order.Amount,
				order.OrderID,
			); err != nil {
				return nil, fmt.Errorf("listing fee activation failed: %w", err)
			}
		}

		_ = s.paymentRepo.RecordTransaction(ctx, order)
		return order, nil
	}

	if verResult.Status == models.PaymentStatusFailed ||
		verResult.Status == models.PaymentStatusCancelled ||
		verResult.Status == models.PaymentStatusExpired {
		order.Status = verResult.Status
		order.FailureReason = &verResult.FailureReason
		_ = s.paymentRepo.UpdateOrderStatus(
			ctx,
			order.OrderID,
			verResult.Status,
			order.ProviderTxID,
			order.PayerPhone,
			&verResult.FailureReason,
			verResult.RawResponse,
			order.VerificationMetadata,
		)
		return order, nil
	}

	return order, nil
}

// SimulatePaymentAction handles developer/demo simulation of success, failed, cancelled, expired
func (s *PaymentService) SimulatePaymentAction(ctx context.Context, orderID, action string) (*models.PaymentOrder, error) {
	order, err := s.paymentRepo.GetOrderByOrderID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("order not found: %w", err)
	}

	if !order.IsDemo && order.Environment == models.PaymentEnvProduction {
		return nil, errors.New("payment simulation is strictly forbidden on production orders")
	}

	actionClean := strings.ToLower(strings.TrimSpace(action))
	switch actionClean {
	case "success", "successful":
		now := time.Now()
		order.Status = models.PaymentStatusSuccessful
		order.VerifiedAt = &now
		simTx := fmt.Sprintf("SIM-TX-%d", time.Now().UnixNano())
		order.ProviderTxID = &simTx

		meta := order.VerificationMetadata
		if meta == nil {
			meta = make(map[string]interface{})
		}
		meta["simulated_action"] = "success"

		if err := s.paymentRepo.UpdateOrderStatus(
			ctx,
			order.OrderID,
			models.PaymentStatusSuccessful,
			order.ProviderTxID,
			order.PayerPhone,
			nil,
			map[string]interface{}{"simulation": true, "action": "success"},
			meta,
		); err != nil {
			return nil, err
		}

		if order.OrderType == models.OrderTypeTierUpgrade && order.TargetTier != nil {
			if errAct := s.paymentRepo.ActivateTier(ctx, order.UserID, order.UserRole, *order.TargetTier, order.DurationDays, order.OrderID); errAct != nil {
				return nil, fmt.Errorf("ActivateTier failed: %w", errAct)
			}
		} else if order.OrderType == models.OrderTypeListingFee && order.ProductID != nil {
			if errAct := s.paymentRepo.ActivateListingFee(ctx, *order.ProductID, order.Amount, order.OrderID); errAct != nil {
				return nil, fmt.Errorf("ActivateListingFee failed: %w", errAct)
			}
		}

		_ = s.paymentRepo.RecordTransaction(ctx, order)
		return order, nil

	case "fail", "failed":
		reason := "Simulated user payment card/account decline"
		order.Status = models.PaymentStatusFailed
		order.FailureReason = &reason
		_ = s.paymentRepo.UpdateOrderStatus(ctx, order.OrderID, models.PaymentStatusFailed, order.ProviderTxID, order.PayerPhone, &reason, nil, nil)
		return order, nil

	case "cancel", "cancelled":
		reason := "Simulated user cancelled transaction at checkout"
		order.Status = models.PaymentStatusCancelled
		order.FailureReason = &reason
		_ = s.paymentRepo.UpdateOrderStatus(ctx, order.OrderID, models.PaymentStatusCancelled, order.ProviderTxID, order.PayerPhone, &reason, nil, nil)
		return order, nil

	case "expire", "expired", "timeout":
		reason := "Simulated payment session expired / timeout"
		order.Status = models.PaymentStatusExpired
		order.FailureReason = &reason
		_ = s.paymentRepo.UpdateOrderStatus(ctx, order.OrderID, models.PaymentStatusExpired, order.ProviderTxID, order.PayerPhone, &reason, nil, nil)
		return order, nil

	default:
		return nil, fmt.Errorf("invalid simulation action: %s (choose: success, failed, cancelled, expired)", action)
	}
}

// HandleWebhook processes inbound webhook callbacks from payment providers
func (s *PaymentService) HandleWebhook(ctx context.Context, providerCode string, rawBody []byte, signature string) (*models.PaymentOrder, error) {
	globalMode, _ := s.paymentRepo.GetGlobalOperatingMode(ctx)
	isProd := globalMode == "production"

	adapter, err := s.resolveAdapter(ctx, providerCode, isProd)
	if err != nil {
		return nil, err
	}

	res, err := adapter.ProcessWebhook(ctx, rawBody, signature)
	if err != nil {
		return nil, fmt.Errorf("webhook processing error: %w", err)
	}

	order, err := s.paymentRepo.GetOrderByOrderID(ctx, res.OrderID)
	if err != nil {
		return nil, fmt.Errorf("webhook order not found: %s", res.OrderID)
	}

	if order.Status == models.PaymentStatusSuccessful {
		return order, nil // Idempotent duplicate delivery
	}

	if res.Status == models.PaymentStatusSuccessful {
		now := time.Now()
		order.Status = models.PaymentStatusSuccessful
		order.VerifiedAt = &now
		if res.ProviderTxID != "" {
			order.ProviderTxID = &res.ProviderTxID
		}
		if res.PayerPhone != "" {
			order.PayerPhone = &res.PayerPhone
		}

		_ = s.paymentRepo.UpdateOrderStatus(
			ctx,
			order.OrderID,
			models.PaymentStatusSuccessful,
			order.ProviderTxID,
			order.PayerPhone,
			nil,
			res.RawPayload,
			order.VerificationMetadata,
		)

		if order.OrderType == models.OrderTypeTierUpgrade && order.TargetTier != nil {
			if errAct := s.paymentRepo.ActivateTier(ctx, order.UserID, order.UserRole, *order.TargetTier, order.DurationDays, order.OrderID); errAct != nil {
				return nil, fmt.Errorf("ActivateTier failed: %w", errAct)
			}
		} else if order.OrderType == models.OrderTypeListingFee && order.ProductID != nil {
			if errAct := s.paymentRepo.ActivateListingFee(ctx, *order.ProductID, order.Amount, order.OrderID); errAct != nil {
				return nil, fmt.Errorf("ActivateListingFee failed: %w", errAct)
			}
		}

		_ = s.paymentRepo.RecordTransaction(ctx, order)
		return order, nil
	}

	order.Status = res.Status
	order.FailureReason = &res.FailureReason
	_ = s.paymentRepo.UpdateOrderStatus(
		ctx,
		order.OrderID,
		res.Status,
		&res.ProviderTxID,
		&res.PayerPhone,
		&res.FailureReason,
		res.RawPayload,
		order.VerificationMetadata,
	)

	return order, nil
}

// SendUpgradeOTP sends a dedicated Fida SMS OTP for tier upgrade
func (s *PaymentService) SendUpgradeOTP(ctx context.Context, userID int, rawPhone, targetTier, role string) (string, error) {
	plainOTP, err := s.userRepo.GenerateOTPWithType(ctx, rawPhone, "tier_upgrade", role)
	if err != nil {
		return "", fmt.Errorf("failed to generate upgrade OTP: %w", err)
	}

	normalized, err := NormalizeEthiopianPhone(rawPhone)
	if err != nil {
		return "", err
	}

	err = s.smsManager.SendOTPMessage(ctx, normalized, plainOTP, "tier_upgrade")
	if err != nil {
		// Log warning but proceed if simulation
		return plainOTP, nil
	}

	return plainOTP, nil
}

// VerifyUpgradeOTP validates the upgrade OTP
func (s *PaymentService) VerifyUpgradeOTP(ctx context.Context, userID int, rawPhone, code, role string) error {
	normalized, err := NormalizeEthiopianPhone(rawPhone)
	if err != nil {
		return err
	}

	return s.userRepo.VerifyOTPForRole(ctx, normalized, code, "tier_upgrade", role)
}

// SetOperatingMode sets the global payment environment (demo or production)
func (s *PaymentService) SetOperatingMode(ctx context.Context, adminID int, mode string) error {
	modeClean := strings.ToLower(strings.TrimSpace(mode))
	if modeClean != "demo" && modeClean != "production" {
		return errors.New("operating mode must be either 'demo' or 'production'")
	}

	// Fail-closed verification if switching to production
	if modeClean == "production" {
		tbCfg, err := s.paymentRepo.GetProviderConfigByCode(ctx, "telebirr")
		if err != nil || tbCfg == nil || tbCfg.AppID == "" || tbCfg.APIEndpoint == "" {
			return errors.New("cannot activate production: Telebirr credentials incomplete (AppID/Endpoint missing)")
		}
		tbKey, _ := s.paymentRepo.GetDecryptedAppKey(ctx, "telebirr")
		if tbKey == "" {
			return errors.New("cannot activate production: Telebirr AppKey is missing")
		}

		etCfg, err := s.paymentRepo.GetProviderConfigByCode(ctx, "ethio_telecom")
		if err != nil || etCfg == nil || etCfg.AppID == "" || etCfg.APIEndpoint == "" {
			return errors.New("cannot activate production: Ethio Telecom credentials incomplete")
		}
	}

	return s.paymentRepo.SetGlobalOperatingMode(ctx, modeClean)
}

func (s *PaymentService) GetPaymentRepo() *repository.PaymentRepo {
	return s.paymentRepo
}
