package repository

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/your-org/used-tools-marketplace/internal/models"
	"github.com/your-org/used-tools-marketplace/pkg/constants"
)

type TransactionRepo struct {
	db            *pgxpool.Pool
	settingsCache *SettingsCache
}

func NewTransactionRepo(db *pgxpool.Pool) *TransactionRepo {
	return &TransactionRepo{
		db:            db,
		settingsCache: NewSettingsCache(db, 1*time.Minute), // 1-minute TTL for settings
	}
}

// CreateAndSell atomically:
//  1. Locks the product row with SELECT FOR UPDATE to prevent concurrent purchases.
//  2. Verifies the product is still 'approved' (re-checks inside the transaction).
//  3. Reads broker_fee_rate from platform_settings (falls back to 2.0% if missing).
//  4. Inserts the transaction row with computed broker_fee and seller_gets.
//  5. Marks the product as 'sold'.
//
// All steps run inside a single DB transaction — either all succeed or all roll back.
func (r *TransactionRepo) CreateAndSell(ctx context.Context, sellerID int, t *models.Transaction) (int, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	// Step 1 — lock the product row so no concurrent request can also buy it.
	var productSellerID int
	var status string
	var expiresAt *time.Time
	var productPrice float64
	//nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
	err = tx.QueryRow(ctx, `
		SELECT seller_id, status, expires_at, price
		FROM products
		WHERE id = $1
		FOR UPDATE`, t.ProductID,
	).Scan(&productSellerID, &status, &expiresAt, &productPrice)
	if err != nil {
		return 0, fmt.Errorf("product not found: %w", err)
	}

	// Security: Verify the seller recording this transaction actually owns the product.
	if sellerID > 0 && productSellerID != sellerID {
		return 0, errors.New("product does not belong to you")
	}

	// Step 2 — re-verify status inside the lock.
	if status != "approved" {
		return 0, errors.New("product is no longer available for purchase")
	}
	if expiresAt != nil && expiresAt.Before(time.Now()) {
		return 0, errors.New("this listing has expired")
	}
	
	// Read currency tolerance as PERCENTAGE from settings (admin-configurable)
	// Default 0.1% tolerance for floating-point rounding errors
	toleranceStr := r.settingsCache.Get(ctx, constants.SettingCurrencyTolerancePercent)
	tolerancePercent := 0.001 // 0.1% default
	if toleranceStr != "" {
		if v, err := strconv.ParseFloat(toleranceStr, 64); err == nil && v > 0 && v < 100 {
			tolerancePercent = v / 100.0 // Convert percentage to decimal (e.g., 0.1 → 0.001)
		}
	}
	
	// Calculate tolerance as percentage of product price (prevents manipulation)
	maxAllowedDifference := productPrice * tolerancePercent
	
	// Strict price match — percentage-based tolerance prevents underpayment attacks
	if math.Abs(t.Amount-productPrice) > maxAllowedDifference {
		return 0, fmt.Errorf(
			"transaction amount %.2f does not match the listed price %.2f (tolerance: %.2f%%)",
			t.Amount, productPrice, tolerancePercent*100,
		)
	}

	// Step 3 — read broker_fee_rate from cache (outside transaction lock for performance).
	// Cache has 1-minute TTL, so admin changes take effect within 60 seconds max.
	// This prevents all concurrent purchases from serializing on a single settings table scan.
	rateStr := r.settingsCache.Get(ctx, "broker_fee_rate")
	rate := 2.0
	if rateStr != "" {
		if v, err := strconv.ParseFloat(rateStr, 64); err == nil && v > 0 && v <= 100 {
			rate = v
		}
	}
	t.BrokerFee = math.Round(t.Amount*(rate/100)*100) / 100
	t.SellerGets = math.Round((t.Amount-t.BrokerFee)*100) / 100

	// Step 4 — insert transaction with computed fees and seller_id for audit trail.
	// Always attribute the transaction to the product's true verified seller.
	//nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
	err = tx.QueryRow(ctx, `
		INSERT INTO transactions (product_id, seller_id, buyer_name, buyer_phone, amount, broker_fee, seller_gets, notes)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, broker_fee, seller_gets`,
		t.ProductID, productSellerID, t.BuyerName, t.BuyerPhone, t.Amount, t.BrokerFee, t.SellerGets, t.Notes,
	).Scan(&t.ID, &t.BrokerFee, &t.SellerGets)
	if err != nil {
		return 0, err
	}

	// Step 5 — mark product sold.
	//nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
	_, err = tx.Exec(ctx, `
		UPDATE products SET status='sold', updated_at=NOW()
		WHERE id = $1`, t.ProductID)
	if err != nil {
		return 0, err
	}

	return productSellerID, tx.Commit(ctx)
}

// UpdateStatus changes a transaction's status. Valid transitions:
// pending → completed | cancelled
func (r *TransactionRepo) UpdateStatus(ctx context.Context, txID int, status string) error {
	tag, err := r.db.Exec(ctx, //nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
		`UPDATE transactions SET status=$1::transaction_status, updated_at=NOW()
		 WHERE id=$2 AND status='pending'::transaction_status`,
		status, txID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("transaction %d not found or already processed (only pending transactions can be updated)", txID)
	}
	return nil
}
