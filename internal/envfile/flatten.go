package envfile

import (
	"fmt"
	"strings"
)

// FlattenOptions controls how nested key structures are flattened.
type FlattenOptions struct {
	// Separator is the string used to join key segments (default: "_").
	Separator string
	// Prefix is an optional prefix prepended to every output key.
	Prefix string
	// Uppercase forces all output keys to uppercase.
	Uppercase bool
}

// DefaultFlattenOptions returns sensible defaults for Flatten.
func DefaultFlattenOptions() FlattenOptions {
	return FlattenOptions{
		Separator: "_",
		Uppercase: true,
	}
}

// Flatten takes a nested map (map keys may contain the separator to represent
// hierarchy) and returns a single-level map with compound keys joined by the
// configured separator.
//
// Example input with separator "__":
//
//	{"DB__HOST": "localhost", "DB__PORT": "5432", "APP_NAME": "vaultpipe"}
//
// Example output (separator "_", prefix "APP"):
//
//	{"APP_DB_HOST": "localhost", "APP_DB_PORT": "5432", "APP_APP_NAME": "vaultpipe"}
func Flatten(input map[string]string, opts FlattenOptions) (map[string]string, error) {
	if opts.Separator == "" {
		opts.Separator = "_"
	}

	result := make(map[string]string, len(input))

	for k, v := range input {
		if k == "" {
			return nil, fmt.Errorf("flatten: empty key is not allowed")
		}

		// Normalise segments: split on common nested separators then rejoin.
		segments := splitSegments(k)
		joined := strings.Join(segments, opts.Separator)

		if opts.Prefix != "" {
			joined = opts.Prefix + opts.Separator + joined
		}

		if opts.Uppercase {
			joined = strings.ToUpper(joined)
		}

		result[joined] = v
	}

	return result, nil
}

// splitSegments splits a key on ".", "/", and "__" into individual segments,
// filtering out any empty parts produced by consecutive separators.
func splitSegments(key string) []string {
	// Normalise compound separators to a single pipe so we can split once.
	norm := strings.ReplaceAll(key, "__", "|")
	norm = strings.ReplaceAll(norm, ".", "|")
	norm = strings.ReplaceAll(norm, "/", "|")

	parts := strings.Split(norm, "|")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if p != "" {
			out = append(out, p)
		}
	}
	if len(out) == 0 {
		return []string{key}
	}
	return out
}
