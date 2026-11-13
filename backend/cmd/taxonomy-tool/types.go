package main

import "regexp"

// CategoryDefinition represents a category in categories.json (compiled artifact)
type CategoryDefinition struct {
	Slug   string `json:"slug"`
	Parent string `json:"parent,omitempty"`
	Icon   string `json:"icon,omitempty"`
	IsLeaf bool   `json:"is_leaf,omitempty"`
}

// LocaleEntry represents category translations and aliases in locales/<lang>.json
type LocaleEntry struct {
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	Aliases     []string `json:"aliases,omitempty"`
}

// LocalizedNames represents multilingual translations for en, am, om, ti
type LocalizedNames struct {
	EN string `json:"en"`
	AM string `json:"am"`
	OM string `json:"om"`
	TI string `json:"ti"`
}

// LocalizedAliases represents multilingual search synonyms
type LocalizedAliases struct {
	EN []string `json:"en,omitempty"`
	AM []string `json:"am,omitempty"`
	OM []string `json:"om,omitempty"`
	TI []string `json:"ti,omitempty"`
}

// NodeMeta is the metadata file (node.json) stored inside each taxonomy folder.
// The folder name IS the slug. The parent is the parent folder.
// The level is determined by folder depth. is_leaf is determined by
// whether sub-folders exist.
type NodeMeta struct {
	Name        LocalizedNames    `json:"name"`
	Icon        string            `json:"icon,omitempty"`
	Description string            `json:"description,omitempty"`
	Aliases     *LocalizedAliases `json:"aliases,omitempty"`
}

// TaxonomyTreeNode represents a hierarchical category node in taxonomy_tree.json
type TaxonomyTreeNode struct {
	Slug        string              `json:"slug"`
	Name        map[string]string   `json:"name"`
	Icon        string              `json:"icon,omitempty"`
	Layer       int                 `json:"layer,omitempty"`
	Description string              `json:"description,omitempty"`
	Aliases     map[string][]string `json:"aliases,omitempty"`
	IsLeaf      bool                `json:"is_leaf"`
	Children    []TaxonomyTreeNode  `json:"children,omitempty"`
}

var slugRegex = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)
