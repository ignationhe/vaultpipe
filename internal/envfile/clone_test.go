package envfile

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempEnvClone(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "clone-src-*.env")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestClone_CopiesAllKeys(t *testing.T) {
	src := writeTempEnvClone(t, "FOO=bar\nBAZ=qux\n")
	dst := filepath.Join(t.TempDir(), "dst.env")

	n, err := Clone(src, dst, DefaultCloneOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 2 {
		t.Errorf("want 2 keys written, got %d", n)
	}

	env, _ := Parse(dst)
	if env["FOO"] != "bar" || env["BAZ"] != "qux" {
		t.Errorf("unexpected env: %v", env)
	}
}

func TestClone_FailsWhenDestExists(t *testing.T) {
	src := writeTempEnvClone(t, "FOO=bar\n")
	dst := writeTempEnvClone(t, "EXISTING=1\n")

	_, err := Clone(src, dst, DefaultCloneOptions())
	if err == nil {
		t.Fatal("expected error when destination exists")
	}
}

func TestClone_OverwriteReplacesFile(t *testing.T) {
	src := writeTempEnvClone(t, "KEY=new\n")
	dst := writeTempEnvClone(t, "KEY=old\n")

	opts := DefaultCloneOptions()
	opts.Overwrite = true
	_, err := Clone(src, dst, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	env, _ := Parse(dst)
	if env["KEY"] != "new" {
		t.Errorf("want KEY=new, got %q", env["KEY"])
	}
}

func TestClone_FilterKeys(t *testing.T) {
	src := writeTempEnvClone(t, "A=1\nB=2\nC=3\n")
	dst := filepath.Join(t.TempDir(), "dst.env")

	opts := DefaultCloneOptions()
	opts.FilterKeys = []string{"A", "C"}
	n, err := Clone(src, dst, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 2 {
		t.Errorf("want 2 keys, got %d", n)
	}

	env, _ := Parse(dst)
	if _, ok := env["B"]; ok {
		t.Error("B should have been filtered out")
	}
}

func TestClone_TransformKey(t *testing.T) {
	src := writeTempEnvClone(t, "foo=bar\n")
	dst := filepath.Join(t.TempDir(), "dst.env")

	opts := DefaultCloneOptions()
	opts.TransformKey = func(k string) string { return strings.ToUpper(k) }
	_, err := Clone(src, dst, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	env, _ := Parse(dst)
	if env["FOO"] != "bar" {
		t.Errorf("want FOO=bar, got %v", env)
	}
}

func TestCloneWithPrefix_PrependsPrefix(t *testing.T) {
	src := writeTempEnvClone(t, "NAME=alice\n")
	dst := filepath.Join(t.TempDir(), "dst.env")

	_, err := CloneWithPrefix(src, dst, "APP_", DefaultCloneOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	env, _ := Parse(dst)
	if env["APP_NAME"] != "alice" {
		t.Errorf("want APP_NAME=alice, got %v", env)
	}
}
