package handlers

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/your-org/used-tools-marketplace/internal/models"
	"github.com/your-org/used-tools-marketplace/internal/service"
	"github.com/your-org/used-tools-marketplace/pkg/response"
)

type OrganizationHandler struct {
	svc *service.OrganizationService
}

func NewOrganizationHandler(svc *service.OrganizationService) *OrganizationHandler {
	return &OrganizationHandler{svc: svc}
}

// POST /api/organizations
func (h *OrganizationHandler) RegisterOrganization(c *fiber.Ctx) error {
	userID, err := mustUserID(c)
	if err != nil {
		return err
	}

	var body struct {
		Name                        string  `json:"name"`
		EntityType                  string  `json:"entity_type"` // COMPANY, GOVERNMENT, OTHER_LEGALLY_APPROVED_ENTITY
		RegistrationNumber          *string `json:"registration_number"`
		TINNumber                   *string `json:"tin_number"`
		OfficialAddress             *string `json:"official_address"`
		OfficialEmail               *string `json:"official_email"`
		OfficialPhone               *string `json:"official_phone"`
		PublishingAuthorityEvidence *string `json:"publishing_authority_evidence"`
	}
	if err := c.BodyParser(&body); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	if strings.TrimSpace(body.Name) == "" {
		return response.BadRequest(c, "organization name is required")
	}
	if body.EntityType == "" {
		body.EntityType = models.EntityCompany
	}

	org := &models.Organization{
		Name:                        body.Name,
		EntityType:                  body.EntityType,
		RegistrationNumber:          body.RegistrationNumber,
		TINNumber:                   body.TINNumber,
		OfficialAddress:             body.OfficialAddress,
		OfficialEmail:               body.OfficialEmail,
		OfficialPhone:               body.OfficialPhone,
		PublishingAuthorityEvidence: body.PublishingAuthorityEvidence,
	}

	if err := h.svc.RegisterOrganization(c.Context(), userID, org); err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.OK(c, fiber.Map{
		"organization": org,
		"message":      "organization registration submitted for verification review",
	})
}

// GET /api/organizations/me
func (h *OrganizationHandler) GetMyOrganization(c *fiber.Ctx) error {
	userID, err := mustUserID(c)
	if err != nil {
		return err
	}

	org, member, err := h.svc.GetUserOrganization(c.Context(), userID)
	if err != nil {
		return response.InternalError(c, err)
	}
	if org == nil {
		return response.OK(c, fiber.Map{
			"has_organization": false,
		})
	}

	return response.OK(c, fiber.Map{
		"has_organization": true,
		"organization":     org,
		"member":           member,
	})
}

// GET /api/admin/organizations
func (h *OrganizationHandler) AdminListOrganizations(c *fiber.Ctx) error {
	status := c.Query("status", "")
	entityType := c.Query("entity_type", "")
	page, limit := parsePagination(c, 50, 100)
	offset := (page - 1) * limit

	list, total, err := h.svc.ListOrganizations(c.Context(), status, entityType, limit, offset)
	if err != nil {
		return response.InternalError(c, err)
	}

	return response.OK(c, fiber.Map{
		"organizations": list,
		"total":         total,
	})
}

// PUT /api/admin/organizations/:id/review
func (h *OrganizationHandler) AdminReviewOrganization(c *fiber.Ctx) error {
	adminID, err := mustUserID(c)
	if err != nil {
		return err
	}

	id, err := parseIDParam(c, "id")
	if err != nil {
		return err
	}

	var body struct {
		Approved bool   `json:"approved"`
		Notes    string `json:"notes"`
	}
	if err := c.BodyParser(&body); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	if err := h.svc.ReviewOrganization(c.Context(), id, adminID, body.Approved, body.Notes); err != nil {
		return response.BadRequest(c, err.Error())
	}

	statusMsg := "rejected"
	if body.Approved {
		statusMsg = "verified"
	}

	return response.OK(c, fiber.Map{
		"message": "organization marked as " + statusMsg,
	})
}

