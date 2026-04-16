package envfile

import (
	"strings"
	"testing"
)

func makeDiffs() []DiffEntry {
	return []DiffEntry{
		{Key: "FOO", Status: StatusAdded, NewValue: "bar"},
		{Key: "OLD", Status: StatusRemoved, OldValue: "gone"},
		{Key: "DB", Status: StatusUpdated, OldValue: "v1", NewValue: "v2"},
		{Key: "SAME", Status: StatusUnchanged, OldValue: "x", NewValue: "x"},
	}
}

func TestSummarize_Counts(t *testing.T) {
	s := Summarize(makeDiffs(), DefaultSummaryOptions())
	if s.Added != 1 || s.Removed != 1 || s.Updated != 1 {
		t.Fatalf("unexpected counts: %+v", s)
	}
	if len(s.Lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(s.Lines))
	}
}

func TestSummarize_LinePrefixes(t *testing.T) {
	s := Summarize(makeDiffs(), DefaultSummaryOptions())
	prefixes := map[string]bool{}
	for _, l := range s.Lines {
		prefixes[string(l[0])] = true
	}
	for _, p := range []string{"+", "-", "~"} {
		if !prefixes[p] {
			t.Errorf("missing prefix %q", p)
		}
	}
}

func TestSummarize_ShowValuesRedacted(t *testing.T) {
	opts := SummaryOptions{ShowValues: true, Redact: true}
	s := Summarize(makeDiffs(), opts)
	for _, l := range s.Lines {
		if strings.Contains(l, "bar") || strings.Contains(l, "gone") || strings.Contains(l, "v1") {
			t.Errorf("line should be redacted: %s", l)
		}
		if !strings.Contains(l, "***") {
			t.Errorf("line should contain redaction marker: %s", l)
		}
	}
}

func TestSummarize_ShowValuesPlain(t *testing.T) {
	opts := SummaryOptions{ShowValues: true, Redact: false}
	s := Summarize(makeDiffs(), opts)
	found := false
	for _, l := range s.Lines {
		if strings.Contains(l, "v1") && strings.Contains(l, "v2") {
			found = true
		}
	}
	if !found {
		t.Error("expected plain updated values in output")
	}
}

func TestSummarize_EmptyDiff(t *testing.T) {
	s := Summarize([]DiffEntry{}, DefaultSummaryOptions())
	if s.Added != 0 || s.Removed != 0 || s.Updated != 0 || len(s.Lines) != 0 {
		t.Error("expected zero summary for empty diff")
	}
}
