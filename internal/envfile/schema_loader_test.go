package envfile

import (
	"os"
	"path/filepath"
	"testing"
)

func writeSchemaFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "schema.json")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestLoadSchema_BasicFields(t *testing.T) {
	p := writeSchemaFile(t, `[
		{"key":"HOST","required":true},
		{"key":"PORT","required":true,"pattern":"^\\d+$"}
	]`)
	schema, err := LoadSchema(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(schema.Fields) != 2 {
		t.Fatalf("expected 2 fields, got %d", len(schema.Fields))
	}
	if schema.Fields[1].Pattern == nil {
		t.Error("expected pattern to be compiled")
	}
}

func TestLoadSchema_FileNotFound(t *testing.T) {
	_, err := LoadSchema("/nonexistent/schema.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoadSchema_InvalidJSON(t *testing.T) {
	p := writeSchemaFile(t, `{not valid json}`)
	_, err := LoadSchema(p)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestLoadSchema_InvalidPattern(t *testing.T) {
	p := writeSchemaFile(t, `[{"key":"X","required":true,"pattern":"[invalid"}]`)
	_, err := LoadSchema(p)
	if err == nil {
		t.Fatal("expected error for invalid regex pattern")
	}
}

func TestLoadSchema_RoundTripValidation(t *testing.T) {
	p := writeSchemaFile(t, `[
		{"key":"APP_ENV","required":true,"pattern":"^(development|staging|production)$"}
	]`)
	schema, err := LoadSchema(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	env := map[string]string{"APP_ENV": "staging"}
	violations := ValidateSchema(env, schema)
	if len(violations) != 0 {
		t.Fatalf("expected no violations, got %v", violations)
	}
	env["APP_ENV"] = "unknown"
	violations = ValidateSchema(env, schema)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
}
