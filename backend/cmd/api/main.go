package main

import (
	"context"
	"flag"
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
	_ = godotenv.Load("backend/.env")
	_ = godotenv.Load("../backend/.env")
	migrateOnly := flag.Bool("migrate-only", false, "Run database migrations and exit")
	seedOnly := flag.Bool("seed-only", false, "Run master data seed and exit")
	flag.Parse()

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

	if *migrateOnly {
		if err := database.RunMigrations(context.Background(), db, ""); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
		os.Exit(0)
	}
	if *seedOnly {
		if err := database.VerifyOrSeed(context.Background(), db, ""); err != nil {
			log.Fatalf("Seed failed: %v", err)
		}
		os.Exit(0)
	}

	// Always verify migrations and seeds on every run with auto-recovery
	startupCtx, startupCancel := context.WithTimeout(context.Background(), 2*time.Minute)
	if err := database.AutoMigrateAndSeed(startupCtx, db); err != nil {
		log.Printf("[HUNTIQ DB] Warning: auto-migrate/seed encountered error: %v (continuing)", err)
	}

	// Ensure critical schema columns exist across any connected database instance
	_, _ = db.Exec(startupCtx, `
		DO $$ BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'account_entity_type') THEN
				CREATE TYPE account_entity_type AS ENUM ('INDIVIDUAL', 'COMPANY', 'GOVERNMENT', 'OTHER_LEGALLY_APPROVED_ENTITY');
			END IF;
			IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'org_verification_status') THEN
				CREATE TYPE org_verification_status AS ENUM ('DRAFT', 'SUBMITTED', 'IDENTITY_PENDING', 'DOCUMENT_REVIEW', 'VERIFIED', 'REJECTED', 'SUSPENDED', 'EXPIRED', 'REVERIFICATION_REQUIRED');
			END IF;
			IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'identity_verification_status') THEN
				CREATE TYPE identity_verification_status AS ENUM ('pending', 'verified', 'rejected', 'escalated', 'expired');
			END IF;
		END $$;

		CREATE TABLE IF NOT EXISTS organizations (
			id SERIAL PRIMARY KEY,
			name VARCHAR(200) NOT NULL,
			entity_type account_entity_type NOT NULL DEFAULT 'INDIVIDUAL',
			registration_number VARCHAR(100),
			tin_number VARCHAR(50),
			official_address TEXT,
			official_email VARCHAR(150),
			official_phone VARCHAR(25),
			publishing_authority_evidence TEXT,
			verification_status org_verification_status NOT NULL DEFAULT 'DRAFT',
			reviewed_by INT,
			reviewed_at TIMESTAMPTZ,
			review_notes TEXT,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);

		ALTER TABLE products ADD COLUMN IF NOT EXISTS custom_url character varying(500);
		ALTER TABLE products ADD COLUMN IF NOT EXISTS social_links jsonb DEFAULT '{}'::jsonb;
		ALTER TABLE products ADD COLUMN IF NOT EXISTS publisher_role VARCHAR(50) NOT NULL DEFAULT 'seller';
		ALTER TABLE products ADD COLUMN IF NOT EXISTS organization_id INT;
		ALTER TABLE products ADD COLUMN IF NOT EXISTS brokerage_disclosure_agreed BOOLEAN NOT NULL DEFAULT true;

		ALTER TABLE users ADD COLUMN IF NOT EXISTS custom_url character varying(500);
		ALTER TABLE users ADD COLUMN IF NOT EXISTS account_entity_type account_entity_type NOT NULL DEFAULT 'INDIVIDUAL';
		ALTER TABLE users ADD COLUMN IF NOT EXISTS identity_verified BOOLEAN NOT NULL DEFAULT false;
		ALTER TABLE users ADD COLUMN IF NOT EXISTS organization_id INT;

		-- Clean any legacy mixed-language category names where Amharic in parentheses leaked into English
		UPDATE category_hierarchy 
		SET name_en = TRIM(REGEXP_REPLACE(name_en, '\s*\([^)]*[\u1200-\u137F][^)]*\)', ''))
		WHERE name_en ~ '[\u1200-\u137F]';

		UPDATE categories 
		SET name = TRIM(REGEXP_REPLACE(name, '\s*\([^)]*[\u1200-\u137F][^)]*\)', ''))
		WHERE name ~ '[\u1200-\u137F]';
	`)
	startupCancel()

	searchConfigCandidates := []string{
		"internal/database/loccatdata/search_config.json",
		"backend/internal/database/loccatdata/search_config.json",
		"../internal/database/loccatdata/search_config.json",
		"../../internal/database/loccatdata/search_config.json",
	}
	searchLoaded := false
	for _, sc := range searchConfigCandidates {
		if err := search.LoadSearchConfig(sc); err == nil {
			log.Printf("[HUNTIQ] Loaded search config from %s", sc)
			searchLoaded = true
			break
		}
	}
	if !searchLoaded {
		log.Printf("[INFO] Search config: default settings loaded")
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
	smsRepo := repository.NewSMSRepo(db)
	ledgerRepo := repository.NewLedgerRepo(db)
	referralTokenRepo := repository.NewReferralTokenRepo(db)
	telebirrRepo := repository.NewTelebirrRepo(db)
	identityRepo := repository.NewIdentityRepo(db)
	orgRepo := repository.NewOrganizationRepo(db)
	recoveryRepo := repository.NewRecoveryRepo(db)
	masterDataRepo := repository.NewMasterDataRepo(db)

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
	ledgerSvc := service.NewLedgerService(ledgerRepo)
	smsMgr := service.NewSMSManager(smsRepo, cfg.DevMode)
	referralTokenSvc := service.NewReferralTokenService(referralTokenRepo, ledgerSvc)
	telebirrCfg := &service.TelebirrConfig{
		AppID:       cfg.TelebirrAppID,
		AppKey:      cfg.TelebirrAppKey,
		ShortCode:   cfg.TelebirrShortCode,
		APIEndpoint: cfg.TelebirrAPIEndpoint,
		NotifyURL:   cfg.TelebirrNotifyURL,
		ReturnURL:   cfg.TelebirrReturnURL,
	}
	telebirrSvc := service.NewTelebirrService(telebirrRepo, referralTokenSvc, ledgerSvc, telebirrCfg, cfg.DevMode)

	referralSvc := service.NewReferralService(referralRepo).
		WithSettingsReader(func(ctx context.Context, k string) string {
			return adminSvc.GetSetting(ctx, k)
		}).
		WithLedgerService(ledgerSvc)
	fpaSvc := service.NewFPAService(fpaRepo)
	identitySvc := service.NewIdentityService(identityRepo, cfg.DevMode)
	orgSvc := service.NewOrganizationService(orgRepo)
	masterDataSvc := service.NewMasterDataService(masterDataRepo)

	txSvc.WithReferralService(referralSvc)
	authSvc.WithReferralService(referralSvc).
		WithAdminService(adminSvc).
		WithSMSManager(smsMgr).
		WithReferralTokenService(referralTokenSvc).
		WithRecoveryRepo(recoveryRepo).
		WithIdentityService(identitySvc).
		WithOrganizationService(orgSvc)
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
	telebirrHandler := handlers.NewTelebirrHandler(telebirrSvc)
	walletHandler := handlers.NewWalletHandler(ledgerSvc, referralTokenSvc, cfg.FrontendURL)
	adminFinancialHandler := handlers.NewAdminFinancialHandler(ledgerSvc, smsRepo)
	identityHandler := handlers.NewIdentityHandler(identitySvc)
	organizationHandler := handlers.NewOrganizationHandler(orgSvc)
	masterDataHandler := handlers.NewMasterDataHandler(masterDataSvc)

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
		TelebirrHandler:       telebirrHandler,
		WalletHandler:         walletHandler,
		AdminFinancialHandler: adminFinancialHandler,
		FPAHandler:            fpaHandler,
		SearchHandler:         searchHandler,
		IdentityHandler:       identityHandler,
		OrganizationHandler:   organizationHandler,
		MasterDataHandler:     masterDataHandler,
	})

	// Background listener with resilience for port binding
	serverErr := make(chan error, 1)
	go func() {
		log.Printf("Server starting on port %s", cfg.AppPort)
		var listenErr error
		for attempt := 1; attempt <= 5; attempt++ {
			listenErr = app.Listen(":" + cfg.AppPort)
			if listenErr == nil || !strings.Contains(listenErr.Error(), "address already in use") {
				break
			}
			log.Printf("Port %s busy (attempt %d/5), waiting 1s...", cfg.AppPort, attempt)
			time.Sleep(1 * time.Second)
		}
		if listenErr != nil {
			serverErr <- listenErr
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
