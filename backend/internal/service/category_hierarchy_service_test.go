package service

import (
	"context"
	"testing"

	"github.com/your-org/used-tools-marketplace/internal/models"
)

func TestGenerateSlug(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Electronics", "electronics"},
		{"Mobile Phones & Tablets", "mobile-phones-tablets"},
		{"  Heavy Machinery / Equipment  ", "heavy-machinery-equipment"},
		{"Commercial & Residential Property 2024", "commercial-residential-property-2024"},
		{"Multiple---Dashes   And   Spaces", "multiple-dashes-and-spaces"},
		{"Special!@#$%^&*()_+Characters", "specialcharacters"},
		{"Ethiopian Birr / ዋጋ", "ethiopian-birr"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := generateSlug(tt.input)
			if got != tt.expected {
				t.Errorf("generateSlug(%q) = %q; want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestIsValidSlug(t *testing.T) {
	tests := []struct {
		slug  string
		valid bool
	}{
		{"electronics", true},
		{"mobile-phones", true},
		{"heavy-machinery-equipment", true},
		{"category-123", true},
		{"", false},
		{"-invalid-leading", false},
		{"invalid-trailing-", false},
		{"invalid--double-dash", false},
		{"Invalid-Caps", false},
		{"invalid spaces", false},
		{"invalid@symbol", false},
	}

	for _, tt := range tests {
		t.Run(tt.slug, func(t *testing.T) {
			got := isValidSlug(tt.slug)
			if got != tt.valid {
				t.Errorf("isValidSlug(%q) = %v; want %v", tt.slug, got, tt.valid)
			}
		})
	}
}

func TestCategorySelfParentValidation(t *testing.T) {
	svc := &CategoryHierarchyService{}
	ctx := context.Background()

	parentID := 42
	cat := &models.CategoryHierarchy{
		ID:       42,
		NameEn:   "Electronics",
		Slug:     "electronics",
		ParentID: &parentID,
	}

	err := svc.UpdateCategory(ctx, cat)
	if err == nil {
		t.Fatalf("expected error when category is its own parent, got nil")
	}
	expectedErrMsg := "category cannot be its own parent"
	if err.Error() != expectedErrMsg {
		t.Errorf("expected error %q, got %q", expectedErrMsg, err.Error())
	}
}

func TestCategoryValidationRequiredFields(t *testing.T) {
	svc := &CategoryHierarchyService{}
	ctx := context.Background()

	// Missing ID
	err := svc.UpdateCategory(ctx, &models.CategoryHierarchy{NameEn: "No ID"})
	if err == nil {
		t.Errorf("expected error for missing ID, got nil")
	}

	// Missing NameEn
	err = svc.UpdateCategory(ctx, &models.CategoryHierarchy{ID: 1, NameEn: ""})
	if err == nil {
		t.Errorf("expected error for missing name_en, got nil")
	}

	// Invalid slug with disallowed characters
	err = svc.UpdateCategory(ctx, &models.CategoryHierarchy{ID: 1, NameEn: "Test", Slug: "---"})
	if err == nil {
		t.Errorf("expected error for invalid slug, got nil")
	}
}

func TestValidateCategoryForListing_Empty(t *testing.T) {
	svc := &CategoryHierarchyService{}
	ctx := context.Background()

	err := svc.ValidateCategoryForListing(ctx, nil, nil)
	if err == nil {
		t.Fatalf("expected error for nil category IDs, got nil")
	}

	zero := 0
	err = svc.ValidateCategoryForListing(ctx, &zero, &zero)
	if err == nil {
		t.Fatalf("expected error for zero category IDs, got nil")
	}
}

