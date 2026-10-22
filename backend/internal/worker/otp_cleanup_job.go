package jobs

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func StartOTPCleanup(ctx context.Context, db *pgxpool.Pool) {
	go func() {
		// Wait 10 minutes after startup — lower priority than expiry and user cleanup
		select {
		case <-time.After(10 * time.Minute):
		case <-ctx.Done():
			return
		}
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()
		for {
			runOTPCleanup(db)
			select {
			case <-ticker.C:
			case <-ctx.Done():
				log.Println("[otp-cleanup] shutting down") //nolint:gosec
				return
			}
		}
	}()
}

func runOTPCleanup(db *pgxpool.Pool) {
	jobID := logJobStart(db, "otp_cleanup")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Delete OTP rows that are either verified or expired and older than 7 days
	tag, err := db.Exec(ctx,
		`DELETE FROM otp_codes
		 WHERE ((status IN ('verified', 'expired') OR expires_at < NOW())
		   AND created_at < NOW() - INTERVAL '7 days')
		    OR (status = 'locked' AND created_at < NOW() - INTERVAL '24 hours')`)
	if err != nil {
		log.Printf("[otp-cleanup] otp_codes cleanup failed: %v", err) //nolint:gosec
		logJobFailed(db, jobID, err.Error())
		return
	}
	otpCount := int(tag.RowsAffected())
	if otpCount > 0 {
		log.Printf("[otp-cleanup] deleted %d stale OTP rows", otpCount) //nolint:gosec
	}

	// Delete user_activity_logs older than 90 days
	tag, err = db.Exec(ctx,
		`DELETE FROM user_activity_logs
		 WHERE created_at < NOW() - INTERVAL '90 days'`)
	if err != nil {
		log.Printf("[otp-cleanup] activity_logs cleanup failed: %v", err) //nolint:gosec
		logJobFailed(db, jobID, err.Error())
		return
	}
	activityCount := int(tag.RowsAffected())
	if activityCount > 0 {
		log.Printf("[otp-cleanup] deleted %d old activity log rows", activityCount) //nolint:gosec
	}
	
	logJobSuccess(db, jobID, otpCount+activityCount)
}
