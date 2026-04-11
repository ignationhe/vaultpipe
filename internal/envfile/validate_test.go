package envfile

import (
	"testing"
)

func TestValidate_ValidKeys(t *testing.T) {
	env := map[string]string{
		"DATABASE_URL": "postgres://localhost/db",
		"_PRIVATE":     "value",
		"key123":       "val",
	}
	result := Validate(env, DefaultValidateOptions())
	if result.HasErrors() {
		t.Fatalf("expected no errors, got: %s", result.Error())
	}
}

func TestValidate_InvalidKeyStartsWithDigit(t *testing.T) {
	env := map[string]string{
		"1BAD_KEY": "value",
	}
	result := Validate(env, DefaultValidateOptions())
	if !result.HasErrors() {
		t.Fatal("expected validation error for key starting with digit")
	}
	if result.Errors[0].Key != "1BAD_KEY" {
		t.Errorf("unexpected key in error: %s", result.Errors[0].Key)
	}
}

func TestValidate_InvalidKeyWithHyphen(t *testing.T) {
	env := map[string]string{
		"BAD-KEY": "value",
	}
	result := Validate(env, DefaultValidateOptions())
	if !result.HasErrors() {
		t.Fatal("expected validation error for key with hyphen")
	}
}

func TestValidate_RequireUppercase(t *testing.T) {
	env := map[string]string{
		"good_key": "value",
		"GOOD_KEY": "value",
	}
	opts := ValidateOptions{RequireUppercase: true}
	result := Validate(env, opts)
	if !result.HasErrors() {
		t.Fatal("expected error for lowercase key")
	}
	// Exactly one error (for "good_key")
	if len(result.Errors) != 1 {
		t.Errorf("expected 1 error, got %d", len(result.Errors))
	}
	if result.Errors[0].Key != "good_key" {
		t.Errorf("unexpected key: %s", result.Errors[0].Key)
	}
}

func TestValidate_ForbidEmpty(t *testing.T) {
	env := map[string]string{
		"PRESENT": "value",
		"EMPTY":   "",
	}
	opts := ValidateOptions{ForbidEmpty: true}
	result := Validate(env, opts)
	if !result.HasErrors() {
		t.Fatal("expected error for empty value")
	}
	if result.Errors[0].Key != "EMPTY" {
		t.Errorf("unexpected key: %s", result.Errors[0].Key)
	}
}

func TestValidate_EmptyMap(t *testing.T) {
	result := Validate(map[string]string{}, DefaultValidateOptions())
	if result.HasErrors() {
		t.Fatalf("expected no errors for empty map, got: %s", result.Error())
	}
}

func TestValidationResult_ErrorString(t *testing.T) {
	r := &ValidationResult{
		Errors: []ValidationError{
			{Key: "A", Message: "bad"},
			{Key: "B", Message: "also bad"},
		},
	}
	s := r.Error()
	if s == "" {
		t.Error("expected non-empty error string")
	}
}

func TestValidate_MultipleInvalidKeys(t *testing.T) {
	env := map[string]string{
		"1INVALID": "value",
		"BAD-KEY":  "value",
		"VALID_KEY": "value",
	}
	result := Validate(env, DefaultValidateOptions())
	if !result.HasErrors() {
		t.Fatal("expected validation errors for invalid keys")
	}
	if len(result.Errors) != 2 {
		t.Errorf("expected 2 errors, got %d", len(result.Errors))
	}
}
