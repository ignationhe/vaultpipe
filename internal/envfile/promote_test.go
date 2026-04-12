package envfile

import (
	"testing"
)

func TestPromote_AllKeysWhenNoneSpecified(t *testing.T) {
	src := map[string]string{"A": "1", "B": "2"}
	dst := map[string]string{}
	out, res, err := Promote(src, dst, DefaultPromoteOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["A"] != "1" || out["B"] != "2" {
		t.Errorf("expected all keys promoted, got %v", out)
	}
	if len(res.Promoted) != 2 {
		t.Errorf("expected 2 promoted, got %d", len(res.Promoted))
	}
}

func TestPromote_SkipsExistingWhenNoOverwrite(t *testing.T) {
	src := map[string]string{"A": "new"}
	dst := map[string]string{"A": "old"}
	opts := DefaultPromoteOptions()
	out, res, err := Promote(src, dst, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["A"] != "old" {
		t.Errorf("expected key to be preserved, got %q", out["A"])
	}
	if len(res.Skipped) != 1 || res.Skipped[0] != "A" {
		t.Errorf("expected A in skipped, got %v", res.Skipped)
	}
}

func TestPromote_OverwritesWhenEnabled(t *testing.T) {
	src := map[string]string{"A": "new"}
	dst := map[string]string{"A": "old"}
	opts := DefaultPromoteOptions()
	opts.Overwrite = true
	out, res, err := Promote(src, dst, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["A"] != "new" {
		t.Errorf("expected key to be overwritten, got %q", out["A"])
	}
	if len(res.Overwritten) != 1 {
		t.Errorf("expected 1 overwritten, got %v", res.Overwritten)
	}
}

func TestPromote_LimitsToSpecifiedKeys(t *testing.T) {
	src := map[string]string{"A": "1", "B": "2", "C": "3"}
	dst := map[string]string{}
	opts := DefaultPromoteOptions()
	opts.Keys = []string{"A", "C"}
	out, res, err := Promote(src, dst, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["B"]; ok {
		t.Errorf("expected B to be excluded from promotion")
	}
	if len(res.Promoted) != 2 {
		t.Errorf("expected 2 promoted, got %d", len(res.Promoted))
	}
}

func TestPromote_NilSourceReturnsError(t *testing.T) {
	_, _, err := Promote(nil, map[string]string{}, DefaultPromoteOptions())
	if err == nil {
		t.Error("expected error for nil source")
	}
}

func TestPromote_NilDestinationCreatesNew(t *testing.T) {
	src := map[string]string{"X": "42"}
	out, _, err := Promote(src, nil, DefaultPromoteOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["X"] != "42" {
		t.Errorf("expected X=42, got %q", out["X"])
	}
}

func TestPromote_SkipsKeysMissingFromSource(t *testing.T) {
	src := map[string]string{"A": "1"}
	dst := map[string]string{}
	opts := DefaultPromoteOptions()
	opts.Keys = []string{"A", "MISSING"}
	out, res, err := Promote(src, dst, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 1 {
		t.Errorf("expected 1 key in output, got %d", len(out))
	}
	if len(res.Promoted) != 1 {
		t.Errorf("expected 1 promoted, got %v", res.Promoted)
	}
}
