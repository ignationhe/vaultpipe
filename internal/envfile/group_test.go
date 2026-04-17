package envfile

import (
	"testing"
)

var baseEnvForGroup = map[string]string{
	"DB_HOST":     "localhost",
	"DB_PORT":     "5432",
	"APP_NAME":    "vaultpipe",
	"APP_VERSION": "1.0",
	"STANDALONE":  "yes",
}

func TestGroup_ByUnderscore(t *testing.T) {
	result := Group(baseEnvForGroup, DefaultGroupOptions())

	if result["DB"]["HOST"] != "localhost" {
		t.Errorf("expected DB/HOST=localhost, got %q", result["DB"]["HOST"])
	}
	if result["DB"]["PORT"] != "5432" {
		t.Errorf("expected DB/PORT=5432, got %q", result["DB"]["PORT"])
	}
	if result["APP"]["NAME"] != "vaultpipe" {
		t.Errorf("expected APP/NAME=vaultpipe")
	}
}

func TestGroup_StandaloneKeyInCatchAll(t *testing.T) {
	result := Group(baseEnvForGroup, DefaultGroupOptions())
	if result["_"]["STANDALONE"] != "yes" {
		t.Errorf("expected STANDALONE in '_' bucket")
	}
}

func TestGroup_KeepPrefix(t *testing.T) {
	opts := DefaultGroupOptions()
	opts.KeepPrefix = true
	result := Group(baseEnvForGroup, opts)

	if result["DB"]["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST key preserved with KeepPrefix")
	}
}

func TestGroup_CustomSeparator(t *testing.T) {
	env := map[string]string{
		"aws.region": "us-east-1",
		"aws.key":    "AKID",
		"plain":      "value",
	}
	opts := GroupOptions{Separator: ".", KeepPrefix: false}
	result := Group(env, opts)

	if result["aws"]["region"] != "us-east-1" {
		t.Errorf("expected aws/region=us-east-1")
	}
	if result["_"]["plain"] != "value" {
		t.Errorf("expected plain in catch-all")
	}
}

func TestGroup_EmptyMapReturnsEmpty(t *testing.T) {
	result := Group(map[string]string{}, DefaultGroupOptions())
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d groups", len(result))
	}
}

func TestGroup_DefaultSeparatorFallback(t *testing.T) {
	// passing zero-value opts should still work
	result := Group(map[string]string{"A_B": "1"}, GroupOptions{})
	if result["A"]["B"] != "1" {
		t.Errorf("expected A/B=1 with zero-value opts")
	}
}
