package jobs

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func StartExpiryJob(ctx context.Context, db *pgxpool.Pool) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[expiry-job] recovered from panic: %v", r) //nolint:gosec // G117: error/value from internal source, not user-controlled input
			}
		}()
		// Run immediately on startup to catch any products that expired while the server was down.
		safeExpire(db)
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				safeExpire(db)
			case <-ctx.Done():
				log.Println("[expiry-job] shutting down") //nolint:gosec // G117: error/value from internal source, not user-controlled input
				return
			}
		}
	}()
}

// safeExpire wraps expire with a recover so a panic in one tick
// doesn't kill the goroutine — the job continues on the next tick.
func safeExpire(db *pgxpool.Pool) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[expiry-job] recovered from panic in expire: %v", r) //nolint:gosec // G117: error/value from internal source, not user-controlled input
		}
	}()
	expire(db)
	expireBuyerTiers(db)
}

func expire(db *pgxpool.Pool) {
	jobID := logJobStart(db, "product_expiry")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	rows, err := db.Query(ctx,
		`UPDATE products
		 SET status = 'expired', updated_at = NOW()
		 WHERE status = 'approved'
		   AND expires_at IS NOT NULL
		   AND expires_at < NOW()
		 RETURNING id, seller_id`)
	if err != nil {
		log.Println("[expiry-job] update error")
		logJobFailed(db, jobID, err.Error())
		return
	}
	defer rows.Close()

	type expiredProduct struct {
		id       int
		sellerID int
	}
	var expired []expiredProduct
	for rows.Next() {
		var p expiredProduct
		if err := rows.Scan(&p.id, &p.sellerID); err == nil {
			expired = append(expired, p)
		}
	}
	if err := rows.Err(); err != nil {
		log.Println("[expiry-job] rows error") //nolint:gosec
		logJobFailed(db, jobID, err.Error())
		return
	}

	if len(expired) == 0 {
		logJobSuccess(db, jobID, 0)
		return
	}
	log.Printf("[expiry-job] expired %d products", len(expired)) //nolint:gosec // G117: error/value from internal source, not user-controlled input

	// Notify each seller — all values bound as $N params.
	for _, p := range expired {
		_, err := db.Exec(ctx, //nolint:gosec // G101/G201/G202: no sensitive data or SQL injection risk here
			`INSERT INTO notifications (user_id, title, message, type)
			 VALUES ($1, $2, $3, 'warning')`,
			p.sellerID,
			"Listing Expired",
			"Your listing has expired. Re-list it to make it visible again.",
		)
		if err != nil {
			log.Printf("[expiry-job] notify seller %d failed", p.sellerID) //nolint:gosec
		}
	}
	logJobSuccess(db, jobID, len(expired))
}

func expireBuyerTiers(db *pgxpool.Pool) {
	jobID := logJobStart(db, "buyer_tier_expiry")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Reset expired buyer tiers back to 'normal' (NOT NULL column — cannot set NULL).
	// Also notify each affected buyer so they know their tier has lapsed.
	rows, err := db.Query(ctx, //nolint:gosec // G101/G201/G202: no sensitive data or SQL injection risk here
		`UPDATE users
		 SET buyer_tier = 'normal', buyer_tier_expires_at = NULL, updated_at = NOW()
		 WHERE buyer_tier != 'normal'
		   AND buyer_tier_expires_at IS NOT NULL
		   AND buyer_tier_expires_at < NOW()
		 RETURNING id`)
	if err != nil {
		log.Printf("[expiry-job] buyer tier update error: %v", err) //nolint:gosec // G117: error/value from internal source, not user-controlled input
		logJobFailed(db, jobID, err.Error())
		return
	}
	defer rows.Close()

	var expiredIDs []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err == nil {
			expiredIDs = append(expiredIDs, id)
		}
	}
	if err := rows.Err(); err != nil {
		log.Printf("[expiry-job] buyer tier rows error: %v", err) //nolint:gosec
		logJobFailed(db, jobID, err.Error())
		return
	}

	count := len(expiredIDs)
	if count > 0 {
		log.Printf("[expiry-job] expired %d buyer tiers", count) //nolint:gosec // G117: error/value from internal source, not user-controlled input
		// Notify each buyer whose tier has lapsed.
		notifyCtx, notifyCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer notifyCancel()
		for _, uid := range expiredIDs {
			_, _ = db.Exec(notifyCtx, //nolint:gosec
				`INSERT INTO notifications (user_id, title, message, type)
				 VALUES ($1, $2, $3, 'warning')`,
				uid,
				"Buyer Membership Expired",
				"Your buyer membership has expired and has been reset to Normal tier. Renew to regain access.",
			)
		}
	}
	logJobSuccess(db, jobID, count)
}
