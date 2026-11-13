package models

import (
	"encoding/json"
	"time"
)

// ClassificationScheme represents an organizational taxonomy system (e.g. Products, Industry, Construction)
type ClassificationScheme struct {
	ID          int             `json:"id"`
	Code        string          `json:"code"`
	NameEn      string          `json:"name_en"`
	NameAm      string          `json:"name_am,omitempty"`
	NameOm      string          `json:"name_om,omitempty"`
	NameTi      string          `json:"name_ti,omitempty"`
	Description string          `json:"description,omitempty"`
	Icon        string          `json:"icon,omitempty"`
	Status      string          `json:"status"` // ACTIVE, SUSPENDED, ARCHIVED, DRAFT
	Version     int             `json:"version"`
	Revision    int             `json:"revision"`
	IsSystem    bool            `json:"is_system"`
	Rules       json.RawMessage `json:"rules,omitempty"`
	CreatedBy   *int            `json:"created_by,omitempty"`
	UpdatedBy   *int            `json:"updated_by,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

// GetName returns localized display name for scheme
func (s ClassificationScheme) GetName(lang string) string {
	switch lang {
	case "am":
		if s.NameAm != "" {
			return s.NameAm
		}
	case "om", "or":
		if s.NameOm != "" {
			return s.NameOm
		}
	case "ti":
		if s.NameTi != "" {
			return s.NameTi
		}
	}
	if s.NameEn != "" {
		return s.NameEn
	}
	return s.Code
}

// ClassificationConcept represents a universal canonical identity independent of hierarchy placement
type ClassificationConcept struct {
	ID            int             `json:"id"`
	Code          string          `json:"code"`
	CanonicalName string          `json:"canonical_name"`
	NameEn        string          `json:"name_en"`
	NameAm        string          `json:"name_am,omitempty"`
	NameOm        string          `json:"name_om,omitempty"`
	NameTi        string          `json:"name_ti,omitempty"`
	Description   string          `json:"description,omitempty"`
	ConceptType   string          `json:"concept_type"` // PRODUCT, SERVICE, MATERIAL, EQUIPMENT, PROPERTY, EMPLOYMENT, etc.
	Status        string          `json:"status"`       // ACTIVE, DRAFT, REVIEW, ARCHIVED, RETIRED
	BrandID       *int            `json:"brand_id,omitempty"`
	Metadata      json.RawMessage `json:"metadata,omitempty"`
	Version       int             `json:"version"`
	CreatedBy     *int            `json:"created_by,omitempty"`
	UpdatedBy     *int            `json:"updated_by,omitempty"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`

	// Relationships and Placements populated on demand
	Placements    []ClassificationPlacement    `json:"placements,omitempty"`
	Relationships []ClassificationRelationship `json:"relationships,omitempty"`
	Aliases       []ClassificationAlias        `json:"aliases,omitempty"`
}

// GetName returns localized display name for concept
func (c ClassificationConcept) GetName(lang string) string {
	switch lang {
	case "am":
		if c.NameAm != "" {
			return c.NameAm
		}
	case "om", "or":
		if c.NameOm != "" {
			return c.NameOm
		}
	case "ti":
		if c.NameTi != "" {
			return c.NameTi
		}
	}
	if c.NameEn != "" {
		return c.NameEn
	}
	return c.CanonicalName
}

// ClassificationPlacement decouples the concept from its position in a scheme hierarchy
type ClassificationPlacement struct {
	ID                  int                    `json:"id"`
	SchemeID            int                    `json:"scheme_id"`
	ConceptID           int                    `json:"concept_id"`
	ParentPlacementID   *int                   `json:"parent_placement_id,omitempty"`
	Role                string                 `json:"role"` // ROOT, SUPERFAMILY, CATEGORY, SUBCATEGORY, LEAF, OTHER
	DisplayNameOverride string                 `json:"display_name_override,omitempty"`
	Slug                string                 `json:"slug"`
	Path                string                 `json:"path"`
	Depth               int                    `json:"depth"`
	DisplayOrder        int                    `json:"display_order"`
	IsLeaf              bool                   `json:"is_leaf"`
	IsActive            bool                   `json:"is_active"`
	Status              string                 `json:"status"` // DRAFT, REVIEW, APPROVED, PUBLISHED, ARCHIVED
	CreatedBy           *int                   `json:"created_by,omitempty"`
	UpdatedBy           *int                   `json:"updated_by,omitempty"`
	CreatedAt           time.Time              `json:"created_at"`
	UpdatedAt           time.Time              `json:"updated_at"`

	// Populated on demand
	Concept     *ClassificationConcept `json:"concept,omitempty"`
	Breadcrumbs []CategoryBreadcrumb   `json:"breadcrumbs,omitempty"`
	Children    []PlacementTreeNode    `json:"children,omitempty"`
}

// GetDisplayName returns the effective name taking override or concept localized name into account
func (p ClassificationPlacement) GetDisplayName(lang string) string {
	if p.DisplayNameOverride != "" {
		return p.DisplayNameOverride
	}
	if p.Concept != nil {
		return p.Concept.GetName(lang)
	}
	return p.Slug
}

// PlacementTreeNode represents a recursive tree node for API responses
type PlacementTreeNode struct {
	ID           int                 `json:"id"`
	SchemeID     int                 `json:"scheme_id"`
	ConceptID    int                 `json:"concept_id"`
	ParentID     *int                `json:"parent_id,omitempty"`
	Name         string              `json:"name"`
	NameEn       string              `json:"name_en"`
	NameAm       string              `json:"name_am,omitempty"`
	NameOm       string              `json:"name_om,omitempty"`
	NameTi       string              `json:"name_ti,omitempty"`
	Slug         string              `json:"slug"`
	Role         string              `json:"role"`
	Depth        int                 `json:"depth"`
	IsLeaf       bool                `json:"is_leaf"`
	IsActive     bool                `json:"is_active"`
	DisplayOrder int                 `json:"display_order"`
	Icon         string              `json:"icon,omitempty"`
	Count        int                 `json:"count"` // Associated listings count
	Children     []PlacementTreeNode `json:"children,omitempty"`
}

// ClassificationAlias represents an alias, alternative name or search synonym
type ClassificationAlias struct {
	ID          int       `json:"id"`
	ConceptID   *int      `json:"concept_id,omitempty"`
	PlacementID *int      `json:"placement_id,omitempty"`
	Alias       string    `json:"alias"`
	Lang        string    `json:"lang"`
	AliasType   string    `json:"alias_type"` // SYNONYM, ALTERNATIVE_NAME, TRANSLITERATION, etc.
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
}

// ClassificationRelationshipType represents a configured relationship semantic
type ClassificationRelationshipType struct {
	Code            string `json:"code"`
	Name            string `json:"name"`
	Description     string `json:"description,omitempty"`
	IsDirected      bool   `json:"is_directed"`
	InverseTypeCode string `json:"inverse_type_code,omitempty"`
	IsActive        bool   `json:"is_active"`
}

// ClassificationRelationship connects concepts across industries and applications
type ClassificationRelationship struct {
	ID                   int                    `json:"id"`
	SourceConceptID      int                    `json:"source_concept_id"`
	TargetConceptID      int                    `json:"target_concept_id"`
	RelationshipTypeCode string                 `json:"relationship_type_code"`
	ContextSchemeID      *int                   `json:"context_scheme_id,omitempty"`
	Status               string                 `json:"status"`
	Metadata             json.RawMessage        `json:"metadata,omitempty"`
	CreatedBy            *int                   `json:"created_by,omitempty"`
	CreatedAt            time.Time              `json:"created_at"`
	UpdatedAt            time.Time              `json:"updated_at"`

	// Populated on demand
	SourceConcept *ClassificationConcept `json:"source_concept,omitempty"`
	TargetConcept *ClassificationConcept `json:"target_concept,omitempty"`
}

// ClassificationAttributeRule assigns dynamic attributes with inheritance
type ClassificationAttributeRule struct {
	ID                       int             `json:"id"`
	PlacementID              *int            `json:"placement_id,omitempty"`
	ConceptID                *int            `json:"concept_id,omitempty"`
	AttributeID              int             `json:"attribute_id"`
	IsRequired               bool            `json:"is_required"`
	IsInherited              bool            `json:"is_inherited"`
	DisplayOrder             int             `json:"display_order"`
	OverrideValidationRules  json.RawMessage `json:"override_validation_rules,omitempty"`
	OverrideUnit             string          `json:"override_unit,omitempty"`
	CreatedAt                time.Time       `json:"created_at"`

	// Populated on demand
	Attribute *AttributeDefinition `json:"attribute,omitempty"`
}

// ClassificationAttributeDTO represents an effective attribute with placement inheritance details
type ClassificationAttributeDTO struct {
	AttributeID     int             `json:"attribute_id"`
	NameEn          string          `json:"name_en"`
	NameAm          string          `json:"name_am,omitempty"`
	NameOm          string          `json:"name_om,omitempty"`
	NameTi          string          `json:"name_ti,omitempty"`
	Slug            string          `json:"slug"`
	DataType        string          `json:"data_type"`
	ValidationRules json.RawMessage `json:"validation_rules,omitempty"`
	Unit            string          `json:"unit,omitempty"`
	IsRequired      bool            `json:"is_required"`
	IsFilterable    bool            `json:"is_filterable"`
	IsSearchable    bool            `json:"is_searchable"`
	IsInherited     bool            `json:"is_inherited"`
	DisplayOrder    int             `json:"display_order"`
}

// ClassificationChangeRequest manages review governance workflow
type ClassificationChangeRequest struct {
	ID              int             `json:"id"`
	EntityType      string          `json:"entity_type"`
	EntityID        *int            `json:"entity_id,omitempty"`
	Action          string          `json:"action"`
	Status          string          `json:"status"`
	ProposedChanges json.RawMessage `json:"proposed_changes"`
	PreviousState   json.RawMessage `json:"previous_state,omitempty"`
	ImpactSummary   json.RawMessage `json:"impact_summary,omitempty"`
	RequestedBy     *int            `json:"requested_by,omitempty"`
	ReviewedBy      *int            `json:"reviewed_by,omitempty"`
	RejectionReason string          `json:"rejection_reason,omitempty"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

// ClassificationAuditLog records immutable classification change events
type ClassificationAuditLog struct {
	ID         int             `json:"id"`
	UserID     *int            `json:"user_id,omitempty"`
	Action     string          `json:"action"`
	EntityType string          `json:"entity_type"`
	EntityID   int             `json:"entity_id"`
	Changes    json.RawMessage `json:"changes"`
	Reason     string          `json:"reason,omitempty"`
	IPAddress  string          `json:"ip_address,omitempty"`
	CreatedAt  time.Time       `json:"created_at"`
}

// ImpactAnalysisResult details downstream effects of a modification
type ImpactAnalysisResult struct {
	EntityID                 int                    `json:"entity_id"`
	EntityType               string                 `json:"entity_type"`
	AffectedPlacementsCount  int                    `json:"affected_placements_count"`
	AffectedListingsCount    int                    `json:"affected_listings_count"`
	AffectedDescendantsCount int                    `json:"affected_descendants_count"`
	AffectedRelationships   int                    `json:"affected_relationships"`
	CanSafelyExecute         bool                   `json:"can_safely_execute"`
	Warnings                 []string               `json:"warnings"`
	Details                  map[string]interface{} `json:"details,omitempty"`
}

// UniversalSearchResult bundles unified search results across concepts, placements and listings
type UniversalSearchResult struct {
	Query         string                  `json:"query"`
	Concepts      []ClassificationConcept `json:"concepts"`
	Placements    []PlacementTreeNode     `json:"placements"`
	Relationships []ClassificationRelationship `json:"relationships,omitempty"`
	TotalResults  int                     `json:"total_results"`
}

// ConceptSearchFilter holds query parameters for concept queries
type ConceptSearchFilter struct {
	SearchQuery string
	ConceptType string
	Status      string
	SchemeCode  string
	Page        int
	Limit       int
	Lang        string
}

