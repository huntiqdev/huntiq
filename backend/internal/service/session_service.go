package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/your-org/used-tools-marketplace/internal/models"
	"github.com/your-org/used-tools-marketplace/internal/security"
)

var (
	ErrInvalidRefreshToken = errors.New("invalid or unauthenticated refresh token")
	ErrRefreshTokenExpired = errors.New("refresh token has expired")
	ErrRefreshTokenReused  = errors.New("security alert: refresh token reuse detected; all sessions revoked")
	ErrUserInactive        = errors.New("user account is inactive or suspended")
)

type SessionRepository interface {
	CreateRefreshToken(ctx context.Context, token *models.UserRefreshToken) error
	GetRefreshTokenByHash(ctx context.Context, tokenHash string) (*models.UserRefreshToken, error)
	RotateRefreshToken(ctx context.Context, oldTokenID uuid.UUID, newToken *models.UserRefreshToken) error
	RevokeFamily(ctx context.Context, familyID uuid.UUID, reason string) error
	RevokeByUserID(ctx context.Context, userID int, reason string) error
	RevokeByTokenHash(ctx context.Context, tokenHash, reason string) error
}

type SessionUserRepository interface {
	GetByID(ctx context.Context, id int) (*models.User, error)
}

type SessionService struct {
	repo        SessionRepository
	userRepo    SessionUserRepository
	jwtSecret   string
	devMode     bool
	accessTTL   time.Duration
	refreshTTL  time.Duration
	graceWindow time.Duration
}

func NewSessionService(
	repo SessionRepository,
	userRepo SessionUserRepository,
	jwtSecret string,
	devMode bool,
) *SessionService {
	return &SessionService{
		repo:        repo,
		userRepo:    userRepo,
		jwtSecret:   jwtSecret,
		devMode:     devMode,
		accessTTL:   15 * time.Minute,
		refreshTTL:  7 * 24 * time.Hour,
		graceWindow: 10 * time.Second,
	}
}

// JWTSecret returns the JWT signing secret
func (s *SessionService) JWTSecret() string {
	return s.jwtSecret
}

// CreateSessionForUserID retrieves the user and establishes a new session
func (s *SessionService) CreateSessionForUserID(
	ctx context.Context,
	userID int,
	ipAddress, userAgent string,
) (accessToken string, refreshToken string, csrfToken string, err error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return "", "", "", err
	}
	return s.CreateSession(ctx, user, ipAddress, userAgent)
}

// GenerateRandomSecret creates a cryptographically secure hex-encoded random string
func GenerateRandomSecret(byteLen int) (string, error) {
	b := make([]byte, byteLen)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// GenerateCSRFToken creates a cryptographically secure URL-safe Anti-CSRF token
func GenerateCSRFToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// HashToken computes the SHA-256 hash of a raw token string
func HashToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

// CreateSession generates a new access token, refresh token family, and CSRF token for a user
func (s *SessionService) CreateSession(
	ctx context.Context,
	user *models.User,
	ipAddress, userAgent string,
) (accessToken string, refreshToken string, csrfToken string, err error) {
	// 1. Generate 15-minute access token
	accessToken, err = security.GenerateJWTWithProfile(
		user.ID,
		user.Name,
		user.Phone,
		user.Role,
		user.AdminRole,
		user.BuyerTier,
		user.SellerPreference,
		user.Capabilities,
		s.jwtSecret,
	)
	if err != nil {
		return "", "", "", err
	}

	// 2. Generate cryptographically random refresh token
	rawRefreshToken, err := GenerateRandomSecret(32)
	if err != nil {
		return "", "", "", err
	}

	// 3. Generate Anti-CSRF token
	csrfToken, err = GenerateCSRFToken()
	if err != nil {
		return "", "", "", err
	}

	// 4. Store refresh token family in database
	familyID := uuid.New()
	tokenRecord := &models.UserRefreshToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		FamilyID:  familyID,
		TokenHash: HashToken(rawRefreshToken),
		IsRevoked: false,
		ExpiresAt: time.Now().Add(s.refreshTTL),
		UserAgent: userAgent,
		IPAddress: ipAddress,
	}

	if err := s.repo.CreateRefreshToken(ctx, tokenRecord); err != nil {
		return "", "", "", err
	}

	return accessToken, rawRefreshToken, csrfToken, nil
}

// RotateSession atomically rotates an existing refresh token with reuse detection
func (s *SessionService) RotateSession(
	ctx context.Context,
	rawRefreshToken string,
	ipAddress, userAgent string,
) (newAccessToken string, newRefreshToken string, newCSRFToken string, err error) {
	if rawRefreshToken == "" {
		return "", "", "", ErrInvalidRefreshToken
	}

	tokenHash := HashToken(rawRefreshToken)
	token, err := s.repo.GetRefreshTokenByHash(ctx, tokenHash)
	if err != nil {
		return "", "", "", ErrInvalidRefreshToken
	}

	// 1. Check Replay / Reuse Detection
	if token.IsRevoked {
		// Concurrency grace check: if token was rotated within the last 10 seconds,
		// allow legitimate concurrent requests (e.g. 2 browser tabs refreshing at the exact same moment)
		if token.RevocationReason == "rotated" && token.RevokedAt != nil && time.Since(*token.RevokedAt) < s.graceWindow {
			log.Printf("[AUTH_CONCURRENCY] Concurrent refresh within %s grace window for user %d", s.graceWindow, token.UserID)
			// Fetch user and issue fresh access token for the sibling tab
			u, uErr := s.userRepo.GetByID(ctx, token.UserID)
			if uErr != nil || !u.IsActive {
				return "", "", "", ErrUserInactive
			}
			newAccess, aErr := security.GenerateJWTWithProfile(
				u.ID, u.Name, u.Phone, u.Role, u.AdminRole, u.BuyerTier, u.SellerPreference, u.Capabilities, s.jwtSecret,
			)
			if aErr != nil {
				return "", "", "", aErr
			}
			newCSRF, _ := GenerateCSRFToken()
			// Return same refresh token within grace period
			return newAccess, rawRefreshToken, newCSRF, nil
		}

		// OUTSIDE grace window: Severe security violation — reuse of consumed token!
		log.Printf("[SECURITY_ALERT] Refresh token reuse detected for user %d, family %s! Revoking all sessions in family.", token.UserID, token.FamilyID)
		_ = s.repo.RevokeFamily(ctx, token.FamilyID, "reuse_detected")
		return "", "", "", ErrRefreshTokenReused
	}

	// 2. Check Expiration
	if time.Now().After(token.ExpiresAt) {
		_ = s.repo.RevokeByTokenHash(ctx, tokenHash, "expired")
		return "", "", "", ErrRefreshTokenExpired
	}

	// 3. Verify User Account is Active
	user, err := s.userRepo.GetByID(ctx, token.UserID)
	if err != nil || !user.IsActive {
		_ = s.repo.RevokeByUserID(ctx, token.UserID, "account_inactive")
		return "", "", "", ErrUserInactive
	}

	// 4. Generate new Token Pair and Anti-CSRF token
	newRawRefreshToken, err := GenerateRandomSecret(32)
	if err != nil {
		return "", "", "", err
	}
	newCSRFToken, err = GenerateCSRFToken()
	if err != nil {
		return "", "", "", err
	}
	newAccessToken, err = security.GenerateJWTWithProfile(
		user.ID, user.Name, user.Phone, user.Role, user.AdminRole, user.BuyerTier, user.SellerPreference, user.Capabilities, s.jwtSecret,
	)
	if err != nil {
		return "", "", "", err
	}

	// 5. Atomically rotate the token in the database
	newTokenRecord := &models.UserRefreshToken{
		ID:              uuid.New(),
		UserID:          user.ID,
		FamilyID:        token.FamilyID, // Keep same family lineage
		TokenHash:       HashToken(newRawRefreshToken),
		ParentTokenHash: token.TokenHash,
		IsRevoked:       false,
		ExpiresAt:       time.Now().Add(s.refreshTTL),
		UserAgent:       userAgent,
		IPAddress:       ipAddress,
	}

	if err := s.repo.RotateRefreshToken(ctx, token.ID, newTokenRecord); err != nil {
		return "", "", "", err
	}

	return newAccessToken, newRawRefreshToken, newCSRFToken, nil
}

// RevokeSession revokes a single session by its refresh token (e.g. user logout)
func (s *SessionService) RevokeSession(ctx context.Context, rawRefreshToken string) error {
	if rawRefreshToken == "" {
		return nil
	}
	return s.repo.RevokeByTokenHash(ctx, HashToken(rawRefreshToken), "logout")
}

// RevokeAllUserSessions revokes all active sessions for a user (e.g. global logout, password change, suspension)
func (s *SessionService) RevokeAllUserSessions(ctx context.Context, userID int, reason string) error {
	return s.repo.RevokeByUserID(ctx, userID, reason)
}

// SetAuthCookies writes hardened, secure cookies to the Fiber response
func (s *SessionService) SetAuthCookies(c *fiber.Ctx, accessToken, refreshToken, csrfToken string) {
	secure := !s.devMode

	// 1. Session Access Token Cookie (15 minutes, HttpOnly)
	c.Cookie(&fiber.Cookie{
		Name:     "huntiq_session",
		Value:    accessToken,
		Path:     "/",
		MaxAge:   int(s.accessTTL.Seconds()),
		HTTPOnly: true,
		Secure:   secure,
		SameSite: "Lax",
	})
	// Backward compatibility alias cookie for Next.js SSR
	c.Cookie(&fiber.Cookie{
		Name:     "hm_token",
		Value:    accessToken,
		Path:     "/",
		MaxAge:   int(s.accessTTL.Seconds()),
		HTTPOnly: true,
		Secure:   secure,
		SameSite: "Lax",
	})

	// 2. Refresh Token Cookie (7 days, HttpOnly)
	if refreshToken != "" {
		c.Cookie(&fiber.Cookie{
			Name:     "huntiq_refresh",
			Value:    refreshToken,
			Path:     "/",
			MaxAge:   int(s.refreshTTL.Seconds()),
			HTTPOnly: true,
			Secure:   secure,
			SameSite: "Lax",
		})
	}

	// 3. Anti-CSRF Token Cookie (Client-readable so frontend sets X-CSRF-Token header)
	if csrfToken != "" {
		c.Cookie(&fiber.Cookie{
			Name:     "huntiq_csrf",
			Value:    csrfToken,
			Path:     "/",
			MaxAge:   int(s.refreshTTL.Seconds()),
			HTTPOnly: false, // Intentionally readable by JavaScript for Anti-CSRF header
			Secure:   secure,
			SameSite: "Lax",
		})
	}
}

// ClearAuthCookies expires all authentication and session cookies on logout
func (s *SessionService) ClearAuthCookies(c *fiber.Ctx) {
	secure := !s.devMode

	cookieNames := []string{"huntiq_session", "hm_token", "huntiq_refresh", "huntiq_csrf", "hm_admin_session"}
	for _, name := range cookieNames {
		c.Cookie(&fiber.Cookie{
			Name:     name,
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			HTTPOnly: name != "huntiq_csrf",
			Secure:   secure,
			SameSite: "Lax",
		})
	}
}
