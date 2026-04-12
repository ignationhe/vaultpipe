package envfile

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Snapshot represents a point-in-time capture of an env file's key-value pairs.
type Snapshot struct {
	Timestamp time.Time         `json:"timestamp"`
	Source    string            `json:"source"`
	Values    map[string]string `json:"values"`
}

// DefaultSnapshotDir is the directory used when none is specified.
const DefaultSnapshotDir = ".vaultpipe/snapshots"

// SaveSnapshot writes the current state of env to a JSON snapshot file.
func SaveSnapshot(env map[string]string, source, dir string) (string, error) {
	if dir == "" {
		dir = DefaultSnapshotDir
	}
	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", fmt.Errorf("snapshot: create dir: %w", err)
	}

	snap := Snapshot{
		Timestamp: time.Now().UTC(),
		Source:    source,
		Values:    env,
	}

	filename := fmt.Sprintf("%s_%s.json",
		filepath.Base(source),
		snap.Timestamp.Format("20060102T150405Z"),
	)
	path := filepath.Join(dir, filename)

	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return "", fmt.Errorf("snapshot: marshal: %w", err)
	}
	if err := os.WriteFile(path, data, 0600); err != nil {
		return "", fmt.Errorf("snapshot: write: %w", err)
	}
	return path, nil
}

// LoadSnapshot reads a snapshot from a JSON file.
func LoadSnapshot(path string) (*Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("snapshot: read: %w", err)
	}
	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return nil, fmt.Errorf("snapshot: parse: %w", err)
	}
	return &snap, nil
}
