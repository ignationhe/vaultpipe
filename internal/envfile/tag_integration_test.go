package envfile

import (
	"os"
	"testing"
)

// TestTag_ThenWrite_RoundTrip tags an env map, writes only tagged keys, then
// parses the result and verifies the expected keys are present.
func TestTag_ThenWrite_RoundTrip(t *testing.T) {
	env := map[string]string{
		"DB_HOST":     "localhost",
		"DB_PASSWORD": "hunter2",
		"APP_PORT":    "9000",
	}

	opts := DefaultTagOptions()
	opts.SkipUntagged = true
	opts.Rules = map[string][]string{
		"database": {"DB_*"},
	}

	entries, err := Tag(env, opts)
	if err != nil {
		t.Fatalf("Tag: %v", err)
	}

	// Build filtered map from tagged entries
	filtered := make(map[string]string, len(entries))
	for _, e := range entries {
		filtered[e.Key] = env[e.Key]
	}

	tmp, err := os.CreateTemp(t.TempDir(), "tagged-*.env")
	if err != nil {
		t.Fatalf("create temp: %v", err)
	}
	tmp.Close()

	if err := Write(filtered, tmp.Name()); err != nil {
		t.Fatalf("Write: %v", err)
	}

	parsed, err := Parse(tmp.Name())
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}

	if len(parsed) != 2 {
		t.Errorf("expected 2 keys, got %d", len(parsed))
	}
	if parsed["DB_HOST"] != "localhost" {
		t.Errorf("DB_HOST mismatch: %q", parsed["DB_HOST"])
	}
	if _, ok := parsed["APP_PORT"]; ok {
		t.Error("APP_PORT should have been excluded")
	}
}

// TestTag_EmptyRulesTagsNothing ensures that with no rules every key is
// returned untagged when SkipUntagged is false.
func TestTag_EmptyRulesTagsNothing(t *testing.T) {
	env := map[string]string{"A": "1", "B": "2"}
	entries, err := Tag(env, DefaultTagOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(entries))
	}
	for _, e := range entries {
		if len(e.Tags) != 0 {
			t.Errorf("expected no tags on %q, got %v", e.Key, e.Tags)
		}
	}
}
