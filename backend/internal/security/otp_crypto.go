package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"
)

// OTPCrypto handles OTP encryption/decryption with injected key
type OTPCrypto struct {
	key []byte
}

// NewOTPCrypto creates a new OTP crypto instance with the provided base64-encoded key
func NewOTPCrypto(keywordBase64 string) (*OTPCrypto, error) {
	decoded, err := base64.StdEncoding.DecodeString(keywordBase64)
	if err != nil || len(decoded) < 32 {
		return nil, errors.New("OTP_ENCRYPT_KEY must be a base64-encoded string decoding to at least 32 bytes")
	}
	return &OTPCrypto{
		key: decoded[:32],
	}, nil
}

// EncryptOTP encrypts a 6-digit OTP code using AES-256-GCM.
// Returns a hex-encoded string: 24-byte nonce || ciphertext || 16-byte tag.
func (o *OTPCrypto) EncryptOTP(plainCode string) (string, error) {
	block, err := aes.NewCipher(o.key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize()) // 12 bytes
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	ciphertext := gcm.Seal(nonce, nonce, []byte(plainCode), nil)
	return hex.EncodeToString(ciphertext), nil
}

// DecryptOTP decrypts a hex-encoded AES-256-GCM ciphertext back to the 6-digit code.
func (o *OTPCrypto) DecryptOTP(hexCipher string) (string, error) {
	data, err := hex.DecodeString(hexCipher)
	if err != nil {
		return "", errors.New("invalid OTP ciphertext encoding")
	}
	block, err := aes.NewCipher(o.key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("ciphertext too short")
	}
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plain, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", errors.New("OTP decryption failed — tampered or wrong key")
	}
	return string(plain), nil
}
