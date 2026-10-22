package security_test

import (
	"testing"

	"github.com/your-org/used-tools-marketplace/internal/security"
)

func TestPasswordHashAndVerify(t *testing.T) {
	password := "SecureAdmin123!"

	hash, err := security.HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	if !security.CheckPassword(hash, password) {
		t.Errorf("CheckPassword should return true for correct password")
	}

	if security.CheckPassword(hash, "WrongPassword!") {
		t.Errorf("CheckPassword should return false for incorrect password")
	}
}
