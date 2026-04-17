package envfile

import "sort"

// GroupOptions controls how keys are grouped.
type GroupOptions struct {
	// Separator splits key into (group, remainder). Default: "_"
	Separator string
	// KeepPrefix retains the group prefix in the key name.
	KeepPrefix bool
}

// DefaultGroupOptions returns sensible defaults.
func DefaultGroupOptions() GroupOptions {
	return GroupOptions{
		Separator:  "_",
		KeepPrefix: false,
	}
}

// Group partitions a flat map into named buckets by the first segment of each
// key, split on opts.Separator. Keys with no separator go into the "_" bucket.
func Group(env map[string]string, opts GroupOptions) map[string]map[string]string {
	if opts.Separator == "" {
		opts.Separator = "_"
	}

	result := map[string]map[string]string{}

	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := env[k]
		groupName, remainder := splitOnFirst(k, opts.Separator)
		if remainder == "" {
			// No separator found — place in catch-all bucket.
			groupName = "_"
			remainder = k
		} else if opts.KeepPrefix {
			remainder = k
		}
		if result[groupName] == nil {
			result[groupName] = map[string]string{}
		}
		result[groupName][remainder] = v
	}
	return result
}

// splitOnFirst splits s on the first occurrence of sep.
// Returns (s, "") when sep is not found.
func splitOnFirst(s, sep string) (string, string) {
	for i := 0; i <= len(s)-len(sep); i++ {
		if s[i:i+len(sep)] == sep {
			return s[:i], s[i+len(sep):]
		}
	}
	return s, ""
}
