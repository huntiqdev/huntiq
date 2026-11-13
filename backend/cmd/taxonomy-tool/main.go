package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	cmd := os.Args[1]
	switch cmd {
	case "compile-roots":
		runCompileRoots()
	case "populate-genuine":
		runPopulateGenuine()
	case "audit":
		runAuditReport()
	default:
		fmt.Printf("Unknown command: %s\n", cmd)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("HUNTIQ Folder-Based Taxonomy Tool")
	fmt.Println("Source of truth: data/taxonomy/category/ (folder tree with node.json files)")
	fmt.Println("Usage:")
	fmt.Println("  go run ./cmd/taxonomy-tool compile-roots       Validate category/ folder tree (check translations, slugs, duplicates)")
	fmt.Println("  go run ./cmd/taxonomy-tool populate-genuine    Walk category/ folders directly and load taxonomy into PostgreSQL")
	fmt.Println("  go run ./cmd/taxonomy-tool audit               Run database category hierarchy integrity audit")
}

func getDBConnectionString() string {
	dsn := os.Getenv("DATABASE_URL")
	if dsn != "" {
		return dsn
	}
	pass := os.Getenv("DB_PASSWORD")
	if pass == "" {
		pass = "postgres"
	}
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		getenv("DB_HOST", "localhost"),
		getenv("DB_PORT", "5432"),
		getenv("DB_USER", "postgres"),
		pass,
		getenv("DB_NAME", "hulubedej_db"),
		getenv("DB_SSLMODE", "disable"),
	)
}

func getenv(key, defVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defVal
}

// getDBPool creates a connection pool for audit/populate commands.
func getDBPool() (*pgxpool.Pool, context.Context) {
	dsn := getDBConnectionString()
	ctx := context.Background()
	db, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("[FATAL] Database connection failed: %v", err)
	}
	return db, ctx
}
