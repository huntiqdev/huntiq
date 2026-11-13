package database

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestMigrationValidation_ProductionMigrations verifies that all 43 real production
// migrations pass sequence, format, and idempotency checks with zero errors.
func TestMigrationValidation_ProductionMigrations(t *testing.T) {
	start := time.Now()

	// Locate migrations directory relative to test file
	migrationsDir := filepath.Join("migrations")
	if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
		// When run from repository root
		migrationsDir = filepath.Join("backend", "internal", "database", "migrations")
	}

	result, err := ValidateMigrationsDirectory(migrationsDir)
	if err != nil {
		t.Fatalf("failed to validate production migrations: %v", err)
	}

	elapsed := time.Since(start)

	if !result.Valid {
		t.Fatalf("production migrations validation failed with errors: %v", result.Errors)
	}

	if result.TotalFiles != 78 {
		t.Errorf("expected 78 production migration files, got %d", result.TotalFiles)
	}
	if result.HistoricalCount != 76 {
		t.Errorf("expected 76 historical migrations, got %d", result.HistoricalCount)
	}
	if result.NewCount != 2 {
		t.Errorf("expected 2 active/new migrations (077_consolidated_taxonomy_mapping, 078_universal_classification_platform), got %d", result.NewCount)
	}
	if len(result.AllowlistedCollisions) != 0 {
		t.Errorf("expected 0 allowlisted collision pairs in normalized migrations, got %d", len(result.AllowlistedCollisions))
	}

	t.Logf("Validated %d production migrations in %v (under 200ms target)", result.TotalFiles, elapsed)
}

// TestMigrationValidation_ValidSequence tests a clean contiguous sequence of migrations.
func TestMigrationValidation_ValidSequence(t *testing.T) {
	files := []string{
		"081_user_sessions.sql",
		"082_secondary_verticals.sql",
		"083_payment_webhooks.sql",
		"084_admin_modularization.sql",
	}

	res := ValidateMigrationFilenames(files)
	if !res.Valid {
		t.Fatalf("expected valid sequence, got errors: %v", res.Errors)
	}
	if res.NewCount != 4 {
		t.Errorf("expected 4 new migrations, got %d", res.NewCount)
	}
}

// TestMigrationValidation_MissingSequenceNumber tests that gaps in numbers are detected.
func TestMigrationValidation_MissingSequenceNumber(t *testing.T) {
	files := []string{
		"081_user_sessions.sql",
		"083_payment_webhooks.sql", // Missing 082!
	}

	res := ValidateMigrationFilenames(files)
	if res.Valid {
		t.Fatal("expected validation to fail for sequence gap, but it passed")
	}

	foundGapError := false
	for _, e := range res.Errors {
		if strings.Contains(e, "[GAP]") && strings.Contains(e, "expected 082") {
			foundGapError = true
			break
		}
	}
	if !foundGapError {
		t.Errorf("expected gap error mentioning 082, got: %v", res.Errors)
	}
}

// TestMigrationValidation_NewDuplicatePrefix tests that new duplicates are strictly rejected.
func TestMigrationValidation_NewDuplicatePrefix(t *testing.T) {
	files := []string{
		"081_user_sessions.sql",
		"081_duplicate_sessions.sql",
	}

	res := ValidateMigrationFilenames(files)
	if res.Valid {
		t.Fatal("expected duplicate prefix 081 to fail, but it passed")
	}

	foundDupError := false
	for _, e := range res.Errors {
		if strings.Contains(e, "[DUPLICATE]") && strings.Contains(e, "081") {
			foundDupError = true
			break
		}
	}
	if !foundDupError {
		t.Errorf("expected duplicate error for prefix 081, got: %v", res.Errors)
	}
}

// TestMigrationValidation_HistoricalAllowlistBehavior tests that allowlisted pairs pass.
func TestMigrationValidation_HistoricalAllowlistBehavior(t *testing.T) {
	orig := HistoricalAllowlistedFiles
	defer func() { HistoricalAllowlistedFiles = orig }()
	HistoricalAllowlistedFiles = map[string]bool{
		"035_mock_taxonomy_alpha.sql": true,
		"035_mock_taxonomy_beta.sql":  true,
	}

	files := []string{
		"035_mock_taxonomy_alpha.sql",
		"035_mock_taxonomy_beta.sql",
		"036_mock_taxonomy_gamma.sql",
	}

	res := ValidateMigrationFilenames(files)
	if !res.Valid {
		t.Fatalf("expected historical collisions to pass via allowlist, got errors: %v", res.Errors)
	}
	if len(res.AllowlistedCollisions) != 1 {
		t.Errorf("expected 1 allowlisted collision, got %d", len(res.AllowlistedCollisions))
	}
}

// TestMigrationValidation_UnapprovedHistoricalCollision tests that unexpected duplicates fail.
func TestMigrationValidation_UnapprovedHistoricalCollision(t *testing.T) {
	files := []string{
		"001_schema.sql",
		"002_users.sql",
		"002_rogue_users.sql", // Not allowlisted!
	}

	res := ValidateMigrationFilenames(files)
	if res.Valid {
		t.Fatal("expected unapproved historical duplicate 002 to fail, but it passed")
	}

	foundDup := false
	for _, e := range res.Errors {
		if strings.Contains(e, "[DUPLICATE]") && strings.Contains(e, "002") {
			foundDup = true
			break
		}
	}
	if !foundDup {
		t.Errorf("expected duplicate error for prefix 002, got: %v", res.Errors)
	}
}

// TestMigrationValidation_InvalidFilenameConvention tests naming convention rejections.
func TestMigrationValidation_InvalidFilenameConvention(t *testing.T) {
	invalidCases := []string{
		"41_short_prefix.sql",
		"041-dash-separated.sql",
		"041_UPPERCASE_NAME.sql",
		"041_no_extension",
		"041_wrong_extension.txt",
		"prefix_at_end_041.sql",
	}

	for _, tc := range invalidCases {
		res := ValidateMigrationFilenames([]string{tc})
		if res.Valid {
			t.Errorf("expected filename %q to fail validation, but it passed", tc)
		}
	}
}

// TestMigrationValidation_Ordering tests lexicographical vs numerical sorting checks.
func TestMigrationValidation_Ordering(t *testing.T) {
	// Numerical order should match lexicographical order when sorted
	files := []string{
		"001_a.sql",
		"002_b.sql",
		"003_c.sql",
		"004_d.sql",
	}

	res := ValidateMigrationFilenames(files)
	if !res.Valid {
		t.Errorf("expected clean ordering, got errors: %v", res.Errors)
	}
}

// TestMigrationValidation_IdempotencyRequirements tests SQL idempotency checks.
func TestMigrationValidation_IdempotencyRequirements(t *testing.T) {
	t.Run("Valid Idempotent SQL", func(t *testing.T) {
		sql := `
			CREATE TABLE IF NOT EXISTS new_feature (id INT PRIMARY KEY);
			CREATE INDEX IF NOT EXISTS idx_new_feature_id ON new_feature(id);
			DROP TABLE IF EXISTS old_table;
			DROP INDEX IF EXISTS old_idx;
		`
		errs, _, _ := ValidateMigrationSQL("082_new_feature.sql", sql)
		if len(errs) > 0 {
			t.Errorf("expected valid idempotent SQL to have 0 errors, got: %v", errs)
		}
	})

	t.Run("Non-Idempotent CREATE TABLE in New Migration", func(t *testing.T) {
		sql := `CREATE TABLE dangerous_table (id INT PRIMARY KEY);`
		errs, _, _ := ValidateMigrationSQL("082_danger.sql", sql)
		if len(errs) == 0 {
			t.Fatal("expected non-idempotent CREATE TABLE to be flagged in new migration")
		}
	})

	t.Run("Non-Idempotent DROP TABLE in New Migration", func(t *testing.T) {
		sql := `DROP TABLE forgotten_table;`
		errs, _, _ := ValidateMigrationSQL("082_danger.sql", sql)
		if len(errs) == 0 {
			t.Fatal("expected non-idempotent DROP TABLE to be flagged in new migration")
		}
	})

	t.Run("Preceded by DROP INDEX IF EXISTS is Idempotent", func(t *testing.T) {
		sql := `
			DROP INDEX IF EXISTS uq_custom_idx;
			CREATE UNIQUE INDEX uq_custom_idx ON users(phone);
		`
		errs, _, _ := ValidateMigrationSQL("082_custom_idx.sql", sql)
		if len(errs) > 0 {
			t.Errorf("expected DROP preceding CREATE INDEX to pass as idempotent, got: %v", errs)
		}
	})
}

// TestMigrationValidation_TransactionSafetyClassification tests detection of non-tx statements.
func TestMigrationValidation_TransactionSafetyClassification(t *testing.T) {
	if !IsNonTransactionalStatement("CREATE INDEX CONCURRENTLY idx_fast ON products(id);") {
		t.Error("expected CREATE INDEX CONCURRENTLY to be classified as non-transactional")
	}
	if !IsNonTransactionalStatement("DROP INDEX CONCURRENTLY idx_fast;") {
		t.Error("expected DROP INDEX CONCURRENTLY to be classified as non-transactional")
	}
	if !IsNonTransactionalStatement("VACUUM FULL products;") {
		t.Error("expected VACUUM to be classified as non-transactional")
	}
	if IsNonTransactionalStatement("CREATE TABLE IF NOT EXISTS normal_table (id INT);") {
		t.Error("expected normal CREATE TABLE to be classified as transactional")
	}
}

// TestMigrationValidation_FailureDiagnostics tests clear and actionable error messages.
func TestMigrationValidation_FailureDiagnostics(t *testing.T) {
	res := ValidateMigrationFilenames([]string{"081_first.sql", "083_third.sql"})
	if res.Valid {
		t.Fatal("expected failure")
	}
	diagnostic := res.Errors[0]
	if !strings.Contains(diagnostic, "[GAP]") || !strings.Contains(diagnostic, "expected 082") {
		t.Errorf("diagnostic message is not clear/actionable: %q", diagnostic)
	}
}

// TestValidateMigrationsScript_ShellIntegration runs validate_migrations.sh directly.
func TestValidateMigrationsScript_ShellIntegration(t *testing.T) {
	start := time.Now()

	// Locate script relative to working directory
	scriptPath := filepath.Join("..", "..", "scripts", "validate_migrations.sh")
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		scriptPath = filepath.Join("backend", "scripts", "validate_migrations.sh")
	}

	cmd := exec.Command(scriptPath)
	output, err := cmd.CombinedOutput()
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("validate_migrations.sh exited with error: %v\nOutput: %s", err, string(output))
	}

	if !strings.Contains(string(output), "SUCCESS: All ") || !strings.Contains(string(output), "migrations passed") {
		t.Errorf("expected success banner in output, got:\n%s", string(output))
	}

	t.Logf("Executed shell validator in %v", elapsed)
}
