package envfile

import (
	"regexp"
	"strings"
)

// RedactOptions controls how secrets are redacted in output.
type RedactOptions struct {
	// Keys is a list of exact key names to redact.
	Keys []string
	// Patterns is a list of regex patterns matched against key names.
	Patterns []string
	// Placeholder is the string used to replace redacted values.
	Placeholder string
}

// DefaultRedactOptions returns sensible defaults for redaction.
func DefaultRedactOptions() RedactOptions {
	return RedactOptions{
		Keys: []string{},
		Patterns: []string{
			`(?i)password`,
			`(?i)secret`,
			`(?i)token`,
			`(?i)api_key`,
			`(?i)private`,
		},
		Placeholder: "***REDACTED***",
	}
}

// Redact returns a copy of secrets with sensitive values replaced by the
// placeholder. Keys matched by exact name or by any of the regex Patterns
// are considered sensitive.
func Redact(secrets map[string]string, opts RedactOptions) map[string]string {
	if opts.Placeholder == "" {
		opts.Placeholder = DefaultRedactOptions().Placeholder
	}

	exactSet := make(map[string]struct{}, len(opts.Keys))
	for _, k := range opts.Keys {
		exactSet[strings.ToUpper(k)] = struct{}{}
	}

	compiled := make([]*regexp.Regexp, 0, len(opts.Patterns))
	for _, p := range opts.Patterns {
		if re, err := regexp.Compile(p); err == nil {
			compiled = append(compiled, re)
		}
	}

	result := make(map[string]string, len(secrets))
	for k, v := range secrets {
		if isSensitive(k, exactSet, compiled) {
			result[k] = opts.Placeholder
		} else {
			result[k] = v
		}
	}
	return result
}

func isSensitive(key string, exact map[string]struct{}, patterns []*regexp.Regexp) bool {
	if _, ok := exact[strings.ToUpper(key)]; ok {
		return true
	}
	for _, re := range patterns {
		if re.MatchString(key) {
			return true
		}
	}
	return false
}
