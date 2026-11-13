package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/your-org/used-tools-marketplace/internal/models"
)

type ClassificationRepo struct {
	db *pgxpool.Pool
}

func NewClassificationRepo(db *pgxpool.Pool) *ClassificationRepo {
	return &ClassificationRepo{db: db}
}

// ── Schemes ──────────────────────────────────────────────────────────────────

func (r *ClassificationRepo) CreateScheme(ctx context.Context, s *models.ClassificationScheme) (*models.ClassificationScheme, error) {
	query := `
		INSERT INTO classification_schemes (
			code, name_en, name_am, name_om, name_ti, description, icon,
			status, version, revision, is_system, rules, created_by, updated_by
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, 1, 1, $9, $10, $11, $11)
		RETURNING id, created_at, updated_at`

	rulesJSON := s.Rules
	if len(rulesJSON) == 0 {
		rulesJSON = json.RawMessage("{}")
	}

	err := r.db.QueryRow(ctx, query,
		s.Code, s.NameEn, s.NameAm, s.NameOm, s.NameTi, s.Description, s.Icon,
		s.Status, s.IsSystem, rulesJSON, s.CreatedBy,
	).Scan(&s.ID, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("create scheme failed: %w", err)
	}
	return s, nil
}

func (r *ClassificationRepo) GetSchemeByCode(ctx context.Context, code string) (*models.ClassificationScheme, error) {
	query := `
		SELECT id, code, name_en, COALESCE(name_am, ''), COALESCE(name_om, ''), COALESCE(name_ti, ''),
		       COALESCE(description, ''), COALESCE(icon, ''), status, version, revision, is_system,
		       rules, created_by, updated_by, created_at, updated_at
		FROM classification_schemes
		WHERE code = $1`

	var s models.ClassificationScheme
	err := r.db.QueryRow(ctx, query, code).Scan(
		&s.ID, &s.Code, &s.NameEn, &s.NameAm, &s.NameOm, &s.NameTi,
		&s.Description, &s.Icon, &s.Status, &s.Version, &s.Revision, &s.IsSystem,
		&s.Rules, &s.CreatedBy, &s.UpdatedBy, &s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &s, nil
}

func (r *ClassificationRepo) GetSchemeByID(ctx context.Context, id int) (*models.ClassificationScheme, error) {
	query := `
		SELECT id, code, name_en, COALESCE(name_am, ''), COALESCE(name_om, ''), COALESCE(name_ti, ''),
		       COALESCE(description, ''), COALESCE(icon, ''), status, version, revision, is_system,
		       rules, created_by, updated_by, created_at, updated_at
		FROM classification_schemes
		WHERE id = $1`

	var s models.ClassificationScheme
	err := r.db.QueryRow(ctx, query, id).Scan(
		&s.ID, &s.Code, &s.NameEn, &s.NameAm, &s.NameOm, &s.NameTi,
		&s.Description, &s.Icon, &s.Status, &s.Version, &s.Revision, &s.IsSystem,
		&s.Rules, &s.CreatedBy, &s.UpdatedBy, &s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &s, nil
}

func (r *ClassificationRepo) ListSchemes(ctx context.Context, status string) ([]models.ClassificationScheme, error) {
	query := `
		SELECT id, code, name_en, COALESCE(name_am, ''), COALESCE(name_om, ''), COALESCE(name_ti, ''),
		       COALESCE(description, ''), COALESCE(icon, ''), status, version, revision, is_system,
		       rules, created_by, updated_by, created_at, updated_at
		FROM classification_schemes`

	args := []interface{}{}
	if status != "" {
		query += ` WHERE status = $1`
		args = append(args, status)
	}
	query += ` ORDER BY id ASC`

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schemes []models.ClassificationScheme
	for rows.Next() {
		var s models.ClassificationScheme
		if err := rows.Scan(
			&s.ID, &s.Code, &s.NameEn, &s.NameAm, &s.NameOm, &s.NameTi,
			&s.Description, &s.Icon, &s.Status, &s.Version, &s.Revision, &s.IsSystem,
			&s.Rules, &s.CreatedBy, &s.UpdatedBy, &s.CreatedAt, &s.UpdatedAt,
		); err != nil {
			return nil, err
		}
		schemes = append(schemes, s)
	}
	return schemes, nil
}

func (r *ClassificationRepo) UpdateScheme(ctx context.Context, s *models.ClassificationScheme) error {
	query := `
		UPDATE classification_schemes
		SET name_en = $1, name_am = $2, name_om = $3, name_ti = $4,
		    description = $5, icon = $6, status = $7, rules = $8,
		    revision = revision + 1, updated_by = $9, updated_at = NOW()
		WHERE id = $10`

	rulesJSON := s.Rules
	if len(rulesJSON) == 0 {
		rulesJSON = json.RawMessage("{}")
	}

	_, err := r.db.Exec(ctx, query,
		s.NameEn, s.NameAm, s.NameOm, s.NameTi, s.Description, s.Icon,
		s.Status, rulesJSON, s.UpdatedBy, s.ID,
	)
	return err
}

// ── Concepts ─────────────────────────────────────────────────────────────────

func (r *ClassificationRepo) CreateConcept(ctx context.Context, c *models.ClassificationConcept) (*models.ClassificationConcept, error) {
	query := `
		INSERT INTO classification_concepts (
			code, canonical_name, name_en, name_am, name_om, name_ti,
			description, concept_type, status, brand_id, metadata, version,
			created_by, updated_by
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, 1, $12, $12)
		RETURNING id, created_at, updated_at`

	meta := c.Metadata
	if len(meta) == 0 {
		meta = json.RawMessage("{}")
	}
	if c.Status == "" {
		c.Status = "ACTIVE"
	}
	if c.ConceptType == "" {
		c.ConceptType = "PRODUCT"
	}

	err := r.db.QueryRow(ctx, query,
		c.Code, c.CanonicalName, c.NameEn, c.NameAm, c.NameOm, c.NameTi,
		c.Description, c.ConceptType, c.Status, c.BrandID, meta, c.CreatedBy,
	).Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("create concept failed: %w", err)
	}
	return c, nil
}

func (r *ClassificationRepo) GetConceptByID(ctx context.Context, id int) (*models.ClassificationConcept, error) {
	query := `
		SELECT id, code, canonical_name, name_en, COALESCE(name_am, ''), COALESCE(name_om, ''), COALESCE(name_ti, ''),
		       COALESCE(description, ''), concept_type, status, brand_id, metadata, version,
		       created_by, updated_by, created_at, updated_at
		FROM classification_concepts
		WHERE id = $1`

	var c models.ClassificationConcept
	err := r.db.QueryRow(ctx, query, id).Scan(
		&c.ID, &c.Code, &c.CanonicalName, &c.NameEn, &c.NameAm, &c.NameOm, &c.NameTi,
		&c.Description, &c.ConceptType, &c.Status, &c.BrandID, &c.Metadata, &c.Version,
		&c.CreatedBy, &c.UpdatedBy, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}

func (r *ClassificationRepo) GetConceptByCode(ctx context.Context, code string) (*models.ClassificationConcept, error) {
	query := `
		SELECT id, code, canonical_name, name_en, COALESCE(name_am, ''), COALESCE(name_om, ''), COALESCE(name_ti, ''),
		       COALESCE(description, ''), concept_type, status, brand_id, metadata, version,
		       created_by, updated_by, created_at, updated_at
		FROM classification_concepts
		WHERE code = $1`

	var c models.ClassificationConcept
	err := r.db.QueryRow(ctx, query, code).Scan(
		&c.ID, &c.Code, &c.CanonicalName, &c.NameEn, &c.NameAm, &c.NameOm, &c.NameTi,
		&c.Description, &c.ConceptType, &c.Status, &c.BrandID, &c.Metadata, &c.Version,
		&c.CreatedBy, &c.UpdatedBy, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}

func (r *ClassificationRepo) SearchConcepts(ctx context.Context, f models.ConceptSearchFilter) ([]models.ClassificationConcept, int, error) {
	where := []string{"1=1"}
	args := []interface{}{}
	idx := 1

	if f.ConceptType != "" {
		where = append(where, fmt.Sprintf("concept_type = $%d", idx))
		args = append(args, f.ConceptType)
		idx++
	}
	if f.Status != "" {
		where = append(where, fmt.Sprintf("status = $%d", idx))
		args = append(args, f.Status)
		idx++
	}
	if f.SearchQuery != "" {
		term := "%" + strings.ToLower(strings.TrimSpace(f.SearchQuery)) + "%"
		where = append(where, fmt.Sprintf("(LOWER(canonical_name) LIKE $%d OR LOWER(name_en) LIKE $%d OR LOWER(code) LIKE $%d OR LOWER(COALESCE(name_am,'')) LIKE $%d)", idx, idx, idx, idx))
		args = append(args, term)
		idx++
	}

	whereClause := strings.Join(where, " AND ")

	countQuery := fmt.Sprintf("SELECT count(*) FROM classification_concepts WHERE %s", whereClause)
	var total int
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	limit := f.Limit
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	offset := f.Page * limit
	if offset < 0 {
		offset = 0
	}

	query := fmt.Sprintf(`
		SELECT id, code, canonical_name, name_en, COALESCE(name_am, ''), COALESCE(name_om, ''), COALESCE(name_ti, ''),
		       COALESCE(description, ''), concept_type, status, brand_id, metadata, version,
		       created_by, updated_by, created_at, updated_at
		FROM classification_concepts
		WHERE %s
		ORDER BY id ASC
		LIMIT $%d OFFSET $%d`, whereClause, idx, idx+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []models.ClassificationConcept
	for rows.Next() {
		var c models.ClassificationConcept
		if err := rows.Scan(
			&c.ID, &c.Code, &c.CanonicalName, &c.NameEn, &c.NameAm, &c.NameOm, &c.NameTi,
			&c.Description, &c.ConceptType, &c.Status, &c.BrandID, &c.Metadata, &c.Version,
			&c.CreatedBy, &c.UpdatedBy, &c.CreatedAt, &c.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		list = append(list, c)
	}
	return list, total, nil
}

func (r *ClassificationRepo) UpdateConcept(ctx context.Context, c *models.ClassificationConcept) error {
	query := `
		UPDATE classification_concepts
		SET canonical_name = $1, name_en = $2, name_am = $3, name_om = $4, name_ti = $5,
		    description = $6, concept_type = $7, status = $8, brand_id = $9, metadata = $10,
		    version = version + 1, updated_by = $11, updated_at = NOW()
		WHERE id = $12`

	meta := c.Metadata
	if len(meta) == 0 {
		meta = json.RawMessage("{}")
	}

	_, err := r.db.Exec(ctx, query,
		c.CanonicalName, c.NameEn, c.NameAm, c.NameOm, c.NameTi,
		c.Description, c.ConceptType, c.Status, c.BrandID, meta, c.UpdatedBy, c.ID,
	)
	return err
}

func (r *ClassificationRepo) CheckDuplicateCandidates(ctx context.Context, name, code string) ([]models.ClassificationConcept, error) {
	query := `
		SELECT id, code, canonical_name, name_en, COALESCE(name_am, ''), COALESCE(name_om, ''), COALESCE(name_ti, ''),
		       COALESCE(description, ''), concept_type, status, brand_id, metadata, version,
		       created_by, updated_by, created_at, updated_at
		FROM classification_concepts
		WHERE LOWER(code) = LOWER($1)
		   OR LOWER(canonical_name) = LOWER($2)
		   OR similarity(name_en, $2) > 0.4
		ORDER BY similarity(name_en, $2) DESC
		LIMIT 5`

	rows, err := r.db.Query(ctx, query, code, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var candidates []models.ClassificationConcept
	for rows.Next() {
		var c models.ClassificationConcept
		if err := rows.Scan(
			&c.ID, &c.Code, &c.CanonicalName, &c.NameEn, &c.NameAm, &c.NameOm, &c.NameTi,
			&c.Description, &c.ConceptType, &c.Status, &c.BrandID, &c.Metadata, &c.Version,
			&c.CreatedBy, &c.UpdatedBy, &c.CreatedAt, &c.UpdatedAt,
		); err != nil {
			return nil, err
		}
		candidates = append(candidates, c)
	}
	return candidates, nil
}

// ── Placements & Tree Traversal ──────────────────────────────────────────────

func (r *ClassificationRepo) CreatePlacement(ctx context.Context, p *models.ClassificationPlacement) (*models.ClassificationPlacement, error) {
	query := `
		INSERT INTO classification_placements (
			scheme_id, concept_id, parent_placement_id, role, display_name_override,
			slug, display_order, is_leaf, is_active, status, created_by, updated_by
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $11)
		RETURNING id, path, depth, created_at, updated_at`

	if p.Role == "" {
		p.Role = "CATEGORY"
	}
	if p.Status == "" {
		p.Status = "PUBLISHED"
	}

	err := r.db.QueryRow(ctx, query,
		p.SchemeID, p.ConceptID, p.ParentPlacementID, p.Role, p.DisplayNameOverride,
		p.Slug, p.DisplayOrder, p.IsLeaf, p.IsActive, p.Status, p.CreatedBy,
	).Scan(&p.ID, &p.Path, &p.Depth, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("create placement failed: %w", err)
	}
	return p, nil
}

func (r *ClassificationRepo) GetPlacementByID(ctx context.Context, id int, lang string) (*models.ClassificationPlacement, error) {
	query := `
		SELECT cp.id, cp.scheme_id, cp.concept_id, cp.parent_placement_id, cp.role,
		       COALESCE(cp.display_name_override, ''), cp.slug, cp.path, cp.depth,
		       cp.display_order, cp.is_leaf, cp.is_active, cp.status,
		       cp.created_by, cp.updated_by, cp.created_at, cp.updated_at,
		       cc.code, cc.canonical_name, cc.name_en, COALESCE(cc.name_am, ''),
		       COALESCE(cc.name_om, ''), COALESCE(cc.name_ti, ''), cc.concept_type,
		       COALESCE(cc.metadata->>'icon', '')
		FROM classification_placements cp
		JOIN classification_concepts cc ON cc.id = cp.concept_id
		WHERE cp.id = $1`

	var p models.ClassificationPlacement
	var c models.ClassificationConcept
	var icon string

	err := r.db.QueryRow(ctx, query, id).Scan(
		&p.ID, &p.SchemeID, &p.ConceptID, &p.ParentPlacementID, &p.Role,
		&p.DisplayNameOverride, &p.Slug, &p.Path, &p.Depth,
		&p.DisplayOrder, &p.IsLeaf, &p.IsActive, &p.Status,
		&p.CreatedBy, &p.UpdatedBy, &p.CreatedAt, &p.UpdatedAt,
		&c.Code, &c.CanonicalName, &c.NameEn, &c.NameAm,
		&c.NameOm, &c.NameTi, &c.ConceptType, &icon,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	c.ID = p.ConceptID
	p.Concept = &c

	breadcrumbs, err := r.GetBreadcrumbs(ctx, p.ID, lang)
	if err == nil {
		p.Breadcrumbs = breadcrumbs
	}

	return &p, nil
}

func (r *ClassificationRepo) GetBreadcrumbs(ctx context.Context, placementID int, lang string) ([]models.CategoryBreadcrumb, error) {
	query := `
		SELECT ancestor.id, 
		       COALESCE(NULLIF(ancestor.display_name_override, ''), cc.name_en) AS name,
		       ancestor.slug, ancestor.depth
		FROM classification_placements target
		JOIN classification_placements ancestor 
		  ON ancestor.scheme_id = target.scheme_id 
		 AND target.path LIKE ancestor.path || '%'
		JOIN classification_concepts cc ON cc.id = ancestor.concept_id
		WHERE target.id = $1
		ORDER BY ancestor.depth ASC`

	rows, err := r.db.Query(ctx, query, placementID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var crumbs []models.CategoryBreadcrumb
	for rows.Next() {
		var b models.CategoryBreadcrumb
		if err := rows.Scan(&b.ID, &b.Name, &b.Slug, &b.Level); err != nil {
			return nil, err
		}
		crumbs = append(crumbs, b)
	}
	return crumbs, nil
}

func (r *ClassificationRepo) GetTree(ctx context.Context, schemeCode, lang string, maxDepth int) ([]models.PlacementTreeNode, error) {
	if schemeCode == "" {
		schemeCode = "SYS_PROD"
	}
	if maxDepth <= 0 {
		maxDepth = 10
	}

	query := `
		SELECT cp.id, cp.scheme_id, cp.concept_id, cp.parent_placement_id,
		       COALESCE(NULLIF(cp.display_name_override, ''), cc.name_en),
		       cc.name_en, COALESCE(cc.name_am, ''), COALESCE(cc.name_om, ''), COALESCE(cc.name_ti, ''),
		       cp.slug, cp.role, cp.depth, cp.is_leaf, cp.is_active, cp.display_order,
		       COALESCE(cc.metadata->>'icon', ''),
		       COALESCE(s.product_count, 0)
		FROM classification_placements cp
		JOIN classification_schemes cs ON cs.id = cp.scheme_id
		JOIN classification_concepts cc ON cc.id = cp.concept_id
		LEFT JOIN (
			SELECT placement_id, count(*) AS product_count
			FROM products
			WHERE status = 'approved' AND placement_id IS NOT NULL
			GROUP BY placement_id
		) s ON s.placement_id = cp.id
		WHERE cs.code = $1 AND cp.is_active = true AND cp.depth <= $2
		ORDER BY cp.depth ASC, cp.display_order ASC, cp.id ASC`

	rows, err := r.db.Query(ctx, query, schemeCode, maxDepth)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	nodesByID := make(map[int]*models.PlacementTreeNode)
	var rootNodes []models.PlacementTreeNode

	for rows.Next() {
		var n models.PlacementTreeNode
		var icon string
		if err := rows.Scan(
			&n.ID, &n.SchemeID, &n.ConceptID, &n.ParentID,
			&n.Name, &n.NameEn, &n.NameAm, &n.NameOm, &n.NameTi,
			&n.Slug, &n.Role, &n.Depth, &n.IsLeaf, &n.IsActive, &n.DisplayOrder,
			&icon, &n.Count,
		); err != nil {
			return nil, err
		}
		n.Icon = icon
		n.Children = make([]models.PlacementTreeNode, 0)

		// Choose display name based on locale
		switch lang {
		case "am":
			if n.NameAm != "" {
				n.Name = n.NameAm
			}
		case "om", "or":
			if n.NameOm != "" {
				n.Name = n.NameOm
			}
		case "ti":
			if n.NameTi != "" {
				n.Name = n.NameTi
			}
		}

		nodesByID[n.ID] = &n
	}

	for _, n := range nodesByID {
		if n.ParentID == nil {
			rootNodes = append(rootNodes, *n)
		} else if parent, ok := nodesByID[*n.ParentID]; ok {
			parent.Children = append(parent.Children, *n)
		} else {
			rootNodes = append(rootNodes, *n)
		}
	}

	return rootNodes, nil
}

func (r *ClassificationRepo) GetRoots(ctx context.Context, schemeCode, lang string) ([]models.PlacementTreeNode, error) {
	if schemeCode == "" {
		schemeCode = "SYS_PROD"
	}
	query := `
		SELECT cp.id, cp.scheme_id, cp.concept_id, cp.parent_placement_id,
		       COALESCE(NULLIF(cp.display_name_override, ''), cc.name_en),
		       cc.name_en, COALESCE(cc.name_am, ''), COALESCE(cc.name_om, ''), COALESCE(cc.name_ti, ''),
		       cp.slug, cp.role, cp.depth, cp.is_leaf, cp.is_active, cp.display_order,
		       COALESCE(cc.metadata->>'icon', ''),
		       COALESCE(s.product_count, 0)
		FROM classification_placements cp
		JOIN classification_schemes cs ON cs.id = cp.scheme_id
		JOIN classification_concepts cc ON cc.id = cp.concept_id
		LEFT JOIN (
			SELECT placement_id, count(*) AS product_count
			FROM products
			WHERE status = 'approved' AND placement_id IS NOT NULL
			GROUP BY placement_id
		) s ON s.placement_id = cp.id
		WHERE cs.code = $1 AND cp.parent_placement_id IS NULL AND cp.is_active = true
		ORDER BY cp.display_order ASC, cp.id ASC`

	rows, err := r.db.Query(ctx, query, schemeCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roots []models.PlacementTreeNode
	for rows.Next() {
		var n models.PlacementTreeNode
		var icon string
		if err := rows.Scan(
			&n.ID, &n.SchemeID, &n.ConceptID, &n.ParentID,
			&n.Name, &n.NameEn, &n.NameAm, &n.NameOm, &n.NameTi,
			&n.Slug, &n.Role, &n.Depth, &n.IsLeaf, &n.IsActive, &n.DisplayOrder,
			&icon, &n.Count,
		); err != nil {
			return nil, err
		}
		n.Icon = icon
		switch lang {
		case "am":
			if n.NameAm != "" {
				n.Name = n.NameAm
			}
		case "om", "or":
			if n.NameOm != "" {
				n.Name = n.NameOm
			}
		case "ti":
			if n.NameTi != "" {
				n.Name = n.NameTi
			}
		}
		roots = append(roots, n)
	}
	return roots, nil
}

func (r *ClassificationRepo) GetChildren(ctx context.Context, parentPlacementID int, lang string) ([]models.PlacementTreeNode, error) {
	query := `
		SELECT cp.id, cp.scheme_id, cp.concept_id, cp.parent_placement_id,
		       COALESCE(NULLIF(cp.display_name_override, ''), cc.name_en),
		       cc.name_en, COALESCE(cc.name_am, ''), COALESCE(cc.name_om, ''), COALESCE(cc.name_ti, ''),
		       cp.slug, cp.role, cp.depth, cp.is_leaf, cp.is_active, cp.display_order,
		       COALESCE(cc.metadata->>'icon', ''),
		       COALESCE(s.product_count, 0)
		FROM classification_placements cp
		JOIN classification_concepts cc ON cc.id = cp.concept_id
		LEFT JOIN (
			SELECT placement_id, count(*) AS product_count
			FROM products
			WHERE status = 'approved' AND placement_id IS NOT NULL
			GROUP BY placement_id
		) s ON s.placement_id = cp.id
		WHERE cp.parent_placement_id = $1 AND cp.is_active = true
		ORDER BY cp.display_order ASC, cp.id ASC`

	rows, err := r.db.Query(ctx, query, parentPlacementID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var children []models.PlacementTreeNode
	for rows.Next() {
		var n models.PlacementTreeNode
		var icon string
		if err := rows.Scan(
			&n.ID, &n.SchemeID, &n.ConceptID, &n.ParentID,
			&n.Name, &n.NameEn, &n.NameAm, &n.NameOm, &n.NameTi,
			&n.Slug, &n.Role, &n.Depth, &n.IsLeaf, &n.IsActive, &n.DisplayOrder,
			&icon, &n.Count,
		); err != nil {
			return nil, err
		}
		n.Icon = icon
		switch lang {
		case "am":
			if n.NameAm != "" {
				n.Name = n.NameAm
			}
		case "om", "or":
			if n.NameOm != "" {
				n.Name = n.NameOm
			}
		case "ti":
			if n.NameTi != "" {
				n.Name = n.NameTi
			}
		}
		children = append(children, n)
	}
	return children, nil
}

func (r *ClassificationRepo) MovePlacement(ctx context.Context, placementID, newParentPlacementID int) error {
	// Trigger update_classification_placement_path will enforce scheme consistency,
	// calculate new depth and path, and cascade_classification_placement_path will update descendants.
	var newParent *int
	if newParentPlacementID > 0 {
		newParent = &newParentPlacementID
	}

	query := `
		UPDATE classification_placements
		SET parent_placement_id = $1, updated_at = NOW()
		WHERE id = $2`

	res, err := r.db.Exec(ctx, query, newParent, placementID)
	if err != nil {
		return fmt.Errorf("move placement failed: %w", err)
	}
	if res.RowsAffected() == 0 {
		return errors.New("placement not found")
	}
	return nil
}

func (r *ClassificationRepo) UpdatePlacement(ctx context.Context, p *models.ClassificationPlacement) error {
	query := `
		UPDATE classification_placements
		SET role = $1, display_name_override = $2, display_order = $3,
		    is_leaf = $4, is_active = $5, status = $6, updated_by = $7, updated_at = NOW()
		WHERE id = $8`

	res, err := r.db.Exec(ctx, query,
		p.Role, p.DisplayNameOverride, p.DisplayOrder, p.IsLeaf,
		p.IsActive, p.Status, p.UpdatedBy, p.ID,
	)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return errors.New("placement not found")
	}
	return nil
}

func (r *ClassificationRepo) DeletePlacement(ctx context.Context, placementID int) error {
	// Safe retirement/archiving rather than hard delete to preserve listing references
	query := `UPDATE classification_placements SET is_active = false, status = 'ARCHIVED', updated_at = NOW() WHERE id = $1`
	res, err := r.db.Exec(ctx, query, placementID)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return errors.New("placement not found")
	}
	return nil
}

// ── Multi-Industry & Multi-Application Relationships ────────────────────────

func (r *ClassificationRepo) CreateRelationship(ctx context.Context, rel *models.ClassificationRelationship) error {
	query := `
		INSERT INTO classification_relationships (
			source_concept_id, target_concept_id, relationship_type_code,
			context_scheme_id, status, metadata, created_by
		) VALUES ($1, $2, $3, $4, 'ACTIVE', $5, $6)
		RETURNING id, created_at, updated_at`

	meta := rel.Metadata
	if len(meta) == 0 {
		meta = json.RawMessage("{}")
	}

	return r.db.QueryRow(ctx, query,
		rel.SourceConceptID, rel.TargetConceptID, rel.RelationshipTypeCode,
		rel.ContextSchemeID, meta, rel.CreatedBy,
	).Scan(&rel.ID, &rel.CreatedAt, &rel.UpdatedAt)
}

func (r *ClassificationRepo) GetRelationships(ctx context.Context, conceptID int, relType string) ([]models.ClassificationRelationship, error) {
	query := `
		SELECT cr.id, cr.source_concept_id, cr.target_concept_id, cr.relationship_type_code,
		       cr.context_scheme_id, cr.status, cr.metadata, cr.created_by, cr.created_at, cr.updated_at,
		       src.code, src.canonical_name, src.name_en,
		       tgt.code, tgt.canonical_name, tgt.name_en
		FROM classification_relationships cr
		JOIN classification_concepts src ON src.id = cr.source_concept_id
		JOIN classification_concepts tgt ON tgt.id = cr.target_concept_id
		WHERE (cr.source_concept_id = $1 OR cr.target_concept_id = $1)`

	args := []interface{}{conceptID}
	if relType != "" {
		query += ` AND cr.relationship_type_code = $2`
		args = append(args, relType)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rels []models.ClassificationRelationship
	for rows.Next() {
		var rel models.ClassificationRelationship
		var src models.ClassificationConcept
		var tgt models.ClassificationConcept
		if err := rows.Scan(
			&rel.ID, &rel.SourceConceptID, &rel.TargetConceptID, &rel.RelationshipTypeCode,
			&rel.ContextSchemeID, &rel.Status, &rel.Metadata, &rel.CreatedBy, &rel.CreatedAt, &rel.UpdatedAt,
			&src.Code, &src.CanonicalName, &src.NameEn,
			&tgt.Code, &tgt.CanonicalName, &tgt.NameEn,
		); err != nil {
			return nil, err
		}
		src.ID = rel.SourceConceptID
		tgt.ID = rel.TargetConceptID
		rel.SourceConcept = &src
		rel.TargetConcept = &tgt
		rels = append(rels, rel)
	}
	return rels, nil
}

func (r *ClassificationRepo) ListRelationshipTypes(ctx context.Context) ([]models.ClassificationRelationshipType, error) {
	query := `
		SELECT code, name, COALESCE(description, ''), is_directed, COALESCE(inverse_type_code, ''), is_active
		FROM classification_relationship_types
		WHERE is_active = true
		ORDER BY code ASC`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var types []models.ClassificationRelationshipType
	for rows.Next() {
		var t models.ClassificationRelationshipType
		if err := rows.Scan(&t.Code, &t.Name, &t.Description, &t.IsDirected, &t.InverseTypeCode, &t.IsActive); err != nil {
			return nil, err
		}
		types = append(types, t)
	}
	return types, nil
}

func (r *ClassificationRepo) DeleteRelationship(ctx context.Context, id int) error {
	query := `DELETE FROM classification_relationships WHERE id = $1`
	res, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return errors.New("relationship not found")
	}
	return nil
}

// ── Dynamic Attributes Inheritance ──────────────────────────────────────────

func (r *ClassificationRepo) GetEffectiveAttributes(ctx context.Context, placementID int) ([]models.ClassificationAttributeDTO, error) {
	// Materialized path ancestor lookup
	query := `
		WITH ancestor_placements AS (
			SELECT unnest(string_to_array(trim(both '/' from path), '/')::int[]) AS ancestor_id
			FROM classification_placements
			WHERE id = $1
		),
		candidate_attributes AS (
			SELECT 
				ad.id AS attribute_id,
				ad.name_en,
				COALESCE(ad.name_am, '') AS name_am,
				COALESCE(ad.name_om, '') AS name_om,
				COALESCE(ad.name_ti, '') AS name_ti,
				ad.slug,
				ad.data_type,
				COALESCE(ar.override_validation_rules, ad.validation_rules) AS validation_rules,
				COALESCE(ar.override_unit, ad.unit, '') AS unit,
				ar.is_required,
				ad.is_filterable,
				ad.is_searchable,
				(ar.placement_id != $1) AS is_inherited,
				ar.display_order,
				ROW_NUMBER() OVER (
					PARTITION BY ad.id 
					ORDER BY 
						CASE WHEN ar.placement_id = $1 THEN 0 ELSE 1 END,
						ar.display_order ASC
				) AS rn
			FROM classification_attribute_rules ar
			JOIN attribute_definitions ad ON ad.id = ar.attribute_id
			WHERE ar.placement_id IN (SELECT ancestor_id FROM ancestor_placements)
		)
		SELECT 
			attribute_id, name_en, name_am, name_om, name_ti, slug,
			data_type, validation_rules, unit, is_required, is_filterable,
			is_searchable, is_inherited, display_order
		FROM candidate_attributes
		WHERE rn = 1
		ORDER BY display_order ASC, attribute_id ASC`

	rows, err := r.db.Query(ctx, query, placementID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attrs []models.ClassificationAttributeDTO
	for rows.Next() {
		var a models.ClassificationAttributeDTO
		if err := rows.Scan(
			&a.AttributeID, &a.NameEn, &a.NameAm, &a.NameOm, &a.NameTi, &a.Slug,
			&a.DataType, &a.ValidationRules, &a.Unit, &a.IsRequired, &a.IsFilterable,
			&a.IsSearchable, &a.IsInherited, &a.DisplayOrder,
		); err != nil {
			return nil, err
		}
		attrs = append(attrs, a)
	}
	return attrs, nil
}

func (r *ClassificationRepo) SaveAttributeRule(ctx context.Context, rule *models.ClassificationAttributeRule) error {
	query := `
		INSERT INTO classification_attribute_rules (
			placement_id, concept_id, attribute_id, is_required, is_inherited,
			display_order, override_validation_rules, override_unit
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (placement_id, attribute_id) DO UPDATE SET
			is_required = EXCLUDED.is_required,
			is_inherited = EXCLUDED.is_inherited,
			display_order = EXCLUDED.display_order,
			override_validation_rules = EXCLUDED.override_validation_rules,
			override_unit = EXCLUDED.override_unit`

	_, err := r.db.Exec(ctx, query,
		rule.PlacementID, rule.ConceptID, rule.AttributeID, rule.IsRequired,
		rule.IsInherited, rule.DisplayOrder, rule.OverrideValidationRules, rule.OverrideUnit,
	)
	return err
}

func (r *ClassificationRepo) DeleteAttributeRule(ctx context.Context, id int) error {
	query := `DELETE FROM classification_attribute_rules WHERE id = $1`
	res, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return errors.New("attribute rule not found")
	}
	return nil
}

// ── Impact Analysis & Auditing ──────────────────────────────────────────────

func (r *ClassificationRepo) AnalyzeImpact(ctx context.Context, entityType string, entityID int) (*models.ImpactAnalysisResult, error) {
	result := &models.ImpactAnalysisResult{
		EntityID:   entityID,
		EntityType: entityType,
		Warnings:   make([]string, 0),
		Details:    make(map[string]interface{}),
	}

	switch entityType {
	case "PLACEMENT":
		// Descendant placements count
		var descCount int
		_ = r.db.QueryRow(ctx, `
			SELECT count(*) FROM classification_placements desc_p
			JOIN classification_placements target ON target.id = $1
			WHERE desc_p.path LIKE target.path || '%' AND desc_p.id <> target.id`, entityID).Scan(&descCount)
		result.AffectedDescendantsCount = descCount

		// Listings count linked directly or indirectly
		var listingCount int
		_ = r.db.QueryRow(ctx, `
			SELECT count(*) FROM products p
			JOIN classification_placements target ON target.id = $1
			JOIN classification_placements sub ON sub.path LIKE target.path || '%'
			WHERE p.placement_id = sub.id`, entityID).Scan(&listingCount)
		result.AffectedListingsCount = listingCount

		if listingCount > 0 {
			result.Warnings = append(result.Warnings, fmt.Sprintf("%d active marketplace listings depend on this category or its subcategories.", listingCount))
			result.CanSafelyExecute = false
		} else {
			result.CanSafelyExecute = true
		}

	case "CONCEPT":
		var placementCount int
		_ = r.db.QueryRow(ctx, `SELECT count(*) FROM classification_placements WHERE concept_id = $1`, entityID).Scan(&placementCount)
		result.AffectedPlacementsCount = placementCount

		var relCount int
		_ = r.db.QueryRow(ctx, `SELECT count(*) FROM classification_relationships WHERE source_concept_id = $1 OR target_concept_id = $1`, entityID).Scan(&relCount)
		result.AffectedRelationships = relCount

		var listingCount int
		_ = r.db.QueryRow(ctx, `SELECT count(*) FROM products WHERE concept_id = $1`, entityID).Scan(&listingCount)
		result.AffectedListingsCount = listingCount

		if listingCount > 0 || placementCount > 0 {
			result.Warnings = append(result.Warnings, fmt.Sprintf("Concept is actively referenced by %d placements and %d listings.", placementCount, listingCount))
			result.CanSafelyExecute = false
		} else {
			result.CanSafelyExecute = true
		}
	}

	return result, nil
}

func (r *ClassificationRepo) LogAudit(ctx context.Context, log *models.ClassificationAuditLog) error {
	query := `
		INSERT INTO classification_audit_logs (
			user_id, action, entity_type, entity_id, changes, reason, ip_address
		) VALUES ($1, $2, $3, $4, $5, $6, $7)`

	changes := log.Changes
	if len(changes) == 0 {
		changes = json.RawMessage("{}")
	}

	_, err := r.db.Exec(ctx, query,
		log.UserID, log.Action, log.EntityType, log.EntityID,
		changes, log.Reason, log.IPAddress,
	)
	return err
}

func (r *ClassificationRepo) ListAuditLogs(ctx context.Context, entityType string, entityID int, limit, offset int) ([]models.ClassificationAuditLog, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	query := `
		SELECT id, user_id, action, entity_type, entity_id, changes, COALESCE(reason, ''), COALESCE(ip_address, ''), created_at
		FROM classification_audit_logs
		WHERE ($1 = '' OR entity_type = $1)
		  AND ($2 = 0 OR entity_id = $2)
		ORDER BY id DESC
		LIMIT $3 OFFSET $4`

	rows, err := r.db.Query(ctx, query, entityType, entityID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.ClassificationAuditLog
	for rows.Next() {
		var l models.ClassificationAuditLog
		if err := rows.Scan(
			&l.ID, &l.UserID, &l.Action, &l.EntityType, &l.EntityID,
			&l.Changes, &l.Reason, &l.IPAddress, &l.CreatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, l)
	}
	return list, nil
}

func (r *ClassificationRepo) CreateChangeRequest(ctx context.Context, cr *models.ClassificationChangeRequest) error {
	query := `
		INSERT INTO classification_change_requests (
			entity_type, entity_id, action, status, proposed_changes,
			previous_state, impact_summary, requested_by
		) VALUES ($1, $2, $3, 'DRAFT', $4, $5, $6, $7)
		RETURNING id, created_at, updated_at`

	impact := cr.ImpactSummary
	if len(impact) == 0 {
		impact = json.RawMessage("{}")
	}

	return r.db.QueryRow(ctx, query,
		cr.EntityType, cr.EntityID, cr.Action, cr.ProposedChanges,
		cr.PreviousState, impact, cr.RequestedBy,
	).Scan(&cr.ID, &cr.CreatedAt, &cr.UpdatedAt)
}

func (r *ClassificationRepo) ReviewChangeRequest(ctx context.Context, id int, status string, reviewerID int, reason string) error {
	query := `
		UPDATE classification_change_requests
		SET status = $1, reviewed_by = $2, rejection_reason = $3, updated_at = NOW()
		WHERE id = $4`

	res, err := r.db.Exec(ctx, query, status, reviewerID, reason, id)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return errors.New("change request not found")
	}
	return nil
}

// ── Universal Deterministic Search ──────────────────────────────────────────

func (r *ClassificationRepo) SearchUniversal(ctx context.Context, query, schemeCode, lang string, limit int) (*models.UniversalSearchResult, error) {
	if limit <= 0 || limit > 50 {
		limit = 20
	}
	cleanTerm := strings.TrimSpace(query)
	result := &models.UniversalSearchResult{
		Query:        query,
		Concepts:     make([]models.ClassificationConcept, 0),
		Placements:   make([]models.PlacementTreeNode, 0),
		TotalResults: 0,
	}

	if cleanTerm == "" {
		return result, nil
	}

	likePattern := "%" + strings.ToLower(cleanTerm) + "%"

	// 1. Search Canonical Concepts
	conceptQuery := `
		SELECT id, code, canonical_name, name_en, COALESCE(name_am, ''), COALESCE(name_om, ''), COALESCE(name_ti, ''),
		       COALESCE(description, ''), concept_type, status, brand_id, metadata, version,
		       created_by, updated_by, created_at, updated_at
		FROM classification_concepts
		WHERE status = 'ACTIVE'
		  AND (
		      LOWER(canonical_name) LIKE $1
		   OR LOWER(name_en) LIKE $1
		   OR similarity(name_en, $2) > 0.3
		   OR id IN (SELECT concept_id FROM classification_aliases WHERE LOWER(alias) LIKE $1)
		  )
		ORDER BY similarity(name_en, $2) DESC, id ASC
		LIMIT $3`

	cRows, err := r.db.Query(ctx, conceptQuery, likePattern, cleanTerm, limit)
	if err == nil {
		defer cRows.Close()
		for cRows.Next() {
			var c models.ClassificationConcept
			if err := cRows.Scan(
				&c.ID, &c.Code, &c.CanonicalName, &c.NameEn, &c.NameAm, &c.NameOm, &c.NameTi,
				&c.Description, &c.ConceptType, &c.Status, &c.BrandID, &c.Metadata, &c.Version,
				&c.CreatedBy, &c.UpdatedBy, &c.CreatedAt, &c.UpdatedAt,
			); err == nil {
				result.Concepts = append(result.Concepts, c)
			}
		}
	}

	// 2. Search Placements in requested scheme
	if schemeCode == "" {
		schemeCode = "SYS_PROD"
	}
	placementQuery := `
		SELECT cp.id, cp.scheme_id, cp.concept_id, cp.parent_placement_id,
		       COALESCE(NULLIF(cp.display_name_override, ''), cc.name_en),
		       cc.name_en, COALESCE(cc.name_am, ''), COALESCE(cc.name_om, ''), COALESCE(cc.name_ti, ''),
		       cp.slug, cp.role, cp.depth, cp.is_leaf, cp.is_active, cp.display_order,
		       COALESCE(cc.metadata->>'icon', ''), 0
		FROM classification_placements cp
		JOIN classification_schemes cs ON cs.id = cp.scheme_id
		JOIN classification_concepts cc ON cc.id = cp.concept_id
		WHERE cs.code = $1 AND cp.is_active = true
		  AND (
		      LOWER(cp.slug) LIKE $2
		   OR LOWER(cc.name_en) LIKE $2
		   OR similarity(cc.name_en, $3) > 0.3
		   OR cc.id IN (SELECT concept_id FROM classification_aliases WHERE LOWER(alias) LIKE $2)
		  )
		ORDER BY similarity(cc.name_en, $3) DESC, cp.depth ASC
		LIMIT $4`

	pRows, err := r.db.Query(ctx, placementQuery, schemeCode, likePattern, cleanTerm, limit)
	if err == nil {
		defer pRows.Close()
		for pRows.Next() {
			var p models.PlacementTreeNode
			var icon string
			if err := pRows.Scan(
				&p.ID, &p.SchemeID, &p.ConceptID, &p.ParentID,
				&p.Name, &p.NameEn, &p.NameAm, &p.NameOm, &p.NameTi,
				&p.Slug, &p.Role, &p.Depth, &p.IsLeaf, &p.IsActive, &p.DisplayOrder,
				&icon, &p.Count,
			); err == nil {
				p.Icon = icon
				result.Placements = append(result.Placements, p)
			}
		}
	}

	result.TotalResults = len(result.Concepts) + len(result.Placements)
	return result, nil
}
