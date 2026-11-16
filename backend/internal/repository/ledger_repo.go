package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/your-org/used-tools-marketplace/internal/models"
)

type LedgerRepo struct {
	db *pgxpool.Pool
}

func NewLedgerRepo(db *pgxpool.Pool) *LedgerRepo {
	return &LedgerRepo{db: db}
}

// GetOrCreateWalletAccount ensures a wallet exists for the given user, creating one if necessary.
func (r *LedgerRepo) GetOrCreateWalletAccount(ctx context.Context, userID int) (*models.WalletAccount, error) {
	var w models.WalletAccount
	query := `
		INSERT INTO wallet_accounts (user_id, currency, available_balance, pending_balance, reserved_balance, total_withdrawn, lifetime_earned, is_locked)
		VALUES ($1, 'ETB', 0.00, 0.00, 0.00, 0.00, 0.00, false)
		ON CONFLICT (user_id) DO UPDATE SET updated_at = NOW()
		RETURNING id, user_id, currency, available_balance, pending_balance, reserved_balance, total_withdrawn, lifetime_earned, is_locked, created_at, updated_at
	`
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&w.ID,
		&w.UserID,
		&w.Currency,
		&w.AvailableBalance,
		&w.PendingBalance,
		&w.ReservedBalance,
		&w.TotalWithdrawn,
		&w.LifetimeEarned,
		&w.IsLocked,
		&w.CreatedAt,
		&w.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &w, nil
}

// GetWalletAccount fetches a wallet account for reading.
func (r *LedgerRepo) GetWalletAccount(ctx context.Context, userID int) (*models.WalletAccount, error) {
	var w models.WalletAccount
	query := `
		SELECT id, user_id, currency, available_balance, pending_balance, reserved_balance, total_withdrawn, lifetime_earned, is_locked, created_at, updated_at
		FROM wallet_accounts
		WHERE user_id = $1
	`
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&w.ID,
		&w.UserID,
		&w.Currency,
		&w.AvailableBalance,
		&w.PendingBalance,
		&w.ReservedBalance,
		&w.TotalWithdrawn,
		&w.LifetimeEarned,
		&w.IsLocked,
		&w.CreatedAt,
		&w.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return r.GetOrCreateWalletAccount(ctx, userID)
		}
		return nil, err
	}
	return &w, nil
}

// RecordLedgerTransaction atomically writes an immutable double-entry ledger entry and adjusts wallet accounts.
func (r *LedgerRepo) RecordLedgerTransaction(ctx context.Context, txEntry *models.LedgerTransaction) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx) // Safe if already committed

	// 1. Check idempotency
	var existingID int64
	err = tx.QueryRow(ctx, `SELECT id FROM ledger_transactions WHERE idempotency_key = $1`, txEntry.IdempotencyKey).Scan(&existingID)
	if err == nil {
		// Already applied idempotently
		return nil
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return err
	}

	// 2. Lock and fetch wallet account
	var w models.WalletAccount
	err = tx.QueryRow(ctx, `
		SELECT id, user_id, currency, available_balance, pending_balance, reserved_balance, total_withdrawn, lifetime_earned, is_locked
		FROM wallet_accounts
		WHERE user_id = $1
		FOR UPDATE
	`, txEntry.UserID).Scan(
		&w.ID,
		&w.UserID,
		&w.Currency,
		&w.AvailableBalance,
		&w.PendingBalance,
		&w.ReservedBalance,
		&w.TotalWithdrawn,
		&w.LifetimeEarned,
		&w.IsLocked,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// Create wallet row within tx
			err = tx.QueryRow(ctx, `
				INSERT INTO wallet_accounts (user_id, currency, available_balance, pending_balance, reserved_balance, total_withdrawn, lifetime_earned, is_locked)
				VALUES ($1, 'ETB', 0.00, 0.00, 0.00, 0.00, 0.00, false)
				RETURNING id, user_id, currency, available_balance, pending_balance, reserved_balance, total_withdrawn, lifetime_earned, is_locked
			`, txEntry.UserID).Scan(
				&w.ID,
				&w.UserID,
				&w.Currency,
				&w.AvailableBalance,
				&w.PendingBalance,
				&w.ReservedBalance,
				&w.TotalWithdrawn,
				&w.LifetimeEarned,
				&w.IsLocked,
			)
			if err != nil {
				return fmt.Errorf("failed to initialize wallet inside tx: %w", err)
			}
		} else {
			return err
		}
	}

	if w.IsLocked {
		return errors.New("wallet is locked: financial operations suspended")
	}

	// 3. Mutate balances according to entry type
	switch txEntry.EntryType {
	case "referral_reward_a", "referral_commission_b":
		// Direct credit to available balance
		if txEntry.CreditAccount == "user_available" {
			w.AvailableBalance += txEntry.Amount
			w.LifetimeEarned += txEntry.Amount
		} else if txEntry.CreditAccount == "user_pending" {
			w.PendingBalance += txEntry.Amount
			w.LifetimeEarned += txEntry.Amount
		}
	case "withdrawal_hold":
		// Move from available to reserved
		if w.AvailableBalance < txEntry.Amount {
			return fmt.Errorf("insufficient available balance: requested %.2f, available %.2f", txEntry.Amount, w.AvailableBalance)
		}
		w.AvailableBalance -= txEntry.Amount
		w.ReservedBalance += txEntry.Amount
	case "withdrawal_settle":
		// Reserved moves out of platform to external payout
		if w.ReservedBalance < txEntry.Amount {
			return fmt.Errorf("insufficient reserved balance for settlement: requested %.2f, reserved %.2f", txEntry.Amount, w.ReservedBalance)
		}
		w.ReservedBalance -= txEntry.Amount
		w.TotalWithdrawn += txEntry.Amount
	case "withdrawal_reversal":
		// Reserved is released back to available
		if w.ReservedBalance < txEntry.Amount {
			return fmt.Errorf("insufficient reserved balance for reversal: requested %.2f, reserved %.2f", txEntry.Amount, w.ReservedBalance)
		}
		w.ReservedBalance -= txEntry.Amount
		w.AvailableBalance += txEntry.Amount
	case "manual_compensation":
		if txEntry.CreditAccount == "user_available" {
			w.AvailableBalance += txEntry.Amount
		} else if txEntry.DebitAccount == "user_available" {
			if w.AvailableBalance < txEntry.Amount {
				return fmt.Errorf("insufficient available balance for compensation debit: requested %.2f, available %.2f", txEntry.Amount, w.AvailableBalance)
			}
			w.AvailableBalance -= txEntry.Amount
		}
	default:
		return fmt.Errorf("unsupported ledger entry type: %s", txEntry.EntryType)
	}

	// 4. Update wallet account row
	_, err = tx.Exec(ctx, `
		UPDATE wallet_accounts
		SET available_balance = $1, pending_balance = $2, reserved_balance = $3,
		    total_withdrawn = $4, lifetime_earned = $5, updated_at = NOW()
		WHERE id = $6
	`, w.AvailableBalance, w.PendingBalance, w.ReservedBalance, w.TotalWithdrawn, w.LifetimeEarned, w.ID)
	if err != nil {
		return fmt.Errorf("failed to update wallet balance: %w", err)
	}

	// 5. Insert ledger entry
	var insertedID int64
	err = tx.QueryRow(ctx, `
		INSERT INTO ledger_transactions (
			ledger_tx_id, user_id, debit_account, credit_account, amount, currency,
			entry_type, reference_id, source_event, status, idempotency_key, notes, actor_id, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, NOW())
		RETURNING id
	`, txEntry.LedgerTxID, txEntry.UserID, txEntry.DebitAccount, txEntry.CreditAccount,
		txEntry.Amount, txEntry.Currency, txEntry.EntryType, txEntry.ReferenceID,
		txEntry.SourceEvent, txEntry.Status, txEntry.IdempotencyKey, txEntry.Notes, txEntry.ActorID).Scan(&insertedID)
	if err != nil {
		return fmt.Errorf("failed to insert ledger record: %w", err)
	}
	txEntry.ID = insertedID

	return tx.Commit(ctx)
}

// GetLedgerTransactions retrieves paginated transactions for a specific user.
func (r *LedgerRepo) GetLedgerTransactions(ctx context.Context, userID int, limit, offset int) ([]models.LedgerTransaction, int, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	var total int
	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM ledger_transactions WHERE user_id = $1`, userID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.db.Query(ctx, `
		SELECT id, ledger_tx_id, user_id, debit_account, credit_account, amount, currency,
		       entry_type, reference_id, source_event, status, idempotency_key, notes, actor_id, created_at
		FROM ledger_transactions
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []models.LedgerTransaction
	for rows.Next() {
		var t models.LedgerTransaction
		if err := rows.Scan(
			&t.ID,
			&t.LedgerTxID,
			&t.UserID,
			&t.DebitAccount,
			&t.CreditAccount,
			&t.Amount,
			&t.Currency,
			&t.EntryType,
			&t.ReferenceID,
			&t.SourceEvent,
			&t.Status,
			&t.IdempotencyKey,
			&t.Notes,
			&t.ActorID,
			&t.CreatedAt,
		); err != nil {
			return nil, 0, err
		}
		list = append(list, t)
	}

	return list, total, nil
}

// GetAllLedgerTransactions retrieves system-wide audit transactions for administrators.
func (r *LedgerRepo) GetAllLedgerTransactions(ctx context.Context, limit, offset int) ([]models.LedgerTransaction, int, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	var total int
	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM ledger_transactions`).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.db.Query(ctx, `
		SELECT id, ledger_tx_id, user_id, debit_account, credit_account, amount, currency,
		       entry_type, reference_id, source_event, status, idempotency_key, notes, actor_id, created_at
		FROM ledger_transactions
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []models.LedgerTransaction
	for rows.Next() {
		var t models.LedgerTransaction
		if err := rows.Scan(
			&t.ID,
			&t.LedgerTxID,
			&t.UserID,
			&t.DebitAccount,
			&t.CreditAccount,
			&t.Amount,
			&t.Currency,
			&t.EntryType,
			&t.ReferenceID,
			&t.SourceEvent,
			&t.Status,
			&t.IdempotencyKey,
			&t.Notes,
			&t.ActorID,
			&t.CreatedAt,
		); err != nil {
			return nil, 0, err
		}
		list = append(list, t)
	}

	return list, total, nil
}

// GetActiveRewardPolicy fetches the current active referral policy.
func (r *LedgerRepo) GetActiveRewardPolicy(ctx context.Context) (*models.ReferralRewardPolicy, error) {
	var p models.ReferralRewardPolicy
	query := `
		SELECT id, policy_name, registration_reward_amount, transaction_reward_percent,
		       min_qualifying_tx_amount, max_reward_per_tx, max_lifetime_per_referred,
		       max_rewards_per_inviter, reward_holding_days, is_active, version,
		       effective_from, effective_to, updated_by, created_at, updated_at
		FROM referral_reward_policies
		WHERE is_active = true
		ORDER BY id DESC
		LIMIT 1
	`
	err := r.db.QueryRow(ctx, query).Scan(
		&p.ID,
		&p.PolicyName,
		&p.RegistrationRewardAmount,
		&p.TransactionRewardPercent,
		&p.MinQualifyingTxAmount,
		&p.MaxRewardPerTx,
		&p.MaxLifetimePerReferred,
		&p.MaxRewardsPerInviter,
		&p.RewardHoldingDays,
		&p.IsActive,
		&p.Version,
		&p.EffectiveFrom,
		&p.EffectiveTo,
		&p.UpdatedBy,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// UpdateRewardPolicy modifies the policy and writes an immutable audit record to policy versions.
func (r *LedgerRepo) UpdateRewardPolicy(ctx context.Context, p *models.ReferralRewardPolicy, changedBy int, changeReason string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// 1. Snapshot current version
	var cur models.ReferralRewardPolicy
	err = tx.QueryRow(ctx, `
		SELECT id, version, registration_reward_amount, transaction_reward_percent, min_qualifying_tx_amount, max_reward_per_tx, reward_holding_days
		FROM referral_reward_policies
		WHERE id = $1
		FOR UPDATE
	`, p.ID).Scan(&cur.ID, &cur.Version, &cur.RegistrationRewardAmount, &cur.TransactionRewardPercent, &cur.MinQualifyingTxAmount, &cur.MaxRewardPerTx, &cur.RewardHoldingDays)
	if err != nil {
		return err
	}

	// 2. Insert into referral_reward_policy_versions
	_, err = tx.Exec(ctx, `
		INSERT INTO referral_reward_policy_versions (
			policy_id, version, registration_reward_amount, transaction_reward_percent,
			min_qualifying_tx_amount, max_reward_per_tx, reward_holding_days, changed_by, change_reason, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW())
	`, cur.ID, cur.Version, cur.RegistrationRewardAmount, cur.TransactionRewardPercent, cur.MinQualifyingTxAmount, cur.MaxRewardPerTx, cur.RewardHoldingDays, changedBy, changeReason)
	if err != nil {
		return fmt.Errorf("failed to save policy audit snapshot: %w", err)
	}

	// 3. Update policy with incremented version
	newVersion := cur.Version + 1
	_, err = tx.Exec(ctx, `
		UPDATE referral_reward_policies
		SET registration_reward_amount = $1, transaction_reward_percent = $2,
		    min_qualifying_tx_amount = $3, max_reward_per_tx = $4,
		    reward_holding_days = $5, version = $6, updated_by = $7, updated_at = NOW()
		WHERE id = $8
	`, p.RegistrationRewardAmount, p.TransactionRewardPercent, p.MinQualifyingTxAmount,
		p.MaxRewardPerTx, p.RewardHoldingDays, newVersion, changedBy, p.ID)
	if err != nil {
		return fmt.Errorf("failed to update policy: %w", err)
	}

	return tx.Commit(ctx)
}
