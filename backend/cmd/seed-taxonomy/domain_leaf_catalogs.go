package main

import (
	"fmt"
	"strings"
)

// LeafCategorySpec defines a concrete terminal category
type LeafCategorySpec struct {
	ParentSubSlug string
	Slug          string
	NameEN        string
	NameAM        string
	NameOM        string
	NameTI        string
}

// SubcategorySpec defines a sub-branch
type SubcategorySpec struct {
	DomainSlug string
	Slug       string
	NameEN     string
	NameAM     string
	NameOM     string
	NameTI     string
	Leaves     []LeafCategorySpec
}

// GenerateFull6000Leaves constructs over 5,000 distinct, valid, domain-specific leaf categories
func GenerateFull6000Leaves() ([]NestedCategoryDef, int) {
	var list []NestedCategoryDef
	seen := make(map[string]bool)

	addNode := func(parent, slug, en, am, om, ti string, isLeaf bool, order int) {
		slug = strings.ToLower(strings.TrimSpace(slug))
		if seen[slug] {
			slug = fmt.Sprintf("%s-%s", parent, slug)
		}
		seen[slug] = true

		list = append(list, NestedCategoryDef{
			ParentSlug:   parent,
			Slug:         slug,
			NameEN:       en,
			NameAM:       am,
			NameOM:       om,
			NameTI:       ti,
			IsLeaf:       isLeaf,
			DisplayOrder: order,
		})
	}

	// First, load all curated nodes from existing files
	existing := GenerateAllCategoriesExpanded()
	for _, n := range existing {
		addNode(n.ParentSlug, n.Slug, n.NameEN, n.NameAM, n.NameOM, n.NameTI, n.IsLeaf, n.DisplayOrder)
	}

	// Now systematically expand detailed commercial sectors to ensure full granularity across all 60 domains
	sectors := GetGranularSectorsData()
	for _, sec := range sectors {
		// Add subcategory
		addNode(sec.DomainSlug, sec.Slug, sec.NameEN, sec.NameAM, sec.NameOM, sec.NameTI, false, 1)

		// Add leaf items
		for idx, leaf := range sec.Leaves {
			addNode(sec.Slug, leaf.Slug, leaf.NameEN, leaf.NameAM, leaf.NameOM, leaf.NameTI, true, idx+1)
		}
	}

	leafCount := 0
	for _, n := range list {
		if n.IsLeaf {
			leafCount++
		}
	}

	return list, leafCount
}
