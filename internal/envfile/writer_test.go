package envfile

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWrite_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	secrets := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
	}

	if err := Write(path, secrets); err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	parsed, err := Parse(path)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	for k, v := range secrets {
		if got := parsed[k]; got != v {
			t.Errorf("key %q: got %q, want %q", k, got, v)
		}
	}
}

func TestWrite_QuotesValuesWithSpaces(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	secrets := map[string]string{
		"GREETING": "hello world",
	}

	if err := Write(path, secrets); err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	parsed, err := Parse(path)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if got := parsed["GREETING"]; got != "hello world" {
		t.Errorf("GREETING: got %q, want %q", got, "hello world")
	}
}

func TestWrite_SetsFilePermissions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	if err := Write(path, map[string]string{"KEY": "val"}); err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat() error = %v", err)
	}

	if perm := info.Mode().Perm(); perm != 0600 {
		t.Errorf("file permissions: got %o, want %o", perm, 0600)
	}
}

func TestMerge_PreservesExistingKeys(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	initial := map[string]string{"EXISTING": "value", "TO_UPDATE": "old"}
	if err := Write(path, initial); err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	if err := Merge(path, map[string]string{"TO_UPDATE": "new", "ADDED": "yes"}); err != nil {
		t.Fatalf("Merge() error = %v", err)
	}

	parsed, err := Parse(path)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	expected := map[string]string{"EXISTING": "value", "TO_UPDATE": "new", "ADDED": "yes"}
	for k, v := range expected {
		if got := parsed[k]; got != v {
			t.Errorf("key %q: got %q, want %q", k, got, v)
		}
	}
}

func TestMerge_CreatesFileIfNotExists(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	if err := Merge(path, map[string]string{"NEW_KEY": "new_val"}); err != nil {
		t.Fatalf("Merge() error = %v", err)
	}

	parsed, err := Parse(path)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if got := parsed["NEW_KEY"]; got != "new_val" {
		t.Errorf("NEW_KEY: got %q, want %q", got, "new_val")
	}
}
