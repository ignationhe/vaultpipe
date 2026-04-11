package envfile

import (
	"strings"
	"testing"
)

func TestLint_WarnOnLowercaseKey(t *testing.T) {
	env := map[string]string{"myKey": "value"}
	opts := DefaultLintOptions()
	issues := Lint(env, opts)
	if !hasIssueFor(issues, "myKey", LintWarning, "lowercase") {
		t.Errorf("expected lowercase warning for key 'myKey', got %v", issues)
	}
}

func TestLint_NoWarnUppercaseKey(t *testing.T) {
	env := map[string]string{"MY_KEY": "value"}
	opts := DefaultLintOptions()
	issues := Lint(env, opts)
	if hasIssueFor(issues, "MY_KEY", LintWarning, "lowercase") {
		t.Errorf("unexpected lowercase warning for key 'MY_KEY'")
	}
}

func TestLint_WarnOnEmptyValue(t *testing.T) {
	env := map[string]string{"EMPTY_KEY": ""}
	opts := DefaultLintOptions()
	issues := Lint(env, opts)
	if !hasIssueFor(issues, "EMPTY_KEY", LintWarning, "empty") {
		t.Errorf("expected empty value warning, got %v", issues)
	}
}

func TestLint_ErrorOnLeadingUnderscore(t *testing.T) {
	env := map[string]string{"_INTERNAL": "secret"}
	opts := DefaultLintOptions()
	opts.ErrorOnLeadingUnderscore = true
	issues := Lint(env, opts)
	if !hasIssueFor(issues, "_INTERNAL", LintError, "underscore") {
		t.Errorf("expected error for leading underscore, got %v", issues)
	}
}

func TestLint_NoErrorLeadingUnderscoreWhenDisabled(t *testing.T) {
	env := map[string]string{"_INTERNAL": "secret"}
	opts := DefaultLintOptions()
	opts.ErrorOnLeadingUnderscore = false
	issues := Lint(env, opts)
	if hasIssueFor(issues, "_INTERNAL", LintError, "underscore") {
		t.Errorf("unexpected underscore error when rule is disabled")
	}
}

func TestLint_WarnOnLongValue(t *testing.T) {
	env := map[string]string{"BIG_SECRET": strings.Repeat("x", 300)}
	opts := DefaultLintOptions()
	issues := Lint(env, opts)
	if !hasIssueFor(issues, "BIG_SECRET", LintWarning, "exceeds") {
		t.Errorf("expected long value warning, got %v", issues)
	}
}

func TestLint_NoIssuesForCleanEnv(t *testing.T) {
	env := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
	}
	opts := DefaultLintOptions()
	issues := Lint(env, opts)
	if len(issues) != 0 {
		t.Errorf("expected no issues for clean env, got %v", issues)
	}
}

func TestLint_IssueStringFormat(t *testing.T) {
	issue := LintIssue{Key: "bad_key", Message: "some problem", Severity: LintWarning}
	s := issue.String()
	if !strings.Contains(s, "warning") || !strings.Contains(s, "bad_key") {
		t.Errorf("unexpected issue string format: %s", s)
	}
}

// hasIssueFor checks whether issues contains an entry matching key, severity,
// and a message that contains the given substring.
func hasIssueFor(issues []LintIssue, key string, sev LintSeverity, msgSubstr string) bool {
	for _, i := range issues {
		if i.Key == key && i.Severity == sev && strings.Contains(i.Message, msgSubstr) {
			return true
		}
	}
	return false
}
