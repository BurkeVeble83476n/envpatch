// Package patch applies a set of changes (additions, updates, removals) to an
// existing env map, producing a new map that reflects the desired state.
package patch

import (
	"errors"
	"fmt"
)

// Op describes the kind of change a PatchEntry represents.
type Op string

const (
	// OpSet adds or updates a key with the given value.
	OpSet Op = "set"
	// OpRemove deletes a key from the target map.
	OpRemove Op = "remove"
)

// Entry is a single instruction in a patch.
type Entry struct {
	Op    Op
	Key   string
	Value string // only meaningful for OpSet
}

// Patch holds an ordered list of entries to apply.
type Patch struct {
	Entries []Entry
}

// Apply applies p to a copy of base and returns the resulting map.
// base must not be nil. An empty Patch is valid and returns a copy of base.
func Apply(base map[string]string, p Patch) (map[string]string, error) {
	if base == nil {
		return nil, errors.New("patch: base map must not be nil")
	}

	out := make(map[string]string, len(base))
	for k, v := range base {
		out[k] = v
	}

	for _, e := range p.Entries {
		if e.Key == "" {
			return nil, fmt.Errorf("patch: entry has empty key (op=%s)", e.Op)
		}
		switch e.Op {
		case OpSet:
			out[e.Key] = e.Value
		case OpRemove:
			delete(out, e.Key)
		default:
			return nil, fmt.Errorf("patch: unknown op %q for key %q", e.Op, e.Key)
		}
	}

	return out, nil
}
