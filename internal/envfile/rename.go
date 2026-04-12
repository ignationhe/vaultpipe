package envfile

import (
	"fmt"
	"regexp"
	"strings"
)

// RenameRule describes a single key rename operation.
type RenameRule struct {
	From    string // exact key name or regex pattern
	To      string // target key name
	Pattern bool   // treat From as a regex
}

// DefaultRenameOptions returns sensible defaults.
func DefaultRenameOptions() RenameOptions {
	return RenameOptions{
		DeleteOriginal: true,
		SkipMissing:    true,
	}
}

// RenameOptions controls Rename behaviour.
type RenameOptions struct {
	Rules          []RenameRule
	DeleteOriginal bool // remove the old key after renaming
	SkipMissing    bool // silently ignore rules whose From key is absent
}

// Rename applies rename rules to env, returning a new map.
func Rename(env map[string]string, opts RenameOptions) (map[string]string, error) {
	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = v
	}

	for _, rule := range opts.Rules {
		if rule.Pattern {
			re, err := regexp.Compile(rule.From)
			if err != nil {
				return nil, fmt.Errorf("rename: invalid pattern %q: %w", rule.From, err)
			}
			for k, v := range env {
				if re.MatchString(k) {
					newKey := re.ReplaceAllString(k, rule.To)
					newKey = strings.ToUpper(newKey)
					out[newKey] = v
					if opts.DeleteOriginal && newKey != k {
						delete(out, k)
					}
				}
			}
			continue
		}

		v, ok := out[rule.From]
		if !ok {
			if opts.SkipMissing {
				continue
			}
			return nil, fmt.Errorf("rename: key %q not found", rule.From)
		}
		out[rule.To] = v
		if opts.DeleteOriginal && rule.To != rule.From {
			delete(out, rule.From)
		}
	}

	return out, nil
}
