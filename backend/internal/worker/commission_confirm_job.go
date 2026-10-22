package jobs

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func StartCommissionConfirmJob(ctx context.Context, db *pgxpool.Pool) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[commission-confirm-job] recovered from panic: %v", r) //nolint:gosec
			}
		}()

		// Run immediately on startup to catch any commissions that expired while server was down
		safeConfirmCommissions(db)

		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				safeConfirmCommissions(db)
			case <-ctx.Done():
				log.Println("[commission-confirm-job] shutting down") //nolint:gosec
				return
			}
		}
	}()
}

// safeConfirmCommissions wraps confirmCommissions with panic recovery so a panic
// in one tick doesn't kill the goroutine — the job continues on the next tick.
func safeConfirmCommissions(db *pgxpool.Pool) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[commission-confirm-job] recovered from panic in confirmCommissions: %v", r) //nolint:gosec
		}
	}()
	confirmCommissions(db)
}

func confirmCommissions(db *pgxpool.Pool) {
	jobID := logJobStart(db, "commission_confirmation")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Atomically update pending commissions whose pending_until has passed
	tag, err := db.Exec(ctx, //nolint:gosec
		`UPDATE referral_commissions
		 SET status = 'confirmed', confirmed_at = NOW()
		 WHERE status = 'pending'
		   AND pending_until IS NOT NULL
		   AND pending_until <= NOW()`)
	if err != nil {
		log.Printf("[commission-confirm-job] update error: %v", err) //nolint:gosec
		logJobFailed(db, jobID, err.Error())
		return
	}

	count := int(tag.RowsAffected())
	if count > 0 {
		log.Printf("[commission-confirm-job] confirmed %d pending commission(s)", count) //nolint:gosec
	}
	logJobSuccess(db, jobID, count)
}
