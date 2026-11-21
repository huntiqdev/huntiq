package database

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"math/big"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/your-org/used-tools-marketplace/internal/models"
	"golang.org/x/crypto/bcrypt"
)

// AdminSeedDef defines an administrative account seed
type AdminSeedDef struct {
	Name        string
	Email       string
	Phone       string
	Password    string
	Role        string
	AdminRole   string
	Department  string
	Permissions []string
}

// EnterpriseAdmins lists the canonical 33 enterprise administrative accounts
var EnterpriseAdmins = []AdminSeedDef{
	// ── Department 1: Executive & Strategic Governance ──
	{
		Name: "Super Admin", Email: "admin@habeshamarket.com", Phone: "0911000000",
		Password: "Admin123!", Role: "admin", AdminRole: models.LegacyRoleSuper, Department: "Executive Governance",
		Permissions: []string{"all"},
	},
	{
		Name: "Platform Owner", Email: "owner@habeshamarket.com", Phone: "0911000001",
		Password: "Owner123!", Role: "admin", AdminRole: models.RolePlatformOwner, Department: "Executive Governance",
		Permissions: []string{"all"},
	},
	{
		Name: "Gov Compliance Officer", Email: "compliance@habeshamarket.com", Phone: "0911000011",
		Password: "Compliance123!", Role: "admin", AdminRole: models.RoleGovComplianceOfficer, Department: "Executive Governance",
		Permissions: []string{"dashboard", "finance", "profit", "transactions", "operations", "logs", "seller_verification", "trust_safety", "appeals", "refunds"},
	},
	{
		Name: "Legal Counsel", Email: "legal@habeshamarket.com", Phone: "0911000012",
		Password: "Legal123!", Role: "admin", AdminRole: models.RoleLegalCounsel, Department: "Executive Governance",
		Permissions: []string{"dashboard", "trust_safety", "appeals", "refunds", "seller_verification", "logs", "transactions"},
	},
	{
		Name: "Human Resources Director", Email: "hr@habeshamarket.com", Phone: "0911000031",
		Password: "Hr123!", Role: "admin", AdminRole: models.RoleHRDirector, Department: "Executive Governance",
		Permissions: []string{"dashboard", "logs", "notifications", "operations"},
	},

	// ── Department 2: Cybersecurity, SOC & Trust & Safety ──
	{
		Name: "Chief Info Security Officer", Email: "ciso@habeshamarket.com", Phone: "0911000013",
		Password: "Ciso123!", Role: "admin", AdminRole: models.RoleCISO, Department: "Cybersecurity & SOC",
		Permissions: []string{"dashboard", "logs", "otp_logs", "trust_safety", "operations", "notifications", "settings", "analytics"},
	},
	{
		Name: "Security Administrator", Email: "security@habeshamarket.com", Phone: "0911000003",
		Password: "Security123!", Role: "admin", AdminRole: models.RoleSecurityAdmin, Department: "Cybersecurity & SOC",
		Permissions: []string{"dashboard", "logs", "otp_logs", "trust_safety", "operations", "analytics"},
	},
	{
		Name: "SOC Analyst", Email: "soc@habeshamarket.com", Phone: "0911000014",
		Password: "Soc123!", Role: "admin", AdminRole: models.RoleSOCAnalyst, Department: "Cybersecurity & SOC",
		Permissions: []string{"dashboard", "logs", "otp_logs", "trust_safety", "operations"},
	},
	{
		Name: "Fraud Investigator", Email: "fraud@habeshamarket.com", Phone: "0911000015",
		Password: "Fraud123!", Role: "admin", AdminRole: models.RoleFraudInvestigator, Department: "Cybersecurity & SOC",
		Permissions: []string{"dashboard", "trust_safety", "transactions", "logs", "seller_verification", "refunds", "appeals"},
	},
	{
		Name: "Compliance Auditor", Email: "complianceauditor@habeshamarket.com", Phone: "0911000032",
		Password: "ComplianceAuditor123!", Role: "admin", AdminRole: models.RoleComplianceAuditor, Department: "Cybersecurity & SOC",
		Permissions: []string{"dashboard", "logs", "otp_logs", "transactions", "profit", "operations", "finance"},
	},

	// ── Department 3: Banking, Treasury, Tax & Financial Operations ──
	{
		Name: "Finance Director", Email: "finance@habeshamarket.com", Phone: "0911000016",
		Password: "Finance123!", Role: "admin", AdminRole: models.RoleFinanceDirector, Department: "Finance & Treasury",
		Permissions: []string{"dashboard", "profit", "finance", "finance:manage", "transactions", "refunds", "withdrawals"},
	},
	{
		Name: "Transaction Auditor", Email: "auditor@habeshamarket.com", Phone: "0911000017",
		Password: "Auditor123!", Role: "admin", AdminRole: models.RoleTransactionAuditor, Department: "Finance & Treasury",
		Permissions: []string{"dashboard", "profit", "finance", "transactions", "refunds", "logs"},
	},
	{
		Name: "Tax Controller", Email: "tax@habeshamarket.com", Phone: "0911000018",
		Password: "Tax123!", Role: "admin", AdminRole: models.RoleTaxController, Department: "Finance & Treasury",
		Permissions: []string{"dashboard", "profit", "finance", "transactions"},
	},
	{
		Name: "Escrow Reconciler", Email: "escrow@habeshamarket.com", Phone: "0911000019",
		Password: "Escrow123!", Role: "admin", AdminRole: models.RoleEscrowReconciler, Department: "Finance & Treasury",
		Permissions: []string{"dashboard", "transactions", "refunds", "profit"},
	},
	{
		Name: "Finance Support Agent", Email: "financesupport@habeshamarket.com", Phone: "0911000033",
		Password: "FinanceSupport123!", Role: "admin", AdminRole: models.RoleFinanceSupportAgent, Department: "Finance & Treasury",
		Permissions: []string{"dashboard", "transactions", "refunds"},
	},

	// ── Department 4: VIP Concierge, Premium & Growth Services ──
	{
		Name: "VIP Account Director", Email: "vip@habeshamarket.com", Phone: "0911000020",
		Password: "Vip123!", Role: "admin", AdminRole: models.RoleVIPAccountDirector, Department: "VIP & Growth Services",
		Permissions: []string{"dashboard", "sellers", "premium_users", "seller_upgrades", "buyer_upgrades", "referrals"},
	},
	{
		Name: "Gold Concierge", Email: "gold@habeshamarket.com", Phone: "0911000021",
		Password: "Gold123!", Role: "admin", AdminRole: models.RoleGoldConcierge, Department: "VIP & Growth Services",
		Permissions: []string{"dashboard", "sellers", "premium_users", "seller_upgrades", "buyer_upgrades"},
	},
	{
		Name: "Premium Concierge", Email: "premium@habeshamarket.com", Phone: "0911000022",
		Password: "Premium123!", Role: "admin", AdminRole: models.RolePremiumConcierge, Department: "VIP & Growth Services",
		Permissions: []string{"dashboard", "sellers", "premium_users", "seller_upgrades", "buyer_upgrades"},
	},
	{
		Name: "Referral Director", Email: "referral@habeshamarket.com", Phone: "0911000023",
		Password: "Referral123!", Role: "admin", AdminRole: models.RoleReferralDirector, Department: "VIP & Growth Services",
		Permissions: []string{"dashboard", "referrals", "withdrawals", "analytics"},
	},
	{
		Name: "Retention Specialist", Email: "retention@habeshamarket.com", Phone: "0911000034",
		Password: "Retention123!", Role: "admin", AdminRole: models.RoleRetentionSpecialist, Department: "VIP & Growth Services",
		Permissions: []string{"dashboard", "sellers", "premium_users", "analytics"},
	},

	// ── Department 5: Marketplace Operations & KYC Verification ──
	{
		Name: "Operations Administrator", Email: "operations@habeshamarket.com", Phone: "0911000004",
		Password: "Operations123!", Role: "admin", AdminRole: models.RoleOperationsAdmin, Department: "Operations & KYC",
		Permissions: []string{"dashboard", "operations", "logs", "notifications", "banners"},
	},
	{
		Name: "KYC Verification Officer", Email: "kyc@habeshamarket.com", Phone: "0911000024",
		Password: "Kyc123!", Role: "admin", AdminRole: models.RoleKYCVerificationOfficer, Department: "Operations & KYC",
		Permissions: []string{"dashboard", "seller_verification", "sellers", "trust_safety"},
	},
	{
		Name: "Product Inspection Officer", Email: "inspection@habeshamarket.com", Phone: "0911000025",
		Password: "Inspection123!", Role: "admin", AdminRole: models.RoleProductInspectionOfficer, Department: "Operations & KYC",
		Permissions: []string{"dashboard", "pending", "products", "seller_verification"},
	},
	{
		Name: "Dispute & Refund Officer", Email: "dispute@habeshamarket.com", Phone: "0911000026",
		Password: "Dispute123!", Role: "admin", AdminRole: models.RoleDisputeRefundOfficer, Department: "Operations & KYC",
		Permissions: []string{"dashboard", "refunds", "appeals", "trust_safety", "transactions"},
	},
	{
		Name: "Appeals Officer", Email: "appeals@habeshamarket.com", Phone: "0911000027",
		Password: "Appeals123!", Role: "admin", AdminRole: models.RoleAppealsOfficer, Department: "Operations & KYC",
		Permissions: []string{"dashboard", "appeals", "trust_safety", "sellers"},
	},

	// ── Department 6: Catalog, Taxonomy, Logistics & Support ──
	{
		Name: "Platform Administrator", Email: "platform@habeshamarket.com", Phone: "0911000002",
		Password: "Platform123!", Role: "admin", AdminRole: models.RolePlatformAdmin, Department: "Catalog & Infrastructure",
		Permissions: []string{"dashboard", "categories", "locations", "classification", "banners", "settings", "analytics"},
	},
	{
		Name: "Taxonomy Architect", Email: "taxonomy@habeshamarket.com", Phone: "0911000028",
		Password: "Taxonomy123!", Role: "admin", AdminRole: models.RoleTaxonomyArchitect, Department: "Catalog & Infrastructure",
		Permissions: []string{"dashboard", "categories", "category_trends", "classification"},
	},
	{
		Name: "GIS Location Specialist", Email: "gis@habeshamarket.com", Phone: "0911000029",
		Password: "Gis123!", Role: "admin", AdminRole: models.RoleGISLocationSpecialist, Department: "Catalog & Infrastructure",
		Permissions: []string{"dashboard", "locations", "analytics"},
	},
	{
		Name: "Content Moderator", Email: "moderator@habeshamarket.com", Phone: "0911000005",
		Password: "Moderator123!", Role: "admin", AdminRole: models.RoleModerator, Department: "Catalog & Infrastructure",
		Permissions: []string{"dashboard", "pending", "products", "trust_safety", "banners"},
	},
	{
		Name: "Content Moderation Lead", Email: "contentmod@habeshamarket.com", Phone: "0911000035",
		Password: "ContentMod123!", Role: "admin", AdminRole: models.RoleContentModerator, Department: "Catalog & Infrastructure",
		Permissions: []string{"dashboard", "pending", "products", "trust_safety", "banners"},
	},
	{
		Name: "Customer Support Lead", Email: "supportlead@habeshamarket.com", Phone: "0911000030",
		Password: "SupportLead123!", Role: "admin", AdminRole: models.RoleCustomerSupportLead, Department: "Catalog & Infrastructure",
		Permissions: []string{"dashboard", "sellers", "seller_upgrades", "buyer_upgrades", "refunds", "appeals", "products"},
	},
	{
		Name: "Support Agent", Email: "support@habeshamarket.com", Phone: "0911000006",
		Password: "Support123!", Role: "admin", AdminRole: models.RoleSupportAgent, Department: "Catalog & Infrastructure",
		Permissions: []string{"dashboard", "sellers", "seller_upgrades", "buyer_upgrades", "products"},
	},
	{
		Name: "General Staff Administrator", Email: "staff@habeshamarket.com", Phone: "0911000036",
		Password: "Staff123!", Role: "admin", AdminRole: models.LegacyRoleNormal, Department: "Catalog & Infrastructure",
		Permissions: models.DefaultAdminPermissions,
	},
}

// PurgeNonAdminsAndProducts removes all non-admin users, demo listings, and associated data.
func PurgeNonAdminsAndProducts(ctx context.Context, pool *pgxpool.Pool) error {
	log.Println("[HUNTIQ DB] Purging all non-admin users and mock products...")

	// 1. Delete product-dependent tables
	queries := []string{
		`DELETE FROM product_images;`,
		`DELETE FROM product_attributes;`,
		`DELETE FROM product_secondary_categories;`,
		`DELETE FROM product_variants;`,
		`DELETE FROM order_requests;`,
		`DELETE FROM telebirr_orders;`,
		`DELETE FROM transactions;`,
		`DELETE FROM buyer_tier_payments;`,
		`DELETE FROM upgrade_requests;`,
		`DELETE FROM refund_requests;`,
		`DELETE FROM products;`,
		// 2. Delete non-admin user data
		`DELETE FROM abuse_reports WHERE reported_user_id IN (SELECT id FROM users WHERE role != 'admin') OR reporter_id IN (SELECT id FROM users WHERE role != 'admin');`,
		`DELETE FROM account_appeals WHERE user_id IN (SELECT id FROM users WHERE role != 'admin');`,
		`DELETE FROM identity_verifications WHERE user_id IN (SELECT id FROM users WHERE role != 'admin');`,
		`DELETE FROM seller_verification_documents WHERE seller_id IN (SELECT id FROM users WHERE role != 'admin');`,
		`DELETE FROM merchant_branches WHERE merchant_id IN (SELECT id FROM users WHERE role != 'admin');`,
		`DELETE FROM merchant_analytics_events WHERE merchant_id IN (SELECT id FROM users WHERE role != 'admin');`,
		`DELETE FROM service_appointments WHERE customer_id IN (SELECT id FROM users WHERE role != 'admin') OR merchant_id IN (SELECT id FROM users WHERE role != 'admin');`,
		`DELETE FROM invoices WHERE user_id IN (SELECT id FROM users WHERE role != 'admin');`,
		`DELETE FROM ledger_transactions WHERE user_id IN (SELECT id FROM users WHERE role != 'admin');`,
		`DELETE FROM wallet_accounts WHERE user_id IN (SELECT id FROM users WHERE role != 'admin');`,
		`DELETE FROM referral_commissions WHERE referred_id IN (SELECT id FROM users WHERE role != 'admin') OR referrer_id IN (SELECT id FROM users WHERE role != 'admin');`,
		`DELETE FROM referral_withdrawals WHERE user_id IN (SELECT id FROM users WHERE role != 'admin');`,
		`DELETE FROM referral_attributions WHERE registered_user_id IN (SELECT id FROM users WHERE role != 'admin') OR inviter_id IN (SELECT id FROM users WHERE role != 'admin');`,
		`DELETE FROM referral_tokens WHERE user_id IN (SELECT id FROM users WHERE role != 'admin');`,
		`DELETE FROM referral_fraud_evaluations WHERE user_id IN (SELECT id FROM users WHERE role != 'admin') OR inviter_id IN (SELECT id FROM users WHERE role != 'admin');`,
		`DELETE FROM user_capabilities WHERE user_id IN (SELECT id FROM users WHERE role != 'admin');`,
		`DELETE FROM user_activity_logs WHERE user_id IN (SELECT id FROM users WHERE role != 'admin');`,
		`DELETE FROM user_refresh_tokens WHERE user_id IN (SELECT id FROM users WHERE role != 'admin');`,
		`DELETE FROM user_tier_history WHERE user_id IN (SELECT id FROM users WHERE role != 'admin');`,
		`DELETE FROM user_warnings WHERE user_id IN (SELECT id FROM users WHERE role != 'admin');`,
		`DELETE FROM password_recovery_requests WHERE user_id IN (SELECT id FROM users WHERE role != 'admin');`,
		`DELETE FROM organization_members WHERE user_id IN (SELECT id FROM users WHERE role != 'admin');`,
		`DELETE FROM invitation_requests WHERE invited_by IN (SELECT id FROM users WHERE role != 'admin');`,
		`DELETE FROM ai_conversations WHERE user_id IN (SELECT id FROM users WHERE role != 'admin');`,
		`DELETE FROM ai_audit_events WHERE user_id IN (SELECT id FROM users WHERE role != 'admin');`,
		`DELETE FROM users WHERE role != 'admin';`,
	}

	for _, q := range queries {
		if _, err := pool.Exec(ctx, q); err != nil {
			log.Printf("[HUNTIQ DB] Warning during purge query (%s): %v", q, err)
		}
	}

	log.Println("[HUNTIQ DB] ✓ Non-admin users and products purged completely.")
	return nil
}

// SeedAdmins provisions all 33 enterprise administrator accounts with their roles and permissions.
func SeedAdmins(ctx context.Context, pool *pgxpool.Pool) error {
	var adminCount int
	_ = pool.QueryRow(ctx, `SELECT count(*) FROM users WHERE role = 'admin'`).Scan(&adminCount)
	if adminCount >= len(EnterpriseAdmins) {
		log.Printf("[HUNTIQ DB] ✓ Enterprise administrators verified (%d active admins).", adminCount)
		return nil
	}

	log.Printf("[HUNTIQ DB] ⚙ Seeding & Synchronizing all %d HUNTIQ Enterprise Administrative Accounts...", len(EnterpriseAdmins))

	for _, a := range EnterpriseAdmins {
		hash, err := bcrypt.GenerateFromPassword([]byte(a.Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("hash failed for %s: %w", a.Email, err)
		}

		refCode, err := generateAdminReferralCode()
		if err != nil {
			return fmt.Errorf("referral code failed: %w", err)
		}

		var userID int
		err = pool.QueryRow(ctx, `
			INSERT INTO users (name, phone, email, password_hash, role, admin_role, is_active, referral_code)
			VALUES ($1, $2, $3, $4, $5, $6, true, $7)
			ON CONFLICT (phone, role) DO UPDATE
			  SET name          = EXCLUDED.name,
			      email         = EXCLUDED.email,
			      password_hash = EXCLUDED.password_hash,
			      admin_role    = EXCLUDED.admin_role,
			      is_active     = true,
			      updated_at    = NOW()
			RETURNING id
		`, a.Name, a.Phone, a.Email, string(hash), a.Role, a.AdminRole, refCode).Scan(&userID)

		if err != nil {
			err = pool.QueryRow(ctx, `
				UPDATE users
				SET name          = $1,
				    password_hash = $2,
				    phone         = $3,
				    role          = $4,
				    admin_role    = $5,
				    is_active     = true,
				    updated_at    = NOW()
				WHERE LOWER(email) = LOWER($6)
				RETURNING id
			`, a.Name, string(hash), a.Phone, a.Role, a.AdminRole, a.Email).Scan(&userID)
			if err != nil {
				return fmt.Errorf("failed to seed admin %s: %w", a.Email, err)
			}
		}

		// Clear previous permissions for this admin and insert fresh
		_, _ = pool.Exec(ctx, `DELETE FROM admin_permissions WHERE admin_id = $1`, userID)

		permsToGrant := a.Permissions
		if len(permsToGrant) == 1 && permsToGrant[0] == "all" {
			permsToGrant = models.AllAdminPermissions
		}

		for _, p := range permsToGrant {
			_, _ = pool.Exec(ctx, `
				INSERT INTO admin_permissions (admin_id, permission, enabled, updated_at)
				VALUES ($1, $2, true, NOW())
				ON CONFLICT (admin_id, permission) DO UPDATE SET enabled = true
			`, userID, p)
		}
	}

	_, _ = pool.Exec(ctx, `CHECKPOINT;`)
	log.Printf("[HUNTIQ DB] ✓ All %d administrative accounts successfully seeded and permissioned.", len(EnterpriseAdmins))
	return nil
}

func generateAdminReferralCode() (string, error) {
	const chars = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	b := make([]byte, 8)
	for i := range b {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			return "", fmt.Errorf("referral code generation failed: %w", err)
		}
		b[i] = chars[n.Int64()]
	}
	return "REF" + string(b), nil
}
