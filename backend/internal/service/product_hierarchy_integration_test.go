package service

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/your-org/used-tools-marketplace/internal/config"
	"github.com/your-org/used-tools-marketplace/internal/database"
	"github.com/your-org/used-tools-marketplace/internal/repository"
	"github.com/your-org/used-tools-marketplace/pkg/dto"
)

func setupProductTestDB(t *testing.T) (*pgxpool.Pool, func()) {
	if os.Getenv("DB_PASSWORD") == "" && os.Getenv("DATABASE_URL") == "" {
		t.Skip("Skipping database-dependent test: DB_PASSWORD or DATABASE_URL not set in environment")
		return nil, nil
	}

	cfg := config.NewTestConfig()
	db, err := database.NewPool(cfg)
	if err != nil {
		t.Skipf("Skipping database-dependent test: DB connection failed: %v", err)
		return nil, nil
	}

	cleanup := func() {
		db.Close()
	}
	return db, cleanup
}

// TestProductLifecycle_HierarchicalLeaf163_DualMapping tests that creating a listing
// with hierarchical leaf category ID 163 (Refrigerators and Freezers, parent 44: Home Appliances)
// automatically maps category_id to 4 (Home & Living legacy root) while preserving
// category_hierarchy_id as 163, satisfying foreign key integrity on both tables.
func TestProductLifecycle_HierarchicalLeaf163_DualMapping(t *testing.T) {
	db, cleanup := setupProductTestDB(t)
	if db == nil {
		return
	}
	defer cleanup()

	ctx := context.Background()

	// Initialize repositories
	productRepo := repository.NewProductRepo(db)
	catRepo := repository.NewCategoryRepo(db)
	catHierarchyRepo := repository.NewCategoryHierarchyRepo(db)
	locationRepo := repository.NewLocationRepo(db)
	userRepo := repository.NewUserRepo(db)

	productSvc := NewProductService(productRepo, catRepo, locationRepo, userRepo).
		WithCategoryHierarchyRepo(catHierarchyRepo)

	// Verify Category 163 exists in database
	valid, active, err := catHierarchyRepo.ValidateCategory(ctx, 163)
	if err != nil || !valid || !active {
		t.Fatalf("Category 163 expected to exist and be active: valid=%v, active=%v, err=%v", valid, active, err)
	}

	// Find active seller
	var sellerID int
	err = db.QueryRow(ctx, "SELECT id FROM users WHERE role = 'seller' AND is_active = true LIMIT 1").Scan(&sellerID)
	if err != nil {
		t.Fatalf("Failed to find active seller in test database: %v", err)
	}

	leafCatID := 163
	createReq := dto.CreateProductRequest{
		Title:               "Samsung Double Door No-Frost Refrigerator Test",
		Price:               55000.00,
		Description:         "Brand new energy efficient refrigerator with multi air flow system.",
		CategoryHierarchyID: &leafCatID,
		ListingType:         "sale",
		Region:              "Addis Ababa",
		ImageURLs:           []string{"/uploads/fridge1.jpg", "/uploads/fridge2.jpg"},
	}

	// 1. Create product
	product, err := productSvc.CreateForSeller(ctx, sellerID, createReq)
	if err != nil {
		t.Fatalf("CreateForSeller failed: %v", err)
	}
	defer func() {
		// Clean up test listing from database
		_, _ = db.Exec(context.Background(), "DELETE FROM products WHERE id = $1", product.ID)
	}()

	// 2. Verify dual category resolution
	if product.CategoryID == nil {
		t.Fatal("expected product.CategoryID to be non-nil")
	}
	if *product.CategoryID != 4 {
		t.Errorf("expected product.CategoryID to map to legacy root 4 (Home & Living), got %d", *product.CategoryID)
	}
	if product.CategoryHierarchyID == nil {
		t.Fatal("expected product.CategoryHierarchyID to be non-nil")
	}
	if *product.CategoryHierarchyID != 163 {
		t.Errorf("expected product.CategoryHierarchyID to retain leaf 163, got %d", *product.CategoryHierarchyID)
	}

	// 3. Verify in PostgreSQL directly
	var dbCatID, dbCatHierarchyID int
	err = db.QueryRow(ctx, "SELECT category_id, category_hierarchy_id FROM products WHERE id = $1", product.ID).
		Scan(&dbCatID, &dbCatHierarchyID)
	if err != nil {
		t.Fatalf("Failed to query created product from database: %v", err)
	}
	if dbCatID != 4 {
		t.Errorf("database category_id: expected 4, got %d", dbCatID)
	}
	if dbCatHierarchyID != 163 {
		t.Errorf("database category_hierarchy_id: expected 163, got %d", dbCatHierarchyID)
	}

	// 4. Update product
	updatedDesc := "Updated refrigerator description with extended warranty."
	updatedPrice := 52000.00
	updateReq := dto.UpdateProductRequest{
		Title:               &product.Title,
		Price:               &updatedPrice,
		Description:         &updatedDesc,
		CategoryHierarchyID: &leafCatID,
	}

	err = productSvc.UpdateForSeller(ctx, product.ID, sellerID, updateReq)
	if err != nil {
		t.Fatalf("UpdateForSeller failed: %v", err)
	}

	// 5. Verify update preserved category integrity
	err = db.QueryRow(ctx, "SELECT price, category_id, category_hierarchy_id FROM products WHERE id = $1", product.ID).
		Scan(&updatedPrice, &dbCatID, &dbCatHierarchyID)
	if err != nil {
		t.Fatalf("Failed to query updated product from database: %v", err)
	}
	if dbCatID != 4 || dbCatHierarchyID != 163 {
		t.Errorf("Update corrupted categories: category_id=%d, category_hierarchy_id=%d", dbCatID, dbCatHierarchyID)
	}

	t.Logf("Successfully verified end-to-end listing create/read/update with leaf category 163 -> legacy root 4 in PostgreSQL")
}
