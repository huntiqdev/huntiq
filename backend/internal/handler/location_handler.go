package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/your-org/used-tools-marketplace/internal/models"
	"github.com/your-org/used-tools-marketplace/internal/service"
	"github.com/your-org/used-tools-marketplace/pkg/response"
)

type LocationHandler struct {
	svc *service.LocationService
}

func NewLocationHandler(svc *service.LocationService) *LocationHandler {
	return &LocationHandler{svc: svc}
}

// ── GPS Geocoding & Hierarchy Public Endpoints ───────────────────────────────

// POST /api/locations/reverse-geocode
// Resolves GPS coordinates to the verified administrative hierarchy with confidence scoring
func (h *LocationHandler) ReverseGeocode(c *fiber.Ctx) error {
	var body struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	}
	if err := c.BodyParser(&body); err != nil {
		return response.BadRequest(c, "invalid JSON payload with latitude and longitude")
	}
	if body.Latitude == 0 && body.Longitude == 0 {
		return response.BadRequest(c, "latitude and longitude are required")
	}

	result, err := h.svc.ReverseGeocode(c.Context(), body.Latitude, body.Longitude)
	if err != nil {
		return response.NotFound(c, err.Error())
	}
	return response.OK(c, result)
}

// GET /api/locations/tree?lang=...
func (h *LocationHandler) GetLocationTree(c *fiber.Ctx) error {
	lang := c.Query("lang", "en")
	tree, err := h.svc.GetLocationTree(c.Context(), lang)
	if err != nil {
		return response.InternalError(c, err)
	}
	return response.OKList(c, tree)
}

// GET /api/locations/children?parent_id=...&lang=...
func (h *LocationHandler) ListChildren(c *fiber.Ctx) error {
	lang := c.Query("lang", "en")
	var parentID *int
	if pStr := c.Query("parent_id"); pStr != "" {
		if id, err := strconv.Atoi(pStr); err == nil {
			parentID = &id
		}
	}
	children, err := h.svc.ListChildren(c.Context(), parentID, lang)
	if err != nil {
		return response.InternalError(c, err)
	}
	return response.OKList(c, children)
}

// GET /api/locations/search-hierarchy?q=...&lang=...
func (h *LocationHandler) SearchHierarchy(c *fiber.Ctx) error {
	q := c.Query("q")
	if q == "" {
		return response.BadRequest(c, "query param 'q' is required")
	}
	lang := c.Query("lang", "en")
	results, err := h.svc.SearchHierarchy(c.Context(), q, lang)
	if err != nil {
		return response.InternalError(c, err)
	}
	return response.OKList(c, results)
}

// ── Admin Hierarchy Management Endpoints ─────────────────────────────────────

// GET /api/admin/location-hierarchy
func (h *LocationHandler) AdminListLocations(c *fiber.Ctx) error {
	var f models.LocationFilter
	if p := c.Query("parent_id"); p != "" {
		if id, err := strconv.Atoi(p); err == nil {
			f.ParentID = &id
		}
	}
	f.LocationType = c.Query("type")
	if lvl := c.Query("level"); lvl != "" {
		if l, err := strconv.Atoi(lvl); err == nil {
			f.Level = &l
		}
	}
	if act := c.Query("active"); act != "" {
		b := act == "true" || act == "1"
		f.IsActive = &b
	}
	f.SearchTerm = c.Query("q")
	f.Page = c.QueryInt("page", 1)
	f.Limit = c.QueryInt("limit", 50)

	locs, total, err := h.svc.ListLocationsAdmin(c.Context(), f)
	if err != nil {
		return response.InternalError(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    locs,
		"meta": fiber.Map{
			"total": total,
			"page":  f.Page,
			"limit": f.Limit,
		},
	})
}

// GET /api/admin/location-hierarchy/:id
func (h *LocationHandler) AdminGetLocation(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id <= 0 {
		return response.BadRequest(c, "invalid location ID")
	}
	loc, err := h.svc.GetLocationByID(c.Context(), id, "en")
	if err != nil {
		return response.NotFound(c, "location not found")
	}
	return response.OK(c, loc)
}

// POST /api/admin/location-hierarchy
func (h *LocationHandler) AdminCreateLocation(c *fiber.Ctx) error {
	var loc models.LocationHierarchy
	if err := c.BodyParser(&loc); err != nil {
		return response.BadRequest(c, "invalid location payload")
	}
	id, err := h.svc.CreateLocation(c.Context(), &loc)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}
	return response.Created(c, map[string]interface{}{
		"id":      id,
		"message": "location created successfully",
	})
}

// PUT /api/admin/location-hierarchy/:id
func (h *LocationHandler) AdminUpdateLocation(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id <= 0 {
		return response.BadRequest(c, "invalid location ID")
	}
	var loc models.LocationHierarchy
	if err := c.BodyParser(&loc); err != nil {
		return response.BadRequest(c, "invalid location payload")
	}
	loc.ID = id
	if err := h.svc.UpdateLocation(c.Context(), &loc); err != nil {
		return response.BadRequest(c, err.Error())
	}
	return response.OK(c, map[string]interface{}{
		"id":      id,
		"message": "location updated successfully",
	})
}

// DELETE /api/admin/location-hierarchy/:id
func (h *LocationHandler) AdminDeleteLocation(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id <= 0 {
		return response.BadRequest(c, "invalid location ID")
	}
	if err := h.svc.DeleteLocation(c.Context(), id); err != nil {
		return response.BadRequest(c, err.Error())
	}
	return response.OK(c, map[string]interface{}{
		"message": "location deleted successfully",
	})
}

// POST /api/admin/location-hierarchy/import
func (h *LocationHandler) AdminImportLocations(c *fiber.Ctx) error {
	var body struct {
		Locations []models.LocationHierarchy `json:"locations"`
	}
	if err := c.BodyParser(&body); err != nil {
		return response.BadRequest(c, "invalid JSON payload with locations array")
	}
	if len(body.Locations) == 0 {
		return response.BadRequest(c, "locations array cannot be empty")
	}
	imported, failed, err := h.svc.ImportLocations(c.Context(), body.Locations)
	if err != nil {
		return response.InternalError(c, err)
	}
	return response.OK(c, map[string]interface{}{
		"imported": imported,
		"failed":   failed,
		"total":    len(body.Locations),
	})
}

// ── Legacy Endpoints (Preserved for 100% Backward Compatibility) ─────────────

// GET /api/locations/regions
func (h *LocationHandler) ListRegions(c *fiber.Ctx) error {
	regions, err := h.svc.ListRegions(c.Context())
	if err != nil {
		return response.InternalError(c, err)
	}
	return response.OKList(c, regions)
}

// GET /api/locations/zones?region=...
func (h *LocationHandler) ListZones(c *fiber.Ctx) error {
	region := c.Query("region")
	if region == "" {
		return response.BadRequest(c, "region query param required")
	}
	zones, err := h.svc.ListZones(c.Context(), region)
	if err != nil {
		return response.InternalError(c, err)
	}
	return response.OKList(c, zones)
}

// GET /api/locations/subcities?region=...&zone=...
func (h *LocationHandler) ListSubcities(c *fiber.Ctx) error {
	region := c.Query("region")
	if region == "" {
		return response.BadRequest(c, "region query param required")
	}
	zone := c.Query("zone", "")
	subcities, err := h.svc.ListSubcities(c.Context(), region, zone)
	if err != nil {
		return response.InternalError(c, err)
	}
	return response.OKList(c, subcities)
}

// GET /api/locations/woredas?region=...&zone=...&subcity=...
func (h *LocationHandler) ListWoredas(c *fiber.Ctx) error {
	region := c.Query("region")
	if region == "" {
		return response.BadRequest(c, "region query param required")
	}
	zone := c.Query("zone", "")
	subcity := c.Query("subcity", "")
	if zone == "" && subcity == "" {
		return response.BadRequest(c, "zone or subcity query param required")
	}
	woredas, err := h.svc.ListWoredas(c.Context(), region, zone, subcity)
	if err != nil {
		return response.InternalError(c, err)
	}
	return response.OKList(c, woredas)
}

// GET /api/locations/search?q=...
func (h *LocationHandler) SearchLocations(c *fiber.Ctx) error {
	q := c.Query("q")
	if q == "" {
		return response.BadRequest(c, "q query param required")
	}
	if len(q) > 100 {
		q = q[:100]
	}
	results, err := h.svc.SearchLocations(c.Context(), q)
	if err != nil {
		return response.InternalError(c, err)
	}
	return response.OKList(c, results)
}
