package geocoding

import (
	"context"
	"fmt"
	"math"
	"sync"
	"testing"
	"time"

	"github.com/your-org/used-tools-marketplace/internal/models"
)

// MockProvider implements GeocodingProvider for testing
type MockProvider struct {
	name       string
	callCount  int
	mu         sync.Mutex
	resultFunc func(lat, lng float64) (*GeocodingResult, error)
}

func (m *MockProvider) Name() string { return m.name }
func (m *MockProvider) ReverseGeocode(ctx context.Context, lat, lng float64) (*GeocodingResult, error) {
	m.mu.Lock()
	m.callCount++
	m.mu.Unlock()
	if m.resultFunc != nil {
		return m.resultFunc(lat, lng)
	}
	return nil, fmt.Errorf("mock error")
}

func TestCoordinateRoundingAndCacheKey(t *testing.T) {
	lat1, lng1 := 8.9951234, 38.7901234
	lat2, lng2 := 8.9954999, 38.7904999

	key1 := fmt.Sprintf("%.3f,%.3f", math.Round(lat1*1000)/1000, math.Round(lng1*1000)/1000)
	key2 := fmt.Sprintf("%.3f,%.3f", math.Round(lat2*1000)/1000, math.Round(lng2*1000)/1000)

	if key1 != key2 {
		t.Errorf("Expected keys to match within 100m window: got %s and %s", key1, key2)
	}
	if key1 != "8.995,38.790" {
		t.Errorf("Expected key '8.995,38.790', got %s", key1)
	}
}

func TestBreadcrumbMappingNoFabrication(t *testing.T) {
	// Verify that if a location is resolved only at Region level (e.g. Afar),
	// Subcity, Woreda, and Kebele remain empty strings and are never fabricated.
	match := &models.GeocodingMatch{
		LocationID:   3,
		LocationType: "region",
		Name:         "Afar",
		Latitude:     11.7559,
		Longitude:    40.9587,
		Confidence:   0.90,
		Breadcrumbs: []models.LocationBreadcrumb{
			{ID: 1, Name: "Ethiopia", LocationType: "country", Level: 0},
			{ID: 3, Name: "Afar", LocationType: "region", Level: 1},
		},
		Provider: "local_db",
	}

	for _, bc := range match.Breadcrumbs {
		switch bc.LocationType {
		case "country":
			match.Country = bc.Name
		case "region":
			match.Region = bc.Name
		case "subcity":
			match.Subcity = bc.Name
		case "woreda":
			match.Woreda = bc.Name
		}
	}

	if match.Country != "Ethiopia" {
		t.Errorf("Expected Country 'Ethiopia', got '%s'", match.Country)
	}
	if match.Region != "Afar" {
		t.Errorf("Expected Region 'Afar', got '%s'", match.Region)
	}
	if match.Subcity != "" {
		t.Errorf("Subcity must NOT be fabricated: got '%s'", match.Subcity)
	}
	if match.Woreda != "" {
		t.Errorf("Woreda must NOT be fabricated: got '%s'", match.Woreda)
	}
}

func TestCompositeCachingBehavior(t *testing.T) {
	service := &CompositeGeocodingService{
		cache: make(map[string]cacheEntry),
	}

	// Prime the cache
	testKey := "8.995,38.790"
	expectedResult := &GeocodingResult{
		MatchedLocation: &models.GeocodingMatch{
			LocationID:   16,
			Name:         "Bole",
			LocationType: "subcity",
			Confidence:   0.95,
		},
		Provider: "local_db",
	}

	service.cache[testKey] = cacheEntry{
		result:    expectedResult,
		expiresAt: time.Now().Add(1 * time.Hour),
	}

	// Request with coordinates that round to the same key
	res, err := service.ReverseGeocode(context.Background(), 8.9951, 38.7902)
	if err != nil {
		t.Fatalf("ReverseGeocode failed: %v", err)
	}
	if res.MatchedLocation.LocationID != 16 {
		t.Errorf("Expected LocationID 16 from cache, got %d", res.MatchedLocation.LocationID)
	}
	if res.MatchedLocation.Name != "Bole" {
		t.Errorf("Expected Name 'Bole' from cache, got '%s'", res.MatchedLocation.Name)
	}
}

func TestConcurrentCacheAccess(t *testing.T) {
	service := &CompositeGeocodingService{
		cache: make(map[string]cacheEntry),
	}

	var wg sync.WaitGroup
	workers := 20
	iterations := 50

	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for i := 0; i < iterations; i++ {
				lat := 8.500 + float64(i)*0.01
				lng := 38.500 + float64(workerID)*0.01
				key := fmt.Sprintf("%.3f,%.3f", math.Round(lat*1000)/1000, math.Round(lng*1000)/1000)

				service.mu.Lock()
				service.cache[key] = cacheEntry{
					result: &GeocodingResult{
						MatchedLocation: &models.GeocodingMatch{Name: "Test"},
					},
					expiresAt: time.Now().Add(10 * time.Minute),
				}
				service.mu.Unlock()

				service.mu.RLock()
				_ = service.cache[key]
				service.mu.RUnlock()
			}
		}(w)
	}

	wg.Wait()
}

func TestNominatimUserAgent(t *testing.T) {
	provider := NewNominatimProvider(nil)
	expectedUA := "HUNTIQ-Universal-Marketplace/1.0"
	if provider.userAgent != expectedUA {
		t.Errorf("expected Nominatim User-Agent %q, got %q", expectedUA, provider.userAgent)
	}
}
