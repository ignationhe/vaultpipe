package envfile

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"text/tabwriter"
)

// AuditFormat controls the output format for audit logs.
type AuditFormat string

const (
	AuditFormatText AuditFormat = "text"
	AuditFormatJSON AuditFormat = "json"
)

// WriteAuditLog writes the audit log to the given writer in the specified format.
func WriteAuditLog(log AuditLog, w io.Writer, format AuditFormat) error {
	switch format {
	case AuditFormatJSON:
		return writeAuditJSON(log, w)
	default:
		return writeAuditText(log, w)
	}
}

func writeAuditText(log AuditLog, w io.Writer) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "TIMESTAMP\tKEY\tACTION\tSOURCE")
	for _, e := range log.Entries {
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n",
			e.Timestamp.Format("2006-01-02T15:04:05Z"),
			e.Key,
			string(e.Action),
			e.Source,
		)
	}
	return tw.Flush()
}

func writeAuditJSON(log AuditLog, w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(log.Entries)
}

// WriteAuditLogToFile writes the audit log to a file, appending if it exists.
func WriteAuditLogToFile(log AuditLog, path string, format AuditFormat) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return fmt.Errorf("audit: open file: %w", err)
	}
	defer f.Close()
	return WriteAuditLog(log, f, format)
}
