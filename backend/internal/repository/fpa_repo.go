package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ScenarioModel represents a database record in financial_scenarios
type ScenarioModel struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TaxRuleModel represents a database record in financial_tax_rules
type TaxRuleModel struct {
	ID              int        `json:"id"`
	Jurisdiction    string     `json:"jurisdiction"`
	RuleCode        string     `json:"rule_code"`
	RuleName        string     `json:"rule_name"`
	Category        string     `json:"category"`
	StandardRate    *float64   `json:"standard_rate"`
	ThresholdAmount *float64   `json:"threshold_amount"`
	EffectiveDate   string     `json:"effective_date"`
	SourceCitation  string     `json:"source_citation"`
	SourceURL       *string    `json:"source_url"`
	IsVerified      bool       `json:"is_verified"`
	ReviewerNotes   *string    `json:"reviewer_notes"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// ActualRecordModel represents historical actual accounting data
type ActualRecordModel struct {
	ID                     int       `json:"id"`
	PeriodMonth            string    `json:"period_month"`
	Currency               string    `json:"currency"`
	ActiveUsers            int       `json:"active_users"`
	CompletedTransactions  int       `json:"completed_transactions"`
	GrossMerchandiseVolume float64   `json:"gross_merchandise_volume"`
	CommissionRevenue      float64   `json:"commission_revenue"`
	ListingFeeRevenue      float64   `json:"listing_fee_revenue"`
	SubscriptionRevenue    float64   `json:"subscription_revenue"`
	OtherRevenue           float64   `json:"other_revenue"`
	InfraCostActual        float64   `json:"infra_cost_actual"`
	StaffingCostActual     float64   `json:"staffing_cost_actual"`
	MarketingCostActual    float64   `json:"marketing_cost_actual"`
	AdminCostActual        float64   `json:"admin_cost_actual"`
	VATCollectedActual     float64   `json:"vat_collected_actual"`
	VATPaidActual          float64   `json:"vat_paid_actual"`
	IncomeTaxPaidActual    float64   `json:"income_tax_paid_actual"`
	ClosingCashBalance     float64   `json:"closing_cash_balance"`
	IsAudited              bool      `json:"is_audited"`
	DataSource             string    `json:"data_source"`
	Notes                  *string   `json:"notes"`
	UpdatedAt              time.Time `json:"updated_at"`
}

// AuditLogModel represents financial audit log entries
type AuditLogModel struct {
	ID         int             `json:"id"`
	UserID     *int            `json:"user_id"`
	ScenarioID *string         `json:"scenario_id"`
	Action     string          `json:"action"`
	Details    json.RawMessage `json:"details"`
	IPAddress  *string         `json:"ip_address"`
	CreatedAt  time.Time       `json:"created_at"`
}

// FPARepo handles PostgreSQL operations for financial planning
type FPARepo struct {
	db *pgxpool.Pool
}

// NewFPARepo initializes a new FPARepo
func NewFPARepo(db *pgxpool.Pool) *FPARepo {
	return &FPARepo{db: db}
}

// ListScenarios returns all configured scenarios
func (r *FPARepo) ListScenarios(ctx context.Context) ([]ScenarioModel, error) {
	if r == nil || r.db == nil {
		return []ScenarioModel{
			{ID: "base", Name: "Base Case Planning", Description: "Standard model", IsActive: true},
			{ID: "conservative", Name: "Conservative Downside Case", Description: "Downside model", IsActive: true},
			{ID: "high_growth", Name: "High Growth Upside Case", Description: "Upside model", IsActive: true},
		}, nil
	}
	rows, err := r.db.Query(ctx, `
		SELECT id, name, COALESCE(description, ''), is_active, created_at, updated_at
		FROM financial_scenarios
		ORDER BY created_at ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []ScenarioModel
	for rows.Next() {
		var s ScenarioModel
		if err := rows.Scan(&s.ID, &s.Name, &s.Description, &s.IsActive, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, s)
	}
	return list, nil
}

// GetAssumptions loads saved assumption JSON for a scenario and section
func (r *FPARepo) GetAssumptions(ctx context.Context, scenarioID, section string) (json.RawMessage, error) {
	if r == nil || r.db == nil {
		return nil, nil
	}
	var raw json.RawMessage
	err := r.db.QueryRow(ctx, `
		SELECT assumptions_json
		FROM financial_assumptions
		WHERE scenario_id = $1 AND section = $2`, scenarioID, section).Scan(&raw)
	return raw, err
}

// SaveAssumptions upserts custom assumptions
func (r *FPARepo) SaveAssumptions(ctx context.Context, scenarioID, section string, payload json.RawMessage, adminID int) error {
	if r == nil || r.db == nil {
		return nil
	}
	_, err := r.db.Exec(ctx, `
		INSERT INTO financial_assumptions (scenario_id, section, assumptions_json, updated_by, updated_at)
		VALUES ($1, $2, $3, $4, NOW())
		ON CONFLICT (scenario_id, section) DO UPDATE
		SET assumptions_json = EXCLUDED.assumptions_json,
		    updated_by = EXCLUDED.updated_by,
		    updated_at = NOW()`, scenarioID, section, payload, adminID)
	return err
}

// ResetAssumptions deletes custom overrides for a scenario, reverting to code defaults
func (r *FPARepo) ResetAssumptions(ctx context.Context, scenarioID string) error {
	if r == nil || r.db == nil {
		return nil
	}
	_, err := r.db.Exec(ctx, `DELETE FROM financial_assumptions WHERE scenario_id = $1`, scenarioID)
	return err
}

// ListTaxRules retrieves all official tax rules with source citations
func (r *FPARepo) ListTaxRules(ctx context.Context) ([]TaxRuleModel, error) {
	if r == nil || r.db == nil {
		return []TaxRuleModel{}, nil
	}
	rows, err := r.db.Query(ctx, `
		SELECT id, jurisdiction, rule_code, rule_name, category, standard_rate, threshold_amount,
		       to_char(effective_date, 'YYYY-MM-DD'), source_citation, source_url, is_verified,
		       reviewer_notes, updated_at
		FROM financial_tax_rules
		ORDER BY category ASC, id ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []TaxRuleModel
	for rows.Next() {
		var t TaxRuleModel
		if err := rows.Scan(
			&t.ID, &t.Jurisdiction, &t.RuleCode, &t.RuleName, &t.Category,
			&t.StandardRate, &t.ThresholdAmount, &t.EffectiveDate,
			&t.SourceCitation, &t.SourceURL, &t.IsVerified, &t.ReviewerNotes, &t.UpdatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, t)
	}
	return list, nil
}

// ListActuals returns recorded monthly accounting actuals
func (r *FPARepo) ListActuals(ctx context.Context) ([]ActualRecordModel, error) {
	if r == nil || r.db == nil {
		return []ActualRecordModel{}, nil
	}
	rows, err := r.db.Query(ctx, `
		SELECT id, period_month, currency, active_users, completed_transactions,
		       gross_merchandise_volume, commission_revenue, listing_fee_revenue,
		       subscription_revenue, other_revenue, infra_cost_actual, staffing_cost_actual,
		       marketing_cost_actual, admin_cost_actual, vat_collected_actual, vat_paid_actual,
		       income_tax_paid_actual, closing_cash_balance, is_audited, data_source, notes, updated_at
		FROM financial_actuals
		ORDER BY period_month DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []ActualRecordModel
	for rows.Next() {
		var a ActualRecordModel
		if err := rows.Scan(
			&a.ID, &a.PeriodMonth, &a.Currency, &a.ActiveUsers, &a.CompletedTransactions,
			&a.GrossMerchandiseVolume, &a.CommissionRevenue, &a.ListingFeeRevenue,
			&a.SubscriptionRevenue, &a.OtherRevenue, &a.InfraCostActual, &a.StaffingCostActual,
			&a.MarketingCostActual, &a.AdminCostActual, &a.VATCollectedActual, &a.VATPaidActual,
			&a.IncomeTaxPaidActual, &a.ClosingCashBalance, &a.IsAudited, &a.DataSource, &a.Notes, &a.UpdatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, a)
	}
	return list, nil
}

// LogAudit records a financial audit log entry
func (r *FPARepo) LogAudit(ctx context.Context, adminID int, scenarioID, action string, details json.RawMessage, ipAddress string) error {
	if r == nil || r.db == nil {
		return nil
	}
	_, err := r.db.Exec(ctx, `
		INSERT INTO financial_audit_logs (user_id, scenario_id, action, details, ip_address, created_at)
		VALUES ($1, $2, $3, $4, $5, NOW())`, adminID, scenarioID, action, details, ipAddress)
	return err
}

// ListAuditLogs retrieves the most recent financial audit trail entries
func (r *FPARepo) ListAuditLogs(ctx context.Context, limit int) ([]AuditLogModel, error) {
	if r == nil || r.db == nil {
		return []AuditLogModel{}, nil
	}
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	rows, err := r.db.Query(ctx, `
		SELECT id, user_id, scenario_id, action, details, ip_address, created_at
		FROM financial_audit_logs
		ORDER BY created_at DESC
		LIMIT $1`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []AuditLogModel
	for rows.Next() {
		var l AuditLogModel
		if err := rows.Scan(&l.ID, &l.UserID, &l.ScenarioID, &l.Action, &l.Details, &l.IPAddress, &l.CreatedAt); err != nil {
			return nil, err
		}
		logs = append(logs, l)
	}
	return logs, nil
}
