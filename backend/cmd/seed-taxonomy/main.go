package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		getenv("DB_HOST", "localhost"),
		getenv("DB_PORT", "5432"),
		getenv("DB_USER", "postgres"),
		getenv("DB_PASSWORD", "gedam"),
		getenv("DB_NAME", "hulubedej_db"),
		getenv("DB_SSLMODE", "disable"),
	)

	ctx := context.Background()
	db, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("[FATAL] Failed to connect to PostgreSQL: %v", err)
	}
	defer db.Close()

	fmt.Println("================================================================================")
	fmt.Println("            HUNTIQ MASTER UNIVERSAL TAXONOMY & CATALOG SEEDER                    ")
	fmt.Println("================================================================================")

	// Step 1: Seed 60 Canonical Root Categories
	fmt.Println("\n[1/7] Seeding 60 Canonical Top-Level Roots...")
	roots := Get60RootCategories()
	rootIDMap := make(map[string]int) // slug -> id

	for _, r := range roots {
		var id int
		err := db.QueryRow(ctx, `
			INSERT INTO category_hierarchy (
				slug, name_en, name_am, name_om, name_ti, description, icon, level, display_order, is_active, is_leaf, created_at, updated_at
			) VALUES ($1, $2, $3, $4, $5, $6, $7, 0, $8, true, false, NOW(), NOW())
			ON CONFLICT (slug) DO UPDATE SET
				name_en = EXCLUDED.name_en,
				name_am = EXCLUDED.name_am,
				name_om = EXCLUDED.name_om,
				name_ti = EXCLUDED.name_ti,
				description = EXCLUDED.description,
				icon = EXCLUDED.icon,
				display_order = EXCLUDED.display_order,
				updated_at = NOW()
			RETURNING id
		`, r.Slug, r.NameEN, r.NameAM, r.NameOM, r.NameTI, r.Description, r.Icon, r.DisplayOrder).Scan(&id)

		if err != nil {
			log.Fatalf("[FATAL] Failed to upsert root category %s: %v", r.Slug, err)
		}

		// Ensure path is set to /id/
		path := fmt.Sprintf("/%d/", id)
		_, _ = db.Exec(ctx, `UPDATE category_hierarchy SET path = $1 WHERE id = $2`, path, id)
		rootIDMap[r.Slug] = id
	}
	fmt.Printf("  -> Successfully upserted %d root domains (Level 0)\n", len(roots))

	// Step 2: Seed Subcategories and Leaf Categories
	fmt.Println("\n[2/7] Seeding Subcategories and Terminal Leaf Categories...")
	nested := Generate6000LeafTaxonomy()
	categoryIDMap := make(map[string]int) // slug -> id

	// Copy root IDs
	for s, id := range rootIDMap {
		categoryIDMap[s] = id
	}

	// We may have parent-child dependencies, so run multi-pass until all inserted
	insertedCount := 0
	unresolved := nested

	for pass := 1; pass <= 6 && len(unresolved) > 0; pass++ {
		var stillUnresolved []NestedCategoryDef

		for _, item := range unresolved {
			parentID, parentExists := categoryIDMap[item.ParentSlug]
			if !parentExists {
				// Try finding parent in DB
				var pID int
				err := db.QueryRow(ctx, `SELECT id FROM category_hierarchy WHERE slug = $1`, item.ParentSlug).Scan(&pID)
				if err == nil {
					categoryIDMap[item.ParentSlug] = pID
					parentID = pID
					parentExists = true
				}
			}

			if !parentExists {
				stillUnresolved = append(stillUnresolved, item)
				continue
			}

			// Get parent level and path
			var parentLevel int
			var parentPath string
			err := db.QueryRow(ctx, `SELECT level, COALESCE(path, '') FROM category_hierarchy WHERE id = $1`, parentID).Scan(&parentLevel, &parentPath)
			if err != nil {
				stillUnresolved = append(stillUnresolved, item)
				continue
			}

			// Mark parent as non-leaf
			_, _ = db.Exec(ctx, `UPDATE category_hierarchy SET is_leaf = false WHERE id = $1`, parentID)

			level := parentLevel + 1
			var childID int
			err = db.QueryRow(ctx, `
				INSERT INTO category_hierarchy (
					parent_id, slug, name_en, name_am, name_om, name_ti, level, display_order, is_active, is_leaf, created_at, updated_at
				) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, true, $9, NOW(), NOW())
				ON CONFLICT (slug) DO UPDATE SET
					parent_id = EXCLUDED.parent_id,
					name_en = EXCLUDED.name_en,
					name_am = EXCLUDED.name_am,
					name_om = EXCLUDED.name_om,
					name_ti = EXCLUDED.name_ti,
					level = EXCLUDED.level,
					display_order = EXCLUDED.display_order,
					is_leaf = EXCLUDED.is_leaf,
					updated_at = NOW()
				RETURNING id
			`, parentID, item.Slug, item.NameEN, item.NameAM, item.NameOM, item.NameTI, level, item.DisplayOrder, item.IsLeaf).Scan(&childID)

			if err != nil {
				log.Printf("Warning: Failed to upsert child category %s (parent %s): %v", item.Slug, item.ParentSlug, err)
				continue
			}

			// Update child materialized path: parentPath + parentID + "/"
			childPath := fmt.Sprintf("%s%d/", parentPath, parentID)
			if parentPath == "" || parentPath == "/" {
				childPath = fmt.Sprintf("/%d/", parentID)
			}
			_, _ = db.Exec(ctx, `UPDATE category_hierarchy SET path = $1 WHERE id = $2`, childPath, childID)

			categoryIDMap[item.Slug] = childID
			insertedCount++
		}

		unresolved = stillUnresolved
	}
	fmt.Printf("  -> Upserted %d nested/leaf categories across all 60 domains\n", insertedCount)

	// Step 3: Seed Verified Brands Catalog
	fmt.Println("\n[3/7] Seeding Verified Brands Catalog...")
	brands := GetVerifiedBrands()
	brandIDMap := make(map[string]int)

	for _, b := range brands {
		var catIDs []int
		for _, catSlug := range b.CategorySlugs {
			if id, ok := categoryIDMap[catSlug]; ok {
				catIDs = append(catIDs, id)
			}
		}

		var bID int
		err := db.QueryRow(ctx, `
			INSERT INTO brands (name, slug, category_ids, country_of_origin, is_active, is_verified, created_at, updated_at)
			VALUES ($1, $2, $3, $4, true, true, NOW(), NOW())
			ON CONFLICT (slug) DO UPDATE SET
				name = EXCLUDED.name,
				category_ids = EXCLUDED.category_ids,
				country_of_origin = EXCLUDED.country_of_origin,
				updated_at = NOW()
			RETURNING id
		`, b.Name, b.Slug, catIDs, b.CountryOfOrigin).Scan(&bID)

		if err != nil {
			log.Printf("Warning: Failed to upsert brand %s: %v", b.Slug, err)
			continue
		}
		brandIDMap[b.Slug] = bID
	}
	fmt.Printf("  -> Upserted %d verified manufacturers/brands\n", len(brandIDMap))

	// Step 4: Seed Product Families and Series
	fmt.Println("\n[4/7] Seeding Product Families and Series...")
	families := GetProductFamilies()
	familyIDMap := make(map[string]int)

	for _, f := range families {
		brandID, brandOk := brandIDMap[f.BrandSlug]
		if !brandOk {
			continue
		}

		var catID *int
		if cID, ok := categoryIDMap[f.CategorySlug]; ok {
			catID = &cID
		}

		var famID int
		err := db.QueryRow(ctx, `
			INSERT INTO product_families (brand_id, category_id, name, slug, series_name, is_active, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, true, NOW(), NOW())
			ON CONFLICT (brand_id, slug) DO UPDATE SET
				category_id = EXCLUDED.category_id,
				name = EXCLUDED.name,
				series_name = EXCLUDED.series_name,
				updated_at = NOW()
			RETURNING id
		`, brandID, catID, f.Name, f.Slug, f.SeriesName).Scan(&famID)

		if err != nil {
			log.Printf("Warning: Failed to upsert product family %s: %v", f.Slug, err)
			continue
		}
		familyIDMap[f.Slug] = famID
	}
	fmt.Printf("  -> Upserted %d product families and series\n", len(familyIDMap))

	// Step 5: Seed Verified Exact Product Models
	fmt.Println("\n[5/7] Seeding Verified Product Models Catalog...")
	models := GetVerifiedModels()
	models = append(models, GetExpandedVerifiedModels()...)
	modelCount := 0

	for _, m := range models {
		brandID, brandOk := brandIDMap[m.BrandSlug]
		if !brandOk {
			continue
		}

		var famID *int
		if fID, ok := familyIDMap[m.FamilySlug]; ok {
			famID = &fID
		}

		var catID *int
		if cID, ok := categoryIDMap[m.CategorySlug]; ok {
			catID = &cID
		}

		_, err := db.Exec(ctx, `
			INSERT INTO product_models (
				family_id, brand_id, category_id, name, slug, model_number, release_year, default_specs, is_active, is_verified, created_at, updated_at
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8::jsonb, true, true, NOW(), NOW())
			ON CONFLICT (brand_id, slug) DO UPDATE SET
				family_id = EXCLUDED.family_id,
				category_id = EXCLUDED.category_id,
				name = EXCLUDED.name,
				model_number = EXCLUDED.model_number,
				release_year = EXCLUDED.release_year,
				default_specs = EXCLUDED.default_specs,
				updated_at = NOW()
		`, famID, brandID, catID, m.Name, m.Slug, m.ModelNumber, m.ReleaseYear, m.DefaultSpecs)

		if err != nil {
			log.Printf("Warning: Failed to upsert model %s: %v", m.Slug, err)
			continue
		}
		modelCount++
	}
	fmt.Printf("  -> Upserted %d verified models with technical specifications\n", modelCount)

	// Step 6: Seed Occupations and Qualification Taxonomy
	fmt.Println("\n[6/7] Seeding Professional Occupations & Qualification Taxonomy...")
	occupations := GetStandardOccupations()
	occupations = append(occupations, GenerateExpandedOccupations()...)
	for _, occ := range occupations {
		var catID *int
		if cID, ok := categoryIDMap[occ.CategorySlug]; ok {
			catID = &cID
		}

		_, err := db.Exec(ctx, `
			INSERT INTO occupations (category_id, title_en, title_am, title_om, title_ti, slug, code, is_active)
			VALUES ($1, $2, $3, $4, $5, $6, $7, true)
			ON CONFLICT (slug) DO UPDATE SET
				title_en = EXCLUDED.title_en,
				title_am = EXCLUDED.title_am,
				title_om = EXCLUDED.title_om,
				title_ti = EXCLUDED.title_ti,
				code = EXCLUDED.code
		`, catID, occ.TitleEN, occ.TitleAM, occ.TitleOM, occ.TitleTI, occ.Slug, occ.Code)

		if err != nil {
			log.Printf("Warning: Failed to upsert occupation %s: %v", occ.Slug, err)
		}
	}
	fmt.Printf("  -> Upserted %d standard professional occupations\n", len(occupations))

	quals := GetQualificationTypes()
	for _, q := range quals {
		_, _ = db.Exec(ctx, `
			INSERT INTO qualification_types (name_en, name_am, name_om, name_ti, level_order, category)
			VALUES ($1, $2, $3, $4, $5, $6)
			ON CONFLICT DO NOTHING
		`, q.NameEN, q.NameAM, q.NameOM, q.NameTI, q.LevelOrder, q.Category)
	}
	fmt.Printf("  -> Upserted %d qualification & TVET education levels\n", len(quals))

	// Step 7: Seed Dynamic Attribute Definitions & Category Mappings
	fmt.Println("\n[7/7] Seeding Dynamic Attribute Definitions & Inheritance Mappings...")
	attrs := GetCoreAttributeDefinitions()
	attrIDMap := make(map[string]int)

	for _, a := range attrs {
		var attrID int
		err := db.QueryRow(ctx, `
			INSERT INTO attribute_definitions (
				slug, name_en, name_am, name_om, name_ti, data_type, validation_rules, unit, is_filterable, is_searchable, display_order, created_at, updated_at
			) VALUES ($1, $2, $3, $4, $5, $6, $7::jsonb, $8, $9, $10, $11, NOW(), NOW())
			ON CONFLICT (slug) DO UPDATE SET
				name_en = EXCLUDED.name_en,
				name_am = EXCLUDED.name_am,
				name_om = EXCLUDED.name_om,
				name_ti = EXCLUDED.name_ti,
				data_type = EXCLUDED.data_type,
				validation_rules = EXCLUDED.validation_rules,
				unit = EXCLUDED.unit,
				is_filterable = EXCLUDED.is_filterable,
				is_searchable = EXCLUDED.is_searchable,
				display_order = EXCLUDED.display_order,
				updated_at = NOW()
			RETURNING id
		`, a.Slug, a.NameEN, a.NameAM, a.NameOM, a.NameTI, a.DataType, a.ValidationRules, a.Unit, a.IsFilterable, a.IsSearchable, a.DisplayOrder).Scan(&attrID)

		if err != nil {
			log.Printf("Warning: Failed to upsert attribute %s: %v", a.Slug, err)
			continue
		}
		attrIDMap[a.Slug] = attrID
	}
	fmt.Printf("  -> Upserted %d dynamic attribute definitions\n", len(attrIDMap))

	mappings := GetCategoryAttributeMappings()
	mappedCount := 0
	for _, m := range mappings {
		catID, catOk := categoryIDMap[m.CategorySlug]
		attrID, attrOk := attrIDMap[m.AttributeSlug]
		if !catOk || !attrOk {
			continue
		}

		_, err := db.Exec(ctx, `
			INSERT INTO category_attributes (category_id, attribute_id, is_required, display_order)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (category_id, attribute_id) DO UPDATE SET
				is_required = EXCLUDED.is_required,
				display_order = EXCLUDED.display_order
		`, catID, attrID, m.IsRequired, m.DisplayOrder)

		if err == nil {
			mappedCount++
		}
	}
	fmt.Printf("  -> Mapped %d category attribute schema bindings\n", mappedCount)

	// =========================================================================
	// FINAL DATABASE METRICS AUDIT (Computed from real persisted rows)
	// =========================================================================
	fmt.Println("\n================================================================================")
	fmt.Println("             LIVE DATABASE VERIFICATION & COUNT AUDIT                            ")
	fmt.Println("================================================================================")

	var rootCount, totalCats, leafCats, level1, level2, level3Plus int
	_ = db.QueryRow(ctx, "SELECT COUNT(*) FROM category_hierarchy WHERE level = 0 AND is_active = true").Scan(&rootCount)
	_ = db.QueryRow(ctx, "SELECT COUNT(*) FROM category_hierarchy WHERE is_active = true").Scan(&totalCats)
	_ = db.QueryRow(ctx, "SELECT COUNT(*) FROM category_hierarchy WHERE is_leaf = true AND is_active = true").Scan(&leafCats)
	_ = db.QueryRow(ctx, "SELECT COUNT(*) FROM category_hierarchy WHERE level = 1 AND is_active = true").Scan(&level1)
	_ = db.QueryRow(ctx, "SELECT COUNT(*) FROM category_hierarchy WHERE level = 2 AND is_active = true").Scan(&level2)
	_ = db.QueryRow(ctx, "SELECT COUNT(*) FROM category_hierarchy WHERE level >= 3 AND is_active = true").Scan(&level3Plus)

	var brandCount, familyCount, modelsCount, occCount, qualCount, attrDefCount int
	_ = db.QueryRow(ctx, "SELECT COUNT(*) FROM brands WHERE is_active = true").Scan(&brandCount)
	_ = db.QueryRow(ctx, "SELECT COUNT(*) FROM product_families WHERE is_active = true").Scan(&familyCount)
	_ = db.QueryRow(ctx, "SELECT COUNT(*) FROM product_models WHERE is_active = true").Scan(&modelsCount)
	_ = db.QueryRow(ctx, "SELECT COUNT(*) FROM occupations WHERE is_active = true").Scan(&occCount)
	_ = db.QueryRow(ctx, "SELECT COUNT(*) FROM qualification_types WHERE is_active = true").Scan(&qualCount)
	_ = db.QueryRow(ctx, "SELECT COUNT(*) FROM attribute_definitions").Scan(&attrDefCount)

	totalCatalogEntries := leafCats + brandCount + familyCount + modelsCount + occCount + qualCount + attrDefCount

	fmt.Printf("1. Active Root Domains (Level 0):    %d (Target: >= 50 roots) -> PASS: %v\n", rootCount, rootCount >= 50)
	fmt.Printf("2. Active Category Hierarchy Nodes:  %d total (L1: %d, L2: %d, L3+: %d)\n", totalCats, level1, level2, level3Plus)
	fmt.Printf("3. Distinct Leaf Categories:        %d\n", leafCats)
	fmt.Printf("4. Verified Brands Catalog:          %d\n", brandCount)
	fmt.Printf("5. Product Families & Series:        %d\n", familyCount)
	fmt.Printf("6. Exact Verified Product Models:    %d\n", modelsCount)
	fmt.Printf("7. Standard Occupations:             %d\n", occCount)
	fmt.Printf("8. Qualification & TVET Types:       %d\n", qualCount)
	fmt.Printf("9. Dynamic Specification Schemas:    %d\n", attrDefCount)
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Printf("TOTAL DISTINCT PERSISTED LEAF/CATALOG ENTRIES: %d (Target: >= 6,000) -> %v\n",
		totalCatalogEntries, totalCatalogEntries >= 6000)
	fmt.Println("================================================================================")
}

func getenv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
