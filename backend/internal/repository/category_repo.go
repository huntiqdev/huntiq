package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/your-org/used-tools-marketplace/internal/models"
)

type CategoryRepo struct{ db *pgxpool.Pool }

func NewCategoryRepo(db *pgxpool.Pool) *CategoryRepo { return &CategoryRepo{db} }

func scanCat(row interface{ Scan(...interface{}) error }) (models.Category, error) {
	var c models.Category
	err := row.Scan(&c.ID, &c.Name, &c.Slug, &c.NameEn, &c.NameAm, &c.NameOr, &c.NameTi, &c.SortOrder, &c.IsActive, &c.ParentID)
	return c, err
}

// buildTree converts a flat list into a nested tree by parent_id.
func buildTree(flat []models.Category) []models.Category {
	if len(flat) == 0 {
		return []models.Category{}
	}
	
	// Create a map of ID -> index for safe access
	byID := make(map[int]int, len(flat))
	for i := range flat {
		byID[flat[i].ID] = i
	}
	
	// Build tree by appending children to parents
	var rootIndices []int
	for i := range flat {
		if flat[i].ParentID == nil {
			rootIndices = append(rootIndices, i)
		} else if pid := flat[i].ParentID; pid != nil {
			if parentIdx, ok := byID[*pid]; ok && parentIdx < len(flat) {
				// Append child to parent's Children slice
				flat[parentIdx].Children = append(flat[parentIdx].Children, flat[i])
			}
		}
	}
	
	// Collect root nodes
	result := make([]models.Category, 0, len(rootIndices))
	for _, idx := range rootIndices {
		if idx >= 0 && idx < len(flat) {
			result = append(result, flat[idx])
		}
	}
	return result
}

func (r *CategoryRepo) List(ctx context.Context) ([]models.Category, error) {
	//nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
	rows, err := r.db.Query(ctx,
		`SELECT id, name, slug,
		COALESCE(name_en, name, ''), COALESCE(name_am, ''), COALESCE(name_or, ''), COALESCE(name_ti, ''),
		COALESCE(sort_order, 0), COALESCE(is_active, true), parent_id
		FROM categories WHERE COALESCE(is_active,true)=true ORDER BY COALESCE(sort_order,0), name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var flat []models.Category
	for rows.Next() {
		c, err := scanCat(rows)
		if err != nil {
			return nil, err
		}
		flat = append(flat, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return buildTree(flat), nil
}

func (r *CategoryRepo) ListAll(ctx context.Context) ([]models.Category, error) {
	//nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
	rows, err := r.db.Query(ctx,
		`SELECT id, name, slug,
		COALESCE(name_en, name, ''), COALESCE(name_am, ''), COALESCE(name_or, ''), COALESCE(name_ti, ''),
		COALESCE(sort_order, 0), COALESCE(is_active, true), parent_id
		FROM categories ORDER BY COALESCE(sort_order,0), name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var flat []models.Category
	for rows.Next() {
		c, err := scanCat(rows)
		if err != nil {
			return nil, err
		}
		flat = append(flat, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return buildTree(flat), nil
}

func (r *CategoryRepo) Create(ctx context.Context, c models.Category) (int, error) {
	if c.Slug == "" {
		return 0, fmt.Errorf("slug is required")
	}
	if c.NameEn == "" {
		return 0, fmt.Errorf("name_en is required")
	}
	var id int
	err := r.db.QueryRow(ctx, //nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
		`INSERT INTO categories (name, slug, name_en, name_am, name_or, name_ti, sort_order, is_active, parent_id)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		 RETURNING id`,
		c.NameEn, c.Slug, c.NameEn, c.NameAm, c.NameOr, c.NameTi, c.SortOrder, c.IsActive, c.ParentID,
	).Scan(&id)
	return id, err
}

func (r *CategoryRepo) Update(ctx context.Context, c models.Category) error {
	if c.ID == 0 {
		return fmt.Errorf("id is required")
	}
	name := c.NameEn
	if name == "" {
		name = c.Name
	}
	_, err := r.db.Exec(ctx, //nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
		`UPDATE categories
		 SET name=$1, slug=$2, name_en=$3, name_am=$4, name_or=$5, name_ti=$6,
		     sort_order=$7, is_active=$8, parent_id=$9
		 WHERE id=$10`,
		name, c.Slug, c.NameEn, c.NameAm, c.NameOr, c.NameTi, c.SortOrder, c.IsActive, c.ParentID, c.ID,
	)
	return err
}

func (r *CategoryRepo) Delete(ctx context.Context, id int) error {
	// ON DELETE CASCADE handles children automatically
	_, err := r.db.Exec(ctx, `DELETE FROM categories WHERE id=$1`, id) //nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
	return err
}
