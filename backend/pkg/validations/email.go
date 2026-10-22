package validations

import "regexp"

var emailPattern = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]{2,}$`)

// ValidateEmail validates email format
func ValidateEmail(email string) bool {
	if email == "" {
		return true // Email is optional in most cases
	}
	return emailPattern.MatchString(email)
}
