package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/your-org/used-tools-marketplace/internal/models"
	"github.com/your-org/used-tools-marketplace/internal/service"
)

type AttributeHandler struct {
	svc        *service.AttributeService
	productSvc *service.ProductService
}

func NewAttributeHandler(svc *service.AttributeService) *AttributeHandler {
	return &AttributeHandler{svc: svc}
}

func (h *AttributeHandler) WithProductService(productSvc *service.ProductService) *AttributeHandler {
	h.productSvc = productSvc
	return h
}

// ListAttributeDefinitions returns all attribute definitions
// GET /api/attributes
func (h *AttributeHandler) ListAttributeDefinitions(c *fiber.Ctx) error {
	attrs, err := h.svc.ListAttributeDefinitions(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "failed to retrieve attributes",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    attrs,
	})
}

// GetCategoryAttributes returns applicable attributes for a category
// GET /api/categories/:id/attributes
func (h *AttributeHandler) GetCategoryAttributes(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid category ID",
		})
	}

	attrs, err := h.svc.GetCategoryAttributes(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "failed to retrieve category attributes",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    attrs,
	})
}

// GetProductAttributes returns all attributes for a product
// GET /api/products/:id/attributes
func (h *AttributeHandler) GetProductAttributes(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid product ID",
		})
	}

	// Verify buyer tier authorization before retrieving specifications
	if h.productSvc != nil {
		prod, err := h.productSvc.GetByID(c.Context(), id)
		if err == nil && prod != nil {
			userTier, _ := c.Locals("buyerTier").(string)
			role, _ := c.Locals("role").(string)
			adminRole, _ := c.Locals("adminRole").(string)
			isAdmin := role == "admin" || adminRole != ""

			if !isTierAuthorized(userTier, prod.RequiredBuyerTier, isAdmin) {
				return c.JSON(fiber.Map{
					"success":       true,
					"data":          []models.ProductAttribute{},
					"locked":        true,
					"required_tier": prod.RequiredBuyerTier,
				})
			}
		}
	}

	attrs, err := h.svc.GetProductAttributes(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "failed to retrieve product attributes",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    attrs,
		"locked":  false,
	})
}

// SetProductAttributes sets multiple product attributes
// POST /api/seller/products/:id/attributes
func (h *AttributeHandler) SetProductAttributes(c *fiber.Ctx) error {
	idParam := c.Params("id")
	productID, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid product ID",
		})
	}

	var req struct {
		CategoryID int                      `json:"category_id"`
		Attributes []models.ProductAttribute `json:"attributes"`
	}
	
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid request body",
		})
	}

	if err := h.svc.SetProductAttributes(c.Context(), productID, req.CategoryID, req.Attributes); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "product attributes updated successfully",
	})
}

// CreateAttributeDefinition creates a new attribute (admin)
// POST /api/admin/attributes
func (h *AttributeHandler) CreateAttributeDefinition(c *fiber.Ctx) error {
	var req models.AttributeDefinition
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid request body",
		})
	}

	id, err := h.svc.CreateAttributeDefinition(c.Context(), &req)
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

// AssignAttributeToCategory assigns an attribute to a category (admin)
// POST /api/admin/categories/:id/attributes
func (h *AttributeHandler) AssignAttributeToCategory(c *fiber.Ctx) error {
	idParam := c.Params("id")
	categoryID, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid category ID",
		})
	}

	var req models.CategoryAttribute
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid request body",
		})
	}

	req.CategoryID = categoryID
	
	if err := h.svc.AssignAttributeToCategory(c.Context(), &req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "attribute assigned to category successfully",
	})
}

// RemoveAttributeFromCategory removes an attribute from a category (admin)
// DELETE /api/admin/categories/:id/attributes/:attribute_id
func (h *AttributeHandler) RemoveAttributeFromCategory(c *fiber.Ctx) error {
	categoryID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid category ID",
		})
	}

	attributeID, err := strconv.Atoi(c.Params("attribute_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid attribute ID",
		})
	}

	if err := h.svc.RemoveAttributeFromCategory(c.Context(), categoryID, attributeID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "attribute removed from category successfully",
	})
}
