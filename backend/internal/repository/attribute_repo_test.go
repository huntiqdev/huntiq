package repository

import (
	"context"
	"testing"
)

func TestAttributeRepo_GetCategoryAttributes_NullableFields(t *testing.T) {
	db, cleanup := setupTestDB(t)
	if db == nil {
		return
	}
	defer cleanup()

	repo := NewAttributeRepo(db)
	ctx := context.Background()

	// Category 14 has attributes with both null and non-null units/validation rules
	attrs, err := repo.GetCategoryAttributes(ctx, 14)
	if err != nil {
		t.Fatalf("GetCategoryAttributes failed: %v", err)
	}

	if len(attrs) == 0 {
		t.Fatalf("expected attributes for category 14, got 0")
	}

	// Verify required attributes exist
	hasBrand := false
	hasScreenSize := false
	for _, attr := range attrs {
		if attr.Slug == "brand" {
			hasBrand = true
			if !attr.IsRequired {
				t.Errorf("expected brand to be required")
			}
		}
		if attr.Slug == "screen-size" {
			hasScreenSize = true
			if attr.Unit != "inches" {
				t.Errorf("expected unit 'inches' for screen-size, got %q", attr.Unit)
			}
		}
	}

	if !hasBrand {
		t.Errorf("expected brand attribute in category 14")
	}
	if !hasScreenSize {
		t.Errorf("expected screen-size attribute in category 14")
	}
}

func TestAttributeRepo_GetCategoryAttributes_NonexistentCategory(t *testing.T) {
	db, cleanup := setupTestDB(t)
	if db == nil {
		return
	}
	defer cleanup()

	repo := NewAttributeRepo(db)
	ctx := context.Background()

	attrs, err := repo.GetCategoryAttributes(ctx, 999999)
	if err != nil {
		t.Fatalf("unexpected error for nonexistent category: %v", err)
	}

	if len(attrs) != 0 {
		t.Errorf("expected 0 attributes for nonexistent category, got %d", len(attrs))
	}
}

func TestAttributeRepo_GetAttributeDefinitionAndBySlug(t *testing.T) {
	db, cleanup := setupTestDB(t)
	if db == nil {
		return
	}
	defer cleanup()

	repo := NewAttributeRepo(db)
	ctx := context.Background()

	// Test by slug (with null unit: e.g. "storage" or "processor")
	attrBySlug, err := repo.GetAttributeBySlug(ctx, "processor")
	if err != nil {
		t.Fatalf("GetAttributeBySlug failed: %v", err)
	}
	if attrBySlug.Slug != "processor" {
		t.Errorf("expected slug 'processor', got %q", attrBySlug.Slug)
	}

	// Test by ID
	attrByID, err := repo.GetAttributeDefinition(ctx, attrBySlug.ID)
	if err != nil {
		t.Fatalf("GetAttributeDefinition failed: %v", err)
	}
	if attrByID.ID != attrBySlug.ID {
		t.Errorf("expected ID %d, got %d", attrBySlug.ID, attrByID.ID)
	}
}
