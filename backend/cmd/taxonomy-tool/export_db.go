package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// DBCategoryRecord represents a row from category_hierarchy with complete metadata
type DBCategoryRecord struct {
	ID              int                    `json:"id"`
	ParentID        *int                   `json:"parent_id"`
	Slug            string                 `json:"slug"`
	Code            string                 `json:"code,omitempty"`
	Level           int                    `json:"level"`
	Path            string                 `json:"path"`
	CanonicalName   string                 `json:"canonical_name"`
	NameEn          string                 `json:"name_en"`
	NameAm          string                 `json:"name_am"`
	NameOm          string                 `json:"name_om"`
	NameTi          string                 `json:"name_ti"`
	Definition      string                 `json:"definition"`
	DisplayOrder    int                    `json:"display_order"`
	Icon            string                 `json:"icon,omitempty"`
	ImageURL        string                 `json:"image_url,omitempty"`
	IsActive        bool                   `json:"is_active"`
	IsLeaf          bool                   `json:"is_leaf"`
	MetaTitle       string                 `json:"meta_title,omitempty"`
	MetaDescription string                 `json:"meta_description,omitempty"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
	Translations    map[string]string      `json:"translations"`
	Aliases         map[string][]string    `json:"aliases"`
	Attributes      []DBCategoryAttribute  `json:"attributes,omitempty"`
	Children        []*DBCategoryRecord    `json:"children,omitempty"`
}

type DBCategoryAttribute struct {
	AttributeID     int                    `json:"attribute_id"`
	Name            string                 `json:"name"`
	Slug            string                 `json:"slug"`
	DataType        string                 `json:"data_type"`
	Unit            string                 `json:"unit,omitempty"`
	IsRequired      bool                   `json:"is_required"`
	IsFilterable    bool                   `json:"is_filterable"`
	IsSearchable    bool                   `json:"is_searchable"`
	DisplayOrder    int                    `json:"display_order"`
	ValidationRules map[string]interface{} `json:"validation_rules,omitempty"`
}

type DBExportManifest struct {
	ExportTimestamp      time.Time          `json:"export_timestamp"`
	DatabaseName         string             `json:"database_name"`
	TotalCategoriesInDB  int                `json:"total_categories_in_db"`
	TotalExported        int                `json:"total_exported"`
	RootFamiliesCount    int                `json:"root_families_count"`
	SuperFamiliesCount   int                `json:"super_families_count"`
	ThirdLevelCount      int                `json:"third_level_count"`
	FourthLevelCount     int                `json:"fourth_level_count"`
	TotalAliasesExported int                `json:"total_aliases_exported"`
	AliasesByLanguage    map[string]int     `json:"aliases_by_language"`
	DomainFilesCount     int                `json:"domain_files_count"`
	DomainFiles          []string           `json:"domain_files"`
	ReconciliationResult string             `json:"reconciliation_result"`
}

func runExportAllDBJSON() {
	exportDir := flag.String("dir", "data/taxonomy/db-canonical-export", "Output directory for exported DB domain JSON files")
	flag.CommandLine.Parse(os.Args[2:])

	ctx := context.Background()
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:gedam@localhost:5432/hulubedej_db?sslmode=disable"
	}

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	if err := os.MkdirAll(*exportDir, 0755); err != nil {
		log.Fatalf("Failed to create export directory: %v", err)
	}

	fmt.Printf("[REVERSE-EXPORT] Starting full reverse export from database: %s\n", dbURL)

	// 1. Fetch total count from DB
	var dbTotal int
	err = pool.QueryRow(ctx, "SELECT count(*) FROM category_hierarchy").Scan(&dbTotal)
	if err != nil {
		log.Fatalf("Failed to query total categories: %v", err)
	}
	fmt.Printf("[REVERSE-EXPORT] Database category_hierarchy contains %d categories.\n", dbTotal)

	// 2. Fetch all aliases grouped by category_id and language
	aliasMap := make(map[int]map[string][]string)
	aliasCountByLang := make(map[string]int)
	totalAliases := 0

	aliasRows, err := pool.Query(ctx, "SELECT category_id, language, alias FROM category_aliases ORDER BY category_id, language, alias")
	if err == nil {
		defer aliasRows.Close()
		for aliasRows.Next() {
			var catID int
			var lang, alias string
			if err := aliasRows.Scan(&catID, &lang, &alias); err == nil {
				if aliasMap[catID] == nil {
					aliasMap[catID] = make(map[string][]string)
				}
				aliasMap[catID][lang] = append(aliasMap[catID][lang], alias)
				aliasCountByLang[lang]++
				totalAliases++
			}
		}
	}
	fmt.Printf("[REVERSE-EXPORT] Loaded %d aliases across languages: %+v\n", totalAliases, aliasCountByLang)

	// 3. Fetch attributes grouped by category_id
	attrMap := make(map[int][]DBCategoryAttribute)
	attrRows, err := pool.Query(ctx, `
		SELECT 
			category_id, attribute_id, attribute_name, attribute_slug, 
			data_type, COALESCE(unit, ''), is_required, is_filterable, is_searchable, 
			display_order, validation_rules
		FROM category_applicable_attributes
		ORDER BY category_id, display_order
	`)
	if err == nil {
		defer attrRows.Close()
		for attrRows.Next() {
			var catID int
			var a DBCategoryAttribute
			var valRules []byte
			if err := attrRows.Scan(
				&catID, &a.AttributeID, &a.Name, &a.Slug,
				&a.DataType, &a.Unit, &a.IsRequired, &a.IsFilterable, &a.IsSearchable,
				&a.DisplayOrder, &valRules,
			); err == nil {
				if len(valRules) > 0 {
					_ = json.Unmarshal(valRules, &a.ValidationRules)
				}
				attrMap[catID] = append(attrMap[catID], a)
			}
		}
	}
	fmt.Printf("[REVERSE-EXPORT] Loaded applicable attributes for %d categories.\n", len(attrMap))

	// 4. Fetch all category records
	catQuery := `
		SELECT 
			id, parent_id, slug, level, COALESCE(path, ''),
			COALESCE(name_en, ''), COALESCE(name_am, ''), COALESCE(name_om, ''), COALESCE(name_ti, ''),
			COALESCE(description, ''), display_order, COALESCE(icon, ''), COALESCE(image_url, ''),
			is_active, is_leaf, COALESCE(meta_title, ''), COALESCE(meta_description, ''),
			created_at, updated_at
		FROM category_hierarchy
		ORDER BY level ASC, display_order ASC, id ASC
	`
	rows, err := pool.Query(ctx, catQuery)
	if err != nil {
		log.Fatalf("Failed to fetch categories: %v", err)
	}
	defer rows.Close()

	recordMap := make(map[int]*DBCategoryRecord)
	var flatRecords []*DBCategoryRecord
	var rootRecords []*DBCategoryRecord

	levelCounts := make(map[int]int)

	for rows.Next() {
		r := &DBCategoryRecord{
			Translations: make(map[string]string),
			Aliases:      make(map[string][]string),
		}
		var nameEn, nameAm, nameOm, nameTi, desc, icon, img, mTitle, mDesc string
		if err := rows.Scan(
			&r.ID, &r.ParentID, &r.Slug, &r.Level, &r.Path,
			&nameEn, &nameAm, &nameOm, &nameTi,
			&desc, &r.DisplayOrder, &icon, &img,
			&r.IsActive, &r.IsLeaf, &mTitle, &mDesc,
			&r.CreatedAt, &r.UpdatedAt,
		); err != nil {
			log.Fatalf("Failed scanning category row: %v", err)
		}

		r.NameEn = nameEn
		r.CanonicalName = nameEn
		r.NameAm = nameAm
		r.NameOm = nameOm
		r.NameTi = nameTi
		r.Definition = desc
		r.Icon = icon
		r.ImageURL = img
		r.MetaTitle = mTitle
		r.MetaDescription = mDesc

		r.Translations["en"] = nameEn
		if nameAm != "" {
			r.Translations["am"] = nameAm
		}
		if nameOm != "" {
			r.Translations["or"] = nameOm
		}
		if nameTi != "" {
			r.Translations["ti"] = nameTi
		}

		if aliases, ok := aliasMap[r.ID]; ok {
			r.Aliases = aliases
		}
		if attrs, ok := attrMap[r.ID]; ok {
			r.Attributes = attrs
		}

		levelCounts[r.Level]++
		recordMap[r.ID] = r
		flatRecords = append(flatRecords, r)

		if r.ParentID == nil {
			rootRecords = append(rootRecords, r)
		}
	}

	// 5. Assemble hierarchy tree
	for _, rec := range flatRecords {
		if rec.ParentID != nil {
			if parent, exists := recordMap[*rec.ParentID]; exists {
				parent.Children = append(parent.Children, rec)
			}
		}
	}

	fmt.Printf("[REVERSE-EXPORT] Assembled %d Root Family trees across %d total categories.\n", len(rootRecords), len(flatRecords))

	// 6. Write each Root Family domain file
	var domainFileNames []string
	for idx, root := range rootRecords {
		domainSlug := root.Slug
		if domainSlug == "" {
			domainSlug = fmt.Sprintf("root_%d", root.ID)
		}
		filename := fmt.Sprintf("%03d_%s.json", idx+1, strings.ReplaceAll(domainSlug, "-", "_"))
		filePath := filepath.Join(*exportDir, filename)

		fileObj := struct {
			SchemaVersion      string            `json:"schema_version"`
			TaxonomyVersion    string            `json:"taxonomy_version"`
			ExportTimestamp    time.Time         `json:"export_timestamp"`
			SourceLanguage     string            `json:"source_language"`
			SupportedLanguages []string          `json:"supported_languages"`
			DomainSlug         string            `json:"domain_slug"`
			RootFamily         *DBCategoryRecord `json:"root_family"`
		}{
			SchemaVersion:      "2.0.0",
			TaxonomyVersion:    "2.0.0",
			ExportTimestamp:    time.Now().UTC(),
			SourceLanguage:     "en",
			SupportedLanguages: []string{"en", "am", "or", "ti"},
			DomainSlug:         root.Slug,
			RootFamily:         root,
		}

		data, err := json.MarshalIndent(fileObj, "", "  ")
		if err != nil {
			log.Fatalf("Failed marshaling domain %s: %v", filename, err)
		}
		if err := os.WriteFile(filePath, data, 0644); err != nil {
			log.Fatalf("Failed writing domain %s: %v", filename, err)
		}
		domainFileNames = append(domainFileNames, filename)
	}

	// 7. Write flat snapshot file with ALL categories
	snapshotPath := filepath.Join(*exportDir, "all_categories_snapshot.json")
	flatData, err := json.MarshalIndent(flatRecords, "", "  ")
	if err != nil {
		log.Fatalf("Failed marshaling snapshot: %v", err)
	}
	if err := os.WriteFile(snapshotPath, flatData, 0644); err != nil {
		log.Fatalf("Failed writing snapshot: %v", err)
	}
	fmt.Printf("[REVERSE-EXPORT] Saved flat snapshot: %s (%d records)\n", snapshotPath, len(flatRecords))

	// 8. Write manifest and reconciliation report
	reconcileStatus := "EXACT_MATCH"
	if len(flatRecords) != dbTotal {
		reconcileStatus = fmt.Sprintf("MISMATCH: db has %d, exported %d", dbTotal, len(flatRecords))
	}

	manifest := DBExportManifest{
		ExportTimestamp:      time.Now().UTC(),
		DatabaseName:         "hulubedej_db",
		TotalCategoriesInDB:  dbTotal,
		TotalExported:        len(flatRecords),
		RootFamiliesCount:    levelCounts[0],
		SuperFamiliesCount:   levelCounts[1],
		ThirdLevelCount:      levelCounts[2],
		FourthLevelCount:     levelCounts[3],
		TotalAliasesExported: totalAliases,
		AliasesByLanguage:    aliasCountByLang,
		DomainFilesCount:     len(domainFileNames),
		DomainFiles:          domainFileNames,
		ReconciliationResult: reconcileStatus,
	}

	manifestPath := filepath.Join(*exportDir, "manifest.json")
	manifestData, _ := json.MarshalIndent(manifest, "", "  ")
	_ = os.WriteFile(manifestPath, manifestData, 0644)

	fmt.Printf("\n=== REVERSE EXPORT COMPLETE ===\n")
	fmt.Printf("Total DB Categories:     %d\n", dbTotal)
	fmt.Printf("Total Exported Records:  %d\n", len(flatRecords))
	fmt.Printf("Level 1 (Roots):         %d\n", levelCounts[0])
	fmt.Printf("Level 2 (SuperFamilies): %d\n", levelCounts[1])
	fmt.Printf("Level 3 (ThirdLevel):    %d\n", levelCounts[2])
	fmt.Printf("Level 4 (FourthLevel):   %d\n", levelCounts[3])
	fmt.Printf("Total Aliases:           %d\n", totalAliases)
	fmt.Printf("Domain Files Created:    %d\n", len(domainFileNames))
	fmt.Printf("Reconciliation Status:   %s ✅\n", reconcileStatus)
}
