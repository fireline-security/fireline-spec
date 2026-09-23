package main

import (
	"os"
	"path/filepath"
	"testing"
)

// TestValidateAll_AllFixturesMatchExpectations exercises the exact same
// fixture-walking/validate logic `go run ./tools/validate` (and `task
// validate`) run: valid/ fixtures must validate cleanly, invalid/ fixtures
// must fail for the reason their filename claims. This makes the existing
// check go test-able (and coverage-measurable) without changing its
// runtime behavior; `task validate`'s `go run` entry point is unchanged.
func TestValidateAll_AllFixturesMatchExpectations(t *testing.T) {
	if err := validateAll("../../schemas/v1/observation.schema.json", "../../fixtures/v1/observation"); err != nil {
		t.Fatalf("validateAll: %v", err)
	}
}

func TestValidateAll_MissingSchema(t *testing.T) {
	if err := validateAll("../../schemas/v1/does-not-exist.json", "../../fixtures/v1/observation"); err == nil {
		t.Fatal("expected an error for a missing schema file, got nil")
	}
}

// TestValidateAll_ReportsFixtureMismatch exercises validateAll's own
// failures>0 branch directly, distinct from checkDir's own (see
// checkdir_test.go); the real fixtures/ tree never has a mismatch, so
// TestValidateAll_AllFixturesMatchExpectations alone can't reach this.
func TestValidateAll_ReportsFixtureMismatch(t *testing.T) {
	dir := t.TempDir()
	validDir := filepath.Join(dir, "valid")
	invalidDir := filepath.Join(dir, "invalid")
	if err := os.MkdirAll(validDir, 0o750); err != nil {
		t.Fatalf("mkdir valid: %v", err)
	}
	if err := os.MkdirAll(invalidDir, 0o750); err != nil {
		t.Fatalf("mkdir invalid: %v", err)
	}
	writeFixture(t, validDir, "actually-invalid.json", `{"schema_version": "v2"}`)

	if err := validateAll("../../schemas/v1/observation.schema.json", dir); err == nil {
		t.Fatal("expected an error for a fixture that doesn't match its own directory's expectation, got nil")
	}
}

// TestRun_Succeeds exercises run() itself, whose schemaPath/fixturesRoot
// constants are relative to the repo root the way `go run ./tools/validate`
// is always invoked from, not to this test binary's own package
// directory, so it needs a real (restored) chdir to match.
func TestRun_Succeeds(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd: %v", err)
	}
	if err := os.Chdir("../.."); err != nil {
		t.Fatalf("Chdir to repo root: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(wd); err != nil {
			t.Fatalf("restore Chdir: %v", err)
		}
	})

	if err := run(); err != nil {
		t.Fatalf("run: %v", err)
	}
}
