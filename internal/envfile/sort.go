package envfile

import (
	"sort"
	"strings"
)

// SortOrder defines the ordering strategy for keys.
type SortOrder string

const (
	SortAsc  SortOrder = "asc"
	SortDesc SortOrder = "desc"
)

// SortOptions controls how Sort behaves.
type SortOptions struct {
	Order      SortOrder
	KeysOnly   []string // if set, only these keys are sorted; others appended after
	IgnoreCase bool
}

// DefaultSortOptions returns sensible defaults.
func DefaultSortOptions() SortOptions {
	return SortOptions{
		Order:      SortAsc,
		IgnoreCase: false,
	}
}

// Sort returns an ordered slice of keys from env according to opts.
// When KeysOnly is set, those keys are sorted first; remaining keys follow unsorted.
func Sort(env map[string]string, opts SortOptions) []string {
	var primary, rest []string

	if len(opts.KeysOnly) > 0 {
		included := toSet(opts.KeysOnly)
		for k := range env {
			if included[k] {
				primary = append(primary, k)
			} else {
				rest = append(rest, k)
			}
		}
	} else {
		for k := range env {
			primary = append(primary, k)
		}
	}

	cmpKey := func(s string) string {
		if opts.IgnoreCase {
			return strings.ToLower(s)
		}
		return s
	}

	sort.Slice(primary, func(i, j int) bool {
		a, b := cmpKey(primary[i]), cmpKey(primary[j])
		if opts.Order == SortDesc {
			return a > b
		}
		return a < b
	})

	sort.Strings(rest)
	return append(primary, rest...)
}
