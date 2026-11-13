package database

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var (
	// FilenameRegex matches standard 3-digit prefixed migration filenames.
	FilenameRegex = regexp.MustCompile(`^[0-9]{3}_[a-z0-9_]+\.sql$`)

	// HistoricalAllowlistedFiles holds any approved migrations with duplicate prefixes.
	// In HUNTIQ canonical normalized migrations, all collisions are eliminated (0 collisions).
	HistoricalAllowlistedFiles = map[string]bool{}

	createTableRegex = regexp.MustCompile(`(?i)^\s*CREATE\s+(?:UNLOGGED\s+)?TABLE\s+`)
	dropTableRegex   = regexp.MustCompile(`(?i)^\s*DROP\s+TABLE\s+`)
	createIndexRegex = regexp.MustCompile(`(?i)^\s*CREATE\s+(?:UNIQUE\s+)?INDEX\s+`)
	dropIndexRegex   = regexp.MustCompile(`(?i)^\s*DROP\s+INDEX\s+`)

	ifNotExistsRegex = regexp.MustCompile(`(?i)IF\s+NOT\s+EXISTS`)
	ifExistsRegex    = regexp.MustCompile(`(?i)IF\s+EXISTS`)

	nonTxConcurrentRegex = regexp.MustCompile(`(?i)(?:CREATE|DROP)\s+(?:UNIQUE\s+)?INDEX\s+CONCURRENTLY`)
	nonTxVacuumRegex     = regexp.MustCompile(`(?i)\bVACUUM\b`)
	nonTxClusterRegex    = regexp.MustCompile(`(?i)\bCLUSTER\s+[a-zA-Z0-9_]+\s+USING\b`)
)

// MigrationValidationResult holds the findings of a migration check.
type MigrationValidationResult struct {
	Valid                 bool
	TotalFiles            int
	HistoricalCount       int
	NewCount              int
	AllowlistedCollisions []string
	NonTransactionalOps   []string
	Warnings              []string
	Errors                []string
}

// IsHistoricalAllowlisted returns true if the filename is an approved historical collision.
func IsHistoricalAllowlisted(filename string) bool {
	return HistoricalAllowlistedFiles[filename]
}

// IsNonTransactionalStatement checks if SQL contains statements that cannot run inside multi-statement transactions.
func IsNonTransactionalStatement(sql string) bool {
	return nonTxConcurrentRegex.MatchString(sql) ||
		nonTxVacuumRegex.MatchString(sql) ||
		nonTxClusterRegex.MatchString(sql)
}

// ValidateMigrationFilenames checks sequence, formatting, gaps, and collisions.
func ValidateMigrationFilenames(filenames []string) *MigrationValidationResult {
	res := &MigrationValidationResult{
		Valid:                 true,
		AllowlistedCollisions: make([]string, 0),
		NonTransactionalOps:   make([]string, 0),
		Warnings:              make([]string, 0),
		Errors:                make([]string, 0),
	}

	if len(filenames) == 0 {
		res.Valid = false
		res.Errors = append(res.Errors, "no migration files provided")
		return res
	}

	// Sort filenames lexicographically
	sortedFiles := make([]string, len(filenames))
	copy(sortedFiles, filenames)
	sort.Strings(sortedFiles)
	res.TotalFiles = len(sortedFiles)

	seenPrefixes := make(map[string][]string)
	prevPrefixInt := -1

	for _, filename := range sortedFiles {
		// 1. Filename format check
		if !FilenameRegex.MatchString(filename) {
			res.Valid = false
			res.Errors = append(res.Errors, fmt.Sprintf("[FORMAT] invalid filename: %q, must match '^[0-9]{3}_[a-z0-9_]+\\.sql$'", filename))
			continue
		}

		prefix := filename[:3]
		prefixInt, err := strconv.Atoi(prefix)
		if err != nil {
			res.Valid = false
			res.Errors = append(res.Errors, fmt.Sprintf("[FORMAT] non-numeric prefix in filename: %q", filename))
			continue
		}

		if prefixInt < 77 {
			res.HistoricalCount++
		} else {
			res.NewCount++
		}

		// 2. Duplicate prefix check
		existing := seenPrefixes[prefix]
		if len(existing) > 0 {
			allAllowlisted := IsHistoricalAllowlisted(filename)
			for _, prevFile := range existing {
				if !IsHistoricalAllowlisted(prevFile) {
					allAllowlisted = false
					break
				}
			}

			if allAllowlisted {
				res.AllowlistedCollisions = append(res.AllowlistedCollisions, fmt.Sprintf("prefix %s: %s & %s", prefix, existing[0], filename))
			} else {
				res.Valid = false
				res.Errors = append(res.Errors, fmt.Sprintf("[DUPLICATE] unapproved duplicate prefix %q found in %v and %q", prefix, existing, filename))
			}
			seenPrefixes[prefix] = append(existing, filename)
		} else {
			seenPrefixes[prefix] = []string{filename}
		}

		// 3. Sequential ordering and gap check
		if prevPrefixInt >= 0 {
			if prefixInt < prevPrefixInt {
				res.Valid = false
				res.Errors = append(res.Errors, fmt.Sprintf("[ORDER] sequence inversion: %q (prefix %d) appears after prefix %d", filename, prefixInt, prevPrefixInt))
			} else if prefixInt > prevPrefixInt {
				expectedNext := prevPrefixInt + 1
				if prefixInt > expectedNext {
					res.Valid = false
					if prefixInt >= 77 {
						res.Errors = append(res.Errors, fmt.Sprintf("[GAP] missing migration sequence: gap between %03d and %03d (expected %03d)", prevPrefixInt, prefixInt, expectedNext))
					} else {
						res.Errors = append(res.Errors, fmt.Sprintf("[GAP] missing canonical migration sequence: gap between %03d and %03d (expected %03d)", prevPrefixInt, prefixInt, expectedNext))
					}
				}
			}
		}
		prevPrefixInt = prefixInt
	}

	return res
}

// ValidateMigrationSQL checks SQL content for idempotency and transaction safety.
func ValidateMigrationSQL(filename string, content string) ([]string, []string, []string) {
	var errs, warns, nonTx []string

	prefixInt := 0
	if len(filename) >= 3 {
		prefixInt, _ = strconv.Atoi(filename[:3])
	}

	// Remove SQL comments to avoid false positives
	lines := strings.Split(content, "\n")
	var strippedLines []string
	for _, l := range lines {
		trimmed := strings.TrimSpace(l)
		if strings.HasPrefix(trimmed, "--") {
			continue
		}
		// Strip inline comment
		if idx := strings.Index(trimmed, "--"); idx >= 0 {
			trimmed = strings.TrimSpace(trimmed[:idx])
		}
		if trimmed != "" {
			strippedLines = append(strippedLines, trimmed)
		}
	}
	cleanContent := strings.Join(strippedLines, "\n")

	// Check non-transactional statements
	if nonTxConcurrentRegex.MatchString(cleanContent) {
		nonTx = append(nonTx, fmt.Sprintf("[NON-TX] %q: contains CONCURRENT index statement (must execute outside multi-statement transaction)", filename))
	}
	if nonTxVacuumRegex.MatchString(cleanContent) {
		nonTx = append(nonTx, fmt.Sprintf("[NON-TX] %q: contains VACUUM statement", filename))
	}

	// Check idempotency per line
	for _, line := range strippedLines {
		// CREATE TABLE
		if createTableRegex.MatchString(line) && !ifNotExistsRegex.MatchString(line) {
			if prefixInt >= 77 {
				errs = append(errs, fmt.Sprintf("[IDEMPOTENCY] %q: table creation without IF NOT EXISTS: %s", filename, line))
			} else {
				warns = append(warns, fmt.Sprintf("[HISTORICAL] %q: historical table statement without IF NOT EXISTS: %s", filename, line))
			}
		}

		// DROP TABLE
		if dropTableRegex.MatchString(line) && !ifExistsRegex.MatchString(line) {
			if prefixInt >= 77 {
				errs = append(errs, fmt.Sprintf("[IDEMPOTENCY] %q: drop table without IF EXISTS: %s", filename, line))
			} else {
				warns = append(warns, fmt.Sprintf("[HISTORICAL] %q: historical drop table without IF EXISTS: %s", filename, line))
			}
		}

		// CREATE INDEX
		if createIndexRegex.MatchString(line) && !ifNotExistsRegex.MatchString(line) {
			// Extract index name
			words := strings.Fields(line)
			var idxName string
			for i, w := range words {
				if strings.EqualFold(w, "INDEX") && i+1 < len(words) {
					next := words[i+1]
					if strings.EqualFold(next, "CONCURRENTLY") && i+2 < len(words) {
						idxName = words[i+2]
					} else {
						idxName = next
					}
					break
				}
			}

			// If DROP INDEX IF EXISTS <idxName> exists in content, it is idempotent
			dropPattern := regexp.MustCompile(fmt.Sprintf(`(?i)DROP\s+INDEX\s+(?:IF\s+EXISTS\s+)?%s`, regexp.QuoteMeta(idxName)))
			if idxName != "" && !dropPattern.MatchString(cleanContent) {
				if prefixInt >= 77 {
					errs = append(errs, fmt.Sprintf("[IDEMPOTENCY] %q: index %q created without IF NOT EXISTS or preceding DROP INDEX IF EXISTS", filename, idxName))
				} else {
					warns = append(warns, fmt.Sprintf("[HISTORICAL] %q: historical index %q created without IF NOT EXISTS", filename, idxName))
				}
			}
		}
	}

	return errs, warns, nonTx
}

// ValidateMigrationsDirectory validates all SQL files within a directory.
func ValidateMigrationsDirectory(dir string) (*MigrationValidationResult, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read migrations directory: %w", err)
	}

	var sqlFiles []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".sql") {
			sqlFiles = append(sqlFiles, e.Name())
		}
	}

	result := ValidateMigrationFilenames(sqlFiles)

	for _, fname := range sqlFiles {
		filePath := filepath.Join(dir, fname)
		contentBytes, err := os.ReadFile(filePath)
		if err != nil {
			result.Valid = false
			result.Errors = append(result.Errors, fmt.Sprintf("read file %q: %v", fname, err))
			continue
		}

		sqlErrs, sqlWarns, nonTx := ValidateMigrationSQL(fname, string(contentBytes))
		if len(sqlErrs) > 0 {
			result.Valid = false
			result.Errors = append(result.Errors, sqlErrs...)
		}
		result.Warnings = append(result.Warnings, sqlWarns...)
		result.NonTransactionalOps = append(result.NonTransactionalOps, nonTx...)
	}

	return result, nil
}
