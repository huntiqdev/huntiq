package service

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/your-org/used-tools-marketplace/internal/models"
	"github.com/your-org/used-tools-marketplace/internal/repository"
)

type LedgerService struct {
	repo *repository.LedgerRepo
}

func NewLedgerService(repo *repository.LedgerRepo) *LedgerService {
	return &LedgerService{repo: repo}
}

// roundCents rounds a monetary amount to 2 decimal places.
func roundCents(val float64) float64 {
	return math.Round(val*100.0) / 100.0
}

// GetWallet returns the current financial state of the user's wallet.
func (s *LedgerService) GetWallet(ctx context.Context, userID int) (*models.WalletAccount, error) {
	if userID <= 0 {
		return nil, errors.New("invalid user ID")
	}
	return s.repo.GetWalletAccount(ctx, userID)
}

// GetLedgerHistory returns paginated ledger transactions for a user.
func (s *LedgerService) GetLedgerHistory(ctx context.Context, userID int, limit, offset int) ([]models.LedgerTransaction, int, error) {
	if userID <= 0 {
		return nil, 0, errors.New("invalid user ID")
	}
	return s.repo.GetLedgerTransactions(ctx, userID, limit, offset)
}

// GetAllLedgerHistory returns system-wide ledger history for administrative audit.
func (s *LedgerService) GetAllLedgerHistory(ctx context.Context, limit, offset int) ([]models.LedgerTransaction, int, error) {
	return s.repo.GetAllLedgerTransactions(ctx, limit, offset)
}

// CreditRegistrationReward credits Reward A to the inviter's wallet upon verified registration.
func (s *LedgerService) CreditRegistrationReward(ctx context.Context, inviterID, referredUserID int) error {
	if inviterID <= 0 || referredUserID <= 0 || inviterID == referredUserID {
		return errors.New("invalid inviter or referred user ID for registration reward")
	}

	policy, err := s.repo.GetActiveRewardPolicy(ctx)
	if err != nil {
		return fmt.Errorf("failed to load reward policy: %w", err)
	}
	if policy == nil || !policy.IsActive || policy.RegistrationRewardAmount <= 0 {
		return nil // Rewards disabled or zero
	}

	amount := roundCents(policy.RegistrationRewardAmount)
	refID := fmt.Sprintf("user:%d", referredUserID)
	notes := fmt.Sprintf("Reward A: Verified phone registration bonus for inviting user %d (Policy v%d)", referredUserID, policy.Version)

	txEntry := &models.LedgerTransaction{
		LedgerTxID:     fmt.Sprintf("LTX-REG-%d-%d", inviterID, referredUserID),
		UserID:         inviterID,
		DebitAccount:   "platform_marketing_pool",
		CreditAccount:  "user_available",
		Amount:         amount,
		Currency:       "ETB",
		EntryType:      "referral_reward_a",
		ReferenceID:    &refID,
		SourceEvent:    "phone_verified_registration",
		Status:         "posted",
		IdempotencyKey: fmt.Sprintf("reg-reward-%d-%d", inviterID, referredUserID),
		Notes:          &notes,
	}

	return s.repo.RecordLedgerTransaction(ctx, txEntry)
}

// CreditTransactionCommission credits Reward B to the inviter when a referred user completes a qualifying transaction.
func (s *LedgerService) CreditTransactionCommission(ctx context.Context, inviterID, referredUserID, txID int, orderAmount float64) error {
	if inviterID <= 0 || referredUserID <= 0 || txID <= 0 || orderAmount <= 0 {
		return nil
	}

	policy, err := s.repo.GetActiveRewardPolicy(ctx)
	if err != nil {
		return fmt.Errorf("failed to load reward policy: %w", err)
	}
	if policy == nil || !policy.IsActive || policy.TransactionRewardPercent <= 0 {
		return nil
	}

	if orderAmount < policy.MinQualifyingTxAmount {
		return nil // Does not meet minimum order threshold
	}

	calculated := roundCents(orderAmount * policy.TransactionRewardPercent / 100.0)
	if policy.MaxRewardPerTx > 0 && calculated > policy.MaxRewardPerTx {
		calculated = policy.MaxRewardPerTx
	}

	if calculated <= 0 {
		return nil
	}

	refID := fmt.Sprintf("tx:%d", txID)
	notes := fmt.Sprintf("Reward B: %.2f%% qualifying transaction commission on order #%d (%.2f ETB)", policy.TransactionRewardPercent, txID, orderAmount)

	txEntry := &models.LedgerTransaction{
		LedgerTxID:     fmt.Sprintf("LTX-TX-%d-%d", txID, inviterID),
		UserID:         inviterID,
		DebitAccount:   "platform_revenue_pool",
		CreditAccount:  "user_available",
		Amount:         calculated,
		Currency:       "ETB",
		EntryType:      "referral_commission_b",
		ReferenceID:    &refID,
		SourceEvent:    "qualifying_transaction_completed",
		Status:         "posted",
		IdempotencyKey: fmt.Sprintf("tx-reward-%d-%d", txID, inviterID),
		Notes:          &notes,
	}

	return s.repo.RecordLedgerTransaction(ctx, txEntry)
}

// ReserveWithdrawal holds user funds atomically for an outgoing payout request.
func (s *LedgerService) ReserveWithdrawal(ctx context.Context, userID int, amount float64, withdrawalID int) error {
	if userID <= 0 || amount <= 0 || withdrawalID <= 0 {
		return errors.New("invalid withdrawal reservation parameters")
	}

	cleanAmount := roundCents(amount)
	refID := fmt.Sprintf("wd:%d", withdrawalID)
	notes := fmt.Sprintf("Hold reserved funds for withdrawal request #%d", withdrawalID)

	txEntry := &models.LedgerTransaction{
		LedgerTxID:     fmt.Sprintf("LTX-WD-HOLD-%d", withdrawalID),
		UserID:         userID,
		DebitAccount:   "user_available",
		CreditAccount:  "user_reserved",
		Amount:         cleanAmount,
		Currency:       "ETB",
		EntryType:      "withdrawal_hold",
		ReferenceID:    &refID,
		SourceEvent:    "withdrawal_requested",
		Status:         "posted",
		IdempotencyKey: fmt.Sprintf("wd-reserve-%d", withdrawalID),
		Notes:          &notes,
	}

	return s.repo.RecordLedgerTransaction(ctx, txEntry)
}

// SettleWithdrawal finalizes the payout deduction upon confirmed external disbursement.
func (s *LedgerService) SettleWithdrawal(ctx context.Context, userID int, amount float64, withdrawalID int, adminID int, reference string) error {
	cleanAmount := roundCents(amount)
	refID := fmt.Sprintf("wd:%d", withdrawalID)
	notes := fmt.Sprintf("Settled payout for withdrawal #%d (Ref: %s)", withdrawalID, reference)

	txEntry := &models.LedgerTransaction{
		LedgerTxID:     fmt.Sprintf("LTX-WD-SETTLE-%d", withdrawalID),
		UserID:         userID,
		DebitAccount:   "user_reserved",
		CreditAccount:  "telebirr_payout_gateway",
		Amount:         cleanAmount,
		Currency:       "ETB",
		EntryType:      "withdrawal_settle",
		ReferenceID:    &refID,
		SourceEvent:    "withdrawal_payout_settled",
		Status:         "posted",
		IdempotencyKey: fmt.Sprintf("wd-settle-%d", withdrawalID),
		Notes:          &notes,
		ActorID:        &adminID,
	}

	return s.repo.RecordLedgerTransaction(ctx, txEntry)
}

// ReleaseWithdrawalReservation releases reserved funds back to user available balance on rejection.
func (s *LedgerService) ReleaseWithdrawalReservation(ctx context.Context, userID int, amount float64, withdrawalID int, adminID int, reason string) error {
	cleanAmount := roundCents(amount)
	refID := fmt.Sprintf("wd:%d", withdrawalID)
	notes := fmt.Sprintf("Released reserved hold for rejected withdrawal #%d: %s", withdrawalID, reason)

	txEntry := &models.LedgerTransaction{
		LedgerTxID:     fmt.Sprintf("LTX-WD-REL-%d", withdrawalID),
		UserID:         userID,
		DebitAccount:   "user_reserved",
		CreditAccount:  "user_available",
		Amount:         cleanAmount,
		Currency:       "ETB",
		EntryType:      "withdrawal_reversal",
		ReferenceID:    &refID,
		SourceEvent:    "withdrawal_rejected",
		Status:         "posted",
		IdempotencyKey: fmt.Sprintf("wd-release-%d", withdrawalID),
		Notes:          &notes,
		ActorID:        &adminID,
	}

	return s.repo.RecordLedgerTransaction(ctx, txEntry)
}

// CompensateBalance records an audited compensating manual entry with maker attribution.
func (s *LedgerService) CompensateBalance(ctx context.Context, userID int, amount float64, isCredit bool, reason string, adminID int) error {
	if userID <= 0 || amount <= 0 || adminID <= 0 || reason == "" {
		return errors.New("invalid parameters for manual compensating ledger entry")
	}

	cleanAmount := roundCents(amount)
	notes := fmt.Sprintf("Manual compensating entry: %s", reason)
	txID := fmt.Sprintf("LTX-COMP-%d-%d", userID, time.Now().UnixNano())

	debitAcc := "platform_reconciliation_pool"
	creditAcc := "user_available"
	if !isCredit {
		debitAcc = "user_available"
		creditAcc = "platform_reconciliation_pool"
	}

	txEntry := &models.LedgerTransaction{
		LedgerTxID:     txID,
		UserID:         userID,
		DebitAccount:   debitAcc,
		CreditAccount:  creditAcc,
		Amount:         cleanAmount,
		Currency:       "ETB",
		EntryType:      "manual_compensation",
		SourceEvent:    "admin_manual_compensation",
		Status:         "posted",
		IdempotencyKey: fmt.Sprintf("comp-%d-%d", userID, time.Now().UnixNano()),
		Notes:          &notes,
		ActorID:        &adminID,
	}

	return s.repo.RecordLedgerTransaction(ctx, txEntry)
}

// GetRewardPolicy returns the current active policy.
func (s *LedgerService) GetRewardPolicy(ctx context.Context) (*models.ReferralRewardPolicy, error) {
	return s.repo.GetActiveRewardPolicy(ctx)
}

// UpdateRewardPolicy modifies the policy rules and stores an immutable audit snapshot.
func (s *LedgerService) UpdateRewardPolicy(ctx context.Context, p *models.ReferralRewardPolicy, adminID int, reason string) error {
	if p.RegistrationRewardAmount < 0 || p.TransactionRewardPercent < 0 || p.MinQualifyingTxAmount < 0 {
		return errors.New("financial reward values cannot be negative")
	}
	return s.repo.UpdateRewardPolicy(ctx, p, adminID, reason)
}

