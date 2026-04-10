package envfile

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestRotate_CreatesRotatedFileAndWritesNew(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, ".env")

	// Seed an existing file.
	if err := Write(src, map[string]string{"OLD": "value"}); err != nil {
		t.Fatalf("seed: %v", err)
	}

	newContent := map[string]string{"NEW": "fresh"}
	ts := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	opts := DefaultRotateOptions()
	opts.Timestamp = ts

	rotated, err := Rotate(src, newContent, opts)
	if err != nil {
		t.Fatalf("Rotate: %v", err)
	}

	expectedRotated := filepath.Join(dir, ".env.20240601T120000Z")
	if rotated != expectedRotated {
		t.Errorf("rotated path = %q, want %q", rotated, expectedRotated)
	}

	// Rotated file must contain old content.
	oldData, err := os.ReadFile(rotated)
	if err != nil {
		t.Fatalf("read rotated: %v", err)
	}
	if string(oldData) == "" {
		t.Error("rotated file is empty")
	}

	// Source file must contain new content.
	parsed, err := Parse(src)
	if err != nil {
		t.Fatalf("parse new: %v", err)
	}
	if parsed["NEW"] != "fresh" {
		t.Errorf("new file NEW = %q, want %q", parsed["NEW"], "fresh")
	}
	if _, ok := parsed["OLD"]; ok {
		t.Error("new file still contains OLD key")
	}
}

func TestRotate_NoRotatedFileWhenSourceMissing(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, ".env")

	opts := DefaultRotateOptions()
	rotated, err := Rotate(src, map[string]string{"K": "v"}, opts)
	if err != nil {
		t.Fatalf("Rotate: %v", err)
	}
	if rotated != "" {
		t.Errorf("expected empty rotated path, got %q", rotated)
	}
	// Source should now exist with new content.
	if _, err := os.Stat(src); err != nil {
		t.Errorf("source not created: %v", err)
	}
}

func TestRotate_PrunesOldRotations(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, ".env")

	// Create 5 pre-existing backup files.
	for i := 0; i < 5; i++ {
		name := filepath.Join(dir, ".env.2024010"+string(rune('1'+i))+"T000000Z")
		_ = os.WriteFile(name, []byte("old"), 0o600)
	}
	_ = Write(src, map[string]string{"A": "1"})

	opts := DefaultRotateOptions() // MaxBackups = 5
	opts.Timestamp = time.Date(2024, 6, 10, 0, 0, 0, 0, time.UTC)

	_, err := Rotate(src, map[string]string{"B": "2"}, opts)
	if err != nil {
		t.Fatalf("Rotate: %v", err)
	}

	entries, _ := filepath.Glob(filepath.Join(dir, ".env.*"))
	if len(entries) > 5 {
		t.Errorf("expected at most 5 backups, got %d", len(entries))
	}
}
