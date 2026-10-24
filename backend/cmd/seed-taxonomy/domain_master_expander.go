package main

import (
	"fmt"
	"strings"
)

// SubSectorTemplate defines a subcategory and its concrete leaves for a domain
type SubSectorTemplate struct {
	SubSlug string
	NameEN  string
	NameAM  string
	NameOM  string
	NameTI  string
	Leaves  [][4]string // slug, nameEN, nameAM, nameOM (nameTI generated)
}

// Generate6000LeafTaxonomy creates over 5,200 distinct, authentic leaf categories across all 60 domains
func Generate6000LeafTaxonomy() []NestedCategoryDef {
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

	// 1. Add all base categories already created
	for _, c := range BuildMasterDomainMatrix() {
		addNode(c.ParentSlug, c.Slug, c.NameEN, c.NameAM, c.NameOM, c.NameTI, c.IsLeaf, c.DisplayOrder)
	}

	// 2. Generate comprehensive leaf coverage for every domain from 01 to 60
	domainCatalogues := GetDomainSystematicCatalogues()
	for _, dom := range domainCatalogues {
		for subIdx, sub := range dom.SubSectors {
			addNode(dom.DomainSlug, sub.SubSlug, sub.NameEN, sub.NameAM, sub.NameOM, sub.NameTI, false, subIdx+1)
			for leafIdx, leaf := range sub.Leaves {
				ti := fmt.Sprintf("ናይ %s", leaf[1])
				addNode(sub.SubSlug, leaf[0], leaf[1], leaf[2], leaf[3], ti, true, leafIdx+1)
			}
		}
	}

	// 3. Add all 6 domain catalog groups (Domains 01 to 60)
	allCatalogs := append([]DomainCatalogDef{}, GetCatalogDomains01To10()...)
	allCatalogs = append(allCatalogs, GetCatalogDomains11To20()...)
	allCatalogs = append(allCatalogs, GetCatalogDomains21To30()...)
	allCatalogs = append(allCatalogs, GetCatalogDomains31To40()...)
	allCatalogs = append(allCatalogs, GetCatalogDomains41To50()...)
	allCatalogs = append(allCatalogs, GetCatalogDomains51To60()...)

	for _, dom := range allCatalogs {
		for subIdx, sub := range dom.Subcategories {
			addNode(dom.RootSlug, sub.Slug, sub.NameEN, sub.NameAM, sub.NameOM, sub.NameTI, false, subIdx+1)
			for leafIdx, leaf := range sub.Items {
				addNode(sub.Slug, leaf.Slug, leaf.NameEN, leaf.NameAM, leaf.NameOM, leaf.NameTI, true, leafIdx+1)
			}
		}
	}

	return list
}

// DomainSystematicCatalog holds domain subsectors
type DomainSystematicCatalog struct {
	DomainSlug string
	SubSectors []SubSectorTemplate
}
