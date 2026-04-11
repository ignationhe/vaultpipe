package envfile

import (
	"fmt"
	"sort"
	"strings"
)

// FormatOptions controls how the env file output is rendered.
type FormatOptions struct {
	// SortKeys sorts keys alphabetically in the output.
	SortKeys bool
	// SectionComment adds a header comment to the output.
	SectionComment string
	// BlankLineBetweenKeys inserts a blank line between each key-value pair.
	BlankLineBetweenKeys bool
}

// DefaultFormatOptions returns sensible defaults for formatting.
func DefaultFormatOptions() FormatOptions {
	return FormatOptions{
		SortKeys:             true,
		SectionComment:       "",
		BlankLineBetweenKeys: false,
	}
}

// Format renders a map of key-value pairs into an env file string
// according to the provided FormatOptions.
func Format(secrets map[string]string, opts FormatOptions) string {
	var sb strings.Builder

	if opts.SectionComment != "" {
		sb.WriteString("# ")
		sb.WriteString(opts.SectionComment)
		sb.WriteString("\n")
	}

	keys := make([]string, 0, len(secrets))
	for k := range secrets {
		keys = append(keys, k)
	}

	if opts.SortKeys {
		sort.Strings(keys)
	}

	for i, k := range keys {
		v := secrets[k]
		if needsQuoting(v) {
			v = fmt.Sprintf("%q", v)
		}
		sb.WriteString(fmt.Sprintf("%s=%s\n", k, v))
		if opts.BlankLineBetweenKeys && i < len(keys)-1 {
			sb.WriteString("\n")
		}
	}

	return sb.String()
}
