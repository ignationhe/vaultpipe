package envfile

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSaveSnapshot_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	env := map[string]string{"FOO": "bar", "BAZ": "qux"}

	path, err := SaveSnapshot(env, ".env", dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("snapshot file not found: %v", err)
	}
}

func TestSaveSnapshot_FileNameContainsSource(t *testing.T) {
	dir := t.TempDir()
	env := map[string]string{"KEY": "val"}

	path, err := SaveSnapshot(env, ".env", dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	base := filepath.Base(path)
	if !strings.HasPrefix(base, ".env_") {
		t.Errorf("expected filename to start with '.env_', got %q", base)
	}
}

func TestSaveSnapshot_ContentsAreValid(t *testing.T) {
	dir := t.TempDir()
	env := map[string]string{"API_KEY": "secret", "PORT": "8080"}

	path, err := SaveSnapshot(env, "prod.env", dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(path)
	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if snap.Values["API_KEY"] != "secret" {
		t.Errorf("expected API_KEY=secret, got %q", snap.Values["API_KEY"])
	}
	if snap.Source != "prod.env" {
		t.Errorf("expected source prod.env, got %q", snap.Source)
	}
}

func TestLoadSnapshot_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	env := map[string]string{"DB_URL": "postgres://localhost/db"}

	path, err := SaveSnapshot(env, ".env", dir)
	if err != nil {
		t.Fatalf("save error: %v", err)
	}

	snap, err := LoadSnapshot(path)
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	if snap.Values["DB_URL"] != "postgres://localhost/db" {
		t.Errorf("unexpected value: %q", snap.Values["DB_URL"])
	}
}

func TestLoadSnapshot_FileNotFound(t *testing.T) {
	_, err := LoadSnapshot("/nonexistent/path/snap.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestSaveSnapshot_DefaultDirUsedWhenEmpty(t *testing.T) {
	// Override working directory so we don't pollute the repo.
	tmp := t.TempDir()
	old, _ := os.Getwd()
	_ = os.Chdir(tmp)
	defer os.Chdir(old)

	env := map[string]string{"X": "1"}
	path, err := SaveSnapshot(env, ".env", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(path, DefaultSnapshotDir) {
		t.Errorf("expected path to contain default dir, got %q", path)
	}
}
