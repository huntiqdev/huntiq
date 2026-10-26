package repository

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SearchRepo struct {
	db *pgxpool.Pool
}

func NewSearchRepo(db *pgxpool.Pool) *SearchRepo {
	return &SearchRepo{db: db}
}

// LogSearch records a search query for analytics
func (r *SearchRepo) LogSearch(ctx context.Context, query string, resultCount int, categoryID *int, filters map[string]interface{}, userID *int, sessionID, ipAddress string) error {
	filtersJSON, _ := json.Marshal(filters)
	
	normalizedQuery := ""
	err := r.db.QueryRow(ctx, `SELECT normalize_query($1)`, query).Scan(&normalizedQuery)
	if err != nil {
		normalizedQuery = query
	}
	
	_, err = r.db.Exec(ctx, `
		INSERT INTO search_analytics (
			query, normalized_query, result_count, category_id,
			filters_used, user_id, session_id, ip_address
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		query, normalizedQuery, resultCount, categoryID,
		filtersJSON, userID, sessionID, ipAddress,
	)
	
	return err
}

// LogClick records when a user clicks a search result
func (r *SearchRepo) LogClick(ctx context.Context, query string, productID int, position int, userID *int, sessionID string) error {
	normalizedQuery := ""
	_ = r.db.QueryRow(ctx, `SELECT normalize_query($1)`, query).Scan(&normalizedQuery)
	if normalizedQuery == "" {
		normalizedQuery = query
	}
	
	_, err := r.db.Exec(ctx, `
		UPDATE search_analytics
		SET clicked_product_id = $1, clicked_position = $2
		WHERE id = (
			SELECT id FROM search_analytics
			WHERE normalized_query = $3
			  AND user_id IS NOT DISTINCT FROM $4
			  AND session_id = $5
			  AND created_at >= NOW() - INTERVAL '1 hour'
			ORDER BY created_at DESC
			LIMIT 1
		)`,
		productID, position, normalizedQuery, userID, sessionID,
	)
	
	return err
}

// GetPopularSearches returns trending search queries
func (r *SearchRepo) GetPopularSearches(ctx context.Context, limit int) ([]map[string]interface{}, error) {
	rows, err := r.db.Query(ctx, `
		SELECT normalized_query, search_count, avg_results, clicks, ctr_percent
		FROM popular_searches
		LIMIT $1`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var results []map[string]interface{}
	for rows.Next() {
		var query string
		var searchCount, clicks int
		var avgResults, ctr float64
		
		if err := rows.Scan(&query, &searchCount, &avgResults, &clicks, &ctr); err != nil {
			continue
		}
		
		results = append(results, map[string]interface{}{
			"query":        query,
			"search_count": searchCount,
			"avg_results":  avgResults,
			"clicks":       clicks,
			"ctr_percent":  ctr,
		})
	}
	
	return results, rows.Err()
}

// GetZeroResultQueries returns queries that found no results
func (r *SearchRepo) GetZeroResultQueries(ctx context.Context, limit int) ([]map[string]interface{}, error) {
	rows, err := r.db.Query(ctx, `
		SELECT query, normalized_query, occurrence_count
		FROM zero_result_queries
		LIMIT $1`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var results []map[string]interface{}
	for rows.Next() {
		var query, normalized string
		var count int
		
		if err := rows.Scan(&query, &normalized, &count); err != nil {
			continue
		}
		
		results = append(results, map[string]interface{}{
			"query":            query,
			"normalized_query": normalized,
			"occurrence_count": count,
		})
	}
	
	return results, rows.Err()
}
