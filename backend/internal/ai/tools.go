package ai

import (
	"context"
	"fmt"
	"strings"

	"github.com/your-org/used-tools-marketplace/internal/repository"
	"github.com/your-org/used-tools-marketplace/pkg/dto"
)

type ProductSearchTool struct {
	productRepo *repository.ProductRepo
}

func NewProductSearchTool(productRepo *repository.ProductRepo) *ProductSearchTool {
	return &ProductSearchTool{productRepo: productRepo}
}

func (t *ProductSearchTool) SearchListings(ctx context.Context, term string, minPrice, maxPrice *float64, limit int) ([]ListingCard, error) {
	if t.productRepo == nil {
		return nil, nil
	}

	if limit <= 0 || limit > 10 {
		limit = 6
	}

	filter := dto.ProductFilter{
		Limit:    limit,
		Page:     1,
		MinPrice: minPrice,
		MaxPrice: maxPrice,
	}

	cleanTerm := strings.TrimSpace(term)
	if cleanTerm != "" {
		filter.Search = &cleanTerm
	}

	products, err := t.productRepo.ListApproved(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list approved products: %w", err)
	}

	// Resilient multi-token fallback: if multi-word term (e.g. "samsung phone") yields 0 rows,
	// retry with key tokens (e.g. "samsung") within the active price filters.
	if len(products) == 0 && strings.Contains(cleanTerm, " ") {
		words := strings.Fields(cleanTerm)
		for _, w := range words {
			wClean := strings.TrimSpace(w)
			if len(wClean) >= 3 && wClean != "phone" && wClean != "item" && wClean != "device" {
				filter.Search = &wClean
				fallbackProducts, errFallback := t.productRepo.ListApproved(ctx, filter)
				if errFallback == nil && len(fallbackProducts) > 0 {
					products = fallbackProducts
					break
				}
			}
		}
	}

	cards := make([]ListingCard, 0, len(products))
	for _, p := range products {
		catName := ""
		if p.Category != nil {
			catName = p.Category.Name
		} else if p.Subcategory != "" {
			catName = p.Subcategory
		}

		locName := ""
		if p.Location != nil {
			parts := []string{}
			if p.Location.Subcity != nil && *p.Location.Subcity != "" {
				parts = append(parts, *p.Location.Subcity)
			}
			if p.Location.Region != "" {
				parts = append(parts, p.Location.Region)
			}
			locName = strings.Join(parts, ", ")
		}

		imgURL := ""
		if len(p.Images) > 0 {
			imgURL = p.Images[0].ImageURL
		}

		listingType := p.ListingType
		if listingType == "" {
			listingType = "sale"
		}

		cards = append(cards, ListingCard{
			ID:           p.ID,
			Title:        p.Title,
			Price:        p.Price,
			Currency:     "ETB",
			CategoryName: catName,
			LocationName: locName,
			Condition:    p.Condition,
			ListingType:  listingType,
			ImageURL:     imgURL,
			ViewCount:    p.ViewCount,
		})
	}

	return cards, nil
}

type TrendTool struct {
	productRepo *repository.ProductRepo
	searchRepo  *repository.SearchRepo
}

func NewTrendTool(productRepo *repository.ProductRepo, searchRepo *repository.SearchRepo) *TrendTool {
	return &TrendTool{
		productRepo: productRepo,
		searchRepo:  searchRepo,
	}
}

func (t *TrendTool) GetTrendSummary(ctx context.Context) string {
	var sb strings.Builder

	if t.productRepo != nil {
		popCats, err := t.productRepo.PopularCategories(ctx, 30, 4)
		if err == nil && len(popCats) > 0 {
			names := make([]string, 0, len(popCats))
			for _, pc := range popCats {
				names = append(names, fmt.Sprintf("%s (%d listings)", pc.Name, pc.NewListings))
			}
			sb.WriteString("High-demand categories: " + strings.Join(names, ", ") + ". ")
		}
	}

	if t.searchRepo != nil {
		popSearches, err := t.searchRepo.GetPopularSearches(ctx, 5)
		if err == nil && len(popSearches) > 0 {
			terms := make([]string, 0, len(popSearches))
			for _, ps := range popSearches {
				if q, ok := ps["query"].(string); ok && q != "" {
					terms = append(terms, q)
				}
			}
			if len(terms) > 0 {
				sb.WriteString("Top user searches: " + strings.Join(terms, ", ") + ".")
			}
		}
	}

	return strings.TrimSpace(sb.String())
}
