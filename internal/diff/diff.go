// Package diff provides functionality to compare two parsed .env files
// and report added, removed, and changed keys.
package diff

import "sort"

// ChangeType represents the kind of change for a key.
type ChangeType string

const (
	Added   ChangeType = "added"
	Removed ChangeType = "removed"
	Changed ChangeType = "changed"
)

// Change describes a single difference between two env files.
type Change struct {
	Key      string
	Type     ChangeType
	OldValue string
	NewValue string
}

// Result holds all differences between a base and target env map.
type Result struct {
	Changes []Change
}

// HasChanges returns true if there are any differences.
func (r *Result) HasChanges() bool {
	return len(r.Changes) > 0
}

// Diff compares base and target env maps and returns a Result.
// Keys present only in target are Added; only in base are Removed;
// present in both with different values are Changed.
func Diff(base, target map[string]string) *Result {
	result := &Result{}

	for key, targetVal := range target {
		baseVal, exists := base[key]
		if !exists {
			result.Changes = append(result.Changes, Change{
				Key:      key,
				Type:     Added,
				NewValue: targetVal,
			})
		} else if baseVal != targetVal {
			result.Changes = append(result.Changes, Change{
				Key:      key,
				Type:     Changed,
				OldValue: baseVal,
				NewValue: targetVal,
			})
		}
	}

	for key, baseVal := range base {
		if _, exists := target[key]; !exists {
			result.Changes = append(result.Changes, Change{
				Key:      key,
				Type:     Removed,
				OldValue: baseVal,
			})
		}
	}

	sort.Slice(result.Changes, func(i, j int) bool {
		return result.Changes[i].Key < result.Changes[j].Key
	})

	return result
}
