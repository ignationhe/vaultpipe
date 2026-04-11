package envfile

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

// ExportFormat defines the output format for exported secrets.
type ExportFormat string

const (
	FormatDotenv ExportFormat = "dotenv"
	FormatExport ExportFormat = "export"
	FormatJSON   ExportFormat = "json"
)

// ExportOptions controls how secrets are rendered during export.
type ExportOptions struct {
	Format ExportFormat
	Sorted bool
}

// DefaultExportOptions returns sensible defaults for Export.
func DefaultExportOptions() ExportOptions {
	return ExportOptions{
		Format: FormatDotenv,
		Sorted: true,
	}
}

// Export writes secrets from the given map to w using the specified format.
func Export(w io.Writer, secrets map[string]string, opts ExportOptions) error {
	keys := make([]string, 0, len(secrets))
	for k := range secrets {
		keys = append(keys, k)
	}
	if opts.Sorted {
		sort.Strings(keys)
	}

	switch opts.Format {
	case FormatExport:
		for _, k := range keys {
			v := secrets[k]
			if needsQuoting(v) {
				v = fmt.Sprintf(`"%s"`, v)
			}
			fmt.Fprintf(w, "export %s=%s\n", k, v)
		}
	case FormatJSON:
		pairs := make([]string, 0, len(keys))
		for _, k := range keys {
			pairs = append(pairs, fmt.Sprintf("  %q: %q", k, secrets[k]))
		}
		fmt.Fprintf(w, "{\n%s\n}\n", strings.Join(pairs, ",\n"))
	default: // FormatDotenv
		for _, k := range keys {
			v := secrets[k]
			if needsQuoting(v) {
				v = fmt.Sprintf(`"%s"`, v)
			}
			fmt.Fprintf(w, "%s=%s\n", k, v)
		}
	}
	return nil
}

// ExportToFile writes secrets to the given file path using the specified options.
func ExportToFile(path string, secrets map[string]string, opts ExportOptions) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("export: open %s: %w", path, err)
	}
	defer f.Close()
	return Export(f, secrets, opts)
}
