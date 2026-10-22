package handlers

import (
	"errors"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/your-org/used-tools-marketplace/pkg/dto"
	"github.com/your-org/used-tools-marketplace/internal/service"
	"github.com/your-org/used-tools-marketplace/pkg/response"
)

// timeNow is a package-level variable so tests can override it.
var timeNow = time.Now

type ProductHandler struct {
	svc     *service.ProductService
	authSvc *service.AuthService
}

func NewProductHandler(svc *service.ProductService, authSvc *service.AuthService) *ProductHandler {
	return &ProductHandler{svc: svc, authSvc: authSvc}
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
	f.Page = parseQueryInt(c, "page", 1)
	f.Limit = parseQueryInt(c, "limit", 20)
	// Hard cap — never return more than 100 products per request
	if f.Limit > 100 {
		f.Limit = 100
	}
	if f.Page < 1 {
		f.Page = 1
	}

	products, err := h.svc.ListApproved(c.Context(), f)
	if err != nil {
		return response.InternalError(c, err)
	}
	return response.OKList(c, products)
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

