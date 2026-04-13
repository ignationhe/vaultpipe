package envfile

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func tempPinPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "pins.json")
}

func TestPin_StoresKeyValue(t *testing.T) {
	env := map[string]string{"DB_PASS": "secret", "API_KEY": "abc123"}
	path := tempPinPath(t)

	pf, err := Pin(env, []string{"DB_PASS"}, path, DefaultPinOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pf.Pins) != 1 {
		t.Fatalf("expected 1 pin, got %d", len(pf.Pins))
	}
	if pf.Pins[0].Key != "DB_PASS" || pf.Pins[0].Value != "secret" {
		t.Errorf("unexpected pin entry: %+v", pf.Pins[0])
	}
}

func TestPin_SkipsMissingKey(t *testing.T) {
	env := map[string]string{"PRESENT": "yes"}
	path := tempPinPath(t)

	pf, err := Pin(env, []string{"MISSING"}, path, DefaultPinOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pf.Pins) != 0 {
		t.Errorf("expected no pins, got %d", len(pf.Pins))
	}
}

func TestPin_NoOverwriteByDefault(t *testing.T) {
	env := map[string]string{"TOKEN": "first"}
	path := tempPinPath(t)

	_, _ = Pin(env, []string{"TOKEN"}, path, DefaultPinOptions())

	env["TOKEN"] = "second"
	pf, err := Pin(env, []string{"TOKEN"}, path, DefaultPinOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pf.Pins[0].Value != "first" {
		t.Errorf("expected original value to be preserved, got %q", pf.Pins[0].Value)
	}
}

func TestPin_OverwriteWhenEnabled(t *testing.T) {
	env := map[string]string{"TOKEN": "first"}
	path := tempPinPath(t)

	_, _ = Pin(env, []string{"TOKEN"}, path, DefaultPinOptions())

	env["TOKEN"] = "second"
	opts := DefaultPinOptions()
	opts.Overwrite = true
	pf, err := Pin(env, []string{"TOKEN"}, path, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pf.Pins[0].Value != "second" {
		t.Errorf("expected updated value, got %q", pf.Pins[0].Value)
	}
	if len(pf.Pins) != 1 {
		t.Errorf("expected exactly 1 pin after overwrite, got %d", len(pf.Pins))
	}
}

func TestPin_PersistsAcrossLoads(t *testing.T) {
	env := map[string]string{"SECRET": "val"}
	path := tempPinPath(t)

	_, _ = Pin(env, []string{"SECRET"}, path, DefaultPinOptions())

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("could not read pin file: %v", err)
	}
	var pf PinFile
	if err := json.Unmarshal(data, &pf); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(pf.Pins) != 1 || pf.Pins[0].Key != "SECRET" {
		t.Errorf("unexpected persisted pins: %+v", pf.Pins)
	}
}

func TestApplyPins_OverlaysValues(t *testing.T) {
	env := map[string]string{"DB_PASS": "live", "OTHER": "keep"}
	path := tempPinPath(t)

	_, _ = Pin(map[string]string{"DB_PASS": "pinned"}, []string{"DB_PASS"}, path, DefaultPinOptions())

	result, err := ApplyPins(env, path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["DB_PASS"] != "pinned" {
		t.Errorf("expected pinned value, got %q", result["DB_PASS"])
	}
	if result["OTHER"] != "keep" {
		t.Errorf("expected OTHER to be preserved, got %q", result["OTHER"])
	}
}

func TestApplyPins_NoPinFile(t *testing.T) {
	env := map[string]string{"KEY": "val"}
	path := tempPinPath(t) // file does not exist

	result, err := ApplyPins(env, path)
	if err != nil {
		t.Fatalf("unexpected error when no pin file: %v", err)
	}
	if result["KEY"] != "val" {
		t.Errorf("expected original value, got %q", result["KEY"])
	}
}
