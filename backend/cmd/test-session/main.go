package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"github.com/your-org/used-tools-marketplace/internal/config"
	"github.com/your-org/used-tools-marketplace/internal/database"
	"github.com/your-org/used-tools-marketplace/internal/repository"
	"github.com/your-org/used-tools-marketplace/internal/security"
	"github.com/your-org/used-tools-marketplace/internal/service"
)

func main() {
	if os.Getenv("DB_PASSWORD") == "" && os.Getenv("DATABASE_URL") == "" {
		log.Fatal("[FATAL] DB_PASSWORD environment variable is required (or set DATABASE_URL)")
	}
	if os.Getenv("JWT_SECRET") == "" {
		os.Setenv("JWT_SECRET", "test-session-ephemeral-jwt-secret-value-32chars!")
	}
	if os.Getenv("OTP_ENCRYPT_KEY") == "" {
		os.Setenv("OTP_ENCRYPT_KEY", "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=")
	}
	if os.Getenv("ALLOW_ORIGIN") == "" {
		os.Setenv("ALLOW_ORIGIN", "http://localhost:3000")
	}
	if os.Getenv("FRONTEND_URL") == "" {
		os.Setenv("FRONTEND_URL", "http://localhost:3000")
	}
	if os.Getenv("APP_BASE_URL") == "" {
		os.Setenv("APP_BASE_URL", "http://localhost:8080")
	}
	if os.Getenv("DB_USER") == "" {
		os.Setenv("DB_USER", "postgres")
	}
	if os.Getenv("DB_NAME") == "" {
		os.Setenv("DB_NAME", "hulubedej_db")
	}

	cfg := config.Load()
	db, err := database.NewPool(cfg)
	if err != nil {
		log.Fatalf("db pool failed: %v", err)
	}
	defer db.Close()

	ctx := context.Background()
	adminPasskeyRepo := repository.NewAdminPasskeyRepo(db)
	userRepo := repository.NewUserRepo(db)
	adminRepo := repository.NewAdminRepo(db)

	waManager, err := security.NewWebAuthnManager("localhost", []string{"http://localhost:3000", "http://localhost:8080"})
	if err != nil {
		log.Fatalf("waManager failed: %v", err)
	}

	adminAuthSvc := service.NewAdminAuthService(adminPasskeyRepo, userRepo, adminRepo, waManager, cfg.JWTSecret)

	// Query active admin
	var adminID int
	err = db.QueryRow(ctx, "SELECT id FROM users WHERE role='admin' AND is_active=true LIMIT 1").Scan(&adminID)
	if err != nil {
		log.Fatalf("query admin failed: %v", err)
	}

	u, err := userRepo.GetByID(ctx, adminID)
	if err != nil {
		log.Fatalf("get admin %d failed: %v", adminID, err)
	}

	rawToken, sess, err := adminAuthSvc.CreateAdminSession(ctx, u.ID, "executive", "127.0.0.1", "browser-agent")
	if err != nil {
		log.Fatalf("create admin session failed: %v", err)
	}

	jwtToken, err := security.GenerateAdminJWT(u.ID, u.AdminRole, sess.ID, rawToken, cfg.JWTSecret)
	if err != nil {
		log.Fatalf("generate admin jwt failed: %v", err)
	}

	out := map[string]interface{}{
		"user_id":       u.ID,
		"admin_role":    u.AdminRole,
		"session_id":    sess.ID,
		"token":         jwtToken,
		"session_token": rawToken,
	}

	enc, _ := json.Marshal(out)
	fmt.Println(string(enc))
}
