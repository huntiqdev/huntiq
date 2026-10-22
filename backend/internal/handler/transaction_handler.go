package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/your-org/used-tools-marketplace/pkg/dto"
	"github.com/your-org/used-tools-marketplace/internal/service"
	"github.com/your-org/used-tools-marketplace/pkg/response"
)

type TransactionHandler struct {
	svc *service.TransactionService
}

func NewTransactionHandler(svc *service.TransactionService) *TransactionHandler {
	return &TransactionHandler{svc: svc}
}

// POST /api/transactions
func (h *TransactionHandler) Create(c *fiber.Ctx) error {
	var req dto.CreateTransactionRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	if req.ProductID <= 0 {
		return response.BadRequest(c, "product_id is required")
	}
	if req.BuyerName == "" || req.BuyerPhone == "" {
		return response.BadRequest(c, "buyer_name and buyer_phone are required")
	}
	if req.Amount <= 0 {
		return response.BadRequest(c, "amount must be greater than zero")
	}
	txn, err := h.svc.Create(c.Context(), req)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}
	return response.OK(c, txn)
}

// PUT /api/admin/transactions/:id/status  — admin marks transaction completed or cancelled
func (h *TransactionHandler) UpdateStatus(c *fiber.Ctx) error {
	id, err := parseIDParam(c, "id")
	if err != nil {
		return err
	}
	var body struct {
		Status string `json:"status"`
	}
	if err := c.BodyParser(&body); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	if err := h.svc.UpdateStatus(c.Context(), id, body.Status); err != nil {
		return response.BadRequest(c, err.Error())
	}
	return response.OK(c, fiber.Map{"message": "transaction status updated"})
}
