package envfile

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExport_DotenvFormat(t *testing.T) {
	secrets := map[string]string{
		"FOO": "bar",
		"BAZ": "qux",
	}
	var buf bytes.Buffer
	opts := DefaultExportOptions()
	if err := Export(&buf, secrets, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "BAZ=qux") {
		t.Errorf("expected BAZ=qux in output, got:\n%s", out)
	}
	if !strings.Contains(out, "FOO=bar") {
		t.Errorf("expected FOO=bar in output, got:\n%s", out)
	}
}

func TestExport_ExportFormat(t *testing.T) {
	secrets := map[string]string{"MY_VAR": "hello world"}
	var buf bytes.Buffer
	opts := ExportOptions{Format: FormatExport, Sorted: true}
	if err := Export(&buf, secrets, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.HasPrefix(out, "export MY_VAR=") {
		t.Errorf("expected export prefix, got: %s", out)
	}
	if !strings.Contains(out, `"hello world"`) {
		t.Errorf("expected quoted value, got: %s", out)
	}
}

func TestExport_JSONFormat(t *testing.T) {
	secrets := map[string]string{"KEY": "value"}
	var buf bytes.Buffer
	opts := ExportOptions{Format: FormatJSON, Sorted: true}
	if err := Export(&buf, secrets, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, `"KEY"`) || !strings.Contains(out, `"value"`) {
		t.Errorf("expected JSON output with KEY/value, got:\n%s", out)
	}
	if !strings.HasPrefix(out, "{") {
		t.Errorf("expected JSON object, got: %s", out)
	}
}

func TestExport_SortedOutput(t *testing.T) {
	secrets := map[string]string{"ZZZ": "last", "AAA": "first", "MMM": "mid"}
	var buf bytes.Buffer
	opts := ExportOptions{Format: FormatDotenv, Sorted: true}
	Export(&buf, secrets, opts)
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if !strings.HasPrefix(lines[0], "AAA") {
		t.Errorf("expected AAA first, got: %s", lines[0])
	}
	if !strings.HasPrefix(lines[2], "ZZZ") {
		t.Errorf("expected ZZZ last, got: %s", lines[2])
	}
}

func TestExportToFile_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.env")
	secrets := map[string]string{"TOKEN": "abc123"}
	if err := ExportToFile(path, secrets, DefaultExportOptions()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("could not read file: %v", err)
	}
	if !strings.Contains(string(data), "TOKEN=abc123") {
		t.Errorf("expected TOKEN=abc123 in file, got: %s", string(data))
	}
}
