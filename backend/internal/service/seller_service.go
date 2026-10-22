package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/your-org/used-tools-marketplace/pkg/dto"
	"github.com/your-org/used-tools-marketplace/internal/models"
	"github.com/your-org/used-tools-marketplace/internal/repository"
)

type SellerService struct {
	repo *repository.SellerRepoExport
}

func NewSellerService(repo *repository.SellerRepoExport) *SellerService {
	return &SellerService{repo: repo}
}

func (s *SellerService) GetStats(ctx context.Context, sellerID int) (*repository.SellerStats, error) {
	return s.repo.GetStats(ctx, sellerID)
}

func (s *SellerService) GetProfile(ctx context.Context, sellerID int) (*models.User, error) {
	return s.repo.GetProfile(ctx, sellerID)
}

func (s *SellerService) UpdateProfile(ctx context.Context, u *models.User) error {
	if err := s.repo.UpdateProfile(ctx, u); err != nil {
		if strings.Contains(err.Error(), "users_email_key") {
			return errors.New("email address already in use by another account")
		}
		return err
	}
	return nil
}

func (s *SellerService) ListTransactions(ctx context.Context, sellerID int) ([]models.Transaction, error) {
	return s.repo.ListTransactionsBySeller(ctx, sellerID)
}

func (s *SellerService) ListProductsWithFee(ctx context.Context, sellerID int) ([]models.Product, error) {
	return s.repo.ListProductsWithFee(ctx, sellerID)
}

func (s *SellerService) DeleteProduct(ctx context.Context, productID, sellerID int) error {
	return s.repo.DeleteProduct(ctx, productID, sellerID)
}

// RecordListingFeePayment stores the seller's payment reference for admin verification.
// Sets fee_payment_pending=true so admin can see it needs review.
// fee_paid remains FALSE until admin explicitly verifies the payment.
// The product stays in 'pending' status until admin approves.
func (s *SellerService) RecordListingFeePayment(ctx context.Context, sellerID int, req dto.ListingFeePaymentRequest) error {
	return s.repo.RecordListingFeePayment(ctx, sellerID, req)
}

// ── Order Requests ────────────────────────────────────────────────────────────

func (s *SellerService) CreateOrderRequest(ctx context.Context, req dto.CreateOrderRequest) (*models.OrderRequest, error) {
	if req.BuyerName == "" || req.BuyerPhone == "" {
		return nil, errors.New("buyer_name and buyer_phone are required")
	}
	if req.ProductID <= 0 {
		return nil, errors.New("product_id is required")
	}
	sellerID, requiredTier, err := s.repo.GetProductSellerIDAndTier(ctx, req.ProductID)
	if err != nil {
		return nil, errors.New("product not found or not available")
	}
	// Enforce order hierarchy:
	// normal product → any buyer
	// gold product   → buyer must be gold or premium
	// premium product → buyer must be premium
	// BUG FIX #5: Updated argument order to match new TierAllowed signature
	if !TierAllowed(req.BuyerTier, requiredTier) {
		return nil, fmt.Errorf("this product requires a %s buyer account to place an order", requiredTier)
	}
	o := &models.OrderRequest{
		ProductID:        req.ProductID,
		SellerID:         sellerID,
		BuyerName:        req.BuyerName,
		BuyerPhone:       req.BuyerPhone,
		BuyerEmail:       req.BuyerEmail,
		Message:          req.Message,
		OfferPrice:       req.OfferPrice,
		BuyerTierAtOrder: req.BuyerTier,
	}
	id, err := s.repo.CreateOrderRequest(ctx, o)
	if err != nil {
		return nil, err
	}
	o.ID = id
	o.Status = "pending"
	return o, nil
}



func (s *SellerService) ListOrders(ctx context.Context, sellerID int) ([]models.OrderRequest, error) {
	return s.repo.ListOrdersBySeller(ctx, sellerID)
}

func (s *SellerService) UpdateOrderStatus(ctx context.Context, orderID, sellerID int, req dto.UpdateOrderRequest) error {
	allowed := map[string]bool{"accepted": true, "declined": true, "completed": true, "cancelled": true}
	if !allowed[req.Status] {
		return errors.New("invalid status")
	}
	return s.repo.UpdateOrderStatus(ctx, orderID, sellerID, req.Status, req.SellerNote)
}

// ── Tier Upgrade Requests ─────────────────────────────────────────────────────

// SubmitUpgradeRequest validates and stores a seller's tier upgrade payment request.
func (s *SellerService) SubmitUpgradeRequest(ctx context.Context, sellerID int, requestedTier, paymentMethod, payerPhone, reference, notes string, amount float64) (int, error) {
	if requestedTier != "gold" && requestedTier != "premium" {
		return 0, errors.New("requested_tier must be gold or premium")
	}
	validMethods := map[string]bool{"telebirr": true, "safaricom": true, "cbe": true}
	if !validMethods[paymentMethod] {
		return 0, errors.New("payment_method must be telebirr, safaricom, or cbe")
	}
	if payerPhone == "" {
		return 0, errors.New("payer_phone is required")
	}
	if reference == "" {
		return 0, errors.New("payment reference is required")
	}
	if amount <= 0 {
		return 0, errors.New("amount must be greater than zero")
	}
	// Verify seller exists and check current tier
	profile, err := s.repo.GetProfile(ctx, sellerID)
	if err != nil {
		return 0, errors.New("seller not found")
	}
	// Prevent downgrade or same-tier request
	if TierRank[requestedTier] <= TierRank[profile.Tier] {
		return 0, fmt.Errorf("you are already on %s tier or higher", profile.Tier)
	}
	return s.repo.SubmitUpgradeRequest(ctx, sellerID, requestedTier, paymentMethod, payerPhone, reference, notes, amount)
}

// GetUpgradeRequests returns the seller's upgrade request history.
func (s *SellerService) GetUpgradeRequests(ctx context.Context, sellerID int) ([]map[string]interface{}, error) {
	return s.repo.GetUpgradeRequests(ctx, sellerID)
}

// ListNotifications returns seller notifications (expiry alerts, etc.).
func (s *SellerService) ListNotifications(ctx context.Context, sellerID int) ([]map[string]interface{}, error) {
	return s.repo.ListNotifications(ctx, sellerID)
}

// MarkNotificationRead marks a seller notification as read.
func (s *SellerService) MarkNotificationRead(ctx context.Context, notifID, sellerID int) error {
	return s.repo.MarkNotificationRead(ctx, notifID, sellerID)
}

// CreateNotification inserts a notification for a seller. Fire-and-forget safe.
func (s *SellerService) CreateNotification(ctx context.Context, sellerID int, title, message, ntype string) error {
	return s.repo.CreateSellerNotification(ctx, sellerID, title, message, ntype)
}

// UpdatePreference updates a seller's listing preference.
func (s *SellerService) UpdatePreference(ctx context.Context, sellerID int, preference string) error {
	return s.repo.UpdatePreference(ctx, sellerID, preference)
}
