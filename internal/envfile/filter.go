package envfile

import "strings"

// FilterOptions controls which keys are included or excluded.
type FilterOptions struct {
	// Prefix, if non-empty, restricts keys to those with this prefix.
	Prefix string
	// Include is an explicit allowlist of keys. Empty means allow all.
	Include []string
	// Exclude is a denylist of keys that should be dropped.
	Exclude []string
}

// Filter returns a new map containing only the entries that satisfy opts.
func Filter(secrets map[string]string, opts FilterOptions) map[string]string {
	includeSet := toSet(opts.Include)
	excludeSet := toSet(opts.Exclude)

	result := make(map[string]string, len(secrets))
	for k, v := range secrets {
		if opts.Prefix != "" && !strings.HasPrefix(k, opts.Prefix) {
			continue
		}
		if len(includeSet) > 0 {
			if _, ok := includeSet[k]; !ok {
				continue
			}
		}
		if _, ok := excludeSet[k]; ok {
			continue
		}
		result[k] = v
	}
	return result
}

// StripPrefix removes a leading prefix from every key in secrets.
// Keys that do not carry the prefix are left unchanged.
func StripPrefix(secrets map[string]string, prefix string) map[string]string {
	if prefix == "" {
		return secrets
	}
	result := make(map[string]string, len(secrets))
	for k, v := range secrets {
		result[strings.TrimPrefix(k, prefix)] = v
	}
	return result
}

// FilterAndStrip is a convenience function that applies Filter followed by
// StripPrefix using opts.Prefix. This is the common pattern where callers
// want only keys matching a prefix and then want the prefix removed from
// the resulting keys.
func FilterAndStrip(secrets map[string]string, opts FilterOptions) map[string]string {
	filtered := Filter(secrets, opts)
	return StripPrefix(filtered, opts.Prefix)
}

func toSet(keys []string) map[string]struct{} {
	s := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		s[k] = struct{}{}
	}
	return s
}
