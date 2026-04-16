package envfile

import (
	"testing"
)

func baseEnvForSort() map[string]string {
	return map[string]string{
		"ZEBRA": "1",
		"APPLE": "2",
		"MANGO": "3",
		"banana": "4",
	}
}

func TestSort_AscendingOrder(t *testing.T) {
	env := baseEnvForSort()
	keys := Sort(env, DefaultSortOptions())
	if keys[0] != "APPLE" || keys[1] != "MANGO" || keys[2] != "ZEBRA" {
		t.Errorf("unexpected order: %v", keys)
	}
}

func TestSort_Descend) {
	env := map[string]string{"A": "1", "B": "2", "C": "3"}
	keys := Sort(env, SortOptions{Order: SortDesc})
	if keys[0] != "C" || keys[1] != "B" || keys[2] != "A" {
		t.Errorf("unexpected desc order: %v", keys)
	}
}

func TestSort_IgnoreCase(t *testing.T) {
	env := map[string]string{"zebra": "1", "Apple": "2", "mango": "3"}
	keys := Sort(env, SortOptions{Order: SortAsc, IgnoreCase: true})
	if keys[0] != "Apple" || keys[1] != "mango" || keys[2] != "zebra" {
		t.Errorf("unexpected case-insensitive order: %v", keys)
	}
}

func TestSort_KeysOnly_SortedFirst(t *testing.T) {
	env := map[string]string{"Z": "1", "A": "2", "M": "3", "EXTRA": "4"}
	opts := SortOptions{Order: SortAsc, KeysOnly: []string{"Z", "A", "M"}}
	keys := Sort(env, opts)
	if keys[0] != "A" || keys[1] != "M" || keys[2] != "Z" {
		t.Errorf("primary keys not sorted: %v", keys)
	}
	if keys[3] != "EXTRA" {
		t.Errorf("expected EXTRA at end, got %v", keys[3])
	}
}

func TestSort_EmptyEnv(t *testing.T) {
	keys := Sort(map[string]string{}, DefaultSortOptions())
	if len(keys) != 0 {
		t.Errorf("expected empty, got %v", keys)
	}
}
