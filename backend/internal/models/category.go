package models

import "time"

// Category represents the legacy flat category structure (kept for backward compatibility)
type Category struct {
	ID        int        `json:"id"`
	Name      string     `json:"name"` // default display name (English)
	Slug      string     `json:"slug"`
	NameEn    string     `json:"name_en"`
	NameAm    string     `json:"name_am"`
	NameOr    string     `json:"name_or"`
	NameTi    string     `json:"name_ti"`
	SortOrder int        `json:"sort_order"`
	IsActive  bool       `json:"is_active"`
	ParentID  *int       `json:"parent_id,omitempty"`
	Children  []Category `json:"children,omitempty"`
}

// CategoryHierarchy represents the new dynamic N-level category system
type CategoryHierarchy struct {
	ID          int        `json:"id"`
	ParentID    *int       `json:"parent_id,omitempty"`
	NameEn      string     `json:"name_en"`
	NameAm      string     `json:"name_am,omitempty"`
	NameOm      string     `json:"name_om,omitempty"` // Afaan Oromo
	Slug        string     `json:"slug"`
	Description string     `json:"description,omitempty"`
	Level       int        `json:"level"`
	Path        string     `json:"path"` // Materialized path: /1/5/12/
	DisplayOrder int       `json:"display_order"`
	Icon        string     `json:"icon,omitempty"`
	IsActive    bool       `json:"is_active"`
	IsLeaf      bool       `json:"is_leaf"` // Can products be attached?
	MetaTitle   string     `json:"meta_title,omitempty"`
	MetaDescription string `json:"meta_description,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	
	// Relationships (populated on demand)
	Parent      *CategoryHierarchy   `json:"parent,omitempty"`
	Children    []CategoryHierarchy  `json:"children,omitempty"`
	Breadcrumbs []CategoryBreadcrumb `json:"breadcrumbs,omitempty"`
}

// CategoryBreadcrumb represents a single breadcrumb in the category path
type CategoryBreadcrumb struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Slug  string `json:"slug"`
	Level int    `json:"level"`
}

// CategoryWithStats includes product counts for analytics
type CategoryWithStats struct {
	CategoryHierarchy
	ProductCount  int `json:"product_count"`
	ApprovedCount int `json:"approved_count"`
	PendingCount  int `json:"pending_count"`
}

// CategoryTree represents a full category tree structure
type CategoryTree struct {
	ID       int            `json:"id"`
	Name     string         `json:"name"`
	Slug     string         `json:"slug"`
	Level    int            `json:"level"`
	IsLeaf   bool           `json:"is_leaf"`
	Count    int            `json:"count"` // Product count
	Children []CategoryTree `json:"children,omitempty"`
}

// PopularCategory holds a category with its recent activity score.
type PopularCategory struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Slug          string `json:"slug"`
	NewListings   int    `json:"new_listings"`
	OrderCount    int    `json:"order_count"`
	ActivityScore int    `json:"activity_score"`
}

// CategoryFilter for querying categories
type CategoryFilter struct {
	ParentID    *int
	Level       *int
	IsActive    *bool
	IsLeaf      *bool
	SearchTerm  string
	IncludeStats bool
	Page        int
	Limit       int
}
