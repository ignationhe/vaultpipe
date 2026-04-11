package envfile

import "testing"

func TestCompare_SameKeys(t *testing.T) {
	a := map[string]string{"FOO": "bar", "BAZ": "qux"}
	b := map[string]string{"FOO": "bar", "BAZ": "qux"}
	r := Compare(a, b)
	if r.HasDifferences() {
		t.Error("expected no differences")
	}
	if len(r.Same) != 2 {
		t.Errorf("expected 2 same keys, got %d", len(r.Same))
	}
}

func TestCompare_OnlyInA(t *testing.T) {
	a := map[string]string{"FOO": "bar", "ONLY_A": "val"}
	b := map[string]string{"FOO": "bar"}
	r := Compare(a, b)
	if !r.HasDifferences() {
		t.Error("expected differences")
	}
	if _, ok := r.OnlyInA["ONLY_A"]; !ok {
		t.Error("expected ONLY_A in OnlyInA")
	}
}

func TestCompare_OnlyInB(t *testing.T) {
	a := map[string]string{"FOO": "bar"}
	b := map[string]string{"FOO": "bar", "ONLY_B": "val"}
	r := Compare(a, b)
	if !r.HasDifferences() {
		t.Error("expected differences")
	}
	if _, ok := r.OnlyInB["ONLY_B"]; !ok {
		t.Error("expected ONLY_B in OnlyInB")
	}
}

func TestCompare_DifferentValues(t *testing.T) {
	a := map[string]string{"FOO": "old"}
	b := map[string]string{"FOO": "new"}
	r := Compare(a, b)
	if !r.HasDifferences() {
		t.Error("expected differences")
	}
	pair, ok := r.Different["FOO"]
	if !ok {
		t.Fatal("expected FOO in Different")
	}
	if pair[0] != "old" || pair[1] != "new" {
		t.Errorf("unexpected pair: %v", pair)
	}
}

func TestCompare_SortedKeys(t *testing.T) {
	a := map[string]string{"Z": "1", "A": "2"}
	b := map[string]string{"M": "3", "A": "2"}
	r := Compare(a, b)
	keys := r.SortedKeys()
	if len(keys) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(keys))
	}
	if keys[0] != "A" || keys[1] != "M" || keys[2] != "Z" {
		t.Errorf("unexpected order: %v", keys)
	}
}

func TestCompare_EmptyMaps(t *testing.T) {
	r := Compare(map[string]string{}, map[string]string{})
	if r.HasDifferences() {
		t.Error("expected no differences for empty maps")
	}
	if len(r.SortedKeys()) != 0 {
		t.Error("expected no keys")
	}
}
