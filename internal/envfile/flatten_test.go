package envfile

import (
	"testing"
)

func TestFlatten_DefaultSeparatorAndUppercase(t *testing.T) {
	input := map[string]string{
		"db.host": "localhost",
		"db.port": "5432",
	}
	opts := DefaultFlattenOptions()
	out, err := Flatten(input, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %q", out["DB_HOST"])
	}
	if out["DB_PORT"] != "5432" {
		t.Errorf("expected DB_PORT=5432, got %q", out["DB_PORT"])
	}
}

func TestFlatten_SlashSeparator(t *testing.T) {
	input := map[string]string{
		"app/name": "vaultpipe",
		"app/version": "1.0",
	}
	opts := DefaultFlattenOptions()
	out, err := Flatten(input, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["APP_NAME"] != "vaultpipe" {
		t.Errorf("expected APP_NAME=vaultpipe, got %q", out["APP_NAME"])
	}
	if out["APP_VERSION"] != "1.0" {
		t.Errorf("expected APP_VERSION=1.0, got %q", out["APP_VERSION"])
	}
}

func TestFlatten_DoubleUnderscoreSeparator(t *testing.T) {
	input := map[string]string{
		"DB__HOST": "db.internal",
	}
	opts := DefaultFlattenOptions()
	out, err := Flatten(input, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DB_HOST"] != "db.internal" {
		t.Errorf("expected DB_HOST=db.internal, got %q", out["DB_HOST"])
	}
}

func TestFlatten_WithPrefix(t *testing.T) {
	input := map[string]string{
		"host": "localhost",
	}
	opts := DefaultFlattenOptions()
	opts.Prefix = "MYAPP"
	out, err := Flatten(input, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["MYAPP_HOST"] != "localhost" {
		t.Errorf("expected MYAPP_HOST=localhost, got %q", out["MYAPP_HOST"])
	}
}

func TestFlatten_LowercaseWhenDisabled(t *testing.T) {
	input := map[string]string{
		"db.host": "localhost",
	}
	opts := DefaultFlattenOptions()
	opts.Uppercase = false
	out, err := Flatten(input, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["db_host"] != "localhost" {
		t.Errorf("expected db_host=localhost, got %q", out["db_host"])
	}
}

func TestFlatten_EmptyKeyReturnsError(t *testing.T) {
	input := map[string]string{
		"": "value",
	}
	_, err := Flatten(input, DefaultFlattenOptions())
	if err == nil {
		t.Error("expected error for empty key, got nil")
	}
}

func TestFlatten_EmptyMapReturnsEmptyMap(t *testing.T) {
	out, err := Flatten(map[string]string{}, DefaultFlattenOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 0 {
		t.Errorf("expected empty map, got %v", out)
	}
}
