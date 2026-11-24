package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/your-org/used-tools-marketplace/internal/models"
)

type MasterDataRepository interface {
	GetExistingLocationsMap(ctx context.Context) (map[int]*models.MasterLocationNode, error)
	GetExistingCategoriesMap(ctx context.Context) (map[int]*models.MasterCategoryNode, error)
	GetLocationsForExport(ctx context.Context, regionCode string) ([]models.MasterLocationNode, error)
	GetCategoriesForExport(ctx context.Context, systemCode, rootSlug string) ([]models.MasterCategoryNode, error)
	UpsertLocationBatch(ctx context.Context, nodes []models.MasterLocationNode, adminID int) error
	UpsertCategoryBatch(ctx context.Context, nodes []models.MasterCategoryNode, adminID int) error
	GetMasterDataStatus(ctx context.Context) (*models.MasterDataStatus, error)
	LogMasterDataAudit(ctx context.Context, adminID int, action string, catID *int, oldVal, newVal interface{}, reason, ip string) error
}

type MasterDataRepo struct {
	db *pgxpool.Pool
}

func NewMasterDataRepo(db *pgxpool.Pool) *MasterDataRepo {
	return &MasterDataRepo{db: db}
}

func (r *MasterDataRepo) GetExistingLocationsMap(ctx context.Context) (map[int]*models.MasterLocationNode, error) {
	query := `
		SELECT 
			id, parent_id, location_type, name_en, COALESCE(name_am,''), COALESCE(name_om,''), COALESCE(name_ti,''),
			slug, COALESCE(admin_code,''), level, COALESCE(path,''), display_order, is_active, is_leaf, is_deprecated,
			replaced_by_id, latitude, longitude, min_lat, max_lat, min_lng, max_lng,
			COALESCE(aliases, '{}'), COALESCE(source, 'authoritative'), COALESCE(verification_status, 'verified')
		FROM location_hierarchy
	`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query locations: %w", err)
	}
	defer rows.Close()

	result := make(map[int]*models.MasterLocationNode)
	for rows.Next() {
		var n models.MasterLocationNode
		var (
			parentID, replacedByID                   *int
			lat, lng, minLat, maxLat, minLng, maxLng *float64
			aliases                                  []string
		)
		err := rows.Scan(
			&n.ID, &parentID, &n.Type, &n.Names.En, &n.Names.Am, &n.Names.Om, &n.Names.Ti,
			&n.Slug, &n.Code, &n.Level, &n.Path, &n.DisplayOrder, &n.IsActive, &n.IsLeaf, &n.IsDeprecated,
			&replacedByID, &lat, &lng, &minLat, &maxLat, &minLng, &maxLng,
			&aliases, &n.Source, &n.VerificationStatus,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan location: %w", err)
		}
		n.ParentID = parentID
		n.ReplacedByID = replacedByID
		n.Aliases = aliases
		if lat != nil || lng != nil || minLat != nil {
			n.Coordinates = &models.LocationCoordinates{
				Latitude:  lat,
				Longitude: lng,
			}
			if minLat != nil || maxLat != nil {
				n.Coordinates.BBox = &models.BBox{
					MinLat: minLat,
					MaxLat: maxLat,
					MinLng: minLng,
					MaxLng: maxLng,
				}
			}
		}
		result[n.ID] = &n
	}
	return result, nil
}

func (r *MasterDataRepo) GetExistingCategoriesMap(ctx context.Context) (map[int]*models.MasterCategoryNode, error) {
	query := `
		SELECT 
			id, parent_id, classification_system_code, name_en, COALESCE(name_am,''), COALESCE(name_om,''), COALESCE(name_ti,''),
			slug, description, level, COALESCE(path,''), display_order, icon, image_url,
			is_active, is_leaf, is_archived, meta_title, meta_description
		FROM category_hierarchy
	`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query categories: %w", err)
	}
	defer rows.Close()

	result := make(map[int]*models.MasterCategoryNode)
	for rows.Next() {
		var c models.MasterCategoryNode
		var parentID *int
		var desc, icon, img, metaT, metaD *string
		err := rows.Scan(
			&c.ID, &parentID, &c.ClassificationSystemCode, &c.Names.En, &c.Names.Am, &c.Names.Om, &c.Names.Ti,
			&c.Slug, &desc, &c.Level, &c.Path, &c.DisplayOrder, &icon, &img,
			&c.IsActive, &c.IsLeaf, &c.IsArchived, &metaT, &metaD,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan category: %w", err)
		}
		c.ParentID = parentID
		c.Description = desc
		c.Icon = icon
		c.ImageURL = img
		c.MetaTitle = metaT
		c.MetaDescription = metaD
		c.Aliases = []string{}
		c.Redirects = []string{}
		result[c.ID] = &c
	}
	return result, nil
}

func (r *MasterDataRepo) GetLocationsForExport(ctx context.Context, regionCode string) ([]models.MasterLocationNode, error) {
	locs, err := r.GetExistingLocationsMap(ctx)
	if err != nil {
		return nil, err
	}

	var list []models.MasterLocationNode
	for _, l := range locs {
		if regionCode == "" || l.Code == regionCode || (l.ParentCode != nil && *l.ParentCode == regionCode) {
			list = append(list, *l)
		}
	}
	return list, nil
}

func (r *MasterDataRepo) GetCategoriesForExport(ctx context.Context, systemCode, rootSlug string) ([]models.MasterCategoryNode, error) {
	cats, err := r.GetExistingCategoriesMap(ctx)
	if err != nil {
		return nil, err
	}

	var list []models.MasterCategoryNode
	for _, c := range cats {
		if systemCode != "" && c.ClassificationSystemCode != systemCode {
			continue
		}
		list = append(list, *c)
	}
	return list, nil
}

func (r *MasterDataRepo) UpsertLocationBatch(ctx context.Context, nodes []models.MasterLocationNode, adminID int) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Sort locations by Level ASC to guarantee parent locations are inserted before children
	sort.SliceStable(nodes, func(i, j int) bool {
		return nodes[i].Level < nodes[j].Level
	})

	upsertSQL := `
		INSERT INTO location_hierarchy (
			id, parent_id, location_type, name_en, name_am, name_om, name_ti,
			slug, admin_code, level, display_order, is_active, is_leaf, is_deprecated,
			replaced_by_id, latitude, longitude, min_lat, max_lat, min_lng, max_lng,
			aliases, source, verification_status, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7,
			$8, $9, $10, $11, $12, $13, $14,
			$15, $16, $17, $18, $19, $20, $21,
			$22, $23, $24, NOW()
		)
		ON CONFLICT (id) DO UPDATE SET
			parent_id = EXCLUDED.parent_id,
			location_type = EXCLUDED.location_type,
			name_en = EXCLUDED.name_en,
			name_am = EXCLUDED.name_am,
			name_om = EXCLUDED.name_om,
			name_ti = EXCLUDED.name_ti,
			slug = EXCLUDED.slug,
			admin_code = EXCLUDED.admin_code,
			level = EXCLUDED.level,
			display_order = EXCLUDED.display_order,
			is_active = EXCLUDED.is_active,
			is_leaf = EXCLUDED.is_leaf,
			is_deprecated = EXCLUDED.is_deprecated,
			replaced_by_id = EXCLUDED.replaced_by_id,
			latitude = EXCLUDED.latitude,
			longitude = EXCLUDED.longitude,
			min_lat = EXCLUDED.min_lat,
			max_lat = EXCLUDED.max_lat,
			min_lng = EXCLUDED.min_lng,
			max_lng = EXCLUDED.max_lng,
			aliases = EXCLUDED.aliases,
			source = EXCLUDED.source,
			verification_status = EXCLUDED.verification_status,
			updated_at = NOW();
	`

	batchSize := 200
	for i := 0; i < len(nodes); i += batchSize {
		end := i + batchSize
		if end > len(nodes) {
			end = len(nodes)
		}
		chunk := nodes[i:end]

		batch := &pgx.Batch{}
		for _, n := range chunk {
			var lat, lng, minLat, maxLat, minLng, maxLng *float64
			if n.Coordinates != nil {
				lat = n.Coordinates.Latitude
				lng = n.Coordinates.Longitude
				if n.Coordinates.BBox != nil {
					minLat = n.Coordinates.BBox.MinLat
					maxLat = n.Coordinates.BBox.MaxLat
					minLng = n.Coordinates.BBox.MinLng
					maxLng = n.Coordinates.BBox.MaxLng
				}
			}

			batch.Queue(upsertSQL,
				n.ID, n.ParentID, n.Type, n.Names.En, n.Names.Am, n.Names.Om, n.Names.Ti,
				n.Slug, n.Code, n.Level, n.DisplayOrder, n.IsActive, n.IsLeaf, n.IsDeprecated,
				n.ReplacedByID, lat, lng, minLat, maxLat, minLng, maxLng,
				n.Aliases, n.Source, n.VerificationStatus,
			)
		}

		br := tx.SendBatch(ctx, batch)
		for j := 0; j < len(chunk); j++ {
			if _, err := br.Exec(); err != nil {
				br.Close()
				return fmt.Errorf("failed to upsert location chunk %d-%d: %w", i, end, err)
			}
		}
		if err := br.Close(); err != nil {
			return fmt.Errorf("failed closing location batch %d-%d: %w", i, end, err)
		}
	}

	return tx.Commit(ctx)
}

func (r *MasterDataRepo) UpsertCategoryBatch(ctx context.Context, nodes []models.MasterCategoryNode, adminID int) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Sort categories by Level ASC to guarantee parent categories are inserted before children
	sort.SliceStable(nodes, func(i, j int) bool {
		return nodes[i].Level < nodes[j].Level
	})

	upsertSQL := `
		INSERT INTO category_hierarchy (
			id, parent_id, classification_system_code, name_en, name_am, name_om, name_ti,
			slug, description, level, path, display_order, icon, image_url,
			is_active, is_leaf, is_archived, meta_title, meta_description, i18n, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7,
			$8, $9, $10, $11, $12, $13, $14,
			$15, $16, $17, $18, $19, $20, NOW()
		)
		ON CONFLICT (id) DO UPDATE SET
			parent_id = EXCLUDED.parent_id,
			classification_system_code = EXCLUDED.classification_system_code,
			name_en = EXCLUDED.name_en,
			name_am = EXCLUDED.name_am,
			name_om = EXCLUDED.name_om,
			name_ti = EXCLUDED.name_ti,
			slug = EXCLUDED.slug,
			description = EXCLUDED.description,
			level = EXCLUDED.level,
			path = EXCLUDED.path,
			display_order = EXCLUDED.display_order,
			icon = EXCLUDED.icon,
			image_url = EXCLUDED.image_url,
			is_active = EXCLUDED.is_active,
			is_leaf = EXCLUDED.is_leaf,
			is_archived = EXCLUDED.is_archived,
			meta_title = EXCLUDED.meta_title,
			meta_description = EXCLUDED.meta_description,
			i18n = EXCLUDED.i18n,
			updated_at = NOW();
	`

	batchSize := 500
	for i := 0; i < len(nodes); i += batchSize {
		end := i + batchSize
		if end > len(nodes) {
			end = len(nodes)
		}
		chunk := nodes[i:end]

		batch := &pgx.Batch{}
		for _, c := range chunk {
			i18nBytes, _ := json.Marshal(map[string]string{
				"en": c.Names.En,
				"am": c.Names.Am,
				"om": c.Names.Om,
				"ti": c.Names.Ti,
			})

			batch.Queue(upsertSQL,
				c.ID, c.ParentID, c.ClassificationSystemCode, c.Names.En, c.Names.Am, c.Names.Om, c.Names.Ti,
				c.Slug, c.Description, c.Level, c.Path, c.DisplayOrder, c.Icon, c.ImageURL,
				c.IsActive, c.IsLeaf, c.IsArchived, c.MetaTitle, c.MetaDescription, string(i18nBytes),
			)
		}

		br := tx.SendBatch(ctx, batch)
		for j := 0; j < len(chunk); j++ {
			if _, err := br.Exec(); err != nil {
				br.Close()
				return fmt.Errorf("failed to upsert category chunk %d-%d: %w", i, end, err)
			}
		}
		if err := br.Close(); err != nil {
			return fmt.Errorf("failed closing category batch %d-%d: %w", i, end, err)
		}
	}

	// Synchronize category search aliases
	aliasBatch := &pgx.Batch{}
	aliasCount := 0
	aliasUpsertSQL := `
		INSERT INTO category_aliases (category_id, alias, language)
		VALUES ($1, $2, 'en')
		ON CONFLICT (category_id, alias) DO NOTHING;
	`
	for _, c := range nodes {
		for _, a := range c.Aliases {
			aClean := strings.TrimSpace(a)
			if aClean != "" && !strings.EqualFold(aClean, c.Names.En) {
				aliasBatch.Queue(aliasUpsertSQL, c.ID, aClean)
				aliasCount++
				if aliasCount%500 == 0 {
					br := tx.SendBatch(ctx, aliasBatch)
					for k := 0; k < aliasBatch.Len(); k++ {
						if _, err := br.Exec(); err != nil {
							br.Close()
							return fmt.Errorf("failed to upsert category alias: %w", err)
						}
					}
					if err := br.Close(); err != nil {
						return fmt.Errorf("failed closing alias batch: %w", err)
					}
					aliasBatch = &pgx.Batch{}
				}
			}
		}
	}
	if aliasBatch.Len() > 0 {
		br := tx.SendBatch(ctx, aliasBatch)
		for k := 0; k < aliasBatch.Len(); k++ {
			if _, err := br.Exec(); err != nil {
				br.Close()
				return fmt.Errorf("failed to upsert final category aliases: %w", err)
			}
		}
		if err := br.Close(); err != nil {
			return fmt.Errorf("failed closing final alias batch: %w", err)
		}
	}

	return tx.Commit(ctx)
}

func (r *MasterDataRepo) GetMasterDataStatus(ctx context.Context) (*models.MasterDataStatus, error) {
	status := &models.MasterDataStatus{
		SupportedLangs: []string{"en", "am", "om", "ti"},
		LastImportAt:   time.Now(),
	}

	var locCount, activeLocs int
	err := r.db.QueryRow(ctx, `SELECT count(*), count(CASE WHEN is_active THEN 1 END) FROM location_hierarchy;`).Scan(&locCount, &activeLocs)
	if err != nil {
		return nil, err
	}
	status.LocationsCount = locCount
	status.ActiveLocations = activeLocs

	var catCount, activeCats int
	err = r.db.QueryRow(ctx, `SELECT count(*), count(CASE WHEN is_active THEN 1 END) FROM category_hierarchy WHERE is_archived = false;`).Scan(&catCount, &activeCats)
	if err != nil {
		return nil, err
	}
	status.CategoriesCount = catCount
	status.ActiveCats = activeCats

	return status, nil
}

func (r *MasterDataRepo) LogMasterDataAudit(ctx context.Context, adminID int, action string, catID *int, oldVal, newVal interface{}, reason, ip string) error {
	var oldBytes, newBytes []byte
	if oldVal != nil {
		oldBytes, _ = json.Marshal(oldVal)
	}
	if newVal != nil {
		newBytes, _ = json.Marshal(newVal)
	}

	var adminParam interface{} = adminID
	if adminID <= 0 {
		adminParam = nil
	}

	query := `
		INSERT INTO taxonomy_audit_logs (admin_id, category_id, action, old_values, new_values, reason, ip_address, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW());
	`
	_, err := r.db.Exec(ctx, query, adminParam, catID, action, oldBytes, newBytes, reason, ip)
	return err
}
