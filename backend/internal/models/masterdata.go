package models

import (
	"encoding/json"
	"time"
)

// FourLanguageNames represents the mandatory 4-language localized names (en, am, om, ti)
type FourLanguageNames struct {
	En string `json:"en"`
	Am string `json:"am"`
	Om string `json:"om"`
	Ti string `json:"ti"`
}

// BBox represents a geographic bounding box
type BBox struct {
	MinLat *float64 `json:"min_lat,omitempty"`
	MaxLat *float64 `json:"max_lat,omitempty"`
	MinLng *float64 `json:"min_lng,omitempty"`
	MaxLng *float64 `json:"max_lng,omitempty"`
}

// LocationCoordinates represents latitude, longitude, and optional bounding box
type LocationCoordinates struct {
	Latitude  *float64 `json:"latitude,omitempty"`
	Longitude *float64 `json:"longitude,omitempty"`
	BBox      *BBox    `json:"bbox,omitempty"`
}

// MasterLocationNode represents an individual location in the master JSON architecture
type MasterLocationNode struct {
	ID                 int                  `json:"id"`
	Code               string               `json:"code"`
	ParentID           *int                 `json:"parent_id"`
	ParentCode         *string              `json:"parent_code"`
	Type               string               `json:"type"`
	Names              FourLanguageNames    `json:"names"`
	Slug               string               `json:"slug"`
	Level              int                  `json:"level"`
	Path               string               `json:"path"`
	DisplayOrder       int                  `json:"display_order"`
	IsActive           bool                 `json:"is_active"`
	IsLeaf             bool                 `json:"is_leaf"`
	IsDeprecated       bool                 `json:"is_deprecated"`
	ReplacedByID       *int                 `json:"replaced_by_id,omitempty"`
	Coordinates        *LocationCoordinates `json:"coordinates,omitempty"`
	Aliases            []string             `json:"aliases"`
	Source             string               `json:"source"`
	VerificationStatus string               `json:"verification_status"`
}

// MasterCategoryNode represents an individual category in the master JSON architecture
type MasterCategoryNode struct {
	ID                       int               `json:"id"`
	ParentID                 *int              `json:"parent_id"`
	ParentCode               *string           `json:"parent_code"`
	ClassificationSystemCode string            `json:"classification_system_code"`
	Names                    FourLanguageNames `json:"names"`
	Slug                     string            `json:"slug"`
	Description              *string           `json:"description,omitempty"`
	Level                    int               `json:"level"`
	Path                     string            `json:"path"`
	DisplayOrder             int               `json:"display_order"`
	Icon                     *string           `json:"icon,omitempty"`
	ImageURL                 *string           `json:"image_url,omitempty"`
	IsActive                 bool              `json:"is_active"`
	IsLeaf                   bool              `json:"is_leaf"`
	IsArchived               bool              `json:"is_archived"`
	MetaTitle                *string           `json:"meta_title,omitempty"`
	MetaDescription          *string           `json:"meta_description,omitempty"`
	Aliases                  []string          `json:"aliases"`
	Redirects                []string          `json:"redirects"`
}

// MasterDataDiffItem describes an individual field or entity change during import dry-run
type MasterDataDiffItem struct {
	EntityID   int         `json:"entity_id"`
	Identifier string      `json:"identifier"`
	Action     string      `json:"action"` // CREATE, UPDATE, NOOP, CONFLICT, ERROR
	Field      string      `json:"field,omitempty"`
	OldValue   interface{} `json:"old_value,omitempty"`
	NewValue   interface{} `json:"new_value,omitempty"`
	Message    string      `json:"message,omitempty"`
}

// MasterDataImportRequest is the payload sent to the master data import API
type MasterDataImportRequest struct {
	DatasetType     string               `json:"dataset_type"` // locations | categories
	Version         string               `json:"version"`
	DryRun          bool                 `json:"dry_run"`
	Scope           string               `json:"scope,omitempty"` // e.g. ET-OR, SYS_PROD
	Records         json.RawMessage      `json:"records,omitempty"`
	LocationRecords []MasterLocationNode `json:"location_records,omitempty"`
	CategoryRecords []MasterCategoryNode `json:"category_records,omitempty"`
}

// MasterDataImportResult summarizes the outcome of a dry-run or transactional publish
type MasterDataImportResult struct {
	DatasetType  string               `json:"dataset_type"`
	DryRun       bool                 `json:"dry_run"`
	TotalRecords int                  `json:"total_records"`
	Created      int                  `json:"created"`
	Updated      int                  `json:"updated"`
	Unchanged    int                  `json:"unchanged"`
	Conflicts    int                  `json:"conflicts"`
	Errors       int                  `json:"errors"`
	Diff         []MasterDataDiffItem `json:"diff,omitempty"`
	Timestamp    time.Time            `json:"timestamp"`
}

// MasterDataStatus provides an operational overview of the master data datasets
type MasterDataStatus struct {
	LocationsCount  int       `json:"locations_count"`
	CategoriesCount int       `json:"categories_count"`
	ActiveLocations int       `json:"active_locations"`
	ActiveCats      int       `json:"active_categories"`
	LastImportAt    time.Time `json:"last_import_at"`
	SupportedLangs  []string  `json:"supported_languages"`
}

