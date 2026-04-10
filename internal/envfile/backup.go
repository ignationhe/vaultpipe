package envfile

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// BackupOptions configures backup behavior.
type BackupOptions struct {
	MaxBackups int
}

// DefaultBackupOptions returns sensible defaults.
func DefaultBackupOptions() BackupOptions {
	return BackupOptions{
		MaxBackups: 5,
	}
}

// Backup creates a timestamped copy of the given .env file.
// It returns the path of the backup file or an error.
func Backup(envPath string, opts BackupOptions) (string, error) {
	data, err := os.ReadFile(envPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil // nothing to back up
		}
		return "", fmt.Errorf("backup: read source: %w", err)
	}

	dir := filepath.Dir(envPath)
	base := filepath.Base(envPath)
	timestamp := time.Now().UTC().Format("20060102T150405Z")
	backupName := fmt.Sprintf("%s.%s.bak", base, timestamp)
	backupPath := filepath.Join(dir, backupName)

	if err := os.WriteFile(backupPath, data, 0600); err != nil {
		return "", fmt.Errorf("backup: write backup: %w", err)
	}

	if opts.MaxBackups > 0 {
		if err := pruneBackups(dir, base, opts.MaxBackups); err != nil {
			return backupPath, fmt.Errorf("backup: prune: %w", err)
		}
	}

	return backupPath, nil
}

// pruneBackups removes the oldest backups if the count exceeds maxBackups.
func pruneBackups(dir, base string, maxBackups int) error {
	pattern := filepath.Join(dir, base+".*.bak")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}

	if len(matches) <= maxBackups {
		return nil
	}

	// matches from Glob are lexicographically sorted; oldest timestamps come first.
	toRemove := matches[:len(matches)-maxBackups]
	for _, f := range toRemove {
		if err := os.Remove(f); err != nil && !os.IsNotExist(err) {
			return err
		}
	}
	return nil
}
