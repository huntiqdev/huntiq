package handlers

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/your-org/used-tools-marketplace/internal/service"
	"github.com/your-org/used-tools-marketplace/pkg/response"
)

type WalletHandler struct {
	ledgerSvc   *service.LedgerService
	tokenSvc    *service.ReferralTokenService
	frontendURL string
}

func NewWalletHandler(ledgerSvc *service.LedgerService, tokenSvc *service.ReferralTokenService, frontendURL string) *WalletHandler {
	if frontendURL == "" {
		frontendURL = "https://huntiq.et"
	}
	return &WalletHandler{
		ledgerSvc:   ledgerSvc,
		tokenSvc:    tokenSvc,
		frontendURL: frontendURL,
	}
}

// GET /api/wallet
func (h *WalletHandler) GetMyWallet(c *fiber.Ctx) error {
	userID, err := mustUserID(c)
	if err != nil {
		return err
	}

	wallet, err := h.ledgerSvc.GetWallet(c.Context(), userID)
	if err != nil {
		return response.InternalError(c, err)
	}

	return response.OK(c, wallet)
}

// GET /api/wallet/ledger
func (h *WalletHandler) GetMyLedger(c *fiber.Ctx) error {
	userID, err := mustUserID(c)
	if err != nil {
		return err
	}

	page, limit := parsePagination(c, 20, 100)
	offset := (page - 1) * limit

	txs, total, err := h.ledgerSvc.GetLedgerHistory(c.Context(), userID, limit, offset)
	if err != nil {
		return response.InternalError(c, err)
	}

	return response.OK(c, fiber.Map{
		"transactions": txs,
		"total":        total,
	})
}

// GET /api/referrals/token-link
func (h *WalletHandler) GetMyReferralLink(c *fiber.Ctx) error {
	userID, err := mustUserID(c)
	if err != nil {
		return err
	}

	referralCode := c.Query("code", "")
	token, url, err := h.tokenSvc.GetUserReferralLink(c.Context(), userID, referralCode, h.frontendURL)
	if err != nil {
		return response.InternalError(c, err)
	}

	return response.OK(c, fiber.Map{
		"token":        token,
		"referral_url": url,
	})
}

// POST /api/referrals/click (public click tracking)
func (h *WalletHandler) TrackReferralClick(c *fiber.Ctx) error {
	var body struct {
		Token      string `json:"token"`
		LandingURL string `json:"landing_url"`
	}
	if err := c.BodyParser(&body); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	if strings.TrimSpace(body.Token) == "" {
		return response.BadRequest(c, "referral token is required")
	}

	ip := c.IP()
	ua := c.Get("User-Agent")

	tokenObj, err := h.tokenSvc.RecordClick(c.Context(), body.Token, ip, ua, body.LandingURL)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.OK(c, fiber.Map{
		"valid":   true,
		"token":   tokenObj.Token,
		"inviter": tokenObj.UserID,
	})
}
