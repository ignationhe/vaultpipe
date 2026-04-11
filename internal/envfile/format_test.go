package envfile

import (
	"strings"
	"testing"
)

func TestFormat_SortsKeysAlphabetically(t *testing.T) {
	secrets := map[string]string{
		"ZEBRA": "last",
		"ALPHA": "first",
		"MANGO": "middle",
	}
	opts := DefaultFormatOptions()
	out := Format(secrets, opts)
	lines := nonEmptyLines(out)
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "ALPHA=") {
		t.Errorf("expected first line to be ALPHA, got %s", lines[0])
	}
	if !strings.HasPrefix(lines[2], "ZEBRA=") {
		t.Errorf("expected last line to be ZEBRA, got %s", lines[2])
	}
}

func TestFormat_SectionComment(t *testing.T) {
	secrets := map[string]string{"FOO": "bar"}
	opts := DefaultFormatOptions()
	opts.SectionComment = "managed by vaultpipe"
	out := Format(secrets, opts)
	if !strings.HasPrefix(out, "# managed by vaultpipe\n") {
		t.Errorf("expected section comment header, got: %s", out)
	}
}

func TestFormat_QuotesValuesWithSpaces(t *testing.T) {
	secrets := map[string]string{"MSG": "hello world"}
	out := Format(secrets, DefaultFormatOptions())
	if !strings.Contains(out, `MSG="hello world"`) {
		t.Errorf("expected quoted value, got: %s", out)
	}
}

func TestFormat_BlankLineBetweenKeys(t *testing.T) {
	secrets := map[string]string{"A": "1", "B": "2"}
	opts := DefaultFormatOptions()
	opts.BlankLineBetweenKeys = true
	out := Format(secrets, opts)
	// Expect a blank line between A and B
	if !strings.Contains(out, "\n\n") {
		t.Errorf("expected blank line between keys, got: %q", out)
	}
}

func TestFormat_EmptyMap(t *testing.T) {
	out := Format(map[string]string{}, DefaultFormatOptions())
	if out != "" {
		t.Errorf("expected empty string for empty map, got: %q", out)
	}
}

func TestFormat_NoSortPreservesInsertionOrder(t *testing.T) {
	// With SortKeys=false we just verify output contains all keys.
	secrets := map[string]string{"Z": "26", "A": "1"}
	opts := DefaultFormatOptions()
	opts.SortKeys = false
	out := Format(secrets, opts)
	if !strings.Contains(out, "Z=26") || !strings.Contains(out, "A=1") {
		t.Errorf("expected both keys in output, got: %s", out)
	}
}

func nonEmptyLines(s string) []string {
	var result []string
	for _, l := range strings.Split(s, "\n") {
		if strings.TrimSpace(l) != "" {
			result = append(result, l)
		}
	}
	return result
}
