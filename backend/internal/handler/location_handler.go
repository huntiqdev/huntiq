package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/your-org/used-tools-marketplace/internal/service"
	"github.com/your-org/used-tools-marketplace/pkg/response"
)

type LocationHandler struct {
	svc *service.LocationService
}

func NewLocationHandler(svc *service.LocationService) *LocationHandler {
	return &LocationHandler{svc: svc}
}

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
// subcity is optional when zone is provided (some regions have no subcity level)
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
	// Cap search query length to prevent abuse
	if len(q) > 100 {
		q = q[:100]
	}
	results, err := h.svc.SearchLocations(c.Context(), q)
	if err != nil {
		return response.InternalError(c, err)
	}
	return response.OKList(c, results)
}
