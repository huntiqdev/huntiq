package handlers

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/your-org/used-tools-marketplace/internal/service"
	"github.com/your-org/used-tools-marketplace/pkg/response"
)

type TelebirrHandler struct {
	svc *service.TelebirrService
}

func NewTelebirrHandler(svc *service.TelebirrService) *TelebirrHandler {
	return &TelebirrHandler{svc: svc}
}

// POST /api/telebirr/orders
func (h *TelebirrHandler) CreateOrder(c *fiber.Ctx) error {
	userID, err := mustUserID(c)
	if err != nil {
		return err
	}

	var body struct {
		Amount    float64 `json:"amount"`
		OrderType string  `json:"order_type"`
		Title     string  `json:"title"`
	}
	if err := c.BodyParser(&body); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	if body.Amount <= 0 {
		return response.BadRequest(c, "amount must be greater than zero")
	}
	if strings.TrimSpace(body.OrderType) == "" {
		body.OrderType = "listing_fee"
	}
	if strings.TrimSpace(body.Title) == "" {
		body.Title = "HUNTIQ Marketplace Payment"
	}

	order, checkoutURL, err := h.svc.CreatePaymentOrder(c.Context(), userID, body.Amount, body.OrderType, body.Title)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.OK(c, fiber.Map{
		"order":        order,
		"checkout_url": checkoutURL,
	})
}

// POST /api/telebirr/webhook
func (h *TelebirrHandler) HandleWebhook(c *fiber.Ctx) error {
	sig := c.Get("X-Telebirr-Signature")
	if sig == "" {
		sig = c.Get("Signature")
	}

	rawBody := c.Body()
	order, err := h.svc.HandleWebhook(c.Context(), rawBody, sig)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    1,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":         0,
		"message":      "success",
		"out_trade_no": order.OutTradeNo,
		"status":       order.Status,
	})
}

// GET /api/telebirr/orders/:outTradeNo
func (h *TelebirrHandler) GetOrderStatus(c *fiber.Ctx) error {
	outTradeNo := c.Params("outTradeNo")
	if strings.TrimSpace(outTradeNo) == "" {
		return response.BadRequest(c, "out_trade_no is required")
	}

	order, err := h.svc.GetOrderStatus(c.Context(), outTradeNo)
	if err != nil {
		return response.NotFound(c, "order not found")
	}

	return response.OK(c, order)
}

// GET /api/telebirr/my-orders
func (h *TelebirrHandler) GetUserOrders(c *fiber.Ctx) error {
	userID, err := mustUserID(c)
	if err != nil {
		return err
	}

	page, limit := parsePagination(c, 20, 100)
	offset := (page - 1) * limit

	orders, total, err := h.svc.GetUserOrders(c.Context(), userID, limit, offset)
	if err != nil {
		return response.InternalError(c, err)
	}

	return response.OK(c, fiber.Map{
		"orders": orders,
		"total":  total,
	})
}
