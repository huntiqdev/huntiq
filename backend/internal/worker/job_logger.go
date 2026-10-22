package jobs

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// logJobStart logs the start of a job and returns the job log ID
func logJobStart(db *pgxpool.Pool, jobName string) int {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	var jobID int
	err := db.QueryRow(ctx, //nolint:gosec
		`INSERT INTO job_execution_log (job_name, status, started_at)
		 VALUES ($1, 'running', NOW())
		 RETURNING id`,
		jobName).Scan(&jobID)
	if err != nil {
		// Logging failure should not break the job
		return 0
	}
	return jobID
}

// logJobSuccess marks a job as successful
func logJobSuccess(db *pgxpool.Pool, jobID int, recordsProcessed int) {
	if jobID == 0 {
		return
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	_, _ = db.Exec(ctx, //nolint:gosec
		`UPDATE job_execution_log
		 SET status = 'success',
		     completed_at = NOW(),
		     duration_ms = EXTRACT(EPOCH FROM (NOW() - started_at)) * 1000,
		     records_processed = $2
		 WHERE id = $1`,
		jobID, recordsProcessed)
}

// logJobFailed marks a job as failed with an error message
func logJobFailed(db *pgxpool.Pool, jobID int, errorMsg string) {
	if jobID == 0 {
		return
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	_, _ = db.Exec(ctx, //nolint:gosec
		`UPDATE job_execution_log
		 SET status = 'failed',
		     completed_at = NOW(),
		     duration_ms = EXTRACT(EPOCH FROM (NOW() - started_at)) * 1000,
		     error_message = $2
		 WHERE id = $1`,
		jobID, errorMsg)
}
