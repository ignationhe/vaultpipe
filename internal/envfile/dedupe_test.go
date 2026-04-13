package envfile

import (
	"testing"
)

func pairs(kv ...string) [][2]string {
	var out [][2]string
	for i := 0; i+1 < len(kv); i += 2 {
		out = append(out, [2]string{kv[i], kv[i+1]})
	}
	return out
}

func TestDedupe_KeepLastByDefault(t *testing.T) {
	res, err := Dedupe(pairs("KEY", "first", "KEY", "second"), DefaultDedupeOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Env["KEY"] != "second" {
		t.Errorf("expected 'second', got %q", res.Env["KEY"])
	}
	if len(res.Duplicates) != 1 || res.Duplicates[0] != "KEY" {
		t.Errorf("expected duplicates=[KEY], got %v", res.Duplicates)
	}
}

func TestDedupe_KeepFirst(t *testing.T) {
	opts := DedupeOptions{Strategy: DedupeKeepFirst}
	res, err := Dedupe(pairs("KEY", "first", "KEY", "second"), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Env["KEY"] != "first" {
		t.Errorf("expected 'first', got %q", res.Env["KEY"])
	}
}

func TestDedupe_ErrorOnDuplicate(t *testing.T) {
	opts := DedupeOptions{Strategy: DedupeError}
	_, err := Dedupe(pairs("KEY", "a", "KEY", "b"), opts)
	if err == nil {
		t.Fatal("expected error for duplicate key, got nil")
	}
}

func TestDedupe_NoDuplicates(t *testing.T) {
	res, err := Dedupe(pairs("A", "1", "B", "2", "C", "3"), DefaultDedupeOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Duplicates) != 0 {
		t.Errorf("expected no duplicates, got %v", res.Duplicates)
	}
	if len(res.Env) != 3 {
		t.Errorf("expected 3 keys, got %d", len(res.Env))
	}
}

func TestDedupe_EmptyInput(t *testing.T) {
	res, err := Dedupe(nil, DefaultDedupeOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Env) != 0 {
		t.Errorf("expected empty map, got %v", res.Env)
	}
}

func TestDedupe_MultipleDuplicateKeys(t *testing.T) {
	res, err := Dedupe(pairs("X", "1", "Y", "a", "X", "2", "Y", "b", "X", "3"), DefaultDedupeOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Env["X"] != "3" {
		t.Errorf("expected X=3, got %q", res.Env["X"])
	}
	if res.Env["Y"] != "b" {
		t.Errorf("expected Y=b, got %q", res.Env["Y"])
	}
	if len(res.Duplicates) != 2 {
		t.Errorf("expected 2 duplicate keys, got %v", res.Duplicates)
	}
}
