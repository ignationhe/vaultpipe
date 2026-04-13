package envfile

import (
	"fmt"
	"os"
	"strings"
)

// InjectOptions controls how secrets are injected into the process environment.
type InjectOptions struct {
	// Prefix filters keys to only those starting with the given prefix.
	Prefix string
	// StripPrefix removes the prefix from keys before injecting.
	StripPrefix bool
	// Overwrite allows overwriting existing OS environment variables.
	Overwrite bool
	// DryRun prints what would be injected without actually calling os.Setenv.
	DryRun bool
}

// DefaultInjectOptions returns sensible defaults for InjectOptions.
func DefaultInjectOptions() InjectOptions {
	return InjectOptions{
		Overwrite: false,
		DryRun:    false,
	}
}

// InjectedKey records the result of a single key injection.
type InjectedKey struct {
	Key      string
	Skipped  bool // true when Overwrite is false and key already exists
	DryRun   bool
}

// Inject takes a map of key/value pairs and sets them in the current process
// environment according to opts. It returns a slice of InjectedKey records
// describing what happened for each key.
func Inject(env map[string]string, opts InjectOptions) ([]InjectedKey, error) {
	results := make([]InjectedKey, 0, len(env))

	for k, v := range env {
		key := k

		if opts.Prefix != "" && !strings.HasPrefix(key, opts.Prefix) {
			continue
		}

		if opts.StripPrefix && opts.Prefix != "" {
			key = strings.TrimPrefix(key, opts.Prefix)
		}

		if key == "" {
			return nil, fmt.Errorf("inject: stripping prefix %q from %q produces empty key", opts.Prefix, k)
		}

		record := InjectedKey{Key: key, DryRun: opts.DryRun}

		if !opts.Overwrite && os.Getenv(key) != "" {
			record.Skipped = true
			results = append(results, record)
			continue
		}

		if !opts.DryRun {
			if err := os.Setenv(key, v); err != nil {
				return nil, fmt.Errorf("inject: setenv %q: %w", key, err)
			}
		}

		results = append(results, record)
	}

	return results, nil
}
