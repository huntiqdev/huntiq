package repository

import (
	"context"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// AsyncActivityLog represents a user activity log entry for async writing.
type AsyncActivityLog struct {
	Event     string
	Phone     string
	Email     string
	UserType  string
	IPAddress string
	Meta      string
}

// AsyncLogger provides non-blocking activity log writes with batching.
// Logs are buffered in-memory and flushed every second or when batch size reaches 100.
type AsyncLogger struct {
	db      *pgxpool.Pool
	logChan chan AsyncActivityLog
	done    chan struct{}
}

// NewAsyncLogger creates an async logger with a 1000-entry buffer.
// A background worker flushes logs to DB every second or every 100 entries.
func NewAsyncLogger(db *pgxpool.Pool) *AsyncLogger {
	al := &AsyncLogger{
		db:      db,
		logChan: make(chan AsyncActivityLog, 1000),
		done:    make(chan struct{}),
	}
	go al.worker()
	return al
}

// Log sends a log entry to the background worker (non-blocking).
// If the buffer is full, the oldest entry is dropped (back-pressure).
func (al *AsyncLogger) Log(entry AsyncActivityLog) {
	select {
	case al.logChan <- entry:
		// Successfully queued
	default:
		// Buffer full - drop log to prevent blocking critical path
		// In production, increment a dropped_logs metric here
		log.Printf("[WARN] Activity log buffer full, dropping entry for event %s", entry.Event)
	}
}

// Close gracefully shuts down the logger, flushing remaining logs.
func (al *AsyncLogger) Close() {
	close(al.logChan)
	<-al.done // Wait for worker to finish
}

// worker runs in background, batching and flushing logs to DB.
func (al *AsyncLogger) worker() {
	defer close(al.done)
	
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	
	batch := make([]AsyncActivityLog, 0, 100)
	
	for {
		select {
		case entry, ok := <-al.logChan:
			if !ok {
				// Channel closed - flush remaining and exit
				if len(batch) > 0 {
					al.flushBatch(batch)
				}
				return
			}
			
			batch = append(batch, entry)
			if len(batch) >= 100 {
				al.flushBatch(batch)
				batch = batch[:0]
			}
			
		case <-ticker.C:
			// Periodic flush (every second)
			if len(batch) > 0 {
				al.flushBatch(batch)
				batch = batch[:0]
			}
		}
	}
}

// flushBatch inserts a batch of log entries using a single multi-row INSERT.
// The placeholder string is built purely from integer arithmetic — no user
// data ever enters the SQL string, only the $N indices.
func (al *AsyncLogger) flushBatch(batch []AsyncActivityLog) {
	if len(batch) == 0 {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Pre-build the VALUES clause from integer offsets only.
	// Separated from args population so no tainted value is in scope
	// when the SQL string is constructed.
	n := len(batch)
	placeholders := make([]string, n)
	for i := range placeholders {
		base := i * 6
		placeholders[i] = "($" + strconv.Itoa(base+1) + ",$" + strconv.Itoa(base+2) +
			",$" + strconv.Itoa(base+3) + ",$" + strconv.Itoa(base+4) +
			",$" + strconv.Itoa(base+5) + ",$" + strconv.Itoa(base+6) + ")"
	}
	sql := `INSERT INTO user_activity_logs (event, phone, email, user_type, ip_address, meta) VALUES ` +
		strings.Join(placeholders, ",") //nolint:gosec // G201: placeholders built from static template with computed $N indices using strconv.Itoa

	// Populate args separately — user values never touch the SQL string above.
	args := make([]interface{}, 0, n*6)
	for _, entry := range batch {
		var phone, email, ip, meta interface{}
		if entry.Phone != "" {
			phone = entry.Phone
		}
		if entry.Email != "" {
			email = entry.Email
		}
		if entry.IPAddress != "" {
			ip = entry.IPAddress
		}
		if entry.Meta != "" {
			meta = entry.Meta
		}
		args = append(args, entry.Event, phone, email, entry.UserType, ip, meta)
	}

	if _, err := al.db.Exec(ctx, sql, args...); err != nil {
		log.Printf("[ERROR] Failed to flush activity logs batch (size=%d)", len(batch))
	}
}

