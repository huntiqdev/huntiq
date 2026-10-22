package service

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/your-org/used-tools-marketplace/internal/models"
	"github.com/your-org/used-tools-marketplace/internal/repository"
)

type AttributeService struct {
	repo *repository.AttributeRepo
}

func NewAttributeService(repo *repository.AttributeRepo) *AttributeService {
	return &AttributeService{repo: repo}
}

// GetAttributeDefinition retrieves an attribute by ID
func (s *AttributeService) GetAttributeDefinition(ctx context.Context, id int) (*models.AttributeDefinition, error) {
	return s.repo.GetAttributeDefinition(ctx, id)
}

// ListAttributeDefinitions retrieves all attributes
func (s *AttributeService) ListAttributeDefinitions(ctx context.Context) ([]models.AttributeDefinition, error) {
	return s.repo.ListAttributeDefinitions(ctx)
}

// GetCategoryAttributes retrieves applicable attributes for a category
func (s *AttributeService) GetCategoryAttributes(ctx context.Context, categoryID int) ([]models.ApplicableAttribute, error) {
	return s.repo.GetCategoryAttributes(ctx, categoryID)
}

// GetProductAttributes retrieves all attributes for a product
func (s *AttributeService) GetProductAttributes(ctx context.Context, productID int) ([]models.ProductAttribute, error) {
	return s.repo.GetProductAttributes(ctx, productID)
}

// SetProductAttributes sets multiple product attributes (validates and saves)
func (s *AttributeService) SetProductAttributes(ctx context.Context, productID int, categoryID int, attrs []models.ProductAttribute) error {
	// Get applicable attributes for the category
	applicable, err := s.repo.GetCategoryAttributes(ctx, categoryID)
	if err != nil {
		return fmt.Errorf("failed to get category attributes: %w", err)
	}
	
	// Build map of applicable attributes for validation
	applicableMap := make(map[int]*models.ApplicableAttribute)
	for i := range applicable {
		applicableMap[applicable[i].AttributeID] = &applicable[i]
	}
	
	// Validate each attribute
	for _, attr := range attrs {
		app, ok := applicableMap[attr.AttributeID]
		if !ok {
			return fmt.Errorf("attribute %d not applicable to this category", attr.AttributeID)
		}
		
		// Check required
		if app.IsRequired && isAttributeEmpty(&attr) {
			return fmt.Errorf("attribute '%s' is required", app.NameEn)
		}
		
		// Validate based on data type and rules
		if err := s.validateAttributeValue(&attr, app); err != nil {
			return fmt.Errorf("validation failed for '%s': %w", app.NameEn, err)
		}
		
		// Set product ID
		attr.ProductID = productID
		
		// Save to database
		if err := s.repo.SetProductAttribute(ctx, &attr); err != nil {
			return fmt.Errorf("failed to save attribute: %w", err)
		}
	}
	
	return nil
}

// CreateAttributeDefinition creates a new attribute (admin)
func (s *AttributeService) CreateAttributeDefinition(ctx context.Context, attr *models.AttributeDefinition) (int, error) {
	// Validate
	if attr.NameEn == "" {
		return 0, fmt.Errorf("name_en is required")
	}
	if attr.Slug == "" {
		attr.Slug = s.generateSlug(attr.NameEn)
	}
	if !s.isValidSlug(attr.Slug) {
		return 0, fmt.Errorf("invalid slug format")
	}
	if !s.isValidDataType(attr.DataType) {
		return 0, fmt.Errorf("invalid data_type: must be text, number, boolean, select, multiselect, range, or date")
	}
	
	return s.repo.CreateAttributeDefinition(ctx, attr)
}

// AssignAttributeToCategory assigns an attribute to a category (admin)
func (s *AttributeService) AssignAttributeToCategory(ctx context.Context, ca *models.CategoryAttribute) error {
	return s.repo.AssignAttributeToCategory(ctx, ca)
}

// RemoveAttributeFromCategory removes an attribute from a category (admin)
func (s *AttributeService) RemoveAttributeFromCategory(ctx context.Context, categoryID, attributeID int) error {
	return s.repo.RemoveAttributeFromCategory(ctx, categoryID, attributeID)
}

// Helper: Validate attribute value based on data type and rules
func (s *AttributeService) validateAttributeValue(attr *models.ProductAttribute, app *models.ApplicableAttribute) error {
	switch app.DataType {
	case "text":
		if attr.ValueText == nil {
			return fmt.Errorf("text value required")
		}
		// Check pattern if provided
		if pattern, ok := app.ValidationRules["pattern"].(string); ok {
			matched, _ := regexp.MatchString(pattern, *attr.ValueText)
			if !matched {
				return fmt.Errorf("value does not match required pattern")
			}
		}
		
	case "number", "range":
		if attr.ValueNumber == nil {
			return fmt.Errorf("number value required")
		}
		// Check min/max
		if min, ok := app.ValidationRules["min"].(float64); ok {
			if *attr.ValueNumber < min {
				return fmt.Errorf("value must be at least %.2f", min)
			}
		}
		if max, ok := app.ValidationRules["max"].(float64); ok {
			if *attr.ValueNumber > max {
				return fmt.Errorf("value must be at most %.2f", max)
			}
		}
		
	case "boolean":
		if attr.ValueBool == nil {
			return fmt.Errorf("boolean value required")
		}
		
	case "date":
		if attr.ValueDate == nil {
			return fmt.Errorf("date value required")
		}
		
	case "select":
		if attr.ValueText == nil {
			return fmt.Errorf("selection required")
		}
		// Check if value is in options
		if options, ok := app.ValidationRules["options"].([]interface{}); ok {
			valid := false
			for _, opt := range options {
				if optStr, ok := opt.(string); ok && optStr == *attr.ValueText {
					valid = true
					break
				}
			}
			if !valid {
				return fmt.Errorf("value must be one of the predefined options")
			}
		}
		
	case "multiselect":
		if attr.ValueJSON == nil {
			return fmt.Errorf("multiselect value required")
		}
		// Validate each selected value
		if values, ok := attr.ValueJSON["values"].([]interface{}); ok {
			if options, ok := app.ValidationRules["options"].([]interface{}); ok {
				for _, val := range values {
					valStr, _ := val.(string)
					valid := false
					for _, opt := range options {
						if optStr, ok := opt.(string); ok && optStr == valStr {
							valid = true
							break
						}
					}
					if !valid {
						return fmt.Errorf("invalid option: %s", valStr)
					}
				}
			}
		}
	}
	
	return nil
}

// Helper: Check if attribute value is empty
func isAttributeEmpty(attr *models.ProductAttribute) bool {
	return attr.ValueText == nil &&
		attr.ValueNumber == nil &&
		attr.ValueBool == nil &&
		attr.ValueDate == nil &&
		attr.ValueJSON == nil
}



// Helper: Generate slug from name
func (s *AttributeService) generateSlug(name string) string {
	slug := strings.ToLower(name)
	reg := regexp.MustCompile("[^a-z0-9\\s-]")
	slug = reg.ReplaceAllString(slug, "")
	slug = strings.ReplaceAll(slug, " ", "-")
	reg = regexp.MustCompile("-+")
	slug = reg.ReplaceAllString(slug, "-")
	return strings.Trim(slug, "-")
}

// Helper: Validate slug format
func (s *AttributeService) isValidSlug(slug string) bool {
	match, _ := regexp.MatchString("^[a-z0-9]+(-[a-z0-9]+)*$", slug)
	return match
}

// Helper: Validate data type
func (s *AttributeService) isValidDataType(dt string) bool {
	validTypes := map[string]bool{
		"text": true, "number": true, "boolean": true,
		"select": true, "multiselect": true, "range": true, "date": true,
	}
	return validTypes[dt]
}
