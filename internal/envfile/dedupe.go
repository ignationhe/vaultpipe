package envfile

import "fmt"

// DedupeStrategy controls how duplicate keys are resolved.
type DedupeStrategy string

const (
	DedupeKeepFirst DedupeStrategy = "first"
	DedupeKeepLast  DedupeStrategy = "last"
	DedupeError     DedupeStrategy = "error"
)

// DedupeOptions configures the Dedupe function.
type DedupeOptions struct {
	// Strategy determines which value wins when a key appears more than once.
	Strategy DedupeStrategy
}

// DefaultDedupeOptions returns sensible defaults.
func DefaultDedupeOptions() DedupeOptions {
	return DedupeOptions{
		Strategy: DedupeKeepLast,
	}
}

// DedupeResult holds the cleaned map and a report of which keys were duplicated.
type DedupeResult struct {
	Env        map[string]string
	Duplicates []string // keys that appeared more than once
}

// Dedupe removes duplicate keys from the provided ordered pairs, applying the
// chosen strategy. The input is represented as a slice of [2]string pairs so
// that insertion order and repetition are both visible.
func Dedupe(pairs [][2]string, opts DedupeOptions) (DedupeResult, error) {
	if opts.Strategy == "" {
		opts = DefaultDedupeOptions()
	}

	seen := make(map[string]int) // key -> count
	result := make(map[string]string)
	var dupes []string

	for _, pair := range pairs {
		k, v := pair[0], pair[1]
		seen[k]++

		switch opts.Strategy {
		case DedupeError:
			if seen[k] > 1 {
				return DedupeResult{}, fmt.Errorf("duplicate key: %q", k)
			}
			result[k] = v
		case DedupeKeepFirst:
			if seen[k] == 1 {
				result[k] = v
			}
		default: // DedupeKeepLast
			result[k] = v
		}
	}

	for k, count := range seen {
		if count > 1 {
			dupes = append(dupes, k)
		}
	}

	return DedupeResult{Env: result, Duplicates: dupes}, nil
}
