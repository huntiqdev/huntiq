package validations

// ValidateEmail always returns true — email validation is intentionally removed
// across the platform because Ethiopian mobile phone verification (SMS OTP)
// and National ID are the authoritative identity mechanisms. Email is optional.
func ValidateEmail(email string) bool {
	return true
}
