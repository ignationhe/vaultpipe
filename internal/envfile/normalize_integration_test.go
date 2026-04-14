package envfile

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNormalize_ThenWrite_RoundTrip(t *testing.T) {
	env := map[string]string{
		"db-host":   "  localhost  ",
		"api--key":  "secret",
		"my__token": "abc123",
	}

	opts := DefaultNormalizeOptions()
	opts.CollapseUnderscores = true

	normalized, err := Normalize(env, opts)
	if err != nil {
		t.Fatalf("Normalize error: %v", err)
	}

	tmp := filepath.Join(t.TempDir(), ".env")
	if err := Write(tmp, normalized); err != nil {
		t.Fatalf("Write error: %v", err)
	}

	parsed, err := Parse(tmp)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	expected := map[string]string{
		"DB_HOST":  "localhost",
		"API__KEY": "secret",
		"MY_TOKEN": "abc123",
	}
	for k, want := range expected {
		if got := parsed[k]; got != want {
			t.Errorf("key %s: expected %q, got %q", k, want, got)
		}
	}
}

func TestNormalize_ThenDiff_DetectsKeyRename(t *testing.T) {
	old := map[string]string{
		"my-key": "value",
	}
	newEnv, err := Normalize(old, DefaultNormalizeOptions())
	if err != nil {
		t.Fatalf("Normalize error: %v", err)
	}

	diffs := Diff(old, newEnv)
	hasChanges := HasChanges(diffs)
	if !hasChanges {
		t.Error("expected diff to detect changes after normalization")
	}
}

func TestNormalize_FilePermissionsPreserved(t *testing.T) {
	env := map[string]string{"key-name": "val"}
	normalized, _ := Normalize(env, DefaultNormalizeOptions())

	tmp := filepath.Join(t.TempDir(), ".env")
	if err := Write(tmp, normalized); err != nil {
		t.Fatalf("Write error: %v", err)
	}

	info, err := os.Stat(tmp)
	if err != nil {
		t.Fatalf("Stat error: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected permissions 0600, got %v", info.Mode().Perm())
	}
}
