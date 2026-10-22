package referral

import "regexp"

// Referral code format: REF + 8 alphanumeric characters (excluding I, O, 0, 1 for clarity)
var CodePattern = regexp.MustCompile(`^REF[A-HJ-NP-Z2-9]{8}$`)

// ValidateReferralCode checks if a referral code matches the expected format
func ValidateReferralCode(code string) bool {
	return CodePattern.MatchString(code)
}
