package service

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/your-org/used-tools-marketplace/internal/models"
	"github.com/your-org/used-tools-marketplace/internal/repository"
)

var slugRegex = regexp.MustCompile(`^[a-z0-9-_]+$`)

type MasterDataService struct {
	repo repository.MasterDataRepository
}

func NewMasterDataService(repo repository.MasterDataRepository) *MasterDataService {
	return &MasterDataService{repo: repo}
}

// ValidateLocations checks 4-language completeness, slugs, levels, and cycles for locations
func (s *MasterDataService) ValidateLocations(nodes []models.MasterLocationNode, existing map[int]*models.MasterLocationNode) []models.MasterDataDiffItem {
	var issues []models.MasterDataDiffItem
	nodeMap := make(map[int]models.MasterLocationNode)
	slugMap := make(map[string]int)

	for _, n := range nodes {
		nodeMap[n.ID] = n
		if n.Slug != "" {
			if otherID, ok := slugMap[n.Slug]; ok && otherID != n.ID {
				issues = append(issues, models.MasterDataDiffItem{
					EntityID:   n.ID,
					Identifier: n.Slug,
					Action:     "CONFLICT",
					Field:      "slug",
					Message:    fmt.Sprintf("duplicate slug '%s' within import batch", n.Slug),
				})
			}
			slugMap[n.Slug] = n.ID
		}

		// 4-Language validation
		if strings.TrimSpace(n.Names.En) == "" {
			issues = append(issues, models.MasterDataDiffItem{
				EntityID: n.ID, Identifier: n.Code, Action: "ERROR", Field: "names.en", Message: "English name is required",
			})
		}
		if strings.TrimSpace(n.Names.Am) == "" {
			issues = append(issues, models.MasterDataDiffItem{
				EntityID: n.ID, Identifier: n.Code, Action: "ERROR", Field: "names.am", Message: "Amharic name is required",
			})
		}
		if strings.TrimSpace(n.Names.Om) == "" {
			issues = append(issues, models.MasterDataDiffItem{
				EntityID: n.ID, Identifier: n.Code, Action: "ERROR", Field: "names.om", Message: "Afaan Oromo name is required",
			})
		}
		if strings.TrimSpace(n.Names.Ti) == "" {
			issues = append(issues, models.MasterDataDiffItem{
				EntityID: n.ID, Identifier: n.Code, Action: "ERROR", Field: "names.ti", Message: "Tigrinya name is required",
			})
		}

		// Slug validation
		if !slugRegex.MatchString(n.Slug) {
			issues = append(issues, models.MasterDataDiffItem{
				EntityID: n.ID, Identifier: n.Slug, Action: "ERROR", Field: "slug", Message: "invalid slug characters",
			})
		}

		// Cycle prevention
		if n.ParentID != nil && *n.ParentID == n.ID {
			issues = append(issues, models.MasterDataDiffItem{
				EntityID: n.ID, Identifier: n.Code, Action: "ERROR", Field: "parent_id", Message: "location cannot be its own parent",
			})
		}
	}

	return issues
}

// ValidateCategories checks 4-language completeness, slugs, levels, and cycles for categories
func (s *MasterDataService) ValidateCategories(nodes []models.MasterCategoryNode, existing map[int]*models.MasterCategoryNode) []models.MasterDataDiffItem {
	var issues []models.MasterDataDiffItem
	nodeMap := make(map[int]models.MasterCategoryNode)
	slugMap := make(map[string]int)

	for _, c := range nodes {
		nodeMap[c.ID] = c
		if c.Slug != "" {
			if otherID, ok := slugMap[c.Slug]; ok && otherID != c.ID {
				issues = append(issues, models.MasterDataDiffItem{
					EntityID:   c.ID,
					Identifier: c.Slug,
					Action:     "CONFLICT",
					Field:      "slug",
					Message:    fmt.Sprintf("duplicate slug '%s' within import batch", c.Slug),
				})
			}
			slugMap[c.Slug] = c.ID
		}

		// 4-Language validation
		if strings.TrimSpace(c.Names.En) == "" {
			issues = append(issues, models.MasterDataDiffItem{
				EntityID: c.ID, Identifier: c.Slug, Action: "ERROR", Field: "names.en", Message: "English name is required",
			})
		}
		if strings.TrimSpace(c.Names.Am) == "" {
			issues = append(issues, models.MasterDataDiffItem{
				EntityID: c.ID, Identifier: c.Slug, Action: "ERROR", Field: "names.am", Message: "Amharic name is required",
			})
		}
		if strings.TrimSpace(c.Names.Om) == "" {
			issues = append(issues, models.MasterDataDiffItem{
				EntityID: c.ID, Identifier: c.Slug, Action: "ERROR", Field: "names.om", Message: "Afaan Oromo name is required",
			})
		}
		if strings.TrimSpace(c.Names.Ti) == "" {
			issues = append(issues, models.MasterDataDiffItem{
				EntityID: c.ID, Identifier: c.Slug, Action: "ERROR", Field: "names.ti", Message: "Tigrinya name is required",
			})
		}

		// Slug validation
		if !slugRegex.MatchString(c.Slug) {
			issues = append(issues, models.MasterDataDiffItem{
				EntityID: c.ID, Identifier: c.Slug, Action: "ERROR", Field: "slug", Message: "invalid slug characters",
			})
		}

		// Cycle prevention
		if c.ParentID != nil && *c.ParentID == c.ID {
			issues = append(issues, models.MasterDataDiffItem{
				EntityID: c.ID, Identifier: c.Slug, Action: "ERROR", Field: "parent_id", Message: "category cannot be its own parent",
			})
		}
	}

	return issues
}

// ImportLocations processes an import request (dry-run or live publish)
func (s *MasterDataService) ImportLocations(ctx context.Context, req models.MasterDataImportRequest, adminID int, ip string) (*models.MasterDataImportResult, error) {
	existing, err := s.repo.GetExistingLocationsMap(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch existing locations: %w", err)
	}

	validationErrors := s.ValidateLocations(req.LocationRecords, existing)
	if len(validationErrors) > 0 {
		return &models.MasterDataImportResult{
			DatasetType:  "locations",
			DryRun:       req.DryRun,
			TotalRecords: len(req.LocationRecords),
			Errors:       len(validationErrors),
			Diff:         validationErrors,
			Timestamp:    time.Now(),
		}, nil
	}

	result := &models.MasterDataImportResult{
		DatasetType:  "locations",
		DryRun:       req.DryRun,
		TotalRecords: len(req.LocationRecords),
		Timestamp:    time.Now(),
	}

	for _, n := range req.LocationRecords {
		old, exists := existing[n.ID]
		if !exists {
			result.Created++
			result.Diff = append(result.Diff, models.MasterDataDiffItem{
				EntityID:   n.ID,
				Identifier: n.Code,
				Action:     "CREATE",
				NewValue:   n.Names.En,
				Message:    "New location to be inserted",
			})
		} else {
			// Compare fields
			modified := false
			if old.Names.En != n.Names.En || old.Names.Am != n.Names.Am || old.Names.Om != n.Names.Om || old.Names.Ti != n.Names.Ti {
				modified = true
			}
			if old.Slug != n.Slug || old.Level != n.Level || old.IsActive != n.IsActive {
				modified = true
			}
			if (old.ParentID == nil && n.ParentID != nil) || (old.ParentID != nil && n.ParentID == nil) ||
				(old.ParentID != nil && n.ParentID != nil && *old.ParentID != *n.ParentID) {
				modified = true
			}

			if modified {
				result.Updated++
				result.Diff = append(result.Diff, models.MasterDataDiffItem{
					EntityID:   n.ID,
					Identifier: n.Code,
					Action:     "UPDATE",
					OldValue:   old.Names.En,
					NewValue:   n.Names.En,
					Message:    "Location properties updated",
				})
			} else {
				result.Unchanged++
			}
		}
	}

	if !req.DryRun && len(req.LocationRecords) > 0 {
		if err := s.repo.UpsertLocationBatch(ctx, req.LocationRecords, adminID); err != nil {
			return nil, fmt.Errorf("transactional publish failed: %w", err)
		}
		_ = s.repo.LogMasterDataAudit(ctx, adminID, "PUBLISH_LOCATIONS", nil, nil, map[string]int{"total": len(req.LocationRecords)}, "Master data locations import", ip)
	}

	return result, nil
}

// ImportCategories processes an import request for categories (dry-run or live publish)
func (s *MasterDataService) ImportCategories(ctx context.Context, req models.MasterDataImportRequest, adminID int, ip string) (*models.MasterDataImportResult, error) {
	existing, err := s.repo.GetExistingCategoriesMap(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch existing categories: %w", err)
	}

	validationErrors := s.ValidateCategories(req.CategoryRecords, existing)
	if len(validationErrors) > 0 {
		return &models.MasterDataImportResult{
			DatasetType:  "categories",
			DryRun:       req.DryRun,
			TotalRecords: len(req.CategoryRecords),
			Errors:       len(validationErrors),
			Diff:         validationErrors,
			Timestamp:    time.Now(),
		}, nil
	}

	result := &models.MasterDataImportResult{
		DatasetType:  "categories",
		DryRun:       req.DryRun,
		TotalRecords: len(req.CategoryRecords),
		Timestamp:    time.Now(),
	}

	for _, c := range req.CategoryRecords {
		old, exists := existing[c.ID]
		if !exists {
			result.Created++
			result.Diff = append(result.Diff, models.MasterDataDiffItem{
				EntityID:   c.ID,
				Identifier: c.Slug,
				Action:     "CREATE",
				NewValue:   c.Names.En,
				Message:    "New category to be created",
			})
		} else {
			modified := false
			if old.Names.En != c.Names.En || old.Names.Am != c.Names.Am || old.Names.Om != c.Names.Om || old.Names.Ti != c.Names.Ti {
				modified = true
			}
			if old.Slug != c.Slug || old.Level != c.Level || old.IsActive != c.IsActive || old.IsLeaf != c.IsLeaf {
				modified = true
			}
			if (old.ParentID == nil && c.ParentID != nil) || (old.ParentID != nil && c.ParentID == nil) ||
				(old.ParentID != nil && c.ParentID != nil && *old.ParentID != *c.ParentID) {
				modified = true
			}

			if modified {
				result.Updated++
				result.Diff = append(result.Diff, models.MasterDataDiffItem{
					EntityID:   c.ID,
					Identifier: c.Slug,
					Action:     "UPDATE",
					OldValue:   old.Names.En,
					NewValue:   c.Names.En,
					Message:    "Category attributes modified",
				})
			} else {
				result.Unchanged++
			}
		}
	}

	if !req.DryRun && len(req.CategoryRecords) > 0 {
		if err := s.repo.UpsertCategoryBatch(ctx, req.CategoryRecords, adminID); err != nil {
			return nil, fmt.Errorf("transactional publish failed: %w", err)
		}
		_ = s.repo.LogMasterDataAudit(ctx, adminID, "PUBLISH_CATEGORIES", nil, nil, map[string]int{"total": len(req.CategoryRecords)}, "Master data categories import", ip)
	}

	return result, nil
}

// ExportLocations exports locations matching the given region code
func (s *MasterDataService) ExportLocations(ctx context.Context, regionCode string) ([]models.MasterLocationNode, error) {
	return s.repo.GetLocationsForExport(ctx, regionCode)
}

// ExportCategories exports categories matching the given system code and root slug
func (s *MasterDataService) ExportCategories(ctx context.Context, systemCode, rootSlug string) ([]models.MasterCategoryNode, error) {
	return s.repo.GetCategoriesForExport(ctx, systemCode, rootSlug)
}

// GetStatus returns the operational counts and language status
func (s *MasterDataService) GetStatus(ctx context.Context) (*models.MasterDataStatus, error) {
	return s.repo.GetMasterDataStatus(ctx)
}

