package envfile

import (
	"fmt"
	"strings"
)

// SplitOptions controls how a flat env map is split into multiple maps.
type SplitOptions struct {
	// Delimiter separates the group prefix from the key name (default "__").
	Delimiter string
	// KeepPrefix retains the prefix in the resulting key when true.
	KeepPrefix bool
	// Ungrouped is the label used for keys that have no matching prefix.
	Ungrouped string
}

// DefaultSplitOptions returns sensible defaults for Split.
func DefaultSplitOptions() SplitOptions {
	return SplitOptions{
		Delimiter:  "__",
		KeepPrefix: false,
		Ungrouped:  "default",
	}
}

// Split partitions a flat env map into named groups based on key prefixes.
// Keys like "APP__PORT" and "APP__HOST" are placed in group "APP" with keys
// "PORT" and "HOST" respectively (unless KeepPrefix is true).
// Keys without a recognised prefix land in the Ungrouped bucket.
//
// If opts.Delimiter is empty, DefaultSplitOptions().Delimiter is used.
func Split(env map[string]string, opts SplitOptions) (map[string]map[string]string, error) {
	if env == nil {
		return nil, fmt.Errorf("split: env map must not be nil")
	}
	if opts.Delimiter == "" {
		opts.Delimiter = DefaultSplitOptions().Delimiter
	}
	if opts.Ungrouped == "" {
		opts.Ungrouped = DefaultSplitOptions().Ungrouped
	}

	result := make(map[string]map[string]string)

	for k, v := range env {
		idx := strings.Index(k, opts.Delimiter)
		if idx <= 0 {
			// No delimiter found or delimiter is at position 0 — ungrouped.
			if result[opts.Ungrouped] == nil {
				result[opts.Ungrouped] = make(map[string]string)
			}
			result[opts.Ungrouped][k] = v
			continue
		}

		group := k[:idx]
		var key string
		if opts.KeepPrefix {
			key = k
		} else {
			key = k[idx+len(opts.Delimiter):]
		}
		if key == "" {
			// Delimiter was trailing — treat as ungrouped.
			if result[opts.Ungrouped] == nil {
				result[opts.Ungrouped] = make(map[string]string)
			}
			result[opts.Ungrouped][k] = v
			continue
		}

		if result[group] == nil {
			result[group] = make(map[string]string)
		}
		result[group][key] = v
	}

	return result, nil
}
