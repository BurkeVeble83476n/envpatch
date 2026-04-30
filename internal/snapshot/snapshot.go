// Package snapshot provides functionality to capture and compare
// .env file states over time, enabling drift detection between
// recorded snapshots and current environment configurations.
package snapshot

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// Snapshot represents a captured state of an .env file at a point in time.
type Snapshot struct {
	CapturedAt time.Time         `json:"captured_at"`
	Label      string            `json:"label"`
	Values     map[string]string `json:"values"`
}

// New creates a new Snapshot from the given env map with an optional label.
func New(env map[string]string, label string) (*Snapshot, error) {
	if env == nil {
		return nil, fmt.Errorf("snapshot: env map must not be nil")
	}
	copy := make(map[string]string, len(env))
	for k, v := range env {
		copy[k] = v
	}
	return &Snapshot{
		CapturedAt: time.Now().UTC(),
		Label:      label,
		Values:     copy,
	}, nil
}

// Save serialises the snapshot as JSON to the given writer.
func Save(s *Snapshot, w io.Writer) error {
	if s == nil {
		return fmt.Errorf("snapshot: cannot save nil snapshot")
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(s)
}

// Load deserialises a snapshot from the given reader.
func Load(r io.Reader) (*Snapshot, error) {
	if r == nil {
		return nil, fmt.Errorf("snapshot: reader must not be nil")
	}
	var s Snapshot
	if err := json.NewDecoder(r).Decode(&s); err != nil {
		return nil, fmt.Errorf("snapshot: decode failed: %w", err)
	}
	return &s, nil
}

// Drift describes a single key-level difference between a snapshot and a
// current env map.
type Drift struct {
	Key      string
	Kind     string // "added", "removed", "changed"
	OldValue string
	NewValue string
}

// Compare returns the list of drifts between the snapshot's recorded values
// and the provided current env map.
func Compare(s *Snapshot, current map[string]string) ([]Drift, error) {
	if s == nil {
		return nil, fmt.Errorf("snapshot: snapshot must not be nil")
	}
	if current == nil {
		return nil, fmt.Errorf("snapshot: current map must not be nil")
	}
	var drifts []Drift
	for k, oldVal := range s.Values {
		newVal, ok := current[k]
		if !ok {
			drifts = append(drifts, Drift{Key: k, Kind: "removed", OldValue: oldVal})
		} else if newVal != oldVal {
			drifts = append(drifts, Drift{Key: k, Kind: "changed", OldValue: oldVal, NewValue: newVal})
		}
	}
	for k, newVal := range current {
		if _, ok := s.Values[k]; !ok {
			drifts = append(drifts, Drift{Key: k, Kind: "added", NewValue: newVal})
		}
	}
	return drifts, nil
}
