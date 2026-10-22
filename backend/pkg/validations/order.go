package validations

// OrderValidationResult contains validation errors for an order request
type OrderValidationResult struct {
	Valid      bool
	NameError  string
	PhoneError string
}

// ValidateOrderRequest validates order request fields
func ValidateOrderRequest(name, phone string, phonePreFilled bool) OrderValidationResult {
	result := OrderValidationResult{Valid: true}

	// If phone is pre-filled from authenticated user, skip validation
	if phonePreFilled {
		return result
	}

	// Name validation
	if name == "" {
		result.Valid = false
		result.NameError = "name_required"
	}

	// Phone validation
	phoneResult := ValidateEthiopianPhone(phone)
	if !phoneResult.Valid {
		result.Valid = false
		result.PhoneError = phoneResult.Error
	}

	return result
}
