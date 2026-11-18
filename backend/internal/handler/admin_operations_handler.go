package handlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/your-org/used-tools-marketplace/internal/models"
	"github.com/your-org/used-tools-marketplace/pkg/response"
)

func (h *AdminHandler) GetOperationsHealth(c *fiber.Ctx) error {
	health, err := h.svc.GetOperationsHealth(c.Context())
	if err != nil {
		return response.InternalError(c, err)
	}
	return response.OK(c, health)
}

// GET /api/admin/operations/jobs
func (h *AdminHandler) GetJobStatusSummary(c *fiber.Ctx) error {
	jobs, err := h.svc.ListJobStatusSummary(c.Context())
	if err != nil {
		return response.InternalError(c, err)
	}
	return response.OKList(c, jobs)
}

// GET /api/admin/operations/failed-jobs
func (h *AdminHandler) GetFailedJobs(c *fiber.Ctx) error {
	limit := parseQueryInt(c, "limit", 50)
	jobs, err := h.svc.ListRecentFailedJobs(c.Context(), limit)
	if err != nil {
		return response.InternalError(c, err)
	}
	return response.OKList(c, jobs)
}

// ── Seller Verification (Role 08: Seller Verification Admin) ──────────────────

// GET /api/admin/verification/documents
func (h *AdminHandler) ListVerificationDocuments(c *fiber.Ctx) error {
	status := c.Query("status", "pending")
	page, limit := parsePagination(c, 50, 100)

	docs, total, err := h.svc.ListVerificationDocuments(c.Context(), status, page, limit)
	if err != nil {
		return response.InternalError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    docs,
		"meta": fiber.Map{
			"total": total,
			"page":  page,
			"limit": limit,
		},
	})
}

// PUT /api/admin/verification/documents/:id/review
func (h *AdminHandler) ReviewVerificationDocument(c *fiber.Ctx) error {
	adminID, err := mustUserID(c)
	if err != nil {
		return err
	}

	docID, err := parseIDParam(c, "id")
	if err != nil {
		return err
	}

	var body struct {
		Approved bool   `json:"approved"`
		Reason   string `json:"reason"`
		Notes    string `json:"notes"`
	}
	if err := c.BodyParser(&body); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	if err := h.svc.ReviewVerificationDocument(c.Context(), docID, adminID, body.Approved, body.Reason, body.Notes); err != nil {
		return response.InternalError(c, err)
	}

	// Log the action
	action := "approve_verification"
	if !body.Approved {
		action = "reject_verification"
	}
	h.svc.InsertAdminLog(c.Context(), models.AdminLog{
		AdminID:   adminID,
		Action:    action,
		Notes:     fmt.Sprintf("Document ID: %d, Reason: %s", docID, body.Reason),
		IPAddress: c.IP(),
	})

	return response.OK(c, fiber.Map{
		"message": "document reviewed successfully",
	})
}

// ── Trust & Safety (Role 21: Fraud & Trust & Safety Admin) ────────────────────

// GET /api/admin/trust/reports
func (h *AdminHandler) ListAbuseReports(c *fiber.Ctx) error {
	filter := models.AbuseReportFilter{
		Status:         c.Query("status", "pending"),
		Priority:       c.Query("priority", "all"),
		Reason:         c.Query("reason", "all"),
		ReportedUserID: parseQueryInt(c, "reported_user_id", 0),
	}
	filter.Page, filter.Limit = parsePagination(c, 50, 100)

	reports, total, err := h.svc.ListAbuseReports(c.Context(), filter)
	if err != nil {
		return response.InternalError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    reports,
		"meta": fiber.Map{
			"total": total,
			"page":  filter.Page,
			"limit": filter.Limit,
		},
	})
}

// PUT /api/admin/trust/reports/:id/review
func (h *AdminHandler) ReviewAbuseReport(c *fiber.Ctx) error {
	adminID, err := mustUserID(c)
	if err != nil {
		return err
	}

	reportID, err := parseIDParam(c, "id")
	if err != nil {
		return err
	}

	var body struct {
		Status      string `json:"status"` // 'investigating', 'resolved', 'dismissed'
		ActionTaken string `json:"action_taken"`
		Action      string `json:"action"`
		Notes       string `json:"notes"`
	}
	if err := c.BodyParser(&body); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	if body.ActionTaken == "" && body.Action != "" {
		body.ActionTaken = body.Action
	}

	if body.Status == "" {
		return response.BadRequest(c, "status is required")
	}

	if err := h.svc.ReviewAbuseReport(c.Context(), reportID, adminID, body.Status, body.ActionTaken, body.Notes); err != nil {
		return response.BadRequest(c, err.Error())
	}

	// Log the action
	h.svc.InsertAdminLog(c.Context(), models.AdminLog{
		AdminID:   adminID,
		Action:    "review_abuse_report",
		Notes:     fmt.Sprintf("Report ID: %d, Status: %s, Action: %s", reportID, body.Status, body.ActionTaken),
		IPAddress: c.IP(),
	})

	return response.OK(c, fiber.Map{
		"message": "abuse report reviewed successfully",
	})
}

// PUT /api/admin/trust/reports/:id/priority
func (h *AdminHandler) UpdateAbuseReportPriority(c *fiber.Ctx) error {
	reportID, err := parseIDParam(c, "id")
	if err != nil {
		return err
	}

	var body struct {
		Priority string `json:"priority"` // 'low', 'normal', 'high', 'critical'
	}
	if err := c.BodyParser(&body); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	if err := h.svc.UpdateAbuseReportPriority(c.Context(), reportID, body.Priority); err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.OK(c, fiber.Map{
		"message": "priority updated successfully",
	})
}

// ── Refund Management (Role 20: Refund & Dispute Admin) ───────────────────────

// GET /api/admin/refunds
func (h *AdminHandler) ListRefundRequests(c *fiber.Ctx) error {
	filter := models.RefundRequestFilter{
		Status:        c.Query("status", "pending"),
		RequesterRole: c.Query("requester_role", "all"),
	}
	filter.Page, filter.Limit = parsePagination(c, 50, 100)

	refunds, total, err := h.svc.ListRefundRequests(c.Context(), filter)
	if err != nil {
		return response.InternalError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    refunds,
		"meta": fiber.Map{
			"total": total,
			"page":  filter.Page,
			"limit": filter.Limit,
		},
	})
}

// PUT /api/admin/refunds/:id/review
func (h *AdminHandler) ReviewRefundRequest(c *fiber.Ctx) error {
	adminID, err := mustUserID(c)
	if err != nil {
		return err
	}

	refundID, err := parseIDParam(c, "id")
	if err != nil {
		return err
	}

	var body struct {
		Approved bool   `json:"approved"`
		Reason   string `json:"reason"`
	}
	if err := c.BodyParser(&body); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	if err := h.svc.ReviewRefundRequest(c.Context(), refundID, adminID, body.Approved, body.Reason); err != nil {
		return response.BadRequest(c, err.Error())
	}

	// Log the action
	action := "approve_refund"
	if !body.Approved {
		action = "reject_refund"
	}
	h.svc.InsertAdminLog(c.Context(), models.AdminLog{
		AdminID:   adminID,
		Action:    action,
		Notes:     fmt.Sprintf("Refund ID: %d, Reason: %s", refundID, body.Reason),
		IPAddress: c.IP(),
	})

	return response.OK(c, fiber.Map{
		"message": "refund request reviewed successfully",
	})
}

// PUT /api/admin/refunds/:id/processed
func (h *AdminHandler) MarkRefundProcessed(c *fiber.Ctx) error {
	adminID, err := mustUserID(c)
	if err != nil {
		return err
	}

	refundID, err := parseIDParam(c, "id")
	if err != nil {
		return err
	}

	var body struct {
		Reference string `json:"reference"`
	}
	if err := c.BodyParser(&body); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	if body.Reference == "" {
		return response.BadRequest(c, "refund reference is required")
	}

	if err := h.svc.MarkRefundProcessed(c.Context(), refundID, body.Reference); err != nil {
		return response.BadRequest(c, err.Error())
	}

	// Log the action
	h.svc.InsertAdminLog(c.Context(), models.AdminLog{
		AdminID:   adminID,
		Action:    "process_refund",
		Notes:     fmt.Sprintf("Refund ID: %d, Reference: %s", refundID, body.Reference),
		IPAddress: c.IP(),
	})

	return response.OK(c, fiber.Map{
		"message": "refund marked as processed",
	})
}

// ── Account Appeals (Role 30: Appeals & Account Recovery Admin) ───────────────

// GET /api/admin/appeals
func (h *AdminHandler) ListAccountAppeals(c *fiber.Ctx) error {
	filter := models.AccountAppealFilter{
		Status:     c.Query("status", "pending"),
		AppealType: c.Query("appeal_type", "all"),
		UserID:     parseQueryInt(c, "user_id", 0),
	}
	filter.Page, filter.Limit = parsePagination(c, 50, 100)

	appeals, total, err := h.svc.ListAccountAppeals(c.Context(), filter)
	if err != nil {
		return response.InternalError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    appeals,
		"meta": fiber.Map{
			"total": total,
			"page":  filter.Page,
			"limit": filter.Limit,
		},
	})
}

// PUT /api/admin/appeals/:id/review
func (h *AdminHandler) ReviewAccountAppeal(c *fiber.Ctx) error {
	adminID, err := mustUserID(c)
	if err != nil {
		return err
	}

	appealID, err := parseIDParam(c, "id")
	if err != nil {
		return err
	}

	var body struct {
		Approved    bool   `json:"approved"`
		Decision    string `json:"decision"`
		ActionTaken string `json:"action_taken"`
	}
	if err := c.BodyParser(&body); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	if body.Decision == "" {
		return response.BadRequest(c, "decision is required")
	}

	if err := h.svc.ReviewAccountAppeal(c.Context(), appealID, adminID, body.Approved, body.Decision, body.ActionTaken); err != nil {
		return response.BadRequest(c, err.Error())
	}

	// Log the action
	action := "approve_appeal"
	if !body.Approved {
		action = "reject_appeal"
	}
	h.svc.InsertAdminLog(c.Context(), models.AdminLog{
		AdminID:   adminID,
		Action:    action,
		Notes:     fmt.Sprintf("Appeal ID: %d, Decision: %s, Action: %s", appealID, body.Decision, body.ActionTaken),
		IPAddress: c.IP(),
	})

	return response.OK(c, fiber.Map{
		"message": "appeal reviewed successfully",
	})
}

// PUT /api/admin/appeals/:id/status
func (h *AdminHandler) UpdateAppealStatus(c *fiber.Ctx) error {
	appealID, err := parseIDParam(c, "id")
	if err != nil {
		return err
	}

	var body struct {
		Status string `json:"status"` // 'pending', 'under_review', 'approved', 'rejected'
	}
	if err := c.BodyParser(&body); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	if err := h.svc.UpdateAppealStatus(c.Context(), appealID, body.Status); err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.OK(c, fiber.Map{
		"message": "appeal status updated successfully",
	})
}
