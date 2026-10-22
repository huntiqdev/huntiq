package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type LocationRepo struct{ db *pgxpool.Pool }

func NewLocationRepo(db *pgxpool.Pool) *LocationRepo { return &LocationRepo{db} }

const sqlListRegions = `
	SELECT DISTINCT region,
	        MAX(region_am) AS region_am,
	        MAX(region_or) AS region_or,
	        MAX(region_ti) AS region_ti
	 FROM locations
	 GROUP BY region
	 ORDER BY region`

const sqlListZones = `
	SELECT DISTINCT zone,
	        MAX(zone_am) AS zone_am,
	        MAX(zone_or) AS zone_or,
	        MAX(zone_ti) AS zone_ti
	 FROM locations
	 WHERE region = $1 AND zone IS NOT NULL
	 GROUP BY zone
	 ORDER BY zone`

const sqlListSubcitiesByZone = `
	SELECT DISTINCT subcity,
	        MAX(subcity_am) AS subcity_am,
	        MAX(subcity_or) AS subcity_or,
	        MAX(subcity_ti) AS subcity_ti
	 FROM locations
	 WHERE region = $1 AND zone = $2 AND subcity IS NOT NULL
	 GROUP BY subcity ORDER BY subcity`

const sqlListSubcitiesByRegion = `
	SELECT DISTINCT subcity,
	        MAX(subcity_am) AS subcity_am,
	        MAX(subcity_or) AS subcity_or,
	        MAX(subcity_ti) AS subcity_ti
	 FROM locations
	 WHERE region = $1 AND subcity IS NOT NULL
	 GROUP BY subcity ORDER BY subcity`

const sqlListWoredasByZoneAndSubcity = `
	SELECT DISTINCT woreda,
	        MAX(woreda_am) AS woreda_am,
	        MAX(woreda_or) AS woreda_or,
	        MAX(woreda_ti) AS woreda_ti
	 FROM locations
	 WHERE region = $1 AND zone = $2 AND subcity = $3 AND woreda IS NOT NULL
	 GROUP BY woreda ORDER BY woreda`

const sqlListWoredasByZone = `
	SELECT DISTINCT woreda,
	        MAX(woreda_am) AS woreda_am,
	        MAX(woreda_or) AS woreda_or,
	        MAX(woreda_ti) AS woreda_ti
	 FROM locations
	 WHERE region = $1 AND zone = $2 AND woreda IS NOT NULL
	 GROUP BY woreda ORDER BY woreda`

const sqlListWoredasBySubcity = `
	SELECT DISTINCT woreda,
	        MAX(woreda_am) AS woreda_am,
	        MAX(woreda_or) AS woreda_or,
	        MAX(woreda_ti) AS woreda_ti
	 FROM locations
	 WHERE region = $1 AND subcity = $2 AND woreda IS NOT NULL
	 GROUP BY woreda ORDER BY woreda`

const sqlSearchLocations = `
	WITH scored AS (
	  SELECT DISTINCT ON (region, COALESCE(zone,''), COALESCE(subcity,''), COALESCE(woreda,''))
	    region, zone, subcity, woreda,
	    region_am, region_or, region_ti,
	    zone_am, zone_or, zone_ti,
	    subcity_am, subcity_or, subcity_ti,
	    woreda_am, woreda_or, woreda_ti,
	    GREATEST(
	      similarity(LOWER(region),        LOWER($1)),
	      similarity(LOWER(COALESCE(region_am,'')), LOWER($1)),
	      similarity(LOWER(COALESCE(region_or,'')), LOWER($1)),
	      similarity(LOWER(COALESCE(region_ti,'')), LOWER($1)),
	      similarity(LOWER(COALESCE(zone,'')),      LOWER($1)),
	      similarity(LOWER(COALESCE(zone_am,'')),   LOWER($1)),
	      similarity(LOWER(COALESCE(zone_or,'')),   LOWER($1)),
	      similarity(LOWER(COALESCE(zone_ti,'')),   LOWER($1)),
	      similarity(LOWER(COALESCE(subcity,'')),   LOWER($1)),
	      similarity(LOWER(COALESCE(subcity_am,'')),LOWER($1)),
	      similarity(LOWER(COALESCE(subcity_or,'')),LOWER($1)),
	      similarity(LOWER(COALESCE(subcity_ti,'')),LOWER($1)),
	      similarity(LOWER(COALESCE(woreda,'')),    LOWER($1)),
	      similarity(LOWER(COALESCE(woreda_am,'')), LOWER($1)),
	      similarity(LOWER(COALESCE(woreda_or,'')), LOWER($1)),
	      similarity(LOWER(COALESCE(woreda_ti,'')), LOWER($1))
	    ) AS score
	  FROM locations
	  WHERE
	    region        ILIKE $2 OR region_am  ILIKE $2 OR region_or  ILIKE $2 OR region_ti  ILIKE $2 OR
	    zone          ILIKE $2 OR zone_am    ILIKE $2 OR zone_or    ILIKE $2 OR zone_ti    ILIKE $2 OR
	    subcity       ILIKE $2 OR subcity_am ILIKE $2 OR subcity_or ILIKE $2 OR subcity_ti ILIKE $2 OR
	    woreda        ILIKE $2 OR woreda_am  ILIKE $2 OR woreda_or  ILIKE $2 OR woreda_ti  ILIKE $2 OR
	    similarity(LOWER(region),        LOWER($1)) > 0.2 OR
	    similarity(LOWER(COALESCE(subcity,'')),   LOWER($1)) > 0.2 OR
	    similarity(LOWER(COALESCE(subcity_am,'')),LOWER($1)) > 0.2 OR
	    similarity(LOWER(COALESCE(woreda,'')),    LOWER($1)) > 0.2
	)
	SELECT region, zone, subcity, woreda,
	       region_am, region_or, region_ti,
	       zone_am, zone_or, zone_ti,
	       subcity_am, subcity_or, subcity_ti,
	       woreda_am, woreda_or, woreda_ti
	FROM scored
	ORDER BY score DESC
	LIMIT 15`

const sqlResolveWoreda  = `SELECT id FROM locations WHERE region ILIKE $1 AND zone ILIKE $2 AND subcity ILIKE $3 AND woreda ILIKE $4 LIMIT 1`
const sqlResolveSubcity = `SELECT id FROM locations WHERE region ILIKE $1 AND zone ILIKE $2 AND subcity ILIKE $3 LIMIT 1`
const sqlResolveZone    = `SELECT id FROM locations WHERE region ILIKE $1 AND zone ILIKE $2 LIMIT 1`
const sqlResolveRegion  = `SELECT id FROM locations WHERE region ILIKE $1 LIMIT 1`

// ListRegions returns distinct regions with all language variants.
func (r *LocationRepo) ListRegions(ctx context.Context) ([]map[string]interface{}, error) {
	//nolint:gosec // G201/G202: static SQL with no parameters, no user input in SQL structure
	rows, err := r.db.Query(ctx, sqlListRegions)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []map[string]interface{}
	for rows.Next() {
		var en string
		var am, or_, ti *string
		if err := rows.Scan(&en, &am, &or_, &ti); err != nil {
			return nil, err
		}
		out = append(out, map[string]interface{}{
			"en": en, "am": am, "or": or_, "ti": ti,
		})
	}
	return out, rows.Err()
}

// ListZones returns distinct zones for a region with all language variants.
func (r *LocationRepo) ListZones(ctx context.Context, region string) ([]map[string]interface{}, error) {
	//nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
	rows, err := r.db.Query(ctx, sqlListZones, region)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanI18nRows(rows)
}

// ListSubcities returns distinct subcities with all language variants.
func (r *LocationRepo) ListSubcities(ctx context.Context, region, zone string) ([]map[string]interface{}, error) {
	var rows interface {
		Next() bool
		Scan(...interface{}) error
		Err() error
		Close()
	}
	var err error
	if zone != "" {
		//nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
		rows, err = r.db.Query(ctx, sqlListSubcitiesByZone, region, zone)
	} else {
		//nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
		rows, err = r.db.Query(ctx, sqlListSubcitiesByRegion, region)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanI18nRows(rows)
}

// ListWoredas returns distinct woredas with all language variants.
// When subcity is empty but zone is provided, returns all woredas in that zone.
func (r *LocationRepo) ListWoredas(ctx context.Context, region, zone, subcity string) ([]map[string]interface{}, error) {
	var rows interface {
		Next() bool
		Scan(...interface{}) error
		Err() error
		Close()
	}
	var err error
	switch {
	case zone != "" && subcity != "":
		//nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
		rows, err = r.db.Query(ctx, sqlListWoredasByZoneAndSubcity, region, zone, subcity)
	case zone != "" && subcity == "":
		//nolint:gosec // G201/G202: all values passed as $N par
		rows, err = r.db.Query(ctx, sqlListWoredasBySubcity, region, subcity)
	default:
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanI18nRows(rows)
}

// SearchLocations uses trigram similarity + prefix + ILIKE across all language columns.
// Returns results ranked by best match score descending.
func (r *LocationRepo) SearchLocations(ctx context.Context, q string) ([]map[string]interface{}, error) {
	if q == "" {
		return nil, nil
	}
	like := "%" + escapeLike(q) + "%"
	rows, err := r.db.Query(ctx, sqlSearchLocations, q, like)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []map[string]interface{}
	for rows.Next() {
		var region string
		var zone, subcity, woreda *string
		var regionAm, regionOr, regionTi *string
		var zoneAm, zoneOr, zoneTi *string
		var subcityAm, subcityOr, subcityTi *string
		var woredaAm, woredaOr, woredaTi *string
		if err := rows.Scan(
			&region, &zone, &subcity, &woreda,
			&regionAm, &regionOr, &regionTi,
			&zoneAm, &zoneOr, &zoneTi,
			&subcityAm, &subcityOr, &subcityTi,
			&woredaAm, &woredaOr, &woredaTi,
		); err != nil {
			return nil, err
		}
		out = append(out, map[string]interface{}{
			"region": region, "zone": zone, "subcity": subcity, "woreda": woreda,
			"region_am": regionAm, "region_or": regionOr, "region_ti": regionTi,
			"zone_am": zoneAm, "zone_or": zoneOr, "zone_ti": zoneTi,
			"subcity_am": subcityAm, "subcity_or": subcityOr, "subcity_ti": subcityTi,
			"woreda_am": woredaAm, "woreda_or": woredaOr, "woreda_ti": woredaTi,
		})
	}
	return out, rows.Err()
}

// ResolveID finds the best-matching location_id for the given hierarchy.
// It tries the most specific match first, falling back up the hierarchy.
func (r *LocationRepo) ResolveID(ctx context.Context, region, zone, subcity, woreda string) (int, error) {
	var id int
	type query struct {
		sql  string
		args []interface{}
		skip bool
	}
	queries := []query{
		{sql: sqlResolveWoreda, args: []interface{}{region, zone, subcity, woreda}, skip: region == "" || zone == "" || subcity == "" || woreda == ""},
		{sql: sqlResolveSubcity, args: []interface{}{region, zone, subcity}, skip: region == "" || zone == "" || subcity == ""},
		{sql: sqlResolveZone, args: []interface{}{region, zone}, skip: region == "" || zone == ""},
		{sql: sqlResolveRegion, args: []interface{}{region}, skip: region == ""},
	}
	for _, q := range queries {
		if q.skip {
			continue
		}
		if err := r.db.QueryRow(ctx, q.sql, q.args...).Scan(&id); err == nil {
			return id, nil
		}
	}
	return 0, fmt.Errorf("location not found for region=%q", region)
}

func scanI18nRows(rows interface {
	Next() bool
	Scan(...interface{}) error
	Err() error
}) ([]map[string]interface{}, error) {
	var out []map[string]interface{}
	for rows.Next() {
		var en string
		var am, or_, ti *string
		if err := rows.Scan(&en, &am, &or_, &ti); err != nil {
			return nil, err
		}
		out = append(out, map[string]interface{}{
			"en": en, "am": am, "or": or_, "ti": ti,
		})
	}
	return out, rows.Err()
}
