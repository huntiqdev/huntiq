package repository_test

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/your-org/used-tools-marketplace/internal/models"
	"github.com/your-org/used-tools-marketplace/internal/repository"
)

func TestGetCategoryTree_NoOrphanRoots(t *testing.T) {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "postgres://postgres@/hulubedej_db?host=/tmp&port=5433"
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		t.Skipf("skipping db integration test: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		t.Skipf("skipping db integration test: %v", err)
	}

	repo := repository.NewCategoryHierarchyRepo(pool)

	tree, err := repo.GetCategoryTree(ctx, "en")
	if err != nil {
		t.Fatalf("GetCategoryTree failed: %v", err)
	}

	t.Logf("Total root categories returned: %d", len(tree))

	// Clean root categories should be returned (8 core marketplace + 50 economic sectors), NOT 16,000!
	if len(tree) < 50 || len(tree) > 70 {
		t.Fatalf("expected 50-70 clean root categories, got %d", len(tree))
	}

	for _, root := range tree {
		if root.Level != 0 {
			t.Errorf("expected level 0 root, got level %d for %s", root.Level, root.Name)
		}
		if root.ParentID != nil && *root.ParentID != 0 {
			t.Errorf("expected nil or 0 parent_id for root, got %d for %s", *root.ParentID, root.Name)
		}
		// Root names must be clean domain titles, never raw leaf products
		if len(root.Children) > 0 {
			t.Logf("Root '%s' has %d superfamilies", root.Name, len(root.Children))
		}
	}
}

func TestSearchCategories_ReturnsFullBreadcrumbs(t *testing.T) {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "postgres://postgres@/hulubedej_db?host=/tmp&port=5433"
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		t.Skipf("skipping db integration test: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		t.Skipf("skipping db integration test: %v", err)
	}

	repo := repository.NewCategoryHierarchyRepo(pool)

	results, err := repo.SearchCategories(ctx, "beer", "en", 10)
	if err != nil {
		t.Fatalf("SearchCategories failed: %v", err)
	}

	if len(results) == 0 {
		t.Fatalf("expected search results for 'beer', got 0")
	}

	foundBeer := false
	for _, res := range results {
		t.Logf("Match: %s (level %d) | Hierarchy: %s", res.Name, res.Level, res.Hierarchy)
		if len(res.Breadcrumbs) > 0 && res.Name == "Beer Products" {
			foundBeer = true
			if res.Breadcrumbs[0].Name != "Food & Beverage" {
				t.Errorf("expected root 'Food & Beverage', got '%s'", res.Breadcrumbs[0].Name)
			}
		}
	}

	if !foundBeer {
		t.Errorf("expected 'Beer Products' with breadcrumbs")
	}

	// Test Amharic & Alias search: "ድሪል" (drill)
	amResults, err := repo.SearchCategories(ctx, "ድሪል", "am", 10)
	if err != nil {
		t.Fatalf("SearchCategories for 'ድሪል' failed: %v", err)
	}
	if len(amResults) == 0 {
		t.Fatalf("expected search results for Amharic alias 'ድሪል', got 0")
	}
	t.Logf("Found %d results for Amharic 'ድሪል'. Top match: %s (%s)", len(amResults), amResults[0].Name, amResults[0].Hierarchy)
}

func TestCameraLenses_FullHierarchyAndAttributes(t *testing.T) {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "postgres://postgres@/hulubedej_db?host=/tmp&port=5433"
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		t.Skipf("skipping db integration test: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		t.Skipf("skipping db integration test: %v", err)
	}

	repo := repository.NewCategoryHierarchyRepo(pool)

	// 1. Search for "camera lenses"
	results, err := repo.SearchCategories(ctx, "camera lenses", "en", 10)
	if err != nil {
		t.Fatalf("SearchCategories for 'camera lenses' failed: %v", err)
	}
	if len(results) == 0 {
		t.Fatalf("expected search results for 'camera lenses', got 0")
	}
	foundLens := false
	for _, r := range results {
		if r.ID == 1380 {
			foundLens = true
			if len(r.Breadcrumbs) < 2 {
				t.Fatalf("expected at least 2 breadcrumb levels for Camera Lenses, got %d", len(r.Breadcrumbs))
			}
			t.Logf("Found Camera Lenses: Hierarchy: %s", r.Hierarchy)
		}
	}
	if !foundLens {
		t.Errorf("expected Camera Lenses (ID 1380) in search results")
	}

	// 2. Verify Children of Camera Lenses (1380) - must NOT contain Enterprise Hardware
	parentID := 1380
	children, err := repo.ListCategories(ctx, models.CategoryFilter{ParentID: &parentID, Lang: "en"})
	if err != nil {
		t.Fatalf("ListCategories(ParentID: 1380) failed: %v", err)
	}
	if len(children) < 5 {
		t.Errorf("expected at least 5 lens subcategories under 1380, got %d", len(children))
	}
	for _, c := range children {
		t.Logf("Subcategory under 1380: %s (%s)", c.NameEn, c.Slug)
		if strings.Contains(c.NameEn, "Enterprise Hardware") || strings.Contains(c.NameEn, "Telecom") {
			t.Errorf("found unwanted telecom fallback dummy category under 1380: %s", c.NameEn)
		}
	}

	// 3. Verify Dynamic Attributes for Camera Lenses (1380)
	var attrCount int
	err = pool.QueryRow(ctx, "SELECT count(*) FROM get_category_attributes(1380)").Scan(&attrCount)
	if err != nil {
		t.Fatalf("get_category_attributes(1380) failed: %v", err)
	}
	if attrCount != 10 {
		t.Errorf("expected 10 attributes for category 1380, got %d", attrCount)
	}

	// 4. Verify child category inherits attributes
	if len(children) > 0 {
		var childAttrCount int
		err = pool.QueryRow(ctx, "SELECT count(*) FROM get_category_attributes($1)", children[0].ID).Scan(&childAttrCount)
		if err != nil {
			t.Fatalf("get_category_attributes(%d) failed: %v", children[0].ID, err)
		}
		if childAttrCount != 10 {
			t.Errorf("expected child category %d to inherit 10 attributes, got %d", children[0].ID, childAttrCount)
		}
	}
}


