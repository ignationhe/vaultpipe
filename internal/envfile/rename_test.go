package envfile

import (
	"testing"
)

func baseEnv() map[string]string {
	return map[string]string{
		"DB_HOST":     "localhost",
		"DB_PORT":     "5432",
		"APP_SECRET":  "s3cr3t",
		"OLD_API_KEY": "key123",
	}
}

func TestRename_ExactKey(t *testing.T) {
	env := baseEnv()
	opts := DefaultRenameOptions()
	opts.Rules = []RenameRule{{From: "DB_HOST", To: "DATABASE_HOST"}}

	out, err := Rename(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DATABASE_HOST"] != "localhost" {
		t.Errorf("expected DATABASE_HOST=localhost, got %q", out["DATABASE_HOST"])
	}
	if _, ok := out["DB_HOST"]; ok {
		t.Error("expected DB_HOST to be deleted")
	}
}

func TestRename_KeepOriginal(t *testing.T) {
	env := baseEnv()
	opts := DefaultRenameOptions()
	opts.DeleteOriginal = false
	opts.Rules = []RenameRule{{From: "DB_HOST", To: "DATABASE_HOST"}}

	out, err := Rename(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DB_HOST"] != "localhost" {
		t.Error("expected DB_HOST to be preserved")
	}
	if out["DATABASE_HOST"] != "localhost" {
		t.Error("expected DATABASE_HOST to be set")
	}
}

func TestRename_MissingKeySkipped(t *testing.T) {
	env := baseEnv()
	opts := DefaultRenameOptions()
	opts.Rules = []RenameRule{{From: "NONEXISTENT", To: "SOMETHING"}}

	_, err := Rename(env, opts)
	if err != nil {
		t.Fatalf("expected no error when SkipMissing=true, got: %v", err)
	}
}

func TestRename_MissingKeyErrors(t *testing.T) {
	env := baseEnv()
	opts := DefaultRenameOptions()
	opts.SkipMissing = false
	opts.Rules = []RenameRule{{From: "NONEXISTENT", To: "SOMETHING"}}

	_, err := Rename(env, opts)
	if err == nil {
		t.Fatal("expected error when SkipMissing=false and key absent")
	}
}

func TestRename_PatternRule(t *testing.T) {
	env := baseEnv()
	opts := DefaultRenameOptions()
	opts.Rules = []RenameRule{{From: "^OLD_(.*)", To: "LEGACY_$1", Pattern: true}}

	out, err := Rename(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["LEGACY_API_KEY"] != "key123" {
		t.Errorf("expected LEGACY_API_KEY=key123, got %q", out["LEGACY_API_KEY"])
	}
	if _, ok := out["OLD_API_KEY"]; ok {
		t.Error("expected OLD_API_KEY to be removed")
	}
}

func TestRename_InvalidPatternReturnsError(t *testing.T) {
	env := baseEnv()
	opts := DefaultRenameOptions()
	opts.Rules = []RenameRule{{From: "[invalid", To: "X", Pattern: true}}

	_, err := Rename(env, opts)
	if err == nil {
		t.Fatal("expected error for invalid regex pattern")
	}
}
