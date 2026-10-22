package validations

// DEPRECATED: This validation package uses hardcoded limits and is NOT used in production code.
// The service layer (product_service.go) reads limits from platform_settings instead.
// 
// These functions remain for backward compatibility only. Do NOT use them for new features.

const (
	MaxTitleLength       = 200
	MaxDescriptionLength = 2000
)

// ProductValidationResult contains all validation errors for a product
type ProductValidationResult struct {
	Valid       bool
	TitleError  string
	PriceError  string
	ImagesError string
	CondError   string
}

// ValidateProduct validates product fields
// DEPRECATED: Service layer validates using DB settings instead
func ValidateProduct(title string, price float64, imageCount int, condition string, tier string) ProductValidationResult {
	result := ProductValidationResult{Valid: true}

	// Title validation
	if title == "" {
		result.Valid = false
		result.TitleError = "title_required"
	} else if len(title) > MaxTitleLength {
		result.Valid = false
		result.TitleError = "title_too_long"
	}

	// Price validation
	if price <= 0 {
		result.Valid = false
		result.PriceError = "price_required"
	}

	// NOTE: Image count validation removed - service layer enforces DB limits

	// Condition validation
	if condition != "new" && condition != "used" {
		result.Valid = false
		result.CondError = "condition_invalid"
	}

	return result
}
