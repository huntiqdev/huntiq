package repository

import (
	"context"
	"fmt"
	"math"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/your-org/used-tools-marketplace/internal/models"
)

type LocationHierarchyRepo struct {
	db *pgxpool.Pool
}

func NewLocationHierarchyRepo(db *pgxpool.Pool) *LocationHierarchyRepo {
	return &LocationHierarchyRepo{db: db}
}

// GetByID retrieves a single location by ID with localized breadcrumbs
func (r *LocationHierarchyRepo) GetByID(ctx context.Context, id int, lang string) (*models.LocationHierarchy, error) {
	query := `
		SELECT id, parent_id, location_type,
		       COALESCE(name_en, ''), COALESCE(name_am, ''), COALESCE(name_om, ''), COALESCE(name_ti, ''),
		       COALESCE(slug, ''), COALESCE(admin_code, ''), level, COALESCE(path, ''), display_order,
		       is_active, is_leaf, is_deprecated, replaced_by_id,
		       latitude, longitude, min_lat, max_lat, min_lng, max_lng,
		       COALESCE(aliases, '{}'), COALESCE(source, ''), COALESCE(verification_status, ''),
		       created_at, updated_at
		FROM location_hierarchy
		WHERE id = $1`

	var loc models.LocationHierarchy
	err := r.db.QueryRow(ctx, query, id).Scan(
		&loc.ID, &loc.ParentID, &loc.LocationType,
		&loc.NameEn, &loc.NameAm, &loc.NameOm, &loc.NameTi,
		&loc.Slug, &loc.AdminCode, &loc.Level, &loc.Path, &loc.DisplayOrder,
		&loc.IsActive, &loc.IsLeaf, &loc.IsDeprecated, &loc.ReplacedByID,
		&loc.Latitude, &loc.Longitude, &loc.MinLat, &loc.MaxLat, &loc.MinLng, &loc.MaxLng,
		&loc.Aliases, &loc.Source, &loc.VerificationStatus,
		&loc.CreatedAt, &loc.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	breadcrumbs, err := r.GetBreadcrumbs(ctx, id, lang)
	if err == nil {
		loc.Breadcrumbs = breadcrumbs
	}

	return &loc, nil
}

// GetBreadcrumbs returns the full hierarchy path from country root to this location
func (r *LocationHierarchyRepo) GetBreadcrumbs(ctx context.Context, locationID int, lang string) ([]models.LocationBreadcrumb, error) {
	nameCol := "name_en"
	switch lang {
	case "am":
		nameCol = "COALESCE(NULLIF(name_am, ''), name_en)"
	case "om":
		nameCol = "COALESCE(NULLIF(name_om, ''), name_en)"
	case "ti":
		nameCol = "COALESCE(NULLIF(name_ti, ''), name_en)"
	}

	query := fmt.Sprintf(`
		SELECT id, %s AS name, slug, level, location_type
		FROM get_location_ancestors($1)
		ORDER BY level ASC`, nameCol)

	rows, err := r.db.Query(ctx, query, locationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var breadcrumbs []models.LocationBreadcrumb
	for rows.Next() {
		var bc models.LocationBreadcrumb
		if err := rows.Scan(&bc.ID, &bc.Name, &bc.Slug, &bc.Level, &bc.LocationType); err != nil {
			return nil, err
		}
		breadcrumbs = append(breadcrumbs, bc)
	}

	return breadcrumbs, rows.Err()
}

// ListChildren returns direct child locations under a given parent ID (or roots if parentID is nil)
func (r *LocationHierarchyRepo) ListChildren(ctx context.Context, parentID *int, lang string) ([]models.LocationHierarchy, error) {
	var query string
	var rows interface {
		Next() bool
		Scan(...interface{}) error
		Err() error
		Close()
	}
	var err error

	baseQuery := `
		SELECT id, parent_id, location_type,
		       COALESCE(name_en, ''), COALESCE(name_am, ''), COALESCE(name_om, ''), COALESCE(name_ti, ''),
		       COALESCE(slug, ''), COALESCE(admin_code, ''), level, COALESCE(path, ''), display_order,
		       is_active, is_leaf, is_deprecated, replaced_by_id,
		       latitude, longitude, min_lat, max_lat, min_lng, max_lng,
		       COALESCE(aliases, '{}'), COALESCE(source, ''), COALESCE(verification_status, ''),
		       created_at, updated_at
		FROM location_hierarchy`

	if parentID == nil {
		query = baseQuery + ` WHERE parent_id IS NULL AND is_active = true ORDER BY display_order ASC, name_en ASC`
		rows, err = r.db.Query(ctx, query)
	} else {
		query = baseQuery + ` WHERE parent_id = $1 AND is_active = true ORDER BY display_order ASC, name_en ASC`
		rows, err = r.db.Query(ctx, query, *parentID)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.LocationHierarchy
	for rows.Next() {
		var loc models.LocationHierarchy
		if err := rows.Scan(
			&loc.ID, &loc.ParentID, &loc.LocationType,
			&loc.NameEn, &loc.NameAm, &loc.NameOm, &loc.NameTi,
			&loc.Slug, &loc.AdminCode, &loc.Level, &loc.Path, &loc.DisplayOrder,
			&loc.IsActive, &loc.IsLeaf, &loc.IsDeprecated, &loc.ReplacedByID,
			&loc.Latitude, &loc.Longitude, &loc.MinLat, &loc.MaxLat, &loc.MinLng, &loc.MaxLng,
			&loc.Aliases, &loc.Source, &loc.VerificationStatus,
			&loc.CreatedAt, &loc.UpdatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, loc)
	}

	return list, rows.Err()
}

// ListTree returns the full hierarchical location tree
func (r *LocationHierarchyRepo) ListTree(ctx context.Context, lang string) ([]models.LocationTree, error) {
	nameCol := "name_en"
	switch lang {
	case "am":
		nameCol = "COALESCE(NULLIF(name_am, ''), name_en)"
	case "om":
		nameCol = "COALESCE(NULLIF(name_om, ''), name_en)"
	case "ti":
		nameCol = "COALESCE(NULLIF(name_ti, ''), name_en)"
	}

	query := fmt.Sprintf(`
		SELECT id, parent_id, %s AS name, slug, location_type, COALESCE(admin_code, ''),
		       level, latitude, longitude, is_leaf
		FROM location_hierarchy
		WHERE is_active = true
		ORDER BY level ASC, display_order ASC, name_en ASC`, nameCol)

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	type flatNode struct {
		ID           int
		ParentID     *int
		Name         string
		Slug         string
		LocationType string
		AdminCode    string
		Level        int
		Latitude     *float64
		Longitude    *float64
		IsLeaf       bool
	}

	var flat []flatNode
	for rows.Next() {
		var fn flatNode
		if err := rows.Scan(
			&fn.ID, &fn.ParentID, &fn.Name, &fn.Slug, &fn.LocationType, &fn.AdminCode,
			&fn.Level, &fn.Latitude, &fn.Longitude, &fn.IsLeaf,
		); err != nil {
			return nil, err
		}
		flat = append(flat, fn)
	}

	nodeMap := make(map[int]*models.LocationTree)
	for _, fn := range flat {
		nodeMap[fn.ID] = &models.LocationTree{
			ID:           fn.ID,
			Name:         fn.Name,
			Slug:         fn.Slug,
			LocationType: fn.LocationType,
			AdminCode:    fn.AdminCode,
			Level:        fn.Level,
			Latitude:     fn.Latitude,
			Longitude:    fn.Longitude,
			IsLeaf:       fn.IsLeaf,
			Children:     []models.LocationTree{},
		}
	}

	// Link children to parents bottom-up (from deepest level to root)
	for i := len(flat) - 1; i >= 0; i-- {
		fn := flat[i]
		if fn.ParentID != nil {
			node := nodeMap[fn.ID]
			if parent, ok := nodeMap[*fn.ParentID]; ok {
				parent.Children = append(parent.Children, *node)
			}
		}
	}

	var roots []models.LocationTree
	for _, fn := range flat {
		if fn.ParentID == nil {
			if node, ok := nodeMap[fn.ID]; ok {
				roots = append(roots, *node)
			}
		}
	}

	return roots, nil
}

// MatchByCoordinates resolves a GPS coordinate to the most specific location node.
// Returns the matched location and a confidence score between 0.0 and 1.0.
func (r *LocationHierarchyRepo) MatchByCoordinates(ctx context.Context, lat, lng float64) (*models.LocationHierarchy, float64, error) {
	// Step 1: Bounding box containment match (highest confidence, prioritized by deepest level)
	bboxQuery := `
		SELECT id
		FROM location_hierarchy
		WHERE is_active = true
		  AND min_lat IS NOT NULL AND max_lat IS NOT NULL
		  AND min_lng IS NOT NULL AND max_lng IS NOT NULL
		  AND $1 >= min_lat AND $1 <= max_lat
		  AND $2 >= min_lng AND $2 <= max_lng
		ORDER BY level DESC, display_order ASC
		LIMIT 1`

	var matchedID int
	err := r.db.QueryRow(ctx, bboxQuery, lat, lng).Scan(&matchedID)
	if err == nil && matchedID > 0 {
		loc, getErr := r.GetByID(ctx, matchedID, "en")
		if getErr == nil {
			return loc, 0.95, nil
		}
	}

	// Step 2: Nearest centroid distance match (if within sensible tolerance)
	// Euclidean approximation in degrees: 0.1 degree ~ 11 km
	centroidQuery := `
		SELECT id,
		       SQRT(POW(latitude - $1, 2) + POW(longitude - $2, 2)) AS deg_dist
		FROM location_hierarchy
		WHERE is_active = true
		  AND latitude IS NOT NULL AND longitude IS NOT NULL
		  AND level >= 1 -- match at least a region, city, or subcity
		ORDER BY deg_dist ASC
		LIMIT 1`

	var degDist float64
	err = r.db.QueryRow(ctx, centroidQuery, lat, lng).Scan(&matchedID, &degDist)
	if err != nil {
		return nil, 0, fmt.Errorf("no location matches coordinates: %w", err)
	}

	// Convert degree distance roughly to kilometers (~111km per degree)
	kmDist := degDist * 111.0

	// Confidence scales down as distance increases
	confidence := 0.90
	if kmDist > 25.0 {
		confidence = math.Max(0.50, 0.90-(kmDist-25.0)*0.01)
	}

	loc, err := r.GetByID(ctx, matchedID, "en")
	if err != nil {
		return nil, 0, err
	}

	return loc, confidence, nil
}

// SearchLocations performs fuzzy trigram search across all 4 localized name columns and aliases
func (r *LocationHierarchyRepo) SearchLocations(ctx context.Context, q string, lang string, limit int) ([]models.LocationHierarchy, error) {
	if limit <= 0 || limit > 50 {
		limit = 15
	}
	cleanQ := strings.TrimSpace(q)
	if cleanQ == "" {
		return nil, nil
	}
	likePattern := "%" + escapeLike(cleanQ) + "%"

	query := `
		SELECT id, parent_id, location_type,
		       COALESCE(name_en, ''), COALESCE(name_am, ''), COALESCE(name_om, ''), COALESCE(name_ti, ''),
		       COALESCE(slug, ''), COALESCE(admin_code, ''), level, COALESCE(path, ''), display_order,
		       is_active, is_leaf, is_deprecated, replaced_by_id,
		       latitude, longitude, min_lat, max_lat, min_lng, max_lng,
		       COALESCE(aliases, '{}'), COALESCE(source, ''), COALESCE(verification_status, ''),
		       created_at, updated_at,
		       GREATEST(
		           similarity(LOWER(name_en), LOWER($1)),
		           similarity(LOWER(COALESCE(name_am, '')), LOWER($1)),
		           similarity(LOWER(COALESCE(name_om, '')), LOWER($1)),
		           similarity(LOWER(COALESCE(name_ti, '')), LOWER($1))
		       ) AS score
		FROM location_hierarchy
		WHERE is_active = true
		  AND (
		      name_en ILIKE $2 OR name_am ILIKE $2 OR name_om ILIKE $2 OR name_ti ILIKE $2
		      OR admin_code ILIKE $2
		      OR array_to_string(aliases, ' ') ILIKE $2
		      OR similarity(LOWER(name_en), LOWER($1)) > 0.2
		      OR similarity(LOWER(COALESCE(name_am, '')), LOWER($1)) > 0.2
		  )
		ORDER BY score DESC, level DESC
		LIMIT $3`

	rows, err := r.db.Query(ctx, query, cleanQ, likePattern, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.LocationHierarchy
	for rows.Next() {
		var loc models.LocationHierarchy
		var score float64
		if err := rows.Scan(
			&loc.ID, &loc.ParentID, &loc.LocationType,
			&loc.NameEn, &loc.NameAm, &loc.NameOm, &loc.NameTi,
			&loc.Slug, &loc.AdminCode, &loc.Level, &loc.Path, &loc.DisplayOrder,
			&loc.IsActive, &loc.IsLeaf, &loc.IsDeprecated, &loc.ReplacedByID,
			&loc.Latitude, &loc.Longitude, &loc.MinLat, &loc.MaxLat, &loc.MinLng, &loc.MaxLng,
			&loc.Aliases, &loc.Source, &loc.VerificationStatus,
			&loc.CreatedAt, &loc.UpdatedAt, &score,
		); err != nil {
			return nil, err
		}
		breadcrumbs, _ := r.GetBreadcrumbs(ctx, loc.ID, lang)
		loc.Breadcrumbs = breadcrumbs
		results = append(results, loc)
	}

	return results, rows.Err()
}

// ListLocations retrieves paginated location records with filtering for the admin management console
func (r *LocationHierarchyRepo) ListLocations(ctx context.Context, f models.LocationFilter) ([]models.LocationHierarchy, int, error) {
	var conditions []string
	args := []interface{}{}
	argIdx := 1

	if f.ParentID != nil {
		conditions = append(conditions, fmt.Sprintf("parent_id = $%d", argIdx))
		args = append(args, *f.ParentID)
		argIdx++
	}
	if f.LocationType != "" {
		conditions = append(conditions, fmt.Sprintf("location_type = $%d", argIdx))
		args = append(args, f.LocationType)
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
	if f.SearchTerm != "" {
		conditions = append(conditions, fmt.Sprintf("(name_en ILIKE $%d OR name_am ILIKE $%d OR name_om ILIKE $%d OR name_ti ILIKE $%d OR admin_code ILIKE $%d)", argIdx, argIdx, argIdx, argIdx, argIdx))
		args = append(args, "%"+escapeLike(f.SearchTerm)+"%")
		argIdx++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM location_hierarchy %s", whereClause)
	var total int
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	if f.Limit <= 0 {
		f.Limit = 50
	}
	if f.Page <= 0 {
		f.Page = 1
	}
	offset := (f.Page - 1) * f.Limit

	dataQuery := fmt.Sprintf(`
		SELECT id, parent_id, location_type,
		       COALESCE(name_en, ''), COALESCE(name_am, ''), COALESCE(name_om, ''), COALESCE(name_ti, ''),
		       COALESCE(slug, ''), COALESCE(admin_code, ''), level, COALESCE(path, ''), display_order,
		       is_active, is_leaf, is_deprecated, replaced_by_id,
		       latitude, longitude, min_lat, max_lat, min_lng, max_lng,
		       COALESCE(aliases, '{}'), COALESCE(source, ''), COALESCE(verification_status, ''),
		       created_at, updated_at
		FROM location_hierarchy
		%s
		ORDER BY level ASC, display_order ASC, name_en ASC
		LIMIT $%d OFFSET $%d`, whereClause, argIdx, argIdx+1)

	args = append(args, f.Limit, offset)
	rows, err := r.db.Query(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var locations []models.LocationHierarchy
	for rows.Next() {
		var loc models.LocationHierarchy
		if err := rows.Scan(
			&loc.ID, &loc.ParentID, &loc.LocationType,
			&loc.NameEn, &loc.NameAm, &loc.NameOm, &loc.NameTi,
			&loc.Slug, &loc.AdminCode, &loc.Level, &loc.Path, &loc.DisplayOrder,
			&loc.IsActive, &loc.IsLeaf, &loc.IsDeprecated, &loc.ReplacedByID,
			&loc.Latitude, &loc.Longitude, &loc.MinLat, &loc.MaxLat, &loc.MinLng, &loc.MaxLng,
			&loc.Aliases, &loc.Source, &loc.VerificationStatus,
			&loc.CreatedAt, &loc.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		locations = append(locations, loc)
	}

	return locations, total, rows.Err()
}

// Create inserts a new location node
func (r *LocationHierarchyRepo) Create(ctx context.Context, loc *models.LocationHierarchy) (int, error) {
	if loc.Aliases == nil {
		loc.Aliases = []string{}
	}
	query := `
		INSERT INTO location_hierarchy (
			parent_id, location_type, name_en, name_am, name_om, name_ti, slug,
			admin_code, display_order, is_active, is_leaf, latitude, longitude,
			min_lat, max_lat, min_lng, max_lng, aliases, source, verification_status
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7,
			$8, $9, $10, $11, $12, $13,
			$14, $15, $16, $17, $18, $19, $20
		) RETURNING id`

	var id int
	err := r.db.QueryRow(ctx, query,
		loc.ParentID, loc.LocationType, loc.NameEn, loc.NameAm, loc.NameOm, loc.NameTi, loc.Slug,
		loc.AdminCode, loc.DisplayOrder, loc.IsActive, loc.IsLeaf, loc.Latitude, loc.Longitude,
		loc.MinLat, loc.MaxLat, loc.MinLng, loc.MaxLng, loc.Aliases, loc.Source, loc.VerificationStatus,
	).Scan(&id)

	return id, err
}

// Update modifies an existing location node
func (r *LocationHierarchyRepo) Update(ctx context.Context, loc *models.LocationHierarchy) error {
	if loc.Aliases == nil {
		loc.Aliases = []string{}
	}
	query := `
		UPDATE location_hierarchy
		SET parent_id = $2, location_type = $3, name_en = $4, name_am = $5,
		    name_om = $6, name_ti = $7, slug = $8, admin_code = $9,
		    display_order = $10, is_active = $11, is_leaf = $12,
		    latitude = $13, longitude = $14, min_lat = $15, max_lat = $16,
		    min_lng = $17, max_lng = $18, aliases = $19, source = $20,
		    verification_status = $21, updated_at = NOW()
		WHERE id = $1`

	tag, err := r.db.Exec(ctx, query,
		loc.ID, loc.ParentID, loc.LocationType, loc.NameEn, loc.NameAm,
		loc.NameOm, loc.NameTi, loc.Slug, loc.AdminCode,
		loc.DisplayOrder, loc.IsActive, loc.IsLeaf,
		loc.Latitude, loc.Longitude, loc.MinLat, loc.MaxLat,
		loc.MinLng, loc.MaxLng, loc.Aliases, loc.Source,
		loc.VerificationStatus,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("location not found")
	}
	return nil
}

// Delete removes a location node if it has no children or attached products
func (r *LocationHierarchyRepo) Delete(ctx context.Context, id int) error {
	var childCount int
	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM location_hierarchy WHERE parent_id = $1`, id).Scan(&childCount)
	if err != nil {
		return err
	}
	if childCount > 0 {
		return fmt.Errorf("cannot delete location: has %d child locations attached", childCount)
	}

	var productCount int
	err = r.db.QueryRow(ctx, `SELECT COUNT(*) FROM products WHERE location_hierarchy_id = $1`, id).Scan(&productCount)
	if err != nil {
		return err
	}
	if productCount > 0 {
		return fmt.Errorf("cannot delete location: has %d active product listings attached", productCount)
	}

	tag, err := r.db.Exec(ctx, `DELETE FROM location_hierarchy WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("location not found")
	}
	return nil
}
