package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/your-org/used-tools-marketplace/internal/config"
	handlers "github.com/your-org/used-tools-marketplace/internal/handler"
	"github.com/your-org/used-tools-marketplace/internal/middleware"
)

type AppDeps struct {
	Config                *config.Config
	MW                    *middleware.Middleware
	ProductHandler        *handlers.ProductHandler
	LocationHandler       *handlers.LocationHandler
	CategoryHandler       *handlers.CategoryHandler
	AttributeHandler      *handlers.AttributeHandler
	ClassificationHandler *handlers.ClassificationHandler
	AuthHandler           *handlers.AuthHandler
	AdminAuthHandler      *handlers.AdminAuthHandler
	AdminHandler          *handlers.AdminHandler
	SellerHandler         *handlers.SellerHandler
	BuyerHandler          *handlers.BuyerHandler
	TransactionHandler    *handlers.TransactionHandler
	BannerHandler         *handlers.BannerHandler
	ReferralHandler       *handlers.ReferralHandler
	TelebirrHandler       *handlers.TelebirrHandler
	WalletHandler         *handlers.WalletHandler
	AdminFinancialHandler *handlers.AdminFinancialHandler
	FPAHandler            *handlers.FPAHandler
	SearchHandler         *handlers.SearchHandler
	IdentityHandler       *handlers.IdentityHandler
	OrganizationHandler   *handlers.OrganizationHandler
	MasterDataHandler     *handlers.MasterDataHandler
	AIHandler             *handlers.AIHandler
}

func registerRoutes(app *fiber.App, d *AppDeps) {
	api := app.Group("/api")
	api.Use(d.MW.OptionalAuth)
	api.Use(d.MW.RequireNoMaintenance)
	api.Use(middleware.EthiopiaGeoAccess(middleware.GeoPolicyConfig{
		DevMode: d.Config.DevMode,
		EnabledCheck: func() bool {
			return d.Config.EnableEthiopiaGeoBlock
		},
		ExemptPaths: []string{
			"/api/telebirr/webhook",
			"/api/chapa/webhook",
			"/api/cbe/webhook",
		},
	}))
	api.Use(middleware.CSRFProtect(middleware.CSRFConfig{
		CookieName:      "huntiq_csrf",
		HeaderName:      "X-CSRF-Token",
		IsOriginAllowed: d.Config.IsAllowedOrigin,
		ExemptPaths: []string{
			"/api/chapa/webhook",
			"/api/telebirr/webhook",
			"/api/cbe/webhook",
			"/api/webhooks",
			"/api/ai/chat",     // AI chat should work for anonymous users
			"/api/ai/health",   // Health check is read-only
		},
	}))

	publicLimiter := limiter.New(limiter.Config{
		Max:          10,
		Expiration:   1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string { return c.IP() },
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"success": false,
				"error":   "too many requests — please wait a minute and try again",
			})
		},
	})

	searchLimiter := limiter.New(limiter.Config{
		Max:          60,
		Expiration:   1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string { return c.IP() },
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"success": false,
				"error":   "search rate limit reached — please slow down",
			})
		},
	})

	maxAuth := 10
	if d.Config != nil && d.Config.DevMode {
		maxAuth = 120
	}
	authLimiter := limiter.New(limiter.Config{
		Max:          maxAuth,
		Expiration:   1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string { return c.IP() },
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"success": false,
				"error":   "too many requests — please wait a minute and try again",
			})
		},
	})

	otpLimiter := limiter.New(limiter.Config{
		Max:          5,
		Expiration:   1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string { return c.IP() },
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"success": false,
				"error":   "too many OTP requests — please wait a minute",
			})
		},
	})

	referralLimiter := limiter.New(limiter.Config{
		Max:        3,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			if id, ok := c.Locals("userID").(int); ok {
				return fmt.Sprintf("referral:%d", id)
			}
			return "referral:" + c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"success": false,
				"error":   "too many referral code attempts — please wait a minute",
			})
		},
	})

	uploadLimiter := limiter.New(limiter.Config{
		Max:        20,
		Expiration: 1 * time.Hour,
		KeyGenerator: func(c *fiber.Ctx) string {
			if id, ok := c.Locals("userID").(int); ok {
				return fmt.Sprintf("upload:%d", id)
			}
			return "upload:" + c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"success": false,
				"error":   "upload limit reached — max 20 images per hour",
			})
		},
	})

	// ── Public Routes ────────────────────────────────────────────────────────
	api.Get("/products", searchLimiter, d.MW.OptionalAuth, d.ProductHandler.List)
	api.Get("/products/:id", d.MW.OptionalAuth, d.ProductHandler.Get)
	api.Get("/categories", d.ProductHandler.ListCategories)
	api.Get("/categories/roots", d.CategoryHandler.GetRoots)
	api.Get("/categories/tree", d.CategoryHandler.GetCategoryTree)
	api.Get("/categories/search", d.CategoryHandler.SearchCategories)
	api.Get("/categories/resolve", d.CategoryHandler.ResolveCategory)
	api.Get("/categories/popular", d.ProductHandler.PopularCategories)
	api.Get("/categories/check-slug", d.CategoryHandler.CheckSlug)
	api.Get("/categories/slug/:slug", d.CategoryHandler.GetCategoryBySlug)
	api.Get("/categories/:id", d.CategoryHandler.GetCategory)
	api.Get("/categories/:id/children", d.CategoryHandler.GetCategoryChildren)
	api.Get("/categories/:id/ancestors", d.CategoryHandler.GetCategoryBreadcrumbs)
	api.Get("/categories/:id/attributes", d.AttributeHandler.GetCategoryAttributes)

	// Universal Classification (Public)
	api.Get("/classification/schemes", d.ClassificationHandler.ListSchemes)
	api.Get("/classification/schemes/:code", d.ClassificationHandler.GetScheme)
	api.Get("/classification/schemes/:code/tree", d.ClassificationHandler.GetTree)
	api.Get("/classification/schemes/:code/roots", d.ClassificationHandler.GetRoots)
	api.Get("/classification/placements/:id", d.ClassificationHandler.GetPlacement)
	api.Get("/classification/placements/:id/children", d.ClassificationHandler.GetChildren)
	api.Get("/classification/placements/:id/attributes", d.ClassificationHandler.GetPlacementAttributes)
	api.Get("/classification/concepts", d.ClassificationHandler.SearchConcepts)
	api.Get("/classification/concepts/:id", d.ClassificationHandler.GetConcept)
	api.Get("/classification/concepts/:id/relationships", d.ClassificationHandler.GetRelationships)
	api.Get("/classification/relationship-types", d.ClassificationHandler.ListRelationshipTypes)
	api.Get("/classification/search", d.ClassificationHandler.SearchUniversal)

	api.Get("/attributes", d.AttributeHandler.ListAttributeDefinitions)
	api.Get("/products/:id/attributes", d.MW.OptionalAuth, d.AttributeHandler.GetProductAttributes)
	api.Post("/orders", publicLimiter, d.MW.OptionalAuth, d.SellerHandler.CreateOrder)
	api.Get("/fee-options", d.ProductHandler.FeeOptions)
	api.Get("/locations/regions", d.LocationHandler.ListRegions)
	api.Get("/locations/zones", d.LocationHandler.ListZones)
	api.Get("/locations/subcities", d.LocationHandler.ListSubcities)
	api.Get("/locations/woredas", d.LocationHandler.ListWoredas)
	api.Get("/locations/search", d.LocationHandler.SearchLocations)
	api.Post("/locations/reverse-geocode", d.LocationHandler.ReverseGeocode)
	api.Get("/locations/tree", d.LocationHandler.GetLocationTree)
	api.Get("/locations/children", d.LocationHandler.ListChildren)
	api.Get("/locations/search-hierarchy", d.LocationHandler.SearchHierarchy)
	api.Get("/banners", d.BannerHandler.ListPublic)
	api.Get("/public/payment-settings", d.AdminHandler.GetPublicPaymentSettings)
	api.Get("/public/settings", d.AdminHandler.GetPublicSettings)
	api.Post("/reports", d.BuyerHandler.CreateAbuseReport)
	api.Post("/appeals", d.BuyerHandler.CreateAccountAppeal)
	api.Post("/contact", publicLimiter, d.BuyerHandler.SubmitContactMessage)
	api.Get("/public/promotion", d.AdminHandler.GetPublicPromotion)

	// ── HUNTIQ AI Marketplace Intelligence Route ──────────────────────
	aiLimiter := limiter.New(limiter.Config{
		Max:        30,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string { return c.IP() },
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"success": false,
				"error":   "Too many AI inquiries. Please wait a moment before asking again.",
			})
		},
	})
	api.Post("/ai/chat", aiLimiter, d.MW.OptionalAuth, d.AIHandler.Chat)
	api.Get("/ai/health", d.AIHandler.Health)

	// ── Auth & Capabilities ──────────────────────────────────────────────────
	api.Post("/auth/otp/send", otpLimiter, d.AuthHandler.SendOTP)
	api.Post("/auth/otp/verify", otpLimiter, d.AuthHandler.VerifyOTP)
	api.Post("/auth/otp/login", authLimiter, d.AuthHandler.OTPLogin)
	api.Post("/auth/demand", d.MW.RequireAuth, d.AuthHandler.SetDemand)
	api.Post("/auth/register", authLimiter, d.AuthHandler.UnifiedRegister)
	api.Post("/auth/login", authLimiter, d.AuthHandler.UnifiedLogin)
	api.Post("/auth/refresh", authLimiter, d.AuthHandler.Refresh)
	api.Post("/auth/logout", d.AuthHandler.Logout)
	api.Post("/auth/logout-all", d.MW.RequireAuth, d.AuthHandler.LogoutAll)
	api.Get("/auth/me", d.MW.RequireAuth, d.AuthHandler.GetMe)
	api.Get("/auth/capabilities", d.MW.RequireAuth, d.AuthHandler.GetCapabilities)
	api.Post("/auth/capabilities/activate", d.MW.RequireAuth, d.AuthHandler.ActivateCapability)
	api.Post("/auth/forgot-password", authLimiter, d.AuthHandler.UnifiedForgotPassword)
	api.Post("/auth/verify-recovery-otp", authLimiter, d.AuthHandler.VerifyRecoveryOTP)
	api.Post("/auth/reset-password", authLimiter, d.AuthHandler.UnifiedResetPassword)

	// ── National Identity Verification (Fayda eKYC) ──────────────────────────
	api.Post("/identity/verify", d.MW.RequireAuth, d.IdentityHandler.VerifyIdentity)
	api.Get("/identity/status", d.MW.RequireAuth, d.IdentityHandler.GetMyIdentityStatus)

	// ── Business & Government Organizations ─────────────────────────────────
	api.Post("/organizations", d.MW.RequireAuth, d.OrganizationHandler.RegisterOrganization)
	api.Get("/organizations/me", d.MW.RequireAuth, d.OrganizationHandler.GetMyOrganization)

	api.Post("/seller/register", authLimiter, d.AuthHandler.Register)
	api.Post("/seller/login", authLimiter, d.AuthHandler.Login)
	api.Post("/seller/check", authLimiter, d.AuthHandler.CheckSellerExists)
	api.Post("/seller/forgot-password", authLimiter, d.AuthHandler.SellerForgotPassword)
	api.Post("/seller/reset-password", authLimiter, d.AuthHandler.SellerResetPassword)
	api.Post("/buyer/register", authLimiter, d.AuthHandler.BuyerRegister)
	api.Post("/buyer/login", authLimiter, d.AuthHandler.BuyerLogin)
	api.Post("/buyer/check", authLimiter, d.AuthHandler.CheckBuyerExists)
	api.Post("/buyer/forgot-password", authLimiter, d.AuthHandler.BuyerForgotPassword)
	api.Post("/buyer/reset-password", authLimiter, d.AuthHandler.BuyerResetPassword)
	api.Post("/admin/login", authLimiter, d.AuthHandler.AdminLogin)
	api.Post("/admin/login/verify", otpLimiter, d.AuthHandler.AdminLoginVerifyOTP)
	api.Post("/admin/forgot-password", authLimiter, d.AuthHandler.AdminForgotPassword)
	api.Post("/admin/reset-password", authLimiter, d.AuthHandler.AdminResetPassword)

	// ── WebAuthn Passkeys ───────────────────────────────────────────────────
	api.Post("/admin/passkeys/login/begin", d.AdminAuthHandler.BeginLogin)
	api.Post("/admin/auth/passkey/login/begin", d.AdminAuthHandler.BeginLogin)
	api.Post("/admin/passkeys/login/finish", d.AdminAuthHandler.FinishLogin)
	api.Post("/admin/auth/passkey/login/verify", d.AdminAuthHandler.FinishLogin)
	api.Post("/admin/passkeys/register/begin", d.MW.RequireAuth, middleware.RequireRole("admin"), d.AdminAuthHandler.BeginRegistration)
	api.Post("/admin/auth/passkey/register/begin", d.MW.RequireAuth, middleware.RequireRole("admin"), d.AdminAuthHandler.BeginRegistration)
	api.Post("/admin/passkeys/register/finish", d.MW.RequireAuth, middleware.RequireRole("admin"), d.AdminAuthHandler.FinishRegistration)
	api.Post("/admin/auth/passkey/register/verify", d.MW.RequireAuth, middleware.RequireRole("admin"), d.AdminAuthHandler.FinishRegistration)
	api.Get("/admin/passkeys", d.MW.RequireAuth, middleware.RequireRole("admin"), d.AdminAuthHandler.ListPasskeys)
	api.Get("/admin/auth/passkey/list", d.MW.RequireAuth, middleware.RequireRole("admin"), d.AdminAuthHandler.ListPasskeys)
	api.Delete("/admin/passkeys/:id", d.MW.RequireAuth, middleware.RequireRole("admin"), d.AdminAuthHandler.RevokePasskey)
	api.Delete("/admin/auth/passkey/:id", d.MW.RequireAuth, middleware.RequireRole("admin"), d.AdminAuthHandler.RevokePasskey)
	api.Post("/admin/reauth/begin", d.MW.RequireAuth, middleware.RequireRole("admin"), d.AdminAuthHandler.BeginReauth)
	api.Post("/admin/reauth/finish", d.MW.RequireAuth, middleware.RequireRole("admin"), d.AdminAuthHandler.FinishReauth)
	api.Post("/admin/auth/reauth", d.MW.RequireAuth, middleware.RequireRole("admin"), d.AdminAuthHandler.ReauthPassword)
	api.Post("/admin/logout", d.AdminAuthHandler.Logout)
	api.Post("/admin/auth/session/revoke", d.MW.RequireAuth, middleware.RequireRole("admin"), d.AdminAuthHandler.RevokeSessionEndpoint)

	// ── Seller Routes ────────────────────────────────────────────────────────
	seller := api.Group("/seller", d.MW.RequireAuth, middleware.RequireRole("seller"))
	seller.Get("/stats", d.SellerHandler.GetStats)
	seller.Get("/fee-preview", d.ProductHandler.FeePreview)
	seller.Get("/products", d.ProductHandler.SellerList)
	seller.Get("/products/detailed", d.SellerHandler.ListProductsDetailed)
	seller.Post("/products", d.ProductHandler.SellerCreate)
	seller.Put("/products/:id", d.ProductHandler.SellerUpdate)
	seller.Delete("/products/:id", d.SellerHandler.DeleteProduct)
	seller.Get("/transactions", d.MW.RequireSellerTier("gold"), d.SellerHandler.ListTransactions)
	seller.Get("/orders", d.MW.RequireSellerTier("gold"), d.SellerHandler.ListOrders)
	seller.Put("/orders/:id", d.MW.RequireSellerTier("gold"), d.SellerHandler.UpdateOrder)
	seller.Get("/profile", d.SellerHandler.GetProfile)
	seller.Put("/profile", d.SellerHandler.UpdateProfile)
	seller.Put("/preference", d.SellerHandler.UpdatePreference)
	seller.Put("/password", d.AuthHandler.SellerChangePassword)
	seller.Post("/listing-fee", d.SellerHandler.PayListingFee)
	seller.Post("/upgrade-tier", d.SellerHandler.RequestUpgrade)
	seller.Get("/upgrade-tier", d.SellerHandler.GetUpgradeRequests)
	seller.Get("/notifications", d.SellerHandler.ListNotifications)
	seller.Put("/notifications/:id/read", d.SellerHandler.MarkNotificationRead)
	seller.Get("/referral", d.ReferralHandler.GetSellerStats)
	seller.Get("/referral/commissions", d.ReferralHandler.GetCommissions)
	seller.Post("/referral/accept", referralLimiter, d.ReferralHandler.AcceptInvitation)
	seller.Get("/referral/balance", d.ReferralHandler.GetBalance)
	seller.Post("/referral/withdraw", d.ReferralHandler.RequestWithdrawal)
	seller.Get("/referral/withdrawals", d.ReferralHandler.GetMyWithdrawals)
	seller.Post("/products/:id/attributes", d.AttributeHandler.SetProductAttributes)
	seller.Get("/verification/documents", d.SellerHandler.ListVerificationDocuments)
	seller.Post("/verification/documents", d.SellerHandler.UploadVerificationDocument)
	seller.Post("/refunds", d.SellerHandler.CreateRefundRequest)

	api.Post("/upload/image", d.MW.RequireAuth, uploadLimiter, handlers.UploadImage)
	api.Put("/auth/change-password", d.MW.RequireAuth, authLimiter, d.AuthHandler.SellerChangePassword)
	api.Post("/transactions", d.MW.RequireAuth, middleware.RequireRole("seller"), authLimiter, d.TransactionHandler.Create)

	// ── Financial Wallet & Telebirr Payments ──────────────────────────────────
	api.Get("/wallet", d.MW.RequireAuth, d.WalletHandler.GetMyWallet)
	api.Get("/wallet/ledger", d.MW.RequireAuth, d.WalletHandler.GetMyLedger)
	api.Get("/referrals/token-link", d.MW.RequireAuth, d.WalletHandler.GetMyReferralLink)
	api.Post("/referrals/click", d.WalletHandler.TrackReferralClick)
	api.Post("/telebirr/webhook", d.TelebirrHandler.HandleWebhook)
	api.Get("/telebirr/orders/:outTradeNo", d.TelebirrHandler.GetOrderStatus)
	api.Post("/telebirr/orders", d.MW.RequireAuth, d.TelebirrHandler.CreateOrder)
	api.Get("/telebirr/my-orders", d.MW.RequireAuth, d.TelebirrHandler.GetUserOrders)

	// ── Buyer Routes ─────────────────────────────────────────────────────────
	buyer := api.Group("/buyer", d.MW.RequireAuth, middleware.RequireRole("buyer"))
	buyer.Post("/upgrade-tier", d.BuyerHandler.UpgradeTier)
	buyer.Get("/upgrade-tier", d.BuyerHandler.GetUpgradeRequests)
	buyer.Get("/profile", d.BuyerHandler.GetProfile)
	buyer.Put("/profile", d.BuyerHandler.UpdateProfile)
	buyer.Get("/orders", d.BuyerHandler.GetOrders)
	buyer.Put("/password", d.BuyerHandler.ChangePassword)
	buyer.Get("/notifications", d.BuyerHandler.ListNotifications)
	buyer.Put("/notifications/:id/read", d.BuyerHandler.MarkNotificationRead)
	buyer.Get("/referral", d.ReferralHandler.GetBuyerStats)
	buyer.Get("/referral/commissions", d.ReferralHandler.GetCommissions)
	buyer.Post("/referral/accept", referralLimiter, d.ReferralHandler.AcceptInvitation)
	buyer.Get("/referral/balance", d.ReferralHandler.GetBalance)
	buyer.Post("/referral/withdraw", d.ReferralHandler.RequestWithdrawal)
	buyer.Get("/referral/withdrawals", d.ReferralHandler.GetMyWithdrawals)
	buyer.Post("/refunds", d.BuyerHandler.CreateRefundRequest)

	// ── Admin Routes ─────────────────────────────────────────────────────────
	admin := api.Group("/admin", d.MW.RequireAuth, middleware.RequireRole("admin"))
	admin.Get("/me", d.AdminHandler.GetMe)
	admin.Get("/me/permissions", d.AdminHandler.GetMyPermissions)
	admin.Put("/password", d.AdminHandler.ChangePassword)

	// Multi-Computer Device Identity Governance
	admin.Get("/devices", d.AdminAuthHandler.ListDevices)
	admin.Post("/devices/register", d.AdminAuthHandler.RegisterDevice)
	admin.Delete("/devices/:id", d.AdminAuthHandler.RevokeDevice)

	admin.Get("/products/pending", d.MW.RequirePermission("pending"), d.AdminHandler.ListPending)
	admin.Get("/products", d.MW.RequirePermission("products"), d.AdminHandler.ListAllProducts)
	admin.Put("/products/:id/review", d.MW.RequirePermission("pending"), d.AdminHandler.ReviewProduct)
	admin.Delete("/products/:id", d.MW.RequirePermission("products"), d.AdminHandler.DeleteProduct)

	admin.Get("/sellers", d.MW.RequirePermission("sellers"), d.AdminHandler.ListSellers)
	admin.Put("/sellers/:id/active", d.MW.RequirePermission("sellers"), d.AdminHandler.SetSellerActive)

	admin.Get("/locations", d.MW.RequirePermission("locations"), d.AdminHandler.ListLocations)
	admin.Post("/locations", d.MW.RequirePermission("locations"), d.AdminHandler.AddLocation)
	admin.Put("/locations/:id", d.MW.RequirePermission("locations"), d.AdminHandler.UpdateLocation)
	admin.Delete("/locations/:id", d.MW.RequirePermission("locations"), d.AdminHandler.DeleteLocation)

	admin.Get("/location-hierarchy", d.MW.RequirePermission("locations"), d.LocationHandler.AdminListLocations)
	admin.Get("/location-hierarchy/:id", d.MW.RequirePermission("locations"), d.LocationHandler.AdminGetLocation)
	admin.Post("/location-hierarchy", d.MW.RequirePermission("locations"), d.LocationHandler.AdminCreateLocation)
	admin.Put("/location-hierarchy/:id", d.MW.RequirePermission("locations"), d.LocationHandler.AdminUpdateLocation)
	admin.Delete("/location-hierarchy/:id", d.MW.RequirePermission("locations"), d.LocationHandler.AdminDeleteLocation)
	admin.Post("/location-hierarchy/import", d.MW.RequirePermission("locations"), d.LocationHandler.AdminImportLocations)
	admin.Post("/location-hierarchy/test-geocode", d.MW.RequirePermission("locations"), d.LocationHandler.ReverseGeocode)

	admin.Get("/analytics/overview", d.MW.RequirePermission("analytics"), d.AdminHandler.AnalyticsOverview)
	admin.Get("/analytics/heatmap", d.MW.RequirePermission("analytics"), d.AdminHandler.AnalyticsHeatmap)
	admin.Get("/analytics/monthly", d.MW.RequirePermission("analytics"), d.AdminHandler.AnalyticsMonthly)
	admin.Get("/analytics/gantt", d.MW.RequirePermission("analytics"), d.AdminHandler.AnalyticsGantt)
	admin.Get("/analytics/fee-revenue", d.MW.RequireAnyPermission("analytics", "finance", "profit"), d.AdminHandler.AnalyticsFeeRevenue)

	admin.Get("/transactions", d.MW.RequirePermission("transactions"), d.AdminHandler.ListTransactions)
	admin.Put("/transactions/:id/status", d.MW.RequirePermission("transactions"), d.TransactionHandler.UpdateStatus)

	admin.Get("/profit/summary", d.MW.RequirePermission("profit"), d.AdminHandler.ProfitSummary)
	admin.Get("/profit/by-tier", d.MW.RequirePermission("profit"), d.AdminHandler.ProfitByTier)
	admin.Get("/profit/monthly", d.MW.RequirePermission("profit"), d.AdminHandler.ProfitMonthly)
	admin.Get("/profit/weekly", d.MW.RequirePermission("profit"), d.AdminHandler.ProfitWeekly)
	admin.Get("/profit/annual", d.MW.RequirePermission("profit"), d.AdminHandler.ProfitAnnual)

	financeView := d.MW.RequireAnyPermission("profit", "finance", "finance:view")
	financeManage := d.MW.RequireAnyPermission("finance:manage", "finance_manage")
	admin.Get("/finance/forecast", financeView, d.FPAHandler.GetForecast)
	admin.Get("/finance/scenarios", financeView, d.FPAHandler.GetScenarios)
	admin.Get("/finance/capacity-20m", financeView, d.FPAHandler.GetCapacityPlanning20M)
	admin.Get("/finance/tax-rules", financeView, d.FPAHandler.GetTaxRules)
	admin.Get("/finance/actuals", financeView, d.FPAHandler.GetActuals)
	admin.Get("/finance/audit-logs", financeView, d.FPAHandler.GetAuditLogs)
	admin.Put("/finance/assumptions", financeManage, d.FPAHandler.UpdateAssumptions)
	admin.Post("/finance/assumptions/reset", financeManage, d.FPAHandler.ResetAssumptions)
	admin.Get("/finance/export", financeView, d.FPAHandler.ExportForecast)

	admin.Get("/fee-payments/pending", d.MW.RequireAnyPermission("products", "transactions", "pending"), d.AdminHandler.ListPendingFeePayments)
	admin.Put("/fee-payments/:id/verify", d.MW.RequireAnyPermission("products", "transactions", "pending"), d.AdminHandler.VerifyListingFeePayment)
	admin.Put("/fee-payments/:id/reject", d.MW.RequireAnyPermission("products", "transactions", "pending"), d.AdminHandler.RejectListingFeePayment)

	admin.Get("/settings", d.MW.RequirePermission("settings"), d.AdminHandler.GetSettings)
	admin.Put("/settings", d.MW.RequirePermission("settings"), d.AdminHandler.UpdateSettings)
	admin.Put("/settings/with-reason", d.MW.RequirePermission("settings"), d.AdminHandler.UpdateSettingsWithReason)
	admin.Get("/settings/log", d.MW.RequirePermission("settings"), d.AdminHandler.SettingsChangeLog)

	admin.Post("/users/:id/warn", d.MW.RequirePermission("sellers"), d.AdminHandler.IssueWarning)
	admin.Put("/users/:id/ban", middleware.RequirePlatformOwner, d.AdminHandler.BanUser)
	admin.Get("/system/health", d.MW.RequireAnyPermission("operations", "settings"), d.AdminHandler.SystemHealth)
	admin.Post("/system/health/repair", middleware.RequirePlatformAdmin, d.AdminHandler.RepairSystemHealth)

	admin.Get("/notifications", d.MW.RequirePermission("notifications"), d.AdminHandler.ListNotifications)
	admin.Put("/notifications/:id/read", d.MW.RequirePermission("notifications"), d.AdminHandler.MarkNotificationRead)
	admin.Post("/notifications", d.MW.RequirePermission("notifications"), d.AdminHandler.CreateNotification)

	admin.Get("/categories", d.MW.RequirePermission("categories"), d.AdminHandler.ListCategories)
	admin.Get("/category-hierarchy", d.MW.RequirePermission("categories"), d.CategoryHandler.ListCategories)
	admin.Get("/category-hierarchy/:id/stats", d.MW.RequirePermission("categories"), d.CategoryHandler.GetCategoryStats)
	admin.Post("/category-hierarchy", d.MW.RequirePermission("categories"), d.CategoryHandler.CreateCategory)
	admin.Put("/category-hierarchy/:id", d.MW.RequirePermission("categories"), d.CategoryHandler.UpdateCategory)
	admin.Put("/category-hierarchy/:id/status", d.MW.RequirePermission("categories"), d.CategoryHandler.SetCategoryStatus)
	admin.Put("/category-hierarchy/:id/reorder", d.MW.RequirePermission("categories"), d.CategoryHandler.ReorderCategory)
	admin.Delete("/category-hierarchy/:id", d.MW.RequirePermission("categories"), d.CategoryHandler.DeleteCategory)

	admin.Get("/classification/schemes", d.MW.RequirePermission("categories"), d.ClassificationHandler.ListSchemes)
	admin.Post("/classification/schemes", d.MW.RequirePermission("categories"), d.ClassificationHandler.CreateScheme)
	admin.Put("/classification/schemes/:id", d.MW.RequirePermission("categories"), d.ClassificationHandler.UpdateScheme)
	admin.Post("/classification/placements", d.MW.RequirePermission("categories"), d.ClassificationHandler.CreatePlacement)
	admin.Put("/classification/placements/:id", d.MW.RequirePermission("categories"), d.ClassificationHandler.UpdatePlacement)
	admin.Put("/classification/placements/:id/move", d.MW.RequirePermission("categories"), d.ClassificationHandler.MovePlacement)
	admin.Delete("/classification/placements/:id", d.MW.RequirePermission("categories"), d.ClassificationHandler.DeletePlacement)
	admin.Post("/classification/concepts", d.MW.RequirePermission("categories"), d.ClassificationHandler.CreateConcept)
	admin.Put("/classification/concepts/:id", d.MW.RequirePermission("categories"), d.ClassificationHandler.UpdateConcept)
	admin.Post("/classification/concepts/check-duplicates", d.MW.RequirePermission("categories"), d.ClassificationHandler.CheckDuplicateCandidates)
	admin.Post("/classification/relationships", d.MW.RequirePermission("categories"), d.ClassificationHandler.CreateRelationship)
	admin.Delete("/classification/relationships/:id", d.MW.RequirePermission("categories"), d.ClassificationHandler.DeleteRelationship)
	admin.Post("/classification/placements/:id/attributes", d.MW.RequirePermission("categories"), d.ClassificationHandler.SaveAttributeRule)
	admin.Delete("/classification/attribute-rules/:id", d.MW.RequirePermission("categories"), d.ClassificationHandler.DeleteAttributeRule)
	admin.Get("/classification/impact-analysis", d.MW.RequirePermission("categories"), d.ClassificationHandler.AnalyzeImpact)
	admin.Get("/classification/audit-logs", d.MW.RequireAnyPermission("logs", "categories"), d.ClassificationHandler.ListAuditLogs)

	admin.Get("/attributes", d.MW.RequirePermission("categories"), d.AttributeHandler.ListAttributeDefinitions)
	admin.Post("/attributes", d.MW.RequirePermission("categories"), d.AttributeHandler.CreateAttributeDefinition)
	admin.Post("/categories/:id/attributes", d.MW.RequirePermission("categories"), d.AttributeHandler.AssignAttributeToCategory)
	admin.Delete("/categories/:id/attributes/:attribute_id", d.MW.RequirePermission("categories"), d.AttributeHandler.RemoveAttributeFromCategory)

	// Master Data Management (Locations & Categories JSON Exchange Architecture)
	admin.Get("/masterdata/status", d.MW.RequireAnyPermission("locations", "categories"), d.MasterDataHandler.GetStatus)
	admin.Post("/masterdata/validate", d.MW.RequireAnyPermission("locations", "categories"), d.MasterDataHandler.ValidatePayload)
	admin.Post("/masterdata/import/locations", d.MW.RequirePermission("locations"), d.MasterDataHandler.ImportLocations)
	admin.Post("/masterdata/import/categories", d.MW.RequirePermission("categories"), d.MasterDataHandler.ImportCategories)
	admin.Get("/masterdata/export/locations", d.MW.RequirePermission("locations"), d.MasterDataHandler.ExportLocations)
	admin.Get("/masterdata/export/categories", d.MW.RequirePermission("categories"), d.MasterDataHandler.ExportCategories)

	admin.Get("/admins", middleware.RequirePlatformOwner, d.AdminHandler.ListAdmins)
	admin.Post("/admins", middleware.RequirePlatformOwner, d.AdminHandler.CreateAdmin)
	admin.Put("/admins/:id/role", middleware.RequirePlatformOwner, d.AdminHandler.UpdateAdminRole)
	admin.Delete("/admins/:id", middleware.RequirePlatformOwner, d.AdminHandler.DeleteAdmin)
	admin.Get("/admins/:id/permissions", middleware.RequirePlatformOwner, d.AdminHandler.GetAdminPermissions)
	admin.Put("/admins/:id/permissions", middleware.RequirePlatformOwner, d.AdminHandler.SetAdminPermissions)

	admin.Get("/banners", d.MW.RequirePermission("banners"), d.BannerHandler.ListAll)
	admin.Post("/banners", d.MW.RequirePermission("banners"), d.BannerHandler.Create)
	admin.Put("/banners/:id", d.MW.RequirePermission("banners"), d.BannerHandler.Update)
	admin.Delete("/banners/:id", d.MW.RequirePermission("banners"), d.BannerHandler.Delete)

	admin.Get("/logs", d.MW.RequireAnyPermission("logs"), d.AdminHandler.ListAdminLogs)
	admin.Get("/otp-logs", d.MW.RequireAnyPermission("otp_logs"), d.AuthHandler.ListOTPLogs)
	admin.Get("/activity-logs", d.MW.RequireAnyPermission("logs", "otp_logs"), d.AuthHandler.ListActivityLogs)

	admin.Put("/users/:id/seller-tier", d.MW.RequireAnyPermission("seller_upgrades", "buyer_upgrades"), d.AdminHandler.SetSellerTier)
	admin.Put("/users/:id/buyer-tier", d.MW.RequireAnyPermission("seller_upgrades", "buyer_upgrades"), d.AdminHandler.SetBuyerTier)
	admin.Get("/users/:id/tier-history", d.MW.RequirePermission("premium_users"), d.AdminHandler.GetTierHistory)
	admin.Get("/premium-users", d.MW.RequirePermission("premium_users"), d.AdminHandler.ListPremiumUsers)

	admin.Get("/upgrade-requests", d.MW.RequirePermission("seller_upgrades"), d.AdminHandler.ListUpgradeRequests)
	admin.Put("/upgrade-requests/:id/approve", d.MW.RequirePermission("seller_upgrades"), d.AdminHandler.ApproveUpgradeRequest)
	admin.Put("/upgrade-requests/:id/reject", d.MW.RequirePermission("seller_upgrades"), d.AdminHandler.RejectUpgradeRequest)

	admin.Get("/buyer-upgrade-requests", d.MW.RequirePermission("buyer_upgrades"), d.AdminHandler.ListBuyerUpgradeRequests)
	admin.Put("/buyer-upgrade-requests/:id/approve", d.MW.RequirePermission("buyer_upgrades"), d.AdminHandler.ApproveBuyerUpgradeRequest)
	admin.Put("/buyer-upgrade-requests/:id/reject", d.MW.RequirePermission("buyer_upgrades"), d.AdminHandler.RejectBuyerUpgradeRequest)

	admin.Get("/categories/popularity", d.MW.RequirePermission("category_trends"), d.AdminHandler.CategoryPopularity)
	admin.Get("/referral/overview", d.MW.RequirePermission("referrals"), d.ReferralHandler.AdminOverview)
	admin.Get("/referral/list", d.MW.RequirePermission("referrals"), d.ReferralHandler.AdminList)
	admin.Put("/referral/users/:id/suspend", d.MW.RequireAnyPermission("referrals", "trust_safety"), d.ReferralHandler.AdminSetSuspended)
	admin.Put("/referral/users/:id/block", d.MW.RequireAnyPermission("referrals", "trust_safety"), d.ReferralHandler.AdminSetBlocked)
	admin.Get("/referral/withdrawals", d.MW.RequirePermissionOrSuper("withdrawals"), d.ReferralHandler.AdminListWithdrawals)
	admin.Put("/referral/withdrawals/:id/approve", d.MW.RequirePermissionOrSuper("withdrawals"), d.ReferralHandler.AdminApproveWithdrawal)
	admin.Put("/referral/withdrawals/:id/reject", d.MW.RequirePermissionOrSuper("withdrawals"), d.ReferralHandler.AdminRejectWithdrawal)
	admin.Put("/referral/withdrawals/:id/paid", d.MW.RequirePermissionOrSuper("withdrawals"), d.ReferralHandler.AdminMarkWithdrawalPaid)

	// ── Financial Ledger Governance & SMS Audit ──────────────────────────────
	admin.Get("/financial/ledger", d.MW.RequireAnyPermission("finance", "analytics"), d.AdminFinancialHandler.GetLedgerAuditTrail)
	admin.Get("/referral/policy", d.MW.RequireAnyPermission("finance", "finance:view", "referrals"), d.AdminFinancialHandler.GetRewardPolicy)
	admin.Put("/referral/policy", d.MW.RequireAnyPermission("finance:manage", "platform_owner"), middleware.RequirePortal("executive"), d.AdminFinancialHandler.UpdateRewardPolicy)
	admin.Get("/sms/logs", d.MW.RequireAnyPermission("logs", "operations"), d.AdminFinancialHandler.GetSMSDeliveryLogs)
	admin.Post("/financial/compensate", d.MW.RequireAnyPermission("finance:manage", "platform_owner"), middleware.RequirePortal("executive"), d.AdminFinancialHandler.CompensateLedger)

	admin.Get("/operations/health", d.MW.RequirePermission("operations"), middleware.RequirePortal("platform"), d.AdminHandler.GetOperationsHealth)
	admin.Get("/operations/jobs", d.MW.RequirePermission("operations"), middleware.RequirePortal("platform"), d.AdminHandler.GetJobStatusSummary)
	admin.Get("/operations/failed-jobs", d.MW.RequirePermission("operations"), middleware.RequirePortal("platform"), d.AdminHandler.GetFailedJobs)

	// ── Identity & Organization Governance ──
	admin.Get("/identity/pending", d.MW.RequireAnyPermission("sellers", "trust_safety", "operations"), d.IdentityHandler.AdminListPending)
	admin.Put("/identity/:id/resolve", d.MW.RequireAnyPermission("sellers", "trust_safety"), d.IdentityHandler.AdminResolveException)
	admin.Get("/organizations", d.MW.RequireAnyPermission("sellers", "operations", "trust_safety"), d.OrganizationHandler.AdminListOrganizations)
	admin.Put("/organizations/:id/review", d.MW.RequireAnyPermission("sellers", "trust_safety"), d.OrganizationHandler.AdminReviewOrganization)

	admin.Get("/verification/documents", d.MW.RequirePermission("seller_verification"), d.AdminHandler.ListVerificationDocuments)
	admin.Put("/verification/documents/:id/review", d.MW.RequirePermission("seller_verification"), d.AdminHandler.ReviewVerificationDocument)

	admin.Get("/trust/reports", d.MW.RequirePermission("trust_safety"), middleware.RequirePortal("security"), d.AdminHandler.ListAbuseReports)
	admin.Put("/trust/reports/:id/review", d.MW.RequirePermission("trust_safety"), middleware.RequirePortal("security"), d.AdminHandler.ReviewAbuseReport)
	admin.Put("/trust/reports/:id/priority", d.MW.RequirePermission("trust_safety"), middleware.RequirePortal("security"), d.AdminHandler.UpdateAbuseReportPriority)

	admin.Get("/refunds", d.MW.RequirePermission("refunds"), d.AdminHandler.ListRefundRequests)
	admin.Put("/refunds/:id/review", d.MW.RequirePermission("refunds"), d.AdminHandler.ReviewRefundRequest)
	admin.Put("/refunds/:id/processed", d.MW.RequirePermission("refunds"), d.AdminHandler.MarkRefundProcessed)

	admin.Get("/appeals", d.MW.RequirePermission("appeals"), d.AdminHandler.ListAccountAppeals)
	admin.Put("/appeals/:id/review", d.MW.RequirePermission("appeals"), d.AdminHandler.ReviewAccountAppeal)
	admin.Put("/appeals/:id/status", d.MW.RequirePermission("appeals"), d.AdminHandler.UpdateAppealStatus)

	admin.Get("/search/popular", d.MW.RequirePermission("analytics"), d.SearchHandler.PopularSearches)
	admin.Get("/search/zero-results", d.MW.RequirePermission("analytics"), d.SearchHandler.ZeroResultQueries)

	if d.Config.DevMode {
		api.Get("/dev/otp", d.AuthHandler.DevPeekOTP)
		api.Get("/dev/otp/decrypt", d.AuthHandler.DevDecryptOTP)
		log.Println("⚠️  DEV MODE: /api/dev/otp and /api/dev/otp/decrypt endpoints enabled")
	}
}

