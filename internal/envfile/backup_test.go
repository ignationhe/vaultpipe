package envfile

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestBackup_CreatesBackupFile(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")
	_ = os.WriteFile(envPath, []byte("KEY=value\n"), 0600)

	backupPath, err := Backup(envPath, DefaultBackupOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if backupPath == "" {
		t.Fatal("expected a backup path, got empty string")
	}

	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		t.Errorf("backup file not found at %s", backupPath)
	}
}

func TestBackup_ContentsMatch(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")
	original := []byte("KEY=value\nFOO=bar\n")
	_ = os.WriteFile(envPath, original, 0600)

	backupPath, err := Backup(envPath, DefaultBackupOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(backupPath)
	if err != nil {
		t.Fatalf("read backup: %v", err)
	}
	if string(data) != string(original) {
		t.Errorf("backup contents mismatch: got %q, want %q", data, original)
	}
}

func TestBackup_ReturnsEmptyWhenSourceMissing(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env.nonexistent")

	backupPath, err := Backup(envPath, DefaultBackupOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if backupPath != "" {
		t.Errorf("expected empty path for missing source, got %q", backupPath)
	}
}

func TestBackup_PrunesOldBackups(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")
	_ = os.WriteFile(envPath, []byte("K=v\n"), 0600)

	opts := BackupOptions{MaxBackups: 3}

	for i := 0; i < 5; i++ {
		time.Sleep(2 * time.Millisecond) // ensure distinct timestamps
		_, err := Backup(envPath, opts)
		if err != nil {
			t.Fatalf("backup %d failed: %v", i, err)
		}
	}

	matches, err := filepath.Glob(filepath.Join(dir, ".env.*.bak"))
	if err != nil {
		t.Fatalf("glob error: %v", err)
	}
	if len(matches) != 3 {
		t.Errorf("expected 3 backups after pruning, got %d", len(matches))
	}
}
