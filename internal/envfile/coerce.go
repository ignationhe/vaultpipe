package envfile

import (
	"fmt"
	"strings"
)

// CoerceOptions controls how values are coerced to canonical string forms.
type CoerceOptions struct {
	// BoolTrue is the canonical string for truthy values (default: "true").
	BoolTrue string
	// BoolFalse is the canonical string for falsy values (default: "false").
	BoolFalse string
	// TruthyValues are input strings considered truthy (case-insensitive).
	TruthyValues []string
	// FalsyValues are input strings considered falsy (case-insensitive).
	FalsyValues []string
	// TrimSpace removes surrounding whitespace before coercion.
	TrimSpace bool
	// ErrorOnUnknown returns an error when a value cannot be coerced.
	ErrorOnUnknown bool
	// Keys restricts coercion to these keys; empty means all keys.
	Keys []string
}

// DefaultCoerceOptions returns sensible defaults for CoerceOptions.
func DefaultCoerceOptions() CoerceOptions {
	return CoerceOptions{
		BoolTrue:       "true",
		BoolFalse:      "false",
		TruthyValues:   []string{"1", "yes", "on", "true", "enabled"},
		FalsyValues:    []string{"0", "no", "off", "false", "disabled"},
		TrimSpace:      true,
		ErrorOnUnknown: false,
	}
}

// Coerce normalises boolean-like values in env to canonical true/false strings.
// Non-boolean values are left untouched unless ErrorOnUnknown is set.
func Coerce(env map[string]string, opts CoerceOptions) (map[string]string, error) {
	keySet := toSet(opts.Keys)

	truthy := make(map[string]struct{}, len(opts.TruthyValues))
	for _, v := range opts.TruthyValues {
		truthy[strings.ToLower(v)] = struct{}{}
	}
	falsy := make(map[string]struct{}, len(opts.FalsyValues))
	for _, v := range opts.FalsyValues {
		falsy[strings.ToLower(v)] = struct{}{}
	}

	out := make(map[string]string, len(env))
	for k, v := range env {
		if len(keySet) > 0 {
			if _, ok := keySet[k]; !ok {
				out[k] = v
				continue
			}
		}

		raw := v
		if opts.TrimSpace {
			raw = strings.TrimSpace(raw)
		}

		lower := strings.ToLower(raw)
		if _, ok := truthy[lower]; ok {
			out[k] = opts.BoolTrue
			continue
		}
		if _, ok := falsy[lower]; ok {
			out[k] = opts.BoolFalse
			continue
		}

		if opts.ErrorOnUnknown {
			return nil, fmt.Errorf("coerce: key %q has non-boolean value %q", k, raw)
		}
		out[k] = v
	}
	return out, nil
}
