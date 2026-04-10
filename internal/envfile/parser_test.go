package envfile

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("failed to write temp env file: %v", err)
	}
	return path
}

func TestParse_BasicKeyValue(t *testing.T) {
	path := writeTempEnv(t, "DB_HOST=localhost\nDB_PORT=5432\n")
	got, err := Parse(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %q", got["DB_HOST"])
	}
	if got["DB_PORT"] != "5432" {
		t.Errorf("expected DB_PORT=5432, got %q", got["DB_PORT"])
	}
}

func TestParse_IgnoresComments(t *testing.T) {
	path := writeTempEnv(t, "# this is a comment\nAPI_KEY=secret\n")
	got, err := Parse(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := got["# this is a comment"]; ok {
		t.Error("comment line should not be parsed as a key")
	}
	if got["API_KEY"] != "secret" {
		t.Errorf("expected API_KEY=secret, got %q", got["API_KEY"])
	}
}

func TestParse_QuotedValues(t *testing.T) {
	path := writeTempEnv(t, `TOKEN="my token value"` + "\n" + `SECRET='another secret'` + "\n")
	got, err := Parse(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["TOKEN"] != "my token value" {
		t.Errorf("expected unquoted value, got %q", got["TOKEN"])
	}
	if got["SECRET"] != "another secret" {
		t.Errorf("expected unquoted value, got %q", got["SECRET"])
	}
}

func TestParse_NonExistentFile(t *testing.T) {
	got, err := Parse("/nonexistent/.env")
	if err != nil {
		t.Fatalf("expected no error for missing file, got: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected empty map for missing file, got %v", got)
	}
}

func TestParse_InvalidLine(t *testing.T) {
	path := writeTempEnv(t, "INVALID_LINE_NO_EQUALS\n")
	_, err := Parse(path)
	if err == nil {
		t.Error("expected error for invalid line format, got nil")
	}
}
