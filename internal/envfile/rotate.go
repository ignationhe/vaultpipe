package envfile

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// RotateOptions controls rotation behaviour.
type RotateOptions struct {
	// BackupDir is where rotated files are stored. Defaults to same dir as source.
	BackupDir string
	// MaxBackups is the maximum number of rotated files to keep (0 = unlimited).
	MaxBackups int
	// Timestamp overrides the rotation timestamp (useful for testing).
	Timestamp time.Time
}

// DefaultRotateOptions returns sensible defaults.
func DefaultRotateOptions() RotateOptions {
	return RotateOptions{
		MaxBackups: 5,
	}
}

// Rotate renames the current env file to a timestamped copy and writes new
// content into the original path. It returns the path of the rotated file, or
// an empty string if the source did not exist.
func Rotate(src string, newContent map[string]string, opts RotateOptions) (string, error) {
	if opts.Timestamp.IsZero() {
		opts.Timestamp = time.Now()
	}

	backupDir := opts.BackupDir
	if backupDir == "" {
		backupDir = filepath.Dir(src)
	}

	rotatedPath := ""

	// Only rotate if the source file already exists.
	if _, err := os.Stat(src); err == nil {
		ts := opts.Timestamp.UTC().Format("20060102T150405Z")
		base := filepath.Base(src)
		rotatedPath = filepath.Join(backupDir, fmt.Sprintf("%s.%s", base, ts))

		data, err := os.ReadFile(src)
		if err != nil {
			return "", fmt.Errorf("rotate: read source: %w", err)
		}
		if err := os.WriteFile(rotatedPath, data, 0o600); err != nil {
			return "", fmt.Errorf("rotate: write rotated file: %w", err)
		}

		if opts.MaxBackups > 0 {
			if err := pruneBackups(backupDir, base, opts.MaxBackups); err != nil {
				return rotatedPath, fmt.Errorf("rotate: prune: %w", err)
			}
		}
	}

	// Write fresh content to the original path.
	if err := Write(src, newContent); err != nil {
		return rotatedPath, fmt.Errorf("rotate: write new content: %w", err)
	}

	return rotatedPath, nil
}
