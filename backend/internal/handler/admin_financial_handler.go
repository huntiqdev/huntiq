package handlers

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/your-org/used-tools-marketplace/internal/models"
	"github.com/your-org/used-tools-marketplace/internal/repository"
	"github.com/your-org/used-tools-marketplace/internal/service"
	"github.com/your-org/used-tools-marketplace/pkg/response"
)

type AdminFinancialHandler struct {
	ledgerSvc *service.LedgerService
	smsRepo   *repository.SMSRepo
}

func NewAdminFinancialHandler(ledgerSvc *service.LedgerService, smsRepo *repository.SMSRepo) *AdminFinancialHandler {
	return &AdminFinancialHandler{
		ledgerSvc: ledgerSvc,
		smsRepo:   smsRepo,
	}
}

// GET /api/admin/financial/ledger
func (h *AdminFinancialHandler) GetLedgerAuditTrail(c *fiber.Ctx) error {
	page, limit := parsePagination(c, 50, 200)
	offset := (page - 1) * limit

	txs, total, err := h.ledgerSvc.GetAllLedgerHistory(c.Context(), limit, offset)
	if err != nil {
		return response.InternalError(c, err)
	}

	return response.OK(c, fiber.Map{
		"transactions": txs,
		"total":        total,
	})
}

// GET /api/admin/referral/policy
func (h *AdminFinancialHandler) GetRewardPolicy(c *fiber.Ctx) error {
	policy, err := h.ledgerSvc.GetRewardPolicy(c.Context())
	if err != nil {
		return response.InternalError(c, err)
	}
	return response.OK(c, policy)
}

// PUT /api/admin/referral/policy
func (h *AdminFinancialHandler) UpdateRewardPolicy(c *fiber.Ctx) error {
	adminID, err := mustUserID(c)
	if err != nil {
		return err
	}

	var body struct {
		ID                       int     `json:"id"`
		RegistrationRewardAmount float64 `json:"registration_reward_amount"`
		TransactionRewardPercent float64 `json:"transaction_reward_percent"`
		MinQualifyingTxAmount    float64 `json:"min_qualifying_tx_amount"`
		MaxRewardPerTx           float64 `json:"max_reward_per_tx"`
		RewardHoldingDays        int     `json:"reward_holding_days"`
		Reason                   string  `json:"reason"`
	}
	if err := c.BodyParser(&body); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	if strings.TrimSpace(body.Reason) == "" {
		return response.BadRequest(c, "audit reason is required for policy update")
	}

	policy := &models.ReferralRewardPolicy{
		ID:                       body.ID,
		RegistrationRewardAmount: body.RegistrationRewardAmount,
		TransactionRewardPercent: body.TransactionRewardPercent,
		MinQualifyingTxAmount:    body.MinQualifyingTxAmount,
		MaxRewardPerTx:           body.MaxRewardPerTx,
		RewardHoldingDays:        body.RewardHoldingDays,
	}

	if err := h.ledgerSvc.UpdateRewardPolicy(c.Context(), policy, adminID, body.Reason); err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.OK(c, fiber.Map{
		"message": "policy updated and version snapshot recorded",
	})
}

// GET /api/admin/sms/logs
func (h *AdminFinancialHandler) GetSMSDeliveryLogs(c *fiber.Ctx) error {
	page, limit := parsePagination(c, 50, 200)
	offset := (page - 1) * limit

	logs, total, err := h.smsRepo.GetDeliveryLogs(c.Context(), limit, offset)
	if err != nil {
		return response.InternalError(c, err)
	}

	return response.OK(c, fiber.Map{
		"logs":  logs,
		"total": total,
	})
}

// POST /api/admin/financial/compensate
func (h *AdminFinancialHandler) CompensateLedger(c *fiber.Ctx) error {
	adminID, err := mustUserID(c)
	if err != nil {
		return err
	}

	var body struct {
		UserID   int     `json:"user_id"`
		Amount   float64 `json:"amount"`
		IsCredit bool    `json:"is_credit"`
		Reason   string  `json:"reason"`
	}
	if err := c.BodyParser(&body); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	if body.UserID <= 0 || body.Amount <= 0 || strings.TrimSpace(body.Reason) == "" {
		return response.BadRequest(c, "user_id, positive amount, and reason are required")
	}

	if err := h.ledgerSvc.CompensateBalance(c.Context(), body.UserID, body.Amount, body.IsCredit, body.Reason, adminID); err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.OK(c, fiber.Map{
		"message": "compensating ledger entry recorded successfully",
	})
}

