package envfile

import (
	"strings"
)

// NormalizeOptions controls how keys and values are normalized.
type NormalizeOptions struct {
	// UppercaseKeys converts all keys to uppercase.
	UppercaseKeys bool
	// TrimValues removes leading/trailing whitespace from values.
	TrimValues bool
	// ReplaceHyphens replaces hyphens in keys with underscores.
	ReplaceHyphens bool
	// CollapseUnderscores replaces runs of underscores with a single one.
	CollapseUnderscores bool
}

// DefaultNormalizeOptions returns sensible defaults.
func DefaultNormalizeOptions() NormalizeOptions {
	return NormalizeOptions{
		UppercaseKeys:       true,
		TrimValues:          true,
		ReplaceHyphens:      true,
		CollapseUnderscores: false,
	}
}

// Normalize applies key/value normalization rules to the given env map.
// It returns a new map with normalized keys and values.
func Normalize(env map[string]string, opts NormalizeOptions) (map[string]string, error) {
	result := make(map[string]string, len(env))

	for k, v := range env {
		nk := normalizeKey(k, opts)
		nv := normalizeValue(v, opts)
		result[nk] = nv
	}

	return result, nil
}

func normalizeKey(k string, opts NormalizeOptions) string {
	if opts.ReplaceHyphens {
		k = strings.ReplaceAll(k, "-", "_")
	}
	if opts.CollapseUnderscores {
		for strings.Contains(k, "__") {
			k = strings.ReplaceAll(k, "__", "_")
		}
	}
	if opts.UppercaseKeys {
		k = strings.ToUpper(k)
	}
	return k
}

func normalizeValue(v string, opts NormalizeOptions) string {
	if opts.TrimValues {
		v = strings.TrimSpace(v)
	}
	return v
}
