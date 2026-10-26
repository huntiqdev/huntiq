package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/your-org/used-tools-marketplace/internal/models"
)

type CategoryHierarchyRepo struct {
	db *pgxpool.Pool
}

func NewCategoryHierarchyRepo(db *pgxpool.Pool) *CategoryHierarchyRepo {
	return &CategoryHierarchyRepo{db: db}
}

// GetByID retrieves a category with its breadcrumb path
func (r *CategoryHierarchyRepo) GetByID(ctx context.Context, id int, lang string) (*models.CategoryHierarchy, error) {
	query := `
		SELECT id, parent_id,
		       COALESCE(name_en, ''), COALESCE(name_am, ''), COALESCE(name_om, ''), COALESCE(name_ti, ''),
		       COALESCE(slug, ''), COALESCE(description, ''),
		       level, COALESCE(path, ''), display_order,
		       COALESCE(icon, ''), COALESCE(image_url, ''), is_active, is_leaf,
		       COALESCE(meta_title, ''), COALESCE(meta_description, ''),
		       created_at, updated_at
		FROM category_hierarchy
		WHERE id = $1`

	var cat models.CategoryHierarchy
	err := r.db.QueryRow(ctx, query, id).Scan(
		&cat.ID, &cat.ParentID, &cat.NameEn, &cat.NameAm, &cat.NameOm, &cat.NameTi,
		&cat.Slug, &cat.Description, &cat.Level, &cat.Path, &cat.DisplayOrder,
		&cat.Icon, &cat.ImageURL, &cat.IsActive, &cat.IsLeaf, &cat.MetaTitle,
		&cat.MetaDescription, &cat.CreatedAt, &cat.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	breadcrumbs, err := r.GetBreadcrumbs(ctx, id, lang)
	if err == nil {
		cat.Breadcrumbs = breadcrumbs
	}

	return &cat, nil
}

// GetBySlug retrieves a category by its slug
func (r *CategoryHierarchyRepo) GetBySlug(ctx context.Context, slug string, lang string) (*models.CategoryHierarchy, error) {
	query := `
		SELECT id, parent_id,
		       COALESCE(name_en, ''), COALESCE(name_am, ''), COALESCE(name_om, ''), COALESCE(name_ti, ''),
		       COALESCE(slug, ''), COALESCE(description, ''),
		       level, COALESCE(path, ''), display_order,
		       COALESCE(icon, ''), COALESCE(image_url, ''), is_active, is_leaf,
		       COALESCE(meta_title, ''), COALESCE(meta_description, ''),
		       created_at, updated_at
		FROM category_hierarchy
		WHERE slug = $1`

	var cat models.CategoryHierarchy
	err := r.db.QueryRow(ctx, query, slug).Scan(
		&cat.ID, &cat.ParentID, &cat.NameEn, &cat.NameAm, &cat.NameOm, &cat.NameTi,
		&cat.Slug, &cat.Description, &cat.Level, &cat.Path, &cat.DisplayOrder,
		&cat.Icon, &cat.ImageURL, &cat.IsActive, &cat.IsLeaf, &cat.MetaTitle,
		&cat.MetaDescription, &cat.CreatedAt, &cat.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	breadcrumbs, err := r.GetBreadcrumbs(ctx, cat.ID, lang)
	if err == nil {
		cat.Breadcrumbs = breadcrumbs
	}

	return &cat, nil
}

// GetRoots returns all active root-level categories
func (r *CategoryHierarchyRepo) GetRoots(ctx context.Context, lang string) ([]models.CategoryHierarchy, error) {
	query := `
		SELECT id, parent_id,
		       COALESCE(name_en, ''), COALESCE(name_am, ''), COALESCE(name_om, ''), COALESCE(name_ti, ''),
		       COALESCE(slug, ''), COALESCE(description, ''),
		       level, COALESCE(path, ''), display_order,
		       COALESCE(icon, ''), COALESCE(image_url, ''), is_active, is_leaf,
		       COALESCE(meta_title, ''), COALESCE(meta_description, ''),
		       created_at, updated_at
		FROM category_hierarchy
		WHERE parent_id IS NULL AND is_active = true
		ORDER BY display_order ASC, name_en ASC`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.CategoryHierarchy
	for rows.Next() {
		var cat models.CategoryHierarchy
		err := rows.Scan(
			&cat.ID, &cat.ParentID, &cat.NameEn, &cat.NameAm, &cat.NameOm, &cat.NameTi,
			&cat.Slug, &cat.Description, &cat.Level, &cat.Path, &cat.DisplayOrder,
			&cat.Icon, &cat.ImageURL, &cat.IsActive, &cat.IsLeaf, &cat.MetaTitle,
			&cat.MetaDescription, &cat.CreatedAt, &cat.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		categories = append(categories, cat)
	}

	return categories, rows.Err()
}

// ListCategories retrieves categories with optional filtering
func (r *CategoryHierarchyRepo) ListCategories(ctx context.Context, f models.CategoryFilter) ([]models.CategoryHierarchy, error) {
	var conditions []string
	args := []interface{}{}
	argIdx := 1

	if f.ParentID != nil {
		conditions = append(conditions, fmt.Sprintf("parent_id = $%d", argIdx))
		args = append(args, *f.ParentID)
		argIdx++
	}
	if f.Level != nil {
		conditions = append(conditions, fmt.Sprintf("level = $%d", argIdx))
		args = append(args, *f.Level)
		argIdx++
	}
	if f.IsActive != nil {
		conditions = append(conditions, fmt.Sprintf("is_active = $%d", argIdx))
		args = append(args, *f.IsActive)
		argIdx++
	}
	if f.IsLeaf != nil {
		conditions = append(conditions, fmt.Sprintf("is_leaf = $%d", argIdx))
		args = append(args, *f.IsLeaf)
		argIdx++
	}
	if f.SearchTerm != "" {
		conditions = append(conditions, fmt.Sprintf("(name_en ILIKE $%d OR name_am ILIKE $%d OR name_om ILIKE $%d OR name_ti ILIKE $%d OR slug ILIKE $%d)", argIdx, argIdx, argIdx, argIdx, argIdx))
		args = append(args, "%"+f.SearchTerm+"%")
		argIdx++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	limit := f.Limit
	if limit <= 0 || limit > 500 {
		limit = 100
	}
	page := f.Page
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	query := fmt.Sprintf(`
		SELECT id, parent_id,
		       COALESCE(name_en, ''), COALESCE(name_am, ''), COALESCE(name_om, ''), COALESCE(name_ti, ''),
		       COALESCE(slug, ''), COALESCE(description, ''),
		       level, COALESCE(path, ''), display_order,
		       COALESCE(icon, ''), COALESCE(image_url, ''), is_active, is_leaf,
		       COALESCE(meta_title, ''), COALESCE(meta_description, ''),
		       created_at, updated_at
		FROM category_hierarchy
		%s
		ORDER BY display_order ASC, name_en ASC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIdx, argIdx+1)

	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.CategoryHierarchy
	for rows.Next() {
		var cat models.CategoryHierarchy
		err := rows.Scan(
			&cat.ID, &cat.ParentID, &cat.NameEn, &cat.NameAm, &cat.NameOm, &cat.NameTi,
			&cat.Slug, &cat.Description, &cat.Level, &cat.Path, &cat.DisplayOrder,
			&cat.Icon, &cat.ImageURL, &cat.IsActive, &cat.IsLeaf, &cat.MetaTitle,
			&cat.MetaDescription, &cat.CreatedAt, &cat.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		categories = append(categories, cat)
	}

	return categories, rows.Err()
}

// GetCategoryTree builds a full hierarchical tree
func (r *CategoryHierarchyRepo) GetCategoryTree(ctx context.Context, lang string) ([]models.CategoryTree, error) {
	query := `
		WITH RECURSIVE category_tree AS (
			-- Root categories
			SELECT c.id, c.parent_id,
			       COALESCE(c.name_en,'') AS name_en,
			       COALESCE(c.name_am,'') AS name_am,
			       COALESCE(c.name_om,'') AS name_om,
			       COALESCE(c.name_ti,'') AS name_ti,
			       COALESCE(c.slug,'') AS slug,
			       c.level, c.is_leaf,
			       COALESCE(c.icon, '') AS icon,
			       COALESCE(c.image_url, '') AS image_url,
			       COALESCE(s.product_count, 0) AS product_count,
			       c.display_order
			FROM category_hierarchy c
			LEFT JOIN category_stats s ON s.category_id = c.id
			WHERE c.parent_id IS NULL AND c.is_active = true
			
			UNION ALL
			
			-- Child categories
			SELECT c.id, c.parent_id,
			       COALESCE(c.name_en,'') AS name_en,
			       COALESCE(c.name_am,'') AS name_am,
			       COALESCE(c.name_om,'') AS name_om,
			       COALESCE(c.name_ti,'') AS name_ti,
			       COALESCE(c.slug,'') AS slug,
			       c.level, c.is_leaf,
			       COALESCE(c.icon, '') AS icon,
			       COALESCE(c.image_url, '') AS image_url,
			       COALESCE(s.product_count, 0) AS product_count,
			       c.display_order
			FROM category_hierarchy c
			INNER JOIN category_tree ct ON c.parent_id = ct.id
			LEFT JOIN category_stats s ON s.category_id = c.id
			WHERE c.is_active = true
		)
		SELECT id, parent_id, name_en, name_am, name_om, name_ti, slug, level, is_leaf, icon, image_url, product_count
		FROM category_tree
		ORDER BY level, display_order, id`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	type flatCat struct {
		ID           int
		ParentID     *int
		NameEn       string
		NameAm       string
		NameOm       string
		NameTi       string
		Slug         string
		Level        int
		IsLeaf       bool
		Icon         string
		ImageURL     string
		ProductCount int
	}

	var flat []flatCat
	for rows.Next() {
		var fc flatCat
		err := rows.Scan(&fc.ID, &fc.ParentID, &fc.NameEn, &fc.NameAm, &fc.NameOm, &fc.NameTi,
			&fc.Slug, &fc.Level, &fc.IsLeaf, &fc.Icon, &fc.ImageURL, &fc.ProductCount)
		if err != nil {
			return nil, err
		}
		flat = append(flat, fc)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	catMap := make(map[int]*models.CategoryTree)

	for _, fc := range flat {
		name := fc.NameEn
		switch lang {
		case "am":
			if fc.NameAm != "" {
				name = fc.NameAm
			}
		case "om":
			if fc.NameOm != "" {
				name = fc.NameOm
			}
		case "ti":
			if fc.NameTi != "" {
				name = fc.NameTi
			}
		}

		node := &models.CategoryTree{
			ID:       fc.ID,
			Name:     name,
			NameEn:   fc.NameEn,
			NameAm:   fc.NameAm,
			NameOm:   fc.NameOm,
			NameTi:   fc.NameTi,
			Slug:     fc.Slug,
			Level:    fc.Level,
			IsLeaf:   fc.IsLeaf,
			Icon:     fc.Icon,
			ImageURL: fc.ImageURL,
			Count:    fc.ProductCount,
			Children: []models.CategoryTree{},
		}
		catMap[fc.ID] = node
	}

	var roots []models.CategoryTree
	for _, fc := range flat {
		node := catMap[fc.ID]
		if node == nil {
			continue
		}
		if fc.ParentID == nil {
			roots = append(roots, *node)
		} else {
			if parent := catMap[*fc.ParentID]; parent != nil {
				parent.Children = append(parent.Children, *node)
			}
		}
	}

	for i, r := range roots {
		if updated, ok := catMap[r.ID]; ok {
			roots[i] = *updated
		}
	}

	return roots, nil
}

// GetBreadcrumbs returns the full path from root to target category
func (r *CategoryHierarchyRepo) GetBreadcrumbs(ctx context.Context, categoryID int, lang string) ([]models.CategoryBreadcrumb, error) {
	nameCol := "a.name_en"
	switch lang {
	case "am":
		nameCol = "COALESCE(NULLIF(ch.name_am, ''), a.name_en)"
	case "om":
		nameCol = "COALESCE(NULLIF(ch.name_om, ''), a.name_en)"
	case "ti":
		nameCol = "COALESCE(NULLIF(ch.name_ti, ''), a.name_en)"
	}

	query := fmt.Sprintf(`
		SELECT a.id, %s AS name, a.slug, a.level
		FROM get_category_ancestors($1) a
		JOIN category_hierarchy ch ON ch.id = a.id
		ORDER BY a.level ASC
	`, nameCol)

	rows, err := r.db.Query(ctx, query, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	breadcrumbs := []models.CategoryBreadcrumb{}
	for rows.Next() {
		var bc models.CategoryBreadcrumb
		err := rows.Scan(&bc.ID, &bc.Name, &bc.Slug, &bc.Level)
		if err != nil {
			return nil, err
		}
		breadcrumbs = append(breadcrumbs, bc)
	}

	return breadcrumbs, rows.Err()
}

// GetAncestors returns all parent category IDs
func (r *CategoryHierarchyRepo) GetAncestors(ctx context.Context, categoryID int) ([]int, error) {
	query := `SELECT ARRAY(SELECT id FROM get_category_ancestors($1))`
	var ancestorIDs []int
	err := r.db.QueryRow(ctx, query, categoryID).Scan(&ancestorIDs)
	return ancestorIDs, err
}

// GetRootAncestor returns the level 0 root ancestor ID and slug for a given category ID
func (r *CategoryHierarchyRepo) GetRootAncestor(ctx context.Context, categoryID int) (int, string, error) {
	query := `SELECT id, slug FROM get_category_ancestors($1) WHERE level = 0 LIMIT 1`
	var id int
	var slug string
	err := r.db.QueryRow(ctx, query, categoryID).Scan(&id, &slug)
	return id, slug, err
}

// GetDescendants returns all child category IDs (recursive)
func (r *CategoryHierarchyRepo) GetDescendants(ctx context.Context, categoryID int) ([]int, error) {
	query := `SELECT ARRAY(SELECT id FROM get_category_descendants($1))`
	var descendantIDs []int
	err := r.db.QueryRow(ctx, query, categoryID).Scan(&descendantIDs)
	return descendantIDs, err
}

// CheckSlugExists checks if a slug is already taken by another category
func (r *CategoryHierarchyRepo) CheckSlugExists(ctx context.Context, slug string, excludeID int) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM category_hierarchy WHERE slug = $1 AND id != $2)`
	var exists bool
	err := r.db.QueryRow(ctx, query, slug, excludeID).Scan(&exists)
	return exists, err
}

// ValidateCategory checks if a category exists and whether it is active
func (r *CategoryHierarchyRepo) ValidateCategory(ctx context.Context, categoryID int) (exists bool, isActive bool, err error) {
	query := `SELECT is_active FROM category_hierarchy WHERE id = $1`
	err = r.db.QueryRow(ctx, query, categoryID).Scan(&isActive)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return false, false, nil
		}
		return false, false, err
	}
	return true, isActive, nil
}

// IsDescendantOf checks if childID is a descendant of ancestorID
func (r *CategoryHierarchyRepo) IsDescendantOf(ctx context.Context, childID int, ancestorID int) (bool, error) {
	if childID == ancestorID {
		return true, nil
	}
	query := `SELECT EXISTS (
		SELECT 1 FROM get_category_descendants($1) WHERE id = $2
	)`
	var isDesc bool
	err := r.db.QueryRow(ctx, query, ancestorID, childID).Scan(&isDesc)
	return isDesc, err
}

// CreateCategory inserts a new category
func (r *CategoryHierarchyRepo) CreateCategory(ctx context.Context, cat *models.CategoryHierarchy) (int, error) {
	query := `
		INSERT INTO category_hierarchy (
			parent_id, name_en, name_am, name_om, name_ti, slug, description,
			display_order, icon, image_url, is_leaf, is_active, meta_title, meta_description
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING id`

	var id int
	err := r.db.QueryRow(ctx, query,
		cat.ParentID, cat.NameEn, cat.NameAm, cat.NameOm, cat.NameTi, cat.Slug,
		cat.Description, cat.DisplayOrder, cat.Icon, cat.ImageURL, cat.IsLeaf,
		cat.IsActive, cat.MetaTitle, cat.MetaDescription,
	).Scan(&id)

	return id, err
}

// UpdateCategory modifies an existing category
func (r *CategoryHierarchyRepo) UpdateCategory(ctx context.Context, cat *models.CategoryHierarchy) error {
	query := `
		UPDATE category_hierarchy
		SET parent_id = $2, name_en = $3, name_am = $4, name_om = $5, name_ti = $6,
		    slug = $7, description = $8, display_order = $9, icon = $10,
		    image_url = $11, is_leaf = $12, meta_title = $13, meta_description = $14,
		    is_active = $15, updated_at = NOW()
		WHERE id = $1`

	tag, err := r.db.Exec(ctx, query,
		cat.ID, cat.ParentID, cat.NameEn, cat.NameAm, cat.NameOm, cat.NameTi,
		cat.Slug, cat.Description, cat.DisplayOrder, cat.Icon,
		cat.ImageURL, cat.IsLeaf, cat.MetaTitle, cat.MetaDescription, cat.IsActive,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("category not found")
	}
	return nil
}

// ReorderCategory updates the display order of a category
func (r *CategoryHierarchyRepo) ReorderCategory(ctx context.Context, id int, displayOrder int) error {
	query := `UPDATE category_hierarchy SET display_order = $2, updated_at = NOW() WHERE id = $1`
	tag, err := r.db.Exec(ctx, query, id, displayOrder)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("category not found")
	}
	return nil
}

// SetCategoryStatus toggles active status of a category
func (r *CategoryHierarchyRepo) SetCategoryStatus(ctx context.Context, id int, isActive bool) error {
	query := `UPDATE category_hierarchy SET is_active = $2, updated_at = NOW() WHERE id = $1`
	tag, err := r.db.Exec(ctx, query, id, isActive)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("category not found")
	}
	return nil
}

// DeleteCategory removes a category (only if no children or products; otherwise suggests deactivation)
func (r *CategoryHierarchyRepo) DeleteCategory(ctx context.Context, id int) error {
	var childCount int
	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM category_hierarchy WHERE parent_id = $1`, id).Scan(&childCount)
	if err != nil {
		return err
	}
	if childCount > 0 {
		return fmt.Errorf("cannot delete category with %d active child subcategories; deactivate it instead", childCount)
	}

	var productCount int
	err = r.db.QueryRow(ctx, `SELECT COUNT(*) FROM products WHERE category_hierarchy_id = $1`, id).Scan(&productCount)
	if err != nil {
		return err
	}
	if productCount > 0 {
		return fmt.Errorf("cannot delete category referenced by %d products; deactivate it instead to preserve listing integrity", productCount)
	}

	tag, err := r.db.Exec(ctx, `DELETE FROM category_hierarchy WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("category not found")
	}
	return nil
}

// GetCategoryStats retrieves stats for a category
func (r *CategoryHierarchyRepo) GetCategoryStats(ctx context.Context, categoryID int) (*models.CategoryWithStats, error) {
	query := `
		SELECT 
			c.id, c.parent_id,
			COALESCE(c.name_en, ''), COALESCE(c.name_am, ''), COALESCE(c.name_om, ''), COALESCE(c.name_ti, ''),
			COALESCE(c.slug, ''), COALESCE(c.description, ''),
			c.level, COALESCE(c.path, ''), c.display_order,
			COALESCE(c.icon, ''), COALESCE(c.image_url, ''), c.is_active, c.is_leaf,
			COALESCE(c.meta_title, ''), COALESCE(meta_description, ''),
			c.created_at, c.updated_at,
			COALESCE(s.product_count, 0) AS product_count,
			COALESCE(s.approved_count, 0) AS approved_count,
			COALESCE(s.pending_count, 0) AS pending_count
		FROM category_hierarchy c
		LEFT JOIN category_stats s ON s.category_id = c.id
		WHERE c.id = $1`

	var cat models.CategoryWithStats
	err := r.db.QueryRow(ctx, query, categoryID).Scan(
		&cat.ID, &cat.ParentID, &cat.NameEn, &cat.NameAm, &cat.NameOm, &cat.NameTi,
		&cat.Slug, &cat.Description, &cat.Level, &cat.Path, &cat.DisplayOrder,
		&cat.Icon, &cat.ImageURL, &cat.IsActive, &cat.IsLeaf, &cat.MetaTitle,
		&cat.MetaDescription, &cat.CreatedAt, &cat.UpdatedAt,
		&cat.ProductCount, &cat.ApprovedCount, &cat.PendingCount,
	)

	return &cat, err
}
