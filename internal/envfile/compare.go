package envfile

import "sort"

// CompareResult holds the result of comparing two env maps.
type CompareResult struct {
	OnlyInA   map[string]string // keys present in A but not B
	OnlyInB   map[string]string // keys present in B but not A
	Different map[string][2]string // keys in both but with different values [a, b]
	Same      map[string]string // keys with identical values in both
}

// SortedKeys returns a sorted slice of all keys across all categories.
func (r *CompareResult) SortedKeys() []string {
	seen := map[string]struct{}{}
	for k := range r.OnlyInA {
		seen[k] = struct{}{}
	}
	for k := range r.OnlyInB {
		seen[k] = struct{}{}
	}
	for k := range r.Different {
		seen[k] = struct{}{}
	}
	for k := range r.Same {
		seen[k] = struct{}{}
	}
	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// HasDifferences returns true if there are any keys that differ between A and B.
func (r *CompareResult) HasDifferences() bool {
	return len(r.OnlyInA) > 0 || len(r.OnlyInB) > 0 || len(r.Different) > 0
}

// Compare compares two env maps and returns a CompareResult.
func Compare(a, b map[string]string) *CompareResult {
	result := &CompareResult{
		OnlyInA:   make(map[string]string),
		OnlyInB:   make(map[string]string),
		Different: make(map[string][2]string),
		Same:      make(map[string]string),
	}
	for k, va := range a {
		if vb, ok := b[k]; ok {
			if va == vb {
				result.Same[k] = va
			} else {
				result.Different[k] = [2]string{va, vb}
			}
		} else {
			result.OnlyInA[k] = va
		}
	}
	for k, vb := range b {
		if _, ok := a[k]; !ok {
			result.OnlyInB[k] = vb
		}
	}
	return result
}
