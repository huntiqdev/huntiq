package cache

import (
	"context"
	"testing"
	"time"
)

func TestMemoryCache_SetGetDelete(t *testing.T) {
	c := NewMemoryCache(100 * time.Millisecond)
	defer c.Close()

	ctx := context.Background()
	key := "test:key:1"
	val := []byte("hello-huntiq")

	// 1. Get non-existent
	if _, ok := c.Get(ctx, key); ok {
		t.Fatalf("expected key not found")
	}

	// 2. Set with TTL
	if err := c.Set(ctx, key, val, 1*time.Second); err != nil {
		t.Fatalf("unexpected error setting key: %v", err)
	}

	// 3. Get existing
	got, ok := c.Get(ctx, key)
	if !ok || string(got) != string(val) {
		t.Fatalf("expected %s, got %s (ok=%v)", string(val), string(got), ok)
	}

	// 4. Delete
	if err := c.Delete(ctx, key); err != nil {
		t.Fatalf("unexpected error deleting key: %v", err)
	}

	if _, ok := c.Get(ctx, key); ok {
		t.Fatalf("expected key deleted")
	}
}

func TestMemoryCache_Expiration(t *testing.T) {
	c := NewMemoryCache(50 * time.Millisecond)
	defer c.Close()

	ctx := context.Background()
	key := "test:expiring"
	val := []byte("temporary")

	if err := c.Set(ctx, key, val, 60*time.Millisecond); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should exist immediately
	if _, ok := c.Get(ctx, key); !ok {
		t.Fatalf("expected key to exist immediately")
	}

	// Wait for expiration
	time.Sleep(100 * time.Millisecond)

	if _, ok := c.Get(ctx, key); ok {
		t.Fatalf("expected key to be expired")
	}
}

func TestCache_FactoryFallbackToMemory(t *testing.T) {
	// Provide unreachable address, must gracefully fall back to MemoryCache without error
	c := New("127.0.0.1:9999", "", 0)
	defer c.Close()

	ctx := context.Background()
	if err := c.Set(ctx, "fallback:key", []byte("val"), 1*time.Second); err != nil {
		t.Fatalf("fallback cache failed: %v", err)
	}

	val, ok := c.Get(ctx, "fallback:key")
	if !ok || string(val) != "val" {
		t.Fatalf("fallback cache get failed: %s, ok=%v", string(val), ok)
	}
}
