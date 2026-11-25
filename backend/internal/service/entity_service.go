package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/your-org/used-tools-marketplace/internal/models"
	"github.com/your-org/used-tools-marketplace/internal/repository"
)

type EntityService struct {
	repo *repository.EntityRepo
}

func NewEntityService(repo *repository.EntityRepo) *EntityService {
	return &EntityService{repo: repo}
}

// ResolveCanonicalEntity resolves user search or category input to canonical entities
func (s *EntityService) ResolveCanonicalEntity(ctx context.Context, term string, lang string, limit int) ([]models.EntityResolutionResult, error) {
	clean := strings.TrimSpace(term)
	if clean == "" {
		return []models.EntityResolutionResult{}, nil
	}
	if limit <= 0 || limit > 20 {
		limit = 5
	}
	if lang == "" {
		lang = "en"
	}
	return s.repo.ResolveCanonicalEntity(ctx, clean, lang, limit)
}

// GetEntityByID retrieves a canonical entity by ID with its relationships and discovery paths
func (s *EntityService) GetEntityByID(ctx context.Context, id int) (*models.CanonicalEntity, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid entity id: %d", id)
	}
	return s.repo.GetEntityByID(ctx, id)
}

// GetEntityBySlug retrieves a canonical entity by slug
func (s *EntityService) GetEntityBySlug(ctx context.Context, slug string) (*models.CanonicalEntity, error) {
	slug = strings.TrimSpace(slug)
	if slug == "" {
		return nil, fmt.Errorf("empty entity slug")
	}
	return s.repo.GetEntityBySlug(ctx, slug)
}

// GetEntityDiscoveryPaths returns all multi-domain discovery paths for an entity
func (s *EntityService) GetEntityDiscoveryPaths(ctx context.Context, entityID int) ([]models.EntityDiscoveryPath, error) {
	if entityID <= 0 {
		return nil, fmt.Errorf("invalid entity id: %d", entityID)
	}
	return s.repo.GetEntityDiscoveryPaths(ctx, entityID)
}

// GetEntityRelationships returns the relationship edges originating from this entity
func (s *EntityService) GetEntityRelationships(ctx context.Context, entityID int) ([]models.EntityRelationship, error) {
	if entityID <= 0 {
		return nil, fmt.Errorf("invalid entity id: %d", entityID)
	}
	return s.repo.GetEntityRelationships(ctx, entityID)
}

// GetEntityAliases returns the registered aliases and spelling variants for this entity
func (s *EntityService) GetEntityAliases(ctx context.Context, entityID int) ([]models.EntityAlias, error) {
	if entityID <= 0 {
		return nil, fmt.Errorf("invalid entity id: %d", entityID)
	}
	return s.repo.GetEntityAliases(ctx, entityID)
}
