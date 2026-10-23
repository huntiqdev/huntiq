package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/your-org/used-tools-marketplace/internal/config"
	"github.com/your-org/used-tools-marketplace/internal/models"
	"github.com/your-org/used-tools-marketplace/internal/security"
)

func main() {
	action := flag.String("action", "", "Action: bootstrap-owner | generate-recovery-token | list-passkeys | status")
	userID := flag.Int("user-id", 0, "Target Administrator User ID")
	email := flag.String("email", "", "Admin Email for bootstrap")
	phone := flag.String("phone", "", "Admin Phone for bootstrap")
	name := flag.String("name", "Platform Owner", "Admin Display Name")
	flag.Parse()

	// Load .env
	_ = godotenv.Load()
	cfg := config.Load()

	pool, err := pgxpool.New(context.Background(), cfg.DSN())
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer pool.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	switch *action {
	case "status":
		var count int
		_ = pool.QueryRow(ctx, "SELECT COUNT(*) FROM users WHERE role='admin' AND admin_role='platform_owner'").Scan(&count)
		fmt.Printf("Platform Owners count: %d\n", count)
		rows, err := pool.Query(ctx, "SELECT id, name, email, phone, admin_role, is_active FROM users WHERE role='admin'")
		if err != nil {
			log.Fatalf("Query failed: %v", err)
		}
		defer rows.Close()
		fmt.Println("\nRegistered Administrators:")
		for rows.Next() {
			var id int
			var n, e, p, role string
			var active bool
			_ = rows.Scan(&id, &n, &e, &p, &role, &active)
			fmt.Printf(" - ID %d: %s <%s> (%s) Role: %s Active: %v\n", id, n, e, p, role, active)
		}

	case "generate-recovery-token":
		if *userID <= 0 {
			log.Fatalf("Error: --user-id is required")
		}
		var user models.User
		err := pool.QueryRow(ctx, "SELECT id, name, email, admin_role FROM users WHERE id=$1 AND role='admin'", *userID).
			Scan(&user.ID, &user.Name, &user.Email, &user.AdminRole)
		if err != nil {
			log.Fatalf("Admin with ID %d not found: %v", *userID, err)
		}

		b := make([]byte, 24)
		_, _ = rand.Read(b)
		token := hex.EncodeToString(b)
		tokenHash := security.HashToken(token)

		// Insert recovery token into admin_auth_challenges
		_, err = pool.Exec(ctx, `
			INSERT INTO admin_auth_challenges (challenge, user_id, portal, ceremony_type, session_data, expires_at, created_at)
			VALUES ($1, $2, 'executive', 'reauth', '{}'::jsonb, NOW() + INTERVAL '15 minutes', NOW())
		`, tokenHash, user.ID)
		if err != nil {
			log.Fatalf("Failed to issue recovery challenge: %v", err)
		}

		// Insert into audit log
		_, _ = pool.Exec(ctx, `
			INSERT INTO admin_logs (admin_id, action, notes, created_at)
			VALUES ($1, 'EMERGENCY_RECOVERY_TOKEN_GENERATED', 'Recovery token generated via offline CLI utility', NOW())
		`, user.ID)

		fmt.Println("=================================================================")
		fmt.Println("HUNTIQ SECURE RECOVERY TOKEN GENERATED")
		fmt.Println("=================================================================")
		fmt.Printf("User:       %s <%s> (Role: %s)\n", user.Name, user.Email, user.AdminRole)
		fmt.Printf("Token:      %s\n", token)
		fmt.Println("Expires in: 15 minutes")
		fmt.Println("Use this token at: /executive-console/login?recovery=" + token)
		fmt.Println("=================================================================")

	case "bootstrap-owner":
		if *email == "" || *phone == "" {
			log.Fatalf("Error: --email and --phone are required for bootstrap")
		}

		// Check if platform_owner already exists
		var existingCount int
		_ = pool.QueryRow(ctx, "SELECT COUNT(*) FROM users WHERE role='admin' AND admin_role='platform_owner'").Scan(&existingCount)
		if existingCount > 0 {
			log.Fatalf("Refusing bootstrap: a Platform Owner already exists in the database. Use emergency recovery instead.")
		}

		b := make([]byte, 16)
		_, _ = rand.Read(b)
		randomPassword := hex.EncodeToString(b)
		pwHash, _ := security.HashPassword(randomPassword)

		var id int
		err := pool.QueryRow(ctx, `
			INSERT INTO users (name, phone, email, password_hash, role, admin_role, is_active, referral_code, created_at)
			VALUES ($1, $2, $3, $4, 'admin', 'platform_owner', true,
			        'REF' || translate(upper(substring(md5(random()::text), 1, 8)), 'IO01', 'ABCD'), NOW())
			RETURNING id
		`, *name, *phone, *email, pwHash).Scan(&id)
		if err != nil {
			log.Fatalf("Failed to create Platform Owner: %v", err)
		}

		_, _ = pool.Exec(ctx, `
			INSERT INTO admin_logs (admin_id, action, notes, created_at)
			VALUES ($1, 'PLATFORM_OWNER_BOOTSTRAP', 'Initial Platform Owner provisioned via offline CLI', NOW())
		`, id)

		fmt.Println("=================================================================")
		fmt.Println("PLATFORM OWNER SUCCESSFULLY PROVISIONED")
		fmt.Println("=================================================================")
		fmt.Printf("ID:    %d\n", id)
		fmt.Printf("Name:  %s\n", *name)
		fmt.Printf("Email: %s\n", *email)
		fmt.Printf("Phone: %s\n", *phone)
		fmt.Printf("Temporary Password: %s\n", randomPassword)
		fmt.Println("Notice: Please log in at /executive-console/login and enroll your WebAuthn passkey immediately.")
		fmt.Println("=================================================================")

	default:
		fmt.Println("Usage: admin-recovery --action=[bootstrap-owner|generate-recovery-token|status] [options]")
		os.Exit(1)
	}
}
