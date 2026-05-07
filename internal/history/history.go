package history

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// Entry records a single pipeline execution event with metadata.
type Entry struct {
	Timestamp time.Time         `json:"timestamp"`
	Operation string            `json:"operation"`
	Changes   map[string]string `json:"changes,omitempty"`
	Message   string            `json:"message,omitempty"`
}

// Log is an ordered sequence of history entries.
type Log struct {
	Entries []Entry `json:"entries"`
}

// New returns an empty Log.
func New() *Log {
	return &Log{}
}

// Record appends a new entry to the log.
func (l *Log) Record(operation, message string, changes map[string]string) {
	l.Entries = append(l.Entries, Entry{
		Timestamp: time.Now().UTC(),
		Operation: operation,
		Message:   message,
		Changes:   changes,
	})
}

// Save serialises the log as JSON to w.
func (l *Log) Save(w io.Writer) error {
	if w == nil {
		return fmt.Errorf("history: writer must not be nil")
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(l)
}

// Load deserialises a Log from r.
func Load(r io.Reader) (*Log, error) {
	if r == nil {
		return nil, fmt.Errorf("history: reader must not be nil")
	}
	var l Log
	if err := json.NewDecoder(r).Decode(&l); err != nil {
		return nil, fmt.Errorf("history: decode: %w", err)
	}
	return &l, nil
}

// Len returns the number of entries in the log.
func (l *Log) Len() int { return len(l.Entries) }
