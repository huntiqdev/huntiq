// seed creates or updates the admin user in the database.
// Usage: go run ./cmd/seed
package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

// adminPhone, adminEmail, adminName are read from env vars to avoid hardcoding PII.
// Set SEED_ADMIN_PHONE, SEED_ADMIN_EMAIL, SEED_ADMIN_NAME before running.
func getAdminIdentity() (phone, email, name string) {
	phone = os.Getenv("SEED_ADMIN_PHONE")
	if phone == "" {
		log.Fatal("[FATAL] SEED_ADMIN_PHONE env var is required")
	}
	email = os.Getenv("SEED_ADMIN_EMAIL")
	if email == "" {
		log.Fatal("[FATAL] SEED_ADMIN_EMAIL env var is required")
	}
	name = os.Getenv("SEED_ADMIN_NAME")
	if name == "" {
		name = "Super Admin"
	}
	return
}

func main() {
	_ = godotenv.Load()

	adminPassword := os.Getenv("SEED_ADMIN_PASSWORD")
	if adminPassword == "" {
		log.Fatal("[FATAL] SEED_ADMIN_PASSWORD env var is required — never hardcode credentials")
	}
	if len(adminPassword) < 8 {
		log.Fatal("[FATAL] SEED_ADMIN_PASSWORD must be at least 8 characters")
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			getenv("DB_HOST", "localhost"),
			getenv("DB_PORT", "5432"),
			getenv("DB_USER", ""),
			getenv("DB_PASSWORD", ""),
			getenv("DB_NAME", ""),
			getenv("DB_SSLMODE", "disable"),
		)
	}

	db, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("connect: %v", err)
	}
	defer db.Close()

	hash, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("hash: %v", err)
	}

	referralCode, err := generateReferralCode()
	if err != nil {
		log.Fatalf("referral code: %v", err)
	}

	adminPhone, adminEmail, adminName := getAdminIdentity()

	_, err = db.Exec(context.Background(), //nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
		`INSERT INTO users (name, phone, email, password_hash, role, admin_role, is_active, referral_code)
		VALUES ($1, $2, $3, $4, 'admin', 'super', true, $5)
		ON CONFLICT (phone, role) DO UPDATE
		  SET password_hash = EXCLUDED.password_hash,
		      email         = EXCLUDED.email,
		      role          = 'admin',
		      admin_role    = 'super',
		      is_active     = true,
		      referral_code = COALESCE(users.referral_code, EXCLUDED.referral_code),
		      updated_at    = NOW()
	`, adminName, adminPhone, adminEmail, string(hash), referralCode)
	if err != nil {
		log.Fatalf("insert: %v", err)
	}

	fmt.Println("✅ Admin user created/updated successfully")
	fmt.Printf("   Email:    %s\n", adminEmail)
	fmt.Printf("   Phone:    %s\n", adminPhone)
	fmt.Println("\n⚠️  Login at /admin/login using email + password")
	fmt.Println("⚠️  Change the password after first login!")
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func generateReferralCode() (string, error) {
	const chars = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	b := make([]byte, 8)
	for i := range b {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			return "", fmt.Errorf("referral code generation failed: %w", err)
		}
		b[i] = chars[n.Int64()]
	}
	return "REF" + string(b), nil
}
