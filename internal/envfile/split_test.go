package envfile

import (
	"testing"
)

func baseEnvForSplit() map[string]string {
	return map[string]string{
		"APP__HOST":    "localhost",
		"APP__PORT":    "8080",
		"DB__HOST":     "db.internal",
		"DB__PASSWORD": "secret",
		"LOG_LEVEL":    "info",
	}
}

func TestSplit_GroupsByPrefix(t *testing.T) {
	result, err := Split(baseEnvForSplit(), DefaultSplitOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["APP"]["HOST"] != "localhost" {
		t.Errorf("expected APP.HOST=localhost, got %q", result["APP"]["HOST"])
	}
	if result["DB"]["PASSWORD"] != "secret" {
		t.Errorf("expected DB.PASSWORD=secret, got %q", result["DB"]["PASSWORD"])
	}
}

func TestSplit_UngroupedKey(t *testing.T) {
	result, err := Split(baseEnvForSplit(), DefaultSplitOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	def := result["default"]
	if def == nil {
		t.Fatal("expected a 'default' group for ungrouped keys")
	}
	if def["LOG_LEVEL"] != "info" {
		t.Errorf("expected LOG_LEVEL=info in default group, got %q", def["LOG_LEVEL"])
	}
}

func TestSplit_KeepPrefix(t *testing.T) {
	opts := DefaultSplitOptions()
	opts.KeepPrefix = true
	result, err := Split(baseEnvForSplit(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := result["APP"]["APP__HOST"]; !ok {
		t.Error("expected APP__HOST key to be preserved when KeepPrefix=true")
	}
}

func TestSplit_CustomDelimiter(t *testing.T) {
	env := map[string]string{
		"APP.HOST": "localhost",
		"APP.PORT": "9090",
		"ORPHAN":   "yes",
	}
	opts := DefaultSplitOptions()
	opts.Delimiter = "."
	result, err := Split(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["APP"]["HOST"] != "localhost" {
		t.Errorf("expected APP.HOST=localhost, got %q", result["APP"]["HOST"])
	}
	if result["default"]["ORPHAN"] != "yes" {
		t.Errorf("expected ORPHAN in default group")
	}
}

func TestSplit_NilEnvReturnsError(t *testing.T) {
	_, err := Split(nil, DefaultSplitOptions())
	if err == nil {
		t.Error("expected error for nil env map")
	}
}

func TestSplit_TrailingDelimiterGoesToUngrouped(t *testing.T) {
	env := map[string]string{
		"APP__": "trailing",
	}
	result, err := Split(env, DefaultSplitOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["default"]["APP__"] != "trailing" {
		t.Errorf("expected trailing-delimiter key in default group")
	}
}
