package envfile

import (
	"testing"
)

func baseEnvForPluck() map[string]string {
	return map[string]string{
		"DB_HOST":     "localhost",
		"DB_PORT":     "5432",
		"API_KEY":     "secret",
		"LOG_LEVEL":   "info",
	}
}

func TestPluck_SelectedKeys(t *testing.T) {
	src := baseEnvForPluck()
	out, err := Pluck(src, PluckOptions{Keys: []string{"DB_HOST", "API_KEY"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(out))
	}
	if out["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %s", out["DB_HOST"])
	}
	if out["API_KEY"] != "secret" {
		t.Errorf("expected API_KEY=secret, got %s", out["API_KEY"])
	}
}

func TestPluck_EmptyKeysReturnsAll(t *testing.T) {
	src := baseEnvForPluck()
	out, err := Pluck(src, PluckOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != len(src) {
		t.Errorf("expected %d keys, got %d", len(src), len(out))
	}
}

func TestPluck_MissingKeySkippedByDefault(t *testing.T) {
	src := baseEnvForPluck()
	out, err := Pluck(src, PluckOptions{Keys: []string{"DB_HOST", "MISSING"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 1 {
		t.Errorf("expected 1 key, got %d", len(out))
	}
}

func TestPluck_ErrorOnMissingKey(t *testing.T) {
	src := baseEnvForPluck()
	_, err := Pluck(src, PluckOptions{Keys: []string{"MISSING"}, ErrorOnMiss: true})
	if err == nil {
		t.Fatal("expected error for missing key, got nil")
	}
}

func TestPluck_DoesNotMutateSrc(t *testing.T) {
	src := baseEnvForPluck()
	out, _ := Pluck(src, PluckOptions{Keys: []string{"DB_HOST"}})
	out["DB_HOST"] = "changed"
	if src["DB_HOST"] != "localhost" {
		t.Error("Pluck mutated the source map")
	}
}
