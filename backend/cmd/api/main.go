package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"

	"github.com/your-org/used-tools-marketplace/internal/cache"
	"github.com/your-org/used-tools-marketplace/internal/config"
	"github.com/your-org/used-tools-marketplace/internal/database"
	handlers "github.com/your-org/used-tools-marketplace/internal/handler"
	"github.com/your-org/used-tools-marketplace/internal/middleware"
	"github.com/your-org/used-tools-marketplace/internal/models"
	"github.com/your-org/used-tools-marketplace/internal/repository"
	"github.com/your-org/used-tools-marketplace/internal/security"
	"github.com/your-org/used-tools-marketplace/internal/service"
	jobs "github.com/your-org/used-tools-marketplace/internal/worker"
	"github.com/your-org/used-tools-marketplace/pkg/search"
)

func main() {
	_ = godotenv.Load()
	cfg := config.Load()
	handlers.AppBaseURL = cfg.AppBaseURL

	otpCrypto, err := security.NewOTPCrypto(cfg.OTPEncryptKey)
	if err != nil {
		log.Fatalf("OTP crypto initialization failed: %v", err)
	}

	db, err := database.NewPool(cfg)
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	defer db.Close()

	if err := search.LoadSearchConfig("data/search_config.json"); err != nil {
		log.Printf("[INFO] Search config: default settings loaded (%v)", err)
	}

	// Background worker context with graceful cancellation
	jobCtx, cancelJobs := context.WithCancel(context.Background())
	defer cancelJobs()
	jobs.StartExpiryJob(jobCtx, db)
	jobs.StartUnverifiedUserCleanup(jobCtx, db)
	jobs.StartOTPCleanup(jobCtx, db)
	jobs.StartCommissionConfirmJob(jobCtx, db)

	// Repositories
	productRepo := repository.NewProductRepo(db)
	categoryRepo := repository.NewCategoryRepo(db)
	categoryHierarchyRepo := repository.NewCategoryHierarchyRepo(db)
	attributeRepo := repository.NewAttributeRepo(db)
	searchRepo := repository.NewSearchRepo(db)
	locationRepo := repository.NewLocationRepo(db)
	locationHierarchyRepo := repository.NewLocationHierarchyRepo(db)
	userRepo := repository.NewUserRepo(db).WithOTPCrypto(otpCrypto)
	txRepo := repository.NewTransactionRepo(db)
	adminRepo := repository.NewAdminRepo(db)
	sellerRepo := repository.NewSellerRepo(db)
	bannerRepo := repository.NewBannerRepo(db)
	referralRepo := repository.NewReferralRepo(db)
	sessionRepo := repository.NewSessionRepo(db)
	classificationRepo := repository.NewClassificationRepo(db)
	adminPasskeyRepo := repository.NewAdminPasskeyRepo(db)
	fpaRepo := repository.NewFPARepo(db)

	// In-memory active user cache
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

	// Services
	productSvc := service.NewProductService(productRepo, categoryRepo, locationRepo, userRepo).
		WithBackendBase(cfg.AppBaseURL).
		WithAttributeRepo(attributeRepo).
		WithCategoryHierarchyRepo(categoryHierarchyRepo)
	locationSvc := service.NewLocationService(locationRepo).WithHierarchy(locationHierarchyRepo)
	categoryHierarchySvc := service.NewCategoryHierarchyService(categoryHierarchyRepo)
	classificationSvc := service.NewClassificationService(classificationRepo)
	attributeSvc := service.NewAttributeService(attributeRepo)
	authSvc := service.NewAuthService(userRepo, cfg.JWTSecret).WithDevMode(cfg.DevMode).WithOTPCrypto(otpCrypto)
	sessionSvc := service.NewSessionService(sessionRepo, userRepo, cfg.JWTSecret, cfg.DevMode)
	authSvc.WithSessionRevoker(sessionSvc.RevokeAllUserSessions)
	txSvc := service.NewTransactionService(txRepo, productRepo, userRepo)
	adminSvc := service.NewAdminService(adminRepo)
	sellerSvc := service.NewSellerService(sellerRepo)
	referralSvc := service.NewReferralService(referralRepo).WithSettingsReader(func(ctx context.Context, k string) string {
		return adminSvc.GetSetting(ctx, k)
	})
	fpaSvc := service.NewFPAService(fpaRepo)

	txSvc.WithReferralService(referralSvc)
	authSvc.WithReferralService(referralSvc).WithAdminService(adminSvc)
	adminSvc.WithReferralService(referralSvc).WithSellerService(sellerSvc)

	// WebAuthn Manager
	rpID := security.DeriveRPID(cfg.FrontendURL)
	waManager, err := security.NewWebAuthnManager(rpID, []string{cfg.FrontendURL, cfg.AppBaseURL, "http://localhost:3000", "http://localhost:8080"})
	if err != nil {
		log.Fatalf("WebAuthn initialization failed: %v", err)
	}
	adminAuthSvc := service.NewAdminAuthService(adminPasskeyRepo, userRepo, adminRepo, waManager, cfg.JWTSecret)
	adminSvc.WithSessionRevoker(func(ctx context.Context, userID int, reason string) error {
		isActiveMu.Lock()
		delete(isActiveCache, userID)
		isActiveMu.Unlock()
		return adminAuthSvc.RevokeAllSessions(ctx, userID, reason)
	})
	authSvc.WithSessionCreator(func(ctx context.Context, userID int, portal, ip, ua string) (string, int, error) {
		rawToken, sess, err := adminAuthSvc.CreateAdminSession(ctx, userID, portal, ip, ua)
		if err != nil {
			return "", 0, err
		}
		return rawToken, sess.ID, nil
	})

	// Middleware
	mw := middleware.New(cfg.JWTSecret).
		WithSessionValidator(adminAuthSvc.ValidateSession).
		WithAdminSessionValidator(adminAuthSvc.ValidateSessionByID).
		WithReauthValidator(adminAuthSvc.VerifyReauth).
		WithIsActiveChecker(func(userID int) (bool, error) {
			isActiveMu.Lock()
			if e, ok := isActiveCache[userID]; ok && time.Now().Before(e.expires) {
				isActiveMu.Unlock()
				return e.active, nil
			}
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
			if u.Role != "admin" && !u.OTPVerified && authSvc.IsOTPRequired(ctx) {
				active = false
			}

			isActiveMu.Lock()
			isActiveCache[userID] = struct {
				active  bool
				expires time.Time
			}{active, time.Now().Add(10 * time.Second)}
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
		}).
		WithPermissionChecker(func(adminID int, permission string) (bool, error) {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			return adminSvc.HasPermission(ctx, adminID, permission)
		}).
		WithMaintenanceChecker(func() (bool, error) {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			return adminSvc.GetSetting(ctx, "maintenance_mode") == "true", nil
		})

	appCache := cache.New(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB)
	defer appCache.Close()

	// Handlers
	productHandler := handlers.NewProductHandler(productSvc, authSvc).WithSearchLogger(searchRepo).WithCategoryHierarchy(categoryHierarchySvc).WithCache(appCache)
	locationHandler := handlers.NewLocationHandler(locationSvc)
	categoryHandler := handlers.NewCategoryHandler(categoryHierarchySvc).WithCache(appCache).WithAdminLogger(func(ctx context.Context, adminID int, action, notes string) {
		adminSvc.InsertAdminLog(ctx, models.AdminLog{AdminID: adminID, Action: action, Notes: notes})
	})
	attributeHandler := handlers.NewAttributeHandler(attributeSvc).WithProductService(productSvc)
	authHandler := handlers.NewAuthHandler(authSvc).WithSessionService(sessionSvc).WithAdminService(adminSvc).WithDevMode(cfg.DevMode)
	adminAuthHandler := handlers.NewAdminAuthHandler(adminAuthSvc, authSvc, adminSvc, !cfg.DevMode)
	txHandler := handlers.NewTransactionHandler(txSvc)
	adminHandler := handlers.NewAdminHandler(adminSvc, authSvc).WithCategoryRepo(categoryRepo).WithProductRepo(productRepo)
	sellerHandler := handlers.NewSellerHandler(sellerSvc).WithAdminService(adminSvc)
	bannerHandler := handlers.NewBannerHandler(bannerRepo)
	buyerHandler := handlers.NewBuyerHandler(authSvc).WithAdminService(adminSvc)
	fpaHandler := handlers.NewFPAHandler(fpaSvc)
	classificationHandler := handlers.NewClassificationHandler(classificationSvc)
	referralHandler := handlers.NewReferralHandler(referralSvc, cfg.FrontendURL).
		WithAdminNotifier(func(ctx context.Context, userID int, amount float64) {
			notifyCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			msg := fmt.Sprintf("User #%d requested withdrawal of ETB %.2f", userID, amount)
			_ = adminSvc.CreateNotification(notifyCtx, "Withdrawal Request", msg, "warning", "")
		}).
		WithAdminLogger(func(ctx context.Context, adminID int, targetUserID *int, action, notes string) {
			adminSvc.InsertAdminLog(ctx, models.AdminLog{AdminID: adminID, TargetUserID: targetUserID, Action: action, Notes: notes})
		})
	searchHandler := handlers.NewSearchHandler(searchRepo)

	// Fiber HTTP Server
	app := fiber.New(fiber.Config{
		BodyLimit:   22 * 1024 * 1024,
		ProxyHeader: fiber.HeaderXForwardedFor,
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
				log.Printf("[ERROR] unhandled (IP: %s): %s %s: %v", c.IP(), c.Method(), c.Path(), err)
			}
			return c.Status(code).JSON(fiber.Map{"success": false, "error": msg})
		},
	})

	app.Use(middleware.CorrelationID)
	app.Use(middleware.RequestLogger)
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
		AllowOriginsFunc: cfg.IsAllowedOrigin,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, X-CSRF-Token",
		AllowMethods:     "GET, POST, PUT, PATCH, DELETE, OPTIONS",
		AllowCredentials: true,
	}))

	app.Static("/uploads", handlers.UploadDir)
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"success": true, "message": "HUNTIQ API is running"})
	})
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// Register all application routes
	registerRoutes(app, &AppDeps{
		Config:                cfg,
		MW:                    mw,
		ProductHandler:        productHandler,
		LocationHandler:       locationHandler,
		CategoryHandler:       categoryHandler,
		AttributeHandler:      attributeHandler,
		ClassificationHandler: classificationHandler,
		AuthHandler:           authHandler,
		AdminAuthHandler:      adminAuthHandler,
		AdminHandler:          adminHandler,
		SellerHandler:         sellerHandler,
		BuyerHandler:          buyerHandler,
		TransactionHandler:    txHandler,
		BannerHandler:         bannerHandler,
		ReferralHandler:       referralHandler,
		FPAHandler:            fpaHandler,
		SearchHandler:         searchHandler,
	})

	// Background listener
	serverErr := make(chan error, 1)
	go func() {
		log.Printf("Server starting on port %s", cfg.AppPort)
		if err := app.Listen(":" + cfg.AppPort); err != nil {
			serverErr <- err
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErr:
		log.Fatalf("Server startup failed: %v", err)
	case sig := <-quit:
		log.Printf("Received signal %v — initiating graceful shutdown...", sig)
	}

	cancelJobs()
	if err := app.ShutdownWithTimeout(10 * time.Second); err != nil {
		log.Printf("Server shutdown warning: %v", err)
	}
	log.Println("Graceful shutdown completed cleanly.")
}
