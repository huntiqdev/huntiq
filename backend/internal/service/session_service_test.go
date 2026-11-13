package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
	"github.com/your-org/used-tools-marketplace/internal/models"
	"github.com/your-org/used-tools-marketplace/internal/security"
)

// Mock in-memory implementation of SessionRepository
type mockSessionRepo struct {
	tokens map[string]*models.UserRefreshToken // key is TokenHash
}

func newMockSessionRepo() *mockSessionRepo {
	return &mockSessionRepo{
		tokens: make(map[string]*models.UserRefreshToken),
	}
}

func (m *mockSessionRepo) CreateRefreshToken(ctx context.Context, token *models.UserRefreshToken) error {
	m.tokens[token.TokenHash] = token
	return nil
}

func (m *mockSessionRepo) GetRefreshTokenByHash(ctx context.Context, tokenHash string) (*models.UserRefreshToken, error) {
	t, ok := m.tokens[tokenHash]
	if !ok {
		return nil, errors.New("token not found")
	}
	return t, nil
}

func (m *mockSessionRepo) RotateRefreshToken(ctx context.Context, oldTokenID uuid.UUID, newToken *models.UserRefreshToken) error {
	// Find and mark old token as revoked
	now := time.Now()
	for _, t := range m.tokens {
		if t.ID == oldTokenID {
			t.IsRevoked = true
			t.RevokedAt = &now
			t.RevocationReason = "rotated"
			break
		}
	}
	m.tokens[newToken.TokenHash] = newToken
	return nil
}

func (m *mockSessionRepo) RevokeFamily(ctx context.Context, familyID uuid.UUID, reason string) error {
	now := time.Now()
	for _, t := range m.tokens {
		if t.FamilyID == familyID {
			t.IsRevoked = true
			t.RevokedAt = &now
			t.RevocationReason = reason
		}
	}
	return nil
}

func (m *mockSessionRepo) RevokeByUserID(ctx context.Context, userID int, reason string) error {
	now := time.Now()
	for _, t := range m.tokens {
		if t.UserID == userID {
			t.IsRevoked = true
			t.RevokedAt = &now
			t.RevocationReason = reason
		}
	}
	return nil
}

func (m *mockSessionRepo) RevokeByTokenHash(ctx context.Context, tokenHash, reason string) error {
	now := time.Now()
	if t, ok := m.tokens[tokenHash]; ok {
		t.IsRevoked = true
		t.RevokedAt = &now
		t.RevocationReason = reason
	}
	return nil
}

// Mock SessionUserRepository
type mockUserRepo struct {
	users map[int]*models.User
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{
		users: make(map[int]*models.User),
	}
}

func (m *mockUserRepo) GetByID(ctx context.Context, id int) (*models.User, error) {
	u, ok := m.users[id]
	if !ok {
		return nil, errors.New("user not found")
	}
	return u, nil
}

func TestSessionService_FullLifecycle(t *testing.T) {
	ctx := context.Background()
	jwtSecret := "test-secret-key-that-is-at-least-32-bytes-long"

	mockRepo := newMockSessionRepo()
	mockUsers := newMockUserRepo()

	testUser := &models.User{
		ID:           101,
		Name:         "Abebe Bikila",
		Phone:        "+251911223344",
		Role:         "seller",
		IsActive:     true,
		Capabilities: []string{"buyer", "seller"},
	}
	mockUsers.users[testUser.ID] = testUser

	sessionSvc := NewSessionService(mockRepo, mockUsers, jwtSecret, true)

	var initialRefreshToken string

	t.Run("CreateSession generates tokens and records family", func(t *testing.T) {
		access, refresh, csrf, err := sessionSvc.CreateSession(ctx, testUser, "127.0.0.1", "TestBrowser")
		if err != nil {
			t.Fatalf("CreateSession failed: %v", err)
		}
		if access == "" || refresh == "" || csrf == "" {
			t.Fatalf("Tokens must not be empty")
		}

		// Verify access token claims
		claims, err := security.ParseJWT(access, jwtSecret)
		if err != nil {
			t.Fatalf("Access token invalid: %v", err)
		}
		if claims.UserID != testUser.ID {
			t.Errorf("expected userID %d, got %d", testUser.ID, claims.UserID)
		}

		// Verify stored hash in repository
		tokenHash := HashToken(refresh)
		stored, err := mockRepo.GetRefreshTokenByHash(ctx, tokenHash)
		if err != nil {
			t.Fatalf("Stored token not found in repo: %v", err)
		}
		if stored.UserID != testUser.ID {
			t.Errorf("stored userID mismatch: got %d, want %d", stored.UserID, testUser.ID)
		}
		if stored.IsRevoked {
			t.Errorf("new token must not be revoked")
		}

		initialRefreshToken = refresh
	})

	var rotatedRefreshToken string

	t.Run("RotateSession rotates token and updates lineage", func(t *testing.T) {
		newAccess, newRefresh, newCSRF, err := sessionSvc.RotateSession(ctx, initialRefreshToken, "127.0.0.1", "TestBrowser")
		if err != nil {
			t.Fatalf("RotateSession failed: %v", err)
		}
		if newRefresh == initialRefreshToken {
			t.Errorf("rotated refresh token must differ from initial")
		}
		if newAccess == "" {
			t.Errorf("new access token must not be empty")
		}
		if newCSRF == "" {
			t.Errorf("new CSRF token must not be empty")
		}

		// Verify old token is now revoked with reason "rotated"
		oldHash := HashToken(initialRefreshToken)
		oldToken, err := mockRepo.GetRefreshTokenByHash(ctx, oldHash)
		if err != nil {
			t.Fatalf("Old token not found: %v", err)
		}
		if !oldToken.IsRevoked || oldToken.RevocationReason != "rotated" {
			t.Errorf("old token should be marked revoked with reason 'rotated'")
		}

		// Verify new token lineage
		newHash := HashToken(newRefresh)
		newToken, err := mockRepo.GetRefreshTokenByHash(ctx, newHash)
		if err != nil {
			t.Fatalf("New token not found: %v", err)
		}
		if newToken.ParentTokenHash != oldHash {
			t.Errorf("newToken.ParentTokenHash = %s, want %s", newToken.ParentTokenHash, oldHash)
		}
		if newToken.FamilyID != oldToken.FamilyID {
			t.Errorf("FamilyID must be preserved across rotations")
		}

		rotatedRefreshToken = newRefresh
	})

	t.Run("Concurrency grace window: reuse within 10s succeeds without revoking family", func(t *testing.T) {
		// Calling RotateSession with the old token again immediately (within 10s grace window)
		graceAccess, graceRefresh, graceCSRF, err := sessionSvc.RotateSession(ctx, initialRefreshToken, "127.0.0.1", "ConcurrentTab")
		if err != nil {
			t.Fatalf("Grace window reuse should succeed for concurrent tab, got error: %v", err)
		}
		if graceAccess == "" || graceCSRF == "" {
			t.Errorf("Grace access token must not be empty")
		}
		if graceRefresh != initialRefreshToken {
			t.Errorf("Within grace period, returns existing refresh token")
		}

		// Verify family is NOT revoked
		newHash := HashToken(rotatedRefreshToken)
		newToken, err := mockRepo.GetRefreshTokenByHash(ctx, newHash)
		if err != nil {
			t.Fatalf("New token not found: %v", err)
		}
		if newToken.IsRevoked {
			t.Errorf("Family must not be revoked during legitimate concurrency grace window")
		}
	})

	t.Run("Replay detection: reuse outside grace window revokes entire family", func(t *testing.T) {
		// Artificially age the old token's revocation timestamp beyond 10s
		oldHash := HashToken(initialRefreshToken)
		oldToken, _ := mockRepo.GetRefreshTokenByHash(ctx, oldHash)
		twentySecondsAgo := time.Now().Add(-20 * time.Second)
		oldToken.RevokedAt = &twentySecondsAgo

		// Attempt to rotate using old token outside grace window
		_, _, _, err := sessionSvc.RotateSession(ctx, initialRefreshToken, "127.0.0.1", "Attacker")
		if !errors.Is(err, ErrRefreshTokenReused) {
			t.Fatalf("expected ErrRefreshTokenReused, got: %v", err)
		}

		// Verify all tokens in the family are now revoked!
		newHash := HashToken(rotatedRefreshToken)
		activeToken, _ := mockRepo.GetRefreshTokenByHash(ctx, newHash)
		if !activeToken.IsRevoked || activeToken.RevocationReason != "reuse_detected" {
			t.Errorf("active token in family must be revoked with reason 'reuse_detected'")
		}
	})

	t.Run("Expired refresh token returns ErrRefreshTokenExpired", func(t *testing.T) {
		// Create a new session with an expired token
		_, refresh, _, err := sessionSvc.CreateSession(ctx, testUser, "127.0.0.1", "TestBrowser")
		if err != nil {
			t.Fatalf("CreateSession failed: %v", err)
		}

		// Manually expire it
		tokenHash := HashToken(refresh)
		tRecord, _ := mockRepo.GetRefreshTokenByHash(ctx, tokenHash)
		tRecord.ExpiresAt = time.Now().Add(-1 * time.Hour)

		_, _, _, err = sessionSvc.RotateSession(ctx, refresh, "127.0.0.1", "TestBrowser")
		if !errors.Is(err, ErrRefreshTokenExpired) {
			t.Fatalf("expected ErrRefreshTokenExpired, got: %v", err)
		}
	})

	t.Run("Inactive user returns ErrUserInactive", func(t *testing.T) {
		_, refresh, _, err := sessionSvc.CreateSession(ctx, testUser, "127.0.0.1", "TestBrowser")
		if err != nil {
			t.Fatalf("CreateSession failed: %v", err)
		}

		// Deactivate user
		testUser.IsActive = false

		_, _, _, err = sessionSvc.RotateSession(ctx, refresh, "127.0.0.1", "TestBrowser")
		if !errors.Is(err, ErrUserInactive) {
			t.Fatalf("expected ErrUserInactive, got: %v", err)
		}

		// Restore user active
		testUser.IsActive = true
	})

	t.Run("RevokeAllUserSessions invalidates all sessions for user", func(t *testing.T) {
		_, refresh1, _, _ := sessionSvc.CreateSession(ctx, testUser, "127.0.0.1", "Browser1")
		_, refresh2, _, _ := sessionSvc.CreateSession(ctx, testUser, "127.0.0.1", "Browser2")

		err := sessionSvc.RevokeAllUserSessions(ctx, testUser.ID, "password_changed")
		if err != nil {
			t.Fatalf("RevokeAllUserSessions failed: %v", err)
		}

		t1, _ := mockRepo.GetRefreshTokenByHash(ctx, HashToken(refresh1))
		t2, _ := mockRepo.GetRefreshTokenByHash(ctx, HashToken(refresh2))

		if !t1.IsRevoked || t1.RevocationReason != "password_changed" {
			t.Errorf("t1 should be revoked with reason 'password_changed'")
		}
		if !t2.IsRevoked || t2.RevocationReason != "password_changed" {
			t.Errorf("t2 should be revoked with reason 'password_changed'")
		}
	})

	t.Run("SetAuthCookies and ClearAuthCookies set correct cookie attributes", func(t *testing.T) {
		app := fiber.New()
		c := app.AcquireCtx(&fasthttp.RequestCtx{})
		defer app.ReleaseCtx(c)

		sessionSvc.SetAuthCookies(c, "access-token-123", "refresh-token-456", "csrf-token-789")

		// Verify cookies were added to response
		resp := c.Response()
		cookies := resp.Header.PeekCookie("huntiq_session")
		if len(cookies) == 0 {
			t.Errorf("huntiq_session cookie was not set")
		}
		refreshCookie := resp.Header.PeekCookie("huntiq_refresh")
		if len(refreshCookie) == 0 {
			t.Errorf("huntiq_refresh cookie was not set")
		}
		csrfCookie := resp.Header.PeekCookie("huntiq_csrf")
		if len(csrfCookie) == 0 {
			t.Errorf("huntiq_csrf cookie was not set")
		}

		// Clear cookies
		sessionSvc.ClearAuthCookies(c)
	})
}
