package envfile

import "fmt"

// PluckOptions controls how keys are extracted.
type PluckOptions struct {
	Keys        []string
	ErrorOnMiss bool
}

// DefaultPluckOptions returns sensible defaults.
func DefaultPluckOptions() PluckOptions {
	return PluckOptions{
		ErrorOnMiss: false,
	}
}

// Pluck returns a new map containing only the specified keys.
// If ErrorOnMiss is true, any key not found in src returns an error.
func Pluck(src map[string]string, opts PluckOptions) (map[string]string, error) {
	if len(opts.Keys) == 0 {
		return copyMapPluck(src), nil
	}

	out := make(map[string]string, len(opts.Keys))
	for _, k := range opts.Keys {
		v, ok := src[k]
		if !ok {
			if opts.ErrorOnMiss {
				return nil, fmt.Errorf("pluck: key %q not found", k)
			}
			continue
		}
		out[k] = v
	}
	return out, nil
}

func copyMapPluck(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}
