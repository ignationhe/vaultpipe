package envfile

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// TemplateOptions controls how template rendering behaves.
type TemplateOptions struct {
	// MissingKey defines behaviour when a referenced key is absent.
	// "error" (default) returns an error; "keep" leaves the placeholder; "empty" substitutes "".
	MissingKey string
}

// DefaultTemplateOptions returns sensible defaults.
func DefaultTemplateOptions() TemplateOptions {
	return TemplateOptions{MissingKey: "error"}
}

var placeholderRe = regexp.MustCompile(`\$\{([A-Za-z_][A-Za-z0-9_]*)\}`)

// Render replaces ${KEY} placeholders in src with values from vars.
// Any key that appears in vars is substituted; behaviour for missing keys is
// controlled by opts.MissingKey.
func Render(src string, vars map[string]string, opts TemplateOptions) (string, error) {
	var renderErr error
	result := placeholderRe.ReplaceAllStringFunc(src, func(match string) string {
		if renderErr != nil {
			return match
		}
		key := placeholderRe.FindStringSubmatch(match)[1]
		val, ok := vars[key]
		if !ok {
			switch opts.MissingKey {
			case "keep":
				return match
			case "empty":
				return ""
			default:
				renderErr = fmt.Errorf("template: missing key %q", key)
				return match
			}
		}
		return val
	})
	if renderErr != nil {
		return "", renderErr
	}
	return result, nil
}

// RenderMap applies Render to every value in the map and returns a new map.
func RenderMap(m map[string]string, vars map[string]string, opts TemplateOptions) (map[string]string, error) {
	out := make(map[string]string, len(m))
	for k, v := range m {
		rendered, err := Render(v, vars, opts)
		if err != nil {
			return nil, fmt.Errorf("key %q: %w", k, err)
		}
		out[k] = rendered
	}
	return out, nil
}

// RenderFile reads a template file, substitutes vars, and writes the result to dst.
func RenderFile(srcPath, dstPath string, vars map[string]string, opts TemplateOptions) error {
	data, err := os.ReadFile(srcPath)
	if err != nil {
		return fmt.Errorf("template: read %s: %w", srcPath, err)
	}
	output, err := Render(string(data), vars, opts)
	if err != nil {
		return err
	}
	if err := os.WriteFile(dstPath, []byte(strings.TrimRight(output, "\n")+"\n"), 0o600); err != nil {
		return fmt.Errorf("template: write %s: %w", dstPath, err)
	}
	return nil
}
