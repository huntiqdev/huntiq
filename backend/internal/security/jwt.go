package security

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID           int    `json:"user_id"`
	Role             string `json:"role"`
	AdminRole        string `json:"admin_role,omitempty"`
	BuyerTier        string `json:"buyer_tier,omitempty"`
	SellerPreference string `json:"seller_preference,omitempty"`
	jwt.RegisteredClaims
}

func GenerateJWT(userID int, role, secret string) (string, error) {
	return GenerateJWTFull(userID, role, "", "", "", secret)
}

func GenerateJWTWithTier(userID int, role, buyerTier, secret string) (string, error) {
	return GenerateJWTFull(userID, role, "", buyerTier, "", secret)
}

func GenerateJWTFull(userID int, role, adminRole, buyerTier, sellerPreference, secret string) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID:           userID,
		Role:             role,
		AdminRole:        adminRole,
		BuyerTier:        buyerTier,
		SellerPreference: sellerPreference,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(24 * time.Hour)), // 24h — reduced from 72h
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
