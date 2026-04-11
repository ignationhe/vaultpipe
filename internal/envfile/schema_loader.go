package envfile

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
)

// schemaFieldJSON is the JSON representation of a SchemaField.
type schemaFieldJSON struct {
	Key      string `json:"key"`
	Required bool   `json:"required"`
	Pattern  string `json:"pattern,omitempty"`
	Desc     string `json:"desc,omitempty"`
}

// LoadSchema reads a JSON schema file and returns a Schema.
// The JSON file should be an array of field objects.
func LoadSchema(path string) (Schema, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Schema{}, fmt.Errorf("schema: read file: %w", err)
	}

	var raw []schemaFieldJSON
	if err := json.Unmarshal(data, &raw); err != nil {
		return Schema{}, fmt.Errorf("schema: parse JSON: %w", err)
	}

	var fields []SchemaField
	for _, r := range raw {
		field := SchemaField{
			Key:      r.Key,
			Required: r.Required,
			Desc:     r.Desc,
		}
		if r.Pattern != "" {
			pat, err := regexp.Compile(r.Pattern)
			if err != nil {
				return Schema{}, fmt.Errorf("schema: invalid pattern for key %q: %w", r.Key, err)
			}
			field.Pattern = pat
		}
		fields = append(fields, field)
	}

	return Schema{Fields: fields}, nil
}
