package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/your-org/used-tools-marketplace/internal/models"
	"github.com/your-org/used-tools-marketplace/internal/repository"
	"github.com/your-org/used-tools-marketplace/internal/service/geocoding"
)

type LocationService struct {
	repo          *repository.LocationRepo
	hierarchyRepo *repository.LocationHierarchyRepo
	geocoder      *geocoding.CompositeGeocodingService
}

func NewLocationService(repo repository.LocationRepository) *LocationService {
	return &LocationService{
		repo: repo,
	}
}

// WithHierarchy injects the LocationHierarchyRepo and initializes the Geocoder
func (s *LocationService) WithHierarchy(hRepo *repository.LocationHierarchyRepo) *LocationService {
	s.hierarchyRepo = hRepo
	if hRepo != nil {
		s.geocoder = geocoding.NewCompositeGeocodingService(hRepo)
	}
	return s
}

// ReverseGeocode resolves GPS coordinates to verified location hierarchy with intelligent caching
func (s *LocationService) ReverseGeocode(ctx context.Context, lat, lng float64) (*geocoding.GeocodingResult, error) {
	if s.geocoder == nil {
		return nil, fmt.Errorf("geocoder service not initialized")
	}
	if lat < -90 || lat > 90 || lng < -180 || lng > 180 {
		return nil, fmt.Errorf("coordinates out of valid range (lat: [-90, 90], lng: [-180, 180])")
	}
	return s.geocoder.ReverseGeocode(ctx, lat, lng)
}

// GetLocationTree returns the complete hierarchical tree
func (s *LocationService) GetLocationTree(ctx context.Context, lang string) ([]models.LocationTree, error) {
	if s.hierarchyRepo == nil {
		return nil, fmt.Errorf("location hierarchy repository not configured")
	}
	return s.hierarchyRepo.ListTree(ctx, lang)
}

// GetLocationByID returns a single location with ancestors/breadcrumbs
func (s *LocationService) GetLocationByID(ctx context.Context, id int, lang string) (*models.LocationHierarchy, error) {
	if s.hierarchyRepo == nil {
		return nil, fmt.Errorf("location hierarchy repository not configured")
	}
	return s.hierarchyRepo.GetByID(ctx, id, lang)
}

// ListChildren returns direct child locations under parentID
func (s *LocationService) ListChildren(ctx context.Context, parentID *int, lang string) ([]models.LocationHierarchy, error) {
	if s.hierarchyRepo == nil {
		return nil, fmt.Errorf("location hierarchy repository not configured")
	}
	return s.hierarchyRepo.ListChildren(ctx, parentID, lang)
}

// SearchHierarchy searches location hierarchy by query string
func (s *LocationService) SearchHierarchy(ctx context.Context, q, lang string) ([]models.LocationHierarchy, error) {
	if s.hierarchyRepo == nil {
		return nil, fmt.Errorf("location hierarchy repository not configured")
	}
	return s.hierarchyRepo.SearchLocations(ctx, q, lang, 20)
}

// ListLocationsAdmin retrieves paginated locations with filtering for admin management
func (s *LocationService) ListLocationsAdmin(ctx context.Context, filter models.LocationFilter) ([]models.LocationHierarchy, int, error) {
	if s.hierarchyRepo == nil {
		return nil, 0, fmt.Errorf("location hierarchy repository not configured")
	}
	return s.hierarchyRepo.ListLocations(ctx, filter)
}

// CreateLocation creates a new location node with validation and slug generation
func (s *LocationService) CreateLocation(ctx context.Context, loc *models.LocationHierarchy) (int, error) {
	if s.hierarchyRepo == nil {
		return 0, fmt.Errorf("location hierarchy repository not configured")
	}
	if strings.TrimSpace(loc.NameEn) == "" {
		return 0, fmt.Errorf("name_en is required")
	}
	if strings.TrimSpace(loc.LocationType) == "" {
		loc.LocationType = "woreda"
	}
	if strings.TrimSpace(loc.Slug) == "" {
		loc.Slug = strings.ToLower(strings.ReplaceAll(strings.TrimSpace(loc.NameEn), " ", "-"))
	}
	return s.hierarchyRepo.Create(ctx, loc)
}

// UpdateLocation updates an existing location node
func (s *LocationService) UpdateLocation(ctx context.Context, loc *models.LocationHierarchy) error {
	if s.hierarchyRepo == nil {
		return fmt.Errorf("location hierarchy repository not configured")
	}
	if loc.ID <= 0 {
		return fmt.Errorf("invalid location ID")
	}
	if strings.TrimSpace(loc.NameEn) == "" {
		return fmt.Errorf("name_en is required")
	}
	return s.hierarchyRepo.Update(ctx, loc)
}

// DeleteLocation removes a location node if safe
func (s *LocationService) DeleteLocation(ctx context.Context, id int) error {
	if s.hierarchyRepo == nil {
		return fmt.Errorf("location hierarchy repository not configured")
	}
	return s.hierarchyRepo.Delete(ctx, id)
}

// ImportLocations imports a batch of location nodes
func (s *LocationService) ImportLocations(ctx context.Context, locs []models.LocationHierarchy) (int, int, error) {
	if s.hierarchyRepo == nil {
		return 0, 0, fmt.Errorf("location hierarchy repository not configured")
	}
	imported := 0
	failed := 0
	for _, l := range locs {
		loc := l
		if strings.TrimSpace(loc.NameEn) == "" {
			failed++
			continue
		}
		if strings.TrimSpace(loc.Slug) == "" {
			loc.Slug = strings.ToLower(strings.ReplaceAll(strings.TrimSpace(loc.NameEn), " ", "-"))
		}
		_, err := s.hierarchyRepo.Create(ctx, &loc)
		if err != nil {
			failed++
		} else {
			imported++
		}
	}
	return imported, failed, nil
}

// ── Legacy Methods (Preserved for 100% Backward Compatibility) ───────────────

func (s *LocationService) ListRegions(ctx context.Context) ([]map[string]interface{}, error) {
	return s.repo.ListRegions(ctx)
}

func (s *LocationService) ListZones(ctx context.Context, region string) ([]map[string]interface{}, error) {
	return s.repo.ListZones(ctx, region)
}

func (s *LocationService) ListSubcities(ctx context.Context, region, zone string) ([]map[string]interface{}, error) {
	return s.repo.ListSubcities(ctx, region, zone)
}

func (s *LocationService) ListWoredas(ctx context.Context, region, zone, subcity string) ([]map[string]interface{}, error) {
	return s.repo.ListWoredas(ctx, region, zone, subcity)
}

func (s *LocationService) SearchLocations(ctx context.Context, q string) ([]map[string]interface{}, error) {
	return s.repo.SearchLocations(ctx, q)
}
