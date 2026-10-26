package service

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/your-org/used-tools-marketplace/internal/models"
	"github.com/your-org/used-tools-marketplace/internal/repository"
)

type CategoryHierarchyService struct {
	repo *repository.CategoryHierarchyRepo
}

func NewCategoryHierarchyService(repo *repository.CategoryHierarchyRepo) *CategoryHierarchyService {
	return &CategoryHierarchyService{repo: repo}
}

// GetByID retrieves a category with breadcrumbs
func (s *CategoryHierarchyService) GetByID(ctx context.Context, id int, lang string) (*models.CategoryHierarchy, error) {
	if lang == "" {
		lang = "en"
	}
	return s.repo.GetByID(ctx, id, lang)
}

// GetBySlug retrieves a category by slug
func (s *CategoryHierarchyService) GetBySlug(ctx context.Context, slug string, lang string) (*models.CategoryHierarchy, error) {
	if lang == "" {
		lang = "en"
	}
	return s.repo.GetBySlug(ctx, slug, lang)
}

// GetRoots returns all active root categories
func (s *CategoryHierarchyService) GetRoots(ctx context.Context, lang string) ([]models.CategoryHierarchy, error) {
	if lang == "" {
		lang = "en"
	}
	return s.repo.GetRoots(ctx, lang)
}

// ListCategories retrieves categories with filtering
func (s *CategoryHierarchyService) ListCategories(ctx context.Context, f models.CategoryFilter) ([]models.CategoryHierarchy, error) {
	if f.Page < 1 {
		f.Page = 1
	}
	if f.Limit < 1 || f.Limit > 500 {
		f.Limit = 100
	}
	return s.repo.ListCategories(ctx, f)
}

// GetCategoryTree builds the full hierarchical tree
func (s *CategoryHierarchyService) GetCategoryTree(ctx context.Context, lang string) ([]models.CategoryTree, error) {
	if lang == "" {
		lang = "en"
	}
	return s.repo.GetCategoryTree(ctx, lang)
}

// GetBreadcrumbs returns the category path
func (s *CategoryHierarchyService) GetBreadcrumbs(ctx context.Context, categoryID int, lang string) ([]models.CategoryBreadcrumb, error) {
	if lang == "" {
		lang = "en"
	}
	return s.repo.GetBreadcrumbs(ctx, categoryID, lang)
}

// GetDescendants returns all child category IDs
func (s *CategoryHierarchyService) GetDescendants(ctx context.Context, categoryID int) ([]int, error) {
	return s.repo.GetDescendants(ctx, categoryID)
}

// CheckSlugAvailability checks whether a slug can be used
func (s *CategoryHierarchyService) CheckSlugAvailability(ctx context.Context, slug string, excludeID int) (bool, error) {
	norm := generateSlug(slug)
	if !isValidSlug(norm) {
		return false, fmt.Errorf("invalid slug format: must be lowercase alphanumeric with hyphens")
	}
	exists, err := s.repo.CheckSlugExists(ctx, norm, excludeID)
	if err != nil {
		return false, err
	}
	return !exists, nil
}

// ReorderCategory updates category sort display order
func (s *CategoryHierarchyService) ReorderCategory(ctx context.Context, id int, displayOrder int) error {
	if id <= 0 {
		return fmt.Errorf("valid category ID is required")
	}
	return s.repo.ReorderCategory(ctx, id, displayOrder)
}

// SetCategoryStatus toggles active state of a category
func (s *CategoryHierarchyService) SetCategoryStatus(ctx context.Context, id int, isActive bool) error {
	if id <= 0 {
		return fmt.Errorf("valid category ID is required")
	}
	return s.repo.SetCategoryStatus(ctx, id, isActive)
}

// ValidateCategoryForListing verifies that a listing's category selection is valid and active
func (s *CategoryHierarchyService) ValidateCategoryForListing(ctx context.Context, categoryID *int, categoryHierarchyID *int) error {
	if categoryID == nil && categoryHierarchyID == nil {
		return fmt.Errorf("category is required")
	}

	targetID := 0
	if categoryHierarchyID != nil && *categoryHierarchyID > 0 {
		targetID = *categoryHierarchyID
	} else if categoryID != nil && *categoryID > 0 {
		targetID = *categoryID
	}
	if targetID <= 0 {
		return fmt.Errorf("a valid marketplace category is required")
	}

	if s.repo == nil {
		return nil
	}
	exists, isActive, err := s.repo.ValidateCategory(ctx, targetID)
	if err != nil {
		return fmt.Errorf("failed to validate category: %w", err)
	}
	if !exists {
		return fmt.Errorf("category ID %d does not exist", targetID)
	}
	if !isActive {
		return fmt.Errorf("category ID %d is inactive and cannot accept listings", targetID)
	}

	// If both root and subcategory IDs are supplied, verify hierarchy
	if categoryID != nil && categoryHierarchyID != nil && *categoryID > 0 && *categoryHierarchyID > 0 && *categoryID != *categoryHierarchyID {
		isDesc, err := s.repo.IsDescendantOf(ctx, *categoryHierarchyID, *categoryID)
		if err != nil {
			return fmt.Errorf("failed to verify category hierarchy: %w", err)
		}
		if !isDesc {
			return fmt.Errorf("category hierarchy mismatch: category %d is not a descendant of parent category %d", *categoryHierarchyID, *categoryID)
		}
	}

	return nil
}

// CreateCategory validates and creates a new category
func (s *CategoryHierarchyService) CreateCategory(ctx context.Context, cat *models.CategoryHierarchy) (int, error) {
	if cat.NameEn == "" {
		return 0, fmt.Errorf("name_en is required")
	}
	if cat.Slug == "" {
		cat.Slug = generateSlug(cat.NameEn)
	} else {
		cat.Slug = generateSlug(cat.Slug)
	}
	
	if !isValidSlug(cat.Slug) {
		return 0, fmt.Errorf("invalid slug format: must be lowercase alphanumeric with hyphens")
	}

	// Check slug uniqueness
	exists, err := s.repo.CheckSlugExists(ctx, cat.Slug, 0)
	if err != nil {
		return 0, fmt.Errorf("slug check failed: %w", err)
	}
	if exists {
		return 0, fmt.Errorf("slug '%s' is already in use by another category", cat.Slug)
	}

	// Prevent circular references
	if cat.ParentID != nil {
		if *cat.ParentID == cat.ID {
			return 0, fmt.Errorf("category cannot be its own parent")
		}
		
		parent, err := s.repo.GetByID(ctx, *cat.ParentID, "en")
		if err != nil {
			return 0, fmt.Errorf("parent category not found")
		}
		
		if parent.IsLeaf {
			return 0, fmt.Errorf("cannot add child to leaf category")
		}
	}

	return s.repo.CreateCategory(ctx, cat)
}

// UpdateCategory validates and updates a category
func (s *CategoryHierarchyService) UpdateCategory(ctx context.Context, cat *models.CategoryHierarchy) error {
	if cat.ID == 0 {
		return fmt.Errorf("category ID is required")
	}
	if cat.NameEn == "" {
		return fmt.Errorf("name_en is required")
	}
	if cat.Slug == "" {
		cat.Slug = generateSlug(cat.NameEn)
	} else {
		cat.Slug = generateSlug(cat.Slug)
	}
	if !isValidSlug(cat.Slug) {
		return fmt.Errorf("invalid slug format")
	}

	// Prevent self-parent reference immediately
	if cat.ParentID != nil && *cat.ParentID == cat.ID {
		return fmt.Errorf("category cannot be its own parent")
	}

	// Check slug uniqueness excluding self
	exists, err := s.repo.CheckSlugExists(ctx, cat.Slug, cat.ID)
	if err != nil {
		return fmt.Errorf("slug check failed: %w", err)
	}
	if exists {
		return fmt.Errorf("slug '%s' is already in use by another category", cat.Slug)
	}

	// Prevent circular references with descendants
	if cat.ParentID != nil {
		descendants, err := s.repo.GetDescendants(ctx, cat.ID)
		if err == nil {
			for _, descID := range descendants {
				if descID == *cat.ParentID {
					return fmt.Errorf("circular reference: cannot move category under its own descendant")
				}
			}
		}
	}

	return s.repo.UpdateCategory(ctx, cat)
}

// DeleteCategory validates and deletes a category
func (s *CategoryHierarchyService) DeleteCategory(ctx context.Context, id int) error {
	if id == 0 {
		return fmt.Errorf("category ID is required")
	}
	return s.repo.DeleteCategory(ctx, id)
}

// GetCategoryStats retrieves category statistics
func (s *CategoryHierarchyService) GetCategoryStats(ctx context.Context, categoryID int) (*models.CategoryWithStats, error) {
	return s.repo.GetCategoryStats(ctx, categoryID)
}

// Helper functions

func generateSlug(name string) string {
	slug := strings.ToLower(strings.TrimSpace(name))
	reg := regexp.MustCompile("[^a-z0-9\\s-]")
	slug = reg.ReplaceAllString(slug, "")
	slug = strings.ReplaceAll(slug, " ", "-")
	reg = regexp.MustCompile("-+")
	slug = reg.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")
	return slug
}

func isValidSlug(slug string) bool {
	match, _ := regexp.MatchString("^[a-z0-9]+(-[a-z0-9]+)*$", slug)
	return match
}
