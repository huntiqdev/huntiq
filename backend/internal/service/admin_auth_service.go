package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/your-org/used-tools-marketplace/internal/models"
	"github.com/your-org/used-tools-marketplace/internal/repository"
	"github.com/your-org/used-tools-marketplace/internal/security"
)

// AdminAuthService coordinates administrative authentication, passkeys, sessions, and zero-trust controls
type AdminAuthService struct {
	passkeyRepo *repository.AdminPasskeyRepo
	userRepo    *repository.UserRepo
	adminRepo   *repository.AdminRepo
	waManager   *security.WebAuthnManager
	jwtSecret   string
}

// NewAdminAuthService creates a new instance of AdminAuthService
func NewAdminAuthService(
	passkeyRepo *repository.AdminPasskeyRepo,
	userRepo *repository.UserRepo,
	adminRepo *repository.AdminRepo,
	waManager *security.WebAuthnManager,
	jwtSecret string,
) *AdminAuthService {
	return &AdminAuthService{
		passkeyRepo: passkeyRepo,
		userRepo:    userRepo,
		adminRepo:   adminRepo,
		waManager:   waManager,
		jwtSecret:   jwtSecret,
	}
}

// GenerateRandomToken generates a cryptographically secure hex token
func GenerateRandomToken(bytesLen int) string {
	b := make([]byte, bytesLen)
	if _, err := rand.Read(b); err != nil {
		panic("crypto/rand unavailable")
	}
	return hex.EncodeToString(b)
}

// CheckThrottling checks progressive throttling across (identifier, IP)
func (s *AdminAuthService) CheckThrottling(ctx context.Context, identifier, ip string) (bool, int, error) {
	// Lookback 5 minutes
	failures, err := s.passkeyRepo.CountRecentFailedAttempts(ctx, identifier, ip, 5*time.Minute)
	if err != nil {
		return true, 0, nil // Fail open on query error to avoid self-DoS, but log
	}

	if failures >= 10 {
		return false, 300, nil // 5 minutes backoff
	}
	if failures >= 5 {
		return false, 60, nil // 1 minute backoff
	}
	return true, 0, nil
}

// RecordAuthResult logs an attempt for progressive rate limiting
func (s *AdminAuthService) RecordAuthResult(ctx context.Context, identifier, ip, portal string, success bool, reason string) {
	_ = s.passkeyRepo.RecordAuthAttempt(ctx, &models.AuthAttempt{
		Identifier:    identifier,
		IPAddress:     ip,
		Portal:        portal,
		Success:       success,
		FailureReason: reason,
	})
}

// BeginPasskeyRegistration initiates WebAuthn enrollment
func (s *AdminAuthService) BeginPasskeyRegistration(ctx context.Context, userID int, portal string) (*protocol.CredentialCreation, *models.AuthChallenge, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, nil, errors.New("user not found")
	}
	if user.Role != "admin" {
		return nil, nil, errors.New("unauthorized: administrative account required")
	}

	creds, err := s.passkeyRepo.ListPasskeysByUserID(ctx, userID)
	if err != nil {
		return nil, nil, err
	}

	adapter := &security.WebAuthnUserAdapter{
		User:        user,
		Credentials: creds,
	}

	creation, sessionData, err := s.waManager.BeginRegistration(adapter)
	if err != nil {
		return nil, nil, fmt.Errorf("begin registration: %w", err)
	}

	sessJSON, err := security.SerializeSessionData(sessionData)
	if err != nil {
		return nil, nil, fmt.Errorf("serialize session: %w", err)
	}

	challengeStr := creation.Response.Challenge.String()
	challengeRecord := &models.AuthChallenge{
		Challenge:    challengeStr,
		UserID:       &userID,
		Portal:       portal,
		CeremonyType: "registration",
		SessionData:  sessJSON,
		ExpiresAt:    time.Now().Add(120 * time.Second), // 2-minute registration window
	}

	if err := s.passkeyRepo.SaveChallenge(ctx, challengeRecord); err != nil {
		return nil, nil, fmt.Errorf("save challenge: %w", err)
	}

	return creation, challengeRecord, nil
}

// FinishPasskeyRegistration verifies and persists a newly enrolled passkey
func (s *AdminAuthService) FinishPasskeyRegistration(
	ctx context.Context,
	userID int,
	challengeStr string,
	authenticatorName string,
	parsedCreation *protocol.ParsedCredentialCreationData,
) (*models.PasskeyCredential, error) {
	challengeRecord, err := s.passkeyRepo.GetAndConsumeChallenge(ctx, challengeStr)
	if err != nil {
		return nil, fmt.Errorf("invalid challenge: %w", err)
	}
	if challengeRecord.UserID == nil || *challengeRecord.UserID != userID {
		return nil, errors.New("challenge not bound to authenticated user")
	}
	if challengeRecord.CeremonyType != "registration" {
		return nil, errors.New("invalid ceremony type")
	}

	sessionData, err := security.DeserializeSessionData(challengeRecord.SessionData)
	if err != nil {
		return nil, fmt.Errorf("deserialize session: %w", err)
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	creds, err := s.passkeyRepo.ListPasskeysByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	adapter := &security.WebAuthnUserAdapter{
		User:        user,
		Credentials: creds,
	}

	cred, err := s.waManager.FinishRegistration(adapter, *sessionData, parsedCreation)
	if err != nil {
		return nil, fmt.Errorf("verify registration: %w", err)
	}

	if authenticatorName == "" {
		authenticatorName = "Passkey Authenticator"
	}

	var transports []string
	for _, t := range cred.Transport {
		transports = append(transports, string(t))
	}

	attachment := "platform"
	if cred.Authenticator.Attachment != "" {
		attachment = string(cred.Authenticator.Attachment)
	}

	passkeyRecord := &models.PasskeyCredential{
		UserID:          userID,
		CredentialID:    cred.ID,
		PublicKey:       cred.PublicKey,
		AttestationType: cred.AttestationType,
		AAGUID:          cred.Authenticator.AAGUID,
		SignCount:       cred.Authenticator.SignCount,
		Attachment:      attachment,
		Name:            authenticatorName,
		Transports:      transports,
		BackupEligible:  cred.Flags.BackupEligible,
		BackupState:     cred.Flags.BackupState,
	}

	if err := s.passkeyRepo.SavePasskey(ctx, passkeyRecord); err != nil {
		return nil, fmt.Errorf("persist passkey: %w", err)
	}

	return passkeyRecord, nil
}

// BeginPasskeyLogin initiates passkey assertion for login
func (s *AdminAuthService) BeginPasskeyLogin(ctx context.Context, identifier, portal string) (*protocol.CredentialAssertion, *models.AuthChallenge, error) {
	var user *models.User
	var creds []models.PasskeyCredential
	var err error

	if identifier != "" {
		user, err = s.userRepo.GetAdminByEmail(ctx, identifier)
		if err != nil {
			// Try by phone
			user, err = s.userRepo.GetByPhone(ctx, identifier)
			if err != nil {
				user = nil
			}
		}
	}

	var adapter *security.WebAuthnUserAdapter
	if user != nil && user.Role == "admin" {
		creds, _ = s.passkeyRepo.ListPasskeysByUserID(ctx, user.ID)
		adapter = &security.WebAuthnUserAdapter{
			User:        user,
			Credentials: creds,
		}
	} else {
		// Generic dummy user for discoverable credentials / privacy
		dummyUser := &models.User{
			ID:    0,
			Email: "admin@huntiq.local",
			Name:  "HUNTIQ Administrator",
			Role:  "admin",
		}
		adapter = &security.WebAuthnUserAdapter{
			User:        dummyUser,
			Credentials: []models.PasskeyCredential{},
		}
	}

	assertion, sessionData, err := s.waManager.BeginLogin(adapter)
	if err != nil {
		return nil, nil, fmt.Errorf("begin login: %w", err)
	}

	sessJSON, err := security.SerializeSessionData(sessionData)
	if err != nil {
		return nil, nil, fmt.Errorf("serialize session: %w", err)
	}

	var uid *int
	if user != nil && user.ID > 0 {
		uid = &user.ID
	}

	challengeStr := assertion.Response.Challenge.String()
	challengeRecord := &models.AuthChallenge{
		Challenge:    challengeStr,
		UserID:       uid,
		Portal:       portal,
		CeremonyType: "authentication",
		SessionData:  sessJSON,
		ExpiresAt:    time.Now().Add(90 * time.Second),
	}

	if err := s.passkeyRepo.SaveChallenge(ctx, challengeRecord); err != nil {
		return nil, nil, fmt.Errorf("save challenge: %w", err)
	}

	return assertion, challengeRecord, nil
}

// FinishPasskeyLogin validates the passkey assertion, enforces user verification, and issues a server session
func (s *AdminAuthService) FinishPasskeyLogin(
	ctx context.Context,
	challengeStr string,
	portal string,
	ipAddress string,
	userAgent string,
	parsedAssertion *protocol.ParsedCredentialAssertionData,
) (string, int, *models.User, error) {
	challengeRecord, err := s.passkeyRepo.GetAndConsumeChallenge(ctx, challengeStr)
	if err != nil {
		return "", 0, nil, fmt.Errorf("invalid challenge: %w", err)
	}

	credID := parsedAssertion.RawID
	passkey, err := s.passkeyRepo.GetPasskeyByCredentialID(ctx, credID)
	if err != nil {
		return "", 0, nil, errors.New("unrecognized passkey credential")
	}

	user, err := s.userRepo.GetByID(ctx, passkey.UserID)
	if err != nil {
		return "", 0, nil, errors.New("administrator not found")
	}
	if user.Role != "admin" || !user.IsActive {
		return "", 0, nil, errors.New("administrator account is inactive or disabled")
	}

	// Portal alignment check
	if err := ValidatePortalAccess(user.AdminRole, portal); err != nil {
		return "", 0, nil, err
	}

	sessionData, err := security.DeserializeSessionData(challengeRecord.SessionData)
	if err != nil {
		return "", 0, nil, fmt.Errorf("deserialize session: %w", err)
	}

	userCreds, err := s.passkeyRepo.ListPasskeysByUserID(ctx, user.ID)
	if err != nil {
		return "", 0, nil, err
	}

	adapter := &security.WebAuthnUserAdapter{
		User:        user,
		Credentials: userCreds,
	}

	verifiedCred, err := s.waManager.FinishLogin(adapter, *sessionData, parsedAssertion)
	if err != nil {
		return "", 0, nil, fmt.Errorf("assertion verification failed: %w", err)
	}

	// Signature counter rollback / clone check
	if verifiedCred.Authenticator.SignCount > 0 && verifiedCred.Authenticator.SignCount <= passkey.SignCount {
		_ = s.passkeyRepo.FlagPasskeyCloneWarning(ctx, passkey.CredentialID)
		return "", 0, nil, errors.New("security warning: potential authenticator cloning or replay detected")
	}

	_ = s.passkeyRepo.UpdatePasskeySignCount(ctx, passkey.CredentialID, verifiedCred.Authenticator.SignCount)

	rawSessionToken, session, err := s.CreateAdminSession(ctx, user.ID, portal, ipAddress, userAgent)
	if err != nil {
		return "", 0, nil, fmt.Errorf("create session: %w", err)
	}

	return rawSessionToken, session.ID, user, nil
}

// CreateAdminSession creates a server-side session for an administrator (from passkey or password login)
func (s *AdminAuthService) CreateAdminSession(ctx context.Context, userID int, portal, ipAddress, userAgent string) (string, *models.AdminSession, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil || user == nil || !user.IsActive || user.Role != "admin" {
		return "", nil, errors.New("administrator account inactive or unauthorized")
	}

	idleTimeout := 1800 // 30 minutes for staff
	absLifetime := 12 * time.Hour
	if user.AdminRole == models.RolePlatformOwner || portal == "executive" {
		idleTimeout = 900 // 15 minutes for executive
		absLifetime = 4 * time.Hour
	}

	rawSessionToken := GenerateRandomToken(32)
	tokenHash := security.HashToken(rawSessionToken)

	session := &models.AdminSession{
		SessionTokenHash:   tokenHash,
		UserID:             user.ID,
		Portal:             portal,
		Role:               user.AdminRole,
		IPAddress:          ipAddress,
		UserAgent:          userAgent,
		IdleTimeoutSeconds: idleTimeout,
		AbsoluteExpiresAt:  time.Now().Add(absLifetime),
	}

	if err := s.passkeyRepo.CreateSession(ctx, session); err != nil {
		return "", nil, fmt.Errorf("create session: %w", err)
	}

	return rawSessionToken, session, nil
}

// ValidateSessionByID validates a session by primary key ID and ensures user matches and is active
func (s *AdminAuthService) ValidateSessionByID(ctx context.Context, sessionID, userID int) (*models.AdminSession, *models.User, error) {
	if sessionID <= 0 {
		return nil, nil, errors.New("invalid session id")
	}

	session, err := s.passkeyRepo.GetSessionByID(ctx, sessionID)
	if err != nil {
		return nil, nil, errors.New("invalid or expired session")
	}

	if session.UserID != userID {
		return nil, nil, errors.New("session user mismatch")
	}

	user, err := s.userRepo.GetByID(ctx, session.UserID)
	if err != nil || user == nil || !user.IsActive || user.Role != "admin" {
		_ = s.passkeyRepo.RevokeSession(ctx, session.ID, "account_inactive_or_deleted")
		return nil, nil, errors.New("account inactive or unauthorized")
	}

	_ = s.passkeyRepo.UpdateSessionActivity(ctx, session.ID)

	return session, user, nil
}

// GenerateAdminJWT creates a token for an administrator bound to a specific server session ID
func (s *AdminAuthService) GenerateAdminJWT(userID int, adminRole string, sessionID int, sessionToken string) (string, error) {
	return security.GenerateAdminJWT(userID, adminRole, sessionID, sessionToken, s.jwtSecret)
}


// ValidatePortalAccess checks whether an admin_role can access the given portal
func ValidatePortalAccess(role, portal string) error {
	switch portal {
	case "executive":
		if role != models.RolePlatformOwner && role != models.LegacyRoleSuper {
			return errors.New("unauthorized: executive console requires Platform Owner role")
		}
	case "platform":
		if role != models.RolePlatformOwner && role != models.RolePlatformAdmin && role != models.LegacyRoleSuper {
			return errors.New("unauthorized: platform console requires Platform Administrator role")
		}
	case "security":
		if role != models.RolePlatformOwner && role != models.RoleSecurityAdmin && role != models.LegacyRoleSuper {
			return errors.New("unauthorized: security console requires Security Administrator role")
		}
	case "staff":
		// Staff portal accepts all legitimate administrative roles
		if role == "" {
			return errors.New("unauthorized: administrative credentials required")
		}
	}
	return nil
}

// ValidateSession verifies a raw session token and returns the session and user
func (s *AdminAuthService) ValidateSession(ctx context.Context, sessionToken string) (*models.AdminSession, *models.User, error) {
	if sessionToken == "" {
		return nil, nil, errors.New("empty session token")
	}

	tokenHash := security.HashToken(sessionToken)
	session, err := s.passkeyRepo.GetSessionByTokenHash(ctx, tokenHash)
	if err != nil {
		return nil, nil, errors.New("invalid or expired session")
	}

	user, err := s.userRepo.GetByID(ctx, session.UserID)
	if err != nil || user == nil || !user.IsActive || user.Role != "admin" {
		_ = s.passkeyRepo.RevokeSession(ctx, session.ID, "account_inactive_or_deleted")
		return nil, nil, errors.New("account inactive or unauthorized")
	}

	// Update session activity asynchronously or synchronously
	_ = s.passkeyRepo.UpdateSessionActivity(ctx, session.ID)

	return session, user, nil
}

// RevokeSession terminates a specific session
func (s *AdminAuthService) RevokeSession(ctx context.Context, sessionToken, reason string) error {
	tokenHash := security.HashToken(sessionToken)
	session, err := s.passkeyRepo.GetSessionByTokenHash(ctx, tokenHash)
	if err != nil {
		return nil // Already revoked or absent
	}
	return s.passkeyRepo.RevokeSession(ctx, session.ID, reason)
}

// RevokeAllSessions invalidates all active sessions for a user
func (s *AdminAuthService) RevokeAllSessions(ctx context.Context, userID int, reason string) error {
	return s.passkeyRepo.RevokeAllUserSessions(ctx, userID, reason)
}

// ListPasskeys retrieves active credentials for an administrator
func (s *AdminAuthService) ListPasskeys(ctx context.Context, userID int) ([]models.PasskeyCredential, error) {
	return s.passkeyRepo.ListPasskeysByUserID(ctx, userID)
}

// RevokePasskey removes a passkey, preventing deletion of the last authenticator for Platform Owner
func (s *AdminAuthService) RevokePasskey(ctx context.Context, passkeyID, userID int) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return errors.New("user not found")
	}

	creds, err := s.passkeyRepo.ListPasskeysByUserID(ctx, userID)
	if err != nil {
		return err
	}

	if (user.AdminRole == models.RolePlatformOwner || user.AdminRole == models.LegacyRoleSuper) && len(creds) <= 1 {
		return errors.New("cannot revoke the last remaining passkey for Platform Owner — register a replacement first")
	}

	return s.passkeyRepo.RevokePasskey(ctx, passkeyID, userID, "revoked_by_admin")
}

// IssueReauthToken creates a fresh high-assurance re-authentication token
func (s *AdminAuthService) IssueReauthToken(ctx context.Context, userID int, actionScope string) (string, error) {
	rawToken := GenerateRandomToken(24)
	tokenHash := security.HashToken(rawToken)

	tokenRecord := &models.ReauthToken{
		UserID:      userID,
		TokenHash:   tokenHash,
		ActionScope: actionScope,
		ExpiresAt:   time.Now().Add(10 * time.Minute), // Strict 10-minute validity
	}

	if err := s.passkeyRepo.SaveReauthToken(ctx, tokenRecord); err != nil {
		return "", err
	}
	return rawToken, nil
}

// VerifyReauth validates and consumes a reauth token for a sensitive action
func (s *AdminAuthService) VerifyReauth(ctx context.Context, token, actionScope string, userID int) (bool, error) {
	if token == "" {
		return false, nil
	}
	tokenHash := security.HashToken(token)
	return s.passkeyRepo.ValidateAndConsumeReauthToken(ctx, tokenHash, userID, actionScope)
}

// ValidateAdminPassword checks whether the password matches the administrator account
func (s *AdminAuthService) ValidateAdminPassword(ctx context.Context, userID int, password string) (bool, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil || user == nil || !user.IsActive || user.Role != "admin" {
		return false, errors.New("administrator not found or inactive")
	}
	return security.CheckPassword(user.PasswordHash, password), nil
}

// RevokeSessionByID marks a session revoked if it belongs to the given userID or caller is authorized
func (s *AdminAuthService) RevokeSessionByID(ctx context.Context, sessionID, userID int, reason string) error {
	session, err := s.passkeyRepo.GetSessionByID(ctx, sessionID)
	if err != nil {
		return nil // Already revoked or absent
	}
	if userID > 0 && session.UserID != userID {
		return errors.New("unauthorized to revoke another user's session")
	}
	return s.passkeyRepo.RevokeSession(ctx, sessionID, reason)
}

