package envfile

import "strings"

// ScopeEntry represents a single key scoped to an environment.
type ScopeEntry struct {
	Key   string
	Value string
	Scope string
}

// ScopeOptions controls how scoping is applied.
type ScopeOptions struct {
	// Scopes is an ordered list of scope names (e.g. ["base", "staging", "prod"]).
	// Later scopes override earlier ones.
	Scopes []string
	// Separator is placed between scope prefix and key (default "__").
	Separator string
	// TargetScope is the scope to resolve into the final flat map.
	TargetScope string
	// StripPrefix removes the scope prefix from output keys when true.
	StripPrefix bool
}

// DefaultScopeOptions returns sensible defaults.
func DefaultScopeOptions() ScopeOptions {
	return ScopeOptions{
		Separator:   "__",
		StripPrefix: true,
	}
}

// Scope resolves a flat env map that contains scope-prefixed keys into a
// single merged map for the requested target scope. Keys without any scope
// prefix are treated as base/global values and are always included unless
// overridden by a scoped key.
//
// Example input:
//
//	{"DB_URL": "base", "PROD__DB_URL": "prod-url", "STAGING__DB_URL": "staging-url"}
//
// With TargetScope="PROD" the output is {"DB_URL": "prod-url"}.
func Scope(env map[string]string, opts ScopeOptions) (map[string]string, error) {
	if opts.Separator == "" {
		opts.Separator = "__"
	}

	target := strings.ToUpper(opts.TargetScope)
	out := make(map[string]string)

	// First pass: collect unscoped (global) keys.
	for k, v := range env {
		if !isScopedKey(k, opts.Scopes, opts.Separator) {
			out[k] = v
		}
	}

	// Second pass: apply matching scoped keys in declared scope order.
	for _, scope := range opts.Scopes {
		prefix := strings.ToUpper(scope) + opts.Separator
		for k, v := range env {
			if !strings.HasPrefix(strings.ToUpper(k), prefix) {
				continue
			}
			baseKey := k[len(prefix):]
			if strings.ToUpper(scope) == target {
				if opts.StripPrefix {
					out[baseKey] = v
				} else {
					out[k] = v
				}
			}
		}
	}

	return out, nil
}

// isScopedKey returns true if the key starts with any known scope prefix.
func isScopedKey(key string, scopes []string, sep string) bool {
	upper := strings.ToUpper(key)
	for _, s := range scopes {
		if strings.HasPrefix(upper, strings.ToUpper(s)+sep) {
			return true
		}
	}
	return false
}
