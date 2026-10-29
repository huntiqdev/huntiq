package repository

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/your-org/used-tools-marketplace/internal/config"
	"github.com/your-org/used-tools-marketplace/internal/database"
	"github.com/your-org/used-tools-marketplace/internal/models"
)

func setupTestDB(t *testing.T) (*pgxpool.Pool, func()) {
	if os.Getenv("DB_PASSWORD") == "" && os.Getenv("DATABASE_URL") == "" {
		t.Skip("Skipping live database test: DB_PASSWORD or DATABASE_URL not set in environment")
		return nil, nil
	}

	cfg := config.NewTestConfig()
	db, err := database.NewPool(cfg)
	if err != nil {
		t.Skipf("Skipping live database test: DB connection failed: %v", err)
		return nil, nil
	}

	cleanup := func() {
		db.Close()
	}
	return db, cleanup
}

func TestCategoryHierarchyRepo_GetBreadcrumbs_ValidCategoryWithAncestors(t *testing.T) {
	db, cleanup := setupTestDB(t)
	if db == nil {
		return
	}
	defer cleanup()

	repo := NewCategoryHierarchyRepo(db)
	ctx := context.Background()

	// Category 14 is "Computers & Laptops", child of Category 1 "Electronics"
	crumbs, err := repo.GetBreadcrumbs(ctx, 14, "en")
	if err != nil {
		t.Fatalf("GetBreadcrumbs failed for valid category with ancestors: %v", err)
	}

	if len(crumbs) < 2 {
		t.Fatalf("expected at least 2 breadcrumbs for category 14, got %d", len(crumbs))
	}

	// Verify root ancestor (first in list)
	root := crumbs[0]
	if root.ID != 1 {
		t.Errorf("expected root ID 1, got %d", root.ID)
	}
	if root.Name != "Electronics" {
		t.Errorf("expected root name 'Electronics', got %q", root.Name)
	}
	if root.Level != 0 {
		t.Errorf("expected root level 0, got %d", root.Level)
	}

	// Verify target category (last in list)
	target := crumbs[len(crumbs)-1]
	if target.ID != 14 {
		t.Errorf("expected target ID 14, got %d", target.ID)
	}
	if target.Name != "Computers & Laptops" {
		t.Errorf("expected target name 'Computers & Laptops', got %q", target.Name)
	}
	if target.Level != 1 {
		t.Errorf("expected target level 1, got %d", target.Level)
	}

	// Test Amharic localization fallback
	amCrumbs, err := repo.GetBreadcrumbs(ctx, 14, "am")
	if err != nil {
		t.Fatalf("GetBreadcrumbs with amharic failed: %v", err)
	}
	if len(amCrumbs) >= 2 {
		if amCrumbs[1].Name != "ኮምፒውተሮችና ላፕቶፖች" {
			t.Errorf("expected Amharic name 'ኮምፒውተሮችና ላፕቶፖች', got %q", amCrumbs[1].Name)
		}
	}
}

func TestCategoryHierarchyRepo_GetBreadcrumbs_RootCategory(t *testing.T) {
	db, cleanup := setupTestDB(t)
	if db == nil {
		return
	}
	defer cleanup()

	repo := NewCategoryHierarchyRepo(db)
	ctx := context.Background()

	// Category 1 is root "Electronics"
	crumbs, err := repo.GetBreadcrumbs(ctx, 1, "en")
	if err != nil {
		t.Fatalf("GetBreadcrumbs failed for root category: %v", err)
	}

	if len(crumbs) != 1 {
		t.Fatalf("expected exactly 1 breadcrumb for root category, got %d", len(crumbs))
	}

	if crumbs[0].ID != 1 {
		t.Errorf("expected ID 1, got %d", crumbs[0].ID)
	}
	if crumbs[0].Name != "Electronics" {
		t.Errorf("expected name 'Electronics', got %q", crumbs[0].Name)
	}
	if crumbs[0].Slug != "electronics" {
		t.Errorf("expected slug 'electronics', got %q", crumbs[0].Slug)
	}
	if crumbs[0].Level != 0 {
		t.Errorf("expected level 0, got %d", crumbs[0].Level)
	}
}

func TestCategoryHierarchyRepo_GetBreadcrumbs_NonexistentCategory(t *testing.T) {
	db, cleanup := setupTestDB(t)
	if db == nil {
		return
	}
	defer cleanup()

	repo := NewCategoryHierarchyRepo(db)
	ctx := context.Background()

	// Nonexistent category ID 999999
	crumbs, err := repo.GetBreadcrumbs(ctx, 999999, "en")
	if err != nil {
		t.Fatalf("GetBreadcrumbs returned unexpected error for nonexistent category: %v", err)
	}

	if len(crumbs) != 0 {
		t.Errorf("expected 0 breadcrumbs for nonexistent category, got %d", len(crumbs))
	}
}

func TestCategoryHierarchyRepo_GetBreadcrumbs_RowOrdering(t *testing.T) {
	db, cleanup := setupTestDB(t)
	if db == nil {
		return
	}
	defer cleanup()

	repo := NewCategoryHierarchyRepo(db)
	ctx := context.Background()

	// Test ordering for a child category
	crumbs, err := repo.GetBreadcrumbs(ctx, 14, "en")
	if err != nil {
		t.Fatalf("GetBreadcrumbs failed: %v", err)
	}

	for i := 0; i < len(crumbs)-1; i++ {
		if crumbs[i].Level >= crumbs[i+1].Level {
			t.Errorf("ancestor ordering violated at index %d: level %d >= level %d",
				i, crumbs[i].Level, crumbs[i+1].Level)
		}
	}
}

func TestCategoryHierarchyRepo_GetCategoryTree_RecursionAndAlphabeticalSiblingOrdering(t *testing.T) {
	db, cleanup := setupTestDB(t)
	if db == nil {
		return
	}
	defer cleanup()

	repo := NewCategoryHierarchyRepo(db)
	ctx := context.Background()

	tree, err := repo.GetCategoryTree(ctx, "en")
	if err != nil {
		t.Fatalf("GetCategoryTree failed: %v", err)
	}

	if len(tree) == 0 {
		t.Fatalf("expected non-empty category tree, got 0 roots")
	}

	// Verify roots are alphabetically sorted among themselves
	for i := 0; i < len(tree)-1; i++ {
		nameA := tree[i].Name
		nameB := tree[i+1].Name
		if strings.ToLower(strings.TrimSpace(nameA)) > strings.ToLower(strings.TrimSpace(nameB)) {
			t.Errorf("Root categories not sorted alphabetically: %q > %q at index %d", nameA, nameB, i)
		}
	}

	totalNodes := 0
	maxDepth := 0
	nodesAtLevel := make(map[int]int)

	var walkAndVerifySiblings func(nodes []models.CategoryTree, parentID *int, depth int)
	walkAndVerifySiblings = func(nodes []models.CategoryTree, parentID *int, depth int) {
		if depth > maxDepth {
			maxDepth = depth
		}
		for i := range nodes {
			totalNodes++
			nodesAtLevel[nodes[i].Level]++

			// Verify parent relationship
			if parentID != nil && nodes[i].ParentID != nil && *nodes[i].ParentID != *parentID {
				t.Errorf("Node %d has parent_id %d, expected %d", nodes[i].ID, *nodes[i].ParentID, *parentID)
			}

			// Verify alphabetical sorting within siblings
			if i < len(nodes)-1 {
				nameA := strings.ToLower(strings.TrimSpace(nodes[i].Name))
				nameB := strings.ToLower(strings.TrimSpace(nodes[i+1].Name))
				if nameA > nameB {
					t.Errorf("Sibling ordering violated under parent %v: %q (ID %d) > %q (ID %d)",
						parentID, nodes[i].Name, nodes[i].ID, nodes[i+1].Name, nodes[i+1].ID)
				}
			}

			// Recurse into children
			if len(nodes[i].Children) > 0 {
				nodeID := nodes[i].ID
				walkAndVerifySiblings(nodes[i].Children, &nodeID, depth+1)
			}
		}
	}

	walkAndVerifySiblings(tree, nil, 0)

	t.Logf("Tree verification: %d total categories, max depth %d, levels: %v", totalNodes, maxDepth, nodesAtLevel)

	// In the database there are 2,380 categories across levels 0, 1, 2, 3
	if totalNodes < 2000 {
		t.Errorf("Expected full tree to retain all active categories (>= 2000), got %d (depth truncation detected)", totalNodes)
	}

	// Verify deep levels exist (levels 2 and 3 must not be truncated)
	if nodesAtLevel[2] == 0 {
		t.Errorf("Expected level 2 categories in tree, got 0")
	}
	if nodesAtLevel[3] == 0 {
		t.Errorf("Expected level 3 categories in tree, got 0")
	}
}

func TestCategoryHierarchyRepo_GetCategoryTree_AmharicAlphabeticalOrdering(t *testing.T) {
	db, cleanup := setupTestDB(t)
	if db == nil {
		return
	}
	defer cleanup()

	repo := NewCategoryHierarchyRepo(db)
	ctx := context.Background()

	tree, err := repo.GetCategoryTree(ctx, "am")
	if err != nil {
		t.Fatalf("GetCategoryTree in Amharic failed: %v", err)
	}

	// Verify roots are sorted by Amharic display name
	for i := 0; i < len(tree)-1; i++ {
		nameA := strings.ToLower(strings.TrimSpace(tree[i].Name))
		nameB := strings.ToLower(strings.TrimSpace(tree[i+1].Name))
		if nameA > nameB {
			t.Errorf("Amharic root ordering violated at index %d: %q (ID %d) > %q (ID %d)",
				i, tree[i].Name, tree[i].ID, tree[i+1].Name, tree[i+1].ID)
		}
	}
}

func TestCategoryHierarchyRepo_GetRoots_AlphabeticalOrdering(t *testing.T) {
	db, cleanup := setupTestDB(t)
	if db == nil {
		return
	}
	defer cleanup()

	repo := NewCategoryHierarchyRepo(db)
	ctx := context.Background()

	rootsEn, err := repo.GetRoots(ctx, "en")
	if err != nil {
		t.Fatalf("GetRoots in English failed: %v", err)
	}
	for i := 0; i < len(rootsEn)-1; i++ {
		nameA := strings.ToLower(strings.TrimSpace(rootsEn[i].NameEn))
		nameB := strings.ToLower(strings.TrimSpace(rootsEn[i+1].NameEn))
		if nameA > nameB {
			t.Errorf("English root ordering violated: %q > %q", nameA, nameB)
		}
	}

	rootsAm, err := repo.GetRoots(ctx, "am")
	if err != nil {
		t.Fatalf("GetRoots in Amharic failed: %v", err)
	}
	for i := 0; i < len(rootsAm)-1; i++ {
		nameA := rootsAm[i].NameAm
		if nameA == "" {
			nameA = rootsAm[i].NameEn
		}
		nameB := rootsAm[i+1].NameAm
		if nameB == "" {
			nameB = rootsAm[i+1].NameEn
		}
		if strings.ToLower(strings.TrimSpace(nameA)) > strings.ToLower(strings.TrimSpace(nameB)) {
			t.Errorf("Amharic root ordering violated: %q > %q", nameA, nameB)
		}
	}
}


