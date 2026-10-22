package geocoding

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/your-org/used-tools-marketplace/internal/models"
	"github.com/your-org/used-tools-marketplace/internal/repository"
)

// GeocodingResult represents the normalized result of a reverse geocoding operation
type GeocodingResult struct {
	MatchedLocation *models.GeocodingMatch `json:"matched_location"`
	RawAddress      string                 `json:"raw_address,omitempty"`
	Provider        string                 `json:"provider"`
}

// GeocodingProvider interface for pluggable geocoding backends
type GeocodingProvider interface {
	Name() string
	ReverseGeocode(ctx context.Context, lat, lng float64) (*GeocodingResult, error)
}

// LocalDBProvider matches coordinates directly against database bounding boxes and centroids
type LocalDBProvider struct {
	repo *repository.LocationHierarchyRepo
}

func NewLocalDBProvider(repo *repository.LocationHierarchyRepo) *LocalDBProvider {
	return &LocalDBProvider{repo: repo}
}

func (p *LocalDBProvider) Name() string { return "local_db" }

func (p *LocalDBProvider) ReverseGeocode(ctx context.Context, lat, lng float64) (*GeocodingResult, error) {
	loc, confidence, err := p.repo.MatchByCoordinates(ctx, lat, lng)
	if err != nil {
		return nil, err
	}

	match := &models.GeocodingMatch{
		LocationID:   loc.ID,
		LocationType: loc.LocationType,
		Name:         loc.NameEn,
		Latitude:     lat,
		Longitude:    lng,
		Confidence:   confidence,
		Breadcrumbs:  loc.Breadcrumbs,
		Provider:     "local_db",
	}

	// Populate specific administrative levels from breadcrumbs
	for _, bc := range loc.Breadcrumbs {
		switch bc.LocationType {
		case "country":
			match.Country = bc.Name
		case "region":
			match.Region = bc.Name
		case "zone":
			match.Zone = bc.Name
		case "city":
			if match.Region == "" {
				match.Region = bc.Name // e.g. Addis Ababa is both region and city
			}
			match.City = bc.Name
		case "subcity":
			match.Subcity = bc.Name
		case "woreda":
			match.Woreda = bc.Name
		}
	}

	// Fallback assignment if current node was not in breadcrumbs loop
	switch loc.LocationType {
	case "city":
		if match.City == "" {
			match.City = loc.NameEn
		}
	case "subcity":
		if match.Subcity == "" {
			match.Subcity = loc.NameEn
		}
	case "woreda":
		if match.Woreda == "" {
			match.Woreda = loc.NameEn
		}
	case "region":
		if match.Region == "" {
			match.Region = loc.NameEn
		}
	}

	var addressParts []string
	if match.Subcity != "" {
		addressParts = append(addressParts, match.Subcity)
	}
	if match.City != "" && match.City != match.Subcity {
		addressParts = append(addressParts, match.City)
	}
	if match.Zone != "" {
		addressParts = append(addressParts, match.Zone)
	}
	if match.Region != "" && match.Region != match.City {
		addressParts = append(addressParts, match.Region)
	}
	addressParts = append(addressParts, "Ethiopia")

	match.RawAddress = strings.Join(addressParts, ", ")

	return &GeocodingResult{
		MatchedLocation: match,
		RawAddress:      match.RawAddress,
		Provider:        "local_db",
	}, nil
}

// NominatimProvider implements OSM reverse geocoding with respectful usage limits
type NominatimProvider struct {
	client    *http.Client
	repo      *repository.LocationHierarchyRepo
	userAgent string
}

func NewNominatimProvider(repo *repository.LocationHierarchyRepo) *NominatimProvider {
	return &NominatimProvider{
		client: &http.Client{
			Timeout: 4 * time.Second,
		},
		repo:      repo,
		userAgent: "SELLBUY-Ethiopian-Marketplace/1.0",
	}
}

func (p *NominatimProvider) Name() string { return "nominatim" }

type osmResponse struct {
	DisplayName string `json:"display_name"`
	Address     struct {
		Subcity       string `json:"suburb"`
		City          string `json:"city"`
		Town          string `json:"town"`
		County        string `json:"county"`
		StateDistrict string `json:"state_district"`
		State         string `json:"state"`
		Country       string `json:"country"`
	} `json:"address"`
}

func (p *NominatimProvider) ReverseGeocode(ctx context.Context, lat, lng float64) (*GeocodingResult, error) {
	url := fmt.Sprintf("https://nominatim.openstreetmap.org/reverse?format=jsonv2&lat=%.6f&lon=%.6f&zoom=14", lat, lng)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", p.userAgent)
	req.Header.Set("Accept-Language", "en")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("nominatim returned status: %d", resp.StatusCode)
	}

	var osm osmResponse
	if err := json.NewDecoder(resp.Body).Decode(&osm); err != nil {
		return nil, err
	}

	// Try to match the suburb or city back to our verified database
	searchCandidates := []string{osm.Address.Subcity, osm.Address.City, osm.Address.Town, osm.Address.State}
	var matchedLoc *models.LocationHierarchy
	for _, cand := range searchCandidates {
		if cand == "" {
			continue
		}
		matches, err := p.repo.SearchLocations(ctx, cand, "en", 1)
		if err == nil && len(matches) > 0 {
			matchedLoc = &matches[0]
			break
		}
	}

	if matchedLoc == nil {
		return nil, fmt.Errorf("no verified location matched address: %s", osm.DisplayName)
	}

	match := &models.GeocodingMatch{
		LocationID:   matchedLoc.ID,
		LocationType: matchedLoc.LocationType,
		Name:         matchedLoc.NameEn,
		Latitude:     lat,
		Longitude:    lng,
		Confidence:   0.85,
		Breadcrumbs:  matchedLoc.Breadcrumbs,
		RawAddress:   osm.DisplayName,
		Provider:     "nominatim",
	}

	for _, bc := range matchedLoc.Breadcrumbs {
		switch bc.LocationType {
		case "region":
			match.Region = bc.Name
		case "zone":
			match.Zone = bc.Name
		case "city":
			match.City = bc.Name
		case "subcity":
			match.Subcity = bc.Name
		case "woreda":
			match.Woreda = bc.Name
		}
	}

	return &GeocodingResult{
		MatchedLocation: match,
		RawAddress:      osm.DisplayName,
		Provider:        "nominatim",
	}, nil
}

// CacheEntry stores a reverse geocode result with timestamp
type cacheEntry struct {
	result    *GeocodingResult
	expiresAt time.Time
}

// CompositeGeocodingService coordinates caching, local DB, and external fallbacks
type CompositeGeocodingService struct {
	localProvider *LocalDBProvider
	osmProvider   *NominatimProvider
	cache         map[string]cacheEntry
	mu            sync.RWMutex
}

func NewCompositeGeocodingService(repo *repository.LocationHierarchyRepo) *CompositeGeocodingService {
	return &CompositeGeocodingService{
		localProvider: NewLocalDBProvider(repo),
		osmProvider:   NewNominatimProvider(repo),
		cache:         make(map[string]cacheEntry),
	}
}

// ReverseGeocode resolves coordinates to verified location hierarchy with intelligent caching
func (s *CompositeGeocodingService) ReverseGeocode(ctx context.Context, lat, lng float64) (*GeocodingResult, error) {
	// Round coordinates to 3 decimals (~110 meters) for high cache hit rate and user privacy
	cacheKey := fmt.Sprintf("%.3f,%.3f", math.Round(lat*1000)/1000, math.Round(lng*1000)/1000)

	s.mu.RLock()
	if entry, ok := s.cache[cacheKey]; ok {
		if time.Now().Before(entry.expiresAt) {
			s.mu.RUnlock()
			return entry.result, nil
		}
	}
	s.mu.RUnlock()

	// 1. Try Local Database Provider (spatial bounding box & centroid match)
	res, err := s.localProvider.ReverseGeocode(ctx, lat, lng)
	if err == nil && res.MatchedLocation != nil && res.MatchedLocation.Confidence >= 0.80 {
		s.saveCache(cacheKey, res)
		return res, nil
	}

	// 2. If confidence was low or local lookup failed, try OSM Nominatim fallback
	if s.osmProvider != nil {
		osmRes, osmErr := s.osmProvider.ReverseGeocode(ctx, lat, lng)
		if osmErr == nil && osmRes.MatchedLocation != nil {
			s.saveCache(cacheKey, osmRes)
			return osmRes, nil
		}
	}

	// 3. If OSM failed but local had a partial match, return local match
	if res != nil && res.MatchedLocation != nil {
		s.saveCache(cacheKey, res)
		return res, nil
	}

	if err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("unable to resolve coordinates %.6f, %.6f to any known administrative location", lat, lng)
}

func (s *CompositeGeocodingService) saveCache(key string, result *GeocodingResult) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cache[key] = cacheEntry{
		result:    result,
		expiresAt: time.Now().Add(24 * time.Hour), // 24-hour cache
	}
}
