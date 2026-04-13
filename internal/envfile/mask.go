package envfile

import (
	"regexp"
	"strings"
)

// MaskOptions controls how values are masked in output.
type MaskOptions struct {
	// Patterns is a list of key patterns (substring match) to mask.
	Patterns []string
	// Placeholder replaces the masked value.
	Placeholder string
	// MaskPartial shows the first N characters before the placeholder.
	MaskPartial int
	// CaseSensitive controls whether pattern matching is case-sensitive.
	CaseSensitive bool
}

// DefaultMaskOptions returns sensible defaults.
func DefaultMaskOptions() MaskOptions {
	return MaskOptions{
		Patterns:    []string{"SECRET", "PASSWORD", "TOKEN", "KEY", "PRIVATE", "CREDENTIAL"},
		Placeholder: "****",
		MaskPartial: 0,
	}
}

// Mask returns a copy of env with sensitive values replaced by the placeholder.
func Mask(env map[string]string, opts MaskOptions) map[string]string {
	if opts.Placeholder == "" {
		opts.Placeholder = "****"
	}

	result := make(map[string]string, len(env))
	for k, v := range env {
		if shouldMask(k, opts) {
			result[k] = maskedValue(v, opts)
		} else {
			result[k] = v
		}
	}
	return result
}

func shouldMask(key string, opts MaskOptions) bool {
	for _, pattern := range opts.Patterns {
		a, b := key, pattern
		if !opts.CaseSensitive {
			a = strings.ToUpper(a)
			b = strings.ToUpper(b)
		}
		matched, _ := regexp.MatchString(regexp.QuoteMeta(b), a)
		if matched {
			return true
		}
	}
	return false
}

func maskedValue(v string, opts MaskOptions) string {
	if opts.MaskPartial > 0 && len(v) > opts.MaskPartial {
		return v[:opts.MaskPartial] + opts.Placeholder
	}
	return opts.Placeholder
}
