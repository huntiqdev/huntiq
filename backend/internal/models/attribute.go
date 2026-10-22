package models

import "time"

// AttributeDefinition represents a global attribute template
type AttributeDefinition struct {
	ID              int                    `json:"id"`
	NameEn          string                 `json:"name_en"`
	NameAm          string                 `json:"name_am,omitempty"`
	NameOm          string                 `json:"name_om,omitempty"`
	Slug            string                 `json:"slug"`
	Description     string                 `json:"description,omitempty"`
	DataType        string                 `json:"data_type"` // text, number, boolean, select, multiselect, range, date
	ValidationRules map[string]interface{} `json:"validation_rules,omitempty"`
	Unit            string                 `json:"unit,omitempty"`
	IsRequired      bool                   `json:"is_required"`
	IsSearchable    bool                   `json:"is_searchable"`
	IsFilterable    bool                   `json:"is_filterable"`
	DisplayOrder    int                    `json:"display_order"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
}

// CategoryAttribute maps an attribute to a category
type CategoryAttribute struct {
	ID           int       `json:"id"`
	CategoryID   int       `json:"category_id"`
	AttributeID  int       `json:"attribute_id"`
	IsRequired   bool      `json:"is_required"`
	IsInherited  bool      `json:"is_inherited"`
	DisplayOrder int       `json:"display_order"`
	CreatedAt    time.Time `json:"created_at"`
	
	// Populated on demand
	Attribute *AttributeDefinition `json:"attribute,omitempty"`
}

// ProductAttribute stores a product's attribute value
type ProductAttribute struct {
	ID          int                    `json:"id"`
	ProductID   int                    `json:"product_id"`
	AttributeID int                    `json:"attribute_id"`
	ValueText   *string                `json:"value_text,omitempty"`
	ValueNumber *float64               `json:"value_number,omitempty"`
	ValueBool   *bool                  `json:"value_boolean,omitempty"`
	ValueDate   *string                `json:"value_date,omitempty"`
	ValueJSON   map[string]interface{} `json:"value_json,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
	
	// Populated on demand
	Attribute *AttributeDefinition `json:"attribute,omitempty"`
}

// ApplicableAttribute represents an attribute with category-specific config
type ApplicableAttribute struct {
	AttributeID     int                    `json:"attribute_id"`
	NameEn          string                 `json:"name_en"`
	NameAm          string                 `json:"name_am,omitempty"`
	NameOm          string                 `json:"name_om,omitempty"`
	Slug            string                 `json:"slug"`
	DataType        string                 `json:"data_type"`
	ValidationRules map[string]interface{} `json:"validation_rules,omitempty"`
	Unit            string                 `json:"unit,omitempty"`
	IsRequired      bool                   `json:"is_required"`
	IsFilterable    bool                   `json:"is_filterable"`
	IsSearchable    bool                   `json:"is_searchable"`
	DisplayOrder    int                    `json:"display_order"`
}

// AttributeFilter represents filter criteria for product attributes
type AttributeFilter struct {
	AttributeID int         `json:"attribute_id"`
	Slug        string      `json:"slug"`
	Operator    string      `json:"operator"` // eq, ne, gt, gte, lt, lte, in, contains
	Value       interface{} `json:"value"`
}

// ProductWithAttributes extends Product with attribute values
type ProductWithAttributes struct {
	Product
	Attributes []ProductAttribute `json:"attributes,omitempty"`
}
