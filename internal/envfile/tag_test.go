package envfile

import (
	"testing"
)

func baseEnvForTag() map[string]string {
	return map[string]string{
		"DB_HOST":     "localhost",
		"DB_PASSWORD": "secret",
		"APP_PORT":    "8080",
		"APP_DEBUG":   "true",
		"LOG_LEVEL":   "info",
	}
}

func TestTag_ExactKeyMatch(t *testing.T) {
	env := baseEnvForTag()
	opts := DefaultTagOptions()
	opts.Rules = map[string][]string{
		"secret": {"DB_PASSWORD"},
	}
	entries, err := Tag(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, e := range entries {
		if e.Key == "DB_PASSWORD" {
			if len(e.Tags) != 1 || e.Tags[0] != "secret" {
				t.Errorf("expected tag 'secret' on DB_PASSWORD, got %v", e.Tags)
			}
			return
		}
	}
	t.Error("DB_PASSWORD not found in entries")
}

func TestTag_GlobPattern(t *testing.T) {
	env := baseEnvForTag()
	opts := DefaultTagOptions()
	opts.Rules = map[string][]string{
		"database": {"DB_*"},
	}
	entries, err := Tag(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tagged := map[string]bool{}
	for _, e := range entries {
		for _, tag := range e.Tags {
			if tag == "database" {
				tagged[e.Key] = true
			}
		}
	}
	if !tagged["DB_HOST"] || !tagged["DB_PASSWORD"] {
		t.Errorf("expected DB_HOST and DB_PASSWORD tagged as 'database', got %v", tagged)
	}
}

func TestTag_SkipUntagged(t *testing.T) {
	env := baseEnvForTag()
	opts := DefaultTagOptions()
	opts.SkipUntagged = true
	opts.Rules = map[string][]string{
		"app": {"APP_*"},
	}
	entries, err := Tag(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(entries))
	}
	for _, e := range entries {
		if !startsWith(e.Key, "APP_") {
			t.Errorf("unexpected key %q in skip-untagged result", e.Key)
		}
	}
}

func TestTag_MultipleTagsOnSameKey(t *testing.T) {
	env := map[string]string{"DB_PASSWORD": "x"}
	opts := DefaultTagOptions()
	opts.Rules = map[string][]string{
		"secret":   {"DB_PASSWORD"},
		"database": {"DB_*"},
	}
	entries, err := Tag(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if len(entries[0].Tags) != 2 {
		t.Errorf("expected 2 tags, got %v", entries[0].Tags)
	}
}

func TestTag_NilEnvReturnsError(t *testing.T) {
	_, err := Tag(nil, DefaultTagOptions())
	if err == nil {
		t.Error("expected error for nil env")
	}
}

func TestTag_InvalidPatternReturnsError(t *testing.T) {
	env := map[string]string{"KEY": "val"}
	opts := DefaultTagOptions()
	opts.Rules = map[string][]string{
		"bad": {"[invalid"},
	}
	_, err := Tag(env, opts)
	if err == nil {
		t.Error("expected error for invalid pattern")
	}
}

func startsWith(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}
