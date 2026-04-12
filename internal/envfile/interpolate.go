package envfile

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// interpolateVarRe matches ${VAR_NAME} and $VAR_NAME patterns.
var interpolateVarRe = regexp.MustCompile(`\$\{([A-Za-z_][A-Za-z0-9_]*)\}|\$([A-Za-z_][A-Za-z0-9_]*)`)

// InterpolateOptions controls how variable interpolation behaves.
type InterpolateOptions struct {
	// FallbackToEnv allows falling back to OS environment variables.
	FallbackToEnv bool
	// ErrorOnMissing returns an error when a referenced variable is not found.
	ErrorOnMissing bool
	// DefaultValue is used when a variable is missing and ErrorOnMissing is false.
	DefaultValue string
}

// DefaultInterpolateOptions returns sensible defaults.
func DefaultInterpolateOptions() InterpolateOptions {
	return InterpolateOptions{
		FallbackToEnv:  true,
		ErrorOnMissing: false,
		DefaultValue:   "",
	}
}

// Interpolate resolves variable references within values of the provided map.
// Values may reference other keys in the same map using ${KEY} or $KEY syntax.
func Interpolate(env map[string]string, opts InterpolateOptions) (map[string]string, error) {
	result := make(map[string]string, len(env))
	for k, v := range env {
		resolved, err := interpolateValue(v, env, opts)
		if err != nil {
			return nil, fmt.Errorf("key %q: %w", k, err)
		}
		result[k] = resolved
	}
	return result, nil
}

func interpolateValue(val string, env map[string]string, opts InterpolateOptions) (string, error) {
	var lastErr error
	result := interpolateVarRe.ReplaceAllStringFunc(val, func(match string) string {
		name := extractVarName(match)
		if v, ok := env[name]; ok {
			return v
		}
		if opts.FallbackToEnv {
			if v, ok := os.LookupEnv(name); ok {
				return v
			}
		}
		if opts.ErrorOnMissing {
			lastErr = fmt.Errorf("undefined variable: %s", name)
			return match
		}
		return opts.DefaultValue
	})
	if lastErr != nil {
		return "", lastErr
	}
	return result, nil
}

func extractVarName(match string) string {
	match = strings.TrimPrefix(match, "$")
	match = strings.TrimPrefix(match, "{")
	match = strings.TrimSuffix(match, "}")
	return match
}
