package handlers

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/your-org/used-tools-marketplace/internal/models"
	"github.com/your-org/used-tools-marketplace/internal/service"
	"github.com/your-org/used-tools-marketplace/pkg/response"
)

type MasterDataHandler struct {
	svc *service.MasterDataService
}

func NewMasterDataHandler(svc *service.MasterDataService) *MasterDataHandler {
	return &MasterDataHandler{svc: svc}
}

// POST /api/admin/masterdata/import/locations
func (h *MasterDataHandler) ImportLocations(c *fiber.Ctx) error {
	var req models.MasterDataImportRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid JSON payload for location import")
	}
	req.DatasetType = "locations"
	if len(req.LocationRecords) == 0 && len(req.Records) > 0 {
		var locs []models.MasterLocationNode
		if err := json.Unmarshal(req.Records, &locs); err == nil {
			req.LocationRecords = locs
		}
	}
	if len(req.LocationRecords) == 0 {
		return response.BadRequest(c, "location_records array cannot be empty")
	}

	adminID, _ := c.Locals("userID").(int)
	result, err := h.svc.ImportLocations(c.Context(), req, adminID, c.IP())
	if err != nil {
		return response.InternalError(c, err)
	}

	return response.OK(c, result)
}

// POST /api/admin/masterdata/import/categories
func (h *MasterDataHandler) ImportCategories(c *fiber.Ctx) error {
	var req models.MasterDataImportRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid JSON payload for category import")
	}
	req.DatasetType = "categories"
	if len(req.CategoryRecords) == 0 && len(req.Records) > 0 {
		var cats []models.MasterCategoryNode
		if err := json.Unmarshal(req.Records, &cats); err == nil {
			req.CategoryRecords = cats
		}
	}
	if len(req.CategoryRecords) == 0 {
		return response.BadRequest(c, "category_records array cannot be empty")
	}

	adminID, _ := c.Locals("userID").(int)
	result, err := h.svc.ImportCategories(c.Context(), req, adminID, c.IP())
	if err != nil {
		return response.InternalError(c, err)
	}

	return response.OK(c, result)
}

// GET /api/admin/masterdata/export/locations?region=ET-OR
func (h *MasterDataHandler) ExportLocations(c *fiber.Ctx) error {
	region := c.Query("region", "")
	locations, err := h.svc.ExportLocations(c.Context(), region)
	if err != nil {
		return response.InternalError(c, err)
	}

	return c.JSON(fiber.Map{
		"dataset_type":  "locations",
		"version":       "1.0.0",
		"generated_at":  time.Now().UTC().Format(time.RFC3339),
		"scope":         region,
		"total_records": len(locations),
		"records":       locations,
	})
}

// GET /api/admin/masterdata/export/categories?system=SYS_PROD&root=electronics
func (h *MasterDataHandler) ExportCategories(c *fiber.Ctx) error {
	system := c.Query("system", "")
	if system == "all" {
		system = ""
	}
	root := c.Query("root", "")
	categories, err := h.svc.ExportCategories(c.Context(), system, root)
	if err != nil {
		return response.InternalError(c, err)
	}

	return c.JSON(fiber.Map{
		"dataset_type":  "categories",
		"version":       "1.0.0",
		"generated_at":  time.Now().UTC().Format(time.RFC3339),
		"scope":         system,
		"total_records": len(categories),
		"records":       categories,
	})
}

// GET /api/admin/masterdata/status
func (h *MasterDataHandler) GetStatus(c *fiber.Ctx) error {
	status, err := h.svc.GetStatus(c.Context())
	if err != nil {
		return response.InternalError(c, err)
	}
	return response.OK(c, status)
}

// POST /api/admin/masterdata/validate
func (h *MasterDataHandler) ValidatePayload(c *fiber.Ctx) error {
	var req models.MasterDataImportRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid JSON payload")
	}

	var errors []models.MasterDataDiffItem
	if req.DatasetType == "locations" {
		errors = h.svc.ValidateLocations(req.LocationRecords, nil)
	} else if req.DatasetType == "categories" {
		errors = h.svc.ValidateCategories(req.CategoryRecords, nil)
	} else {
		return response.BadRequest(c, "invalid dataset_type: must be 'locations' or 'categories'")
	}

	valid := len(errors) == 0
	return c.JSON(fiber.Map{
		"valid":          valid,
		"errors_count":   len(errors),
		"validation_log": errors,
	})
}

// Dummy helper for strconv to prevent unused import
var _ = strconv.Itoa

