package envfile

import (
	"os"
	"path/filepath"
	"testing"
)

// TestClone_ThenDiff verifies that cloning and then diffing the two files
// produces no changes when no transformation is applied.
func TestClone_ThenDiff(t *testing.T) {
	src := writeTempEnvClone(t, "DB_HOST=localhost\nDB_PORT=5432\n")
	dst := filepath.Join(t.TempDir(), "clone.env")

	_, err := Clone(src, dst, DefaultCloneOptions())
	if err != nil {
		t.Fatalf("clone: %v", err)
	}

	srcEnv, _ := Parse(src)
	dstEnv, _ := Parse(dst)

	diffs := Diff(srcEnv, dstEnv)
	if HasChanges(diffs) {
		t.Errorf("expected no changes after clone, got: %v", diffs)
	}
}

// TestClone_ThenMerge verifies that cloning into an existing file via Overwrite
// and then merging with extra keys preserves both sets.
func TestClone_ThenMerge(t *testing.T) {
	src := writeTempEnvClone(t, "SECRET=abc123\n")
	dst := filepath.Join(t.TempDir(), "merged.env")

	// Write a pre-existing key to dst.
	if err := os.WriteFile(dst, []byte("EXISTING=yes\n"), 0o600); err != nil {
		t.Fatalf("setup: %v", err)
	}

	opts := DefaultCloneOptions()
	opts.Overwrite = true
	// Only clone SECRET; EXISTING will be lost (clone replaces).
	_, err := Clone(src, dst, opts)
	if err != nil {
		t.Fatalf("clone: %v", err)
	}

	// Now merge EXISTING back.
	extra := map[string]string{"EXISTING": "yes"}
	if err := Merge(dst, extra); err != nil {
		t.Fatalf("merge: %v", err)
	}

	env, _ := Parse(dst)
	if env["SECRET"] != "abc123" {
		t.Errorf("want SECRET=abc123, got %q", env["SECRET"])
	}
	if env["EXISTING"] != "yes" {
		t.Errorf("want EXISTING=yes, got %q", env["EXISTING"])
	}
}
