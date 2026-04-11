package envfile

import (
	"fmt"
	"strings"
)

// LintSeverity represents the severity level of a lint issue.
type LintSeverity string

const (
	LintWarning LintSeverity = "warning"
	LintError   LintSeverity = "error"
)

// LintIssue describes a single lint finding for a key/value pair.
type LintIssue struct {
	Key      string
	Message  string
	Severity LintSeverity
}

func (i LintIssue) String() string {
	return fmt.Sprintf("[%s] %s: %s", i.Severity, i.Key, i.Message)
}

// LintOptions controls which lint rules are applied.
type LintOptions struct {
	// WarnOnLowercase emits a warning when a key contains lowercase letters.
	WarnOnLowercase bool
	// WarnOnEmptyValue emits a warning when a key has an empty value.
	WarnOnEmptyValue bool
	// ErrorOnDuplicate is not applicable at map level (keys are unique), kept
	// for future raw-file linting.
	ErrorOnLeadingUnderscore bool
	// WarnOnLongValue emits a warning when a value exceeds MaxValueLength.
	WarnOnLongValue bool
	MaxValueLength  int
}

// DefaultLintOptions returns a sensible default configuration.
func DefaultLintOptions() LintOptions {
	return LintOptions{
		WarnOnLowercase:          true,
		WarnOnEmptyValue:         true,
		ErrorOnLeadingUnderscore: false,
		WarnOnLongValue:          true,
		MaxValueLength:           256,
	}
}

// Lint inspects the provided env map and returns a slice of LintIssues.
// An empty slice means no issues were found.
func Lint(env map[string]string, opts LintOptions) []LintIssue {
	var issues []LintIssue

	for k, v := range env {
		if opts.WarnOnLowercase && k != strings.ToUpper(k) {
			issues = append(issues, LintIssue{
				Key:      k,
				Message:  "key contains lowercase letters; consider using UPPER_SNAKE_CASE",
				Severity: LintWarning,
			})
		}

		if opts.WarnOnEmptyValue && v == "" {
			issues = append(issues, LintIssue{
				Key:      k,
				Message:  "value is empty",
				Severity: LintWarning,
			})
		}

		if opts.ErrorOnLeadingUnderscore && strings.HasPrefix(k, "_") {
			issues = append(issues, LintIssue{
				Key:      k,
				Message:  "key starts with an underscore, which is reserved for internal use",
				Severity: LintError,
			})
		}

		if opts.WarnOnLongValue && opts.MaxValueLength > 0 && len(v) > opts.MaxValueLength {
			issues = append(issues, LintIssue{
				Key:      k,
				Message:  fmt.Sprintf("value length %d exceeds maximum %d", len(v), opts.MaxValueLength),
				Severity: LintWarning,
			})
		}
	}

	return issues
}
