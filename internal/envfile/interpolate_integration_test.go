package envfile

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInterpolate_ThenWrite_RoundTrip(t *testing.T) {
	env := map[string]string{
		"BASE_URL": "https://example.com",
		"API_URL":  "${BASE_URL}/api",
		"HEALTH":   "${BASE_URL}/health",
	}
	opts := DefaultInterpolateOptions()
	opts.FallbackToEnv = false

	resolved, err := Interpolate(env, opts)
	if err != nil {
		t.Fatalf("interpolate: %v", err)
	}

	tmp := filepath.Join(t.TempDir(), ".env")
	if err := Write(tmp, resolved); err != nil {
		t.Fatalf("write: %v", err)
	}

	parsed, err := Parse(tmp)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	if parsed["API_URL"] != "https://example.com/api" {
		t.Errorf("API_URL: got %q", parsed["API_URL"])
	}
	if parsed["HEALTH"] != "https://example.com/health" {
		t.Errorf("HEALTH: got %q", parsed["HEALTH"])
	}
}

func TestInterpolate_WithOSEnvFallback_Integration(t *testing.T) {
	os.Setenv("_VAULTPIPE_TEST_REGION", "us-east-1")
	defer os.Unsetenv("_VAULTPIPE_TEST_REGION")

	env := map[string]string{
		"BUCKET": "my-bucket-${_VAULTPIPE_TEST_REGION}",
	}
	opts := DefaultInterpolateOptions()

	result, err := Interpolate(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["BUCKET"] != "my-bucket-us-east-1" {
		t.Errorf("expected expanded bucket name, got %q", result["BUCKET"])
	}
}
