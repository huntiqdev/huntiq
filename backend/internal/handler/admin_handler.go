package handlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/your-org/used-tools-marketplace/internal/models"
	"github.com/your-org/used-tools-marketplace/internal/repository"
	"github.com/your-org/used-tools-marketplace/internal/service"
	"github.com/your-org/used-tools-marketplace/pkg/response"
	"github.com/your-org/used-tools-marketplace/pkg/validations"
	"strings"
)

type AdminHandler struct {
	svc      *service.AdminService
	authSvc  *service.AuthService
	catRepo  repository.CategoryRepository
	prodRepo repository.ProductRepository
}

func NewAdminHandler(svc *service.AdminService, authSvc ...*service.AuthService) *AdminHandler {
	h := &AdminHandler{svc: svc}
	if len(authSvc) > 0 {
		h.authSvc = authSvc[0]
	}
	return h
}

func (h *AdminHandler) WithCategoryRepo(r repository.CategoryRepository) *AdminHandler {
	h.catRepo = r
	return h
}

func (h *AdminHandler) WithProductRepo(r repository.ProductRepository) *AdminHandler {
	h.prodRepo = r
	return h
}

func adminRole(c *fiber.Ctx) string {
	r, _ := c.Locals("adminRole").(string)
	return r
}

// PUT /api/admin/password
func (h *AdminHandler) ChangePassword(c *fiber.Ctx) error {
	adminID, err := mustUserID(c)
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
	if h.authSvc == nil {
		return response.InternalError(c, fmt.Errorf("password change not configured"))
	}
	if err := h.authSvc.ChangePassword(c.Context(), adminID, body.CurrentPassword, body.NewPassword); err != nil {
		return response.BadRequest(c, err.Error())
	}
	return response.OK(c, fiber.Map{"message": "password updated"})
}

// GET /api/admin/products/pending
func (h *AdminHandler) ListPending(c *fiber.Ctx) error {
	products, err := h.svc.ListPending(c.Context())
	if err != nil {
		return response.InternalError(c, err)
	}
	return response.OKList(c, products)
}

// GET /api/admin/products
func (h *AdminHandler) ListAllProducts(c *fiber.Ctx) error {
	page, limit := parsePagination(c, 50, 200)
	var products []models.Product
	var total int
	var err error
	if status := c.Query("status"); status != "" {
		validStatuses := map[string]bool{
			"pending": true, "approved": true, "rejected": true,
			"sold": true, "expired": true,
		}
		if !validStatuses[status] {
			return response.BadRequest(c, "invalid status — must be pending, approved, rejected, sold, or expired")
		}
		products, total, err = h.svc.ListProductsByStatus(c.Context(), status, page, limit)
	} else {
		products, total, err = h.svc.ListAllProducts(c.Context(), page, limit)
	}
	if err != nil {
		return response.InternalError(c, err)
	}
	if products == nil {
		products = []models.Product{}
	}
	return response.OK(c, fiber.Map{"data": products, "total": total, "page": page, "limit": limit})
}

// PUT /api/admin/products/:id/review
func (h *AdminHandler) ReviewProduct(c *fiber.Ctx) error {
	adminID, err := mustUserID(c)
	if err != nil {
		return err
	}
	productID, err := parseIDParam(c, "id")
	if err != nil {
		return err
	}
	var body struct {
		Status string `json:"status"`
		Notes  string `json:"notes"`
	}
	if err := c.BodyParser(&body); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	if err := h.svc.ReviewProduct(c.Context(), productID, adminID, body.Status, body.Notes); err != nil {
		return response.BadRequest(c, err.Error())
	}
	action := "approve"
	if body.Status == "rejected" {
		action = "reject"
	}
	h.svc.InsertAdminLog(c.Context(), models.AdminLog{
		AdminID:   adminID,
		ProductID: &productID,
		Action:    action,
		Notes:     body.Notes,
		IPAddress: c.IP(),
	})
	return response.OK(c, fiber.Map{"message": "product reviewed"})
}

// DELETE /api/admin/products/:id
func (h *AdminHandler) DeleteProduct(c *fiber.Ctx) error {
	adminID, err := mustUserID(c)
	if err != nil {
		return err
	}
	id, err := parseIDParam(c, "id")
	if err != nil {
		return err
	}
	if err := h.svc.DeleteProduct(c.Context(), id); err != nil {
		return response.InternalError(c, err)
	}
	h.svc.InsertAdminLog(c.Context(), models.AdminLog{
		AdminID:   adminID,
		ProductID: &id,
		Action:    "delete",
		IPAddress: c.IP(),
	})
	return response.OK(c, fiber.Map{"message": "product deleted"})
}

// GET /api/admin/sellers
func (h *AdminHandler) ListSellers(c *fiber.Ctx) error {
	page, limit := parsePagination(c, 50, 200)
	search := strings.TrimSpace(c.Query("search"))
	status := strings.TrimSpace(c.Query("status"))
	tier := strings.TrimSpace(c.Query("tier"))
	strategy := strings.TrimSpace(c.Query("strategy"))
	if strategy == "" {
		strategy = strings.TrimSpace(c.Query("preference"))
	}
	sellers, total, err := h.svc.ListSellersFiltered(c.Context(), page, limit, search, status, tier, strategy)
	if err != nil {
		return response.InternalError(c, err)
	}
	if sellers == nil {
		sellers = []models.User{}
	}
	// Hide sensitive profit data for non-platform owner
	if adminRole(c) != models.RolePlatformOwner && adminRole(c) != models.LegacyRoleSuper {
		for i := range sellers {
			sellers[i].NetProfit = 0
		}
	}
	return response.OK(c, fiber.Map{"data": sellers, "total": total, "page": page, "limit": limit})
}

// PUT /api/admin/sellers/:id/active
func (h *AdminHandler) SetSellerActive(c *fiber.Ctx) error {
	adminID, err := mustUserID(c)
	if err != nil {
		return err
	}
	id, err := parseIDParam(c, "id")
	if err != nil {
		return err
	}
	var body struct {
		Active bool   `json:"active"`
		Reason string `json:"reason"`
	}
	if err := c.BodyParser(&body); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	// Reason is required when deactivating
	if !body.Active && strings.TrimSpace(body.Reason) == "" {
		return response.BadRequest(c, "reason is required when deactivating a seller")
	}
	if err := h.svc.SetSellerActive(c.Context(), id, body.Active, strings.TrimSpace(body.Reason)); err != nil {
		return response.BadRequest(c, err.Error())
	}
	action := "activate"
	if !body.Active {
		action = "deactivate"
	}
	h.svc.InsertAdminLog(c.Context(), models.AdminLog{
		AdminID:      adminID,
		TargetUserID: &id,
		Action:       action,
		Notes:        strings.TrimSpace(body.Reason),
		IPAddress:    c.IP(),
	})
	return response.OK(c, fiber.Map{"message": "seller updated"})
}

// GET /api/admin/transactions  — permission-gated
func (h *AdminHandler) ListTransactions(c *fiber.Ctx) error {
	page, limit := parsePagination(c, 50, 200)
	search := strings.TrimSpace(c.Query("search"))
	status := strings.TrimSpace(c.Query("status"))
	method := strings.TrimSpace(c.Query("method"))
	txType := strings.TrimSpace(c.Query("type"))
	fromDate := strings.TrimSpace(c.Query("from_date"))
	toDate := strings.TrimSpace(c.Query("to_date"))

	txns, total, err := h.svc.ListAllTransactionsFiltered(c.Context(), page, limit, search, status, method, txType, fromDate, toDate)
	if err != nil {
		return response.InternalError(c, err)
	}
	if txns == nil {
		txns = []models.Transaction{}
	}
	return response.OK(c, fiber.Map{"data": txns, "total": total, "page": page, "limit": limit})
}

// GET /api/admin/locations
func (h *AdminHandler) ListLocations(c *fiber.Ctx) error {
	locs, err := h.svc.ListLocations(c.Context())
	if err != nil {
		return response.InternalError(c, err)
	}
	return response.OKList(c, locs)
}

// POST /api/admin/locations
func (h *AdminHandler) AddLocation(c *fiber.Ctx) error {
	var l models.Location
	if err := c.BodyParser(&l); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	id, err := h.svc.AddLocation(c.Context(), l)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}
	l.ID = id
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"success": true, "data": l})
}

// PUT /api/admin/locations/:id
func (h *AdminHandler) UpdateLocation(c *fiber.Ctx) error {
	id, err := parseIDParam(c, "id")
	if err != nil {
		return err
	}
	var l models.Location
	if err := c.BodyParser(&l); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	l.ID = id
	if err := h.svc.UpdateLocation(c.Context(), l); err != nil {
		return response.BadRequest(c, err.Error())
	}
	return response.OK(c, fiber.Map{"message": "location updated"})
}

// DELETE /api/admin/locations/:id
func (h *AdminHandler) DeleteLocation(c *fiber.Ctx) error {
	id, err := parseIDParam(c, "id")
	if err != nil {
		return err
	}
	if err := h.svc.DeleteLocation(c.Context(), id); err != nil {
		return response.InternalError(c, err)
	}
	return response.OK(c, fiber.Map{"message": "location deleted"})
}

// ── Analytics ─────────────────────────────────────────────────────────────────

// GET /api/admin/analytics/overview
func (h *AdminHandler) ListCategories(c *fiber.Ctx) error {
	if h.catRepo == nil {
		return response.InternalError(c, fiber.NewError(500, "category repo not configured"))
	}
	cats, err := h.catRepo.ListAll(c.Context())
	if err != nil {
		return response.InternalError(c, err)
	}
	return response.OKList(c, cats)
}

// GET /api/admin/me — returns the current admin's profile
// Uses authSvc.GetProfile (single GetByID query) — avoids loading all admins.
