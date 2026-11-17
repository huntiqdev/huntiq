package models

import "time"

const (
	EntityIndividual = "INDIVIDUAL"
	EntityCompany    = "COMPANY"
	EntityGovernment = "GOVERNMENT"
	EntityOther      = "OTHER_LEGALLY_APPROVED_ENTITY"

	OrgStatusDraft                  = "DRAFT"
	OrgStatusSubmitted              = "SUBMITTED"
	OrgStatusIdentityPending        = "IDENTITY_PENDING"
	OrgStatusDocumentReview         = "DOCUMENT_REVIEW"
	OrgStatusVerified               = "VERIFIED"
	OrgStatusRejected               = "REJECTED"
	OrgStatusSuspended              = "SUSPENDED"
	OrgStatusExpired               = "EXPIRED"
	OrgStatusReverificationRequired = "REVERIFICATION_REQUIRED"

	OrgRoleOwner     = "owner"
	OrgRoleAdmin     = "admin"
	OrgRolePublisher = "publisher"
	OrgRoleFinance   = "finance"
	OrgRoleViewer    = "viewer"
)

// Organization represents a corporate or government publishing institution.
type Organization struct {
	ID                         int        `json:"id"`
	Name                       string     `json:"name"`
	EntityType                 string     `json:"entity_type"` // COMPANY, GOVERNMENT, OTHER_LEGALLY_APPROVED_ENTITY
	RegistrationNumber         *string    `json:"registration_number,omitempty"`
	TINNumber                  *string    `json:"tin_number,omitempty"`
	OfficialAddress            *string    `json:"official_address,omitempty"`
	OfficialEmail              *string    `json:"official_email,omitempty"`
	OfficialPhone              *string    `json:"official_phone,omitempty"`
	PublishingAuthorityEvidence *string   `json:"publishing_authority_evidence,omitempty"`
	VerificationStatus         string     `json:"verification_status"`
	ReviewedBy                 *int       `json:"reviewed_by,omitempty"`
	ReviewedAt                 *time.Time `json:"reviewed_at,omitempty"`
	ReviewNotes                *string    `json:"review_notes,omitempty"`
	CreatedAt                  time.Time  `json:"created_at"`
	UpdatedAt                  time.Time  `json:"updated_at"`
}

// OrganizationMember links an authorized user to an organization with a specific permission role.
type OrganizationMember struct {
	ID             int       `json:"id"`
	OrganizationID int       `json:"organization_id"`
	UserID         int       `json:"user_id"`
	OrgRole        string    `json:"org_role"` // owner, admin, publisher, finance, viewer
	IsActive       bool      `json:"is_active"`
	JoinedAt       time.Time `json:"joined_at"`
}

