package security

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID           int      `json:"user_id"`
	Name             string   `json:"name,omitempty"`
	Phone            string   `json:"phone,omitempty"`
	Role             string   `json:"role"`
	AdminRole        string   `json:"admin_role,omitempty"`
	BuyerTier        string   `json:"buyer_tier,omitempty"`
	SellerPreference string   `json:"seller_preference,omitempty"`
	Capabilities     []string `json:"capabilities,omitempty"`
	SessionID        int      `json:"session_id,omitempty"`
	SessionToken     string   `json:"session_token,omitempty"`
	jwt.RegisteredClaims
}

// HasCapability checks if the JWT token contains the given capability.
// When Capabilities are present in the token, they are authoritative.
func (c *Claims) HasCapability(cap string) bool {
	if len(c.Capabilities) > 0 {
		for _, capItem := range c.Capabilities {
			if capItem == cap {
				return true
			}
		}
		return false
	}

	// Fallback logic only for legacy tokens generated before multi-capability support
	if cap == "buyer" {
		return true // all registered accounts are buyers
	}
	if cap == "seller" && c.Role == "seller" {
		return true
	}
	if cap == "landlord" && (c.SellerPreference == "rent" || c.SellerPreference == "both") {
		return true
	}
	return false
}

func GenerateJWT(userID int, role, secret string) (string, error) {
	return GenerateJWTFull(userID, role, "", "", "", secret)
}

func GenerateJWTWithTier(userID int, role, buyerTier, secret string) (string, error) {
	return GenerateJWTFull(userID, role, "", buyerTier, "", secret)
}

func GenerateJWTFull(userID int, role, adminRole, buyerTier, sellerPreference, secret string) (string, error) {
	var caps []string
	if role == "seller" {
		caps = []string{"buyer", "seller"}
		if sellerPreference == "rent" || sellerPreference == "both" {
			caps = append(caps, "landlord")
		}
	} else if role == "buyer" {
		caps = []string{"buyer"}
	} else if role == "admin" {
		caps = []string{"buyer", "seller", "admin"}
	}
	return GenerateJWTWithCapabilities(userID, role, adminRole, buyerTier, sellerPreference, caps, secret)
}

func GenerateJWTWithCapabilities(userID int, role, adminRole, buyerTier, sellerPreference string, capabilities []string, secret string) (string, error) {
	return GenerateJWTWithProfile(userID, "", "", role, adminRole, buyerTier, sellerPreference, capabilities, secret)
}

func GenerateJWTWithProfile(userID int, name, phone, role, adminRole, buyerTier, sellerPreference string, capabilities []string, secret string) (string, error) {
	now := time.Now()
	if len(capabilities) == 0 {
		if role == "seller" {
			capabilities = []string{"buyer", "seller"}
			if sellerPreference == "rent" || sellerPreference == "both" {
				capabilities = append(capabilities, "landlord")
			}
		} else if role == "buyer" {
			capabilities = []string{"buyer"}
		} else if role == "admin" {
			capabilities = []string{"buyer", "seller", "admin"}
		}
	}
	// Ensure 'buyer' is always in capabilities for any registered user
	hasBuyer := false
	for _, c := range capabilities {
		if c == "buyer" {
			hasBuyer = true
			break
		}
	}
	if !hasBuyer && role != "admin" {
		capabilities = append([]string{"buyer"}, capabilities...)
	}

	claims := Claims{
		UserID:           userID,
		Name:             name,
		Phone:            phone,
		Role:             role,
		AdminRole:        adminRole,
		BuyerTier:        buyerTier,
		SellerPreference: sellerPreference,
		Capabilities:     capabilities,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(24 * time.Hour)), // 24h
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
}

// GenerateAdminJWT creates a token for an administrator bound to a specific server session ID.
func GenerateAdminJWT(userID int, adminRole string, sessionID int, sessionToken, secret string) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID:       userID,
		Role:         "admin",
		AdminRole:    adminRole,
		Capabilities: []string{"buyer", "seller", "admin"},
		SessionID:    sessionID,
		SessionToken: sessionToken,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(12 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
}


func ParseJWT(tokenStr, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
