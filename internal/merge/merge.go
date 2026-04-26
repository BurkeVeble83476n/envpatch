// Package merge provides functionality for merging .env files.
// It combines a base environment file with an overlay, allowing
// selective overrides while preserving keys not present in the overlay.
package merge

import (
	"fmt"
	"sort"
)

// Strategy defines how conflicts are resolved during a merge.
type Strategy int

const (
	// StrategyOverlay replaces base values with overlay values when keys conflict.
	StrategyOverlay Strategy = iota
	// StrategyKeepBase retains base values when keys conflict.
	StrategyKeepBase
)

// Result holds the outcome of a merge operation.
type Result struct {
	// Merged is the final key-value map after merging.
	Merged map[string]string
	// Added contains keys that exist only in the overlay.
	Added []string
	// Overridden contains keys that existed in base and were replaced by overlay.
	Overridden []string
}

// Merge combines base and overlay maps according to the given strategy.
// Keys present only in base are always preserved.
// Keys present only in overlay are always added.
// Conflicting keys are resolved by the strategy.
func Merge(base, overlay map[string]string, strategy Strategy) (*Result, error) {
	if base == nil {
		return nil, fmt.Errorf("merge: base map must not be nil")
	}
	if overlay == nil {
		return nil, fmt.Errorf("merge: overlay map must not be nil")
	}

	merged := make(map[string]string, len(base))
	for k, v := range base {
		merged[k] = v
	}

	var added, overridden []string

	for k, v := range overlay {
		if _, exists := merged[k]; exists {
			if strategy == StrategyOverlay {
				merged[k] = v
				overridden = append(overridden, k)
			}
		} else {
			merged[k] = v
			added = append(added, k)
		}
	}

	sort.Strings(added)
	sort.Strings(overridden)

	return &Result{
		Merged:     merged,
		Added:      added,
		Overridden: overridden,
	}, nil
}
