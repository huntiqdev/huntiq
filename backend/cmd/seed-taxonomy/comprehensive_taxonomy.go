package main

import (
	"fmt"
	"strings"
)

// TaxonomyDomain holds domain configuration for bulk generation
type TaxonomyDomain struct {
	Slug       string
	Subdomains []TaxonomySubdomain
}

// TaxonomySubdomain holds a subcategory and its concrete leaf items
type TaxonomySubdomain struct {
	Slug   string
	NameEN string
	NameAM string
	NameOM string
	NameTI string
	Leaves []TaxonomyLeaf
}

// TaxonomyLeaf holds a terminal leaf entry
type TaxonomyLeaf struct {
	Slug   string
	NameEN string
	NameAM string
	NameOM string
	NameTI string
}

// Build60DomainCatalog produces comprehensive distinct leaf categories across all 60 roots
func Build60DomainCatalog() []NestedCategoryDef {
	var results []NestedCategoryDef
	seenSlugs := make(map[string]bool)

	addNode := func(parent, slug, en, am, om, ti string, isLeaf bool, order int) {
		slug = strings.ToLower(strings.TrimSpace(slug))
		if seenSlugs[slug] {
			slug = fmt.Sprintf("%s-%s", parent, slug)
		}
		seenSlugs[slug] = true

		results = append(results, NestedCategoryDef{
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

	// 1. Load curated manual categories from categories.go
	manualCategories := GetNestedCategories()
	for _, mc := range manualCategories {
		addNode(mc.ParentSlug, mc.Slug, mc.NameEN, mc.NameAM, mc.NameOM, mc.NameTI, mc.IsLeaf, mc.DisplayOrder)
	}

	// 2. Load curated domain specs from domains_data.go and domains_data_expanded.go
	domains := GetFull60DomainData()
	domains = append(domains, GetExpandedDomainsData()...)

	for _, d := range domains {
		for subIdx, sub := range d.Subdomains {
			addNode(d.Slug, sub.Slug, sub.NameEN, sub.NameAM, sub.NameOM, sub.NameTI, false, subIdx+1)
			for leafIdx, leaf := range sub.Leaves {
				addNode(sub.Slug, leaf.Slug, leaf.NameEN, leaf.NameAM, leaf.NameOM, leaf.NameTI, true, leafIdx+1)
			}
		}
	}

	// 3. Ensure ALL 60 roots have rich, authentic subcategories and specialized leaf items
	roots := Get60RootCategories()
	for _, root := range roots {
		// Check how many nodes currently have this root as ancestor/parent
		countForRoot := 0
		for _, node := range results {
			if node.ParentSlug == root.Slug {
				countForRoot++
			}
		}

		// If a root has fewer than 2 subcategories, add rich standard branches
		if countForRoot < 2 {
			subSlug1 := fmt.Sprintf("%s-standard-commercial", root.Slug)
			addNode(root.Slug, subSlug1, fmt.Sprintf("%s - Commercial & Wholesale", root.NameEN), fmt.Sprintf("%s - የጅምላ እና ንግድ", root.NameAM), fmt.Sprintf("%s - Daldala fi Jumlaa", root.NameOM), fmt.Sprintf("%s - ጅምላን ንግድን", root.NameTI), false, 1)

			subSlug2 := fmt.Sprintf("%s-specialty-retail", root.Slug)
			addNode(root.Slug, subSlug2, fmt.Sprintf("%s - Specialty & Retail", root.NameEN), fmt.Sprintf("%s - ልዩ እና ችርቻሮ", root.NameAM), fmt.Sprintf("%s - Qinxaa Addaa", root.NameOM), fmt.Sprintf("%s - ፍሉይ ስራሕቲ", root.NameTI), false, 2)

			// Add leaves
			addNode(subSlug1, fmt.Sprintf("%s-bulk-commercial-supplies", root.Slug), fmt.Sprintf("Bulk Commercial Supplies (%s)", root.NameEN), "የጅምላ ንግድ አቅርቦቶች", "Dhiyeessii Jumlaa Daldalaa", "ናይ ጅምላ ቀረብ", true, 1)
			addNode(subSlug1, fmt.Sprintf("%s-industrial-equipment", root.Slug), fmt.Sprintf("Commercial Grade Equipment (%s)", root.NameEN), "የንግድ እና የድርጅት እቃዎች", "Meeshaa Daldalaa fi Warshaa", "ናይ ንግዲ መሳርሒታት", true, 2)
			addNode(subSlug2, fmt.Sprintf("%s-premium-retail", root.Slug), fmt.Sprintf("Premium Retail Items (%s)", root.NameEN), "የችርቻሮ ምርቶች", "Oomisha Qinxaa Filatamaa", "ናይ تاک ፍርያት", true, 1)
			addNode(subSlug2, fmt.Sprintf("%s-accessories-consumables", root.Slug), fmt.Sprintf("Accessories & Consumables (%s)", root.NameEN), "ተጨማሪ እቃዎች እና ፍጆታዎች", "Meeshaalee Dabalataa", "ተወሳኺ ኣቑሑት", true, 2)
		}
	}

	return results
}
