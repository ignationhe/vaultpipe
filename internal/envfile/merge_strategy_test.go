package envfile

import (
	"testing"
)

func TestMergeWith_NoConflict(t *testing.T) {
	base := map[string]string{"A": "1"}
	incoming := map[string]string{"B": "2"}

	result := MergeWith(base, incoming, DefaultMergeOptions())

	if result["A"] != "1" {
		t.Errorf("expected A=1, got %q", result["A"])
	}
	if result["B"] != "2" {
		t.Errorf("expected B=2, got %q", result["B"])
	}
}

func TestMergeWith_KeepExistingOnConflict(t *testing.T) {
	base := map[string]string{"KEY": "original"}
	incoming := map[string]string{"KEY": "new"}

	opts := MergeOptions{Strategy: StrategyKeepExisting}
	result := MergeWith(base, incoming, opts)

	if result["KEY"] != "original" {
		t.Errorf("expected original, got %q", result["KEY"])
	}
}

func TestMergeWith_OverwriteOnConflict(t *testing.T) {
	base := map[string]string{"KEY": "original"}
	incoming := map[string]string{"KEY": "new"}

	opts := MergeOptions{Strategy: StrategyOverwrite}
	result := MergeWith(base, incoming, opts)

	if result["KEY"] != "new" {
		t.Errorf("expected new, got %q", result["KEY"])
	}
}

func TestMergeWith_KeepBothOnConflict(t *testing.T) {
	base := map[string]string{"KEY": "original"}
	incoming := map[string]string{"KEY": "new"}

	opts := MergeOptions{Strategy: StrategyKeepBoth, ConflictSuffix: "_INCOMING"}
	result := MergeWith(base, incoming, opts)

	if result["KEY"] != "original" {
		t.Errorf("expected original for KEY, got %q", result["KEY"])
	}
	if result["KEY_INCOMING"] != "new" {
		t.Errorf("expected new for KEY_INCOMING, got %q", result["KEY_INCOMING"])
	}
}

func TestMergeWith_DefaultSuffixWhenEmpty(t *testing.T) {
	base := map[string]string{"X": "a"}
	incoming := map[string]string{"X": "b"}

	opts := MergeOptions{Strategy: StrategyKeepBoth, ConflictSuffix: ""}
	result := MergeWith(base, incoming, opts)

	if _, ok := result["X_NEW"]; !ok {
		t.Error("expected X_NEW key to exist with default suffix")
	}
}

func TestMergeWith_DoesNotMutateBase(t *testing.T) {
	base := map[string]string{"A": "1"}
	incoming := map[string]string{"A": "2"}

	opts := MergeOptions{Strategy: StrategyOverwrite}
	MergeWith(base, incoming, opts)

	if base["A"] != "1" {
		t.Error("base map was mutated")
	}
}
