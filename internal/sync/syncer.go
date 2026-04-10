package sync

import (
	"errors"
	"fmt"
	"os"

	"github.com/yourusername/vaultpipe/internal/envfile"
)

// VaultClient is the interface the syncer relies on to fetch secrets.
type VaultClient interface {
	GetSecrets(path string) (map[string]string, error)
}

// Options configures a Sync run.
type Options struct {
	VaultPath  string
	EnvFile    string
	DryRun     bool
	Filter     envfile.FilterOptions
	StripPrefix string
}

// Result summarises what happened during a sync.
type Result struct {
	Diff    []envfile.Change
	Written bool
}

// Syncer orchestrates fetching secrets and writing them to an env file.
type Syncer struct {
	client VaultClient
}

// New creates a Syncer backed by client.
func New(client VaultClient) *Syncer {
	return &Syncer{client: client}
}

// Sync fetches secrets from Vault and merges them into the env file.
func (s *Syncer) Sync(opts Options) (Result, error) {
	secrets, err := s.client.GetSecrets(opts.VaultPath)
	if err != nil {
		return Result{}, fmt.Errorf("vault: %w", err)
	}

	// Apply key filtering and optional prefix stripping.
	secrets = envfile.Filter(secrets, opts.Filter)
	if opts.StripPrefix != "" {
		secrets = envfile.StripPrefix(secrets, opts.StripPrefix)
	}

	existing := map[string]string{}
	if parsed, err := envfile.Parse(opts.EnvFile); err == nil {
		existing = parsed
	} else if !isNotExist(err) {
		return Result{}, fmt.Errorf("parse env file: %w", err)
	}

	changes := envfile.Diff(existing, secrets)
	result := Result{Diff: changes}

	if opts.DryRun || !envfile.HasChanges(changes) {
		return result, nil
	}

	merged := envfile.Merge(existing, secrets)
	if err := envfile.Write(opts.EnvFile, merged); err != nil {
		return result, fmt.Errorf("write env file: %w", err)
	}
	result.Written = true
	return result, nil
}

func isNotExist(err error) bool {
	return errors.Is(err, os.ErrNotExist)
}
