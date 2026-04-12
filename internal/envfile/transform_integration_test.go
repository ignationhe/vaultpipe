package envfile

import (
	"os"
	"testing"
)

func TestTransform_ThenWrite_RoundTrip(t *testing.T) {
	src := map[string]string{
		"DB_HOST": "  localhost  ",
		"DB_PORT": "5432",
		"SECRET":  "mysecret",
	}

	opts := DefaultTransformOptions()
	opts.Rules["DB_HOST"] = TrimSpaceValues()
	opts.Rules["SECRET"] = UppercaseValues()

	transformed, err := Transform(src, opts)
	if err != nil {
		t.Fatalf("Transform: %v", err)
	}

	tmp, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatalf("CreateTemp: %v", err)
	}
	tmp.Close()

	if err := Write(transformed, tmp.Name()); err != nil {
		t.Fatalf("Write: %v", err)
	}

	parsed, err := Parse(tmp.Name())
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}

	if parsed["DB_HOST"] != "localhost" {
		t.Errorf("DB_HOST: expected 'localhost', got %q", parsed["DB_HOST"])
	}
	if parsed["SECRET"] != "MYSECRET" {
		t.Errorf("SECRET: expected 'MYSECRET', got %q", parsed["SECRET"])
	}
	if parsed["DB_PORT"] != "5432" {
		t.Errorf("DB_PORT: expected '5432', got %q", parsed["DB_PORT"])
	}
}

func TestTransform_WildcardTrimThenDiff(t *testing.T) {
	old := map[string]string{"FOO": "bar", "BAZ": "qux"}
	src := map[string]string{"FOO": "  bar  ", "BAZ": "  qux  "}

	opts := DefaultTransformOptions()
	opts.Rules["*"] = TrimSpaceValues()

	transformed, err := Transform(src, opts)
	if err != nil {
		t.Fatalf("Transform: %v", err)
	}

	diffs := Diff(old, transformed)
	if HasChanges(diffs) {
		t.Errorf("expected no changes after trim, but got diffs: %+v", diffs)
	}
}

// TestTransform_WriteAndReparse_PreservesAllKeys verifies that every key
// present in the source map survives a full Write → Parse round-trip without
// any keys being silently dropped.
func TestTransform_WriteAndReparse_PreservesAllKeys(t *testing.T) {
	src := map[string]string{
		"ALPHA": "one",
		"BETA":  "two",
		"GAMMA": "three",
	}

	tmp, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatalf("CreateTemp: %v", err)
	}
	tmp.Close()

	if err := Write(src, tmp.Name()); err != nil {
		t.Fatalf("Write: %v", err)
	}

	parsed, err := Parse(tmp.Name())
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}

	for key, want := range src {
		got, ok := parsed[key]
		if !ok {
			t.Errorf("key %q missing after round-trip", key)
			continue
		}
		if got != want {
			t.Errorf("key %q: expected %q, got %q", key, want, got)
		}
	}
}
