package audit

import (
	"fmt"
	"io"
	"time"

	"github.com/user/envpatch/internal/diff"
)

// EntryKind describes the type of audit event.
type EntryKind string

const (
	KindMerge    EntryKind = "merge"
	KindPatch    EntryKind = "patch"
	KindValidate EntryKind = "validate"
	KindSnapshot EntryKind = "snapshot"
)

// Entry represents a single audit log record.
type Entry struct {
	Timestamp time.Time
	Kind      EntryKind
	Message   string
	Changes   []diff.Change
}

// Logger writes audit entries to an io.Writer.
type Logger struct {
	w io.Writer
}

// New creates a Logger that writes to w.
// Returns an error if w is nil.
func New(w io.Writer) (*Logger, error) {
	if w == nil {
		return nil, fmt.Errorf("audit: writer must not be nil")
	}
	return &Logger{w: w}, nil
}

// Record writes an audit entry to the underlying writer.
func (l *Logger) Record(kind EntryKind, message string, changes []diff.Change) error {
	entry := Entry{
		Timestamp: time.Now().UTC(),
		Kind:      kind,
		Message:   message,
		Changes:   changes,
	}
	_, err := fmt.Fprintf(l.w, "[%s] %s: %s (%d change(s))\n",
		entry.Timestamp.Format(time.RFC3339),
		entry.Kind,
		entry.Message,
		len(entry.Changes),
	)
	if err != nil {
		return fmt.Errorf("audit: failed to write entry: %w", err)
	}
	for _, c := range entry.Changes {
		_, err = fmt.Fprintf(l.w, "  key=%s op=%s\n", c.Key, c.Op)
		if err != nil {
			return fmt.Errorf("audit: failed to write change: %w", err)
		}
	}
	return nil
}
