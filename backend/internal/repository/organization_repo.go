package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/your-org/used-tools-marketplace/internal/models"
)

type OrganizationRepo struct {
	db *pgxpool.Pool
}

func NewOrganizationRepo(db *pgxpool.Pool) *OrganizationRepo {
	return &OrganizationRepo{db: db}
}

// CreateOrganization creates a company or government entity and links the registering user as the owner.
func (r *OrganizationRepo) CreateOrganization(ctx context.Context, org *models.Organization, ownerUserID int) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := `
		INSERT INTO organizations (
			name, entity_type, registration_number, tin_number, official_address,
			official_email, official_phone, publishing_authority_evidence,
			verification_status, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW()
		) RETURNING id, created_at, updated_at
	`
	err = tx.QueryRow(ctx, query,
		org.Name, org.EntityType, org.RegistrationNumber, org.TINNumber, org.OfficialAddress,
		org.OfficialEmail, org.OfficialPhone, org.PublishingAuthorityEvidence,
		models.OrgStatusSubmitted,
	).Scan(&org.ID, &org.CreatedAt, &org.UpdatedAt)
	if err != nil {
		return err
	}

	// Add user as organization owner
	_, err = tx.Exec(ctx, `
		INSERT INTO organization_members (organization_id, user_id, org_role, is_active, joined_at)
		VALUES ($1, $2, 'owner', true, NOW())
	`, org.ID, ownerUserID)
	if err != nil {
		return fmt.Errorf("failed to add owner membership: %w", err)
	}

	// Update user's entity type and organization reference
	_, err = tx.Exec(ctx, `
		UPDATE users
		SET organization_id = $1, account_entity_type = $2, business_name = $3, updated_at = NOW()
		WHERE id = $4
	`, org.ID, org.EntityType, org.Name, ownerUserID)
	if err != nil {
		return fmt.Errorf("failed to link user organization: %w", err)
	}

	return tx.Commit(ctx)
}

// GetOrganizationByID retrieves an organization by ID.
func (r *OrganizationRepo) GetOrganizationByID(ctx context.Context, id int) (*models.Organization, error) {
	var o models.Organization
	query := `
		SELECT id, name, entity_type, registration_number, tin_number, official_address,
		       official_email, official_phone, publishing_authority_evidence, verification_status,
		       reviewed_by, reviewed_at, review_notes, created_at, updated_at
		FROM organizations
		WHERE id = $1
	`
	err := r.db.QueryRow(ctx, query, id).Scan(
		&o.ID, &o.Name, &o.EntityType, &o.RegistrationNumber, &o.TINNumber, &o.OfficialAddress,
		&o.OfficialEmail, &o.OfficialPhone, &o.PublishingAuthorityEvidence, &o.VerificationStatus,
		&o.ReviewedBy, &o.ReviewedAt, &o.ReviewNotes, &o.CreatedAt, &o.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &o, nil
}

// GetUserMembership returns the active organization membership for a user.
func (r *OrganizationRepo) GetUserMembership(ctx context.Context, userID int) (*models.Organization, *models.OrganizationMember, error) {
	var m models.OrganizationMember
	var o models.Organization

	query := `
		SELECT m.id, m.organization_id, m.user_id, m.org_role, m.is_active, m.joined_at,
		       o.id, o.name, o.entity_type, o.registration_number, o.tin_number, o.official_address,
		       o.official_email, o.official_phone, o.publishing_authority_evidence, o.verification_status,
		       o.reviewed_by, o.reviewed_at, o.review_notes, o.created_at, o.updated_at
		FROM organization_members m
		JOIN organizations o ON o.id = m.organization_id
		WHERE m.user_id = $1 AND m.is_active = true
		LIMIT 1
	`
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&m.ID, &m.OrganizationID, &m.UserID, &m.OrgRole, &m.IsActive, &m.JoinedAt,
		&o.ID, &o.Name, &o.EntityType, &o.RegistrationNumber, &o.TINNumber, &o.OfficialAddress,
		&o.OfficialEmail, &o.OfficialPhone, &o.PublishingAuthorityEvidence, &o.VerificationStatus,
		&o.ReviewedBy, &o.ReviewedAt, &o.ReviewNotes, &o.CreatedAt, &o.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil, nil
		}
		return nil, nil, err
	}
	return &o, &m, nil
}

// ReviewOrganization records an authorized approval or rejection of an organization.
func (r *OrganizationRepo) ReviewOrganization(ctx context.Context, orgID int, status string, reviewerID int, notes string) error {
	var isMember bool
	err := r.db.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM organization_members WHERE organization_id=$1 AND user_id=$2)`, orgID, reviewerID).Scan(&isMember)
	if err == nil && isMember {
		return errors.New("maker-checker separation: administrator associated with this organization cannot review its verification status")
	}

	_, err = r.db.Exec(ctx, `
		UPDATE organizations
		SET verification_status = $1, reviewed_by = $2, review_notes = $3,
		    reviewed_at = NOW(), updated_at = NOW()
		WHERE id = $4
	`, status, reviewerID, notes, orgID)
	return err
}

// ListOrganizations returns paginated organizations with optional status and entity filtering.
func (r *OrganizationRepo) ListOrganizations(ctx context.Context, status, entityType string, limit, offset int) ([]models.Organization, int, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	whereClause := "WHERE 1=1"
	var args []interface{}
	idx := 1

	if status != "" {
		whereClause += fmt.Sprintf(" AND verification_status = $%d", idx)
		args = append(args, status)
		idx++
	}
	if entityType != "" {
		whereClause += fmt.Sprintf(" AND entity_type = $%d", idx)
		args = append(args, entityType)
		idx++
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM organizations %s", whereClause)
	var total int
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := fmt.Sprintf(`
		SELECT id, name, entity_type, registration_number, tin_number, official_address,
		       official_email, official_phone, publishing_authority_evidence, verification_status,
		       reviewed_by, reviewed_at, review_notes, created_at, updated_at
		FROM organizations
		%s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, idx, idx+1)

	args = append(args, limit, offset)
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []models.Organization
	for rows.Next() {
		var o models.Organization
		if err := rows.Scan(
			&o.ID, &o.Name, &o.EntityType, &o.RegistrationNumber, &o.TINNumber, &o.OfficialAddress,
			&o.OfficialEmail, &o.OfficialPhone, &o.PublishingAuthorityEvidence, &o.VerificationStatus,
			&o.ReviewedBy, &o.ReviewedAt, &o.ReviewNotes, &o.CreatedAt, &o.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		list = append(list, o)
	}

	return list, total, nil
}

