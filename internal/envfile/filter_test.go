package envfile

import (
	"testing"
)

var sampleSecrets = map[string]string{
	"APP_HOST":    "localhost",
	"APP_PORT":    "8080",
	"DB_HOST":     "db.local",
	"DB_PASSWORD": "secret",
	"DEBUG":       "true",
}

func TestFilter_ByPrefix(t *testing.T) {
	result := Filter(sampleSecrets, FilterOptions{Prefix: "APP_"})
	if len(result) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(result))
	}
	if _, ok := result["APP_HOST"]; !ok {
		t.Error("expected APP_HOST")
	}
	if _, ok := result["APP_PORT"]; !ok {
		t.Error("expected APP_PORT")
	}
}

func TestFilter_ByInclude(t *testing.T) {
	result := Filter(sampleSecrets, FilterOptions{Include: []string{"DEBUG", "DB_HOST"}})
	if len(result) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(result))
	}
	if result["DEBUG"] != "true" {
		t.Errorf("unexpected value for DEBUG: %s", result["DEBUG"])
	}
}

func TestFilter_ByExclude(t *testing.T) {
	result := Filter(sampleSecrets, FilterOptions{Exclude: []string{"DB_PASSWORD", "DEBUG"}})
	if len(result) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(result))
	}
	if _, ok := result["DB_PASSWORD"]; ok {
		t.Error("DB_PASSWORD should have been excluded")
	}
}

func TestFilter_PrefixAndExcludeCombined(t *testing.T) {
	result := Filter(sampleSecrets, FilterOptions{
		Prefix:  "DB_",
		Exclude: []string{"DB_PASSWORD"},
	})
	if len(result) != 1 {
		t.Fatalf("expected 1 key, got %d", len(result))
	}
	if _, ok := result["DB_HOST"]; !ok {
		t.Error("expected DB_HOST")
	}
}

func TestFilter_EmptyOptsReturnsAll(t *testing.T) {
	result := Filter(sampleSecrets, FilterOptions{})
	if len(result) != len(sampleSecrets) {
		t.Fatalf("expected %d keys, got %d", len(sampleSecrets), len(result))
	}
}

func TestStripPrefix_RemovesPrefix(t *testing.T) {
	input := map[string]string{"APP_HOST": "localhost", "APP_PORT": "8080"}
	result := StripPrefix(input, "APP_")
	if result["HOST"] != "localhost" {
		t.Errorf("expected HOST=localhost, got %s", result["HOST"])
	}
	if result["PORT"] != "8080" {
		t.Errorf("expected PORT=8080, got %s", result["PORT"])
	}
}

func TestStripPrefix_EmptyPrefixNoOp(t *testing.T) {
	input := map[string]string{"KEY": "val"}
	result := StripPrefix(input, "")
	if result["KEY"] != "val" {
		t.Error("expected KEY to be unchanged")
	}
}
