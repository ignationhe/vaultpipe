package envfile

import (
	"os"
	"testing"
)

func TestInterpolate_BasicSubstitution(t *testing.T) {
	env := map[string]string{
		"HOST": "localhost",
		"DSN":  "postgres://${HOST}/db",
	}
	opts := DefaultInterpolateOptions()
	opts.FallbackToEnv = false
	result, err := Interpolate(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["DSN"] != "postgres://localhost/db" {
		t.Errorf("expected substituted DSN, got %q", result["DSN"])
	}
}

func TestInterpolate_DollarSyntax(t *testing.T) {
	env := map[string]string{
		"NAME": "world",
		"MSG":  "hello $NAME",
	}
	opts := DefaultInterpolateOptions()
	opts.FallbackToEnv = false
	result, err := Interpolate(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["MSG"] != "hello world" {
		t.Errorf("got %q", result["MSG"])
	}
}

func TestInterpolate_FallbackToOSEnv(t *testing.T) {
	os.Setenv("_TEST_IVAR", "fromOS")
	defer os.Unsetenv("_TEST_IVAR")

	env := map[string]string{
		"VAL": "${_TEST_IVAR}",
	}
	opts := DefaultInterpolateOptions()
	opts.FallbackToEnv = true
	result, err := Interpolate(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["VAL"] != "fromOS" {
		t.Errorf("expected fromOS, got %q", result["VAL"])
	}
}

func TestInterpolate_MissingKeyUsesDefault(t *testing.T) {
	env := map[string]string{
		"VAL": "${MISSING}",
	}
	opts := DefaultInterpolateOptions()
	opts.FallbackToEnv = false
	opts.DefaultValue = "EMPTY"
	result, err := Interpolate(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["VAL"] != "EMPTY" {
		t.Errorf("expected EMPTY, got %q", result["VAL"])
	}
}

func TestInterpolate_ErrorOnMissing(t *testing.T) {
	env := map[string]string{
		"VAL": "${UNDEFINED}",
	}
	opts := DefaultInterpolateOptions()
	opts.FallbackToEnv = false
	opts.ErrorOnMissing = true
	_, err := Interpolate(env, opts)
	if err == nil {
		t.Fatal("expected error for missing variable")
	}
}

func TestInterpolate_NoReferencesUnchanged(t *testing.T) {
	env := map[string]string{
		"PLAIN": "just-a-value",
	}
	opts := DefaultInterpolateOptions()
	result, err := Interpolate(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["PLAIN"] != "just-a-value" {
		t.Errorf("expected unchanged value, got %q", result["PLAIN"])
	}
}
