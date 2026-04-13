package envfile

import (
	"testing"
)

func baseEnvForPatch() map[string]string {
	return map[string]string{
		"APP_HOST": "localhost",
		"APP_PORT": "8080",
		"DB_PASS":  "secret",
	}
}

func TestPatch_SetAddsNewKey(t *testing.T) {
	out, err := Patch(baseEnvForPatch(), []PatchRule{
		{Op: PatchSet, Key: "NEW_KEY", Value: "hello"},
	}, DefaultPatchOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["NEW_KEY"] != "hello" {
		t.Errorf("expected NEW_KEY=hello, got %q", out["NEW_KEY"])
	}
}

func TestPatch_SetOverwritesExistingKey(t *testing.T) {
	out, err := Patch(baseEnvForPatch(), []PatchRule{
		{Op: PatchSet, Key: "APP_PORT", Value: "9090"},
	}, DefaultPatchOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["APP_PORT"] != "9090" {
		t.Errorf("expected APP_PORT=9090, got %q", out["APP_PORT"])
	}
}

func TestPatch_DeleteRemovesKey(t *testing.T) {
	out, err := Patch(baseEnvForPatch(), []PatchRule{
		{Op: PatchDelete, Key: "DB_PASS"},
	}, DefaultPatchOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["DB_PASS"]; ok {
		t.Error("expected DB_PASS to be deleted")
	}
}

func TestPatch_DeleteMissingKeyIgnoredByDefault(t *testing.T) {
	_, err := Patch(baseEnvForPatch(), []PatchRule{
		{Op: PatchDelete, Key: "DOES_NOT_EXIST"},
	}, DefaultPatchOptions())
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestPatch_DeleteMissingKeyErrorsWhenNotIgnoring(t *testing.T) {
	opts := PatchOptions{IgnoreMissing: false}
	_, err := Patch(baseEnvForPatch(), []PatchRule{
		{Op: PatchDelete, Key: "DOES_NOT_EXIST"},
	}, opts)
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestPatch_RenameMovesKey(t *testing.T) {
	out, err := Patch(baseEnvForPatch(), []PatchRule{
		{Op: PatchRename, Key: "APP_HOST", NewKey: "SERVICE_HOST"},
	}, DefaultPatchOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["SERVICE_HOST"] != "localhost" {
		t.Errorf("expected SERVICE_HOST=localhost, got %q", out["SERVICE_HOST"])
	}
	if _, ok := out["APP_HOST"]; ok {
		t.Error("expected APP_HOST to be removed after rename")
	}
}

func TestPatch_UnknownOpReturnsError(t *testing.T) {
	_, err := Patch(baseEnvForPatch(), []PatchRule{
		{Op: "upsert", Key: "X"},
	}, DefaultPatchOptions())
	if err == nil {
		t.Fatal("expected error for unknown op")
	}
}

func TestPatch_OriginalMapUnmodified(t *testing.T) {
	env := baseEnvForPatch()
	_, err := Patch(env, []PatchRule{
		{Op: PatchSet, Key: "APP_PORT", Value: "1234"},
		{Op: PatchDelete, Key: "DB_PASS"},
	}, DefaultPatchOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["APP_PORT"] != "8080" {
		t.Error("original map should not be modified")
	}
	if _, ok := env["DB_PASS"]; !ok {
		t.Error("original map should still contain DB_PASS")
	}
}
