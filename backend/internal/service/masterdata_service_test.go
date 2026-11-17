package service

import (
	"context"
	"testing"

	"github.com/your-org/used-tools-marketplace/internal/models"
)

type mockMasterDataRepo struct {
	locMap map[int]*models.MasterLocationNode
	catMap map[int]*models.MasterCategoryNode
}

func (m *mockMasterDataRepo) GetExistingLocationsMap(ctx context.Context) (map[int]*models.MasterLocationNode, error) {
	return m.locMap, nil
}
func (m *mockMasterDataRepo) GetExistingCategoriesMap(ctx context.Context) (map[int]*models.MasterCategoryNode, error) {
	return m.catMap, nil
}
func (m *mockMasterDataRepo) GetLocationsForExport(ctx context.Context, regionCode string) ([]models.MasterLocationNode, error) {
	var list []models.MasterLocationNode
	for _, l := range m.locMap {
		list = append(list, *l)
	}
	return list, nil
}
func (m *mockMasterDataRepo) GetCategoriesForExport(ctx context.Context, systemCode, rootSlug string) ([]models.MasterCategoryNode, error) {
	var list []models.MasterCategoryNode
	for _, c := range m.catMap {
		list = append(list, *c)
	}
	return list, nil
}
func (m *mockMasterDataRepo) UpsertLocationBatch(ctx context.Context, nodes []models.MasterLocationNode, adminID int) error {
	for _, n := range nodes {
		node := n
		m.locMap[n.ID] = &node
	}
	return nil
}
func (m *mockMasterDataRepo) UpsertCategoryBatch(ctx context.Context, nodes []models.MasterCategoryNode, adminID int) error {
	for _, c := range nodes {
		cat := c
		m.catMap[c.ID] = &cat
	}
	return nil
}
func (m *mockMasterDataRepo) GetMasterDataStatus(ctx context.Context) (*models.MasterDataStatus, error) {
	return &models.MasterDataStatus{
		LocationsCount:  len(m.locMap),
		CategoriesCount: len(m.catMap),
		SupportedLangs:  []string{"en", "am", "om", "ti"},
	}, nil
}
func (m *mockMasterDataRepo) LogMasterDataAudit(ctx context.Context, adminID int, action string, catID *int, oldVal, newVal interface{}, reason, ip string) error {
	return nil
}

func TestMasterData_ValidateLocations_FourLanguages(t *testing.T) {
	svc := NewMasterDataService(&mockMasterDataRepo{})

	invalidNode := models.MasterLocationNode{
		ID:    999,
		Code:  "ET-TEST",
		Type:  "woreda",
		Slug:  "test-woreda",
		Names: models.FourLanguageNames{En: "Test", Am: "ቴስት", Om: "" /* missing om */, Ti: "ቴስት"},
	}

	diffs := svc.ValidateLocations([]models.MasterLocationNode{invalidNode}, nil)
	if len(diffs) != 1 {
		t.Fatalf("expected 1 validation error for missing Afaan Oromo, got %d", len(diffs))
	}
	if diffs[0].Field != "names.om" {
		t.Errorf("expected error field 'names.om', got '%s'", diffs[0].Field)
	}
}

func TestMasterData_ValidateCategories_CycleDetection(t *testing.T) {
	svc := NewMasterDataService(&mockMasterDataRepo{})
	parentID := 100

	cycleNode := models.MasterCategoryNode{
		ID:                       100,
		ParentID:                 &parentID, // self-referencing cycle
		ClassificationSystemCode: "SYS_PROD",
		Slug:                     "cycle-cat",
		Names: models.FourLanguageNames{
			En: "Cycle", Am: "ሳይክል", Om: "Saayikilii", Ti: "ሳይክል",
		},
	}

	diffs := svc.ValidateCategories([]models.MasterCategoryNode{cycleNode}, nil)
	foundCycleError := false
	for _, d := range diffs {
		if d.Field == "parent_id" && d.Message == "category cannot be its own parent" {
			foundCycleError = true
			break
		}
	}
	if !foundCycleError {
		t.Errorf("expected cycle detection error for self-referencing parent")
	}
}

func TestMasterData_ValidateCategories_SlugAndBatchDuplicates(t *testing.T) {
	svc := NewMasterDataService(&mockMasterDataRepo{})

	nodes := []models.MasterCategoryNode{
		{
			ID: 10, Slug: "unique-slug-1", ClassificationSystemCode: "SYS_PROD",
			Names: models.FourLanguageNames{En: "One", Am: "አንድ", Om: "Tokko", Ti: "ሓደ"},
		},
		{
			ID: 11, Slug: "unique-slug-1", // duplicate slug in batch!
			ClassificationSystemCode: "SYS_PROD",
			Names: models.FourLanguageNames{En: "Two", Am: "ሁለት", Om: "Lama", Ti: "ክልተ"},
		},
	}

	diffs := svc.ValidateCategories(nodes, nil)
	foundDuplicate := false
	for _, d := range diffs {
		if d.Field == "slug" && d.Action == "CONFLICT" {
			foundDuplicate = true
			break
		}
	}
	if !foundDuplicate {
		t.Errorf("expected CONFLICT error for duplicate slug in batch")
	}
}

func TestMasterData_ImportCategories_DryRun(t *testing.T) {
	mockRepo := &mockMasterDataRepo{
		catMap: map[int]*models.MasterCategoryNode{
			1: {
				ID: 1, Slug: "electronics", Names: models.FourLanguageNames{En: "Electronics", Am: "ኤሌክትሮኒክስ", Om: "Elektirooniksii", Ti: "ኤሌክትሮኒክስ"},
			},
		},
	}
	svc := NewMasterDataService(mockRepo)

	req := models.MasterDataImportRequest{
		DatasetType: "categories",
		DryRun:      true,
		CategoryRecords: []models.MasterCategoryNode{
			// Unchanged
			{
				ID: 1, Slug: "electronics", Names: models.FourLanguageNames{En: "Electronics", Am: "ኤሌክትሮኒክስ", Om: "Elektirooniksii", Ti: "ኤሌክትሮኒክስ"},
			},
			// Modified
			{
				ID: 2, Slug: "vehicles", Names: models.FourLanguageNames{En: "Vehicles New", Am: "ተሽከርካሪዎች", Om: "Konkolaattota", Ti: "ተሽከርከርቲ"},
			},
			// New
			{
				ID: 3, Slug: "new-cat", Names: models.FourLanguageNames{En: "New Category", Am: "አዲስ", Om: "Haaraa", Ti: "ሓድሽ"},
			},
		},
	}

	mockRepo.catMap[2] = &models.MasterCategoryNode{
		ID: 2, Slug: "vehicles", Names: models.FourLanguageNames{En: "Vehicles Old", Am: "ተሽከርካሪዎች", Om: "Konkolaattota", Ti: "ተሽከርከርቲ"},
	}

	res, err := svc.ImportCategories(context.Background(), req, 1, "127.0.0.1")
	if err != nil {
		t.Fatalf("dry-run import failed: %v", err)
	}

	if res.Unchanged != 1 {
		t.Errorf("expected 1 unchanged record, got %d", res.Unchanged)
	}
	if res.Updated != 1 {
		t.Errorf("expected 1 updated record, got %d", res.Updated)
	}
	if res.Created != 1 {
		t.Errorf("expected 1 created record, got %d", res.Created)
	}
	if res.Errors != 0 {
		t.Errorf("expected 0 errors, got %d", res.Errors)
	}
	if !res.DryRun {
		t.Errorf("expected DryRun=true in result")
	}
}

func TestMasterData_ImportLocations_Publish(t *testing.T) {
	mockRepo := &mockMasterDataRepo{
		locMap: make(map[int]*models.MasterLocationNode),
	}
	svc := NewMasterDataService(mockRepo)

	req := models.MasterDataImportRequest{
		DatasetType: "locations",
		DryRun:      false,
		LocationRecords: []models.MasterLocationNode{
			{
				ID: 1, Code: "ET", Type: "country", Slug: "ethiopia", Level: 0, IsActive: true,
				Names: models.FourLanguageNames{En: "Ethiopia", Am: "ኢትዮጵያ", Om: "Itoophiyaa", Ti: "ኢትዮጵያ"},
			},
		},
	}

	res, err := svc.ImportLocations(context.Background(), req, 42, "127.0.0.1")
	if err != nil {
		t.Fatalf("publish failed: %v", err)
	}

	if res.Created != 1 {
		t.Errorf("expected 1 created, got %d", res.Created)
	}
	if len(mockRepo.locMap) != 1 {
		t.Errorf("expected 1 location in repository, got %d", len(mockRepo.locMap))
	}
}
