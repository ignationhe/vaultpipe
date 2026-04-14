package envfile

import (
	"testing"
)

func TestNormalize_UppercaseKeys(t *testing.T) {
	env := map[string]string{"db_host": "localhost", "api_key": "secret"}
	opts := DefaultNormalizeOptions()
	opts.ReplaceHyphens = false

	result, err := Normalize(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %q", result["DB_HOST"])
	}
	if result["API_KEY"] != "secret" {
		t.Errorf("expected API_KEY=secret, got %q", result["API_KEY"])
	}
}

func TestNormalize_ReplaceHyphens(t *testing.T) {
	env := map[string]string{"my-key": "value", "another-key": "val2"}
	opts := DefaultNormalizeOptions()

	result, err := Normalize(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := result["MY_KEY"]; !ok {
		t.Errorf("expected MY_KEY to exist in result")
	}
	if _, ok := result["ANOTHER_KEY"]; !ok {
		t.Errorf("expected ANOTHER_KEY to exist in result")
	}
}

func TestNormalize_TrimValues(t *testing.T) {
	env := map[string]string{"KEY": "  spaced  "}
	opts := DefaultNormalizeOptions()

	result, err := Normalize(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["KEY"] != "spaced" {
		t.Errorf("expected trimmed value, got %q", result["KEY"])
	}
}

func TestNormalize_CollapseUnderscores(t *testing.T) {
	env := map[string]string{"MY__KEY": "val", "A___B": "val2"}
	opts := DefaultNormalizeOptions()
	opts.CollapseUnderscores = true

	result, err := Normalize(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := result["MY_KEY"]; !ok {
		t.Errorf("expected MY_KEY, got keys: %v", result)
	}
	if _, ok := result["A_B"]; !ok {
		t.Errorf("expected A_B, got keys: %v", result)
	}
}

func TestNormalize_NoOpWhenAllDisabled(t *testing.T) {
	env := map[string]string{"my-key": "  val  "}
	opts := NormalizeOptions{
		UppercaseKeys:       false,
		TrimValues:          false,
		ReplaceHyphens:      false,
		CollapseUnderscores: false,
	}

	result, err := Normalize(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["my-key"] != "  val  " {
		t.Errorf("expected unchanged key/value, got key=%q val=%q", "my-key", result["my-key"])
	}
}

func TestNormalize_EmptyMapReturnsEmpty(t *testing.T) {
	result, err := Normalize(map[string]string{}, DefaultNormalizeOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty result, got %v", result)
	}
}
