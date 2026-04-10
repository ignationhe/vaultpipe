package envfile

// ChangeType represents the type of change for a key.
type ChangeType string

const (
	ChangeAdded   ChangeType = "added"
	ChangeRemoved ChangeType = "removed"
	ChangeUpdated ChangeType = "updated"
	ChangeUnchanged ChangeType = "unchanged"
)

// Change represents a single key-level diff entry.
type Change struct {
	Key      string
	OldValue string
	NewValue string
	Type     ChangeType
}

// Diff compares two env maps (existing vs incoming) and returns a slice of
// Change entries describing what would be added, removed, updated, or left
// unchanged.
func Diff(existing, incoming map[string]string) []Change {
	var changes []Change

	// Check for additions and updates.
	for key, newVal := range incoming {
		oldVal, exists := existing[key]
		if !exists {
			changes = append(changes, Change{
				Key:      key,
				NewValue: newVal,
				Type:     ChangeAdded,
			})
		} else if oldVal != newVal {
			changes = append(changes, Change{
				Key:      key,
				OldValue: oldVal,
				NewValue: newVal,
				Type:     ChangeUpdated,
			})
		} else {
			changes = append(changes, Change{
				Key:      key,
				OldValue: oldVal,
				NewValue: newVal,
				Type:     ChangeUnchanged,
			})
		}
	}

	// Check for removals.
	for key, oldVal := range existing {
		if _, exists := incoming[key]; !exists {
			changes = append(changes, Change{
				Key:      key,
				OldValue: oldVal,
				Type:     ChangeRemoved,
			})
		}
	}

	return changes
}

// HasChanges returns true if any of the provided changes are additions,
// updates, or removals.
func HasChanges(changes []Change) bool {
	for _, c := range changes {
		if c.Type != ChangeUnchanged {
			return true
		}
	}
	return false
}
