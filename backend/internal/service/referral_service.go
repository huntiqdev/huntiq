package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/your-org/used-tools-marketplace/internal/repository"
	"github.com/your-org/used-tools-marketplace/pkg/constants"
)

// ReferralService handles referral code logic and commission recording.
type ReferralService struct {
	repo           *repository.ReferralRepo
	settingsReader func(ctx context.Context, key string) string
}

func NewReferralService(repo repository.ReferralRepository) *ReferralService {
	return &ReferralService{repo: repo}
}

// WithSettingsReader injects the settings reader function
func (s *ReferralService) WithSettingsReader(reader func(ctx context.Context, key string) string) *ReferralService {
	s.settingsReader = reader
	return s
}

// GetStats returns referral dashboard data for a user.
func (s *ReferralService) GetStats(ctx context.Context, userID int) (*repository.ReferralStats, error) {
	return s.repo.GetStats(ctx, userID)
}

// GetCommissions returns paginated commission history for an inviter.
func (s *ReferralService) GetCommissions(ctx context.Context, inviterID, page, limit int) ([]repository.CommissionEntry, int, error) {
	return s.repo.GetCommissions(ctx, inviterID, page, limit)
}

// AcceptInvitation links a registered user to an inviter via referral code.
func (s *ReferralService) AcceptInvitation(ctx context.Context, userID int, referralCode string) error {
	if referralCode == "" {
		return errors.New("referral code is required")
	}
	return s.repo.AcceptInvitationByCode(ctx, userID, referralCode)
}

// SetInvitedByOnRegister links a new user to an inviter at registration time (fire-and-forget safe).
func (s *ReferralService) SetInvitedByOnRegister(ctx context.Context, userID int, referralCode string) {
	if referralCode == "" {
		return
	}
	if err := s.repo.SetInvitedByOnRegister(ctx, userID, referralCode); err != nil {
		log.Printf("[ReferralService] SetInvitedByOnRegister user=%d code=%s: %v", userID, referralCode, err) //nolint:gosec // G117: error/value from internal source, not user-controlled input
	}
}

// MaybeRecordCommission checks if the acting user has an inviter and records a 2% commission.
//
// Rules enforced here:
//   - inviter's referral code must not be suspended
//   - normal buyers (role="buyer", tier="normal") generate ZERO commission — only paid buyers do
//   - commission is idempotent: ON CONFLICT DO NOTHING prevents double-recording
//
// txType: 'product_purchase' | 'gold_upgrade' | 'premium_upgrade' |
//
//	'buyer_gold_upgrade' | 'buyer_premium_upgrade' |
//	'seller_registration' | 'seller_renewal'
//
// actingUserRole: "seller" or "buyer" — determines eligibility
// actingUserTier: buyer tier ("normal"/"gold"/"premium") — normal buyers don't generate commission
// sourceID: ID of the transaction/upgrade request for deduplication
func (s *ReferralService) MaybeRecordCommission(
	ctx context.Context,
	actingUserID int,
	txAmount float64,
	txType string,
	actingUserRole string,
	actingUserTier string, // only relevant for buyers
	sourceID *int,
) {
	// Rule: normal buyers never generate commission for their inviter
	if actingUserRole == "buyer" && (actingUserTier == "normal" || actingUserTier == "") {
		return
	}

	inviterID, suspended, err := s.repo.GetInviterID(ctx, actingUserID)
	if err != nil || inviterID == 0 {
		return // no inviter
	}
	if suspended {
		return // inviter's code is suspended — no commission
	}
	if err := s.repo.RecordCommission(ctx, inviterID, actingUserID, txAmount, txType, actingUserRole, sourceID); err != nil {
		log.Printf("[ReferralService] RecordCommission inviter=%d invited=%d type=%s: %v", //nolint:gosec // G117: error/value from internal source, not user-controlled input
			inviterID, actingUserID, txType, err)
	}
}

// RecordPendingCommission records a commission with 'pending' status that will be
// confirmed after the chargeback window expires (typically 30 days).
// This prevents commission withdrawal before payment finalization.
//
// SECURITY: Protects against chargeback fraud where attacker:
//  1. Creates fake upgrades → Earns commission
//  2. Withdraws commission
//  3. Disputes payments with bank → Gets refund
//  4. Result: Keeps commission + gets money back = FREE MONEY
//
// SOLUTION: Commission stays 'pending' for 30 days. If payment is disputed,
// commission status can be set to 'reversed' before it becomes withdrawable.
func (s *ReferralService) RecordPendingCommission(
	ctx context.Context,
	actingUserID int,
	txAmount float64,
	txType string,
	actingUserRole string,
	actingUserTier string,
	sourceID *int,
	pendingDuration time.Duration, // typically 30 days
) {
	// Same eligibility rules as MaybeRecordCommission
	if actingUserRole == "buyer" && (actingUserTier == "normal" || actingUserTier == "") {
		return
	}

	inviterID, suspended, err := s.repo.GetInviterID(ctx, actingUserID)
	if err != nil || inviterID == 0 {
		return // no inviter
	}
	if suspended {
		return // inviter's code is suspended
	}

	// Calculate pending_until timestamp
	pendingUntil := time.Now().Add(pendingDuration)

	if err := s.repo.RecordPendingCommission(
		ctx, inviterID, actingUserID, txAmount, txType, actingUserRole, sourceID, pendingUntil,
	); err != nil {
		log.Printf("[ReferralService] RecordPendingCommission inviter=%d invited=%d type=%s pending_until=%v: %v",
			inviterID, actingUserID, txType, pendingUntil, err)
	}
}

// AdminGetOverview returns aggregate referral stats for the admin profit page.
func (s *ReferralService) AdminGetOverview(ctx context.Context) (map[string]interface{}, error) {
	return s.repo.AdminGetReferralOverview(ctx)
}

// AdminListReferrals returns paginated referral relationships.
func (s *ReferralService) AdminListReferrals(ctx context.Context, page, limit int, search string) ([]repository.AdminReferralRow, int, error) {
	return s.repo.AdminListReferrals(ctx, page, limit, search)
}

// AdminSetSuspended suspends or unsuspends a user's referral code.
func (s *ReferralService) AdminSetSuspended(ctx context.Context, userID int, suspended bool) error {
	return s.repo.SetReferralSuspended(ctx, userID, suspended)
}

// AdminSetPermanentlyBlocked permanently blocks or unblocks a user from the referral system.
// Blocked users cannot earn commissions, request withdrawals, or accept invitations.
func (s *ReferralService) AdminSetPermanentlyBlocked(ctx context.Context, userID int, blocked bool, reason string) error {
	if blocked && reason == "" {
		return errors.New("block reason is required")
	}
	return s.repo.SetPermanentlyBlocked(ctx, userID, blocked, reason)
}

// GetUserBalance returns a user's referral balance.
func (s *ReferralService) GetUserBalance(ctx context.Context, userID int) (*repository.UserBalance, error) {
	return s.repo.GetUserBalance(ctx, userID)
}

// RequestWithdrawal validates amount against platform settings then creates a pending request.
func (s *ReferralService) RequestWithdrawal(ctx context.Context, userID int, amount float64, paymentMethod, payerPhone string) (int, error) {
	minAmount := constants.DefaultMinReferralWithdrawal
	maxAmount := constants.DefaultMaxReferralWithdrawal

	if s.settingsReader != nil {
		if v := s.settingsReader(ctx, constants.SettingMinReferralWithdrawal); v != "" {
			if f, err := strconv.ParseFloat(v, 64); err == nil && f > 0 {
				minAmount = f
			}
		}
		if v := s.settingsReader(ctx, constants.SettingMaxReferralWithdrawal); v != "" {
			if f, err := strconv.ParseFloat(v, 64); err == nil && f > 0 {
				maxAmount = f
			}
		}
	}
	if amount < minAmount {
		return 0, fmt.Errorf("minimum withdrawal amount is ETB %.0f", minAmount)
	}
	if amount > maxAmount {
		return 0, fmt.Errorf("maximum withdrawal amount is ETB %.0f per request", maxAmount)
	}
	return s.repo.RequestWithdrawal(ctx, userID, amount, paymentMethod, payerPhone)
}

// GetMyWithdrawals returns a user's own withdrawal history.
func (s *ReferralService) GetMyWithdrawals(ctx context.Context, userID int) ([]repository.WithdrawalRequest, error) {
	return s.repo.GetMyWithdrawals(ctx, userID)
}

// AdminListWithdrawals returns paginated withdrawal requests for admin review.
func (s *ReferralService) AdminListWithdrawals(ctx context.Context, status string, page, limit int) ([]repository.WithdrawalRequest, int, error) {
	return s.repo.AdminListWithdrawals(ctx, status, page, limit)
}

// AdminApproveWithdrawal marks a withdrawal approved.
func (s *ReferralService) AdminApproveWithdrawal(ctx context.Context, withdrawalID, adminID int, note string) error {
	return s.repo.AdminApproveWithdrawal(ctx, withdrawalID, adminID, note)
}

// AdminRejectWithdrawal rejects a withdrawal with a reason.
func (s *ReferralService) AdminRejectWithdrawal(ctx context.Context, withdrawalID, adminID int, reason string) error {
	return s.repo.AdminRejectWithdrawal(ctx, withdrawalID, adminID, reason)
}

// AdminMarkWithdrawalPaid marks a withdrawal as paid and records the reference.
func (s *ReferralService) AdminMarkWithdrawalPaid(ctx context.Context, withdrawalID, adminID int, reference string) error {
	return s.repo.AdminMarkWithdrawalPaid(ctx, withdrawalID, adminID, reference)
}
