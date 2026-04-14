package envfile

import (
	"strings"
)

// TrimOptions configures the Trim operation.
type TrimOptions struct {
	// Keys limits trimming to specific keys; empty means all keys.
	Keys []string
	// TrimLeft removes leading whitespace/characters from values.
	TrimLeft bool
	// TrimRight removes trailing whitespace/characters from values.
	TrimRight bool
	// Cutset is the set of characters to trim; defaults to whitespace.
	Cutset string
}

// DefaultTrimOptions returns sensible defaults for Trim.
func DefaultTrimOptions() TrimOptions {
	return TrimOptions{
		TrimLeft:  true,
		TrimRight: true,
		Cutset:    "",
	}
}

// Trim removes leading and/or trailing characters from env values.
// When Cutset is empty, Unicode whitespace is trimmed.
// When Keys is non-empty, only those keys are affected.
func Trim(env map[string]string, opts TrimOptions) (map[string]string, error) {
	keySet := toSet(opts.Keys)
	out := make(map[string]string, len(env))

	for k, v := range env {
		if len(keySet) > 0 && !keySet[k] {
			out[k] = v
			continue
		}

		trimmed := v
		if opts.Cutset == "" {
			switch {
			case opts.TrimLeft && opts.TrimRight:
				trimmed = strings.TrimSpace(v)
			case opts.TrimLeft:
				trimmed = strings.TrimLeftFunc(v, func(r rune) bool { return r == ' ' || r == '\t' })
			case opts.TrimRight:
				trimmed = strings.TrimRightFunc(v, func(r rune) bool { return r == ' ' || r == '\t' })
			}
		} else {
			switch {
			case opts.TrimLeft && opts.TrimRight:
				trimmed = strings.Trim(v, opts.Cutset)
			case opts.TrimLeft:
				trimmed = strings.TrimLeft(v, opts.Cutset)
			case opts.TrimRight:
				trimmed = strings.TrimRight(v, opts.Cutset)
			}
		}
		out[k] = trimmed
	}
	return out, nil
}
