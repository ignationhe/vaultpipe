package envfile

import (
	"fmt"
	"regexp"
)

// SchemaField describes a single expected key in an env file.
type SchemaField struct {
	Key      string
	Required bool
	Pattern  *regexp.Regexp // optional value pattern
	Desc     string
}

// Schema is a collection of expected fields.
type Schema struct {
	Fields []SchemaField
}

// SchemaViolation represents a single schema mismatch.
type SchemaViolation struct {
	Key     string
	Message string
}

func (v SchemaViolation) Error() string {
	return fmt.Sprintf("schema violation for %q: %s", v.Key, v.Message)
}

// ValidateSchema checks the provided env map against the schema.
// It returns a slice of violations (never nil, may be empty).
func ValidateSchema(env map[string]string, schema Schema) []SchemaViolation {
	var violations []SchemaViolation

	for _, field := range schema.Fields {
		val, exists := env[field.Key]

		if !exists || val == "" {
			if field.Required {
				violations = append(violations, SchemaViolation{
					Key:     field.Key,
					Message: "required key is missing or empty",
				})
			}
			continue
		}

		if field.Pattern != nil && !field.Pattern.MatchString(val) {
			violations = append(violations, SchemaViolation{
				Key:     field.Key,
				Message: fmt.Sprintf("value %q does not match pattern %s", val, field.Pattern.String()),
			})
		}
	}

	return violations
}
