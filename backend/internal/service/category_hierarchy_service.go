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

// ListCategories retrieves categories with filtering
func (s *CategoryHierarchyService) ListCategories(ctx context.Context, f models.CategoryFilter) ([]models.CategoryHierarchy, error) {
	if f.Page < 1 {
		f.Page = 1
	}
	if f.Limit < 1 || f.Limit > 100 {
		f.Limit = 50
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

// CreateCategory validates and creates a new category
func (s *CategoryHierarchyService) CreateCategory(ctx context.Context, cat *models.CategoryHierarchy) (int, error) {
	// Validate required fields
	if cat.NameEn == "" {
		return 0, fmt.Errorf("name_en is required")
	}
	if cat.Slug == "" {
		cat.Slug = generateSlug(cat.NameEn)
	}
	
	// Validate slug format
	if !isValidSlug(cat.Slug) {
		return 0, fmt.Errorf("invalid slug format: must be lowercase alphanumeric with hyphens")
	}

	// Prevent circular references
	if cat.ParentID != nil {
		if *cat.ParentID == cat.ID {
			return 0, fmt.Errorf("category cannot be its own parent")
		}
		
		// Verify parent exists
		parent, err := s.repo.GetByID(ctx, *cat.ParentID, "en")
		if err != nil {
			return 0, fmt.Errorf("parent category not found")
		}
		
		// Verify parent is not a descendant (prevent circular reference)
		descendants, err := s.repo.GetDescendants(ctx, cat.ID)
		if err == nil {
			for _, descID := range descendants {
				if descID == *cat.ParentID {
					return 0, fmt.Errorf("circular reference detected: parent cannot be a descendant")
				}
			}
		}
		
		// If parent is not a leaf, update it (leaf categories cannot have children)
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
		return fmt.Errorf("slug is required")
	}
	if !isValidSlug(cat.Slug) {
		return fmt.Errorf("invalid slug format")
	}

	// Prevent circular references
	if cat.ParentID != nil {
		if *cat.ParentID == cat.ID {
			return fmt.Errorf("category cannot be its own parent")
		}
		
		// Check if new parent is a descendant of this category
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
	// Convert to lowercase
	slug := strings.ToLower(name)
	// Remove special characters, keep alphanumeric and spaces
	reg := regexp.MustCompile("[^a-z0-9\\s-]")
	slug = reg.ReplaceAllString(slug, "")
	// Replace spaces with hyphens
	slug = strings.ReplaceAll(slug, " ", "-")
	// Remove multiple consecutive hyphens
	reg = regexp.MustCompile("-+")
	slug = reg.ReplaceAllString(slug, "-")
	// Trim hyphens from edges
	slug = strings.Trim(slug, "-")
	return slug
}

func isValidSlug(slug string) bool {
	// Slug must be lowercase alphanumeric with hyphens only
	match, _ := regexp.MatchString("^[a-z0-9]+(-[a-z0-9]+)*$", slug)
	return match
}
