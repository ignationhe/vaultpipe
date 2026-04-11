package envfile

import (
	"regexp"
	"testing"
)

func TestValidateSchema_AllRequiredPresent(t *testing.T) {
	env := map[string]string{"HOST": "localhost", "PORT": "8080"}
	schema := Schema{
		Fields: []SchemaField{
			{Key: "HOST", Required: true},
			{Key: "PORT", Required: true},
		},
	}
	violations := ValidateSchema(env, schema)
	if len(violations) != 0 {
		t.Fatalf("expected no violations, got %v", violations)
	}
}

func TestValidateSchema_MissingRequiredKey(t *testing.T) {
	env := map[string]string{"HOST": "localhost"}
	schema := Schema{
		Fields: []SchemaField{
			{Key: "PORT", Required: true},
		},
	}
	violations := ValidateSchema(env, schema)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Key != "PORT" {
		t.Errorf("expected violation for PORT, got %s", violations[0].Key)
	}
}

func TestValidateSchema_PatternMismatch(t *testing.T) {
	env := map[string]string{"PORT": "not-a-number"}
	schema := Schema{
		Fields: []SchemaField{
			{Key: "PORT", Required: true, Pattern: regexp.MustCompile(`^\d+$`)},
		},
	}
	violations := ValidateSchema(env, schema)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
}

func TestValidateSchema_PatternMatch(t *testing.T) {
	env := map[string]string{"PORT": "3000"}
	schema := Schema{
		Fields: []SchemaField{
			{Key: "PORT", Required: true, Pattern: regexp.MustCompile(`^\d+$`)},
		},
	}
	violations := ValidateSchema(env, schema)
	if len(violations) != 0 {
		t.Fatalf("expected no violations, got %v", violations)
	}
}

func TestValidateSchema_OptionalMissingKeyNoViolation(t *testing.T) {
	env := map[string]string{}
	schema := Schema{
		Fields: []SchemaField{
			{Key: "DEBUG", Required: false},
		},
	}
	violations := ValidateSchema(env, schema)
	if len(violations) != 0 {
		t.Fatalf("expected no violations for optional missing key, got %v", violations)
	}
}

func TestValidateSchema_MultipleViolations(t *testing.T) {
	env := map[string]string{"PORT": "abc"}
	schema := Schema{
		Fields: []SchemaField{
			{Key: "HOST", Required: true},
			{Key: "PORT", Required: true, Pattern: regexp.MustCompile(`^\d+$`)},
		},
	}
	violations := ValidateSchema(env, schema)
	if len(violations) != 2 {
		t.Fatalf("expected 2 violations, got %d", len(violations))
	}
}
