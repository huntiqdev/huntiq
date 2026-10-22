// Package validations provides centralized validation logic shared across the application.
package validations

import (
	"regexp"
	"strings"
)

var (
	// Ethiopian phone number patterns
	ethioPhonePattern = regexp.MustCompile(`^(\+?251|0)?([79]\d{8})$`)
)

// PhoneResult represents the result of phone validation
type PhoneResult struct {
	Valid      bool
	Normalized string // +251XXXXXXXXX format
	Error      string
}

// ValidateEthiopianPhone validates and normalizes Ethiopian phone numbers.
// Supports formats: 0912345678, 0712345678, +251912345678, 251912345678, 912345678
// Returns normalized format: +251XXXXXXXXX
func ValidateEthiopianPhone(phone string) PhoneResult {
	if phone == "" {
		return PhoneResult{Valid: false, Error: "phone_required"}
	}

	// Clean input: remove spaces, dashes, parentheses
	cleaned := strings.ReplaceAll(phone, " ", "")
	cleaned = strings.ReplaceAll(cleaned, "-", "")
	cleaned = strings.ReplaceAll(cleaned, "(", "")
	cleaned = strings.ReplaceAll(cleaned, ")", "")

	// Match pattern
	matches := ethioPhonePattern.FindStringSubmatch(cleaned)
	if len(matches) < 3 {
		return PhoneResult{Valid: false, Error: "phone_invalid"}
	}

	// Normalize to +251XXXXXXXXX
	normalized := "+251" + matches[2]

	return PhoneResult{
		Valid:      true,
		Normalized: normalized,
		Error:      "",
	}
}

// IsValidPhone checks if a phone number is valid (boolean helper)
func IsValidPhone(phone string) bool {
	result := ValidateEthiopianPhone(phone)
	return result.Valid
}
