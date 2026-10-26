package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/your-org/used-tools-marketplace/internal/config"
	"github.com/your-org/used-tools-marketplace/internal/database"
	"github.com/your-org/used-tools-marketplace/internal/middleware"
	"github.com/your-org/used-tools-marketplace/internal/models"
	"github.com/your-org/used-tools-marketplace/internal/repository"
	"github.com/your-org/used-tools-marketplace/internal/security"
	"github.com/your-org/used-tools-marketplace/internal/service"
)

func setupTransactionManagementApp(t *testing.T) (*fiber.App, *AdminHandler, *service.AdminAuthService, *pgxpool.Pool, string, func()) {
	if os.Getenv("DB_PASSWORD") == "" && os.Getenv("DATABASE_URL") == "" {
		t.Skip("Skipping live database test: DB_PASSWORD or DATABASE_URL not set in environment")
		return nil, nil, nil, nil, "", nil
	}

	cfg := config.NewTestConfig()
	db, err := database.NewPool(cfg)
	if err != nil {
		t.Skipf("Skipping live database test: DB connection failed: %v", err)
		return nil, nil, nil, nil, "", nil
	}

	adminRepo := repository.NewAdminRepo(db)
	userRepo := repository.NewUserRepo(db)
	adminPasskeyRepo := repository.NewAdminPasskeyRepo(db)

	waManager, _ := security.NewWebAuthnManager("localhost", []string{"http://localhost:3000", "http://localhost:8080"})
	adminAuthSvc := service.NewAdminAuthService(adminPasskeyRepo, userRepo, adminRepo, waManager, cfg.JWTSecret)
	adminSvc := service.NewAdminService(adminRepo)

	mw := middleware.New(cfg.JWTSecret).
		WithSessionValidator(adminAuthSvc.ValidateSession).
		WithAdminSessionValidator(adminAuthSvc.ValidateSessionByID).
		WithIsActiveChecker(func(userID int) (bool, error) {
			u, err := userRepo.GetByID(context.Background(), userID)
			if err != nil {
				return false, err
			}
			return u.IsActive, nil
		}).
		WithPermissionChecker(func(adminID int, permission string) (bool, error) {
			u, err := userRepo.GetByID(context.Background(), adminID)
			if err != nil {
				return false, err
			}
			// Platform owner, platform admin, and operations admin have transactions permission
			if u.AdminRole == models.RolePlatformOwner || u.AdminRole == models.LegacyRoleSuper ||
				u.AdminRole == models.RolePlatformAdmin || u.AdminRole == models.RoleOperationsAdmin {
				return true, nil
			}
			return false, nil
		})

	adminHandler := NewAdminHandler(adminSvc)

	app := fiber.New()
	admin := app.Group("/api/admin", mw.RequireAuth, middleware.RequireRole("admin"))
	admin.Get("/transactions", mw.RequirePermission("transactions"), adminHandler.ListTransactions)

	cleanup := func() {
		db.Close()
	}

	return app, adminHandler, adminAuthSvc, db, cfg.JWTSecret, cleanup
}

func TestAdminTransactions_Endpoints(t *testing.T) {
	app, _, adminAuthSvc, db, jwtSecret, cleanup := setupTransactionManagementApp(t)
	if app == nil {
		return
	}
	defer cleanup()

	ctx := context.Background()
	ts := time.Now().UnixNano()

	// 1. Create Test Platform Owner Admin
	ownerPhone := fmt.Sprintf("+25191%07d", ts%10000000)
	var ownerID int
	err := db.QueryRow(ctx, `
		INSERT INTO users (name, phone, role, admin_role, is_active, rating, referral_code)
		VALUES ('Tx Owner Test', $1, 'admin', $2, true, 5.0, $3)
		RETURNING id`, ownerPhone, models.RolePlatformOwner, "REF"+ownerPhone[len(ownerPhone)-6:]).Scan(&ownerID)
	if err != nil {
		t.Fatalf("Failed to insert test owner admin: %v", err)
	}
	defer db.Exec(ctx, "DELETE FROM users WHERE id = $1", ownerID)

	rawOwner, sessOwner, err := adminAuthSvc.CreateAdminSession(ctx, ownerID, "executive", "127.0.0.1", "test-agent")
	if err != nil {
		t.Fatalf("Failed to create owner admin session: %v", err)
	}
	defer db.Exec(ctx, "DELETE FROM admin_sessions WHERE user_id = $1", ownerID)

	ownerToken, _ := security.GenerateAdminJWT(ownerID, models.RolePlatformOwner, sessOwner.ID, rawOwner, jwtSecret)

	// 2. Create Test Operations Admin (Permitted for transactions)
	opsPhone := fmt.Sprintf("+25192%07d", (ts+1)%10000000)
	var opsID int
	err = db.QueryRow(ctx, `
		INSERT INTO users (name, phone, role, admin_role, is_active, rating, referral_code)
		VALUES ('Tx Ops Test', $1, 'admin', $2, true, 5.0, $3)
		RETURNING id`, opsPhone, models.RoleOperationsAdmin, "REF"+opsPhone[len(opsPhone)-6:]).Scan(&opsID)
	if err != nil {
		t.Fatalf("Failed to insert ops admin: %v", err)
	}
	defer db.Exec(ctx, "DELETE FROM users WHERE id = $1", opsID)

	rawOps, sessOps, err := adminAuthSvc.CreateAdminSession(ctx, opsID, "operations", "127.0.0.1", "test-agent")
	if err != nil {
		t.Fatalf("Failed to create ops admin session: %v", err)
	}
	defer db.Exec(ctx, "DELETE FROM admin_sessions WHERE user_id = $1", opsID)

	opsToken, _ := security.GenerateAdminJWT(opsID, models.RoleOperationsAdmin, sessOps.ID, rawOps, jwtSecret)

	// 3. Create Test Moderator Admin (Forbidden for transactions)
	modPhone := fmt.Sprintf("+25193%07d", (ts+2)%10000000)
	var modID int
	err = db.QueryRow(ctx, `
		INSERT INTO users (name, phone, role, admin_role, is_active, rating, referral_code)
		VALUES ('Tx Mod Test', $1, 'admin', $2, true, 5.0, $3)
		RETURNING id`, modPhone, models.RoleModerator, "REF"+modPhone[len(modPhone)-6:]).Scan(&modID)
	if err != nil {
		t.Fatalf("Failed to insert mod admin: %v", err)
	}
	defer db.Exec(ctx, "DELETE FROM users WHERE id = $1", modID)

	rawMod, sessMod, err := adminAuthSvc.CreateAdminSession(ctx, modID, "moderation", "127.0.0.1", "test-agent")
	if err != nil {
		t.Fatalf("Failed to create mod admin session: %v", err)
	}
	defer db.Exec(ctx, "DELETE FROM admin_sessions WHERE user_id = $1", modID)

	modToken, _ := security.GenerateAdminJWT(modID, models.RoleModerator, sessMod.ID, rawMod, jwtSecret)

	// 4. Create Test Seller & Buyer Users for Transaction Records
	sellerPhone := fmt.Sprintf("+25194%07d", (ts+3)%10000000)
	var testSellerID int
	err = db.QueryRow(ctx, `
		INSERT INTO users (name, phone, business_name, role, is_active, rating, referral_code)
		VALUES ('Test Seller Tx', $1, 'Tx Enterprise Store', 'seller', true, 4.9, $2)
		RETURNING id`, sellerPhone, "REF"+sellerPhone[len(sellerPhone)-6:]).Scan(&testSellerID)
	if err != nil {
		t.Fatalf("Failed to insert test seller: %v", err)
	}
	defer db.Exec(ctx, "DELETE FROM users WHERE id = $1", testSellerID)

	// 5. Insert Deterministic Test Transactions
	refA := fmt.Sprintf("TXN-TEST-A-%d", ts%1000000)
	refB := fmt.Sprintf("TXN-TEST-B-%d", ts%1000000)
	refC := fmt.Sprintf("TXN-TEST-C-%d", ts%1000000)

	var txIDA, txIDB, txIDC int
	err = db.QueryRow(ctx, `
		INSERT INTO transactions (
			user_id, seller_id, buyer_name, buyer_phone, amount, fee_amount, broker_fee, seller_gets,
			transaction_type, payment_method, payment_reference, payer_phone, status, notes
		) VALUES (
			$1, $2, 'Abebe Kebede', '+251911000001', 5000.00, 100.00, 100.00, 4900.00,
			'listing_fee', 'telebirr', $3, '+251911000001', 'completed', 'Audit note Telebirr'
		) RETURNING id`, testSellerID, testSellerID, refA).Scan(&txIDA)
	if err != nil {
		t.Fatalf("Failed to insert test transaction A: %v", err)
	}
	defer db.Exec(ctx, "DELETE FROM transactions WHERE id = $1", txIDA)

	err = db.QueryRow(ctx, `
		INSERT INTO transactions (
			user_id, seller_id, buyer_name, buyer_phone, amount, fee_amount, broker_fee, seller_gets,
			transaction_type, payment_method, payment_reference, payer_phone, status, notes
		) VALUES (
			$1, $2, 'Sara Tesfaye', '+251922000002', 12500.00, 250.00, 250.00, 12250.00,
			'product_sale', 'cbe', $3, '+251922000002', 'pending', 'Audit note CBE'
		) RETURNING id`, testSellerID, testSellerID, refB).Scan(&txIDB)
	if err != nil {
		t.Fatalf("Failed to insert test transaction B: %v", err)
	}
	defer db.Exec(ctx, "DELETE FROM transactions WHERE id = $1", txIDB)

	err = db.QueryRow(ctx, `
		INSERT INTO transactions (
			user_id, seller_id, buyer_name, buyer_phone, amount, fee_amount, broker_fee, seller_gets,
			transaction_type, payment_method, payment_reference, payer_phone, status, notes
		) VALUES (
			$1, $2, 'Dawit Alemu', '+251933000003', 3200.00, 64.00, 64.00, 3136.00,
			'product_sale', 'safaricom', $3, '+251933000003', 'cancelled', 'Audit note Safaricom'
		) RETURNING id`, testSellerID, testSellerID, refC).Scan(&txIDC)
	if err != nil {
		t.Fatalf("Failed to insert test transaction C: %v", err)
	}
	defer db.Exec(ctx, "DELETE FROM transactions WHERE id = $1", txIDC)

	// ── Test Suites ─────────────────────────────────────────────────────────────

	t.Run("ListTransactions_Unauthenticated_Returns401", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/admin/transactions", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Request error: %v", err)
		}
		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected 401 Unauthorized, got %d", resp.StatusCode)
		}
	})

	t.Run("ListTransactions_ForbiddenRole_Returns403", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/admin/transactions", nil)
		req.Header.Set("Authorization", "Bearer "+modToken)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Request error: %v", err)
		}
		if resp.StatusCode != http.StatusForbidden {
			t.Errorf("Expected 403 Forbidden for Moderator on transactions endpoint, got %d", resp.StatusCode)
		}
	})

	t.Run("ListTransactions_AuthorizedOperationsAdmin_Returns200WithPagination", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/admin/transactions?page=1&limit=25", nil)
		req.Header.Set("Authorization", "Bearer "+opsToken)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Request error: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected 200 OK for Operations Admin, got %d", resp.StatusCode)
		}

		var body struct {
			Success bool `json:"success"`
			Data    struct {
				Data  []models.Transaction `json:"data"`
				Total int                  `json:"total"`
				Page  int                  `json:"page"`
				Limit int                  `json:"limit"`
			} `json:"data"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}
		if !body.Success {
			t.Errorf("Expected success = true")
		}
		if body.Data.Page != 1 || body.Data.Limit != 25 {
			t.Errorf("Expected page=1, limit=25, got page=%d, limit=%d", body.Data.Page, body.Data.Limit)
		}
		if body.Data.Total < 3 {
			t.Errorf("Expected total >= 3, got %d", body.Data.Total)
		}
	})

	t.Run("ListTransactions_AuthorizedPlatformOwner_Returns200", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/admin/transactions", nil)
		req.Header.Set("Authorization", "Bearer "+ownerToken)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Request error: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected 200 OK for Platform Owner, got %d", resp.StatusCode)
		}
	})

	t.Run("ListTransactions_FilterSearchByReference", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/admin/transactions?search="+refA, nil)
		req.Header.Set("Authorization", "Bearer "+ownerToken)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Request error: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected 200 OK, got %d", resp.StatusCode)
		}

		var body struct {
			Data struct {
				Data  []models.Transaction `json:"data"`
				Total int                  `json:"total"`
			} `json:"data"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
			t.Fatalf("Failed to decode: %v", err)
		}
		if body.Data.Total != 1 || len(body.Data.Data) != 1 {
			t.Fatalf("Expected exactly 1 matched transaction for reference %s, got total=%d, data len=%d", refA, body.Data.Total, len(body.Data.Data))
		}
		tx := body.Data.Data[0]
		if tx.ID != txIDA || tx.PaymentReference != refA {
			t.Errorf("Expected transaction ID %d with ref %s, got ID %d with ref %s", txIDA, refA, tx.ID, tx.PaymentReference)
		}
		if tx.Amount != 5000.00 || tx.BrokerFee != 100.00 || tx.PaymentMethod != "telebirr" {
			t.Errorf("Unexpected transaction fields: amount=%.2f, fee=%.2f, method=%s", tx.Amount, tx.BrokerFee, tx.PaymentMethod)
		}
	})

	t.Run("ListTransactions_FilterSearchByID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/admin/transactions?search=%d", txIDB), nil)
		req.Header.Set("Authorization", "Bearer "+ownerToken)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Request error: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected 200 OK, got %d", resp.StatusCode)
		}

		var body struct {
			Data struct {
				Data  []models.Transaction `json:"data"`
				Total int                  `json:"total"`
			} `json:"data"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
			t.Fatalf("Failed to decode: %v", err)
		}
		if body.Data.Total < 1 {
			t.Fatalf("Expected at least 1 match for search=%d, got %d", txIDB, body.Data.Total)
		}
		found := false
		for _, tx := range body.Data.Data {
			if tx.ID == txIDB {
				found = true
				if tx.BuyerName != "Sara Tesfaye" || tx.PaymentMethod != "cbe" {
					t.Errorf("Unexpected transaction details: %+v", tx)
				}
				break
			}
		}
		if !found {
			t.Errorf("Expected to find txID %d in search results", txIDB)
		}
	})

	t.Run("ListTransactions_FilterByStatus", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/admin/transactions?status=completed&search="+refA, nil)
		req.Header.Set("Authorization", "Bearer "+ownerToken)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Request error: %v", err)
		}
		var body struct {
			Data struct {
				Data  []models.Transaction `json:"data"`
				Total int                  `json:"total"`
			} `json:"data"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
			t.Fatalf("Failed to decode: %v", err)
		}
		if body.Data.Total != 1 || body.Data.Data[0].Status != "completed" {
			t.Errorf("Expected 1 completed transaction, got total=%d", body.Data.Total)
		}

		// Searching for refB (which is 'pending') with status=completed should return 0
		reqMismatched := httptest.NewRequest(http.MethodGet, "/api/admin/transactions?status=completed&search="+refB, nil)
		reqMismatched.Header.Set("Authorization", "Bearer "+ownerToken)
		respMismatched, err := app.Test(reqMismatched)
		if err != nil {
			t.Fatalf("Request error: %v", err)
		}
		var bodyMismatched struct {
			Data struct {
				Data  []models.Transaction `json:"data"`
				Total int                  `json:"total"`
			} `json:"data"`
		}
		if err := json.NewDecoder(respMismatched.Body).Decode(&bodyMismatched); err != nil {
			t.Fatalf("Failed to decode: %v", err)
		}
		if bodyMismatched.Data.Total != 0 || len(bodyMismatched.Data.Data) != 0 {
			t.Errorf("Expected 0 results for pending transaction under completed status filter, got %d", bodyMismatched.Data.Total)
		}
	})

	t.Run("ListTransactions_FilterByPaymentMethod", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/admin/transactions?method=cbe&search="+refB, nil)
		req.Header.Set("Authorization", "Bearer "+ownerToken)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Request error: %v", err)
		}
		var body struct {
			Data struct {
				Data  []models.Transaction `json:"data"`
				Total int                  `json:"total"`
			} `json:"data"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
			t.Fatalf("Failed to decode: %v", err)
		}
		if body.Data.Total != 1 || body.Data.Data[0].PaymentMethod != "cbe" {
			t.Errorf("Expected 1 cbe transaction, got total=%d", body.Data.Total)
		}
	})

	t.Run("ListTransactions_FilterByType", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/admin/transactions?type=listing_fee&search="+refA, nil)
		req.Header.Set("Authorization", "Bearer "+ownerToken)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Request error: %v", err)
		}
		var body struct {
			Data struct {
				Data  []models.Transaction `json:"data"`
				Total int                  `json:"total"`
			} `json:"data"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
			t.Fatalf("Failed to decode: %v", err)
		}
		if body.Data.Total != 1 || body.Data.Data[0].TransactionType != "listing_fee" {
			t.Errorf("Expected 1 listing_fee transaction, got total=%d", body.Data.Total)
		}
	})

	t.Run("ListTransactions_FilterByDateRange", func(t *testing.T) {
		today := time.Now().Format("2006-01-02")
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/admin/transactions?from_date=%s&to_date=%s&search=%s", today, today, refA), nil)
		req.Header.Set("Authorization", "Bearer "+ownerToken)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Request error: %v", err)
		}
		var body struct {
			Data struct {
				Data  []models.Transaction `json:"data"`
				Total int                  `json:"total"`
			} `json:"data"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
			t.Fatalf("Failed to decode: %v", err)
		}
		if body.Data.Total != 1 {
			t.Errorf("Expected 1 transaction created today, got %d", body.Data.Total)
		}
	})

	t.Run("ListTransactions_EmptyResults_ReturnsEmptyArray", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/admin/transactions?search=NON_EXISTENT_REFERENCE_999999", nil)
		req.Header.Set("Authorization", "Bearer "+ownerToken)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Request error: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected 200 OK for empty result query, got %d", resp.StatusCode)
		}
		var body struct {
			Success bool `json:"success"`
			Data    struct {
				Data  []models.Transaction `json:"data"`
				Total int                  `json:"total"`
			} `json:"data"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
			t.Fatalf("Failed to decode: %v", err)
		}
		if body.Data.Total != 0 {
			t.Errorf("Expected total = 0, got %d", body.Data.Total)
		}
		if body.Data.Data == nil || len(body.Data.Data) != 0 {
			t.Errorf("Expected non-nil empty slice `[]`, got %+v", body.Data.Data)
		}
	})

	t.Run("ListTransactions_MalformedParameters_HandledSafely", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/admin/transactions?page=-5&limit=999999&from_date=not-a-date&to_date=invalid", nil)
		req.Header.Set("Authorization", "Bearer "+ownerToken)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Request error: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected 200 OK with sanitized pagination bounds, got %d", resp.StatusCode)
		}
		var body struct {
			Data struct {
				Page  int `json:"page"`
				Limit int `json:"limit"`
			} `json:"data"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
			t.Fatalf("Failed to decode: %v", err)
		}
		// Bound checks: page must be >= 1, limit must be default (50) because 999999 > maxLimit (200)
		if body.Data.Page != 1 {
			t.Errorf("Expected page clamped to 1, got %d", body.Data.Page)
		}
		if body.Data.Limit != 50 {
			t.Errorf("Expected limit clamped to default 50, got %d", body.Data.Limit)
		}
	})

	t.Run("ListTransactions_SensitiveFields_NoCredentialsOrSecretsExposed", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/admin/transactions?search="+refA, nil)
		req.Header.Set("Authorization", "Bearer "+ownerToken)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Request error: %v", err)
		}

		var raw map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
			t.Fatalf("Failed to decode raw response: %v", err)
		}

		rawBytes, _ := json.Marshal(raw)
		rawStr := string(rawBytes)

		// Verify that no internal passwords, secrets, tokens, or private keys leak in payload
		forbiddenKeywords := []string{"password", "password_hash", "token_hash", "otp_hash", "secret", "private_key"}
		for _, kw := range forbiddenKeywords {
			if containsKeyword(rawStr, kw) {
				t.Errorf("Transaction response leaked sensitive keyword %q: %s", kw, rawStr)
			}
		}
	})
}

func containsKeyword(haystack, needle string) bool {
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(haystack), &m); err != nil {
		return false
	}
	return checkMapKeys(m, needle)
}

func checkMapKeys(m map[string]interface{}, needle string) bool {
	for k, v := range m {
		if k == needle {
			return true
		}
		if subMap, ok := v.(map[string]interface{}); ok {
			if checkMapKeys(subMap, needle) {
				return true
			}
		}
		if subList, ok := v.([]interface{}); ok {
			for _, item := range subList {
				if itemMap, ok := item.(map[string]interface{}); ok {
					if checkMapKeys(itemMap, needle) {
						return true
					}
				}
			}
		}
	}
	return false
}
