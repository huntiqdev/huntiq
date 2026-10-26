package repository

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestAdminRepo_CreateAdminNotification_SemanticMapping(t *testing.T) {
	db, cleanup := setupTestDB(t)
	if db == nil {
		return
	}
	defer cleanup()

	repo := NewAdminRepo(db)
	ctx := context.Background()

	testCases := []struct {
		inputRole     string
		inputType     string
		expectedPrior string
	}{
		{"super", "critical", "high"},
		{"super", "error", "high"},
		{"admin", "warning", "medium"},
		{"", "info", "normal"},
		{"super", "low", "low"},
		{"super", "custom_unknown", "normal"},
	}

	for i, tc := range testCases {
		testTitle := fmt.Sprintf("Test Notification %d - %s", i, tc.inputType)
		testMsg := fmt.Sprintf("Testing semantic mapping for type %s at %d", tc.inputType, time.Now().UnixNano())

		err := repo.CreateAdminNotification(ctx, testTitle, testMsg, tc.inputType, tc.inputRole)
		if err != nil {
			t.Fatalf("CreateAdminNotification failed for type %q: %v", tc.inputType, err)
		}

		// Query database directly to verify type, priority, and target_role columns
		var storedType, storedPriority string
		var storedTargetRole *string
		query := `SELECT type, priority, target_role FROM admin_notifications WHERE title = $1 ORDER BY id DESC LIMIT 1`
		err = db.QueryRow(ctx, query, testTitle).Scan(&storedType, &storedPriority, &storedTargetRole)
		if err != nil {
			t.Fatalf("Failed to query created notification: %v", err)
		}

		if storedType != tc.inputType {
			t.Errorf("expected stored type %q, got %q", tc.inputType, storedType)
		}
		if storedPriority != tc.expectedPrior {
			t.Errorf("expected stored priority %q for type %q, got %q", tc.expectedPrior, tc.inputType, storedPriority)
		}
		if tc.inputRole != "" {
			if storedTargetRole == nil || *storedTargetRole != tc.inputRole {
				t.Errorf("expected target_role %q, got %v", tc.inputRole, storedTargetRole)
			}
		} else {
			if storedTargetRole != nil {
				t.Errorf("expected NULL target_role, got %v", *storedTargetRole)
			}
		}

		// Clean up test notification
		_, _ = db.Exec(ctx, `DELETE FROM admin_notifications WHERE title = $1`, testTitle)
	}
}

func TestAdminRepo_AdminNotifications_ReadByAndListing(t *testing.T) {
	db, cleanup := setupTestDB(t)
	if db == nil {
		return
	}
	defer cleanup()

	repo := NewAdminRepo(db)
	ctx := context.Background()

	testTitle := fmt.Sprintf("Regression Notif %d", time.Now().UnixNano())
	testMsg := "Testing read_by tracking"
	err := repo.CreateAdminNotification(ctx, testTitle, testMsg, "info", "super")
	if err != nil {
		t.Fatalf("CreateAdminNotification failed: %v", err)
	}
	defer func() {
		_, _ = db.Exec(ctx, `DELETE FROM admin_notifications WHERE title = $1`, testTitle)
	}()

	var notifID int
	err = db.QueryRow(ctx, `SELECT id FROM admin_notifications WHERE title = $1`, testTitle).Scan(&notifID)
	if err != nil {
		t.Fatalf("failed to get notif ID: %v", err)
	}

	admin1 := 8881
	admin2 := 8882

	// List for admin1
	notifs, err := repo.ListAdminNotifications(ctx, admin1, "super")
	if err != nil {
		t.Fatalf("ListAdminNotifications failed: %v", err)
	}

	var found *map[string]interface{}
	for _, n := range notifs {
		if n["id"] == notifID {
			found = &n
			break
		}
	}
	if found == nil {
		t.Fatalf("created notification not found in ListAdminNotifications")
	}
	if (*found)["is_read"].(bool) != false {
		t.Errorf("expected is_read to be false initially for admin1")
	}

	// Mark read by admin1
	err = repo.MarkNotificationRead(ctx, notifID, admin1)
	if err != nil {
		t.Fatalf("MarkNotificationRead failed: %v", err)
	}

	// List again for admin1
	notifs, err = repo.ListAdminNotifications(ctx, admin1, "super")
	if err != nil {
		t.Fatalf("ListAdminNotifications after mark read failed: %v", err)
	}
	for _, n := range notifs {
		if n["id"] == notifID {
			if n["is_read"].(bool) != true {
				t.Errorf("expected is_read to be true for admin1 after MarkNotificationRead")
			}
			break
		}
	}

	// List for admin2 — should still be unread for admin2
	notifsAdmin2, err := repo.ListAdminNotifications(ctx, admin2, "super")
	if err != nil {
		t.Fatalf("ListAdminNotifications for admin2 failed: %v", err)
	}
	for _, n := range notifsAdmin2 {
		if n["id"] == notifID {
			if n["is_read"].(bool) != false {
				t.Errorf("expected is_read to be false for admin2")
			}
			break
		}
	}
}
