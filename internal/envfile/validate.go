package envfile

import (
	"fmt"
	"regexp"
	"strings"
)

// ValidationError represents a single validation issue found in an env map.
type ValidationError struct {
	Key     string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("key %q: %s", e.Key, e.Message)
}

// ValidationResult holds all errors found during validation.
type ValidationResult struct {
	Errors []ValidationError
}

func (r *ValidationResult) HasErrors() bool {
	return len(r.Errors) > 0
}

func (r *ValidationResult) Error() string {
	msgs := make([]string, len(r.Errors))
	for i, e := range r.Errors {
		msgs[i] = e.Error()
	}
	return strings.Join(msgs, "; ")
}

// validKeyRe matches POSIX-style env var names: uppercase letters, digits, underscores,
// must not start with a digit.
var validKeyRe = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

// ValidateOptions controls which checks are performed.
type ValidateOptions struct {
	// RequireUppercase enforces that all keys are UPPER_CASE.
	RequireUppercase bool
	// ForbidEmpty rejects keys whose value is an empty string.
	ForbidEmpty bool
}

// DefaultValidateOptions returns sensible defaults.
func DefaultValidateOptions() ValidateOptions {
	return ValidateOptions{
		RequireUppercase: false,
		ForbidEmpty:      false,
	}
}

// Validate checks every key/value pair in env against naming rules and the
// supplied options. It always returns a ValidationResult (never nil).
func Validate(env map[string]string, opts ValidateOptions) *ValidationResult {
	result := &ValidationResult{}

	for k, v := range env {
		if !validKeyRe.MatchString(k) {
			result.Errors = append(result.Errors, ValidationError{
				Key:     k,
				Message: "invalid identifier (must match [A-Za-z_][A-Za-z0-9_]*)",
			})
			continue
		}

		if opts.RequireUppercase && k != strings.ToUpper(k) {
			result.Errors = append(result.Errors, ValidationError{
				Key:     k,
				Message: "key must be UPPER_CASE",
			})
		}

		if opts.ForbidEmpty && v == "" {
			result.Errors = append(result.Errors, ValidationError{
				Key:     k,
				Message: "value must not be empty",
			})
		}
	}

	return result
}
