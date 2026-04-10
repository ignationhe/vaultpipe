package sync

import (
	"fmt"

	"github.com/yourusername/vaultpipe/internal/envfile"
	"github.com/yourusername/vaultpipe/internal/vault"
)

// Result holds the outcome of a sync operation.
type Result struct {
	Path    string
	Diff    []envfile.DiffEntry
	Written bool
}

// Options configures the behaviour of a sync run.
type Options struct {
	// DryRun prevents any writes; only the diff is computed.
	DryRun bool
}

// Syncer orchestrates fetching secrets from Vault and writing them to a
// local .env file.
type Syncer struct {
	client *vault.Client
}

// New creates a Syncer backed by the supplied Vault client.
func New(client *vault.Client) *Syncer {
	return &Syncer{client: client}
}

// Sync fetches secrets at vaultPath and merges them into envPath.
// When opts.DryRun is true the file is never modified.
func (s *Syncer) Sync(vaultPath, envPath string, opts Options) (Result, error) {
	secrets, err := s.client.GetSecrets(vaultPath)
	if err != nil {
		return Result{}, fmt.Errorf("fetching secrets from vault: %w", err)
	}

	existing, err := envfile.Parse(envPath)
	if err != nil {
		// Treat a missing file as an empty map so we still write it.
		existing = map[string]string{}
	}

	diffs := envfile.Diff(existing, secrets)

	result := Result{
		Path: envPath,
		Diff: diffs,
	}

	if opts.DryRun || !envfile.HasChanges(diffs) {
		return result, nil
	}

	merged := envfile.Merge(existing, secrets)
	if err := envfile.Write(envPath, merged); err != nil {
		return result, fmt.Errorf("writing env file: %w", err)
	}

	result.Written = true
	return result, nil
}
