package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/your-org/used-tools-marketplace/internal/models"
	"github.com/your-org/used-tools-marketplace/pkg/dto"
	"github.com/your-org/used-tools-marketplace/internal/service"
	"github.com/your-org/used-tools-marketplace/pkg/response"
	"github.com/your-org/used-tools-marketplace/pkg/validations"
)

// BuyerHandler handles all buyer-specific HTTP endpoints.
// Auth endpoints (register, login, forgot/reset password, OTP) remain in AuthHandler
// because they are shared with the seller flow and run before a JWT exists.
type BuyerHandler struct {
	authSvc  *service.AuthService
	adminSvc *service.AdminService
}

func NewBuyerHandler(authSvc *service.AuthService) *BuyerHandler {
	return &BuyerHandler{authSvc: authSvc}
}

func (h *BuyerHandler) WithAdminService(a *service.AdminService) *BuyerHandler {
	h.adminSvc = a
	return h
}

// GET /api/buyer/profile
func (h *BuyerHandler) GetProfile(c *fiber.Ctx) error {
	userID, err := mustUserID(c)
	if err != nil {
		return err
	}
	u, err := h.authSvc.GetProfile(c.Context(), userID)
	if err != nil {
		return response.NotFound(c, "profile not found")
	}
	return response.OK(c, u)
}

// PUT /api/buyer/profile
func (h *BuyerHandler) UpdateProfile(c *fiber.Ctx) error {
	userID, err := mustUserID(c)
	if err != nil {
		return err
	}
	var body struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	if err := c.BodyParser(&body); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	if body.Name == "" {
		return response.BadRequest(c, "name is required")
	}
	if body.Email != "" {
		if !validations.ValidateEmail(body.Email) {
			return response.BadRequest(c, "invalid email format")
		}
	}
	if err := h.authSvc.UpdateBuyerProfile(c.Context(), userID, body.Name, body.Email); err != nil {
		return response.BadRequest(c, err.Error())
	}
	return response.OK(c, fiber.Map{"message": "profile updated"})
}

// PUT /api/buyer/password
func (h *BuyerHandler) ChangePassword(c *fiber.Ctx) error {
	userID, err := mustUserID(c)
	if err != nil {
		return err
	}
	var body struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}
	if err := c.BodyParser(&body); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	if body.CurrentPassword == "" || body.NewPassword == "" {
		return response.BadRequest(c, "current_password and new_password are required")
	}
	if valid, errKey := validations.ValidatePassword(body.NewPassword); !valid {
		if errKey == "password_min" {
			return response.BadRequest(c, "new password must be at least 8 characters")
		}
		return response.BadRequest(c, "invalid password")
	}
	if err := h.authSvc.ChangePassword(c.Context(), userID, body.CurrentPassword, body.NewPassword); err != nil {
		return response.BadRequest(c, err.Error())
	}
	return response.OK(c, fiber.Map{"message": "password updated"})
}

// GET /api/buyer/orders
func (h *BuyerHandler) GetOrders(c *fiber.Ctx) error {
	userID, err := mustUserID(c)
	if err != nil {
		return err
	}
	orders, err := h.authSvc.GetBuyerOrders(c.Context(), userID)
	if err != nil {
		return response.InternalError(c, err)
	}
	if orders == nil {
		orders = []map[string]interface{}{}
	}
	return response.OK(c, orders)
}

// POST /api/buyer/upgrade-tier — submit tier upgrade request after payment
func (h *BuyerHandler) UpgradeTier(c *fiber.Ctx) error {
	userID, err := mustUserID(c)
	if err != nil {
		return err
	}
	var req dto.BuyerUpgradeRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	if h.adminSvc != nil {
		if val := h.adminSvc.GetSetting(c.Context(), "allow_buyer_tier_upgrade"); val == "false" {
			return response.BadRequest(c, "buyer tier upgrades are currently closed — contact admin")
		}
	}
	if err := h.authSvc.UpgradeBuyerTier(c.Context(), userID, req); err != nil {
		return response.BadRequest(c, err.Error())
	}
	return response.OK(c, fiber.Map{"message": "upgrade request submitted — pending admin approval"})
}

// GET /api/buyer/notifications
func (h *BuyerHandler) ListNotifications(c *fiber.Ctx) error {
	userID, err := mustUserID(c)
	if err != nil {
		return err
	}
	notifs, err := h.authSvc.GetBuyerNotifications(c.Context(), userID)
	if err != nil {
		return response.InternalError(c, err)
	}
	if notifs == nil {
		notifs = []map[string]interface{}{}
	}
	return response.OKList(c, notifs)
}

// PUT /api/buyer/notifications/:id/read
func (h *BuyerHandler) MarkNotificationRead(c *fiber.Ctx) error {
	userID, err := mustUserID(c)
	if err != nil {
		return err
	}
	id, err := parseIDParam(c, "id")
	if err != nil {
		return err
	}
	if err := h.authSvc.MarkBuyerNotificationRead(c.Context(), id, userID); err != nil {
		return response.InternalError(c, err)
	}
	return response.OK(c, fiber.Map{"message": "marked as read"})
}

// GET /api/buyer/upgrade-tier — list buyer's own upgrade requests
func (h *BuyerHandler) GetUpgradeRequests(c *fiber.Ctx) error {
	userID, err := mustUserID(c)
	if err != nil {
		return err
	}
	requests, err := h.authSvc.GetBuyerUpgradeRequestsByUser(c.Context(), userID)
	if err != nil {
		return response.InternalError(c, err)
	}
	if requests == nil {
		requests = []map[string]interface{}{}
	}
	return response.OKList(c, requests)
}

// POST /api/buyer/refunds
func (h *BuyerHandler) CreateRefundRequest(c *fiber.Ctx) error {
	buyerID, err := mustUserID(c)
	if err != nil {
		return err
	}
	var body struct {
		TransactionID *int     `json:"transaction_id"`
		ProductID     *int     `json:"product_id"`
		Amount        float64  `json:"amount"`
		Reason        string   `json:"reason"`
		Description   string   `json:"description"`
		EvidenceURLs  []string `json:"evidence_urls"`
	}
	if err := c.BodyParser(&body); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	if body.Amount <= 0 || body.Reason == "" {
		return response.BadRequest(c, "amount and reason are required")
	}
	req := &models.RefundRequest{
		TransactionID: body.TransactionID,
		ProductID:     body.ProductID,
		RequesterID:   buyerID,
		RequesterRole: "buyer",
		Amount:        body.Amount,
		Reason:        body.Reason,
		Description:   body.Description,
		EvidenceURLs:  body.EvidenceURLs,
	}
	id, err := h.adminSvc.CreateRefundRequest(c.Context(), req)
	if err != nil {
		return response.InternalError(c, err)
	}
	req.ID = id
	req.Status = "pending"
	return response.Created(c, req)
}

// POST /api/reports
func (h *BuyerHandler) CreateAbuseReport(c *fiber.Ctx) error {
	var body struct {
		ReportedUserID    *int     `json:"reported_user_id"`
		ReportedProductID *int     `json:"reported_product_id"`
		Reason            string   `json:"reason"`
		Description       string   `json:"description"`
		EvidenceURLs      []string `json:"evidence_urls"`
		Priority          string   `json:"priority"`
	}
	if err := c.BodyParser(&body); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	if body.Reason == "" || body.Description == "" {
		return response.BadRequest(c, "reason and description are required")
	}
	var reporterID *int
	role := "anonymous"
	if uid, err := mustUserID(c); err == nil && uid > 0 {
		reporterID = &uid
		role = "user"
	}
	rep := &models.AbuseReport{
		ReporterID:        reporterID,
		ReporterRole:      role,
		ReportedUserID:    body.ReportedUserID,
		ReportedProductID: body.ReportedProductID,
		Reason:            body.Reason,
		Description:       body.Description,
		EvidenceURLs:      body.EvidenceURLs,
		Priority:          body.Priority,
	}
	id, err := h.adminSvc.CreateAbuseReport(c.Context(), rep)
	if err != nil {
		return response.InternalError(c, err)
	}
	rep.ID = id
	rep.Status = "pending"
	return response.Created(c, rep)
}

// POST /api/appeals
func (h *BuyerHandler) CreateAccountAppeal(c *fiber.Ctx) error {
	var body struct {
		UserID            int      `json:"user_id"`
		AppealType        string   `json:"appeal_type"`
		RelatedProductID  *int     `json:"related_product_id"`
		RelatedDocumentID *int     `json:"related_document_id"`
		Reason            string   `json:"reason"`
		EvidenceURLs      []string `json:"evidence_urls"`
	}
	if err := c.BodyParser(&body); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	uid := body.UserID
	if uid <= 0 {
		if authUID, err := mustUserID(c); err == nil && authUID > 0 {
			uid = authUID
		}
	}
	if uid <= 0 || body.AppealType == "" || body.Reason == "" {
		return response.BadRequest(c, "user_id, appeal_type, and reason are required")
	}
	ap := &models.AccountAppeal{
		UserID:            uid,
		AppealType:        body.AppealType,
		RelatedProductID:  body.RelatedProductID,
		RelatedDocumentID: body.RelatedDocumentID,
		Reason:            body.Reason,
		EvidenceURLs:      body.EvidenceURLs,
	}
	id, err := h.adminSvc.CreateAccountAppeal(c.Context(), ap)
	if err != nil {
		return response.InternalError(c, err)
	}
	ap.ID = id
	ap.Status = "pending"
	return response.Created(c, ap)
}
