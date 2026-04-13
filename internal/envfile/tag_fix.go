package envfile

import "sort"

// This file exists solely to satisfy the sort import used in tag.go
// without modifying that file. The sortedKeys helper is defined in tag.go
// and uses sort.Strings, so we ensure the import is resolved here.
// (In a real module the helper would live in a shared util file.)

var _ = sort.Strings // ensure import is used
