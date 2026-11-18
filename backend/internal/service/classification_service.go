package service

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/your-org/used-tools-marketplace/internal/models"
	"github.com/your-org/used-tools-marketplace/internal/repository"
)

var (
	slugFormatRegex = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)
)

type ClassificationService struct {
	repo *repository.ClassificationRepo
}

func NewClassificationService(repo *repository.ClassificationRepo) *ClassificationService {
	return &ClassificationService{repo: repo}
}

// ── Schemes ──────────────────────────────────────────────────────────────────

func (s *ClassificationService) CreateScheme(ctx context.Context, scheme *models.ClassificationScheme) (*models.ClassificationScheme, error) {
	scheme.Code = strings.ToUpper(strings.TrimSpace(scheme.Code))
	if scheme.Code == "" {
		return nil, errors.New("scheme code is required")
	}
	if strings.TrimSpace(scheme.NameEn) == "" {
		return nil, errors.New("english name is required")
	}
	if scheme.Status == "" {
		scheme.Status = "ACTIVE"
	}
	return s.repo.CreateScheme(ctx, scheme)
}

func (s *ClassificationService) GetSchemeByCode(ctx context.Context, code string) (*models.ClassificationScheme, error) {
	code = strings.ToUpper(strings.TrimSpace(code))
	scheme, err := s.repo.GetSchemeByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	if scheme == nil {
		return nil, errors.New("scheme not found")
	}
	return scheme, nil
}

func (s *ClassificationService) ListSchemes(ctx context.Context, status string) ([]models.ClassificationScheme, error) {
	return s.repo.ListSchemes(ctx, status)
}

func (s *ClassificationService) UpdateScheme(ctx context.Context, scheme *models.ClassificationScheme) error {
	if scheme.ID <= 0 {
		return errors.New("invalid scheme ID")
	}
	if strings.TrimSpace(scheme.NameEn) == "" {
		return errors.New("english name is required")
	}
	return s.repo.UpdateScheme(ctx, scheme)
}

// ── Concepts ─────────────────────────────────────────────────────────────────

func (s *ClassificationService) CreateConcept(ctx context.Context, concept *models.ClassificationConcept) (*models.ClassificationConcept, error) {
	concept.CanonicalName = strings.TrimSpace(concept.CanonicalName)
	if concept.CanonicalName == "" {
		concept.CanonicalName = concept.NameEn
	}
	if concept.CanonicalName == "" {
		return nil, errors.New("canonical name is required")
	}

	if concept.Code == "" {
		concept.Code = toKebabCase(concept.CanonicalName)
	} else {
		concept.Code = strings.ToLower(strings.TrimSpace(concept.Code))
	}

	if !slugFormatRegex.MatchString(concept.Code) {
		return nil, fmt.Errorf("invalid concept code format %q (must be kebab-case)", concept.Code)
	}

	if concept.NameEn == "" {
		concept.NameEn = concept.CanonicalName
	}
	if concept.ConceptType == "" {
		concept.ConceptType = "PRODUCT"
	}

	// Duplicate candidate check
	candidates, err := s.repo.CheckDuplicateCandidates(ctx, concept.CanonicalName, concept.Code)
	if err == nil && len(candidates) > 0 {
		for _, c := range candidates {
			if strings.EqualFold(c.Code, concept.Code) {
				return nil, fmt.Errorf("concept with code %q already exists (ID: %d)", concept.Code, c.ID)
			}
		}
	}

	return s.repo.CreateConcept(ctx, concept)
}

func (s *ClassificationService) GetConceptByID(ctx context.Context, id int) (*models.ClassificationConcept, error) {
	concept, err := s.repo.GetConceptByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if concept == nil {
		return nil, errors.New("concept not found")
	}

	// Attach relationships
	rels, err := s.repo.GetRelationships(ctx, id, "")
	if err == nil {
		concept.Relationships = rels
	}

	return concept, nil
}

func (s *ClassificationService) SearchConcepts(ctx context.Context, filter models.ConceptSearchFilter) ([]models.ClassificationConcept, int, error) {
	return s.repo.SearchConcepts(ctx, filter)
}

func (s *ClassificationService) UpdateConcept(ctx context.Context, concept *models.ClassificationConcept) error {
	if concept.ID <= 0 {
		return errors.New("invalid concept ID")
	}
	return s.repo.UpdateConcept(ctx, concept)
}

func (s *ClassificationService) CheckDuplicates(ctx context.Context, name, code string) ([]models.ClassificationConcept, error) {
	return s.repo.CheckDuplicateCandidates(ctx, name, code)
}

// ── Placements & Tree ────────────────────────────────────────────────────────

func (s *ClassificationService) CreatePlacement(ctx context.Context, placement *models.ClassificationPlacement) (*models.ClassificationPlacement, error) {
	if placement.SchemeID <= 0 {
		return nil, errors.New("scheme_id is required")
	}
	if placement.ConceptID <= 0 {
		return nil, errors.New("concept_id is required")
	}

	// Fetch concept to get default slug if omitted
	if placement.Slug == "" {
		concept, err := s.repo.GetConceptByID(ctx, placement.ConceptID)
		if err != nil || concept == nil {
			return nil, errors.New("concept not found")
		}
		placement.Slug = concept.Code
	} else {
		placement.Slug = strings.ToLower(strings.TrimSpace(placement.Slug))
	}

	if !slugFormatRegex.MatchString(placement.Slug) {
		return nil, fmt.Errorf("invalid placement slug format %q (must be kebab-case)", placement.Slug)
	}

	if placement.Role == "" {
		if placement.ParentPlacementID == nil {
			placement.Role = "ROOT"
		} else {
			placement.Role = "CATEGORY"
		}
	}

	return s.repo.CreatePlacement(ctx, placement)
}

func (s *ClassificationService) GetPlacementByID(ctx context.Context, id int, lang string) (*models.ClassificationPlacement, error) {
	placement, err := s.repo.GetPlacementByID(ctx, id, lang)
	if err != nil {
		return nil, err
	}
	if placement == nil {
		return nil, errors.New("placement not found")
	}
	return placement, nil
}

func (s *ClassificationService) GetTree(ctx context.Context, schemeCode, lang string, maxDepth int) ([]models.PlacementTreeNode, error) {
	return s.repo.GetTree(ctx, schemeCode, lang, maxDepth)
}

func (s *ClassificationService) GetRoots(ctx context.Context, schemeCode, lang string) ([]models.PlacementTreeNode, error) {
	return s.repo.GetRoots(ctx, schemeCode, lang)
}

func (s *ClassificationService) GetChildren(ctx context.Context, parentPlacementID int, lang string) ([]models.PlacementTreeNode, error) {
	return s.repo.GetChildren(ctx, parentPlacementID, lang)
}

func (s *ClassificationService) MovePlacement(ctx context.Context, placementID, newParentPlacementID int) error {
	if placementID <= 0 {
		return errors.New("invalid placement ID")
	}
	if placementID == newParentPlacementID {
		return errors.New("cannot move placement to become its own parent")
	}

	// Verify cycle: newParent cannot be a descendant of placementID
	target, err := s.repo.GetPlacementByID(ctx, placementID, "en")
	if err != nil || target == nil {
		return errors.New("target placement not found")
	}

	if newParentPlacementID > 0 {
		newParent, err := s.repo.GetPlacementByID(ctx, newParentPlacementID, "en")
		if err != nil || newParent == nil {
			return errors.New("new parent placement not found")
		}
		if newParent.SchemeID != target.SchemeID {
			return errors.New("cannot move placement across different schemes")
		}
		if strings.HasPrefix(newParent.Path, target.Path) {
			return errors.New("cannot move category inside one of its own descendants (hierarchy cycle)")
		}
	}

	return s.repo.MovePlacement(ctx, placementID, newParentPlacementID)
}

func (s *ClassificationService) UpdatePlacement(ctx context.Context, placement *models.ClassificationPlacement) error {
	if placement.ID <= 0 {
		return errors.New("invalid placement ID")
	}
	return s.repo.UpdatePlacement(ctx, placement)
}

func (s *ClassificationService) DeletePlacement(ctx context.Context, placementID int) error {
	if placementID <= 0 {
		return errors.New("invalid placement ID")
	}

	// Run impact check
	impact, err := s.repo.AnalyzeImpact(ctx, "PLACEMENT", placementID)
	if err == nil && !impact.CanSafelyExecute {
		return fmt.Errorf("cannot archive category: %s", strings.Join(impact.Warnings, "; "))
	}

	return s.repo.DeletePlacement(ctx, placementID)
}

// ── Multi-Industry & Multi-Application Relationships ────────────────────────

func (s *ClassificationService) CreateRelationship(ctx context.Context, rel *models.ClassificationRelationship) error {
	if rel.SourceConceptID <= 0 || rel.TargetConceptID <= 0 {
		return errors.New("both source and target concept IDs are required")
	}
	if rel.SourceConceptID == rel.TargetConceptID {
		return errors.New("cannot create self-referencing relationship")
	}
	if rel.RelationshipTypeCode == "" {
		return errors.New("relationship_type_code is required")
	}
	return s.repo.CreateRelationship(ctx, rel)
}

func (s *ClassificationService) GetRelationships(ctx context.Context, conceptID int, relType string) ([]models.ClassificationRelationship, error) {
	return s.repo.GetRelationships(ctx, conceptID, relType)
}

func (s *ClassificationService) ListRelationshipTypes(ctx context.Context) ([]models.ClassificationRelationshipType, error) {
	return s.repo.ListRelationshipTypes(ctx)
}

func (s *ClassificationService) DeleteRelationship(ctx context.Context, id int) error {
	return s.repo.DeleteRelationship(ctx, id)
}

// ── Dynamic Attributes Inheritance ──────────────────────────────────────────

func (s *ClassificationService) GetEffectiveAttributes(ctx context.Context, placementID int) ([]models.ClassificationAttributeDTO, error) {
	return s.repo.GetEffectiveAttributes(ctx, placementID)
}

func (s *ClassificationService) SaveAttributeRule(ctx context.Context, rule *models.ClassificationAttributeRule) error {
	if rule.AttributeID <= 0 {
		return errors.New("attribute_id is required")
	}
	if rule.PlacementID == nil && rule.ConceptID == nil {
		return errors.New("either placement_id or concept_id must be specified")
	}
	return s.repo.SaveAttributeRule(ctx, rule)
}

func (s *ClassificationService) DeleteAttributeRule(ctx context.Context, id int) error {
	return s.repo.DeleteAttributeRule(ctx, id)
}

// ── Impact Analysis & Auditing ──────────────────────────────────────────────

func (s *ClassificationService) AnalyzeImpact(ctx context.Context, entityType string, entityID int) (*models.ImpactAnalysisResult, error) {
	return s.repo.AnalyzeImpact(ctx, entityType, entityID)
}

func (s *ClassificationService) ListAuditLogs(ctx context.Context, entityType string, entityID int, limit, offset int) ([]models.ClassificationAuditLog, error) {
	return s.repo.ListAuditLogs(ctx, entityType, entityID, limit, offset)
}

func (s *ClassificationService) CreateChangeRequest(ctx context.Context, cr *models.ClassificationChangeRequest) error {
	if cr.EntityType == "" || cr.Action == "" {
		return errors.New("entity_type and action are required")
	}
	return s.repo.CreateChangeRequest(ctx, cr)
}

func (s *ClassificationService) ReviewChangeRequest(ctx context.Context, id int, status string, reviewerID int, reason string) error {
	if status != "APPROVED" && status != "REJECTED" && status != "PUBLISHED" {
		return errors.New("invalid review status")
	}
	return s.repo.ReviewChangeRequest(ctx, id, status, reviewerID, reason)
}

// ── Universal Search ────────────────────────────────────────────────────────

func (s *ClassificationService) SearchUniversal(ctx context.Context, query, schemeCode, lang string, limit int) (*models.UniversalSearchResult, error) {
	return s.repo.SearchUniversal(ctx, query, schemeCode, lang, limit)
}

// Helper to convert arbitrary string into kebab-case
func toKebabCase(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	reg := regexp.MustCompile(`[^a-z0-9]+`)
	s = reg.ReplaceAllString(s, "-")
	return strings.Trim(s, "-")
}
