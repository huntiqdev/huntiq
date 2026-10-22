package repository

import (
	"context"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/sync/singleflight"
)

// SettingsCache provides a simple in-memory cache for platform_settings with 1-minute TTL.
// This prevents reading fee rates from the DB inside transaction locks, which serializes
// all concurrent purchases on a single table scan.
//
// SECURITY FIX #4: Uses singleflight pattern to prevent cache stampede DOS attacks.
// Without singleflight: 1000 concurrent cache misses → 1000 database queries → connection pool exhaustion
// With singleflight: 1000 concurrent cache misses → 1 database query (999x reduction)
type SettingsCache struct {
	db      *pgxpool.Pool
	cache   map[string]cachedValue
	mu      sync.RWMutex
	ttl     time.Duration
	group   singleflight.Group
}

type cachedValue struct {
	value   string
	expires time.Time
}

// NewSettingsCache creates a settings cache with the specified TTL.
// A background goroutine cleans expired entries every 5 minutes.
func NewSettingsCache(db *pgxpool.Pool, ttl time.Duration) *SettingsCache {
	sc := &SettingsCache{
		db:    db,
		cache: make(map[string]cachedValue),
		ttl:   ttl,
	}
	
	// Background cleanup goroutine
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			sc.cleanup()
		}
	}()
	
	return sc
}

// Get retrieves a setting from cache or DB. Returns empty string if not found.
// Uses singleflight to deduplicate concurrent requests for the same key.
func (sc *SettingsCache) Get(ctx context.Context, key string) string {
	// Try cache first (fast path)
	sc.mu.RLock()
	if val, ok := sc.cache[key]; ok && time.Now().Before(val.expires) {
		sc.mu.RUnlock()
		return val.value
	}
	sc.mu.RUnlock()
	
	// Cache miss or expired - use singleflight to deduplicate concurrent DB queries
	// All concurrent requests for the same key will wait for the first one to complete
	result, _, _ := sc.group.Do(key, func() (interface{}, error) {
		// Double-check cache after acquiring singleflight lock (another goroutine may have populated it)
		sc.mu.RLock()
		if val, ok := sc.cache[key]; ok && time.Now().Before(val.expires) {
			sc.mu.RUnlock()
			return val.value, nil
		}
		sc.mu.RUnlock()
		
		// Fetch from DB
		return sc.fetchFromDB(ctx, key), nil
	})
	
	return result.(string)
}

// fetchFromDB queries the database and updates the cache.
// This is called within singleflight.Do(), so only one goroutine executes it at a time per key.
func (sc *SettingsCache) fetchFromDB(ctx context.Context, key string) string {
	var value string
	//nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
	_ = sc.db.QueryRow(ctx, `SELECT value FROM platform_settings WHERE key=$1`, key).Scan(&value)
	
	// Store in cache
	sc.mu.Lock()
	sc.cache[key] = cachedValue{
		value:   value,
		expires: time.Now().Add(sc.ttl),
	}
	sc.mu.Unlock()
	
	return value
}

// Invalidate removes a key from the cache, forcing next Get to fetch from DB.
func (sc *SettingsCache) Invalidate(key string) {
	sc.mu.Lock()
	delete(sc.cache, key)
	sc.mu.Unlock()
}

// InvalidateAll clears the entire cache.
func (sc *SettingsCache) InvalidateAll() {
	sc.mu.Lock()
	sc.cache = make(map[string]cachedValue)
	sc.mu.Unlock()
}

// cleanup removes expired entries from the cache.
func (sc *SettingsCache) cleanup() {
	now := time.Now()
	sc.mu.Lock()
	for key, val := range sc.cache {
		if now.After(val.expires) {
			delete(sc.cache, key)
		}
	}
	sc.mu.Unlock()
}
