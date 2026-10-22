package security

import (
	"crypto/rand"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(plain string) (string, error) {
	// Cost 12 is the minimum recommended for financial applications.
	// bcrypt.DefaultCost (10) is too fast for brute-force resistance.
	b, err := bcrypt.GenerateFromPassword([]byte(plain), 12)
	return string(b), err
}

func CheckPassword(hash, plain string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain)) == nil
}

// RandomPassword generates a cryptographically random 16-char hex password.
// Returns an error instead of a weak fallback — callers must handle entropy failure.
func RandomPassword() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		// crypto/rand failure is a fatal system-level error (e.g. /dev/urandom unavailable).
		// Panicking here is safer than returning a predictable timestamp-based value (CWE-338).
		panic("crypto/rand unavailable — cannot generate secure password")
	}
	return hex.EncodeToString(b)
}
