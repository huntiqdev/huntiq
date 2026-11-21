// seed purges non-admins and seeds/synchronizes all 33 enterprise administrator accounts.
// Usage: go run ./cmd/seed
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/your-org/used-tools-marketplace/internal/config"
	"github.com/your-org/used-tools-marketplace/internal/database"
)

func main() {
	_ = godotenv.Load()
	cfg := config.Load()

	pool, err := pgxpool.New(context.Background(), cfg.DSN())
	if err != nil {
		log.Fatalf("connect failed: %v", err)
	}
	defer pool.Close()

	ctx := context.Background()

	// 1. Purge any sellers, buyers, products, orders, transactions
	if err := database.PurgeNonAdminsAndProducts(ctx, pool); err != nil {
		log.Fatalf("purge failed: %v", err)
	}

	// 2. Seed all 33 enterprise admins
	if err := database.SeedAdmins(ctx, pool); err != nil {
		log.Fatalf("seed admins failed: %v", err)
	}

	fmt.Printf("\n🎉 Marketplace successfully initialized with ZERO products, ZERO sellers, ZERO buyers, and %d seeded enterprise administrators!\n", len(database.EnterpriseAdmins))
	os.Exit(0)
}
