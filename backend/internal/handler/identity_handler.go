package handlers

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/your-org/used-tools-marketplace/internal/service"
	"github.com/your-org/used-tools-marketplace/pkg/response"
)

type IdentityHandler struct {
	svc *service.IdentityService
}

func NewIdentityHandler(svc *service.IdentityService) *IdentityHandler {
	return &IdentityHandler{svc: svc}
}

// POST /api/identity/verify
func (h *IdentityHandler) VerifyIdentity(c *fiber.Ctx) error {
	userID, err := mustUserID(c)
	if err != nil {
		return err
	}

	var body struct {
		FIN          string `json:"fin"`
		FullName     string `json:"full_legal_name"`
		Phone        string `json:"phone"`
		ConsentGiven bool   `json:"consent_given"`
		AuthCode     string `json:"auth_code"`
	}
	if err := c.BodyParser(&body); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	if strings.TrimSpace(body.FIN) == "" {
		return response.BadRequest(c, "national identity credential (FIN) is required")
	}
	if strings.TrimSpace(body.FullName) == "" {
		return response.BadRequest(c, "legal full name is required")
	}
	if !body.ConsentGiven {
		return response.BadRequest(c, "explicit consent to national identity verification is required")
	}

	req := &service.FaydaVerificationRequest{
		UserID:       userID,
		FIN:          body.FIN,
		FullName:     body.FullName,
		Phone:        body.Phone,
		ConsentGiven: body.ConsentGiven,
		AuthCode:     body.AuthCode,
	}

	verification, err := h.svc.ProcessIndividualIdentityVerification(c.Context(), req)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.OK(c, fiber.Map{
		"verification": verification,
		"message":      "identity verification processed successfully",
	})
}

// GET /api/identity/status
func (h *IdentityHandler) GetMyIdentityStatus(c *fiber.Ctx) error {
	userID, err := mustUserID(c)
	if err != nil {
		return err
	}

	record, err := h.svc.GetUserIdentityVerification(c.Context(), userID)
	if err != nil {
		return response.InternalError(c, err)
	}
	if record == nil {
		return response.OK(c, fiber.Map{
			"status":            "unverified",
			"identity_verified": false,
		})
	}

	return response.OK(c, fiber.Map{
		"status":            record.Status,
		"identity_verified": record.Status == "verified",
		"id_masked":         record.IDMasked,
		"provider_code":     record.ProviderCode,
		"auth_reference":    record.AuthReference,
		"phone_matched":     record.PhoneMatched,
		"verified_at":       record.VerifiedAt,
	})
}

// GET /api/admin/identity/pending
func (h *IdentityHandler) AdminListPending(c *fiber.Ctx) error {
	page, limit := parsePagination(c, 50, 100)
	offset := (page - 1) * limit

	list, total, err := h.svc.ListPendingVerifications(c.Context(), limit, offset)
	if err != nil {
		return response.InternalError(c, err)
	}

	return response.OK(c, fiber.Map{
		"verifications": list,
		"total":         total,
	})
}

// PUT /api/admin/identity/:id/resolve
func (h *IdentityHandler) AdminResolveException(c *fiber.Ctx) error {
	adminID, err := mustUserID(c)
	if err != nil {
		return err
	}

	id, err := parseIDParam(c, "id")
	if err != nil {
		return err
	}

	var body struct {
		Status string `json:"status"` // verified or rejected
		Reason string `json:"reason"`
	}
	if err := c.BodyParser(&body); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	if err := h.svc.ResolveException(c.Context(), int64(id), body.Status, adminID, body.Reason); err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.OK(c, fiber.Map{
		"message": "identity verification exception resolved successfully",
	})
}

