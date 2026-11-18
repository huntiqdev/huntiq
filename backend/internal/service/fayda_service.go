package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/your-org/used-tools-marketplace/internal/models"
	"github.com/your-org/used-tools-marketplace/internal/repository"
)

// MaskNationalID produces a privacy-preserving masked credential representation (e.g., ******1234).
func MaskNationalID(idStr string) string {
	clean := strings.TrimSpace(idStr)
	if len(clean) < 4 {
		return "******"
	}
	suffix := clean[len(clean)-4:]
	return "******" + suffix
}

// HashNationalID computes a one-way deterministic SHA-256 hash for duplicate account detection.
func HashNationalID(idStr string) string {
	clean := strings.ToUpper(strings.TrimSpace(idStr))
	h := sha256.Sum256([]byte(clean))
	return hex.EncodeToString(h[:])
}

// FaydaVerificationRequest represents inbound payload for National ID verification.
type FaydaVerificationRequest struct {
	UserID       int    `json:"user_id"`
	FIN          string `json:"fin"` // Fayda Identification Number
	FullName     string `json:"full_legal_name"`
	Phone        string `json:"phone"`
	ConsentGiven bool   `json:"consent_given"`
	AuthCode     string `json:"auth_code,omitempty"`
}

// FaydaVerificationResponse represents authoritative provider verification response.
type FaydaVerificationResponse struct {
	IsValid          bool                   `json:"is_valid"`
	LegalName        string                 `json:"legal_name"`
	IDMasked         string                 `json:"id_masked"`
	PhoneMatched     bool                   `json:"phone_matched"`
	PhoneMatchStatus string                 `json:"phone_match_status"`
	AuthReference    string                 `json:"auth_reference"`
	Metadata         map[string]interface{} `json:"metadata"`
}

// IdentityVerificationProvider defines the contract for official eKYC verification adapters.
type IdentityVerificationProvider interface {
	VerifyIdentity(ctx context.Context, req *FaydaVerificationRequest) (*FaydaVerificationResponse, error)
	HealthCheck(ctx context.Context) error
	ProviderCode() string
}

// MockFaydaDevProvider simulates official Fayda eKYC in development with realistic verification outcomes.
type MockFaydaDevProvider struct{}

func NewMockFaydaDevProvider() *MockFaydaDevProvider {
	return &MockFaydaDevProvider{}
}

func (p *MockFaydaDevProvider) ProviderCode() string {
	return "fayda_dev_mock"
}

func (p *MockFaydaDevProvider) HealthCheck(ctx context.Context) error {
	return nil
}

func (p *MockFaydaDevProvider) VerifyIdentity(ctx context.Context, req *FaydaVerificationRequest) (*FaydaVerificationResponse, error) {
	cleanFIN := strings.TrimSpace(req.FIN)
	if len(cleanFIN) != 12 {
		return nil, errors.New("invalid Fayda ID format: FIN must be exactly 12 digits")
	}
	for _, ch := range cleanFIN {
		if ch < '0' || ch > '9' {
			return nil, errors.New("invalid Fayda ID format: FIN must contain only numbers")
		}
	}
	if strings.TrimSpace(req.FullName) == "" {
		return nil, errors.New("legal full name is required")
	}
	if !req.ConsentGiven {
		return nil, errors.New("explicit consent to identity verification is required")
	}

	authRef := fmt.Sprintf("FAYDA-AUTH-%d-%s", time.Now().Unix(), cleanFIN[len(cleanFIN)-4:])
	masked := MaskNationalID(cleanFIN)

	return &FaydaVerificationResponse{
		IsValid:          true,
		LegalName:        req.FullName,
		IDMasked:         masked,
		PhoneMatched:     false,
		PhoneMatchStatus: "unsupported_by_provider",
		AuthReference:    authRef,
		Metadata: map[string]interface{}{
			"auth_mechanism":  "official_ekyc_simulated",
			"assurance_level": "substantive",
			"verified_at":     time.Now().Format(time.RFC3339),
		},
	}, nil
}

// OfficialFaydaOIDCProvider connects to Ethiopia's National ID Program (NIDP) MOSIP eKYC gateway.
type OfficialFaydaOIDCProvider struct {
	clientID     string
	clientSecret string
	endpointURL  string
}

func NewOfficialFaydaOIDCProvider(clientID, clientSecret, endpointURL string) *OfficialFaydaOIDCProvider {
	return &OfficialFaydaOIDCProvider{
		clientID:     clientID,
		clientSecret: clientSecret,
		endpointURL:  endpointURL,
	}
}

func (p *OfficialFaydaOIDCProvider) ProviderCode() string {
	return "fayda"
}

func (p *OfficialFaydaOIDCProvider) HealthCheck(ctx context.Context) error {
	if p.endpointURL == "" {
		return errors.New("fayda eKYC endpoint URL is not configured")
	}
	return nil
}

func (p *OfficialFaydaOIDCProvider) VerifyIdentity(ctx context.Context, req *FaydaVerificationRequest) (*FaydaVerificationResponse, error) {
	if p.endpointURL == "" || p.clientID == "" {
		// When official production credentials have not yet been provisioned by NIDP,
		// the account must not be falsely approved.
		return nil, errors.New("fayda_unavailable: official NIDP production connection is pending merchant provisioning")
	}

	// Production eKYC handshake with NIDP OAuth2 / OIDC assertion verification
	// (Moses / OpenID Connect UserInfo exchange with cryptographic JWT signature verification)
	return nil, errors.New("fayda_service_pending: awaiting NIDP production network authorization")
}

// IdentityService orchestrates identity verification, fraud prevention, and audit records.
type IdentityService struct {
	repo     *repository.IdentityRepo
	provider IdentityVerificationProvider
	devMode  bool
}

func NewIdentityService(repo *repository.IdentityRepo, devMode bool) *IdentityService {
	var prov IdentityVerificationProvider
	if devMode {
		prov = NewMockFaydaDevProvider()
	} else {
		prov = NewOfficialFaydaOIDCProvider("", "", "")
	}
	return &IdentityService{
		repo:     repo,
		provider: prov,
		devMode:  devMode,
	}
}

// WithProvider sets a custom verification provider adapter.
func (s *IdentityService) WithProvider(p IdentityVerificationProvider) *IdentityService {
	s.provider = p
	return s
}

// ProcessIndividualIdentityVerification processes National ID verification with privacy preservation.
func (s *IdentityService) ProcessIndividualIdentityVerification(ctx context.Context, req *FaydaVerificationRequest) (*models.IdentityVerification, error) {
	if req.UserID <= 0 {
		return nil, errors.New("invalid user ID")
	}
	if strings.TrimSpace(req.FIN) == "" {
		return nil, errors.New("national identity credential (FIN) is required")
	}
	if strings.TrimSpace(req.FullName) == "" {
		return nil, errors.New("legal full name is required")
	}
	if !req.ConsentGiven {
		return nil, errors.New("explicit consent to identity verification and data processing is required")
	}

	finHash := HashNationalID(req.FIN)
	finMasked := MaskNationalID(req.FIN)

	// Call authoritative eKYC provider
	res, err := s.provider.VerifyIdentity(ctx, req)
	if err != nil {
		log.Printf("[IdentityService] Verification provider failure for user %d: %v", req.UserID, err)
		// If external provider is temporarily unavailable, create a pending/escalated record without granting fake verification
		pendingVerification := &models.IdentityVerification{
			UserID:               req.UserID,
			ProviderCode:         s.provider.ProviderCode(),
			CredentialType:       "national_id",
			CredentialIDHash:     finHash,
			IDMasked:             finMasked,
			FullLegalName:        req.FullName,
			PhoneMatched:         false,
			PhoneMatchStatus:     "unsupported_by_provider",
			AuthReference:        fmt.Sprintf("PENDING-%d", time.Now().UnixNano()),
			Status:               models.IdentityStatusPending,
			ConsentRecorded:      true,
			ConsentTimestamp:     time.Now(),
			VerificationMetadata: map[string]interface{}{"provider_error": err.Error()},
		}
		_ = s.repo.CreateIdentityVerification(ctx, pendingVerification)
		return pendingVerification, fmt.Errorf("identity verification pending: %w", err)
	}

	now := time.Now()
	status := models.IdentityStatusPending
	if res.IsValid {
		status = models.IdentityStatusVerified
	}

	verification := &models.IdentityVerification{
		UserID:               req.UserID,
		ProviderCode:         s.provider.ProviderCode(),
		CredentialType:       "national_id",
		CredentialIDHash:     finHash,
		IDMasked:             res.IDMasked,
		FullLegalName:        res.LegalName,
		PhoneMatched:         res.PhoneMatched,
		PhoneMatchStatus:     res.PhoneMatchStatus,
		AuthReference:        res.AuthReference,
		Status:               status,
		ConsentRecorded:      true,
		ConsentTimestamp:     now,
		VerificationMetadata: res.Metadata,
		VerifiedAt:           &now,
	}

	if err := s.repo.CreateIdentityVerification(ctx, verification); err != nil {
		return nil, fmt.Errorf("failed to store identity verification record: %w", err)
	}

	return verification, nil
}

// GetUserIdentityVerification fetches the verification record for a user.
func (s *IdentityService) GetUserIdentityVerification(ctx context.Context, userID int) (*models.IdentityVerification, error) {
	return s.repo.GetIdentityVerificationByUserID(ctx, userID)
}

// ResolveException allows authorized KYC Verification Officers to approve or reject escalated identity verifications.
func (s *IdentityService) ResolveException(ctx context.Context, id int64, status string, reviewerID int, reason string) error {
	if status != models.IdentityStatusVerified && status != models.IdentityStatusRejected {
		return errors.New("invalid decision status: must be 'verified' or 'rejected'")
	}
	if strings.TrimSpace(reason) == "" {
		return errors.New("documented review reason is required")
	}
	return s.repo.UpdateIdentityStatus(ctx, id, status, reviewerID, reason)
}

// ListPendingVerifications lists pending identity records for administrative review.
func (s *IdentityService) ListPendingVerifications(ctx context.Context, limit, offset int) ([]models.IdentityVerification, int, error) {
	return s.repo.ListPendingIdentityVerifications(ctx, limit, offset)
}
