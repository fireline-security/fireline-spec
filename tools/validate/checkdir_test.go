package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

func mustCompileSchema(t *testing.T) *jsonschema.Schema {
	t.Helper()
	compiler := jsonschema.NewCompiler()
	schema, err := compiler.Compile("../../schemas/v1/observation.schema.json")
	if err != nil {
		t.Fatalf("compile schema: %v", err)
	}
	return schema
}

func writeFixture(t *testing.T, dir, name, body string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0o600); err != nil {
		t.Fatalf("write fixture: %v", err)
	}
}

// TestCheckDir_UnexpectedlyInvalid and the tests below call checkDir
// directly with temp-dir fixtures: the real fixtures/ directory only hits
// the happy path, so validateAll's normal run never exercises checkDir's
// mismatch/decode/walk error branches.
func TestCheckDir_UnexpectedlyInvalid(t *testing.T) {
	schema := mustCompileSchema(t)
	dir := t.TempDir()
	writeFixture(t, dir, "bad.json", `{"schema_version": "v2"}`)

	if failures := checkDir(schema, dir, true); failures != 1 {
		t.Errorf("got %d failures, want 1", failures)
	}
}

func TestCheckDir_UnexpectedlyValid(t *testing.T) {
	schema := mustCompileSchema(t)
	dir := t.TempDir()
	writeFixture(t, dir, "good.json", `{
		"schema_version": "v1",
		"source": {"tool": "trivy", "rule_id": "x"},
		"identity_components": {"package": "openssl"},
		"severity_raw": "HIGH",
		"raw_payload": {}
	}`)

	if failures := checkDir(schema, dir, false); failures != 1 {
		t.Errorf("got %d failures, want 1", failures)
	}
}

func TestCheckDir_MalformedJSON(t *testing.T) {
	schema := mustCompileSchema(t)
	dir := t.TempDir()
	writeFixture(t, dir, "malformed.json", "not json")

	if failures := checkDir(schema, dir, true); failures != 1 {
		t.Errorf("got %d failures, want 1", failures)
	}
}

func TestCheckDir_UnopenableFile(t *testing.T) {
	schema := mustCompileSchema(t)
	dir := t.TempDir()
	path := filepath.Join(dir, "unreadable.json")
	if err := os.WriteFile(path, []byte(`{}`), 0o000); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	if failures := checkDir(schema, dir, true); failures != 1 {
		t.Errorf("got %d failures, want 1 (file exists but can't be opened)", failures)
	}
}

func TestCheckDir_NonexistentDir(t *testing.T) {
	schema := mustCompileSchema(t)
	missing := filepath.Join(t.TempDir(), "does-not-exist")

	if failures := checkDir(schema, missing, true); failures != 1 {
		t.Errorf("got %d failures, want 1", failures)
	}
}
