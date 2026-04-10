package envfile

import (
	"testing"
)

func TestDiff_DetectsAddedKey(t *testing.T) {
	existing := map[string]string{"FOO": "bar"}
	incoming := map[string]string{"FOO": "bar", "NEW_KEY": "value"}

	changes := Diff(existing, incoming)
	if !hasChangeOfType(changes, "NEW_KEY", ChangeAdded) {
		t.Error("expected NEW_KEY to be marked as added")
	}
}

func TestDiff_DetectsRemovedKey(t *testing.T) {
	existing := map[string]string{"FOO": "bar", "OLD_KEY": "gone"}
	incoming := map[string]string{"FOO": "bar"}

	changes := Diff(existing, incoming)
	if !hasChangeOfType(changes, "OLD_KEY", ChangeRemoved) {
		t.Error("expected OLD_KEY to be marked as removed")
	}
}

func TestDiff_DetectsUpdatedKey(t *testing.T) {
	existing := map[string]string{"FOO": "old"}
	incoming := map[string]string{"FOO": "new"}

	changes := Diff(existing, incoming)
	if !hasChangeOfType(changes, "FOO", ChangeUpdated) {
		t.Error("expected FOO to be marked as updated")
	}

	for _, c := range changes {
		if c.Key == "FOO" {
			if c.OldValue != "old" || c.NewValue != "new" {
				t.Errorf("unexpected old/new values: %q / %q", c.OldValue, c.NewValue)
			}
		}
	}
}

func TestDiff_DetectsUnchangedKey(t *testing.T) {
	existing := map[string]string{"FOO": "same"}
	incoming := map[string]string{"FOO": "same"}

	changes := Diff(existing, incoming)
	if !hasChangeOfType(changes, "FOO", ChangeUnchanged) {
		t.Error("expected FOO to be marked as unchanged")
	}
}

func TestHasChanges_ReturnsFalseWhenAllUnchanged(t *testing.T) {
	changes := []Change{
		{Key: "A", Type: ChangeUnchanged},
		{Key: "B", Type: ChangeUnchanged},
	}
	if HasChanges(changes) {
		t.Error("expected HasChanges to return false")
	}
}

func TestHasChanges_ReturnsTrueWhenChangePresent(t *testing.T) {
	changes := []Change{
		{Key: "A", Type: ChangeUnchanged},
		{Key: "B", Type: ChangeAdded},
	}
	if !HasChanges(changes) {
		t.Error("expected HasChanges to return true")
	}
}

// hasChangeOfType is a helper to find a change entry by key and type.
func hasChangeOfType(changes []Change, key string, ct ChangeType) bool {
	for _, c := range changes {
		if c.Key == key && c.Type == ct {
			return true
		}
	}
	return false
}
