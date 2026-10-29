package handlers

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/your-org/used-tools-marketplace/internal/cache"
	"github.com/your-org/used-tools-marketplace/internal/service"
	"github.com/your-org/used-tools-marketplace/internal/service/aiclient"
	"github.com/your-org/used-tools-marketplace/pkg/dto"
	"github.com/your-org/used-tools-marketplace/pkg/response"
)

// timeNow is a package-level variable so tests can override it.
var timeNow = time.Now

// SearchLogger abstracts search analytics recording for testability and loose coupling.
type SearchLogger interface {
	LogSearch(ctx context.Context, query string, resultCount int, categoryID *int, filters map[string]interface{}, userID *int, sessionID, ipAddress string) error
}

type searchLogEvent struct {
	query       string
	resultCount int
	categoryID  *int
	filters     map[string]interface{}
	userID      *int
	sessionID   string
	ipAddress   string
}

type ProductHandler struct {
	svc          *service.ProductService
	authSvc      *service.AuthService
	aiClient     *aiclient.Client
	cache        cache.Cache
	searchChan   chan searchLogEvent
	logger       SearchLogger
	droppedCount uint64
}

func NewProductHandler(svc *service.ProductService, authSvc *service.AuthService) *ProductHandler {
	return &ProductHandler{svc: svc, authSvc: authSvc}
}

// DroppedSearchEvents returns the number of search events dropped due to queue saturation.
func (h *ProductHandler) DroppedSearchEvents() uint64 {
	return atomic.LoadUint64(&h.droppedCount)
}

// WithAIClient attaches the AI service client for intelligent query understanding.
func (h *ProductHandler) WithAIClient(client *aiclient.Client) *ProductHandler {
	h.aiClient = client
	return h
}

// WithCache attaches the Cache instance for AI query and search acceleration.
func (h *ProductHandler) WithCache(c cache.Cache) *ProductHandler {
	h.cache = c
	return h
}

// WithSearchLogger wires a bounded non-blocking worker for search query recording.
func (h *ProductHandler) WithSearchLogger(logger SearchLogger) *ProductHandler {
	if logger == nil {
		return h
	}
	h.logger = logger
	h.searchChan = make(chan searchLogEvent, 500)
	go func() {
		for ev := range h.searchChan {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			if err := h.logger.LogSearch(ctx, ev.query, ev.resultCount, ev.categoryID, ev.filters, ev.userID, ev.sessionID, ev.ipAddress); err != nil {
				log.Printf("search analytics logging failed: %v", err)
			}
			cancel()
		}
	}()
	return h
}

// isTierAuthorized verifies if the requester's buyer tier satisfies the product's required buyer tier.
func isTierAuthorized(userTier, requiredTier string, isAdmin bool) bool {
	if isAdmin {
		return true
	}
	if requiredTier == "" || requiredTier == "normal" {
		return true
	}
	tierRank := map[string]int{
		"guest":   0,
		"":        0,
		"normal":  1,
		"gold":    2,
		"premium": 3,
	}
	return tierRank[userTier] >= tierRank[requiredTier]
}

// GET /api/products
func (h *ProductHandler) List(c *fiber.Ctx) error {
	f := dto.ProductFilter{}

	if v := c.Query("category_id"); v != "" {
		id, err := strconv.Atoi(v)
		if err != nil {
			return response.BadRequest(c, "invalid category_id")
		}
		f.CategoryID = &id
	}
	if v := c.Query("category_hierarchy_id"); v != "" {
		id, err := strconv.Atoi(v)
		if err != nil {
			return response.BadRequest(c, "invalid category_hierarchy_id")
		}
		f.CategoryHierarchyID = &id
	}
	f.IncludeDescendants = true
	if v := c.Query("include_descendants"); v == "false" || v == "0" {
		f.IncludeDescendants = false
	}
	if v := c.Query("min_price"); v != "" {
		p, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return response.BadRequest(c, "invalid min_price")
		}
		f.MinPrice = &p
	}
	if v := c.Query("max_price"); v != "" {
		p, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return response.BadRequest(c, "invalid max_price")
		}
		f.MaxPrice = &p
	}
	if v := c.Query("location_id"); v != "" {
		id, err := strconv.Atoi(v)
		if err != nil {
			return response.BadRequest(c, "invalid location_id")
		}
		f.LocationID = &id
	}
	if v := c.Query("location_hierarchy_id"); v != "" {
		id, err := strconv.Atoi(v)
		if err != nil {
			return response.BadRequest(c, "invalid location_hierarchy_id")
		}
		f.LocationHierarchyID = &id
	}
	// Location hierarchy params — most specific takes precedence
	if v := c.Query("woreda"); v != "" {
		f.Woreda = &v
	} else if v := c.Query("subcity"); v != "" {
		f.Subcity = &v
	} else if v := c.Query("zone"); v != "" {
		f.Zone = &v
	} else if v := c.Query("region"); v != "" {
		f.Region = &v
	}
	if v := c.Query("search"); v != "" {
		v = truncate(v, 200)
		f.Search = &v
	}
	if v := c.Query("condition"); v != "" {
		switch v {
		case "new", "used":
			f.Condition = &v
		}
	}
	if v := c.Query("subcategory"); v != "" {
		v = truncate(v, 100)
		f.Subcategory = &v
	}
	if v := c.Query("item"); v != "" {
		v = truncate(v, 100)
		f.Item = &v
	}
	// Listing type filter: 'sale' or 'rent'
	if v := c.Query("listing_type"); v != "" {
		if v == "sale" || v == "rent" {
			f.ListingType = &v
		}
	}
	if v := c.Query("sort_by"); v != "" {
		switch v {
		case "price_asc", "price_desc", "newest":
			f.SortBy = &v
		}
	} else if v := c.Query("sort"); v != "" {
		switch v {
		case "price_asc", "price_desc", "newest":
			f.SortBy = &v
		}
	}
	f.Page = parseQueryInt(c, "page", 1)
	f.Limit = parseQueryInt(c, "limit", 20)
	// Hard cap — never return more than 100 products per request
	if f.Limit > 100 {
		f.Limit = 100
	}
	if f.Page < 1 {
		f.Page = 1
	}

	var searchMeta fiber.Map
	searchMode := strings.ToLower(strings.TrimSpace(c.Query("mode", "manual")))
	locale := c.Query("locale", "")
	lang := c.Query("lang", "")
	if locale == "" && lang != "" {
		locale = lang
	}
	prefLang := c.Query("preferred_language", "")
	if prefLang == "" {
		prefLang = locale
	}

	// Query understanding via internal AI worker (only if searchMode == "ai" and configured)
	if searchMode == "ai" && f.Search != nil && strings.TrimSpace(*f.Search) != "" && h.aiClient != nil && h.aiClient.IsEnabled() {
		var aiResp *aiclient.QueryUnderstandingResponse
		cleanQueryRaw := strings.ToLower(strings.TrimSpace(*f.Search))
		cacheKey := fmt.Sprintf("huntiq:ai:q:v1:%s:%s:%x", locale, prefLang, sha256.Sum256([]byte(cleanQueryRaw)))

		// 1. Fast Cache Lookup (< 1ms)
		if h.cache != nil {
			if cachedBytes, found := h.cache.Get(c.Context(), cacheKey); found {
				var cached aiclient.QueryUnderstandingResponse
				if err := json.Unmarshal(cachedBytes, &cached); err == nil {
					aiResp = &cached
				}
			}
		}

		// 2. Cache miss: Invocate Python AI worker and store for 24 hours
		if aiResp == nil {
			params := aiclient.UnderstandQueryParams{
				Query:             *f.Search,
				Locale:            locale,
				Mode:              "ai",
				PreferredLanguage: prefLang,
			}
			if resp, err := h.aiClient.UnderstandQueryWithParams(c.Context(), params); err == nil && resp != nil {
				aiResp = resp
				if h.cache != nil {
					if marshaled, err := json.Marshal(aiResp); err == nil {
						_ = h.cache.Set(c.Context(), cacheKey, marshaled, 24*time.Hour)
					}
				}
			}
		}

		if aiResp != nil {
			searchMeta = fiber.Map{
				"ai_mode":             true,
				"normalized_query":    aiResp.NormalizedQuery,
				"cleaned_search_term": aiResp.CleanedSearchTerm,
				"sort_by":             aiResp.SortBy,
				"primary_intent":      aiResp.PrimaryIntent,
				"localized_answer":    aiResp.LocalizedAnswer,
				"detected_locale":     aiResp.DetectedLocale,
				"deep_parameters":     aiResp.DeepParameters,
				"pii_detected":        aiResp.PIIDetected,
				"is_safe":             aiResp.IsSafe,
				"safety_flags":         aiResp.SafetyFlags,
				"reasoning_chain":      aiResp.ReasoningChain,
			}

			// Clean search term from conversational filler and noise for database tsquery
			if aiResp.CleanedSearchTerm != "" {
				cleanTerm := aiResp.CleanedSearchTerm
				f.Search = &cleanTerm
			} else if aiResp.PIIDetected && aiResp.NormalizedQuery != "" {
				cleanQ := aiResp.NormalizedQuery
				f.Search = &cleanQ
			}

			// Sort by extraction from AI signals
			if f.SortBy == nil && aiResp.SortBy != nil && *aiResp.SortBy != "" {
				f.SortBy = aiResp.SortBy
			}

			// Intent-based listing type inference (rent vs sale)
			if f.ListingType == nil {
				if aiResp.PrimaryIntent == "rent_property" {
					lt := "rent"
					f.ListingType = &lt
				} else if aiResp.PrimaryIntent == "buy_product" || aiResp.PrimaryIntent == "vehicle_search" {
					lt := "sale"
					f.ListingType = &lt
				}
			}
			// Extract price boundaries if not explicitly provided by client
			if f.MaxPrice == nil && aiResp.ExtractedEntities.PriceRange.MaxPrice != nil {
				f.MaxPrice = aiResp.ExtractedEntities.PriceRange.MaxPrice
			}
			if f.MinPrice == nil && aiResp.ExtractedEntities.PriceRange.MinPrice != nil {
				f.MinPrice = aiResp.ExtractedEntities.PriceRange.MinPrice
			}
			// Extract condition if not explicitly provided
			if f.Condition == nil && aiResp.ExtractedEntities.Condition != nil {
				f.Condition = aiResp.ExtractedEntities.Condition
			}
			// Extract location if not explicitly provided
			if f.LocationID == nil && f.LocationHierarchyID == nil && f.Region == nil && f.Subcity == nil && aiResp.ExtractedEntities.Location != nil {
				locStr := strings.TrimSpace(*aiResp.ExtractedEntities.Location)
				if locStr != "" {
					if strings.EqualFold(locStr, "Addis Ababa") || strings.EqualFold(locStr, "Oromia") || strings.EqualFold(locStr, "Amhara") || strings.EqualFold(locStr, "Tigray") {
						f.Region = &locStr
					} else {
						f.Subcity = &locStr
					}
				}
			}
			// Extract seller tier filter if detected by AI
			if f.SellerTier == nil && aiResp.ExtractedEntities.Attributes != nil {
				if stVal, ok := aiResp.ExtractedEntities.Attributes["seller_tier"].(string); ok && stVal != "" {
					f.SellerTier = &stVal
				}
			}

			// Route to suggested top category if client did not select a category
			if f.CategoryID == nil && f.CategoryHierarchyID == nil && len(aiResp.SuggestedCategoryIDs) > 0 {
				catID := aiResp.SuggestedCategoryIDs[0]
				f.CategoryHierarchyID = &catID
			}
		} else {
			// Deterministic Go AI Fallback
			rawSearch := *f.Search
			hasIntent, fallbackArgs := aiclient.ParseFallbackSearchArgs(nil, rawSearch)
			cleanTerm := aiclient.CleanConversationalSearchTerm(rawSearch)
			if cleanTerm != "" && cleanTerm != "all" && cleanTerm != "products" && cleanTerm != "product" && cleanTerm != "items" {
				f.Search = &cleanTerm
			} else {
				f.Search = nil
			}

			if hasIntent {
				if f.Condition == nil {
					if cond, ok := fallbackArgs["condition"].(string); ok && cond != "" {
						f.Condition = &cond
					}
				}
				if f.ListingType == nil {
					if lt, ok := fallbackArgs["listing_type"].(string); ok && lt != "" {
						f.ListingType = &lt
					}
				}
				if f.MaxPrice == nil {
					if mp, ok := fallbackArgs["max_price"].(float64); ok && mp > 0 {
						f.MaxPrice = &mp
					}
				}
				if f.MinPrice == nil {
					if minP, ok := fallbackArgs["min_price"].(float64); ok && minP > 0 {
						f.MinPrice = &minP
					}
				}
				if f.LocationID == nil && f.LocationHierarchyID == nil && f.Region == nil && f.Subcity == nil {
					if loc, ok := fallbackArgs["location"].(string); ok && loc != "" {
						if strings.EqualFold(loc, "Addis Ababa") || strings.EqualFold(loc, "Oromia") || strings.EqualFold(loc, "Amhara") || strings.EqualFold(loc, "Tigray") {
							f.Region = &loc
						} else {
							f.Subcity = &loc
						}
					}
				}
			}

			searchMeta = fiber.Map{
				"ai_mode":             true,
				"normalized_query":    cleanTerm,
				"cleaned_search_term": cleanTerm,
				"primary_intent":      "search_marketplace",
			}
		}
	} else if searchMode == "manual" {
		searchMeta = fiber.Map{
			"ai_mode": false,
		}
		if f.Search != nil && *f.Search != "" {
			cleanTerm := aiclient.CleanConversationalSearchTerm(*f.Search)
			if cleanTerm != "" && cleanTerm != "all" && cleanTerm != "products" && cleanTerm != "product" {
				f.Search = &cleanTerm
			}
		}
	}

	// Query explicit seller_tier param if passed by client
	if f.SellerTier == nil {
		if st := c.Query("seller_tier", ""); st != "" {
			f.SellerTier = &st
		}
	}

	totalMatching, _ := h.svc.CountApproved(c.Context(), f)
	products, err := h.svc.ListApproved(c.Context(), f)
	if err != nil {
		return response.InternalError(c, err)
	}

	if searchMeta == nil {
		searchMeta = fiber.Map{}
	}
	searchMeta["total_matching"] = totalMatching
	searchMeta["returned_count"] = len(products)
	searchMeta["page"] = f.Page
	searchMeta["limit"] = f.Limit

	// For count_listings intent, provide accurate grounded count
	if intentStr, ok := searchMeta["primary_intent"].(string); ok && intentStr == "count_listings" {
		itemLabel := ""
		if cl, ok := searchMeta["cleaned_search_term"].(string); ok && cl != "" {
			itemLabel = cl
		} else if f.Search != nil {
			itemLabel = *f.Search
		}
		if itemLabel == "" {
			itemLabel = "items"
		}
		switch prefLang {
		case "am":
			searchMeta["localized_answer"] = fmt.Sprintf("በዳታቤዝ ውስጥ በአሁኑ ሰዓት %d ንቁ የ'%s' ዝርዝር ምርቶች ይገኛሉ።", totalMatching, itemLabel)
		case "om":
			searchMeta["localized_answer"] = fmt.Sprintf("Kuusaa gabaa keessatti yeroo ammaa oomishaaleen '%s' %d ta'an argamaniiru.", itemLabel, totalMatching)
		case "ti":
			searchMeta["localized_answer"] = fmt.Sprintf("ኣብ ዳታቤዝ ዕዳጋ ኣብዚ እዋን %d ንጡፋት ናይ '%s' ዝርዝራት ኣለዉ።", totalMatching, itemLabel)
		default:
			searchMeta["localized_answer"] = fmt.Sprintf("Found %d active '%s' listings matching your criteria in the marketplace.", totalMatching, itemLabel)
		}
	}

	// Resolve requester tier & admin status from context
	userTier, _ := c.Locals("buyerTier").(string)
	role, _ := c.Locals("role").(string)
	adminRole, _ := c.Locals("adminRole").(string)
	isAdmin := role == "admin" || adminRole != ""

	for i := range products {
		if !isTierAuthorized(userTier, products[i].RequiredBuyerTier, isAdmin) {
			products[i].SpecificationsLocked = true
			products[i].Attributes = nil // Strip restricted specifications from search/list results
		} else {
			products[i].SpecificationsLocked = false
		}
	}

	// Asynchronously record search analytics on the first page of search results
	if h.searchChan != nil && f.Search != nil && strings.TrimSpace(*f.Search) != "" && f.Page == 1 {
		var uid *int
		if uVal := c.Locals("userID"); uVal != nil {
			if id, ok := uVal.(int); ok && id > 0 {
				uid = &id
			}
		}
		filters := map[string]interface{}{}
		if f.CategoryID != nil {
			filters["category_id"] = *f.CategoryID
		}
		if f.CategoryHierarchyID != nil {
			filters["category_hierarchy_id"] = *f.CategoryHierarchyID
		}
		if f.Condition != nil {
			filters["condition"] = *f.Condition
		}
		if f.ListingType != nil {
			filters["listing_type"] = *f.ListingType
		}
		ev := searchLogEvent{
			query:       strings.TrimSpace(*f.Search),
			resultCount: len(products),
			categoryID:  f.CategoryID,
			filters:     filters,
			userID:      uid,
			sessionID:   c.Get("X-Session-ID"),
			ipAddress:   c.IP(),
		}
		select {
		case h.searchChan <- ev:
		default:
			// Queue full under burst traffic; drop non-blocking to protect response latency
			atomic.AddUint64(&h.droppedCount, 1)
		}
	}
	// Real-Time Deep Product Market Analysis on actual database listings
	if searchMeta != nil && len(products) > 0 && h.aiClient != nil {
		candidates := make([]aiclient.ProductCandidateSummary, 0, len(products))
		for _, p := range products {
			var locStr *string
			if p.Location != nil {
				l := p.Location.Region
				if p.Location.Subcity != nil && *p.Location.Subcity != "" {
					l = *p.Location.Subcity + ", " + l
				}
				locStr = &l
			}
			candidates = append(candidates, aiclient.ProductCandidateSummary{
				ID:          p.ID,
				Title:       p.Title,
				Price:       p.Price,
				SellerTier:  p.SellerTier,
				Condition:   p.Condition,
				Location:    locStr,
				ListingType: p.ListingType,
			})
		}
		qStr := ""
		if f.Search != nil {
			qStr = *f.Search
		}
		analysisReq := aiclient.ProductAnalysisRequest{
			Query:             qStr,
			Locale:            locale,
			PreferredLanguage: prefLang,
			TotalMatching:     totalMatching,
			Products:          candidates,
		}
		if analysisResp, err := h.aiClient.AnalyzeProducts(c.Context(), analysisReq); err == nil && analysisResp != nil {
			searchMeta["real_time_analysis"] = analysisResp
			if analysisResp.MarketVerdict != "" {
				searchMeta["market_verdict"] = analysisResp.MarketVerdict
			}
		} else {
			// Instantaneous in-process deterministic analysis fallback (< 0.01ms)
			detResp := aiclient.ComputeDeterministicAnalysis(analysisReq)
			searchMeta["real_time_analysis"] = detResp
			if detResp.MarketVerdict != "" {
				searchMeta["market_verdict"] = detResp.MarketVerdict
			}
		}
	}

	return response.OKListWithMeta(c, products, searchMeta)
}

// GET /api/products/:id  — public, only approved and non-expired products visible
func (h *ProductHandler) Get(c *fiber.Ctx) error {
	id, err := parseIDParam(c, "id")
	if err != nil {
		return err
	}
	product, err := h.svc.GetByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return response.NotFound(c, "product not found")
		}
		return response.InternalError(c, err)
	}
	if product.Status != "approved" || !product.SellerActive {
		return response.NotFound(c, "product not found")
	}
	// Also hide products whose listing has expired but the sweep job hasn't run yet
	if product.ExpiresAt != nil && product.ExpiresAt.Before(timeNow()) {
		return response.NotFound(c, "product not found")
	}

	// Resolve requester tier & admin status from context
	userTier, _ := c.Locals("buyerTier").(string)
	role, _ := c.Locals("role").(string)
	adminRole, _ := c.Locals("adminRole").(string)
	isAdmin := role == "admin" || adminRole != ""

	// Server-enforced specification locking: unauthorized users receive clean public data without restricted specs
	if !isTierAuthorized(userTier, product.RequiredBuyerTier, isAdmin) {
		product.SpecificationsLocked = true
		product.Attributes = nil // Strip restricted specifications from response payload
	} else {
		product.SpecificationsLocked = false
	}

	return response.OK(c, product)
}

// GET /api/categories
func (h *ProductHandler) ListCategories(c *fiber.Ctx) error {
	cats, err := h.svc.ListCategories(c.Context())
	if err != nil {
		return response.InternalError(c, err)
	}
	return response.OKList(c, cats)
}

// GET /api/categories/popular?days=7&limit=8
// Returns categories ranked by activity (new listings + orders) in the last N days.
func (h *ProductHandler) PopularCategories(c *fiber.Ctx) error {
	days := parseQueryInt(c, "days", 7)
	limit := parseQueryInt(c, "limit", 8)
	if days < 1 || days > 90 {
		days = 7
	}
	if limit < 1 || limit > 20 {
		limit = 8
	}
	popular, err := h.svc.PopularCategories(c.Context(), days, limit)
	if err != nil {
		return response.InternalError(c, err)
	}
	return response.OKList(c, popular)
}

// GET /api/seller/fee-preview?price=X&duration_days=Y
func (h *ProductHandler) FeePreview(c *fiber.Ctx) error {
	sellerID, err := mustUserID(c)
	if err != nil {
		return err
	}
	price, err := strconv.ParseFloat(c.Query("price", "0"), 64)
	if err != nil || price <= 0 {
		return response.BadRequest(c, "invalid price")
	}
	duration := parseQueryInt(c, "duration_days", 30)
	if duration <= 0 {
		duration = 30
	}
	breakdown, err := h.svc.PreviewListingFee(c.Context(), sellerID, price, duration)
	if err != nil {
		return response.InternalError(c, err)
	}
	return response.OK(c, breakdown)
}

// GET /api/fee-options — public, returns all duration options with fee breakdown for a given price+tier
func (h *ProductHandler) FeeOptions(c *fiber.Ctx) error {
	price, err := strconv.ParseFloat(c.Query("price", "0"), 64)
	if err != nil || price <= 0 {
		return response.BadRequest(c, "invalid price")
	}
	tierStr := c.Query("tier", "normal")
	validTiers := map[string]bool{"normal": true, "gold": true, "premium": true}
	if !validTiers[tierStr] {
		tierStr = "normal"
	}
	tier := service.SellerTier(tierStr)
	type option struct {
		DurationDays int    `json:"duration_days"`
		Label        string `json:"label"`
		service.FeeBreakdown
	}
	labels := map[int]string{2: "2 days", 7: "1 week", 30: "1 month", 180: "6 months", 365: "1 year"}
	breakdowns := h.svc.FeeOptions(c.Context(), tier, price)
	var opts []option
	for _, b := range breakdowns {
		opts = append(opts, option{DurationDays: b.DurationDays, Label: labels[b.DurationDays], FeeBreakdown: b})
	}
	return response.OKList(c, opts)
}

// POST /api/seller/products
func (h *ProductHandler) SellerCreate(c *fiber.Ctx) error {
	sellerID, err := mustUserID(c)
	if err != nil {
		return err
	}
	var req dto.CreateProductRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	product, err := h.svc.CreateForSeller(c.Context(), sellerID, req)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}
	if h.cache != nil {
		_ = h.cache.Delete(c.Context(), fmt.Sprintf("cache:product:%d", product.ID))
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"success": true, "data": product})
}

// PUT /api/seller/products/:id
func (h *ProductHandler) SellerUpdate(c *fiber.Ctx) error {
	sellerID, err := mustUserID(c)
	if err != nil {
		return err
	}
	id, err := parseIDParam(c, "id")
	if err != nil {
		return err
	}
	var req dto.UpdateProductRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	if err := h.svc.UpdateForSeller(c.Context(), id, sellerID, req); err != nil {
		return response.BadRequest(c, err.Error())
	}
	if h.cache != nil {
		_ = h.cache.Delete(c.Context(), fmt.Sprintf("cache:product:%d", id))
	}
	return response.OK(c, fiber.Map{"message": "product updated"})
}

// GET /api/seller/products
func (h *ProductHandler) SellerList(c *fiber.Ctx) error {
	sellerID, err := mustUserID(c)
	if err != nil {
		return err
	}
	products, err := h.svc.ListForSeller(c.Context(), sellerID)
	if err != nil {
		return response.InternalError(c, err)
	}
	return response.OKList(c, products)
}

