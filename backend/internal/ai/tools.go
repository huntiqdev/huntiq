package ai

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/your-org/used-tools-marketplace/internal/models"
	"github.com/your-org/used-tools-marketplace/pkg/dto"
)

// ProductSearchProvider defines the interface required to retrieve real listings.
type ProductSearchProvider interface {
	ListApproved(ctx context.Context, f dto.ProductFilter) ([]models.Product, error)
	GetByID(ctx context.Context, id int) (*models.Product, error)
}

// CategorySearchProvider defines the interface required to query category taxonomy.
type CategorySearchProvider interface {
	SearchCategories(ctx context.Context, term string, lang string, limit int) ([]models.CategorySearchResult, error)
}

// LocationSearchProvider defines the interface required to query locations.
type LocationSearchProvider interface {
	SearchLocations(ctx context.Context, q string, lang string, limit int) ([]models.LocationHierarchy, error)
}

// ToolRegistry manages authorized, narrowly scoped tools for HUNTIQ AI.
type ToolRegistry struct {
	products   ProductSearchProvider
	categories CategorySearchProvider
	locations  LocationSearchProvider
	appBaseURL string
}

func NewToolRegistry(p ProductSearchProvider, c CategorySearchProvider, l LocationSearchProvider, appBaseURL string) *ToolRegistry {
	return &ToolRegistry{
		products:   p,
		categories: c,
		locations:  l,
		appBaseURL: appBaseURL,
	}
}

// MapProductToListingCard converts a canonical database models.Product to an authorized ListingCardDTO.
func (t *ToolRegistry) MapProductToListingCard(p models.Product) ListingCardDTO {
	card := ListingCardDTO{
		ID:          p.ID,
		Title:       p.Title,
		Price:       p.Price,
		Currency:    "ETB",
		Condition:   p.Condition,
		ListingType: p.ListingType,
		URL:         fmt.Sprintf("/products/%d", p.ID),
		CreatedAt:   p.CreatedAt,
	}

	if p.Category != nil {
		card.CategoryName = p.Category.Name
	}
	if p.Location != nil {
		parts := []string{p.Location.Region}
		if p.Location.Subcity != nil && *p.Location.Subcity != "" {
			parts = append(parts, *p.Location.Subcity)
		} else if p.Location.Zone != nil && *p.Location.Zone != "" {
			parts = append(parts, *p.Location.Zone)
		}
		card.LocationName = strings.Join(parts, ", ")
	}

	// Extract real authorized images
	var images []string
	for _, img := range p.Images {
		if img.DeletedAt == nil && !img.IsBlurred && img.ImageURL != "" {
			fullURL := img.ImageURL
			if !strings.HasPrefix(fullURL, "http://") && !strings.HasPrefix(fullURL, "https://") {
				if !strings.HasPrefix(fullURL, "/") {
					fullURL = "/" + fullURL
				}
			}
			images = append(images, fullURL)
		}
	}

	if len(images) > 0 {
		card.PrimaryImage = images[0]
		card.ImageURLs = images
	} else {
		// Neutral fallback placeholder if no image exists
		card.PrimaryImage = "/placeholder-image.png"
		card.ImageURLs = []string{"/placeholder-image.png"}
	}

	return card
}

// SearchPublicListings executes an authorized marketplace listing search.
func (t *ToolRegistry) SearchPublicListings(ctx context.Context, query string, minPrice, maxPrice *float64, limit int) ([]ListingCardDTO, error) {
	if limit <= 0 || limit > 20 {
		limit = 10
	}

	f := dto.ProductFilter{
		Limit: limit,
		Page:  1,
	}

	q := strings.TrimSpace(query)
	if q != "" {
		f.Search = &q
	}
	if minPrice != nil && *minPrice > 0 {
		f.MinPrice = minPrice
	}
	if maxPrice != nil && *maxPrice > 0 {
		f.MaxPrice = maxPrice
	}

	sort := "trending"
	f.SortBy = &sort

	if t.products == nil {
		return nil, nil
	}

	products, err := t.products.ListApproved(ctx, f)
	if err != nil {
		return nil, fmt.Errorf("search public listings: %w", err)
	}

	cards := make([]ListingCardDTO, 0, len(products))
	for _, p := range products {
		cards = append(cards, t.MapProductToListingCard(p))
	}

	return cards, nil
}

// GetPublicListingDetails retrieves verified details for a single approved listing.
func (t *ToolRegistry) GetPublicListingDetails(ctx context.Context, id int) (*ListingCardDTO, error) {
	if t.products == nil || id <= 0 {
		return nil, nil
	}

	p, err := t.products.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if p == nil || p.Status != "approved" {
		return nil, fmt.Errorf("listing #%d is not available", id)
	}

	card := t.MapProductToListingCard(*p)
	return &card, nil
}

// ExtractSearchFilters extracts structured search intent (search term, max price, min price) from natural text.
func ExtractSearchFilters(text string) (string, *float64, *float64) {
	clean := strings.ToLower(text)

	var minPrice, maxPrice *float64

	// Regex for "under X" or "below X" or "less than X"
	reUnder := regexp.MustCompile(`(?i)(?:under|below|less\s+than|maximum|max)\s+([0-9]+(?:,[0-9]+)*(?:\.[0-9]+)?)\s*(?:etb|birr|k)?`)
	if match := reUnder.FindStringSubmatch(clean); len(match) > 1 {
		numStr := strings.ReplaceAll(match[1], ",", "")
		if val, err := strconv.ParseFloat(numStr, 64); err == nil {
			if strings.HasSuffix(strings.ToLower(match[0]), "k") {
				val *= 1000
			}
			maxPrice = &val
		}
	}

	// Regex for "above X" or "more than X" or "minimum|min X"
	reAbove := regexp.MustCompile(`(?i)(?:above|over|more\s+than|minimum|min)\s+([0-9]+(?:,[0-9]+)*(?:\.[0-9]+)?)\s*(?:etb|birr|k)?`)
	if match := reAbove.FindStringSubmatch(clean); len(match) > 1 {
		numStr := strings.ReplaceAll(match[1], ",", "")
		if val, err := strconv.ParseFloat(numStr, 64); err == nil {
			if strings.HasSuffix(strings.ToLower(match[0]), "k") {
				val *= 1000
			}
			minPrice = &val
		}
	}

	// Clean out conversational filler words to isolate entity search keywords
	fillers := []string{
		"find me", "find", "show me", "show", "search for", "search", "looking for", "i want",
		"i need", "can you show", "available", "please", "around", "near", "in", "for sale",
		"for rent", "buy", "sell", "cheap", "good", "best", "እፈልጋለሁ", "አሳየኝ", "ፈልግ",
	}

	searchTerm := clean
	for _, f := range fillers {
		searchTerm = strings.ReplaceAll(searchTerm, f, " ")
	}
	// Strip extracted price clause
	if maxPrice != nil {
		searchTerm = reUnder.ReplaceAllString(searchTerm, " ")
	}
	if minPrice != nil {
		searchTerm = reAbove.ReplaceAllString(searchTerm, " ")
	}

	searchTerm = strings.TrimSpace(regexp.MustCompile(`\s+`).ReplaceAllString(searchTerm, " "))

	return searchTerm, minPrice, maxPrice
}
