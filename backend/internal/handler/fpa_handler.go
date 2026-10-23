package handlers

import (
	"encoding/json"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/your-org/used-tools-marketplace/internal/service"
	"github.com/your-org/used-tools-marketplace/pkg/response"
)

// FPAHandler handles incoming HTTP requests for Financial Planning & Analysis
type FPAHandler struct {
	svc *service.FPAService
}

// NewFPAHandler creates a new FPAHandler instance
func NewFPAHandler(svc *service.FPAService) *FPAHandler {
	return &FPAHandler{svc: svc}
}

// GetForecast returns the full forecast report for a scenario and horizon
// GET /api/admin/finance/forecast?scenario=base&horizon=60
func (h *FPAHandler) GetForecast(c *fiber.Ctx) error {
	scenario := c.Query("scenario", "base")
	horizonStr := c.Query("horizon", "60")
	horizon, err := strconv.Atoi(horizonStr)
	if err != nil || horizon <= 0 {
		horizon = 60
	}

	report, err := h.svc.GetForecast(c.Context(), scenario, horizon)
	if err != nil {
		return response.InternalError(c, err)
	}
	return response.OK(c, report)
}

// GetScenarios returns Base, Conservative, and High-Growth scenarios side-by-side
// GET /api/admin/finance/scenarios?horizon=60
func (h *FPAHandler) GetScenarios(c *fiber.Ctx) error {
	horizonStr := c.Query("horizon", "60")
	horizon, err := strconv.Atoi(horizonStr)
	if err != nil || horizon <= 0 {
		horizon = 60
	}

	comparison, err := h.svc.GetScenarioComparison(c.Context(), horizon)
	if err != nil {
		return response.InternalError(c, err)
	}
	return response.OK(c, comparison)
}

// GetCapacityPlanning20M returns technical demand and bottleneck analysis for the 20M DAU scenario
// GET /api/admin/finance/capacity-20m
func (h *FPAHandler) GetCapacityPlanning20M(c *fiber.Ctx) error {
	res, err := h.svc.GetCapacityPlanning20M(c.Context())
	if err != nil {
		return response.InternalError(c, err)
	}
	return response.OK(c, res)
}

// GetTaxRules returns official Ethiopian statutory tax rules with source citations
// GET /api/admin/finance/tax-rules
func (h *FPAHandler) GetTaxRules(c *fiber.Ctx) error {
	rules, err := h.svc.GetTaxRules(c.Context())
	if err != nil {
		return response.InternalError(c, err)
	}
	return response.OKList(c, rules)
}

// GetActuals returns historical actuals benchmark records
// GET /api/admin/finance/actuals
func (h *FPAHandler) GetActuals(c *fiber.Ctx) error {
	actuals, err := h.svc.GetActuals(c.Context())
	if err != nil {
		return response.InternalError(c, err)
	}
	return response.OKList(c, actuals)
}

// GetAuditLogs returns financial assumption change logs
// GET /api/admin/finance/audit-logs
func (h *FPAHandler) GetAuditLogs(c *fiber.Ctx) error {
	limitStr := c.Query("limit", "50")
	limit, _ := strconv.Atoi(limitStr)
	logs, err := h.svc.GetAuditLogs(c.Context(), limit)
	if err != nil {
		return response.InternalError(c, err)
	}
	return response.OKList(c, logs)
}

type updateAssumptionsReq struct {
	ScenarioID string          `json:"scenario_id"`
	Section    string          `json:"section"`
	Payload    json.RawMessage `json:"payload"`
}

// UpdateAssumptions updates a specific assumption category
// PUT /api/admin/finance/assumptions
func (h *FPAHandler) UpdateAssumptions(c *fiber.Ctx) error {
	adminID := c.Locals("userID").(int)
	ip := c.IP()

	var req updateAssumptionsReq
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request payload: "+err.Error())
	}

	if req.ScenarioID == "" || req.Section == "" || len(req.Payload) == 0 {
		return response.BadRequest(c, "scenario_id, section, and payload are required")
	}

	if err := h.svc.UpdateAssumptionsSection(c.Context(), adminID, req.ScenarioID, req.Section, req.Payload, ip); err != nil {
		return response.InternalError(c, err)
	}

	return response.OK(c, map[string]interface{}{
		"success": true,
		"message": "Financial assumptions updated successfully",
	})
}

type resetAssumptionsReq struct {
	ScenarioID string `json:"scenario_id"`
}

// ResetAssumptions reverts custom overrides to system defaults
// POST /api/admin/finance/assumptions/reset
func (h *FPAHandler) ResetAssumptions(c *fiber.Ctx) error {
	adminID := c.Locals("userID").(int)
	ip := c.IP()

	var req resetAssumptionsReq
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	if req.ScenarioID == "" {
		req.ScenarioID = "base"
	}

	if err := h.svc.ResetAssumptions(c.Context(), adminID, req.ScenarioID, ip); err != nil {
		return response.InternalError(c, err)
	}

	return response.OK(c, map[string]interface{}{
		"success": true,
		"message": "Financial assumptions reset to system defaults",
	})
}

// ExportForecast generates a downloadable CSV or JSON export
// GET /api/admin/finance/export?scenario=base&horizon=60&format=csv
func (h *FPAHandler) ExportForecast(c *fiber.Ctx) error {
	scenario := c.Query("scenario", "base")
	horizonStr := c.Query("horizon", "60")
	horizon, _ := strconv.Atoi(horizonStr)
	if horizon <= 0 {
		horizon = 60
	}
	format := c.Query("format", "csv")

	if format == "json" {
		rep, err := h.svc.GetForecast(c.Context(), scenario, horizon)
		if err != nil {
			return response.InternalError(c, err)
		}
		c.Set("Content-Type", "application/json")
		c.Set("Content-Disposition", `attachment; filename="huntiq_financial_forecast.json"`)
		return c.JSON(rep)
	}

	csvContent, err := h.svc.ExportForecastCSV(c.Context(), scenario, horizon)
	if err != nil {
		return response.InternalError(c, err)
	}

	c.Set("Content-Type", "text/csv; charset=utf-8")
	c.Set("Content-Disposition", `attachment; filename="huntiq_financial_forecast.csv"`)
	return c.SendString(csvContent)
}
