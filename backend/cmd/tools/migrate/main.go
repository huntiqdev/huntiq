package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/your-org/used-tools-marketplace/internal/config"
	"github.com/your-org/used-tools-marketplace/internal/database"
)

func main() {
	_ = godotenv.Load()
	dryRun := flag.Bool("dry-run", false, "Preview migrations without executing them")
	dirFlag := flag.String("dir", "", "Custom path to migrations directory")
	flag.Parse()

	// Locate migrations directory
	migrationsDir := *dirFlag
	if migrationsDir == "" {
		candidates := []string{
			"internal/database/migrations",
			"backend/internal/database/migrations",
			"../internal/database/migrations",
		}
		for _, c := range candidates {
			if info, err := os.Stat(c); err == nil && info.IsDir() {
				migrationsDir = c
				break
			}
		}
	}

	if migrationsDir == "" {
		log.Fatalf("Could not locate migrations directory. Use -dir flag.")
	}

	fmt.Printf("=== HUNTIQ Native Go Migration Tool ===\n")
	fmt.Printf("Migrations Directory: %s\n", migrationsDir)

	// Step 1: Pre-execution validation using Go validator
	fmt.Println("--> Running pre-execution sequence and idempotency validation...")
	valResult, err := database.ValidateMigrationsDirectory(migrationsDir)
	if err != nil {
		log.Fatalf("Validator error: %v", err)
	}

	if !valResult.Valid {
		fmt.Fprintf(os.Stderr, "ERROR: Migration validation failed with %d error(s):\n", len(valResult.Errors))
		for _, e := range valResult.Errors {
			fmt.Fprintf(os.Stderr, "  - %s\n", e)
		}
		os.Exit(1)
	}
	fmt.Printf("--> Validation passed: %d migration files verified (0 collisions, strict order).\n", valResult.TotalFiles)

	// Step 2: Database connection
	cfg := config.Load()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	pool, err := database.NewPoolNoSchema(cfg)
	if err != nil {
		log.Fatalf("Database connection failed (%s:%s): %v", cfg.DBHost, cfg.DBPort, err)
	}
	defer pool.Close()

	// Step 3: Ensure schema_migrations table exists
	_, err = pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
	`)
	if err != nil {
		log.Fatalf("Failed to ensure schema_migrations table: %v", err)
	}

	// Step 4: Query applied migrations
	rows, err := pool.Query(ctx, `SELECT version FROM schema_migrations;`)
	if err != nil {
		log.Fatalf("Failed to query applied migrations: %v", err)
	}
	defer rows.Close()

	applied := make(map[string]bool)
	for rows.Next() {
		var ver string
		if err := rows.Scan(&ver); err == nil {
			applied[ver] = true
		}
	}

	// Step 5: Read and sort migration files
	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		log.Fatalf("Failed to read migrations directory: %v", err)
	}

	var sqlFiles []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".sql") {
			sqlFiles = append(sqlFiles, e.Name())
		}
	}
	sort.Strings(sqlFiles)

	appliedCount := 0
	skippedCount := 0

	for _, fname := range sqlFiles {
		if applied[fname] {
			fmt.Printf("  [SKIPPED] %s (already applied)\n", fname)
			skippedCount++
			continue
		}

		if *dryRun {
			fmt.Printf("  [WOULD APPLY] %s\n", fname)
			appliedCount++
			continue
		}

		fmt.Printf("  [APPLYING] %s ...\n", fname)
		contentBytes, err := os.ReadFile(filepath.Join(migrationsDir, fname))
		if err != nil {
			log.Fatalf("Failed to read migration file %s: %v", fname, err)
		}
		content := string(contentBytes)

		hasExplicitTx := strings.Contains(strings.ToUpper(content), "BEGIN;") || strings.Contains(strings.ToUpper(content), "COMMIT;")

		if database.IsNonTransactionalStatement(content) || hasExplicitTx {
			// Non-transactional statements or migrations with explicit BEGIN/COMMIT run directly on pool
			if _, err := pool.Exec(ctx, content); err != nil {
				log.Fatalf("Error executing migration %s: %v", fname, err)
			}
			if _, err := pool.Exec(ctx, `INSERT INTO schema_migrations (version) VALUES ($1);`, fname); err != nil {
				log.Fatalf("Error recording migration %s: %v", fname, err)
			}
		} else {
			// Run transactional migrations inside a transaction
			tx, err := pool.Begin(ctx)
			if err != nil {
				log.Fatalf("Failed to begin transaction for %s: %v", fname, err)
			}
			if _, err := tx.Exec(ctx, content); err != nil {
				_ = tx.Rollback(ctx)
				log.Fatalf("Error executing migration %s: %v", fname, err)
			}
			if _, err := tx.Exec(ctx, `INSERT INTO schema_migrations (version) VALUES ($1);`, fname); err != nil {
				_ = tx.Rollback(ctx)
				log.Fatalf("Error recording migration %s: %v", fname, err)
			}
			if err := tx.Commit(ctx); err != nil {
				log.Fatalf("Error committing migration %s: %v", fname, err)
			}
		}

		fmt.Printf("  [DONE]     %s\n", fname)
		appliedCount++
	}

	fmt.Printf("\n=== Migration Finished ===\n")
	fmt.Printf("Total Files: %d | Applied: %d | Skipped: %d\n", len(sqlFiles), appliedCount, skippedCount)
}
