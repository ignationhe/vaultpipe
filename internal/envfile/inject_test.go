package envfile

import (
	"os"
	"testing"
)

func TestInject_SetsEnvVars(t *testing.T) {
	t.Setenv("INJECT_FOO", "") // ensure clean state
	env := map[string]string{"INJECT_FOO": "bar"}
	opts := DefaultInjectOptions()

	results, err := Inject(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if got := os.Getenv("INJECT_FOO"); got != "bar" {
		t.Errorf("expected INJECT_FOO=bar, got %q", got)
	}
}

func TestInject_SkipsExistingWithoutOverwrite(t *testing.T) {
	t.Setenv("INJECT_EXISTING", "original")
	env := map[string]string{"INJECT_EXISTING": "new"}
	opts := DefaultInjectOptions()

	results, err := Inject(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !results[0].Skipped {
		t.Error("expected key to be skipped")
	}
	if got := os.Getenv("INJECT_EXISTING"); got != "original" {
		t.Errorf("expected original value, got %q", got)
	}
}

func TestInject_OverwriteReplaces(t *testing.T) {
	t.Setenv("INJECT_OW", "old")
	env := map[string]string{"INJECT_OW": "new"}
	opts := DefaultInjectOptions()
	opts.Overwrite = true

	_, err := Inject(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := os.Getenv("INJECT_OW"); got != "new" {
		t.Errorf("expected new value, got %q", got)
	}
}

func TestInject_DryRunDoesNotSet(t *testing.T) {
	os.Unsetenv("INJECT_DRY")
	env := map[string]string{"INJECT_DRY": "value"}
	opts := DefaultInjectOptions()
	opts.DryRun = true

	results, err := Inject(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || results[0].DryRun != true {
		t.Error("expected dry run record")
	}
	if got := os.Getenv("INJECT_DRY"); got != "" {
		t.Errorf("expected empty env, got %q", got)
	}
}

func TestInject_PrefixFiltersAndStrips(t *testing.T) {
	os.Unsetenv("BAR")
	env := map[string]string{
		"APP_BAR": "baz",
		"OTHER":   "nope",
	}
	opts := DefaultInjectOptions()
	opts.Prefix = "APP_"
	opts.StripPrefix = true

	results, err := Inject(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Key != "BAR" {
		t.Errorf("expected key BAR, got %q", results[0].Key)
	}
	if got := os.Getenv("BAR"); got != "baz" {
		t.Errorf("expected BAR=baz, got %q", got)
	}
}

func TestInject_EmptyKeyAfterStripReturnsError(t *testing.T) {
	env := map[string]string{"APP_": "val"}
	opts := DefaultInjectOptions()
	opts.Prefix = "APP_"
	opts.StripPrefix = true

	_, err := Inject(env, opts)
	if err == nil {
		t.Error("expected error for empty key after strip")
	}
}
