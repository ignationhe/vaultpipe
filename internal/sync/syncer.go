package sync

import (
	"context"
	"fmt"
	"log"

	"github.com/user/vaultpipe/internal/envfile"
)

// VaultClient is the interface for fetching secrets from Vault.
type VaultClient interface {
	GetSecrets(ctx context.Context, path string) (map[string]string, error)
}

// Options configures Syncer behavior.
type Options struct {
	DryRun bool
	Backup bool
}

// Syncer orchestrates pulling secrets from Vault and writing them to an env file.
type Syncer struct {
	vault  VaultClient
	opts   Options
}

// New creates a new Syncer.
func New(vault VaultClient, opts Options) *Syncer {
	return &Syncer{vault: vault, opts: opts}
}

// Sync fetches secrets from the given Vault path and merges them into envPath.
// It returns the diff so callers can display a summary.
func (s *Syncer) Sync(ctx context.Context, vaultPath, envPath string) ([]envfile.DiffEntry, error) {
	secrets, err := s.vault.GetSecrets(ctx, vaultPath)
	if err != nil {
		return nil, fmt.Errorf("sync: fetch secrets: %w", err)
	}

	existing, err := envfile.Parse(envPath)
	if err != nil && !isNotExist(err) {
		return nil, fmt.Errorf("sync: parse env file: %w", err)
	}

	diffs := envfile.Diff(existing, secrets)
	if !envfile.HasChanges(diffs) {
		log.Println("vaultpipe: no changes detected")
		return diffs, nil
	}

	if s.opts.DryRun {
		log.Println("vaultpipe: dry-run mode — no files written")
		return diffs, nil
	}

	if s.opts.Backup {
		backupPath, err := envfile.Backup(envPath, envfile.DefaultBackupOptions())
		if err != nil {
			return nil, fmt.Errorf("sync: backup: %w", err)
		}
		if backupPath != "" {
			log.Printf("vaultpipe: backup written to %s", backupPath)
		}
	}

	if err := envfile.Merge(envPath, secrets); err != nil {
		return nil, fmt.Errorf("sync: write env file: %w", err)
	}

	return diffs, nil
}

func isNotExist(err error) bool {
	if err == nil {
		return false
	}
	return os.IsNotExist(err)
}
