package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/your-org/used-tools-marketplace/internal/config"
	"github.com/your-org/used-tools-marketplace/internal/database"
)

func main() {
	if os.Getenv("DB_PASSWORD") == "" && os.Getenv("DATABASE_URL") == "" {
		log.Fatal("[FATAL] DB_PASSWORD environment variable is required (or set DATABASE_URL)")
	}
	if os.Getenv("JWT_SECRET") == "" {
		os.Setenv("JWT_SECRET", "audit-tool-ephemeral-jwt-secret-value-32chars!")
	}
	if os.Getenv("OTP_ENCRYPT_KEY") == "" {
		os.Setenv("OTP_ENCRYPT_KEY", "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=")
	}
	if os.Getenv("ALLOW_ORIGIN") == "" {
		os.Setenv("ALLOW_ORIGIN", "http://localhost:3000")
	}
	if os.Getenv("FRONTEND_URL") == "" {
		os.Setenv("FRONTEND_URL", "http://localhost:3000")
	}
	if os.Getenv("APP_BASE_URL") == "" {
		os.Setenv("APP_BASE_URL", "http://localhost:8080")
	}
	if os.Getenv("DB_USER") == "" {
		os.Setenv("DB_USER", "postgres")
	}
	if os.Getenv("DB_NAME") == "" {
		os.Setenv("DB_NAME", "hulubedej_db")
	}

	cfg := config.Load()
	db, err := database.NewPool(cfg)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer db.Close()
	ctx := context.Background()

	fmt.Println("=== HUNTIQ DATABASE TAXONOMY AUDIT ===")

	// 1. Check legacy categories table
	var catCount, catActive, catRoot, catChild int
	_ = db.QueryRow(ctx, "SELECT COUNT(*) FROM categories").Scan(&catCount)
	_ = db.QueryRow(ctx, "SELECT COUNT(*) FROM categories WHERE is_active = true").Scan(&catActive)
	_ = db.QueryRow(ctx, "SELECT COUNT(*) FROM categories WHERE parent_id IS NULL AND is_active = true").Scan(&catRoot)
	_ = db.QueryRow(ctx, "SELECT COUNT(*) FROM categories WHERE parent_id IS NOT NULL AND is_active = true").Scan(&catChild)
	fmt.Printf("\n[Legacy categories table]\nTotal: %d | Active: %d | Active Roots: %d | Active Children: %d\n", catCount, catActive, catRoot, catChild)

	rows, err := db.Query(ctx, "SELECT id, name, slug, parent_id, is_active FROM categories ORDER BY id")
	if err == nil {
		fmt.Println("All rows in categories:")
		for rows.Next() {
			var id int
			var name, slug string
			var parentID *int
			var isActive bool
			if err := rows.Scan(&id, &name, &slug, &parentID, &isActive); err == nil {
				pStr := "ROOT"
				if parentID != nil {
					pStr = fmt.Sprintf("parent=%d", *parentID)
				}
				fmt.Printf("  - [#%d] slug=%-25s name=%-30s active=%-5t %s\n", id, slug, name, isActive, pStr)
			}
		}
		rows.Close()
	}

	// 2. Check category_hierarchy table
	var chCount, chActive, chLeaf int
	_ = db.QueryRow(ctx, "SELECT COUNT(*) FROM category_hierarchy").Scan(&chCount)
	_ = db.QueryRow(ctx, "SELECT COUNT(*) FROM category_hierarchy WHERE is_active = true").Scan(&chActive)
	_ = db.QueryRow(ctx, "SELECT COUNT(*) FROM category_hierarchy WHERE is_leaf = true AND is_active = true").Scan(&chLeaf)
	fmt.Printf("\n[category_hierarchy table]\nTotal: %d | Active: %d | Active Leaves: %d\n", chCount, chActive, chLeaf)

	hRows, err := db.Query(ctx, "SELECT level, COUNT(*), COUNT(*) FILTER (WHERE is_leaf = true) FROM category_hierarchy WHERE is_active = true GROUP BY level ORDER BY level")
	if err == nil {
		fmt.Println("Hierarchy depth distribution:")
		for hRows.Next() {
			var lvl, cnt, leafCnt int
			_ = hRows.Scan(&lvl, &cnt, &leafCnt)
			fmt.Printf("  - Level %d: %d nodes (%d leaves)\n", lvl, cnt, leafCnt)
		}
		hRows.Close()
	}

	rootRows, err := db.Query(ctx, "SELECT id, slug, name_en, is_leaf, display_order FROM category_hierarchy WHERE level = 0 AND is_active = true ORDER BY display_order, id")
	if err == nil {
		fmt.Println("Root categories (level = 0):")
		for rootRows.Next() {
			var id, disp int
			var slug, nameEn string
			var isLeaf bool
			_ = rootRows.Scan(&id, &slug, &nameEn, &isLeaf, &disp)
			var childCount int
			_ = db.QueryRow(ctx, "SELECT COUNT(*) FROM category_hierarchy WHERE parent_id = $1 AND is_active = true", id).Scan(&childCount)
			fmt.Printf("  - [#%d] slug=%-25s name=%-35s leaf=%-5t children=%d\n", id, slug, nameEn, isLeaf, childCount)
		}
		rootRows.Close()
	}

	// 3. Check attribute_definitions & category_attributes
	var attrDefCount, catAttrCount int
	_ = db.QueryRow(ctx, "SELECT COUNT(*) FROM attribute_definitions").Scan(&attrDefCount)
	_ = db.QueryRow(ctx, "SELECT COUNT(*) FROM category_attributes").Scan(&catAttrCount)
	fmt.Printf("\n[Dynamic Attribute System]\nAttribute Definitions: %d | Category-Attribute Mappings: %d\n", attrDefCount, catAttrCount)

	attrRows, err := db.Query(ctx, "SELECT id, slug, name_en, data_type, is_required, is_filterable, unit FROM attribute_definitions ORDER BY id")
	if err == nil {
		fmt.Println("Attribute definitions:")
		for attrRows.Next() {
			var id int
			var slug, nameEn, dataType string
			var isReq, isFilt bool
			var unit *string
			_ = attrRows.Scan(&id, &slug, &nameEn, &dataType, &isReq, &isFilt, &unit)
			uStr := ""
			if unit != nil {
				uStr = fmt.Sprintf("unit=%s", *unit)
			}
			fmt.Printf("  - [#%d] slug=%-20s type=%-12s req=%-5t filt=%-5t %s name=%s\n", id, slug, dataType, isReq, isFilt, uStr, nameEn)
		}
		attrRows.Close()
	}

	// 4. Check products table taxonomy usage
	var prodTotal, prodWithCat, prodWithHier, prodWithSubcat, prodWithItem int
	_ = db.QueryRow(ctx, "SELECT COUNT(*) FROM products").Scan(&prodTotal)
	_ = db.QueryRow(ctx, "SELECT COUNT(*) FROM products WHERE category_id IS NOT NULL").Scan(&prodWithCat)
	_ = db.QueryRow(ctx, "SELECT COUNT(*) FROM products WHERE category_hierarchy_id IS NOT NULL").Scan(&prodWithHier)
	_ = db.QueryRow(ctx, "SELECT COUNT(*) FROM products WHERE subcategory IS NOT NULL AND subcategory != ''").Scan(&prodWithSubcat)
	_ = db.QueryRow(ctx, "SELECT COUNT(*) FROM products WHERE item IS NOT NULL AND item != ''").Scan(&prodWithItem)
	fmt.Printf("\n[Products Table Taxonomy Usage]\nTotal Products: %d | with category_id: %d | with category_hierarchy_id: %d | with subcategory: %d | with item: %d\n",
		prodTotal, prodWithCat, prodWithHier, prodWithSubcat, prodWithItem)

	// 5. Check if brand/model or jobs tables exist
	checkTables := []string{"brands", "models", "product_models", "brand_dictionary", "synonym_dictionary", "jobs", "job_vacancies", "resumes", "cvs", "candidate_profiles", "qualifications", "category_aliases"}
	fmt.Println("\n[Specialized Schema Check]:")
	for _, tbl := range checkTables {
		var exists bool
		_ = db.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = $1)", tbl).Scan(&exists)
		var count int
		if exists {
			_ = db.QueryRow(ctx, fmt.Sprintf("SELECT COUNT(*) FROM %s", tbl)).Scan(&count)
			fmt.Printf("  - Table %-20s : EXISTS (rows: %d)\n", tbl, count)
		} else {
			fmt.Printf("  - Table %-20s : FALSE\n", tbl)
		}
	}
}
