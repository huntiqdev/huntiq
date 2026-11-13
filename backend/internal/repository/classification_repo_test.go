package repository

import (
	"context"
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/your-org/used-tools-marketplace/internal/config"
	"github.com/your-org/used-tools-marketplace/internal/database"
	"github.com/your-org/used-tools-marketplace/internal/models"
)

func setupClassificationTestDB(t *testing.T) (*pgxpool.Pool, *ClassificationRepo, func()) {
	if os.Getenv("DATABASE_URL") == "" && os.Getenv("DB_PASSWORD") == "" {
		t.Skip("Skipping classification DB test: DATABASE_URL not set")
		return nil, nil, nil
	}

	cfg := config.NewTestConfig()
	db, err := database.NewPool(cfg)
	if err != nil {
		t.Skipf("Skipping classification DB test: DB connection failed: %v", err)
		return nil, nil, nil
	}

	repo := NewClassificationRepo(db)
	cleanup := func() {
		db.Close()
	}

	return db, repo, cleanup
}

func TestClassificationRepo_Schemes(t *testing.T) {
	db, repo, cleanup := setupClassificationTestDB(t)
	if db == nil {
		return
	}
	defer cleanup()

	ctx := context.Background()

	// 1. List seeded schemes
	schemes, err := repo.ListSchemes(ctx, "ACTIVE")
	if err != nil {
		t.Fatalf("ListSchemes failed: %v", err)
	}
	if len(schemes) < 7 {
		t.Errorf("Expected at least 7 seeded schemes, got %d", len(schemes))
	}

	// 2. Get scheme by code
	prodScheme, err := repo.GetSchemeByCode(ctx, "SYS_PROD")
	if err != nil {
		t.Fatalf("GetSchemeByCode(SYS_PROD) failed: %v", err)
	}
	if prodScheme.Code != "SYS_PROD" {
		t.Errorf("Expected code SYS_PROD, got %s", prodScheme.Code)
	}

	// 3. Create a custom scheme
	testSchemeCode := "TEST_SCHEME_XYZ"
	// Clean up if already exists from prior run
	_, _ = db.Exec(ctx, "DELETE FROM classification_schemes WHERE code = $1", testSchemeCode)

	rulesJSON, _ := json.Marshal(map[string]interface{}{"allow_multi_parent": false, "max_depth": 5})
	createdScheme, err := repo.CreateScheme(ctx, &models.ClassificationScheme{
		Code:        testSchemeCode,
		NameEn:      "Test Industrial Scheme",
		NameAm:      "የሙከራ ምድብ",
		Description: "A custom scheme for automated integration testing",
		Status:      "ACTIVE",
		IsSystem:    false,
		Rules:       rulesJSON,
	})
	if err != nil {
		t.Fatalf("CreateScheme failed: %v", err)
	}
	if createdScheme.ID == 0 {
		t.Errorf("Expected non-zero ID for created scheme")
	}

	// 4. Update the custom scheme
	createdScheme.Description = "Updated test description"
	err = repo.UpdateScheme(ctx, createdScheme)
	if err != nil {
		t.Fatalf("UpdateScheme failed: %v", err)
	}

	fetched, err := repo.GetSchemeByID(ctx, createdScheme.ID)
	if err != nil {
		t.Fatalf("GetSchemeByID failed: %v", err)
	}
	if fetched.Description != "Updated test description" {
		t.Errorf("Expected updated description, got %s", fetched.Description)
	}

	// 5. Clean up custom scheme
	_, _ = db.Exec(ctx, "DELETE FROM classification_schemes WHERE code = $1", testSchemeCode)
}

func TestClassificationRepo_Concepts(t *testing.T) {
	db, repo, cleanup := setupClassificationTestDB(t)
	if db == nil {
		return
	}
	defer cleanup()

	ctx := context.Background()

	testConceptCode := "test_drill_press_xyz"
	// Clean up if exists
	_, _ = db.Exec(ctx, "DELETE FROM classification_concepts WHERE code = $1", testConceptCode)

	metaJSON, _ := json.Marshal(map[string]interface{}{"power_source": "electric"})
	// 1. Create canonical concept
	concept, err := repo.CreateConcept(ctx, &models.ClassificationConcept{
		Code:          testConceptCode,
		CanonicalName: "Heavy Duty Drill Press",
		NameEn:        "Heavy Duty Drill Press",
		NameAm:        "የከባድ ሥራ ድሪል ፕሬስ",
		Description:   "Stationary drilling machine",
		ConceptType:   "EQUIPMENT",
		Status:        "ACTIVE",
		Metadata:      metaJSON,
	})
	if err != nil {
		t.Fatalf("CreateConcept failed: %v", err)
	}
	if concept.ID == 0 {
		t.Fatalf("Expected non-zero concept ID")
	}

	// 2. GetConceptByID
	fetched, err := repo.GetConceptByID(ctx, concept.ID)
	if err != nil {
		t.Fatalf("GetConceptByID failed: %v", err)
	}
	if fetched.CanonicalName != "Heavy Duty Drill Press" {
		t.Errorf("Expected 'Heavy Duty Drill Press', got %s", fetched.CanonicalName)
	}

	// 3. GetConceptByCode
	byCode, err := repo.GetConceptByCode(ctx, testConceptCode)
	if err != nil {
		t.Fatalf("GetConceptByCode failed: %v", err)
	}
	if byCode.ID != concept.ID {
		t.Errorf("ID mismatch: %d vs %d", byCode.ID, concept.ID)
	}

	// 4. CheckDuplicateCandidates (fuzzy similarity)
	candidates, err := repo.CheckDuplicateCandidates(ctx, "Drill Press", testConceptCode)
	if err != nil {
		t.Fatalf("CheckDuplicateCandidates failed: %v", err)
	}
	found := false
	for _, c := range candidates {
		if c.ID == concept.ID {
			found = true
			break
		}
	}
	if !found {
		t.Logf("Candidate search did not return newly created concept directly; candidates=%d", len(candidates))
	}

	// Clean up
	_, _ = db.Exec(ctx, "DELETE FROM classification_concepts WHERE id = $1", concept.ID)
}

func TestClassificationRepo_Placements_TreeAndInheritance(t *testing.T) {
	db, repo, cleanup := setupClassificationTestDB(t)
	if db == nil {
		return
	}
	defer cleanup()

	ctx := context.Background()

	scheme, err := repo.GetSchemeByCode(ctx, "SYS_PROD")
	if err != nil {
		t.Fatalf("GetSchemeByCode failed: %v", err)
	}

	// 1. Create a Root Concept and Placement
	rootConcept, err := repo.CreateConcept(ctx, &models.ClassificationConcept{
		Code:          "test_root_concept_xyz",
		CanonicalName: "Test Root Machinery",
		NameEn:        "Test Root Machinery",
		ConceptType:   "EQUIPMENT",
		Status:        "ACTIVE",
	})
	if err != nil {
		t.Fatalf("Create root concept failed: %v", err)
	}
	defer func() {
		_, _ = db.Exec(ctx, "DELETE FROM classification_concepts WHERE id = $1", rootConcept.ID)
	}()

	rootPlacement, err := repo.CreatePlacement(ctx, &models.ClassificationPlacement{
		SchemeID:  scheme.ID,
		ConceptID: rootConcept.ID,
		Role:      "ROOT",
		Slug:      "test-root-machinery-xyz",
		Status:    "PUBLISHED",
		IsActive:  true,
	})
	if err != nil {
		t.Fatalf("Create root placement failed: %v", err)
	}
	defer func() {
		_, _ = db.Exec(ctx, "DELETE FROM classification_placements WHERE id = $1", rootPlacement.ID)
	}()

	if rootPlacement.Depth != 0 {
		t.Errorf("Expected depth 0 for root, got %d", rootPlacement.Depth)
	}
	if !strings.HasPrefix(rootPlacement.Path, "/") || !strings.HasSuffix(rootPlacement.Path, "/") {
		t.Errorf("Expected path format /<id>/, got %s", rootPlacement.Path)
	}

	// 2. Create Child Placement
	childConcept, err := repo.CreateConcept(ctx, &models.ClassificationConcept{
		Code:          "test_child_concept_xyz",
		CanonicalName: "Test Child Lathes",
		NameEn:        "Test Child Lathes",
		ConceptType:   "EQUIPMENT",
		Status:        "ACTIVE",
	})
	if err != nil {
		t.Fatalf("Create child concept failed: %v", err)
	}
	defer func() {
		_, _ = db.Exec(ctx, "DELETE FROM classification_concepts WHERE id = $1", childConcept.ID)
	}()

	parentID := rootPlacement.ID
	childPlacement, err := repo.CreatePlacement(ctx, &models.ClassificationPlacement{
		SchemeID:          scheme.ID,
		ConceptID:         childConcept.ID,
		ParentPlacementID: &parentID,
		Role:              "CATEGORY",
		Slug:              "test-child-lathes-xyz",
		Status:            "PUBLISHED",
		IsActive:          true,
	})
	if err != nil {
		t.Fatalf("Create child placement failed: %v", err)
	}
	defer func() {
		_, _ = db.Exec(ctx, "DELETE FROM classification_placements WHERE id = $1", childPlacement.ID)
	}()

	if childPlacement.Depth != 1 {
		t.Errorf("Expected child depth 1, got %d", childPlacement.Depth)
	}
	if !strings.Contains(childPlacement.Path, rootPlacement.Path) {
		t.Errorf("Expected child path %s to contain root path %s", childPlacement.Path, rootPlacement.Path)
	}

	// 3. Test Breadcrumbs
	crumbs, err := repo.GetBreadcrumbs(ctx, childPlacement.ID, "en")
	if err != nil {
		t.Fatalf("GetBreadcrumbs failed: %v", err)
	}
	if len(crumbs) < 2 {
		t.Errorf("Expected at least 2 breadcrumb items, got %d", len(crumbs))
	}
	if crumbs[0].ID != rootPlacement.ID {
		t.Errorf("Expected first breadcrumb to be root ID %d, got %d", rootPlacement.ID, crumbs[0].ID)
	}

	// 4. Test Scheme Tree retrieval
	tree, err := repo.GetTree(ctx, scheme.Code, "en", 3)
	if err != nil {
		t.Fatalf("GetTree failed: %v", err)
	}
	if len(tree) == 0 {
		t.Errorf("Expected non-empty tree for scheme %s", scheme.Code)
	}

	// 5. Test GetRoots
	roots, err := repo.GetRoots(ctx, scheme.Code, "en")
	if err != nil {
		t.Fatalf("GetRoots failed: %v", err)
	}
	if len(roots) == 0 {
		t.Errorf("Expected roots to be returned")
	}

	// 6. Test GetChildren
	children, err := repo.GetChildren(ctx, rootPlacement.ID, "en")
	if err != nil {
		t.Fatalf("GetChildren failed: %v", err)
	}
	if len(children) != 1 || children[0].ID != childPlacement.ID {
		t.Errorf("Expected 1 child matching childPlacement ID %d, got %d", childPlacement.ID, len(children))
	}
}

func TestClassificationRepo_AttributeInheritance(t *testing.T) {
	db, repo, cleanup := setupClassificationTestDB(t)
	if db == nil {
		return
	}
	defer cleanup()

	ctx := context.Background()
	scheme, err := repo.GetSchemeByCode(ctx, "SYS_PROD")
	if err != nil {
		t.Fatalf("GetSchemeByCode failed: %v", err)
	}

	// 1. Create attribute definitions in attribute_definitions
	var attrRootID, attrChildID int
	err = db.QueryRow(ctx, `
		INSERT INTO attribute_definitions (name_en, slug, data_type, is_required)
		VALUES ('Warranty Period', 'test_warranty_xyz', 'number', false)
		RETURNING id`).Scan(&attrRootID)
	if err != nil {
		t.Fatalf("Insert root attribute definition failed: %v", err)
	}
	defer func() {
		_, _ = db.Exec(ctx, "DELETE FROM attribute_definitions WHERE id = $1", attrRootID)
	}()

	err = db.QueryRow(ctx, `
		INSERT INTO attribute_definitions (name_en, slug, data_type, is_required)
		VALUES ('Operating Voltage', 'test_voltage_xyz', 'select', true)
		RETURNING id`).Scan(&attrChildID)
	if err != nil {
		t.Fatalf("Insert child attribute definition failed: %v", err)
	}
	defer func() {
		_, _ = db.Exec(ctx, "DELETE FROM attribute_definitions WHERE id = $1", attrChildID)
	}()

	// 2. Setup Root + Child placements
	rootConcept, _ := repo.CreateConcept(ctx, &models.ClassificationConcept{
		Code: "test_attr_root_c", CanonicalName: "Test Attr Root", NameEn: "Test Attr Root", ConceptType: "PRODUCT", Status: "ACTIVE",
	})
	defer func() { _, _ = db.Exec(ctx, "DELETE FROM classification_concepts WHERE id = $1", rootConcept.ID) }()

	rootPlacement, _ := repo.CreatePlacement(ctx, &models.ClassificationPlacement{
		SchemeID: scheme.ID, ConceptID: rootConcept.ID, Role: "ROOT", Slug: "test-attr-root", Status: "PUBLISHED", IsActive: true,
	})
	defer func() { _, _ = db.Exec(ctx, "DELETE FROM classification_placements WHERE id = $1", rootPlacement.ID) }()

	childConcept, _ := repo.CreateConcept(ctx, &models.ClassificationConcept{
		Code: "test_attr_child_c", CanonicalName: "Test Attr Child", NameEn: "Test Attr Child", ConceptType: "PRODUCT", Status: "ACTIVE",
	})
	defer func() { _, _ = db.Exec(ctx, "DELETE FROM classification_concepts WHERE id = $1", childConcept.ID) }()

	parentID := rootPlacement.ID
	childPlacement, _ := repo.CreatePlacement(ctx, &models.ClassificationPlacement{
		SchemeID: scheme.ID, ConceptID: childConcept.ID, ParentPlacementID: &parentID, Role: "CATEGORY", Slug: "test-attr-child", Status: "PUBLISHED", IsActive: true,
	})
	defer func() { _, _ = db.Exec(ctx, "DELETE FROM classification_placements WHERE id = $1", childPlacement.ID) }()

	// 3. Add attribute rule to Root
	rootPlacementID := rootPlacement.ID
	err = repo.SaveAttributeRule(ctx, &models.ClassificationAttributeRule{
		PlacementID:  &rootPlacementID,
		AttributeID:  attrRootID,
		IsRequired:   false,
		IsInherited:  true,
		DisplayOrder: 1,
	})
	if err != nil {
		t.Fatalf("SaveAttributeRule root failed: %v", err)
	}

	// 4. Add attribute rule to Child
	childPlacementID := childPlacement.ID
	err = repo.SaveAttributeRule(ctx, &models.ClassificationAttributeRule{
		PlacementID:  &childPlacementID,
		AttributeID:  attrChildID,
		IsRequired:   true,
		IsInherited:  true,
		DisplayOrder: 2,
	})
	if err != nil {
		t.Fatalf("SaveAttributeRule child failed: %v", err)
	}

	// 5. Query Effective Attributes for Child Placement
	effectiveAttrs, err := repo.GetEffectiveAttributes(ctx, childPlacement.ID)
	if err != nil {
		t.Fatalf("GetEffectiveAttributes failed: %v", err)
	}

	hasWarranty := false
	hasVoltage := false
	for _, a := range effectiveAttrs {
		if a.Slug == "test_warranty_xyz" {
			hasWarranty = true
			if !a.IsInherited {
				t.Errorf("Expected test_warranty_xyz to be marked IsInherited=true")
			}
		}
		if a.Slug == "test_voltage_xyz" {
			hasVoltage = true
			if a.IsInherited {
				t.Errorf("Expected test_voltage_xyz to be direct (IsInherited=false)")
			}
		}
	}

	if !hasWarranty {
		t.Errorf("Expected child placement to inherit 'test_warranty_xyz' from root ancestor")
	}
	if !hasVoltage {
		t.Errorf("Expected child placement to have direct attribute 'test_voltage_xyz'")
	}
}

func TestClassificationRepo_Relationships(t *testing.T) {
	db, repo, cleanup := setupClassificationTestDB(t)
	if db == nil {
		return
	}
	defer cleanup()

	ctx := context.Background()

	// 1. List relationship types
	types, err := repo.ListRelationshipTypes(ctx)
	if err != nil {
		t.Fatalf("ListRelationshipTypes failed: %v", err)
	}
	if len(types) == 0 {
		t.Fatalf("Expected seeded relationship types, got 0")
	}

	// 2. Create source & target concepts
	sourceConcept, _ := repo.CreateConcept(ctx, &models.ClassificationConcept{
		Code: "rel_source_c", CanonicalName: "Diamond Core Bit", NameEn: "Diamond Core Bit", ConceptType: "EQUIPMENT", Status: "ACTIVE",
	})
	defer func() { _, _ = db.Exec(ctx, "DELETE FROM classification_concepts WHERE id = $1", sourceConcept.ID) }()

	targetConcept, _ := repo.CreateConcept(ctx, &models.ClassificationConcept{
		Code: "rel_target_c", CanonicalName: "Core Drill Rig", NameEn: "Core Drill Rig", ConceptType: "EQUIPMENT", Status: "ACTIVE",
	})
	defer func() { _, _ = db.Exec(ctx, "DELETE FROM classification_concepts WHERE id = $1", targetConcept.ID) }()

	// 3. Create Relationship: Diamond Core Bit COMPATIBLE_WITH Core Drill Rig
	rel := &models.ClassificationRelationship{
		SourceConceptID:      sourceConcept.ID,
		RelationshipTypeCode: "COMPATIBLE_WITH",
		TargetConceptID:      targetConcept.ID,
		Status:               "ACTIVE",
	}
	err = repo.CreateRelationship(ctx, rel)
	if err != nil {
		t.Fatalf("CreateRelationship failed: %v", err)
	}

	// 4. Query relationships for source concept
	rels, err := repo.GetRelationships(ctx, sourceConcept.ID, "")
	if err != nil {
		t.Fatalf("GetRelationships failed: %v", err)
	}
	if len(rels) != 1 {
		t.Fatalf("Expected 1 relationship, got %d", len(rels))
	}
	if rels[0].RelationshipTypeCode != "COMPATIBLE_WITH" {
		t.Errorf("Expected relationship COMPATIBLE_WITH, got %s", rels[0].RelationshipTypeCode)
	}
}

func TestClassificationRepo_ImpactAnalysis(t *testing.T) {
	db, repo, cleanup := setupClassificationTestDB(t)
	if db == nil {
		return
	}
	defer cleanup()

	ctx := context.Background()

	// Placement 1 is Electronics root
	impact, err := repo.AnalyzeImpact(ctx, "PLACEMENT", 1)
	if err != nil {
		t.Fatalf("AnalyzeImpact failed: %v", err)
	}
	if impact.EntityID != 1 {
		t.Errorf("Expected EntityID 1, got %d", impact.EntityID)
	}
	if impact.AffectedDescendantsCount < 2 {
		t.Errorf("Expected multiple descendants for Electronics root, got %d", impact.AffectedDescendantsCount)
	}
	t.Logf("Electronics root impact: %d descendants, %d active listings, safe=%t",
		impact.AffectedDescendantsCount, impact.AffectedListingsCount, impact.CanSafelyExecute)
}

func TestClassificationRepo_UniversalSearch(t *testing.T) {
	db, repo, cleanup := setupClassificationTestDB(t)
	if db == nil {
		return
	}
	defer cleanup()

	ctx := context.Background()

	// Search for "electronics" or "phone"
	result, err := repo.SearchUniversal(ctx, "phone", "SYS_PROD", "en", 10)
	if err != nil {
		t.Fatalf("SearchUniversal failed: %v", err)
	}
	t.Logf("SearchUniversal 'phone': %d concepts matched, %d placements matched",
		len(result.Concepts), len(result.Placements))
}
