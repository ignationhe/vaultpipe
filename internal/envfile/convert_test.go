package envfile

import (
	"strings"
	"testing"
)

func TestConvert_DotenvFormat(t *testing.T) {
	env := map[string]string{"FOO": "bar", "BAZ": "qux"}
	opts := DefaultConvertOptions()
	out, err := Convert(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "FOO=bar") {
		t.Errorf("expected FOO=bar in output, got:\n%s", out)
	}
	if !strings.Contains(out, "BAZ=qux") {
		t.Errorf("expected BAZ=qux in output, got:\n%s", out)
	}
}

func TestConvert_YAMLFormat(t *testing.T) {
	env := map[string]string{"HOST": "localhost"}
	opts := ConvertOptions{Format: FormatYAML, Sort: true}
	out, err := Convert(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "HOST: localhost") {
		t.Errorf("expected YAML line, got:\n%s", out)
	}
}

func TestConvert_TOMLFormat(t *testing.T) {
	env := map[string]string{"PORT": "8080"}
	opts := ConvertOptions{Format: FormatTOML, Sort: true}
	out, err := Convert(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, `PORT = "8080"`) {
		t.Errorf("expected TOML line, got:\n%s", out)
	}
}

func TestConvert_UnknownFormatReturnsError(t *testing.T) {
	env := map[string]string{"K": "v"}
	_, err := Convert(env, ConvertOptions{Format: "xml"})
	if err == nil {
		t.Fatal("expected error for unknown format")
	}
}

func TestConvert_PrefixApplied(t *testing.T) {
	env := map[string]string{"NAME": "alice"}
	opts := ConvertOptions{Format: FormatDotenv, Sort: true, Prefix: "APP_"}
	out, err := Convert(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "APP_NAME=alice") {
		t.Errorf("expected prefixed key, got:\n%s", out)
	}
}

func TestConvert_QuotesValuesWithSpaces(t *testing.T) {
	env := map[string]string{"MSG": "hello world"}
	opts := DefaultConvertOptions()
	out, err := Convert(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, `MSG="hello world"`) {
		t.Errorf("expected quoted value, got:\n%s", out)
	}
}

func TestConvert_SortedOutput(t *testing.T) {
	env := map[string]string{"ZEBRA": "1", "ALPHA": "2", "MIDDLE": "3"}
	opts := ConvertOptions{Format: FormatDotenv, Sort: true}
	out, err := Convert(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	alphaIdx := strings.Index(out, "ALPHA")
	middleIdx := strings.Index(out, "MIDDLE")
	zebraIdx := strings.Index(out, "ZEBRA")
	if !(alphaIdx < middleIdx && middleIdx < zebraIdx) {
		t.Errorf("expected sorted output, got:\n%s", out)
	}
}
