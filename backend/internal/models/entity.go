package models

import (
	"encoding/json"
	"time"
)

// EntityNature represents physical vs non-physical concept classes
const (
	EntityNatureMaterial    = "MATERIAL"
	EntityNatureNonMaterial = "NON_MATERIAL"
)

// Canonical entity classification types
const (
	EntityTypeProduct             = "PRODUCT"
	EntityTypeToolEquipment       = "TOOL_EQUIPMENT"
	EntityTypeVehicle             = "VEHICLE"
	EntityTypeProperty            = "PROPERTY"
	EntityTypeService             = "SERVICE"
	EntityTypeProfession          = "PROFESSION"
	EntityTypeSkill               = "SKILL"
	EntityTypeRawMaterial         = "RAW_MATERIAL"
	EntityTypeBusinessOpportunity = "BUSINESS_OPPORTUNITY"
)

// CanonicalEntity represents an authoritative canonical concept in the HUNTIQ Universal Entity Engine
type CanonicalEntity struct {
	ID                     int                    `json:"id"`
	CanonicalNameEn        string                 `json:"canonical_name_en"`
	CanonicalNameAm        string                 `json:"canonical_name_am,omitempty"`
	CanonicalNameOm        string                 `json:"canonical_name_om,omitempty"`
	CanonicalNameTi        string                 `json:"canonical_name_ti,omitempty"`
	CanonicalSlug          string                 `json:"canonical_slug"`
	EntityNature           string                 `json:"entity_nature"` // MATERIAL or NON_MATERIAL
	EntityType             string                 `json:"entity_type"`   // PRODUCT, TOOL_EQUIPMENT, VEHICLE, etc.
	ShortDefinitionEn      string                 `json:"short_definition_en"`
	ShortDefinitionAm      string                 `json:"short_definition_am,omitempty"`
	ShortDefinitionOm      string                 `json:"short_definition_om,omitempty"`
	ShortDefinitionTi      string                 `json:"short_definition_ti,omitempty"`
	DetailedDefinitionEn   string                 `json:"detailed_definition_en,omitempty"`
	DetailedDefinitionAm   string                 `json:"detailed_definition_am,omitempty"`
	TypicalExamples        json.RawMessage        `json:"typical_examples,omitempty"`
	TypicalUseCases        json.RawMessage        `json:"typical_use_cases,omitempty"`
	FunctionalCapabilities json.RawMessage        `json:"functional_capabilities,omitempty"`
	RelevantAttributes     json.RawMessage        `json:"relevant_attributes,omitempty"`
	DisambiguationRules    json.RawMessage        `json:"disambiguation_rules,omitempty"`
	PrimaryCategoryID      *int                   `json:"primary_category_id,omitempty"`
	ParentEntityID         *int                   `json:"parent_entity_id,omitempty"`
	Status                 string                 `json:"status"`
	CreatedAt              time.Time              `json:"created_at"`
	UpdatedAt              time.Time              `json:"updated_at"`

	// Enriched fields
	PrimaryCategory *CategoryHierarchy    `json:"primary_category,omitempty"`
	ParentEntity    *CanonicalEntity      `json:"parent_entity,omitempty"`
	Aliases         []EntityAlias         `json:"aliases,omitempty"`
	Relationships   []EntityRelationship  `json:"relationships,omitempty"`
	DiscoveryPaths  []EntityDiscoveryPath `json:"discovery_paths,omitempty"`
}

// GetName returns localized canonical name with fallback
func (e CanonicalEntity) GetName(lang string) string {
	switch lang {
	case "am":
		if e.CanonicalNameAm != "" {
			return e.CanonicalNameAm
		}
	case "om", "or":
		if e.CanonicalNameOm != "" {
			return e.CanonicalNameOm
		}
	case "ti":
		if e.CanonicalNameTi != "" {
			return e.CanonicalNameTi
		}
	}
	if e.CanonicalNameEn != "" {
		return e.CanonicalNameEn
	}
	return e.CanonicalSlug
}

// GetShortDefinition returns localized short definition with fallback
func (e CanonicalEntity) GetShortDefinition(lang string) string {
	switch lang {
	case "am":
		if e.ShortDefinitionAm != "" {
			return e.ShortDefinitionAm
		}
	case "om", "or":
		if e.ShortDefinitionOm != "" {
			return e.ShortDefinitionOm
		}
	case "ti":
		if e.ShortDefinitionTi != "" {
			return e.ShortDefinitionTi
		}
	}
	return e.ShortDefinitionEn
}

// EntityAlias represents a known synonym, misspelling, abbreviation, or regional term
type EntityAlias struct {
	ID              int       `json:"id"`
	EntityID        int       `json:"entity_id"`
	AliasText       string    `json:"alias_text"`
	NormalizedAlias string    `json:"normalized_alias"`
	Language        string    `json:"language"`
	AliasType       string    `json:"alias_type"`
	ConfidenceScore float64   `json:"confidence_score"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
}

// EntityRelationship represents a directed edge in the Knowledge Graph
type EntityRelationship struct {
	ID               int                `json:"id"`
	SourceEntityID   int                `json:"source_entity_id"`
	TargetEntityID   *int               `json:"target_entity_id,omitempty"`
	TargetCategoryID *int               `json:"target_category_id,omitempty"`
	RelationshipType string             `json:"relationship_type"`
	ConfidenceScore  float64            `json:"confidence_score"`
	Metadata         json.RawMessage    `json:"metadata,omitempty"`
	Notes            string             `json:"notes,omitempty"`
	Status           string             `json:"status"`
	CreatedAt        time.Time          `json:"created_at"`
	TargetEntity     *CanonicalEntity   `json:"target_entity,omitempty"`
	TargetCategory   *CategoryHierarchy `json:"target_category,omitempty"`
}

// EntityDiscoveryPath represents an automatic navigation or search route
type EntityDiscoveryPath struct {
	PathType         string  `json:"path_type"` // PRIMARY_CATEGORY, RELATED_CATEGORY, RELATED_ENTITY
	TargetID         int     `json:"target_id"`
	TargetName       string  `json:"target_name"`
	RelationshipType string  `json:"relationship_type"`
	Depth            int     `json:"depth"`
	Confidence       float64 `json:"confidence"`
}

// EntityResolutionResult models the output of resolve_canonical_entity()
type EntityResolutionResult struct {
	EntityID            int                   `json:"entity_id"`
	CanonicalNameEn     string                `json:"canonical_name_en"`
	CanonicalNameAm     string                `json:"canonical_name_am,omitempty"`
	CanonicalSlug       string                `json:"canonical_slug"`
	EntityNature        string                `json:"entity_nature"`
	EntityType          string                `json:"entity_type"`
	ShortDefinitionEn   string                `json:"short_definition_en"`
	ShortDefinitionAm   string                `json:"short_definition_am,omitempty"`
	PrimaryCategoryID   *int                  `json:"primary_category_id,omitempty"`
	Score               float64               `json:"score"`
	ConfidenceTier      string                `json:"confidence_tier"` // HIGH, MEDIUM, LOW
	MatchMethod         string                `json:"match_method"`
	SuggestedCorrection string                `json:"suggested_correction,omitempty"`
	MasterEntity        *CanonicalEntity      `json:"master_entity,omitempty"`
	DiscoveryPaths      []EntityDiscoveryPath `json:"discovery_paths,omitempty"`
}
