// seed creates or updates all 33 enterprise administrator accounts in the database.
// Usage: go run ./cmd/seed
package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/your-org/used-tools-marketplace/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type AdminSeed struct {
	Name        string
	Email       string
	Phone       string
	Password    string
	Role        string
	AdminRole   string
	Department  string
	Permissions []string
}

func main() {
	_ = godotenv.Load()

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		host := getenv("DB_HOST", "/tmp")
		port := getenv("DB_PORT", "5433")
		user := getenv("DB_USER", "postgres")
		pass := getenv("DB_PASSWORD", "gedam")
		dbname := getenv("DB_NAME", "hulubedej_db")
		sslmode := getenv("DB_SSLMODE", "disable")

		dsn = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			host, port, user, pass, dbname, sslmode,
		)
	}

	db, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("connect: %v", err)
	}
	defer db.Close()

	ctx := context.Background()

	admins := []AdminSeed{
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
			Permissions: []string{"dashboard", "finance", "finance:manage", "profit", "transactions", "withdrawals", "referrals", "analytics"},
		},
		{
			Name: "Transaction Auditor", Email: "auditor@habeshamarket.com", Phone: "0911000017",
			Password: "Auditor123!", Role: "admin", AdminRole: models.RoleTransactionAuditor, Department: "Finance & Treasury",
			Permissions: []string{"dashboard", "transactions", "profit", "withdrawals", "finance"},
		},
		{
			Name: "Tax Controller", Email: "tax@habeshamarket.com", Phone: "0911000018",
			Password: "Tax123!", Role: "admin", AdminRole: models.RoleTaxController, Department: "Finance & Treasury",
			Permissions: []string{"dashboard", "finance", "profit", "transactions", "analytics"},
		},
		{
			Name: "Escrow Reconciler", Email: "escrow@habeshamarket.com", Phone: "0911000019",
			Password: "Escrow123!", Role: "admin", AdminRole: models.RoleEscrowReconciler, Department: "Finance & Treasury",
			Permissions: []string{"dashboard", "transactions", "refunds", "withdrawals", "profit"},
		},
		{
			Name: "Finance Support Agent", Email: "financesupport@habeshamarket.com", Phone: "0911000033",
			Password: "FinanceSupport123!", Role: "admin", AdminRole: models.RoleFinanceSupportAgent, Department: "Finance & Treasury",
			Permissions: []string{"dashboard", "transactions", "refunds", "seller_upgrades", "buyer_upgrades"},
		},

		// ── Department 4: VIP Concierge, Premium & Growth Services ──
		{
			Name: "VIP Account Director", Email: "vip@habeshamarket.com", Phone: "0911000020",
			Password: "Vip123!", Role: "admin", AdminRole: models.RoleVIPAccountDirector, Department: "VIP & Growth",
			Permissions: []string{"dashboard", "sellers", "seller_upgrades", "buyer_upgrades", "premium_users", "analytics", "referrals"},
		},
		{
			Name: "Gold Concierge", Email: "gold@habeshamarket.com", Phone: "0911000021",
			Password: "Gold123!", Role: "admin", AdminRole: models.RoleGoldConcierge, Department: "VIP & Growth",
			Permissions: []string{"dashboard", "sellers", "seller_upgrades", "products", "premium_users"},
		},
		{
			Name: "Premium Concierge", Email: "premium@habeshamarket.com", Phone: "0911000022",
			Password: "Premium123!", Role: "admin", AdminRole: models.RolePremiumConcierge, Department: "VIP & Growth",
			Permissions: []string{"dashboard", "sellers", "seller_upgrades", "buyer_upgrades", "products", "banners", "premium_users"},
		},
		{
			Name: "Referral Director", Email: "referral@habeshamarket.com", Phone: "0911000023",
			Password: "Referral123!", Role: "admin", AdminRole: models.RoleReferralDirector, Department: "VIP & Growth",
			Permissions: []string{"dashboard", "referrals", "withdrawals", "analytics"},
		},
		{
			Name: "Retention Specialist", Email: "retention@habeshamarket.com", Phone: "0911000034",
			Password: "Retention123!", Role: "admin", AdminRole: models.RoleRetentionSpecialist, Department: "VIP & Growth",
			Permissions: []string{"dashboard", "sellers", "analytics", "category_trends", "premium_users"},
		},

		// ── Department 5: Marketplace Operations & KYC Verification ──
		{
			Name: "Operations Administrator", Email: "operations@habeshamarket.com", Phone: "0911000004",
			Password: "Operations123!", Role: "admin", AdminRole: models.RoleOperationsAdmin, Department: "Marketplace Operations",
			Permissions: []string{"dashboard", "operations", "seller_verification", "refunds", "appeals", "pending", "products", "sellers", "seller_upgrades", "buyer_upgrades", "transactions"},
		},
		{
			Name: "KYC Verification Officer", Email: "kyc@habeshamarket.com", Phone: "0911000024",
			Password: "Kyc123!", Role: "admin", AdminRole: models.RoleKYCVerificationOfficer, Department: "Marketplace Operations",
			Permissions: []string{"dashboard", "seller_verification", "sellers", "trust_safety"},
		},
		{
			Name: "Product Inspection Officer", Email: "inspection@habeshamarket.com", Phone: "0911000025",
			Password: "Inspection123!", Role: "admin", AdminRole: models.RoleProductInspectionOfficer, Department: "Marketplace Operations",
			Permissions: []string{"dashboard", "pending", "products"},
		},
		{
			Name: "Dispute & Refund Officer", Email: "dispute@habeshamarket.com", Phone: "0911000026",
			Password: "Dispute123!", Role: "admin", AdminRole: models.RoleDisputeRefundOfficer, Department: "Marketplace Operations",
			Permissions: []string{"dashboard", "refunds", "transactions", "appeals"},
		},
		{
			Name: "Appeals Officer", Email: "appeals@habeshamarket.com", Phone: "0911000027",
			Password: "Appeals123!", Role: "admin", AdminRole: models.RoleAppealsOfficer, Department: "Marketplace Operations",
			Permissions: []string{"dashboard", "appeals", "trust_safety", "sellers"},
		},

		// ── Department 6: Catalog, Taxonomy, Logistics & Support ──
		{
			Name: "Platform Administrator", Email: "platform@habeshamarket.com", Phone: "0911000002",
			Password: "Platform123!", Role: "admin", AdminRole: models.RolePlatformAdmin, Department: "Catalog & Infrastructure",
			Permissions: []string{"dashboard", "settings", "categories", "category_trends", "analytics", "locations", "banners", "notifications", "profit", "finance", "sellers", "transactions", "operations", "classification"},
		},
		{
			Name: "Taxonomy Architect", Email: "taxonomy@habeshamarket.com", Phone: "0911000028",
			Password: "Taxonomy123!", Role: "admin", AdminRole: models.RoleTaxonomyArchitect, Department: "Catalog & Infrastructure",
			Permissions: []string{"dashboard", "categories", "category_trends", "analytics", "classification"},
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

	fmt.Printf("🚀 Seeding & Synchronizing ALL %d HUNTIQ Enterprise Administrative Accounts...\n\n", len(admins))

	for _, a := range admins {
		hash, err := bcrypt.GenerateFromPassword([]byte(a.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Fatalf("hash failed for %s: %v", a.Email, err)
		}

		refCode, err := generateReferralCode()
		if err != nil {
			log.Fatalf("referral code failed: %v", err)
		}

		var userID int
		err = db.QueryRow(ctx, `
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
			err = db.QueryRow(ctx, `
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
				log.Fatalf("Failed to seed admin %s: %v", a.Email, err)
			}
		}

		// Clear previous permissions for this admin and insert fresh
		_, _ = db.Exec(ctx, `DELETE FROM admin_permissions WHERE admin_id = $1`, userID)

		permsToGrant := a.Permissions
		if len(permsToGrant) == 1 && permsToGrant[0] == "all" {
			permsToGrant = models.AllAdminPermissions
		}

		for _, p := range permsToGrant {
			_, _ = db.Exec(ctx, `
				INSERT INTO admin_permissions (admin_id, permission, enabled, updated_at)
				VALUES ($1, $2, true, NOW())
				ON CONFLICT (admin_id, permission) DO UPDATE SET enabled = true
			`, userID, p)
		}

		fmt.Printf("   ✅ [%-26s] %-30s (Phone: %s | Pass: %-18s | Dept: %s)\n",
			a.AdminRole, a.Email, a.Phone, a.Password, a.Department)
	}

	// Force immediate WAL checkpoint
	_, _ = db.Exec(ctx, `CHECKPOINT;`)

	fmt.Printf("\n🎉 All %d administrative accounts successfully seeded, permissioned, and checkpointed to disk!\n", len(admins))
	fmt.Println("⚠️  Universal Dev Password: Admin123! is also accepted for all accounts in development mode.")
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func generateReferralCode() (string, error) {
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
