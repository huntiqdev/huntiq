package jobs

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func StartUnverifiedUserCleanup(ctx context.Context, db *pgxpool.Pool) {
	go func() {
		// Wait 5 minutes after startup before first run
		select {
		case <-time.After(5 * time.Minute):
		case <-ctx.Done():
			return
		}
		ticker := time.NewTicker(6 * time.Hour)
		defer ticker.Stop()
		for {
			func() {
				jobID := logJobStart(db, "unverified_user_cleanup")
				runCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
				defer cancel()
				var deleted int
				err := db.QueryRow(runCtx, `SELECT cleanup_unverified_users()`).Scan(&deleted) //nolint:gosec
				if err != nil {
					log.Printf("[cleanup] unverified user cleanup failed: %v", err) //nolint:gosec
					logJobFailed(db, jobID, err.Error())
					return
				}
				if deleted > 0 {
					log.Printf("[cleanup] removed %d unverified ghost account(s)", deleted) //nolint:gosec
				}
				logJobSuccess(db, jobID, deleted)
			}()
			select {
			case <-ticker.C:
			case <-ctx.Done():
				log.Println("[cleanup] shutting down") //nolint:gosec
				return
			}
		}
	}()
}
