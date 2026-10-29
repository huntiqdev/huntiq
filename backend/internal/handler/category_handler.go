package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/your-org/used-tools-marketplace/internal/models"
	"github.com/your-org/used-tools-marketplace/internal/service"
)

type CategoryHandler struct {
	svc *service.CategoryHierarchyService
}

func NewCategoryHandler(svc *service.CategoryHierarchyService) *CategoryHandler {
	return &CategoryHandler{svc: svc}
}

// GetCategoryTree returns the full hierarchical category tree
// GET /api/categories/tree?lang=en
func (h *CategoryHandler) GetCategoryTree(c *fiber.Ctx) error {
	lang := c.Query("lang", "en")
	
	tree, err := h.svc.GetCategoryTree(c.Context(), lang)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "failed to retrieve category tree",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    tree,
	})
}

// GetCategory returns a single category with breadcrumbs
// GET /api/categories/:id?lang=en
func (h *CategoryHandler) GetCategory(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid category ID",
		})
	}

	lang := c.Query("lang", "en")
	
	cat, err := h.svc.GetByID(c.Context(), id, lang)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   "category not found",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    cat,
	})
}

// GetCategoryBySlug returns a category by slug
// GET /api/categories/slug/:slug?lang=en
func (h *CategoryHandler) GetCategoryBySlug(c *fiber.Ctx) error {
	slug := c.Params("slug")
	lang := c.Query("lang", "en")
	
	cat, err := h.svc.GetBySlug(c.Context(), slug, lang)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   "category not found",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    cat,
	})
}

// ListCategories returns a filtered list of categories
// GET /api/categories?parent_id=1&level=2&is_active=true&page=1&limit=50
func (h *CategoryHandler) ListCategories(c *fiber.Ctx) error {
	filter := models.CategoryFilter{
		Page:  1,
		Limit: 50,
	}

	if parentID := c.QueryInt("parent_id", -1); parentID >= 0 {
		filter.ParentID = &parentID
	}
	if level := c.QueryInt("level", -1); level >= 0 {
		filter.Level = &level
	}
	if isActive := c.Query("is_active"); isActive != "" {
		active := isActive == "true"
		filter.IsActive = &active
	}
	if isLeaf := c.Query("is_leaf"); isLeaf != "" {
		leaf := isLeaf == "true"
		filter.IsLeaf = &leaf
	}
	if search := c.Query("search"); search != "" {
		filter.SearchTerm = search
	}
	if lang := c.Query("lang"); lang != "" {
		filter.Lang = lang
	}
	if page := c.QueryInt("page", 1); page > 0 {
		filter.Page = page
	}
	if limit := c.QueryInt("limit", 50); limit > 0 && limit <= 500 {
		filter.Limit = limit
	}

	categories, err := h.svc.ListCategories(c.Context(), filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "failed to retrieve categories",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    categories,
	})
}

// GetCategoryChildren returns direct children of a category
// GET /api/categories/:id/children?lang=en
func (h *CategoryHandler) GetCategoryChildren(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid category ID",
		})
	}

	lang := c.Query("lang", "en")
	filter := models.CategoryFilter{
		ParentID: &id,
		Lang:     lang,
		Page:     1,
		Limit:    500,
	}
	
	active := true
	filter.IsActive = &active

	children, err := h.svc.ListCategories(c.Context(), filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "failed to retrieve child categories",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    children,
	})
}

// GetCategoryBreadcrumbs returns the breadcrumb path for a category
// GET /api/categories/:id/ancestors?lang=en
func (h *CategoryHandler) GetCategoryBreadcrumbs(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid category ID",
		})
	}

	lang := c.Query("lang", "en")
	
	breadcrumbs, err := h.svc.GetBreadcrumbs(c.Context(), id, lang)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "failed to retrieve breadcrumbs",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    breadcrumbs,
	})
}

// CreateCategory creates a new category (admin only)
// POST /api/admin/categories
func (h *CategoryHandler) CreateCategory(c *fiber.Ctx) error {
	var req models.CategoryHierarchy
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid request body",
		})
	}

	id, err := h.svc.CreateCategory(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"id": id,
		},
	})
}

// UpdateCategory updates an existing category (admin only)
// PUT /api/admin/categories/:id
func (h *CategoryHandler) UpdateCategory(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid category ID",
		})
	}

	var req models.CategoryHierarchy
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid request body",
		})
	}

	req.ID = id
	
	if err := h.svc.UpdateCategory(c.Context(), &req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "category updated successfully",
	})
}

// DeleteCategory deletes a category (admin only)
// DELETE /api/admin/categories/:id
func (h *CategoryHandler) DeleteCategory(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid category ID",
		})
	}

	if err := h.svc.DeleteCategory(c.Context(), id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "category deleted successfully",
	})
}

// GetCategoryStats returns category statistics (admin only)
// GET /api/admin/categories/:id/stats
func (h *CategoryHandler) GetCategoryStats(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid category ID",
		})
	}

	stats, err := h.svc.GetCategoryStats(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   "category not found",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    stats,
	})
}

// GetRoots returns all active root/pillar categories (level = 1)
// GET /api/categories/roots?lang=en
func (h *CategoryHandler) GetRoots(c *fiber.Ctx) error {
	lang := c.Query("lang", "en")
	roots, err := h.svc.GetRoots(c.Context(), lang)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "failed to retrieve root categories",
		})
	}
	return c.JSON(fiber.Map{
		"success": true,
		"data":    roots,
	})
}

// CheckSlug checks if a category slug is available
// GET /api/categories/check-slug?slug=xyz&exclude_id=12
func (h *CategoryHandler) CheckSlug(c *fiber.Ctx) error {
	slug := c.Query("slug")
	if slug == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "slug query parameter is required",
		})
	}
	excludeID := c.QueryInt("exclude_id", 0)
	available, err := h.svc.CheckSlugAvailability(c.Context(), slug, excludeID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "failed to check slug availability",
		})
	}
	return c.JSON(fiber.Map{
		"success":   true,
		"available": available,
		"slug":      slug,
	})
}

// ReorderCategory updates a category's display order (admin only)
// PUT /api/admin/categories/:id/reorder
func (h *CategoryHandler) ReorderCategory(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid category ID",
		})
	}
	var body struct {
		DisplayOrder int `json:"display_order"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid request body",
		})
	}
	if err := h.svc.ReorderCategory(c.Context(), id, body.DisplayOrder); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"success": true,
		"message": "category order updated successfully",
	})
}

// SetCategoryStatus toggles active/inactive state of a category (admin only)
// PUT /api/admin/categories/:id/status
func (h *CategoryHandler) SetCategoryStatus(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid category ID",
		})
	}
	var body struct {
		IsActive bool `json:"is_active"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid request body",
		})
	}
	if err := h.svc.SetCategoryStatus(c.Context(), id, body.IsActive); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"success": true,
		"message": "category status updated successfully",
	})
}
