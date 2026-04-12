package envfile

import (
	"os"
	"testing"
)

func TestRename_ThenWrite_RoundTrip(t *testing.T) {
	src := map[string]string{
		"OLD_HOST": "db.internal",
		"OLD_PORT": "3306",
	}

	opts := DefaultRenameOptions()
	opts.Rules = []RenameRule{
		{From: "OLD_HOST", To: "NEW_HOST"},
		{From: "OLD_PORT", To: "NEW_PORT"},
	}

	renamed, err := Rename(src, opts)
	if err != nil {
		t.Fatalf("Rename error: %v", err)
	}

	tmp, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatal(err)
	}
	tmp.Close()

	if err := Write(renamed, tmp.Name()); err != nil {
		t.Fatalf("Write error: %v", err)
	}

	parsed, err := Parse(tmp.Name())
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if parsed["NEW_HOST"] != "db.internal" {
		t.Errorf("expected NEW_HOST=db.internal, got %q", parsed["NEW_HOST"])
	}
	if parsed["NEW_PORT"] != "3306" {
		t.Errorf("expected NEW_PORT=3306, got %q", parsed["NEW_PORT"])
	}
	if _, ok := parsed["OLD_HOST"]; ok {
		t.Error("OLD_HOST should not appear in output")
	}
}

func TestRename_PatternThenWrite(t *testing.T) {
	src := map[string]string{
		"LEGACY_DB_URL":  "postgres://old",
		"LEGACY_DB_NAME": "mydb",
		"KEEP_THIS":      "yes",
	}

	opts := DefaultRenameOptions()
	opts.Rules = []RenameRule{
		{From: "^LEGACY_(.*)", To: "CURRENT_$1", Pattern: true},
	}

	renamed, err := Rename(src, opts)
	if err != nil {
		t.Fatalf("Rename error: %v", err)
	}

	if renamed["KEEP_THIS"] != "yes" {
		t.Error("KEEP_THIS should be preserved")
	}
	if renamed["CURRENT_DB_URL"] != "postgres://old" {
		t.Errorf("expected CURRENT_DB_URL, got %v", renamed)
	}
}
