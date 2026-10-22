package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/your-org/used-tools-marketplace/pkg/dto"
	"github.com/your-org/used-tools-marketplace/internal/models"
	"github.com/your-org/used-tools-marketplace/internal/service"
	"github.com/your-org/used-tools-marketplace/pkg/response"
	"github.com/your-org/used-tools-marketplace/pkg/validations"
)

type SellerHandler struct {
	svc      *service.SellerService
	adminSvc *service.AdminService
}

func NewSellerHandler(svc *service.SellerService) *SellerHandler {
	return &SellerHandler{svc: svc}
}

func (h *SellerHandler) WithAdminService(a *service.AdminService) *SellerHandler {
	h.adminSvc = a
	return h
}

// POST /api/seller/listing-fee  — seller records their listing fee payment
func (h *SellerHandler) PayListingFee(c *fiber.Ctx) error {
	sellerID, err := mustUserID(c)
	if err != nil {
		return err
	}
	var req dto.ListingFeePaymentRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	if req.ProductID <= 0 {
		return response.BadRequest(c, "product_id is required")
	}
	if req.Amount <= 0 {
		return response.BadRequest(c, "amount must be greater than zero")
	}
	if req.Reference == "" {
		return response.BadRequest(c, "payment reference is required")
	}
	if req.PayerPhone == "" {
		return response.BadRequest(c, "payer_phone is required")
	}
	validMethods := map[string]bool{"telebirr": true, "safaricom": true, "cbe": true}
	if !validMethods[req.PaymentMethod] {
		return response.BadRequest(c, "payment_method must be telebirr, safaricom, or cbe")
	}
	if err := h.svc.RecordListingFeePayment(c.Context(), sellerID, req); err != nil {
		return response.BadRequest(c, err.Error())
	}
	return response.OK(c, fiber.Map{"message": "payment recorded — product is pending review"})
}

// GET /api/seller/stats
func (h *SellerHandler) GetStats(c *fiber.Ctx) error {
	sellerID, err := mustUserID(c)
	if err != nil {
		return err
	}
	stats, err := h.svc.GetStats(c.Context(), sellerID)
	if err != nil {
		return response.InternalError(c, err)
	}
	return response.OK(c, stats)
}

// GET /api/seller/profile
func (h *SellerHandler) GetProfile(c *fiber.Ctx) error {
	sellerID, err := mustUserID(c)
	if err != nil {
		return err
	}
	u, err := h.svc.GetProfile(c.Context(), sellerID)
	if err != nil {
		return response.NotFound(c, "profile not found")
	}
	return response.OK(c, u)
}

// PUT /api/seller/profile
func (h *SellerHandler) UpdateProfile(c *fiber.Ctx) error {
	sellerID, err := mustUserID(c)
	if err != nil {
		return err
	}
	// Only accept safe profile fields — never bind to full models.User to prevent role/tier injection
	var body struct {
		Name             string   `json:"name"`
		BusinessName     string   `json:"business_name"`
		Email            string   `json:"email"`
		AdditionalPhones []string `json:"additional_phones"`
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

	// Get current user to check tier
	profile, err := h.svc.GetProfile(c.Context(), sellerID)
	if err != nil {
		return response.InternalError(c, err)
	}

	u := &models.User{
		ID:               sellerID,
		Name:             body.Name,
		BusinessName:     body.BusinessName,
		Email:            body.Email,
		AdditionalPhones: body.AdditionalPhones,
		Tier:             profile.Tier, // needed for validation
	}
	if err := h.svc.UpdateProfile(c.Context(), u); err != nil {
		return response.BadRequest(c, err.Error())
	}
	return response.OK(c, fiber.Map{"message": "profile updated"})
}

// PUT /api/seller/preference
func (h *SellerHandler) UpdatePreference(c *fiber.Ctx) error {
	sellerID, err := mustUserID(c)
	if err != nil {
		return err
	}
	var body struct {
		SellerPreference string `json:"seller_preference"`
	}
	if err := c.BodyParser(&body); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	
	validPrefs := map[string]bool{"sale": true, "rent": true, "both": true}
	if !validPrefs[body.SellerPreference] {
		return response.BadRequest(c, "seller_preference must be 'sale', 'rent', or 'both'")
	}
	
	if err := h.svc.UpdatePreference(c.Context(), sellerID, body.SellerPreference); err != nil {
		return response.InternalError(c, err)
	}
	
	return response.OK(c, fiber.Map{"message": "preference updated"})
}

// GET /api/seller/transactions
func (h *SellerHandler) ListTransactions(c *fiber.Ctx) error {
	sellerID, err := mustUserID(c)
	if err != nil {
		return err
	}
	txns, err := h.svc.ListTransactions(c.Context(), sellerID)
	if err != nil {
		return response.InternalError(c, err)
	}
	return response.OKList(c, txns)
}

// GET /api/seller/products/detailed  — includes fee + expiry fields
func (h *SellerHandler) ListProductsDetailed(c *fiber.Ctx) error {
	sellerID, err := mustUserID(c)
	if err != nil {
		return err
	}
	products, err := h.svc.ListProductsWithFee(c.Context(), sellerID)
	if err != nil {
		return response.InternalError(c, err)
	}
	return response.OKList(c, products)
}

// DELETE /api/seller/products/:id
func (h *SellerHandler) DeleteProduct(c *fiber.Ctx) error {
	sellerID, err := mustUserID(c)
	if err != nil {
		return err
	}
	id, err := parseIDParam(c, "id")
	if err != nil {
		return err
	}
	if err := h.svc.DeleteProduct(c.Context(), id, sellerID); err != nil {
		return response.BadRequest(c, err.Error())
	}
	return response.OK(c, fiber.Map{"message": "product deleted"})
}

// ── Order Requests ────────────────────────────────────────────────────────────

// POST /api/orders  — public, any buyer can submit
func (h *SellerHandler) CreateOrder(c *fiber.Ctx) error {
	var req dto.CreateOrderRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	// Inject buyer tier from JWT if authenticated, else default to normal
	if bt, ok := c.Locals("buyerTier").(string); ok && bt != "" {
		req.BuyerTier = bt
	} else {
		req.BuyerTier = "normal"
	}
	order, err := h.svc.CreateOrderRequest(c.Context(), req)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"success": true, "data": order})
}

// GET /api/seller/orders
func (h *SellerHandler) ListOrders(c *fiber.Ctx) error {
	sellerID, err := mustUserID(c)
	if err != nil {
		return err
	}
	orders, err := h.svc.ListOrders(c.Context(), sellerID)
	if err != nil {
		return response.InternalError(c, err)
	}
	return response.OKList(c, orders)
}

// PUT /api/seller/orders/:id
func (h *SellerHandler) UpdateOrder(c *fiber.Ctx) error {
	sellerID, err := mustUserID(c)
	if err != nil {
		return err
	}
	id, err := parseIDParam(c, "id")
	if err != nil {
		return err
	}
	var req dto.UpdateOrderRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	if err := h.svc.UpdateOrderStatus(c.Context(), id, sellerID, req); err != nil {
		return response.BadRequest(c, err.Error())
	}
	return response.OK(c, fiber.Map{"message": "order updated"})
}

// ── Tier Upgrade Requests ─────────────────────────────────────────────────────

// POST /api/seller/upgrade-tier
func (h *SellerHandler) RequestUpgrade(c *fiber.Ctx) error {
	sellerID, err := mustUserID(c)
	if err != nil {
		return err
	}
	var body struct {
		RequestedTier string  `json:"requested_tier"`
		PaymentMethod string  `json:"payment_method"`
		PayerPhone    string  `json:"payer_phone"`
		Reference     string  `json:"reference"`
		Amount        float64 `json:"amount"`
		Notes         string  `json:"notes"`
	}
	if err := c.BodyParser(&body); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	// Validate tier before using it to build a settings key
	if body.RequestedTier != "gold" && body.RequestedTier != "premium" {
		return response.BadRequest(c, "requested_tier must be gold or premium")
	}
	// Check platform setting — admin can close upgrade path per tier
	if h.adminSvc != nil {
		settingKey := "allow_seller_upgrade_to_" + body.RequestedTier
		val := h.adminSvc.GetSetting(c.Context(), settingKey)
		if val == "false" {
			return response.BadRequest(c, body.RequestedTier+" tier upgrades are currently closed — contact admin")
		}
	}
	id, err := h.svc.SubmitUpgradeRequest(c.Context(), sellerID,
		body.RequestedTier, body.PaymentMethod, body.PayerPhone,
		body.Reference, body.Notes, body.Amount)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    fiber.Map{"id": id, "message": "Upgrade request submitted — admin will review and activate your tier within 24 hours"},
	})
}

// GET /api/seller/upgrade-tier
func (h *SellerHandler) GetUpgradeRequests(c *fiber.Ctx) error {
	sellerID, err := mustUserID(c)
	if err != nil {
		return err
	}
	requests, err := h.svc.GetUpgradeRequests(c.Context(), sellerID)
	if err != nil {
		return response.InternalError(c, err)
	}
	return response.OKList(c, requests)
}

// GET /api/seller/notifications
func (h *SellerHandler) ListNotifications(c *fiber.Ctx) error {
	sellerID, err := mustUserID(c)
	if err != nil {
		return err
	}
	notifs, err := h.svc.ListNotifications(c.Context(), sellerID)
	if err != nil {
		return response.InternalError(c, err)
	}
	if notifs == nil {
		notifs = []map[string]interface{}{}
	}
	return response.OKList(c, notifs)
}

// PUT /api/seller/notifications/:id/read
func (h *SellerHandler) MarkNotificationRead(c *fiber.Ctx) error {
	sellerID, err := mustUserID(c)
	if err != nil {
		return err
	}
	id, err := parseIDParam(c, "id")
	if err != nil {
		return err
	}
	if err := h.svc.MarkNotificationRead(c.Context(), id, sellerID); err != nil {
		return response.InternalError(c, err)
	}
	return response.OK(c, fiber.Map{"message": "marked as read"})
}

// GET /api/seller/verification/documents
func (h *SellerHandler) ListVerificationDocuments(c *fiber.Ctx) error {
	sellerID, err := mustUserID(c)
	if err != nil {
		return err
	}
	if h.adminSvc == nil {
		return response.OK(c, []interface{}{})
	}
	docs, err := h.adminSvc.ListSellerVerificationDocuments(c.Context(), sellerID)
	if err != nil {
		return response.InternalError(c, err)
	}
	return response.OK(c, docs)
}

// POST /api/seller/verification/documents
func (h *SellerHandler) UploadVerificationDocument(c *fiber.Ctx) error {
	sellerID, err := mustUserID(c)
	if err != nil {
		return err
	}
	var body struct {
		DocumentType  string `json:"document_type"`
		FileURL       string `json:"file_url"`
		FileName      string `json:"file_name"`
		FileSizeBytes int    `json:"file_size_bytes"`
	}
	if err := c.BodyParser(&body); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	if body.DocumentType == "" || body.FileURL == "" {
		return response.BadRequest(c, "document_type and file_url are required")
	}
	doc := &models.SellerVerificationDocument{
		SellerID:      sellerID,
		DocumentType:  body.DocumentType,
		FileURL:       body.FileURL,
		FileName:      body.FileName,
		FileSizeBytes: body.FileSizeBytes,
	}
	id, err := h.adminSvc.CreateVerificationDocument(c.Context(), doc)
	if err != nil {
		return response.InternalError(c, err)
	}
	doc.ID = id
	doc.Status = "pending"
	return response.Created(c, doc)
}

// POST /api/seller/refunds
func (h *SellerHandler) CreateRefundRequest(c *fiber.Ctx) error {
	sellerID, err := mustUserID(c)
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
		RequesterID:   sellerID,
		RequesterRole: "seller",
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
