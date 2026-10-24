package main

import (
	"fmt"
	"strings"
)

// DomainMatrixEntry holds a root domain with granular sub-branches and leaf items
type DomainMatrixEntry struct {
	DomainSlug string
	SubBranches []SubBranchDef
}

type SubBranchDef struct {
	SubSlug string
	NameEN  string
	NameAM  string
	NameOM  string
	NameTI  string
	Leaves  []LeafDef
}

type LeafDef struct {
	Slug   string
	NameEN string
	NameAM string
	NameOM string
	NameTI string
}

// BuildMasterDomainMatrix constructs the full comprehensive taxonomy with thousands of valid leaves
func BuildMasterDomainMatrix() []NestedCategoryDef {
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

	// 1. Add all categories from existing generators
	for _, c := range GenerateUniversal6000Taxonomy() {
		addNode(c.ParentSlug, c.Slug, c.NameEN, c.NameAM, c.NameOM, c.NameTI, c.IsLeaf, c.DisplayOrder)
	}

	// 2. Add full matrix covering all domains up to 5,000+ leaves
	matrix := GetFullSystematicDomainMatrix()
	for _, d := range matrix {
		for subIdx, sub := range d.SubBranches {
			addNode(d.DomainSlug, sub.SubSlug, sub.NameEN, sub.NameAM, sub.NameOM, sub.NameTI, false, subIdx+1)
			for leafIdx, leaf := range sub.Leaves {
				addNode(sub.SubSlug, leaf.Slug, leaf.NameEN, leaf.NameAM, leaf.NameOM, leaf.NameTI, true, leafIdx+1)
			}
		}
	}

	return list
}
