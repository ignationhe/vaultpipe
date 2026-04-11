package envfile

// MergeStrategy controls how conflicts are resolved when merging two env maps.
type MergeStrategy int

const (
	// StrategyKeepExisting preserves the existing value on conflict.
	StrategyKeepExisting MergeStrategy = iota
	// StrategyOverwrite replaces the existing value with the incoming value.
	StrategyOverwrite
	// StrategyKeepBoth retains the existing key and adds the incoming key with a suffix.
	StrategyKeepBoth
)

// MergeOptions configures the behaviour of MergeWith.
type MergeOptions struct {
	Strategy MergeStrategy
	// ConflictSuffix is appended to the incoming key when StrategyKeepBoth is used.
	// Defaults to "_NEW" when empty.
	ConflictSuffix string
}

// DefaultMergeOptions returns sensible defaults.
func DefaultMergeOptions() MergeOptions {
	return MergeOptions{
		Strategy:       StrategyKeepExisting,
		ConflictSuffix: "_NEW",
	}
}

// MergeWith merges incoming into base according to opts.
// base is never mutated; a new map is returned.
func MergeWith(base, incoming map[string]string, opts MergeOptions) map[string]string {
	if opts.ConflictSuffix == "" {
		opts.ConflictSuffix = "_NEW"
	}

	result := make(map[string]string, len(base))
	for k, v := range base {
		result[k] = v
	}

	for k, v := range incoming {
		if _, exists := result[k]; !exists {
			// No conflict — always add.
			result[k] = v
			continue
		}

		switch opts.Strategy {
		case StrategyOverwrite:
			result[k] = v
		case StrategyKeepBoth:
			result[k+opts.ConflictSuffix] = v
		default: // StrategyKeepExisting
			// do nothing
		}
	}

	return result
}
