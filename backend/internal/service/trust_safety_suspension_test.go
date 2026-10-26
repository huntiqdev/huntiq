package service

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/your-org/used-tools-marketplace/internal/config"
	"github.com/your-org/used-tools-marketplace/internal/database"
	"github.com/your-org/used-tools-marketplace/internal/repository"
)

func setupTestDB(t *testing.T) (*pgxpool.Pool, func()) {
	if os.Getenv("DB_PASSWORD") == "" && os.Getenv("DATABASE_URL") == "" {
		t.Skip("Skipping database-dependent test: DB_PASSWORD or DATABASE_URL not set in environment")
		return nil, nil
	}

	cfg := config.NewTestConfig()
	db, err := database.NewPool(cfg)
	if err != nil {
		t.Skipf("Skipping database-dependent test: DB connection failed: %v", err)
		return nil, nil
	}

	cleanup := func() {
		db.Close()
	}
	return db, cleanup
}

func createTestUser(t *testing.T, db *pgxpool.Pool, phone, role, adminRole string) int {
	var userID int
	err := db.QueryRow(context.Background(),
		`INSERT INTO users (name, phone, role, admin_role, is_active, is_banned, referral_code, created_at, updated_at)
		 VALUES ($1, $2, $3, NULLIF($4,'')::admin_role, TRUE, FALSE, $5, NOW(), NOW())
		 RETURNING id`,
		"Test User "+phone, phone, role, adminRole, "REF"+phone[len(phone)-6:],
	).Scan(&userID)
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	// Add buyer and seller capabilities
	_, _ = db.Exec(context.Background(),
		`INSERT INTO user_capabilities (user_id, capability, status, activated_at)
		 VALUES ($1, 'buyer', 'active', NOW()), ($1, 'seller', 'active', NOW())
		 ON CONFLICT (user_id, capability) DO UPDATE SET status = 'active'`,
		userID)

	return userID
}

func deleteTestUser(db *pgxpool.Pool, userID int) {
	_, _ = db.Exec(context.Background(), `DELETE FROM abuse_reports WHERE reported_user_id=$1 OR reporter_id=$1`, userID)
	_, _ = db.Exec(context.Background(), `DELETE FROM user_capabilities WHERE user_id=$1`, userID)
	_, _ = db.Exec(context.Background(), `DELETE FROM users WHERE id=$1`, userID)
}

func TestReviewAbuseReport_SuspensionEnforcement(t *testing.T) {
	db, cleanup := setupTestDB(t)
	if db == nil {
		return
	}
	defer cleanup()

	ctx := context.Background()
	adminRepo := repository.NewAdminRepo(db)
	adminSvc := NewAdminService(adminRepo)

	ts := time.Now().UnixNano()
	adminPhone := fmt.Sprintf("+25191%07d", ts%10000000)
	targetPhone := fmt.Sprintf("+25192%07d", (ts+1)%10000000)

	adminID := createTestUser(t, db, adminPhone, "admin", "super")
	defer deleteTestUser(db, adminID)

	targetUserID := createTestUser(t, db, targetPhone, "seller", "")
	defer deleteTestUser(db, targetUserID)

	// Create an abuse report targeting targetUserID
	var reportID int
	err := db.QueryRow(ctx,
		`INSERT INTO abuse_reports (reporter_id, reported_user_id, reason, description, status, priority)
		 VALUES ($1, $2, 'spam', 'Repeated spam messages', 'pending', 'high')
		 RETURNING id`,
		adminID, targetUserID,
	).Scan(&reportID)
	if err != nil {
		t.Fatalf("failed to insert abuse report: %v", err)
	}

	// Wire session revoker tracker
	var revokedUsers []struct {
		UserID int
		Reason string
	}
	adminSvc.WithSessionRevoker(func(ctx context.Context, userID int, reason string) error {
		revokedUsers = append(revokedUsers, struct {
			UserID int
			Reason string
		}{UserID: userID, Reason: reason})
		return nil
	})

	// 1. Review abuse report with "suspended" action
	err = adminSvc.ReviewAbuseReport(ctx, reportID, adminID, "resolved", "suspended", "Excessive abuse confirmed")
	if err != nil {
		t.Fatalf("ReviewAbuseReport with suspended failed: %v", err)
	}

	// 2. Verify target user state in DB
	var isActive, isBanned, referralSuspended bool
	var deactivationReason string
	err = db.QueryRow(ctx,
		`SELECT is_active, is_banned, referral_suspended, COALESCE(deactivation_reason, '') 
		 FROM users WHERE id=$1`, targetUserID,
	).Scan(&isActive, &isBanned, &referralSuspended, &deactivationReason)
	if err != nil {
		t.Fatalf("failed to query target user: %v", err)
	}

	if isActive {
		t.Errorf("expected is_active=false for suspended user, got true")
	}
	if isBanned {
		t.Errorf("expected is_banned=false for suspended user (must preserve distinction from permanent ban), got true")
	}
	if !referralSuspended {
		t.Errorf("expected referral_suspended=true for suspended user, got false")
	}
	if !strings.Contains(deactivationReason, "Account suspended via abuse report") || !strings.Contains(deactivationReason, "Excessive abuse confirmed") {
		t.Errorf("expected deactivation reason to contain report note, got: %s", deactivationReason)
	}

	// 3. Verify user capabilities are marked 'suspended'
	var activeCapsCount int
	err = db.QueryRow(ctx,
		`SELECT COUNT(*) FROM user_capabilities WHERE user_id=$1 AND status='active'`,
		targetUserID,
	).Scan(&activeCapsCount)
	if err != nil {
		t.Fatalf("failed to query user capabilities: %v", err)
	}
	if activeCapsCount != 0 {
		t.Errorf("expected 0 active capabilities for suspended user, got %d", activeCapsCount)
	}

	// 4. Verify session revoker was called
	if len(revokedUsers) != 1 {
		t.Fatalf("expected 1 session revoker call, got %d", len(revokedUsers))
	}
	if revokedUsers[0].UserID != targetUserID || revokedUsers[0].Reason != "user_suspended" {
		t.Errorf("unexpected revoker call: %+v", revokedUsers[0])
	}
}

func TestReviewAbuseReport_BanVsSuspendDistinction(t *testing.T) {
	db, cleanup := setupTestDB(t)
	if db == nil {
		return
	}
	defer cleanup()

	ctx := context.Background()
	adminRepo := repository.NewAdminRepo(db)
	adminSvc := NewAdminService(adminRepo)

	ts := time.Now().UnixNano()
	adminID := createTestUser(t, db, fmt.Sprintf("+25191%07d", ts%10000000), "admin", "super")
	defer deleteTestUser(db, adminID)

	bannedUserID := createTestUser(t, db, fmt.Sprintf("+25192%07d", (ts+10)%10000000), "seller", "")
	defer deleteTestUser(db, bannedUserID)

	suspendedUserID := createTestUser(t, db, fmt.Sprintf("+25193%07d", (ts+20)%10000000), "seller", "")
	defer deleteTestUser(db, suspendedUserID)

	// Direct Ban vs Suspend
	if err := adminSvc.BanUser(ctx, bannedUserID, adminID, "Severe violations"); err != nil {
		t.Fatalf("BanUser failed: %v", err)
	}
	if err := adminSvc.SuspendUser(ctx, suspendedUserID, adminID, "Temporary hold"); err != nil {
		t.Fatalf("SuspendUser failed: %v", err)
	}

	// Check Banned User
	var bActive, bBanned bool
	_ = db.QueryRow(ctx, "SELECT is_active, is_banned FROM users WHERE id=$1", bannedUserID).Scan(&bActive, &bBanned)
	if bActive || !bBanned {
		t.Errorf("Banned user must have is_active=false, is_banned=true, got is_active=%v, is_banned=%v", bActive, bBanned)
	}

	// Check Suspended User
	var sActive, sBanned bool
	_ = db.QueryRow(ctx, "SELECT is_active, is_banned FROM users WHERE id=$1", suspendedUserID).Scan(&sActive, &sBanned)
	if sActive || sBanned {
		t.Errorf("Suspended user must have is_active=false, is_banned=false, got is_active=%v, is_banned=%v", sActive, sBanned)
	}
}
