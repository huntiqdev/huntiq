package service

import (
	"context"
	"errors"
	"strings"

	"github.com/your-org/used-tools-marketplace/internal/models"
	"github.com/your-org/used-tools-marketplace/internal/repository"
)

type OrganizationService struct {
	repo *repository.OrganizationRepo
}

func NewOrganizationService(repo *repository.OrganizationRepo) *OrganizationService {
	return &OrganizationService{repo: repo}
}

// RegisterOrganization registers a new company or government account with structured verification.
func (s *OrganizationService) RegisterOrganization(ctx context.Context, userID int, org *models.Organization) error {
	if userID <= 0 {
		return errors.New("invalid user ID")
	}
	cleanName := strings.TrimSpace(org.Name)
	if cleanName == "" {
		return errors.New("organization legal name is required")
	}
	org.Name = cleanName

	// Validate entity classification
	switch org.EntityType {
	case models.EntityCompany, models.EntityGovernment, models.EntityOther:
		// Valid
	default:
		return errors.New("invalid entity_type: must be COMPANY, GOVERNMENT, or OTHER_LEGALLY_APPROVED_ENTITY")
	}

	// Government account requirements
	if org.EntityType == models.EntityGovernment {
		if org.PublishingAuthorityEvidence == nil || strings.TrimSpace(*org.PublishingAuthorityEvidence) == "" {
			return errors.New("evidence of government publishing authority is mandatory for government accounts")
		}
		if org.OfficialEmail == nil || strings.TrimSpace(*org.OfficialEmail) == "" {
			return errors.New("official institutional email is mandatory for government organization verification")
		}
	}

	// Company account requirements
	if org.EntityType == models.EntityCompany {
		if org.RegistrationNumber == nil || strings.TrimSpace(*org.RegistrationNumber) == "" {
			return errors.New("business registration number / commercial license number is required")
		}
		if org.TINNumber == nil || strings.TrimSpace(*org.TINNumber) == "" {
			return errors.New("tax identification number (TIN) is required for company verification")
		}
	}

	return s.repo.CreateOrganization(ctx, org, userID)
}

// ReviewOrganization handles administrative approval or rejection with audit recording.
func (s *OrganizationService) ReviewOrganization(ctx context.Context, orgID, reviewerID int, approved bool, notes string) error {
	if orgID <= 0 || reviewerID <= 0 {
		return errors.New("invalid organization ID or reviewer ID")
	}
	if strings.TrimSpace(notes) == "" {
		return errors.New("documented review reason is required")
	}

	status := models.OrgStatusRejected
	if approved {
		status = models.OrgStatusVerified
	}

	return s.repo.ReviewOrganization(ctx, orgID, status, reviewerID, notes)
}

// GetOrganization retrieves an organization by ID.
func (s *OrganizationService) GetOrganization(ctx context.Context, orgID int) (*models.Organization, error) {
	return s.repo.GetOrganizationByID(ctx, orgID)
}

// GetUserOrganization retrieves the organization membership of a user.
func (s *OrganizationService) GetUserOrganization(ctx context.Context, userID int) (*models.Organization, *models.OrganizationMember, error) {
	return s.repo.GetUserMembership(ctx, userID)
}

// ListOrganizations lists organizations with optional status and entity filtering.
func (s *OrganizationService) ListOrganizations(ctx context.Context, status, entityType string, limit, offset int) ([]models.Organization, int, error) {
	return s.repo.ListOrganizations(ctx, status, entityType, limit, offset)
}

