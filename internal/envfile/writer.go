package envfile

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

// Write writes a map of key-value pairs to the given file path in .env format.
// If the file already exists, it will be overwritten.
func Write(path string, secrets map[string]string) error {
	var sb strings.Builder

	keys := make([]string, 0, len(secrets))
	for k := range secrets {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := secrets[k]
		if needsQuoting(v) {
			v = fmt.Sprintf("%q", v)
		}
		sb.WriteString(fmt.Sprintf("%s=%s\n", k, v))
	}

	return os.WriteFile(path, []byte(sb.String()), 0600)
}

// Merge applies the provided secrets on top of the existing .env file at path.
// Keys present in secrets will be added or updated; keys not in secrets are preserved.
func Merge(path string, secrets map[string]string) error {
	existing, err := Parse(path)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("reading existing env file: %w", err)
	}
	if existing == nil {
		existing = make(map[string]string)
	}

	for k, v := range secrets {
		existing[k] = v
	}

	return Write(path, existing)
}

// needsQuoting returns true if the value contains characters that require quoting.
func needsQuoting(v string) bool {
	for _, c := range v {
		if c == ' ' || c == '\t' || c == '#' || c == '"' || c == '\'' || c == '\n' {
			return true
		}
	}
	return false
}
