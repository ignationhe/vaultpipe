package envfile

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPatch_ThenWrite_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	initial := map[string]string{
		"HOST": "localhost",
		"PORT": "5432",
		"PASS": "old",
	}
	if err := Write(path, initial); err != nil {
		t.Fatalf("Write: %v", err)
	}

	parsed, err := Parse(path)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}

	patched, err := Patch(parsed, []PatchRule{
		{Op: PatchSet, Key: "PASS", Value: "new-secret"},
		{Op: PatchRename, Key: "HOST", NewKey: "DB_HOST"},
		{Op: PatchDelete, Key: "PORT"},
	}, DefaultPatchOptions())
	if err != nil {
		t.Fatalf("Patch: %v", err)
	}

	if err := Write(path, patched); err != nil {
		t.Fatalf("Write after patch: %v", err)
	}

	final, err := Parse(path)
	if err != nil {
		t.Fatalf("Parse after write: %v", err)
	}

	if final["PASS"] != "new-secret" {
		t.Errorf("expected PASS=new-secret, got %q", final["PASS"])
	}
	if final["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %q", final["DB_HOST"])
	}
	if _, ok := final["HOST"]; ok {
		t.Error("HOST should have been renamed to DB_HOST")
	}
	if _, ok := final["PORT"]; ok {
		t.Error("PORT should have been deleted")
	}
}

func TestPatch_ThenDiff_DetectsChanges(t *testing.T) {
	env := map[string]string{
		"A": "1",
		"B": "2",
		"C": "3",
	}

	patched, err := Patch(env, []PatchRule{
		{Op: PatchSet, Key: "A", Value: "99"},
		{Op: PatchDelete, Key: "C"},
	}, DefaultPatchOptions())
	if err != nil {
		t.Fatalf("Patch: %v", err)
	}

	diffs := Diff(env, patched)
	if !HasChanges(diffs) {
		t.Fatal("expected changes after patch")
	}

	statuses := map[string]string{}
	for _, d := range diffs {
		statuses[d.Key] = string(d.Status)
	}
	if statuses["A"] != "updated" {
		t.Errorf("expected A=updated, got %q", statuses["A"])
	}
	if statuses["C"] != "removed" {
		t.Errorf("expected C=removed, got %q", statuses["C"])
	}

	_ = os.Getenv("CI") // suppress unused import warning
}
