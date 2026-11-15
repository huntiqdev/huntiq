package service

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/your-org/used-tools-marketplace/internal/config"
	"github.com/your-org/used-tools-marketplace/internal/database"
	"github.com/your-org/used-tools-marketplace/internal/models"
	"github.com/your-org/used-tools-marketplace/internal/repository"
)

func setupClassificationServiceTest(t *testing.T) (*pgxpool.Pool, *ClassificationService, func()) {
	if os.Getenv("DATABASE_URL") == "" && os.Getenv("DB_PASSWORD") == "" {
		t.Skip("Skipping classification service test: DATABASE_URL not set")
		return nil, nil, nil
	}

	cfg := config.NewTestConfig()
	db, err := database.NewPool(cfg)
	if err != nil {
		t.Skipf("Skipping classification service test: DB connection failed: %v", err)
		return nil, nil, nil
	}

	repo := repository.NewClassificationRepo(db)
	svc := NewClassificationService(repo)
	cleanup := func() {
		db.Close()
	}

	return db, svc, cleanup
}

func TestClassificationService_SchemeValidation(t *testing.T) {
	_, svc, cleanup := setupClassificationServiceTest(t)
	if svc == nil {
		return
	}
	defer cleanup()

	ctx := context.Background()

	// 1. Missing code
	_, err := svc.CreateScheme(ctx, &models.ClassificationScheme{
		NameEn: "Missing Code Scheme",
	})
	if err == nil {
		t.Errorf("Expected error for missing scheme code, got nil")
	}

	// 2. Missing name
	_, err = svc.CreateScheme(ctx, &models.ClassificationScheme{
		Code: "VALID_CODE",
	})
	if err == nil {
		t.Errorf("Expected error for missing scheme name_en, got nil")
	}
}

func TestClassificationService_ConceptValidation(t *testing.T) {
	_, svc, cleanup := setupClassificationServiceTest(t)
	if svc == nil {
		return
	}
	defer cleanup()

	ctx := context.Background()

	// 1. Missing canonical name
	_, err := svc.CreateConcept(ctx, &models.ClassificationConcept{
		Code: "test_code",
	})
	if err == nil {
		t.Errorf("Expected error for missing canonical_name, got nil")
	}

	// 2. Auto-code generation from name
	conceptName := fmt.Sprintf("Hydraulic Breaker Attachment %d", time.Now().UnixNano())
	concept, err := svc.CreateConcept(ctx, &models.ClassificationConcept{
		CanonicalName: conceptName,
		ConceptType:   "EQUIPMENT",
	})
	if err != nil {
		t.Fatalf("CreateConcept with auto-code failed: %v", err)
	}
	if concept.Code == "" {
		t.Errorf("Expected code to be auto-generated from canonical name")
	}
}

func TestClassificationService_MovePlacement_CycleAndBoundaryGuards(t *testing.T) {
	db, svc, cleanup := setupClassificationServiceTest(t)
	if svc == nil {
		return
	}
	defer cleanup()

	ctx := context.Background()

	// 1. Invalid ID checks
	err := svc.MovePlacement(ctx, 0, 1)
	if err == nil {
		t.Errorf("Expected error for placement ID <= 0")
	}

	// 2. Self-parenting check
	err = svc.MovePlacement(ctx, 5, 5)
	if err == nil || err.Error() != "cannot move placement to become its own parent" {
		t.Errorf("Expected self-parenting error, got: %v", err)
	}

	// 3. Move parent into its child (Hierarchy cycle guard)
	// In seed data: Placement 1 is Electronics (path /1/), Placement 9 is Mobile Phones (path /1/9/)
	err = svc.MovePlacement(ctx, 1, 9)
	if err == nil {
		t.Errorf("Expected cycle detection error when moving parent into its child, got nil")
	} else if err.Error() != "cannot move category inside one of its own descendants (hierarchy cycle)" {
		t.Logf("Got cycle error message: %v", err)
	}

	// 4. Cross-scheme move guard
	// Create a placement in a different scheme (e.g. SYS_IND)
	repo := repository.NewClassificationRepo(db)
	indScheme, err := repo.GetSchemeByCode(ctx, "SYS_IND")
	if err != nil {
		t.Fatalf("GetSchemeByCode(SYS_IND) failed: %v", err)
	}

	indConcept, _ := repo.CreateConcept(ctx, &models.ClassificationConcept{
		Code: "ind_test_concept_c", CanonicalName: "Test Ind Concept", NameEn: "Test Ind Concept", ConceptType: "PRODUCT", Status: "ACTIVE",
	})
	defer func() { _, _ = db.Exec(ctx, "DELETE FROM classification_concepts WHERE id = $1", indConcept.ID) }()

	indPlacement, err := repo.CreatePlacement(ctx, &models.ClassificationPlacement{
		SchemeID: indScheme.ID, ConceptID: indConcept.ID, Role: "ROOT", Slug: "test-ind-placement", Status: "PUBLISHED", IsActive: true,
	})
	if err != nil {
		t.Fatalf("CreatePlacement in SYS_IND failed: %v", err)
	}
	defer func() { _, _ = db.Exec(ctx, "DELETE FROM classification_placements WHERE id = $1", indPlacement.ID) }()

	// Try moving Placement 9 (in SYS_PROD) under indPlacement (in SYS_IND)
	err = svc.MovePlacement(ctx, 9, indPlacement.ID)
	if err == nil {
		t.Errorf("Expected error when moving placement across different schemes, got nil")
	} else if err.Error() != "cannot move placement across different schemes" {
		t.Logf("Got cross-scheme error message: %v", err)
	}
}

func TestClassificationService_RelationshipValidation(t *testing.T) {
	_, svc, cleanup := setupClassificationServiceTest(t)
	if svc == nil {
		return
	}
	defer cleanup()

	ctx := context.Background()

	// 1. Missing IDs
	err := svc.CreateRelationship(ctx, &models.ClassificationRelationship{
		SourceConceptID: 0,
		TargetConceptID: 5,
	})
	if err == nil {
		t.Errorf("Expected error for missing source concept ID")
	}

	// 2. Self-referencing
	err = svc.CreateRelationship(ctx, &models.ClassificationRelationship{
		SourceConceptID: 5,
		TargetConceptID: 5,
	})
	if err == nil {
		t.Errorf("Expected error for self-referencing relationship")
	}

	// 3. Missing relationship type
	err = svc.CreateRelationship(ctx, &models.ClassificationRelationship{
		SourceConceptID: 5,
		TargetConceptID: 6,
	})
	if err == nil {
		t.Errorf("Expected error for missing relationship_type_code")
	}
}

func TestClassificationService_AttributeRuleValidation(t *testing.T) {
	_, svc, cleanup := setupClassificationServiceTest(t)
	if svc == nil {
		return
	}
	defer cleanup()

	ctx := context.Background()

	// 1. Missing attribute ID
	err := svc.SaveAttributeRule(ctx, &models.ClassificationAttributeRule{
		AttributeID: 0,
	})
	if err == nil {
		t.Errorf("Expected error for missing attribute_id")
	}

	// 2. Neither placement_id nor concept_id specified
	err = svc.SaveAttributeRule(ctx, &models.ClassificationAttributeRule{
		AttributeID: 1,
	})
	if err == nil {
		t.Errorf("Expected error when neither placement_id nor concept_id is specified")
	}
}
