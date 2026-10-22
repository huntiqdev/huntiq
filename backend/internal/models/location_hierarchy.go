package models

import "time"

// LocationHierarchy represents an administrative location node in the data-driven hierarchy.
// Supports Country -> Region -> Zone/City -> Subcity/Woreda -> Kebele.
type LocationHierarchy struct {
	ID                 int                  `json:"id"`
	ParentID           *int                 `json:"parent_id,omitempty"`
	LocationType       string               `json:"location_type"` // country, region, zone, city, subcity, woreda, kebele
	NameEn             string               `json:"name_en"`
	NameAm             string               `json:"name_am,omitempty"`
	NameOm             string               `json:"name_om,omitempty"`
	NameTi             string               `json:"name_ti,omitempty"`
	Slug               string               `json:"slug"`
	AdminCode          string               `json:"admin_code,omitempty"`
	Level              int                  `json:"level"`
	Path               string               `json:"path"` // Materialized path: /1/5/12/
	DisplayOrder       int                  `json:"display_order"`
	IsActive           bool                 `json:"is_active"`
	IsLeaf             bool                 `json:"is_leaf"`
	IsDeprecated       bool                 `json:"is_deprecated"`
	ReplacedByID       *int                 `json:"replaced_by_id,omitempty"`
	Latitude           *float64             `json:"latitude,omitempty"`
	Longitude          *float64             `json:"longitude,omitempty"`
	MinLat             *float64             `json:"min_lat,omitempty"`
	MaxLat             *float64             `json:"max_lat,omitempty"`
	MinLng             *float64             `json:"min_lng,omitempty"`
	MaxLng             *float64             `json:"max_lng,omitempty"`
	Aliases            []string             `json:"aliases,omitempty"`
	Source             string               `json:"source,omitempty"`
	VerificationStatus string               `json:"verification_status,omitempty"`
	CreatedAt          time.Time            `json:"created_at"`
	UpdatedAt          time.Time            `json:"updated_at"`

	// Relationships & Breadcrumbs (populated on demand)
	Children    []LocationHierarchy `json:"children,omitempty"`
	Breadcrumbs []LocationBreadcrumb `json:"breadcrumbs,omitempty"`
}

// LocationBreadcrumb represents a single breadcrumb item along the parent path
type LocationBreadcrumb struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Slug         string `json:"slug"`
	Level        int    `json:"level"`
	LocationType string `json:"location_type"`
}

// LocationTree represents a nested node in the location tree
type LocationTree struct {
	ID           int            `json:"id"`
	Name         string         `json:"name"`
	Slug         string         `json:"slug"`
	LocationType string         `json:"location_type"`
	AdminCode    string         `json:"admin_code,omitempty"`
	Level        int            `json:"level"`
	Latitude     *float64       `json:"latitude,omitempty"`
	Longitude    *float64       `json:"longitude,omitempty"`
	IsLeaf       bool           `json:"is_leaf"`
	Children     []LocationTree `json:"children,omitempty"`
}

// GeocodingMatch represents the result of reverse geocoding resolved to the location hierarchy
type GeocodingMatch struct {
	LocationID   int                  `json:"location_id"`
	LocationType string               `json:"location_type"`
	Name         string               `json:"name"`
	Country      string               `json:"country"`
	Region       string               `json:"region"`
	Zone         string               `json:"zone,omitempty"`
	City         string               `json:"city,omitempty"`
	Subcity      string               `json:"subcity,omitempty"`
	Woreda       string               `json:"woreda,omitempty"`
	Latitude     float64              `json:"latitude"`
	Longitude    float64              `json:"longitude"`
	Confidence   float64              `json:"confidence"` // 0.0 to 1.0
	Breadcrumbs  []LocationBreadcrumb `json:"breadcrumbs"`
	RawAddress   string               `json:"raw_address,omitempty"`
	Provider     string               `json:"provider"`
}

// LocationFilter options for listing locations
type LocationFilter struct {
	ParentID     *int   `json:"parent_id"`
	LocationType string `json:"location_type"`
	Level        *int   `json:"level"`
	IsActive     *bool  `json:"is_active"`
	SearchTerm   string `json:"search_term"`
	Page         int    `json:"page"`
	Limit        int    `json:"limit"`
}
