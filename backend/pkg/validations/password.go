package validations

import "unicode"

const MinPasswordLength = 8

// ValidatePassword validates password requirements.
// Returns (valid bool, errorKey string).
//
// Requirements:
//   - At least 8 characters
//   - At least one uppercase letter
//   - At least one lowercase letter
//   - At least one digit
//
// Note: special character is strongly encouraged but not required
// so as not to break SMS-originated passwords.
func ValidatePassword(password string) (bool, string) {
	if len(password) < MinPasswordLength {
		return false, "password_min"
	}

	var hasUpper, hasLower, hasDigit bool
	for _, ch := range password {
		switch {
		case unicode.IsUpper(ch):
			hasUpper = true
		case unicode.IsLower(ch):
			hasLower = true
		case unicode.IsDigit(ch):
			hasDigit = true
		}
	}

	if !hasUpper {
		return false, "password_needs_upper"
	}
	if !hasLower {
		return false, "password_needs_lower"
	}
	if !hasDigit {
		return false, "password_needs_digit"
	}

	return true, ""
}
