package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/your-org/used-tools-marketplace/internal/models"
	"github.com/your-org/used-tools-marketplace/internal/service"
)

type ClassificationHandler struct {
	svc *service.ClassificationService
}

func NewClassificationHandler(svc *service.ClassificationService) *ClassificationHandler {
	return &ClassificationHandler{svc: svc}
}

// ── Scheme Endpoints ─────────────────────────────────────────────────────────

// GET /api/classification/schemes
func (h *ClassificationHandler) ListSchemes(c *fiber.Ctx) error {
	status := c.Query("status", "ACTIVE")
	schemes, err := h.svc.ListSchemes(c.Context(), status)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"success": true,
		"data":    schemes,
	})
}

// GET /api/classification/schemes/:code
func (h *ClassificationHandler) GetScheme(c *fiber.Ctx) error {
	code := c.Params("code")
	scheme, err := h.svc.GetSchemeByCode(c.Context(), code)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"success": true,
		"data":    scheme,
	})
}

// POST /api/admin/classification/schemes
func (h *ClassificationHandler) CreateScheme(c *fiber.Ctx) error {
	var req models.ClassificationScheme
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid request body",
		})
	}
	if uid, ok := c.Locals("userID").(int); ok {
		req.CreatedBy = &uid
	}

	created, err := h.svc.CreateScheme(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    created,
	})
}

// PUT /api/admin/classification/schemes/:id
func (h *ClassificationHandler) UpdateScheme(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid scheme ID",
		})
	}
	var req models.ClassificationScheme
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid request body",
		})
	}
	req.ID = id
	if uid, ok := c.Locals("userID").(int); ok {
		req.UpdatedBy = &uid
	}

	if err := h.svc.UpdateScheme(c.Context(), &req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"success": true,
		"message": "scheme updated successfully",
	})
}

// ── Hierarchy & Placement Endpoints ─────────────────────────────────────────

// GET /api/classification/schemes/:code/tree
func (h *ClassificationHandler) GetTree(c *fiber.Ctx) error {
	code := c.Params("code")
	lang := c.Query("lang", "en")
	maxDepth := c.QueryInt("depth", 10)

	tree, err := h.svc.GetTree(c.Context(), code, lang, maxDepth)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"success": true,
		"data":    tree,
	})
}

// GET /api/classification/schemes/:code/roots
func (h *ClassificationHandler) GetRoots(c *fiber.Ctx) error {
	code := c.Params("code")
	lang := c.Query("lang", "en")

	roots, err := h.svc.GetRoots(c.Context(), code, lang)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"success": true,
		"data":    roots,
	})
}

// GET /api/classification/placements/:id
func (h *ClassificationHandler) GetPlacement(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid placement ID",
		})
	}
	lang := c.Query("lang", "en")

	placement, err := h.svc.GetPlacementByID(c.Context(), id, lang)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"success": true,
		"data":    placement,
	})
}

// GET /api/classification/placements/:id/children
func (h *ClassificationHandler) GetChildren(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid parent placement ID",
		})
	}
	lang := c.Query("lang", "en")

	children, err := h.svc.GetChildren(c.Context(), id, lang)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"success": true,
		"data":    children,
	})
}

// GET /api/classification/placements/:id/attributes
func (h *ClassificationHandler) GetPlacementAttributes(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid placement ID",
		})
	}

	attrs, err := h.svc.GetEffectiveAttributes(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"success": true,
		"data":    attrs,
	})
}

// POST /api/admin/classification/placements
func (h *ClassificationHandler) CreatePlacement(c *fiber.Ctx) error {
	var req models.ClassificationPlacement
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid request body",
		})
	}
	if uid, ok := c.Locals("userID").(int); ok {
		req.CreatedBy = &uid
	}

	created, err := h.svc.CreatePlacement(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    created,
	})
}

// PUT /api/admin/classification/placements/:id
func (h *ClassificationHandler) UpdatePlacement(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid placement ID",
		})
	}
	var req models.ClassificationPlacement
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid request body",
		})
	}
	req.ID = id
	if uid, ok := c.Locals("userID").(int); ok {
		req.UpdatedBy = &uid
	}

	if err := h.svc.UpdatePlacement(c.Context(), &req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"success": true,
		"message": "placement updated successfully",
	})
}

// PUT /api/admin/classification/placements/:id/move
func (h *ClassificationHandler) MovePlacement(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid placement ID",
		})
	}
	type MoveReq struct {
		NewParentID int `json:"new_parent_id"`
	}
	var req MoveReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid request body",
		})
	}

	if err := h.svc.MovePlacement(c.Context(), id, req.NewParentID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"success": true,
		"message": "placement moved successfully",
	})
}

// DELETE /api/admin/classification/placements/:id
func (h *ClassificationHandler) DeletePlacement(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid placement ID",
		})
	}

	if err := h.svc.DeletePlacement(c.Context(), id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"success": true,
		"message": "placement archived successfully",
	})
}

// ── Concepts Endpoints ───────────────────────────────────────────────────────

// GET /api/classification/concepts
func (h *ClassificationHandler) SearchConcepts(c *fiber.Ctx) error {
	filter := models.ConceptSearchFilter{
		SearchQuery: c.Query("q"),
		ConceptType: c.Query("type"),
		Status:      c.Query("status", "ACTIVE"),
		Page:        c.QueryInt("page", 0),
		Limit:       c.QueryInt("limit", 20),
		Lang:        c.Query("lang", "en"),
	}

	concepts, total, err := h.svc.SearchConcepts(c.Context(), filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"success": true,
		"data":    concepts,
		"total":   total,
		"page":    filter.Page,
		"limit":   filter.Limit,
	})
}

// GET /api/classification/concepts/:id
func (h *ClassificationHandler) GetConcept(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid concept ID",
		})
	}

	concept, err := h.svc.GetConceptByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"success": true,
		"data":    concept,
	})
}

// POST /api/admin/classification/concepts
func (h *ClassificationHandler) CreateConcept(c *fiber.Ctx) error {
	var req models.ClassificationConcept
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid request body",
		})
	}
	if uid, ok := c.Locals("userID").(int); ok {
		req.CreatedBy = &uid
	}

	created, err := h.svc.CreateConcept(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    created,
	})
}

// PUT /api/admin/classification/concepts/:id
func (h *ClassificationHandler) UpdateConcept(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid concept ID",
		})
	}
	var req models.ClassificationConcept
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid request body",
		})
	}
	req.ID = id
	if uid, ok := c.Locals("userID").(int); ok {
		req.UpdatedBy = &uid
	}

	if err := h.svc.UpdateConcept(c.Context(), &req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"success": true,
		"message": "concept updated successfully",
	})
}

// POST /api/admin/classification/concepts/check-duplicates
func (h *ClassificationHandler) CheckDuplicateCandidates(c *fiber.Ctx) error {
	type DupReq struct {
		Name string `json:"name"`
		Code string `json:"code"`
	}
	var req DupReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid request body",
		})
	}

	candidates, err := h.svc.CheckDuplicates(c.Context(), req.Name, req.Code)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"success": true,
		"data":    candidates,
	})
}

// ── Multi-Industry & Multi-Application Relationships ────────────────────────

// POST /api/admin/classification/relationships
func (h *ClassificationHandler) CreateRelationship(c *fiber.Ctx) error {
	var req models.ClassificationRelationship
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid request body",
		})
	}
	if uid, ok := c.Locals("userID").(int); ok {
		req.CreatedBy = &uid
	}

	if err := h.svc.CreateRelationship(c.Context(), &req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    req,
	})
}

// GET /api/classification/concepts/:id/relationships
func (h *ClassificationHandler) GetRelationships(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid concept ID",
		})
	}
	relType := c.Query("type")

	rels, err := h.svc.GetRelationships(c.Context(), id, relType)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"success": true,
		"data":    rels,
	})
}

// GET /api/classification/relationship-types
func (h *ClassificationHandler) ListRelationshipTypes(c *fiber.Ctx) error {
	types, err := h.svc.ListRelationshipTypes(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"success": true,
		"data":    types,
	})
}

// DELETE /api/admin/classification/relationships/:id
func (h *ClassificationHandler) DeleteRelationship(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid relationship ID",
		})
	}

	if err := h.svc.DeleteRelationship(c.Context(), id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"success": true,
		"message": "relationship deleted successfully",
	})
}

// ── Attributes Endpoints ────────────────────────────────────────────────────

// POST /api/admin/classification/placements/:id/attributes
func (h *ClassificationHandler) SaveAttributeRule(c *fiber.Ctx) error {
	placementID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid placement ID",
		})
	}
	var req models.ClassificationAttributeRule
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid request body",
		})
	}
	req.PlacementID = &placementID

	if err := h.svc.SaveAttributeRule(c.Context(), &req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"success": true,
		"message": "attribute rule saved successfully",
	})
}

// DELETE /api/admin/classification/attribute-rules/:id
func (h *ClassificationHandler) DeleteAttributeRule(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid attribute rule ID",
		})
	}

	if err := h.svc.DeleteAttributeRule(c.Context(), id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"success": true,
		"message": "attribute rule deleted successfully",
	})
}

// ── Impact Analysis, Governance & Search ────────────────────────────────────

// GET /api/admin/classification/impact-analysis
func (h *ClassificationHandler) AnalyzeImpact(c *fiber.Ctx) error {
	entityType := c.Query("type", "PLACEMENT")
	entityID := c.QueryInt("id")
	if entityID <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "id parameter is required",
		})
	}

	impact, err := h.svc.AnalyzeImpact(c.Context(), entityType, entityID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"success": true,
		"data":    impact,
	})
}

// GET /api/admin/classification/audit-logs
func (h *ClassificationHandler) ListAuditLogs(c *fiber.Ctx) error {
	entityType := c.Query("type")
	entityID := c.QueryInt("id")
	limit := c.QueryInt("limit", 20)
	page := c.QueryInt("page", 0)

	logs, err := h.svc.ListAuditLogs(c.Context(), entityType, entityID, limit, page*limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"success": true,
		"data":    logs,
	})
}

// GET /api/classification/search
func (h *ClassificationHandler) SearchUniversal(c *fiber.Ctx) error {
	q := c.Query("q")
	scheme := c.Query("scheme", "SYS_PROD")
	lang := c.Query("lang", "en")
	limit := c.QueryInt("limit", 20)

	results, err := h.svc.SearchUniversal(c.Context(), q, scheme, lang, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"success": true,
		"data":    results,
	})
}

