package service

import (
	"context"

	"github.com/your-org/used-tools-marketplace/internal/repository"
)

type LocationService struct {
	repo *repository.LocationRepo
}

func NewLocationService(repo repository.LocationRepository) *LocationService {
	return &LocationService{repo: repo}
}

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
