package validations

// LocationValidationResult contains validation errors for location
type LocationValidationResult struct {
	Valid       bool
	RegionError string
}

// ValidateLocation validates location fields
func ValidateLocation(regionID int) LocationValidationResult {
	result := LocationValidationResult{Valid: true}

	if regionID <= 0 {
		result.Valid = false
		result.RegionError = "region_required"
	}

	return result
}
