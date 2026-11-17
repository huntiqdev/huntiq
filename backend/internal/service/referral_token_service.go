package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/your-org/used-tools-marketplace/internal/models"
	"github.com/your-org/used-tools-marketplace/internal/repository"
)

type ReferralTokenService struct {
	tokenRepo *repository.ReferralTokenRepo
	ledgerSvc *LedgerService
}

func NewReferralTokenService(tokenRepo *repository.ReferralTokenRepo, ledgerSvc *LedgerService) *ReferralTokenService {
	return &ReferralTokenService{
		tokenRepo: tokenRepo,
		ledgerSvc: ledgerSvc,
	}
}

// hashIdentifier computes a privacy-preserving SHA-256 hash.
func hashIdentifier(val string) string {
	clean := strings.TrimSpace(val)
	if clean == "" {
		return ""
	}
	hash := sha256.Sum256([]byte(clean))
	return hex.EncodeToString(hash[:])
}

// GetUserReferralLink retrieves or creates the user's canonical token and formatted registration URL.
func (s *ReferralTokenService) GetUserReferralLink(ctx context.Context, userID int, referralCode, baseURL string) (string, string, error) {
	if userID <= 0 {
		return "", "", errors.New("invalid user ID")
	}
	if baseURL == "" {
		baseURL = "https://huntiq.et"
	}
	baseURL = strings.TrimRight(baseURL, "/")

	tokenObj, err := s.tokenRepo.GetOrCreateReferralToken(ctx, userID, referralCode)
	if err != nil {
		return "", "", err
	}

	url := fmt.Sprintf("%s/register?ref=%s", baseURL, tokenObj.Token)
	return tokenObj.Token, url, nil
}

// RecordClick captures an inbound referral click event with IP/UA hashes.
func (s *ReferralTokenService) RecordClick(ctx context.Context, tokenStr, ip, userAgent, landingURL string) (*models.ReferralToken, error) {
	cleanToken := strings.TrimSpace(tokenStr)
	if cleanToken == "" {
		return nil, errors.New("token is required")
	}

	tokenObj, err := s.tokenRepo.GetByToken(ctx, cleanToken)
	if err != nil {
		// Also allow fallback lookup by referral code
		tokenObj, err = s.tokenRepo.GetByReferralCode(ctx, cleanToken)
		if err != nil {
			return nil, errors.New("invalid referral token or code")
		}
	}

	if !tokenObj.IsActive {
		return nil, errors.New("referral token is inactive")
	}

	ipHash := hashIdentifier(ip)
	var uaHashPtr *string
	if uaHash := hashIdentifier(userAgent); uaHash != "" {
		uaHashPtr = &uaHash
	}

	if landingURL == "" {
		landingURL = "/register"
	}

	attr := &models.ReferralAttribution{
		TokenID:       tokenObj.ID,
		InviterID:     tokenObj.UserID,
		IPHash:        ipHash,
		UserAgentHash: uaHashPtr,
		LandingURL:    landingURL,
	}

	_, err = s.tokenRepo.RecordAttribution(ctx, attr)
	if err != nil {
		// Log but do not fail click tracking
	}

	_ = s.tokenRepo.IncrementClickCount(ctx, tokenObj.ID)

	return tokenObj, nil
}

// ProcessRegistrationReferral resolves the token/code, permanently locks the attribution to the new user,
// sets users.invited_by, and awards Reward A (verified phone registration reward) via LedgerService.
func (s *ReferralTokenService) ProcessRegistrationReferral(ctx context.Context, tokenOrCode, ip, userAgent string, newUserID int) error {
	clean := strings.TrimSpace(tokenOrCode)
	if clean == "" || newUserID <= 0 {
		return nil
	}

	// 1. Resolve token
	tokenObj, err := s.tokenRepo.GetByToken(ctx, clean)
	if err != nil {
		// Try resolving by referral code
		tokenObj, err = s.tokenRepo.GetByReferralCode(ctx, clean)
		if err != nil {
			return nil // No valid token found, ignore gracefully
		}
	}

	if !tokenObj.IsActive || tokenObj.UserID == newUserID {
		return nil // Avoid self-referral or inactive token
	}

	inviterID := tokenObj.UserID

	// 2. Permanently set invited_by on the new user
	err = s.tokenRepo.SetUserInvitedBy(ctx, newUserID, inviterID)
	if err != nil {
		return fmt.Errorf("failed to link invited_by: %w", err)
	}

	// 3. Find latest click attribution for this token and IP, or record one if not present
	ipHash := hashIdentifier(ip)
	attr, err := s.tokenRepo.FindLatestActiveAttribution(ctx, tokenObj.ID, ipHash)
	if err == nil && attr != nil {
		_ = s.tokenRepo.LockAttributionOnRegistration(ctx, attr.ID, newUserID)
	} else {
		// Create direct attribution
		newAttr := &models.ReferralAttribution{
			TokenID:   tokenObj.ID,
			InviterID: inviterID,
			IPHash:    ipHash,
			LandingURL: "/register",
		}
		attrID, aErr := s.tokenRepo.RecordAttribution(ctx, newAttr)
		if aErr == nil {
			_ = s.tokenRepo.LockAttributionOnRegistration(ctx, attrID, newUserID)
		}
	}

	// 4. Increment conversion counter on the referral token
	_ = s.tokenRepo.IncrementConversionCount(ctx, tokenObj.ID)

	// 5. Credit Reward A (Verified Phone Registration Bonus) to inviter wallet
	if s.ledgerSvc != nil {
		return s.ledgerSvc.CreditRegistrationReward(ctx, inviterID, newUserID)
	}

	return nil
}

// GetUserInviter returns the inviter ID for a user.
func (s *ReferralTokenService) GetUserInviter(ctx context.Context, userID int) (int, error) {
	return s.tokenRepo.GetUserInviter(ctx, userID)
}

