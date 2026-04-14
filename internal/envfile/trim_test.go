package envfile

import (
	"testing"
)

func TestTrim_DefaultTrimsAllWhitespace(t *testing.T) {
	env := map[string]string{
		"FOO": "  hello  ",
		"BAR": "\tworld\t",
	}
	opts := DefaultTrimOptions()
	out, err := Trim(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["FOO"] != "hello" {
		t.Errorf("FOO: got %q, want %q", out["FOO"], "hello")
	}
	if out["BAR"] != "world" {
		t.Errorf("BAR: got %q, want %q", out["BAR"], "world")
	}
}

func TestTrim_OnlyTrimLeft(t *testing.T) {
	env := map[string]string{"KEY": "  value  "}
	opts := TrimOptions{TrimLeft: true, TrimRight: false}
	out, _ := Trim(env, opts)
	if out["KEY"] != "value  " {
		t.Errorf("got %q, want %q", out["KEY"], "value  ")
	}
}

func TestTrim_OnlyTrimRight(t *testing.T) {
	env := map[string]string{"KEY": "  value  "}
	opts := TrimOptions{TrimLeft: false, TrimRight: true}
	out, _ := Trim(env, opts)
	if out["KEY"] != "  value" {
		t.Errorf("got %q, want %q", out["KEY"], "  value")
	}
}

func TestTrim_LimitsToSpecifiedKeys(t *testing.T) {
	env := map[string]string{
		"TRIMMED":    "  yes  ",
		"UNTOUCHED":  "  no  ",
	}
	opts := TrimOptions{Keys: []string{"TRIMMED"}, TrimLeft: true, TrimRight: true}
	out, _ := Trim(env, opts)
	if out["TRIMMED"] != "yes" {
		t.Errorf("TRIMMED: got %q, want %q", out["TRIMMED"], "yes")
	}
	if out["UNTOUCHED"] != "  no  " {
		t.Errorf("UNTOUCHED: got %q, want %q", out["UNTOUCHED"], "  no  ")
	}
}

func TestTrim_CustomCutset(t *testing.T) {
	env := map[string]string{"KEY": "***secret***"}
	opts := TrimOptions{TrimLeft: true, TrimRight: true, Cutset: "*"}
	out, _ := Trim(env, opts)
	if out["KEY"] != "secret" {
		t.Errorf("got %q, want %q", out["KEY"], "secret")
	}
}

func TestTrim_NoOpWhenBothFalse(t *testing.T) {
	env := map[string]string{"KEY": "  value  "}
	opts := TrimOptions{TrimLeft: false, TrimRight: false}
	out, _ := Trim(env, opts)
	if out["KEY"] != "  value  " {
		t.Errorf("got %q, want %q", out["KEY"], "  value  ")
	}
}

func TestTrim_EmptyEnvReturnsEmpty(t *testing.T) {
	out, err := Trim(map[string]string{}, DefaultTrimOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 0 {
		t.Errorf("expected empty map, got %v", out)
	}
}
