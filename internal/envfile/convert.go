package envfile

import (
	"fmt"
	"strings"
)

// ConvertFormat represents a supported conversion target format.
type ConvertFormat string

const (
	FormatDotenv ConvertFormat = "dotenv"
	FormatYAML   ConvertFormat = "yaml"
	FormatTOML   ConvertFormat = "toml"
)

// DefaultConvertOptions returns sensible defaults for Convert.
func DefaultConvertOptions() ConvertOptions {
	return ConvertOptions{
		Format: FormatDotenv,
		Sort:   true,
	}
}

// ConvertOptions controls the behaviour of Convert.
type ConvertOptions struct {
	Format ConvertFormat
	Sort   bool
	Prefix string // optional prefix to prepend to every key
}

// Convert serialises env map to the requested text format.
// Returns the formatted string or an error for unknown formats.
func Convert(env map[string]string, opts ConvertOptions) (string, error) {
	keys := sortedKeys(env, opts.Sort)

	switch opts.Format {
	case FormatDotenv, "":
		return convertDotenv(env, keys, opts.Prefix), nil
	case FormatYAML:
		return convertYAML(env, keys, opts.Prefix), nil
	case FormatTOML:
		return convertTOML(env, keys, opts.Prefix), nil
	default:
		return "", fmt.Errorf("convert: unknown format %q", opts.Format)
	}
}

func applyPrefix(key, prefix string) string {
	if prefix == "" {
		return key
	}
	return prefix + key
}

func convertDotenv(env map[string]string, keys []string, prefix string) string {
	var sb strings.Builder
	for _, k := range keys {
		v := env[k]
		fullKey := applyPrefix(k, prefix)
		if needsQuoting(v) {
			fmt.Fprintf(&sb, "%s=\"%s\"\n", fullKey, v)
		} else {
			fmt.Fprintf(&sb, "%s=%s\n", fullKey, v)
		}
	}
	return sb.String()
}

func convertYAML(env map[string]string, keys []string, prefix string) string {
	var sb strings.Builder
	for _, k := range keys {
		v := env[k]
		fullKey := applyPrefix(k, prefix)
		if needsQuoting(v) {
			fmt.Fprintf(&sb, "%s: \"%s\"\n", fullKey, v)
		} else {
			fmt.Fprintf(&sb, "%s: %s\n", fullKey, v)
		}
	}
	return sb.String()
}

func convertTOML(env map[string]string, keys []string, prefix string) string {
	var sb strings.Builder
	for _, k := range keys {
		v := env[k]
		fullKey := applyPrefix(k, prefix)
		fmt.Fprintf(&sb, "%s = \"%s\"\n", fullKey, v)
	}
	return sb.String()
}

// sortedKeys returns map keys either sorted or in iteration order.
func sortedKeys(env map[string]string, sort bool) []string {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	if sort {
		sliceSort(keys)
	}
	return keys
}

func sliceSort(s []string) {
	for i := 1; i < len(s); i++ {
		for j := i; j > 0 && s[j] < s[j-1]; j-- {
			s[j], s[j-1] = s[j-1], s[j]
		}
	}
}
