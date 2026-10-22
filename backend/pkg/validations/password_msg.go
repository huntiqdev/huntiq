package validations

// PasswordErrorMessage converts a validation error key to a human-readable message.
func PasswordErrorMessage(errKey string) string {
	switch errKey {
	case "password_min":
		return "password must be at least 8 characters"
	case "password_needs_upper":
		return "password must contain at least one uppercase letter"
	case "password_needs_lower":
		return "password must contain at least one lowercase letter"
	case "password_needs_digit":
		return "password must contain at least one number"
	default:
		return "invalid password"
	}
}
