package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/your-org/used-tools-marketplace/internal/models"
)

type EntityRepo struct {
	db *pgxpool.Pool
}

func NewEntityRepo(db *pgxpool.Pool) *EntityRepo {
	return &EntityRepo{db: db}
}

// ResolveCanonicalEntity resolves user query text into canonical entities via multi-stage fuzzy/phonetic/trigram matching
func (r *EntityRepo) ResolveCanonicalEntity(ctx context.Context, term string, lang string, limit int) ([]models.EntityResolutionResult, error) {
	if term == "" {
		return []models.EntityResolutionResult{}, nil
	}
	if limit <= 0 || limit > 20 {
		limit = 5
	}
	if lang == "" {
		lang = "en"
	}

	query := `
		SELECT 
			out_entity_id,
			out_canonical_name_en,
			COALESCE(out_canonical_name_am, ''),
			out_canonical_slug,
			out_entity_nature,
			out_entity_type,
			out_short_definition_en,
			COALESCE(out_short_definition_am, ''),
			out_primary_category_id,
			out_score,
			out_confidence_tier,
			out_match_method,
			COALESCE(out_suggested_correction, '')
		FROM public.resolve_canonical_entity($1, $2, $3)
	`

	rows, err := r.db.Query(ctx, query, term, lang, limit)
	if err != nil {
		return nil, fmt.Errorf("resolve canonical entity query failed: %w", err)
	}
	defer rows.Close()

	var results []models.EntityResolutionResult
	for rows.Next() {
		var res models.EntityResolutionResult
		var primaryCatID sql.NullInt32

		err := rows.Scan(
			&res.EntityID,
			&res.CanonicalNameEn,
			&res.CanonicalNameAm,
			&res.CanonicalSlug,
			&res.EntityNature,
			&res.EntityType,
			&res.ShortDefinitionEn,
			&res.ShortDefinitionAm,
			&primaryCatID,
			&res.Score,
			&res.ConfidenceTier,
			&res.MatchMethod,
			&res.SuggestedCorrection,
		)
		if err != nil {
			return nil, fmt.Errorf("scan entity resolution row failed: %w", err)
		}

		if primaryCatID.Valid {
			idVal := int(primaryCatID.Int32)
			res.PrimaryCategoryID = &idVal
		}

		results = append(results, res)
	}

	return results, nil
}

// GetEntityByID returns the full canonical entity specification by ID
func (r *EntityRepo) GetEntityByID(ctx context.Context, id int) (*models.CanonicalEntity, error) {
	query := `
		SELECT 
			id,
			canonical_name_en,
			COALESCE(canonical_name_am, ''),
			COALESCE(canonical_name_om, ''),
			COALESCE(canonical_name_ti, ''),
			canonical_slug,
			entity_nature,
			entity_type,
			short_definition_en,
			COALESCE(short_definition_am, ''),
			COALESCE(short_definition_om, ''),
			COALESCE(short_definition_ti, ''),
			COALESCE(detailed_definition_en, ''),
			COALESCE(detailed_definition_am, ''),
			typical_examples,
			typical_use_cases,
			functional_capabilities,
			relevant_attributes,
			disambiguation_rules,
			primary_category_id,
			parent_entity_id,
			status,
			created_at,
			updated_at
		FROM public.canonical_entities
		WHERE id = $1
	`

	row := r.db.QueryRow(ctx, query, id)

	var e models.CanonicalEntity
	var primCat, parentEnt sql.NullInt32
	var exBytes, useBytes, capBytes, attrBytes, disBytes []byte

	err := row.Scan(
		&e.ID,
		&e.CanonicalNameEn,
		&e.CanonicalNameAm,
		&e.CanonicalNameOm,
		&e.CanonicalNameTi,
		&e.CanonicalSlug,
		&e.EntityNature,
		&e.EntityType,
		&e.ShortDefinitionEn,
		&e.ShortDefinitionAm,
		&e.ShortDefinitionOm,
		&e.ShortDefinitionTi,
		&e.DetailedDefinitionEn,
		&e.DetailedDefinitionAm,
		&exBytes,
		&useBytes,
		&capBytes,
		&attrBytes,
		&disBytes,
		&primCat,
		&parentEnt,
		&e.Status,
		&e.CreatedAt,
		&e.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("get entity by id %d: %w", id, err)
	}

	if primCat.Valid {
		idVal := int(primCat.Int32)
		e.PrimaryCategoryID = &idVal
	}
	if parentEnt.Valid {
		idVal := int(parentEnt.Int32)
		e.ParentEntityID = &idVal
	}

	e.TypicalExamples = json.RawMessage(exBytes)
	e.TypicalUseCases = json.RawMessage(useBytes)
	e.FunctionalCapabilities = json.RawMessage(capBytes)
	e.RelevantAttributes = json.RawMessage(attrBytes)
	e.DisambiguationRules = json.RawMessage(disBytes)

	// Fetch discovery paths & relationships
	paths, _ := r.GetEntityDiscoveryPaths(ctx, e.ID)
	e.DiscoveryPaths = paths

	rels, _ := r.GetEntityRelationships(ctx, e.ID)
	e.Relationships = rels

	aliases, _ := r.GetEntityAliases(ctx, e.ID)
	e.Aliases = aliases

	return &e, nil
}

// GetEntityBySlug returns the canonical entity by unique slug
func (r *EntityRepo) GetEntityBySlug(ctx context.Context, slug string) (*models.CanonicalEntity, error) {
	var id int
	err := r.db.QueryRow(ctx, "SELECT id FROM public.canonical_entities WHERE canonical_slug = $1", slug).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("entity slug %q not found: %w", slug, err)
	}
	return r.GetEntityByID(ctx, id)
}

// GetEntityDiscoveryPaths retrieves multi-domain discovery paths for an entity
func (r *EntityRepo) GetEntityDiscoveryPaths(ctx context.Context, entityID int) ([]models.EntityDiscoveryPath, error) {
	query := `
		SELECT 
			out_path_type,
			out_target_id,
			out_target_name,
			out_relationship_type,
			out_depth,
			out_confidence
		FROM public.get_entity_discovery_paths($1)
	`
	rows, err := r.db.Query(ctx, query, entityID)
	if err != nil {
		return nil, fmt.Errorf("get entity discovery paths failed: %w", err)
	}
	defer rows.Close()

	var paths []models.EntityDiscoveryPath
	for rows.Next() {
		var p models.EntityDiscoveryPath
		if err := rows.Scan(&p.PathType, &p.TargetID, &p.TargetName, &p.RelationshipType, &p.Depth, &p.Confidence); err != nil {
			return nil, err
		}
		paths = append(paths, p)
	}
	return paths, nil
}

// GetEntityRelationships returns the relationship edges originating from this entity
func (r *EntityRepo) GetEntityRelationships(ctx context.Context, entityID int) ([]models.EntityRelationship, error) {
	query := `
		SELECT 
			id,
			source_entity_id,
			target_entity_id,
			target_category_id,
			relationship_type,
			confidence_score,
			metadata,
			COALESCE(notes, ''),
			status,
			created_at
		FROM public.entity_relationships
		WHERE source_entity_id = $1 AND status = 'ACTIVE'
		ORDER BY confidence_score DESC, id ASC
	`
	rows, err := r.db.Query(ctx, query, entityID)
	if err != nil {
		return nil, fmt.Errorf("get entity relationships failed: %w", err)
	}
	defer rows.Close()

	var rels []models.EntityRelationship
	for rows.Next() {
		var rel models.EntityRelationship
		var targetEnt, targetCat sql.NullInt32
		var metaBytes []byte

		err := rows.Scan(
			&rel.ID,
			&rel.SourceEntityID,
			&targetEnt,
			&targetCat,
			&rel.RelationshipType,
			&rel.ConfidenceScore,
			&metaBytes,
			&rel.Notes,
			&rel.Status,
			&rel.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		if targetEnt.Valid {
			idVal := int(targetEnt.Int32)
			rel.TargetEntityID = &idVal
		}
		if targetCat.Valid {
			idVal := int(targetCat.Int32)
			rel.TargetCategoryID = &idVal
		}
		rel.Metadata = json.RawMessage(metaBytes)
		rels = append(rels, rel)
	}
	return rels, nil
}

// GetEntityAliases returns the registered aliases and spelling variants for this entity
func (r *EntityRepo) GetEntityAliases(ctx context.Context, entityID int) ([]models.EntityAlias, error) {
	query := `
		SELECT 
			id,
			entity_id,
			alias_text,
			normalized_alias,
			language,
			alias_type,
			confidence_score,
			status,
			created_at
		FROM public.entity_aliases
		WHERE entity_id = $1 AND status = 'ACTIVE'
		ORDER BY confidence_score DESC, alias_type ASC
	`
	rows, err := r.db.Query(ctx, query, entityID)
	if err != nil {
		return nil, fmt.Errorf("get entity aliases failed: %w", err)
	}
	defer rows.Close()

	var aliases []models.EntityAlias
	for rows.Next() {
		var a models.EntityAlias
		err := rows.Scan(
			&a.ID,
			&a.EntityID,
			&a.AliasText,
			&a.NormalizedAlias,
			&a.Language,
			&a.AliasType,
			&a.ConfidenceScore,
			&a.Status,
			&a.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		aliases = append(aliases, a)
	}
	return aliases, nil
}
