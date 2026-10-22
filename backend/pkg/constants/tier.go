package constants

// Tier names
const (
	TierNormal  = "normal"
	TierGold    = "gold"
	TierPremium = "premium"
)

// TierRank defines the tier hierarchy (0 = lowest, 2 = highest)
var TierRank = map[string]int{
	TierNormal:  0,
	TierGold:    1,
	TierPremium: 2,
}

// ValidTiers is a list of all valid tier values
var ValidTiers = []string{TierNormal, TierGold, TierPremium}

// IsTierValid checks if a tier string is valid
func IsTierValid(tier string) bool {
	_, exists := TierRank[tier]
	return exists
}

// MeetsRequiredTier checks if userTier meets or exceeds requiredTier
func MeetsRequiredTier(userTier, requiredTier string) bool {
	userRank := TierRank[userTier]
	requiredRank := TierRank[requiredTier]
	return userRank >= requiredRank
}

// CanUpgradeTo checks if a user can upgrade from currentTier to targetTier
func CanUpgradeTo(currentTier, targetTier string) bool {
	currentRank := TierRank[currentTier]
	targetRank := TierRank[targetTier]
	return targetRank > currentRank
}
