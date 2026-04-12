package envfile

import "fmt"

// PromoteOptions configures how secrets are promoted between environments.
type PromoteOptions struct {
	// FromEnv is the source environment label (e.g. "staging").
	FromEnv string
	// ToEnv is the destination environment label (e.g. "production").
	ToEnv string
	// Keys limits promotion to specific keys; empty means promote all.
	Keys []string
	// Overwrite controls whether existing keys in the destination are replaced.
	Overwrite bool
}

// DefaultPromoteOptions returns sensible defaults for PromoteOptions.
func DefaultPromoteOptions() PromoteOptions {
	return PromoteOptions{
		Overwrite: false,
	}
}

// PromoteResult holds the outcome of a promotion operation.
type PromoteResult struct {
	Promoted  []string
	Skipped   []string
	Overwritten []string
}

// Promote copies selected keys from src into dst according to opts.
// It returns a PromoteResult describing what changed.
func Promote(src, dst map[string]string, opts PromoteOptions) (map[string]string, PromoteResult, error) {
	if src == nil {
		return nil, PromoteResult{}, fmt.Errorf("promote: source map must not be nil")
	}
	if dst == nil {
		dst = make(map[string]string)
	}

	keys := opts.Keys
	if len(keys) == 0 {
		for k := range src {
			keys = append(keys, k)
		}
	}

	out := copyMap(dst)
	var result PromoteResult

	for _, k := range keys {
		v, ok := src[k]
		if !ok {
			continue
		}
		if existing, exists := out[k]; exists {
			if !opts.Overwrite {
				result.Skipped = append(result.Skipped, k)
				continue
			}
			if existing != v {
				result.Overwritten = append(result.Overwritten, k)
			}
		} else {
			result.Promoted = append(result.Promoted, k)
		}
		out[k] = v
	}

	return out, result, nil
}
