package envfile

import "fmt"

// ChainStep represents a named transformation step applied to an env map.
type ChainStep struct {
	Name string
	Fn   func(map[string]string) (map[string]string, error)
}

// ChainResult holds the output of each step in the chain.
type ChainResult struct {
	Step   string
	Output map[string]string
	Err    error
}

// ChainOptions configures pipeline execution behaviour.
type ChainOptions struct {
	// StopOnError halts the chain on the first step that returns an error.
	StopOnError bool
	// Trace captures intermediate results for every step when true.
	Trace bool
}

// DefaultChainOptions returns sensible defaults.
func DefaultChainOptions() ChainOptions {
	return ChainOptions{
		StopOnError: true,
		Trace:       false,
	}
}

// Chain runs a sequence of transformation steps against an initial env map.
// It returns the final map, a slice of ChainResult (populated when Trace is
// enabled or a step errors), and the first error encountered (if any).
func Chain(
	initial map[string]string,
	steps []ChainStep,
	opts ChainOptions,
) (map[string]string, []ChainResult, error) {
	current := copyMap(initial)
	var results []ChainResult

	for _, step := range steps {
		if step.Fn == nil {
			err := fmt.Errorf("chain: step %q has nil function", step.Name)
			results = append(results, ChainResult{Step: step.Name, Err: err})
			if opts.StopOnError {
				return current, results, err
			}
			continue
		}

		out, err := step.Fn(current)
		if err != nil {
			r := ChainResult{Step: step.Name, Err: err}
			if opts.Trace {
				r.Output = copyMap(current)
			}
			results = append(results, r)
			if opts.StopOnError {
				return current, results, fmt.Errorf("chain: step %q failed: %w", step.Name, err)
			}
			continue
		}

		current = out
		if opts.Trace {
			results = append(results, ChainResult{Step: step.Name, Output: copyMap(current)})
		}
	}

	return current, results, nil
}

func copyMap(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}
