package envfile

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScope_ThenWrite_RoundTrip(t *testing.T) {
	env := map[string]string{
		"LOG_LEVEL":        "info",
		"PROD__LOG_LEVEL":  "warn",
		"PROD__API_SECRET": "s3cr3t",
	}

	opts := DefaultScopeOptions()
	opts.Scopes = []string{"prod"}
	opts.TargetScope = "prod"

	resolved, err := Scope(env, opts)
	if err != nil {
		t.Fatalf("Scope error: %v", err)
	}

	tmp := filepath.Join(t.TempDir(), ".env.prod")
	if err := Write(tmp, resolved); err != nil {
		t.Fatalf("Write error: %v", err)
	}

	parsed, err := Parse(tmp)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if parsed["LOG_LEVEL"] != "warn" {
		t.Errorf("expected warn, got %q", parsed["LOG_LEVEL"])
	}
	if parsed["API_SECRET"] != "s3cr3t" {
		t.Errorf("expected s3cr3t, got %q", parsed["API_SECRET"])
	}
}

func TestScope_ThenDiff_DetectsOverrides(t *testing.T) {
	base := map[string]string{
		"DB_URL":    "base-db",
		"LOG_LEVEL": "info",
	}

	env := map[string]string{
		"DB_URL":         "base-db",
		"LOG_LEVEL":      "info",
		"PROD__DB_URL":   "prod-db",
		"PROD__LOG_LEVEL": "error",
	}

	opts := DefaultScopeOptions()
	opts.Scopes = []string{"prod"}
	opts.TargetScope = "prod"

	resolved, err := Scope(env, opts)
	if err != nil {
		t.Fatalf("Scope error: %v", err)
	}

	diffs := Diff(base, resolved)
	if !HasChanges(diffs) {
		t.Error("expected changes between base and prod-scoped env")
	}

	updated := 0
	for _, d := range diffs {
		if d.Status == StatusUpdated {
			updated++
		}
	}
	if updated != 2 {
		t.Errorf("expected 2 updated keys, got %d", updated)
	}
}

func TestScope_FilePermissionsAfterWrite(t *testing.T) {
	env := map[string]string{"DEV__KEY": "val", "KEY": "global"}
	opts := DefaultScopeOptions()
	opts.Scopes = []string{"dev"}
	opts.TargetScope = "dev"

	resolved, _ := Scope(env, opts)
	tmp := filepath.Join(t.TempDir(), ".env")
	_ = Write(tmp, resolved)

	info, err := os.Stat(tmp)
	if err != nil {
		t.Fatalf("stat error: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected 0600 permissions, got %v", info.Mode().Perm())
	}
}
