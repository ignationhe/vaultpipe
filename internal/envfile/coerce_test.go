package envfile

import (
	"testing"
)

func baseEnvForCoerce() map[string]string {
	return map[string]string{
		"FEATURE_FLAG": "yes",
		"DEBUG":        "1",
		"VERBOSE":      "off",
		"RETRIES":      "3",
		"NAME":         "alice",
	}
}

func TestCoerce_TruthyToTrue(t *testing.T) {
	env := map[string]string{"A": "yes", "B": "1", "C": "on", "D": "enabled"}
	opts := DefaultCoerceOptions()
	out, err := Coerce(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for k, v := range out {
		if v != "true" {
			t.Errorf("key %q: expected \"true\", got %q", k, v)
		}
	}
}

func TestCoerce_FalsyToFalse(t *testing.T) {
	env := map[string]string{"A": "no", "B": "0", "C": "off", "D": "disabled"}
	opts := DefaultCoerceOptions()
	out, err := Coerce(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for k, v := range out {
		if v != "false" {
			t.Errorf("key %q: expected \"false\", got %q", k, v)
		}
	}
}

func TestCoerce_NonBooleanPreserved(t *testing.T) {
	env := baseEnvForCoerce()
	opts := DefaultCoerceOptions()
	out, err := Coerce(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["RETRIES"] != "3" {
		t.Errorf("expected RETRIES=3, got %q", out["RETRIES"])
	}
	if out["NAME"] != "alice" {
		t.Errorf("expected NAME=alice, got %q", out["NAME"])
	}
}

func TestCoerce_ErrorOnUnknown(t *testing.T) {
	env := map[string]string{"PORT": "8080"}
	opts := DefaultCoerceOptions()
	opts.ErrorOnUnknown = true
	_, err := Coerce(env, opts)
	if err == nil {
		t.Fatal("expected error for non-boolean value, got nil")
	}
}

func TestCoerce_LimitsToSpecifiedKeys(t *testing.T) {
	env := map[string]string{"A": "yes", "B": "yes"}
	opts := DefaultCoerceOptions()
	opts.Keys = []string{"A"}
	out, err := Coerce(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["A"] != "true" {
		t.Errorf("expected A=true, got %q", out["A"])
	}
	if out["B"] != "yes" {
		t.Errorf("expected B=yes (unchanged), got %q", out["B"])
	}
}

func TestCoerce_CaseInsensitive(t *testing.T) {
	env := map[string]string{"X": "YES", "Y": "FALSE"}
	opts := DefaultCoerceOptions()
	out, err := Coerce(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["X"] != "true" {
		t.Errorf("expected X=true, got %q", out["X"])
	}
	if out["Y"] != "false" {
		t.Errorf("expected Y=false, got %q", out["Y"])
	}
}

func TestCoerce_TrimSpaceBeforeCoerce(t *testing.T) {
	env := map[string]string{"FLAG": "  yes  "}
	opts := DefaultCoerceOptions()
	out, err := Coerce(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["FLAG"] != "true" {
		t.Errorf("expected FLAG=true after trim, got %q", out["FLAG"])
	}
}
