package repository

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/your-org/used-tools-marketplace/internal/config"
	"github.com/your-org/used-tools-marketplace/internal/database"
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
