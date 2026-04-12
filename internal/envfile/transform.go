package envfile

import (
	"fmt"
	"strings"
)

// TransformFunc is a function that transforms a single env value.
type TransformFunc func(key, value string) (string, error)

// TransformOptions controls how Transform applies mutations to env maps.
type TransformOptions struct {
	// Rules maps key patterns (exact or wildcard suffix "*") to TransformFuncs.
	Rules map[string]TransformFunc
	// SkipErrors continues processing even if a transform returns an error.
	SkipErrors bool
}

// DefaultTransformOptions returns an empty, non-strict TransformOptions.
func DefaultTransformOptions() TransformOptions {
	return TransformOptions{
		Rules:      make(map[string]TransformFunc),
		SkipErrors: false,
	}
}

// Transform applies registered transform rules to each key in src.
// Exact key matches take priority over wildcard ("*") matches.
// Returns a new map with transformed values; src is not modified.
func Transform(src map[string]string, opts TransformOptions) (map[string]string, error) {
	if opts.Rules == nil {
		opts.Rules = make(map[string]TransformFunc)
	}

	out := make(map[string]string, len(src))
	for k, v := range src {
		out[k] = v
	}

	for k, v := range out {
		fn, ok := opts.Rules[k]
		if !ok {
			fn, ok = opts.Rules["*"]
		}
		if !ok {
			continue
		}
		result, err := fn(k, v)
		if err != nil {
			if opts.SkipErrors {
				continue
			}
			return nil, fmt.Errorf("transform error on key %q: %w", k, err)
		}
		out[k] = result
	}
	return out, nil
}

// UppercaseValues returns a TransformFunc that uppercases the value.
func UppercaseValues() TransformFunc {
	return func(_, v string) (string, error) {
		return strings.ToUpper(v), nil
	}
}

// LowercaseValues returns a TransformFunc that lowercases the value.
func LowercaseValues() TransformFunc {
	return func(_, v string) (string, error) {
		return strings.ToLower(v), nil
	}
}

// TrimSpaceValues returns a TransformFunc that trims whitespace from values.
func TrimSpaceValues() TransformFunc {
	return func(_, v string) (string, error) {
		return strings.TrimSpace(v), nil
	}
}
