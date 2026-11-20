package ai

import (
	"context"
	"testing"
	"time"

	"github.com/your-org/used-tools-marketplace/internal/models"
	"github.com/your-org/used-tools-marketplace/pkg/dto"
)

type mockProductSearcher struct {
	products []models.Product
}

func (m *mockProductSearcher) ListApproved(ctx context.Context, f dto.ProductFilter) ([]models.Product, error) {
	return m.products, nil
}

func (m *mockProductSearcher) GetByID(ctx context.Context, id int) (*models.Product, error) {
	for _, p := range m.products {
		if p.ID == id {
			return &p, nil
		}
	}
	return nil, nil
}

func TestExtractSearchFilters(t *testing.T) {
	tests := []struct {
		input       string
		expectedMax *float64
		expectedMin *float64
	}{
		{
			input: "Find laptops under 50000 ETB",
			expectedMax: func() *float64 {
				v := 50000.0
				return &v
			}(),
			expectedMin: nil,
		},
		{
			input: "Apartment for rent above 25k birr",
			expectedMax: nil,
			expectedMin: func() *float64 {
				v := 25000.0
				return &v
			}(),
		},
	}

	for _, tc := range tests {
		_, minP, maxP := ExtractSearchFilters(tc.input)
		if tc.expectedMax != nil {
			if maxP == nil || *maxP != *tc.expectedMax {
				t.Errorf("input '%s': expected max price %v, got %v", tc.input, *tc.expectedMax, maxP)
			}
		}
		if tc.expectedMin != nil {
			if minP == nil || *minP != *tc.expectedMin {
				t.Errorf("input '%s': expected min price %v, got %v", tc.input, *tc.expectedMin, minP)
			}
		}
	}
}

func TestToolRegistry_MapProductToListingCard(t *testing.T) {
	deletedTime := time.Now()
	product := models.Product{
		ID:          101,
		Title:       "MacBook Pro M2 16GB",
		Price:       95000,
		Condition:   "used_like_new",
		ListingType: "sale",
		Category: &models.Category{
			ID:   5,
			Name: "Laptops & Computers",
		},
		Location: &models.Location{
			Region: "Addis Ababa",
			Subcity: func() *string {
				s := "Bole"
				return &s
			}(),
		},
		Images: []models.ProductImage{
			{
				ID:        1,
				ImageURL:  "/uploads/products/macbook-main.jpg",
				IsBlurred: false,
			},
			{
				ID:        2,
				ImageURL:  "/uploads/products/macbook-blur.jpg",
				IsBlurred: true, // Should be ignored
			},
			{
				ID:        3,
				ImageURL:  "/uploads/products/macbook-deleted.jpg",
				DeletedAt: &deletedTime, // Should be ignored
			},
			{
				ID:        4,
				ImageURL:  "/uploads/products/macbook-side.jpg",
				IsBlurred: false,
			},
		},
	}

	reg := NewToolRegistry(nil, nil, nil, "http://localhost:8080")
	card := reg.MapProductToListingCard(product)

	if card.ID != 101 {
		t.Errorf("expected ID 101, got %d", card.ID)
	}
	if card.Title != "MacBook Pro M2 16GB" {
		t.Errorf("expected title match, got %s", card.Title)
	}
	if card.Price != 95000 {
		t.Errorf("expected price 95000, got %f", card.Price)
	}
	if card.CategoryName != "Laptops & Computers" {
		t.Errorf("expected category 'Laptops & Computers', got '%s'", card.CategoryName)
	}
	if card.LocationName != "Addis Ababa, Bole" {
		t.Errorf("expected location 'Addis Ababa, Bole', got '%s'", card.LocationName)
	}
	if card.PrimaryImage != "/uploads/products/macbook-main.jpg" {
		t.Errorf("expected primary image '/uploads/products/macbook-main.jpg', got '%s'", card.PrimaryImage)
	}
	// Verify filtered images: only 2 valid images (excluding blurred and deleted)
	if len(card.ImageURLs) != 2 {
		t.Errorf("expected 2 valid images, got %d", len(card.ImageURLs))
	}
}

func TestToolRegistry_SearchPublicListings(t *testing.T) {
	mockP := &mockProductSearcher{
		products: []models.Product{
			{
				ID:          201,
				Title:       "Dell XPS 15",
				Price:       60000,
				Status:      "approved",
				ListingType: "sale",
			},
		},
	}

	reg := NewToolRegistry(mockP, nil, nil, "http://localhost:8080")
	cards, err := reg.SearchPublicListings(context.Background(), "Dell", nil, nil, 5)
	if err != nil {
		t.Fatalf("unexpected search error: %v", err)
	}
	if len(cards) != 1 {
		t.Fatalf("expected 1 result, got %d", len(cards))
	}
	if cards[0].ID != 201 {
		t.Errorf("expected ID 201, got %d", cards[0].ID)
	}
	if cards[0].PrimaryImage != "/placeholder-image.png" {
		t.Errorf("expected fallback placeholder, got %s", cards[0].PrimaryImage)
	}
}
