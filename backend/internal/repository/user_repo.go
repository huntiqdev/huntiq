package repository

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/your-org/used-tools-marketplace/internal/models"
	"github.com/your-org/used-tools-marketplace/internal/security"
	"github.com/your-org/used-tools-marketplace/pkg/constants"
)


// sanitizeLog removes newlines and carriage returns from user input before logging
// to prevent log injection attacks (CWE-117).
func sanitizeLog(s string) string {
	return strings.NewReplacer("\n", "", "\r", "", "\t", " ").Replace(s)
}

type UserRepo struct {
	db          *pgxpool.Pool
	asyncLogger *AsyncLogger
	otpCrypto   *security.OTPCrypto
}

func NewUserRepo(db *pgxpool.Pool) *UserRepo {
	return &UserRepo{
		db:          db,
		asyncLogger: NewAsyncLogger(db),
	}
}

// WithOTPCrypto injects the OTP encryption/decryption handler
func (r *UserRepo) WithOTPCrypto(otpCrypto *security.OTPCrypto) *UserRepo {
	r.otpCrypto = otpCrypto
	return r
}

// generateReferralCode produces a unique "REF" + 8 char code using the same
// charset as migration 021 (no I, O, 0, 1 to avoid visual confusion).
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

func scanUser(row interface {
	Scan(...interface{}) error
}) (*models.User, error) {
	u := &models.User{}
	err := row.Scan(
		&u.ID, &u.Name, &u.Phone, &u.BusinessName, &u.Email, &u.PasswordHash,
		&u.Rating, &u.Role, &u.AdminRole, &u.IsActive, &u.DeactivationReason,
		&u.Tier, &u.TotalValueSold,
		&u.BuyerTier, &u.OTPVerified, &u.SellerPreference,
		&u.CreatedAt, &u.UpdatedAt,
	)
	return u, err
}

func (r *UserRepo) Create(ctx context.Context, u *models.User) (int, error) {
	var id int
	role := u.Role
	if role == "" {
		role = "seller"
	}
	sellerPref := u.SellerPreference
	if sellerPref == "" {
		sellerPref = "sale" // default
	}
	refCode, err := generateReferralCode()
	if err != nil {
		return 0, err
	}
	//nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
	err = r.db.QueryRow(ctx, `
		INSERT INTO users (name, phone, business_name, email, password_hash, role, is_active, referral_code, seller_preference)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id`,
		u.Name, u.Phone, u.BusinessName, u.Email, u.PasswordHash, role, u.IsActive, refCode, sellerPref,
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	// Always activate 'buyer' capability for all registered users
	_, _ = r.db.Exec(ctx, `
		INSERT INTO user_capabilities (user_id, capability, status) 
		VALUES ($1, 'buyer', 'active') 
		ON CONFLICT (user_id, capability) DO NOTHING`, id)

	// If registered as seller or has seller/landlord preference, activate respective capabilities
	if role == "seller" || u.HasCapability("seller") {
		_, _ = r.db.Exec(ctx, `
			INSERT INTO user_capabilities (user_id, capability, status) 
			VALUES ($1, 'seller', 'active') 
			ON CONFLICT (user_id, capability) DO NOTHING`, id)
	}
	if sellerPref == "rent" || sellerPref == "both" || u.HasCapability("landlord") {
		_, _ = r.db.Exec(ctx, `
			INSERT INTO user_capabilities (user_id, capability, status) 
			VALUES ($1, 'landlord', 'active') 
			ON CONFLICT (user_id, capability) DO NOTHING`, id)
	}

	return id, nil
}

// CreateOrGet inserts a new user or returns the existing user's ID on (phone, role) conflict.
// This is safe for concurrent public-sell requests for the same phone number.
func (r *UserRepo) CreateOrGet(ctx context.Context, u *models.User) (int, error) {
	role := u.Role
	if role == "" {
		role = "seller"
	}
	sellerPref := u.SellerPreference
	if sellerPref == "" {
		sellerPref = "sale"
	}
	refCode, err := generateReferralCode()
	if err != nil {
		return 0, err
	}
	var id int
	//nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
	err = r.db.QueryRow(ctx, `
		INSERT INTO users (name, phone, business_name, email, password_hash, role, is_active, referral_code, seller_preference)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (phone, role) DO UPDATE SET updated_at = NOW()
		RETURNING id`,
		u.Name, u.Phone, u.BusinessName, u.Email, u.PasswordHash, role, u.IsActive, refCode, sellerPref,
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	_, _ = r.db.Exec(ctx, `INSERT INTO user_capabilities (user_id, capability, status) VALUES ($1, 'buyer', 'active') ON CONFLICT DO NOTHING`, id)
	if role == "seller" {
		_, _ = r.db.Exec(ctx, `INSERT INTO user_capabilities (user_id, capability, status) VALUES ($1, 'seller', 'active') ON CONFLICT DO NOTHING`, id)
	}
	return id, nil
}

// GetCapabilities returns the active capability names for a user
func (r *UserRepo) GetCapabilities(ctx context.Context, userID int) ([]string, error) {
	rows, err := r.db.Query(ctx, `
		SELECT capability FROM user_capabilities 
		WHERE user_id = $1 AND status = 'active' 
		ORDER BY capability ASC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var caps []string
	for rows.Next() {
		var c string
		if err := rows.Scan(&c); err == nil {
			caps = append(caps, c)
		}
	}
	if len(caps) == 0 {
		caps = []string{"buyer"}
	}
	return caps, nil
}

// ActivateCapability adds or restores an active capability for a user
func (r *UserRepo) ActivateCapability(ctx context.Context, userID int, capability string, metadata map[string]interface{}) error {
	metaJSON, _ := json.Marshal(metadata)
	_, err := r.db.Exec(ctx, `
		INSERT INTO user_capabilities (user_id, capability, status, activated_at, metadata)
		VALUES ($1, $2, 'active', NOW(), $3)
		ON CONFLICT (user_id, capability) DO UPDATE
		SET status = 'active', activated_at = NOW(), metadata = COALESCE(user_capabilities.metadata, '{}'::jsonb) || EXCLUDED.metadata`,
		userID, capability, metaJSON)
	if err != nil {
		return err
	}
	// If activating seller or landlord, update seller_preference and role where helpful for backward compatibility
	if capability == "seller" {
		_, _ = r.db.Exec(ctx, `UPDATE users SET role = 'seller' WHERE id = $1 AND role = 'buyer'`, userID)
	} else if capability == "landlord" {
		_, _ = r.db.Exec(ctx, `UPDATE users SET seller_preference = CASE WHEN seller_preference = 'sale' THEN 'both' ELSE 'rent' END WHERE id = $1`, userID)
	}
	if bn, ok := metadata["business_name"].(string); ok && bn != "" {
		_, _ = r.db.Exec(ctx, `UPDATE users SET business_name = $1 WHERE id = $2 AND (business_name IS NULL OR business_name = '')`, bn, userID)
	}
	return nil
}

// SuspendCapability marks a capability as suspended
func (r *UserRepo) SuspendCapability(ctx context.Context, userID int, capability string, reason string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE user_capabilities
		SET status = 'suspended',
		    metadata = jsonb_set(COALESCE(metadata, '{}'::jsonb), '{suspend_reason}', to_jsonb($3::text))
		WHERE user_id = $1 AND capability = $2`,
		userID, capability, reason)
	return err
}

// PopulateCapabilities populates the Capabilities slice on a models.User
func (r *UserRepo) PopulateCapabilities(ctx context.Context, u *models.User) {
	if u == nil {
		return
	}
	caps, err := r.GetCapabilities(ctx, u.ID)
	if err == nil && len(caps) > 0 {
		u.Capabilities = caps
	} else {
		u.Capabilities = []string{"buyer"}
		if u.Role == "seller" {
			u.Capabilities = append(u.Capabilities, "seller")
		}
	}
}

func (r *UserRepo) GetByPhone(ctx context.Context, phone string) (*models.User, error) {
	u, err := scanUser(r.db.QueryRow(ctx,
		"SELECT id, name, phone, COALESCE(business_name,''), COALESCE(email,''), password_hash, rating, role, COALESCE(admin_role::varchar,''), is_active, COALESCE(deactivation_reason,''), COALESCE(tier::varchar,'normal'), COALESCE(total_value_sold,0), COALESCE(buyer_tier::varchar,'normal'), COALESCE(otp_verified,false), COALESCE(seller_preference,'sale'), created_at, updated_at FROM users WHERE phone = $1 ORDER BY created_at DESC LIMIT 1", phone))
	if err == nil && u != nil {
		r.PopulateCapabilities(ctx, u)
	}
	return u, err
}

// GetByPhoneAndRole returns a user matching both phone and role.
func (r *UserRepo) GetByPhoneAndRole(ctx context.Context, phone, role string) (*models.User, error) {
	u, err := scanUser(r.db.QueryRow(ctx,
		"SELECT id, name, phone, COALESCE(business_name,''), COALESCE(email,''), password_hash, rating, role, COALESCE(admin_role::varchar,''), is_active, COALESCE(deactivation_reason,''), COALESCE(tier::varchar,'normal'), COALESCE(total_value_sold,0), COALESCE(buyer_tier::varchar,'normal'), COALESCE(otp_verified,false), COALESCE(seller_preference,'sale'), created_at, updated_at FROM users WHERE phone = $1 AND role = $2", phone, role))
	if err == nil && u != nil {
		r.PopulateCapabilities(ctx, u)
	}
	return u, err
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	u, err := scanUser(r.db.QueryRow(ctx,
		"SELECT id, name, phone, COALESCE(business_name,''), COALESCE(email,''), password_hash, rating, role, COALESCE(admin_role::varchar,''), is_active, COALESCE(deactivation_reason,''), COALESCE(tier::varchar,'normal'), COALESCE(total_value_sold,0), COALESCE(buyer_tier::varchar,'normal'), COALESCE(otp_verified,false), COALESCE(seller_preference,'sale'), created_at, updated_at FROM users WHERE email = $1", email))
	if err == nil && u != nil {
		r.PopulateCapabilities(ctx, u)
	}
	return u, err
}


// GetAdminByEmail is a safe version for admin login that handles missing/null enum columns gracefully.
func (r *UserRepo) GetAdminByEmail(ctx context.Context, email string) (*models.User, error) {
	u := &models.User{}
	//nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
	err := r.db.QueryRow(ctx, `
		SELECT id, name, phone, COALESCE(business_name,''), COALESCE(email,''), password_hash,
		       rating, role,
		       COALESCE(admin_role::varchar, ''),
		       is_active,
		       COALESCE(tier::varchar,'normal'), COALESCE(total_value_sold,0),
		       COALESCE(buyer_tier::varchar,'normal'), COALESCE(otp_verified,false),
		       created_at, updated_at
		FROM users WHERE email = $1 AND role = 'admin'`, email,
	).Scan(
		&u.ID, &u.Name, &u.Phone, &u.BusinessName, &u.Email, &u.PasswordHash,
		&u.Rating, &u.Role, &u.AdminRole, &u.IsActive,
		&u.Tier, &u.TotalValueSold,
		&u.BuyerTier, &u.OTPVerified,
		&u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *UserRepo) GetByID(ctx context.Context, id int) (*models.User, error) {
	u, err := scanUser(r.db.QueryRow(ctx,
		"SELECT id, name, phone, COALESCE(business_name,''), COALESCE(email,''), password_hash, rating, role, COALESCE(admin_role::varchar,''), is_active, COALESCE(deactivation_reason,''), COALESCE(tier::varchar,'normal'), COALESCE(total_value_sold,0), COALESCE(buyer_tier::varchar,'normal'), COALESCE(otp_verified,false), COALESCE(seller_preference,'sale'), created_at, updated_at FROM users WHERE id = $1", id))
	if err == nil && u != nil {
		r.PopulateCapabilities(ctx, u)
	}
	return u, err
}

func (r *UserRepo) UpdatePassword(ctx context.Context, id int, hash string) error {
	_, err := r.db.Exec(ctx, `UPDATE users SET password_hash=$1, updated_at=NOW() WHERE id=$2`, hash, id) //nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
	return err
}

func (r *UserRepo) UpdateSellerPreference(ctx context.Context, id int, preference string) error {
	_, err := r.db.Exec(ctx, `UPDATE users SET seller_preference=$1, updated_at=NOW() WHERE id=$2 AND role='seller'`, preference, id) //nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
	return err
}

func (r *UserRepo) GetTier(ctx context.Context, id int) (string, float64, error) {
	var tier string
	var totalValueSold float64
	err := r.db.QueryRow(ctx, //nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
		`SELECT COALESCE(tier::text,'normal'), COALESCE(total_value_sold,0) FROM users WHERE id=$1`, id,
	).Scan(&tier, &totalValueSold)
	return tier, totalValueSold, err
}

func (r *UserRepo) UpdateTotalValueSold(ctx context.Context, id int, addAmount float64) error {
	// Read tier thresholds from platform_settings — fall back to defaults
	// if the keys are missing (e.g. fresh install before admin configures them).
	goldThreshold := constants.DefaultGoldThreshold
	premiumThreshold := constants.DefaultPremiumThreshold
	if v, err := r.GetSetting(ctx, constants.SettingGoldThreshold); err == nil && v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil && f > 0 {
			goldThreshold = f
		}
	}
	if v, err := r.GetSetting(ctx, constants.SettingPremiumThreshold); err == nil && v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil && f > 0 {
			premiumThreshold = f
		}
	}
	//nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
	_, err := r.db.Exec(ctx, `
		UPDATE users
		SET total_value_sold = COALESCE(total_value_sold,0) + $1,
		    tier = CASE
		      WHEN COALESCE(total_value_sold,0) + $1 >= $3 THEN 'premium'::seller_tier
		      WHEN COALESCE(total_value_sold,0) + $1 >= $2 THEN 
		          CASE WHEN tier = 'premium' THEN 'premium'::seller_tier ELSE 'gold'::seller_tier END
		      ELSE tier -- Never demote automatically on sale
		    END,
		    updated_at = NOW()
		WHERE id = $4`, addAmount, goldThreshold, premiumThreshold, id)
	return err
}

// ── OTP ───────────────────────────────────────────────────────────────────────

// GenerateOTP creates a 6-digit OTP for phone-based users (sellers/buyers).
// Defaults user_type to "seller" — use GenerateOTPWithType to specify explicitly.
func (r *UserRepo) GenerateOTP(ctx context.Context, phone, purpose string) (string, error) {
	return r.generateOTP(ctx, phone, "", purpose, "seller")
}

// GenerateOTPWithType creates a 6-digit OTP for phone-based users with an explicit user_type.
func (r *UserRepo) GenerateOTPWithType(ctx context.Context, phone, purpose, userType string) (string, error) {
	return r.generateOTP(ctx, phone, "", purpose, userType)
}

// GenerateEmailOTP creates a 6-digit OTP for email-based users (admins).
func (r *UserRepo) GenerateEmailOTP(ctx context.Context, email, purpose, userType string) (string, error) {
	return r.generateOTP(ctx, "", email, purpose, userType)
}

func (r *UserRepo) generateOTP(ctx context.Context, phone, email, purpose, userType string) (string, error) {
	// All operations run inside a transaction so the rate-limit count check
	// and the INSERT are atomic — prevents double-tap Resend bypassing the limit.
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return "", fmt.Errorf("OTP generation failed: %w", err)
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	// GLOBAL rate limit: max 5 OTP requests per identity per hour (all purposes combined).
	// Prevents brute-force by cycling through different purposes (login, register, reset, etc).
	var globalCount int
	var scanErr error
	if phone != "" {
		//nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
		scanErr = tx.QueryRow(ctx, `
			SELECT COUNT(*) FROM otp_codes
			WHERE phone=$1
			  AND created_at > NOW() - INTERVAL '1 hour'`,
			phone).Scan(&globalCount)
	} else {
		//nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
		scanErr = tx.QueryRow(ctx, `
			SELECT COUNT(*) FROM otp_codes
			WHERE email=$1
			  AND created_at > NOW() - INTERVAL '1 hour'`,
			email).Scan(&globalCount)
	}
	if scanErr != nil {
		return "", fmt.Errorf("rate limit check failed: %w", scanErr)
	}
	if globalCount >= 5 {
		return "", fmt.Errorf("maximum OTP requests exceeded — please wait before requesting another code")
	}

	// Per-purpose rate limit: max 3 OTP requests per identity per purpose per hour.
	// Count ALL rows in the window regardless of lock state — prevents bypass
	// via triggering a lock to reset the counter.
	var recentCount int
	if phone != "" {
		//nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
		scanErr = tx.QueryRow(ctx, `
			SELECT COUNT(*) FROM otp_codes
			WHERE phone=$1 AND purpose=$2
			  AND created_at > NOW() - INTERVAL '1 hour'`,
			phone, purpose).Scan(&recentCount)
	} else {
		//nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
		scanErr = tx.QueryRow(ctx, `
			SELECT COUNT(*) FROM otp_codes
			WHERE email=$1 AND purpose=$2
			  AND created_at > NOW() - INTERVAL '1 hour'`,
			email, purpose).Scan(&recentCount)
	}
	if scanErr != nil {
		return "", fmt.Errorf("rate limit check failed: %w", scanErr)
	}
	if recentCount >= 3 {
		return "", fmt.Errorf("too many OTP requests for this action — please wait before requesting another code")
	}

	n, err := rand.Int(rand.Reader, big.NewInt(1_000_000))
	if err != nil {
		return "", fmt.Errorf("OTP generation failed: %w", err)
	}
	plainCode := fmt.Sprintf("%06d", n.Int64()) //nolint:gosec // G117/G201: internal values only, no user-controlled input reaches this format string

	if r.otpCrypto == nil {
		return "", fmt.Errorf("OTP crypto not configured")
	}
	encryptedCode, err := r.otpCrypto.EncryptOTP(plainCode)
	if err != nil {
		return "", fmt.Errorf("OTP encryption failed: %w", err)
	}

	// Read expiry from platform_settings — admin-controlled, defaults to 15 minutes.
	expiryMinutes := 15
	if v, err2 := r.GetSetting(ctx, "otp_expiry_minutes"); err2 == nil && v != "" {
		if n, err3 := strconv.Atoi(v); err3 == nil && n >= 1 && n <= 60 {
			expiryMinutes = n
		}
	}
	// Compute expiry timestamp in Go to avoid passing expiryMinutes into SQL arithmetic
	expiresAt := time.Now().Add(time.Duration(expiryMinutes) * time.Minute)
	// Invalidate any existing unused OTPs for this identity+purpose
	if phone != "" {
		_, _ = tx.Exec(ctx, `UPDATE otp_codes SET used=TRUE, locked_until=NULL WHERE phone=$1 AND purpose=$2 AND used=FALSE`, phone, purpose)
		_, err = tx.Exec(ctx,
			`INSERT INTO otp_codes (phone, email, code, purpose, user_type, expires_at) VALUES ($1, $2, $3, $4, $5, $6)`,
			phone, nil, encryptedCode, purpose, userType, expiresAt)
	} else {
		_, _ = tx.Exec(ctx, `UPDATE otp_codes SET used=TRUE, locked_until=NULL WHERE email=$1 AND purpose=$2 AND used=FALSE`, email, purpose)
		_, err = tx.Exec(ctx,
			`INSERT INTO otp_codes (phone, email, code, purpose, user_type, expires_at) VALUES ($1, $2, $3, $4, $5, $6)`,
			nil, email, encryptedCode, purpose, userType, expiresAt)
	}
	if err != nil {
		return "", err
	}
	if err := tx.Commit(ctx); err != nil {
		return "", err
	}
	return plainCode, nil
}

// VerifyOTP checks phone-based OTP. Increments attempts on failure (max 5).
// Returns the user_type stored on the OTP row so callers can log it correctly.
func (r *UserRepo) VerifyOTP(ctx context.Context, phone, code, purpose string) error {
	return r.verifyOTP(ctx, phone, "", code, purpose)
}

// VerifyOTPForRole checks phone-based OTP and validates it was issued for the
// given role (user_type). Prevents cross-role OTP reuse — e.g. a seller reset
// OTP cannot be used to reset a buyer password on the same phone.
func (r *UserRepo) VerifyOTPForRole(ctx context.Context, phone, code, purpose, role string) error {
	// Fetch the latest active OTP row for this phone+purpose+role
	var rowID int
	var encryptedCode string
	var rowAttempts int
	//nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
	rowErr := r.db.QueryRow(ctx, `
		SELECT id, code, attempts FROM otp_codes
		WHERE phone=$1 AND purpose=$2 AND user_type=$3
		  AND used=FALSE AND expires_at > NOW()
		  AND (locked_until IS NULL OR locked_until <= NOW())
		ORDER BY created_at DESC LIMIT 1`,
		phone, purpose, role).Scan(&rowID, &encryptedCode, &rowAttempts)
	if rowErr != nil {
		var lockedUntil *string
		_ = r.db.QueryRow(ctx, //nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
			`SELECT locked_until::text FROM otp_codes
			WHERE phone=$1 AND purpose=$2 AND user_type=$3
			  AND used=FALSE AND locked_until > NOW()
			ORDER BY created_at DESC LIMIT 1`,
			phone, purpose, role).Scan(&lockedUntil)
		if lockedUntil != nil {
			return fmt.Errorf("too many failed attempts — locked for 30 minutes")
		}
		return fmt.Errorf("invalid or expired OTP")
	}
	if r.otpCrypto == nil {
		return fmt.Errorf("OTP crypto not configured")
	}
	plainStored, err := r.otpCrypto.DecryptOTP(encryptedCode)
	if err != nil {
		return fmt.Errorf("invalid or expired OTP")
	}

	// Use constant-time comparison to prevent timing-based OTP enumeration (CWE-208).
	if subtle.ConstantTimeCompare([]byte(plainStored), []byte(code)) != 1 {
		newAttempts := rowAttempts + 1
		if newAttempts >= 5 {
			_, _ = r.db.Exec(ctx, //nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
				`UPDATE otp_codes SET attempts=$1, locked_until=NOW()+INTERVAL '30 minutes' WHERE id=$2`,
				newAttempts, rowID)
			return fmt.Errorf("too many failed attempts — locked for 30 minutes")
		}
		_, _ = r.db.Exec(ctx, `UPDATE otp_codes SET attempts=$1 WHERE id=$2`, newAttempts, rowID) //nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
		return fmt.Errorf("invalid OTP — %d attempt(s) remaining", 5-newAttempts)
	}
	_, err = r.db.Exec(ctx, `UPDATE otp_codes SET used=TRUE WHERE id=$1`, rowID) //nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
	return err
}

// VerifyEmailOTP checks email-based OTP. Increments attempts on failure (max 5).
func (r *UserRepo) VerifyEmailOTP(ctx context.Context, email, code, purpose string) error {
	return r.verifyOTP(ctx, "", email, code, purpose)
}

func (r *UserRepo) verifyOTP(ctx context.Context, phone, email, code, purpose string) error {
	// Run the full verify inside a transaction with SELECT FOR UPDATE to prevent
	// concurrent requests from both reading the same OTP row as valid (race condition).
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("invalid or expired OTP")
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	var rowID int
	var encryptedCode string
	var rowAttempts int
	var expiresAt time.Time
	var lockedUntil *time.Time
	var rowErr error
	if phone != "" {
		//nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
		rowErr = tx.QueryRow(ctx, `
			SELECT id, code, attempts, expires_at, locked_until
			FROM otp_codes
			WHERE phone=$1 AND purpose=$2 AND used=FALSE
			ORDER BY created_at DESC LIMIT 1
			FOR UPDATE`,
			phone, purpose).Scan(&rowID, &encryptedCode, &rowAttempts, &expiresAt, &lockedUntil)
	} else {
		//nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
		rowErr = tx.QueryRow(ctx, `
			SELECT id, code, attempts, expires_at, locked_until
			FROM otp_codes
			WHERE email=$1 AND purpose=$2 AND used=FALSE
			ORDER BY created_at DESC LIMIT 1
			FOR UPDATE`,
			email, purpose).Scan(&rowID, &encryptedCode, &rowAttempts, &expiresAt, &lockedUntil)
	}
	if rowErr != nil {
		return fmt.Errorf("invalid or expired OTP")
	}
	// Check lock first — gives a better message than "expired"
	if lockedUntil != nil && lockedUntil.After(time.Now()) {
		return fmt.Errorf("too many failed attempts — locked for 30 minutes")
	}
	if time.Now().After(expiresAt) {
		return fmt.Errorf("invalid or expired OTP")
	}

	if r.otpCrypto == nil {
		return fmt.Errorf("OTP crypto not configured")
	}
	plainStored, err := r.otpCrypto.DecryptOTP(encryptedCode)
	if err != nil {
		return fmt.Errorf("invalid or expired OTP")
	}

	// Constant-time comparison prevents timing-based OTP enumeration (CWE-208).
	if subtle.ConstantTimeCompare([]byte(plainStored), []byte(code)) != 1 {
		newAttempts := rowAttempts + 1
		if newAttempts >= 5 {
			_, _ = tx.Exec(ctx, //nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
				`UPDATE otp_codes SET attempts=$1, locked_until=NOW()+INTERVAL '30 minutes' WHERE id=$2`,
				newAttempts, rowID)
			_ = tx.Commit(ctx)
			return fmt.Errorf("too many failed attempts — locked for 30 minutes")
		}
		_, _ = tx.Exec(ctx, `UPDATE otp_codes SET attempts=$1 WHERE id=$2`, newAttempts, rowID) //nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
		_ = tx.Commit(ctx)
		return fmt.Errorf("invalid OTP — %d attempt(s) remaining", 5-newAttempts)
	}

	_, err = tx.Exec(ctx, `UPDATE otp_codes SET used=TRUE WHERE id=$1`, rowID) //nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// ListOTPLogs returns paginated OTP log entries for the admin panel.
func (r *UserRepo) ListOTPLogs(ctx context.Context, f OTPLogFilter) ([]OTPLogEntry, int, error) {
	limit := f.Limit
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	page := f.Page
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	var phoneArg, emailArg, userTypeArg, statusArg interface{}
	if f.Phone != "" {
		phoneArg = likeWrap(f.Phone)
	}
	if f.Email != "" {
		emailArg = likeWrap(f.Email)
	}
	if f.UserType != "" {
		userTypeArg = f.UserType
	}
	if f.Status != "" {
		statusArg = f.Status
	}

	var total int
	err := r.db.QueryRow(ctx, `
		SELECT COUNT(*) FROM otp_log_view WHERE 1=1
		AND ($1::text IS NULL OR phone ILIKE $1 ESCAPE '\')
		AND ($2::text IS NULL OR email ILIKE $2 ESCAPE '\')
		AND ($3::text IS NULL OR user_type = $3)
		AND ($4::text IS NULL OR status = $4)`,
		phoneArg, emailArg, userTypeArg, statusArg).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.db.Query(ctx, `
		SELECT id, phone, email, purpose, user_type, used, attempts, expires_at, created_at, status, COALESCE(encrypted_code,'')
		FROM otp_log_view WHERE 1=1
		AND ($1::text IS NULL OR phone ILIKE $1 ESCAPE '\')
		AND ($2::text IS NULL OR email ILIKE $2 ESCAPE '\')
		AND ($3::text IS NULL OR user_type = $3)
		AND ($4::text IS NULL OR status = $4)
		ORDER BY created_at DESC LIMIT $5 OFFSET $6`,
		phoneArg, emailArg, userTypeArg, statusArg, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var out []OTPLogEntry
	for rows.Next() {
		var e OTPLogEntry
		if err := rows.Scan(&e.ID, &e.Phone, &e.Email, &e.Purpose, &e.UserType,
			&e.Used, &e.Attempts, &e.ExpiresAt, &e.CreatedAt, &e.Status, &e.EncryptedCode); err != nil {
			return nil, 0, err
		}
		// Decrypt for super admin only — prevents non-super admins from reading OTP codes
		if f.ShowPlainOTP && e.EncryptedCode != "" && r.otpCrypto != nil {
			if plain, err := r.otpCrypto.DecryptOTP(e.EncryptedCode); err == nil {
				e.PlainOTP = plain
			}
		}
		e.EncryptedCode = "" // never send raw ciphertext to frontend
		out = append(out, e)
	}
	return out, total, rows.Err()
}

// ── Buyer tier ────────────────────────────────────────────────────────────────

const (
	BuyerGoldCost    = 10_000.0
	BuyerPremiumCost = 50_000.0
)

// UpgradeBuyerTier upgrades the buyer's tier and records the payment.
// SubmitBuyerUpgradeRequest stores a pending upgrade request — tier is NOT applied until admin approves.
func (r *UserRepo) SubmitBuyerUpgradeRequest(ctx context.Context, userID int, tier string, amount float64, reference, payerPhone, paymentMethod string, durationDays int) (int, error) {
	var id int
	//nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
	err := r.db.QueryRow(ctx, `
		INSERT INTO upgrade_requests (user_id, role, requested_tier, payment_method, payer_phone, reference, amount, duration_days)
		VALUES ($1, 'buyer', $2, $3, $4, $5, $6, $7)
		RETURNING id`,
		userID, tier, paymentMethod, payerPhone, reference, amount, durationDays).Scan(&id)
	if err != nil {
		// Unique partial index violation — one pending request per user per role
		if strings.Contains(err.Error(), "uq_upgrade_requests_one_pending") {
			return 0, fmt.Errorf("you already have a pending upgrade request — wait for admin review before submitting another")
		}
		return 0, err
	}
	return id, nil
}

// ApproveBuyerUpgradeRequest applies the tier upgrade and marks the request approved.
// Delegates to upgrade_requests table — same path as admin approval.
func (r *UserRepo) ApproveBuyerUpgradeRequest(ctx context.Context, requestID, adminID int) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	var buyerID int
	var tier string
	var durationDays *int
	//nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
	err = tx.QueryRow(ctx, `
		SELECT user_id, requested_tier, duration_days
		FROM upgrade_requests WHERE id=$1 AND role='buyer' AND status='pending'`,
		requestID).Scan(&buyerID, &tier, &durationDays)
	if err != nil {
		return fmt.Errorf("upgrade request not found or already processed")
	}

	if durationDays != nil {
		//nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
		_, err = tx.Exec(ctx,
			`UPDATE users SET buyer_tier=$1::buyer_tier, buyer_tier_paid_at=NOW(),
			  buyer_tier_expires_at=NOW() + ($3 * INTERVAL '1 day'), updated_at=NOW()
			 WHERE id=$2`,
			tier, buyerID, *durationDays)
	} else {
		//nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
		_, err = tx.Exec(ctx,
			`UPDATE users SET buyer_tier=$1::buyer_tier, buyer_tier_paid_at=NOW(),
			  buyer_tier_expires_at=NULL, updated_at=NOW()
			 WHERE id=$2`,
			tier, buyerID)
	}
	if err != nil {
		return err
	}
	//nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
	_, err = tx.Exec(ctx, `
		UPDATE upgrade_requests
		SET status='approved', reviewed_by=$1, reviewed_at=NOW(), updated_at=NOW()
		WHERE id=$2`, adminID, requestID)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// RejectBuyerUpgradeRequest marks the request rejected with a reason.
func (r *UserRepo) RejectBuyerUpgradeRequest(ctx context.Context, requestID, adminID int, reason string) error {
	//nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
	tag, err := r.db.Exec(ctx, `
		UPDATE upgrade_requests
		SET status='rejected', reviewed_by=$1, reviewed_at=NOW(), reject_reason=$2, updated_at=NOW()
		WHERE id=$3 AND role='buyer' AND status='pending'`,
		adminID, reason, requestID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("upgrade request not found or already processed")
	}
	return nil
}

// ListBuyerUpgradeRequests returns paginated buyer upgrade requests filtered by status.
func (r *UserRepo) ListBuyerUpgradeRequests(ctx context.Context, status string, page, limit int) ([]map[string]interface{}, int, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	var statusArg interface{}
	if status == "pending" || status == "approved" || status == "rejected" {
		statusArg = status
	}

	var total int
	if err := r.db.QueryRow(ctx, `
		SELECT COUNT(*) FROM upgrade_requests r
		WHERE r.role='buyer'
		AND ($1::text IS NULL OR r.status = $1)`,
		statusArg).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count query failed: %w", err)
	}

	rows, err := r.db.Query(ctx, `
		SELECT r.id, r.user_id, u.name, u.phone, r.requested_tier,
		       r.payment_method, r.payer_phone, r.reference, r.amount, r.duration_days,
		       r.status, r.reject_reason, r.reviewed_by, r.reviewed_at, r.created_at
		FROM upgrade_requests r JOIN users u ON u.id = r.user_id
		WHERE r.role='buyer'
		AND ($1::text IS NULL OR r.status = $1)
		ORDER BY r.created_at DESC LIMIT $2 OFFSET $3`,
		statusArg, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var (
			id, buyerID                int
			name, phone, reqTier       string
			payMethod, payerPhone, ref string
			amt                        float64
			durationDays               *int
			st                         string
			rejectReason               *string
			reviewedBy                 *int
			reviewedAt, createdAt      *time.Time
		)
		if err := rows.Scan(&id, &buyerID, &name, &phone, &reqTier,
			&payMethod, &payerPhone, &ref, &amt, &durationDays,
			&st, &rejectReason, &reviewedBy, &reviewedAt, &createdAt); err != nil {
			return nil, 0, err
		}
		results = append(results, map[string]interface{}{
			"id": id, "buyer_id": buyerID, "buyer_name": name, "buyer_phone": phone,
			"requested_tier": reqTier, "payment_method": payMethod,
			"payer_phone": payerPhone, "reference": ref, "amount": amt, "duration_days": durationDays,
			"status": st, "reject_reason": rejectReason,
			"reviewed_by": reviewedBy, "reviewed_at": reviewedAt, "created_at": createdAt,
		})
	}
	if results == nil {
		results = []map[string]interface{}{}
	}
	return results, total, rows.Err()
}

// GetBuyerTier returns the buyer_tier for a user.
func (r *UserRepo) GetBuyerTier(ctx context.Context, userID int) (string, error) {
	var tier string
	err := r.db.QueryRow(ctx, //nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
		`SELECT COALESCE(buyer_tier::text,'normal') FROM users WHERE id=$1`, userID,
	).Scan(&tier)
	return tier, err
}

func (r *UserRepo) ListOrdersByBuyerPhone(ctx context.Context, phone string) ([]map[string]interface{}, error) {
	rows, err := r.db.Query(ctx, `
		SELECT o.id, o.product_id, o.seller_id, o.buyer_name, o.buyer_phone,
		       COALESCE(o.message,''), o.status, COALESCE(o.seller_note,''),
		       o.created_at, o.updated_at,
		       COALESCE(p.title,''), COALESCE(p.price,0)
		FROM order_requests o
		LEFT JOIN products p ON p.id = o.product_id
		WHERE o.buyer_phone = $1
		ORDER BY o.created_at DESC
		LIMIT 100`, phone)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []map[string]interface{}
	for rows.Next() {
		var id, productID, sellerID int
		var buyerName, buyerPhone, message, status, sellerNote string
		var createdAt, updatedAt interface{}
		var productTitle string
		var productPrice float64
		if err := rows.Scan(&id, &productID, &sellerID, &buyerName, &buyerPhone,
			&message, &status, &sellerNote, &createdAt, &updatedAt,
			&productTitle, &productPrice); err != nil {
			return nil, err
		}
		out = append(out, map[string]interface{}{
			"id": id, "product_id": productID, "seller_id": sellerID,
			"buyer_name": buyerName, "buyer_phone": buyerPhone,
			"message": message, "status": status, "seller_note": sellerNote,
			"created_at": createdAt, "updated_at": updatedAt,
			"product_title": productTitle, "product_price": productPrice,
		})
	}
	return out, rows.Err()
}

// UpdateUnverifiedRegistration updates name, email, and password on an existing
// unverified row (otp_verified=false). Called when a user re-submits the
// registration form after OTP expiry so the new password they typed is stored.
func (r *UserRepo) UpdateUnverifiedRegistration(ctx context.Context, phone, role, passwordHash, name, email string) error {
	_, err := r.db.Exec(ctx, //nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
		`UPDATE users SET name=$1, email=$2, password_hash=$3, updated_at=NOW()
		 WHERE phone=$4 AND role=$5 AND otp_verified=FALSE`,
		name, email, passwordHash, phone, role)
	return err
}

// SetOTPVerified marks a user's phone as OTP-verified (any role).
func (r *UserRepo) SetOTPVerified(ctx context.Context, phone string) error {
	_, err := r.db.Exec(ctx, //nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
		`UPDATE users SET otp_verified=TRUE, updated_at=NOW() WHERE phone=$1`, phone)
	return err
}

// SetOTPVerifiedForRole marks otp_verified on the specific role row.
// Needed when the same phone exists as both seller and buyer.
func (r *UserRepo) SetOTPVerifiedForRole(ctx context.Context, phone, role string) error {
	// Map OTP user_type to DB role
	dbRole := role
	if role == "seller" || role == "buyer" {
		dbRole = role
	} else {
		// Unknown type — fall back to updating all rows for this phone
		return r.SetOTPVerified(ctx, phone)
	}
	_, err := r.db.Exec(ctx, //nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
		`UPDATE users SET otp_verified=TRUE, updated_at=NOW() WHERE phone=$1 AND role=$2`, phone, dbRole)
	return err
}

// RecordFirstLogin sets last_login_at=NOW() only if it was NULL (first ever login).
// Returns true if this was the first login.
func (r *UserRepo) RecordFirstLogin(ctx context.Context, userID int) (bool, error) {
	tag, err := r.db.Exec(ctx, //nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
		`UPDATE users SET last_login_at=NOW(), updated_at=NOW()
		 WHERE id=$1 AND last_login_at IS NULL`, userID)
	if err != nil {
		return false, err
	}
	return tag.RowsAffected() > 0, nil
}

func (r *UserRepo) DeleteOTPLogs(ctx context.Context, status string) (int64, error) {
	var tag interface{ RowsAffected() int64 }
	var err error
	if status == "all" {
		tag, err = r.db.Exec(ctx, `DELETE FROM otp_codes`) //nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
	} else {
		tag, err = r.db.Exec(ctx, `DELETE FROM otp_codes WHERE used=TRUE OR expires_at < NOW()`) //nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
	}
	if err != nil {
		return 0, err
	}
	return tag.RowsAffected(), nil
}

// DevPeekOTP decrypts and returns the latest active OTP for a phone — DEV ONLY.
func (r *UserRepo) DevPeekOTP(ctx context.Context, phone string) (string, error) {
	var encryptedCode string
	//nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
	err := r.db.QueryRow(ctx, `
		SELECT code FROM otp_codes
		WHERE phone=$1 AND used=FALSE AND expires_at > NOW()
		ORDER BY created_at DESC LIMIT 1`, phone,
	).Scan(&encryptedCode)
	if err != nil {
		// Sanitize phone before embedding in error message to prevent log injection
		return "", fmt.Errorf("no active OTP for %s: %w", sanitizeLog(phone), err)
	}
	if r.otpCrypto == nil {
		return "", fmt.Errorf("OTP crypto not configured")
	}
	return r.otpCrypto.DecryptOTP(encryptedCode)
}

// GetSetting reads a single value from platform_settings.
// Returns "" if the key does not exist (no error).
func (r *UserRepo) GetSetting(ctx context.Context, key string) (string, error) {
	var value string
	err := r.db.QueryRow(ctx, `SELECT value FROM platform_settings WHERE key=$1`, key).Scan(&value) //nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
	if err != nil {
		return "", nil
	}
	return value, nil
}

// GetOTPUserType returns the user_type stored on the latest active (unused, unexpired) OTP for a phone.
// Used by the service layer to populate activity log entries correctly.
// Filters to used=FALSE AND expires_at > NOW() so stale OTPs from previous sessions
// cannot cause the wrong role row to be marked as verified.
func (r *UserRepo) GetOTPUserType(ctx context.Context, phone string) string {
	var userType string
	_ = r.db.QueryRow(ctx, //nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
		`SELECT COALESCE(user_type,'seller') FROM otp_codes
		 WHERE phone=$1 AND used=FALSE AND expires_at > NOW()
		 ORDER BY created_at DESC LIMIT 1`,
		phone).Scan(&userType)
	if userType == "" {
		return "seller"
	}
	return userType
}

// LogActivity inserts one row into user_activity_logs. Fire-and-forget safe (async).
func (r *UserRepo) LogActivity(ctx context.Context, event, phone, email, userType, ip, meta string) error {
	r.asyncLogger.Log(AsyncActivityLog{
		Event:     event,
		Phone:     phone,
		Email:     email,
		UserType:  userType,
		IPAddress: ip,
		Meta:      meta,
	})
	return nil // Always returns nil since logging is async
}

// ListActivityLogs returns paginated user_activity_logs for the analytics panel.
func (r *UserRepo) ListActivityLogs(ctx context.Context, f ActivityLogFilter) ([]ActivityLog, int, error) {
	limit := f.Limit
	if limit <= 0 || limit > 200 {
		limit = 100
	}
	page := f.Page
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	var eventArg, phoneArg, userTypeArg interface{}
	var dateFromArg, dateToArg interface{}
	if f.Event != "" {
		eventArg = f.Event
	}
	if f.Phone != "" {
		phoneArg = likeWrap(f.Phone)
	}
	if f.UserType != "" {
		userTypeArg = f.UserType
	}
	if f.DateFrom != nil {
		dateFromArg = f.DateFrom
	}
	if f.DateTo != nil {
		dateToArg = f.DateTo.Add(time.Hour*23 + time.Minute*59 + time.Second*59)
	}

	var total int
	if err := r.db.QueryRow(ctx, `
		SELECT COUNT(*) FROM user_activity_logs WHERE 1=1
		AND ($1::text IS NULL OR event = $1)
		AND ($2::text IS NULL OR phone ILIKE $2 ESCAPE '\')
		AND ($3::text IS NULL OR user_type = $3)
		AND ($4::timestamptz IS NULL OR created_at >= $4)
		AND ($5::timestamptz IS NULL OR created_at <= $5)`,
		eventArg, phoneArg, userTypeArg, dateFromArg, dateToArg).Scan(&total); err != nil {
		return nil, 0, err
	}

	rows, err := r.db.Query(ctx, `
		SELECT id, event, COALESCE(phone,''), COALESCE(email,''), user_type, COALESCE(ip_address,''), COALESCE(meta,''), created_at
		FROM user_activity_logs WHERE 1=1
		AND ($1::text IS NULL OR event = $1)
		AND ($2::text IS NULL OR phone ILIKE $2 ESCAPE '\')
		AND ($3::text IS NULL OR user_type = $3)
		AND ($4::timestamptz IS NULL OR created_at >= $4)
		AND ($5::timestamptz IS NULL OR created_at <= $5)
		ORDER BY created_at DESC LIMIT $6 OFFSET $7`,
		eventArg, phoneArg, userTypeArg, dateFromArg, dateToArg, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	
	var out []ActivityLog
	for rows.Next() {
		var e ActivityLog
		if err := rows.Scan(&e.ID, &e.Event, &e.Phone, &e.Email, &e.UserType, &e.IPAddress, &e.Meta, &e.CreatedAt); err != nil {
			return nil, 0, err
		}
		out = append(out, e)
	}
	return out, total, rows.Err()
}

// UpdateBuyerProfile updates name and email for a buyer.
func (r *UserRepo) UpdateBuyerProfile(ctx context.Context, userID int, name, email string) error {
	if email != "" {
		//nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
		_, err := r.db.Exec(ctx,
			`UPDATE users SET name=$1, email=$2, updated_at=NOW() WHERE id=$3 AND role='buyer'`,
			name, email, userID)
		return err
	}
	_, err := r.db.Exec(ctx,
		`UPDATE users SET name=$1, updated_at=NOW() WHERE id=$2 AND role='buyer'`,
		name, userID)
	return err
}

func (r *UserRepo) ListBuyerUpgradeRequestsByUser(ctx context.Context, userID int) ([]map[string]interface{}, error) {
	//nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
	rows, err := r.db.Query(ctx, `
		SELECT id, requested_tier, payment_method, payer_phone, reference,
		       amount, status, reject_reason, created_at, reviewed_at
		FROM upgrade_requests
		WHERE user_id=$1 AND role='buyer'
		ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []map[string]interface{}
	for rows.Next() {
		var (
			rowID                    int
			tier, method, phone, ref string
			amount                   float64
			status, rejectReason     string
			createdAt, reviewedAt    interface{}
		)
		if err := rows.Scan(&rowID, &tier, &method, &phone, &ref, &amount, &status, &rejectReason, &createdAt, &reviewedAt); err != nil {
			return nil, err
		}
		out = append(out, map[string]interface{}{
			"id":             rowID,
			"requested_tier": tier,
			"payment_method": method,
			"payer_phone":    phone,
			"reference":      ref,
			"amount":         amount,
			"status":         status,
			"reject_reason":  rejectReason,
			"created_at":     createdAt,
			"reviewed_at":    reviewedAt,
		})
	}
	return out, rows.Err()
}

func (r *UserRepo) ListBuyerNotifications(ctx context.Context, userID int) ([]map[string]interface{}, error) {
	//nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
	rows, err := r.db.Query(ctx, `
		SELECT id, title, message, type, is_read, created_at
		FROM notifications
		WHERE user_id=$1
		ORDER BY created_at DESC
		LIMIT 50`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []map[string]interface{}
	for rows.Next() {
		var id int
		var title, msg, ntype string
		var isRead bool
		var createdAt interface{}
		if err := rows.Scan(&id, &title, &msg, &ntype, &isRead, &createdAt); err != nil {
			return nil, err
		}
		out = append(out, map[string]interface{}{
			"id":         id,
			"title":      title,
			"message":    msg,
			"type":       ntype,
			"is_read":    isRead,
			"created_at": createdAt,
		})
	}
	return out, rows.Err()
}

// MarkBuyerNotificationRead marks a notification as read, scoped to the buyer.
func (r *UserRepo) MarkBuyerNotificationRead(ctx context.Context, notifID, userID int) error {
	//nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
	_, err := r.db.Exec(ctx,
		`UPDATE notifications SET is_read=TRUE WHERE id=$1 AND user_id=$2`,
		notifID, userID)
	return err
}

// SetDemand updates the user's current demand (FIND or LIST)
func (r *UserRepo) SetDemand(ctx context.Context, userID int, demand string) error {
	_, err := r.db.Exec(ctx,
		`UPDATE users SET current_demand=$1, updated_at=NOW() WHERE id=$2`,
		demand, userID)
	return err
}
