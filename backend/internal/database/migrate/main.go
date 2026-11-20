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
			"../migrations",
			"migrations",
			"../../internal/database/migrations",
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

	// Step 2: Connect to PostgreSQL (Strict Port 5432)
	cfg := config.Load()
	pool, err := database.NewPoolNoSchema(cfg)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer pool.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// Step 3: Ensure schema_migrations table exists
	_, err = pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
	`)
	if err != nil {
		log.Fatalf("Failed to initialize schema_migrations table: %v", err)
	}

	// Step 4: Fetch applied migrations
	rows, err := pool.Query(ctx, `SELECT version FROM schema_migrations;`)
	if err != nil {
		log.Fatalf("Failed to query applied migrations: %v", err)
	}
	defer rows.Close()

	applied := make(map[string]bool)
	for rows.Next() {
		var v string
		if err := rows.Scan(&v); err == nil {
			applied[v] = true
		}
	}

	// Step 5: Read and sort migration files
	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		log.Fatalf("Failed to read migrations directory: %v", err)
	}

	var files []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".sql") {
			files = append(files, e.Name())
		}
	}
	sort.Strings(files)

	// Step 6: Apply unapplied migrations
	var appliedCount int
	var skippedCount int

	for _, fname := range files {
		if applied[fname] {
			skippedCount++
			if *dryRun {
				fmt.Printf("  [SKIPPED] %s (already applied)\n", fname)
			}
			continue
		}

		fmt.Printf("  --> Applying: %s ... ", fname)
		if *dryRun {
			fmt.Println("[DRY-RUN - WOULD APPLY]")
			appliedCount++
			continue
		}

		content, err := os.ReadFile(filepath.Join(migrationsDir, fname))
		if err != nil {
			fmt.Printf("FAILED reading file: %v\n", err)
			os.Exit(1)
		}

		sqlContent := string(content)
		// Non-transactional migrations (e.g. CREATE INDEX CONCURRENTLY, ALTER TYPE ADD VALUE)
		isNonTx := strings.Contains(strings.ToUpper(sqlContent), "CONCURRENTLY") ||
			strings.Contains(strings.ToUpper(sqlContent), "ALTER TYPE")

		if isNonTx {
			// Run directly on pool connection
			if _, err := pool.Exec(ctx, sqlContent); err != nil {
				fmt.Printf("FAILED (non-tx): %v\n", err)
				os.Exit(1)
			}
			_, _ = pool.Exec(ctx, `INSERT INTO schema_migrations (version) VALUES ($1);`, fname)
		} else {
			// Run within an atomic transaction
			tx, err := pool.Begin(ctx)
			if err != nil {
				fmt.Printf("FAILED starting transaction: %v\n", err)
				os.Exit(1)
			}

			if _, err := tx.Exec(ctx, sqlContent); err != nil {
				_ = tx.Rollback(ctx)
				fmt.Printf("FAILED: %v\n", err)
				os.Exit(1)
			}

			if _, err := tx.Exec(ctx, `INSERT INTO schema_migrations (version) VALUES ($1);`, fname); err != nil {
				_ = tx.Rollback(ctx)
				fmt.Printf("FAILED recording version: %v\n", err)
				os.Exit(1)
			}

			if err := tx.Commit(ctx); err != nil {
				fmt.Printf("FAILED committing transaction: %v\n", err)
				os.Exit(1)
			}
		}

		fmt.Println("OK")
		appliedCount++
	}

	fmt.Printf("\n=== Migration Finished ===\n")
	fmt.Printf("Total Files: %d | Applied: %d | Skipped: %d\n", len(files), appliedCount, skippedCount)
}
