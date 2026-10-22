package config

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
	DBSSLMode     string
	DBMaxConns    int
	DBMinConns    int
	AppPort       string
	AppBaseURL    string
	FrontendURL   string // public frontend URL used for referral links
	JWTSecret     string
	AllowOrigin   string
	OTPEncryptKey string
	DevMode       bool
	UploadDir     string // directory for uploaded files
	MaxUploadMB   int    // max file upload size in MB
}

func Load() *Config {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("[FATAL] JWT_SECRET environment variable is required")
	}
	if len(secret) < 32 {
		log.Fatal("[FATAL] JWT_SECRET must be at least 32 characters long")
	}
	otpKey := os.Getenv("OTP_ENCRYPT_KEY")
	if len(otpKey) < 32 {
		log.Fatal("[FATAL] OTP_ENCRYPT_KEY environment variable is required and must be at least 32 characters long")
	}
	if _, err := base64.StdEncoding.DecodeString(otpKey); err != nil {
		log.Fatal("[FATAL] OTP_ENCRYPT_KEY must be a valid base64-encoded string")
	}

	allowOrigin := os.Getenv("ALLOW_ORIGIN")
	if allowOrigin == "" {
		log.Fatal("[FATAL] ALLOW_ORIGIN environment variable is required")
	}

	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		log.Fatal("[FATAL] FRONTEND_URL environment variable is required")
	}

	devMode := os.Getenv("APP_ENV") == "development"
	if devMode {
		log.Println("[WARN] APP_ENV=development — /api/dev/* endpoints are enabled. Do NOT use in production.")
	}

	appBaseURL := getEnv("APP_BASE_URL", "")
	if appBaseURL == "" {
		log.Fatal("[FATAL] APP_BASE_URL environment variable is required")
	}

	dbUser := getEnv("DB_USER", "")
	if dbUser == "" && os.Getenv("DATABASE_URL") == "" {
		log.Fatal("[FATAL] DB_USER environment variable is required (or set DATABASE_URL)")
	}
	dbName := getEnv("DB_NAME", "")
	if dbName == "" && os.Getenv("DATABASE_URL") == "" {
		log.Fatal("[FATAL] DB_NAME environment variable is required (or set DATABASE_URL)")
	}

	return &Config{
		DBHost:        getEnv("DB_HOST", "localhost"),
		DBPort:        getEnv("DB_PORT", "5432"),
		DBUser:        getEnv("DB_USER", ""),
		DBPassword:    getEnv("DB_PASSWORD", ""),
		DBName:        getEnv("DB_NAME", ""),
		DBSSLMode:     getEnv("DB_SSLMODE", "disable"),
		DBMaxConns:    getEnvInt("DB_MAX_CONNS", 25),
		DBMinConns:    getEnvInt("DB_MIN_CONNS", 5),
		AppPort:       getEnv("APP_PORT", "8080"),
		AppBaseURL:    appBaseURL,
		FrontendURL:   frontendURL,
		JWTSecret:     secret,
		AllowOrigin:   allowOrigin,
		OTPEncryptKey: otpKey,
		DevMode:       devMode,
		UploadDir:     getEnv("UPLOAD_DIR", "./uploads"),
		MaxUploadMB:   getEnvInt("MAX_UPLOAD_MB", 20),
	}
}

// IsAllowedOrigin returns true if the request origin is permitted.
// Allows:
//   - The exact ALLOW_ORIGIN value (e.g. https://hulubebet.vercel.app)
//   - Vercel preview URLs matching <project>-<hash>-<team>.vercel.app (only when ALLOW_ORIGIN is a vercel.app URL)
//   - localhost on any port for local development only (DevMode)
func (c *Config) IsAllowedOrigin(origin string) bool {
	// Exact match — production domain
	if origin == c.AllowOrigin {
		return true
	}
	// Vercel preview deployments: hulubebet-<hash>-<team>.vercel.app
	// Pattern: origin must start with "https://<project>-" and end with ".vercel.app"
	// The team slug is required so "hulubebet-evil.vercel.app" does NOT match.
	if strings.Contains(origin, ".vercel.app") && strings.Contains(c.AllowOrigin, ".vercel.app") {
		project := strings.TrimPrefix(c.AllowOrigin, "https://")
		project = strings.TrimSuffix(project, ".vercel.app")
		originHost := strings.TrimPrefix(origin, "https://")
		// Must match: <project>-<something>-<something>.vercel.app (at least 2 dashes after project)
		prefix := project + "-"
		if strings.HasPrefix(originHost, prefix) {
			rest := strings.TrimPrefix(originHost, prefix)
			rest = strings.TrimSuffix(rest, ".vercel.app")
			// Require at least one more dash (hash-team pattern), blocking bare "hulubebet-evil"
			if strings.Contains(rest, "-") {
				return true
			}
		}
	}
	// Local development only — never allow localhost in production
	if c.DevMode {
		if strings.HasPrefix(origin, "http://localhost:") || strings.HasPrefix(origin, "http://127.0.0.1:") {
			return true
		}
	}
	return false
}

func (c *Config) DSN() string {
	// If DATABASE_URL is set (e.g. from Nhost/Render/Heroku), use it directly.
	if url := os.Getenv("DATABASE_URL"); url != "" {
		return url
	}
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode,
	)
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			return n
		}
	}
	return fallback
}
