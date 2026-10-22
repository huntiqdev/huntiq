package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/joho/godotenv"

	handlers "github.com/your-org/used-tools-marketplace/internal/handler"
	jobs "github.com/your-org/used-tools-marketplace/internal/worker"
	"github.com/your-org/used-tools-marketplace/internal/middleware"
	"github.com/your-org/used-tools-marketplace/internal/models"
	"github.com/your-org/used-tools-marketplace/internal/repository"
	"github.com/your-org/used-tools-marketplace/internal/service"
	"github.com/your-org/used-tools-marketplace/internal/security"
	"github.com/your-org/used-tools-marketplace/internal/config"
	"github.com/your-org/used-tools-marketplace/internal/database"
)

func main() {
	_ = godotenv.Load()

	cfg := config.Load()
	handlers.AppBaseURL = cfg.AppBaseURL

	// Initialize OTP crypto (no global state)
	otpCrypto, err := security.NewOTPCrypto(cfg.OTPEncryptKey)
	if err != nil {
		log.Fatalf("OTP crypto initialization failed: %v", err)
	}

	db, err := database.NewPool(cfg)
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	defer db.Close()

	// Context for background jobs — cancelled on shutdown
	jobCtx, jobCancel := context.WithCancel(context.Background())
	defer jobCancel()

	jobs.StartExpiryJob(jobCtx, db)
	jobs.StartUnverifiedUserCleanup(jobCtx, db)
	jobs.StartOTPCleanup(jobCtx, db)
	jobs.StartCommissionConfirmJob(jobCtx, db)

	// Repositories with OTP crypto injected
	productRepo := repository.NewProductRepo(db)
	categoryRepo := repository.NewCategoryRepo(db)
	categoryHierarchyRepo := repository.NewCategoryHierarchyRepo(db)
	attributeRepo := repository.NewAttributeRepo(db)
	searchRepo := repository.NewSearchRepo(db)
	locationRepo := repository.NewLocationRepo(db)
	userRepo := repository.NewUserRepo(db).WithOTPCrypto(otpCrypto)
	txRepo := repository.NewTransactionRepo(db)
	adminRepo := repository.NewAdminRepo(db)
	sellerRepo := repository.NewSellerRepo(db)
	bannerRepo := repository.NewBannerRepo(db)
	referralRepo := repository.NewReferralRepo(db)
	
	// Search repo for analytics (used by product handler)
	_ = searchRepo // Reserved for future search analytics integration
	
	// Search handler for admin analytics
	searchHandler := handlers.NewSearchHandler(searchRepo)

	// Wire is_active check so deactivated users can't use stale JWTs.
	// Results are cached for 30 s to avoid a DB query on every authenticated request.
	// Cache is bounded to 10,000 entries — oldest entries are evicted when full.
	// A background goroutine evicts expired entries every 5 minutes.
	const isActiveCacheMax = 10_000
	var (
		isActiveMu    sync.Mutex
		isActiveCache = map[int]struct {
			active  bool
			expires time.Time
		}{}
	)
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				now := time.Now()
				isActiveMu.Lock()
				for id, e := range isActiveCache {
					if now.After(e.expires) {
						delete(isActiveCache, id)
					}
				}
				isActiveMu.Unlock()
			case <-jobCtx.Done():
				return
			}
		}
	}()
	// IsActiveChecker is wired after services are initialized (see below) so authSvc is in scope.

	// Services with OTP crypto injected
	productSvc := service.NewProductService(productRepo, categoryRepo, locationRepo, userRepo).
		WithBackendBase(cfg.AppBaseURL).
		WithAttributeRepo(attributeRepo)
	locationSvc := service.NewLocationService(locationRepo)
	categoryHierarchySvc := service.NewCategoryHierarchyService(categoryHierarchyRepo)
	attributeSvc := service.NewAttributeService(attributeRepo)
	authSvc := service.NewAuthService(userRepo, cfg.JWTSecret).
		WithDevMode(cfg.DevMode).
		WithOTPCrypto(otpCrypto)
	txSvc := service.NewTransactionService(txRepo, productRepo, userRepo)
	adminSvc := service.NewAdminService(adminRepo)
	sellerSvc := service.NewSellerService(sellerRepo)
	referralSvc := service.NewReferralService(referralRepo).
		WithSettingsReader(func(ctx context.Context, key string) string {
			return adminSvc.GetSetting(ctx, key)
		})

	// Inject referral service into dependent services
	txSvc.WithReferralService(referralSvc)
	authSvc.WithReferralService(referralSvc)
	adminSvc.WithReferralService(referralSvc)
	adminSvc.WithSellerService(sellerSvc)

	// Inject admin service into auth so first-login notifications work
	authSvc.WithAdminService(adminSvc)

	// Initialize middleware with dependencies (no global state!)
	mw := middleware.New(cfg.JWTSecret).
		WithIsActiveChecker(func(userID int) (bool, error) {
			isActiveMu.Lock()
			if e, ok := isActiveCache[userID]; ok && time.Now().Before(e.expires) {
				isActiveMu.Unlock()
				return e.active, nil
			}
			// Evict ~50% of entries randomly to prevent wiping the entire cache (thundering herd)
			// under a burst of unique user IDs (e.g. scraper traffic).
			if len(isActiveCache) >= isActiveCacheMax {
				i := 0
				for id := range isActiveCache {
					if i%2 == 0 {
						delete(isActiveCache, id)
					}
					i++
				}
			}
			isActiveMu.Unlock()

			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			u, err := userRepo.GetByID(ctx, userID)
			if err != nil {
				return false, err
			}

			active := u.IsActive
			// Strictly enforce OTP verification for all authenticated routes.
			// If the user is not OTP-verified, return a specific error code so the
			// frontend can redirect to the OTP step instead of the login page.
			if u.Role != "admin" && !u.OTPVerified && authSvc.IsOTPRequired(ctx) {
				active = false
			}

			isActiveMu.Lock()
			isActiveCache[userID] = struct {
				active  bool
				expires time.Time
			}{active, time.Now().Add(10 * time.Second)} // 10s matches frontend JWT grace period in core.ts
			isActiveMu.Unlock()
			return active, nil
		}).
		WithSellerTierChecker(func(sellerID int) (string, error) {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			u, err := userRepo.GetByID(ctx, sellerID)
			if err != nil {
				return "normal", err
			}
			if u.Tier == "" {
				return "normal", nil
			}
			return u.Tier, nil
		})

	// Handlers
	productHandler := handlers.NewProductHandler(productSvc, authSvc)
	locationHandler := handlers.NewLocationHandler(locationSvc)
	categoryHandler := handlers.NewCategoryHandler(categoryHierarchySvc)
	attributeHandler := handlers.NewAttributeHandler(attributeSvc)
	authHandler := handlers.NewAuthHandler(authSvc).WithAdminService(adminSvc).WithDevMode(cfg.DevMode)
	txHandler := handlers.NewTransactionHandler(txSvc)
	adminHandler := handlers.NewAdminHandler(adminSvc, authSvc).WithCategoryRepo(categoryRepo).WithProductRepo(productRepo)
	sellerHandler := handlers.NewSellerHandler(sellerSvc).WithAdminService(adminSvc)
	bannerHandler := handlers.NewBannerHandler(bannerRepo)
	buyerHandler  := handlers.NewBuyerHandler(authSvc).WithAdminService(adminSvc)
	referralHandler := handlers.NewReferralHandler(referralSvc, cfg.FrontendURL).
		WithAdminNotifier(func(ctx context.Context, userID int, amount float64) {
			notifyCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			msg := fmt.Sprintf("User #%d has requested a referral commission withdrawal of ETB %.2f. Review in the Referrals panel.", userID, amount) //nolint:gosec // G117/G201: internal values only, no user-controlled input reaches this format string
			// target_role="" means all admins see it; withdrawals-permission admins can act on it.
			_ = adminSvc.CreateNotification(notifyCtx, "Withdrawal Request", msg, "warning", "")
		}).
		WithAdminLogger(func(ctx context.Context, adminID int, targetUserID *int, action, notes string) {
			adminSvc.InsertAdminLog(ctx, models.AdminLog{
				AdminID:      adminID,
				TargetUserID: targetUserID,
				Action:       action,
				Notes:        notes,
				IPAddress:    "", // IP not available in this context without extra plumbing
			})
		})

	// Wire withdrawal settings reader so ReferralService reads limits from DB
	// (already injected above via WithSettingsReader — no global state needed)

	// Inject permission and maintenance checkers into middleware
	mw.WithPermissionChecker(func(adminID int, permission string) (bool, error) {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		return adminSvc.HasPermission(ctx, adminID, permission)
	}).WithMaintenanceChecker(func() (bool, error) {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		v := adminSvc.GetSetting(ctx, "maintenance_mode")
		return v == "true", nil
	})

	app := fiber.New(fiber.Config{
		BodyLimit:   22 * 1024 * 1024, // 22 MB — allows up to 20 MB image uploads with multipart overhead //nolint:mnd
		ProxyHeader: fiber.HeaderXForwardedFor, // trust X-Forwarded-For from reverse proxy (Render/Heroku)
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			msg := "an internal error occurred"
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
				msg = e.Message
				if code == fiber.StatusNotFound && strings.HasPrefix(msg, "Cannot ") {
					msg = "endpoint not found"
				} else if code == fiber.StatusRequestEntityTooLarge {
					msg = "image is too large — maximum size is 20 MB"
				}
			} else {
				// Unknown error — log it but never expose to client
				log.Printf("[ERROR] unhandled (IP: %s): %s %s: %v", c.IP(), c.Method(), c.Path(), err) //nolint:gosec // G117: error/value from internal source, not user-controlled input
			}
			return c.Status(code).JSON(fiber.Map{"success": false, "error": msg})
		},
	})

	app.Use(middleware.CorrelationID)
	app.Use(middleware.RequestLogger)
	// Security headers — applied to every response
	app.Use(func(c *fiber.Ctx) error {
		c.Set("X-Content-Type-Options", "nosniff")
		c.Set("X-Frame-Options", "DENY")
		c.Set("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Set("Content-Security-Policy", "default-src 'none'; frame-ancestors 'none'")
		if !cfg.DevMode {
			c.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}
		return c.Next()
	})
	app.Use(cors.New(cors.Config{
		// AllowOriginsFunc allows the primary domain plus all Vercel preview deployments.
		// Preview URLs follow the pattern: <project>-<hash>-<team>.vercel.app
		// This lets you test any Vercel preview build without updating ALLOW_ORIGIN.
		AllowOriginsFunc: func(origin string) bool {
			return cfg.IsAllowedOrigin(origin)
		},
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
		AllowCredentials: true,
	}))

	// Serve uploaded images as static files
	// UPLOAD_DIR env var overrides default ./uploads (use persistent disk on Render)
	uploadDir := handlers.UploadDir
	app.Static("/uploads", uploadDir)

	// Root endpoint
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"success": true, "message": "HUNTIQ API is running"})
	})

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	api := app.Group("/api")
	api.Use(mw.OptionalAuth)
	api.Use(mw.RequireNoMaintenance)

	// Public rate limiter — keyed by IP for unauthenticated endpoints
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

	// ── Public ──────────────────────────────────────────────────────────────
	api.Get("/products", productHandler.List)
	api.Get("/products/:id", productHandler.Get)
	api.Get("/categories", productHandler.ListCategories)
	api.Get("/categories/popular", productHandler.PopularCategories)
	
	// Category Hierarchy (Public)
	api.Get("/categories/tree", categoryHandler.GetCategoryTree)
	api.Get("/categories/slug/:slug", categoryHandler.GetCategoryBySlug)
	api.Get("/categories/:id", categoryHandler.GetCategory)
	api.Get("/categories/:id/children", categoryHandler.GetCategoryChildren)
	api.Get("/categories/:id/ancestors", categoryHandler.GetCategoryBreadcrumbs)
	api.Get("/categories/:id/attributes", attributeHandler.GetCategoryAttributes)
	
	// Attributes (Public)
	api.Get("/attributes", attributeHandler.ListAttributeDefinitions)
	api.Get("/products/:id/attributes", attributeHandler.GetProductAttributes)
	
	api.Post("/orders", publicLimiter, mw.OptionalAuth, sellerHandler.CreateOrder)
	api.Get("/fee-options", productHandler.FeeOptions)
	api.Get("/locations/regions", locationHandler.ListRegions)
	api.Get("/locations/zones", locationHandler.ListZones)
	api.Get("/locations/subcities", locationHandler.ListSubcities)
	api.Get("/locations/woredas", locationHandler.ListWoredas)
	api.Get("/locations/search", locationHandler.SearchLocations)

	// Public banners (buyer/seller hero)
	api.Get("/banners", bannerHandler.ListPublic)

	// Public payment settings (merchant phones — safe to expose, no auth)
	api.Get("/public/payment-settings", adminHandler.GetPublicPaymentSettings)

	// Trust & Safety / Appeals (Public / User submission)
	api.Post("/reports", buyerHandler.CreateAbuseReport)
	api.Post("/appeals", buyerHandler.CreateAccountAppeal)

	// Public promotion banner (admin-controlled announcement box)
	api.Get("/public/promotion", adminHandler.GetPublicPromotion)

	// ── Auth ─────────────────────────────────────────────────────────────────
	// Rate limiters — keyed by IP to prevent brute-force and mass registration
	authLimiter := limiter.New(limiter.Config{
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
	api.Post("/auth/otp/send", otpLimiter, authHandler.SendOTP)
	api.Post("/auth/otp/verify", otpLimiter, authHandler.VerifyOTP)
	api.Post("/auth/otp/login", authLimiter, authHandler.OTPLogin)
	api.Post("/auth/demand", mw.RequireAuth, authHandler.SetDemand)
	api.Post("/seller/register", authLimiter, authHandler.Register)
	api.Post("/seller/login", authLimiter, authHandler.Login)
	api.Post("/seller/check", authLimiter, authHandler.CheckSellerExists)
	api.Post("/seller/forgot-password", authLimiter, authHandler.SellerForgotPassword)
	api.Post("/seller/reset-password", authLimiter, authHandler.SellerResetPassword)
	api.Post("/buyer/register", authLimiter, authHandler.BuyerRegister)
	api.Post("/buyer/login", authLimiter, authHandler.BuyerLogin)
	api.Post("/buyer/check", authLimiter, authHandler.CheckBuyerExists)
	api.Post("/buyer/forgot-password", authLimiter, authHandler.BuyerForgotPassword)
	api.Post("/buyer/reset-password", authLimiter, authHandler.BuyerResetPassword)
	api.Post("/admin/login", authLimiter, authHandler.AdminLogin)
	api.Post("/admin/login/verify", otpLimiter, authHandler.AdminLoginVerifyOTP)
	api.Post("/admin/forgot-password", authLimiter, authHandler.AdminForgotPassword)
	api.Post("/admin/reset-password", authLimiter, authHandler.AdminResetPassword)

	// ── Seller (JWT + role=seller) ───────────────────────────────────────────
	seller := api.Group("/seller", mw.RequireAuth, middleware.RequireRole("seller"))
	seller.Get("/stats", sellerHandler.GetStats)
	seller.Get("/fee-preview", productHandler.FeePreview)
	seller.Get("/products", productHandler.SellerList)
	seller.Get("/products/detailed", sellerHandler.ListProductsDetailed)
	seller.Post("/products", productHandler.SellerCreate)
	seller.Put("/products/:id", productHandler.SellerUpdate)
	seller.Delete("/products/:id", sellerHandler.DeleteProduct)
	seller.Get("/transactions", mw.RequireSellerTier("gold"), sellerHandler.ListTransactions)
	seller.Get("/orders", mw.RequireSellerTier("gold"), sellerHandler.ListOrders)
	seller.Put("/orders/:id", mw.RequireSellerTier("gold"), sellerHandler.UpdateOrder)
	seller.Get("/profile", sellerHandler.GetProfile)
	seller.Put("/profile", sellerHandler.UpdateProfile)
	seller.Put("/preference", sellerHandler.UpdatePreference)
	seller.Put("/password", authHandler.SellerChangePassword)
	seller.Post("/listing-fee", sellerHandler.PayListingFee)
	seller.Post("/upgrade-tier", sellerHandler.RequestUpgrade)
	seller.Get("/upgrade-tier", sellerHandler.GetUpgradeRequests)
	seller.Get("/notifications", sellerHandler.ListNotifications)
	seller.Put("/notifications/:id/read", sellerHandler.MarkNotificationRead)

	// Referral — seller
	// Rate limiter for referral code validation — prevents brute-force attacks on 34^8 keyspace
	referralLimiter := limiter.New(limiter.Config{
		Max:        3,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			if id, ok := c.Locals("userID").(int); ok {
				return fmt.Sprintf("referral:%d", id) //nolint:gosec // G117/G201: internal values only, no user-controlled input reaches this format string
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
	seller.Get("/referral", referralHandler.GetSellerStats)
	seller.Get("/referral/commissions", referralHandler.GetCommissions)
	seller.Post("/referral/accept", referralLimiter, referralHandler.AcceptInvitation)
	seller.Get("/referral/balance", referralHandler.GetBalance)
	seller.Post("/referral/withdraw", referralHandler.RequestWithdrawal)
	seller.Get("/referral/withdrawals", referralHandler.GetMyWithdrawals)
	
	// Product attributes (Seller)
	seller.Post("/products/:id/attributes", attributeHandler.SetProductAttributes)

	// Verification documents (Seller)
	seller.Get("/verification/documents", sellerHandler.ListVerificationDocuments)
	seller.Post("/verification/documents", sellerHandler.UploadVerificationDocument)
	seller.Post("/refunds", sellerHandler.CreateRefundRequest)

	// Image upload — sellers only, rate-limited to 20 uploads/hour per seller
	uploadLimiter := limiter.New(limiter.Config{
		Max:        20,
		Expiration: 1 * time.Hour,
		KeyGenerator: func(c *fiber.Ctx) string {
			if id, ok := c.Locals("userID").(int); ok {
				return fmt.Sprintf("upload:%d", id) //nolint:gosec // G117/G201: internal values only, no user-controlled input reaches this format string
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
	api.Post("/upload/image", mw.RequireAuth, middleware.RequireRole("seller"), uploadLimiter, handlers.UploadImage)

	// ── Buyer (JWT + role=buyer) ─────────────────────────────────────────────
	buyer := api.Group("/buyer", mw.RequireAuth, middleware.RequireRole("buyer"))
	buyer.Post("/upgrade-tier", buyerHandler.UpgradeTier)
	buyer.Get("/upgrade-tier", buyerHandler.GetUpgradeRequests)
	buyer.Get("/profile", buyerHandler.GetProfile)
	buyer.Put("/profile", buyerHandler.UpdateProfile)
	buyer.Get("/orders", buyerHandler.GetOrders)
	buyer.Put("/password", buyerHandler.ChangePassword)
	buyer.Get("/notifications", buyerHandler.ListNotifications)
	buyer.Put("/notifications/:id/read", buyerHandler.MarkNotificationRead)

	// Referral — buyer
	buyer.Get("/referral", referralHandler.GetBuyerStats)
	buyer.Get("/referral/commissions", referralHandler.GetCommissions)
	buyer.Post("/referral/accept", referralLimiter, referralHandler.AcceptInvitation)
	buyer.Get("/referral/balance", referralHandler.GetBalance)
	buyer.Post("/referral/withdraw", referralHandler.RequestWithdrawal)
	buyer.Get("/referral/withdrawals", referralHandler.GetMyWithdrawals)
	buyer.Post("/refunds", buyerHandler.CreateRefundRequest)

	api.Post("/transactions", mw.RequireAuth, middleware.RequireRole("seller"), authLimiter, txHandler.Create)

	// ── Admin (JWT + role=admin) ─────────────────────────────────────────────
	admin := api.Group("/admin", mw.RequireAuth, middleware.RequireRole("admin"))

	// Profile — any admin, no permission gate
	admin.Get("/me", adminHandler.GetMe)
	admin.Get("/me/permissions", adminHandler.GetMyPermissions) // normal admin fetches own permissions
	admin.Put("/password", adminHandler.ChangePassword)

	// Products
	admin.Get("/products/pending", mw.RequirePermission("pending"), adminHandler.ListPending)
	admin.Get("/products", mw.RequirePermission("products"), adminHandler.ListAllProducts)
	admin.Put("/products/:id/review", mw.RequirePermission("pending"), adminHandler.ReviewProduct)
	admin.Delete("/products/:id", mw.RequirePermission("products"), adminHandler.DeleteProduct)

	// Sellers
	admin.Get("/sellers", mw.RequirePermission("sellers"), adminHandler.ListSellers)
	admin.Put("/sellers/:id/active", mw.RequirePermission("sellers"), adminHandler.SetSellerActive)

	// Locations
	admin.Get("/locations", mw.RequirePermission("locations"), adminHandler.ListLocations)
	admin.Post("/locations", mw.RequirePermission("locations"), adminHandler.AddLocation)
	admin.Put("/locations/:id", mw.RequirePermission("locations"), adminHandler.UpdateLocation)
	admin.Delete("/locations/:id", mw.RequirePermission("locations"), adminHandler.DeleteLocation)

	// Analytics
	admin.Get("/analytics/overview", mw.RequirePermission("analytics"), adminHandler.AnalyticsOverview)
	admin.Get("/analytics/heatmap", mw.RequirePermission("analytics"), adminHandler.AnalyticsHeatmap)
	admin.Get("/analytics/monthly", mw.RequirePermission("analytics"), adminHandler.AnalyticsMonthly)
	admin.Get("/analytics/gantt", mw.RequirePermission("analytics"), adminHandler.AnalyticsGantt)
	admin.Get("/analytics/fee-revenue", middleware.RequireAdminLevel("super"), adminHandler.AnalyticsFeeRevenue)

	// Transactions — permission-gated (sensitive, off by default for new admins)
	admin.Get("/transactions", mw.RequirePermission("transactions"), adminHandler.ListTransactions)
	admin.Put("/transactions/:id/status", mw.RequirePermission("transactions"), txHandler.UpdateStatus)

	// Profit — permission-gated (sensitive, off by default)
	admin.Get("/profit/summary", mw.RequirePermission("profit"), adminHandler.ProfitSummary)
	admin.Get("/profit/by-tier", mw.RequirePermission("profit"), adminHandler.ProfitByTier)
	admin.Get("/profit/monthly", mw.RequirePermission("profit"), adminHandler.ProfitMonthly)
	admin.Get("/profit/weekly", mw.RequirePermission("profit"), adminHandler.ProfitWeekly)
	admin.Get("/profit/annual", mw.RequirePermission("profit"), adminHandler.ProfitAnnual)

	// Listing fee payment verification — admin must verify before revenue is counted
	admin.Get("/fee-payments/pending", mw.RequirePermission("products"), adminHandler.ListPendingFeePayments)
	admin.Put("/fee-payments/:id/verify", mw.RequirePermission("products"), adminHandler.VerifyListingFeePayment)
	admin.Put("/fee-payments/:id/reject", mw.RequirePermission("products"), adminHandler.RejectListingFeePayment)

	// Settings
	admin.Get("/settings", mw.RequirePermission("settings"), adminHandler.GetSettings)
	admin.Put("/settings", mw.RequirePermission("settings"), adminHandler.UpdateSettings)
	admin.Put("/settings/with-reason", middleware.RequireAdminLevel("super"), adminHandler.UpdateSettingsWithReason)
	admin.Get("/settings/log", middleware.RequireAdminLevel("super"), adminHandler.SettingsChangeLog)

	// User Governance
	admin.Post("/users/:id/warn", mw.RequirePermission("sellers"), adminHandler.IssueWarning)
	admin.Put("/users/:id/ban", middleware.RequireAdminLevel("super"), adminHandler.BanUser)

	// System Health
	admin.Get("/system/health", middleware.RequireAdminLevel("super"), adminHandler.SystemHealth)
	admin.Post("/system/health/repair", middleware.RequireAdminLevel("super"), adminHandler.RepairSystemHealth)

	// Notifications
	admin.Get("/notifications", mw.RequirePermission("notifications"), adminHandler.ListNotifications)
	admin.Put("/notifications/:id/read", mw.RequirePermission("notifications"), adminHandler.MarkNotificationRead)
	admin.Post("/notifications", middleware.RequireAdminLevel("super"), adminHandler.CreateNotification)

	// Categories (Legacy)
	admin.Get("/categories", mw.RequirePermission("categories"), adminHandler.ListCategories)
	admin.Post("/categories", middleware.RequireAdminLevel("super"), adminHandler.CreateCategory)
	admin.Put("/categories/:id", middleware.RequireAdminLevel("super"), adminHandler.UpdateCategory)
	admin.Delete("/categories/:id", middleware.RequireAdminLevel("super"), adminHandler.DeleteCategory)
	
	// Category Hierarchy Management (Admin)
	admin.Get("/category-hierarchy", mw.RequirePermission("categories"), categoryHandler.ListCategories)
	admin.Get("/category-hierarchy/:id/stats", mw.RequirePermission("categories"), categoryHandler.GetCategoryStats)
	admin.Post("/category-hierarchy", middleware.RequireAdminLevel("super"), categoryHandler.CreateCategory)
	admin.Put("/category-hierarchy/:id", middleware.RequireAdminLevel("super"), categoryHandler.UpdateCategory)
	admin.Delete("/category-hierarchy/:id", middleware.RequireAdminLevel("super"), categoryHandler.DeleteCategory)
	
	// Attribute Management (Admin)
	admin.Get("/attributes", mw.RequirePermission("categories"), attributeHandler.ListAttributeDefinitions)
	admin.Post("/attributes", middleware.RequireAdminLevel("super"), attributeHandler.CreateAttributeDefinition)
	admin.Post("/categories/:id/attributes", middleware.RequireAdminLevel("super"), attributeHandler.AssignAttributeToCategory)
	admin.Delete("/categories/:id/attributes/:attribute_id", middleware.RequireAdminLevel("super"), attributeHandler.RemoveAttributeFromCategory)

	// Sub-admin management — super admin only
	admin.Get("/admins", middleware.RequireAdminLevel("super"), adminHandler.ListAdmins)
	admin.Post("/admins", middleware.RequireAdminLevel("super"), adminHandler.CreateAdmin)
	admin.Put("/admins/:id/role", middleware.RequireAdminLevel("super"), adminHandler.UpdateAdminRole)
	admin.Delete("/admins/:id", middleware.RequireAdminLevel("super"), adminHandler.DeleteAdmin)
	admin.Get("/admins/:id/permissions", middleware.RequireAdminLevel("super"), adminHandler.GetAdminPermissions)
	admin.Put("/admins/:id/permissions", middleware.RequireAdminLevel("super"), adminHandler.SetAdminPermissions)

	// Hero banners — permission-gated
	admin.Get("/banners", mw.RequirePermission("banners"), bannerHandler.ListAll)
	admin.Post("/banners", mw.RequirePermission("banners"), bannerHandler.Create)
	admin.Put("/banners/:id", mw.RequirePermission("banners"), bannerHandler.Update)
	admin.Delete("/banners/:id", mw.RequirePermission("banners"), bannerHandler.Delete)

	// Admin activity logs — permission-gated (sensitive, off by default)
	admin.Get("/logs", middleware.RequireAdminLevel("super"), adminHandler.ListAdminLogs)
	admin.Delete("/logs", middleware.RequireAdminLevel("super"), adminHandler.DeleteAdminLogs)

	// OTP logs — permission-gated; bulk delete is super admin only (destructive)
	admin.Get("/otp-logs", mw.RequirePermission("otp_logs"), authHandler.ListOTPLogs)
	admin.Delete("/otp-logs", middleware.RequireAdminLevel("super"), authHandler.DeleteOTPLogs)

	// User activity logs — super admin only (full funnel: register/login/otp events)
	admin.Get("/activity-logs", middleware.RequireAdminLevel("super"), authHandler.ListActivityLogs)

	// Tier management — super admin only (structural change)
	admin.Put("/users/:id/seller-tier", middleware.RequireAdminLevel("super"), adminHandler.SetSellerTier)
	admin.Put("/users/:id/buyer-tier", middleware.RequireAdminLevel("super"), adminHandler.SetBuyerTier)
	admin.Get("/users/:id/tier-history", mw.RequirePermission("premium_users"), adminHandler.GetTierHistory)
	admin.Get("/premium-users", mw.RequirePermission("premium_users"), adminHandler.ListPremiumUsers)

	// Seller upgrade requests — seller_upgrades permission
	admin.Get("/upgrade-requests", mw.RequirePermission("seller_upgrades"), adminHandler.ListUpgradeRequests)
	admin.Put("/upgrade-requests/:id/approve", mw.RequirePermission("seller_upgrades"), adminHandler.ApproveUpgradeRequest)
	admin.Put("/upgrade-requests/:id/reject", mw.RequirePermission("seller_upgrades"), adminHandler.RejectUpgradeRequest)

	// Buyer upgrade requests — buyer_upgrades permission
	admin.Get("/buyer-upgrade-requests", mw.RequirePermission("buyer_upgrades"), adminHandler.ListBuyerUpgradeRequests)
	admin.Put("/buyer-upgrade-requests/:id/approve", mw.RequirePermission("buyer_upgrades"), adminHandler.ApproveBuyerUpgradeRequest)
	admin.Put("/buyer-upgrade-requests/:id/reject", mw.RequirePermission("buyer_upgrades"), adminHandler.RejectBuyerUpgradeRequest)

	// Category popularity (admin view with full stats)
	admin.Get("/categories/popularity", mw.RequirePermission("category_trends"), adminHandler.CategoryPopularity)

	// Referral overview (super admin only — shown on profit page)
	admin.Get("/referral/overview", middleware.RequireAdminLevel("super"), referralHandler.AdminOverview)
	admin.Get("/referral/list", middleware.RequireAdminLevel("super"), referralHandler.AdminList)
	admin.Put("/referral/users/:id/suspend", middleware.RequireAdminLevel("super"), referralHandler.AdminSetSuspended)
	admin.Put("/referral/users/:id/block", middleware.RequireAdminLevel("super"), referralHandler.AdminSetBlocked)
	// Withdrawal management — super admin OR admin with 'withdrawals' permission
	admin.Get("/referral/withdrawals", mw.RequirePermissionOrSuper("withdrawals"), referralHandler.AdminListWithdrawals)
	admin.Put("/referral/withdrawals/:id/approve", mw.RequirePermissionOrSuper("withdrawals"), referralHandler.AdminApproveWithdrawal)
	admin.Put("/referral/withdrawals/:id/reject", mw.RequirePermissionOrSuper("withdrawals"), referralHandler.AdminRejectWithdrawal)
	admin.Put("/referral/withdrawals/:id/paid", mw.RequirePermissionOrSuper("withdrawals"), referralHandler.AdminMarkWithdrawalPaid)

	// System Health & Repair routes are registered above under /system/health

	// ── Priority 1 Enhancements: New Admin Features ───────────────────────────────
	// Operations Monitoring (Role 02: Platform Operations Admin)
	admin.Get("/operations/health", mw.RequirePermission("operations"), adminHandler.GetOperationsHealth)
	admin.Get("/operations/jobs", mw.RequirePermission("operations"), adminHandler.GetJobStatusSummary)
	admin.Get("/operations/failed-jobs", mw.RequirePermission("operations"), adminHandler.GetFailedJobs)

	// Seller Verification (Role 08: Seller Verification Admin)
	admin.Get("/verification/documents", mw.RequirePermission("seller_verification"), adminHandler.ListVerificationDocuments)
	admin.Put("/verification/documents/:id/review", mw.RequirePermission("seller_verification"), adminHandler.ReviewVerificationDocument)

	// Trust & Safety (Role 21: Fraud & Trust & Safety Admin)
	admin.Get("/trust/reports", mw.RequirePermission("trust_safety"), adminHandler.ListAbuseReports)
	admin.Put("/trust/reports/:id/review", mw.RequirePermission("trust_safety"), adminHandler.ReviewAbuseReport)
	admin.Put("/trust/reports/:id/priority", mw.RequirePermission("trust_safety"), adminHandler.UpdateAbuseReportPriority)

	// Refund Management (Role 20: Refund & Dispute Admin)
	admin.Get("/refunds", mw.RequirePermission("refunds"), adminHandler.ListRefundRequests)
	admin.Put("/refunds/:id/review", mw.RequirePermission("refunds"), adminHandler.ReviewRefundRequest)
	admin.Put("/refunds/:id/processed", mw.RequirePermission("refunds"), adminHandler.MarkRefundProcessed)

	// Account Appeals (Role 30: Appeals & Account Recovery Admin)
	admin.Get("/appeals", mw.RequirePermission("appeals"), adminHandler.ListAccountAppeals)
	admin.Put("/appeals/:id/review", mw.RequirePermission("appeals"), adminHandler.ReviewAccountAppeal)
	admin.Put("/appeals/:id/status", mw.RequirePermission("appeals"), adminHandler.UpdateAppealStatus)

	// Search Analytics (admin — analytics permission)
	admin.Get("/search/popular", mw.RequirePermission("analytics"), searchHandler.PopularSearches)
	admin.Get("/search/zero-results", mw.RequirePermission("analytics"), searchHandler.ZeroResultQueries)

	// ── Dev helpers (only in development mode) ──────────────────────────────
	if cfg.DevMode {
		api.Get("/dev/otp", authHandler.DevPeekOTP)
		api.Get("/dev/otp/decrypt", authHandler.DevDecryptOTP)
		log.Println("⚠️  DEV MODE: /api/dev/otp and /api/dev/otp/decrypt endpoints enabled — disable in production (set APP_ENV=production)") //nolint:gosec // G117: error/value from internal source, not user-controlled input
	}

	log.Printf("Server starting on port %s", cfg.AppPort) //nolint:gosec // G117: error/value from internal source, not user-controlled input
	log.Fatal(app.Listen(":" + cfg.AppPort))
}
