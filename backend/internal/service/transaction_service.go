package service

import (
	"context"
	"errors"
	"log"

	"github.com/your-org/used-tools-marketplace/pkg/dto"
	"github.com/your-org/used-tools-marketplace/internal/models"
	"github.com/your-org/used-tools-marketplace/internal/repository"
	"github.com/your-org/used-tools-marketplace/pkg/validations"
)

type TransactionService struct {
	repo        *repository.TransactionRepo
	userRepo    *repository.UserRepo
	referralSvc *ReferralService
}

func NewTransactionService(repo repository.TransactionRepository, _ repository.ProductRepository, userRepo repository.UserRepository) *TransactionService {
	// productRepo is no longer needed here — the pre-check is now inside the DB transaction.
	return &TransactionService{repo: repo, userRepo: userRepo}
}

// WithReferralService injects the referral service for commission recording.
func (s *TransactionService) WithReferralService(r *ReferralService) *TransactionService {
	s.referralSvc = r
	return s
}

func (s *TransactionService) Create(ctx context.Context, sellerID int, req dto.CreateTransactionRequest) (*models.Transaction, error) {
	if req.BuyerName == "" || req.BuyerPhone == "" {
		return nil, errors.New("buyer_name and buyer_phone are required")
	}
	if req.ProductID <= 0 {
		return nil, errors.New("product_id is required")
	}
	if req.Amount <= 0 {
		return nil, errors.New("amount must be greater than zero")
	}

	phoneRes := validations.ValidateEthiopianPhone(req.BuyerPhone)
	if !phoneRes.Valid {
		return nil, errors.New("invalid Ethiopian phone number for buyer_phone")
	}
	req.BuyerPhone = phoneRes.Normalized

	// Security: Prevent self-dealing wash trading by using the seller's own phone number
	if s.userRepo != nil && sellerID > 0 {
		sellerUser, err := s.userRepo.GetByID(ctx, sellerID)
		if err == nil && sellerUser != nil {
			if sellerUser.Phone == req.BuyerPhone {
				return nil, errors.New("cannot create a transaction with your own phone number as the buyer")
			}
		}
	}

	t := &models.Transaction{
		ProductID:  req.ProductID,
		BuyerName:  req.BuyerName,
		BuyerPhone: req.BuyerPhone,
		Amount:     req.Amount,
		Notes:      req.Notes,
	}

	// CreateAndSell is atomic: locks the product, verifies ownership and approval,
	// inserts the transaction, and marks the product sold — all in one DB transaction.
	actualSellerID, err := s.repo.CreateAndSell(ctx, sellerID, t)
	if err != nil {
		return nil, err
	}

	// Non-fatal: tier progression failure doesn't invalidate the completed sale.
	// Use seller_gets (net payout), not the gross amount — total_value_sold tracks
	// what the seller actually earned, which is the correct basis for tier thresholds.
	if err := s.userRepo.UpdateTotalValueSold(ctx, actualSellerID, t.SellerGets); err != nil {
		log.Printf("[TransactionService] failed to update seller %d total_value_sold: %v", actualSellerID, err) //nolint:gosec // G117: error/value from internal source, not user-controlled input
	}

	// Non-fatal: record referral commission for the seller's inviter.
	// Commission base = broker_fee (platform's cut), NOT the full sale amount.
	// This ensures the platform keeps its margin and the inviter gets a share of the fee only.
	if s.referralSvc != nil {
		srcID := t.ID
		safeGo(func() {
			s.referralSvc.MaybeRecordCommission(
				context.Background(), actualSellerID, t.BrokerFee,
				"product_purchase", "seller", "",
				&srcID,
			)
		})
	}

	return t, nil
}

// UpdateStatus allows admin to mark a transaction completed or cancelled.
func (s *TransactionService) UpdateStatus(ctx context.Context, txID int, status string) error {
	valid := map[string]bool{"completed": true, "cancelled": true}
	if !valid[status] {
		return errors.New("status must be completed or cancelled")
	}
	return s.repo.UpdateStatus(ctx, txID, status)
}
