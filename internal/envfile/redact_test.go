package envfile

import (
	"testing"
)

func TestRedact_DefaultPatternsObfuscateSensitiveKeys(t *testing.T) {
	secrets := map[string]string{
		"DB_PASSWORD": "s3cr3t",
		"API_KEY":     "abc123",
		"AUTH_TOKEN":  "tok_xyz",
		"DB_HOST":     "localhost",
		"APP_SECRET":  "mysecret",
	}

	result := Redact(secrets, DefaultRedactOptions())

	redacted := []string{"DB_PASSWORD", "API_KEY", "AUTH_TOKEN", "APP_SECRET"}
	for _, k := range redacted {
		if result[k] != "***REDACTED***" {
			t.Errorf("expected %s to be redacted, got %q", k, result[k])
		}
	}

	if result["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST to be preserved, got %q", result["DB_HOST"])
	}
}

func TestRedact_ExactKeyMatch(t *testing.T) {
	secrets := map[string]string{
		"MY_CUSTOM_KEY": "sensitive",
		"SAFE_KEY":      "safe",
	}

	opts := RedactOptions{
		Keys:        []string{"MY_CUSTOM_KEY"},
		Patterns:    []string{},
		Placeholder: "[hidden]",
	}

	result := Redact(secrets, opts)

	if result["MY_CUSTOM_KEY"] != "[hidden]" {
		t.Errorf("expected MY_CUSTOM_KEY to be redacted, got %q", result["MY_CUSTOM_KEY"])
	}
	if result["SAFE_KEY"] != "safe" {
		t.Errorf("expected SAFE_KEY to be preserved, got %q", result["SAFE_KEY"])
	}
}

func TestRedact_CaseInsensitiveExactMatch(t *testing.T) {
	secrets := map[string]string{
		"db_password": "hunter2",
	}

	opts := RedactOptions{
		Keys:        []string{"DB_PASSWORD"},
		Patterns:    []string{},
		Placeholder: "***REDACTED***",
	}

	result := Redact(secrets, opts)

	if result["db_password"] != "***REDACTED***" {
		t.Errorf("expected db_password to be redacted, got %q", result["db_password"])
	}
}

func TestRedact_EmptyOptsPreservesAll(t *testing.T) {
	secrets := map[string]string{
		"FOO": "bar",
		"BAZ": "qux",
	}

	opts := RedactOptions{
		Keys:        []string{},
		Patterns:    []string{},
		Placeholder: "***REDACTED***",
	}

	result := Redact(secrets, opts)

	for k, v := range secrets {
		if result[k] != v {
			t.Errorf("expected %s=%q to be preserved, got %q", k, v, result[k])
		}
	}
}

func TestRedact_DefaultPlaceholderUsedWhenEmpty(t *testing.T) {
	secrets := map[string]string{
		"MY_TOKEN": "tok_abc",
	}

	opts := RedactOptions{
		Keys:        []string{"MY_TOKEN"},
		Patterns:    []string{},
		Placeholder: "",
	}

	result := Redact(secrets, opts)

	if result["MY_TOKEN"] != "***REDACTED***" {
		t.Errorf("expected default placeholder, got %q", result["MY_TOKEN"])
	}
}
