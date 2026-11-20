package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	FilenameRegex = regexp.MustCompile(`^[0-9]{3}_[a-z0-9_]+\.sql$`)

	createTableRegex = regexp.MustCompile(`(?i)^\s*CREATE\s+(?:UNLOGGED\s+)?TABLE\s+`)
	dropTableRegex   = regexp.MustCompile(`(?i)^\s*DROP\s+TABLE\s+`)
	createIndexRegex = regexp.MustCompile(`(?i)^\s*CREATE\s+(?:UNIQUE\s+)?INDEX\s+`)
	dropIndexRegex   = regexp.MustCompile(`(?i)^\s*DROP\s+INDEX\s+`)

	ifNotExistsRegex = regexp.MustCompile(`(?i)IF\s+NOT\s+EXISTS`)
	ifExistsRegex    = regexp.MustCompile(`(?i)IF\s+EXISTS`)

	nonTxConcurrentRegex = regexp.MustCompile(`(?i)(?:CREATE|DROP)\s+(?:UNIQUE\s+)?INDEX\s+CONCURRENTLY`)
	nonTxVacuumRegex     = regexp.MustCompile(`(?i)\bVACUUM\b`)
	nonTxClusterRegex    = regexp.MustCompile(`(?i)\bCLUSTER\s+[a-zA-Z0-9_]+\s+USING\b`)
)

type MigrationValidationResult struct {
	Valid                 bool
	TotalFiles            int
	HistoricalCount       int
	NewCount              int
	AllowlistedCollisions []string
	NonTransactionalOps   []string
	Warnings              []string
	Errors                []string
}

func IsNonTransactionalStatement(sql string) bool {
	return nonTxConcurrentRegex.MatchString(sql) ||
		nonTxVacuumRegex.MatchString(sql) ||
		nonTxClusterRegex.MatchString(sql)
}

func ValidateMigrationFilenames(filenames []string) *MigrationValidationResult {
	res := &MigrationValidationResult{
		Valid:                 true,
		AllowlistedCollisions: make([]string, 0),
		NonTransactionalOps:   make([]string, 0),
		Warnings:              make([]string, 0),
		Errors:                make([]string, 0),
	}

	if len(filenames) == 0 {
		res.Valid = false
		res.Errors = append(res.Errors, "no migration files provided")
		return res
	}

	sortedFiles := make([]string, len(filenames))
	copy(sortedFiles, filenames)
	sort.Strings(sortedFiles)
	res.TotalFiles = len(sortedFiles)

	seenPrefixes := make(map[string][]string)
	prevPrefixInt := -1

	for _, filename := range sortedFiles {
		if !FilenameRegex.MatchString(filename) {
			res.Valid = false
			res.Errors = append(res.Errors, fmt.Sprintf("[FORMAT] invalid filename: %q, must match '^[0-9]{3}_[a-z0-9_]+\\.sql$'", filename))
			continue
		}

		prefix := filename[:3]
		prefixInt, err := strconv.Atoi(prefix)
		if err != nil {
			res.Valid = false
			res.Errors = append(res.Errors, fmt.Sprintf("[FORMAT] non-numeric prefix in filename: %q", filename))
			continue
		}

		if prefixInt <= 49 {
			res.HistoricalCount++
		} else {
			res.NewCount++
		}

		existing := seenPrefixes[prefix]
		if len(existing) > 0 {
			res.Valid = false
			res.Errors = append(res.Errors, fmt.Sprintf("[DUPLICATE] unapproved duplicate prefix %q found in %v and %q", prefix, existing, filename))
			seenPrefixes[prefix] = append(existing, filename)
		} else {
			seenPrefixes[prefix] = []string{filename}
		}

		if prevPrefixInt >= 0 {
			if prefixInt < prevPrefixInt {
				res.Valid = false
				res.Errors = append(res.Errors, fmt.Sprintf("[ORDER] sequence inversion: %q appears after prefix %d", filename, prevPrefixInt))
			} else if prefixInt > prevPrefixInt+1 {
				res.Valid = false
				res.Errors = append(res.Errors, fmt.Sprintf("[GAP] missing migration sequence between %03d and %03d (expected %03d)", prevPrefixInt, prefixInt, prevPrefixInt+1))
			}
		}
		prevPrefixInt = prefixInt
	}

	return res
}

func ValidateMigrationSQL(filename string, content string) ([]string, []string, []string) {
	var errs, warns, nonTx []string
	prefixInt := 0
	if len(filename) >= 3 {
		prefixInt, _ = strconv.Atoi(filename[:3])
	}

	lines := strings.Split(content, "\n")
	var strippedLines []string
	for _, l := range lines {
		trimmed := strings.TrimSpace(l)
		if strings.HasPrefix(trimmed, "--") {
			continue
		}
		if idx := strings.Index(trimmed, "--"); idx >= 0 {
			trimmed = strings.TrimSpace(trimmed[:idx])
		}
		if trimmed != "" {
			strippedLines = append(strippedLines, trimmed)
		}
	}
	cleanContent := strings.Join(strippedLines, "\n")

	if nonTxConcurrentRegex.MatchString(cleanContent) {
		nonTx = append(nonTx, fmt.Sprintf("[NON-TX] %q: CONCURRENT index", filename))
	}
	if nonTxVacuumRegex.MatchString(cleanContent) {
		nonTx = append(nonTx, fmt.Sprintf("[NON-TX] %q: VACUUM statement", filename))
	}

	for _, line := range strippedLines {
		if createTableRegex.MatchString(line) && !ifNotExistsRegex.MatchString(line) {
			if prefixInt >= 77 {
				errs = append(errs, fmt.Sprintf("[IDEMPOTENCY] %q: table creation without IF NOT EXISTS: %s", filename, line))
			} else {
				warns = append(warns, fmt.Sprintf("[HISTORICAL] %q: historical table statement without IF NOT EXISTS: %s", filename, line))
			}
		}
		if dropTableRegex.MatchString(line) && !ifExistsRegex.MatchString(line) {
			if prefixInt >= 77 {
				errs = append(errs, fmt.Sprintf("[IDEMPOTENCY] %q: drop table without IF EXISTS: %s", filename, line))
			} else {
				warns = append(warns, fmt.Sprintf("[HISTORICAL] %q: historical drop table without IF EXISTS: %s", filename, line))
			}
		}
	}

	return errs, warns, nonTx
}

func ValidateMigrationsDirectory(dirPath string) (*MigrationValidationResult, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("read migrations directory %s: %w", dirPath, err)
	}

	var sqlFiles []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sql") {
			sqlFiles = append(sqlFiles, entry.Name())
		}
	}

	res := ValidateMigrationFilenames(sqlFiles)

	for _, filename := range sqlFiles {
		filePath := filepath.Join(dirPath, filename)
		content, err := os.ReadFile(filePath)
		if err != nil {
			res.Valid = false
			res.Errors = append(res.Errors, fmt.Sprintf("cannot read file %s: %v", filename, err))
			continue
		}
		sqlErrs, sqlWarns, sqlNonTx := ValidateMigrationSQL(filename, string(content))
		if len(sqlErrs) > 0 {
			res.Valid = false
			res.Errors = append(res.Errors, sqlErrs...)
		}
		res.Warnings = append(res.Warnings, sqlWarns...)
		res.NonTransactionalOps = append(res.NonTransactionalOps, sqlNonTx...)
	}

	return res, nil
}

func LocateMigrationsDir(custom string) string {
	if custom != "" {
		if info, err := os.Stat(custom); err == nil && info.IsDir() {
			return custom
		}
	}
	candidates := []string{
		"internal/database/migrations",
		"backend/internal/database/migrations",
		"../internal/database/migrations",
		"../../internal/database/migrations",
		"migrations",
	}
	for _, c := range candidates {
		if info, err := os.Stat(c); err == nil && info.IsDir() {
			return c
		}
	}
	return "internal/database/migrations"
}

// RunMigrations validates and executes all pending migrations with auto-recovery.
func RunMigrations(ctx context.Context, pool *pgxpool.Pool, customDir string) error {
	migrationsDir := LocateMigrationsDir(customDir)
	log.Printf("[HUNTIQ DB] Validating migrations in %s...", migrationsDir)

	valResult, err := ValidateMigrationsDirectory(migrationsDir)
	if err != nil {
		return fmt.Errorf("migration validation failed: %w", err)
	}
	if !valResult.Valid {
		return fmt.Errorf("migration sequence validation errors: %v", valResult.Errors)
	}

	// Ensure schema_migrations table exists with auto-recovery
	var lastErr error
	for attempt := 1; attempt <= 3; attempt++ {
		_, lastErr = pool.Exec(ctx, `
			CREATE TABLE IF NOT EXISTS schema_migrations (
				version VARCHAR(255) PRIMARY KEY,
				applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
			);
		`)
		if lastErr == nil {
			break
		}
		log.Printf("[HUNTIQ DB] Warning: schema_migrations init attempt %d failed: %v, retrying...", attempt, lastErr)
		time.Sleep(time.Duration(attempt) * 500 * time.Millisecond)
	}
	if lastErr != nil {
		return fmt.Errorf("initialize schema_migrations: %w", lastErr)
	}

	rows, err := pool.Query(ctx, `SELECT version FROM schema_migrations;`)
	if err != nil {
		return fmt.Errorf("query applied migrations: %w", err)
	}
	defer rows.Close()

	applied := make(map[string]bool)
	for rows.Next() {
		var v string
		if err := rows.Scan(&v); err == nil {
			applied[v] = true
		}
	}

	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("read migrations dir: %w", err)
	}

	var pending []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".sql") {
			if !applied[e.Name()] {
				pending = append(pending, e.Name())
			}
		}
	}
	sort.Strings(pending)

	if len(pending) == 0 {
		log.Printf("[HUNTIQ DB] ✓ Migrations up to date (%d verified, 0 pending).", valResult.TotalFiles)
		return nil
	}

	log.Printf("[HUNTIQ DB] Applying %d pending migration(s)...", len(pending))
	for _, filename := range pending {
		sqlBytes, err := os.ReadFile(filepath.Join(migrationsDir, filename))
		if err != nil {
			return fmt.Errorf("read migration %s: %w", filename, err)
		}
		sqlContent := string(sqlBytes)

		// Execute with auto-recovery / retry
		var execErr error
		for attempt := 1; attempt <= 3; attempt++ {
			if IsNonTransactionalStatement(sqlContent) {
				_, execErr = pool.Exec(ctx, sqlContent)
			} else {
				tx, err := pool.Begin(ctx)
				if err != nil {
					execErr = err
				} else {
					if _, err := tx.Exec(ctx, sqlContent); err != nil {
						_ = tx.Rollback(ctx)
						execErr = err
					} else {
						execErr = tx.Commit(ctx)
					}
				}
			}
			if execErr == nil {
				break
			}
			log.Printf("[HUNTIQ DB] Retrying migration %s (attempt %d/3) due to: %v", filename, attempt, execErr)
			time.Sleep(time.Duration(attempt) * 500 * time.Millisecond)
		}
		if execErr != nil {
			return fmt.Errorf("migration %s failed: %w", filename, execErr)
		}

		_, _ = pool.Exec(ctx, `INSERT INTO schema_migrations (version, applied_at) VALUES ($1, NOW()) ON CONFLICT (version) DO NOTHING`, filename)
		log.Printf("[HUNTIQ DB]   ✓ %s applied", filename)
	}

	log.Printf("[HUNTIQ DB] ✓ All %d migrations successfully applied and verified.", valResult.TotalFiles)
	return nil
}
