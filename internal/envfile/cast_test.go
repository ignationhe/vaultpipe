package envfile

import (
	"testing"
)

func TestCast_StringDefault(t *testing.T) {
	env := map[string]string{"NAME": "alice"}
	results, err := Cast(env, DefaultCastOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Cast != "alice" {
		t.Errorf("expected 'alice', got %v", results[0].Cast)
	}
	if results[0].CastType != CastString {
		t.Errorf("expected CastString, got %v", results[0].CastType)
	}
}

func TestCast_IntRule(t *testing.T) {
	env := map[string]string{"PORT": "8080"}
	opts := DefaultCastOptions()
	opts.Rules = []CastRule{{Key: "PORT", Type: CastInt}}

	results, err := Cast(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].Cast.(int) != 8080 {
		t.Errorf("expected 8080, got %v", results[0].Cast)
	}
}

func TestCast_FloatRule(t *testing.T) {
	env := map[string]string{"RATIO": "3.14"}
	opts := DefaultCastOptions()
	opts.Rules = []CastRule{{Key: "RATIO", Type: CastFloat}}

	results, err := Cast(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	f, ok := results[0].Cast.(float64)
	if !ok || f < 3.13 || f > 3.15 {
		t.Errorf("expected ~3.14, got %v", results[0].Cast)
	}
}

func TestCast_BoolRule(t *testing.T) {
	env := map[string]string{"DEBUG": "true"}
	opts := DefaultCastOptions()
	opts.Rules = []CastRule{{Key: "DEBUG", Type: CastBool}}

	results, err := Cast(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].Cast.(bool) != true {
		t.Errorf("expected true, got %v", results[0].Cast)
	}
}

func TestCast_InvalidIntReturnsError(t *testing.T) {
	env := map[string]string{"PORT": "not-a-number"}
	opts := DefaultCastOptions()
	opts.Rules = []CastRule{{Key: "PORT", Type: CastInt}}

	_, err := Cast(env, opts)
	if err == nil {
		t.Fatal("expected error for invalid int cast")
	}
}

func TestCast_SkipUnknown(t *testing.T) {
	env := map[string]string{"KNOWN": "1", "UNKNOWN": "x"}
	opts := DefaultCastOptions()
	opts.SkipUnknown = true
	opts.Rules = []CastRule{{Key: "KNOWN", Type: CastInt}}

	results, err := Cast(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("expected 1 result (UNKNOWN skipped), got %d", len(results))
	}
	if results[0].Key != "KNOWN" {
		t.Errorf("expected key KNOWN, got %s", results[0].Key)
	}
}

func TestCast_RawPreserved(t *testing.T) {
	env := map[string]string{"COUNT": "42"}
	opts := DefaultCastOptions()
	opts.Rules = []CastRule{{Key: "COUNT", Type: CastInt}}

	results, err := Cast(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].Raw != "42" {
		t.Errorf("expected raw '42', got %q", results[0].Raw)
	}
}
