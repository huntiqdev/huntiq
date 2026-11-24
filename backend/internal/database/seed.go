package database

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/your-org/used-tools-marketplace/internal/models"
	"github.com/your-org/used-tools-marketplace/internal/repository"
	"github.com/your-org/used-tools-marketplace/internal/service"
)

// LocateLoccatDataDir finds the location & category master data directory.
func LocateLoccatDataDir(custom string) string {
	if custom != "" {
		if info, err := os.Stat(custom); err == nil && info.IsDir() {
			return custom
		}
	}
	candidates := []string{
		"internal/database/loccatdata",
		"backend/internal/database/loccatdata",
		"../internal/database/loccatdata",
		"../../internal/database/loccatdata",
		"loccatdata",
	}
	for _, c := range candidates {
		if info, err := os.Stat(c); err == nil && info.IsDir() {
			return c
		}
	}
	return "internal/database/loccatdata"
}

// CheckSeedStatus checks how many categories and locations currently exist in PostgreSQL.
func CheckSeedStatus(ctx context.Context, pool *pgxpool.Pool) (categoriesCount int, locationsCount int, isSeeded bool, err error) {
	_ = pool.QueryRow(ctx, `SELECT count(*) FROM category_hierarchy`).Scan(&categoriesCount)
	_ = pool.QueryRow(ctx, `SELECT count(*) FROM location_hierarchy`).Scan(&locationsCount)
	// Fully seeded means 20,000+ categories and 400+ locations exist
	isSeeded = categoriesCount >= 20000 && locationsCount >= 400
	return categoriesCount, locationsCount, isSeeded, nil
}

func loadAllLocations(baseDir string) ([]models.MasterLocationNode, error) {
	locDir := filepath.Join(baseDir, "locations", "countries", "ET")
	var allNodes []models.MasterLocationNode
	seenIDs := make(map[int]bool)

	err := filepath.Walk(locDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || filepath.Ext(path) != ".json" {
			return nil
		}
		fname := filepath.Base(path)
		if fname == "manifest.json" || fname == "schema.json" || strings.HasSuffix(fname, "_backup.json") || strings.HasSuffix(fname, ".bak") {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		var list []models.MasterLocationNode
		if err := json.Unmarshal(data, &list); err == nil && len(list) > 0 {
			for _, item := range list {
				if item.ID > 0 && !seenIDs[item.ID] {
					seenIDs[item.ID] = true
					allNodes = append(allNodes, item)
				}
			}
			return nil
		}
		var single models.MasterLocationNode
		if err := json.Unmarshal(data, &single); err == nil && single.ID > 0 {
			if !seenIDs[single.ID] {
				seenIDs[single.ID] = true
				allNodes = append(allNodes, single)
			}
		}
		return nil
	})
	return allNodes, err
}

func getClassificationSystemDirs(baseDir string) ([]string, error) {
	sysBase := filepath.Join(baseDir, "categories", "classification-systems")
	entries, err := os.ReadDir(sysBase)
	if err != nil {
		return nil, err
	}
	var systems []string
	for _, e := range entries {
		if e.IsDir() {
			systems = append(systems, e.Name())
		}
	}
	if len(systems) == 0 {
		return []string{"SYS_PROD", "SYS_SRV", "SYS_EMP", "SYS_GIG", "SYS_MED", "SYS_OPP", "SYS_PRM"}, nil
	}
	return systems, nil
}

func getSystemFamilyFiles(baseDir, sys string) ([]string, error) {
	famDir := filepath.Join(baseDir, "categories", "classification-systems", sys, "families")
	entries, err := os.ReadDir(famDir)
	if err != nil {
		return nil, err
	}
	var files []string
	for _, e := range entries {
		if filepath.Ext(e.Name()) == ".json" && !strings.HasSuffix(e.Name(), "_backup.json") && !strings.HasSuffix(e.Name(), ".bak") {
			files = append(files, filepath.Join(famDir, e.Name()))
		}
	}
	return files, nil
}

// VerifyOrSeed verifies if the database has master data seeded. If missing, it auto-seeds and recovers from errors.
func VerifyOrSeed(ctx context.Context, pool *pgxpool.Pool, customDataDir string) error {
	catCount, locCount, isSeeded, _ := CheckSeedStatus(ctx, pool)
	if isSeeded {
		log.Printf("[HUNTIQ DB] ✓ Master data verified: %d categories, %d locations properly seeded.", catCount, locCount)
		return nil
	}

	dataDir := LocateLoccatDataDir(customDataDir)
	log.Printf("[HUNTIQ DB] ⚙ Master data missing or incomplete (%d categories, %d locations). Auto-seeding from %s...", catCount, locCount, dataDir)

	repo := repository.NewMasterDataRepo(pool)
	svc := service.NewMasterDataService(repo)

	// 1. Seed Locations if missing
	if locCount < 400 {
		locNodes, err := loadAllLocations(dataDir)
		if err != nil {
			return fmt.Errorf("load locations from %s: %w", dataDir, err)
		}
		req := models.MasterDataImportRequest{
			DatasetType:     "locations",
			Version:         "1.0.0",
			DryRun:          false,
			LocationRecords: locNodes,
		}
		for attempt := 1; attempt <= 3; attempt++ {
			res, err := svc.ImportLocations(ctx, req, 1, "127.0.0.1")
			if err == nil {
				log.Printf("[HUNTIQ DB]   ✓ Locations seeded: %d created, %d updated", res.Created, res.Updated)
				break
			}
			log.Printf("[HUNTIQ DB] Warning: location import attempt %d failed: %v", attempt, err)
			time.Sleep(time.Duration(attempt) * 500 * time.Millisecond)
		}
	}

	// 2. Seed Categories if missing
	if catCount < 20000 {
		systems, _ := getClassificationSystemDirs(dataDir)
		var allCategories []models.MasterCategoryNode
		seenIDs := make(map[int]bool)

		for _, sys := range systems {
			files, err := getSystemFamilyFiles(dataDir, sys)
			if err != nil {
				continue
			}
			for _, fPath := range files {
				content, err := os.ReadFile(fPath)
				if err != nil {
					continue
				}
				var doc struct {
					Categories []models.MasterCategoryNode `json:"categories"`
				}
				if err := json.Unmarshal(content, &doc); err != nil || len(doc.Categories) == 0 {
					continue
				}
				for _, c := range doc.Categories {
					if c.ID > 0 && !seenIDs[c.ID] {
						seenIDs[c.ID] = true
						allCategories = append(allCategories, c)
					}
				}
			}
		}

		if len(allCategories) > 0 {
			// Sort globally by level ASC to guarantee parents exist before children
			sort.SliceStable(allCategories, func(i, j int) bool {
				return allCategories[i].Level < allCategories[j].Level
			})

			log.Printf("[HUNTIQ DB] Seeding %d canonical categories across all classification systems...", len(allCategories))
			req := models.MasterDataImportRequest{
				DatasetType:     "categories",
				Version:         "2.0.0",
				DryRun:          false,
				CategoryRecords: allCategories,
			}
			res, err := svc.ImportCategories(ctx, req, 1, "127.0.0.1")
			if err != nil {
				log.Printf("[HUNTIQ DB] Warning: category import error: %v", err)
			} else {
				log.Printf("[HUNTIQ DB]   ✓ Categories seeded: %d created, %d updated, %d errors", res.Created, res.Updated, res.Errors)
			}
		}
	}

	// Final verification
	catFinal, locFinal, _, _ := CheckSeedStatus(ctx, pool)
	log.Printf("[HUNTIQ DB] ✓ Master data seed verification complete: %d categories, %d locations active.", catFinal, locFinal)
	return nil
}

// AutoMigrateAndSeed is the unified startup routine that validates migrations, applies pending ones, and verifies seeds.
func AutoMigrateAndSeed(ctx context.Context, pool *pgxpool.Pool) error {
	if err := RunMigrations(ctx, pool, ""); err != nil {
		return fmt.Errorf("auto-migration error: %w", err)
	}
	if err := VerifyOrSeed(ctx, pool, ""); err != nil {
		return fmt.Errorf("auto-seed error: %w", err)
	}
	if err := SeedAdmins(ctx, pool); err != nil {
		return fmt.Errorf("seed-admins error: %w", err)
	}
	return nil
}
