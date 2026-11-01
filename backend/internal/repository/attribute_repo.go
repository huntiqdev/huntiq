package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/your-org/used-tools-marketplace/internal/models"
)

type AttributeRepo struct {
	db *pgxpool.Pool
}

func NewAttributeRepo(db *pgxpool.Pool) *AttributeRepo {
	return &AttributeRepo{db: db}
}

// GetAttributeDefinition retrieves an attribute by ID
func (r *AttributeRepo) GetAttributeDefinition(ctx context.Context, id int) (*models.AttributeDefinition, error) {
	query := `
		SELECT id,
		       COALESCE(name_en,'') AS name_en,
		       COALESCE(name_am,'') AS name_am,
		       COALESCE(name_om,'') AS name_om,
		       slug,
		       COALESCE(description,'') AS description,
		       data_type,
		       validation_rules,
		       COALESCE(unit,'') AS unit,
		       is_required, is_searchable, is_filterable,
		       display_order, created_at, updated_at
		FROM attribute_definitions
		WHERE id = $1`
	
	var attr models.AttributeDefinition
	var validationJSON []byte
	
	err := r.db.QueryRow(ctx, query, id).Scan(
		&attr.ID, &attr.NameEn, &attr.NameAm, &attr.NameOm, &attr.Slug,
		&attr.Description, &attr.DataType, &validationJSON, &attr.Unit,
		&attr.IsRequired, &attr.IsSearchable, &attr.IsFilterable,
		&attr.DisplayOrder, &attr.CreatedAt, &attr.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	
	if len(validationJSON) > 0 {
		json.Unmarshal(validationJSON, &attr.ValidationRules)
	}
	
	return &attr, nil
}

// GetAttributeBySlug retrieves an attribute definition by slug
func (r *AttributeRepo) GetAttributeBySlug(ctx context.Context, slug string) (*models.AttributeDefinition, error) {
	query := `
		SELECT id,
		       COALESCE(name_en,'') AS name_en,
		       COALESCE(name_am,'') AS name_am,
		       COALESCE(name_om,'') AS name_om,
		       slug,
		       COALESCE(description,'') AS description,
		       data_type,
		       validation_rules,
		       COALESCE(unit,'') AS unit,
		       is_required, is_searchable, is_filterable,
		       display_order, created_at, updated_at
		FROM attribute_definitions
		WHERE slug = $1`
	
	var attr models.AttributeDefinition
	var validationJSON []byte
	
	err := r.db.QueryRow(ctx, query, slug).Scan(
		&attr.ID, &attr.NameEn, &attr.NameAm, &attr.NameOm, &attr.Slug,
		&attr.Description, &attr.DataType, &validationJSON, &attr.Unit,
		&attr.IsRequired, &attr.IsSearchable, &attr.IsFilterable,
		&attr.DisplayOrder, &attr.CreatedAt, &attr.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	
	if len(validationJSON) > 0 {
		json.Unmarshal(validationJSON, &attr.ValidationRules)
	}
	
	return &attr, nil
}

// ListAttributeDefinitions retrieves all attribute definitions
func (r *AttributeRepo) ListAttributeDefinitions(ctx context.Context) ([]models.AttributeDefinition, error) {
	query := `
		SELECT id,
		       COALESCE(name_en,'') AS name_en,
		       COALESCE(name_am,'') AS name_am,
		       COALESCE(name_om,'') AS name_om,
		       slug,
		       COALESCE(description,'') AS description,
		       data_type,
		       validation_rules,
		       COALESCE(unit,'') AS unit,
		       is_required, is_searchable, is_filterable,
		       display_order, created_at, updated_at
		FROM attribute_definitions
		ORDER BY display_order, name_en`
	
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var attrs []models.AttributeDefinition
	for rows.Next() {
		var attr models.AttributeDefinition
		var validationJSON []byte
		
		err := rows.Scan(
			&attr.ID, &attr.NameEn, &attr.NameAm, &attr.NameOm, &attr.Slug,
			&attr.Description, &attr.DataType, &validationJSON, &attr.Unit,
			&attr.IsRequired, &attr.IsSearchable, &attr.IsFilterable,
			&attr.DisplayOrder, &attr.CreatedAt, &attr.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		
		if len(validationJSON) > 0 {
			json.Unmarshal(validationJSON, &attr.ValidationRules)
		}
		
		attrs = append(attrs, attr)
	}
	
	return attrs, rows.Err()
}

// GetCategoryAttributes retrieves attributes for a category
func (r *AttributeRepo) GetCategoryAttributes(ctx context.Context, categoryID int) ([]models.ApplicableAttribute, error) {
	query := `
		SELECT attribute_id, name_en, slug, data_type,
		       validation_rules, COALESCE(unit, ''), is_required,
		       is_filterable, display_order
		FROM get_category_attributes($1)`
	
	rows, err := r.db.Query(ctx, query, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	attrs := []models.ApplicableAttribute{}
	for rows.Next() {
		var attr models.ApplicableAttribute
		var validationJSON []byte
		
		err := rows.Scan(
			&attr.AttributeID, &attr.NameEn, &attr.Slug, &attr.DataType,
			&validationJSON, &attr.Unit, &attr.IsRequired,
			&attr.IsFilterable, &attr.DisplayOrder,
		)
		if err != nil {
			return nil, err
		}
		
		if len(validationJSON) > 0 {
			json.Unmarshal(validationJSON, &attr.ValidationRules)
		}
		
		attrs = append(attrs, attr)
	}
	
	return attrs, rows.Err()
}

// GetProductAttributes retrieves all attributes for a product
func (r *AttributeRepo) GetProductAttributes(ctx context.Context, productID int) ([]models.ProductAttribute, error) {
	query := `
		SELECT pa.id, pa.product_id, pa.attribute_id,
		       pa.value_text, pa.value_number, pa.value_boolean,
		       pa.value_date, pa.value_json, pa.created_at, pa.updated_at,
		       COALESCE(ad.name_en,''), COALESCE(ad.name_am,''), COALESCE(ad.name_om,''),
		       ad.slug, ad.data_type, COALESCE(ad.unit,'')
		FROM product_attributes pa
		JOIN attribute_definitions ad ON ad.id = pa.attribute_id
		WHERE pa.product_id = $1
		ORDER BY ad.display_order`
	
	rows, err := r.db.Query(ctx, query, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	attrs := []models.ProductAttribute{}
	for rows.Next() {
		var attr models.ProductAttribute
		var valueJSON []byte
		attr.Attribute = &models.AttributeDefinition{}
		
		err := rows.Scan(
			&attr.ID, &attr.ProductID, &attr.AttributeID,
			&attr.ValueText, &attr.ValueNumber, &attr.ValueBool,
			&attr.ValueDate, &valueJSON, &attr.CreatedAt, &attr.UpdatedAt,
			&attr.Attribute.NameEn, &attr.Attribute.NameAm, &attr.Attribute.NameOm,
			&attr.Attribute.Slug, &attr.Attribute.DataType, &attr.Attribute.Unit,
		)
		if err != nil {
			return nil, err
		}
		
		if len(valueJSON) > 0 {
			json.Unmarshal(valueJSON, &attr.ValueJSON)
		}
		
		attrs = append(attrs, attr)
	}
	
	return attrs, rows.Err()
}

// SetProductAttribute sets or updates a product attribute value
func (r *AttributeRepo) SetProductAttribute(ctx context.Context, attr *models.ProductAttribute) error {
	valueJSON, _ := json.Marshal(attr.ValueJSON)
	
	query := `
		INSERT INTO product_attributes (
			product_id, attribute_id, value_text, value_number,
			value_boolean, value_date, value_json
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (product_id, attribute_id) DO UPDATE
		SET value_text = EXCLUDED.value_text,
		    value_number = EXCLUDED.value_number,
		    value_boolean = EXCLUDED.value_boolean,
		    value_date = EXCLUDED.value_date,
		    value_json = EXCLUDED.value_json,
		    updated_at = NOW()`
	
	_, err := r.db.Exec(ctx, query,
		attr.ProductID, attr.AttributeID, attr.ValueText, attr.ValueNumber,
		attr.ValueBool, attr.ValueDate, valueJSON,
	)
	
	return err
}

// DeleteProductAttributes removes all attributes for a product
func (r *AttributeRepo) DeleteProductAttributes(ctx context.Context, productID int) error {
	_, err := r.db.Exec(ctx, `DELETE FROM product_attributes WHERE product_id = $1`, productID)
	return err
}

// CreateAttributeDefinition creates a new attribute definition (admin)
func (r *AttributeRepo) CreateAttributeDefinition(ctx context.Context, attr *models.AttributeDefinition) (int, error) {
	validationJSON, _ := json.Marshal(attr.ValidationRules)
	
	query := `
		INSERT INTO attribute_definitions (
			name_en, name_am, name_om, slug, description, data_type,
			validation_rules, unit, is_required, is_searchable, is_filterable, display_order
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id`
	
	var id int
	err := r.db.QueryRow(ctx, query,
		attr.NameEn, attr.NameAm, attr.NameOm, attr.Slug, attr.Description,
		attr.DataType, validationJSON, attr.Unit, attr.IsRequired,
		attr.IsSearchable, attr.IsFilterable, attr.DisplayOrder,
	).Scan(&id)
	
	return id, err
}

// AssignAttributeToCategory assigns an attribute to a category
func (r *AttributeRepo) AssignAttributeToCategory(ctx context.Context, ca *models.CategoryAttribute) error {
	query := `
		INSERT INTO category_attributes (
			category_id, attribute_id, is_required, is_inherited, display_order
		) VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (category_id, attribute_id) DO UPDATE
		SET is_required = EXCLUDED.is_required,
		    display_order = EXCLUDED.display_order`
	
	_, err := r.db.Exec(ctx, query,
		ca.CategoryID, ca.AttributeID, ca.IsRequired, ca.IsInherited, ca.DisplayOrder,
	)
	
	return err
}

// RemoveAttributeFromCategory removes an attribute from a category
func (r *AttributeRepo) RemoveAttributeFromCategory(ctx context.Context, categoryID, attributeID int) error {
	tag, err := r.db.Exec(ctx,
		`DELETE FROM category_attributes WHERE category_id = $1 AND attribute_id = $2`,
		categoryID, attributeID,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("attribute mapping not found for category")
	}
	return nil
}
