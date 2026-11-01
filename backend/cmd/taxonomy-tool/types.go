package main

// TaxonomyFile represents the root structure of a canonical taxonomy JSON document.
type TaxonomyFile struct {
	SchemaVersion      string         `json:"schema_version"`
	TaxonomyVersion    string         `json:"taxonomy_version"`
	SourceLanguage     string         `json:"source_language"`
	SupportedLanguages []string       `json:"supported_languages"`
	DomainCode         string         `json:"domain_code"`
	RootFamily         TaxonomyNode   `json:"root_family"`
}

// TaxonomyNode represents a single hierarchical classification entity at any depth.
type TaxonomyNode struct {
	ID             string              `json:"id"`
	Slug           string              `json:"slug"`
	Code           string              `json:"code"`
	Level          int                 `json:"level"`
	CanonicalName  string              `json:"canonical_name"`
	Definition     string              `json:"definition"`
	InclusionRules []string            `json:"inclusion_rules,omitempty"`
	ExclusionRules []string            `json:"exclusion_rules,omitempty"`
	Examples       []string            `json:"examples,omitempty"`
	Status         string              `json:"status"`
	IsLeaf         bool                `json:"is_leaf"`
	Icon           string              `json:"icon,omitempty"`
	Translations   map[string]string   `json:"translations,omitempty"`
	Aliases        map[string][]string `json:"aliases,omitempty"`
	Identifiers    []TaxonomyIdentifier `json:"identifiers,omitempty"`
	Attributes     []TaxonomyAttribute `json:"attributes,omitempty"`
	Children       []TaxonomyNode      `json:"children,omitempty"`
}

// TaxonomyIdentifier represents standardized external classification identifiers.
type TaxonomyIdentifier struct {
	Type string `json:"type"` // UNSPSC, GS1_GPC, TAXONOMIC, CAS, ISCO_08, HS_CODE
	Code string `json:"code"`
	Name string `json:"name,omitempty"`
}

// TaxonomyAttribute defines category-specific structured specification attributes.
type TaxonomyAttribute struct {
	Slug            string                 `json:"slug"`
	NameEN          string                 `json:"name_en"`
	DataType        string                 `json:"data_type"`
	Unit            string                 `json:"unit,omitempty"`
	IsRequired      bool                   `json:"is_required"`
	IsFilterable    bool                   `json:"is_filterable"`
	ValidationRules map[string]interface{} `json:"validation_rules,omitempty"`
}
