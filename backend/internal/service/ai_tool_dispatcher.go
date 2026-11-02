package service

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"sync"
	"time"

	"github.com/your-org/used-tools-marketplace/internal/models"
	"github.com/your-org/used-tools-marketplace/internal/repository"
	"github.com/your-org/used-tools-marketplace/internal/service/aiclient"
	"github.com/your-org/used-tools-marketplace/pkg/dto"
)

// AIToolDispatcher executes authorized, typed marketplace tools on behalf of the AI.
// The LLM has zero direct SQL access; all requests pass through validated repository methods.
type AIToolDispatcher struct {
	productRepo *repository.ProductRepo
}

// NewAIToolDispatcher creates a new tool dispatcher instance.
func NewAIToolDispatcher(productRepo *repository.ProductRepo) *AIToolDispatcher {
	return &AIToolDispatcher{
		productRepo: productRepo,
	}
}

// ListingSummaryDTO represents a clean, grounded marketplace listing returned to the LLM.
type ListingSummaryDTO struct {
	ID                  int      `json:"id"`
	Title               string   `json:"title"`
	Price               float64  `json:"price"`
	Currency            string   `json:"currency"`
	ListingType         string   `json:"listing_type"`                    // 'sale' or 'rent'
	RentalPrice         *float64 `json:"rental_price,omitempty"`          // periodic rental price
	RentalDurationValue *int     `json:"rental_duration_value,omitempty"` // duration count
	RentalDurationUnit  *string  `json:"rental_duration_unit,omitempty"`  // month, etc.
	Condition           string   `json:"condition"`
	SellerTier          string   `json:"seller_tier"`
	City                string   `json:"city"`
	Subcity             string   `json:"subcity,omitempty"`
	DetailURL           string   `json:"detail_url"`
	ImageURL            *string  `json:"image_url,omitempty"`
	ContactAction       string   `json:"contact_action"`      // 'Buy / Contact Seller' or 'Rent / Contact Landlord'
	AvailabilityStatus  string   `json:"availability_status"` // 'Available' or 'Currently Rented'
}

func buildListingSummary(p *models.Product, cityName, subcityName string, imgURL *string) ListingSummaryDTO {
	contactAction := "Buy / Contact Seller"
	if p.ListingType == "rent" {
		contactAction = "Rent / Contact Landlord"
	}
	avail := "Available"
	if p.IsCurrentlyRented {
		avail = "Currently Rented"
	}

	return ListingSummaryDTO{
		ID:                  p.ID,
		Title:               p.Title,
		Price:               p.Price,
		Currency:            "ETB",
		ListingType:         p.ListingType,
		RentalPrice:         p.RentalPrice,
		RentalDurationValue: p.RentalDurationValue,
		RentalDurationUnit:  p.RentalDurationUnit,
		Condition:           p.Condition,
		SellerTier:          p.SellerTier,
		City:                cityName,
		Subcity:             subcityName,
		DetailURL:           fmt.Sprintf("/products/%d", p.ID),
		ImageURL:            imgURL,
		ContactAction:       contactAction,
		AvailabilityStatus:  avail,
	}
}

// MultiCategorySearchResultDTO holds isolated results for one category in a multi-category search.
type MultiCategorySearchResultDTO struct {
	CategoryName string              `json:"category_name"`
	CategorySlug string              `json:"category_slug,omitempty"`
	SearchTerm   string              `json:"search_term"`
	TotalCount   int                 `json:"total_count"`
	Listings     []ListingSummaryDTO `json:"listings"`
	Error        string              `json:"error,omitempty"`
}

// SearchToolResult holds the grounded search execution result.
type SearchToolResult struct {
	TotalMatching int                 `json:"total_matching"`
	ListingsCount int                 `json:"listings_count"`
	Listings      []ListingSummaryDTO `json:"listings"`
	SearchScope   string              `json:"search_scope"`
}

// MarketplaceStatsResult holds authoritative price statistics across live listings.
type MarketplaceStatsResult struct {
	SearchTerm     string  `json:"search_term"`
	TotalListings  int     `json:"total_listings"`
	MinPrice       float64 `json:"min_price"`
	AvgPrice       float64 `json:"avg_price"`
	MaxPrice       float64 `json:"max_price"`
	Currency       string  `json:"currency"`
	LocationFilter string  `json:"location_filter,omitempty"`
}

// ExecuteToolCall validates, authorizes, and runs the requested tool.
func (d *AIToolDispatcher) ExecuteToolCall(ctx context.Context, call aiclient.ChatToolCall) (string, []ListingSummaryDTO, error) {
	switch call.Function.Name {
	case "search_marketplace_listings":
		return d.executeSearchListings(ctx, call.Function.Arguments)
	case "get_listing_details":
		return d.executeGetListingDetails(ctx, call.Function.Arguments)
	case "compare_listings":
		return d.executeCompareListings(ctx, call.Function.Arguments)
	case "get_marketplace_stats":
		return d.executeMarketplaceStats(ctx, call.Function.Arguments)
	default:
		return fmt.Sprintf(`{"error": "unsupported tool '%s'"}`, call.Function.Name), nil, nil
	}
}

func (d *AIToolDispatcher) executeSearchListings(ctx context.Context, args map[string]interface{}) (string, []ListingSummaryDTO, error) {
	searchTerm, _ := args["search_term"].(string)
	filter := dto.ProductFilter{
		Page:  1,
		Limit: 5,
	}

	cleanTerm := aiclient.CleanConversationalSearchTerm(searchTerm)
	if cleanTerm != "" && cleanTerm != "all" && cleanTerm != "products" && cleanTerm != "product" && cleanTerm != "item" && cleanTerm != "items" {
		filter.Search = &cleanTerm
	}

	if minP, ok := args["min_price"].(float64); ok && minP > 0 {
		filter.MinPrice = &minP
	}
	if maxP, ok := args["max_price"].(float64); ok && maxP > 0 {
		filter.MaxPrice = &maxP
	}
	if cond, ok := args["condition"].(string); ok && (cond == "new" || cond == "used") {
		filter.Condition = &cond
	}
	if tier, ok := args["seller_tier"].(string); ok && tier != "" {
		filter.SellerTier = &tier
	}
	if lType, ok := args["listing_type"].(string); ok && (lType == "sale" || lType == "rent") {
		filter.ListingType = &lType
	}
	if loc, ok := args["location"].(string); ok && loc != "" {
		// Use as region or subcity filter
		filter.Subcity = &loc
	}
	if lim, ok := args["limit"].(float64); ok && lim > 0 {
		filter.Limit = int(math.Min(lim, 10))
	}

	totalMatching, err := d.productRepo.CountApproved(ctx, filter)
	if err != nil {
		totalMatching = 0
	}

	products, err := d.productRepo.ListApproved(ctx, filter)
	if err != nil {
		return fmt.Sprintf(`{"error": "failed to search listings: %s"}`, err.Error()), nil, nil
	}

	// Graceful relaxation: if over-constrained filters (price, condition, location) return 0 results,
	// search by keyword alone so the assistant can present available market options
	if len(products) == 0 && (filter.MaxPrice != nil || filter.MinPrice != nil || filter.Condition != nil || filter.Subcity != nil) {
		relaxedFilter := dto.ProductFilter{
			Search: filter.Search,
			Page:   1,
			Limit:  filter.Limit,
		}
		if relaxedProducts, rErr := d.productRepo.ListApproved(ctx, relaxedFilter); rErr == nil && len(relaxedProducts) > 0 {
			products = relaxedProducts
			if rCount, cErr := d.productRepo.CountApproved(ctx, relaxedFilter); cErr == nil {
				totalMatching = rCount
			}
		}
	}

	summaries := make([]ListingSummaryDTO, 0, len(products))
	for _, p := range products {
		cityName := "Addis Ababa"
		subcityName := ""
		if p.Location != nil {
			if p.Location.Region != "" {
				cityName = p.Location.Region
			}
			if p.Location.Subcity != nil {
				subcityName = *p.Location.Subcity
			}
		}

		var imgURL *string
		if len(p.Images) > 0 && p.Images[0].ImageURL != "" {
			url := p.Images[0].ImageURL
			imgURL = &url
		}

		summaries = append(summaries, buildListingSummary(&p, cityName, subcityName, imgURL))
	}

	res := SearchToolResult{
		TotalMatching: totalMatching,
		ListingsCount: len(summaries),
		Listings:      summaries,
		SearchScope:   fmt.Sprintf("Live approved listings for '%s'", searchTerm),
	}

	resBytes, err := json.Marshal(res)
	if err != nil {
		return `{"error": "serialization error"}`, nil, nil
	}

	return string(resBytes), summaries, nil
}

func (d *AIToolDispatcher) executeGetListingDetails(ctx context.Context, args map[string]interface{}) (string, []ListingSummaryDTO, error) {
	productIDFloat, ok := args["product_id"].(float64)
	if !ok || productIDFloat <= 0 {
		return `{"error": "invalid product_id argument"}`, nil, nil
	}
	productID := int(productIDFloat)

	product, err := d.productRepo.GetByID(ctx, productID)
	if err != nil || product == nil {
		return fmt.Sprintf(`{"error": "product ID %d not found or inactive"}`, productID), nil, nil
	}

	cityName := "Addis Ababa"
	subcityName := ""
	if product.Location != nil {
		if product.Location.Region != "" {
			cityName = product.Location.Region
		}
		if product.Location.Subcity != nil {
			subcityName = *product.Location.Subcity
		}
	}

	var imgURL *string
	if len(product.Images) > 0 && product.Images[0].ImageURL != "" {
		url := product.Images[0].ImageURL
		imgURL = &url
	}

	summary := buildListingSummary(product, cityName, subcityName, imgURL)

	payload := map[string]interface{}{
		"id":                    product.ID,
		"title":                 product.Title,
		"price":                 product.Price,
		"currency":              "ETB",
		"listing_type":          product.ListingType,
		"rental_price":          product.RentalPrice,
		"rental_duration_value": product.RentalDurationValue,
		"rental_duration_unit":  product.RentalDurationUnit,
		"deposit_amount":        product.DepositAmount,
		"available_from":        product.AvailableFrom,
		"is_currently_rented":   product.IsCurrentlyRented,
		"contact_action":        summary.ContactAction,
		"availability_status":   summary.AvailabilityStatus,
		"condition":             product.Condition,
		"seller_tier":           product.SellerTier,
		"description":           product.Description,
		"status":                product.Status,
		"city":                  cityName,
		"subcity":               subcityName,
		"attributes":            product.Attributes,
		"created_at":            product.CreatedAt,
	}

	resBytes, _ := json.Marshal(payload)
	return string(resBytes), []ListingSummaryDTO{summary}, nil
}

func (d *AIToolDispatcher) executeCompareListings(ctx context.Context, args map[string]interface{}) (string, []ListingSummaryDTO, error) {
	rawIDs, ok := args["product_ids"].([]interface{})
	if !ok || len(rawIDs) < 2 {
		return `{"error": "compare_listings requires an array of 2 or more product_ids"}`, nil, nil
	}

	var summaries []ListingSummaryDTO
	var detailedListings []map[string]interface{}

	for _, raw := range rawIDs {
		idFloat, ok := raw.(float64)
		if !ok {
			continue
		}
		pID := int(idFloat)
		prod, err := d.productRepo.GetByID(ctx, pID)
		if err != nil || prod == nil {
			continue
		}

		cityName := "Addis Ababa"
		if prod.Location != nil && prod.Location.Region != "" {
			cityName = prod.Location.Region
		}

		var imgURL *string
		if len(prod.Images) > 0 && prod.Images[0].ImageURL != "" {
			url := prod.Images[0].ImageURL
			imgURL = &url
		}

		s := buildListingSummary(prod, cityName, "", imgURL)
		summaries = append(summaries, s)

		detailedListings = append(detailedListings, map[string]interface{}{
			"id":                  prod.ID,
			"title":               prod.Title,
			"price":               prod.Price,
			"listing_type":        prod.ListingType,
			"rental_price":        prod.RentalPrice,
			"condition":           prod.Condition,
			"seller_tier":         prod.SellerTier,
			"city":                cityName,
			"contact_action":      s.ContactAction,
			"availability_status": s.AvailabilityStatus,
			"description":         prod.Description,
		})
	}

	resBytes, _ := json.Marshal(map[string]interface{}{
		"comparison_count": len(detailedListings),
		"items":            detailedListings,
	})

	return string(resBytes), summaries, nil
}

func (d *AIToolDispatcher) executeMarketplaceStats(ctx context.Context, args map[string]interface{}) (string, []ListingSummaryDTO, error) {
	searchTerm, _ := args["search_term"].(string)
	loc, _ := args["location"].(string)

	filter := dto.ProductFilter{
		Page:  1,
		Limit: 50,
	}
	if searchTerm != "" {
		filter.Search = &searchTerm
	}
	if loc != "" {
		filter.Subcity = &loc
	}

	products, err := d.productRepo.ListApproved(ctx, filter)
	if err != nil || len(products) == 0 {
		return fmt.Sprintf(`{"search_term": "%s", "total_listings": 0, "message": "No active listings found"}`, searchTerm), nil, nil
	}

	minP := products[0].Price
	maxP := products[0].Price
	sumP := 0.0

	for _, p := range products {
		if p.Price < minP {
			minP = p.Price
		}
		if p.Price > maxP {
			maxP = p.Price
		}
		sumP += p.Price
	}

	avgP := sumP / float64(len(products))

	stats := MarketplaceStatsResult{
		SearchTerm:     searchTerm,
		TotalListings:  len(products),
		MinPrice:       minP,
		AvgPrice:       math.Round(avgP),
		MaxPrice:       maxP,
		Currency:       "ETB",
		LocationFilter: strings.TrimSpace(loc),
	}

	resBytes, _ := json.Marshal(stats)
	return string(resBytes), nil, nil
}

// ExecuteMultiSearch performs bounded, isolated parallel searches across multiple category requests.
// Each search has an independent timeout and failure in one does not affect others.
func (d *AIToolDispatcher) ExecuteMultiSearch(ctx context.Context, searches []ExtractedSearchRequest) []MultiCategorySearchResultDTO {
	if len(searches) == 0 {
		return nil
	}

	results := make([]MultiCategorySearchResultDTO, len(searches))
	var wg sync.WaitGroup

	for i, req := range searches {
		wg.Add(1)
		go func(idx int, sReq ExtractedSearchRequest) {
			defer wg.Done()

			searchCtx, cancel := context.WithTimeout(ctx, 1500*time.Millisecond)
			defer cancel()

			res := MultiCategorySearchResultDTO{
				CategoryName: sReq.CategoryName,
				CategorySlug: sReq.CategorySlug,
				SearchTerm:   sReq.SearchTerm,
				Listings:     make([]ListingSummaryDTO, 0),
			}

			filter := dto.ProductFilter{
				Page:  1,
				Limit: 4, // Max 4 per category to avoid clutter
			}

			cleanTerm := aiclient.CleanConversationalSearchTerm(sReq.SearchTerm)
			if sReq.Brand != "" {
				if cleanTerm == "" || cleanTerm == "all" || cleanTerm == "products" || cleanTerm == "items" {
					cleanTerm = sReq.Brand
				} else if !strings.Contains(strings.ToLower(cleanTerm), strings.ToLower(sReq.Brand)) {
					cleanTerm = sReq.Brand + " " + cleanTerm
				}
			}
			if cleanTerm != "" && cleanTerm != "all" && cleanTerm != "products" && cleanTerm != "items" {
				filter.Search = &cleanTerm
			}

			if sReq.MinPrice != nil && *sReq.MinPrice > 0 {
				filter.MinPrice = sReq.MinPrice
			}
			if sReq.MaxPrice != nil && *sReq.MaxPrice > 0 {
				filter.MaxPrice = sReq.MaxPrice
			}
			if sReq.ListingType != nil && (*sReq.ListingType == "sale" || *sReq.ListingType == "rent") {
				filter.ListingType = sReq.ListingType
			}
			if sReq.Location != "" {
				loc := sReq.Location
				filter.Subcity = &loc
			}
			if sReq.Condition != "" {
				cond := sReq.Condition
				filter.Condition = &cond
			}

			total, err := d.productRepo.CountApproved(searchCtx, filter)
			if err != nil {
				total = 0
			}
			res.TotalCount = total

			products, err := d.productRepo.ListApproved(searchCtx, filter)
			if err != nil {
				res.Error = err.Error()
				results[idx] = res
				return
			}

			// Relaxation if 0 results and extra filters applied
			if len(products) == 0 && (filter.MaxPrice != nil || filter.MinPrice != nil || filter.Condition != nil || filter.Subcity != nil) {
				relaxedFilter := dto.ProductFilter{
					Search: filter.Search,
					Page:   1,
					Limit:  filter.Limit,
				}
				if relaxedProducts, rErr := d.productRepo.ListApproved(searchCtx, relaxedFilter); rErr == nil && len(relaxedProducts) > 0 {
					products = relaxedProducts
					if rCount, cErr := d.productRepo.CountApproved(searchCtx, relaxedFilter); cErr == nil {
						res.TotalCount = rCount
					}
				}
			}

			summaries := make([]ListingSummaryDTO, 0, len(products))
			for _, p := range products {
				cityName := "Addis Ababa"
				subcityName := ""
				if p.Location != nil {
					if p.Location.Region != "" {
						cityName = p.Location.Region
					}
					if p.Location.Subcity != nil {
						subcityName = *p.Location.Subcity
					}
				}

				var imgURL *string
				if len(p.Images) > 0 && p.Images[0].ImageURL != "" {
					url := p.Images[0].ImageURL
					imgURL = &url
				}

				summaries = append(summaries, buildListingSummary(&p, cityName, subcityName, imgURL))
			}

			res.Listings = summaries
			results[idx] = res
		}(i, req)
	}

	wg.Wait()
	return results
}

