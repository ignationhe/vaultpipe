package envfile

import (
	"fmt"
	"time"
)

// AuditAction describes what happened to a key.
type AuditAction string

const (
	AuditAdded   AuditAction = "added"
	AuditRemoved AuditAction = "removed"
	AuditUpdated AuditAction = "updated"
	AuditUnchanged AuditAction = "unchanged"
)

// AuditEntry records a single change event.
type AuditEntry struct {
	Timestamp time.Time
	Key       string
	Action    AuditAction
	OldValue  string
	NewValue  string
	Source    string
}

// AuditLog holds a collection of audit entries.
type AuditLog struct {
	Entries []AuditEntry
}

// DefaultAuditOptions returns options with sensible defaults.
func DefaultAuditOptions() AuditOptions {
	return AuditOptions{
		IncludeUnchanged: false,
		RedactValues:     true,
	}
}

// AuditOptions controls audit log behaviour.
type AuditOptions struct {
	IncludeUnchanged bool
	RedactValues     bool
	Source           string
}

// Audit builds an AuditLog from a diff result.
func Audit(diffs []DiffEntry, opts AuditOptions) AuditLog {
	log := AuditLog{}
	ts := time.Now().UTC()

	for _, d := range diffs {
		if d.Status == StatusUnchanged && !opts.IncludeUnchanged {
			continue
		}

		old := d.OldValue
		new := d.NewValue
		if opts.RedactValues {
			old = redactValue(old)
			new = redactValue(new)
		}

		log.Entries = append(log.Entries, AuditEntry{
			Timestamp: ts,
			Key:       d.Key,
			Action:    statusToAction(d.Status),
			OldValue:  old,
			NewValue:  new,
			Source:    opts.Source,
		})
	}
	return log
}

func redactValue(v string) string {
	if v == "" {
		return ""
	}
	return "***"
}

func statusToAction(s DiffStatus) AuditAction {
	switch s {
	case StatusAdded:
		return AuditAdded
	case StatusRemoved:
		return AuditRemoved
	case StatusUpdated:
		return AuditUpdated
	default:
		return AuditUnchanged
	}
}

// Summary returns a human-readable summary of the audit log.
func (a AuditLog) Summary() string {
	var added, removed, updated, unchanged int
	for _, e := range a.Entries {
		switch e.Action {
		case AuditAdded:
			added++
		case AuditRemoved:
			removed++
		case AuditUpdated:
			updated++
		case AuditUnchanged:
			unchanged++
		}
	}
	return fmt.Sprintf("added=%d removed=%d updated=%d unchanged=%d", added, removed, updated, unchanged)
}
