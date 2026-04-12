package envfile

import (
	"errors"
	"strings"
	"testing"
)

func TestTransform_ExactKeyRule(t *testing.T) {
	src := map[string]string{"FOO": "hello", "BAR": "world"}
	opts := DefaultTransformOptions()
	opts.Rules["FOO"] = UppercaseValues()

	out, err := Transform(src, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["FOO"] != "HELLO" {
		t.Errorf("expected HELLO, got %q", out["FOO"])
	}
	if out["BAR"] != "world" {
		t.Errorf("BAR should be unchanged, got %q", out["BAR"])
	}
}

func TestTransform_WildcardRule(t *testing.T) {
	src := map[string]string{"A": "  trimme  ", "B": " spaces "}
	opts := DefaultTransformOptions()
	opts.Rules["*"] = TrimSpaceValues()

	out, err := Transform(src, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for k, v := range out {
		if strings.TrimSpace(v) != v {
			t.Errorf("key %q value %q still has spaces", k, v)
		}
	}
}

func TestTransform_ExactTakesPriorityOverWildcard(t *testing.T) {
	src := map[string]string{"KEY": "value"}
	opts := DefaultTransformOptions()
	opts.Rules["*"] = UppercaseValues()
	opts.Rules["KEY"] = LowercaseValues()

	out, err := Transform(src, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["KEY"] != "value" {
		t.Errorf("expected lowercase 'value', got %q", out["KEY"])
	}
}

func TestTransform_ErrorStopsProcessing(t *testing.T) {
	src := map[string]string{"X": "val"}
	opts := DefaultTransformOptions()
	opts.Rules["X"] = func(k, v string) (string, error) {
		return "", errors.New("boom")
	}

	_, err := Transform(src, opts)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestTransform_SkipErrorsContinues(t *testing.T) {
	src := map[string]string{"X": "val", "Y": "ok"}
	opts := DefaultTransformOptions()
	opts.SkipErrors = true
	opts.Rules["X"] = func(k, v string) (string, error) {
		return "", errors.New("boom")
	}
	opts.Rules["Y"] = UppercaseValues()

	out, err := Transform(src, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["X"] != "val" {
		t.Errorf("X should be unchanged on skip, got %q", out["X"])
	}
	if out["Y"] != "OK" {
		t.Errorf("Y should be uppercased, got %q", out["Y"])
	}
}

func TestTransform_NoRulesReturnsClone(t *testing.T) {
	src := map[string]string{"A": "1", "B": "2"}
	opts := DefaultTransformOptions()

	out, err := Transform(src, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != len(src) {
		t.Errorf("expected %d keys, got %d", len(src), len(out))
	}
	src["A"] = "mutated"
	if out["A"] == "mutated" {
		t.Error("Transform should return a copy, not share src")
	}
}
