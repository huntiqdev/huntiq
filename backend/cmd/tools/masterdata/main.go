package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/your-org/used-tools-marketplace/internal/config"
	"github.com/your-org/used-tools-marketplace/internal/models"
	"github.com/your-org/used-tools-marketplace/internal/repository"
	"github.com/your-org/used-tools-marketplace/internal/service"
)

func main() {
	action := flag.String("action", "status", "Action to perform: status, validate, import, export")
	dataset := flag.String("dataset", "all", "Dataset scope: all, locations, categories")
	dryRun := flag.Bool("dry-run", true, "Execute as dry-run without database mutations")
	dataDir := flag.String("dir", "data/master", "Path to master data directory")
	flag.Parse()

	cfg := config.Load()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	pool, err := pgxpool.New(ctx, cfg.DSN())
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer pool.Close()

	repo := repository.NewMasterDataRepo(pool)
	svc := service.NewMasterDataService(repo)

	switch *action {
	case "status":
		status, err := svc.GetStatus(ctx)
		if err != nil {
			log.Fatalf("Failed to fetch status: %v", err)
		}
		fmt.Printf("=== HUNTIQ Master Data Status ===\n")
		fmt.Printf("Locations: %d (Active: %d)\n", status.LocationsCount, status.ActiveLocations)
		fmt.Printf("Categories: %d (Active: %d)\n", status.CategoriesCount, status.ActiveCats)
		fmt.Printf("Languages Supported: %v\n", status.SupportedLangs)
		fmt.Printf("Timestamp: %s\n", status.LastImportAt.Format(time.RFC3339))

	case "validate":
		fmt.Println("=== Validating Master Data Files ===")
		locCount := validateLocations(svc, *dataDir)
		catCount := validateCategories(svc, *dataDir)
		fmt.Printf("Validation passed: %d locations and %d categories verified.\n", locCount, catCount)

	case "import":
		fmt.Printf("=== Importing Master Data (DryRun: %v) ===\n", *dryRun)
		if *dataset == "all" || *dataset == "locations" {
			importLocations(ctx, svc, *dataDir, *dryRun)
		}
		if *dataset == "all" || *dataset == "categories" {
			importCategories(ctx, svc, *dataDir, *dryRun)
		}

	case "export":
		fmt.Println("=== Exporting Master Data from PostgreSQL ===")
		exportLocations(ctx, svc, *dataDir)
		exportCategories(ctx, svc, *dataDir)
		fmt.Println("Export completed successfully.")

	default:
		log.Fatalf("Unknown action '%s'. Use status, validate, import, or export.", *action)
	}
}

func validateLocations(svc *service.MasterDataService, baseDir string) int {
	countryFile := filepath.Join(baseDir, "locations", "countries", "ET", "country.json")
	data, err := os.ReadFile(countryFile)
	if err != nil {
		log.Fatalf("Failed to read %s: %v", countryFile, err)
	}
	var country models.MasterLocationNode
	if err := json.Unmarshal(data, &country); err != nil {
		log.Fatalf("Invalid JSON in %s: %v", countryFile, err)
	}
	diffs := svc.ValidateLocations([]models.MasterLocationNode{country}, nil)
	if len(diffs) > 0 {
		log.Fatalf("Location validation failed on root: %v", diffs)
	}
	return 1
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

func validateCategories(svc *service.MasterDataService, baseDir string) int {
	systems, err := getClassificationSystemDirs(baseDir)
	if err != nil {
		log.Fatalf("Failed to read classification systems: %v", err)
	}
	total := 0
	for _, sys := range systems {
		famDir := filepath.Join(baseDir, "categories", "classification-systems", sys, "families")
		entries, err := os.ReadDir(famDir)
		if err != nil {
			continue
		}
		for _, e := range entries {
			if filepath.Ext(e.Name()) != ".json" {
				continue
			}
			fPath := filepath.Join(famDir, e.Name())
			data, err := os.ReadFile(fPath)
			if err != nil {
				log.Fatalf("Failed to read %s: %v", fPath, err)
			}
			var doc struct {
				Categories []models.MasterCategoryNode `json:"categories"`
			}
			if err := json.Unmarshal(data, &doc); err != nil {
				log.Fatalf("Invalid JSON in %s: %v", fPath, err)
			}
			diffs := svc.ValidateCategories(doc.Categories, nil)
			if len(diffs) > 0 {
				log.Fatalf("Validation error in %s/%s: %v", sys, e.Name(), diffs)
			}
			total += len(doc.Categories)
		}
	}
	return total
}

func importLocations(ctx context.Context, svc *service.MasterDataService, baseDir string, dryRun bool) {
	// Read country
	countryFile := filepath.Join(baseDir, "locations", "countries", "ET", "country.json")
	data, err := os.ReadFile(countryFile)
	if err != nil {
		log.Fatalf("Failed to read country file: %v", err)
	}
	var country models.MasterLocationNode
	_ = json.Unmarshal(data, &country)

	req := models.MasterDataImportRequest{
		DatasetType:     "locations",
		Version:         "1.0.0",
		DryRun:          dryRun,
		LocationRecords: []models.MasterLocationNode{country},
	}

	res, err := svc.ImportLocations(ctx, req, 1, "127.0.0.1")
	if err != nil {
		log.Fatalf("Location import error: %v", err)
	}
	fmt.Printf("Locations: Total=%d, Created=%d, Updated=%d, Unchanged=%d, Errors=%d\n",
		res.TotalRecords, res.Created, res.Updated, res.Unchanged, res.Errors)
}

func importCategories(ctx context.Context, svc *service.MasterDataService, baseDir string, dryRun bool) {
	systems, err := getClassificationSystemDirs(baseDir)
	if err != nil {
		log.Fatalf("Failed to read classification systems: %v", err)
	}
	var allCats []models.MasterCategoryNode
	for _, sys := range systems {
		famDir := filepath.Join(baseDir, "categories", "classification-systems", sys, "families")
		entries, err := os.ReadDir(famDir)
		if err != nil {
			continue
		}
		for _, e := range entries {
			if filepath.Ext(e.Name()) != ".json" {
				continue
			}
			data, err := os.ReadFile(filepath.Join(famDir, e.Name()))
			if err != nil {
				continue
			}
			var doc struct {
				Categories []models.MasterCategoryNode `json:"categories"`
			}
			if err := json.Unmarshal(data, &doc); err == nil {
				allCats = append(allCats, doc.Categories...)
			}
		}
	}

	req := models.MasterDataImportRequest{
		DatasetType:     "categories",
		Version:         "1.0.0",
		DryRun:          dryRun,
		CategoryRecords: allCats,
	}

	res, err := svc.ImportCategories(ctx, req, 1, "127.0.0.1")
	if err != nil {
		log.Fatalf("Category import error: %v", err)
	}
	fmt.Printf("Categories: Total=%d, Created=%d, Updated=%d, Unchanged=%d, Errors=%d\n",
		res.TotalRecords, res.Created, res.Updated, res.Unchanged, res.Errors)
}

func exportLocations(ctx context.Context, svc *service.MasterDataService, baseDir string) {
	locs, err := svc.ExportLocations(ctx, "")
	if err != nil {
		log.Fatalf("Location export error: %v", err)
	}
	fmt.Printf("Exported %d locations from PostgreSQL.\n", len(locs))
}

func exportCategories(ctx context.Context, svc *service.MasterDataService, baseDir string) {
	systems, err := getClassificationSystemDirs(baseDir)
	if err != nil {
		systems = []string{"SYS_PROD", "SYS_SRV", "SYS_EMP", "SYS_GIG", "SYS_MED", "SYS_OPP", "SYS_PRM"}
	}
	total := 0
	for _, sys := range systems {
		cats, err := svc.ExportCategories(ctx, sys, "")
		if err != nil {
			log.Fatalf("Category export error for %s: %v", sys, err)
		}
		total += len(cats)
		fmt.Printf("Classification %s: Exported %d categories from PostgreSQL.\n", sys, len(cats))
	}
	fmt.Printf("Total exported categories: %d across %d systems.\n", total, len(systems))
}
