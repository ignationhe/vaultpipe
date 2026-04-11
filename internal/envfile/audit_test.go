package envfile

import (
	"testing"
)

func makeTestDiffs() []DiffEntry {
	return []DiffEntry{
		{Key: "DB_HOST", Status: StatusAdded, OldValue: "", NewValue: "localhost"},
		{Key: "DB_PASS", Status: StatusUpdated, OldValue: "old", NewValue: "new"},
		{Key: "API_KEY", Status: StatusRemoved, OldValue: "secret", NewValue: ""},
		{Key: "APP_ENV", Status: StatusUnchanged, OldValue: "prod", NewValue: "prod"},
	}
}

func TestAudit_ExcludesUnchangedByDefault(t *testing.T) {
	diffs := makeTestDiffs()
	log := Audit(diffs, DefaultAuditOptions())
	for _, e := range log.Entries {
		if e.Action == AuditUnchanged {
			t.Errorf("expected unchanged entries to be excluded, got key %s", e.Key)
		}
	}
	if len(log.Entries) != 3 {
		t.Errorf("expected 3 entries, got %d", len(log.Entries))
	}
}

func TestAudit_IncludesUnchangedWhenRequested(t *testing.T) {
	diffs := makeTestDiffs()
	opts := DefaultAuditOptions()
	opts.IncludeUnchanged = true
	log := Audit(diffs, opts)
	if len(log.Entries) != 4 {
		t.Errorf("expected 4 entries, got %d", len(log.Entries))
	}
}

func TestAudit_RedactsValuesByDefault(t *testing.T) {
	diffs := makeTestDiffs()
	log := Audit(diffs, DefaultAuditOptions())
	for _, e := range log.Entries {
		if e.Action == AuditAdded && e.NewValue != "***" {
			t.Errorf("expected redacted new value for %s, got %q", e.Key, e.NewValue)
		}
		if e.Action == AuditRemoved && e.OldValue != "***" {
			t.Errorf("expected redacted old value for %s, got %q", e.Key, e.OldValue)
		}
	}
}

func TestAudit_PlainValuesWhenRedactDisabled(t *testing.T) {
	diffs := makeTestDiffs()
	opts := DefaultAuditOptions()
	opts.RedactValues = false
	log := Audit(diffs, opts)
	for _, e := range log.Entries {
		if e.Key == "DB_HOST" && e.NewValue != "localhost" {
			t.Errorf("expected plain value, got %q", e.NewValue)
		}
	}
}

func TestAudit_SetsSource(t *testing.T) {
	diffs := makeTestDiffs()
	opts := DefaultAuditOptions()
	opts.Source = "vault://secret/app"
	log := Audit(diffs, opts)
	for _, e := range log.Entries {
		if e.Source != "vault://secret/app" {
			t.Errorf("expected source to be set, got %q", e.Source)
		}
	}
}

func TestAuditLog_Summary(t *testing.T) {
	diffs := makeTestDiffs()
	opts := DefaultAuditOptions()
	opts.IncludeUnchanged = true
	log := Audit(diffs, opts)
	summary := log.Summary()
	expected := "added=1 removed=1 updated=1 unchanged=1"
	if summary != expected {
		t.Errorf("expected %q, got %q", expected, summary)
	}
}

func TestAudit_EmptyDiffs(t *testing.T) {
	log := Audit([]DiffEntry{}, DefaultAuditOptions())
	if len(log.Entries) != 0 {
		t.Errorf("expected empty log, got %d entries", len(log.Entries))
	}
	if log.Summary() != "added=0 removed=0 updated=0 unchanged=0" {
		t.Errorf("unexpected summary: %s", log.Summary())
	}
}
