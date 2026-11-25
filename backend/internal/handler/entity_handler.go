package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/your-org/used-tools-marketplace/internal/service"
)

type EntityHandler struct {
	svc *service.EntityService
}

func NewEntityHandler(svc *service.EntityService) *EntityHandler {
	return &EntityHandler{svc: svc}
}

func (h *EntityHandler) RegisterRoutes(router fiber.Router) {
	group := router.Group("/entities")
	group.Get("/resolve", h.ResolveEntity)
	group.Get("/slug/:slug", h.GetEntityBySlug)
	group.Get("/:id/graph", h.GetEntityGraph)
	group.Get("/:id/discovery-paths", h.GetEntityDiscoveryPaths)
	group.Get("/:id", h.GetEntityByID)
}

// ResolveEntity resolves user query or search text into canonical entities
// GET /api/entities/resolve?q=phne&lang=en&limit=5
func (h *EntityHandler) ResolveEntity(c *fiber.Ctx) error {
	q := c.Query("q")
	if q == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "query parameter 'q' is required",
		})
	}
	lang := c.Query("lang", "en")
	limit := c.QueryInt("limit", 5)

	results, err := h.svc.ResolveCanonicalEntity(c.Context(), q, lang, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "failed to resolve canonical entity",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"query":   q,
		"data":    results,
	})
}

// GetEntityByID returns the master canonical entity by integer ID
// GET /api/entities/:id
func (h *EntityHandler) GetEntityByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid entity id",
		})
	}

	entity, err := h.svc.GetEntityByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   "canonical entity not found",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    entity,
	})
}

// GetEntityBySlug returns the master canonical entity by unique slug
// GET /api/entities/slug/:slug
func (h *EntityHandler) GetEntityBySlug(c *fiber.Ctx) error {
	slug := c.Params("slug")
	if slug == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "entity slug is required",
		})
	}

	entity, err := h.svc.GetEntityBySlug(c.Context(), slug)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   "canonical entity not found",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    entity,
	})
}

// GetEntityGraph returns the relationship graph edges for the entity
// GET /api/entities/:id/graph
func (h *EntityHandler) GetEntityGraph(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid entity id",
		})
	}

	entity, err := h.svc.GetEntityByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   "canonical entity not found",
		})
	}

	rels, err := h.svc.GetEntityRelationships(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "failed to retrieve entity relationships",
		})
	}

	return c.JSON(fiber.Map{
		"success":       true,
		"entity_id":     entity.ID,
		"entity_name":   entity.CanonicalNameEn,
		"entity_nature": entity.EntityNature,
		"entity_type":   entity.EntityType,
		"relationships": rels,
	})
}

// GetEntityDiscoveryPaths returns all cross-domain discovery paths for an entity
// GET /api/entities/:id/discovery-paths
func (h *EntityHandler) GetEntityDiscoveryPaths(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid entity id",
		})
	}

	paths, err := h.svc.GetEntityDiscoveryPaths(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "failed to retrieve entity discovery paths",
		})
	}

	return c.JSON(fiber.Map{
		"success":   true,
		"entity_id": id,
		"paths":     paths,
	})
}
