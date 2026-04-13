package envfile

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// DefaultCloneOptions returns sensible defaults for Clone.
func DefaultCloneOptions() CloneOptions {
	return CloneOptions{
		Overwrite:    false,
		TransformKey: nil,
		FilterKeys:   nil,
	}
}

// CloneOptions controls the behaviour of Clone.
type CloneOptions struct {
	// Overwrite allows the destination file to be replaced if it already exists.
	Overwrite bool
	// TransformKey, when set, is called for every key before writing.
	// Returning an empty string drops the key from the clone.
	TransformKey func(key string) string
	// FilterKeys limits which keys are copied. nil means copy all.
	FilterKeys []string
}

// Clone reads secrets from src, optionally filters / transforms them, and
// writes the result to dst. It returns the number of keys written.
func Clone(src, dst string, opts CloneOptions) (int, error) {
	if dst == "" {
		return 0, fmt.Errorf("clone: destination path must not be empty")
	}

	if !opts.Overwrite {
		if _, err := os.Stat(dst); err == nil {
			return 0, fmt.Errorf("clone: destination %q already exists (use Overwrite to replace)", dst)
		}
	}

	env, err := Parse(src)
	if err != nil {
		return 0, fmt.Errorf("clone: parse source: %w", err)
	}

	allowed := toSet(opts.FilterKeys)

	out := make(map[string]string, len(env))
	for k, v := range env {
		if len(allowed) > 0 {
			if _, ok := allowed[k]; !ok {
				continue
			}
		}
		newKey := k
		if opts.TransformKey != nil {
			newKey = opts.TransformKey(k)
			if newKey == "" {
				continue
			}
		}
		out[newKey] = v
	}

	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return 0, fmt.Errorf("clone: mkdir: %w", err)
	}

	if err := Write(dst, out); err != nil {
		return 0, fmt.Errorf("clone: write destination: %w", err)
	}

	return len(out), nil
}

// CloneWithPrefix is a convenience wrapper that prefixes every key in dst.
func CloneWithPrefix(src, dst, prefix string, opts CloneOptions) (int, error) {
	opts.TransformKey = func(k string) string {
		return strings.ToUpper(prefix + k)
	}
	return Clone(src, dst, opts)
}
