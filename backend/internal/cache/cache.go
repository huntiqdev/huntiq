package cache

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

// Cache defines the contract for key-value caching with TTL.
type Cache interface {
	Get(ctx context.Context, key string) ([]byte, bool)
	Set(ctx context.Context, key string, val []byte, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
	Close() error
}

// MemoryCache is a high-performance in-memory cache with concurrent read/write protection and TTL eviction.
type MemoryCache struct {
	mu    sync.RWMutex
	items map[string]memItem
	stop  chan struct{}
}

type memItem struct {
	value     []byte
	expiresAt time.Time
}

// NewMemoryCache initializes an in-memory cache with periodic background eviction.
func NewMemoryCache(cleanupInterval time.Duration) *MemoryCache {
	mc := &MemoryCache{
		items: make(map[string]memItem),
		stop:  make(chan struct{}),
	}
	go mc.startCleanup(cleanupInterval)
	return mc
}

func (m *MemoryCache) Get(ctx context.Context, key string) ([]byte, bool) {
	m.mu.RLock()
	item, found := m.items[key]
	m.mu.RUnlock()

	if !found {
		return nil, false
	}
	if !item.expiresAt.IsZero() && time.Now().After(item.expiresAt) {
		m.Delete(ctx, key)
		return nil, false
	}
	// Return a defensive copy to prevent external mutation
	copied := make([]byte, len(item.value))
	copy(copied, item.value)
	return copied, true
}

func (m *MemoryCache) Set(ctx context.Context, key string, val []byte, ttl time.Duration) error {
	var expiresAt time.Time
	if ttl > 0 {
		expiresAt = time.Now().Add(ttl)
	}
	copied := make([]byte, len(val))
	copy(copied, val)

	m.mu.Lock()
	m.items[key] = memItem{
		value:     copied,
		expiresAt: expiresAt,
	}
	m.mu.Unlock()
	return nil
}

func (m *MemoryCache) Delete(ctx context.Context, key string) error {
	m.mu.Lock()
	delete(m.items, key)
	m.mu.Unlock()
	return nil
}

func (m *MemoryCache) Close() error {
	close(m.stop)
	return nil
}

func (m *MemoryCache) startCleanup(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			now := time.Now()
			m.mu.Lock()
			for k, v := range m.items {
				if !v.expiresAt.IsZero() && now.After(v.expiresAt) {
					delete(m.items, k)
				}
			}
			m.mu.Unlock()
		case <-m.stop:
			return
		}
	}
}

// RedisCache wraps a Redis client with strict timeouts and non-blocking failure tolerance.
type RedisCache struct {
	client *redis.Client
}

func (r *RedisCache) Get(ctx context.Context, key string) ([]byte, bool) {
	reqCtx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	val, err := r.client.Get(reqCtx, key).Bytes()
	if err != nil {
		return nil, false
	}
	return val, true
}

func (r *RedisCache) Set(ctx context.Context, key string, val []byte, ttl time.Duration) error {
	reqCtx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	return r.client.Set(reqCtx, key, val, ttl).Err()
}

func (r *RedisCache) Delete(ctx context.Context, key string) error {
	reqCtx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	return r.client.Del(reqCtx, key).Err()
}

func (r *RedisCache) Close() error {
	return r.client.Close()
}

// New creates a resilient Cache. If redisAddr is provided and reachable, it uses Redis.
// Otherwise, it transparently falls back to a fast in-memory cache.
func New(redisAddr, redisPassword string, db int) Cache {
	if redisAddr != "" {
		rdb := redis.NewClient(&redis.Options{
			Addr:         redisAddr,
			Password:     redisPassword,
			DB:           db,
			DialTimeout:  200 * time.Millisecond,
			ReadTimeout:  100 * time.Millisecond,
			WriteTimeout: 100 * time.Millisecond,
			PoolSize:     50,
		})

		ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
		defer cancel()

		if err := rdb.Ping(ctx).Err(); err == nil {
			log.Printf("[CACHE] Connected to Redis cluster at %s", redisAddr)
			return &RedisCache{client: rdb}
		} else {
			log.Printf("[CACHE] Redis at %s unreachable (%v); falling back to in-memory cache", redisAddr, err)
			_ = rdb.Close()
		}
	}

	log.Printf("[CACHE] In-memory high-speed cache initialized (TTL eviction active)")
	return NewMemoryCache(1 * time.Minute)
}
