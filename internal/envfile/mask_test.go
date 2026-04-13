package envfile

import (
	"testing"
)

func baseEnvForMask() map[string]string {
	return map[string]string{
		"DB_PASSWORD":  "supersecret",
		"API_TOKEN":    "tok_abc123",
		"APP_NAME":     "vaultpipe",
		"PRIVATE_KEY":  "-----BEGIN RSA",
		"LOG_LEVEL":    "debug",
	}
}

func TestMask_DefaultPatternsMaskSensitiveKeys(t *testing.T) {
	env := baseEnvForMask()
	opts := DefaultMaskOptions()
	result := Mask(env, opts)

	if result["DB_PASSWORD"] != "****" {
		t.Errorf("expected DB_PASSWORD masked, got %q", result["DB_PASSWORD"])
	}
	if result["API_TOKEN"] != "****" {
		t.Errorf("expected API_TOKEN masked, got %q", result["API_TOKEN"])
	}
	if result["PRIVATE_KEY"] != "****" {
		t.Errorf("expected PRIVATE_KEY masked, got %q", result["PRIVATE_KEY"])
	}
}

func TestMask_NonSensitiveKeysPreserved(t *testing.T) {
	env := baseEnvForMask()
	opts := DefaultMaskOptions()
	result := Mask(env, opts)

	if result["APP_NAME"] != "vaultpipe" {
		t.Errorf("expected APP_NAME preserved, got %q", result["APP_NAME"])
	}
	if result["LOG_LEVEL"] != "debug" {
		t.Errorf("expected LOG_LEVEL preserved, got %q", result["LOG_LEVEL"])
	}
}

func TestMask_CustomPlaceholder(t *testing.T) {
	env := map[string]string{"MY_SECRET": "abc"}
	opts := DefaultMaskOptions()
	opts.Placeholder = "[REDACTED]"
	result := Mask(env, opts)

	if result["MY_SECRET"] != "[REDACTED]" {
		t.Errorf("expected [REDACTED], got %q", result["MY_SECRET"])
	}
}

func TestMask_PartialReveal(t *testing.T) {
	env := map[string]string{"API_TOKEN": "tok_abc123"}
	opts := DefaultMaskOptions()
	opts.MaskPartial = 3
	result := Mask(env, opts)

	if result["API_TOKEN"] != "tok****" {
		t.Errorf("expected tok****, got %q", result["API_TOKEN"])
	}
}

func TestMask_CaseSensitiveDoesNotMatchLowercase(t *testing.T) {
	env := map[string]string{"db_password": "secret"}
	opts := DefaultMaskOptions()
	opts.CaseSensitive = true
	result := Mask(env, opts)

	if result["db_password"] != "secret" {
		t.Errorf("expected value preserved with case-sensitive match, got %q", result["db_password"])
	}
}

func TestMask_EmptyEnvReturnsEmpty(t *testing.T) {
	result := Mask(map[string]string{}, DefaultMaskOptions())
	if len(result) != 0 {
		t.Errorf("expected empty map, got %d entries", len(result))
	}
}
