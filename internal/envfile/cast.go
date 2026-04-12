package envfile

import (
	"fmt"
	"strconv"
	"strings"
)

// CastType represents the target type for casting an env value.
type CastType string

const (
	CastString CastType = "string"
	CastInt    CastType = "int"
	CastFloat  CastType = "float"
	CastBool   CastType = "bool"
)

// CastRule defines a casting rule for a specific key.
type CastRule struct {
	Key  string
	Type CastType
}

// CastResult holds the original string value and its cast representation.
type CastResult struct {
	Key      string
	Raw      string
	Cast     interface{}
	CastType CastType
}

// CastOptions configures the Cast operation.
type CastOptions struct {
	// Rules maps keys to their desired cast type.
	Rules []CastRule
	// SkipUnknown ignores keys not present in Rules instead of keeping them as strings.
	SkipUnknown bool
}

// DefaultCastOptions returns sensible defaults.
func DefaultCastOptions() CastOptions {
	return CastOptions{
		Rules:       []CastRule{},
		SkipUnknown: false,
	}
}

// Cast converts env map values to typed representations based on provided rules.
// Keys not matching any rule are returned as CastString unless SkipUnknown is true.
func Cast(env map[string]string, opts CastOptions) ([]CastResult, error) {
	ruleMap := make(map[string]CastType, len(opts.Rules))
	for _, r := range opts.Rules {
		ruleMap[r.Key] = r.Type
	}

	results := make([]CastResult, 0, len(env))

	for k, v := range env {
		ct, ok := ruleMap[k]
		if !ok {
			if opts.SkipUnknown {
				continue
			}
			ct = CastString
		}

		casted, err := castValue(v, ct)
		if err != nil {
			return nil, fmt.Errorf("cast: key %q value %q cannot be cast to %s: %w", k, v, ct, err)
		}

		results = append(results, CastResult{
			Key:      k,
			Raw:      v,
			Cast:     casted,
			CastType: ct,
		})
	}

	return results, nil
}

func castValue(v string, ct CastType) (interface{}, error) {
	switch ct {
	case CastInt:
		return strconv.Atoi(strings.TrimSpace(v))
	case CastFloat:
		return strconv.ParseFloat(strings.TrimSpace(v), 64)
	case CastBool:
		return strconv.ParseBool(strings.TrimSpace(v))
	default:
		return v, nil
	}
}
