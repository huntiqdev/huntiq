package main

import (
	"fmt"
	"strings"
)

// TaxonomyCategorySpec holds a definition for bulk category hierarchy generation
type TaxonomyCategorySpec struct {
	ParentSlug string
	Slug       string
	NameEN     string
	NameAM     string
	NameOM     string
	NameTI     string
	IsLeaf     bool
}

// GenerateUniversal6000Taxonomy generates thousands of distinct, valid, useful marketplace categories across all 60 domains
func GenerateUniversal6000Taxonomy() []NestedCategoryDef {
	var list []NestedCategoryDef
	seenSlugs := make(map[string]bool)

	addNode := func(parent, slug, en, am, om, ti string, isLeaf bool, order int) {
		slug = strings.ToLower(strings.TrimSpace(slug))
		if seenSlugs[slug] {
			slug = fmt.Sprintf("%s-%s", parent, slug)
		}
		seenSlugs[slug] = true

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

	// 1. Add all categories from existing catalog generators
	for _, c := range GenerateAllCategoriesExpanded() {
		addNode(c.ParentSlug, c.Slug, c.NameEN, c.NameAM, c.NameOM, c.NameTI, c.IsLeaf, c.DisplayOrder)
	}

	// 2. Add granular sector data
	for _, sec := range GetGranularSectorsData() {
		addNode(sec.DomainSlug, sec.Slug, sec.NameEN, sec.NameAM, sec.NameOM, sec.NameTI, false, 1)
		for idx, leaf := range sec.Leaves {
			addNode(sec.Slug, leaf.Slug, leaf.NameEN, leaf.NameAM, leaf.NameOM, leaf.NameTI, true, idx+1)
		}
	}

	// 3. Systematically expand every domain with its standard commercial product families & leaf specializations
	domainDefinitions := GetUniversalDomainLeafDatabase()
	for _, dDef := range domainDefinitions {
		for subIdx, sub := range dDef.Subcategories {
			addNode(dDef.RootSlug, sub.Slug, sub.NameEN, sub.NameAM, sub.NameOM, sub.NameTI, false, subIdx+1)
			for leafIdx, leaf := range sub.Items {
				addNode(sub.Slug, leaf.Slug, leaf.NameEN, leaf.NameAM, leaf.NameOM, leaf.NameTI, true, leafIdx+1)
			}
		}
	}

	return list
}
