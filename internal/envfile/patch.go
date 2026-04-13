package envfile

import "fmt"

// PatchOp represents the type of patch operation.
type PatchOp string

const (
	PatchSet    PatchOp = "set"
	PatchDelete PatchOp = "delete"
	PatchRename PatchOp = "rename"
)

// PatchRule defines a single patch instruction.
type PatchRule struct {
	Op      PatchOp
	Key     string
	Value   string // used by set
	NewKey  string // used by rename
}

// DefaultPatchOptions returns a PatchOptions with safe defaults.
func DefaultPatchOptions() PatchOptions {
	return PatchOptions{
		IgnoreMissing: true,
	}
}

// PatchOptions controls Patch behaviour.
type PatchOptions struct {
	// IgnoreMissing skips delete/rename rules for keys not present in the map.
	IgnoreMissing bool
}

// Patch applies a sequence of PatchRules to env, returning a new map.
func Patch(env map[string]string, rules []PatchRule, opts PatchOptions) (map[string]string, error) {
	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = v
	}

	for _, r := range rules {
		switch r.Op {
		case PatchSet:
			out[r.Key] = r.Value

		case PatchDelete:
			if _, ok := out[r.Key]; !ok {
				if !opts.IgnoreMissing {
					return nil, fmt.Errorf("patch delete: key %q not found", r.Key)
				}
				continue
			}
			delete(out, r.Key)

		case PatchRename:
			val, ok := out[r.Key]
			if !ok {
				if !opts.IgnoreMissing {
					return nil, fmt.Errorf("patch rename: key %q not found", r.Key)
				}
				continue
			}
			if r.NewKey == "" {
				return nil, fmt.Errorf("patch rename: new_key must not be empty for key %q", r.Key)
			}
			delete(out, r.Key)
			out[r.NewKey] = val

		default:
			return nil, fmt.Errorf("patch: unknown op %q", r.Op)
		}
	}
	return out, nil
}
