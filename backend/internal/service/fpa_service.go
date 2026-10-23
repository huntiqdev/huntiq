package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/your-org/used-tools-marketplace/internal/fpa"
	"github.com/your-org/used-tools-marketplace/internal/repository"
)

// FPAService coordinates business financial forecasting, persistence, and audit logging
type FPAService struct {
	repo *repository.FPARepo
}

// NewFPAService creates an instance of FPAService
func NewFPAService(repo *repository.FPARepo) *FPAService {
	return &FPAService{repo: repo}
}

// GetForecast compiles a full forecast report using database-persisted or default assumptions
func (s *FPAService) GetForecast(ctx context.Context, scenarioID string, horizon int) (*fpa.ForecastReport, error) {
	if scenarioID == "" {
		scenarioID = "base"
	}
	cfg := fpa.GetBaselineScenarioConfig(scenarioID)

	// Apply database customizations if present
	cfg = s.applyDatabaseAssumptions(ctx, scenarioID, cfg)

	if horizon > 0 {
		cfg.Growth.ForecastHorizonMonths = horizon
	}

	report := fpa.GenerateForecastReport(cfg)
	return report, nil
}

// GetScenarioComparison runs Base, Conservative, and High-Growth scenarios side-by-side
func (s *FPAService) GetScenarioComparison(ctx context.Context, horizon int) (map[string]*fpa.ForecastReport, error) {
	if horizon <= 0 {
		horizon = 60
	}

	scenarios := []string{"conservative", "base", "high_growth"}
	comparison := make(map[string]*fpa.ForecastReport)

	for _, sc := range scenarios {
		rep, err := s.GetForecast(ctx, sc, horizon)
		if err != nil {
			return nil, err
		}
		comparison[sc] = rep
	}

	return comparison, nil
}

// GetCapacityPlanning20M returns the technical scale evaluation for 20M DAU
func (s *FPAService) GetCapacityPlanning20M(ctx context.Context) (fpa.CapacityPlanning20M, error) {
	cfg := fpa.GetBaselineScenarioConfig("base")
	cfg = s.applyDatabaseAssumptions(ctx, "base", cfg)
	return fpa.Assess20MCapacityScenario(cfg.Workload, cfg.Infra), nil
}

// GetTaxRules returns official Ethiopian statutory tax rules with source citations
func (s *FPAService) GetTaxRules(ctx context.Context) ([]repository.TaxRuleModel, error) {
	return s.repo.ListTaxRules(ctx)
}

// GetActuals returns historical accounting benchmarks
func (s *FPAService) GetActuals(ctx context.Context) ([]repository.ActualRecordModel, error) {
	return s.repo.ListActuals(ctx)
}

// GetAuditLogs returns financial modification logs
func (s *FPAService) GetAuditLogs(ctx context.Context, limit int) ([]repository.AuditLogModel, error) {
	return s.repo.ListAuditLogs(ctx, limit)
}

// UpdateAssumptionsSection updates a specific financial section, validating and recording an audit trail
func (s *FPAService) UpdateAssumptionsSection(ctx context.Context, adminID int, scenarioID, section string, raw json.RawMessage, ipAddress string) error {
	if scenarioID == "" || section == "" {
		return fmt.Errorf("scenarioID and section must not be empty")
	}

	// Validate JSON syntax
	var js map[string]interface{}
	if err := json.Unmarshal(raw, &js); err != nil {
		return fmt.Errorf("invalid assumptions JSON: %w", err)
	}

	if err := s.repo.SaveAssumptions(ctx, scenarioID, section, raw, adminID); err != nil {
		return err
	}

	auditDetails, _ := json.Marshal(map[string]interface{}{
		"section": section,
		"keys":    len(js),
	})
	_ = s.repo.LogAudit(ctx, adminID, scenarioID, "ASSUMPTIONS_UPDATED", auditDetails, ipAddress)
	return nil
}

// ResetAssumptions removes custom database overrides, reverting to code defaults
func (s *FPAService) ResetAssumptions(ctx context.Context, adminID int, scenarioID, ipAddress string) error {
	if err := s.repo.ResetAssumptions(ctx, scenarioID); err != nil {
		return err
	}
	auditDetails, _ := json.Marshal(map[string]string{
		"action": "reverted_to_system_defaults",
	})
	_ = s.repo.LogAudit(ctx, adminID, scenarioID, "ASSUMPTIONS_RESET", auditDetails, ipAddress)
	return nil
}

// ExportForecastCSV generates an RFC-4180 CSV export
func (s *FPAService) ExportForecastCSV(ctx context.Context, scenarioID string, horizon int) (string, error) {
	rep, err := s.GetForecast(ctx, scenarioID, horizon)
	if err != nil {
		return "", err
	}
	return fpa.GenerateCSVExport(rep), nil
}

// applyDatabaseAssumptions applies any saved customizations on top of baseline config
func (s *FPAService) applyDatabaseAssumptions(ctx context.Context, scenarioID string, cfg fpa.FullScenarioConfig) fpa.FullScenarioConfig {
	sections := []string{"growth", "workload", "infra", "staffing", "revenue", "tax", "cash_flow"}
	for _, sec := range sections {
		raw, err := s.repo.GetAssumptions(ctx, scenarioID, sec)
		if err == nil && len(raw) > 0 {
			switch sec {
			case "growth":
				_ = json.Unmarshal(raw, &cfg.Growth)
			case "workload":
				_ = json.Unmarshal(raw, &cfg.Workload)
			case "infra":
				_ = json.Unmarshal(raw, &cfg.Infra)
			case "staffing":
				_ = json.Unmarshal(raw, &cfg.Staffing)
			case "revenue":
				_ = json.Unmarshal(raw, &cfg.Revenue)
			case "tax":
				_ = json.Unmarshal(raw, &cfg.Tax)
			case "cash_flow":
				_ = json.Unmarshal(raw, &cfg.CashFlow)
			}
		}
	}
	return cfg
}
