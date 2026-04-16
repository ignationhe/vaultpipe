package envfile

import "fmt"

// SummaryOptions controls how a diff summary is rendered.
type SummaryOptions struct {
	ShowValues bool
	Redact     bool
}

// DefaultSummaryOptions returns sensible defaults.
func DefaultSummaryOptions() SummaryOptions {
	return SummaryOptions{
		ShowValues: false,
		Redact:     true,
	}
}

// DiffSummary holds counts and human-readable lines for a diff.
type DiffSummary struct {
	Added   int
	Removed int
	Updated int
	Lines   []string
}

// Summarize converts a []DiffEntry into a DiffSummary.
func Summarize(diffs []DiffEntry, opts SummaryOptions) DiffSummary {
	s := DiffSummary{}
	for _, d := range diffs {
		switch d.Status {
		case StatusAdded:
			s.Added++
			s.Lines = append(s.Lines, formatLine("+", d, opts))
		case StatusRemoved:
			s.Removed++
			s.Lines = append(s.Lines, formatLine("-", d, opts))
		case StatusUpdated:
			s.Updated++
			s.Lines = append(s.Lines, formatLine("~", d, opts))
		}
	}
	return s
}

func formatLine(prefix string, d DiffEntry, opts SummaryOptions) string {
	if !opts.ShowValues {
		return fmt.Sprintf("%s %s", prefix, d.Key)
	}
	old, new_ := d.OldValue, d.NewValue
	if opts.Redact {
		if old != "" {
			old = "***"
		}
		if new_ != "" {
			new_ = "***"
		}
	}
	switch d.Status {
	case StatusAdded:
		return fmt.Sprintf("%s %s=%s", prefix, d.Key, new_)
	case StatusRemoved:
		return fmt.Sprintf("%s %s=%s", prefix, d.Key, old)
	default:
		return fmt.Sprintf("%s %s: %s -> %s", prefix, d.Key, old, new_)
	}
}
