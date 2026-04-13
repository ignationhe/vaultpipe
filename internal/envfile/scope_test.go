package envfile

import (
	"testing"
)

func baseEnvForScope() map[string]string {
	return map[string]string{
		"APP_NAME":         "myapp",
		"DB_URL":           "base-db",
		"STAGING__DB_URL":  "staging-db",
		"PROD__DB_URL":     "prod-db",
		"PROD__SECRET_KEY": "prod-secret",
	}
}

func TestScope_ResolvesTargetScopeOverridesBase(t *testing.T) {
	opts := DefaultScopeOptions()
	opts.Scopes = []string{"staging", "prod"}
	opts.TargetScope = "prod"

	out, err := Scope(baseEnvForScope(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DB_URL"] != "prod-db" {
		t.Errorf("expected prod-db, got %q", out["DB_URL"])
	}
	if out["APP_NAME"] != "myapp" {
		t.Errorf("expected global APP_NAME to survive, got %q", out["APP_NAME"])
	}
}

func TestScope_StagingDoesNotIncludeProdKeys(t *testing.T) {
	opts := DefaultScopeOptions()
	opts.Scopes = []string{"staging", "prod"}
	opts.TargetScope = "staging"

	out, err := Scope(baseEnvForScope(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DB_URL"] != "staging-db" {
		t.Errorf("expected staging-db, got %q", out["DB_URL"])
	}
	if _, ok := out["SECRET_KEY"]; ok {
		t.Error("staging should not include PROD__SECRET_KEY")
	}
}

func TestScope_GlobalKeysIncludedWhenNoOverride(t *testing.T) {
	opts := DefaultScopeOptions()
	opts.Scopes = []string{"dev", "prod"}
	opts.TargetScope = "dev"

	out, err := Scope(baseEnvForScope(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["APP_NAME"] != "myapp" {
		t.Errorf("expected global APP_NAME, got %q", out["APP_NAME"])
	}
	// base DB_URL should survive since dev has no override
	if out["DB_URL"] != "base-db" {
		t.Errorf("expected base-db, got %q", out["DB_URL"])
	}
}

func TestScope_StripPrefixFalseKeepsOriginalKey(t *testing.T) {
	opts := DefaultScopeOptions()
	opts.Scopes = []string{"prod"}
	opts.TargetScope = "prod"
	opts.StripPrefix = false

	out, err := Scope(baseEnvForScope(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["PROD__DB_URL"]; !ok {
		t.Error("expected PROD__DB_URL key to be present when StripPrefix=false")
	}
}

func TestScope_EmptyTargetReturnsOnlyGlobals(t *testing.T) {
	opts := DefaultScopeOptions()
	opts.Scopes = []string{"staging", "prod"}
	opts.TargetScope = ""

	out, err := Scope(baseEnvForScope(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["APP_NAME"] != "myapp" {
		t.Errorf("expected APP_NAME, got %q", out["APP_NAME"])
	}
	if _, ok := out["DB_URL"]; !ok {
		// base DB_URL is unscoped so it should be present
		t.Error("expected base DB_URL to be present")
	}
}
