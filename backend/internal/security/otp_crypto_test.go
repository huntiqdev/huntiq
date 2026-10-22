package security_test

import (
	"encoding/base64"
	"testing"

	"github.com/your-org/used-tools-marketplace/internal/security"
)

func TestOTPCrypto_EncryptDecrypt(t *testing.T) {
	// Generate valid 32-byte key base64
	key := make([]byte, 32)
	for i := range key {
		key[i] = byte(i + 1)
	}
	keyBase64 := base64.StdEncoding.EncodeToString(key)

	crypto, err := security.NewOTPCrypto(keyBase64)
	if err != nil {
		t.Fatalf("failed to create OTPCrypto: %v", err)
	}

	code := "849201"
	encrypted, err := crypto.EncryptOTP(code)
	if err != nil {
		t.Fatalf("encryption failed: %v", err)
	}

	if encrypted == code {
		t.Errorf("encrypted code should not match plain code")
	}

	decrypted, err := crypto.DecryptOTP(encrypted)
	if err != nil {
		t.Fatalf("decryption failed: %v", err)
	}

	if decrypted != code {
		t.Errorf("expected decrypted %s, got %s", code, decrypted)
	}
}

func TestOTPCrypto_InvalidKey(t *testing.T) {
	_, err := security.NewOTPCrypto("too-short")
	if err == nil {
		t.Errorf("expected error with invalid key")
	}
}
